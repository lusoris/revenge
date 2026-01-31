# otter Cache

> Source: https://pkg.go.dev/github.com/maypok86/otter/v2
> Fetched: 2026-01-30T23:48:44.386233+00:00
> Content-Hash: a1dc7fbc8e8991af
> Type: html

---

Overview

¶

Package otter contains in-memory caching functionality.

A

Cache

is similar to a hash table, but it also has additional support for policies to bound the map.

Cache

instances should always be configured and created using

Options

.

The

Cache

also has

Cache.Get

/

Cache.BulkGet

/

Cache.Refresh

/

Cache.Refresh

methods
which allows the cache to populate itself on a miss and offers refresh capabilities.

Additional functionality such as bounding by the entry's size, deletion notifications, statistics,
and eviction policies are described in the

Options

.

See

https://maypok86.github.io/otter/user-guide/v2/getting-started/

for more information about otter.

Index

¶

Constants

func LoadCacheFrom[K comparable, V any](c *Cache[K, V], r io.Reader) error

func LoadCacheFromFile[K comparable, V any](c *Cache[K, V], filePath string) error

func SaveCacheTo[K comparable, V any](c *Cache[K, V], w io.Writer) error

func SaveCacheToFile[K comparable, V any](c *Cache[K, V], filePath string) error

type BulkLoader

type BulkLoaderFunc

func (blf BulkLoaderFunc[K, V]) BulkLoad(ctx context.Context, keys []K) (map[K]V, error)

func (blf BulkLoaderFunc[K, V]) BulkReload(ctx context.Context, keys []K, oldValues []V) (map[K]V, error)

type Cache

func Must[K comparable, V any](o *Options[K, V]) *Cache[K, V]

func New[K comparable, V any](o *Options[K, V]) (*Cache[K, V], error)

func (c *Cache[K, V]) All() iter.Seq2[K, V]

func (c *Cache[K, V]) BulkGet(ctx context.Context, keys []K, bulkLoader BulkLoader[K, V]) (map[K]V, error)

func (c *Cache[K, V]) BulkRefresh(ctx context.Context, keys []K, bulkLoader BulkLoader[K, V]) <-chan []RefreshResult[K, V]

func (c *Cache[K, V]) CleanUp()

func (c *Cache[K, V]) Coldest() iter.Seq[Entry[K, V]]

func (c *Cache[K, V]) Compute(key K, remappingFunc func(oldValue V, found bool) (newValue V, op ComputeOp)) (actualValue V, ok bool)

func (c *Cache[K, V]) ComputeIfAbsent(key K, mappingFunc func() (newValue V, cancel bool)) (actualValue V, ok bool)

func (c *Cache[K, V]) ComputeIfPresent(key K, remappingFunc func(oldValue V) (newValue V, op ComputeOp)) (actualValue V, ok bool)

func (c *Cache[K, V]) EstimatedSize() int

func (c *Cache[K, V]) Get(ctx context.Context, key K, loader Loader[K, V]) (V, error)

func (c *Cache[K, V]) GetEntry(key K) (Entry[K, V], bool)

func (c *Cache[K, V]) GetEntryQuietly(key K) (Entry[K, V], bool)

func (c *Cache[K, V]) GetIfPresent(key K) (V, bool)

func (c *Cache[K, V]) GetMaximum() uint64

func (c *Cache[K, V]) Hottest() iter.Seq[Entry[K, V]]

func (c *Cache[K, V]) Invalidate(key K) (value V, invalidated bool)

func (c *Cache[K, V]) InvalidateAll()

func (c *Cache[K, V]) IsRecordingStats() bool

func (c *Cache[K, V]) IsWeighted() bool

func (c *Cache[K, V]) Keys() iter.Seq[K]

func (c *Cache[K, V]) Refresh(ctx context.Context, key K, loader Loader[K, V]) <-chan RefreshResult[K, V]

func (c *Cache[K, V]) Set(key K, value V) (V, bool)

func (c *Cache[K, V]) SetExpiresAfter(key K, expiresAfter time.Duration)

func (c *Cache[K, V]) SetIfAbsent(key K, value V) (V, bool)

func (c *Cache[K, V]) SetMaximum(maximum uint64)

func (c *Cache[K, V]) SetRefreshableAfter(key K, refreshableAfter time.Duration)

func (c *Cache[K, V]) Stats() stats.Stats

func (c *Cache[K, V]) StopAllGoroutines() bool

func (c *Cache[K, V]) Values() iter.Seq[V]

func (c *Cache[K, V]) WeightedSize() uint64

type Clock

type ComputeOp

func (co ComputeOp) String() string

type DeletionCause

func (dc DeletionCause) IsEviction() bool

func (dc DeletionCause) String() string

type DeletionEvent

func (de DeletionEvent[K, V]) WasEvicted() bool

type Entry

func (e Entry[K, V]) ExpiresAfter() time.Duration

func (e Entry[K, V]) ExpiresAt() time.Time

func (e Entry[K, V]) HasExpired() bool

func (e Entry[K, V]) RefreshableAfter() time.Duration

func (e Entry[K, V]) RefreshableAt() time.Time

func (e Entry[K, V]) SnapshotAt() time.Time

type ExpiryCalculator

func ExpiryAccessing[K comparable, V any](duration time.Duration) ExpiryCalculator[K, V]

func ExpiryAccessingFunc[K comparable, V any](f func(entry Entry[K, V]) time.Duration) ExpiryCalculator[K, V]

func ExpiryCreating[K comparable, V any](duration time.Duration) ExpiryCalculator[K, V]

func ExpiryCreatingFunc[K comparable, V any](f func(entry Entry[K, V]) time.Duration) ExpiryCalculator[K, V]

func ExpiryWriting[K comparable, V any](duration time.Duration) ExpiryCalculator[K, V]

func ExpiryWritingFunc[K comparable, V any](f func(entry Entry[K, V]) time.Duration) ExpiryCalculator[K, V]

type Loader

type LoaderFunc

func (lf LoaderFunc[K, V]) Load(ctx context.Context, key K) (V, error)

func (lf LoaderFunc[K, V]) Reload(ctx context.Context, key K, oldValue V) (V, error)

type Logger

type NoopLogger

func (nl *NoopLogger) Error(ctx context.Context, msg string, err error)

func (nl *NoopLogger) Warn(ctx context.Context, msg string, err error)

type Options

type RefreshCalculator

func RefreshCreating[K comparable, V any](duration time.Duration) RefreshCalculator[K, V]

