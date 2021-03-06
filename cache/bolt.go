package cache

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/crdt"
	"github.com/johnny-morrice/godless/log"
	"github.com/johnny-morrice/godless/proto"
	"github.com/pkg/errors"

	pb "github.com/gogo/protobuf/proto"
)

type BoltOptions struct {
	DBOptions    *bolt.Options
	FilePath     string
	Mode         os.FileMode
	Db           *bolt.DB
	MaxCacheSize int
}

type BoltFactory struct {
	BoltOptions
}

func (factory BoltFactory) MakeCache() (api.Cache, error) {
	const failMsg = "BoltFactory.MakeCache failed"

	cache := boltCache{
		db:      factory.Db,
		maxSize: factory.MaxCacheSize,
	}

	err := cache.initBuckets()

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

	return cache, nil
}

func (factory BoltFactory) MakeMemoryImage() (api.MemoryImage, error) {
	const failMsg = "BoltFactory.MakeMemoryImage failed"

	memImg := boltMemoryImage{db: factory.Db}

	err := memImg.initBuckets()

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

	return memImg, nil
}

func MakeBoltFactory(options BoltOptions) (BoltFactory, error) {
	const failMsg = "MakeBoltCacheFactory"

	if options.MaxCacheSize <= 0 {
		options.MaxCacheSize = __DEFAULT_BUFFER_SIZE
	}

	db, err := connectBolt(options)

	if err != nil {
		return BoltFactory{}, errors.Wrap(err, failMsg)
	}

	factory := BoltFactory{
		BoltOptions: options,
	}
	factory.Db = db

	return factory, nil
}

type boltCache struct {
	db      *bolt.DB
	maxSize int
}

func (cache boltCache) initBuckets() error {
	return createAllBucketsIfNotExists(cache.db, BOLT_NAMESPACE_CACHE_BUCKET, BOLT_INDEX_CACHE_BUCKET, BOLT_HEAD_CACHE_BUCKET)
}

func (cache boltCache) GetHead() (crdt.IPFSPath, error) {
	const failMsg = "boltCache.GetHead failed"

	var head crdt.IPFSPath
	err := cache.viewHead(func(bucket *bolt.Bucket) error {
		value := bucket.Get(BOLT_HEAD_CACHE_KEY)
		head = crdt.IPFSPath(value)
		return nil
	})

	if err != nil {
		return crdt.NIL_PATH, errors.Wrap(err, failMsg)
	}

	if !crdt.IsNilPath(head) {
		log.Info("Found HEAD in Bolt: %s", head)
	}

	return head, nil
}

func (cache boltCache) SetHead(head crdt.IPFSPath) error {
	const failMsg = "boltCache.SetHead failed"

	value := []byte(head)

	err := cache.updateHead(func(bucket *bolt.Bucket) error {
		return bucket.Put(BOLT_HEAD_CACHE_KEY, value)
	})

	if err != nil {
		return errors.Wrap(err, failMsg)
	}

	log.Info("Wrote HEAD to Bolt: %s", head)

	return nil
}

func (cache boltCache) GetIndex(indexAddr crdt.IPFSPath) (crdt.Index, error) {
	const failMsg = "boltCache.GetIndex failed"

	indexMessage := &proto.IndexMessage{}
	key := []byte(indexAddr)
	err := cache.viewIndex(func(bucket *bolt.Bucket) error {
		return getCacheItemValue(bucket, key, indexMessage)
	})

	if err != nil {
		return crdt.EmptyIndex(), errors.Wrap(err, failMsg)
	}

	// TODO handle the invalid entries.
	index, _ := crdt.ReadIndexMessage(indexMessage)

	log.Info("Found index in Bolt: %s", indexAddr)

	return index, nil
}

func (cache boltCache) SetIndex(indexAddr crdt.IPFSPath, index crdt.Index) error {
	const failMsg = "boltCache.SetIndex failed"

	// TODO handle the invalid entries.
	indexMessage, _ := crdt.MakeIndexMessage(index)
	key := []byte(indexAddr)

	err := cache.updateIndex(func(bucket *bolt.Bucket) error {
		return putCacheItem(bucket, key, indexMessage)
	})

	if err != nil {
		return errors.Wrap(err, failMsg)
	}

	log.Info("Wrote index to Bolt: %s", indexAddr)

	return nil
}

