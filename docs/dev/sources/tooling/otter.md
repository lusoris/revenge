# otter Cache

> Auto-fetched from [https://pkg.go.dev/github.com/maypok86/otter](https://pkg.go.dev/github.com/maypok86/otter)
> Last Updated: 2026-01-29T20:11:42.883418+00:00

---

Index
¶
Constants
Variables
type Builder
func MustBuilder[K comparable, V any](capacity int) *Builder[K, V]
func NewBuilder[K comparable, V any](capacity int) (*Builder[K, V], error)
func (b *Builder[K, V]) Build() (Cache[K, V], error)
func (b *Builder[K, V]) CollectStats() *Builder[K, V]
func (b *Builder[K, V]) Cost(costFunc func(key K, value V) uint32) *Builder[K, V]
func (b *Builder[K, V]) DeletionListener(deletionListener func(key K, value V, cause DeletionCause)) *Builder[K, V]
func (b *Builder[K, V]) InitialCapacity(initialCapacity int) *Builder[K, V]
func (b *Builder[K, V]) WithTTL(ttl time.Duration) *ConstTTLBuilder[K, V]
func (b *Builder[K, V]) WithVariableTTL() *VariableTTLBuilder[K, V]
type Cache
func (bs Cache) Capacity() int
func (bs Cache) Clear()
func (bs Cache) Close()
func (bs Cache) Delete(key K)
func (bs Cache) DeleteByFunc(f func(key K, value V) bool)
func (bs Cache) Extension() Extension[K, V]
func (bs Cache) Get(key K) (V, bool)
func (bs Cache) Has(key K) bool
func (bs Cache) Range(f func(key K, value V) bool)
func (c Cache[K, V]) Set(key K, value V) bool
func (c Cache[K, V]) SetIfAbsent(key K, value V) bool
func (bs Cache) Size() int
func (bs Cache) Stats() Stats
type CacheWithVariableTTL
func (bs CacheWithVariableTTL) Capacity() int
func (bs CacheWithVariableTTL) Clear()
func (bs CacheWithVariableTTL) Close()
func (bs CacheWithVariableTTL) Delete(key K)
func (bs CacheWithVariableTTL) DeleteByFunc(f func(key K, value V) bool)
func (bs CacheWithVariableTTL) Extension() Extension[K, V]
func (bs CacheWithVariableTTL) Get(key K) (V, bool)
func (bs CacheWithVariableTTL) Has(key K) bool
func (bs CacheWithVariableTTL) Range(f func(key K, value V) bool)
func (c CacheWithVariableTTL[K, V]) Set(key K, value V, ttl time.Duration) bool
func (c CacheWithVariableTTL[K, V]) SetIfAbsent(key K, value V, ttl time.Duration) bool
func (bs CacheWithVariableTTL) Size() int
func (bs CacheWithVariableTTL) Stats() Stats
type ConstTTLBuilder
func (b *ConstTTLBuilder[K, V]) Build() (Cache[K, V], error)
func (b *ConstTTLBuilder[K, V]) CollectStats() *ConstTTLBuilder[K, V]
func (b *ConstTTLBuilder[K, V]) Cost(costFunc func(key K, value V) uint32) *ConstTTLBuilder[K, V]
func (b *ConstTTLBuilder[K, V]) DeletionListener(deletionListener func(key K, value V, cause DeletionCause)) *ConstTTLBuilder[K, V]
func (b *ConstTTLBuilder[K, V]) InitialCapacity(initialCapacity int) *ConstTTLBuilder[K, V]
type DeletionCause
type Entry
func (e Entry[K, V]) Cost() uint32
func (e Entry[K, V]) Expiration() int64
func (e Entry[K, V]) HasExpired() bool
func (e Entry[K, V]) Key() K
func (e Entry[K, V]) TTL() time.Duration
func (e Entry[K, V]) Value() V
type Extension
func (e Extension[K, V]) GetEntry(key K) (Entry[K, V], bool)
func (e Extension[K, V]) GetEntryQuietly(key K) (Entry[K, V], bool)
func (e Extension[K, V]) GetQuietly(key K) (V, bool)
type Stats
func (s Stats) EvictedCost() int64
func (s Stats) EvictedCount() int64
func (s Stats) Hits() int64
func (s Stats) Misses() int64
func (s Stats) Ratio() float64
func (s Stats) RejectedSets() int64
type VariableTTLBuilder
func (b *VariableTTLBuilder[K, V]) Build() (CacheWithVariableTTL[K, V], error)
func (b *VariableTTLBuilder[K, V]) CollectStats() *VariableTTLBuilder[K, V]
func (b *VariableTTLBuilder[K, V]) Cost(costFunc func(key K, value V) uint32) *VariableTTLBuilder[K, V]
func (b *VariableTTLBuilder[K, V]) DeletionListener(deletionListener func(key K, value V, cause DeletionCause)) *VariableTTLBuilder[K, V]
func (b *VariableTTLBuilder[K, V]) InitialCapacity(initialCapacity int) *VariableTTLBuilder[K, V]
Constants
¶
View Source
const (
// Explicit the entry was manually deleted by the user.
Explicit =
core
.
Explicit
// Replaced the entry itself was not actually deleted, but its value was replaced by the user.
Replaced =
core
.
Replaced
// Size the entry was evicted due to size constraints.
Size =
core
.
Size
// Expired the entry's expiration timestamp has passed.
Expired =
core
.
Expired
)
Variables
¶
View Source
var (
// ErrIllegalCapacity means that a non-positive capacity has been passed to the NewBuilder.
ErrIllegalCapacity =
errors
.
New
("capacity should be positive")
// ErrIllegalInitialCapacity means that a non-positive capacity has been passed to the Builder.InitialCapacity.
ErrIllegalInitialCapacity =
errors
.
New
("initial capacity should be positive")
// ErrNilCostFunc means that a nil cost func has been passed to the Builder.Cost.
ErrNilCostFunc =
errors
.
New
("setCostFunc func should not be nil")
// ErrIllegalTTL means that a non-positive ttl has been passed to the Builder.WithTTL.
ErrIllegalTTL =
errors
.
New
("ttl should be positive")
)
Functions
¶
This section is empty.
Types
¶
type
Builder
¶
type Builder[K
comparable
, V
any
] struct {
// contains filtered or unexported fields
}
Builder is a one-shot builder for creating a cache instance.
func
MustBuilder
¶
func MustBuilder[K
comparable
, V
any
](capacity
int
) *
Builder
[K, V]
MustBuilder creates a builder and sets the future cache capacity.
Panics if capacity <= 0.
func
NewBuilder
¶
func NewBuilder[K
comparable
, V
any
](capacity
int
) (*
Builder
[K, V],
error
)
NewBuilder creates a builder and sets the future cache capacity.
Returns an error if capacity <= 0.
func (*Builder[K, V])
Build
¶
func (b *
Builder
[K, V]) Build() (
Cache
[K, V],
error
)
Build creates a configured cache or
returns an error if invalid parameters were passed to the builder.
func (*Builder[K, V])
CollectStats
¶
func (b *
Builder
[K, V]) CollectStats() *
Builder
[K, V]
CollectStats determines whether statistics should be calculated when the cache is running.
By default, statistics calculating is disabled.
func (*Builder[K, V])
Cost
¶
func (b *
Builder
[K, V]) Cost(costFunc func(key K, value V)
uint32
) *
Builder
[K, V]
Cost sets a function to dynamically calculate the cost of an item.
By default, this function always returns 1.
func (*Builder[K, V])
DeletionListener
¶
added in
v1.2.0
func (b *
Builder
[K, V]) DeletionListener(deletionListener func(key K, value V, cause
DeletionCause
)) *
Builder
[K, V]
DeletionListener specifies a listener instance that caches should notify each time an entry is deleted for any
DeletionCause cause. The cache will invoke this listener in the background goroutine
after the entry's deletion operation has completed.
func (*Builder[K, V])
InitialCapacity
¶
added in
v1.1.0
func (b *
Builder
[K, V]) InitialCapacity(initialCapacity
int
) *
Builder
[K, V]
InitialCapacity sets the minimum total size for the internal data structures. Providing a large enough estimate
at construction time avoids the need for expensive resizing operations later, but setting this
value unnecessarily high wastes memory.
func (*Builder[K, V])
WithTTL
¶
func (b *
Builder
[K, V]) WithTTL(ttl
time
.
Duration
) *
ConstTTLBuilder
[K, V]
WithTTL specifies that each item should be automatically removed from the cache once a fixed duration
has elapsed after the item's creation.
func (*Builder[K, V])
WithVariableTTL
¶
func (b *
Builder
[K, V]) WithVariableTTL() *
VariableTTLBuilder
[K, V]
WithVariableTTL specifies that each item should be automatically removed from the cache once a duration has
elapsed after the item's creation. Items are expired based on the custom ttl specified for each item separately.
You should prefer WithTTL to this option whenever possible.
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
Cache is a structure performs a best-effort bounding of a hash table using eviction algorithm
to determine which entries to evict when the capacity is exceeded.
func (Cache)
Capacity
¶
func (bs Cache) Capacity()
int
Capacity returns the cache capacity.
func (Cache)
Clear
¶
func (bs Cache) Clear()
Clear clears the hash table, all policies, buffers, etc.
NOTE: this operation must be performed when no requests are made to the cache otherwise the behavior is undefined.
func (Cache)
Close
¶
func (bs Cache) Close()
Close clears the hash table, all policies, buffers, etc and stop all goroutines.
NOTE: this operation must be performed when no requests are made to the cache otherwise the behavior is undefined.
func (Cache)
Delete
¶
func (bs Cache) Delete(key K)
Delete removes the association for this key from the cache.
func (Cache)
DeleteByFunc
¶
added in
v1.1.0
func (bs Cache) DeleteByFunc(f func(key K, value V)
bool
)
DeleteByFunc removes the association for this key from the cache when the given function returns true.
func (Cache)
Extension
¶
added in
v1.2.0
func (bs Cache) Extension()
Extension
[K, V]
Extension returns access to inspect and perform low-level operations on this cache based on its runtime
characteristics. These operations are optional and dependent on how the cache was constructed
and what abilities the implementation exposes.
func (Cache)
Get
¶
func (bs Cache) Get(key K) (V,
bool
)
Get returns the value associated with the key in this cache.
func (Cache)
Has
¶
func (bs Cache) Has(key K)
bool
Has checks if there is an entry with the given key in the cache.
func (Cache)
Range
¶
func (bs Cache) Range(f func(key K, value V)
bool
)
Range iterates over all entries in the cache.
Iteration stops early when the given function returns false.
func (Cache[K, V])
Set
¶
func (c
Cache
[K, V]) Set(key K, value V)
bool
Set associates the value with the key in this cache.
If it returns false, then the key-value pair had too much cost and the Set was dropped.
func (Cache[K, V])
SetIfAbsent
¶
func (c
Cache
[K, V]) SetIfAbsent(key K, value V)
bool
SetIfAbsent if the specified key is not already associated with a value associates it with the given value.
If the specified key is not already associated with a value, then it returns false.
Also, it returns false if the key-value pair had too much cost and the SetIfAbsent was dropped.
func (Cache)
Size
¶
func (bs Cache) Size()
int
Size returns the current number of entries in the cache.
func (Cache)
Stats
¶
func (bs Cache) Stats()
Stats
Stats returns a current snapshot of this cache's cumulative statistics.
type
CacheWithVariableTTL
¶
type CacheWithVariableTTL[K
comparable
, V
any
] struct {
// contains filtered or unexported fields
}
CacheWithVariableTTL is a structure performs a best-effort bounding of a hash table using eviction algorithm
to determine which entries to evict when the capacity is exceeded.
func (CacheWithVariableTTL)
Capacity
¶
func (bs CacheWithVariableTTL) Capacity()
int
Capacity returns the cache capacity.
func (CacheWithVariableTTL)
Clear
¶
func (bs CacheWithVariableTTL) Clear()
Clear clears the hash table, all policies, buffers, etc.
NOTE: this operation must be performed when no requests are made to the cache otherwise the behavior is undefined.
func (CacheWithVariableTTL)
Close
¶
func (bs CacheWithVariableTTL) Close()
Close clears the hash table, all policies, buffers, etc and stop all goroutines.
NOTE: this operation must be performed when no requests are made to the cache otherwise the behavior is undefined.
func (CacheWithVariableTTL)
Delete
¶
func (bs CacheWithVariableTTL) Delete(key K)
Delete removes the association for this key from the cache.
func (CacheWithVariableTTL)
DeleteByFunc
¶
added in
v1.1.0
func (bs CacheWithVariableTTL) DeleteByFunc(f func(key K, value V)
bool
)
DeleteByFunc removes the association for this key from the cache when the given function returns true.
func (CacheWithVariableTTL)
Extension
¶
added in
v1.2.0
func (bs CacheWithVariableTTL) Extension()
Extension
[K, V]
Extension returns access to inspect and perform low-level operations on this cache based on its runtime
characteristics. These operations are optional and dependent on how the cache was constructed
and what abilities the implementation exposes.
func (CacheWithVariableTTL)
Get
¶
func (bs CacheWithVariableTTL) Get(key K) (V,
bool
)
Get returns the value associated with the key in this cache.
func (CacheWithVariableTTL)
Has
¶
func (bs CacheWithVariableTTL) Has(key K)
bool
Has checks if there is an entry with the given key in the cache.
func (CacheWithVariableTTL)
Range
¶
func (bs CacheWithVariableTTL) Range(f func(key K, value V)
bool
)
Range iterates over all entries in the cache.
Iteration stops early when the given function returns false.
func (CacheWithVariableTTL[K, V])
Set
¶
func (c
CacheWithVariableTTL
[K, V]) Set(key K, value V, ttl
time
.
Duration
)
bool
Set associates the value with the key in this cache and sets the custom ttl for this key-value pair.
If it returns false, then the key-value pair had too much cost and the Set was dropped.
func (CacheWithVariableTTL[K, V])
SetIfAbsent
¶
func (c
CacheWithVariableTTL
[K, V]) SetIfAbsent(key K, value V, ttl
time
.
Duration
)
bool
SetIfAbsent if the specified key is not already associated with a value associates it with the given value
and sets the custom ttl for this key-value pair.
If the specified key is not already associated with a value, then it returns false.
Also, it returns false if the key-value pair had too much cost and the SetIfAbsent was dropped.
func (CacheWithVariableTTL)
Size
¶
func (bs CacheWithVariableTTL) Size()
int
Size returns the current number of entries in the cache.
func (CacheWithVariableTTL)
Stats
¶
func (bs CacheWithVariableTTL) Stats()
Stats
Stats returns a current snapshot of this cache's cumulative statistics.
type
ConstTTLBuilder
¶
type ConstTTLBuilder[K
comparable
, V
any
] struct {
// contains filtered or unexported fields
}
ConstTTLBuilder is a one-shot builder for creating a cache instance.
func (*ConstTTLBuilder[K, V])
Build
¶
func (b *
ConstTTLBuilder
[K, V]) Build() (
Cache
[K, V],
error
)
Build creates a configured cache or
returns an error if invalid parameters were passed to the builder.
func (*ConstTTLBuilder[K, V])
CollectStats
¶
func (b *
ConstTTLBuilder
[K, V]) CollectStats() *
ConstTTLBuilder
[K, V]
CollectStats determines whether statistics should be calculated when the cache is running.
By default, statistics calculating is disabled.
func (*ConstTTLBuilder[K, V])
Cost
¶
func (b *
ConstTTLBuilder
[K, V]) Cost(costFunc func(key K, value V)
uint32
) *
ConstTTLBuilder
[K, V]
Cost sets a function to dynamically calculate the cost of an item.
By default, this function always returns 1.
func (*ConstTTLBuilder[K, V])
DeletionListener
¶
added in
v1.2.0
func (b *
ConstTTLBuilder
[K, V]) DeletionListener(deletionListener func(key K, value V, cause
DeletionCause
)) *
ConstTTLBuilder
[K, V]
DeletionListener specifies a listener instance that caches should notify each time an entry is deleted for any
DeletionCause cause. The cache will invoke this listener in the background goroutine
after the entry's deletion operation has completed.
func (*ConstTTLBuilder[K, V])
InitialCapacity
¶
added in
v1.1.0
func (b *
ConstTTLBuilder
[K, V]) InitialCapacity(initialCapacity
int
) *
ConstTTLBuilder
[K, V]
InitialCapacity sets the minimum total size for the internal data structures. Providing a large enough estimate
at construction time avoids the need for expensive resizing operations later, but setting this
value unnecessarily high wastes memory.
type
DeletionCause
¶
added in
v1.2.0
type DeletionCause =
core
.
DeletionCause
DeletionCause the cause why a cached entry was deleted.
type
Entry
¶
added in
v1.2.0
type Entry[K
comparable
, V
any
] struct {
// contains filtered or unexported fields
}
Entry is a key-value pair that may include policy metadata for the cached entry.
It is an immutable snapshot of the cached data at the time of this entry's creation, and it will not
reflect changes afterward.
func (Entry[K, V])
Cost
¶
added in
v1.2.0
func (e
Entry
[K, V]) Cost()
uint32
Cost returns the entry's cost.
If the cache was not configured with a cost then this value is always 1.
func (Entry[K, V])
Expiration
¶
added in
v1.2.0
func (e
Entry
[K, V]) Expiration()
int64
Expiration returns the entry's expiration time as a unix time,
the number of seconds elapsed since January 1, 1970 UTC.
If the cache was not configured with an expiration policy then this value is always 0.
func (Entry[K, V])
HasExpired
¶
added in
v1.2.0
func (e
Entry
[K, V]) HasExpired()
bool
HasExpired returns true if the entry has expired.
func (Entry[K, V])
Key
¶
added in
v1.2.0
func (e
Entry
[K, V]) Key() K
Key returns the entry's key.
func (Entry[K, V])
TTL
¶
added in
v1.2.0
func (e
Entry
[K, V]) TTL()
time
.
Duration
TTL returns the entry's ttl.
If the cache was not configured with an expiration policy then this value is always -1.
If the entry is expired then this value is always 0.
func (Entry[K, V])
Value
¶
added in
v1.2.0
func (e
Entry
[K, V]) Value() V
Value returns the entry's value.
type
Extension
¶
added in
v1.2.0
type Extension[K
comparable
, V
any
] struct {
// contains filtered or unexported fields
}
Extension is an access point for inspecting and performing low-level operations based on the cache's runtime
characteristics. These operations are optional and dependent on how the cache was constructed
and what abilities the implementation exposes.
func (Extension[K, V])
GetEntry
¶
added in
v1.2.0
func (e
Extension
[K, V]) GetEntry(key K) (
Entry
[K, V],
bool
)
GetEntry returns the cache entry associated with the key in this cache.
func (Extension[K, V])
GetEntryQuietly
¶
added in
v1.2.0
func (e
Extension
[K, V]) GetEntryQuietly(key K) (
Entry
[K, V],
bool
)
GetEntryQuietly returns the cache entry associated with the key in this cache.
Unlike GetEntry, this function does not produce any side effects
such as updating statistics or the eviction policy.
func (Extension[K, V])
GetQuietly
¶
added in
v1.2.0
func (e
Extension
[K, V]) GetQuietly(key K) (V,
bool
)
GetQuietly returns the value associated with the key in this cache.
Unlike Get in the cache, this function does not produce any side effects
such as updating statistics or the eviction policy.
type
Stats
¶
type Stats struct {
// contains filtered or unexported fields
}
Stats is a statistics snapshot.
func (Stats)
EvictedCost
¶
added in
v1.1.0
func (s
Stats
) EvictedCost()
int64
EvictedCost returns the sum of costs of evicted entries.
func (Stats)
EvictedCount
¶
added in
v1.1.0
func (s
Stats
) EvictedCount()
int64
EvictedCount returns the number of evicted entries.
func (Stats)
Hits
¶
func (s
Stats
) Hits()
int64
Hits returns the number of cache hits.
func (Stats)
Misses
¶
func (s
Stats
) Misses()
int64
Misses returns the number of cache misses.
func (Stats)
Ratio
¶
func (s
Stats
) Ratio()
float64
Ratio returns the cache hit ratio.
func (Stats)
RejectedSets
¶
added in
v1.1.0
func (s
Stats
) RejectedSets()
int64
RejectedSets returns the number of rejected sets.
type
VariableTTLBuilder
¶
type VariableTTLBuilder[K
comparable
, V
any
] struct {
// contains filtered or unexported fields
}
VariableTTLBuilder is a one-shot builder for creating a cache instance.
func (*VariableTTLBuilder[K, V])
Build
¶
func (b *
VariableTTLBuilder
[K, V]) Build() (
CacheWithVariableTTL
[K, V],
error
)
Build creates a configured cache or
returns an error if invalid parameters were passed to the builder.
func (*VariableTTLBuilder[K, V])
CollectStats
¶
func (b *
VariableTTLBuilder
[K, V]) CollectStats() *
VariableTTLBuilder
[K, V]
CollectStats determines whether statistics should be calculated when the cache is running.
By default, statistics calculating is disabled.
func (*VariableTTLBuilder[K, V])
Cost
¶
func (b *
VariableTTLBuilder
[K, V]) Cost(costFunc func(key K, value V)
uint32
) *
VariableTTLBuilder
[K, V]
Cost sets a function to dynamically calculate the cost of an item.
By default, this function always returns 1.
func (*VariableTTLBuilder[K, V])
DeletionListener
¶
added in
v1.2.0
func (b *
VariableTTLBuilder
[K, V]) DeletionListener(deletionListener func(key K, value V, cause
DeletionCause
)) *
VariableTTLBuilder
[K, V]
DeletionListener specifies a listener instance that caches should notify each time an entry is deleted for any
DeletionCause cause. The cache will invoke this listener in the background goroutine
after the entry's deletion operation has completed.
func (*VariableTTLBuilder[K, V])
InitialCapacity
¶
added in
v1.1.0
func (b *
VariableTTLBuilder
[K, V]) InitialCapacity(initialCapacity
int
) *
VariableTTLBuilder
[K, V]
InitialCapacity sets the minimum total size for the internal data structures. Providing a large enough estimate
at construction time avoids the need for expensive resizing operations later, but setting this
value unnecessarily high wastes memory.