func RefreshCreatingFunc[K comparable, V any](f func(entry Entry[K, V]) time.Duration) RefreshCalculator[K, V]

func RefreshWriting[K comparable, V any](duration time.Duration) RefreshCalculator[K, V]

func RefreshWritingFunc[K comparable, V any](f func(entry Entry[K, V]) time.Duration) RefreshCalculator[K, V]

type RefreshResult

Constants

¶

View Source

const (

// ErrNotFound should be returned from a Loader.Load/Loader.Reload to indicate that an entry is

// missing at the underlying data source. This helps the cache to determine

// if an entry should be deleted.

//

// NOTE: this only applies to Cache.Get/Cache.Refresh/Loader.Load/Loader.Reload. For Cache.BulkGet/Cache.BulkRefresh,

// this works implicitly if you return a map without the key.

ErrNotFound strError = "otter: the entry was not found in the data source"
)

Variables

¶

This section is empty.

Functions

¶

func

LoadCacheFrom

¶

added in

v2.1.0

func LoadCacheFrom[K

comparable

, V

any

](c *

Cache

[K, V], r

io

.

Reader

)

error

LoadCacheFrom loads cache data from the given

io.Reader

.

See SaveCacheToFile for saving cache data to file.

func

LoadCacheFromFile

¶

added in

v2.1.0

func LoadCacheFromFile[K

comparable

, V

any

](c *

Cache

[K, V], filePath

string

)

error

LoadCacheFromFile loads cache data from the given filePath.

See SaveCacheToFile for saving cache data to file.

func

SaveCacheTo

¶

added in

v2.1.0

func SaveCacheTo[K

comparable

, V

any

](c *

Cache

[K, V], w

io

.

Writer

)

error

SaveCacheTo atomically saves cache data to the given

io.Writer

.

SaveCacheToFile may be called concurrently with other operations on the cache.

The saved data may be loaded with LoadCacheFrom.

WARNING: Beware that this operation is performed within the eviction policy's exclusive lock.
While the operation is in progress further eviction maintenance will be halted.

func

SaveCacheToFile

¶

added in

v2.1.0

func SaveCacheToFile[K

comparable

, V

any

](c *

Cache

[K, V], filePath

string

)

error

SaveCacheToFile atomically saves cache data to the given filePath.

SaveCacheToFile may be called concurrently with other operations on the cache.

The saved data may be loaded with LoadCacheFromFile.

WARNING: Beware that this operation is performed within the eviction policy's exclusive lock.
While the operation is in progress further eviction maintenance will be halted.

Types

¶

type

BulkLoader

¶

type BulkLoader[K

comparable

, V

any

] interface {

// BulkLoad computes or retrieves the values corresponding to keys.

// This method is called by Cache.BulkGet.

//

// If the returned map doesn't contain all requested keys, then the entries it does

// contain will be cached, and Cache.BulkGet will return the partial results. If the returned map

// contains extra keys not present in keys then all returned entries will be cached, but

// only the entries for keys, will be returned from Cache.BulkGet.

//

// WARNING: loading must not attempt to update any mappings of this cache directly.

BulkLoad(ctx

context

.

Context

, keys []K) (map[K]V,

error

)

// BulkReload computes or retrieves replacement values corresponding to already-cached keys.

// If the replacement value is not found, then the mapping will be removed.

// This method is called when an existing cache entry is refreshed by Cache.BulkGet, or through a call to Cache.BulkRefresh.

//

// If the returned map doesn't contain all requested keys, then the entries it does

// contain will be cached. If the returned map

// contains extra keys not present in keys then all returned entries will be cached.

//

// WARNING: loading must not attempt to update any mappings of this cache directly

// or block waiting for other cache operations to complete.

//

// NOTE: all errors returned by this method will be logged (using Logger) and then swallowed.

BulkReload(ctx

context

.

Context

, keys []K, oldValues []V) (map[K]V,

error

)
}

BulkLoader computes or retrieves values, based on the keys, for use in populating a

Cache

.

type

BulkLoaderFunc

¶

type BulkLoaderFunc[K

comparable

, V

any

] func(ctx

context

.

Context

, keys []K) (map[K]V,

error

)

BulkLoaderFunc is an adapter to allow the use of ordinary functions as loaders.
If f is a function with the appropriate signature, BulkLoaderFunc(f) is a

BulkLoader

that calls f.

func (BulkLoaderFunc[K, V])

BulkLoad

¶

func (blf

BulkLoaderFunc

[K, V]) BulkLoad(ctx

context

.

Context

, keys []K) (map[K]V,

error

)

BulkLoad calls f(ctx, keys).

func (BulkLoaderFunc[K, V])

BulkReload

¶

func (blf

BulkLoaderFunc

[K, V]) BulkReload(ctx

context

.

Context

, keys []K, oldValues []V) (map[K]V,

error

)

BulkReload calls f(ctx, keys).

type

Cache

¶

type Cache[K

comparable

, V

any

] struct {

// contains filtered or unexported fields

}

Cache is an in-memory cache implementation that supports full concurrency of retrievals and multiple ways to bound the cache.

func

Must

¶

func Must[K

comparable

, V

any

](o *

Options

[K, V]) *

Cache

[K, V]

Must creates a configured

Cache

instance or
panics if invalid parameters were specified.

This method does not alter the state of the

Options

instance, so it can be invoked
again to create multiple independent caches.

func

New

¶

func New[K

comparable

, V

any

](o *

Options

[K, V]) (*

Cache

[K, V],

error

)

New creates a configured

Cache

instance or
returns an error if invalid parameters were specified.

This method does not alter the state of the

Options

instance, so it can be invoked
again to create multiple independent caches.

func (*Cache[K, V])

All

¶

func (c *

Cache

[K, V]) All()

iter

.

Seq2

[K, V]

All returns an iterator over all key-value pairs in the cache.
The iteration order is not specified and is not guaranteed to be the same from one call to the next.

Iterator is at least weakly consistent: he is safe for concurrent use,
but if the cache is modified (including by eviction) after the iterator is
created, it is undefined which of the changes (if any) will be reflected in that iterator.

func (*Cache[K, V])

BulkGet

¶

func (c *

Cache

[K, V]) BulkGet(ctx

context

.

Context

, keys []K, bulkLoader

BulkLoader

[K, V]) (map[K]V,

error

)

BulkGet returns the value associated with key in this cache, obtaining that value from loader if necessary.
The method improves upon the conventional "if cached, return; otherwise create, cache and return" pattern.

