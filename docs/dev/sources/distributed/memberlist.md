# hashicorp/memberlist

> Source: https://pkg.go.dev/github.com/hashicorp/memberlist
> Fetched: 2026-01-30T23:56:02.682511+00:00
> Content-Hash: 99cd58803c241e6c
> Type: html

---

Overview

¶

memberlist is a library that manages cluster
membership and member failure detection using a gossip based protocol.

The use cases for such a library are far-reaching: all distributed systems
require membership, and memberlist is a re-usable solution to managing
cluster membership and node failure detection.

memberlist is eventually consistent but converges quickly on average.
The speed at which it converges can be heavily tuned via various knobs
on the protocol. Node failures are detected and network partitions are partially
tolerated by attempting to communicate to potentially dead nodes through
multiple routes.

Index

¶

Constants

func AddLabelHeaderToPacket(buf []byte, label string) ([]byte, error)

func AddLabelHeaderToStream(conn net.Conn, label string) error

func LogAddress(addr net.Addr) string

func LogConn(conn net.Conn) string

func LogStringAddress(addr string) string

func ParseCIDRs(v []string) ([]net.IPNet, error)

func RemoveLabelHeaderFromPacket(buf []byte) (newBuf []byte, label string, err error)

func RemoveLabelHeaderFromStream(conn net.Conn) (net.Conn, string, error)

func ValidateKey(key []byte) error

type Address

func (a *Address) String() string

type AliveDelegate

type Broadcast

type ChannelEventDelegate

func (c *ChannelEventDelegate) NotifyJoin(n *Node)

func (c *ChannelEventDelegate) NotifyLeave(n *Node)

func (c *ChannelEventDelegate) NotifyUpdate(n *Node)

type Config

func DefaultLANConfig() *Config

func DefaultLocalConfig() *Config

func DefaultWANConfig() *Config

func (conf *Config) BuildVsnArray() []uint8

func (c *Config) EncryptionEnabled() bool

func (c *Config) IPAllowed(ip net.IP) error

func (c *Config) IPMustBeChecked() bool

type ConflictDelegate

type Delegate

type EventDelegate

type IngestionAwareTransport

deprecated

type Keyring

func NewKeyring(keys [][]byte, primaryKey []byte) (*Keyring, error)

func (k *Keyring) AddKey(key []byte) error

func (k *Keyring) GetKeys() [][]byte

func (k *Keyring) GetPrimaryKey() (key []byte)

func (k *Keyring) RemoveKey(key []byte) error

func (k *Keyring) UseKey(key []byte) error

type Memberlist

func Create(conf *Config) (*Memberlist, error)

func (m *Memberlist) GetHealthScore() int

func (m *Memberlist) Join(existing []string) (int, error)

func (m *Memberlist) Leave(timeout time.Duration) error

func (m *Memberlist) LocalNode() *Node

func (m *Memberlist) Members() []*Node

func (m *Memberlist) NumMembers() (alive int)

func (m *Memberlist) Ping(node string, addr net.Addr) (time.Duration, error)

func (m *Memberlist) ProtocolVersion() uint8

func (m *Memberlist) SendBestEffort(to *Node, msg []byte) error

func (m *Memberlist) SendReliable(to *Node, msg []byte) error

func (m *Memberlist) SendTo(to net.Addr, msg []byte) error

deprecated

func (m *Memberlist) SendToAddress(a Address, msg []byte) error

func (m *Memberlist) SendToTCP(to *Node, msg []byte) error

deprecated

func (m *Memberlist) SendToUDP(to *Node, msg []byte) error

deprecated

func (m *Memberlist) Shutdown() error

func (m *Memberlist) UpdateNode(timeout time.Duration) error

type MergeDelegate

type MockAddress

func (a *MockAddress) Network() string

func (a *MockAddress) String() string

type MockNetwork

func (n *MockNetwork) NewTransport(name string) *MockTransport

type MockTransport

func (t *MockTransport) DialAddressTimeout(a Address, timeout time.Duration) (net.Conn, error)

func (t *MockTransport) DialTimeout(addr string, timeout time.Duration) (net.Conn, error)

func (t *MockTransport) FinalAdvertiseAddr(string, int) (net.IP, int, error)

func (t *MockTransport) IngestPacket(conn net.Conn, addr net.Addr, now time.Time, shouldClose bool) error

func (t *MockTransport) IngestStream(conn net.Conn) error

func (t *MockTransport) PacketCh() <-chan *Packet

func (t *MockTransport) Shutdown() error

func (t *MockTransport) StreamCh() <-chan net.Conn

func (t *MockTransport) WriteTo(b []byte, addr string) (time.Time, error)

func (t *MockTransport) WriteToAddress(b []byte, a Address) (time.Time, error)

type NamedBroadcast

type NetTransport

func NewNetTransport(config *NetTransportConfig) (*NetTransport, error)

func (t *NetTransport) DialAddressTimeout(a Address, timeout time.Duration) (net.Conn, error)

func (t *NetTransport) DialTimeout(addr string, timeout time.Duration) (net.Conn, error)

func (t *NetTransport) FinalAdvertiseAddr(ip string, port int) (net.IP, int, error)

func (t *NetTransport) GetAutoBindPort() int

func (t *NetTransport) IngestPacket(conn net.Conn, addr net.Addr, now time.Time, shouldClose bool) error

func (t *NetTransport) IngestStream(conn net.Conn) error

func (t *NetTransport) PacketCh() <-chan *Packet

func (t *NetTransport) Shutdown() error

func (t *NetTransport) StreamCh() <-chan net.Conn

func (t *NetTransport) WriteTo(b []byte, addr string) (time.Time, error)

func (t *NetTransport) WriteToAddress(b []byte, a Address) (time.Time, error)

type NetTransportConfig

type NoPingResponseError

func (f NoPingResponseError) Error() string

type Node

