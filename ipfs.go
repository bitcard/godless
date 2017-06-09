package godless

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/pkg/errors"
)

type IPFSPath string

func (path IPFSPath) Path() string {
	return string(path)
}

func castIPFSPath(addr RemoteStoreAddress) IPFSPath {
	path, ok := addr.(IPFSPath)

	if !ok {
		panic("addr was not IPFSPath")
	}

	return path
}

type IPFSRecord struct {
	Namespace Namespace
}

func makeIpfsRecord(namespace Namespace) *IPFSRecord {
	return &IPFSRecord{
		Namespace: namespace,
	}
}

func (record *IPFSRecord) encode(w io.Writer) error {
	return EncodeNamespace(record.Namespace, w)
}

func (record *IPFSRecord) decode(r io.Reader) error {
	ns, err := DecodeNamespace(r)

	if err != nil {
		return err
	}

	record.Namespace = ns
	return nil
}

type encoder interface {
	encode(io.Writer) error
}

type decoder interface {
	decode(io.Reader) error
}

type IPFSIndex struct {
	Index RemoteNamespaceIndex
}

func makeIpfsIndex(index RemoteNamespaceIndex) *IPFSIndex {
	return &IPFSIndex{
		Index: index,
	}
}

func (index *IPFSIndex) encode(w io.Writer) error {
	return EncodeIndex(index.Index, w)
}

func (index *IPFSIndex) decode(r io.Reader) error {
	dx, err := DecodeIndex(r)

	if err != nil {
		return err
	}

	index.Index = dx
	return nil
}

// TODO Don't use Shell directly - invent an interface.  This would enable mocking.
type IPFSPeer struct {
	Offline bool
	Url     string
	Client  *http.Client
	Shell   *ipfs.Shell
}

func MakeIPFSPeer(url string, offline bool) RemoteStore {
	peer := &IPFSPeer{
		Url:     url,
		Client:  defaultHttpClient(),
		Offline: offline,
	}

	return peer
}

func (peer *IPFSPeer) Connect() error {
	peer.Shell = ipfs.NewShellWithClient(peer.Url, peer.Client)

	if !peer.Shell.IsUp() {
		return fmt.Errorf("IPFSPeer is not up at '%v'", peer.Url)
	}

	return nil
}

func (peer *IPFSPeer) Disconnect() error {
	// Nothing to do.
	return nil
}

func (peer *IPFSPeer) SubscribeAddrStream(topic RemoteStoreAddress) (<-chan RemoteStoreAddress, <-chan error) {
	stream := make(chan RemoteStoreAddress)
	errch := make(chan error)

	go func() {
		defer close(stream)
		defer close(errch)

		topicText := topic.Path()
		subscription, launchErr := peer.Shell.PubSubSubscribe(topicText)

		if launchErr != nil {
			errch <- launchErr
			return
		}

		for {
			record, recordErr := subscription.Next()

			if recordErr != nil {
				errch <- recordErr
				return
			}

			pubsubPeer := record.From()
			bs := record.Data()
			addr := IPFSPath(string(bs))

			stream <- addr
			loginfo("Subscription update: '%v' from '%v'", addr, pubsubPeer)
		}

	}()

	return stream, errch
}

func (peer *IPFSPeer) AddIndex(index RemoteNamespaceIndex) (RemoteStoreAddress, error) {
	const failMsg = "IPFSPeer.AddIndex failed"

	chunk := makeIpfsIndex(index)

	path, addErr := peer.add(chunk)

	if addErr != nil {
		return nil, errors.Wrap(addErr, failMsg)
	}

	return path, nil
}

func (peer *IPFSPeer) CatIndex(addr RemoteStoreAddress) (RemoteNamespaceIndex, error) {
	path := castIPFSPath(addr)

	chunk := &IPFSIndex{}
	caterr := peer.cat(path, chunk)

	if caterr != nil {
		return EMPTY_INDEX, errors.Wrap(caterr, "IPFSPeer.CatNamespace failed")
	}

	return chunk.Index, nil
}

func (peer *IPFSPeer) AddNamespace(record RemoteNamespaceRecord) (RemoteStoreAddress, error) {
	chunk := makeIpfsRecord(record.Namespace)

	path, err := peer.add(chunk)

	if err != nil {
		return nil, errors.Wrap(err, "IPFSPeer.AddNamespace failed")
	}

	return path, nil
}

func (peer *IPFSPeer) CatNamespace(addr RemoteStoreAddress) (RemoteNamespaceRecord, error) {
	path := castIPFSPath(addr)

	chunk := &IPFSRecord{}
	caterr := peer.cat(path, chunk)

	if caterr != nil {
		return EMPTY_RECORD, errors.Wrap(caterr, "IPFSPeer.CatNamespace failed")
	}

	record := RemoteNamespaceRecord{Namespace: chunk.Namespace}
	return record, nil
}

func (peer *IPFSPeer) add(chunk encoder) (IPFSPath, error) {
	const failMsg = "IPFSPeer.add failed"
	buff := &bytes.Buffer{}
	err := chunk.encode(buff)

	if err != nil {
		return "", errors.Wrap(err, failMsg)
	}

	path, sherr := peer.Shell.Add(buff)

	if sherr != nil {
		return "", errors.Wrap(err, failMsg)
	}

	return IPFSPath(path), nil
}

func (peer *IPFSPeer) cat(path IPFSPath, out decoder) error {
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
		logwarn("error draining reader: %v", drainerr)
	}

	if len(remainder) != 0 {
		logwarn("remaining bits after gob: %v", remainder)
	}

	return nil
}