If another call to Get (BulkGet) is currently loading the value for key,
simply waits for that goroutine to finish and returns its loaded value. Note that
multiple goroutines can concurrently load values for distinct keys.

No observable state associated with this cache is modified until loading completes.

NOTE: duplicate elements in keys will be ignored.

WARNING: When performing a refresh (see

RefreshCalculator

),
the

BulkLoader

will receive a context wrapped in

context.WithoutCancel

.
If you need to control refresh cancellation, you can use closures or values stored in the context.

WARNING:

BulkLoader

must not attempt to update any mappings of this cache directly.

WARNING: For any given key, every bulkLoader used with it should compute the same value.
Otherwise, a call that passes one bulkLoader may return the result of another call
with a differently behaving bulkLoader. For example, a call that requests a short timeout
for an RPC may wait for a similar call that requests a long timeout, or a call by an
unprivileged user may return a resource accessible only to a privileged user making a similar call.

func (*Cache[K, V])

BulkRefresh

¶

func (c *

Cache

[K, V]) BulkRefresh(ctx

context

.

Context

, keys []K, bulkLoader

BulkLoader

[K, V]) <-chan []

RefreshResult

[K, V]

BulkRefresh loads a new value for each key, asynchronously. While the new value is loading the
previous value (if any) will continue to be returned by any Get unless it is evicted.
If the new value is loaded successfully, it will replace the previous value in the cache;
If refreshing returned an error, the previous value will remain,
and the error will be logged using

Logger

and swallowed. If another goroutine is currently
loading the value for key, then this method does not perform an additional load.

Cache

will call BulkLoader.BulkReload for existing keys, and BulkLoader.BulkLoad otherwise.
Loading is asynchronous by delegating to the configured Executor.

BulkRefresh returns a channel that will receive the results when they are ready. The returned channel will not be closed.

NOTE: duplicate elements in keys will be ignored.

WARNING: When performing a refresh (see

RefreshCalculator

),
the

BulkLoader

will receive a context wrapped in

context.WithoutCancel

.
If you need to control refresh cancellation, you can use closures or values stored in the context.

WARNING: If the cache was constructed without

RefreshCalculator

, then BulkRefresh will return the nil channel.

WARNING: BulkLoader.BulkLoad and BulkLoader.BulkReload must not attempt to update any mappings of this cache directly.

WARNING: For any given key, every bulkLoader used with it should compute the same value.
Otherwise, a call that passes one bulkLoader may return the result of another call
with a differently behaving loader. For example, a call that requests a short timeout
for an RPC may wait for a similar call that requests a long timeout, or a call by an
unprivileged user may return a resource accessible only to a privileged user making a similar call.

func (*Cache[K, V])

CleanUp

¶

func (c *

Cache

[K, V]) CleanUp()

CleanUp performs any pending maintenance operations needed by the cache. Exactly which activities are
performed -- if any -- is implementation-dependent.

func (*Cache[K, V])

Coldest

¶

added in

v2.1.0

func (c *

Cache

[K, V]) Coldest()

iter

.

Seq

[

Entry

[K, V]]

Coldest returns an iterator for ordered traversal of the cache entries. The order of
iteration is from the entries least likely to be retained (coldest) to the entries most
likely to be retained (hottest). This order is determined by the eviction policy's best guess
at the start of the iteration.

WARNING: Beware that this iteration is performed within the eviction policy's exclusive lock, so the
iteration should be short and simple. While the iteration is in progress further eviction
maintenance will be halted.

func (*Cache[K, V])

Compute

¶

added in

v2.1.0

func (c *

Cache

[K, V]) Compute(
	key K,
	remappingFunc func(oldValue V, found

bool

) (newValue V, op

ComputeOp

),
) (actualValue V, ok

bool

)

Compute either sets the computed new value for the key,
invalidates the value for the key, or does nothing, based on
the returned

ComputeOp

. When the op returned by remappingFunc
is

WriteOp

, the value is updated to the new value. If
it is

InvalidateOp

, the entry is removed from the cache
altogether. And finally, if the op is

CancelOp

then the
entry is left as-is. In other words, if it did not already
exist, it is not created, and if it did exist, it is not
updated. This is useful to synchronously execute some
operation on the value without incurring the cost of
updating the cache every time.

The ok result indicates whether the entry is present in the cache after the compute operation.
The actualValue result contains the value of the cache
if a corresponding entry is present, or the zero value
otherwise. You can think of these results as equivalent to regular key-value lookups in a map.

This call locks a hash table bucket while the compute function
is executed. It means that modifications on other entries in
the bucket will be blocked until the remappingFunc executes. Consider
this when the function includes long-running operations.

func (*Cache[K, V])

ComputeIfAbsent

¶

added in

v2.1.0

func (c *

Cache

[K, V]) ComputeIfAbsent(
	key K,
	mappingFunc func() (newValue V, cancel

bool

),
) (actualValue V, ok

bool

)

ComputeIfAbsent returns the existing value for the key if
present. Otherwise, it tries to compute the value using the
provided function. If mappingFunc returns true as the cancel value, the computation is cancelled and the zero value
for type V is returned.

The ok result indicates whether the entry is present in the cache after the compute operation.
The actualValue result contains the value of the cache
if a corresponding entry is present, or the zero value
otherwise. You can think of these results as equivalent to regular key-value lookups in a map.

This call locks a hash table bucket while the compute function
is executed. It means that modifications on other entries in
the bucket will be blocked until the valueFn executes. Consider
this when the function includes long-running operations.

func (*Cache[K, V])

ComputeIfPresent

¶

added in

v2.1.0

func (c *

Cache

[K, V]) ComputeIfPresent(
	key K,
	remappingFunc func(oldValue V) (newValue V, op

ComputeOp

),
) (actualValue V, ok

bool

)

ComputeIfPresent returns the zero value for type V if the key is not found.
Otherwise, it tries to compute the value using the provided function.

ComputeIfPresent either sets the computed new value for the key,
invalidates the value for the key, or does nothing, based on
the returned

ComputeOp

. When the op returned by remappingFunc
is

WriteOp

, the value is updated to the new value. If
it is

InvalidateOp

, the entry is removed from the cache
altogether. And finally, if the op is

CancelOp

then the
entry is left as-is. In other words, if it did not already
exist, it is not created, and if it did exist, it is not
updated. This is useful to synchronously execute some
operation on the value without incurring the cost of
updating the cache every time.

