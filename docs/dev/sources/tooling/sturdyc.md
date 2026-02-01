# sturdyc

> Source: https://pkg.go.dev/github.com/viccon/sturdyc
> Fetched: 2026-02-01T11:42:11.390079+00:00
> Content-Hash: bae021928bc4988a
> Type: html

---

### Index ¶

  * Variables
  * func FindCutoff(times []time.Time, percentile float64) time.Time
  * func GetOrFetch[V, T any](ctx context.Context, c *Client[T], key string, fetchFn FetchFn[V]) (V, error)
  * func GetOrFetchBatch[V, T any](ctx context.Context, c *Client[T], ids []string, keyFn KeyFn, ...) (map[string]V, error)
  * func Passthrough[T, V any](ctx context.Context, c *Client[T], key string, fetchFn FetchFn[V]) (V, error)
  * func PassthroughBatch[V, T any](ctx context.Context, c *Client[T], ids []string, keyFn KeyFn, ...) (map[string]V, error)
  * type BatchFetchFn
  * type BatchResponse
  * type Client
  *     * func New[T any](capacity, numShards int, ttl time.Duration, evictionPercentage int, ...) *Client[T]
  *     * func (c *Client[T]) BatchKeyFn(prefix string) KeyFn
    * func (c *Client[T]) Delete(key string)
    * func (c *Client[T]) Get(key string) (T, bool)
    * func (c *Client[T]) GetMany(keys []string) map[string]T
    * func (c *Client[T]) GetManyKeyFn(ids []string, keyFn KeyFn) map[string]T
    * func (c *Client[T]) GetOrFetch(ctx context.Context, key string, fetchFn FetchFn[T]) (T, error)
    * func (c *Client[T]) GetOrFetchBatch(ctx context.Context, ids []string, keyFn KeyFn, fetchFn BatchFetchFn[T]) (map[string]T, error)
    * func (c *Client[T]) NumKeysInflight() int
    * func (c *Client[T]) Passthrough(ctx context.Context, key string, fetchFn FetchFn[T]) (T, error)
    * func (c *Client[T]) PassthroughBatch(ctx context.Context, ids []string, keyFn KeyFn, fetchFn BatchFetchFn[T]) (map[string]T, error)
    * func (c *Client[T]) PermutatedBatchKeyFn(prefix string, permutationStruct interface{}) KeyFn
    * func (c *Client[T]) PermutatedKey(prefix string, permutationStruct interface{}) string
    * func (c *Client[T]) ScanKeys() []string
    * func (c *Client[T]) Set(key string, value T) bool
    * func (c *Client[T]) SetMany(records map[string]T) bool
    * func (c *Client[T]) SetManyKeyFn(records map[string]T, cacheKeyFn KeyFn) bool
    * func (c *Client[T]) Size() int
    * func (c *Client[T]) StoreMissingRecord(key string) bool
  * type Clock
  * type Config
  * type DistributedMetricsRecorder
  * type DistributedStorage
  * type DistributedStorageWithDeletions
  * type FetchFn
  * type KeyFn
  * type Logger
  * type MetricsRecorder
  * type NoopLogger
  *     * func (l *NoopLogger) Error(_ string, _ ...any)
    * func (l *NoopLogger) Warn(_ string, _ ...any)
  * type Option
  *     * func WithClock(clock Clock) Option
    * func WithDistributedMetrics(metricsRecorder DistributedMetricsRecorder) Option
    * func WithDistributedStorage(storage DistributedStorage) Option
    * func WithDistributedStorageEarlyRefreshes(storage DistributedStorageWithDeletions, refreshAfter time.Duration) Option
    * func WithEarlyRefreshes(...) Option
    * func WithEvictionInterval(interval time.Duration) Option
    * func WithLog(log Logger) Option
    * func WithMetrics(recorder MetricsRecorder) Option
    * func WithMissingRecordStorage() Option
    * func WithNoContinuousEvictions() Option
    * func WithRefreshCoalescing(bufferSize int, bufferDuration time.Duration) Option
    * func WithRelativeTimeKeyFormat(truncation time.Duration) Option
  * type RealClock
  *     * func NewClock() *RealClock
  *     * func (c *RealClock) NewTicker(d time.Duration) (<-chan time.Time, func())
    * func (c *RealClock) NewTimer(d time.Duration) (<-chan time.Time, func() bool)
    * func (c *RealClock) Now() time.Time
    * func (c *RealClock) Since(t time.Time) time.Duration
  * type TestClock
  *     * func NewTestClock(time time.Time) *TestClock
  *     * func (c *TestClock) Add(d time.Duration)
    * func (c *TestClock) NewTicker(d time.Duration) (<-chan time.Time, func())
    * func (c *TestClock) NewTimer(d time.Duration) (<-chan time.Time, func() bool)
    * func (c *TestClock) Now() time.Time
    * func (c *TestClock) Set(t time.Time)
    * func (c *TestClock) Since(t time.Time) time.Duration



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/viccon/sturdyc/blob/v1.1.5/errors.go#L5)
    
    
    var (
    
    	// ErrNotFound should be returned from a FetchFn to indicate that a record is
    	// missing at the underlying data source. This helps the cache to determine
    	// if a record should be deleted or stored as a missing record if you have
    	// that functionality enabled. Missing records are refreshed like any other
    	// record, and if your FetchFn returns a value for it, the record will no
    	// longer be considered missing. Please note that this only applies to
    	// client.GetOrFetch and client.Passthrough. For client.GetOrFetchBatch and
    	// client.PassthroughBatch, this works implicitly if you return
    	// a map without the ID, and have store missing records enabled.
    	ErrNotFound = [errors](/errors).[New](/errors#New)("sturdyc: err not found")
    	// ErrMissingRecord is returned by client.GetOrFetch and client.Passthrough when a record has been marked
    	// as missing. The cache will still try to refresh the record in the background if it's being requested.
    	ErrMissingRecord = [errors](/errors).[New](/errors#New)("sturdyc: the record has been marked as missing in the cache")
    	// ErrOnlyCachedRecords can be returned when you're using the cache with
    	// early refreshes or distributed storage functionality. It indicates that
    	// the records *should* have been refreshed from the underlying data source,
    	// but the operation failed. It is up to you to decide whether you want to
    	// proceed with the records that were retrieved from the cache. Note: For
    	// batch operations, this might contain only part of the batch. For example,
    	// if you requested keys 1-10, and we had IDs 1-3 in the cache, but the
    	// request to fetch records 4-10 failed.
    	ErrOnlyCachedRecords = [errors](/errors).[New](/errors#New)("sturdyc: failed to fetch the records that were not in the cache")
    	// ErrInvalidType is returned when you try to use one of the generic
    	// package level functions but the type assertion fails.
    	ErrInvalidType = [errors](/errors).[New](/errors#New)("sturdyc: invalid response type")
    )

### Functions ¶

####  func [FindCutoff](https://github.com/viccon/sturdyc/blob/v1.1.5/quickselect.go#L34) ¶
    
    
    func FindCutoff(times [][time](/time).[Time](/time#Time), percentile [float64](/builtin#float64)) [time](/time).[Time](/time#Time)

FindCutoff returns the time that is the k-th smallest time in the slice. 

####  func [GetOrFetch](https://github.com/viccon/sturdyc/blob/v1.1.5/fetch.go#L124) ¶
    
    
    func GetOrFetch[V, T [any](/builtin#any)](ctx [context](/context).[Context](/context#Context), c *Client[T], key [string](/builtin#string), fetchFn FetchFn[V]) (V, [error](/builtin#error))

GetOrFetch is a convenience function that performs type assertion on the result of client.GetOrFetch. 

Parameters: 
    
    
    ctx - The context to be used for the request.
    c - The cache client.
    key - The key to be fetched.
    fetchFn - Used to retrieve the data from the underlying data source if the key is not found in the cache.
    

Returns: 
    
    
    The value corresponding to the key and an error if one occurred.
    

Type Parameters: 
    
    
    V - The type returned by the fetchFn. Must be assignable to T.
    T - The type stored in the cache.
    

####  func [GetOrFetchBatch](https://github.com/viccon/sturdyc/blob/v1.1.5/fetch.go#L232) ¶
    
    
    func GetOrFetchBatch[V, T [any](/builtin#any)](ctx [context](/context).[Context](/context#Context), c *Client[T], ids [][string](/builtin#string), keyFn KeyFn, fetchFn BatchFetchFn[V]) (map[[string](/builtin#string)]V, [error](/builtin#error))

####  func [Passthrough](https://github.com/viccon/sturdyc/blob/v1.1.5/passthrough.go#L50) ¶
    
    
    func Passthrough[T, V [any](/builtin#any)](ctx [context](/context).[Context](/context#Context), c *Client[T], key [string](/builtin#string), fetchFn FetchFn[V]) (V, [error](/builtin#error))

Passthrough is a convenience function that performs type assertion on the result of client.PassthroughBatch. 

Parameters: 
    
    
    ctx - The context to be used for the request.
    c - The cache client.
    key - The key to be fetched.
    fetchFn - Used to retrieve the data from the underlying data source.
    

Returns: 
    
    
    The value and an error if one occurred and the key was not found in the cache.
    

Type Parameters: 
    
    
    V - The type returned by the fetchFn. Must be assignable to T.
    T - The type stored in the cache.
    

####  func [PassthroughBatch](https://github.com/viccon/sturdyc/blob/v1.1.5/passthrough.go#L102) ¶
    
    
    func PassthroughBatch[V, T [any](/builtin#any)](ctx [context](/context).[Context](/context#Context), c *Client[T], ids [][string](/builtin#string), keyFn KeyFn, fetchFn BatchFetchFn[V]) (map[[string](/builtin#string)]V, [error](/builtin#error))

PassthroughBatch is a convenience function that performs type assertion on the result of client.PassthroughBatch. 

Parameters: 
    
    
    ctx - The context to be used for the request.
    c - The cache client.
    ids - The list of IDs to be fetched.
    keyFn - Used to prefix each ID in order to create a unique cache key.
    fetchFn - Used to retrieve the data from the underlying data source.
    

Returns: 
    
    
    A map of ids to their corresponding values and an error if one occurred.
    

Type Parameters: 
    
    
    V - The type returned by the fetchFn. Must be assignable to T.
    T - The type stored in the cache.
    

### Types ¶

####  type [BatchFetchFn](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L16) ¶
    
    
    type BatchFetchFn[T [any](/builtin#any)] func(ctx [context](/context).[Context](/context#Context), ids [][string](/builtin#string)) (map[[string](/builtin#string)]T, [error](/builtin#error))

BatchFetchFn represents a function that can be used to fetch multiple records from a data source. 

####  type [BatchResponse](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L18) ¶
    
    
    type BatchResponse[T [any](/builtin#any)] map[[string](/builtin#string)]T

####  type [Client](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L55) ¶
    
    
    type Client[T [any](/builtin#any)] struct {
    	*Config
    	// contains filtered or unexported fields
    }

Client represents a cache client that can be used to store and retrieve values. 

####  func [New](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L72) ¶
    
    
    func New[T [any](/builtin#any)](capacity, numShards [int](/builtin#int), ttl [time](/time).[Duration](/time#Duration), evictionPercentage [int](/builtin#int), opts ...Option) *Client[T]

New creates a new Client instance with the specified configuration. 
    
    
    `capacity` defines the maximum number of entries that the cache can store.
    `numShards` Is used to set the number of shards. Has to be greater than 0.
    `ttl` Sets the time to live for each entry in the cache. Has to be greater than 0.
    `evictionPercentage` Percentage of items to evict when the cache exceeds its capacity.
    `opts` allows for additional configurations to be applied to the cache client.
    

####  func (*Client[T]) [BatchKeyFn](https://github.com/viccon/sturdyc/blob/v1.1.5/keys.go#L186) ¶
    
    
    func (c *Client[T]) BatchKeyFn(prefix [string](/builtin#string)) KeyFn

BatchKeyFn provides a function that can be used in conjunction with "GetOrFetchBatch". It takes in a prefix and returns a function that will use the prefix, add a -ID- separator, and then append the ID as a suffix for each item. 

Parameters: 
    
    
    prefix - The prefix to be used for each cache key.
    

Returns: 
    
    
    A function that takes an ID and returns a cache key string with the given prefix and ID.
    

Example usage: 
    
    
    fn := c.BatchKeyFn("some-prefix")
    key := fn("1234") // some-prefix-ID-1234
    

####  func (*Client[T]) [Delete](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L294) ¶
    
    
    func (c *Client[T]) Delete(key [string](/builtin#string))

Delete removes a single entry from the cache. 

Parameters: 
    
    
    key: The key of the entry to be removed.
    

####  func (*Client[T]) [Get](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L148) ¶
    
    
    func (c *Client[T]) Get(key [string](/builtin#string)) (T, [bool](/builtin#bool))

Get retrieves a single value from the cache. 

Parameters: 
    
    
    key - The key to be retrieved.
    

Returns: 
    
    
    The value corresponding to the key and a boolean indicating if the value was found.
    

####  func (*Client[T]) [GetMany](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L164) ¶
    
    
    func (c *Client[T]) GetMany(keys [][string](/builtin#string)) map[[string](/builtin#string)]T

GetMany retrieves multiple values from the cache. 

Parameters: 
    
    
    keys - The list of keys to be retrieved.
    

Returns: 
    
    
    A map of keys to their corresponding values.
    

####  func (*Client[T]) [GetManyKeyFn](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L188) ¶
    
    
    func (c *Client[T]) GetManyKeyFn(ids [][string](/builtin#string), keyFn KeyFn) map[[string](/builtin#string)]T

GetManyKeyFn follows the same API as GetOrFetchBatch and PassthroughBatch. You provide it with a slice of IDs and a keyFn, which is applied to create the cache key. The returned map uses the IDs as keys instead of the cache key. If you've used ScanKeys to retrieve the actual keys, you can retrieve the records using GetMany instead. 

Parameters: 
    
    
    ids - The list of IDs to be retrieved.
    keyFn - A function that generates the cache key for each ID.
    

Returns: 
    
    
    A map of IDs to their corresponding values.
    

####  func (*Client[T]) [GetOrFetch](https://github.com/viccon/sturdyc/blob/v1.1.5/fetch.go#L103) ¶
    
    
    func (c *Client[T]) GetOrFetch(ctx [context](/context).[Context](/context#Context), key [string](/builtin#string), fetchFn FetchFn[T]) (T, [error](/builtin#error))

GetOrFetch attempts to retrieve the specified key from the cache. If the value is absent, it invokes the fetchFn function to obtain it and then stores the result. Additionally, when early refreshes are enabled, GetOrFetch determines if the record needs refreshing and, if necessary, schedules this task for background execution. 

Parameters: 
    
    
    ctx - The context to be used for the request.
    key - The key to be fetched.
    fetchFn - Used to retrieve the data from the underlying data source if the key is not found in the cache.
    

Returns: 
    
    
    The value corresponding to the key and an error if one occurred.
    

####  func (*Client[T]) [GetOrFetchBatch](https://github.com/viccon/sturdyc/blob/v1.1.5/fetch.go#L208) ¶
    
    
    func (c *Client[T]) GetOrFetchBatch(ctx [context](/context).[Context](/context#Context), ids [][string](/builtin#string), keyFn KeyFn, fetchFn BatchFetchFn[T]) (map[[string](/builtin#string)]T, [error](/builtin#error))

GetOrFetchBatch attempts to retrieve the specified ids from the cache. If any of the values are absent, it invokes the fetchFn function to obtain them and then stores the result. Additionally, when background refreshes are enabled, GetOrFetch determines if any of the records need refreshing and, if necessary, schedules this to be performed in the background. 

Parameters: 
    
    
    ctx - The context to be used for the request.
    ids - The list of IDs to be fetched.
    keyFn - Used to generate the cache key for each ID.
    fetchFn - Used to retrieve the data from the underlying data source if any IDs are not found in the cache.
    

Returns: 
    
    
    A map of IDs to their corresponding values and an error if one occurred.
    

####  func (*Client[T]) [NumKeysInflight](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L304) ¶
    
    
    func (c *Client[T]) NumKeysInflight() [int](/builtin#int)

NumKeysInflight returns the number of keys that are currently being fetched. 

Returns: 
    
    
    An integer representing the total number of keys that are currently being fetched.
    

####  func (*Client[T]) [Passthrough](https://github.com/viccon/sturdyc/blob/v1.1.5/passthrough.go#L19) ¶
    
    
    func (c *Client[T]) Passthrough(ctx [context](/context).[Context](/context#Context), key [string](/builtin#string), fetchFn FetchFn[T]) (T, [error](/builtin#error))

Passthrough attempts to retrieve the latest data by calling the provided fetchFn. If fetchFn encounters an error, the cache is used as a fallback. 

Parameters: 
    
    
    ctx - The context to be used for the request.
    key - The key to be fetched.
    fetchFn - Used to retrieve the data from the underlying data source.
    

Returns: 
    
    
    The value and an error if one occurred and the key was not found in the cache.
    

####  func (*Client[T]) [PassthroughBatch](https://github.com/viccon/sturdyc/blob/v1.1.5/passthrough.go#L69) ¶
    
    
    func (c *Client[T]) PassthroughBatch(ctx [context](/context).[Context](/context#Context), ids [][string](/builtin#string), keyFn KeyFn, fetchFn BatchFetchFn[T]) (map[[string](/builtin#string)]T, [error](/builtin#error))

PassthroughBatch attempts to retrieve the latest data by calling the provided fetchFn. If fetchFn encounters an error, the cache is used as a fallback. 

Parameters: 
    
    
    ctx - The context to be used for the request.
    ids - The list of IDs to be fetched.
    keyFn - Used to prefix each ID in order to create a unique cache key.
    fetchFn - Used to retrieve the data from the underlying data source.
    

Returns: 
    
    
    A map of IDs to their corresponding values, and an error if one occurred and
    none of the IDs were found in the cache.
    

####  func (*Client[T]) [PermutatedBatchKeyFn](https://github.com/viccon/sturdyc/blob/v1.1.5/keys.go#L218) ¶
    
    
    func (c *Client[T]) PermutatedBatchKeyFn(prefix [string](/builtin#string), permutationStruct interface{}) KeyFn

PermutatedBatchKeyFn provides a function that can be used in conjunction with GetOrFetchBatch. It takes a prefix and a struct where the fields are concatenated with the ID in order to make a unique cache key. Passing anything but a struct for "permutationStruct" will result in a panic. The cache will only use the EXPORTED fields of the struct to construct the key. The permutation struct should be FLAT, with no nested structs. The fields can be any of the basic types, as well as slices and time.Time values. 

Parameters: 
    
    
    prefix - The prefix for the cache key.
    permutationStruct - A struct whose fields are concatenated to form a unique cache key. Only exported fields are used.
    

Returns: 
    
    
    A function that takes an ID and returns a cache key string with the given prefix, permutation struct fields, and ID.
    

Example usage: 
    
    
    type queryParams struct {
    	City               string
    	Country            string
    }
    params := queryParams{"Stockholm", "Sweden"}
    cacheKeyFunc := c.PermutatedBatchKeyFn("prefix", params)
    key := cacheKeyFunc("1") // prefix-Stockholm-Sweden-ID-1
    

####  func (*Client[T]) [PermutatedKey](https://github.com/viccon/sturdyc/blob/v1.1.5/keys.go#L104) ¶
    
    
    func (c *Client[T]) PermutatedKey(prefix [string](/builtin#string), permutationStruct interface{}) [string](/builtin#string)

PermutatedKey takes a prefix and a struct where the fields are concatenated in order to create a unique cache key. Passing anything but a struct for "permutationStruct" will result in a panic. The cache will only use the EXPORTED fields of the struct to construct the key. The permutation struct should be FLAT, with no nested structs. The fields can be any of the basic types, as well as slices and time.Time values. 

Parameters: 
    
    
    prefix - The prefix for the cache key.
    permutationStruct - A struct whose fields are concatenated to form a unique cache key.
    Only exported fields are used.
    

Returns: 
    
    
    A string to be used as the cache key.
    

Example usage: 
    
    
    type queryParams struct {
    	City               string
    	Country            string
    }
    params := queryParams{"Stockholm", "Sweden"}
    key := c.PermutatedKey("prefix",, params) // prefix-Stockholm-Sweden-1
    

####  func (*Client[T]) [ScanKeys](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L268) ¶
    
    
    func (c *Client[T]) ScanKeys() [][string](/builtin#string)

ScanKeys returns a list of all keys in the cache. 

Returns: 
    
    
    A slice of strings representing all the keys in the cache.
    

####  func (*Client[T]) [Set](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L208) ¶
    
    
    func (c *Client[T]) Set(key [string](/builtin#string), value T) [bool](/builtin#bool)

Set writes a single value to the cache. 

Parameters: 
    
    
    key - The key to be set.
    value - The value to be associated with the key.
    

Returns: 
    
    
    A boolean indicating if the set operation triggered an eviction.
    

####  func (*Client[T]) [SetMany](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L229) ¶
    
    
    func (c *Client[T]) SetMany(records map[[string](/builtin#string)]T) [bool](/builtin#bool)

SetMany writes a map of key-value pairs to the cache. 

Parameters: 
    
    
    records - A map of keys to values to be set in the cache.
    

Returns: 
    
    
    A boolean indicating if any of the set operations triggered an eviction.
    

####  func (*Client[T]) [SetManyKeyFn](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L252) ¶
    
    
    func (c *Client[T]) SetManyKeyFn(records map[[string](/builtin#string)]T, cacheKeyFn KeyFn) [bool](/builtin#bool)

SetManyKeyFn follows the same API as GetOrFetchBatch and PassthroughBatch. It takes a map of records where the keyFn is applied to each key in the map before it's stored in the cache. 

Parameters: 
    
    
    records - A map of IDs to values to be set in the cache.
    cacheKeyFn - A function that generates the cache key for each ID.
    

Returns: 
    
    
    A boolean indicating if any of the set operations triggered an eviction.
    

####  func (*Client[T]) [Size](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L281) ¶
    
    
    func (c *Client[T]) Size() [int](/builtin#int)

Size returns the number of entries in the cache. 

Returns: 
    
    
    An integer representing the total number of entries in the cache.
    

####  func (*Client[T]) [StoreMissingRecord](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L214) ¶
    
    
    func (c *Client[T]) StoreMissingRecord(key [string](/builtin#string)) [bool](/builtin#bool)

StoreMissingRecord writes a single value to the cache. Returns true if it triggered an eviction. 

####  type [Clock](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L10) ¶
    
    
    type Clock interface {
    	Now() [time](/time).[Time](/time#Time)
    	NewTicker(d [time](/time).[Duration](/time#Duration)) (<-chan [time](/time).[Time](/time#Time), func())
    	NewTimer(d [time](/time).[Duration](/time#Duration)) (<-chan [time](/time).[Time](/time#Time), func() [bool](/builtin#bool))
    	Since(t [time](/time).[Time](/time#Time)) [time](/time).[Duration](/time#Duration)
    }

Clock is an abstraction for time.Time package that allows for testing. 

####  type [Config](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L25) ¶
    
    
    type Config struct {
    	// contains filtered or unexported fields
    }

Config represents the configuration that can be applied to the cache. 

####  type [DistributedMetricsRecorder](https://github.com/viccon/sturdyc/blob/v1.1.5/metrics.go#L28) ¶
    
    
    type DistributedMetricsRecorder interface {
    	MetricsRecorder
    	// DistributedCacheHit is called for every key that results in a cache hit.
    	DistributedCacheHit()
    	// DistributedCacheMiss is called for every key that results in a cache miss.
    	DistributedCacheMiss()
    	// DistributedRefresh is called when we retrieve a record from
    	// the distributed storage that should be refreshed.
    	DistributedRefresh()
    	// DistributedMissingRecord is called when we retrieve a record from the
    	// distributed storage that has been marked as a missing record.
    	DistributedMissingRecord()
    	// DistributedFallback is called when you are using a distributed storage
    	// with early refreshes, and the call for a value was supposed to refresh it,
    	// but the call failed. When that happens, the cache fallbacks to the latest
    	// value from the distributed storage.
    	DistributedFallback()
    }

####  type [DistributedStorage](https://github.com/viccon/sturdyc/blob/v1.1.5/distribution.go#L22) ¶
    
    
    type DistributedStorage interface {
    	Get(ctx [context](/context).[Context](/context#Context), key [string](/builtin#string)) ([][byte](/builtin#byte), [bool](/builtin#bool))
    	Set(ctx [context](/context).[Context](/context#Context), key [string](/builtin#string), value [][byte](/builtin#byte))
    	GetBatch(ctx [context](/context).[Context](/context#Context), keys [][string](/builtin#string)) map[[string](/builtin#string)][][byte](/builtin#byte)
    	SetBatch(ctx [context](/context).[Context](/context#Context), records map[[string](/builtin#string)][][byte](/builtin#byte))
    }

DistributedStorage is an abstraction that the cache interacts with in order to keep the distributed storage and in-memory cache in sync. Please note that you are responsible for setting the TTL and eviction policy of this storage. 

####  type [DistributedStorageWithDeletions](https://github.com/viccon/sturdyc/blob/v1.1.5/distribution.go#L35) ¶
    
    
    type DistributedStorageWithDeletions interface {
    	DistributedStorage
    	Delete(ctx [context](/context).[Context](/context#Context), key [string](/builtin#string))
    	DeleteBatch(ctx [context](/context).[Context](/context#Context), keys [][string](/builtin#string))
    }

DistributedStorageWithDeletions is an abstraction that the cache interacts with when you want to use a distributed storage with early refreshes. Please note that you are responsible for setting the TTL and eviction policy of this storage. The cache will only call the delete functions when it performs a refresh and notices that the record has been deleted at the underlying data source. 

####  type [FetchFn](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L13) ¶
    
    
    type FetchFn[T [any](/builtin#any)] func(ctx [context](/context).[Context](/context#Context)) (T, [error](/builtin#error))

FetchFn Fetch represents a function that can be used to fetch a single record from a data source. 

####  type [KeyFn](https://github.com/viccon/sturdyc/blob/v1.1.5/cache.go#L22) ¶
    
    
    type KeyFn func(id [string](/builtin#string)) [string](/builtin#string)

KeyFn is called invoked for each record that a batch fetch operation returns. It is used to create unique cache keys. 

####  type [Logger](https://github.com/viccon/sturdyc/blob/v1.1.5/log.go#L3) ¶
    
    
    type Logger interface {
    	Warn(msg [string](/builtin#string), args ...[any](/builtin#any))
    	Error(msg [string](/builtin#string), args ...[any](/builtin#any))
    }

####  type [MetricsRecorder](https://github.com/viccon/sturdyc/blob/v1.1.5/metrics.go#L3) ¶
    
    
    type MetricsRecorder interface {
    	// CacheHit is called for every key that results in a cache hit.
    	CacheHit()
    	// CacheMiss is called for every key that results in a cache miss.
    	CacheMiss()
    	// AsynchronousRefresh is called when a get operation results in an asynchronous refresh.
    	AsynchronousRefresh()
    	// SynchronousRefresh is called when a get operation results in a synchronous refresh.
    	SynchronousRefresh()
    	// MissingRecord is called every time the cache is asked to
    	// look up a key which has been marked as missing.
    	MissingRecord()
    	// ForcedEviction is called when the cache reaches its capacity, and has to
    	// evict keys in order to write a new one.
    	ForcedEviction()
    	// EntriesEvicted is called when the cache evicts keys from a shard.
    	EntriesEvicted([int](/builtin#int))
    	// ShardIndex is called to report which shard it was that performed an operation.
    	ShardIndex([int](/builtin#int))
    	// CacheBatchRefreshSize is called to report the size of the batch refresh.
    	CacheBatchRefreshSize(size [int](/builtin#int))
    	// ObserveCacheSize is called to report the size of the cache.
    	ObserveCacheSize(callback func() [int](/builtin#int))
    }

####  type [NoopLogger](https://github.com/viccon/sturdyc/blob/v1.1.5/log.go#L8) ¶
    
    
    type NoopLogger struct{}

####  func (*NoopLogger) [Error](https://github.com/viccon/sturdyc/blob/v1.1.5/log.go#L11) ¶
    
    
    func (l *NoopLogger) Error(_ [string](/builtin#string), _ ...[any](/builtin#any))

####  func (*NoopLogger) [Warn](https://github.com/viccon/sturdyc/blob/v1.1.5/log.go#L10) ¶
    
    
    func (l *NoopLogger) Warn(_ [string](/builtin#string), _ ...[any](/builtin#any))

####  type [Option](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L5) ¶
    
    
    type Option func(*Config)

####  func [WithClock](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L16) ¶
    
    
    func WithClock(clock Clock) Option

WithClock can be used to change the clock that the cache uses. This is useful for testing. 

####  func [WithDistributedMetrics](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L144) ¶
    
    
    func WithDistributedMetrics(metricsRecorder DistributedMetricsRecorder) Option

WithDistributedMetrics instructs the cache to report additional metrics regarding its interaction with the distributed storage. 

####  func [WithDistributedStorage](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L113) ¶
    
    
    func WithDistributedStorage(storage DistributedStorage) Option

WithDistributedStorage allows you to use the cache with a distributed key-value store. The "GetOrFetch" and "GetOrFetchBatch" functions will check this store first and only proceed to the underlying data source if the key is missing. When a record is retrieved from the underlying data source, it is written both to memory and to the distributed storage. You are responsible for setting TTL and eviction policies for the distributed storage. Sturdyc will only read and write records. 

####  func [WithDistributedStorageEarlyRefreshes](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L134) ¶
    
    
    func WithDistributedStorageEarlyRefreshes(storage DistributedStorageWithDeletions, refreshAfter [time](/time).[Duration](/time#Duration)) Option

WithDistributedStorageEarlyRefreshes is the distributed equivalent of the "WithEarlyRefreshes" option. It allows distributed records to be refreshed before their TTL expires. If a refresh fails, the cache will fall back to what was returned by the distributed storage. This ensures that data can be served for the duration of the TTL even if an upstream system goes down. To use this functionality, you need to implement an interface with two additional methods for deleting records compared to the simpler "WithDistributedStorage" option. This is because a distributed cache that is used with this option might have low refresh durations but high TTLs. If a record is deleted from the underlying data source, it needs to be propagated to the distributed storage before the TTL expires. However, please note that you are still responsible for managing the TTL and eviction policies for the distributed storage. Sturdyc will only delete records that have been removed at the underlying data source. 

####  func [WithEarlyRefreshes](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L63) ¶
    
    
    func WithEarlyRefreshes(minAsyncRefreshTime, maxAsyncRefreshTime, syncRefreshTime, retryBaseDelay [time](/time).[Duration](/time#Duration)) Option

WithEarlyRefreshes instructs the cache to refresh the keys that are in active rotation, thereby preventing them from ever expiring. This can have a significant impact on your application's latency as you're able to continuously serve frequently used keys from memory. An asynchronous background refresh gets scheduled when a key is requested again after a random time between minRefreshTime and maxRefreshTime has passed. This is an important distinction because it means that the cache won't just naively refresh every key it's ever seen. The third argument to this function will also allow you to provide a duration for when a refresh should become synchronous. If any of the refreshes were to fail, you'll get the latest data from the cache for the duration of the TTL. 

####  func [WithEvictionInterval](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L26) ¶
    
    
    func WithEvictionInterval(interval [time](/time).[Duration](/time#Duration)) Option

WithEvictionInterval sets the interval at which the cache scans a shard to evict expired entries. Setting this to a higher value will increase cache performance and is advised if you don't think you'll exceed the capacity. If the capacity is reached, the cache will still trigger an eviction. 

####  func [WithLog](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L100) ¶
    
    
    func WithLog(log Logger) Option

WithLog allows you to set a custom logger for the cache. The cache isn't chatty, and will only log warnings and errors that would be a nightmare to debug. If you absolutely don't want any logs, you can pass in the sturdyc.NoopLogger. 

####  func [WithMetrics](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L8) ¶
    
    
    func WithMetrics(recorder MetricsRecorder) Option

WithMetrics is used to make the cache report metrics. 

####  func [WithMissingRecordStorage](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L46) ¶
    
    
    func WithMissingRecordStorage() Option

WithMissingRecordStorage allows the cache to mark keys as missing from the underlying data source. This allows you to stop streams of outgoing requests for requests that don't exist. The keys will still have the same TTL and refresh durations as any of the other record in the cache. 

####  func [WithNoContinuousEvictions](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L36) ¶
    
    
    func WithNoContinuousEvictions() Option

WithNoContinuousEvictions improves cache performance when the cache capacity is unlikely to be exceeded. While this setting disables the continuous eviction job, it still allows for the eviction of the least recently used items once the cache reaches its full capacity. 

####  func [WithRefreshCoalescing](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L79) ¶
    
    
    func WithRefreshCoalescing(bufferSize [int](/builtin#int), bufferDuration [time](/time).[Duration](/time#Duration)) Option

WithRefreshCoalescing will make the cache refresh data from batchable endpoints more efficiently. It is going to create a buffer for each cache key permutation, and gather IDs until the bufferSize is reached, or the bufferDuration has passed. 

NOTE: This requires the WithEarlyRefreshes functionality to be enabled. 

####  func [WithRelativeTimeKeyFormat](https://github.com/viccon/sturdyc/blob/v1.1.5/options.go#L90) ¶
    
    
    func WithRelativeTimeKeyFormat(truncation [time](/time).[Duration](/time#Duration)) Option

WithRelativeTimeKeyFormat allows you to control the truncation of time.Time values that are being passed in to the cache key functions. 

####  type [RealClock](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L18) ¶
    
    
    type RealClock struct{}

RealClock provides functions that wraps the real time.Time package. 

####  func [NewClock](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L21) ¶
    
    
    func NewClock() *RealClock

NewClock returns a new RealClock. 

####  func (*RealClock) [NewTicker](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L31) ¶
    
    
    func (c *RealClock) NewTicker(d [time](/time).[Duration](/time#Duration)) (<-chan [time](/time).[Time](/time#Time), func())

NewTicker returns the channel and stop function from the ticker from the standard library. 

####  func (*RealClock) [NewTimer](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L37) ¶
    
    
    func (c *RealClock) NewTimer(d [time](/time).[Duration](/time#Duration)) (<-chan [time](/time).[Time](/time#Time), func() [bool](/builtin#bool))

NewTimer returns the channel and stop function from the timer from the standard library. 

####  func (*RealClock) [Now](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L26) ¶
    
    
    func (c *RealClock) Now() [time](/time).[Time](/time#Time)

Now wraps time.Now() from the standard library. 

####  func (*RealClock) [Since](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L43) ¶
    
    
    func (c *RealClock) Since(t [time](/time).[Time](/time#Time)) [time](/time).[Duration](/time#Duration)

Since wraps time.Since() from the standard library. 

####  type [TestClock](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L61) ¶
    
    
    type TestClock struct {
    	// contains filtered or unexported fields
    }

TestClock is a clock that satisfies the Clock interface. It should only be used for testing. 

####  func [NewTestClock](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L69) ¶
    
    
    func NewTestClock(time [time](/time).[Time](/time#Time)) *TestClock

NewTestClock returns a new TestClock with the specified time. 

####  func (*TestClock) [Add](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L113) ¶
    
    
    func (c *TestClock) Add(d [time](/time).[Duration](/time#Duration))

Add adds the duration to the internal time of the test clock and triggers any timers or tickers that should fire. 

####  func (*TestClock) [NewTicker](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L126) ¶
    
    
    func (c *TestClock) NewTicker(d [time](/time).[Duration](/time#Duration)) (<-chan [time](/time).[Time](/time#Time), func())

NewTicker creates a new ticker that will fire every time the internal clock advances by the specified duration. 

####  func (*TestClock) [NewTimer](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L143) ¶
    
    
    func (c *TestClock) NewTimer(d [time](/time).[Duration](/time#Duration)) (<-chan [time](/time).[Time](/time#Time), func() [bool](/builtin#bool))

NewTimer creates a new timer that will fire once the internal time of the clock has been advanced passed the specified duration. 

####  func (*TestClock) [Now](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L118) ¶
    
    
    func (c *TestClock) Now() [time](/time).[Time](/time#Time)

Now returns the internal time of the test clock. 

####  func (*TestClock) [Set](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L78) ¶
    
    
    func (c *TestClock) Set(t [time](/time).[Time](/time#Time))

Set sets the internal time of the test clock and triggers any timers or tickers that should fire. 

####  func (*TestClock) [Since](https://github.com/viccon/sturdyc/blob/v1.1.5/clock.go#L166) ¶
    
    
    func (c *TestClock) Since(t [time](/time).[Time](/time#Time)) [time](/time).[Duration](/time#Duration)

Since returns the duration between the internal time of the clock and the specified time. 
