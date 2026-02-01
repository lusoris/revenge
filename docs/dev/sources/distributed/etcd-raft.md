# etcd Raft

> Source: https://pkg.go.dev/go.etcd.io/raft/v3
> Fetched: 2026-02-01T11:50:39.885070+00:00
> Content-Hash: 4e8473d5951d18e6
> Type: html

---

### Overview ¶

- Usage
- Usage with Asynchronous Storage Writes
- Implementation notes
- MessageType

Package raft sends and receives messages in the Protocol Buffer format defined in the raftpb package.

Raft is a protocol with which a cluster of nodes can maintain a replicated state machine. The state machine is kept in sync through the use of a replicated log. For more details on Raft, see "In Search of an Understandable Consensus Algorithm" (<https://raft.github.io/raft.pdf>) by Diego Ongaro and John Ousterhout.

A simple example application, _raftexample_, is also available to help illustrate how to use this package in practice: <https://github.com/etcd-io/etcd/tree/main/contrib/raftexample>

#### Usage ¶

The primary object in raft is a Node. You either start a Node from scratch using raft.StartNode or start a Node from some initial state using raft.RestartNode.

To start a node from scratch:

    storage := raft.NewMemoryStorage()
    c := &Config{
      ID:              0x01,
      ElectionTick:    10,
      HeartbeatTick:   1,
      Storage:         storage,
      MaxSizePerMsg:   4096,
      MaxInflightMsgs: 256,
    }
    n := raft.StartNode(c, []raft.Peer{{ID: 0x02}, {ID: 0x03}})
    

To restart a node from previous state:

    storage := raft.NewMemoryStorage()
    
    // recover the in-memory storage from persistent
    // snapshot, state and entries.
    storage.ApplySnapshot(snapshot)
    storage.SetHardState(state)
    storage.Append(entries)
    
    c := &Config{
      ID:              0x01,
      ElectionTick:    10,
      HeartbeatTick:   1,
      Storage:         storage,
      MaxSizePerMsg:   4096,
      MaxInflightMsgs: 256,
    }
    
    // restart raft without peer information.
    // peer information is already included in the storage.
    n := raft.RestartNode(c)
    

Now that you are holding onto a Node you have a few responsibilities:

First, you must read from the Node.Ready() channel and process the updates it contains. These steps may be performed in parallel, except as noted in step 2.

1. Write HardState, Entries, and Snapshot to persistent storage if they are not empty. Note that when writing an Entry with Index i, any previously-persisted entries with Index >= i must be discarded.

2. Send all Messages to the nodes named in the To field. It is important that no messages be sent until the latest HardState has been persisted to disk, and all Entries written by any previous Ready batch (Messages may be sent while entries from the same batch are being persisted). To reduce the I/O latency, an optimization can be applied to make leader write to disk in parallel with its followers (as explained at section 10.2.1 in Raft thesis). If any Message has type MsgSnap, call Node.ReportSnapshot() after it has been sent (these messages may be large).

Note: Marshalling messages is not thread-safe; it is important that you make sure that no new entries are persisted while marshalling. The easiest way to achieve this is to serialize the messages directly inside your main raft loop.

1. Apply Snapshot (if any) and CommittedEntries to the state machine. If any committed Entry has Type EntryConfChange, call Node.ApplyConfChange() to apply it to the node. The configuration change may be cancelled at this point by setting the NodeID field to zero before calling ApplyConfChange (but ApplyConfChange must be called one way or the other, and the decision to cancel must be based solely on the state machine and not external information such as the observed health of the node).

2. Call Node.Advance() to signal readiness for the next batch of updates. This may be done at any time after step 1, although all updates must be processed in the order they were returned by Ready.

Second, all persisted log entries must be made available via an implementation of the Storage interface. The provided MemoryStorage type can be used for this (if you repopulate its state upon a restart), or you can supply your own disk-backed implementation.

Third, when you receive a message from another node, pass it to Node.Step:

    func recvRaftRPC(ctx context.Context, m raftpb.Message) {
     n.Step(ctx, m)
    }
    

Finally, you need to call Node.Tick() at regular intervals (probably via a time.Ticker). Raft has two important timeouts: heartbeat and the election timeout. However, internally to the raft package time is represented by an abstract "tick".

The total state machine handling loop will look something like this:

    for {
      select {
      case <-s.Ticker:
        n.Tick()
      case rd := <-s.Node.Ready():
        saveToStorage(rd.State, rd.Entries, rd.Snapshot)
        send(rd.Messages)
        if !raft.IsEmptySnap(rd.Snapshot) {
          processSnapshot(rd.Snapshot)
        }
        for _, entry := range rd.CommittedEntries {
          process(entry)
          if entry.Type == raftpb.EntryConfChange {
            var cc raftpb.ConfChange
            cc.Unmarshal(entry.Data)
            s.Node.ApplyConfChange(cc)
          }
        }
        s.Node.Advance()
      case <-s.done:
        return
      }
    }
    

To propose changes to the state machine from your node take your application data, serialize it into a byte slice and call:

    n.Propose(ctx, data)
    

If the proposal is committed, data will appear in committed entries with type raftpb.EntryNormal. There is no guarantee that a proposed command will be committed; you may have to re-propose after a timeout.

To add or remove a node in a cluster, build ConfChange struct 'cc' and call:

    n.ProposeConfChange(ctx, cc)
    

After config change is committed, some committed entry with type raftpb.EntryConfChange will be returned. You must apply it to node through:

    var cc raftpb.ConfChange
    cc.Unmarshal(data)
    n.ApplyConfChange(cc)
    

Note: An ID represents a unique node in a cluster for all time. A given ID MUST be used only once even if the old node has been removed. This means that for example IP addresses make poor node IDs since they may be reused. Node IDs must be non-zero.

#### Usage with Asynchronous Storage Writes ¶

The library can be configured with an alternate interface for local storage writes that can provide better performance in the presence of high proposal concurrency by minimizing interference between proposals. This feature is called AsynchronousStorageWrites, and can be enabled using the flag on the Config struct with the same name.

When Asynchronous Storage Writes is enabled, the responsibility of code using the library is different from what was presented above. Users still read from the Node.Ready() channel. However, they process the updates it contains in a different manner. Users no longer consult the HardState, Entries, and Snapshot fields (steps 1 and 3 above). They also no longer call Node.Advance() to indicate that they have processed all entries in the Ready (step 4 above). Instead, all local storage operations are also communicated through messages present in the Ready.Message slice.

The local storage messages come in two flavors. The first flavor is log append messages, which target a LocalAppendThread and carry Entries, HardState, and a Snapshot. The second flavor is entry application messages, which target a LocalApplyThread and carry CommittedEntries. Messages to the same target must be reliably processed in order. Messages to different targets can be processed in any order.

Each local storage message carries a slice of response messages that must delivered after the corresponding storage write has been completed. These responses may target the same node or may target other nodes.

With Asynchronous Storage Writes enabled, the total state machine handling loop will look something like this:

    for {
      select {
      case <-s.Ticker:
        n.Tick()
      case rd := <-s.Node.Ready():
        for _, m := range rd.Messages {
          switch m.To {
          case raft.LocalAppendThread:
            toAppend <- m
          case raft.LocalApplyThread:
            toApply <-m
          default:
            sendOverNetwork(m)
          }
        }
      case <-s.done:
        return
      }
    }
    

Usage of Asynchronous Storage Writes will typically also contain a pair of storage handler threads, one for log writes (append) and one for entry application to the local state machine (apply). Those will look something like:

    // append thread
    go func() {
      for {
        select {
        case m := <-toAppend:
          saveToStorage(m.State, m.Entries, m.Snapshot)
          send(m.Responses)
        case <-s.done:
          return
        }
      }
    }
    
    // apply thread
    go func() {
      for {
        select {
        case m := <-toApply:
          for _, entry := range m.CommittedEntries {
            process(entry)
            if entry.Type == raftpb.EntryConfChange {
              var cc raftpb.ConfChange
              cc.Unmarshal(entry.Data)
              s.Node.ApplyConfChange(cc)
            }
          }
          send(m.Responses)
        case <-s.done:
          return
        }
      }
    }
    

#### Implementation notes ¶

This implementation is up to date with the final Raft thesis (<https://github.com/ongardie/dissertation/blob/master/stanford.pdf>), although our implementation of the membership change protocol differs somewhat from that described in chapter 4. The key invariant that membership changes happen one node at a time is preserved, but in our implementation the membership change takes effect when its entry is applied, not when it is added to the log (so the entry is committed under the old membership instead of the new). This is equivalent in terms of safety, since the old and new configurations are guaranteed to overlap.

To ensure that we do not attempt to commit two membership changes at once by matching log positions (which would be unsafe since they should have different quorum requirements), we simply disallow any proposed membership change while any uncommitted change appears in the leader's log.

This approach introduces a problem when you try to remove a member from a two-member cluster: If one of the members dies before the other one receives the commit of the confchange entry, then the member cannot be removed any more since the cluster cannot make progress. For this reason it is highly recommended to use three or more nodes in every cluster.

#### MessageType ¶

Package raft sends and receives message in Protocol Buffer format (defined in raftpb package). Each state (follower, candidate, leader) implements its own 'step' method ('stepFollower', 'stepCandidate', 'stepLeader') when advancing with the given raftpb.Message. Each step is determined by its raftpb.MessageType. Note that every step is checked by one common method 'Step' that safety-checks the terms of node and incoming message to prevent stale log entries:

    'MsgHup' is used for election. If a node is a follower or candidate, the
    'tick' function in 'raft' struct is set as 'tickElection'. If a follower or
    candidate has not received any heartbeat before the election timeout, it
    passes 'MsgHup' to its Step method and becomes (or remains) a candidate to
    start a new election.
    
    'MsgBeat' is an internal type that signals the leader to send a heartbeat of
    the 'MsgHeartbeat' type. If a node is a leader, the 'tick' function in
    the 'raft' struct is set as 'tickHeartbeat', and triggers the leader to
    send periodic 'MsgHeartbeat' messages to its followers.
    
    'MsgProp' proposes to append data to its log entries. This is a special
    type to redirect proposals to leader. Therefore, send method overwrites
    raftpb.Message's term with its HardState's term to avoid attaching its
    local term to 'MsgProp'. When 'MsgProp' is passed to the leader's 'Step'
    method, the leader first calls the 'appendEntry' method to append entries
    to its log, and then calls 'bcastAppend' method to send those entries to
    its peers. When passed to candidate, 'MsgProp' is dropped. When passed to
    follower, 'MsgProp' is stored in follower's mailbox(msgs) by the send
    method. It is stored with sender's ID and later forwarded to leader by
    rafthttp package.
    
    'MsgApp' contains log entries to replicate. A leader calls bcastAppend,
    which calls sendAppend, which sends soon-to-be-replicated logs in 'MsgApp'
    type. When 'MsgApp' is passed to candidate's Step method, candidate reverts
    back to follower, because it indicates that there is a valid leader sending
    'MsgApp' messages. Candidate and follower respond to this message in
    'MsgAppResp' type.
    
    'MsgAppResp' is response to log replication request('MsgApp'). When
    'MsgApp' is passed to candidate or follower's Step method, it responds by
    calling 'handleAppendEntries' method, which sends 'MsgAppResp' to raft
    mailbox.
    
    'MsgVote' requests votes for election. When a node is a follower or
    candidate and 'MsgHup' is passed to its Step method, then the node calls
    'campaign' method to campaign itself to become a leader. Once 'campaign'
    method is called, the node becomes candidate and sends 'MsgVote' to peers
    in cluster to request votes. When passed to leader or candidate's Step
    method and the message's Term is lower than leader's or candidate's,
    'MsgVote' will be rejected ('MsgVoteResp' is returned with Reject true).
    If leader or candidate receives 'MsgVote' with higher term, it will revert
    back to follower. When 'MsgVote' is passed to follower, it votes for the
    sender only when sender's last term is greater than MsgVote's term or
    sender's last term is equal to MsgVote's term but sender's last committed
    index is greater than or equal to follower's.
    
    'MsgVoteResp' contains responses from voting request. When 'MsgVoteResp' is
    passed to candidate, the candidate calculates how many votes it has won. If
    it's more than majority (quorum), it becomes leader and calls 'bcastAppend'.
    If candidate receives majority of votes of denials, it reverts back to
    follower.
    
    'MsgPreVote' and 'MsgPreVoteResp' are used in an optional two-phase election
    protocol. When Config.PreVote is true, a pre-election is carried out first
    (using the same rules as a regular election), and no node increases its term
    number unless the pre-election indicates that the campaigning node would win.
    This minimizes disruption when a partitioned node rejoins the cluster.
    
    'MsgSnap' requests to install a snapshot message. When a node has just
    become a leader or the leader receives 'MsgProp' message, it calls
    'bcastAppend' method, which then calls 'sendAppend' method to each
    follower. In 'sendAppend', if a leader fails to get term or entries,
    the leader requests snapshot by sending 'MsgSnap' type message.
    
    'MsgSnapStatus' tells the result of snapshot install message. When a
    follower rejected 'MsgSnap', it indicates the snapshot request with
    'MsgSnap' had failed from network issues which causes the network layer
    to fail to send out snapshots to its followers. Then leader considers
    follower's progress as probe. When 'MsgSnap' were not rejected, it
    indicates that the snapshot succeeded and the leader sets follower's
    progress to probe and resumes its log replication.
    
    'MsgHeartbeat' sends heartbeat from leader. When 'MsgHeartbeat' is passed
    to candidate and message's term is higher than candidate's, the candidate
    reverts back to follower and updates its committed index from the one in
    this heartbeat. And it sends the message to its mailbox. When
    'MsgHeartbeat' is passed to follower's Step method and message's term is
    higher than follower's, the follower updates its leaderID with the ID
    from the message.
    
    'MsgHeartbeatResp' is a response to 'MsgHeartbeat'. When 'MsgHeartbeatResp'
    is passed to leader's Step method, the leader knows which follower
    responded. And only when the leader's last committed index is greater than
    follower's Match index, the leader runs 'sendAppend` method.
    
    'MsgUnreachable' tells that request(message) wasn't delivered. When
    'MsgUnreachable' is passed to leader's Step method, the leader discovers
    that the follower that sent this 'MsgUnreachable' is not reachable, often
    indicating 'MsgApp' is lost. When follower's progress state is replicate,
    the leader sets it back to probe.
    
    'MsgStorageAppend' is a message from a node to its local append storage
    thread to write entries, hard state, and/or a snapshot to stable storage.
    The message will carry one or more responses, one of which will be a
    'MsgStorageAppendResp' back to itself. The responses can also contain
    'MsgAppResp', 'MsgVoteResp', and 'MsgPreVoteResp' messages. Used with
    AsynchronousStorageWrites.
    
    'MsgStorageApply' is a message from a node to its local apply storage
    thread to apply committed entries. The message will carry one response,
    which will be a 'MsgStorageApplyResp' back to itself. Used with
    AsynchronousStorageWrites.
    

### Index ¶

- Constants
- Variables
- func DescribeConfState(state pb.ConfState) string
- func DescribeEntries(ents []pb.Entry, f EntryFormatter) string
- func DescribeEntry(e pb.Entry, f EntryFormatter) string
- func DescribeHardState(hs pb.HardState) string
- func DescribeMessage(m pb.Message, f EntryFormatter) string
- func DescribeReady(rd Ready, f EntryFormatter) string
- func DescribeSnapshot(snap pb.Snapshot) string
- func DescribeSoftState(ss SoftState) string
- func IsEmptyHardState(st pb.HardState) bool
- func IsEmptySnap(sp pb.Snapshot) bool
- func IsLocalMsg(msgt pb.MessageType) bool
- func IsLocalMsgTarget(id uint64) bool
- func IsResponseMsg(msgt pb.MessageType) bool
- func MustSync(st, prevst pb.HardState, entsnum int) bool
- func ResetDefaultLogger()
- func SetLogger(l Logger)
- type BasicStatus
- type CampaignType
- type Config
- type DefaultLogger
-     * func (l *DefaultLogger) Debug(v ...interface{})
  - func (l *DefaultLogger) Debugf(format string, v ...interface{})
  - func (l *DefaultLogger) EnableDebug()
  - func (l *DefaultLogger) EnableTimestamps()
  - func (l *DefaultLogger) Error(v ...interface{})
  - func (l *DefaultLogger) Errorf(format string, v ...interface{})
  - func (l *DefaultLogger) Fatal(v ...interface{})
  - func (l *DefaultLogger) Fatalf(format string, v ...interface{})
  - func (l *DefaultLogger) Info(v ...interface{})
  - func (l *DefaultLogger) Infof(format string, v ...interface{})
  - func (l *DefaultLogger) Panic(v ...interface{})
  - func (l *DefaultLogger) Panicf(format string, v ...interface{})
  - func (l *DefaultLogger) Warning(v ...interface{})
  - func (l *DefaultLogger) Warningf(format string, v ...interface{})
- type EntryFormatter
- type Logger
- type MemoryStorage
-     * func NewMemoryStorage() *MemoryStorage
-     * func (ms *MemoryStorage) Append(entries []pb.Entry) error
  - func (ms *MemoryStorage) ApplySnapshot(snap pb.Snapshot) error
  - func (ms *MemoryStorage) Compact(compactIndex uint64) error
  - func (ms *MemoryStorage) CreateSnapshot(i uint64, cs*pb.ConfState, data []byte) (pb.Snapshot, error)
  - func (ms *MemoryStorage) Entries(lo, hi, maxSize uint64) ([]pb.Entry, error)
  - func (ms *MemoryStorage) FirstIndex() (uint64, error)
  - func (ms *MemoryStorage) InitialState() (pb.HardState, pb.ConfState, error)
  - func (ms *MemoryStorage) LastIndex() (uint64, error)
  - func (ms *MemoryStorage) SetHardState(st pb.HardState) error
  - func (ms *MemoryStorage) Snapshot() (pb.Snapshot, error)
  - func (ms *MemoryStorage) Term(i uint64) (uint64, error)
- type Node
-     * func RestartNode(c *Config) Node
  - func StartNode(c *Config, peers []Peer) Node
- type Peer
- type ProgressType
- type RawNode
-     * func NewRawNode(config *Config) (*RawNode, error)
-     * func (rn *RawNode) Advance(_ Ready)
  - func (rn *RawNode) ApplyConfChange(cc pb.ConfChangeI)*pb.ConfState
  - func (rn *RawNode) BasicStatus() BasicStatus
  - func (rn *RawNode) Bootstrap(peers []Peer) error
  - func (rn *RawNode) Campaign() error
  - func (rn *RawNode) ForgetLeader() error
  - func (rn *RawNode) HasReady() bool
  - func (rn *RawNode) Propose(data []byte) error
  - func (rn *RawNode) ProposeConfChange(cc pb.ConfChangeI) error
  - func (rn *RawNode) ReadIndex(rctx []byte)
  - func (rn *RawNode) Ready() Ready
  - func (rn *RawNode) ReportSnapshot(id uint64, status SnapshotStatus)
  - func (rn *RawNode) ReportUnreachable(id uint64)
  - func (rn *RawNode) Status() Status
  - func (rn *RawNode) Step(m pb.Message) error
  - func (rn *RawNode) Tick()
  - func (rn *RawNode) TickQuiesced()
  - func (rn *RawNode) TransferLeader(transferee uint64)
  - func (rn *RawNode) WithProgress(visitor func(id uint64, typ ProgressType, pr tracker.Progress))
- type ReadOnlyOption
- type ReadState
- type Ready
- type SnapshotStatus
- type SoftState
- type StateType
-     * func (st StateType) MarshalJSON() ([]byte, error)
  - func (st StateType) String() string
- type Status
-     * func (s Status) MarshalJSON() ([]byte, error)
  - func (s Status) String() string
- type Storage
- type TraceLogger
- type TracingEvent

### Examples ¶

- Node

### Constants ¶

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/raft.go#L34)

    const (
     // None is a placeholder node ID used when there is no leader.
     None [uint64](/builtin#uint64) = 0
     // LocalAppendThread is a reference to a local thread that saves unstable
     // log entries and snapshots to stable storage. The identifier is used as a
     // target for MsgStorageAppend messages when AsyncStorageWrites is enabled.
     LocalAppendThread [uint64](/builtin#uint64) = [math](/math).[MaxUint64](/math#MaxUint64)
     // LocalApplyThread is a reference to a local thread that applies committed
     // log entries to the local state machine. The identifier is used as a
     // target for MsgStorageApply messages when AsyncStorageWrites is enabled.
     LocalApplyThread [uint64](/builtin#uint64) = [math](/math).[MaxUint64](/math#MaxUint64) - 1
    )

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/state_trace_nop.go#L24)

    const StateTraceDeployed = [false](/builtin#false)

### Variables ¶

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L26)

    var ErrCompacted = [errors](/errors).[New](/errors#New)("requested index is unavailable due to compaction")

ErrCompacted is returned by Storage.Entries/Compact when a requested index is unavailable because it predates the last snapshot.

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/raft.go#L86)

    var ErrProposalDropped = [errors](/errors).[New](/errors#New)("raft proposal dropped")

ErrProposalDropped is returned when the proposal is ignored by some cases, so that the proposer can be notified and fail fast.

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L30)

    var ErrSnapOutOfDate = [errors](/errors).[New](/errors#New)("requested index is older than the existing snapshot")

ErrSnapOutOfDate is returned by Storage.CreateSnapshot when a requested index is older than the existing snapshot.

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L38)

    var ErrSnapshotTemporarilyUnavailable = [errors](/errors).[New](/errors#New)("snapshot is temporarily unavailable")

ErrSnapshotTemporarilyUnavailable is returned by the Storage interface when the required snapshot is temporarily unavailable.

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L25)

    var ErrStepLocalMsg = [errors](/errors).[New](/errors#New)("raft: cannot step raft local message")

ErrStepLocalMsg is returned when try to step a local raft message

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L29)

    var ErrStepPeerNotFound = [errors](/errors).[New](/errors#New)("raft: cannot step as peer not found")

ErrStepPeerNotFound is returned when try to step a response message but there is no peer found in raft.trk for that node.

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L31)

    var (
    
     // ErrStopped is returned by methods on Nodes that have been stopped.
     ErrStopped = [errors](/errors).[New](/errors#New)("raft: stopped")
    )

[View Source](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L34)

    var ErrUnavailable = [errors](/errors).[New](/errors#New)("requested entry at index is unavailable")

ErrUnavailable is returned by Storage interface when the requested log entries are unavailable.

### Functions ¶

#### func [DescribeConfState](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L95) ¶

    func DescribeConfState(state [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfState](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfState)) [string](/builtin#string)

#### func [DescribeEntries](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L244) ¶

    func DescribeEntries(ents [][pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Entry](/go.etcd.io/raft/v3@v3.6.0/raftpb#Entry), f EntryFormatter) [string](/builtin#string)

DescribeEntries calls DescribeEntry for each Entry, adding a newline to each.

#### func [DescribeEntry](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L206) ¶

    func DescribeEntry(e [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Entry](/go.etcd.io/raft/v3@v3.6.0/raftpb#Entry), f EntryFormatter) [string](/builtin#string)

DescribeEntry returns a concise human-readable description of an Entry for debugging.

#### func [DescribeHardState](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L81) ¶

    func DescribeHardState(hs [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState)) [string](/builtin#string)

#### func [DescribeMessage](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L150) ¶

    func DescribeMessage(m [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Message](/go.etcd.io/raft/v3@v3.6.0/raftpb#Message), f EntryFormatter) [string](/builtin#string)

DescribeMessage returns a concise human-readable description of a Message for debugging.

#### func [DescribeReady](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L107) ¶

    func DescribeReady(rd Ready, f EntryFormatter) [string](/builtin#string)

#### func [DescribeSnapshot](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L102) ¶

    func DescribeSnapshot(snap [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Snapshot](/go.etcd.io/raft/v3@v3.6.0/raftpb#Snapshot)) [string](/builtin#string)

#### func [DescribeSoftState](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L91) ¶

    func DescribeSoftState(ss SoftState) [string](/builtin#string)

#### func [IsEmptyHardState](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L122) ¶

    func IsEmptyHardState(st [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState)) [bool](/builtin#bool)

IsEmptyHardState returns true if the given HardState is empty.

#### func [IsEmptySnap](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L127) ¶

    func IsEmptySnap(sp [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Snapshot](/go.etcd.io/raft/v3@v3.6.0/raftpb#Snapshot)) [bool](/builtin#bool)

IsEmptySnap returns true if the given Snapshot is empty.

#### func [IsLocalMsg](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L57) ¶

    func IsLocalMsg(msgt [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[MessageType](/go.etcd.io/raft/v3@v3.6.0/raftpb#MessageType)) [bool](/builtin#bool)

#### func [IsLocalMsgTarget](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L65) ¶

    func IsLocalMsgTarget(id [uint64](/builtin#uint64)) [bool](/builtin#bool)

#### func [IsResponseMsg](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L61) ¶

    func IsResponseMsg(msgt [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[MessageType](/go.etcd.io/raft/v3@v3.6.0/raftpb#MessageType)) [bool](/builtin#bool)

#### func [MustSync](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L193) ¶

    func MustSync(st, prevst [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState), entsnum [int](/builtin#int)) [bool](/builtin#bool)

MustSync returns true if the hard state and count of Raft entries indicate that a synchronous write to persistent storage is required.

#### func [ResetDefaultLogger](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L51) ¶

    func ResetDefaultLogger()

#### func [SetLogger](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L45) ¶

    func SetLogger(l Logger)

### Types ¶

#### type [BasicStatus](https://github.com/etcd-io/raft/blob/v3.6.0/status.go#L33) ¶

    type BasicStatus struct {
     ID [uint64](/builtin#uint64)
    
     [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState)
     SoftState
    
     Applied [uint64](/builtin#uint64)
    
     LeadTransferee [uint64](/builtin#uint64)
    }

BasicStatus contains basic information about the Raft peer. It does not allocate.

#### type [CampaignType](https://github.com/etcd-io/raft/blob/v3.6.0/raft.go#L107) ¶

    type CampaignType [string](/builtin#string)

CampaignType represents the type of campaigning the reason we use the type of string instead of uint64 is because it's simpler to compare and fill in raft entries

#### type [Config](https://github.com/etcd-io/raft/blob/v3.6.0/raft.go#L124) ¶

    type Config struct {
     // ID is the identity of the local raft. ID cannot be 0.
     ID [uint64](/builtin#uint64)
    
     // ElectionTick is the number of Node.Tick invocations that must pass between
     // elections. That is, if a follower does not receive any message from the
     // leader of current term before ElectionTick has elapsed, it will become
     // candidate and start an election. ElectionTick must be greater than
     // HeartbeatTick. We suggest ElectionTick = 10 * HeartbeatTick to avoid
     // unnecessary leader switching.
     ElectionTick [int](/builtin#int)
     // HeartbeatTick is the number of Node.Tick invocations that must pass between
     // heartbeats. That is, a leader sends heartbeat messages to maintain its
     // leadership every HeartbeatTick ticks.
     HeartbeatTick [int](/builtin#int)
    
     // Storage is the storage for raft. raft generates entries and states to be
     // stored in storage. raft reads the persisted entries and states out of
     // Storage when it needs. raft reads out the previous state and configuration
     // out of storage when restarting.
     Storage Storage
     // Applied is the last applied index. It should only be set when restarting
     // raft. raft will not return entries to the application smaller or equal to
     // Applied. If Applied is unset when restarting, raft might return previous
     // applied entries. This is a very application dependent configuration.
     Applied [uint64](/builtin#uint64)
    
     // AsyncStorageWrites configures the raft node to write to its local storage
     // (raft log and state machine) using a request/response message passing
     // interface instead of the default Ready/Advance function call interface.
     // Local storage messages can be pipelined and processed asynchronously
     // (with respect to Ready iteration), facilitating reduced interference
     // between Raft proposals and increased batching of log appends and state
     // machine application. As a result, use of asynchronous storage writes can
     // reduce end-to-end commit latency and increase maximum throughput.
     //
     // When true, the Ready.Message slice will include MsgStorageAppend and
     // MsgStorageApply messages. The messages will target a LocalAppendThread
     // and a LocalApplyThread, respectively. Messages to the same target must be
     // reliably processed in order. In other words, they can't be dropped (like
     // messages over the network) and those targeted at the same thread can't be
     // reordered. Messages to different targets can be processed in any order.
     //
     // MsgStorageAppend carries Raft log entries to append, election votes /
     // term changes / updated commit indexes to persist, and snapshots to apply.
     // All writes performed in service of a MsgStorageAppend must be durable
     // before response messages are delivered. However, if the MsgStorageAppend
     // carries no response messages, durability is not required. The message
     // assumes the role of the Entries, HardState, and Snapshot fields in Ready.
     //
     // MsgStorageApply carries committed entries to apply. Writes performed in
     // service of a MsgStorageApply need not be durable before response messages
     // are delivered. The message assumes the role of the CommittedEntries field
     // in Ready.
     //
     // Local messages each carry one or more response messages which should be
     // delivered after the corresponding storage write has been completed. These
     // responses may target the same node or may target other nodes. The storage
     // threads are not responsible for understanding the response messages, only
     // for delivering them to the correct target after performing the storage
     // write.
     AsyncStorageWrites [bool](/builtin#bool)
    
     // MaxSizePerMsg limits the max byte size of each append message. Smaller
     // value lowers the raft recovery cost(initial probing and message lost
     // during normal operation). On the other side, it might affect the
     // throughput during normal replication. Note: math.MaxUint64 for unlimited,
     // 0 for at most one entry per message.
     MaxSizePerMsg [uint64](/builtin#uint64)
     // MaxCommittedSizePerReady limits the size of the committed entries which
     // can be applying at the same time.
     //
     // Despite its name (preserved for compatibility), this quota applies across
     // Ready structs to encompass all outstanding entries in unacknowledged
     // MsgStorageApply messages when AsyncStorageWrites is enabled.
     MaxCommittedSizePerReady [uint64](/builtin#uint64)
     // MaxUncommittedEntriesSize limits the aggregate byte size of the
     // uncommitted entries that may be appended to a leader's log. Once this
     // limit is exceeded, proposals will begin to return ErrProposalDropped
     // errors. Note: 0 for no limit.
     MaxUncommittedEntriesSize [uint64](/builtin#uint64)
     // MaxInflightMsgs limits the max number of in-flight append messages during
     // optimistic replication phase. The application transportation layer usually
     // has its own sending buffer over TCP/UDP. Setting MaxInflightMsgs to avoid
     // overflowing that sending buffer. TODO (xiangli): feedback to application to
     // limit the proposal rate?
     MaxInflightMsgs [int](/builtin#int)
     // MaxInflightBytes limits the number of in-flight bytes in append messages.
     // Complements MaxInflightMsgs. Ignored if zero.
     //
     // This effectively bounds the bandwidth-delay product. Note that especially
     // in high-latency deployments setting this too low can lead to a dramatic
     // reduction in throughput. For example, with a peer that has a round-trip
     // latency of 100ms to the leader and this setting is set to 1 MB, there is a
     // throughput limit of 10 MB/s for this group. With RTT of 400ms, this drops
     // to 2.5 MB/s. See Little's law to understand the maths behind.
     MaxInflightBytes [uint64](/builtin#uint64)
    
     // CheckQuorum specifies if the leader should check quorum activity. Leader
     // steps down when quorum is not active for an electionTimeout.
     CheckQuorum [bool](/builtin#bool)
    
     // PreVote enables the Pre-Vote algorithm described in raft thesis section
     // 9.6. This prevents disruption when a node that has been partitioned away
     // rejoins the cluster.
     PreVote [bool](/builtin#bool)
    
     // ReadOnlyOption specifies how the read only request is processed.
     //
     // ReadOnlySafe guarantees the linearizability of the read only request by
     // communicating with the quorum. It is the default and suggested option.
     //
     // ReadOnlyLeaseBased ensures linearizability of the read only request by
     // relying on the leader lease. It can be affected by clock drift.
     // If the clock drift is unbounded, leader might keep the lease longer than it
     // should (clock can move backward/pause without any bound). ReadIndex is not safe
     // in that case.
     // CheckQuorum MUST be enabled if ReadOnlyOption is ReadOnlyLeaseBased.
     ReadOnlyOption ReadOnlyOption
    
     // Logger is the logger used for raft log. For multinode which can host
     // multiple raft group, each raft group can have its own logger
     Logger Logger
    
     // DisableProposalForwarding set to true means that followers will drop
     // proposals, rather than forwarding them to the leader. One use case for
     // this feature would be in a situation where the Raft leader is used to
     // compute the data of a proposal, for example, adding a timestamp from a
     // hybrid logical clock to data in a monotonically increasing way. Forwarding
     // should be disabled to prevent a follower with an inaccurate hybrid
     // logical clock from assigning the timestamp and then forwarding the data
     // to the leader.
     DisableProposalForwarding [bool](/builtin#bool)
    
     // DisableConfChangeValidation turns off propose-time verification of
     // configuration changes against the currently active configuration of the
     // raft instance. These checks are generally sensible (cannot leave a joint
     // config unless in a joint config, et cetera) but they have false positives
     // because the active configuration may not be the most recent
     // configuration. This is because configurations are activated during log
     // application, and even the leader can trail log application by an
     // unbounded number of entries.
     // Symmetrically, the mechanism has false negatives - because the check may
     // not run against the "actual" config that will be the predecessor of the
     // newly proposed one, the check may pass but the new config may be invalid
     // when it is being applied. In other words, the checks are best-effort.
     //
     // Users should *not* use this option unless they have a reliable mechanism
     // (above raft) that serializes and verifies configuration changes. If an
     // invalid configuration change enters the log and gets applied, a panic
     // will result.
     //
     // This option may be removed once false positives are no longer possible.
     // See: <https://github.com/etcd-io/raft/issues/80>
     DisableConfChangeValidation [bool](/builtin#bool)
    
     // StepDownOnRemoval makes the leader step down when it is removed from the
     // group or demoted to a learner.
     //
     // This behavior will become unconditional in the future. See:
     // <https://github.com/etcd-io/raft/issues/83>
     StepDownOnRemoval [bool](/builtin#bool)
    
     // raft state tracer
     TraceLogger TraceLogger
    }

Config contains the parameters to start a raft.

#### type [DefaultLogger](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L73) ¶

    type DefaultLogger struct {
     *[log](/log).[Logger](/log#Logger)
     // contains filtered or unexported fields
    }

DefaultLogger is a default implementation of the Logger interface.

#### func (*DefaultLogger) [Debug](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L86) ¶

    func (l *DefaultLogger) Debug(v ...interface{})

#### func (*DefaultLogger) [Debugf](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L92) ¶

    func (l *DefaultLogger) Debugf(format [string](/builtin#string), v ...interface{})

#### func (*DefaultLogger) [EnableDebug](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L82) ¶

    func (l *DefaultLogger) EnableDebug()

#### func (*DefaultLogger) [EnableTimestamps](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L78) ¶

    func (l *DefaultLogger) EnableTimestamps()

#### func (*DefaultLogger) [Error](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L106) ¶

    func (l *DefaultLogger) Error(v ...interface{})

#### func (*DefaultLogger) [Errorf](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L110) ¶

    func (l *DefaultLogger) Errorf(format [string](/builtin#string), v ...interface{})

#### func (*DefaultLogger) [Fatal](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L122) ¶

    func (l *DefaultLogger) Fatal(v ...interface{})

#### func (*DefaultLogger) [Fatalf](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L127) ¶

    func (l *DefaultLogger) Fatalf(format [string](/builtin#string), v ...interface{})

#### func (*DefaultLogger) [Info](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L98) ¶

    func (l *DefaultLogger) Info(v ...interface{})

#### func (*DefaultLogger) [Infof](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L102) ¶

    func (l *DefaultLogger) Infof(format [string](/builtin#string), v ...interface{})

#### func (*DefaultLogger) [Panic](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L132) ¶

    func (l *DefaultLogger) Panic(v ...interface{})

#### func (*DefaultLogger) [Panicf](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L136) ¶

    func (l *DefaultLogger) Panicf(format [string](/builtin#string), v ...interface{})

#### func (*DefaultLogger) [Warning](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L114) ¶

    func (l *DefaultLogger) Warning(v ...interface{})

#### func (*DefaultLogger) [Warningf](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L118) ¶

    func (l *DefaultLogger) Warningf(format [string](/builtin#string), v ...interface{})

#### type [EntryFormatter](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L146) ¶

    type EntryFormatter func([][byte](/builtin#byte)) [string](/builtin#string)

EntryFormatter can be implemented by the application to provide human-readable formatting of entry data. Nil is a valid EntryFormatter and will use a default format.

#### type [Logger](https://github.com/etcd-io/raft/blob/v3.6.0/logger.go#L25) ¶

    type Logger interface {
     Debug(v ...interface{})
     Debugf(format [string](/builtin#string), v ...interface{})
    
     Error(v ...interface{})
     Errorf(format [string](/builtin#string), v ...interface{})
    
     Info(v ...interface{})
     Infof(format [string](/builtin#string), v ...interface{})
    
     Warning(v ...interface{})
     Warningf(format [string](/builtin#string), v ...interface{})
    
     Fatal(v ...interface{})
     Fatalf(format [string](/builtin#string), v ...interface{})
    
     Panic(v ...interface{})
     Panicf(format [string](/builtin#string), v ...interface{})
    }

#### type [MemoryStorage](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L98) ¶

    type MemoryStorage struct {
     // Protects access to all fields. Most methods of MemoryStorage are
     // run on the raft goroutine, but Append() is run on an application
     // goroutine.
     [sync](/sync).[Mutex](/sync#Mutex)
     // contains filtered or unexported fields
    }

MemoryStorage implements the Storage interface backed by an in-memory array.

#### func [NewMemoryStorage](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L113) ¶

    func NewMemoryStorage() *MemoryStorage

NewMemoryStorage creates an empty MemoryStorage.

#### func (*MemoryStorage) [Append](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L277) ¶

    func (ms *MemoryStorage) Append(entries [][pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Entry](/go.etcd.io/raft/v3@v3.6.0/raftpb#Entry)) [error](/builtin#error)

Append the new entries to storage. TODO (xiangli): ensure the entries are continuous and entries[0].Index > ms.entries[0].Index

#### func (*MemoryStorage) [ApplySnapshot](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L207) ¶

    func (ms *MemoryStorage) ApplySnapshot(snap [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Snapshot](/go.etcd.io/raft/v3@v3.6.0/raftpb#Snapshot)) [error](/builtin#error)

ApplySnapshot overwrites the contents of this Storage object with those of the given snapshot.

#### func (*MemoryStorage) [Compact](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L251) ¶

    func (ms *MemoryStorage) Compact(compactIndex [uint64](/builtin#uint64)) [error](/builtin#error)

Compact discards all log entries prior to compactIndex. It is the application's responsibility to not attempt to compact an index greater than raftLog.applied.

#### func (*MemoryStorage) [CreateSnapshot](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L227) ¶

    func (ms *MemoryStorage) CreateSnapshot(i [uint64](/builtin#uint64), cs *[pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfState](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfState), data [][byte](/builtin#byte)) ([pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Snapshot](/go.etcd.io/raft/v3@v3.6.0/raftpb#Snapshot), [error](/builtin#error))

CreateSnapshot makes a snapshot which can be retrieved with Snapshot() and can be used to reconstruct the state at that point. If any configuration changes have been made since the last compaction, the result of the last ApplyConfChange must be passed in.

#### func (*MemoryStorage) [Entries](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L135) ¶

    func (ms *MemoryStorage) Entries(lo, hi, maxSize [uint64](/builtin#uint64)) ([][pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Entry](/go.etcd.io/raft/v3@v3.6.0/raftpb#Entry), [error](/builtin#error))

Entries implements the Storage interface.

#### func (*MemoryStorage) [FirstIndex](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L186) ¶

    func (ms *MemoryStorage) FirstIndex() ([uint64](/builtin#uint64), [error](/builtin#error))

FirstIndex implements the Storage interface.

#### func (*MemoryStorage) [InitialState](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L121) ¶

    func (ms *MemoryStorage) InitialState() ([pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState), [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfState](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfState), [error](/builtin#error))

InitialState implements the Storage interface.

#### func (*MemoryStorage) [LastIndex](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L174) ¶

    func (ms *MemoryStorage) LastIndex() ([uint64](/builtin#uint64), [error](/builtin#error))

LastIndex implements the Storage interface.

#### func (*MemoryStorage) [SetHardState](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L127) ¶

    func (ms *MemoryStorage) SetHardState(st [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState)) [error](/builtin#error)

SetHardState saves the current HardState.

#### func (*MemoryStorage) [Snapshot](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L198) ¶

    func (ms *MemoryStorage) Snapshot() ([pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Snapshot](/go.etcd.io/raft/v3@v3.6.0/raftpb#Snapshot), [error](/builtin#error))

Snapshot implements the Storage interface.

#### func (*MemoryStorage) [Term](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L159) ¶

    func (ms *MemoryStorage) Term(i [uint64](/builtin#uint64)) ([uint64](/builtin#uint64), [error](/builtin#error))

Term implements the Storage interface.

#### type [Node](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L132) ¶

    type Node interface {
     // Tick increments the internal logical clock for the Node by a single tick. Election
     // timeouts and heartbeat timeouts are in units of ticks.
     Tick()
     // Campaign causes the Node to transition to candidate state and start campaigning to become leader.
     Campaign(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)
     // Propose proposes that data be appended to the log. Note that proposals can be lost without
     // notice, therefore it is user's job to ensure proposal retries.
     Propose(ctx [context](/context).[Context](/context#Context), data [][byte](/builtin#byte)) [error](/builtin#error)
     // ProposeConfChange proposes a configuration change. Like any proposal, the
     // configuration change may be dropped with or without an error being
     // returned. In particular, configuration changes are dropped unless the
     // leader has certainty that there is no prior unapplied configuration
     // change in its log.
     //
     // The method accepts either a pb.ConfChange (deprecated) or pb.ConfChangeV2
     // message. The latter allows arbitrary configuration changes via joint
     // consensus, notably including replacing a voter. Passing a ConfChangeV2
     // message is only allowed if all Nodes participating in the cluster run a
     // version of this library aware of the V2 API. See pb.ConfChangeV2 for
     // usage details and semantics.
     ProposeConfChange(ctx [context](/context).[Context](/context#Context), cc [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfChangeI](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfChangeI)) [error](/builtin#error)
    
     // Step advances the state machine using the given message. ctx.Err() will be returned, if any.
     Step(ctx [context](/context).[Context](/context#Context), msg [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Message](/go.etcd.io/raft/v3@v3.6.0/raftpb#Message)) [error](/builtin#error)
    
     // Ready returns a channel that returns the current point-in-time state.
     // Users of the Node must call Advance after retrieving the state returned by Ready (unless
     // async storage writes is enabled, in which case it should never be called).
     //
     // NOTE: No committed entries from the next Ready may be applied until all committed entries
     // and snapshots from the previous one have finished.
     Ready() <-chan Ready
    
     // Advance notifies the Node that the application has saved progress up to the last Ready.
     // It prepares the node to return the next available Ready.
     //
     // The application should generally call Advance after it applies the entries in last Ready.
     //
     // However, as an optimization, the application may call Advance while it is applying the
     // commands. For example. when the last Ready contains a snapshot, the application might take
     // a long time to apply the snapshot data. To continue receiving Ready without blocking raft
     // progress, it can call Advance before finishing applying the last ready.
     //
     // NOTE: Advance must not be called when using AsyncStorageWrites. Response messages from the
     // local append and apply threads take its place.
     Advance()
     // ApplyConfChange applies a config change (previously passed to
     // ProposeConfChange) to the node. This must be called whenever a config
     // change is observed in Ready.CommittedEntries, except when the app decides
     // to reject the configuration change (i.e. treats it as a noop instead), in
     // which case it must not be called.
     //
     // Returns an opaque non-nil ConfState protobuf which must be recorded in
     // snapshots.
     ApplyConfChange(cc [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfChangeI](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfChangeI)) *[pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfState](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfState)
    
     // TransferLeadership attempts to transfer leadership to the given transferee.
     TransferLeadership(ctx [context](/context).[Context](/context#Context), lead, transferee [uint64](/builtin#uint64))
    
     // ForgetLeader forgets a follower's current leader, changing it to None. It
     // remains a leaderless follower in the current term, without campaigning.
     //
     // This is useful with PreVote+CheckQuorum, where followers will normally not
     // grant pre-votes if they've heard from the leader in the past election
     // timeout interval. Leaderless followers can grant pre-votes immediately, so
     // if a quorum of followers have strong reason to believe the leader is dead
     // (for example via a side-channel or external failure detector) and forget it
     // then they can elect a new leader immediately, without waiting out the
     // election timeout. They will also revert to normal followers if they hear
     // from the leader again, or transition to candidates on an election timeout.
     //
     // For example, consider a three-node cluster where 1 is the leader and 2+3
     // have just received a heartbeat from it. If 2 and 3 believe the leader has
     // now died (maybe they know that an orchestration system shut down 1's VM),
     // we can instruct 2 to forget the leader and 3 to campaign. 2 will then be
     // able to grant 3's pre-vote and elect 3 as leader immediately (normally 2
     // would reject the vote until an election timeout passes because it has heard
     // from the leader recently). However, 3 can not campaign unilaterally, a
     // quorum have to agree that the leader is dead, which avoids disrupting the
     // leader if individual nodes are wrong about it being dead.
     //
     // This does nothing with ReadOnlyLeaseBased, since it would allow a new
     // leader to be elected without the old leader knowing.
     ForgetLeader(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)
    
     // ReadIndex request a read state. The read state will be set in the ready.
     // Read state has a read index. Once the application advances further than the read
     // index, any linearizable read requests issued before the read request can be
     // processed safely. The read state will have the same rctx attached.
     // Note that request can be lost without notice, therefore it is user's job
     // to ensure read index retries.
     ReadIndex(ctx [context](/context).[Context](/context#Context), rctx [][byte](/builtin#byte)) [error](/builtin#error)
    
     // Status returns the current status of the raft state machine.
     Status() Status
     // ReportUnreachable reports the given node is not reachable for the last send.
     ReportUnreachable(id [uint64](/builtin#uint64))
     // ReportSnapshot reports the status of the sent snapshot. The id is the raft ID of the follower
     // who is meant to receive the snapshot, and the status is SnapshotFinish or SnapshotFailure.
     // Calling ReportSnapshot with SnapshotFinish is a no-op. But, any failure in applying a
     // snapshot (for e.g., while streaming it from leader to follower), should be reported to the
     // leader with SnapshotFailure. When leader sends a snapshot to a follower, it pauses any raft
     // log probes until the follower can apply the snapshot and advance its state. If the follower
     // can't do that, for e.g., due to a crash, it could end up in a limbo, never getting any
     // updates from the leader. Therefore, it is crucial that the application ensures that any
     // failure in snapshot sending is caught and reported back to the leader; so it can resume raft
     // log probing in the follower.
     ReportSnapshot(id [uint64](/builtin#uint64), status SnapshotStatus)
     // Stop performs any necessary termination of the Node.
     Stop()
    }

Node represents a node in a raft cluster.

Example ¶

    package main
    
    import (
     pb "go.etcd.io/raft/v3/raftpb"
    )
    
    func applyToStore(_ []pb.Entry)      {}
    func sendMessages(_ []pb.Message)    {}
    func saveStateToDisk(_ pb.HardState) {}
    func saveToDisk(_ []pb.Entry)        {}
    
    func main() {
     c := &Config{}
     n := StartNode(c, nil)
     defer n.Stop()
    
     // stuff to n happens in other goroutines
    
     // the last known state
     var prev pb.HardState
     for {
      // Ready blocks until there is new state ready.
      rd := <-n.Ready()
      if !isHardStateEqual(prev, rd.HardState) {
       saveStateToDisk(rd.HardState)
       prev = rd.HardState
      }
    
      saveToDisk(rd.Entries)
      go applyToStore(rd.CommittedEntries)
      sendMessages(rd.Messages)
     }
    }
    

Share Format Run

#### func [RestartNode](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L281) ¶

    func RestartNode(c *Config) Node

RestartNode is similar to StartNode but does not take a list of peers. The current membership of the cluster will be restored from the Storage. If the caller has an existing state machine, pass in the last log index that has been applied to it; otherwise use zero.

#### func [StartNode](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L271) ¶

    func StartNode(c *Config, peers []Peer) Node

StartNode returns a new Node given configuration and a list of raft peers. It appends a ConfChangeAddNode entry for each given peer to the initial log.

Peers must not be zero length; call RestartNode in that case.

#### type [Peer](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L245) ¶

    type Peer struct {
     ID      [uint64](/builtin#uint64)
     Context [][byte](/builtin#byte)
    }

#### type [ProgressType](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L510) ¶

    type ProgressType [byte](/builtin#byte)

ProgressType indicates the type of replica a Progress corresponds to.

    const (
     // ProgressTypePeer accompanies a Progress for a regular peer replica.
     ProgressTypePeer ProgressType = [iota](/builtin#iota)
     // ProgressTypeLearner accompanies a Progress for a learner replica.
     ProgressTypeLearner
    )

#### type [RawNode](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L34) ¶

    type RawNode struct {
     // contains filtered or unexported fields
    }

RawNode is a thread-unsafe Node. The methods of this struct correspond to the methods of Node and are described more fully there.

#### func [NewRawNode](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L51) ¶

    func NewRawNode(config *Config) (*RawNode, [error](/builtin#error))

NewRawNode instantiates a RawNode from the given configuration.

See Bootstrap() for bootstrapping an initial state; this replaces the former 'peers' argument to this method (with identical behavior). However, It is recommended that instead of calling Bootstrap, applications bootstrap their state manually by setting up a Storage that has a first index > 1 and which stores the desired ConfState as its InitialState.

#### func (*RawNode) [Advance](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L482) ¶

    func (rn *RawNode) Advance(_ Ready)

Advance notifies the RawNode that the application has applied and saved progress in the last Ready results.

NOTE: Advance must not be called when using AsyncStorageWrites. Response messages from the local append and apply threads take its place.

#### func (*RawNode) [ApplyConfChange](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L112) ¶

    func (rn *RawNode) ApplyConfChange(cc [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfChangeI](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfChangeI)) *[pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfState](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfState)

ApplyConfChange applies a config change to the local node. The app must call this when it applies a configuration change, except when it decides to reject the configuration change, in which case no call must take place.

#### func (*RawNode) [BasicStatus](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L505) ¶

    func (rn *RawNode) BasicStatus() BasicStatus

BasicStatus returns a BasicStatus. Notably this does not contain the Progress map; see WithProgress for an allocation-free way to inspect it.

#### func (*RawNode) [Bootstrap](https://github.com/etcd-io/raft/blob/v3.6.0/bootstrap.go#L30) ¶

    func (rn *RawNode) Bootstrap(peers []Peer) [error](/builtin#error)

Bootstrap initializes the RawNode for first use by appending configuration changes for the supplied peers. This method returns an error if the Storage is nonempty.

It is recommended that instead of calling this method, applications bootstrap their state manually by setting up a Storage that has a first index > 1 and which stores the desired ConfState as its InitialState.

#### func (*RawNode) [Campaign](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L83) ¶

    func (rn *RawNode) Campaign() [error](/builtin#error)

Campaign causes this RawNode to transition to candidate state.

#### func (*RawNode) [ForgetLeader](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L552) ¶

    func (rn *RawNode) ForgetLeader() [error](/builtin#error)

ForgetLeader forgets a follower's current leader, changing it to None. See (Node).ForgetLeader for details.

#### func (*RawNode) [HasReady](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L453) ¶

    func (rn *RawNode) HasReady() [bool](/builtin#bool)

HasReady called when RawNode user need to check if any Ready pending.

#### func (*RawNode) [Propose](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L90) ¶

    func (rn *RawNode) Propose(data [][byte](/builtin#byte)) [error](/builtin#error)

Propose proposes data be appended to the raft log.

#### func (*RawNode) [ProposeConfChange](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L101) ¶

    func (rn *RawNode) ProposeConfChange(cc [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfChangeI](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfChangeI)) [error](/builtin#error)

ProposeConfChange proposes a config change. See (Node).ProposeConfChange for details.

#### func (*RawNode) [ReadIndex](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L560) ¶

    func (rn *RawNode) ReadIndex(rctx [][byte](/builtin#byte))

ReadIndex requests a read state. The read state will be set in ready. Read State has a read index. Once the application advances further than the read index, any linearizable read requests issued before the read request can be processed safely. The read state will have the same rctx attached.

#### func (*RawNode) [Ready](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L133) ¶

    func (rn *RawNode) Ready() Ready

Ready returns the outstanding work that the application needs to handle. This includes appending and applying entries or a snapshot, updating the HardState, and sending messages. The returned Ready() _must_ be handled and subsequently passed back via Advance().

#### func (*RawNode) [ReportSnapshot](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L539) ¶

    func (rn *RawNode) ReportSnapshot(id [uint64](/builtin#uint64), status SnapshotStatus)

ReportSnapshot reports the status of the sent snapshot.

#### func (*RawNode) [ReportUnreachable](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L534) ¶

    func (rn *RawNode) ReportUnreachable(id [uint64](/builtin#uint64))

ReportUnreachable reports the given node is not reachable for the last send.

#### func (*RawNode) [Status](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L498) ¶

    func (rn *RawNode) Status() Status

Status returns the current status of the given group. This allocates, see BasicStatus and WithProgress for allocation-friendlier choices.

#### func (*RawNode) [Step](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L118) ¶

    func (rn *RawNode) Step(m [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Message](/go.etcd.io/raft/v3@v3.6.0/raftpb#Message)) [error](/builtin#error)

Step advances the state machine using the given message.

#### func (*RawNode) [Tick](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L64) ¶

    func (rn *RawNode) Tick()

Tick advances the internal logical clock by a single tick.

#### func (*RawNode) [TickQuiesced](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L78) ¶

    func (rn *RawNode) TickQuiesced()

TickQuiesced advances the internal logical clock by a single tick without performing any other state machine processing. It allows the caller to avoid periodic heartbeats and elections when all of the peers in a Raft group are known to be at the same state. Expected usage is to periodically invoke Tick or TickQuiesced depending on whether the group is "active" or "quiesced".

WARNING: Be very careful about using this method as it subverts the Raft state machine. You should probably be using Tick instead.

DEPRECATED: This method will be removed in a future release.

#### func (*RawNode) [TransferLeader](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L546) ¶

    func (rn *RawNode) TransferLeader(transferee [uint64](/builtin#uint64))

TransferLeader tries to transfer leadership to the given transferee.

#### func (*RawNode) [WithProgress](https://github.com/etcd-io/raft/blob/v3.6.0/rawnode.go#L521) ¶

    func (rn *RawNode) WithProgress(visitor func(id [uint64](/builtin#uint64), typ ProgressType, pr [tracker](/go.etcd.io/raft/v3@v3.6.0/tracker).[Progress](/go.etcd.io/raft/v3@v3.6.0/tracker#Progress)))

WithProgress is a helper to introspect the Progress for this node and its peers.

#### type [ReadOnlyOption](https://github.com/etcd-io/raft/blob/v3.6.0/raft.go#L56) ¶

    type ReadOnlyOption [int](/builtin#int)
    
    
    const (
     // ReadOnlySafe guarantees the linearizability of the read only request by
     // communicating with the quorum. It is the default and suggested option.
     ReadOnlySafe ReadOnlyOption = [iota](/builtin#iota)
     // ReadOnlyLeaseBased ensures linearizability of the read only request by
     // relying on the leader lease. It can be affected by clock drift.
     // If the clock drift is unbounded, leader might keep the lease longer than it
     // should (clock can move backward/pause without any bound). ReadIndex is not safe
     // in that case.
     ReadOnlyLeaseBased
    )

#### type [ReadState](https://github.com/etcd-io/raft/blob/v3.6.0/read_only.go#L24) ¶

    type ReadState struct {
     Index      [uint64](/builtin#uint64)
     RequestCtx [][byte](/builtin#byte)
    }

ReadState provides state for read only query. It's caller's responsibility to call ReadIndex first before getting this state from ready, it's also caller's duty to differentiate if this state is what it requests through RequestCtx, eg. given a unique id as RequestCtx

#### type [Ready](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L52) ¶

    type Ready struct {
     // The current volatile state of a Node.
     // SoftState will be nil if there is no update.
     // It is not required to consume or store SoftState.
     *SoftState
    
     // The current state of a Node to be saved to stable storage BEFORE
     // Messages are sent.
     //
     // HardState will be equal to empty state if there is no update.
     //
     // If async storage writes are enabled, this field does not need to be acted
     // on immediately. It will be reflected in a MsgStorageAppend message in the
     // Messages slice.
     [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState)
    
     // ReadStates can be used for node to serve linearizable read requests locally
     // when its applied index is greater than the index in ReadState.
     // Note that the readState will be returned when raft receives msgReadIndex.
     // The returned is only valid for the request that requested to read.
     ReadStates []ReadState
    
     // Entries specifies entries to be saved to stable storage BEFORE
     // Messages are sent.
     //
     // If async storage writes are enabled, this field does not need to be acted
     // on immediately. It will be reflected in a MsgStorageAppend message in the
     // Messages slice.
     Entries [][pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Entry](/go.etcd.io/raft/v3@v3.6.0/raftpb#Entry)
    
     // Snapshot specifies the snapshot to be saved to stable storage.
     //
     // If async storage writes are enabled, this field does not need to be acted
     // on immediately. It will be reflected in a MsgStorageAppend message in the
     // Messages slice.
     Snapshot [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Snapshot](/go.etcd.io/raft/v3@v3.6.0/raftpb#Snapshot)
    
     // CommittedEntries specifies entries to be committed to a
     // store/state-machine. These have previously been appended to stable
     // storage.
     //
     // If async storage writes are enabled, this field does not need to be acted
     // on immediately. It will be reflected in a MsgStorageApply message in the
     // Messages slice.
     CommittedEntries [][pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Entry](/go.etcd.io/raft/v3@v3.6.0/raftpb#Entry)
    
     // Messages specifies outbound messages.
     //
     // If async storage writes are not enabled, these messages must be sent
     // AFTER Entries are appended to stable storage.
     //
     // If async storage writes are enabled, these messages can be sent
     // immediately as the messages that have the completion of the async writes
     // as a precondition are attached to the individual MsgStorage{Append,Apply}
     // messages instead.
     //
     // If it contains a MsgSnap message, the application MUST report back to raft
     // when the snapshot has been received or has failed by calling ReportSnapshot.
     Messages [][pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Message](/go.etcd.io/raft/v3@v3.6.0/raftpb#Message)
    
     // MustSync indicates whether the HardState and Entries must be durably
     // written to disk or if a non-durable write is permissible.
     MustSync [bool](/builtin#bool)
    }

Ready encapsulates the entries and messages that are ready to read, be saved to stable storage, committed or sent to other peers. All fields in Ready are read-only.

#### type [SnapshotStatus](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L24) ¶

    type SnapshotStatus [int](/builtin#int)
    
    
    const (
     SnapshotFinish  SnapshotStatus = 1
     SnapshotFailure SnapshotStatus = 2
    )

#### type [SoftState](https://github.com/etcd-io/raft/blob/v3.6.0/node.go#L40) ¶

    type SoftState struct {
     Lead      [uint64](/builtin#uint64) // must use atomic operations to access; keep 64-bit aligned.
     RaftState StateType
    }

SoftState provides state that is useful for logging and debugging. The state is volatile and does not need to be persisted to the WAL.

#### type [StateType](https://github.com/etcd-io/raft/blob/v3.6.0/raft.go#L110) ¶

    type StateType [uint64](/builtin#uint64)

StateType represents the role of a node in a cluster.

    const (
     StateFollower StateType = [iota](/builtin#iota)
     StateCandidate
     StateLeader
     StatePreCandidate
    )

Possible values for StateType.

#### func (StateType) [MarshalJSON](https://github.com/etcd-io/raft/blob/v3.6.0/util.go#L25) ¶

    func (st StateType) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

#### func (StateType) [String](https://github.com/etcd-io/raft/blob/v3.6.0/raft.go#L119) ¶

    func (st StateType) String() [string](/builtin#string)

#### type [Status](https://github.com/etcd-io/raft/blob/v3.6.0/status.go#L26) ¶

    type Status struct {
     BasicStatus
     Config   [tracker](/go.etcd.io/raft/v3@v3.6.0/tracker).[Config](/go.etcd.io/raft/v3@v3.6.0/tracker#Config)
     Progress map[[uint64](/builtin#uint64)][tracker](/go.etcd.io/raft/v3@v3.6.0/tracker).[Progress](/go.etcd.io/raft/v3@v3.6.0/tracker#Progress)
    }

Status contains information about this Raft peer and its view of the system. The Progress is only populated on the leader.

#### func (Status) [MarshalJSON](https://github.com/etcd-io/raft/blob/v3.6.0/status.go#L80) ¶

    func (s Status) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON translates the raft status into JSON. TODO: try to simplify this by introducing ID type into raft

#### func (Status) [String](https://github.com/etcd-io/raft/blob/v3.6.0/status.go#L99) ¶

    func (s Status) String() [string](/builtin#string)

#### type [Storage](https://github.com/etcd-io/raft/blob/v3.6.0/storage.go#L46) ¶

    type Storage interface {
    
     // InitialState returns the saved HardState and ConfState information.
     InitialState() ([pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[HardState](/go.etcd.io/raft/v3@v3.6.0/raftpb#HardState), [pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[ConfState](/go.etcd.io/raft/v3@v3.6.0/raftpb#ConfState), [error](/builtin#error))
    
     // Entries returns a slice of consecutive log entries in the range [lo, hi),
     // starting from lo. The maxSize limits the total size of the log entries
     // returned, but Entries returns at least one entry if any.
     //
     // The caller of Entries owns the returned slice, and may append to it. The
     // individual entries in the slice must not be mutated, neither by the Storage
     // implementation nor the caller. Note that raft may forward these entries
     // back to the application via Ready struct, so the corresponding handler must
     // not mutate entries either (see comments in Ready struct).
     //
     // Since the caller may append to the returned slice, Storage implementation
     // must protect its state from corruption that such appends may cause. For
     // example, common ways to do so are:
     //  - allocate the slice before returning it (safest option),
     //  - return a slice protected by Go full slice expression, which causes
     //  copying on appends (see MemoryStorage).
     //
     // Returns ErrCompacted if entry lo has been compacted, or ErrUnavailable if
     // encountered an unavailable entry in [lo, hi).
     Entries(lo, hi, maxSize [uint64](/builtin#uint64)) ([][pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Entry](/go.etcd.io/raft/v3@v3.6.0/raftpb#Entry), [error](/builtin#error))
    
     // Term returns the term of entry i, which must be in the range
     // [FirstIndex()-1, LastIndex()]. The term of the entry before
     // FirstIndex is retained for matching purposes even though the
     // rest of that entry may not be available.
     Term(i [uint64](/builtin#uint64)) ([uint64](/builtin#uint64), [error](/builtin#error))
     // LastIndex returns the index of the last entry in the log.
     LastIndex() ([uint64](/builtin#uint64), [error](/builtin#error))
     // FirstIndex returns the index of the first log entry that is
     // possibly available via Entries (older entries have been incorporated
     // into the latest Snapshot; if storage only contains the dummy entry the
     // first log entry is not available).
     FirstIndex() ([uint64](/builtin#uint64), [error](/builtin#error))
     // Snapshot returns the most recent snapshot.
     // If snapshot is temporarily unavailable, it should return ErrSnapshotTemporarilyUnavailable,
     // so raft state machine could know that Storage needs some time to prepare
     // snapshot and call Snapshot later.
     Snapshot() ([pb](/go.etcd.io/raft/v3@v3.6.0/raftpb).[Snapshot](/go.etcd.io/raft/v3@v3.6.0/raftpb#Snapshot), [error](/builtin#error))
    }

Storage is an interface that may be implemented by the application to retrieve log entries from storage.

If any Storage method returns an error, the raft instance will become inoperable and refuse to participate in elections; the application is responsible for cleanup and recovery in this case.

#### type [TraceLogger](https://github.com/etcd-io/raft/blob/v3.6.0/state_trace_nop.go#L26) ¶

    type TraceLogger interface{}

#### type [TracingEvent](https://github.com/etcd-io/raft/blob/v3.6.0/state_trace_nop.go#L28) ¶

    type TracingEvent struct{}
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