The ok result indicates whether the entry is present in the cache after the compute operation.
The actualValue result contains the value of the cache
if a corresponding entry is present, or the zero value
otherwise. You can think of these results as equivalent to regular key-value lookups in a map.

This call locks a hash table bucket while the compute function
is executed. It means that modifications on other entries in
the bucket will be blocked until the valueFn executes. Consider
this when the function includes long-running operations.

func (*Cache[K, V])

EstimatedSize

¶

func (c *

Cache

[K, V]) EstimatedSize()

int

EstimatedSize returns the approximate number of entries in this cache. The value returned is an estimate; the
actual count may differ if there are concurrent insertions or deletions, or if some entries are
pending deletion due to expiration. In the case of stale entries
this inaccuracy can be mitigated by performing a CleanUp first.

func (*Cache[K, V])

Get

¶

func (c *

Cache

[K, V]) Get(ctx

context

.

Context

, key K, loader

Loader

[K, V]) (V,

error

)

Get returns the value associated with key in this cache, obtaining that value from loader if necessary.
The method improves upon the conventional "if cached, return; otherwise create, cache and return" pattern.

Get can return an

ErrNotFound

error if the

Loader

returns it.
This means that the entry was not found in the data source.
This allows the cache to recognize when a record is missing from the data source
and subsequently delete the cached entry.
It also enables proper metric collection, as the cache doesn't classify

ErrNotFound

as a load error.

If another call to Get is currently loading the value for key,
simply waits for that goroutine to finish and returns its loaded value. Note that
multiple goroutines can concurrently load values for distinct keys.

No observable state associated with this cache is modified until loading completes.

WARNING: When performing a refresh (see

RefreshCalculator

),
the

Loader

will receive a context wrapped in

context.WithoutCancel

.
If you need to control refresh cancellation, you can use closures or values stored in the context.

WARNING:

Loader

must not attempt to update any mappings of this cache directly.

WARNING: For any given key, every loader used with it should compute the same value.
Otherwise, a call that passes one loader may return the result of another call
with a differently behaving loader. For example, a call that requests a short timeout
for an RPC may wait for a similar call that requests a long timeout, or a call by an
unprivileged user may return a resource accessible only to a privileged user making a similar call.

func (*Cache[K, V])

GetEntry

¶

func (c *

Cache

[K, V]) GetEntry(key K) (

Entry

[K, V],

bool

)

GetEntry returns the cache entry associated with the key in this cache.

func (*Cache[K, V])

GetEntryQuietly

¶

func (c *

Cache

[K, V]) GetEntryQuietly(key K) (

Entry

[K, V],

bool

)

GetEntryQuietly returns the cache entry associated with the key in this cache.

Unlike GetEntry, this function does not produce any side effects
such as updating statistics or the eviction policy.

func (*Cache[K, V])

GetIfPresent

¶

func (c *

Cache

[K, V]) GetIfPresent(key K) (V,

bool

)

GetIfPresent returns the value associated with the key in this cache.

func (*Cache[K, V])

GetMaximum

¶

func (c *

Cache

[K, V]) GetMaximum()

uint64

GetMaximum returns the maximum total weighted or unweighted size of this cache, depending on how the
cache was constructed. If this cache does not use a (weighted) size bound, then the method will return math.MaxUint64.

func (*Cache[K, V])

Hottest

¶

added in

v2.1.0

func (c *

Cache

[K, V]) Hottest()

iter

.

Seq

[

Entry

[K, V]]

Hottest returns an iterator for ordered traversal of the cache entries. The order of
iteration is from the entries most likely to be retained (hottest) to the entries least
likely to be retained (coldest). This order is determined by the eviction policy's best guess
at the start of the iteration.

WARNING: Beware that this iteration is performed within the eviction policy's exclusive lock, so the
iteration should be short and simple. While the iteration is in progress further eviction
maintenance will be halted.

func (*Cache[K, V])

Invalidate

¶

func (c *

Cache

[K, V]) Invalidate(key K) (value V, invalidated

bool

)

Invalidate discards any cached value for the key.

Returns previous value if any. The invalidated result reports whether the key was
present.

func (*Cache[K, V])

InvalidateAll

¶

func (c *

Cache

[K, V]) InvalidateAll()

InvalidateAll discards all entries in the cache. The behavior of this operation is undefined for an entry
that is being loaded (or reloaded) and is otherwise not present.

func (*Cache[K, V])

IsRecordingStats

¶

added in

v2.2.0

func (c *

Cache

[K, V]) IsRecordingStats()

bool

IsRecordingStats returns whether the cache statistics are being accumulated.

func (*Cache[K, V])

IsWeighted

¶

added in

v2.2.0

func (c *

Cache

[K, V]) IsWeighted()

bool

IsWeighted returns whether the cache is bounded by a maximum size or maximum weight.

func (*Cache[K, V])

Keys

¶

added in

v2.1.0

func (c *

Cache

[K, V]) Keys()

iter

.

Seq

[K]

Keys returns an iterator over all keys in the cache.
The iteration order is not specified and is not guaranteed to be the same from one call to the next.

Iterator is at least weakly consistent: he is safe for concurrent use,
but if the cache is modified (including by eviction) after the iterator is
created, it is undefined which of the changes (if any) will be reflected in that iterator.

func (*Cache[K, V])

Refresh

¶

func (c *

Cache

[K, V]) Refresh(ctx

context

.

Context

, key K, loader

Loader

[K, V]) <-chan

RefreshResult

[K, V]

Refresh loads a new value for the key, asynchronously. While the new value is loading the
previous value (if any) will continue to be returned by any Get unless it is evicted.
If the new value is loaded successfully, it will replace the previous value in the cache;
If refreshing returned an error, the previous value will remain,
and the error will be logged using

Logger

(if it's not

ErrNotFound

) and swallowed. If another goroutine is currently
loading the value for key, then this method does not perform an additional load.

Cache

will call Loader.Reload if the cache currently contains a value for the key,
and Loader.Load otherwise.
Loading is asynchronous by delegating to the configured Executor.

Refresh returns a channel that will receive the result when it is ready. The returned channel will not be closed.

WARNING: When performing a refresh (see

RefreshCalculator

),
the

Loader

will receive a context wrapped in

context.WithoutCancel

.
If you need to control refresh cancellation, you can use closures or values stored in the context.

WARNING: If the cache was constructed without

RefreshCalculator

, then Refresh will return the nil channel.

WARNING: Loader.Load and Loader.Reload must not attempt to update any mappings of this cache directly.

