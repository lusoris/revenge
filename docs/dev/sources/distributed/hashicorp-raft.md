# HashiCorp Raft

> Source: https://pkg.go.dev/github.com/hashicorp/raft
> Fetched: 2026-01-30T23:55:52.609112+00:00
> Content-Hash: ca5d2dffa99690ae
> Type: html

---

Index

¶

Constants

Variables

func BootstrapCluster(conf *Config, logs LogStore, stable StableStore, snaps SnapshotStore, ...) error

func EncodeConfiguration(configuration Configuration) []byte

func HasExistingState(logs LogStore, stable StableStore, snaps SnapshotStore) (bool, error)

func MakeCluster(n int, t *testing.T, conf *Config) *cluster

func MakeClusterCustom(t *testing.T, opts *MakeClusterOpts) *cluster

func MakeClusterNoBootstrap(n int, t *testing.T, conf *Config) *cluster

func NewInmemTransport(addr ServerAddress) (ServerAddress, *InmemTransport)

func NewInmemTransportWithTimeout(addr ServerAddress, timeout time.Duration) (ServerAddress, *InmemTransport)

func RecoverCluster(conf *Config, fsm FSM, logs LogStore, stable StableStore, snaps SnapshotStore, ...) error

func ValidateConfig(config *Config) error

type AppendEntriesRequest

func (r *AppendEntriesRequest) GetRPCHeader() RPCHeader

type AppendEntriesResponse

func (r *AppendEntriesResponse) GetRPCHeader() RPCHeader

type AppendFuture

type AppendPipeline

type ApplyFuture

type BatchingFSM

type Config

func DefaultConfig() *Config

type Configuration

func DecodeConfiguration(buf []byte) Configuration

func GetConfiguration(conf *Config, fsm FSM, logs LogStore, stable StableStore, snaps SnapshotStore, ...) (Configuration, error)

func ReadConfigJSON(path string) (Configuration, error)

func ReadPeersJSON(path string) (Configuration, error)

func (c *Configuration) Clone() (copy Configuration)

type ConfigurationChangeCommand

func (c ConfigurationChangeCommand) String() string

type ConfigurationFuture

type ConfigurationStore

type CountingReader

type DiscardSnapshotSink

func (d *DiscardSnapshotSink) Cancel() error

func (d *DiscardSnapshotSink) Close() error

func (d *DiscardSnapshotSink) ID() string

func (d *DiscardSnapshotSink) Write(b []byte) (int, error)

type DiscardSnapshotStore

func NewDiscardSnapshotStore() *DiscardSnapshotStore

func (d *DiscardSnapshotStore) Create(version SnapshotVersion, index, term uint64, configuration Configuration, ...) (SnapshotSink, error)

func (d *DiscardSnapshotStore) List() ([]*SnapshotMeta, error)

func (d *DiscardSnapshotStore) Open(id string) (*SnapshotMeta, io.ReadCloser, error)

type FSM

type FSMSnapshot

type FailedHeartbeatObservation

type FileSnapshotSink

func (s *FileSnapshotSink) Cancel() error

func (s *FileSnapshotSink) Close() error

func (s *FileSnapshotSink) ID() string

func (s *FileSnapshotSink) Write(b []byte) (int, error)

type FileSnapshotStore

func FileSnapTest(t *testing.T) (string, *FileSnapshotStore)

func NewFileSnapshotStore(base string, retain int, logOutput io.Writer) (*FileSnapshotStore, error)

func NewFileSnapshotStoreWithLogger(base string, retain int, logger hclog.Logger) (*FileSnapshotStore, error)

func (f *FileSnapshotStore) Create(version SnapshotVersion, index, term uint64, configuration Configuration, ...) (SnapshotSink, error)

func (f *FileSnapshotStore) List() ([]*SnapshotMeta, error)

func (f *FileSnapshotStore) Open(id string) (*SnapshotMeta, io.ReadCloser, error)

func (f *FileSnapshotStore) ReapSnapshots() error

type FilterFn

type Future

type IndexFuture

type InmemSnapshotSink

func (s *InmemSnapshotSink) Cancel() error

func (s *InmemSnapshotSink) Close() error

func (s *InmemSnapshotSink) ID() string

func (s *InmemSnapshotSink) Write(p []byte) (n int, err error)

type InmemSnapshotStore

func NewInmemSnapshotStore() *InmemSnapshotStore

func (m *InmemSnapshotStore) Create(version SnapshotVersion, index, term uint64, configuration Configuration, ...) (SnapshotSink, error)

func (m *InmemSnapshotStore) List() ([]*SnapshotMeta, error)

func (m *InmemSnapshotStore) Open(id string) (*SnapshotMeta, io.ReadCloser, error)

type InmemStore

func NewInmemStore() *InmemStore

func (i *InmemStore) DeleteRange(min, max uint64) error

func (i *InmemStore) FirstIndex() (uint64, error)

func (i *InmemStore) Get(key []byte) ([]byte, error)

func (i *InmemStore) GetLog(index uint64, log *Log) error

func (i *InmemStore) GetUint64(key []byte) (uint64, error)

func (i *InmemStore) LastIndex() (uint64, error)

func (i *InmemStore) Set(key []byte, val []byte) error

func (i *InmemStore) SetUint64(key []byte, val uint64) error

func (i *InmemStore) StoreLog(log *Log) error

func (i *InmemStore) StoreLogs(logs []*Log) error

type InmemTransport

func (i *InmemTransport) AppendEntries(id ServerID, target ServerAddress, args *AppendEntriesRequest, ...) error

func (i *InmemTransport) AppendEntriesPipeline(id ServerID, target ServerAddress) (AppendPipeline, error)

func (i *InmemTransport) Close() error

func (i *InmemTransport) Connect(peer ServerAddress, t Transport)

func (i *InmemTransport) Consumer() <-chan RPC

func (i *InmemTransport) DecodePeer(buf []byte) ServerAddress

func (i *InmemTransport) Disconnect(peer ServerAddress)

func (i *InmemTransport) DisconnectAll()

func (i *InmemTransport) EncodePeer(id ServerID, p ServerAddress) []byte

func (i *InmemTransport) InstallSnapshot(id ServerID, target ServerAddress, args *InstallSnapshotRequest, ...) error

func (i *InmemTransport) LocalAddr() ServerAddress

func (i *InmemTransport) RequestPreVote(id ServerID, target ServerAddress, args *RequestPreVoteRequest, ...) error

func (i *InmemTransport) RequestVote(id ServerID, target ServerAddress, args *RequestVoteRequest, ...) error

func (i *InmemTransport) SetHeartbeatHandler(cb func(RPC))

func (i *InmemTransport) TimeoutNow(id ServerID, target ServerAddress, args *TimeoutNowRequest, ...) error

type InstallSnapshotRequest

func (r *InstallSnapshotRequest) GetRPCHeader() RPCHeader

type InstallSnapshotResponse

func (r *InstallSnapshotResponse) GetRPCHeader() RPCHeader

type LeaderObservation

type LeadershipTransferFuture

type Log

type LogCache

func NewLogCache(capacity int, store LogStore) (*LogCache, error)

func (c *LogCache) DeleteRange(min, max uint64) error

func (c *LogCache) FirstIndex() (uint64, error)

func (c *LogCache) GetLog(idx uint64, log *Log) error

func (c *LogCache) IsMonotonic() bool

func (c *LogCache) LastIndex() (uint64, error)

func (c *LogCache) StoreLog(log *Log) error

func (c *LogCache) StoreLogs(logs []*Log) error

type LogStore

type LogType

func (lt LogType) String() string

type LoopbackTransport

type MakeClusterOpts

type MockFSM

func (m *MockFSM) Apply(log *Log) interface{}

func (m *MockFSM) Logs() [][]byte

func (m *MockFSM) Restore(inp io.ReadCloser) error

func (m *MockFSM) Snapshot() (FSMSnapshot, error)

type MockFSMConfigStore

func (m *MockFSMConfigStore) StoreConfiguration(index uint64, config Configuration)

type MockMonotonicLogStore

func (m *MockMonotonicLogStore) DeleteRange(min uint64, max uint64) error

func (m *MockMonotonicLogStore) FirstIndex() (uint64, error)

func (m *MockMonotonicLogStore) GetLog(index uint64, log *Log) error

func (m *MockMonotonicLogStore) IsMonotonic() bool

func (m *MockMonotonicLogStore) LastIndex() (uint64, error)

func (m *MockMonotonicLogStore) StoreLog(log *Log) error

func (m *MockMonotonicLogStore) StoreLogs(logs []*Log) error

type MockSnapshot

func (m *MockSnapshot) Persist(sink SnapshotSink) error

func (m *MockSnapshot) Release()

type MonotonicLogStore

type NetworkTransport

func NewNetworkTransport(stream StreamLayer, maxPool int, timeout time.Duration, logOutput io.Writer) *NetworkTransport

func NewNetworkTransportWithConfig(config *NetworkTransportConfig) *NetworkTransport

func NewNetworkTransportWithLogger(stream StreamLayer, maxPool int, timeout time.Duration, logger hclog.Logger) *NetworkTransport

func NewTCPTransport(bindAddr string, advertise net.Addr, maxPool int, timeout time.Duration, ...) (*NetworkTransport, error)

func NewTCPTransportWithConfig(bindAddr string, advertise net.Addr, config *NetworkTransportConfig) (*NetworkTransport, error)

func NewTCPTransportWithLogger(bindAddr string, advertise net.Addr, maxPool int, timeout time.Duration, ...) (*NetworkTransport, error)

func (n *NetworkTransport) AppendEntries(id ServerID, target ServerAddress, args *AppendEntriesRequest, ...) error

func (n *NetworkTransport) AppendEntriesPipeline(id ServerID, target ServerAddress) (AppendPipeline, error)

func (n *NetworkTransport) Close() error

func (n *NetworkTransport) CloseStreams()

func (n *NetworkTransport) Consumer() <-chan RPC

func (n *NetworkTransport) DecodePeer(buf []byte) ServerAddress

func (n *NetworkTransport) EncodePeer(id ServerID, p ServerAddress) []byte

func (n *NetworkTransport) InstallSnapshot(id ServerID, target ServerAddress, args *InstallSnapshotRequest, ...) error

func (n *NetworkTransport) IsShutdown() bool

func (n *NetworkTransport) LocalAddr() ServerAddress

func (n *NetworkTransport) RequestPreVote(id ServerID, target ServerAddress, args *RequestPreVoteRequest, ...) error

func (n *NetworkTransport) RequestVote(id ServerID, target ServerAddress, args *RequestVoteRequest, ...) error

func (n *NetworkTransport) SetHeartbeatHandler(cb func(rpc RPC))

func (n *NetworkTransport) TimeoutNow(id ServerID, target ServerAddress, args *TimeoutNowRequest, ...) error

type NetworkTransportConfig

type Observation

type Observer

func NewObserver(channel chan Observation, blocking bool, filter FilterFn) *Observer

func (or *Observer) GetNumDropped() uint64

func (or *Observer) GetNumObserved() uint64

type PeerObservation

type ProtocolVersion

type RPC

func (r *RPC) Respond(resp interface{}, err error)

type RPCHeader

type RPCResponse

type Raft

func NewRaft(conf *Config, fsm FSM, logs LogStore, stable StableStore, snaps SnapshotStore, ...) (*Raft, error)

func (r *Raft) AddNonvoter(id ServerID, address ServerAddress, prevIndex uint64, timeout time.Duration) IndexFuture

func (r *Raft) AddPeer(peer ServerAddress) Future

deprecated

func (r *Raft) AddVoter(id ServerID, address ServerAddress, prevIndex uint64, timeout time.Duration) IndexFuture

func (r *Raft) AppliedIndex() uint64

func (r *Raft) Apply(cmd []byte, timeout time.Duration) ApplyFuture

func (r *Raft) ApplyLog(log Log, timeout time.Duration) ApplyFuture

func (r *Raft) Barrier(timeout time.Duration) Future

func (r *Raft) BootstrapCluster(configuration Configuration) Future

func (r *Raft) CommitIndex() uint64

func (r *Raft) CurrentTerm() uint64

func (r *Raft) DemoteVoter(id ServerID, prevIndex uint64, timeout time.Duration) IndexFuture

func (r *Raft) DeregisterObserver(or *Observer)

func (r *Raft) GetConfiguration() ConfigurationFuture

func (r *Raft) LastContact() time.Time

func (r *Raft) LastIndex() uint64

func (r *Raft) Leader() ServerAddress

func (r *Raft) LeaderCh() <-chan bool

func (r *Raft) LeaderWithID() (ServerAddress, ServerID)

func (r *Raft) LeadershipTransfer() Future

func (r *Raft) LeadershipTransferToServer(id ServerID, address ServerAddress) Future

func (r *Raft) RegisterObserver(or *Observer)

func (r *Raft) ReloadConfig(rc ReloadableConfig) error

func (r *Raft) ReloadableConfig() ReloadableConfig

func (r *Raft) RemovePeer(peer ServerAddress) Future

deprecated

func (r *Raft) RemoveServer(id ServerID, prevIndex uint64, timeout time.Duration) IndexFuture

func (r *Raft) Restore(meta *SnapshotMeta, reader io.Reader, timeout time.Duration) error

func (r *Raft) Shutdown() Future

func (r *Raft) Snapshot() SnapshotFuture