func (cache boltCache) GetNamespace(namespaceAddr crdt.IPFSPath) (crdt.Namespace, error) {
	const failMsg = "boltCache.GetNamespace"

	namespaceMessage := &proto.NamespaceMessage{}
	key := []byte(namespaceAddr)
	err := cache.viewNamespace(func(bucket *bolt.Bucket) error {
		return getCacheItemValue(bucket, key, namespaceMessage)
	})

	if err != nil {
		return crdt.EmptyNamespace(), errors.Wrap(err, failMsg)
	}

	namespace, invalid := crdt.ReadNamespaceMessage(namespaceMessage)

	if len(invalid) > 0 {
		log.Error("Bolt ignoring %d invalid entries", len(invalid))
	}

	log.Info("Found Namespace in Bolt: %s", namespaceAddr)

	return namespace, nil
}

func (cache boltCache) SetNamespace(namespaceAddr crdt.IPFSPath, namespace crdt.Namespace) error {
	const failMsg = "boltCache.SetNamespace"

	// TODO handle invalid entries
	namespaceMessage, _ := crdt.MakeNamespaceMessage(namespace)
	key := []byte(namespaceAddr)

	err := cache.updateNamespace(func(bucket *bolt.Bucket) error {
		return putCacheItem(bucket, key, namespaceMessage)
	})

	if err != nil {
		return errors.Wrap(err, failMsg)
	}

	log.Info("Wrote Namespace to Bolt: %s", namespaceAddr)

	return nil
}

func (cache boltCache) viewNamespace(viewer func(bucket *bolt.Bucket) error) error {
	return cache.db.View(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_NAMESPACE_CACHE_BUCKET)

		if err != nil {
			return err
		}

		return viewer(bucket)
	})
}

func (cache boltCache) updateNamespace(updater func(bucket *bolt.Bucket) error) error {
	return cache.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_NAMESPACE_CACHE_BUCKET)

		if err != nil {
			return err
		}

		err = cache.evictOldest(bucket)

		if err != nil {
			return err
		}

		return updater(bucket)
	})
}

func (cache boltCache) viewIndex(viewer func(bucket *bolt.Bucket) error) error {
	return cache.db.View(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_INDEX_CACHE_BUCKET)

		if err != nil {
			return err
		}

		return viewer(bucket)
	})
}

func (cache boltCache) updateIndex(updater func(bucket *bolt.Bucket) error) error {
	return cache.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_INDEX_CACHE_BUCKET)

		if err != nil {
			return err
		}

		err = cache.evictOldest(bucket)

		if err != nil {
			return err
		}

		return updater(bucket)
	})
}

func (cache boltCache) viewHead(viewer func(bucket *bolt.Bucket) error) error {
	return cache.db.View(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_HEAD_CACHE_BUCKET)

		if err != nil {
			return err
		}

		return viewer(bucket)
	})
}

func (cache boltCache) updateHead(updater func(bucket *bolt.Bucket) error) error {
	return cache.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_HEAD_CACHE_BUCKET)

		if err != nil {
			return err
		}

		return updater(bucket)
	})
}

func (cache boltCache) CloseCache() error {
	err := cache.db.Close()
	log.Info("Closed boltCache")
	return err
}

