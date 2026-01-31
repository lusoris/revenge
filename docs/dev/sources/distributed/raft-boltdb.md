# hashicorp/raft-boltdb

> Source: https://pkg.go.dev/github.com/hashicorp/raft-boltdb/v2
> Fetched: 2026-01-31T16:03:13.940804+00:00
> Content-Hash: b2144ee4a823d9d1
> Type: html

---

### Index ¶

  * Variables
  * type BoltStore
  *     * func MigrateToV2(source, destination string) (*BoltStore, error)
    * func New(options Options) (*BoltStore, error)
    * func NewBoltStore(path string) (*BoltStore, error)
  *     * func (b *BoltStore) Close() error
    * func (b *BoltStore) DeleteRange(min, max uint64) error
    * func (b *BoltStore) FirstIndex() (uint64, error)
    * func (b *BoltStore) Get(k []byte) ([]byte, error)
    * func (b *BoltStore) GetLog(idx uint64, log *raft.Log) error
    * func (b *BoltStore) GetUint64(key []byte) (uint64, error)
    * func (b *BoltStore) LastIndex() (uint64, error)
    * func (b *BoltStore) RunMetrics(ctx context.Context, interval time.Duration)
    * func (b *BoltStore) Set(k, v []byte) error
    * func (b *BoltStore) SetUint64(key []byte, val uint64) error
    * func (b *BoltStore) Stats() bbolt.Stats
    * func (b *BoltStore) StoreLog(log *raft.Log) error
    * func (b *BoltStore) StoreLogs(logs []*raft.Log) error
    * func (b *BoltStore) Sync() error
  * type Options



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L24)
    
    
    var (
    
    	// An error indicating a given key does not exist
    	ErrKeyNotFound = [errors](/errors).[New](/errors#New)("not found")
    )

### Functions ¶

This section is empty.

### Types ¶

####  type [BoltStore](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L36) ¶
    
    
    type BoltStore struct {
    	// contains filtered or unexported fields
    }

BoltStore provides access to Bbolt for Raft to store and retrieve log entries. It also provides key/value storage, and can be used as a LogStore and StableStore. 

####  func [MigrateToV2](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L315) ¶
    
    
    func MigrateToV2(source, destination [string](/builtin#string)) (*BoltStore, [error](/builtin#error))

MigrateToV2 reads in the source file path of a BoltDB file and outputs all the data migrated to a Bbolt destination file 

####  func [New](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L80) ¶
    
    
    func New(options Options) (*BoltStore, [error](/builtin#error))

New uses the supplied options to open the Bbolt and prepare it for use as a raft backend. 

####  func [NewBoltStore](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L75) ¶
    
    
    func NewBoltStore(path [string](/builtin#string)) (*BoltStore, [error](/builtin#error))

NewBoltStore takes a file path and returns a connected Raft backend. 

####  func (*BoltStore) [Close](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L130) ¶
    
    
    func (b *BoltStore) Close() [error](/builtin#error)

Close is used to gracefully close the DB connection. 

####  func (*BoltStore) [DeleteRange](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L234) ¶
    
    
    func (b *BoltStore) DeleteRange(min, max [uint64](/builtin#uint64)) [error](/builtin#error)

DeleteRange is used to delete logs within a given range inclusively. 

####  func (*BoltStore) [FirstIndex](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L135) ¶
    
    
    func (b *BoltStore) FirstIndex() ([uint64](/builtin#uint64), [error](/builtin#error))

FirstIndex returns the first known index from the Raft log. 

####  func (*BoltStore) [Get](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L276) ¶
    
    
    func (b *BoltStore) Get(k [][byte](/builtin#byte)) ([][byte](/builtin#byte), [error](/builtin#error))

Get is used to retrieve a value from the k/v store by key 

####  func (*BoltStore) [GetLog](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L167) ¶
    
    
    func (b *BoltStore) GetLog(idx [uint64](/builtin#uint64), log *[raft](/github.com/hashicorp/raft).[Log](/github.com/hashicorp/raft#Log)) [error](/builtin#error)

GetLog is used to retrieve a log from Bbolt at a given index. 

####  func (*BoltStore) [GetUint64](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L298) ¶
    
    
    func (b *BoltStore) GetUint64(key [][byte](/builtin#byte)) ([uint64](/builtin#uint64), [error](/builtin#error))

GetUint64 is like Get, but handles uint64 values 

####  func (*BoltStore) [LastIndex](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L151) ¶
    
    
    func (b *BoltStore) LastIndex() ([uint64](/builtin#uint64), [error](/builtin#error))

LastIndex returns the last known index from the Raft log. 

####  func (*BoltStore) [RunMetrics](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/metrics.go#L20) ¶
    
    
    func (b *BoltStore) RunMetrics(ctx [context](/context).[Context](/context#Context), interval [time](/time).[Duration](/time#Duration))

RunMetrics should be executed in a go routine and will periodically emit metrics on the given interval until the context has been cancelled. 

####  func (*BoltStore) [Set](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L260) ¶
    
    
    func (b *BoltStore) Set(k, v [][byte](/builtin#byte)) [error](/builtin#error)

Set is used to set a key/value set outside of the raft log 

####  func (*BoltStore) [SetUint64](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L293) ¶
    
    
    func (b *BoltStore) SetUint64(key [][byte](/builtin#byte), val [uint64](/builtin#uint64)) [error](/builtin#error)

SetUint64 is like Set, but handles uint64 values 

####  func (*BoltStore) [Stats](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L125) ¶
    
    
    func (b *BoltStore) Stats() [bbolt](/go.etcd.io/bbolt).[Stats](/go.etcd.io/bbolt#Stats)

####  func (*BoltStore) [StoreLog](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L186) ¶
    
    
    func (b *BoltStore) StoreLog(log *[raft](/github.com/hashicorp/raft).[Log](/github.com/hashicorp/raft#Log)) [error](/builtin#error)

StoreLog is used to store a single raft log 

####  func (*BoltStore) [StoreLogs](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L191) ¶
    
    
    func (b *BoltStore) StoreLogs(logs []*[raft](/github.com/hashicorp/raft).[Log](/github.com/hashicorp/raft#Log)) [error](/builtin#error)

StoreLogs is used to store a set of raft logs 

####  func (*BoltStore) [Sync](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L309) ¶
    
    
    func (b *BoltStore) Sync() [error](/builtin#error)

Sync performs an fsync on the database file handle. This is not necessary under normal operation unless NoSync is enabled, in which this forces the database file to sync against the disk. 

####  type [Options](https://github.com/hashicorp/raft-boltdb/blob/v2.3.1/v2/bolt_store.go#L47) ¶
    
    
    type Options struct {
    	// Path is the file path to the Bbolt to use
    	Path [string](/builtin#string)
    
    	// BoltOptions contains any specific Bbolt options you might
    	// want to specify [e.g. open timeout]
    	BoltOptions *[bbolt](/go.etcd.io/bbolt).[Options](/go.etcd.io/bbolt#Options)
    
    	// NoSync causes the database to skip fsync calls after each
    	// write to the log. This is unsafe, so it should be used
    	// with caution.
    	NoSync [bool](/builtin#bool)
    
    	// MsgpackUseNewTimeFormat when set to true, force the underlying msgpack
    	// codec to use the new format of time.Time when encoding (used in
    	// go-msgpack v1.1.5 by default). Decoding is not affected, as all
    	// go-msgpack v2.1.0+ decoders know how to decode both formats.
    	MsgpackUseNewTimeFormat [bool](/builtin#bool)
    }

Options contains all the configuration used to open the Bbolt 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