func (n *Node) Address() string

func (n *Node) FullAddress() Address

func (n *Node) String() string

type NodeAwareTransport

type NodeEvent

type NodeEventType

type NodeStateType

type Packet

type PingDelegate

type TransmitLimitedQueue

func (q *TransmitLimitedQueue) GetBroadcasts(overhead, limit int) [][]byte

func (q *TransmitLimitedQueue) NumQueued() int

func (q *TransmitLimitedQueue) Prune(maxRetain int)

func (q *TransmitLimitedQueue) QueueBroadcast(b Broadcast)

func (q *TransmitLimitedQueue) Reset()

type Transport

type UniqueBroadcast

Constants

¶

View Source

const (

ProtocolVersionMin

uint8

= 1

// Version 3 added support for TCP pings but we kept the default

// protocol version at 2 to ease transition to this new feature.

// A memberlist speaking version 2 of the protocol will attempt

// to TCP ping another memberlist who understands version 3 or

// greater.

//

// Version 4 added support for nacks as part of indirect probes.

// A memberlist speaking version 2 of the protocol will expect

// nacks from another memberlist who understands version 4 or

// greater, and likewise nacks will be sent to memberlists who

// understand version 4 or greater.

ProtocolVersion2Compatible = 2

ProtocolVersionMax = 5

)

This is the minimum and maximum protocol version that we can
_understand_. We're allowed to speak at any version within this
range. This range is inclusive.

View Source

const LabelMaxSize = 255

LabelMaxSize is the maximum length of a packet or stream label.

View Source

const (

MetaMaxSize = 512

// Maximum size for node meta data

)

Variables

¶

This section is empty.

Functions

¶

func

AddLabelHeaderToPacket

¶

added in

v0.3.0

func AddLabelHeaderToPacket(buf []

byte

, label

string

) ([]

byte

,

error

)

AddLabelHeaderToPacket prefixes outgoing packets with the correct header if
the label is not empty.

func

AddLabelHeaderToStream

¶

added in

v0.3.0

func AddLabelHeaderToStream(conn

net

.

Conn

, label

string

)

error

AddLabelHeaderToStream prefixes outgoing streams with the correct header if
the label is not empty.

func

LogAddress

¶

func LogAddress(addr

net

.

Addr

)

string

func

LogConn

¶

func LogConn(conn

net

.

Conn

)

string

func

LogStringAddress

¶

added in

v0.2.0

func LogStringAddress(addr

string

)

string

func

ParseCIDRs

¶

added in

v0.2.1

func ParseCIDRs(v []

string

) ([]

net

.

IPNet

,

error

)

ParseCIDRs return a possible empty list of all Network that have been parsed
In case of error, it returns succesfully parsed CIDRs and the last error found

func

RemoveLabelHeaderFromPacket

¶

added in

v0.3.0

func RemoveLabelHeaderFromPacket(buf []

byte

) (newBuf []

byte

, label

string

, err

error

)

RemoveLabelHeaderFromPacket removes any label header from the provided
packet and returns it along with the remaining packet contents.

func

RemoveLabelHeaderFromStream

¶

added in

v0.3.0

func RemoveLabelHeaderFromStream(conn

net

.

Conn

) (

net

.

Conn

,

string

,

error

)

RemoveLabelHeaderFromStream removes any label header from the beginning of
the stream if present and returns it along with an updated conn with that
header removed.

Note that on error it is the caller's responsibility to close the
connection.

func

ValidateKey

¶

func ValidateKey(key []

byte

)

error

ValidateKey will check to see if the key is valid and returns an error if not.

key should be either 16, 24, or 32 bytes to select AES-128,
AES-192, or AES-256.

Types

¶

type

Address

¶

added in

v0.2.0

type Address struct {

// Addr is a network address as a string, similar to Dial. This usually is

// in the form of "host:port". This is required.

Addr

string

// Name is the name of the node being addressed. This is optional but

// transports may require it.

Name

string

}

func (*Address)

String

¶

added in

v0.2.0

func (a *

Address

) String()

string

type

AliveDelegate

¶

type AliveDelegate interface {

// NotifyAlive is invoked when a message about a live

// node is received from the network.  Returning a non-nil

// error prevents the node from being considered a peer.

NotifyAlive(peer *

Node

)

error

}

AliveDelegate is used to involve a client in processing
a node "alive" message. When a node joins, either through
a UDP gossip or TCP push/pull, we update the state of
that node via an alive message. This can be used to filter
a node out and prevent it from being considered a peer
using application specific logic.

type

Broadcast

¶

type Broadcast interface {

// Invalidates checks if enqueuing the current broadcast

// invalidates a previous broadcast

Invalidates(b

Broadcast

)

bool

// Returns a byte form of the message

Message() []

byte

// Finished is invoked when the message will no longer

// be broadcast, either due to invalidation or to the

// transmit limit being reached

Finished()
}

Broadcast is something that can be broadcasted via gossip to
the memberlist cluster.

type

ChannelEventDelegate

¶

type ChannelEventDelegate struct {

Ch chan<-

NodeEvent

}

ChannelEventDelegate is used to enable an application to receive
events about joins and leaves over a channel instead of a direct
function call.

Care must be taken that events are processed in a timely manner from
the channel, since this delegate will block until an event can be sent.

func (*ChannelEventDelegate)

NotifyJoin

¶

func (c *

ChannelEventDelegate

) NotifyJoin(n *

Node

)

func (*ChannelEventDelegate)

NotifyLeave

¶

func (c *

ChannelEventDelegate

) NotifyLeave(n *

Node

)

func (*ChannelEventDelegate)

NotifyUpdate

¶

func (c *

ChannelEventDelegate

) NotifyUpdate(n *

Node

)

type

Config

¶