func (cache boltCache) evictOldest(bucket *bolt.Bucket) error {
	const failMsg = "evictOldest failed"

	stats := bucket.Stats()

	if stats.KeyN <= cache.maxSize {
		return nil
	}

	var oldest boltCacheItem

	err := bucket.ForEach(func(k []byte, v []byte) error {
		inner := bucket.Bucket(k)

		if inner == nil {
			return nil
		}

		item := boltCacheItem{key: k}
		err := item.getTimestamp(inner)

		if err != nil {
			return err
		}

		if oldest.key == nil {
			oldest = item
			return nil
		}

		if item.older(oldest.timestamp) {
			oldest = item
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, failMsg)
	}

	if oldest.key == nil {
		return errors.New("Corrupt cache")
	}

	err = bucket.DeleteBucket(oldest.key)

	return errors.Wrap(err, failMsg)
}

type boltMemoryImage struct {
	db *bolt.DB
}

func (memimg boltMemoryImage) initBuckets() error {
	return createAllBucketsIfNotExists(memimg.db, BOLT_MEMORY_IMAGE_BUCKET)
}

func (memimg boltMemoryImage) GetIndex() (crdt.Index, error) {
	const failMsg = "boltMemoryIndex.GetIndex failed"
	indexMessage := &proto.IndexMessage{}

	err := memimg.view(func(bucket *bolt.Bucket) error {
		indexBytes := bucket.Get(BOLT_MEMORY_IMAGE_INDEX_KEY)

		if indexBytes == nil {
			return nil
		}

		return pb.Unmarshal(indexBytes, indexMessage)
	})

	if err != nil {
		return crdt.EmptyIndex(), errors.Wrap(err, failMsg)
	}

	// TODO handle the invalid entries.
	index, _ := crdt.ReadIndexMessage(indexMessage)

	log.Info("Read Bolt MemoryImage")

	return index, nil
}

func (memimg boltMemoryImage) JoinIndex(index crdt.Index) error {
	const failMsg = "boltMemoryIndex.JoinIndex failed"

	err := memimg.update(func(bucket *bolt.Bucket) error {
		currentMessage := &proto.IndexMessage{}
		currentIndexBytes := bucket.Get(BOLT_MEMORY_IMAGE_INDEX_KEY)

		currentIndex := crdt.EmptyIndex()
		if currentIndexBytes != nil {
			// TODO handle the invalid entries.
			pbErr := pb.Unmarshal(currentIndexBytes, currentMessage)

			if pbErr != nil {
				return pbErr
			}

			currentIndex, _ = crdt.ReadIndexMessage(currentMessage)
		}

		joinedIndex := currentIndex.JoinIndex(index)

		// TODO handle the invalid entries.
		joinedMessage, _ := crdt.MakeIndexMessage(joinedIndex)
		return putMessage(bucket, BOLT_MEMORY_IMAGE_INDEX_KEY, joinedMessage)
	})

	if err != nil {
		return errors.Wrap(err, failMsg)
	}

	log.Info("Updated Bolt MemoryImage")

	return nil
}

func (memimg boltMemoryImage) view(viewer func(bucket *bolt.Bucket) error) error {
	return memimg.db.View(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_MEMORY_IMAGE_BUCKET)

		if err != nil {
			return err
		}

		return viewer(bucket)
	})
}

func (memimg boltMemoryImage) update(updater func(bucket *bolt.Bucket) error) error {
	return memimg.db.Update(func(transaction *bolt.Tx) error {
		bucket, err := getBucket(transaction, BOLT_MEMORY_IMAGE_BUCKET)

		if err != nil {
			return err
		}

		return updater(bucket)
	})
}

func (memimg boltMemoryImage) CloseMemoryImage() error {
	err := memimg.db.Close()
	log.Info("Closed boltMemoryImage")
	return err
}

func connectBolt(options BoltOptions) (*bolt.DB, error) {
	return bolt.Open(options.FilePath, options.Mode, options.DBOptions)
}

func getBucket(transaction *bolt.Tx, bucketName []byte) (*bolt.Bucket, error) {
	bucket := transaction.Bucket(bucketName)

	if bucket == nil {
		bucketNameText := string(bucketName)
		return nil, fmt.Errorf("No bucket for: %s", bucketNameText)
	}

	return bucket, nil
}

func createAllBucketsIfNotExists(db *bolt.DB, bucketName ...[]byte) error {
	const failMsg = "createAllBucketsIfNotExists"
	return db.Update(func(transaction *bolt.Tx) error {
		for _, name := range bucketName {
			_, err := transaction.CreateBucketIfNotExists(name)

			if err != nil {
				return errors.Wrap(err, failMsg)
			}
		}

		return nil
	})

	return nil
}

type boltCacheItem struct {
	timestamp
	key []byte
}

func (item *boltCacheItem) putTimestamp(bucket *bolt.Bucket) error {
	err := bucket.Put(TIMESTAMP_KEY, slice64(item.timestamp.seconds))

	if err != nil {
		msg := fmt.Sprintf("Failed to add timestamp at Bolt key: %s", string(item.key))
		return errors.Wrap(err, msg)
	}

	err = bucket.Put(NANO_TIMESTAMP_KEY, slice64(item.timestamp.nanoseconds))

	if err != nil {
		msg := fmt.Sprintf("Failed to add nanosecond timestamp at Bolt key: %s", string(item.key))
		return errors.Wrap(err, msg)
	}

	return nil
}