func (r *Raft) State() RaftState

func (r *Raft) Stats() map[string]string

func (r *Raft) String() string

func (r *Raft) VerifyLeader() Future

type RaftState

func (s RaftState) String() string

type ReadCloserWrapper

type ReloadableConfig

type RequestPreVoteRequest

func (r *RequestPreVoteRequest) GetRPCHeader() RPCHeader

type RequestPreVoteResponse

func (r *RequestPreVoteResponse) GetRPCHeader() RPCHeader

type RequestVoteRequest

func (r *RequestVoteRequest) GetRPCHeader() RPCHeader

type RequestVoteResponse

func (r *RequestVoteResponse) GetRPCHeader() RPCHeader

type ResumedHeartbeatObservation

type Server

type ServerAddress

func NewInmemAddr() ServerAddress

type ServerAddressProvider

type ServerID

type ServerSuffrage

func (s ServerSuffrage) String() string

type SnapshotFuture

type SnapshotMeta

type SnapshotSink

type SnapshotStore

type SnapshotVersion

type StableStore

type StreamLayer

type TCPStreamLayer

func (t *TCPStreamLayer) Accept() (c net.Conn, err error)

func (t *TCPStreamLayer) Addr() net.Addr

func (t *TCPStreamLayer) Close() (err error)

func (t *TCPStreamLayer) Dial(address ServerAddress, timeout time.Duration) (net.Conn, error)

type TimeoutNowRequest

func (r *TimeoutNowRequest) GetRPCHeader() RPCHeader

type TimeoutNowResponse

func (r *TimeoutNowResponse) GetRPCHeader() RPCHeader

type Transport

type WithClose

type WithPeers

type WithPreVote

type WithRPCHeader

type WrappingFSM

Constants

¶

View Source

const (

// DefaultTimeoutScale is the default TimeoutScale in a NetworkTransport.

DefaultTimeoutScale = 256 * 1024

// 256KB

// DefaultMaxRPCsInFlight is the default value used for pipelining configuration

// if a zero value is passed. See

https://github.com/hashicorp/raft/pull/541

// for rationale. Note, if this is changed we should update the doc comments

// below for NetworkTransportConfig.MaxRPCsInFlight.

DefaultMaxRPCsInFlight = 2
)

View Source

const (

// SuggestedMaxDataSize of the data in a raft log entry, in bytes.

//

// The value is based on current architecture, default timing, etc. Clients can

// ignore this value if they want as there is no actual hard checking

// within the library. As the library is enhanced this value may change

// over time to reflect current suggested maximums.

//

// Applying log entries with data greater than this size risks RPC IO taking

// too long and preventing timely heartbeat signals.  These signals are sent in serial

// in current transports, potentially causing leadership instability.

SuggestedMaxDataSize = 512 * 1024
)

Variables

¶

View Source

var (

// ErrLeader is returned when an operation can't be completed on a

// leader node.

ErrLeader =

errors

.

New

("node is the leader")

// ErrNotLeader is returned when an operation can't be completed on a

// follower or candidate node.

ErrNotLeader =

errors

.

New

("node is not the leader")

// ErrNotVoter is returned when an operation can't be completed on a

// non-voter node.

ErrNotVoter =

errors

.

New

("node is not a voter")

// ErrLeadershipLost is returned when a leader fails to commit a log entry

// because it's been deposed in the process.

ErrLeadershipLost =

errors

.

New

("leadership lost while committing log")

// ErrAbortedByRestore is returned when a leader fails to commit a log

// entry because it's been superseded by a user snapshot restore.

ErrAbortedByRestore =

errors

.

New

("snapshot restored while committing log")

// ErrRaftShutdown is returned when operations are requested against an

// inactive Raft.

ErrRaftShutdown =

errors

.

New

("raft is already shutdown")

// ErrEnqueueTimeout is returned when a command fails due to a timeout.

ErrEnqueueTimeout =

errors

.

New

("timed out enqueuing operation")

// ErrNothingNewToSnapshot is returned when trying to create a snapshot

// but there's nothing new committed to the FSM since we started.

ErrNothingNewToSnapshot =

errors

.

New

("nothing new to snapshot")

// ErrUnsupportedProtocol is returned when an operation is attempted

// that's not supported by the current protocol version.

ErrUnsupportedProtocol =

errors

.

New

("operation not supported with current protocol version")

// ErrCantBootstrap is returned when attempt is made to bootstrap a

// cluster that already has state present.

ErrCantBootstrap =

errors

.

New

("bootstrap only works on new clusters")

// ErrLeadershipTransferInProgress is returned when the leader is rejecting

// client requests because it is attempting to transfer leadership.

ErrLeadershipTransferInProgress =

errors

.

New

("leadership transfer in progress")
)

View Source

var (

// ErrTransportShutdown is returned when operations on a transport are

// invoked after it's been terminated.

ErrTransportShutdown =

errors

.

New

("transport shutdown")

// ErrPipelineShutdown is returned when the pipeline is closed.

ErrPipelineShutdown =

errors

.

New

("append pipeline closed")
)

View Source

var (

// ErrLogNotFound indicates a given log entry is not available.

ErrLogNotFound =

errors

.

New

("log not found")

// ErrPipelineReplicationNotSupported can be returned by the transport to

// signal that pipeline replication is not supported in general, and that

// no error message should be produced.

ErrPipelineReplicationNotSupported =

errors

.

New

("pipeline replication not supported")
)

Functions

¶

func

BootstrapCluster

¶

added in

v1.0.0

func BootstrapCluster(conf *

Config

, logs

LogStore

, stable

StableStore

,
	snaps

SnapshotStore

, trans

Transport

, configuration

Configuration

,
)

error

BootstrapCluster initializes a server's storage with the given cluster
configuration. This should only be called at the beginning of time for the
cluster with an identical configuration listing all Voter servers. There is
no need to bootstrap Nonvoter and Staging servers.

A cluster can only be bootstrapped once from a single participating Voter
server. Any further attempts to bootstrap will return an error that can be
safely ignored.

One approach is to bootstrap a single server with a configuration
listing just itself as a Voter, then invoke AddVoter() on it to add other
servers to the cluster.

func

EncodeConfiguration

¶

added in

v1.1.2

func EncodeConfiguration(configuration

Configuration

) []

byte

EncodeConfiguration serializes a Configuration using MsgPack, or panics on
errors.

func

HasExistingState

¶

added in

v1.0.0

func HasExistingState(logs

LogStore

, stable

StableStore

, snaps

SnapshotStore

) (

bool

,

error

)

HasExistingState returns true if the server has any existing state (logs,
knowledge of a current term, or any snapshots).

func

MakeCluster

¶

added in

v1.1.1

func MakeCluster(n

int

, t *

testing

.

T

, conf *

Config

) *cluster

NOTE: This is exposed for middleware testing purposes and is not a stable API

func

MakeClusterCustom

¶

added in

v1.1.1

func MakeClusterCustom(t *

testing

.

T

, opts *

MakeClusterOpts

) *cluster

NOTE: This is exposed for middleware testing purposes and is not a stable API

func

MakeClusterNoBootstrap

¶

added in

v1.1.1

func MakeClusterNoBootstrap(n

int

, t *

testing

.

T

, conf *

Config

) *cluster

NOTE: This is exposed for middleware testing purposes and is not a stable API

func

NewInmemTransport

¶

func NewInmemTransport(addr

ServerAddress

) (

ServerAddress

, *

InmemTransport

)

NewInmemTransport is used to initialize a new transport
and generates a random local address if none is specified

func

NewInmemTransportWithTimeout

¶

added in

v1.0.1

func NewInmemTransportWithTimeout(addr

ServerAddress

, timeout

time

.

Duration

) (

ServerAddress

, *

InmemTransport

)

NewInmemTransportWithTimeout is used to initialize a new transport and
generates a random local address if none is specified. The given timeout
will be used to decide how long to wait for a connected peer to process the
RPCs that we're sending it. See also Connect() and Consumer().

func

RecoverCluster

¶

added in

v1.0.0

func RecoverCluster(conf *

Config

, fsm

FSM

, logs

LogStore

, stable

StableStore

,
	snaps

SnapshotStore

, trans

Transport

, configuration

Configuration

,
)

error

RecoverCluster is used to manually force a new configuration in order to
recover from a loss of quorum where the current configuration cannot be
restored, such as when several servers die at the same time. This works by
reading all the current state for this server, creating a snapshot with the
supplied configuration, and then truncating the Raft log. This is the only
safe way to force a given configuration without actually altering the log to
insert any new entries, which could cause conflicts with other servers with
different state.

WARNING! This operation implicitly commits all entries in the Raft log, so
in general this is an extremely unsafe operation. If you've lost your other
servers and are performing a manual recovery, then you've also lost the
commit information, so this is likely the best you can do, but you should be
aware that calling this can cause Raft log entries that were in the process
of being replicated but not yet be committed to be committed.

Note the FSM passed here is used for the snapshot operations and will be
left in a state that should not be used by the application. Be sure to
discard this FSM and any associated state and provide a fresh one when
calling NewRaft later.

A typical way to recover the cluster is to shut down all servers and then
run RecoverCluster on every server using an identical configuration. When
the cluster is then restarted, and election should occur and then Raft will
resume normal operation. If it's desired to make a particular server the
leader, this can be used to inject a new configuration with that server as
the sole voter, and then join up other new clean-state peer servers using
the usual APIs in order to bring the cluster back into a known state.

func

ValidateConfig

¶

func ValidateConfig(config *

Config

)

error

ValidateConfig is used to validate a sane configuration

Types

¶

type

AppendEntriesRequest

¶

type AppendEntriesRequest struct {

RPCHeader

// Provide the current term and leader

Term

uint64

// Deprecated: use RPCHeader.Addr instead

Leader []

byte

// Provide the previous entries for integrity checking

PrevLogEntry

uint64

PrevLogTerm

uint64

// New entries to commit

Entries []*

Log

// Commit index on the leader

LeaderCommitIndex

uint64

}

AppendEntriesRequest is the command used to append entries to the
replicated log.

func (*AppendEntriesRequest)

GetRPCHeader

¶

added in

v1.0.0