WARNING: For any given key, every loader used with it should compute the same value.
Otherwise, a call that passes one loader may return the result of another call
with a differently behaving loader. For example, a call that requests a short timeout
for an RPC may wait for a similar call that requests a long timeout, or a call by an
unprivileged user may return a resource accessible only to a privileged user making a similar call.

func (*Cache[K, V])

Set

¶

func (c *

Cache

[K, V]) Set(key K, value V) (V,

bool

)

Set associates the value with the key in this cache.

If the specified key is not already associated with a value, then it returns new value and true.

If the specified key is already associated with a value, then it returns existing value and false.

func (*Cache[K, V])

SetExpiresAfter

¶

func (c *

Cache

[K, V]) SetExpiresAfter(key K, expiresAfter

time

.

Duration

)

SetExpiresAfter specifies that the entry should be automatically removed from the cache once the duration has
elapsed. The expiration policy determines when the entry's age is reset.

func (*Cache[K, V])

SetIfAbsent

¶

func (c *

Cache

[K, V]) SetIfAbsent(key K, value V) (V,

bool

)

SetIfAbsent if the specified key is not already associated with a value associates it with the given value.

If the specified key is not already associated with a value, then it returns new value and true.

If the specified key is already associated with a value, then it returns existing value and false.

func (*Cache[K, V])

SetMaximum

¶

func (c *

Cache

[K, V]) SetMaximum(maximum

uint64

)

SetMaximum specifies the maximum total size of this cache. This value may be interpreted as the weighted
or unweighted threshold size based on how this cache was constructed. If the cache currently
exceeds the new maximum size this operation eagerly evict entries until the cache shrinks to
the appropriate size.

func (*Cache[K, V])

SetRefreshableAfter

¶

func (c *

Cache

[K, V]) SetRefreshableAfter(key K, refreshableAfter

time

.

Duration

)

SetRefreshableAfter specifies that each entry should be eligible for reloading once a fixed duration has elapsed.
The refresh policy determines when the entry's age is reset.

func (*Cache[K, V])

Stats

¶

added in

v2.2.0

func (c *

Cache

[K, V]) Stats()

stats

.

Stats

Stats returns a current snapshot of this cache's cumulative statistics.
All statistics are initialized to zero and are monotonically increasing over the lifetime of the cache.
Due to the performance penalty of maintaining statistics,
some implementations may not record the usage history immediately or at all.

NOTE: If your

stats.Recorder

implementation doesn't also implement

stats.Snapshoter

,
this method will always return a zero-value snapshot.

func (*Cache[K, V])

StopAllGoroutines

¶

added in

v2.3.0

func (c *

Cache

[K, V]) StopAllGoroutines()

bool

StopAllGoroutines stops all goroutines launched by the cache.
It returns true if the call stops goroutines, false if goroutines have already been stopped.

NOTE: In the vast majority of cases, you do not need to call this method, and it will be called automatically.
However, it can sometimes be useful, for example when using synctest.
You can think of this method as analogous to the Stop method of

time.Timer

.

NOTE: This method only stops the goroutines and does not invalidate entries in the cache.
To invalidate entries, you can use the Invalidate / InvalidateAll methods.

func (*Cache[K, V])

Values

¶

added in

v2.1.0

func (c *

Cache

[K, V]) Values()

iter

.

Seq

[V]

Values returns an iterator over all values in the cache.
The iteration order is not specified and is not guaranteed to be the same from one call to the next.

Iterator is at least weakly consistent: he is safe for concurrent use,
but if the cache is modified (including by eviction) after the iterator is
created, it is undefined which of the changes (if any) will be reflected in that iterator.

func (*Cache[K, V])

WeightedSize

¶

func (c *

Cache

[K, V]) WeightedSize()

uint64

WeightedSize returns the approximate accumulated weight of entries in this cache. If this cache does not
use a weighted size bound, then the method will return 0.

type

Clock

¶

added in

v2.1.0

type Clock interface {

// NowNano returns the number of nanoseconds elapsed since this clock's fixed point of reference.

//

// By default, time.Now().UnixNano() is used.

NowNano()

int64

// Tick returns a channel that delivers “ticks” of a clock at intervals.

//

// The cache uses this method only for proactive expiration and calls Tick(time.Second) in a separate goroutine.

//

// By default, [time.Tick] is used.

Tick(duration

time

.

Duration

) <-chan

time

.

Time

}

Clock is a time source that

Returns a time value representing the number of nanoseconds elapsed since some
fixed but arbitrary point in time

Returns a channel that delivers “ticks” of a clock at intervals.

type

ComputeOp

¶

added in

v2.1.0

type ComputeOp

int

ComputeOp tells the Compute methods what to do.

const (

// CancelOp signals to Compute to not do anything as a result

// of executing the lambda. If the entry was not present in

// the map, nothing happens, and if it was present, the

// returned value is ignored.

CancelOp

ComputeOp

=

iota

// WriteOp signals to Compute to update the entry to the

// value returned by the lambda, creating it if necessary.

WriteOp

// InvalidateOp signals to Compute to always discard the entry

// from the cache.

InvalidateOp
)

func (ComputeOp)

String

¶

added in

v2.1.0

func (co

ComputeOp

) String()

string

String implements

fmt.Stringer

interface.

type

DeletionCause

¶

type DeletionCause

int

DeletionCause the cause why a cached entry was deleted.

const (

// CauseInvalidation means that the entry was manually deleted by the user.

CauseInvalidation

DeletionCause

=

iota

+ 1

// CauseReplacement means that the entry itself was not actually deleted, but its value was replaced by the user.

CauseReplacement

// CauseOverflow means that the entry was evicted due to size constraints.

CauseOverflow

// CauseExpiration means that the entry's expiration timestamp has passed.

CauseExpiration
)

func (DeletionCause)

IsEviction

¶

func (dc

DeletionCause

) IsEviction()

bool

IsEviction returns true if there was an automatic deletion due to eviction
(the cause is neither

CauseInvalidation

nor

CauseReplacement

).

func (DeletionCause)

String

¶

func (dc

DeletionCause

) String()

string

String implements

fmt.Stringer

interface.

type

DeletionEvent

¶

type DeletionEvent[K

comparable

, V

any

] struct {

// Key is the key corresponding to the deleted entry.

Key K

// Value is the value corresponding to the deleted entry.

Value V

// Cause is the cause for which entry was deleted.

Cause

DeletionCause

}

DeletionEvent is an event of the deletion of a single entry.

