# River Documentation

> Auto-fetched from [https://riverqueue.com/docs](https://riverqueue.com/docs)
> Last Updated: 2026-01-29T20:11:31.145308+00:00

---

Getting started
Learn how to install River packages for Go, run migrations to get River's database schema in place, and create an initial worker and client to start inserting and working jobs.
Prerequisites
River requires an existing PostgreSQL database, and is most commonly used with
pgx
. River is tested using the three most recent major versions of PostgreSQL.
Installation
To install River, run the following in the directory of a Go project (where a
go.mod
file is present):
Terminal window
go
get
github.com/riverqueue/river
go
get
github.com/riverqueue/river/riverdriver/riverpgxv5
Alternatively, the
riverdatabasesql
driver can be used instead of
riverpgxv5
for compatibility with Go's built-in
database/sql
. See
inserting jobs with Bun
or
GORM
.
Running migrations
River persists jobs to a Postgres database, and needs a small set of tables created to insert jobs and carry out
leader election
. It's bundled with a command line tool which executes migrations, and which future-proofs River in case other migration steps need to be run in future versions.
From the same directory as above, install the River CLI:
Terminal window
go
install
github.com/riverqueue/river/cmd/river@latest
With the
DATABASE_URL
of a target database (looks like
postgres://host:5432/db
), migrate up:
Terminal window
river
migrate-up
--database-url
"
$DATABASE_URL
"
See also
migrations
.
Job args and workers
Each kind of job in River requires two types: a
JobArgs
struct and a
Worker[T JobArgs]
. The
JobArgs
struct has two purposes:
It defines the structured arguments for your worker. These arguments are serialized to JSON before the job is stored in the database.
It defines a
Kind() string
method that will be used to uniquely identify the kind of job in the database.
Here is a simple
Worker
and
JobArgs
setup for a
SortWorker
which will sort and print a list of strings provided in its arguments:
type
SortArgs
struct
{
// Strings is a slice of strings to sort.
Strings
[]
string
`json:"strings"`
}
func
(
SortArgs
)
Kind
()
string
{
return
"sort"
}
type
SortWorker
struct
{
// An embedded WorkerDefaults sets up default methods to fulfill the rest of
// the Worker interface:
river
.
WorkerDefaults
[
SortArgs
]
}
func
(
w
*
SortWorker
)
Work
(
ctx
context
.
Context
,
job
*
river
.
Job
[
SortArgs
])
error
{
sort
.
Strings
(
job
.
Args
.
Strings
)
fmt
.
Printf
(
"Sorted strings: %+v
\n
"
,
job
.
Args
.
Strings
)
return
nil
}
Generics
River utilizes Go generics to simplify your Worker definitions. This means that your worker only needs to deal with fully structured and typed set of arguments. As in the example above, a
Worker
has a 1:1 relationship with the
JobArgs
type it handles.
Registering workers
Jobs are uniquely identified by their "kind" string. Workers are registered on start up so that River knows how to assign jobs to workers:
workers
:=
river
.
NewWorkers
()
// AddWorker panics if the worker is already registered or invalid:
river
.
AddWorker
(
workers
,
&
SortWorker
{})
AddWorker
panics in case of invalid configuration. Given its succinct syntax and that bad configuration should prevent a worker process from booting, panicking is probably a reasonable compromise for most applications. However, for those who find it distastely,
AddWorkerSafely
is also provided:
workers
:=
river
.
NewWorkers
()
if
err
:=
river
.
AddWorkerSafely
(
workers
,
&
SortWorker
{});
err
!=
nil
{
panic
(
"handle this error"
)
}
Starting a client
A River
Client
provides an interface for job insertion and manages job processing and
maintenance services
. A client is created with a database pool,
driver
, and config struct containing a
Workers
bundle and other settings. Here's a client
Client
working one queue (
"default"
) with up to 100 worker goroutines at a time:
dbPool
,
err
:=
pgxpool
.
New
(
ctx
,
os
.
Getenv
(
"DATABASE_URL"
))
if
err
!=
nil
{
// handle error
}
riverClient
,
err
:=
river
.
NewClient
(
riverpgxv5
.
New
(
dbPool
),
&
river
.
Config
{
Queues
:
map
[
string
]
river
.
QueueConfig
{
river
.
QueueDefault
:
{
MaxWorkers
:
100
},
},
Workers
:
workers
,
})
if
err
!=
nil
{
// handle error
}
// Run the client inline. All executed jobs will inherit from ctx:
if
err
:=
riverClient
.
Start
(
ctx
);
err
!=
nil
{
// handle error
}
Stopping
The client should also be stopped on program shutdown:
// Stop fetching new work and wait for active jobs to finish.
if
err
:=
riverClient
.
Stop
(
ctx
);
err
!=
nil
{
// handle error
}
There are some complexities around ensuring clients stop cleanly, but also in a timely manner. Read
Graceful shutdown
for more details on River's stop modes.
Insert-Only clients
A common pattern is to have frontend processes which only insert jobs but do not work them, and a separate pool of workers which only work jobs. River supports this through the use of an insert-only
Client
.
An insert-only client is one that has not been started with
Start()
. For insert-only clients, the
Queues
and
Workers
fields from
Config
can be ommitted; however the
Workers
config allows the
Client
to validate that it is only inserting jobs whose worker is configured and may be worth keeping in place even on insert-only clients.
Inserting jobs
Client.InsertTx
is used in conjunction with an instance of job args to insert a job to work on a transaction:
_
,
err
=
riverClient
.
InsertTx
(
ctx
,
tx
,
SortArgs
{
Strings
:
[]
string
{
"whale"
,
"tiger"
,
"bear"
,
},
},
nil
)
if
err
!=
nil
{
// handle error
}
See the
InsertAndWork
example
for complete code.
Client.Insert
that doesn't take a transaction is also available, although as described in
Transactional enqueuing
, inserting jobs in transactions is usually more appropriate to avoid bugs.
_
,
err
=
riverClient
.
Insert
(
ctx
,
SortArgs
{
Strings
:
[]
string
{
"whale"
,
"tiger"
,
"bear"
,
},
},
nil
)
if
err
!=
nil
{
// handle error
}
See also
Batch job insertion
.
Next
Migrations