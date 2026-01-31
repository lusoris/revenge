# rueidis

> Source: https://pkg.go.dev/github.com/redis/rueidis
> Fetched: 2026-01-31T10:55:59.053575+00:00
> Content-Hash: 1ca7ae640d519aac
> Type: html

---

### Overview ¶

Package rueidis is a fast Golang Redis RESP3 client that does auto pipelining and supports client side caching. 

### Index ¶

  * Constants
  * Variables
  * func BinaryString(bs []byte) string
  * func DecodeSliceOfJSON[T any](result RedisResult, dest *[]T) error
  * func IsParseErr(err error) bool
  * func IsRedisBusyGroup(err error) bool
  * func IsRedisNil(err error) bool
  * func JSON(in any) string
  * func JsonMGet(client Client, ctx context.Context, keys []string, path string) (ret map[string]RedisMessage, err error)
  * func JsonMGetCache(client Client, ctx context.Context, ttl time.Duration, keys []string, ...) (ret map[string]RedisMessage, err error)
  * func JsonMSet(client Client, ctx context.Context, kvs map[string]string, path string) map[string]error
  * func MDel(client Client, ctx context.Context, keys []string) map[string]error
  * func MGet(client Client, ctx context.Context, keys []string) (ret map[string]RedisMessage, err error)
  * func MGetCache(client Client, ctx context.Context, ttl time.Duration, keys []string) (ret map[string]RedisMessage, err error)
  * func MSet(client Client, ctx context.Context, kvs map[string]string) map[string]error
  * func MSetNX(client Client, ctx context.Context, kvs map[string]string) map[string]error
  * func ToVector32(s string) []float32
  * func ToVector64(s string) []float64
  * func VectorString32(v []float32) string
  * func VectorString64(v []float64) string
  * func WithOnSubscriptionHook(ctx context.Context, hook func(PubSubSubscription)) context.Context
  * type AuthCredentials
  * type AuthCredentialsContext
  * type Builder
  * type CacheEntry
  * type CacheStore
  *     * func NewSimpleCacheAdapter(store SimpleCache) CacheStore
  * type CacheStoreOption
  * type Cacheable
  * type CacheableTTL
  *     * func CT(cmd Cacheable, ttl time.Duration) CacheableTTL
  * type Client
  *     * func NewClient(option ClientOption) (client Client, err error)
  * type ClientMode
  * type ClientOption
  *     * func MustParseURL(str string) ClientOption
    * func ParseURL(str string) (opt ClientOption, err error)
  * type ClusterOption
  * type Commands
  * type Completed
  * type CoreClient
  * type DedicatedClient
  * type FtSearchDoc
  * type GeoLocation
  * type Incomplete
  * type KeyValues
  * type KeyZScores
  * type Lua
  *     * func NewLuaScript(script string, opts ...LuaOption) *Lua
    * func NewLuaScriptNoSha(script string) *Lua
    * func NewLuaScriptNoShaRetryable(script string) *Lua
    * func NewLuaScriptReadOnly(script string, opts ...LuaOption) *Lua
    * func NewLuaScriptReadOnlyNoSha(script string) *Lua
    * func NewLuaScriptRetryable(script string, opts ...LuaOption) *Lua
  *     * func (s *Lua) Exec(ctx context.Context, c Client, keys, args []string) (resp RedisResult)
    * func (s *Lua) ExecMulti(ctx context.Context, c Client, multi ...LuaExec) (resp []RedisResult)
  * type LuaExec
  * type LuaOption
  *     * func WithLoadSHA1(enabled bool) LuaOption
  * type MultiRedisResultStream
  * type NewCacheStoreFn
  * type NodeInfo
  * type PubSubHooks
  * type PubSubMessage
  * type PubSubSubscription
  * type ReadNodeSelectorFunc
  *     * func AZAffinityNodeSelector(clientAZ string) ReadNodeSelectorFunc
    * func AZAffinityReplicasAndPrimaryNodeSelector(clientAZ string) ReadNodeSelectorFunc
    * func PreferReplicaNodeSelector() ReadNodeSelectorFunc
  * type RedirectMode
  * type RedisError
  *     * func IsRedisErr(err error) (ret *RedisError, ok bool)
  *     * func (r *RedisError) Error() string
    * func (r *RedisError) IsAsk() (addr string, ok bool)
    * func (r *RedisError) IsBusyGroup() bool
    * func (r *RedisError) IsClusterDown() bool
    * func (r *RedisError) IsLoading() bool
    * func (r *RedisError) IsMoved() (addr string, ok bool)
    * func (r *RedisError) IsNil() bool
    * func (r *RedisError) IsNoScript() bool
    * func (r *RedisError) IsRedirect() (addr string, ok bool)
    * func (r *RedisError) IsTryAgain() bool
  * type RedisMessage
  *     * func (m *RedisMessage) AsBool() (val bool, err error)
    * func (m *RedisMessage) AsBoolSlice() ([]bool, error)
    * func (m *RedisMessage) AsBytes() (bs []byte, err error)
    * func (m *RedisMessage) AsFloat64() (val float64, err error)
    * func (m *RedisMessage) AsFloatSlice() ([]float64, error)
    * func (m *RedisMessage) AsFtAggregate() (total int64, docs []map[string]string, err error)
    * func (m *RedisMessage) AsFtAggregateCursor() (cursor, total int64, docs []map[string]string, err error)
    * func (m *RedisMessage) AsFtSearch() (total int64, docs []FtSearchDoc, err error)
    * func (m *RedisMessage) AsGeosearch() ([]GeoLocation, error)
    * func (m *RedisMessage) AsInt64() (val int64, err error)
    * func (m *RedisMessage) AsIntMap() (map[string]int64, error)
    * func (m *RedisMessage) AsIntSlice() ([]int64, error)
    * func (m *RedisMessage) AsLMPop() (kvs KeyValues, err error)
    * func (m *RedisMessage) AsMap() (map[string]RedisMessage, error)
    * func (m *RedisMessage) AsReader() (reader io.Reader, err error)
    * func (m *RedisMessage) AsScanEntry() (e ScanEntry, err error)
    * func (m *RedisMessage) AsStrMap() (map[string]string, error)
    * func (m *RedisMessage) AsStrSlice() ([]string, error)
    * func (m *RedisMessage) AsUint64() (val uint64, err error)
    * func (m *RedisMessage) AsXRange() ([]XRangeEntry, error)
    * func (m *RedisMessage) AsXRangeEntry() (XRangeEntry, error)
    * func (m *RedisMessage) AsXRangeSlice() (XRangeSlice, error)
    * func (m *RedisMessage) AsXRangeSlices() ([]XRangeSlice, error)
    * func (m *RedisMessage) AsXRead() (ret map[string][]XRangeEntry, err error)
    * func (m *RedisMessage) AsXReadSlices() (map[string][]XRangeSlice, error)
    * func (m *RedisMessage) AsZMPop() (kvs KeyZScores, err error)
    * func (m *RedisMessage) AsZScore() (s ZScore, err error)
    * func (m *RedisMessage) AsZScores() ([]ZScore, error)
    * func (m *RedisMessage) CacheMarshal(buf []byte) []byte
    * func (m *RedisMessage) CachePTTL() int64
    * func (m *RedisMessage) CachePXAT() int64
    * func (m *RedisMessage) CacheSize() int
    * func (m *RedisMessage) CacheTTL() (ttl int64)
    * func (m *RedisMessage) CacheUnmarshalView(buf []byte) error
    * func (m *RedisMessage) DecodeJSON(v any) (err error)
    * func (m *RedisMessage) Error() error
    * func (m *RedisMessage) IsArray() bool
    * func (m *RedisMessage) IsBool() bool
    * func (m *RedisMessage) IsCacheHit() bool
    * func (m *RedisMessage) IsFloat64() bool
    * func (m *RedisMessage) IsInt64() bool
    * func (m *RedisMessage) IsMap() bool
    * func (m *RedisMessage) IsNil() bool
    * func (m *RedisMessage) IsString() bool
    * func (m *RedisMessage) String() string
    * func (m *RedisMessage) ToAny() (any, error)
    * func (m *RedisMessage) ToArray() ([]RedisMessage, error)
    * func (m *RedisMessage) ToBool() (val bool, err error)
    * func (m *RedisMessage) ToFloat64() (val float64, err error)
    * func (m *RedisMessage) ToInt64() (val int64, err error)
    * func (m *RedisMessage) ToMap() (map[string]RedisMessage, error)
    * func (m *RedisMessage) ToString() (val string, err error)
  * type RedisResult
  *     * func (r RedisResult) AsBool() (v bool, err error)
    * func (r RedisResult) AsBoolSlice() (v []bool, err error)
    * func (r RedisResult) AsBytes() (v []byte, err error)
    * func (r RedisResult) AsFloat64() (v float64, err error)
    * func (r RedisResult) AsFloatSlice() (v []float64, err error)
    * func (r RedisResult) AsFtAggregate() (total int64, docs []map[string]string, err error)
    * func (r RedisResult) AsFtAggregateCursor() (cursor, total int64, docs []map[string]string, err error)
    * func (r RedisResult) AsFtSearch() (total int64, docs []FtSearchDoc, err error)
    * func (r RedisResult) AsGeosearch() (locations []GeoLocation, err error)
    * func (r RedisResult) AsInt64() (v int64, err error)
    * func (r RedisResult) AsIntMap() (v map[string]int64, err error)
    * func (r RedisResult) AsIntSlice() (v []int64, err error)
    * func (r RedisResult) AsLMPop() (v KeyValues, err error)
    * func (r RedisResult) AsMap() (v map[string]RedisMessage, err error)
    * func (r RedisResult) AsReader() (v io.Reader, err error)
    * func (r RedisResult) AsScanEntry() (v ScanEntry, err error)
    * func (r RedisResult) AsStrMap() (v map[string]string, err error)
    * func (r RedisResult) AsStrSlice() (v []string, err error)
    * func (r RedisResult) AsUint64() (v uint64, err error)
    * func (r RedisResult) AsXRange() (v []XRangeEntry, err error)
    * func (r RedisResult) AsXRangeEntry() (v XRangeEntry, err error)
    * func (r RedisResult) AsXRangeSlice() (v XRangeSlice, err error)
    * func (r RedisResult) AsXRangeSlices() (v []XRangeSlice, err error)
    * func (r RedisResult) AsXRead() (v map[string][]XRangeEntry, err error)
    * func (r RedisResult) AsXReadSlices() (v map[string][]XRangeSlice, err error)
    * func (r RedisResult) AsZMPop() (v KeyZScores, err error)
    * func (r RedisResult) AsZScore() (v ZScore, err error)
    * func (r RedisResult) AsZScores() (v []ZScore, err error)
    * func (r RedisResult) CachePTTL() int64
    * func (r RedisResult) CachePXAT() int64
    * func (r RedisResult) CacheTTL() int64
    * func (r RedisResult) DecodeJSON(v any) (err error)
    * func (r RedisResult) Error() (err error)
    * func (r RedisResult) IsCacheHit() bool
    * func (r RedisResult) NonRedisError() error
    * func (r *RedisResult) String() string
    * func (r RedisResult) ToAny() (v any, err error)
    * func (r RedisResult) ToArray() (v []RedisMessage, err error)
    * func (r RedisResult) ToBool() (v bool, err error)
    * func (r RedisResult) ToFloat64() (v float64, err error)
    * func (r RedisResult) ToInt64() (v int64, err error)
    * func (r RedisResult) ToMap() (v map[string]RedisMessage, err error)
    * func (r RedisResult) ToMessage() (v RedisMessage, err error)
    * func (r RedisResult) ToString() (v string, err error)
  * type RedisResultStream
  *     * func (s *RedisResultStream) Error() error
    * func (s *RedisResultStream) HasNext() bool
    * func (s *RedisResultStream) WriteTo(w io.Writer) (n int64, err error)
  * type ReplicaInfo
  * type ReplicaSelectorFunc
  * type RetryDelayFn
  * type ScanEntry
  * type Scanner
  *     * func NewScanner(next func(cursor uint64) (ScanEntry, error)) *Scanner
  *     * func (s *Scanner) Err() error
    * func (s *Scanner) Iter() iter.Seq[string]
    * func (s *Scanner) Iter2() iter.Seq2[string, string]
  * type SentinelOption
  * type SimpleCache
  * type StandaloneOption
  * type XRangeEntry
  * type XRangeFieldValue
  * type XRangeSlice
  * type ZScore



### Examples ¶

  * Client (DedicateCAS)
  * Client (DedicatedCAS)
  * Client (Do)
  * Client (DoCache)
  * Client (Scan)
  * IsRedisNil
  * Lua (Exec)
  * NewClient (Cluster)
  * NewClient (Sentinel)
  * NewClient (Single)



