# hashicorp/raft-boltdb

> Source: https://pkg.go.dev/github.com/hashicorp/raft-boltdb/v2
> Fetched: 2026-01-30T23:56:00.080794+00:00
> Content-Hash: 5a2db7dbbff45d8a
> Type: html

---

Index

¶

Variables

type BoltStore

func MigrateToV2(source, destination string) (*BoltStore, error)

func New(options Options) (*BoltStore, error)

func NewBoltStore(path string) (*BoltStore, error)

func (b *BoltStore) Close() error

func (b *BoltStore) DeleteRange(min, max uint64) error

func (b *BoltStore) FirstIndex() (uint64, error)

func (b *BoltStore) Get(k []byte) ([]byte, error)

func (b *BoltStore) GetLog(idx uint64, log *raft.Log) error

func (b *BoltStore) GetUint64(key []byte) (uint64, error)

func (b *BoltStore) LastIndex() (uint64, error)

func (b *BoltStore) RunMetrics(ctx context.Context, interval time.Duration)

func (b *BoltStore) Set(k, v []byte) error

func (b *BoltStore) SetUint64(key []byte, val uint64) error

func (b *BoltStore) Stats() bbolt.Stats

func (b *BoltStore) StoreLog(log *raft.Log) error

func (b *BoltStore) StoreLogs(logs []*raft.Log) error

func (b *BoltStore) Sync() error

type Options

Constants

¶

This section is empty.

Variables

¶

View Source

var (

// An error indicating a given key does not exist

ErrKeyNotFound =

errors

.

New

("not found")
)

Functions

¶

This section is empty.

Types

¶

type

BoltStore

¶

type BoltStore struct {

// contains filtered or unexported fields

}

BoltStore provides access to Bbolt for Raft to store and retrieve
log entries. It also provides key/value storage, and can be used as
a LogStore and StableStore.

func

MigrateToV2

¶

func MigrateToV2(source, destination

string

) (*

BoltStore

,

error

)

MigrateToV2 reads in the source file path of a BoltDB file
and outputs all the data migrated to a Bbolt destination file

func

New

¶

func New(options

Options

) (*

BoltStore

,

error

)

New uses the supplied options to open the Bbolt and prepare it for use as a raft backend.

func

NewBoltStore

¶

func NewBoltStore(path

string

) (*

BoltStore

,

error

)

NewBoltStore takes a file path and returns a connected Raft backend.

func (*BoltStore)

Close

¶

func (b *

BoltStore

) Close()

error

Close is used to gracefully close the DB connection.

func (*BoltStore)

DeleteRange

¶

func (b *

BoltStore

) DeleteRange(min, max

uint64

)

error

DeleteRange is used to delete logs within a given range inclusively.

func (*BoltStore)

FirstIndex

¶

func (b *

BoltStore

) FirstIndex() (

uint64

,

error

)

FirstIndex returns the first known index from the Raft log.

func (*BoltStore)

Get

¶

func (b *

BoltStore

) Get(k []

byte

) ([]

byte

,

error

)

Get is used to retrieve a value from the k/v store by key

func (*BoltStore)

GetLog

¶

func (b *

BoltStore

) GetLog(idx

uint64

, log *

raft

.

Log

)

error

GetLog is used to retrieve a log from Bbolt at a given index.

func (*BoltStore)

GetUint64

¶

func (b *

BoltStore

) GetUint64(key []

byte

) (

uint64

,

error

)

GetUint64 is like Get, but handles uint64 values

func (*BoltStore)

LastIndex

¶

func (b *

BoltStore

) LastIndex() (

uint64

,

error

)

LastIndex returns the last known index from the Raft log.

func (*BoltStore)

RunMetrics

¶

func (b *

BoltStore

) RunMetrics(ctx

context

.

Context

, interval

time

.

Duration

)

RunMetrics should be executed in a go routine and will periodically emit
metrics on the given interval until the context has been cancelled.

func (*BoltStore)

Set

¶

func (b *

BoltStore

) Set(k, v []

byte

)

error

Set is used to set a key/value set outside of the raft log

func (*BoltStore)

SetUint64

¶

func (b *

BoltStore

) SetUint64(key []

byte

, val

uint64

)

error

SetUint64 is like Set, but handles uint64 values

func (*BoltStore)

Stats

¶

func (b *

BoltStore

) Stats()

bbolt

.

Stats

func (*BoltStore)

StoreLog

¶

func (b *

BoltStore

) StoreLog(log *

raft

.

Log

)

error

StoreLog is used to store a single raft log

func (*BoltStore)

StoreLogs

¶

func (b *

BoltStore

) StoreLogs(logs []*

raft

.

Log

)

error

StoreLogs is used to store a set of raft logs

func (*BoltStore)

Sync

¶

func (b *

BoltStore

) Sync()

error

Sync performs an fsync on the database file handle. This is not necessary
under normal operation unless NoSync is enabled, in which this forces the
database file to sync against the disk.

type

Options

¶

type Options struct {

// Path is the file path to the Bbolt to use

Path

string

// BoltOptions contains any specific Bbolt options you might

// want to specify [e.g. open timeout]

BoltOptions *

bbolt

.

Options

// NoSync causes the database to skip fsync calls after each

// write to the log. This is unsafe, so it should be used

// with caution.

NoSync

bool

// MsgpackUseNewTimeFormat when set to true, force the underlying msgpack

// codec to use the new format of time.Time when encoding (used in

// go-msgpack v1.1.5 by default). Decoding is not affected, as all

// go-msgpack v2.1.0+ decoders know how to decode both formats.

MsgpackUseNewTimeFormat

bool

}

Options contains all the configuration used to open the Bbolt