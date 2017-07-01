package ipfs

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/crdt"

	"github.com/johnny-morrice/godless/log"
	"github.com/pkg/errors"
)

type IpfsRecord struct {
	Namespace crdt.Namespace
}

func makeIpfsRecord(namespace crdt.Namespace) *IpfsRecord {
	return &IpfsRecord{
		Namespace: namespace,
	}
}

func (record *IpfsRecord) encode(w io.Writer) error {
	invalid, err := crdt.EncodeNamespace(record.Namespace, w)

	record.logInvalid(invalid)

	return err
}

func (record *IpfsRecord) decode(r io.Reader) error {
	ns, invalid, err := crdt.DecodeNamespace(r)

	record.logInvalid(invalid)

	if err != nil {
		return err
	}

	record.Namespace = ns
	return nil
}

func (record *IpfsRecord) logInvalid(invalid []crdt.InvalidNamespaceEntry) {
	invalidCount := len(invalid)

	if invalidCount > 0 {
		log.Error("IPFSRecord: %d invalid entries", invalidCount)
	}
}

type encoder interface {
	encode(io.Writer) error
}

type decoder interface {
	decode(io.Reader) error
}

type IPFSIndex struct {
	Index crdt.Index
}

func makeIpfsIndex(index crdt.Index) *IPFSIndex {
	return &IPFSIndex{
		Index: index,
	}
}

func (index *IPFSIndex) encode(w io.Writer) error {
	invalid, err := crdt.EncodeIndex(index.Index, w)

	index.logInvalid(invalid)

	return err
}

func (index *IPFSIndex) decode(r io.Reader) error {
	dx, invalid, err := crdt.DecodeIndex(r)

	index.logInvalid(invalid)

	if err != nil {
		return err
	}

	// TODO should cache the invalid details.
	if len(invalid) > 0 {
		log.Warn("IPFSIndex Decoded invalid index entries")
	}

	index.Index = dx
	return nil
}

func (index *IPFSIndex) logInvalid(invalid []crdt.InvalidIndexEntry) {
	invalidCount := len(invalid)

	if invalidCount > 0 {
		log.Error("IPFSRecord: %d invalid entries", invalidCount)
	}
}

type ipfsCloser struct {
	sync.Mutex
	closed bool
}

func (closer *ipfsCloser) isClosed() bool {
	closer.Lock()
	defer closer.Unlock()

	return closer.closed
}

func (closer *ipfsCloser) closeIpfs() {
	closer.Lock()
	defer closer.Unlock()

	closer.closed = true
}

type IpfsRemoteStore struct {
	Shell  api.DataPeer
	closer ipfsCloser
}

func MakeIpfsRemoteStore(peer api.DataPeer) api.RemoteStore {
	return &IpfsRemoteStore{Shell: peer}
}

func (peer *IpfsRemoteStore) Connect() error {
	return peer.Shell.Connect()
}

func (peer *IpfsRemoteStore) Disconnect() error {
	peer.closer.closeIpfs()
	return peer.Shell.Disconnect()
}

func (peer *IpfsRemoteStore) validateShell() error {
	if peer.Shell == nil {
		return peer.Connect()
	}

	return peer.validateConnection()
}

func (peer *IpfsRemoteStore) validateConnection() error {
	if !peer.Shell.IsUp() {
		return errors.New("IPFSPeer is not up")
	}

	return nil
}

func (peer *IpfsRemoteStore) PublishAddr(addr crdt.Link, topics []api.PubSubTopic) error {
	const failMsg = "IPFSPeer.PublishAddr failed"

	if verr := peer.validateShell(); verr != nil {
		return verr
	}

	publishValue, printErr := crdt.PrintLink(addr)

	if printErr != nil {
		return errors.Wrap(printErr, failMsg)
	}

	for _, t := range topics {
		topicText := string(t)
		log.Info("Publishing to topic: %s", t)
		pubsubErr := peer.Shell.PubSubPublish(topicText, string(publishValue))

		if pubsubErr != nil {
			log.Warn("Pubsub failed (topic %s): %s", t, pubsubErr.Error())
			continue
		}

		log.Info("Published to topic: %s", t)
	}

	return nil
}

