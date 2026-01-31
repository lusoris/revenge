# pgxpool Connection Pool

> Source: https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool
> Fetched: 2026-01-30T23:50:11.235371+00:00
> Content-Hash: db6cb2c6bcf0a9db
> Type: html

---

Overview

¶

Creating a Pool

Package pgxpool is a concurrency-safe connection pool for pgx.

pgxpool implements a nearly identical interface to pgx connections.

Creating a Pool

¶

The primary way of creating a pool is with

pgxpool.New

:

pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

The database connection string can be in URL or keyword/value format. PostgreSQL settings, pgx settings, and pool settings can be
specified here. In addition, a config struct can be created by

ParseConfig

.

config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
if err != nil {
    // ...
}
config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
    // do something with every new connection
}

pool, err := pgxpool.NewWithConfig(context.Background(), config)

A pool returns without waiting for any connections to be established. Acquire a connection immediately after creating
the pool to check if a connection can successfully be established.

Index

¶

type AcquireTracer

type Config

func ParseConfig(connString string) (*Config, error)

func (c *Config) ConnString() string

func (c *Config) Copy() *Config

type Conn

func (c *Conn) Begin(ctx context.Context) (pgx.Tx, error)

func (c *Conn) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)

func (c *Conn) Conn() *pgx.Conn

func (c *Conn) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, ...) (int64, error)

func (c *Conn) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

func (c *Conn) Hijack() *pgx.Conn

func (c *Conn) Ping(ctx context.Context) error