func (DeletionEvent[K, V])

WasEvicted

¶

func (de

DeletionEvent

[K, V]) WasEvicted()

bool

WasEvicted returns true if there was an automatic deletion due to eviction (the cause is neither

CauseInvalidation

nor

CauseReplacement

).

type

Entry

¶

type Entry[K

comparable

, V

any

] struct {

// Key is the entry's key.

Key K

// Value is the entry's value.

Value V

// Weight returns the entry's weight.

//

// If the cache was not configured with a weight then this value is always 1.

Weight

uint32

// ExpiresAtNano is the entry's expiration time as a unix time,

// the number of nanoseconds elapsed since January 1, 1970 UTC.

//

// If the cache was not configured with an expiration policy then this value is always math.MaxInt64.

ExpiresAtNano

int64

// RefreshableAtNano is the time after which the entry will be reloaded as a unix time,

// the number of nanoseconds elapsed since January 1, 1970 UTC.

//

// If the cache was not configured with a refresh policy then this value is always math.MaxInt64.

RefreshableAtNano

int64

// SnapshotAtNano is the time when this snapshot of the entry was taken as a unix time,

// the number of nanoseconds elapsed since January 1, 1970 UTC.

//

// If the cache was not configured with a time-based policy then this value is always 0.

SnapshotAtNano

int64

}

Entry is a key-value pair that may include policy metadata for the cached entry.

It is an immutable snapshot of the cached data at the time of this entry's creation, and it will not
reflect changes afterward.

func (Entry[K, V])

ExpiresAfter

¶

func (e

Entry

[K, V]) ExpiresAfter()

time

.

Duration

ExpiresAfter returns the fixed duration used to determine if an entry should be automatically removed due
to elapsing this time bound. An entry is considered fresh if its age is less than this
duration, and stale otherwise. The expiration policy determines when the entry's age is reset.

If the cache was not configured with an expiration policy then this value is always

math.MaxInt64

.

func (Entry[K, V])

ExpiresAt

¶

func (e

Entry

[K, V]) ExpiresAt()

time

.

Time

ExpiresAt returns the entry's expiration time.

If the cache was not configured with an expiration policy then this value is roughly

math.MaxInt64

nanoseconds away from the SnapshotAt.

func (Entry[K, V])

HasExpired

¶

func (e

Entry

[K, V]) HasExpired()

bool

HasExpired returns true if the entry has expired.

func (Entry[K, V])

RefreshableAfter

¶

func (e

Entry

[K, V]) RefreshableAfter()

time

.

Duration

RefreshableAfter returns the fixed duration used to determine if an entry should be eligible for reloading due
to elapsing this time bound. An entry is considered fresh if its age is less than this
duration, and stale otherwise. The refresh policy determines when the entry's age is reset.

If the cache was not configured with a refresh policy then this value is always

math.MaxInt64

.

func (Entry[K, V])

RefreshableAt

¶

func (e

Entry

[K, V]) RefreshableAt()

time

.

Time

RefreshableAt is the time after which the entry will be reloaded.

If the cache was not configured with a refresh policy then this value is roughly

math.MaxInt64

nanoseconds away from the SnapshotAt.

func (Entry[K, V])

SnapshotAt

¶

func (e

Entry

[K, V]) SnapshotAt()

time

.

Time

SnapshotAt is the time when this snapshot of the entry was taken.

If the cache was not configured with a time-based policy then this value is always 1970-01-01 00:00:00 UTC.

type

ExpiryCalculator

¶

type ExpiryCalculator[K

comparable

, V

any

] interface {

// ExpireAfterCreate specifies that the entry should be automatically removed from the cache once the duration has

// elapsed after the entry's creation. To indicate no expiration, an entry may be given an

// excessively long period.

//

// NOTE: ExpiresAtNano and RefreshableAtNano are not initialized at this stage.

ExpireAfterCreate(entry

Entry

[K, V])

time

.

Duration

// ExpireAfterUpdate specifies that the entry should be automatically removed from the cache once the duration has

// elapsed after the replacement of its value. To indicate no expiration, an entry may be given an

// excessively long period. The entry.ExpiresAfter() may be returned to not modify the expiration time.

ExpireAfterUpdate(entry

Entry

[K, V], oldValue V)

time

.

Duration

// ExpireAfterRead specifies that the entry should be automatically removed from the cache once the duration has

// elapsed after its last read. To indicate no expiration, an entry may be given an excessively

// long period. The entry.ExpiresAfter() may be returned to not modify the expiration time.

ExpireAfterRead(entry

Entry

[K, V])

time

.

Duration

}

ExpiryCalculator calculates when cache entries expire. A single expiration time is retained so that the lifetime
of an entry may be extended or reduced by subsequent evaluations.

func

ExpiryAccessing

¶

func ExpiryAccessing[K

comparable

, V

any

](duration

time

.

Duration

)

ExpiryCalculator

[K, V]

ExpiryAccessing returns an

ExpiryCalculator

that specifies that the entry should be automatically deleted from
the cache once the duration has elapsed after the entry's creation, replacement of its value,
or after it was last read.

func

ExpiryAccessingFunc

¶

func ExpiryAccessingFunc[K

comparable

, V

any

](f func(entry

Entry

[K, V])

time

.

Duration

)

ExpiryCalculator

[K, V]

ExpiryAccessingFunc returns an

ExpiryCalculator

that specifies that the entry should be automatically deleted from
the cache once the duration has elapsed after the entry's creation, replacement of its value,
or after it was last read.

func

ExpiryCreating

¶

func ExpiryCreating[K

comparable

, V

any

](duration

time

.

Duration

)

ExpiryCalculator

[K, V]

ExpiryCreating returns an

ExpiryCalculator

that specifies that the entry should be automatically deleted from
the cache once the duration has elapsed after the entry's creation. The expiration time is
not modified when the entry is updated or read.

func

ExpiryCreatingFunc

¶

func ExpiryCreatingFunc[K

comparable

, V

any

](f func(entry

Entry

[K, V])

time

.

Duration

)

ExpiryCalculator

[K, V]

ExpiryCreatingFunc returns an

ExpiryCalculator

that specifies that the entry should be automatically deleted from
the cache once the duration has elapsed after the entry's creation. The expiration time is
not modified when the entry is updated or read.

func

ExpiryWriting

¶

func ExpiryWriting[K

comparable

, V

any

](duration

time

.

Duration

)

ExpiryCalculator

[K, V]

ExpiryWriting returns an

ExpiryCalculator