func (peer *IpfsRemoteStore) SubscribeAddrStream(topic api.PubSubTopic) (<-chan crdt.Link, <-chan error) {
	stream := make(chan crdt.Link)
	errch := make(chan error)

	tidy := func() {
		close(stream)
		close(errch)
	}

	if verr := peer.validateShell(); verr != nil {
		go func() {
			errch <- verr
			defer tidy()
		}()

		return stream, errch
	}

	go func() {
		defer tidy()

		topicText := string(topic)

		var subscription api.PubSubSubscription

	RESTART:
		for {
			if peer.closer.isClosed() {
				return
			}

			var launchErr error
			log.Info("(Re)starting subscription on %s", topic)
			subscription, launchErr = peer.Shell.PubSubSubscribe(topicText)

			if launchErr != nil {
				log.Error("Subcription launch failed, retrying: %s", launchErr.Error())
				continue
			}

			for {
				if peer.closer.isClosed() {
					return
				}

				log.Info("Fetching next subscription message on %s...", topic)
				record, recordErr := subscription.Next()

				if recordErr != nil {
					log.Error("Subscription read failed (topic %s), continuing: %s", topic, recordErr.Error())
					continue RESTART
				}

				pubsubPeer := record.From()
				bs := record.Data()
				addr, err := crdt.ParseLink(crdt.LinkText(bs))

				if err != nil {
					log.Warn("Bad link from peer (topic %s): %v", topic, pubsubPeer)
					continue
				}

				stream <- addr
				log.Info("Subscription update: '%s'", addr.Path())
			}
		}

	}()

	return stream, errch
}

func (peer *IpfsRemoteStore) AddIndex(index crdt.Index) (crdt.IPFSPath, error) {
	const failMsg = "IPFSPeer.AddIndex failed"

	log.Info("Adding index to IPFS...")

	if verr := peer.validateShell(); verr != nil {
		return crdt.NIL_PATH, verr
	}

	chunk := makeIpfsIndex(index)

	path, addErr := peer.add(chunk)

	if addErr != nil {
		return crdt.NIL_PATH, errors.Wrap(addErr, failMsg)
	}

	log.Info("Added index")

	return path, nil
}

func (peer *IpfsRemoteStore) CatIndex(addr crdt.IPFSPath) (crdt.Index, error) {
	log.Info("Catting index from IPFS at: %s ...", addr)

	if verr := peer.validateShell(); verr != nil {
		return crdt.EmptyIndex(), verr
	}

	chunk := &IPFSIndex{}
	caterr := peer.cat(addr, chunk)

	if caterr != nil {
		return crdt.EmptyIndex(), errors.Wrap(caterr, "IPFSPeer.CatNamespace failed")
	}

	log.Info("Catted index")

	return chunk.Index, nil
}

func (peer *IpfsRemoteStore) AddNamespace(namespace crdt.Namespace) (crdt.IPFSPath, error) {
	log.Info("Adding Namespace to IPFS...")

	if verr := peer.validateShell(); verr != nil {
		return crdt.NIL_PATH, verr
	}

	chunk := makeIpfsRecord(namespace)

	path, err := peer.add(chunk)

	if err != nil {
		return crdt.NIL_PATH, errors.Wrap(err, "IPFSPeer.AddNamespace failed")
	}

	log.Info("Added namespace")

	return path, nil
}

func (peer *IpfsRemoteStore) CatNamespace(addr crdt.IPFSPath) (crdt.Namespace, error) {
	log.Info("Catting namespace from IPFS at: %s ...", addr)

	if verr := peer.validateShell(); verr != nil {
		return crdt.EmptyNamespace(), verr
	}

	chunk := &IpfsRecord{}
	caterr := peer.cat(addr, chunk)

	if caterr != nil {
		return crdt.EmptyNamespace(), errors.Wrap(caterr, "IPFSPeer.CatNamespace failed")
	}

	log.Info("Catted namespace")

	return chunk.Namespace, nil
}

func (peer *IpfsRemoteStore) add(chunk encoder) (crdt.IPFSPath, error) {
	const failMsg = "IPFSPeer.add failed"
	buff := &bytes.Buffer{}
	err := chunk.encode(buff)

	if err != nil {
		return crdt.NIL_PATH, errors.Wrap(err, failMsg)
	}

	path, err := peer.Shell.Add(buff)

	if err != nil {
		return crdt.NIL_PATH, errors.Wrap(err, failMsg)
	}

	return crdt.IPFSPath(path), nil
}

func (peer *IpfsRemoteStore) cat(path crdt.IPFSPath, out decoder) error {
	const failMsg = "IPFSPeer.cat failed"
	reader, err := peer.Shell.Cat(string(path))

	if err != nil {
		return errors.Wrap(err, failMsg)
	}

	defer reader.Close()

	err = out.decode(reader)

	if err != nil {
		return errors.Wrap(err, failMsg)
	}

	// According to IPFS binding docs we must drain the reader.
	remainder, drainerr := ioutil.ReadAll(reader)

	if drainerr != nil {
		log.Warn("error draining reader: %s", drainerr.Error())
	}

	if len(remainder) != 0 {
		log.Warn("remaining bits after gob: %d", remainder)
	}

	return nil
}