func (c *Conn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

func (c *Conn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

func (c *Conn) Release()

func (c *Conn) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

type Pool

func New(ctx context.Context, connString string) (*Pool, error)

func NewWithConfig(ctx context.Context, config *Config) (*Pool, error)

func (p *Pool) Acquire(ctx context.Context) (c *Conn, err error)

func (p *Pool) AcquireAllIdle(ctx context.Context) []*Conn

func (p *Pool) AcquireFunc(ctx context.Context, f func(*Conn) error) error

func (p *Pool) Begin(ctx context.Context) (pgx.Tx, error)

func (p *Pool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)

func (p *Pool) Close()

func (p *Pool) Config() *Config

func (p *Pool) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, ...) (int64, error)

func (p *Pool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

func (p *Pool) Ping(ctx context.Context) error

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

func (p *Pool) Reset()

func (p *Pool) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

func (p *Pool) Stat() *Stat

type ReleaseTracer

type ShouldPingParams

type Stat

func (s *Stat) AcquireCount() int64

func (s *Stat) AcquireDuration() time.Duration

func (s *Stat) AcquiredConns() int32

func (s *Stat) CanceledAcquireCount() int64

func (s *Stat) ConstructingConns() int32

func (s *Stat) EmptyAcquireCount() int64

func (s *Stat) EmptyAcquireWaitTime() time.Duration

func (s *Stat) IdleConns() int32

func (s *Stat) MaxConns() int32

func (s *Stat) MaxIdleDestroyCount() int64

func (s *Stat) MaxLifetimeDestroyCount() int64

func (s *Stat) NewConnsCount() int64

func (s *Stat) TotalConns() int32

type TraceAcquireEndData

type TraceAcquireStartData

type TraceReleaseData

type Tx

func (tx *Tx) Begin(ctx context.Context) (pgx.Tx, error)

func (tx *Tx) Commit(ctx context.Context) error

func (tx *Tx) Conn() *pgx.Conn

func (tx *Tx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, ...) (int64, error)

func (tx *Tx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

func (tx *Tx) LargeObjects() pgx.LargeObjects

func (tx *Tx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)

func (tx *Tx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

func (tx *Tx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

func (tx *Tx) Rollback(ctx context.Context) error

func (tx *Tx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

This section is empty.

Types

¶

type

AcquireTracer

¶

added in

v5.6.0

type AcquireTracer interface {

// TraceAcquireStart is called at the beginning of Acquire.

// The returned context is used for the rest of the call and will be passed to the TraceAcquireEnd.

TraceAcquireStart(ctx

context

.

Context

, pool *

Pool

, data

TraceAcquireStartData

)

context

.

Context

// TraceAcquireEnd is called when a connection has been acquired.

TraceAcquireEnd(ctx

context

.

Context

, pool *

Pool

, data

TraceAcquireEndData

)
}

AcquireTracer traces Acquire.

type

Config

¶

type Config struct {

ConnConfig *

pgx

.

ConnConfig

// BeforeConnect is called before a new connection is made. It is passed a copy of the underlying pgx.ConnConfig and

// will not impact any existing open connections.

BeforeConnect func(

context

.

Context

, *

pgx

.

ConnConfig

)

error

// AfterConnect is called after a connection is established, but before it is added to the pool.

AfterConnect func(

context

.

Context

, *

pgx

.

Conn

)

error

// BeforeAcquire is called before a connection is acquired from the pool. It must return true to allow the

// acquisition or false to indicate that the connection should be destroyed and a different connection should be

// acquired.

//

// Deprecated: Use PrepareConn instead. If both PrepareConn and BeforeAcquire are set, PrepareConn will take

// precedence, ignoring BeforeAcquire.

BeforeAcquire func(

context

.

Context

, *

pgx

.

Conn

)

bool

// PrepareConn is called before a connection is acquired from the pool. If this function returns true, the connection

// is considered valid, otherwise the connection is destroyed. If the function returns a non-nil error, the instigating

// query will fail with the returned error.

//

// Specifically, this means that:

//

// 	- If it returns true and a nil error, the query proceeds as normal.

// 	- If it returns true and an error, the connection will be returned to the pool, and the instigating query will fail with the returned error.

// 	- If it returns false, and an error, the connection will be destroyed, and the query will fail with the returned error.

// 	- If it returns false and a nil error, the connection will be destroyed, and the instigating query will be retried on a new connection.

PrepareConn func(

context

.

Context

, *

pgx

.

Conn

) (

bool

,

error

)

// AfterRelease is called after a connection is released, but before it is returned to the pool. It must return true to

// return the connection to the pool or false to destroy the connection.

AfterRelease func(*

pgx

.

Conn

)

bool

// BeforeClose is called right before a connection is closed and removed from the pool.

BeforeClose func(*

pgx

.

Conn

)

// ShouldPing is called after a connection is acquired from the pool. If it returns true, the connection is pinged to check for liveness.

// If this func is not set, the default behavior is to ping connections that have been idle for at least 1 second.

ShouldPing func(

context

.

Context

,

ShouldPingParams

)

bool

// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.

MaxConnLifetime

time

.

Duration

// MaxConnLifetimeJitter is the duration after MaxConnLifetime to randomly decide to close a connection.

// This helps prevent all connections from being closed at the exact same time, starving the pool.

MaxConnLifetimeJitter

time

.

Duration

// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.

MaxConnIdleTime

time

.

Duration

// PingTimeout is the maximum amount of time to wait for a connection to pong before considering it as unhealthy and

// destroying it. If zero, the default is no timeout.

PingTimeout

time

.

Duration

// MaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU().

MaxConns

int32

// MinConns is the minimum size of the pool. After connection closes, the pool might dip below MinConns. A low

// number of MinConns might mean the pool is empty after MaxConnLifetime until the health check has a chance

// to create new connections.

MinConns

int32

// MinIdleConns is the minimum number of idle connections in the pool. You can increase this to ensure that

// there are always idle connections available. This can help reduce tail latencies during request processing,

// as you can avoid the latency of establishing a new connection while handling requests. It is superior

// to MinConns for this purpose.

// Similar to MinConns, the pool might temporarily dip below MinIdleConns after connection closes.

MinIdleConns

int32

// HealthCheckPeriod is the duration between checks of the health of idle connections.

HealthCheckPeriod

time

.

Duration

// contains filtered or unexported fields

}

Config is the configuration struct for creating a pool. It must be created by

ParseConfig

and then it can be
modified.

func

ParseConfig

¶

func ParseConfig(connString

string

) (*

Config

,

error

)

ParseConfig builds a Config from connString. It parses connString with the same behavior as

pgx.ParseConfig

with the
addition of the following variables:

pool_max_conns: integer greater than 0 (default 4)

pool_min_conns: integer 0 or greater (default 0)

pool_max_conn_lifetime: duration string (default 1 hour)

pool_max_conn_idle_time: duration string (default 30 minutes)

pool_health_check_period: duration string (default 1 minute)

pool_max_conn_lifetime_jitter: duration string (default 0)

See Config for definitions of these arguments.

# Example Keyword/Value
user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10 pool_max_conn_lifetime=1h30m

# Example URL
postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10&pool_max_conn_lifetime=1h30m

func (*Config)

ConnString

¶

func (c *

Config

) ConnString()

string

ConnString returns the connection string as parsed by pgxpool.ParseConfig into pgxpool.Config.

func (*Config)

Copy

¶

func (c *

Config

) Copy() *

Config

Copy returns a deep copy of the config that is safe to use and modify.
The only exception is the tls.Config:
according to the tls.Config docs it must not be modified after creation.

type

Conn

¶

type Conn struct {

// contains filtered or unexported fields

}

Conn is an acquired *pgx.Conn from a Pool.

func (*Conn)

Begin

¶

func (c *

Conn

) Begin(ctx

context

.

Context

) (

pgx

.

Tx

,

error

)

Begin starts a transaction block from the *Conn without explicitly setting a transaction mode (see BeginTx with TxOptions if transaction mode is required).

func (*Conn)

BeginTx

¶

func (c *

Conn

) BeginTx(ctx

context

.

Context

, txOptions

pgx

.

TxOptions

) (

pgx

.

Tx

,

error

)

BeginTx starts a transaction block from the *Conn with txOptions determining the transaction mode.

func (*Conn)

Conn

¶

func (c *

Conn

) Conn() *

pgx

.

Conn

func (*Conn)

CopyFrom

¶

func (c *

Conn

) CopyFrom(ctx

context

.

Context

, tableName

pgx

.

Identifier

, columnNames []

string

, rowSrc

pgx

.

CopyFromSource

) (

int64

,

error

)

func (*Conn)

Exec

¶

func (c *

Conn

) Exec(ctx

context

.

Context

, sql

string

, arguments ...

any

) (

pgconn

.

CommandTag

,

error

)

func (*Conn)

Hijack

¶

func (c *

Conn

) Hijack() *

pgx

.

Conn

Hijack assumes ownership of the connection from the pool. Caller is responsible for closing the connection. Hijack
will panic if called on an already released or hijacked connection.

func (*Conn)

Ping

¶

func (c *

Conn

) Ping(ctx

context

.

Context

)

error

func (*Conn)

Query

¶

func (c *

Conn

) Query(ctx

context

.

Context

, sql

string

, args ...

any

) (

pgx

.

Rows

,

error

)

func (*Conn)

QueryRow

¶

func (c *

Conn

) QueryRow(ctx

context

.

Context

, sql

string

, args ...

any

)

pgx

.

Row

func (*Conn)

Release

¶

func (c *

Conn

) Release()

Release returns c to the pool it was acquired from. Once Release has been called, other methods must not be called.
However, it is safe to call Release multiple times. Subsequent calls after the first will be ignored.

func (*Conn)

SendBatch

¶

func (c *

Conn

) SendBatch(ctx

context

.

Context

, b *

pgx

.

Batch

)

pgx

.

BatchResults

type

Pool

¶

type Pool struct {

// contains filtered or unexported fields

}

Pool allows for connection reuse.

func

New

¶

func New(ctx

context

.

Context

, connString

string

) (*

Pool

,

error

)

New creates a new Pool. See

ParseConfig

for information on connString format.

func

NewWithConfig

¶

func NewWithConfig(ctx

context

.

Context

, config *

Config

) (*

Pool

,

error

)

NewWithConfig creates a new Pool. config must have been created by

ParseConfig

.

func (*Pool)

Acquire

¶

func (p *

Pool

) Acquire(ctx

context

.

Context

) (c *

Conn

, err

error

)

Acquire returns a connection (*Conn) from the Pool

func (*Pool)

AcquireAllIdle

¶

func (p *

Pool

) AcquireAllIdle(ctx

context

.

Context

) []*

Conn

AcquireAllIdle atomically acquires all currently idle connections. Its intended use is for health check and
keep-alive functionality. It does not update pool statistics.

func (*Pool)

AcquireFunc

¶

func (p *

Pool

) AcquireFunc(ctx

context

.

Context

, f func(*

Conn

)

error

)

error

AcquireFunc acquires a *Conn and calls f with that *Conn. ctx will only affect the Acquire. It has no effect on the
call of f. The return value is either an error acquiring the *Conn or the return value of f. The *Conn is
automatically released after the call of f.

func (*Pool)

Begin

¶

func (p *

Pool

) Begin(ctx

context

.

Context

) (

pgx

.

Tx

,

error

)

Begin acquires a connection from the Pool and starts a transaction. Unlike database/sql, the context only affects the begin command. i.e. there is no
auto-rollback on context cancellation. Begin initiates a transaction block without explicitly setting a transaction mode for the block (see BeginTx with TxOptions if transaction mode is required).
*pgxpool.Tx is returned, which implements the pgx.Tx interface.
Commit or Rollback must be called on the returned transaction to finalize the transaction block.

func (*Pool)

BeginTx

¶

func (p *

Pool

) BeginTx(ctx

context

.

Context

, txOptions

pgx

.

TxOptions

) (

pgx

.

Tx

,

error

)

BeginTx acquires a connection from the Pool and starts a transaction with pgx.TxOptions determining the transaction mode.
Unlike database/sql, the context only affects the begin command. i.e. there is no auto-rollback on context cancellation.
*pgxpool.Tx is returned, which implements the pgx.Tx interface.
Commit or Rollback must be called on the returned transaction to finalize the transaction block.

func (*Pool)

Close

¶

func (p *

Pool

) Close()

Close closes all connections in the pool and rejects future Acquire calls. Blocks until all connections are returned
to pool and closed.

func (*Pool)

Config

¶

func (p *

Pool

) Config() *

Config

Config returns a copy of config that was used to initialize this pool.

func (*Pool)

CopyFrom

¶

func (p *

Pool

) CopyFrom(ctx

context

.

Context

, tableName

pgx

.

Identifier

, columnNames []

string

, rowSrc

pgx

.

CopyFromSource

) (

int64

,

error

)

func (*Pool)

Exec

¶

func (p *

Pool

) Exec(ctx

context

.

Context

, sql

string

, arguments ...

any

) (

pgconn

.

CommandTag

,

error

)

Exec acquires a connection from the Pool and executes the given SQL.
SQL can be either a prepared statement name or an SQL string.
Arguments should be referenced positionally from the SQL string as $1, $2, etc.
The acquired connection is returned to the pool when the Exec function returns.

func (*Pool)

Ping

¶

func (p *

Pool

) Ping(ctx

context

.

Context

)

error

Ping acquires a connection from the Pool and executes an empty sql statement against it.
If the sql returns without error, the database Ping is considered successful, otherwise, the error is returned.

func (*Pool)

Query

¶

func (p *

Pool

) Query(ctx

context

.

Context

, sql

string

, args ...

any

) (

pgx

.

Rows

,

error

)

Query acquires a connection and executes a query that returns pgx.Rows.
Arguments should be referenced positionally from the SQL string as $1, $2, etc.
See pgx.Rows documentation to close the returned Rows and return the acquired connection to the Pool.

If there is an error, the returned pgx.Rows will be returned in an error state.
If preferred, ignore the error returned from Query and handle errors using the returned pgx.Rows.

For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
needed. See the documentation for those types for details.

func (*Pool)

QueryRow

¶

func (p *

Pool

) QueryRow(ctx

context

.

Context

, sql

string

, args ...

any

)

pgx

.

Row

QueryRow acquires a connection and executes a query that is expected
to return at most one row (pgx.Row). Errors are deferred until pgx.Row's
Scan method is called. If the query selects no rows, pgx.Row's Scan will
return ErrNoRows. Otherwise, pgx.Row's Scan scans the first selected row
and discards the rest. The acquired connection is returned to the Pool when
pgx.Row's Scan method is called.

Arguments should be referenced positionally from the SQL string as $1, $2, etc.

For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
needed. See the documentation for those types for details.

func (*Pool)

Reset

¶

func (p *

Pool

) Reset()

Reset closes all connections, but leaves the pool open. It is intended for use when an error is detected that would
disrupt all connections (such as a network interruption or a server state change).

It is safe to reset a pool while connections are checked out. Those connections will be closed when they are returned
to the pool.

func (*Pool)

SendBatch

¶

func (p *

Pool

) SendBatch(ctx

context

.

Context

, b *

pgx

.

Batch

)

pgx

.

BatchResults

func (*Pool)

Stat

¶

func (p *

Pool

) Stat() *

Stat

Stat returns a pgxpool.Stat struct with a snapshot of Pool statistics.

type

ReleaseTracer

¶

added in

v5.6.0

type ReleaseTracer interface {

// TraceRelease is called at the beginning of Release.

TraceRelease(pool *

Pool

, data

TraceReleaseData

)
}

ReleaseTracer traces Release.

type

ShouldPingParams

¶

added in

v5.7.6

type ShouldPingParams struct {

Conn         *

pgx

.

Conn

IdleDuration

time

.

Duration

}

ShouldPingParams are the parameters passed to ShouldPing.

type

Stat

¶

type Stat struct {

// contains filtered or unexported fields

}

Stat is a snapshot of Pool statistics.

func (*Stat)

AcquireCount

¶

func (s *

Stat

) AcquireCount()

int64

AcquireCount returns the cumulative count of successful acquires from the pool.

func (*Stat)

AcquireDuration

¶

func (s *

Stat

) AcquireDuration()

time

.

Duration

AcquireDuration returns the total duration of all successful acquires from
the pool.

func (*Stat)

AcquiredConns

¶

func (s *

Stat

) AcquiredConns()

int32

AcquiredConns returns the number of currently acquired connections in the pool.

func (*Stat)

CanceledAcquireCount

¶

func (s *

Stat

) CanceledAcquireCount()

int64

CanceledAcquireCount returns the cumulative count of acquires from the pool
that were canceled by a context.

func (*Stat)

ConstructingConns

¶

func (s *

Stat

) ConstructingConns()

int32

ConstructingConns returns the number of conns with construction in progress in
the pool.

func (*Stat)

EmptyAcquireCount

¶

func (s *

Stat

) EmptyAcquireCount()

int64

EmptyAcquireCount returns the cumulative count of successful acquires from the pool
that waited for a resource to be released or constructed because the pool was
empty.

func (*Stat)

EmptyAcquireWaitTime

¶

added in

v5.7.3

func (s *

Stat

) EmptyAcquireWaitTime()

time

.

Duration

EmptyAcquireWaitTime returns the cumulative time waited for successful acquires
from the pool for a resource to be released or constructed because the pool was
empty.

func (*Stat)

IdleConns

¶

func (s *

Stat

) IdleConns()

int32

IdleConns returns the number of currently idle conns in the pool.

func (*Stat)

MaxConns

¶

func (s *

Stat

) MaxConns()

int32

MaxConns returns the maximum size of the pool.

func (*Stat)

MaxIdleDestroyCount

¶

func (s *

Stat

) MaxIdleDestroyCount()

int64

MaxIdleDestroyCount returns the cumulative count of connections destroyed because
they exceeded MaxConnIdleTime.

func (*Stat)

MaxLifetimeDestroyCount

¶

func (s *

Stat

) MaxLifetimeDestroyCount()

int64

MaxLifetimeDestroyCount returns the cumulative count of connections destroyed
because they exceeded MaxConnLifetime.

func (*Stat)

NewConnsCount

¶

func (s *

Stat

) NewConnsCount()

int64

NewConnsCount returns the cumulative count of new connections opened.

func (*Stat)

TotalConns

¶

func (s *

Stat

) TotalConns()

int32

TotalConns returns the total number of resources currently in the pool.
The value is the sum of ConstructingConns, AcquiredConns, and
IdleConns.

type

TraceAcquireEndData

¶

added in

v5.6.0

type TraceAcquireEndData struct {

Conn *

pgx

.

Conn

Err

error

}

type

TraceAcquireStartData

¶

added in

v5.6.0

type TraceAcquireStartData struct{}

type

TraceReleaseData

¶

added in

v5.6.0

type TraceReleaseData struct {

Conn *

pgx

.

Conn

}

type

Tx

¶

type Tx struct {

// contains filtered or unexported fields

}

Tx represents a database transaction acquired from a Pool.

func (*Tx)

Begin

¶

func (tx *

Tx

) Begin(ctx

context

.

Context

) (

pgx

.

Tx

,

error

)

Begin starts a pseudo nested transaction implemented with a savepoint.

func (*Tx)

Commit

¶

func (tx *

Tx

) Commit(ctx

context

.

Context

)

error

Commit commits the transaction and returns the associated connection back to the Pool. Commit will return an error
where errors.Is(ErrTxClosed) is true if the Tx is already closed, but is otherwise safe to call multiple times. If
the commit fails with a rollback status (e.g. the transaction was already in a broken state) then ErrTxCommitRollback
will be returned.

func (*Tx)

Conn

¶

func (tx *

Tx

) Conn() *

pgx

.

Conn

func (*Tx)

CopyFrom

¶

func (tx *

Tx

) CopyFrom(ctx

context

.

Context

, tableName

pgx

.

Identifier

, columnNames []

string

, rowSrc

pgx

.

CopyFromSource

) (

int64

,

error

)

func (*Tx)

Exec

¶

func (tx *

Tx

) Exec(ctx

context

.

Context

, sql

string

, arguments ...

any

) (

pgconn

.

CommandTag

,

error

)

func (*Tx)

LargeObjects

¶

func (tx *

Tx

) LargeObjects()

pgx

.

LargeObjects

func (*Tx)

Prepare

¶

func (tx *

Tx

) Prepare(ctx

context

.

Context

, name, sql

string

) (*

pgconn

.

StatementDescription

,

error

)

Prepare creates a prepared statement with name and sql. If the name is empty,
an anonymous prepared statement will be used. sql can contain placeholders
for bound parameters. These placeholders are referenced positionally as $1, $2, etc.

Prepare is idempotent; i.e. it is safe to call Prepare multiple times with the same
name and sql arguments. This allows a code path to Prepare and Query/Exec without
needing to first check whether the statement has already been prepared.

func (*Tx)

Query

¶

func (tx *

Tx

) Query(ctx

context

.

Context

, sql

string

, args ...

any

) (

pgx

.

Rows

,

error

)

func (*Tx)

QueryRow

¶

func (tx *

Tx

) QueryRow(ctx

context

.

Context

, sql

string

, args ...

any

)

pgx

.

Row

func (*Tx)

Rollback

¶

func (tx *

Tx

) Rollback(ctx

context

.

Context

)

error

Rollback rolls back the transaction and returns the associated connection back to the Pool. Rollback will return
where an error where errors.Is(ErrTxClosed) is true if the Tx is already closed, but is otherwise safe to call
multiple times. Hence, defer tx.Rollback() is safe even if tx.Commit() will be called first in a non-error condition.

func (*Tx)

SendBatch

¶

func (tx *

Tx

) SendBatch(ctx

context

.

Context

, b *

pgx

.

Batch

)

pgx

.

BatchResults