func (item *boltCacheItem) getTimestamp(bucket *bolt.Bucket) error {
	timestamp := bucket.Get(TIMESTAMP_KEY)

	if timestamp == nil {
		return fmt.Errorf("Failed to Get timestamp at Bolt key: %s", string(item.key))
	}

	nano := bucket.Get(NANO_TIMESTAMP_KEY)

	if nano == nil {
		return fmt.Errorf("Failed to Get nano timestamp at Bolt key: %s", string(item.key))
	}

	item.timestamp.seconds = deslice64(timestamp)
	item.timestamp.nanoseconds = deslice64(nano)

	return nil
}

func putMessage(bucket *bolt.Bucket, key []byte, value pb.Message) error {
	keyText := string(key)
	valueBytes, err := pb.Marshal(value)

	if err != nil {
		msg := fmt.Sprintf("Failed to Marshal protobuf message for Bolt key: %s", keyText)
		return errors.Wrap(err, msg)
	}

	err = bucket.Put(key, valueBytes)

	if err != nil {
		msg := fmt.Sprintf("Failed to Put value at Bolt key: %s", keyText)
		return errors.Wrap(err, msg)
	}

	return nil
}

func putCacheItem(bucket *bolt.Bucket, key []byte, value pb.Message) error {
	keyText := string(key)

	if isKeyPresent(bucket, key) {
		return nil
	}

	valueBytes, err := pb.Marshal(value)

	if err != nil {
		msg := fmt.Sprintf("Failed to Marshal protobuf message for Bolt key: %s", keyText)
		return errors.Wrap(err, msg)
	}

	innerBucket, err := bucket.CreateBucket(key)

	if err != nil {
		msg := fmt.Sprintf("Failed to create inner bucket for Bolt key: %s", keyText)
		return errors.Wrap(err, msg)
	}

	item := boltCacheItem{
		key:       key,
		timestamp: makeTimestamp(),
	}

	err = item.putTimestamp(innerBucket)

	if err != nil {
		return err
	}

	err = innerBucket.Put(DATA_KEY, valueBytes)

	if err != nil {
		msg := fmt.Sprintf("Failed to Put value at Bolt key: %s", keyText)
		return errors.Wrap(err, msg)
	}

	return nil
}

func getCacheItemValue(bucket *bolt.Bucket, key []byte, value pb.Message) error {
	keyText := string(key)
	innerBucket := bucket.Bucket(key)

	if innerBucket == nil {
		return fmt.Errorf("Failed to get inner Bucket at Bolt key: %s", keyText)
	}

	valueBytes := innerBucket.Get(DATA_KEY)

	if valueBytes == nil {
		return fmt.Errorf("Failed to Get value at Bolt key: %s", keyText)
	}

	err := pb.Unmarshal(valueBytes, value)

	if err != nil {
		msg := fmt.Sprintf("Failed to Unmarshal protobuf message for Bolt key: %s", keyText)
		return errors.Wrap(err, msg)
	}

	return nil
}

func isKeyPresent(bucket *bolt.Bucket, key []byte) bool {
	bkt := bucket.Bucket(key)
	value := bucket.Get(key)
	return bkt != nil || value != nil
}

func deslice64(bs []byte) int64 {
	return int64(__BYTE_ORDER.Uint64(bs))
}

func slice64(n int64) []byte {
	bs := make([]byte, 8)
	__BYTE_ORDER.PutUint64(bs, uint64(n))
	return bs
}

var __BYTE_ORDER = binary.LittleEndian

var TIMESTAMP_KEY = []byte("timestamp")
var NANO_TIMESTAMP_KEY = []byte("nano_timestamp")
var DATA_KEY = []byte("data")
var BOLT_HEAD_CACHE_KEY = []byte("head")
var BOLT_HEAD_CACHE_BUCKET = []byte("head_cache")
var BOLT_NAMESPACE_CACHE_BUCKET = []byte("namespace_cache")
var BOLT_INDEX_CACHE_BUCKET = []byte("index_cache")
var BOLT_MEMORY_IMAGE_INDEX_KEY = []byte("current_index")
var BOLT_MEMORY_IMAGE_BUCKET = []byte("memory_image")