func (r *

AppendEntriesRequest

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

AppendEntriesResponse

¶

type AppendEntriesResponse struct {

RPCHeader

// Newer term if leader is out of date

Term

uint64

// Last Log is a hint to help accelerate rebuilding slow nodes

LastLog

uint64

// We may not succeed if we have a conflicting entry

Success

bool

// There are scenarios where this request didn't succeed

// but there's no need to wait/back-off the next attempt.

NoRetryBackoff

bool

}

AppendEntriesResponse is the response returned from an
AppendEntriesRequest.

func (*AppendEntriesResponse)

GetRPCHeader

¶

added in

v1.0.0

func (r *

AppendEntriesResponse

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

AppendFuture

¶

type AppendFuture interface {

Future

// Start returns the time that the append request was started.

// It is always OK to call this method.

Start()

time

.

Time

// Request holds the parameters of the AppendEntries call.

// It is always OK to call this method.

Request() *

AppendEntriesRequest

// Response holds the results of the AppendEntries call.

// This method must only be called after the Error

// method returns, and will only be valid on success.

Response() *

AppendEntriesResponse

}

AppendFuture is used to return information about a pipelined AppendEntries request.

type

AppendPipeline

¶

type AppendPipeline interface {

// AppendEntries is used to add another request to the pipeline.

// The send may block which is an effective form of back-pressure.

AppendEntries(args *

AppendEntriesRequest

, resp *

AppendEntriesResponse

) (

AppendFuture

,

error

)

// Consumer returns a channel that can be used to consume

// response futures when they are ready.

Consumer() <-chan

AppendFuture

// Close closes the pipeline and cancels all inflight RPCs

Close()

error

}

AppendPipeline is used for pipelining AppendEntries requests. It is used
to increase the replication throughput by masking latency and better
utilizing bandwidth.

type

ApplyFuture

¶

type ApplyFuture interface {

IndexFuture

// Response returns the FSM response as returned by the FSM.Apply method. This

// must not be called until after the Error method has returned.

// Note that if FSM.Apply returns an error, it will be returned by Response,

// and not by the Error method, so it is always important to check Response

// for errors from the FSM.

Response() interface{}
}

ApplyFuture is used for Apply and can return the FSM response.

type

BatchingFSM

¶

added in

v1.1.2

type BatchingFSM interface {

// ApplyBatch is invoked once a batch of log entries has been committed and

// are ready to be applied to the FSM. ApplyBatch will take in an array of

// log entries. These log entries will be in the order they were committed,

// will not have gaps, and could be of a few log types. Clients should check

// the log type prior to attempting to decode the data attached. Presently

// the LogCommand and LogConfiguration types will be sent.

//

// The returned slice must be the same length as the input and each response

// should correlate to the log at the same index of the input. The returned

// values will be made available in the ApplyFuture returned by Raft.Apply

// method if that method was called on the same Raft node as the FSM.

ApplyBatch([]*

Log

) []interface{}

FSM

}

BatchingFSM extends the FSM interface to add an ApplyBatch function. This can
optionally be implemented by clients to enable multiple logs to be applied to
the FSM in batches. Up to MaxAppendEntries could be sent in a batch.

type

Config

¶

type Config struct {

// ProtocolVersion allows a Raft server to inter-operate with older

// Raft servers running an older version of the code. This is used to

// version the wire protocol as well as Raft-specific log entries that

// the server uses when _speaking_ to other servers. There is currently

// no auto-negotiation of versions so all servers must be manually

// configured with compatible versions. See ProtocolVersionMin and

// ProtocolVersionMax for the versions of the protocol that this server

// can _understand_.

ProtocolVersion

ProtocolVersion

// HeartbeatTimeout specifies the time in follower state without contact

// from a leader before we attempt an election.

HeartbeatTimeout

time

.

Duration

// ElectionTimeout specifies the time in candidate state without contact

// from a leader before we attempt an election.

ElectionTimeout

time

.

Duration

// CommitTimeout specifies the time without an Apply operation before the

// leader sends an AppendEntry RPC to followers, to ensure a timely commit of

// log entries.

// Due to random staggering, may be delayed as much as 2x this value.

CommitTimeout

time

.

Duration

// MaxAppendEntries controls the maximum number of append entries

// to send at once. We want to strike a balance between efficiency

// and avoiding waste if the follower is going to reject because of

// an inconsistent log.

MaxAppendEntries

int

// BatchApplyCh indicates whether we should buffer applyCh

// to size MaxAppendEntries. This enables batch log commitment,

// but breaks the timeout guarantee on Apply. Specifically,

// a log can be added to the applyCh buffer but not actually be

// processed until after the specified timeout.

BatchApplyCh

bool

// If we are a member of a cluster, and RemovePeer is invoked for the

// local node, then we forget all peers and transition into the follower state.

// If ShutdownOnRemove is set, we additional shutdown Raft. Otherwise,

// we can become a leader of a cluster containing only this node.

ShutdownOnRemove

bool

// TrailingLogs controls how many logs we leave after a snapshot. This is used

// so that we can quickly replay logs on a follower instead of being forced to

// send an entire snapshot. The value passed here is the initial setting used.

// This can be tuned during operation using ReloadConfig.

TrailingLogs

uint64

// SnapshotInterval controls how often we check if we should perform a

// snapshot. We randomly stagger between this value and 2x this value to avoid

// the entire cluster from performing a snapshot at once. The value passed

// here is the initial setting used. This can be tuned during operation using

// ReloadConfig.

SnapshotInterval

time

.

Duration

// SnapshotThreshold controls how many outstanding logs there must be before

// we perform a snapshot. This is to prevent excessive snapshotting by

// replaying a small set of logs instead. The value passed here is the initial

// setting used. This can be tuned during operation using ReloadConfig.

SnapshotThreshold

uint64

// LeaderLeaseTimeout is used to control how long the "lease" lasts

// for being the leader without being able to contact a quorum

// of nodes. If we reach this interval without contact, we will

// step down as leader.

LeaderLeaseTimeout

time

.

Duration

// LocalID is a unique ID for this server across all time. When running with

// ProtocolVersion < 3, you must set this to be the same as the network

// address of your transport.

LocalID

ServerID

// NotifyCh is used to provide a channel that will be notified of leadership

// changes. Raft will block writing to this channel, so it should either be

// buffered or aggressively consumed.

NotifyCh chan<-

bool

// LogOutput is used as a sink for logs, unless Logger is specified.

// Defaults to os.Stderr.

LogOutput

io

.

Writer

// LogLevel represents a log level. If the value does not match a known

// logging level hclog.NoLevel is used.

LogLevel

string

// Logger is a user-provided logger. If nil, a logger writing to

// LogOutput with LogLevel is used.

Logger

hclog

.

Logger

// NoSnapshotRestoreOnStart controls if raft will restore a snapshot to the

// FSM on start. This is useful if your FSM recovers from other mechanisms

// than raft snapshotting. Snapshot metadata will still be used to initialize

// raft's configuration and index values.

NoSnapshotRestoreOnStart

bool

// PreVoteDisabled deactivate the pre-vote feature when set to true

PreVoteDisabled

bool

// NoLegacyTelemetry allows to skip the legacy metrics to avoid duplicates.

// legacy metrics are those that have `_peer_name` as metric suffix instead as labels.

// e.g: raft_replication_heartbeat_peer0

NoLegacyTelemetry

bool

// contains filtered or unexported fields

}

Config provides any necessary configuration for the Raft server.

func

DefaultConfig

¶

func DefaultConfig() *

Config

DefaultConfig returns a Config with usable defaults.

type

Configuration

¶

added in

v1.0.0

type Configuration struct {

Servers []

Server

}

Configuration tracks which servers are in the cluster, and whether they have
votes. This should include the local server, if it's a member of the cluster.
The servers are listed no particular order, but each should only appear once.
These entries are appended to the log during membership changes.

func

DecodeConfiguration

¶

added in

v1.1.2

func DecodeConfiguration(buf []

byte

)

Configuration

DecodeConfiguration deserializes a Configuration using MsgPack, or panics on
errors.

func

GetConfiguration

¶

added in

v1.1.2

func GetConfiguration(conf *

Config

, fsm

FSM

, logs

LogStore

, stable

StableStore

,
	snaps

SnapshotStore

, trans

Transport

,
) (

Configuration

,

error

)

GetConfiguration returns the persisted configuration of the Raft cluster
without starting a Raft instance or connecting to the cluster. This function
has identical behavior to Raft.GetConfiguration.

func

ReadConfigJSON

¶

added in

v1.0.0

func ReadConfigJSON(path

string

) (

Configuration

,

error

)

ReadConfigJSON reads a new-style peers.json and returns a configuration
structure. This can be used to perform manual recovery when running protocol
versions that use server IDs.

func

ReadPeersJSON

¶

added in

v1.0.0

func ReadPeersJSON(path

string

) (

Configuration

,

error

)

ReadPeersJSON consumes a legacy peers.json file in the format of the old JSON
peer store and creates a new-style configuration structure. This can be used
to migrate this data or perform manual recovery when running protocol versions
that can interoperate with older, unversioned Raft servers. This should not be
used once server IDs are in use, because the old peers.json file didn't have
support for these, nor non-voter suffrage types.

func (*Configuration)

Clone

¶

added in

v1.0.0

func (c *

Configuration

) Clone() (copy

Configuration

)

Clone makes a deep copy of a Configuration.

type

ConfigurationChangeCommand

¶

added in

v1.0.0

type ConfigurationChangeCommand

uint8

ConfigurationChangeCommand is the different ways to change the cluster
configuration.

const (

// AddVoter adds a server with Suffrage of Voter.

AddVoter

ConfigurationChangeCommand

=

iota

// AddNonvoter makes a server Nonvoter unless its Staging or Voter.

AddNonvoter

// DemoteVoter makes a server Nonvoter unless its absent.

DemoteVoter

// RemoveServer removes a server entirely from the cluster membership.

RemoveServer

// Promote changes a server from Staging to Voter. The command will be a

// no-op if the server is not Staging.

// Deprecated: use AddVoter instead.

Promote

// AddStaging makes a server a Voter.

// Deprecated: AddStaging was actually AddVoter. Use AddVoter instead.

AddStaging = 0

// explicit 0 to preserve the old value.

)

func (ConfigurationChangeCommand)

String

¶

added in

v1.0.0

func (c

ConfigurationChangeCommand

) String()

string

type

ConfigurationFuture

¶

added in

v1.0.0

type ConfigurationFuture interface {

IndexFuture

// Configuration contains the latest configuration. This must

// not be called until after the Error method has returned.

Configuration()

Configuration

}

ConfigurationFuture is used for GetConfiguration and can return the
latest configuration in use by Raft.

type

ConfigurationStore

¶

added in

v1.1.1

type ConfigurationStore interface {

// ConfigurationStore is a superset of the FSM functionality

FSM

// StoreConfiguration is invoked once a log entry containing a configuration

// change is committed. It takes the index at which the configuration was

// written and the configuration value.

StoreConfiguration(index

uint64

, configuration

Configuration

)
}

ConfigurationStore provides an interface that can optionally be implemented by FSMs
to store configuration updates made in the replicated log. In general this is only
necessary for FSMs that mutate durable state directly instead of applying changes
in memory and snapshotting periodically. By storing configuration changes, the
persistent FSM state can behave as a complete snapshot, and be able to recover
without an external snapshot just for persisting the raft configuration.

type

CountingReader

¶

added in

v1.3.5

type CountingReader interface {

io

.

Reader

Count()

int64

}

type

DiscardSnapshotSink

¶

type DiscardSnapshotSink struct{}

DiscardSnapshotSink is used to fulfill the SnapshotSink interface
while always discarding the . This is useful for when the log
should be truncated but no snapshot should be retained. This
should never be used for production use, and is only suitable
for testing.

func (*DiscardSnapshotSink)

Cancel

¶

func (d *

DiscardSnapshotSink

) Cancel()

error

Cancel returns successfully with a nil error

func (*DiscardSnapshotSink)

Close

¶

func (d *

DiscardSnapshotSink

) Close()

error

Close returns a nil error

func (*DiscardSnapshotSink)

ID

¶

func (d *

DiscardSnapshotSink

) ID()

string

ID returns "discard" for DiscardSnapshotSink

func (*DiscardSnapshotSink)

Write

¶

func (d *

DiscardSnapshotSink

) Write(b []

byte

) (

int

,

error

)

Write returns successfully with the length of the input byte slice
to satisfy the WriteCloser interface

type

DiscardSnapshotStore

¶

type DiscardSnapshotStore struct{}

DiscardSnapshotStore is used to successfully snapshot while
always discarding the snapshot. This is useful for when the
log should be truncated but no snapshot should be retained.
This should never be used for production use, and is only
suitable for testing.

func

NewDiscardSnapshotStore

¶

func NewDiscardSnapshotStore() *

DiscardSnapshotStore

NewDiscardSnapshotStore is used to create a new DiscardSnapshotStore.

func (*DiscardSnapshotStore)

Create

¶

func (d *

DiscardSnapshotStore

) Create(version

SnapshotVersion

, index, term

uint64

,
	configuration

Configuration

, configurationIndex

uint64

, trans

Transport

) (

SnapshotSink

,

error

)

Create returns a valid type implementing the SnapshotSink which
always discards the snapshot.

func (*DiscardSnapshotStore)

List

¶

func (d *

DiscardSnapshotStore

) List() ([]*

SnapshotMeta

,

error

)

List returns successfully with a nil for []*SnapshotMeta.

func (*DiscardSnapshotStore)

Open

¶

func (d *

DiscardSnapshotStore

) Open(id

string

) (*

SnapshotMeta

,

io

.

ReadCloser

,

error

)

Open returns an error since the DiscardSnapshotStore does not
support opening snapshots.

type

FSM

¶

type FSM interface {

// Apply is called once a log entry is committed by a majority of the cluster.

//

// Apply should apply the log to the FSM. Apply must be deterministic and

// produce the same result on all peers in the cluster.

//

// The returned value is returned to the client as the ApplyFuture.Response.

Apply(*

Log

) interface{}

// Snapshot returns an FSMSnapshot used to: support log compaction, to

// restore the FSM to a previous state, or to bring out-of-date followers up

// to a recent log index.

//

// The Snapshot implementation should return quickly, because Apply can not

// be called while Snapshot is running. Generally this means Snapshot should

// only capture a pointer to the state, and any expensive IO should happen

// as part of FSMSnapshot.Persist.

//

// Apply and Snapshot are always called from the same thread, but Apply will

// be called concurrently with FSMSnapshot.Persist. This means the FSM should

// be implemented to allow for concurrent updates while a snapshot is happening.

//

// Clients of this library should make no assumptions about whether a returned

// Snapshot() will actually be stored by Raft. In fact it's quite possible that

// any Snapshot returned by this call will be discarded, and that

// FSMSnapshot.Persist will never be called. Raft will always call

// FSMSnapshot.Release however.

Snapshot() (

FSMSnapshot

,

error

)

// Restore is used to restore an FSM from a snapshot. It is not called

// concurrently with any other command. The FSM must discard all previous

// state before restoring the snapshot.

Restore(snapshot

io

.

ReadCloser

)

error

}

FSM is implemented by clients to make use of the replicated log.

type

FSMSnapshot

¶

type FSMSnapshot interface {

// Persist should dump all necessary state to the WriteCloser 'sink',

// and call sink.Close() when finished or call sink.Cancel() on error.

Persist(sink

SnapshotSink

)

error

// Release is invoked when we are finished with the snapshot.

Release()
}

FSMSnapshot is returned by an FSM in response to a Snapshot
It must be safe to invoke FSMSnapshot methods with concurrent
calls to Apply.

type

FailedHeartbeatObservation

¶

added in

v1.3.0

type FailedHeartbeatObservation struct {

PeerID

ServerID

LastContact

time

.

Time

}

FailedHeartbeatObservation is sent when a node fails to heartbeat with the leader

type

FileSnapshotSink

¶

type FileSnapshotSink struct {

// contains filtered or unexported fields

}

FileSnapshotSink implements SnapshotSink with a file.

func (*FileSnapshotSink)

Cancel

¶

func (s *

FileSnapshotSink

) Cancel()

error

Cancel is used to indicate an unsuccessful end.

func (*FileSnapshotSink)

Close

¶

func (s *

FileSnapshotSink

) Close()

error

Close is used to indicate a successful end.

func (*FileSnapshotSink)

ID

¶

func (s *

FileSnapshotSink

) ID()

string

ID returns the ID of the snapshot, can be used with Open()
after the snapshot is finalized.

func (*FileSnapshotSink)

Write

¶

func (s *

FileSnapshotSink

) Write(b []

byte

) (

int

,

error

)

Write is used to append to the state file. We write to the
buffered IO object to reduce the amount of context switches.

type

FileSnapshotStore

¶

type FileSnapshotStore struct {

// contains filtered or unexported fields

}

FileSnapshotStore implements the SnapshotStore interface and allows
snapshots to be made on the local disk.

func

FileSnapTest

¶

added in

v1.1.1

func FileSnapTest(t *

testing

.

T

) (

string

, *

FileSnapshotStore

)

NOTE: This is exposed for middleware testing purposes and is not a stable API

func

NewFileSnapshotStore

¶

func NewFileSnapshotStore(base

string

, retain

int

, logOutput

io

.

Writer

) (*

FileSnapshotStore

,

error

)

NewFileSnapshotStore creates a new FileSnapshotStore based
on a base directory. The `retain` parameter controls how many
snapshots are retained. Must be at least 1.

func

NewFileSnapshotStoreWithLogger

¶

func NewFileSnapshotStoreWithLogger(base

string

, retain

int

, logger

hclog

.

Logger

) (*

FileSnapshotStore

,

error

)

NewFileSnapshotStoreWithLogger creates a new FileSnapshotStore based
on a base directory. The `retain` parameter controls how many
snapshots are retained. Must be at least 1.

func (*FileSnapshotStore)

Create

¶

func (f *

FileSnapshotStore

) Create(version

SnapshotVersion

, index, term

uint64

,
	configuration

Configuration

, configurationIndex

uint64

, trans

Transport

) (

SnapshotSink

,

error

)

Create is used to start a new snapshot

func (*FileSnapshotStore)

List

¶

func (f *

FileSnapshotStore

) List() ([]*

SnapshotMeta

,

error

)

List returns available snapshots in the store.

func (*FileSnapshotStore)

Open

¶

func (f *

FileSnapshotStore

) Open(id

string

) (*

SnapshotMeta

,

io

.

ReadCloser

,

error

)

Open takes a snapshot ID and returns a ReadCloser for that snapshot.

func (*FileSnapshotStore)

ReapSnapshots

¶

func (f *

FileSnapshotStore

) ReapSnapshots()

error

ReapSnapshots reaps any snapshots beyond the retain count.

type

FilterFn

¶

type FilterFn func(o *

Observation

)

bool

FilterFn is a function that can be registered in order to filter observations.
The function reports whether the observation should be included - if
it returns false, the observation will be filtered out.

type

Future

¶

type Future interface {

// Error blocks until the future arrives and then returns the error status

// of the future. This may be called any number of times - all calls will

// return the same value, however is not OK to call this method twice

// concurrently on the same Future instance.

// Error will only return generic errors related to raft, such

// as ErrLeadershipLost, or ErrRaftShutdown. Some operations, such as

// ApplyLog, may also return errors from other methods.

Error()

error

}

Future is used to represent an action that may occur in the future.

type

IndexFuture

¶

added in

v1.0.0

type IndexFuture interface {

Future

// Index holds the index of the newly applied log entry.

// This must not be called until after the Error method has returned.

Index()

uint64

}

IndexFuture is used for future actions that can result in a raft log entry
being created.

type

InmemSnapshotSink

¶

added in

v1.0.0

type InmemSnapshotSink struct {

// contains filtered or unexported fields

}

InmemSnapshotSink implements SnapshotSink in memory

func (*InmemSnapshotSink)

Cancel

¶

added in

v1.0.0

func (s *

InmemSnapshotSink

) Cancel()

error

Cancel returns successfully with a nil error

func (*InmemSnapshotSink)

Close

¶

added in

v1.0.0

func (s *

InmemSnapshotSink

) Close()

error

Close updates the Size and is otherwise a no-op

func (*InmemSnapshotSink)

ID

¶

added in

v1.0.0

func (s *

InmemSnapshotSink

) ID()

string

ID returns the ID of the SnapshotMeta

func (*InmemSnapshotSink)

Write

¶

added in

v1.0.0

func (s *

InmemSnapshotSink

) Write(p []

byte

) (n

int

, err

error

)

Write appends the given bytes to the snapshot contents

type

InmemSnapshotStore

¶

added in

v1.0.0

type InmemSnapshotStore struct {

sync

.

RWMutex

// contains filtered or unexported fields

}

InmemSnapshotStore implements the SnapshotStore interface and
retains only the most recent snapshot

func

NewInmemSnapshotStore

¶

added in

v1.0.0

func NewInmemSnapshotStore() *

InmemSnapshotStore

NewInmemSnapshotStore creates a blank new InmemSnapshotStore

func (*InmemSnapshotStore)

Create

¶

added in

v1.0.0

func (m *

InmemSnapshotStore

) Create(version

SnapshotVersion

, index, term

uint64

,
	configuration

Configuration

, configurationIndex

uint64

, trans

Transport

) (

SnapshotSink

,

error

)

Create replaces the stored snapshot with a new one using the given args

func (*InmemSnapshotStore)

List

¶

added in

v1.0.0

func (m *

InmemSnapshotStore

) List() ([]*

SnapshotMeta

,

error

)

List returns the latest snapshot taken

func (*InmemSnapshotStore)

Open

¶

added in

v1.0.0

func (m *

InmemSnapshotStore

) Open(id

string

) (*

SnapshotMeta

,

io

.

ReadCloser

,

error

)

Open wraps an io.ReadCloser around the snapshot contents

type

InmemStore

¶

type InmemStore struct {

// contains filtered or unexported fields

}

InmemStore implements the LogStore and StableStore interface.
It should NOT EVER be used for production. It is used only for
unit tests. Use the MDBStore implementation instead.

func

NewInmemStore

¶

func NewInmemStore() *

InmemStore

NewInmemStore returns a new in-memory backend. Do not ever
use for production. Only for testing.

func (*InmemStore)

DeleteRange

¶

func (i *

InmemStore

) DeleteRange(min, max

uint64

)

error

DeleteRange implements the LogStore interface.

func (*InmemStore)

FirstIndex

¶

func (i *

InmemStore

) FirstIndex() (

uint64

,

error

)

FirstIndex implements the LogStore interface.

func (*InmemStore)

Get

¶

func (i *

InmemStore

) Get(key []

byte

) ([]

byte

,

error

)

Get implements the StableStore interface.

func (*InmemStore)

GetLog

¶

func (i *

InmemStore

) GetLog(index

uint64

, log *

Log

)

error

GetLog implements the LogStore interface.

func (*InmemStore)

GetUint64

¶

func (i *

InmemStore

) GetUint64(key []

byte

) (

uint64

,

error

)

GetUint64 implements the StableStore interface.

func (*InmemStore)

LastIndex

¶

func (i *

InmemStore

) LastIndex() (

uint64

,

error

)

LastIndex implements the LogStore interface.

func (*InmemStore)

Set

¶

func (i *

InmemStore

) Set(key []

byte

, val []

byte

)

error

Set implements the StableStore interface.

func (*InmemStore)

SetUint64

¶

func (i *

InmemStore

) SetUint64(key []

byte

, val

uint64

)

error

SetUint64 implements the StableStore interface.

func (*InmemStore)

StoreLog

¶

func (i *

InmemStore

) StoreLog(log *

Log

)

error

StoreLog implements the LogStore interface.

func (*InmemStore)

StoreLogs

¶

func (i *

InmemStore

) StoreLogs(logs []*

Log

)

error

StoreLogs implements the LogStore interface.

type

InmemTransport

¶

type InmemTransport struct {

sync

.

RWMutex

// contains filtered or unexported fields

}

InmemTransport Implements the Transport interface, to allow Raft to be
tested in-memory without going over a network.

func (*InmemTransport)

AppendEntries

¶

func (i *

InmemTransport

) AppendEntries(id

ServerID

, target

ServerAddress

, args *

AppendEntriesRequest

, resp *

AppendEntriesResponse

)

error

AppendEntries implements the Transport interface.

func (*InmemTransport)

AppendEntriesPipeline

¶

func (i *

InmemTransport

) AppendEntriesPipeline(id

ServerID

, target

ServerAddress

) (

AppendPipeline

,

error

)

AppendEntriesPipeline returns an interface that can be used to pipeline
AppendEntries requests.

func (*InmemTransport)

Close

¶

func (i *

InmemTransport

) Close()

error

Close is used to permanently disable the transport

func (*InmemTransport)

Connect

¶

func (i *

InmemTransport

) Connect(peer

ServerAddress

, t

Transport

)

Connect is used to connect this transport to another transport for
a given peer name. This allows for local routing.

func (*InmemTransport)

Consumer

¶

func (i *

InmemTransport

) Consumer() <-chan

RPC

Consumer implements the Transport interface.

func (*InmemTransport)

DecodePeer

¶

func (i *

InmemTransport

) DecodePeer(buf []

byte

)

ServerAddress

DecodePeer implements the Transport interface.

func (*InmemTransport)

Disconnect

¶

func (i *

InmemTransport

) Disconnect(peer

ServerAddress

)

Disconnect is used to remove the ability to route to a given peer.

func (*InmemTransport)

DisconnectAll

¶

func (i *

InmemTransport

) DisconnectAll()

DisconnectAll is used to remove all routes to peers.

func (*InmemTransport)

EncodePeer

¶

func (i *

InmemTransport

) EncodePeer(id

ServerID

, p

ServerAddress

) []

byte

EncodePeer implements the Transport interface.

func (*InmemTransport)

InstallSnapshot

¶

func (i *

InmemTransport

) InstallSnapshot(id

ServerID

, target

ServerAddress

, args *

InstallSnapshotRequest

, resp *

InstallSnapshotResponse

, data

io

.

Reader

)

error

InstallSnapshot implements the Transport interface.

func (*InmemTransport)

LocalAddr

¶

func (i *

InmemTransport

) LocalAddr()

ServerAddress

LocalAddr implements the Transport interface.

func (*InmemTransport)

RequestPreVote

¶

added in

v1.7.0

func (i *

InmemTransport

) RequestPreVote(id

ServerID

, target

ServerAddress

, args *

RequestPreVoteRequest

, resp *

RequestPreVoteResponse

)

error

func (*InmemTransport)

RequestVote

¶

func (i *

InmemTransport

) RequestVote(id

ServerID

, target

ServerAddress

, args *

RequestVoteRequest

, resp *

RequestVoteResponse

)

error

RequestVote implements the Transport interface.

func (*InmemTransport)

SetHeartbeatHandler

¶

func (i *

InmemTransport

) SetHeartbeatHandler(cb func(

RPC

))

SetHeartbeatHandler is used to set optional fast-path for
heartbeats, not supported for this transport.

func (*InmemTransport)

TimeoutNow

¶

added in

v1.1.0

func (i *

InmemTransport

) TimeoutNow(id

ServerID

, target

ServerAddress

, args *

TimeoutNowRequest

, resp *

TimeoutNowResponse

)

error

TimeoutNow implements the Transport interface.

type

InstallSnapshotRequest

¶

type InstallSnapshotRequest struct {

RPCHeader

SnapshotVersion

SnapshotVersion

Term

uint64

Leader []

byte

// These are the last index/term included in the snapshot

LastLogIndex

uint64

LastLogTerm

uint64

// Peer Set in the snapshot.

// but remains here in case we receive an InstallSnapshot from a leader

// that's running old code.

// Deprecated: This is deprecated in favor of Configuration

Peers []

byte

// Cluster membership.

Configuration []

byte

// Log index where 'Configuration' entry was originally written.

ConfigurationIndex

uint64

// Size of the snapshot

Size

int64

}

InstallSnapshotRequest is the command sent to a Raft peer to bootstrap its
log (and state machine) from a snapshot on another peer.

func (*InstallSnapshotRequest)

GetRPCHeader

¶

added in

v1.0.0

func (r *

InstallSnapshotRequest

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

InstallSnapshotResponse

¶

type InstallSnapshotResponse struct {

RPCHeader

Term

uint64

Success

bool

}

InstallSnapshotResponse is the response returned from an
InstallSnapshotRequest.

func (*InstallSnapshotResponse)

GetRPCHeader

¶

added in

v1.0.0

func (r *

InstallSnapshotResponse

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

LeaderObservation

¶

type LeaderObservation struct {

// DEPRECATED The LeaderAddr field should now be used

Leader

ServerAddress

LeaderAddr

ServerAddress

LeaderID

ServerID

}

LeaderObservation is used for the data when leadership changes.

type

LeadershipTransferFuture

¶

added in

v1.1.0

type LeadershipTransferFuture interface {

Future

}

LeadershipTransferFuture is used for waiting on a user-triggered leadership
transfer to complete.

type

Log

¶

type Log struct {

// Index holds the index of the log entry.

Index

uint64

// Term holds the election term of the log entry.

Term

uint64

// Type holds the type of the log entry.

Type

LogType

// Data holds the log entry's type-specific data.

Data []

byte

// Extensions holds an opaque byte slice of information for middleware. It

// is up to the client of the library to properly modify this as it adds

// layers and remove those layers when appropriate. This value is a part of

// the log, so very large values could cause timing issues.

//

// N.B. It is _up to the client_ to handle upgrade paths. For instance if

// using this with go-raftchunking, the client should ensure that all Raft

// peers are using a version that can handle that extension before ever

// actually triggering chunking behavior. It is sometimes sufficient to

// ensure that non-leaders are upgraded first, then the current leader is

// upgraded, but a leader changeover during this process could lead to

// trouble, so gating extension behavior via some flag in the client

// program is also a good idea.

Extensions []

byte

// AppendedAt stores the time the leader first appended this log to it's

// LogStore. Followers will observe the leader's time. It is not used for

// coordination or as part of the replication protocol at all. It exists only

// to provide operational information for example how many seconds worth of

// logs are present on the leader which might impact follower's ability to

// catch up after restoring a large snapshot. We should never rely on this

// being in the past when appending on a follower or reading a log back since

// the clock skew can mean a follower could see a log with a future timestamp.

// In general too the leader is not required to persist the log before

// delivering to followers although the current implementation happens to do

// this.

AppendedAt

time

.

Time

}

Log entries are replicated to all members of the Raft cluster
and form the heart of the replicated state machine.

type

LogCache

¶

type LogCache struct {

// contains filtered or unexported fields

}

LogCache wraps any LogStore implementation to provide an
in-memory ring buffer. This is used to cache access to
the recently written entries. For implementations that do not
cache themselves, this can provide a substantial boost by
avoiding disk I/O on recent entries.

func

NewLogCache

¶

func NewLogCache(capacity

int

, store

LogStore

) (*

LogCache

,

error

)

NewLogCache is used to create a new LogCache with the
given capacity and backend store.

func (*LogCache)

DeleteRange

¶

func (c *

LogCache

) DeleteRange(min, max

uint64

)

error

func (*LogCache)

FirstIndex

¶

func (c *

LogCache

) FirstIndex() (

uint64

,

error

)

func (*LogCache)

GetLog

¶

func (c *

LogCache

) GetLog(idx

uint64

, log *

Log

)

error

func (*LogCache)

IsMonotonic

¶

added in

v1.4.0

func (c *

LogCache

) IsMonotonic()

bool

IsMonotonic implements the MonotonicLogStore interface. This is a shim to
expose the underlying store as monotonically indexed or not.

func (*LogCache)

LastIndex

¶

func (c *

LogCache

) LastIndex() (

uint64

,

error

)

func (*LogCache)

StoreLog

¶

func (c *

LogCache

) StoreLog(log *

Log

)

error

func (*LogCache)

StoreLogs

¶

func (c *

LogCache

) StoreLogs(logs []*

Log

)

error

type

LogStore

¶

type LogStore interface {

// FirstIndex returns the first index written. 0 for no entries.

FirstIndex() (

uint64

,

error

)

// LastIndex returns the last index written. 0 for no entries.

LastIndex() (

uint64

,

error

)

// GetLog gets a log entry at a given index.

GetLog(index

uint64

, log *

Log

)

error

// StoreLog stores a log entry.

StoreLog(log *

Log

)

error

// StoreLogs stores multiple log entries. By default the logs stored may not be contiguous with previous logs (i.e. may have a gap in Index since the last log written). If an implementation can't tolerate this it may optionally implement `MonotonicLogStore` to indicate that this is not allowed. This changes Raft's behaviour after restoring a user snapshot to remove all previous logs instead of relying on a "gap" to signal the discontinuity between logs before the snapshot and logs after.

StoreLogs(logs []*

Log

)

error

// DeleteRange deletes a range of log entries. The range is inclusive.

DeleteRange(min, max

uint64

)

error

}

LogStore is used to provide an interface for storing
and retrieving logs in a durable fashion.

type

LogType

¶

type LogType

uint8

LogType describes various types of log entries.

const (

// LogCommand is applied to a user FSM.

LogCommand

LogType

=

iota

// LogNoop is used to assert leadership.

LogNoop

// LogAddPeerDeprecated is used to add a new peer. This should only be used with

// older protocol versions designed to be compatible with unversioned

// Raft servers. See comments in config.go for details.

LogAddPeerDeprecated

// LogRemovePeerDeprecated is used to remove an existing peer. This should only be

// used with older protocol versions designed to be compatible with

// unversioned Raft servers. See comments in config.go for details.

LogRemovePeerDeprecated

// LogBarrier is used to ensure all preceding operations have been

// applied to the FSM. It is similar to LogNoop, but instead of returning

// once committed, it only returns once the FSM manager acks it. Otherwise,

// it is possible there are operations committed but not yet applied to

// the FSM.

LogBarrier

// LogConfiguration establishes a membership change configuration. It is

// created when a server is added, removed, promoted, etc. Only used

// when protocol version 1 or greater is in use.

LogConfiguration
)

func (LogType)

String

¶

added in

v1.3.0

func (lt

LogType

) String()

string

String returns LogType as a human readable string.

type

LoopbackTransport

¶

type LoopbackTransport interface {

Transport

// Embedded transport reference

WithPeers

// Embedded peer management

WithClose

// with a close routine

WithPreVote

// with a prevote

}

LoopbackTransport is an interface that provides a loopback transport suitable for testing
e.g. InmemTransport. It's there so we don't have to rewrite tests.

type

MakeClusterOpts

¶

added in

v1.1.1

type MakeClusterOpts struct {

Peers

int

Bootstrap

bool

Conf            *

Config

ConfigStoreFSM

bool

MakeFSMFunc     func()

FSM

LongstopTimeout

time

.

Duration

MonotonicLogs

bool

}

NOTE: This is exposed for middleware testing purposes and is not a stable API

type

MockFSM

¶

added in

v1.1.1

type MockFSM struct {

sync

.

Mutex

// contains filtered or unexported fields

}

MockFSM is an implementation of the FSM interface, and just stores
the logs sequentially.

NOTE: This is exposed for middleware testing purposes and is not a stable API

func (*MockFSM)

Apply

¶

added in

v1.1.1

func (m *

MockFSM

) Apply(log *

Log

) interface{}

NOTE: This is exposed for middleware testing purposes and is not a stable API

func (*MockFSM)

Logs

¶

added in

v1.1.1

func (m *

MockFSM

) Logs() [][]

byte

NOTE: This is exposed for middleware testing purposes and is not a stable API

func (*MockFSM)

Restore

¶

added in

v1.1.1

func (m *

MockFSM

) Restore(inp

io

.

ReadCloser

)

error

NOTE: This is exposed for middleware testing purposes and is not a stable API

func (*MockFSM)

Snapshot

¶

added in

v1.1.1

func (m *

MockFSM

) Snapshot() (

FSMSnapshot

,

error

)

NOTE: This is exposed for middleware testing purposes and is not a stable API

type

MockFSMConfigStore

¶

added in

v1.1.1

type MockFSMConfigStore struct {

FSM

}

NOTE: This is exposed for middleware testing purposes and is not a stable API

func (*MockFSMConfigStore)

StoreConfiguration

¶

added in

v1.1.1

func (m *

MockFSMConfigStore

) StoreConfiguration(index

uint64

, config

Configuration

)

NOTE: This is exposed for middleware testing purposes and is not a stable API

type

MockMonotonicLogStore

¶

added in

v1.4.0

type MockMonotonicLogStore struct {

// contains filtered or unexported fields

}

MockMonotonicLogStore is a LogStore wrapper for testing the
MonotonicLogStore interface.

func (*MockMonotonicLogStore)

DeleteRange

¶

added in

v1.4.0

func (m *

MockMonotonicLogStore

) DeleteRange(min

uint64

, max

uint64

)

error

DeleteRange implements the LogStore interface.

func (*MockMonotonicLogStore)

FirstIndex

¶

added in

v1.4.0

func (m *

MockMonotonicLogStore

) FirstIndex() (

uint64

,

error

)

FirstIndex implements the LogStore interface.

func (*MockMonotonicLogStore)

GetLog

¶

added in

v1.4.0

func (m *

MockMonotonicLogStore

) GetLog(index

uint64

, log *

Log

)

error

GetLog implements the LogStore interface.

func (*MockMonotonicLogStore)

IsMonotonic

¶

added in

v1.4.0

func (m *

MockMonotonicLogStore

) IsMonotonic()

bool

IsMonotonic implements the MonotonicLogStore interface.

func (*MockMonotonicLogStore)

LastIndex

¶

added in

v1.4.0

func (m *

MockMonotonicLogStore

) LastIndex() (

uint64

,

error

)

LastIndex implements the LogStore interface.

func (*MockMonotonicLogStore)

StoreLog

¶

added in

v1.4.0

func (m *

MockMonotonicLogStore

) StoreLog(log *

Log

)

error

StoreLog implements the LogStore interface.

func (*MockMonotonicLogStore)

StoreLogs

¶

added in

v1.4.0

func (m *

MockMonotonicLogStore

) StoreLogs(logs []*

Log

)

error

StoreLogs implements the LogStore interface.

type

MockSnapshot

¶

added in

v1.1.1

type MockSnapshot struct {

// contains filtered or unexported fields

}

NOTE: This is exposed for middleware testing purposes and is not a stable API

func (*MockSnapshot)

Persist

¶

added in

v1.1.1

func (m *

MockSnapshot

) Persist(sink

SnapshotSink

)

error

NOTE: This is exposed for middleware testing purposes and is not a stable API

func (*MockSnapshot)

Release

¶

added in

v1.1.1

func (m *

MockSnapshot

) Release()

NOTE: This is exposed for middleware testing purposes and is not a stable API

type

MonotonicLogStore

¶

added in

v1.4.0

type MonotonicLogStore interface {

IsMonotonic()

bool

}

MonotonicLogStore is an optional interface for LogStore implementations that
cannot tolerate gaps in between the Index values of consecutive log entries. For example,
this may allow more efficient indexing because the Index values are densely populated. If true is
returned, Raft will avoid relying on gaps to trigger re-synching logs on followers after a
snapshot is restored. The LogStore must have an efficient implementation of
DeleteLogs for the case where all logs are removed, as this must be called after snapshot restore when gaps are not allowed.
We avoid deleting all records for LogStores that do not implement MonotonicLogStore
because although it's always correct to do so, it has a major negative performance impact on the BoltDB store that is currently
the most widely used.

type

NetworkTransport

¶

type NetworkTransport struct {

TimeoutScale

int

// contains filtered or unexported fields

}

NetworkTransport provides a network based transport that can be
used to communicate with Raft on remote machines. It requires
an underlying stream layer to provide a stream abstraction, which can
be simple TCP, TLS, etc.

This transport is very simple and lightweight. Each RPC request is
framed by sending a byte that indicates the message type, followed
by the MsgPack encoded request.

The response is an error string followed by the response object,
both are encoded using MsgPack.

InstallSnapshot is special, in that after the RPC request we stream
the entire state. That socket is not re-used as the connection state
is not known if there is an error.

func

NewNetworkTransport

¶

func NewNetworkTransport(
	stream

StreamLayer

,
	maxPool

int

,
	timeout

time

.

Duration

,
	logOutput

io

.

Writer

,
) *

NetworkTransport

NewNetworkTransport creates a new network transport with the given dialer
and listener. The maxPool controls how many connections we will pool. The
timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply
the timeout by (SnapshotSize / TimeoutScale).

func

NewNetworkTransportWithConfig

¶

added in

v1.0.0

func NewNetworkTransportWithConfig(
	config *

NetworkTransportConfig

,
) *

NetworkTransport

NewNetworkTransportWithConfig creates a new network transport with the given config struct

func

NewNetworkTransportWithLogger

¶

func NewNetworkTransportWithLogger(
	stream

StreamLayer

,
	maxPool

int

,
	timeout

time

.

Duration

,
	logger

hclog

.

Logger

,
) *

NetworkTransport

NewNetworkTransportWithLogger creates a new network transport with the given logger, dialer
and listener. The maxPool controls how many connections we will pool. The
timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply
the timeout by (SnapshotSize / TimeoutScale).

func

NewTCPTransport

¶

func NewTCPTransport(
	bindAddr

string

,
	advertise

net

.

Addr

,
	maxPool

int

,
	timeout

time

.

Duration

,
	logOutput

io

.

Writer

,
) (*

NetworkTransport

,

error

)

NewTCPTransport returns a NetworkTransport that is built on top of
a TCP streaming transport layer.

func

NewTCPTransportWithConfig

¶

added in

v1.0.0

func NewTCPTransportWithConfig(
	bindAddr

string

,
	advertise

net

.

Addr

,
	config *

NetworkTransportConfig

,
) (*

NetworkTransport

,

error

)

NewTCPTransportWithConfig returns a NetworkTransport that is built on top of
a TCP streaming transport layer, using the given config struct.

func

NewTCPTransportWithLogger

¶

func NewTCPTransportWithLogger(
	bindAddr

string

,
	advertise

net

.

Addr

,
	maxPool

int

,
	timeout

time

.

Duration

,
	logger

hclog

.

Logger

,
) (*

NetworkTransport

,

error

)

NewTCPTransportWithLogger returns a NetworkTransport that is built on top of
a TCP streaming transport layer, with log output going to the supplied Logger

func (*NetworkTransport)

AppendEntries

¶

func (n *

NetworkTransport

) AppendEntries(id

ServerID

, target

ServerAddress

, args *

AppendEntriesRequest

, resp *

AppendEntriesResponse

)

error

AppendEntries implements the Transport interface.

func (*NetworkTransport)

AppendEntriesPipeline

¶

func (n *

NetworkTransport

) AppendEntriesPipeline(id

ServerID

, target

ServerAddress

) (

AppendPipeline

,

error

)

AppendEntriesPipeline returns an interface that can be used to pipeline
AppendEntries requests.

func (*NetworkTransport)

Close

¶

func (n *

NetworkTransport

) Close()

error

Close is used to stop the network transport.

func (*NetworkTransport)

CloseStreams

¶

added in

v1.0.1

func (n *

NetworkTransport

) CloseStreams()

CloseStreams closes the current streams.

func (*NetworkTransport)

Consumer

¶

func (n *

NetworkTransport

) Consumer() <-chan

RPC

Consumer implements the Transport interface.

func (*NetworkTransport)

DecodePeer

¶

func (n *

NetworkTransport

) DecodePeer(buf []

byte

)

ServerAddress

DecodePeer implements the Transport interface.

func (*NetworkTransport)

EncodePeer

¶

func (n *

NetworkTransport

) EncodePeer(id

ServerID

, p

ServerAddress

) []

byte

EncodePeer implements the Transport interface.

func (*NetworkTransport)

InstallSnapshot

¶

func (n *

NetworkTransport

) InstallSnapshot(id

ServerID

, target

ServerAddress

, args *

InstallSnapshotRequest

, resp *

InstallSnapshotResponse

, data

io

.

Reader

)

error

InstallSnapshot implements the Transport interface.

func (*NetworkTransport)

IsShutdown

¶

func (n *

NetworkTransport

) IsShutdown()

bool

IsShutdown is used to check if the transport is shutdown.

func (*NetworkTransport)

LocalAddr

¶

func (n *

NetworkTransport

) LocalAddr()

ServerAddress

LocalAddr implements the Transport interface.

func (*NetworkTransport)

RequestPreVote

¶

added in

v1.7.0

func (n *

NetworkTransport

) RequestPreVote(id

ServerID

, target

ServerAddress

, args *

RequestPreVoteRequest

, resp *

RequestPreVoteResponse

)

error

RequestPreVote implements the Transport interface.

func (*NetworkTransport)

RequestVote

¶

func (n *

NetworkTransport

) RequestVote(id

ServerID

, target

ServerAddress

, args *

RequestVoteRequest

, resp *

RequestVoteResponse

)

error

RequestVote implements the Transport interface.

func (*NetworkTransport)

SetHeartbeatHandler

¶

func (n *

NetworkTransport

) SetHeartbeatHandler(cb func(rpc

RPC

))

SetHeartbeatHandler is used to set up a heartbeat handler
as a fast-pass. This is to avoid head-of-line blocking from
disk IO.

func (*NetworkTransport)

TimeoutNow

¶

added in

v1.1.0

func (n *

NetworkTransport

) TimeoutNow(id

ServerID

, target

ServerAddress

, args *

TimeoutNowRequest

, resp *

TimeoutNowResponse

)

error

TimeoutNow implements the Transport interface.

type

NetworkTransportConfig

¶

added in

v1.0.0

type NetworkTransportConfig struct {

// ServerAddressProvider is used to override the target address when establishing a connection to invoke an RPC

ServerAddressProvider

ServerAddressProvider

Logger

hclog

.

Logger

// Dialer

Stream

StreamLayer

// MaxPool controls how many connections we will pool

MaxPool

int

// MaxRPCsInFlight controls the pipelining "optimization" when replicating

// entries to followers.

//

// Setting this to 1 explicitly disables pipelining since no overlapping of

// request processing is allowed. If set to 1 the pipelining code path is

// skipped entirely and every request is entirely synchronous.

//

// If zero is set (or left as default), DefaultMaxRPCsInFlight is used which

// is currently 2. A value of 2 overlaps the preparation and sending of the

// next request while waiting for the previous response, but avoids additional

// queuing.

//

// Historically this was internally fixed at (effectively) 130 however

// performance testing has shown that in practice the pipelining optimization

// combines badly with batching and actually has a very large negative impact

// on commit latency when throughput is high, whilst having very little

// benefit on latency or throughput in any other case! See

// [#541](

https://github.com/hashicorp/raft/pull/541

) for more analysis of the

// performance impacts.

//

// Increasing this beyond 2 is likely to be beneficial only in very

// high-latency network conditions. HashiCorp doesn't recommend using our own

// products this way.

//

// To maintain the behavior from before version 1.4.1 exactly, set this to

// 130. The old internal constant was 128 but was used directly as a channel

// buffer size. Since we send before blocking on the channel and unblock the

// channel as soon as the receiver is done with the earliest outstanding

// request, even an unbuffered channel (buffer=0) allows one request to be

// sent while waiting for the previous one (i.e. 2 inflight). so the old

// buffer actually allowed 130 RPCs to be inflight at once.

MaxRPCsInFlight

int

// Timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply

// the timeout by (SnapshotSize / TimeoutScale).

Timeout

time

.

Duration

// MsgpackUseNewTimeFormat when set to true, force the underlying msgpack

// codec to use the new format of time.Time when encoding (used in

// go-msgpack v1.1.5 by default). Decoding is not affected, as all

// go-msgpack v2.1.0+ decoders know how to decode both formats.

MsgpackUseNewTimeFormat

bool

}

NetworkTransportConfig encapsulates configuration for the network transport layer.

type

Observation

¶

type Observation struct {

// Raft holds the Raft instance generating the observation.

Raft *

Raft

// Data holds observation-specific data. Possible types are

// RequestVoteRequest

// RaftState

// PeerObservation

// LeaderObservation

Data interface{}
}

Observation is sent along the given channel to observers when an event occurs.

type

Observer

¶

type Observer struct {

// contains filtered or unexported fields

}

Observer describes what to do with a given observation.

func

NewObserver

¶

func NewObserver(channel chan

Observation

, blocking

bool

, filter

FilterFn

) *

Observer

NewObserver creates a new observer that can be registered
to make observations on a Raft instance. Observations
will be sent on the given channel if they satisfy the
given filter.

If blocking is true, the observer will block when it can't
send on the channel, otherwise it may discard events.

func (*Observer)

GetNumDropped

¶

func (or *

Observer

) GetNumDropped()

uint64

GetNumDropped returns the number of dropped observations due to blocking.

func (*Observer)

GetNumObserved

¶

func (or *

Observer

) GetNumObserved()

uint64

GetNumObserved returns the number of observations.

type

PeerObservation

¶

added in

v1.1.0

type PeerObservation struct {

Removed

bool

Peer

Server

}

PeerObservation is sent to observers when peers change.

type

ProtocolVersion

¶

added in

v1.0.0

type ProtocolVersion

int

Version History

ProtocolVersion is the version of the protocol (which includes RPC messages
as well as Raft-specific log entries) that this server can _understand_. Use
the ProtocolVersion member of the Config object to control the version of
the protocol to use when _speaking_ to other servers. Note that depending on
the protocol version being spoken, some otherwise understood RPC messages
may be refused. See dispositionRPC for details of this logic.

There are notes about the upgrade path in the description of the versions
below. If you are starting a fresh cluster then there's no reason not to
jump right to the latest protocol version. If you need to interoperate with
older, version 0 Raft servers you'll need to drive the cluster through the
different versions in order.

The version details are complicated, but here's a summary of what's required
to get from a version 0 cluster to version 3:

In version N of your app that starts using the new Raft library with
versioning, set ProtocolVersion to 1.

Make version N+1 of your app require version N as a prerequisite (all
servers must be upgraded). For version N+1 of your app set ProtocolVersion
to 2.

Similarly, make version N+2 of your app require version N+1 as a
prerequisite. For version N+2 of your app, set ProtocolVersion to 3.

During this upgrade, older cluster members will still have Server IDs equal
to their network addresses. To upgrade an older member and give it an ID, it
needs to leave the cluster and re-enter:

Remove the server from the cluster with RemoveServer, using its network
address as its ServerID.

Update the server's config to use a UUID or something else that is
not tied to the machine as the ServerID (restarting the server).

Add the server back to the cluster with AddVoter, using its new ID.

You can do this during the rolling upgrade from N+1 to N+2 of your app, or
as a rolling change at any time after the upgrade.

Version History

¶

0: Original Raft library before versioning was added. Servers running this

version of the Raft library use AddPeerDeprecated/RemovePeerDeprecated
for all configuration changes, and have no support for LogConfiguration.

1: First versioned protocol, used to interoperate with old servers, and begin

the migration path to newer versions of the protocol. Under this version
all configuration changes are propagated using the now-deprecated
RemovePeerDeprecated Raft log entry. This means that server IDs are always
set to be the same as the server addresses (since the old log entry type
cannot transmit an ID), and only AddPeer/RemovePeer APIs are supported.
Servers running this version of the protocol can understand the new
LogConfiguration Raft log entry but will never generate one so they can
remain compatible with version 0 Raft servers in the cluster.

2: Transitional protocol used when migrating an existing cluster to the new

server ID system. Server IDs are still set to be the same as server
addresses, but all configuration changes are propagated using the new
LogConfiguration Raft log entry type, which can carry full ID information.
This version supports the old AddPeer/RemovePeer APIs as well as the new
ID-based AddVoter/RemoveServer APIs which should be used when adding
version 3 servers to the cluster later. This version sheds all
interoperability with version 0 servers, but can interoperate with newer
Raft servers running with protocol version 1 since they can understand the
new LogConfiguration Raft log entry, and this version can still understand
their RemovePeerDeprecated Raft log entries. We need this protocol version
as an intermediate step between 1 and 3 so that servers will propagate the
ID information that will come from newly-added (or -rolled) servers using
protocol version 3, but since they are still using their address-based IDs
from the previous step they will still be able to track commitments and
their own voting status properly. If we skipped this step, servers would
be started with their new IDs, but they wouldn't see themselves in the old
address-based configuration, so none of the servers would think they had a
vote.

3: Protocol adding full support for server IDs and new ID-based server APIs

(AddVoter, AddNonvoter, etc.), old AddPeer/RemovePeer APIs are no longer
supported. Version 2 servers should be swapped out by removing them from
the cluster one-by-one and re-adding them with updated configuration for
this protocol version, along with their server ID. The remove/add cycle
is required to populate their server ID. Note that removing must be done
by ID, which will be the old server's address.

const (

// ProtocolVersionMin is the minimum protocol version

ProtocolVersionMin

ProtocolVersion

= 0

// ProtocolVersionMax is the maximum protocol version

ProtocolVersionMax = 3
)

type

RPC

¶

type RPC struct {

Command  interface{}

Reader

io

.

Reader

// Set only for InstallSnapshot

RespChan chan<-

RPCResponse

}

RPC has a command, and provides a response mechanism.

func (*RPC)

Respond

¶

func (r *

RPC

) Respond(resp interface{}, err

error

)

Respond is used to respond with a response, error or both

type

RPCHeader

¶

added in

v1.0.0

type RPCHeader struct {

// ProtocolVersion is the version of the protocol the sender is

// speaking.

ProtocolVersion

ProtocolVersion

// ID is the ServerID of the node sending the RPC Request or Response

ID []

byte

// Addr is the ServerAddr of the node sending the RPC Request or Response

Addr []

byte

}

RPCHeader is a common sub-structure used to pass along protocol version and
other information about the cluster. For older Raft implementations before
versioning was added this will default to a zero-valued structure when read
by newer Raft versions.

type

RPCResponse

¶

type RPCResponse struct {

Response interface{}

Error

error

}

RPCResponse captures both a response and a potential error.

type

Raft

¶

type Raft struct {

// contains filtered or unexported fields

}

Raft implements a Raft node.

func

NewRaft

¶

func NewRaft(conf *

Config

, fsm

FSM

, logs

LogStore

, stable

StableStore

, snaps

SnapshotStore

, trans

Transport

) (*

Raft

,

error

)

NewRaft is used to construct a new Raft node. It takes a configuration, as well
as implementations of various interfaces that are required. If we have any
old state, such as snapshots, logs, peers, etc, all those will be restored
when creating the Raft node.

func (*Raft)

AddNonvoter

¶

added in

v1.0.0

func (r *

Raft

) AddNonvoter(id

ServerID

, address

ServerAddress

, prevIndex

uint64

, timeout

time

.

Duration

)

IndexFuture

AddNonvoter will add the given server to the cluster but won't assign it a
vote. The server will receive log entries, but it won't participate in
elections or log entry commitment. If the server is already in the cluster,
this updates the server's address. This must be run on the leader or it will
fail. For prevIndex and timeout, see AddVoter.

func (*Raft)

AddPeer

deprecated

func (r *

Raft

) AddPeer(peer

ServerAddress

)

Future

AddPeer to the cluster configuration. Must be run on the leader, or it will fail.

Deprecated: Use AddVoter/AddNonvoter instead.

func (*Raft)

AddVoter

¶

added in

v1.0.0

func (r *

Raft

) AddVoter(id

ServerID

, address

ServerAddress

, prevIndex

uint64

, timeout

time

.

Duration

)

IndexFuture

AddVoter will add the given server to the cluster as a staging server. If the
server is already in the cluster as a voter, this updates the server's address.
This must be run on the leader or it will fail. The leader will promote the
staging server to a voter once that server is ready. If nonzero, prevIndex is
the index of the only configuration upon which this change may be applied; if
another configuration entry has been added in the meantime, this request will
fail. If nonzero, timeout is how long this server should wait before the
configuration change log entry is appended.

func (*Raft)

AppliedIndex

¶

func (r *

Raft

) AppliedIndex()

uint64

AppliedIndex returns the last index applied to the FSM. This is generally
lagging behind the last index, especially for indexes that are persisted but
have not yet been considered committed by the leader. NOTE - this reflects
the last index that was sent to the application's FSM over the apply channel
but DOES NOT mean that the application's FSM has yet consumed it and applied
it to its internal state. Thus, the application's state may lag behind this
index.

func (*Raft)

Apply

¶

func (r *

Raft

) Apply(cmd []

byte

, timeout

time

.

Duration

)

ApplyFuture

Apply is used to apply a command to the FSM in a highly consistent
manner. This returns a future that can be used to wait on the application.
An optional timeout can be provided to limit the amount of time we wait
for the command to be started. This must be run on the leader or it
will fail.

If the node discovers it is no longer the leader while applying the command,
it will return ErrLeadershipLost. There is no way to guarantee whether the
write succeeded or failed in this case. For example, if the leader is
partitioned it can't know if a quorum of followers wrote the log to disk. If
at least one did, it may survive into the next leader's term.

If a user snapshot is restored while the command is in-flight, an
ErrAbortedByRestore is returned. In this case the write effectively failed
since its effects will not be present in the FSM after the restore.

func (*Raft)

ApplyLog

¶

added in

v1.1.1

func (r *

Raft

) ApplyLog(log

Log

, timeout

time

.

Duration

)

ApplyFuture

ApplyLog performs Apply but takes in a Log directly. The only values
currently taken from the submitted Log are Data and Extensions. See
Apply for details on error cases.

func (*Raft)

Barrier

¶

func (r *

Raft

) Barrier(timeout

time

.

Duration

)

Future

Barrier is used to issue a command that blocks until all preceding
operations have been applied to the FSM. It can be used to ensure the
FSM reflects all queued writes. An optional timeout can be provided to
limit the amount of time we wait for the command to be started. This
must be run on the leader, or it will fail.

func (*Raft)

BootstrapCluster

¶

added in

v1.0.0

func (r *

Raft

) BootstrapCluster(configuration

Configuration

)

Future

BootstrapCluster is equivalent to non-member BootstrapCluster but can be
called on an un-bootstrapped Raft instance after it has been created. This
should only be called at the beginning of time for the cluster with an
identical configuration listing all Voter servers. There is no need to
bootstrap Nonvoter and Staging servers.

A cluster can only be bootstrapped once from a single participating Voter
server. Any further attempts to bootstrap will return an error that can be
safely ignored.

One sane approach is to bootstrap a single server with a configuration
listing just itself as a Voter, then invoke AddVoter() on it to add other
servers to the cluster.

func (*Raft)

CommitIndex

¶

added in

v1.6.0

func (r *

Raft

) CommitIndex()

uint64

CommitIndex returns the committed index.
This API maybe helpful for server to implement the read index optimization
as described in the Raft paper.

func (*Raft)

CurrentTerm

¶

added in

v1.7.2

func (r *

Raft

) CurrentTerm()

uint64

CurrentTerm returns the current term.

func (*Raft)

DemoteVoter

¶

added in

v1.0.0

func (r *

Raft

) DemoteVoter(id

ServerID

, prevIndex

uint64

, timeout

time

.

Duration

)

IndexFuture

DemoteVoter will take away a server's vote, if it has one. If present, the
server will continue to receive log entries, but it won't participate in
elections or log entry commitment. If the server is not in the cluster, this
does nothing. This must be run on the leader or it will fail. For prevIndex
and timeout, see AddVoter.

func (*Raft)

DeregisterObserver

¶

func (r *

Raft

) DeregisterObserver(or *

Observer

)

DeregisterObserver deregisters an observer.

func (*Raft)

GetConfiguration

¶

added in

v1.0.0

func (r *

Raft

) GetConfiguration()

ConfigurationFuture

GetConfiguration returns the latest configuration. This may not yet be
committed. The main loop can access this directly.

func (*Raft)

LastContact

¶

func (r *

Raft

) LastContact()

time

.

Time

LastContact returns the time of last contact by a leader.
This only makes sense if we are currently a follower.

func (*Raft)

LastIndex

¶

func (r *

Raft

) LastIndex()

uint64

LastIndex returns the last index in stable storage,
either from the last log or from the last snapshot.

func (*Raft)

Leader

¶

func (r *

Raft

) Leader()

ServerAddress

Leader is used to return the current leader of the cluster.
Deprecated: use LeaderWithID instead
It may return empty string if there is no current leader
or the leader is unknown.
Deprecated: use LeaderWithID instead.

func (*Raft)

LeaderCh

¶

func (r *

Raft

) LeaderCh() <-chan

bool

LeaderCh is used to get a channel which delivers signals on acquiring or
losing leadership. It sends true if we become the leader, and false if we
lose it.

Receivers can expect to receive a notification only if leadership
transition has occurred.

If receivers aren't ready for the signal, signals may drop and only the
latest leadership transition. For example, if a receiver receives subsequent
`true` values, they may deduce that leadership was lost and regained while
the receiver was processing first leadership transition.

func (*Raft)

LeaderWithID

¶

added in

v1.3.7

func (r *

Raft

) LeaderWithID() (

ServerAddress

,

ServerID

)

LeaderWithID is used to return the current leader address and ID of the cluster.
It may return empty strings if there is no current leader
or the leader is unknown.

func (*Raft)

LeadershipTransfer

¶

added in

v1.1.0

func (r *

Raft

) LeadershipTransfer()

Future

LeadershipTransfer will transfer leadership to a server in the cluster.
This can only be called from the leader, or it will fail. The leader will
stop accepting client requests, make sure the target server is up to date
and starts the transfer with a TimeoutNow message. This message has the same
effect as if the election timeout on the target server fires. Since
it is unlikely that another server is starting an election, it is very
likely that the target server is able to win the election.  Note that raft
protocol version 3 is not sufficient to use LeadershipTransfer. A recent
version of that library has to be used that includes this feature.  Using
transfer leadership is safe however in a cluster where not every node has
the latest version. If a follower cannot be promoted, it will fail
gracefully.

func (*Raft)

LeadershipTransferToServer

¶

added in

v1.1.0

func (r *

Raft

) LeadershipTransferToServer(id

ServerID

, address

ServerAddress

)

Future

LeadershipTransferToServer does the same as LeadershipTransfer but takes a
server in the arguments in case a leadership should be transitioned to a
specific server in the cluster.  Note that raft protocol version 3 is not
sufficient to use LeadershipTransfer. A recent version of that library has
to be used that includes this feature. Using transfer leadership is safe
however in a cluster where not every node has the latest version. If a
follower cannot be promoted, it will fail gracefully.

func (*Raft)

RegisterObserver

¶

func (r *

Raft

) RegisterObserver(or *

Observer

)

RegisterObserver registers a new observer.

func (*Raft)

ReloadConfig

¶

added in

v1.3.0

func (r *

Raft

) ReloadConfig(rc

ReloadableConfig

)

error

ReloadConfig updates the configuration of a running raft node. If the new
configuration is invalid an error is returned and no changes made to the
instance. All fields will be copied from rc into the new configuration, even
if they are zero valued.

func (*Raft)

ReloadableConfig

¶

added in

v1.3.0

func (r *

Raft

) ReloadableConfig()

ReloadableConfig

ReloadableConfig returns the current state of the reloadable fields in Raft's
configuration. This is useful for programs to discover the current state for
reporting to users or tests. It is safe to call from any goroutine. It is
intended for reporting and testing purposes primarily; external
synchronization would be required to safely use this in a read-modify-write
pattern for reloadable configuration options.

func (*Raft)

RemovePeer

deprecated

func (r *

Raft

) RemovePeer(peer

ServerAddress

)

Future

Deprecated: Use RemoveServer instead.

func (*Raft)

RemoveServer

¶

added in

v1.0.0

func (r *

Raft

) RemoveServer(id

ServerID

, prevIndex

uint64

, timeout

time

.

Duration

)

IndexFuture

RemoveServer will remove the given server from the cluster. If the current
leader is being removed, it will cause a new election to occur. This must be
run on the leader or it will fail. For prevIndex and timeout, see AddVoter.

func (*Raft)

Restore

¶

added in

v1.0.0

func (r *

Raft

) Restore(meta *

SnapshotMeta

, reader

io

.

Reader

, timeout

time

.

Duration

)

error

Restore is used to manually force Raft to consume an external snapshot, such
as if restoring from a backup. We will use the current Raft configuration,
not the one from the snapshot, so that we can restore into a new cluster. We
will also use the max of the index of the snapshot, or the current index,
and then add 1 to that, so we force a new state with a hole in the Raft log,
so that the snapshot will be sent to followers and used for any new joiners.
This can only be run on the leader, and blocks until the restore is complete
or an error occurs.

WARNING! This operation has the leader take on the state of the snapshot and
then sets itself up so that it replicates that to its followers though the
install snapshot process. This involves a potentially dangerous period where
the leader commits ahead of its followers, so should only be used for disaster
recovery into a fresh cluster, and should not be used in normal operations.

func (*Raft)

Shutdown

¶

func (r *

Raft

) Shutdown()

Future

Shutdown is used to stop the Raft background routines.
This is not a graceful operation. Provides a future that
can be used to block until all background routines have exited.

func (*Raft)

Snapshot

¶

func (r *

Raft

) Snapshot()

SnapshotFuture

Snapshot is used to manually force Raft to take a snapshot. Returns a future
that can be used to block until complete, and that contains a function that
can be used to open the snapshot.

func (*Raft)

State

¶

func (r *

Raft

) State()

RaftState

State returns the state of this raft peer.

func (*Raft)

Stats

¶

func (r *

Raft

) Stats() map[

string

]

string

Stats is used to return a map of various internal stats. This
should only be used for informative purposes or debugging.

Keys are: "state", "term", "last_log_index", "last_log_term",
"commit_index", "applied_index", "fsm_pending",
"last_snapshot_index", "last_snapshot_term",
"latest_configuration", "last_contact", and "num_peers".

The value of "state" is a numeric constant representing one of
the possible leadership states the node is in at any given time.
the possible states are: "Follower", "Candidate", "Leader", "Shutdown".

The value of "latest_configuration" is a string which contains
the id of each server, its suffrage status, and its address.

The value of "last_contact" is either "never" if there
has been no contact with a leader, "0" if the node is in the
leader state, or the time since last contact with a leader
formatted as a string.

The value of "num_peers" is the number of other voting servers in the
cluster, not including this node. If this node isn't part of the
configuration then this will be "0".

All other values are uint64s, formatted as strings.

func (*Raft)

String

¶

func (r *

Raft

) String()

string

String returns a string representation of this Raft node.

func (*Raft)

VerifyLeader

¶

func (r *

Raft

) VerifyLeader()

Future

VerifyLeader is used to ensure this peer is still the leader. It may be used
to prevent returning stale data from the FSM after the peer has lost
leadership.

type

RaftState

¶

type RaftState

uint32

RaftState captures the state of a Raft node: Follower, Candidate, Leader,
or Shutdown.

const (

// Follower is the initial state of a Raft node.

Follower

RaftState

=

iota

// Candidate is one of the valid states of a Raft node.

Candidate

// Leader is one of the valid states of a Raft node.

Leader

// Shutdown is the terminal state of a Raft node.

Shutdown
)

func (RaftState)

String

¶

func (s

RaftState

) String()

string

type

ReadCloserWrapper

¶

added in

v1.3.9

type ReadCloserWrapper interface {

io

.

ReadCloser

WrappedReadCloser()

io

.

ReadCloser

}

ReadCloserWrapper allows access to an underlying ReadCloser from a wrapper.

type

ReloadableConfig

¶

added in

v1.3.0

type ReloadableConfig struct {

// TrailingLogs controls how many logs we leave after a snapshot. This is used

// so that we can quickly replay logs on a follower instead of being forced to

// send an entire snapshot. The value passed here updates the setting at runtime

// which will take effect as soon as the next snapshot completes and truncation

// occurs.

TrailingLogs

uint64

// SnapshotInterval controls how often we check if we should perform a snapshot.

// We randomly stagger between this value and 2x this value to avoid the entire

// cluster from performing a snapshot at once.

SnapshotInterval

time

.

Duration

// SnapshotThreshold controls how many outstanding logs there must be before

// we perform a snapshot. This is to prevent excessive snapshots when we can

// just replay a small set of logs.

SnapshotThreshold

uint64

// HeartbeatTimeout specifies the time in follower state without

// a leader before we attempt an election.

HeartbeatTimeout

time

.

Duration

// ElectionTimeout specifies the time in candidate state without

// a leader before we attempt an election.

ElectionTimeout

time

.

Duration

}

ReloadableConfig is the subset of Config that may be reconfigured during
runtime using raft.ReloadConfig. We choose to duplicate fields over embedding
or accepting a Config but only using specific fields to keep the API clear.
Reconfiguring some fields is potentially dangerous so we should only
selectively enable it for fields where that is allowed.

type

RequestPreVoteRequest

¶

added in

v1.7.0

type RequestPreVoteRequest struct {

RPCHeader

// Provide the term and our id

Term

uint64

// Used to ensure safety

LastLogIndex

uint64

LastLogTerm

uint64

}

RequestPreVoteRequest is the command used by a candidate to ask a Raft peer
for a vote in an election.

func (*RequestPreVoteRequest)

GetRPCHeader

¶

added in

v1.7.0

func (r *

RequestPreVoteRequest

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

RequestPreVoteResponse

¶

added in

v1.7.0

type RequestPreVoteResponse struct {

RPCHeader

// Newer term if leader is out of date.

Term

uint64

// Is the vote granted.

Granted

bool

}

RequestPreVoteResponse is the response returned from a RequestPreVoteRequest.

func (*RequestPreVoteResponse)

GetRPCHeader

¶

added in

v1.7.0

func (r *

RequestPreVoteResponse

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

RequestVoteRequest

¶

type RequestVoteRequest struct {

RPCHeader

// Provide the term and our id

Term

uint64

// Deprecated: use RPCHeader.Addr instead

Candidate []

byte

// Used to ensure safety

LastLogIndex

uint64

LastLogTerm

uint64

// Used to indicate to peers if this vote was triggered by a leadership

// transfer. It is required for leadership transfer to work, because servers

// wouldn't vote otherwise if they are aware of an existing leader.

LeadershipTransfer

bool

}

RequestVoteRequest is the command used by a candidate to ask a Raft peer
for a vote in an election.

func (*RequestVoteRequest)

GetRPCHeader

¶

added in

v1.0.0

func (r *

RequestVoteRequest

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

RequestVoteResponse

¶

type RequestVoteResponse struct {

RPCHeader

// Newer term if leader is out of date.

Term

uint64

// Peers is deprecated, but required by servers that only understand

// protocol version 0. This is not populated in protocol version 2

// and later.

Peers []

byte

// Is the vote granted.

Granted

bool

}

RequestVoteResponse is the response returned from a RequestVoteRequest.

func (*RequestVoteResponse)

GetRPCHeader

¶

added in

v1.0.0

func (r *

RequestVoteResponse

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

ResumedHeartbeatObservation

¶

added in

v1.3.2

type ResumedHeartbeatObservation struct {

PeerID

ServerID

}

ResumedHeartbeatObservation is sent when a node resumes to heartbeat with the leader following failures

type

Server

¶

added in

v1.0.0

type Server struct {

// Suffrage determines whether the server gets a vote.

Suffrage

ServerSuffrage

// ID is a unique string identifying this server for all time.

ID

ServerID

// Address is its network address that a transport can contact.

Address

ServerAddress

}

Server tracks the information about a single server in a configuration.

type

ServerAddress

¶

added in

v1.0.0

type ServerAddress

string

ServerAddress is a network address for a server that a transport can contact.

func

NewInmemAddr

¶

func NewInmemAddr()

ServerAddress

NewInmemAddr returns a new in-memory addr with
a randomly generate UUID as the ID.

type

ServerAddressProvider

¶

added in

v1.0.0

type ServerAddressProvider interface {

ServerAddr(id

ServerID

) (

ServerAddress

,

error

)

}

ServerAddressProvider is a target address to which we invoke an RPC when establishing a connection

type

ServerID

¶

added in

v1.0.0

type ServerID

string

ServerID is a unique string identifying a server for all time.

type

ServerSuffrage

¶

added in

v1.0.0

type ServerSuffrage

int

ServerSuffrage determines whether a Server in a Configuration gets a vote.

const (

// Voter is a server whose vote is counted in elections and whose match index

// is used in advancing the leader's commit index.

Voter

ServerSuffrage

=

iota

// Nonvoter is a server that receives log entries but is not considered for

// elections or commitment purposes.

Nonvoter

// Staging is a server that acts like a Nonvoter. A configuration change

// with a ConfigurationChangeCommand of Promote can change a Staging server

// into a Voter.

// Deprecated: use Nonvoter instead.

Staging
)

Note: Don't renumber these, since the numbers are written into the log.

func (ServerSuffrage)

String

¶

added in

v1.0.0

func (s

ServerSuffrage

) String()

string

type

SnapshotFuture

¶

added in

v1.0.0

type SnapshotFuture interface {

Future

// Open is a function you can call to access the underlying snapshot and

// its metadata. This must not be called until after the Error method

// has returned.

Open() (*

SnapshotMeta

,

io

.

ReadCloser

,

error

)
}

SnapshotFuture is used for waiting on a user-triggered snapshot to complete.

type

SnapshotMeta

¶

type SnapshotMeta struct {

// Version is the version number of the snapshot metadata. This does not cover

// the application's data in the snapshot, that should be versioned

// separately.

Version

SnapshotVersion

// ID is opaque to the store, and is used for opening.

ID

string

// Index and Term store when the snapshot was taken.

Index

uint64

Term

uint64

// Peers is deprecated and used to support version 0 snapshots, but will

// be populated in version 1 snapshots as well to help with upgrades.

Peers []

byte

// Configuration and ConfigurationIndex are present in version 1

// snapshots and later.

Configuration

Configuration

ConfigurationIndex

uint64

// Size is the size of the snapshot in bytes.

Size

int64

}

SnapshotMeta is for metadata of a snapshot.

type

SnapshotSink

¶

type SnapshotSink interface {

io

.

WriteCloser

ID()

string

Cancel()

error

}

SnapshotSink is returned by StartSnapshot. The FSM will Write state
to the sink and call Close on completion. On error, Cancel will be invoked.

type

SnapshotStore

¶

type SnapshotStore interface {

// Create is used to begin a snapshot at a given index and term, and with

// the given committed configuration. The version parameter controls

// which snapshot version to create.

Create(version

SnapshotVersion

, index, term

uint64

, configuration

Configuration

,
		configurationIndex

uint64

, trans

Transport

) (

SnapshotSink

,

error

)

// List is used to list the available snapshots in the store.

// It should return then in descending order, with the highest index first.

List() ([]*

SnapshotMeta

,

error

)

// Open takes a snapshot ID and provides a ReadCloser. Once close is

// called it is assumed the snapshot is no longer needed.

Open(id

string

) (*

SnapshotMeta

,

io

.

ReadCloser

,

error

)
}

SnapshotStore interface is used to allow for flexible implementations
of snapshot storage and retrieval. For example, a client could implement
a shared state store such as S3, allowing new nodes to restore snapshots
without streaming from the leader.

type

SnapshotVersion

¶

added in

v1.0.0

type SnapshotVersion

int

Version History

SnapshotVersion is the version of snapshots that this server can understand.
Currently, it is always assumed that the server generates the latest version,
though this may be changed in the future to include a configurable version.

Version History

¶

0: Original Raft library before versioning was added. The peers portion of

these snapshots is encoded in the legacy format which requires decodePeers
to parse. This version of snapshots should only be produced by the
unversioned Raft library.

1: New format which adds support for a full configuration structure and its

associated log index, with support for server IDs and non-voting server
modes. To ease upgrades, this also includes the legacy peers structure but
that will never be used by servers that understand version 1 snapshots.
Since the original Raft library didn't enforce any versioning, we must
include the legacy peers structure for this version, but we can deprecate
it in the next snapshot version.

const (

// SnapshotVersionMin is the minimum snapshot version

SnapshotVersionMin

SnapshotVersion

= 0

// SnapshotVersionMax is the maximum snapshot version

SnapshotVersionMax = 1
)

type

StableStore

¶

type StableStore interface {

Set(key []

byte

, val []

byte

)

error

// Get returns the value for key, or an empty byte slice if key was not found.

Get(key []

byte

) ([]

byte

,

error

)

SetUint64(key []

byte

, val

uint64

)

error

// GetUint64 returns the uint64 value for key, or 0 if key was not found.

GetUint64(key []

byte

) (

uint64

,

error

)
}

StableStore is used to provide stable storage
of key configurations to ensure safety.

type

StreamLayer

¶

type StreamLayer interface {

net

.

Listener

// Dial is used to create a new outgoing connection

Dial(address

ServerAddress

, timeout

time

.

Duration

) (

net

.

Conn

,

error

)
}

StreamLayer is used with the NetworkTransport to provide
the low level stream abstraction.

type

TCPStreamLayer

¶

type TCPStreamLayer struct {

// contains filtered or unexported fields

}

TCPStreamLayer implements StreamLayer interface for plain TCP.

func (*TCPStreamLayer)

Accept

¶

func (t *

TCPStreamLayer

) Accept() (c

net

.

Conn

, err

error

)

Accept implements the net.Listener interface.

func (*TCPStreamLayer)

Addr

¶

func (t *

TCPStreamLayer

) Addr()

net

.

Addr

Addr implements the net.Listener interface.

func (*TCPStreamLayer)

Close

¶

func (t *

TCPStreamLayer

) Close() (err

error

)

Close implements the net.Listener interface.

func (*TCPStreamLayer)

Dial

¶

func (t *

TCPStreamLayer

) Dial(address

ServerAddress

, timeout

time

.

Duration

) (

net

.

Conn

,

error

)

Dial implements the StreamLayer interface.

type

TimeoutNowRequest

¶

added in

v1.1.0

type TimeoutNowRequest struct {

RPCHeader

}

TimeoutNowRequest is the command used by a leader to signal another server to
start an election.

func (*TimeoutNowRequest)

GetRPCHeader

¶

added in

v1.1.0

func (r *

TimeoutNowRequest

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

TimeoutNowResponse

¶

added in

v1.1.0

type TimeoutNowResponse struct {

RPCHeader

}

TimeoutNowResponse is the response to TimeoutNowRequest.

func (*TimeoutNowResponse)

GetRPCHeader

¶

added in

v1.1.0

func (r *

TimeoutNowResponse

) GetRPCHeader()

RPCHeader

GetRPCHeader - See WithRPCHeader.

type

Transport

¶

type Transport interface {

// Consumer returns a channel that can be used to

// consume and respond to RPC requests.

Consumer() <-chan

RPC

// LocalAddr is used to return our local address to distinguish from our peers.

LocalAddr()

ServerAddress

// AppendEntriesPipeline returns an interface that can be used to pipeline

// AppendEntries requests.

AppendEntriesPipeline(id

ServerID

, target

ServerAddress

) (

AppendPipeline

,

error

)

// AppendEntries sends the appropriate RPC to the target node.

AppendEntries(id

ServerID

, target

ServerAddress

, args *

AppendEntriesRequest

, resp *

AppendEntriesResponse

)

error

// RequestVote sends the appropriate RPC to the target node.

RequestVote(id

ServerID

, target

ServerAddress

, args *

RequestVoteRequest

, resp *

RequestVoteResponse

)

error

// InstallSnapshot is used to push a snapshot down to a follower. The data is read from

// the ReadCloser and streamed to the client.

InstallSnapshot(id

ServerID

, target

ServerAddress

, args *

InstallSnapshotRequest

, resp *

InstallSnapshotResponse

, data

io

.

Reader

)

error

// EncodePeer is used to serialize a peer's address.

EncodePeer(id

ServerID

, addr

ServerAddress

) []

byte

// DecodePeer is used to deserialize a peer's address.

DecodePeer([]

byte

)

ServerAddress

// SetHeartbeatHandler is used to setup a heartbeat handler

// as a fast-pass. This is to avoid head-of-line blocking from

// disk IO. If a Transport does not support this, it can simply

// ignore the call, and push the heartbeat onto the Consumer channel.

SetHeartbeatHandler(cb func(rpc

RPC

))

// TimeoutNow is used to start a leadership transfer to the target node.

TimeoutNow(id

ServerID

, target

ServerAddress

, args *

TimeoutNowRequest

, resp *

TimeoutNowResponse

)

error

}

Transport provides an interface for network transports
to allow Raft to communicate with other nodes.

type

WithClose

¶

type WithClose interface {

// Close permanently closes a transport, stopping

// any associated goroutines and freeing other resources.

Close()

error

}

WithClose is an interface that a transport may provide which
allows a transport to be shut down cleanly when a Raft instance
shuts down.

It is defined separately from Transport as unfortunately it wasn't in the
original interface specification.

type

WithPeers

¶

type WithPeers interface {

Connect(peer

ServerAddress

, t

Transport

)

// Connect a peer

Disconnect(peer

ServerAddress

)

// Disconnect a given peer

DisconnectAll()

// Disconnect all peers, possibly to reconnect them later

}

WithPeers is an interface that a transport may provide which allows for connection and
disconnection. Unless the transport is a loopback transport, the transport specified to
"Connect" is likely to be nil.

type

WithPreVote

¶

added in

v1.7.0

type WithPreVote interface {

// RequestPreVote sends the appropriate RPC to the target node.

RequestPreVote(id

ServerID

, target

ServerAddress

, args *

RequestPreVoteRequest

, resp *

RequestPreVoteResponse

)

error

}

WithPreVote is an interface that a transport may provide which
allows a transport to support a PreVote request.

It is defined separately from Transport as unfortunately it wasn't in the
original interface specification.

type

WithRPCHeader

¶

added in

v1.0.0

type WithRPCHeader interface {

GetRPCHeader()

RPCHeader

}

WithRPCHeader is an interface that exposes the RPC header.

type

WrappingFSM

¶

added in

v1.1.1

type WrappingFSM interface {

Underlying()

FSM

}

NOTE: This is exposed for middleware testing purposes and is not a stable API