that specifies that the entry should be automatically deleted from
the cache once the duration has elapsed after the entry's creation or replacement of its value.
The expiration time is not modified when the entry is read.

func

ExpiryWritingFunc

¶

func ExpiryWritingFunc[K

comparable

, V

any

](f func(entry

Entry

[K, V])

time

.

Duration

)

ExpiryCalculator

[K, V]

ExpiryWritingFunc returns an

ExpiryCalculator

that specifies that the entry should be automatically deleted from
the cache once the duration has elapsed after the entry's creation or replacement of its value.
The expiration time is not modified when the entry is read.

type

Loader

¶

type Loader[K

comparable

, V

any

] interface {

// Load computes or retrieves the value corresponding to key.

//

// WARNING: loading must not attempt to update any mappings of this cache directly.

//

// NOTE: The Loader implementation should always return ErrNotFound

// if the entry was not found in the data source.

Load(ctx

context

.

Context

, key K) (V,

error

)

// Reload computes or retrieves a replacement value corresponding to an already-cached key.

// If the replacement value is not found, then the mapping will be removed if ErrNotFound is returned.

// This method is called when an existing cache entry is refreshed by Cache.Get, or through a call to Cache.Refresh.

//

// WARNING: loading must not attempt to update any mappings of this cache directly

// or block waiting for other cache operations to complete.

//

// NOTE: all errors returned by this method will be logged (using Logger) and then swallowed.

//

// NOTE: The Loader implementation should always return ErrNotFound

// if the entry was not found in the data source.

Reload(ctx

context

.

Context

, key K, oldValue V) (V,

error

)
}

Loader computes or retrieves values, based on a key, for use in populating a

Cache

.

type

LoaderFunc

¶

type LoaderFunc[K

comparable

, V

any

] func(ctx

context

.

Context

, key K) (V,

error

)

LoaderFunc is an adapter to allow the use of ordinary functions as loaders.
If f is a function with the appropriate signature, LoaderFunc(f) is a

Loader

that calls f.

func (LoaderFunc[K, V])

Load

¶

func (lf

LoaderFunc

[K, V]) Load(ctx

context

.

Context

, key K) (V,

error

)

Load calls f(ctx, key).

func (LoaderFunc[K, V])

Reload

¶

func (lf

LoaderFunc

[K, V]) Reload(ctx

context

.

Context

, key K, oldValue V) (V,

error

)

Reload calls f(ctx, key).

type

Logger

¶

type Logger interface {

// Warn logs a message at the warn level with an error.

Warn(ctx

context

.

Context

, msg

string

, err

error

)

// Error logs a message at the error level with an error.

Error(ctx

context

.

Context

, msg

string

, err

error

)
}

Logger is the interface used to get log output from otter.

type

NoopLogger

¶

type NoopLogger struct{}

NoopLogger is a stub implementation of

Logger

interface. It may be useful if error logging is not necessary.

func (*NoopLogger)

Error

¶

func (nl *

NoopLogger

) Error(ctx

context

.

Context

, msg

string

, err

error

)

func (*NoopLogger)

Warn

¶

func (nl *

NoopLogger

) Warn(ctx

context

.

Context

, msg

string

, err

error

)

type

Options

¶

type Options[K

comparable

, V

any

] struct {

// MaximumSize specifies the maximum number of entries the cache may contain.

//

// This option cannot be used in conjunction with MaximumWeight.

//

// NOTE: the cache may evict an entry before this limit is exceeded or temporarily exceed the threshold while evicting.

// As the cache size grows close to the maximum, the cache evicts entries that are less likely to be used again.

// For example, the cache may evict an entry because it hasn't been used recently or very often.

MaximumSize

int

// MaximumWeight specifies the maximum weight of entries the cache may contain. Weight is determined using the

// callback specified with Weigher.

// Use of this method requires specifying an option Weigher prior to calling New.

//

// This option cannot be used in conjunction with MaximumSize.

//

// NOTE: the cache may evict an entry before this limit is exceeded or temporarily exceed the threshold while evicting.

// As the cache size grows close to the maximum, the cache evicts entries that are less likely to be used again.

// For example, the cache may evict an entry because it hasn't been used recently or very often.

//

// NOTE: weight is only used to determine whether the cache is over capacity; it has no effect

// on selecting which entry should be evicted next.

MaximumWeight

uint64

// StatsRecorder accumulates statistics during the operation of a Cache.

//

// NOTE: If your stats.Recorder implementation doesn't also implement stats.Snapshoter,

// Cache.Stats method will always return a zero-value snapshot.

StatsRecorder

stats

.

Recorder

// InitialCapacity specifies the minimum total size for the internal data structures. Providing a large enough estimate

// at construction time avoids the need for expensive resizing operations later, but setting this

// value unnecessarily high wastes memory.

InitialCapacity

int

// Weigher specifies the weigher to use in determining the weight of entries. Entry weight is taken into

// consideration by MaximumWeight when determining which entries to evict, and use

// of this method requires specifying an option MaximumWeight prior to calling New.

// Weights are measured and recorded when entries are inserted into or updated in

// the cache, and are thus effectively static during the lifetime of a cache entry.

//

// When the weight of an entry is zero it will not be considered for size-based eviction (though

// it still may be evicted by other means).

Weigher func(key K, value V)

uint32

// ExpiryCalculator specifies that each entry should be automatically removed from the cache once a duration has

// elapsed after the entry's creation, the most recent replacement of its value, or its last read.

// The expiration time is reset by all cache read and write operations.

ExpiryCalculator

ExpiryCalculator

[K, V]

// OnDeletion specifies a handler instance that caches should notify each time an entry is deleted for any

// DeletionCause reason. The cache will invoke this handler on the configured Executor

// after the entry's deletion operation has completed.

//

// An OnAtomicDeletion may be preferred when the handler should be invoked

// as part of the atomic operation to delete the entry.

OnDeletion func(e

DeletionEvent

[K, V])

// OnAtomicDeletion specifies a handler that caches should notify each time an entry is deleted for any

// DeletionCause. The cache will invoke this handler during the atomic operation to delete the entry.

//

// A OnDeletion may be preferred when the handler should be performed outside the atomic operation to

// delete the entry, or be delegated to the configured Executor.

OnAtomicDeletion func(e

DeletionEvent

[K, V])

// RefreshCalculator specifies that active entries are eligible for automatic refresh once a duration has

// elapsed after the entry's creation, the most recent replacement of its value, or the most recent entry's reload.

// The semantics of refreshes are specified in Cache.Refresh,

// and are performed by calling Loader.Reload in a separate background goroutine.

//

// Automatic refreshes are performed when the first stale request for an entry occurs. The request

// triggering the refresh will make an asynchronous call to Loader.Reload to get a new value.

// Until refresh is completed, requests will continue to return the old value.

//

// NOTE: all errors returned during refresh will be logged (using Logger) and then swallowed.

RefreshCalculator

RefreshCalculator

[K, V]

// Executor specifies the executor to use when running asynchronous tasks. The executor is delegated to

// when sending deletion events, when asynchronous computations are performed by

// Cache.Refresh/Cache.BulkRefresh or for refreshes in Cache.Get/Cache.BulkGet, if RefreshCalculator was specified,

// or when performing periodic maintenance. By default, goroutines are used.

//

// The primary intent of this method is to facilitate testing of caches which have been configured

// with OnDeletion or utilize asynchronous computations. A test may instead prefer

// to configure the cache to execute tasks directly on the same goroutine.

//

// Beware that configuring a cache with an executor that discards tasks or never runs them may

// experience non-deterministic behavior.

Executor func(fn func())

// Clock specifies a nanosecond-precision time source for use in determining when entries should be

// expired or refreshed. By default, time.Now().UnixNano() is used.

//

// The primary intent of this option is to facilitate testing of caches which have been configured

// with ExpiryCalculator or RefreshCalculator.

//

// NOTE: this clock is not used when recording statistics.

Clock

Clock

// Logger specifies the Logger implementation that will be used for logging warning and errors.

//

// The cache will use slog.Default() by default.

Logger

Logger

}