### Constants ¶

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L45)
    
    
    const (
    	// DefaultCacheBytes is the default value of ClientOption.CacheSizeEachConn, which is 128 MiB
    	DefaultCacheBytes = 128 * (1 << 20)
    	// DefaultRingScale is the default value of ClientOption.RingScaleEachConn, which results into having a ring of size 2^10 for each connection
    	DefaultRingScale = 10
    	// DefaultPoolSize is the default value of ClientOption.BlockingPoolSize
    	DefaultPoolSize = 1024
    	// DefaultBlockingPipeline is the default value of ClientOption.BlockingPipeline
    	DefaultBlockingPipeline = 2000
    	// DefaultDialTimeout is the default value of ClientOption.Dialer.Timeout
    	DefaultDialTimeout = 5 * [time](/time).[Second](/time#Second)
    	// DefaultTCPKeepAlive is the default value of ClientOption.Dialer.KeepAlive
    	DefaultTCPKeepAlive = 1 * [time](/time).[Second](/time#Second)
    	// DefaultReadBuffer is the default value of bufio.NewReaderSize for each connection, which is 0.5MiB
    	DefaultReadBuffer = 1 << 19
    	// DefaultWriteBuffer is the default value of bufio.NewWriterSize for each connection, which is 0.5MiB
    	DefaultWriteBuffer = 1 << 19
    	// MaxPipelineMultiplex is the maximum meaningful value for ClientOption.PipelineMultiplex
    	MaxPipelineMultiplex = 8
    	// <https://github.com/valkey-io/valkey/blob/1a34a4ff7f101bb6b17a0b5e9aa3bf7d6bd29f68/src/networking.c#L4118-L4124>
    	ClientModeCluster    ClientMode = "cluster"
    	ClientModeSentinel   ClientMode = "sentinel"
    	ClientModeStandalone ClientMode = "standalone"
    )

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L23)
    
    
    const LibName = "rueidis"

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L24)
    
    
    const LibVer = "1.0.71"