type Config struct {

// The name of this node. This must be unique in the cluster.

Name

string

// Transport is a hook for providing custom code to communicate with

// other nodes. If this is left nil, then memberlist will by default

// make a NetTransport using BindAddr and BindPort from this structure.

Transport

Transport

// Label is an optional set of bytes to include on the outside of each

// packet and stream.

//

// If gossip encryption is enabled and this is set it is treated as GCM

// authenticated data.

Label

string

// SkipInboundLabelCheck skips the check that inbound packets and gossip

// streams need to be label prefixed.

SkipInboundLabelCheck

bool

// Configuration related to what address to bind to and ports to

// listen on. The port is used for both UDP and TCP gossip. It is

// assumed other nodes are running on this port, but they do not need

// to.

BindAddr

string

BindPort

int

// Configuration related to what address to advertise to other

// cluster members. Used for nat traversal.

AdvertiseAddr

string

AdvertisePort

int

// ProtocolVersion is the configured protocol version that we

// will _speak_. This must be between ProtocolVersionMin and

// ProtocolVersionMax.

ProtocolVersion

uint8

// TCPTimeout is the timeout for establishing a stream connection with

// a remote node for a full state sync, and for stream read and write

// operations. This is a legacy name for backwards compatibility, but

// should really be called StreamTimeout now that we have generalized

// the transport.

TCPTimeout

time

.

Duration

// IndirectChecks is the number of nodes that will be asked to perform

// an indirect probe of a node in the case a direct probe fails. Memberlist

// waits for an ack from any single indirect node, so increasing this

// number will increase the likelihood that an indirect probe will succeed

// at the expense of bandwidth.

IndirectChecks

int

// RetransmitMult is the multiplier for the number of retransmissions

// that are attempted for messages broadcasted over gossip. The actual

// count of retransmissions is calculated using the formula:

//

//   Retransmits = RetransmitMult * log(N+1)

//

// This allows the retransmits to scale properly with cluster size. The

// higher the multiplier, the more likely a failed broadcast is to converge

// at the expense of increased bandwidth.

RetransmitMult

int

// SuspicionMult is the multiplier for determining the time an

// inaccessible node is considered suspect before declaring it dead.

// The actual timeout is calculated using the formula:

//

//   SuspicionTimeout = SuspicionMult * log(N+1) * ProbeInterval

//

// This allows the timeout to scale properly with expected propagation

// delay with a larger cluster size. The higher the multiplier, the longer

// an inaccessible node is considered part of the cluster before declaring

// it dead, giving that suspect node more time to refute if it is indeed

// still alive.

SuspicionMult

int

// SuspicionMaxTimeoutMult is the multiplier applied to the

// SuspicionTimeout used as an upper bound on detection time. This max

// timeout is calculated using the formula:

//

// SuspicionMaxTimeout = SuspicionMaxTimeoutMult * SuspicionTimeout

//

// If everything is working properly, confirmations from other nodes will

// accelerate suspicion timers in a manner which will cause the timeout

// to reach the base SuspicionTimeout before that elapses, so this value

// will typically only come into play if a node is experiencing issues

// communicating with other nodes. It should be set to a something fairly

// large so that a node having problems will have a lot of chances to

// recover before falsely declaring other nodes as failed, but short

// enough for a legitimately isolated node to still make progress marking

// nodes failed in a reasonable amount of time.

SuspicionMaxTimeoutMult

int

// PushPullInterval is the interval between complete state syncs.

// Complete state syncs are done with a single node over TCP and are

// quite expensive relative to standard gossiped messages. Setting this

// to zero will disable state push/pull syncs completely.

//

// Setting this interval lower (more frequent) will increase convergence

// speeds across larger clusters at the expense of increased bandwidth

// usage.

PushPullInterval

time

.

Duration

// ProbeInterval and ProbeTimeout are used to configure probing

// behavior for memberlist.

//

// ProbeInterval is the interval between random node probes. Setting

// this lower (more frequent) will cause the memberlist cluster to detect

// failed nodes more quickly at the expense of increased bandwidth usage.

//

// ProbeTimeout is the timeout to wait for an ack from a probed node

// before assuming it is unhealthy. This should be set to 99-percentile

// of RTT (round-trip time) on your network.

ProbeInterval

time

.

Duration

ProbeTimeout

time

.

Duration

// DisableTcpPings will turn off the fallback TCP pings that are attempted

// if the direct UDP ping fails. These get pipelined along with the

// indirect UDP pings.

DisableTcpPings

bool

// DisableTcpPingsForNode is like DisableTcpPings, but lets you control

// whether to perform TCP pings on a node-by-node basis.

DisableTcpPingsForNode func(nodeName

string

)

bool

// AwarenessMaxMultiplier will increase the probe interval if the node

// becomes aware that it might be degraded and not meeting the soft real

// time requirements to reliably probe other nodes.

AwarenessMaxMultiplier

int

// GossipInterval and GossipNodes are used to configure the gossip

// behavior of memberlist.

//

// GossipInterval is the interval between sending messages that need

// to be gossiped that haven't been able to piggyback on probing messages.

// If this is set to zero, non-piggyback gossip is disabled. By lowering

// this value (more frequent) gossip messages are propagated across

// the cluster more quickly at the expense of increased bandwidth.

//

// GossipNodes is the number of random nodes to send gossip messages to

// per GossipInterval. Increasing this number causes the gossip messages

// to propagate across the cluster more quickly at the expense of

// increased bandwidth.

//

// GossipToTheDeadTime is the interval after which a node has died that

// we will still try to gossip to it. This gives it a chance to refute.

GossipInterval

time

.

Duration

GossipNodes

int

GossipToTheDeadTime

time

.

Duration

// GossipVerifyIncoming controls whether to enforce encryption for incoming

// gossip. It is used for upshifting from unencrypted to encrypted gossip on

// a running cluster.

GossipVerifyIncoming

bool

// GossipVerifyOutgoing controls whether to enforce encryption for outgoing

// gossip. It is used for upshifting from unencrypted to encrypted gossip on

// a running cluster.

GossipVerifyOutgoing

bool

// EnableCompression is used to control message compression. This can

// be used to reduce bandwidth usage at the cost of slightly more CPU

// utilization. This is only available starting at protocol version 1.

EnableCompression

bool

// SecretKey is used to initialize the primary encryption key in a keyring.

// The primary encryption key is the only key used to encrypt messages and

// the first key used while attempting to decrypt messages. Providing a

// value for this primary key will enable message-level encryption and

// verification, and automatically install the key onto the keyring.

// The value should be either 16, 24, or 32 bytes to select AES-128,

// AES-192, or AES-256.

SecretKey []

byte

// The keyring holds all of the encryption keys used internally. It is

// automatically initialized using the SecretKey and SecretKeys values.

Keyring *

Keyring

// Delegate and Events are delegates for receiving and providing

// data to memberlist via callback mechanisms. For Delegate, see

// the Delegate interface. For Events, see the EventDelegate interface.

//

// The DelegateProtocolMin/Max are used to guarantee protocol-compatibility

// for any custom messages that the delegate might do (broadcasts,

// local/remote state, etc.). If you don't set these, then the protocol

// versions will just be zero, and version compliance won't be done.

Delegate

Delegate

DelegateProtocolVersion

uint8

DelegateProtocolMin

uint8

DelegateProtocolMax

uint8

Events

EventDelegate

Conflict

ConflictDelegate

Merge

MergeDelegate

Ping

PingDelegate

Alive

AliveDelegate

// DNSConfigPath points to the system's DNS config file, usually located

// at /etc/resolv.conf. It can be overridden via config for easier testing.

DNSConfigPath

string

// LogOutput is the writer where logs should be sent. If this is not

// set, logging will go to stderr by default. You cannot specify both LogOutput

// and Logger at the same time.

LogOutput

io

.

Writer

// Logger is a custom logger which you provide. If Logger is set, it will use

// this for the internal logger. If Logger is not set, it will fall back to the

// behavior for using LogOutput. You cannot specify both LogOutput and Logger

// at the same time.

Logger *

log

.

Logger

// Size of Memberlist's internal channel which handles UDP messages. The

// size of this determines the size of the queue which Memberlist will keep

// while UDP messages are handled.

HandoffQueueDepth

int

// Maximum number of bytes that memberlist will put in a packet (this

// will be for UDP packets by default with a NetTransport). A safe value

// for this is typically 1400 bytes (which is the default). However,

// depending on your network's MTU (Maximum Transmission Unit) you may

// be able to increase this to get more content into each gossip packet.

// This is a legacy name for backward compatibility but should really be

// called PacketBufferSize now that we have generalized the transport.

UDPBufferSize

int

// DeadNodeReclaimTime controls the time before a dead node's name can be

// reclaimed by one with a different address or port. By default, this is 0,

// meaning nodes cannot be reclaimed this way.

DeadNodeReclaimTime

time

.

Duration

// RequireNodeNames controls if the name of a node is required when sending

// a message to that node.

RequireNodeNames

bool

// CIDRsAllowed If nil, allow any connection (default), otherwise specify all networks

// allowed to connect (you must specify IPv6/IPv4 separately)

// Using [] will block all connections.

CIDRsAllowed []

net

.

IPNet

// MetricLabels is a map of optional labels to apply to all metrics emitted.

MetricLabels []metrics.Label

// QueueCheckInterval is the interval at which we check the message

// queue to apply the warning and max depth.

QueueCheckInterval

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

func

DefaultLANConfig

¶

func DefaultLANConfig() *

Config

DefaultLANConfig returns a sane set of configurations for Memberlist.
It uses the hostname as the node name, and otherwise sets very conservative
values that are sane for most LAN environments. The default configuration
errs on the side of caution, choosing values that are optimized
for higher convergence at the cost of higher bandwidth usage. Regardless,
these values are a good starting point when getting started with memberlist.

func

DefaultLocalConfig

¶

func DefaultLocalConfig() *

Config

DefaultLocalConfig works like DefaultConfig, however it returns a configuration
that is optimized for a local loopback environments. The default configuration is
still very conservative and errs on the side of caution.

func

DefaultWANConfig

¶

func DefaultWANConfig() *

Config

DefaultWANConfig works like DefaultConfig, however it returns a configuration
that is optimized for most WAN environments. The default configuration is
still very conservative and errs on the side of caution.

func (*Config)

BuildVsnArray

¶

added in

v0.1.4

func (conf *

Config

) BuildVsnArray() []

uint8

BuildVsnArray creates the array of Vsn

func (*Config)

EncryptionEnabled

¶

func (c *

Config

) EncryptionEnabled()

bool

Returns whether or not encryption is enabled

func (*Config)

IPAllowed

¶

added in

v0.2.1

func (c *

Config

) IPAllowed(ip

net

.

IP

)

error

IPAllowed return an error if access to memberlist is denied

func (*Config)

IPMustBeChecked

¶

added in

v0.2.1

func (c *

Config

) IPMustBeChecked()

bool

IPMustBeChecked return true if IPAllowed must be called

type

ConflictDelegate

¶

type ConflictDelegate interface {

// NotifyConflict is invoked when a name conflict is detected

NotifyConflict(existing, other *

Node

)
}

ConflictDelegate is a used to inform a client that
a node has attempted to join which would result in a
name conflict. This happens if two clients are configured
with the same name but different addresses.

type

Delegate

¶

type Delegate interface {

// NodeMeta is used to retrieve meta-data about the current node

// when broadcasting an alive message. It's length is limited to

// the given byte size. This metadata is available in the Node structure.

NodeMeta(limit

int

) []

byte

// NotifyMsg is called when a user-data message is received.

// Care should be taken that this method does not block, since doing

// so would block the entire UDP packet receive loop. Additionally, the byte

// slice may be modified after the call returns, so it should be copied if needed

NotifyMsg([]

byte

)

// GetBroadcasts is called when user data messages can be broadcast.

// It can return a list of buffers to send. Each buffer should assume an

// overhead as provided with a limit on the total byte size allowed.

// The total byte size of the resulting data to send must not exceed

// the limit. Care should be taken that this method does not block,

// since doing so would block the entire UDP packet receive loop.

GetBroadcasts(overhead, limit

int

) [][]

byte

// LocalState is used for a TCP Push/Pull. This is sent to

// the remote side in addition to the membership information. Any

// data can be sent here. See MergeRemoteState as well. The `join`

// boolean indicates this is for a join instead of a push/pull.

LocalState(join

bool

) []

byte

// MergeRemoteState is invoked after a TCP Push/Pull. This is the

// state received from the remote side and is the result of the

// remote side's LocalState call. The 'join'

// boolean indicates this is for a join instead of a push/pull.

MergeRemoteState(buf []

byte

, join

bool

)
}

Delegate is the interface that clients must implement if they want to hook
into the gossip layer of Memberlist. All the methods must be thread-safe,
as they can and generally will be called concurrently.

type

EventDelegate

¶

type EventDelegate interface {

// NotifyJoin is invoked when a node is detected to have joined.

// The Node argument must not be modified.

NotifyJoin(*

Node

)

// NotifyLeave is invoked when a node is detected to have left.

// The Node argument must not be modified.

NotifyLeave(*

Node

)

// NotifyUpdate is invoked when a node is detected to have

// updated, usually involving the meta data. The Node argument

// must not be modified.

NotifyUpdate(*

Node

)
}

EventDelegate is a simpler delegate that is used only to receive
notifications about members joining and leaving. The methods in this
delegate may be called by multiple goroutines, but never concurrently.
This allows you to reason about ordering.

type

IngestionAwareTransport

deprecated

added in

v0.2.0

type IngestionAwareTransport interface {

IngestPacket(conn

net

.

Conn

, addr

net

.

Addr

, now

time

.

Time

, shouldClose

bool

)

error

IngestStream(conn

net

.

Conn

)

error

}

IngestionAwareTransport is not used.

Deprecated: IngestionAwareTransport is not used and may be removed in a future
version. Define the interface locally instead of referencing this exported
interface.

type

Keyring

¶

type Keyring struct {

// contains filtered or unexported fields

}

func

NewKeyring

¶

func NewKeyring(keys [][]

byte

, primaryKey []

byte

) (*

Keyring

,

error

)

NewKeyring constructs a new container for a set of encryption keys. The
keyring contains all key data used internally by memberlist.

While creating a new keyring, you must do one of:

Omit keys and primary key, effectively disabling encryption

Pass a set of keys plus the primary key

Pass only a primary key

If only a primary key is passed, then it will be automatically added to the
keyring. If creating a keyring with multiple keys, one key must be designated
primary by passing it as the primaryKey. If the primaryKey does not exist in
the list of secondary keys, it will be automatically added at position 0.

A key should be either 16, 24, or 32 bytes to select AES-128,
AES-192, or AES-256.

func (*Keyring)

AddKey

¶

func (k *

Keyring

) AddKey(key []

byte

)

error

AddKey will install a new key on the ring. Adding a key to the ring will make
it available for use in decryption. If the key already exists on the ring,
this function will just return noop.

key should be either 16, 24, or 32 bytes to select AES-128,
AES-192, or AES-256.

func (*Keyring)

GetKeys

¶

func (k *

Keyring

) GetKeys() [][]

byte

GetKeys returns the current set of keys on the ring.

func (*Keyring)

GetPrimaryKey

¶

func (k *

Keyring

) GetPrimaryKey() (key []

byte

)

GetPrimaryKey returns the key on the ring at position 0. This is the key used
for encrypting messages, and is the first key tried for decrypting messages.

func (*Keyring)

RemoveKey

¶

func (k *

Keyring

) RemoveKey(key []

byte

)

error

RemoveKey drops a key from the keyring. This will return an error if the key
requested for removal is currently at position 0 (primary key).

func (*Keyring)

UseKey

¶

func (k *

Keyring

) UseKey(key []

byte

)

error

UseKey changes the key used to encrypt messages. This is the only key used to
encrypt messages, so peers should know this key before this method is called.

type

Memberlist

¶

type Memberlist struct {

// contains filtered or unexported fields

}

func

Create

¶

func Create(conf *

Config

) (*

Memberlist

,

error

)

Create will create a new Memberlist using the given configuration.
This will not connect to any other node (see Join) yet, but will start
all the listeners to allow other nodes to join this memberlist.
After creating a Memberlist, the configuration given should not be
modified by the user anymore.

func (*Memberlist)

GetHealthScore

¶

func (m *

Memberlist

) GetHealthScore()

int

GetHealthScore gives this instance's idea of how well it is meeting the soft
real-time requirements of the protocol. Lower numbers are better, and zero
means "totally healthy".

func (*Memberlist)

Join

¶

func (m *

Memberlist

) Join(existing []

string

) (

int

,

error

)

Join is used to take an existing Memberlist and attempt to join a cluster
by contacting all the given hosts and performing a state sync. Initially,
the Memberlist only contains our own state, so doing this will cause
remote nodes to become aware of the existence of this node, effectively
joining the cluster.

This returns the number of hosts successfully contacted and an error if
none could be reached. If an error is returned, the node did not successfully
join the cluster.

func (*Memberlist)

Leave

¶

func (m *

Memberlist

) Leave(timeout

time

.

Duration

)

error

Leave will broadcast a leave message but will not shutdown the background
listeners, meaning the node will continue participating in gossip and state
updates.

This will block until the leave message is successfully broadcasted to
a member of the cluster, if any exist or until a specified timeout
is reached.

This method is safe to call multiple times, but must not be called
after the cluster is already shut down.

func (*Memberlist)

LocalNode

¶

func (m *

Memberlist

) LocalNode() *

Node

LocalNode is used to return the local Node

func (*Memberlist)

Members

¶

func (m *

Memberlist

) Members() []*

Node

Members returns a list of all known live nodes. The node structures
returned must not be modified. If you wish to modify a Node, make a
copy first.

func (*Memberlist)

NumMembers

¶

func (m *

Memberlist

) NumMembers() (alive

int

)

NumMembers returns the number of alive nodes currently known. Between
the time of calling this and calling Members, the number of alive nodes
may have changed, so this shouldn't be used to determine how many
members will be returned by Members.

func (*Memberlist)

Ping

¶

func (m *

Memberlist

) Ping(node

string

, addr

net

.

Addr

) (

time

.

Duration

,

error

)

Ping initiates a ping to the node with the specified name.

func (*Memberlist)

ProtocolVersion

¶

func (m *

Memberlist

) ProtocolVersion()

uint8

ProtocolVersion returns the protocol version currently in use by
this memberlist.

func (*Memberlist)

SendBestEffort

¶

func (m *

Memberlist

) SendBestEffort(to *

Node

, msg []

byte

)

error

SendBestEffort uses the unreliable packet-oriented interface of the transport
to target a user message at the given node (this does not use the gossip
mechanism). The maximum size of the message depends on the configured
UDPBufferSize for this memberlist instance.

func (*Memberlist)

SendReliable

¶

func (m *

Memberlist

) SendReliable(to *

Node

, msg []

byte

)

error

SendReliable uses the reliable stream-oriented interface of the transport to
target a user message at the given node (this does not use the gossip
mechanism). Delivery is guaranteed if no error is returned, and there is no
limit on the size of the message.

func (*Memberlist)

SendTo

deprecated

func (m *

Memberlist

) SendTo(to

net

.

Addr

, msg []

byte

)

error

Deprecated: SendTo is deprecated in favor of SendBestEffort, which requires a node to
target. If you don't have a node then use SendToAddress.

func (*Memberlist)

SendToAddress

¶

added in

v0.2.0

func (m *

Memberlist

) SendToAddress(a

Address

, msg []

byte

)

error

func (*Memberlist)

SendToTCP

deprecated

func (m *

Memberlist

) SendToTCP(to *

Node

, msg []

byte

)

error

Deprecated: SendToTCP is deprecated in favor of SendReliable.

func (*Memberlist)

SendToUDP

deprecated

func (m *

Memberlist

) SendToUDP(to *

Node

, msg []

byte

)

error

Deprecated: SendToUDP is deprecated in favor of SendBestEffort.

func (*Memberlist)

Shutdown

¶

func (m *

Memberlist

) Shutdown()

error

Shutdown will stop any background maintenance of network activity
for this memberlist, causing it to appear "dead". A leave message
will not be broadcasted prior, so the cluster being left will have
to detect this node's shutdown using probing. If you wish to more
gracefully exit the cluster, call Leave prior to shutting down.

This method is safe to call multiple times.

func (*Memberlist)

UpdateNode

¶

func (m *

Memberlist

) UpdateNode(timeout

time

.

Duration

)

error

UpdateNode is used to trigger re-advertising the local node. This is
primarily used with a Delegate to support dynamic updates to the local
meta data.  This will block until the update message is successfully
broadcasted to a member of the cluster, if any exist or until a specified
timeout is reached.

type

MergeDelegate

¶

type MergeDelegate interface {

// NotifyMerge is invoked when a merge could take place.

// Provides a list of the nodes known by the peer. If

// the return value is non-nil, the merge is canceled.

NotifyMerge(peers []*

Node

)

error

}

MergeDelegate is used to involve a client in
a potential cluster merge operation. Namely, when
a node does a TCP push/pull (as part of a join),
the delegate is involved and allowed to cancel the join
based on custom logic. The merge delegate is NOT invoked
as part of the push-pull anti-entropy.

type

MockAddress

¶

type MockAddress struct {

// contains filtered or unexported fields

}

MockAddress is a wrapper which adds the net.Addr interface to our mock
address scheme.

func (*MockAddress)

Network

¶

func (a *

MockAddress

) Network()

string

See net.Addr.

func (*MockAddress)

String

¶

func (a *

MockAddress

) String()

string

See net.Addr.

type

MockNetwork

¶

type MockNetwork struct {

// contains filtered or unexported fields

}

MockNetwork is used as a factory that produces MockTransport instances which
are uniquely addressed and wired up to talk to each other.

func (*MockNetwork)

NewTransport

¶

func (n *

MockNetwork

) NewTransport(name

string

) *

MockTransport

NewTransport returns a new MockTransport with a unique address, wired up to
talk to the other transports in the MockNetwork.

type

MockTransport

¶

type MockTransport struct {

// contains filtered or unexported fields

}

MockTransport directly plumbs messages to other transports its MockNetwork.

func (*MockTransport)

DialAddressTimeout

¶

added in

v0.2.0

func (t *

MockTransport

) DialAddressTimeout(a

Address

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

See NodeAwareTransport.

func (*MockTransport)

DialTimeout

¶

func (t *

MockTransport

) DialTimeout(addr

string

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

See Transport.

func (*MockTransport)

FinalAdvertiseAddr

¶

func (t *

MockTransport

) FinalAdvertiseAddr(

string

,

int

) (

net

.

IP

,

int

,

error

)

See Transport.

func (*MockTransport)

IngestPacket

¶

added in

v0.2.0

func (t *

MockTransport

) IngestPacket(conn

net

.

Conn

, addr

net

.

Addr

, now

time

.

Time

, shouldClose

bool

)

error

See NodeAwareTransport.

func (*MockTransport)

IngestStream

¶

added in

v0.2.0

func (t *

MockTransport

) IngestStream(conn

net

.

Conn

)

error

See NodeAwareTransport.

func (*MockTransport)

PacketCh

¶

func (t *

MockTransport

) PacketCh() <-chan *

Packet

See Transport.

func (*MockTransport)

Shutdown

¶

func (t *

MockTransport

) Shutdown()

error

See Transport.

func (*MockTransport)

StreamCh

¶

func (t *

MockTransport

) StreamCh() <-chan

net

.

Conn

See Transport.

func (*MockTransport)

WriteTo

¶

func (t *

MockTransport

) WriteTo(b []

byte

, addr

string

) (

time

.

Time

,

error

)

See Transport.

func (*MockTransport)

WriteToAddress

¶

added in

v0.2.0

func (t *

MockTransport

) WriteToAddress(b []

byte

, a

Address

) (

time

.

Time

,

error

)

See NodeAwareTransport.

type

NamedBroadcast

¶

added in

v0.1.1

type NamedBroadcast interface {

Broadcast

// The unique identity of this broadcast message.

Name()

string

}

NamedBroadcast is an optional extension of the Broadcast interface that
gives each message a unique string name, and that is used to optimize

You shoud ensure that Invalidates() checks the same uniqueness as the
example below:

func (b *foo) Invalidates(other Broadcast) bool {
	nb, ok := other.(NamedBroadcast)
	if !ok {
		return false
	}
	return b.Name() == nb.Name()
}

Invalidates() isn't currently used for NamedBroadcasts, but that may change
in the future.

type

NetTransport

¶

type NetTransport struct {

// contains filtered or unexported fields

}

NetTransport is a Transport implementation that uses connectionless UDP for
packet operations, and ad-hoc TCP connections for stream operations.

func

NewNetTransport

¶

func NewNetTransport(config *

NetTransportConfig

) (*

NetTransport

,

error

)

NewNetTransport returns a net transport with the given configuration. On
success all the network listeners will be created and listening.

func (*NetTransport)

DialAddressTimeout

¶

added in

v0.2.0

func (t *

NetTransport

) DialAddressTimeout(a

Address

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

See NodeAwareTransport.

func (*NetTransport)

DialTimeout

¶

func (t *

NetTransport

) DialTimeout(addr

string

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

See Transport.

func (*NetTransport)

FinalAdvertiseAddr

¶

func (t *

NetTransport

) FinalAdvertiseAddr(ip

string

, port

int

) (

net

.

IP

,

int

,

error

)

See Transport.

func (*NetTransport)

GetAutoBindPort

¶

func (t *

NetTransport

) GetAutoBindPort()

int

GetAutoBindPort returns the bind port that was automatically given by the
kernel, if a bind port of 0 was given.

func (*NetTransport)

IngestPacket

¶

added in

v0.2.0

func (t *

NetTransport

) IngestPacket(conn

net

.

Conn

, addr

net

.

Addr

, now

time

.

Time

, shouldClose

bool

)

error

See IngestionAwareTransport.

func (*NetTransport)

IngestStream

¶

added in

v0.2.0

func (t *

NetTransport

) IngestStream(conn

net

.

Conn

)

error

See IngestionAwareTransport.

func (*NetTransport)

PacketCh

¶

func (t *

NetTransport

) PacketCh() <-chan *

Packet

See Transport.

func (*NetTransport)

Shutdown

¶

func (t *

NetTransport

) Shutdown()

error

See Transport.

func (*NetTransport)

StreamCh

¶

func (t *

NetTransport

) StreamCh() <-chan

net

.

Conn

See Transport.

func (*NetTransport)

WriteTo

¶

func (t *

NetTransport

) WriteTo(b []

byte

, addr

string

) (

time

.

Time

,

error

)

See Transport.

func (*NetTransport)

WriteToAddress

¶

added in

v0.2.0

func (t *

NetTransport

) WriteToAddress(b []

byte

, a

Address

) (

time

.

Time

,

error

)

See NodeAwareTransport.

type

NetTransportConfig

¶

type NetTransportConfig struct {

// BindAddrs is a list of addresses to bind to for both TCP and UDP

// communications.

BindAddrs []

string

// BindPort is the port to listen on, for each address above.

BindPort

int

// Logger is a logger for operator messages.

Logger *

log

.

Logger

// MetricLabels is a map of optional labels to apply to all metrics

// emitted by this transport.

MetricLabels []

metrics

.

Label

}

NetTransportConfig is used to configure a net transport.

type

NoPingResponseError

¶

type NoPingResponseError struct {

// contains filtered or unexported fields

}

NoPingResponseError is used to indicate a 'ping' packet was
successfully issued but no response was received

func (NoPingResponseError)

Error

¶

func (f

NoPingResponseError

) Error()

string

type

Node

¶

type Node struct {

Name

string

Addr

net

.

IP

Port

uint16

Meta  []

byte

// Metadata from the delegate for this node.

State

NodeStateType

// State of the node.

PMin

uint8

// Minimum protocol version this understands

PMax

uint8

// Maximum protocol version this understands

PCur

uint8

// Current version node is speaking

DMin

uint8

// Min protocol version for the delegate to understand

DMax

uint8

// Max protocol version for the delegate to understand

DCur

uint8

// Current version delegate is speaking

}

Node represents a node in the cluster.

func (*Node)

Address

¶

func (n *

Node

) Address()

string

Address returns the host:port form of a node's address, suitable for use
with a transport.

func (*Node)

FullAddress

¶

added in

v0.2.0

func (n *

Node

) FullAddress()

Address

FullAddress returns the node name and host:port form of a node's address,
suitable for use with a transport.

func (*Node)

String

¶

added in

v0.1.1

func (n *

Node

) String()

string

String returns the node name

type

NodeAwareTransport

¶

added in

v0.2.0

type NodeAwareTransport interface {

Transport

WriteToAddress(b []

byte

, addr

Address

) (

time

.

Time

,

error

)

DialAddressTimeout(addr

Address

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

type

NodeEvent

¶

type NodeEvent struct {

Event

NodeEventType

Node  *

Node

}

NodeEvent is a single event related to node activity in the memberlist.
The Node member of this struct must not be directly modified. It is passed
as a pointer to avoid unnecessary copies. If you wish to modify the node,
make a copy first.

type

NodeEventType

¶

type NodeEventType

int

NodeEventType are the types of events that can be sent from the
ChannelEventDelegate.

const (

NodeJoin

NodeEventType

=

iota

NodeLeave

NodeUpdate

)

type

NodeStateType

¶

added in

v0.2.1

type NodeStateType

int

const (

StateAlive

NodeStateType

=

iota

StateSuspect

StateDead

StateLeft

)

type

Packet

¶

type Packet struct {

// Buf has the raw contents of the packet.

Buf []

byte

// From has the address of the peer. This is an actual net.Addr so we

// can expose some concrete details about incoming packets.

From

net

.

Addr

// Timestamp is the time when the packet was received. This should be

// taken as close as possible to the actual receipt time to help make an

// accurate RTT measurement during probes.

Timestamp

time

.

Time

}

Packet is used to provide some metadata about incoming packets from peers
over a packet connection, as well as the packet payload.

type

PingDelegate

¶

type PingDelegate interface {

// AckPayload is invoked when an ack is being sent; the returned bytes will be appended to the ack

AckPayload() []

byte

// NotifyPing is invoked when an ack for a ping is received

NotifyPingComplete(other *

Node

, rtt

time

.

Duration

, payload []

byte

)
}

PingDelegate is used to notify an observer how long it took for a ping message to
complete a round trip.  It can also be used for writing arbitrary byte slices
into ack messages. Note that in order to be meaningful for RTT estimates, this
delegate does not apply to indirect pings, nor fallback pings sent over TCP.

type

TransmitLimitedQueue

¶

type TransmitLimitedQueue struct {

// NumNodes returns the number of nodes in the cluster. This is

// used to determine the retransmit count, which is calculated

// based on the log of this.

NumNodes func()

int

// RetransmitMult is the multiplier used to determine the maximum

// number of retransmissions attempted.

RetransmitMult

int

// contains filtered or unexported fields

}

TransmitLimitedQueue is used to queue messages to broadcast to
the cluster (via gossip) but limits the number of transmits per
message. It also prioritizes messages with lower transmit counts
(hence newer messages).

func (*TransmitLimitedQueue)

GetBroadcasts

¶

func (q *

TransmitLimitedQueue

) GetBroadcasts(overhead, limit

int

) [][]

byte

GetBroadcasts is used to get a number of broadcasts, up to a byte limit
and applying a per-message overhead as provided.

func (*TransmitLimitedQueue)

NumQueued

¶

func (q *

TransmitLimitedQueue

) NumQueued()

int

NumQueued returns the number of queued messages

func (*TransmitLimitedQueue)

Prune

¶

func (q *

TransmitLimitedQueue

) Prune(maxRetain

int

)

Prune will retain the maxRetain latest messages, and the rest
will be discarded. This can be used to prevent unbounded queue sizes

func (*TransmitLimitedQueue)

QueueBroadcast

¶

func (q *

TransmitLimitedQueue

) QueueBroadcast(b

Broadcast

)

QueueBroadcast is used to enqueue a broadcast

func (*TransmitLimitedQueue)

Reset

¶

func (q *

TransmitLimitedQueue

) Reset()

Reset clears all the queued messages. Should only be used for tests.

type

Transport

¶

type Transport interface {

// FinalAdvertiseAddr is given the user's configured values (which

// might be empty) and returns the desired IP and port to advertise to

// the rest of the cluster.

FinalAdvertiseAddr(ip

string

, port

int

) (

net

.

IP

,

int

,

error

)

// WriteTo is a packet-oriented interface that fires off the given

// payload to the given address in a connectionless fashion. This should

// return a time stamp that's as close as possible to when the packet

// was transmitted to help make accurate RTT measurements during probes.

//

// This is similar to net.PacketConn, though we didn't want to expose

// that full set of required methods to keep assumptions about the

// underlying plumbing to a minimum. We also treat the address here as a

// string, similar to Dial, so it's network neutral, so this usually is

// in the form of "host:port".

WriteTo(b []

byte

, addr

string

) (

time

.

Time

,

error

)

// PacketCh returns a channel that can be read to receive incoming

// packets from other peers. How this is set up for listening is left as

// an exercise for the concrete transport implementations.

PacketCh() <-chan *

Packet

// DialTimeout is used to create a connection that allows us to perform

// two-way communication with a peer. This is generally more expensive

// than packet connections so is used for more infrequent operations

// such as anti-entropy or fallback probes if the packet-oriented probe

// failed.

DialTimeout(addr

string

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

// StreamCh returns a channel that can be read to handle incoming stream

// connections from other peers. How this is set up for listening is

// left as an exercise for the concrete transport implementations.

StreamCh() <-chan

net

.

Conn

// Shutdown is called when memberlist is shutting down; this gives the

// transport a chance to clean up any listeners.

Shutdown()

error

}

Transport is used to abstract over communicating with other peers. The packet
interface is assumed to be best-effort and the stream interface is assumed to
be reliable.

type

UniqueBroadcast

¶

added in

v0.1.1

type UniqueBroadcast interface {

Broadcast

// UniqueBroadcast is just a marker method for this interface.

UniqueBroadcast()
}

UniqueBroadcast is an optional interface that indicates that each message is
intrinsically unique and there is no need to scan the broadcast queue for
duplicates.

You should ensure that Invalidates() always returns false if implementing
this interface. Invalidates() isn't currently used for UniqueBroadcasts, but
that may change in the future.