Options should be passed to

New

/

Must

to construct a

Cache

having a combination of the following features:

automatic loading of entries into the cache

size-based eviction when a maximum is exceeded based on frequency and recency

time-based expiration of entries, measured since last access or last write

asynchronously refresh when the first stale request for an entry occurs

notification of deleted entries

accumulation of cache access statistics

These features are all optional; caches can be created using all or none of them. By default,
cache instances created using

Options

will not perform any type of eviction.

cache := otter.Must(&Options[string, string]{
	MaximumSize:      10_000,
  	ExpiryCalculator: otter.ExpiryWriting[string, string(10 * time.Minute),
	StatsRecorder:    stats.NewCounter(),
})

Entries are automatically evicted from the cache when any of MaximumSize, MaximumWeight,
ExpiryCalculator are specified.

If MaximumSize or MaximumWeight is specified, entries may be evicted on each cache modification.

If ExpiryCalculator is specified, then entries may be evicted on
each cache modification, on occasional cache accesses, or on calls to

Cache.CleanUp

.
Expired entries may be counted by

Cache.EstimatedSize

, but will never be visible to read or write operations.

Certain cache configurations will result in the accrual of periodic maintenance tasks that
will be performed during write operations, or during occasional read operations in the absence of writes.
The

Cache.CleanUp

method of the returned cache will also perform maintenance, but
calling it should not be necessary with a high-throughput cache. Only caches built with
MaximumSize, MaximumWeight, ExpiryCalculator perform periodic maintenance.

type

RefreshCalculator

¶

type RefreshCalculator[K

comparable

, V

any

] interface {

// RefreshAfterCreate returns the duration after which the entry is eligible for an automatic refresh after the

// entry's creation. To indicate no refresh, an entry may be given an excessively long period.

RefreshAfterCreate(entry

Entry

[K, V])

time

.

Duration

// RefreshAfterUpdate returns the duration after which the entry is eligible for an automatic refresh after the

// replacement of the entry's value due to an explicit update.

// The entry.RefreshableAfter() may be returned to not modify the refresh time.

RefreshAfterUpdate(entry

Entry

[K, V], oldValue V)

time

.

Duration

// RefreshAfterReload returns the duration after which the entry is eligible for an automatic refresh after the

// replacement of the entry's value due to a reload.

// The entry.RefreshableAfter() may be returned to not modify the refresh time.

RefreshAfterReload(entry

Entry

[K, V], oldValue V)

time

.

Duration

// RefreshAfterReloadFailure returns the duration after which the entry is eligible for an automatic refresh after the

// value failed to be reloaded.

// The entry.RefreshableAfter() may be returned to not modify the refresh time.

RefreshAfterReloadFailure(entry

Entry

[K, V], err

error

)

time

.

Duration

}

RefreshCalculator calculates when cache entries will be reloaded. A single refresh time is retained so that the lifetime
of an entry may be extended or reduced by subsequent evaluations.

func

RefreshCreating

¶

func RefreshCreating[K

comparable

, V

any

](duration

time

.

Duration

)

RefreshCalculator

[K, V]

RefreshCreating returns a

RefreshCalculator

that specifies that the entry should be automatically reloaded
once the duration has elapsed after the entry's creation.
The refresh time is not modified when the entry is updated or reloaded.

func

RefreshCreatingFunc

¶

func RefreshCreatingFunc[K

comparable

, V

any

](f func(entry

Entry

[K, V])

time

.

Duration

)

RefreshCalculator

[K, V]

RefreshCreatingFunc returns a

RefreshCalculator

that specifies that the entry should be automatically reloaded
once the duration has elapsed after the entry's creation.
The refresh time is not modified when the entry is updated or reloaded.

func

RefreshWriting

¶

func RefreshWriting[K

comparable

, V

any

](duration

time

.

Duration

)

RefreshCalculator

[K, V]

RefreshWriting returns a

RefreshCalculator

that specifies that the entry should be automatically reloaded
once the duration has elapsed after the entry's creation or the most recent replacement of its value.
The refresh time is not modified when the reload fails.

func

RefreshWritingFunc

¶

func RefreshWritingFunc[K

comparable

, V

any

](f func(entry

Entry

[K, V])

time

.

Duration

)

RefreshCalculator

[K, V]

RefreshWritingFunc returns a

RefreshCalculator

that specifies that the entry should be automatically reloaded
once the duration has elapsed after the entry's creation or the most recent replacement of its value.
The refresh time is not modified when the reload fails.

type

RefreshResult

¶

type RefreshResult[K

comparable

, V

any

] struct {

// Key is the key corresponding to the refreshed entry.

Key K

// Value is the value corresponding to the refreshed entry.

Value V

// Err is the error that Loader / BulkLoader returned.

Err

error

}

RefreshResult holds the results of

Cache.Refresh

/

Cache.BulkRefresh

, so they can be passed
on a channel.