### Variables ¶

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L70)
    
    
    var (
    	// ErrClosing means the Client.Close had been called
    	ErrClosing = [errors](/errors).[New](/errors#New)("rueidis client is closing or unable to connect redis")
    	// ErrNoAddr means the ClientOption.InitAddress is empty
    	ErrNoAddr = [errors](/errors).[New](/errors#New)("no alive address in InitAddress")
    	// ErrNoCache means your redis does not support client-side caching and must set ClientOption.DisableCache to true
    	ErrNoCache = [errors](/errors).[New](/errors#New)("ClientOption.DisableCache must be true for redis not supporting client-side caching or not supporting RESP3")
    	// ErrRESP2PubSubMixed means your redis does not support RESP3 and rueidis can't handle SUBSCRIBE/PSUBSCRIBE/SSUBSCRIBE in mixed case
    	ErrRESP2PubSubMixed = [errors](/errors).[New](/errors#New)("rueidis does not support SUBSCRIBE/PSUBSCRIBE/SSUBSCRIBE mixed with other commands in RESP2")
    	// ErrBlockingPubSubMixed rueidis can't handle SUBSCRIBE/PSUBSCRIBE/SSUBSCRIBE mixed with other blocking commands
    	ErrBlockingPubSubMixed = [errors](/errors).[New](/errors#New)("rueidis does not support SUBSCRIBE/PSUBSCRIBE/SSUBSCRIBE mixed with other blocking commands")
    	// ErrDoCacheAborted means redis abort EXEC request or connection closed
    	ErrDoCacheAborted = [errors](/errors).[New](/errors#New)("failed to fetch the cache because EXEC was aborted by redis or connection closed")
    	// ErrReplicaOnlyNotSupported means ReplicaOnly flag is not supported by
    	// the current client
    	ErrReplicaOnlyNotSupported = [errors](/errors).[New](/errors#New)("ReplicaOnly is not supported for single client")
    	// ErrNoSendToReplicas means the SendToReplicas function must be provided for a standalone client with replicas.
    	ErrNoSendToReplicas = [errors](/errors).[New](/errors#New)("no SendToReplicas provided for standalone client with replicas")
    	// ErrWrongPipelineMultiplex means wrong value for ClientOption.PipelineMultiplex
    	ErrWrongPipelineMultiplex = [errors](/errors).[New](/errors#New)("ClientOption.PipelineMultiplex must not be bigger than MaxPipelineMultiplex")
    	// ErrDedicatedClientRecycled means the caller attempted to use the dedicated client which has been already recycled (after canceled/closed).
    	ErrDedicatedClientRecycled = [errors](/errors).[New](/errors#New)("dedicated client should not be used after recycled")
    	// DisableClientSetInfo is the value that can be used for ClientOption.ClientSetInfo to disable making the CLIENT SETINFO command
    	DisableClientSetInfo = [make](/builtin#make)([][string](/builtin#string), 0)
    )

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L641)
    
    
    var ErrCacheUnmarshal = [errors](/errors).[New](/errors#New)("cache unmarshal error")

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L20)
    
    
    var ErrInvalidShardsRefreshInterval = [errors](/errors).[New](/errors#New)("ShardsRefreshInterval must be greater than or equal to 0")

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L348)
    
    
    var ErrMSetNXNotSet = [errors](/errors).[New](/errors#New)("MSETNX: no key was set")

ErrMSetNXNotSet is used in the MSetNX helper when the underlying MSETNX response is 0. Ref: <https://redis.io/commands/msetnx/>

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L18)
    
    
    var ErrNoSlot = [errors](/errors).[New](/errors#New)("the slot has no redis node")

ErrNoSlot indicates that there is no redis node owning the key slot. 

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L19)
    
    
    var ErrReplicaOnlyConflict = [errors](/errors).[New](/errors#New)("ReplicaOnly conflicts with SendToReplicas option")

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L22)
    
    
    var ErrReplicaOnlyConflictWithReadNodeSelector = [errors](/errors).[New](/errors#New)("ReplicaOnly conflicts with ReadNodeSelector option")

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L21)
    
    
    var ErrReplicaOnlyConflictWithReplicaSelector = [errors](/errors).[New](/errors#New)("ReplicaOnly conflicts with ReplicaSelector option")

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L23)
    
    
    var ErrReplicaSelectorConflictWithReadNodeSelector = [errors](/errors).[New](/errors#New)("either set ReplicaSelector or ReadNodeSelector, not both")

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L24)
    
    
    var ErrSendToReplicasNotSet = [errors](/errors).[New](/errors#New)("SendToReplicas must be set when ReplicaSelector is set")

[View Source](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L22)
    
    
    var Nil = &RedisError{typ: typeNull}

Nil represents a Redis Nil message 

### Functions ¶

####  func [BinaryString](https://github.com/redis/rueidis/blob/v1.0.71/binary.go#L17) ¶
    
    
    func BinaryString(bs [][byte](/builtin#byte)) [string](/builtin#string)

BinaryString convert the provided []byte into a string without a copy. It does what strings.Builder.String() does. Redis Strings are binary safe; this means that it is safe to store any []byte into Redis directly. Users can use this BinaryString helper to insert a []byte as the part of redis command. For example: 
    
    
    client.B().Set().Key(rueidis.BinaryString([]byte{0})).Value(rueidis.BinaryString([]byte{0})).Build()
    

To read back the []byte of the string returned from the Redis, it is recommended to use the RedisMessage.AsReader. 

####  func [DecodeSliceOfJSON](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L161) ¶ added in v1.0.34
    
    
    func DecodeSliceOfJSON[T [any](/builtin#any)](result RedisResult, dest *[]T) [error](/builtin#error)

DecodeSliceOfJSON is a helper that struct-scans each RedisMessage into dest, which must be a slice of the pointer. 

####  func [IsParseErr](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L34) ¶ added in v1.0.40
    
    
    func IsParseErr(err [error](/builtin#error)) [bool](/builtin#bool)

IsParseErr checks if the error is a parse error 

####  func [IsRedisBusyGroup](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L39) ¶ added in v1.0.32
    
    
    func IsRedisBusyGroup(err [error](/builtin#error)) [bool](/builtin#bool)

IsRedisBusyGroup checks if it is a redis BUSYGROUP message. 

####  func [IsRedisNil](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L29) ¶
    
    
    func IsRedisNil(err [error](/builtin#error)) [bool](/builtin#bool)

IsRedisNil is a handy method to check if the error is a redis nil response. All redis nil responses returned as an error. 

Example ¶
    
    
    client, err := NewClient(ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
    	panic(err)
    }
    defer client.Close()
    
    _, err = client.Do(context.Background(), client.B().Get().Key("not_exists").Build()).ToString()
    if err != nil && IsRedisNil(err) {
    	fmt.Printf("it is a nil response")
    }
    

####  func [JSON](https://github.com/redis/rueidis/blob/v1.0.71/binary.go#L73) ¶
    
    
    func JSON(in [any](/builtin#any)) [string](/builtin#string)

JSON convert the provided parameter into a JSON string. Users can use this JSON helper to work with RedisJSON commands. For example: 
    
    
    client.B().JsonSet().Key("a").Path("$.myField").Value(rueidis.JSON("str")).Build()
    

####  func [JsonMGet](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L128) ¶
    
    
    func JsonMGet(client Client, ctx [context](/context).[Context](/context#Context), keys [][string](/builtin#string), path [string](/builtin#string)) (ret map[[string](/builtin#string)]RedisMessage, err [error](/builtin#error))

JsonMGet is a helper that consults redis directly with multiple keys by grouping keys within the same slot into JSON.MGETs or multiple JSON.GETs 

####  func [JsonMGetCache](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L115) ¶
    
    
    func JsonMGetCache(client Client, ctx [context](/context).[Context](/context#Context), ttl [time](/time).[Duration](/time#Duration), keys [][string](/builtin#string), path [string](/builtin#string)) (ret map[[string](/builtin#string)]RedisMessage, err [error](/builtin#error))

JsonMGetCache is a helper that consults the client-side caches with multiple keys by grouping keys within the same slot into multiple JSON.GETs 

####  func [JsonMSet](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L142) ¶ added in v1.0.3
    
    
    func JsonMSet(client Client, ctx [context](/context).[Context](/context#Context), kvs map[[string](/builtin#string)][string](/builtin#string), path [string](/builtin#string)) map[[string](/builtin#string)][error](/builtin#error)

JsonMSet is a helper that consults redis directly with multiple keys by grouping keys within the same slot into JSON.MSETs or multiple JSON.SETs 

####  func [MDel](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L77) ¶ added in v1.0.3
    
    
    func MDel(client Client, ctx [context](/context).[Context](/context#Context), keys [][string](/builtin#string)) map[[string](/builtin#string)][error](/builtin#error)

MDel is a helper that consults the redis directly with multiple keys by grouping keys within the same slot into DELs 

####  func [MGet](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L44) ¶
    
    
    func MGet(client Client, ctx [context](/context).[Context](/context#Context), keys [][string](/builtin#string)) (ret map[[string](/builtin#string)]RedisMessage, err [error](/builtin#error))

MGet is a helper that consults the redis directly with multiple keys by grouping keys within the same slot into MGET or multiple GETs 

####  func [MGetCache](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L14) ¶
    
    
    func MGetCache(client Client, ctx [context](/context).[Context](/context#Context), ttl [time](/time).[Duration](/time#Duration), keys [][string](/builtin#string)) (ret map[[string](/builtin#string)]RedisMessage, err [error](/builtin#error))

MGetCache is a helper that consults the client-side caches with multiple keys by grouping keys within the same slot into multiple GETs 

####  func [MSet](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L58) ¶
    
    
    func MSet(client Client, ctx [context](/context).[Context](/context#Context), kvs map[[string](/builtin#string)][string](/builtin#string)) map[[string](/builtin#string)][error](/builtin#error)

MSet is a helper that consults the redis directly with multiple keys by grouping keys within the same slot into MSETs or multiple SETs 

####  func [MSetNX](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L96) ¶
    
    
    func MSetNX(client Client, ctx [context](/context).[Context](/context#Context), kvs map[[string](/builtin#string)][string](/builtin#string)) map[[string](/builtin#string)][error](/builtin#error)

MSetNX is a helper that consults the redis directly with multiple keys by grouping keys within the same slot into MSETNXs or multiple SETNXs 

####  func [ToVector32](https://github.com/redis/rueidis/blob/v1.0.71/binary.go#L36) ¶
    
    
    func ToVector32(s [string](/builtin#string)) [][float32](/builtin#float32)

ToVector32 reverts VectorString32. User can use this to convert redis response back to []float32. 

####  func [ToVector64](https://github.com/redis/rueidis/blob/v1.0.71/binary.go#L60) ¶
    
    
    func ToVector64(s [string](/builtin#string)) [][float64](/builtin#float64)

ToVector64 reverts VectorString64. User can use this to convert redis response back to []float64. 

####  func [VectorString32](https://github.com/redis/rueidis/blob/v1.0.71/binary.go#L26) ¶
    
    
    func VectorString32(v [][float32](/builtin#float32)) [string](/builtin#string)

VectorString32 convert the provided []float32 into a string. Users can use this to build vector search queries: 
    
    
    client.B().FtSearch().Index("idx").Query("*=>[KNN 5 @vec $V]").
        Params().Nargs(2).NameValue().NameValue("V", rueidis.VectorString32([]float32{1})).
        Dialect(2).Build()
    

####  func [VectorString64](https://github.com/redis/rueidis/blob/v1.0.71/binary.go#L50) ¶
    
    
    func VectorString64(v [][float64](/builtin#float64)) [string](/builtin#string)

VectorString64 convert the provided []float64 into a string. Users can use this to build vector search queries: 
    
    
    client.B().FtSearch().Index("idx").Query("*=>[KNN 5 @vec $V]").
        Params().Nargs(2).NameValue().NameValue("V", rueidis.VectorString64([]float64{1})).
        Dialect(2).Build()
    

####  func [WithOnSubscriptionHook](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L825) ¶ added in v1.0.61
    
    
    func WithOnSubscriptionHook(ctx [context](/context).[Context](/context#Context), hook func(PubSubSubscription)) [context](/context).[Context](/context#Context)

WithOnSubscriptionHook attaches a subscription confirmation hook to the provided context and returns a new context for the Receive method. 

The hook is invoked each time the server sends a subscribe or unsubscribe confirmation, allowing callers to observe the state of a Pub/Sub subscription during the lifetime of a Receive invocation. 

The hook may be called multiple times because the client can resubscribe after a reconnection. Therefore, the hook implementation must be safe to run more than once. Also, there should not be any blocking operations or another `client.Do()` in the hook since it runs in the same goroutine as the pipeline. Otherwise, the pipeline will be blocked. 

### Types ¶

####  type [AuthCredentials](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L461) ¶ added in v1.0.19
    
    
    type AuthCredentials struct {
    	Username [string](/builtin#string)
    	Password [string](/builtin#string)
    }

AuthCredentials is the output of AuthCredentialsFn 

####  type [AuthCredentialsContext](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L456) ¶ added in v1.0.19
    
    
    type AuthCredentialsContext struct {
    	Address [net](/net).[Addr](/net#Addr)
    }

AuthCredentialsContext is the parameter container of AuthCredentialsFn 

####  type [Builder](https://github.com/redis/rueidis/blob/v1.0.71/cmds.go#L6) ¶ added in v1.0.14
    
    
    type Builder = [cmds](/github.com/redis/rueidis@v1.0.71/internal/cmds).[Builder](/github.com/redis/rueidis@v1.0.71/internal/cmds#Builder)

Builder represents a command builder. It should only be created from the client.B() method. 

####  type [CacheEntry](https://github.com/redis/rueidis/blob/v1.0.71/cache.go#L45) ¶
    
    
    type CacheEntry interface {
    	Wait(ctx [context](/context).[Context](/context#Context)) (RedisMessage, [error](/builtin#error))
    }

CacheEntry should be used to wait for a single-flight response when cache missed. 

####  type [CacheStore](https://github.com/redis/rueidis/blob/v1.0.71/cache.go#L21) ¶
    
    
    type CacheStore interface {
    	// Flight is called when DoCache and DoMultiCache, with the requested client side ttl and the current time.
    	// It should look up the store in a single-flight manner and return one of the following three combinations:
    	// Case 1: (empty RedisMessage, nil CacheEntry)     <- when cache missed, and rueidis will send the request to redis.
    	// Case 2: (empty RedisMessage, non-nil CacheEntry) <- when cache missed, and rueidis will use CacheEntry.Wait to wait for response.
    	// Case 3: (non-empty RedisMessage, nil CacheEntry) <- when cache hit
    	Flight(key, cmd [string](/builtin#string), ttl [time](/time).[Duration](/time#Duration), now [time](/time).[Time](/time#Time)) (v RedisMessage, e CacheEntry)
    	// Update is called when receiving the response of the request sent by the above Flight Case 1 from redis.
    	// It should not only update the store but also deliver the response to all CacheEntry.Wait and return a desired client side PXAT of the response.
    	// Note that the server side expire time can be retrieved from RedisMessage.CachePXAT.
    	Update(key, cmd [string](/builtin#string), val RedisMessage) (pxat [int64](/builtin#int64))
    	// Cancel is called when the request sent by the above Flight Case 1 failed.
    	// It should not only deliver the error to all CacheEntry.Wait but also remove the CacheEntry from the store.
    	Cancel(key, cmd [string](/builtin#string), err [error](/builtin#error))
    	// Delete is called when receiving invalidation notifications from redis.
    	// If the keys are nil, then it should delete all non-pending cached entries under all keys.
    	// If the keys are not nil, then it should delete all non-pending cached entries under those keys.
    	Delete(keys []RedisMessage)
    	// Close is called when the connection between redis is broken.
    	// It should flush all cached entries and deliver the error to all pending CacheEntry.Wait.
    	Close(err [error](/builtin#error))
    }

CacheStore is the store interface for the client side caching More detailed interface requirement can be found in cache_test.go 

####  func [NewSimpleCacheAdapter](https://github.com/redis/rueidis/blob/v1.0.71/cache.go#L58) ¶ added in v1.0.1
    
    
    func NewSimpleCacheAdapter(store SimpleCache) CacheStore

NewSimpleCacheAdapter converts a SimpleCache into CacheStore 

####  type [CacheStoreOption](https://github.com/redis/rueidis/blob/v1.0.71/cache.go#L13) ¶
    
    
    type CacheStoreOption struct {
    	// CacheSizeEachConn is redis client side cache size that bind to each TCP connection to a single redis instance.
    	// The default is DefaultCacheBytes.
    	CacheSizeEachConn [int](/builtin#int)
    }

CacheStoreOption will be passed to NewCacheStoreFn 

####  type [Cacheable](https://github.com/redis/rueidis/blob/v1.0.71/cmds.go#L16) ¶
    
    
    type Cacheable = [cmds](/github.com/redis/rueidis@v1.0.71/internal/cmds).[Cacheable](/github.com/redis/rueidis@v1.0.71/internal/cmds#Cacheable)

Cacheable represents a completed Redis command which supports server-assisted client side caching, and it should be created by the Cache() of command builder. 

####  type [CacheableTTL](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L450) ¶
    
    
    type CacheableTTL struct {
    	Cmd Cacheable
    	TTL [time](/time).[Duration](/time#Duration)
    }

CacheableTTL is a parameter container of DoMultiCache 

####  func [CT](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L445) ¶
    
    
    func CT(cmd Cacheable, ttl [time](/time).[Duration](/time#Duration)) CacheableTTL

CT is a shorthand constructor for CacheableTTL 

####  type [Client](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L345) ¶
    
    
    type Client interface {
    	CoreClient
    
    	// DoCache is similar to Do, but it uses opt-in client side caching and requires a client side TTL.
    	// The explicit client side TTL specifies the maximum TTL on the client side.
    	// If the key's TTL on the server is smaller than the client side TTL, the client side TTL will be capped.
    	//  client.Do(ctx, client.B().Get().Key("k").Cache(), time.Minute).ToString()
    	// The above example will send the following command to redis if the cache misses:
    	//  CLIENT CACHING YES
    	//  PTTL k
    	//  GET k
    	// The in-memory cache size is configured by ClientOption.CacheSizeEachConn.
    	// The cmd parameter is recycled after passing into DoCache() and should not be reused.
    	DoCache(ctx [context](/context).[Context](/context#Context), cmd Cacheable, ttl [time](/time).[Duration](/time#Duration)) (resp RedisResult)
    
    	// DoMultiCache is similar to DoCache but works with multiple cacheable commands across different slots.
    	// It will first group commands by slots and will send only cache missed commands to redis.
    	DoMultiCache(ctx [context](/context).[Context](/context#Context), multi ...CacheableTTL) (resp []RedisResult)
    
    	// DoStream send a command to redis through a dedicated connection acquired from a connection pool.
    	// It returns a RedisResultStream, but it does not read the command response until the RedisResultStream.WriteTo is called.
    	// After the RedisResultStream.WriteTo is called, the underlying connection is then recycled.
    	// DoStream should only be used when you want to stream redis response directly to an io.Writer without additional allocation,
    	// otherwise, the normal Do() should be used instead.
    	// Also note that DoStream can only work with commands returning string, integer, or float response.
    	DoStream(ctx [context](/context).[Context](/context#Context), cmd Completed) RedisResultStream
    
    	// DoMultiStream is similar to DoStream, but pipelines multiple commands to redis.
    	// It returns a MultiRedisResultStream, and users should call MultiRedisResultStream.WriteTo as many times as the number of commands sequentially
    	// to read each command response from redis. After all responses are read, the underlying connection is then recycled.
    	// DoMultiStream should only be used when you want to stream redis responses directly to an io.Writer without additional allocation,
    	// otherwise, the normal DoMulti() should be used instead.
    	// DoMultiStream does not support multiple key slots when connecting to a redis cluster.
    	DoMultiStream(ctx [context](/context).[Context](/context#Context), multi ...Completed) MultiRedisResultStream
    
    	// Dedicated acquire a connection from the blocking connection pool, no one else can use the connection
    	// during Dedicated. The main usage of Dedicated is CAS operations, which is WATCH + MULTI + EXEC.
    	// However, one should try to avoid CAS operation but use a Lua script instead, because occupying a connection
    	// is not good for performance.
    	Dedicated(fn func(DedicatedClient) [error](/builtin#error)) (err [error](/builtin#error))
    
    	// Dedicate does the same as Dedicated, but it exposes DedicatedClient directly
    	// and requires user to invoke cancel() manually to put connection back to the pool.
    	Dedicate() (client DedicatedClient, cancel func())
    
    	// Nodes returns each redis node this client known as rueidis.Client. This is useful if you want to
    	// send commands to some specific redis nodes in the cluster.
    	Nodes() map[[string](/builtin#string)]Client
    	// Mode returns the current mode of the client, which indicates whether the client is operating
    	// in standalone, sentinel, or cluster mode.
    	// This can be useful for determining the type of Redis deployment the client is connected to
    	// and for making decisions based on the deployment type.
    	Mode() ClientMode
    }

Client is the redis client interface for both single redis instance and redis cluster. It should be created from the NewClient() 

Example (DedicateCAS) ¶
    
    
    client, err := NewClient(ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
    	panic(err)
    }
    defer client.Close()
    
    c, cancel := client.Dedicate()
    defer cancel()
    
    ctx := context.Background()
    
    // watch keys first
    if err := c.Do(ctx, c.B().Watch().Key("k1", "k2").Build()).Error(); err != nil {
    	panic(err)
    }
    // perform read here
    values, err := c.Do(ctx, c.B().Mget().Key("k1", "k2").Build()).ToArray()
    if err != nil {
    	panic(err)
    }
    v1, _ := values[0].ToString()
    v2, _ := values[1].ToString()
    // perform write with MULTI EXEC
    for _, resp := range c.DoMulti(
    	ctx,
    	c.B().Multi().Build(),
    	c.B().Set().Key("k1").Value(v1+"1").Build(),
    	c.B().Set().Key("k2").Value(v2+"2").Build(),
    	c.B().Exec().Build(),
    ) {
    	if err := resp.Error(); err != nil {
    		panic(err)
    	}
    }
    

Example (DedicatedCAS) ¶
    
    
    client, err := NewClient(ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
    	panic(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    client.Dedicated(func(client DedicatedClient) error {
    	// watch keys first
    	if err := client.Do(ctx, client.B().Watch().Key("k1", "k2").Build()).Error(); err != nil {
    		return err
    	}
    	// perform read here
    	values, err := client.Do(ctx, client.B().Mget().Key("k1", "k2").Build()).ToArray()
    	if err != nil {
    		return err
    	}
    	v1, _ := values[0].ToString()
    	v2, _ := values[1].ToString()
    	// perform write with MULTI EXEC
    	for _, resp := range client.DoMulti(
    		ctx,
    		client.B().Multi().Build(),
    		client.B().Set().Key("k1").Value(v1+"1").Build(),
    		client.B().Set().Key("k2").Value(v2+"2").Build(),
    		client.B().Exec().Build(),
    	) {
    		if err := resp.Error(); err != nil {
    			return err
    		}
    	}
    	return nil
    })
    

Example (Do) ¶
    
    
    client, err := NewClient(ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
    	panic(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    client.Do(ctx, client.B().Set().Key("k").Value("1").Build()).Error()
    
    client.Do(ctx, client.B().Get().Key("k").Build()).ToString()
    
    client.Do(ctx, client.B().Get().Key("k").Build()).AsInt64()
    
    client.Do(ctx, client.B().Hmget().Key("h").Field("a", "b").Build()).ToArray()
    
    client.Do(ctx, client.B().Scard().Key("s").Build()).ToInt64()
    
    client.Do(ctx, client.B().Smembers().Key("s").Build()).AsStrSlice()
    

Example (DoCache) ¶
    
    
    client, err := NewClient(ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
    	panic(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    client.DoCache(ctx, client.B().Get().Key("k").Cache(), time.Minute).ToString()
    
    client.DoCache(ctx, client.B().Get().Key("k").Cache(), time.Minute).AsInt64()
    
    client.DoCache(ctx, client.B().Hmget().Key("h").Field("a", "b").Cache(), time.Minute).ToArray()
    
    client.DoCache(ctx, client.B().Scard().Key("s").Cache(), time.Minute).ToInt64()
    
    client.DoCache(ctx, client.B().Smembers().Key("s").Cache(), time.Minute).AsStrSlice()
    

Example (Scan) ¶
    
    
    client, err := NewClient(ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
    	panic(err)
    }
    defer client.Close()
    
    for _, c := range client.Nodes() { // loop over all your redis nodes
    	var scan ScanEntry
    	for more := true; more; more = scan.Cursor != 0 {
    		if scan, err = c.Do(context.Background(), c.B().Scan().Cursor(scan.Cursor).Build()).AsScanEntry(); err != nil {
    			panic(err)
    		}
    		fmt.Println(scan.Elements)
    	}
    }
    

####  func [NewClient](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L469) ¶
    
    
    func NewClient(option ClientOption) (client Client, err [error](/builtin#error))

NewClient uses ClientOption to initialize the Client for both a cluster client and a single client. It will first try to connect as a cluster client. If the len(ClientOption.InitAddress) == 1 and the address does not enable cluster mode, the NewClient() will use single client instead. 

Example (Cluster) ¶
    
    
    client, _ := NewClient(ClientOption{
    	InitAddress: []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003"},
    	ShuffleInit: true,
    })
    defer client.Close()
    

Example (Sentinel) ¶
    
    
    client, _ := NewClient(ClientOption{
    	InitAddress: []string{"127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"},
    	Sentinel: SentinelOption{
    		MasterSet: "my_master",
    	},
    })
    defer client.Close()
    

Example (Single) ¶
    
    
    client, _ := NewClient(ClientOption{
    	InitAddress: []string{"127.0.0.1:6379"},
    })
    defer client.Close()
    

####  type [ClientMode](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L342) ¶ added in v1.0.56
    
    
    type ClientMode [string](/builtin#string)

####  type [ClientOption](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L101) ¶
    
    
    type ClientOption struct {
    	TLSConfig *[tls](/crypto/tls).[Config](/crypto/tls#Config)
    
    	// DialFn allows for a custom function to be used to create net.Conn connections
    	// Deprecated: use DialCtxFn instead.
    	DialFn func([string](/builtin#string), *[net](/net).[Dialer](/net#Dialer), *[tls](/crypto/tls).[Config](/crypto/tls#Config)) (conn [net](/net).[Conn](/net#Conn), err [error](/builtin#error))
    
    	// DialCtxFn allows for a custom function to be used to create net.Conn connections
    	DialCtxFn func([context](/context).[Context](/context#Context), [string](/builtin#string), *[net](/net).[Dialer](/net#Dialer), *[tls](/crypto/tls).[Config](/crypto/tls#Config)) (conn [net](/net).[Conn](/net#Conn), err [error](/builtin#error))
    
    	// NewCacheStoreFn allows a custom client side caching store for each connection
    	NewCacheStoreFn NewCacheStoreFn
    
    	// OnInvalidations is a callback function in case of client-side caching invalidation received.
    	// Note that this function must be fast; otherwise other redis messages will be blocked.
    	OnInvalidations func([]RedisMessage)
    
    	// SendToReplicas is a function that returns true if the command should be sent to replicas.
    	// NOTE: This function can't be used with the ReplicaOnly option.
    	SendToReplicas func(cmd Completed) [bool](/builtin#bool)
    
    	// AuthCredentialsFn allows for setting the AUTH username and password dynamically on each connection attempt to
    	// support rotating credentials
    	AuthCredentialsFn func(AuthCredentialsContext) (AuthCredentials, [error](/builtin#error))
    
    	// RetryDelay is the function that returns the delay that should be used before retrying the attempt.
    	// The default is an exponential backoff with a maximum delay of 1 second.
    	// Only used when DisableRetry is false.
    	RetryDelay RetryDelayFn
    
    	// Deprecated: use ReadNodeSelector instead.
    	// ReplicaSelector selects a replica node when `SendToReplicas` returns true.
    	// If the function is set, the client will send the selected command to the replica node.
    	// The Returned value is the index of the replica node in the replica slice.
    	// If the returned value is out of range, the primary node will be selected.
    	// If the primary node does not have any replica, the primary node will be selected
    	// and the function will not be called.
    	// Currently only used for a cluster client.
    	// Each ReplicaInfo must not be modified.
    	// NOTE: This function can't be used with ReplicaOnly option.
    	// NOTE: This function must be used with the SendToReplicas function.
    	ReplicaSelector ReplicaSelectorFunc
    
    	// ReadNodeSelector returns index of node selected for a read only command.
    	// If set, ReadNodeSelector is prioritized over ReplicaSelector.
    	// If the returned index is out of range, the primary node will be selected.
    	// The function is called only when SendToReplicas returns true.
    	// Each NodeInfo must not be modified.
    	// NOTE: This function can't be used with ReplicaSelector option.
    	ReadNodeSelector ReadNodeSelectorFunc
    
    	// Sentinel options, including MasterSet and Auth options
    	Sentinel SentinelOption
    
    	// TCP & TLS
    	// Dialer can be used to customize how rueidis connect to a redis instance via TCP, including
    	// - Timeout, the default is DefaultDialTimeout
    	// - KeepAlive, the default is DefaultTCPKeepAlive
    	// The Dialer.KeepAlive interval is used to detect an unresponsive idle tcp connection.
    	// OS takes at least (tcp_keepalive_probes+1)*Dialer.KeepAlive time to conclude an idle connection to be unresponsive.
    	// For example, DefaultTCPKeepAlive = 1s and the default of tcp_keepalive_probes on Linux is 9.
    	// Therefore, it takes at least 10s to kill an idle and unresponsive tcp connection on Linux by default.
    	Dialer [net](/net).[Dialer](/net#Dialer)
    
    	// Redis AUTH parameters
    	Username   [string](/builtin#string)
    	Password   [string](/builtin#string)
    	ClientName [string](/builtin#string)
    
    	// ClientSetInfo will assign various info attributes to the current connection.
    	// Note that ClientSetInfo should have exactly 2 values, the lib name and the lib version respectively.
    	ClientSetInfo [][string](/builtin#string)
    
    	// InitAddress point to redis nodes.
    	// Rueidis will connect to them one by one and issue a CLUSTER SLOT command to initialize the cluster client until success.
    	// If len(InitAddress) == 1 and the address is not running in cluster mode, rueidis will fall back to the single client mode.
    	// If ClientOption.Sentinel.MasterSet is set, then InitAddress will be used to connect sentinels
    	// You can bypass this behavior by using ClientOption.ForceSingleClient.
    	InitAddress [][string](/builtin#string)
    
    	// ClientTrackingOptions will be appended to the CLIENT TRACKING ON command when the connection is established.
    	// The default is []string{"OPTIN"}
    	ClientTrackingOptions [][string](/builtin#string)
    
    	// Standalone is the option for the standalone client.
    	Standalone StandaloneOption
    
    	SelectDB [int](/builtin#int)
    
    	// CacheSizeEachConn is redis client side cache size that bind to each TCP connection to a single redis instance.
    	// The default is DefaultCacheBytes.
    	CacheSizeEachConn [int](/builtin#int)
    
    	// RingScaleEachConn sets the size of the ring buffer in each connection to (2 ^ RingScaleEachConn).
    	// The default is RingScaleEachConn, which results in having a ring of size 2^10 for each connection.
    	// Reducing this value can reduce the memory consumption of each connection at the cost of potential throughput degradation.
    	// Values smaller than 8 are typically not recommended.
    	RingScaleEachConn [int](/builtin#int)
    
    	// ReadBufferEachConn is the size of the bufio.NewReaderSize for each connection, default to DefaultReadBuffer (0.5 MiB).
    	ReadBufferEachConn [int](/builtin#int)
    	// WriteBufferEachConn is the size of the bufio.NewWriterSize for each connection, default to DefaultWriteBuffer (0.5 MiB).
    	WriteBufferEachConn [int](/builtin#int)
    
    	// BlockingPoolCleanup is the duration for cleaning up idle connections.
    	// If BlockingPoolCleanup is 0, then idle connections will not be cleaned up.
    	BlockingPoolCleanup [time](/time).[Duration](/time#Duration)
    	// BlockingPoolMinSize is the minimum size of the connection pool
    	// shared by blocking commands (ex BLPOP, XREAD with BLOCK).
    	// Only relevant if BlockingPoolCleanup is not 0. This parameter limits
    	// the number of idle connections that can be removed by BlockingPoolCleanup.
    	BlockingPoolMinSize [int](/builtin#int)
    
    	// BlockingPoolSize is the size of the connection pool shared by blocking commands (ex BLPOP, XREAD with BLOCK).
    	// The default is DefaultPoolSize.
    	BlockingPoolSize [int](/builtin#int)
    	// BlockingPipeline is the threshold of a pipeline that will be treated as blocking commands when exceeding it.
    	BlockingPipeline [int](/builtin#int)
    
    	// PipelineMultiplex determines how many tcp connections used to pipeline commands to one redis instance.
    	// The default for single and sentinel clients is 2, which means 4 connections (2^2).
    	// The default for cluster clients is 0, which means 1 connection (2^0).
    	PipelineMultiplex [int](/builtin#int)
    
    	// ConnWriteTimeout is a read/write timeout for each connection. If specified,
    	// it is used to control the maximum duration waits for responses to pipeline commands.
    	// Also, ConnWriteTimeout is applied net.Conn.SetDeadline and periodic PINGs,
    	// since the Dialer.KeepAlive will not be triggered if there is data in the outgoing buffer.
    	// ConnWriteTimeout should be set to detect local congestion or unresponsive redis server.
    	// This default is ClientOption.Dialer.KeepAlive * (9+1), where 9 is the default of tcp_keepalive_probes on Linux.
    	ConnWriteTimeout [time](/time).[Duration](/time#Duration)
    
    	// ConnLifetime is a lifetime for each connection. If specified,
    	// connections will close after passing lifetime. Note that the connection which a dedicated client and blocking use is not closed.
    	ConnLifetime [time](/time).[Duration](/time#Duration)
    
    	// MaxFlushDelay when greater than zero pauses pipeline write loop for some time (not larger than MaxFlushDelay)
    	// after each flushing of data to the connection. This gives the pipeline a chance to collect more commands to send
    	// to Redis. Adding this delay increases latency, reduces throughput – but in most cases may significantly reduce
    	// application and Redis CPU utilization due to less executed system calls. By default, Rueidis flushes data to the
    	// connection without extra delays. Depending on network latency and application-specific conditions, the value
    	// of MaxFlushDelay may vary, something like 20 microseconds should not affect latency/throughput a lot but still
    	// produce notable CPU usage reduction under load. Ref: <https://github.com/redis/rueidis/issues/156>
    	MaxFlushDelay [time](/time).[Duration](/time#Duration)
    
    	// ClusterOption is the options for the redis cluster client.
    	ClusterOption ClusterOption
    
    	// DisableTCPNoDelay turns on Nagle's algorithm in pipelining mode by using conn.SetNoDelay(false).
    	// Turning this on can result in lower p99 latencies and lower CPU usages if all your requests are small.
    	// But if you have large requests or fast network, this might degrade the performance. Ref: <https://github.com/redis/rueidis/pull/650>
    	DisableTCPNoDelay [bool](/builtin#bool)
    
    	// ShuffleInit is a handy flag that shuffles the InitAddress after passing to the NewClient() if it is true
    	ShuffleInit [bool](/builtin#bool)
    	// ClientNoTouch controls whether commands alter LRU/LFU stats
    	ClientNoTouch [bool](/builtin#bool)
    	// DisableRetry disables retrying read-only commands under network errors
    	DisableRetry [bool](/builtin#bool)
    	// DisableCache falls back Client.DoCache/Client.DoMultiCache to Client.Do/Client.DoMulti
    	DisableCache [bool](/builtin#bool)
    	// DisableAutoPipelining makes rueidis.Client always pick a connection from the BlockingPool to serve each request.
    	DisableAutoPipelining [bool](/builtin#bool)
    	// AlwaysPipelining makes rueidis.Client always pipeline redis commands even if they are not issued concurrently.
    	AlwaysPipelining [bool](/builtin#bool)
    	// AlwaysRESP2 makes rueidis.Client always uses RESP2; otherwise, it will try using RESP3 first.
    	AlwaysRESP2 [bool](/builtin#bool)
    	//  ForceSingleClient force the usage of a single client connection, without letting the lib guessing
    	//  if redis instance is a cluster or a single redis instance.
    	ForceSingleClient [bool](/builtin#bool)
    
    	// ReplicaOnly indicates that this client will only try to connect to readonly replicas of redis setup.
    	ReplicaOnly [bool](/builtin#bool)
    
    	// ClientNoEvict sets the client eviction mode for the current connection.
    	// When turned on and client eviction is configured,
    	// the current connection will be excluded from the client eviction process
    	// even if we're above the configured client eviction threshold.
    	ClientNoEvict [bool](/builtin#bool)
    
    	// EnableReplicaAZInfo enables the client to load the replica node's availability zone.
    	// If true, the client will set the `AZ` field in `ReplicaInfo`.
    	EnableReplicaAZInfo [bool](/builtin#bool)
    
    	// AZFromInfo forces the `availability_zone` field to be taken from an INFO command instead of HELLO.
    	// Primarily used for AWS MemoryDB.
    	AZFromInfo [bool](/builtin#bool)
    }

ClientOption should be passed to NewClient to construct a Client 

####  func [MustParseURL](https://github.com/redis/rueidis/blob/v1.0.71/url.go#L113) ¶ added in v1.0.17
    
    
    func MustParseURL(str [string](/builtin#string)) ClientOption

####  func [ParseURL](https://github.com/redis/rueidis/blob/v1.0.71/url.go#L21) ¶ added in v1.0.17
    
    
    func ParseURL(str [string](/builtin#string)) (opt ClientOption, err [error](/builtin#error))

ParseURL parses a redis URL into ClientOption. <https://github.com/redis/redis-specifications/blob/master/uri/redis.txt> Example: 

redis://<user>:<password>@<host>:<port>/<db_number> redis://<user>:<password>@<host>:<port>?addr=<host2>:<port2>&addr=<host3>:<port3> unix://<user>:<password>@</path/to/redis.sock>?db=<db_number>

####  type [ClusterOption](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L307) ¶ added in v1.0.47
    
    
    type ClusterOption struct {
    	// ShardsRefreshInterval is the interval to scan the cluster topology.
    	// If the value is zero, refreshment will be disabled.
    	// Cluster topology cache refresh happens always in the background after a successful scan.
    	ShardsRefreshInterval [time](/time).[Duration](/time#Duration)
    
    	// MaxMovedRedirections is the maximum number of times to retry a command when receiving MOVED|ASK responses.
    	// If set to 0 (default), MOVED|ASK retries will continue until the context timeout.
    	// If set to a positive value, the client will return an error after that many MOVED|ASK redirects.
    	// This helps prevent infinite redirect loops in case of cluster misconfiguration.
    	MaxMovedRedirections [int](/builtin#int)
    }

ClusterOption is the options for the redis cluster client. 

####  type [Commands](https://github.com/redis/rueidis/blob/v1.0.71/cmds.go#L34) ¶
    
    
    type Commands []Completed

Commands is an exported alias to []Completed. This allows users to store commands for later usage, for example: 
    
    
    c, release := client.Dedicate()
    defer release()
    
    cmds := make(rueidis.Commands, 0, 10)
    for i := 0; i < 10; i++ {
        cmds = append(cmds, c.B().Set().Key(strconv.Itoa(i)).Value(strconv.Itoa(i)).Build())
    }
    for _, resp := range c.DoMulti(ctx, cmds...) {
        if err := resp.Error(); err != nil {
        panic(err)
    }
    

However, please know that once commands are processed by the Do() or DoMulti(), they are recycled and should not be reused. 

####  type [Completed](https://github.com/redis/rueidis/blob/v1.0.71/cmds.go#L12) ¶
    
    
    type Completed = [cmds](/github.com/redis/rueidis@v1.0.71/internal/cmds).[Completed](/github.com/redis/rueidis@v1.0.71/internal/cmds#Completed)

Completed represents a completed Redis command. It should only be created from the Build() of a command builder. 

####  type [CoreClient](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L419) ¶ added in v1.0.14
    
    
    type CoreClient interface {
    	// B is the getter function to the command builder for the client
    	// If the client is a cluster client, the command builder also prohibits cross-key slots in one command.
    	B() Builder
    	// Do is the method sending user's redis command building from the B() to a redis node.
    	//  client.Do(ctx, client.B().Get().Key("k").Build()).ToString()
    	// All concurrent non-blocking commands will be pipelined automatically and have better throughput.
    	// Blocking commands will use another separated connection pool.
    	// The cmd parameter is recycled after passing into Do() and should not be reused.
    	Do(ctx [context](/context).[Context](/context#Context), cmd Completed) (resp RedisResult)
    	// DoMulti takes multiple redis commands and sends them together, reducing RTT from the user code.
    	// The multi parameters are recycled after passing into DoMulti() and should not be reused.
    	DoMulti(ctx [context](/context).[Context](/context#Context), multi ...Completed) (resp []RedisResult)
    	// Receive accepts SUBSCRIBE, SSUBSCRIBE, PSUBSCRIBE command and a message handler.
    	// Receive will block and then return value only when the following cases:
    	//   1. return nil when received any unsubscribe/punsubscribe message related to the provided `subscribe` command.
    	//   2. return ErrClosing when the client is closed manually.
    	//   3. return ctx.Err() when the `ctx` is done.
    	//   4. return non-nil err when the provided `subscribe` command failed.
    	Receive(ctx [context](/context).[Context](/context#Context), subscribe Completed, fn func(msg PubSubMessage)) [error](/builtin#error)
    	// Close will make further calls to the client be rejected with ErrClosing,
    	// and Close will wait until all pending calls finished.
    	Close()
    }

CoreClient is the minimum interface shared by the Client and the DedicatedClient. 

####  type [DedicatedClient](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L403) ¶
    
    
    type DedicatedClient interface {
    	CoreClient
    
    	// SetPubSubHooks is an alternative way to processing Pub/Sub messages instead of using Receive.
    	// SetPubSubHooks is non-blocking and allows users to subscribe/unsubscribe channels later.
    	// Note that the hooks will be called sequentially but in another goroutine.
    	// The return value will be either:
    	//   1. an error channel, if the hooks passed in are not zero, or
    	//   2. nil, if the hooks passed in are zero. (used for reset hooks)
    	// In the former case, the error channel is guaranteed to be close when the hooks will not be called anymore
    	// and has at most one error describing the reason why the hooks will not be called anymore.
    	// Users can use the error channel to detect disconnection.
    	SetPubSubHooks(hooks PubSubHooks) <-chan [error](/builtin#error)
    }

DedicatedClient is obtained from Client.Dedicated() and it will be bound to a single redis connection, and no other commands can be pipelined into this connection during Client.Dedicated(). If the DedicatedClient is obtained from a cluster client, the first command to it must have a Key() to identify the redis node. 

####  type [FtSearchDoc](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1308) ¶
    
    
    type FtSearchDoc struct {
    	Doc   map[[string](/builtin#string)][string](/builtin#string)
    	Key   [string](/builtin#string)
    	Score [float64](/builtin#float64)
    }

####  type [GeoLocation](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1437) ¶ added in v1.0.8
    
    
    type GeoLocation struct {
    	Name                      [string](/builtin#string)
    	Longitude, Latitude, Dist [float64](/builtin#float64)
    	GeoHash                   [int64](/builtin#int64)
    }

####  type [Incomplete](https://github.com/redis/rueidis/blob/v1.0.71/cmds.go#L9) ¶ added in v1.0.18
    
    
    type Incomplete = [cmds](/github.com/redis/rueidis@v1.0.71/internal/cmds).[Incomplete](/github.com/redis/rueidis@v1.0.71/internal/cmds#Incomplete)

Incomplete represents an incomplete Redis command. It should then be completed by calling Build(). 

####  type [KeyValues](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1272) ¶
    
    
    type KeyValues struct {
    	Key    [string](/builtin#string)
    	Values [][string](/builtin#string)
    }

####  type [KeyZScores](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1290) ¶
    
    
    type KeyZScores struct {
    	Key    [string](/builtin#string)
    	Values []ZScore
    }

####  type [Lua](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L86) ¶
    
    
    type Lua struct {
    	// contains filtered or unexported fields
    }

Lua represents a redis lua script. It should be created from the NewLuaScript() or NewLuaScriptReadOnly(). 

Example (Exec) ¶
    
    
    client, err := NewClient(ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
    	panic(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    
    script := NewLuaScript("return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}")
    
    _, _ = script.Exec(ctx, client, []string{"k1", "k2"}, []string{"a1", "a2"}).ToArray()
    

####  func [NewLuaScript](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L29) ¶
    
    
    func NewLuaScript(script [string](/builtin#string), opts ...LuaOption) *Lua

NewLuaScript creates a Lua instance whose Lua.Exec uses EVALSHA and EVAL. By default, SHA-1 is calculated client-side. Use WithLoadSHA1(true) option to load SHA-1 from Redis instead. 

####  func [NewLuaScriptNoSha](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L42) ¶ added in v1.0.60
    
    
    func NewLuaScriptNoSha(script [string](/builtin#string)) *Lua

NewLuaScriptNoSha creates a Lua instance whose Lua.Exec uses EVAL only (never EVALSHA). No SHA-1 is calculated or loaded. The script is sent to the server every time. Use this when you want to avoid SHA-1 entirely (e.g., to fully avoid hash collision concerns). 

####  func [NewLuaScriptNoShaRetryable](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L62) ¶ added in v1.0.70
    
    
    func NewLuaScriptNoShaRetryable(script [string](/builtin#string)) *Lua

NewLuaScriptNoShaRetryable creates a retryable Lua instance whose Lua.Exec uses EVAL only (never EVALSHA). No SHA-1 is calculated or loaded. The script is sent to the server every time. Use this when you want to avoid SHA-1 entirely (e.g., to fully avoid hash collision concerns). 

####  func [NewLuaScriptReadOnly](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L35) ¶
    
    
    func NewLuaScriptReadOnly(script [string](/builtin#string), opts ...LuaOption) *Lua

NewLuaScriptReadOnly creates a Lua instance whose Lua.Exec uses EVALSHA_RO and EVAL_RO. By default, SHA-1 is calculated client-side. Use WithLoadSHA1(true) option to load SHA-1 from Redis instead. 

####  func [NewLuaScriptReadOnlyNoSha](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L49) ¶ added in v1.0.60
    
    
    func NewLuaScriptReadOnlyNoSha(script [string](/builtin#string)) *Lua

NewLuaScriptReadOnlyNoSha creates a Lua instance whose Lua.Exec uses EVAL_RO only (never EVALSHA_RO). No SHA-1 is calculated or loaded. The script is sent to the server every time. Use this when you want to avoid SHA-1 entirely (e.g., to fully avoid hash collision concerns). 

####  func [NewLuaScriptRetryable](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L55) ¶ added in v1.0.70
    
    
    func NewLuaScriptRetryable(script [string](/builtin#string), opts ...LuaOption) *Lua

NewLuaScriptRetryable creates a retryable Lua instance whose Lua.Exec uses EVALSHA and EVAL. By default, SHA-1 is calculated client-side. Use WithLoadSHA1(true) option to load SHA-1 from Redis instead. 

####  func (*Lua) [Exec](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L103) ¶
    
    
    func (s *Lua) Exec(ctx [context](/context).[Context](/context#Context), c Client, keys, args [][string](/builtin#string)) (resp RedisResult)

Exec the script to the given Client. It will first try with the EVALSHA/EVALSHA_RO and then EVAL/EVAL_RO if the first try failed. If Lua is initialized with disabled SHA1, it will use EVAL/EVAL_RO without the EVALSHA/EVALSHA_RO attempt. If Lua is initialized with SHA-1 loading, it will call SCRIPT LOAD once to obtain the SHA-1 from Redis. Cross-slot keys are prohibited if the Client is a cluster client. 

####  func (*Lua) [ExecMulti](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L169) ¶
    
    
    func (s *Lua) ExecMulti(ctx [context](/context).[Context](/context#Context), c Client, multi ...LuaExec) (resp []RedisResult)

ExecMulti exec the script multiple times by the provided LuaExec to the given Client. For regular constructors, it will SCRIPT LOAD to all redis nodes and then use EVALSHA/EVALSHA_RO. For NoSha constructors, it will use EVAL/EVAL_RO only without any script loading. Cross-slot keys within the single LuaExec are prohibited if the Client is a cluster client. 

####  type [LuaExec](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L160) ¶
    
    
    type LuaExec struct {
    	Keys [][string](/builtin#string)
    	Args [][string](/builtin#string)
    }

LuaExec is a single execution unit of Lua.ExecMulti. 

####  type [LuaOption](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L15) ¶ added in v1.0.68
    
    
    type LuaOption func(*Lua)

LuaOption is a functional option for configuring Lua script behavior. 

####  func [WithLoadSHA1](https://github.com/redis/rueidis/blob/v1.0.71/lua.go#L21) ¶ added in v1.0.68
    
    
    func WithLoadSHA1(enabled [bool](/builtin#bool)) LuaOption

WithLoadSHA1 allows enabling loading of SHA-1 from Redis via SCRIPT LOAD instead of calculating it on the client side. When enabled, the SHA-1 hash is not calculated client-side (important for FIPS compliance). Instead, on first execution, SCRIPT LOAD is called to obtain the SHA-1 from Redis, which is then used for EVALSHA commands in subsequent executions. 

####  type [MultiRedisResultStream](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L1142) ¶ added in v1.0.29
    
    
    type MultiRedisResultStream = RedisResultStream

####  type [NewCacheStoreFn](https://github.com/redis/rueidis/blob/v1.0.71/cache.go#L10) ¶
    
    
    type NewCacheStoreFn func(CacheStoreOption) CacheStore

NewCacheStoreFn can be provided in ClientOption for using a custom CacheStore implementation 

####  type [NodeInfo](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L333) ¶ added in v1.0.65
    
    
    type NodeInfo struct {
    	Addr [string](/builtin#string)
    	AZ   [string](/builtin#string)
    	// contains filtered or unexported fields
    }

NodeInfo is the information of a replica node in a redis cluster. 

####  type [PubSubHooks](https://github.com/redis/rueidis/blob/v1.0.71/pubsub.go#L29) ¶
    
    
    type PubSubHooks struct {
    	// OnMessage will be called when receiving "message" and "pmessage" event.
    	OnMessage func(m PubSubMessage)
    	// OnSubscription will be called when receiving "subscribe", "unsubscribe", "psubscribe" and "punsubscribe" event.
    	OnSubscription func(s PubSubSubscription)
    }

PubSubHooks can be registered into DedicatedClient to process pubsub messages without using Client.Receive 

####  type [PubSubMessage](https://github.com/redis/rueidis/blob/v1.0.71/pubsub.go#L9) ¶
    
    
    type PubSubMessage struct {
    	// Pattern is only available with pmessage.
    	Pattern [string](/builtin#string)
    	// Channel is the channel the message belongs to
    	Channel [string](/builtin#string)
    	// Message is the message content
    	Message [string](/builtin#string)
    }

PubSubMessage represents a pubsub message from redis 

####  type [PubSubSubscription](https://github.com/redis/rueidis/blob/v1.0.71/pubsub.go#L19) ¶
    
    
    type PubSubSubscription struct {
    	// Kind is "subscribe", "unsubscribe", "ssubscribe", "sunsubscribe", "psubscribe" or "punsubscribe"
    	Kind [string](/builtin#string)
    	// Channel is the event subject.
    	Channel [string](/builtin#string)
    	// Count is the current number of subscriptions for a connection.
    	Count [int64](/builtin#int64)
    }

PubSubSubscription represent a pubsub "subscribe", "unsubscribe", "ssubscribe", "sunsubscribe", "psubscribe" or "punsubscribe" event. 

####  type [ReadNodeSelectorFunc](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L97) ¶ added in v1.0.71
    
    
    type ReadNodeSelectorFunc func(slot [uint16](/builtin#uint16), nodes []NodeInfo) [int](/builtin#int)

Define distinct types for safety. 

####  func [AZAffinityNodeSelector](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L411) ¶ added in v1.0.71
    
    
    func AZAffinityNodeSelector(clientAZ [string](/builtin#string)) ReadNodeSelectorFunc

AZAffinityNodeSelector prioritizes replicas in the same AZ using Round-Robin. 

####  func [AZAffinityReplicasAndPrimaryNodeSelector](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L420) ¶ added in v1.0.71
    
    
    func AZAffinityReplicasAndPrimaryNodeSelector(clientAZ [string](/builtin#string)) ReadNodeSelectorFunc

AZAffinityReplicasAndPrimaryNodeSelector prioritizes: 1\. Same-AZ Replicas 2\. Same-AZ Primary 3\. Any Replica 4\. Primary 

####  func [PreferReplicaNodeSelector](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L398) ¶ added in v1.0.71
    
    
    func PreferReplicaNodeSelector() ReadNodeSelectorFunc

PreferReplicaNodeSelector prioritizes reading from any replica using Round-Robin. If no replicas are available, it falls back to the primary. 

####  type [RedirectMode](https://github.com/redis/rueidis/blob/v1.0.71/cluster.go#L1628) ¶
    
    
    type RedirectMode [int](/builtin#int)
    
    
    const (
    	RedirectNone RedirectMode = [iota](/builtin#iota)
    	RedirectMove
    	RedirectAsk
    	RedirectRetry
    )

####  type [RedisError](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L53) ¶
    
    
    type RedisError RedisMessage

RedisError is an error response or a nil message from the redis instance 

####  func [IsRedisErr](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L47) ¶ added in v1.0.3
    
    
    func IsRedisErr(err [error](/builtin#error)) (ret *RedisError, ok [bool](/builtin#bool))

IsRedisErr is a handy method to check if the error is a redis ERR response. 

####  func (*RedisError) [Error](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L63) ¶
    
    
    func (r *RedisError) Error() [string](/builtin#string)

####  func (*RedisError) [IsAsk](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L84) ¶
    
    
    func (r *RedisError) IsAsk() (addr [string](/builtin#string), ok [bool](/builtin#bool))

IsAsk checks if it is a redis ASK message and returns ask address. 

####  func (*RedisError) [IsBusyGroup](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L129) ¶ added in v1.0.32
    
    
    func (r *RedisError) IsBusyGroup() [bool](/builtin#bool)

IsBusyGroup checks if it is a redis BUSYGROUP message. 

####  func (*RedisError) [IsClusterDown](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L119) ¶
    
    
    func (r *RedisError) IsClusterDown() [bool](/builtin#bool)

IsClusterDown checks if it is a redis CLUSTERDOWN message and returns ask address. 

####  func (*RedisError) [IsLoading](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L114) ¶ added in v1.0.49
    
    
    func (r *RedisError) IsLoading() [bool](/builtin#bool)

IsLoading checks if it is a redis LOADING message 

####  func (*RedisError) [IsMoved](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L76) ¶
    
    
    func (r *RedisError) IsMoved() (addr [string](/builtin#string), ok [bool](/builtin#bool))

IsMoved checks if it is a redis MOVED message and returns the moved address. 

####  func (*RedisError) [IsNil](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L71) ¶
    
    
    func (r *RedisError) IsNil() [bool](/builtin#bool)

IsNil checks if it is a redis nil message. 

####  func (*RedisError) [IsNoScript](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L124) ¶
    
    
    func (r *RedisError) IsNoScript() [bool](/builtin#bool)

IsNoScript checks if it is a redis NOSCRIPT message. 

####  func (*RedisError) [IsRedirect](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L92) ¶ added in v1.0.67
    
    
    func (r *RedisError) IsRedirect() (addr [string](/builtin#string), ok [bool](/builtin#bool))

IsRedirect checks if it is a redis REDIRECT message and returns redirect address. 

####  func (*RedisError) [IsTryAgain](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L109) ¶
    
    
    func (r *RedisError) IsTryAgain() [bool](/builtin#bool)

IsTryAgain checks if it is a redis TRYAGAIN message and returns ask address. 

####  type [RedisMessage](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L571) ¶
    
    
    type RedisMessage struct {
    	// contains filtered or unexported fields
    }

RedisMessage is a redis response message, it may be a nil response 

####  func (*RedisMessage) [AsBool](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L817) ¶
    
    
    func (m *RedisMessage) AsBool() (val [bool](/builtin#bool), err [error](/builtin#error))

AsBool checks if the message is a non-nil response and parses it as bool 

####  func (*RedisMessage) [AsBoolSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L953) ¶ added in v1.0.33
    
    
    func (m *RedisMessage) AsBoolSlice() ([][bool](/builtin#bool), [error](/builtin#error))

AsBoolSlice checks if the message is a redis array/set response and converts it to []bool. Redis nil elements and other non-boolean elements will be represented as false. 

####  func (*RedisMessage) [AsBytes](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L775) ¶
    
    
    func (m *RedisMessage) AsBytes() (bs [][byte](/builtin#byte), err [error](/builtin#error))

AsBytes check if the message is a redis string response and return it as an immutable []byte 

####  func (*RedisMessage) [AsFloat64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L838) ¶
    
    
    func (m *RedisMessage) AsFloat64() (val [float64](/builtin#float64), err [error](/builtin#error))

AsFloat64 check if the message is a redis string response and parse it as float64 

####  func (*RedisMessage) [AsFloatSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L933) ¶
    
    
    func (m *RedisMessage) AsFloatSlice() ([][float64](/builtin#float64), [error](/builtin#error))

AsFloatSlice check if the message is a redis array/set response and convert to []float64. redis nil element and other non-float elements will be present as zero. 

####  func (*RedisMessage) [AsFtAggregate](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1386) ¶ added in v1.0.14
    
    
    func (m *RedisMessage) AsFtAggregate() (total [int64](/builtin#int64), docs []map[[string](/builtin#string)][string](/builtin#string), err [error](/builtin#error))

####  func (*RedisMessage) [AsFtAggregateCursor](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1427) ¶ added in v1.0.14
    
    
    func (m *RedisMessage) AsFtAggregateCursor() (cursor, total [int64](/builtin#int64), docs []map[[string](/builtin#string)][string](/builtin#string), err [error](/builtin#error))

####  func (*RedisMessage) [AsFtSearch](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1314) ¶
    
    
    func (m *RedisMessage) AsFtSearch() (total [int64](/builtin#int64), docs []FtSearchDoc, err [error](/builtin#error))

####  func (*RedisMessage) [AsGeosearch](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1443) ¶ added in v1.0.8
    
    
    func (m *RedisMessage) AsGeosearch() ([]GeoLocation, [error](/builtin#error))

####  func (*RedisMessage) [AsInt64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L793) ¶
    
    
    func (m *RedisMessage) AsInt64() (val [int64](/builtin#int64), err [error](/builtin#error))

AsInt64 check if the message is a redis string response and parse it as int64 

####  func (*RedisMessage) [AsIntMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1246) ¶
    
    
    func (m *RedisMessage) AsIntMap() (map[[string](/builtin#string)][int64](/builtin#int64), [error](/builtin#error))

AsIntMap check if the message is a redis map/array/set response and convert to map[string]int64. redis nil element and other non-integer elements will be present as zero. 

####  func (*RedisMessage) [AsIntSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L913) ¶
    
    
    func (m *RedisMessage) AsIntSlice() ([][int64](/builtin#int64), [error](/builtin#error))

AsIntSlice check if the message is a redis array/set response and convert to []int64. redis nil element and other non-integer elements will be present as zero. 

####  func (*RedisMessage) [AsLMPop](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1277) ¶
    
    
    func (m *RedisMessage) AsLMPop() (kvs KeyValues, err [error](/builtin#error))

####  func (*RedisMessage) [AsMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1214) ¶
    
    
    func (m *RedisMessage) AsMap() (map[[string](/builtin#string)]RedisMessage, [error](/builtin#error))

AsMap check if the message is a redis array/set response and convert to map[string]RedisMessage 

####  func (*RedisMessage) [AsReader](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L766) ¶
    
    
    func (m *RedisMessage) AsReader() (reader [io](/io).[Reader](/io#Reader), err [error](/builtin#error))

AsReader check if the message is a redis string response and wrap it with the strings.NewReader 

####  func (*RedisMessage) [AsScanEntry](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1198) ¶
    
    
    func (m *RedisMessage) AsScanEntry() (e ScanEntry, err [error](/builtin#error))

AsScanEntry check if the message is a redis array/set response of length 2 and convert to ScanEntry. 

####  func (*RedisMessage) [AsStrMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1227) ¶
    
    
    func (m *RedisMessage) AsStrMap() (map[[string](/builtin#string)][string](/builtin#string), [error](/builtin#error))

AsStrMap check if the message is a redis map/array/set response and convert to map[string]string. redis nil element and other non-string elements will be present as zero. 

####  func (*RedisMessage) [AsStrSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L899) ¶
    
    
    func (m *RedisMessage) AsStrSlice() ([][string](/builtin#string), [error](/builtin#error))

AsStrSlice check if the message is a redis array/set response and convert to []string. redis nil element and other non-string elements will be present as zero. 

####  func (*RedisMessage) [AsUint64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L805) ¶
    
    
    func (m *RedisMessage) AsUint64() (val [uint64](/builtin#uint64), err [error](/builtin#error))

AsUint64 check if the message is a redis string response and parse it as uint64 

####  func (*RedisMessage) [AsXRange](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L998) ¶
    
    
    func (m *RedisMessage) AsXRange() ([]XRangeEntry, [error](/builtin#error))

AsXRange check if the message is a redis array/set response and convert to []XRangeEntry 

####  func (*RedisMessage) [AsXRangeEntry](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L972) ¶
    
    
    func (m *RedisMessage) AsXRangeEntry() (XRangeEntry, [error](/builtin#error))

AsXRangeEntry check if the message is a redis array/set response of length 2 and convert to XRangeEntry 

####  func (*RedisMessage) [AsXRangeSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1056) ¶ added in v1.0.61
    
    
    func (m *RedisMessage) AsXRangeSlice() (XRangeSlice, [error](/builtin#error))

AsXRangeSlice converts a RedisMessage to XRangeSlice (preserves order and duplicates) 

####  func (*RedisMessage) [AsXRangeSlices](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1093) ¶ added in v1.0.61
    
    
    func (m *RedisMessage) AsXRangeSlices() ([]XRangeSlice, [error](/builtin#error))

AsXRangeSlices converts multiple XRange entries to slice format 

####  func (*RedisMessage) [AsXRead](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1015) ¶
    
    
    func (m *RedisMessage) AsXRead() (ret map[[string](/builtin#string)][]XRangeEntry, err [error](/builtin#error))

AsXRead converts XREAD/XREADGRUOP response to map[string][]XRangeEntry 

####  func (*RedisMessage) [AsXReadSlices](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1110) ¶ added in v1.0.61
    
    
    func (m *RedisMessage) AsXReadSlices() (map[[string](/builtin#string)][]XRangeSlice, [error](/builtin#error))

AsXReadSlices converts XREAD/XREADGROUP response to use slice format 

####  func (*RedisMessage) [AsZMPop](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1295) ¶
    
    
    func (m *RedisMessage) AsZMPop() (kvs KeyZScores, err [error](/builtin#error))

####  func (*RedisMessage) [AsZScore](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1158) ¶
    
    
    func (m *RedisMessage) AsZScore() (s ZScore, err [error](/builtin#error))

AsZScore converts ZPOPMAX and ZPOPMIN command with count 1 response to a single ZScore 

####  func (*RedisMessage) [AsZScores](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1167) ¶
    
    
    func (m *RedisMessage) AsZScores() ([]ZScore, [error](/builtin#error))

AsZScores converts ZRANGE WITHSCORES, ZDIFF WITHSCORES and ZPOPMAX/ZPOPMIN command with count > 1 responses to []ZScore 

####  func (*RedisMessage) [CacheMarshal](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L680) ¶ added in v1.0.52
    
    
    func (m *RedisMessage) CacheMarshal(buf [][byte](/builtin#byte)) [][byte](/builtin#byte)

CacheMarshal writes serialized RedisMessage to the provided buffer. If the provided buffer is nil, CacheMarshal will allocate one. Note that an output format is not compatible with different client versions. 

####  func (*RedisMessage) [CachePTTL](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1557) ¶
    
    
    func (m *RedisMessage) CachePTTL() [int64](/builtin#int64)

CachePTTL returns the remaining PTTL in seconds of client side cache 

####  func (*RedisMessage) [CachePXAT](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1569) ¶
    
    
    func (m *RedisMessage) CachePXAT() [int64](/builtin#int64)

CachePXAT returns the remaining PXAT in seconds of client side cache 

####  func (*RedisMessage) [CacheSize](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L673) ¶ added in v1.0.52
    
    
    func (m *RedisMessage) CacheSize() [int](/builtin#int)

CacheSize returns the buffer size needed by the CacheMarshal. 

####  func (*RedisMessage) [CacheTTL](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1545) ¶
    
    
    func (m *RedisMessage) CacheTTL() (ttl [int64](/builtin#int64))

CacheTTL returns the remaining TTL in seconds of client side cache 

####  func (*RedisMessage) [CacheUnmarshalView](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L692) ¶ added in v1.0.52
    
    
    func (m *RedisMessage) CacheUnmarshalView(buf [][byte](/builtin#byte)) [error](/builtin#error)

CacheUnmarshalView construct the RedisMessage from the buffer produced by CacheMarshal. Note that the buffer can't be reused after CacheUnmarshalView since it uses unsafe.String on top of the buffer. 

####  func (*RedisMessage) [DecodeJSON](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L784) ¶
    
    
    func (m *RedisMessage) DecodeJSON(v [any](/builtin#any)) (err [error](/builtin#error))

DecodeJSON check if the message is a redis string response and treat it as JSON, then unmarshal it into the provided value 

####  func (*RedisMessage) [Error](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L740) ¶
    
    
    func (m *RedisMessage) Error() [error](/builtin#error)

Error check if the message is a redis error response, including nil response 

####  func (*RedisMessage) [IsArray](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L730) ¶
    
    
    func (m *RedisMessage) IsArray() [bool](/builtin#bool)

IsArray check if the message is a redis array response 

####  func (*RedisMessage) [IsBool](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L725) ¶
    
    
    func (m *RedisMessage) IsBool() [bool](/builtin#bool)

IsBool check if the message is a redis RESP3 bool response 

####  func (*RedisMessage) [IsCacheHit](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1540) ¶
    
    
    func (m *RedisMessage) IsCacheHit() [bool](/builtin#bool)

IsCacheHit check if the message is from the client side cache 

####  func (*RedisMessage) [IsFloat64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L715) ¶
    
    
    func (m *RedisMessage) IsFloat64() [bool](/builtin#bool)

IsFloat64 check if the message is a redis RESP3 double response 

####  func (*RedisMessage) [IsInt64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L710) ¶
    
    
    func (m *RedisMessage) IsInt64() [bool](/builtin#bool)

IsInt64 check if the message is a redis RESP3 int response 

####  func (*RedisMessage) [IsMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L735) ¶
    
    
    func (m *RedisMessage) IsMap() [bool](/builtin#bool)

IsMap check if the message is a redis RESP3 map response 

####  func (*RedisMessage) [IsNil](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L705) ¶
    
    
    func (m *RedisMessage) IsNil() [bool](/builtin#bool)

IsNil check if the message is a redis nil response 

####  func (*RedisMessage) [IsString](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L720) ¶
    
    
    func (m *RedisMessage) IsString() [bool](/builtin#bool)

IsString check if the message is a redis string response 

####  func (*RedisMessage) [String](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1619) ¶ added in v1.0.17
    
    
    func (m *RedisMessage) String() [string](/builtin#string)

String returns the human-readable representation of RedisMessage 

####  func (*RedisMessage) [ToAny](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1501) ¶
    
    
    func (m *RedisMessage) ToAny() ([any](/builtin#any), [error](/builtin#error))

ToAny turns the message into go any value 

####  func (*RedisMessage) [ToArray](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L886) ¶
    
    
    func (m *RedisMessage) ToArray() ([]RedisMessage, [error](/builtin#error))

ToArray check if the message is a redis array/set response and return it 

####  func (*RedisMessage) [ToBool](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L862) ¶
    
    
    func (m *RedisMessage) ToBool() (val [bool](/builtin#bool), err [error](/builtin#error))

ToBool check if the message is a redis RESP3 bool response and return it 

####  func (*RedisMessage) [ToFloat64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L874) ¶
    
    
    func (m *RedisMessage) ToFloat64() (val [float64](/builtin#float64), err [error](/builtin#error))

ToFloat64 check if the message is a redis RESP3 double response and return it 

####  func (*RedisMessage) [ToInt64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L850) ¶
    
    
    func (m *RedisMessage) ToInt64() (val [int64](/builtin#int64), err [error](/builtin#error))

ToInt64 check if the message is a redis RESP3 int response and return it 

####  func (*RedisMessage) [ToMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1489) ¶
    
    
    func (m *RedisMessage) ToMap() (map[[string](/builtin#string)]RedisMessage, [error](/builtin#error))

ToMap check if the message is a redis RESP3 map response and return it 

####  func (*RedisMessage) [ToString](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L754) ¶
    
    
    func (m *RedisMessage) ToString() (val [string](/builtin#string), err [error](/builtin#error))

ToString check if the message is a redis string response and return it 

####  type [RedisResult](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L143) ¶
    
    
    type RedisResult struct {
    	// contains filtered or unexported fields
    }

RedisResult is the return struct from Client.Do or Client.DoCache it contains either a redis response or an underlying error (ex. network timeout). 

####  func (RedisResult) [AsBool](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L264) ¶
    
    
    func (r RedisResult) AsBool() (v [bool](/builtin#bool), err [error](/builtin#error))

AsBool delegates to RedisMessage.AsBool 

####  func (RedisResult) [AsBoolSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L324) ¶ added in v1.0.33
    
    
    func (r RedisResult) AsBoolSlice() (v [][bool](/builtin#bool), err [error](/builtin#error))

AsBoolSlice delegates to RedisMessage.AsBoolSlice 

####  func (RedisResult) [AsBytes](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L224) ¶
    
    
    func (r RedisResult) AsBytes() (v [][byte](/builtin#byte), err [error](/builtin#error))

AsBytes delegates to RedisMessage.AsBytes 

####  func (RedisResult) [AsFloat64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L274) ¶
    
    
    func (r RedisResult) AsFloat64() (v [float64](/builtin#float64), err [error](/builtin#error))

AsFloat64 delegates to RedisMessage.AsFloat64 

####  func (RedisResult) [AsFloatSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L314) ¶
    
    
    func (r RedisResult) AsFloatSlice() (v [][float64](/builtin#float64), err [error](/builtin#error))

AsFloatSlice delegates to RedisMessage.AsFloatSlice 

####  func (RedisResult) [AsFtAggregate](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L440) ¶ added in v1.0.14
    
    
    func (r RedisResult) AsFtAggregate() (total [int64](/builtin#int64), docs []map[[string](/builtin#string)][string](/builtin#string), err [error](/builtin#error))

####  func (RedisResult) [AsFtAggregateCursor](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L449) ¶ added in v1.0.14
    
    
    func (r RedisResult) AsFtAggregateCursor() (cursor, total [int64](/builtin#int64), docs []map[[string](/builtin#string)][string](/builtin#string), err [error](/builtin#error))

####  func (RedisResult) [AsFtSearch](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L431) ¶
    
    
    func (r RedisResult) AsFtSearch() (total [int64](/builtin#int64), docs []FtSearchDoc, err [error](/builtin#error))

####  func (RedisResult) [AsGeosearch](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L458) ¶ added in v1.0.8
    
    
    func (r RedisResult) AsGeosearch() (locations []GeoLocation, err [error](/builtin#error))

####  func (RedisResult) [AsInt64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L244) ¶
    
    
    func (r RedisResult) AsInt64() (v [int64](/builtin#int64), err [error](/builtin#error))

AsInt64 delegates to RedisMessage.AsInt64 

####  func (RedisResult) [AsIntMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L488) ¶
    
    
    func (r RedisResult) AsIntMap() (v map[[string](/builtin#string)][int64](/builtin#int64), err [error](/builtin#error))

AsIntMap delegates to RedisMessage.AsIntMap 

####  func (RedisResult) [AsIntSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L304) ¶
    
    
    func (r RedisResult) AsIntSlice() (v [][int64](/builtin#int64), err [error](/builtin#error))

AsIntSlice delegates to RedisMessage.AsIntSlice 

####  func (RedisResult) [AsLMPop](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L413) ¶
    
    
    func (r RedisResult) AsLMPop() (v KeyValues, err [error](/builtin#error))

####  func (RedisResult) [AsMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L468) ¶
    
    
    func (r RedisResult) AsMap() (v map[[string](/builtin#string)]RedisMessage, err [error](/builtin#error))

AsMap delegates to RedisMessage.AsMap 

####  func (RedisResult) [AsReader](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L214) ¶
    
    
    func (r RedisResult) AsReader() (v [io](/io).[Reader](/io#Reader), err [error](/builtin#error))

AsReader delegates to RedisMessage.AsReader 

####  func (RedisResult) [AsScanEntry](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L498) ¶
    
    
    func (r RedisResult) AsScanEntry() (v ScanEntry, err [error](/builtin#error))

AsScanEntry delegates to RedisMessage.AsScanEntry. 

####  func (RedisResult) [AsStrMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L478) ¶
    
    
    func (r RedisResult) AsStrMap() (v map[[string](/builtin#string)][string](/builtin#string), err [error](/builtin#error))

AsStrMap delegates to RedisMessage.AsStrMap 

####  func (RedisResult) [AsStrSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L294) ¶
    
    
    func (r RedisResult) AsStrSlice() (v [][string](/builtin#string), err [error](/builtin#error))

AsStrSlice delegates to RedisMessage.AsStrSlice 

####  func (RedisResult) [AsUint64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L254) ¶
    
    
    func (r RedisResult) AsUint64() (v [uint64](/builtin#uint64), err [error](/builtin#error))

AsUint64 delegates to RedisMessage.AsUint64 

####  func (RedisResult) [AsXRange](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L344) ¶
    
    
    func (r RedisResult) AsXRange() (v []XRangeEntry, err [error](/builtin#error))

AsXRange delegates to RedisMessage.AsXRange 

####  func (RedisResult) [AsXRangeEntry](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L334) ¶
    
    
    func (r RedisResult) AsXRangeEntry() (v XRangeEntry, err [error](/builtin#error))

AsXRangeEntry delegates to RedisMessage.AsXRangeEntry 

####  func (RedisResult) [AsXRangeSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L384) ¶ added in v1.0.61
    
    
    func (r RedisResult) AsXRangeSlice() (v XRangeSlice, err [error](/builtin#error))

AsXRangeSlice delegates to RedisMessage.AsXRangeSlice 

####  func (RedisResult) [AsXRangeSlices](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L394) ¶ added in v1.0.61
    
    
    func (r RedisResult) AsXRangeSlices() (v []XRangeSlice, err [error](/builtin#error))

AsXRangeSlices delegates to RedisMessage.AsXRangeSlices 

####  func (RedisResult) [AsXRead](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L374) ¶
    
    
    func (r RedisResult) AsXRead() (v map[[string](/builtin#string)][]XRangeEntry, err [error](/builtin#error))

AsXRead delegates to RedisMessage.AsXRead 

####  func (RedisResult) [AsXReadSlices](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L404) ¶ added in v1.0.61
    
    
    func (r RedisResult) AsXReadSlices() (v map[[string](/builtin#string)][]XRangeSlice, err [error](/builtin#error))

AsXReadSlices delegates to RedisMessage.AsXReadSlices 

####  func (RedisResult) [AsZMPop](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L422) ¶
    
    
    func (r RedisResult) AsZMPop() (v KeyZScores, err [error](/builtin#error))

####  func (RedisResult) [AsZScore](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L354) ¶
    
    
    func (r RedisResult) AsZScore() (v ZScore, err [error](/builtin#error))

AsZScore delegates to RedisMessage.AsZScore 

####  func (RedisResult) [AsZScores](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L364) ¶
    
    
    func (r RedisResult) AsZScores() (v []ZScore, err [error](/builtin#error))

AsZScores delegates to RedisMessage.AsZScores 

####  func (RedisResult) [CachePTTL](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L538) ¶
    
    
    func (r RedisResult) CachePTTL() [int64](/builtin#int64)

CachePTTL delegates to RedisMessage.CachePTTL 

####  func (RedisResult) [CachePXAT](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L543) ¶
    
    
    func (r RedisResult) CachePXAT() [int64](/builtin#int64)

CachePXAT delegates to RedisMessage.CachePXAT 

####  func (RedisResult) [CacheTTL](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L533) ¶
    
    
    func (r RedisResult) CacheTTL() [int64](/builtin#int64)

CacheTTL delegates to RedisMessage.CacheTTL 

####  func (RedisResult) [DecodeJSON](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L234) ¶
    
    
    func (r RedisResult) DecodeJSON(v [any](/builtin#any)) (err [error](/builtin#error))

DecodeJSON delegates to RedisMessage.DecodeJSON 

####  func (RedisResult) [Error](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L154) ¶
    
    
    func (r RedisResult) Error() (err [error](/builtin#error))

Error returns either underlying error or redis error or nil 

####  func (RedisResult) [IsCacheHit](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L528) ¶
    
    
    func (r RedisResult) IsCacheHit() [bool](/builtin#bool)

IsCacheHit delegates to RedisMessage.IsCacheHit 

####  func (RedisResult) [NonRedisError](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L149) ¶
    
    
    func (r RedisResult) NonRedisError() [error](/builtin#error)

NonRedisError can be used to check if there is an underlying error (ex. network timeout). 

####  func (*RedisResult) [String](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L548) ¶ added in v1.0.17
    
    
    func (r *RedisResult) String() [string](/builtin#string)

String returns human-readable representation of RedisResult 

####  func (RedisResult) [ToAny](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L518) ¶
    
    
    func (r RedisResult) ToAny() (v [any](/builtin#any), err [error](/builtin#error))

ToAny delegates to RedisMessage.ToAny 

####  func (RedisResult) [ToArray](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L284) ¶
    
    
    func (r RedisResult) ToArray() (v []RedisMessage, err [error](/builtin#error))

ToArray delegates to RedisMessage.ToArray 

####  func (RedisResult) [ToBool](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L184) ¶
    
    
    func (r RedisResult) ToBool() (v [bool](/builtin#bool), err [error](/builtin#error))

ToBool delegates to RedisMessage.ToBool 

####  func (RedisResult) [ToFloat64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L194) ¶
    
    
    func (r RedisResult) ToFloat64() (v [float64](/builtin#float64), err [error](/builtin#error))

ToFloat64 delegates to RedisMessage.ToFloat64 

####  func (RedisResult) [ToInt64](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L174) ¶
    
    
    func (r RedisResult) ToInt64() (v [int64](/builtin#int64), err [error](/builtin#error))

ToInt64 delegates to RedisMessage.ToInt64 

####  func (RedisResult) [ToMap](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L508) ¶
    
    
    func (r RedisResult) ToMap() (v map[[string](/builtin#string)]RedisMessage, err [error](/builtin#error))

ToMap delegates to RedisMessage.ToMap 

####  func (RedisResult) [ToMessage](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L164) ¶
    
    
    func (r RedisResult) ToMessage() (v RedisMessage, err [error](/builtin#error))

ToMessage retrieves the RedisMessage 

####  func (RedisResult) [ToString](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L204) ¶
    
    
    func (r RedisResult) ToString() (v [string](/builtin#string), err [error](/builtin#error))

ToString delegates to RedisMessage.ToString 

####  type [RedisResultStream](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L1144) ¶ added in v1.0.29
    
    
    type RedisResultStream struct {
    	// contains filtered or unexported fields
    }

####  func (*RedisResultStream) [Error](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L1158) ¶ added in v1.0.29
    
    
    func (s *RedisResultStream) Error() [error](/builtin#error)

Error returns the error happened when sending commands to redis or reading response from redis. Usually a user is not required to use this function because the error is also reported by the WriteTo. 

####  func (*RedisResultStream) [HasNext](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L1152) ¶ added in v1.0.29
    
    
    func (s *RedisResultStream) HasNext() [bool](/builtin#bool)

HasNext can be used in a for loop condition to check if a further WriteTo call is needed. 

####  func (*RedisResultStream) [WriteTo](https://github.com/redis/rueidis/blob/v1.0.71/pipe.go#L1165) ¶ added in v1.0.29
    
    
    func (s *RedisResultStream) WriteTo(w [io](/io).[Writer](/io#Writer)) (n [int64](/builtin#int64), err [error](/builtin#error))

WriteTo reads a redis response from redis and then write it to the given writer. This function is not thread-safe and should be called sequentially to read multiple responses. An io.EOF error will be reported if all responses are read. 

####  type [ReplicaInfo](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L340) ¶ added in v1.0.52
    
    
    type ReplicaInfo = NodeInfo

ReplicaInfo is the information of a replica node in a redis cluster. 

####  type [ReplicaSelectorFunc](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L98) ¶ added in v1.0.71
    
    
    type ReplicaSelectorFunc func(slot [uint16](/builtin#uint16), replicas []NodeInfo) [int](/builtin#int)

####  type [RetryDelayFn](https://github.com/redis/rueidis/blob/v1.0.71/retry.go#L18) ¶ added in v1.0.48
    
    
    type RetryDelayFn func(attempts [int](/builtin#int), cmd Completed, err [error](/builtin#error)) [time](/time).[Duration](/time#Duration)

RetryDelayFn returns the delay that should be used before retrying the attempt. Will return a negative delay if the delay could not be determined or does not retry. 

####  type [ScanEntry](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1192) ¶
    
    
    type ScanEntry struct {
    	Elements [][string](/builtin#string)
    	Cursor   [uint64](/builtin#uint64)
    }

ScanEntry is the element type of both SCAN, SSCAN, HSCAN and ZSCAN command response. 

####  type [Scanner](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L350) ¶ added in v1.0.63
    
    
    type Scanner struct {
    	// contains filtered or unexported fields
    }

####  func [NewScanner](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L355) ¶ added in v1.0.63
    
    
    func NewScanner(next func(cursor [uint64](/builtin#uint64)) (ScanEntry, [error](/builtin#error))) *Scanner

####  func (*Scanner) [Err](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L392) ¶ added in v1.0.63
    
    
    func (s *Scanner) Err() [error](/builtin#error)

####  func (*Scanner) [Iter](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L368) ¶ added in v1.0.63
    
    
    func (s *Scanner) Iter() [iter](/iter).[Seq](/iter#Seq)[[string](/builtin#string)]

####  func (*Scanner) [Iter2](https://github.com/redis/rueidis/blob/v1.0.71/helper.go#L380) ¶ added in v1.0.63
    
    
    func (s *Scanner) Iter2() [iter](/iter).[Seq2](/iter#Seq2)[[string](/builtin#string), [string](/builtin#string)]

####  type [SentinelOption](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L291) ¶
    
    
    type SentinelOption struct {
    	// TCP & TLS, same as ClientOption but for connecting sentinel
    	Dialer    [net](/net).[Dialer](/net#Dialer)
    	TLSConfig *[tls](/crypto/tls).[Config](/crypto/tls#Config)
    
    	// MasterSet is the redis master set name monitored by sentinel cluster.
    	// If this field is set, then ClientOption.InitAddress will be used to connect to the sentinel cluster.
    	MasterSet [string](/builtin#string)
    
    	// Redis AUTH parameters for sentinel
    	Username   [string](/builtin#string)
    	Password   [string](/builtin#string)
    	ClientName [string](/builtin#string)
    }

SentinelOption contains MasterSet, 

####  type [SimpleCache](https://github.com/redis/rueidis/blob/v1.0.71/cache.go#L50) ¶ added in v1.0.1
    
    
    type SimpleCache interface {
    	Get(key [string](/builtin#string)) RedisMessage
    	Set(key [string](/builtin#string), val RedisMessage)
    	Del(key [string](/builtin#string))
    	Flush()
    }

SimpleCache is an alternative interface should be paired with NewSimpleCacheAdapter to construct a CacheStore 

####  type [StandaloneOption](https://github.com/redis/rueidis/blob/v1.0.71/rueidis.go#L321) ¶ added in v1.0.57
    
    
    type StandaloneOption struct {
    	// ReplicaAddress is the list of replicas for the primary node.
    	// Note that these addresses must be online and cannot be promoted.
    	// An example use case is the reader endpoint provided by cloud vendors.
    	ReplicaAddress [][string](/builtin#string)
    	// EnableRedirect enables the CLIENT CAPA redirect feature for Valkey 8+
    	// When enabled, the client will send CLIENT CAPA redirect during connection
    	// initialization and handle REDIRECT responses from the server.
    	EnableRedirect [bool](/builtin#bool)
    }

StandaloneOption is the options for the standalone client. 

####  type [XRangeEntry](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L966) ¶
    
    
    type XRangeEntry struct {
    	FieldValues map[[string](/builtin#string)][string](/builtin#string)
    	ID          [string](/builtin#string)
    }

XRangeEntry is the element type of both XRANGE and XREVRANGE command response array 

####  type [XRangeFieldValue](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1050) ¶ added in v1.0.61
    
    
    type XRangeFieldValue struct {
    	Field [string](/builtin#string)
    	Value [string](/builtin#string)
    }

####  type [XRangeSlice](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1045) ¶ added in v1.0.61
    
    
    type XRangeSlice struct {
    	ID          [string](/builtin#string)
    	FieldValues []XRangeFieldValue
    }

New slice-based structures that preserve order and duplicates 

####  type [ZScore](https://github.com/redis/rueidis/blob/v1.0.71/message.go#L1142) ¶
    
    
    type ZScore struct {
    	Member [string](/builtin#string)
    	Score  [float64](/builtin#float64)
    }

ZScore is the element type of ZRANGE WITHSCORES, ZDIFF WITHSCORES and ZPOPMAX command response 
