# River Job Queue

> Auto-fetched from [https://pkg.go.dev/github.com/riverqueue/river](https://pkg.go.dev/github.com/riverqueue/river)
> Last Updated: 2026-01-29T20:11:28.911075+00:00

---

Overview
¶
Job args and workers
Registering workers
Starting a client
Inserting jobs
Other features
Development
Package river is a robust high-performance job processing system for Go and
Postgres.
See
homepage
,
docs
, and
godoc
, as well as the
River UI
.
Being built for Postgres, River encourages the use of the same database for
application data and job queue. By enqueueing jobs transactionally along with
other database changes, whole classes of distributed systems problems are
avoided. Jobs are guaranteed to be enqueued if their transaction commits, are
removed if their transaction rolls back, and aren't visible for work _until_
commit. See
transactional enqueueing
for more background on this philosophy.
Job args and workers
¶
Jobs are defined in struct pairs, with an implementation of
`JobArgs`
and one
of
`Worker`
.
Job args contain `json` annotations and define how jobs are serialized to and
from the database, along with a "kind", a stable string that uniquely identifies
the job.
type SortArgs struct {
// Strings is a slice of strings to sort.
Strings []string `json:"strings"`
}

func (SortArgs) Kind() string { return "sort" }
Workers expose a `Work` function that dictates how jobs run.
type SortWorker struct {
// An embedded WorkerDefaults sets up default methods to fulfill the rest of
// the Worker interface:
river.WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[SortArgs]) error {
sort.Strings(job.Args.Strings)
fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
return nil
}
Registering workers
¶
Jobs are uniquely identified by their "kind" string. Workers are registered on
start up so that River knows how to assign jobs to workers:
workers := river.NewWorkers()
// AddWorker panics if the worker is already registered or invalid:
river.AddWorker(workers, &SortWorker{})
Starting a client
¶
A River
`Client`
provides an interface for job insertion and manages job
processing and
maintenance services
. A client's created with a database pool,
driver
, and config struct containing a `Workers` bundle and other settings.
Here's a client `Client` working one queue (`"default"`) with up to 100 worker
goroutines at a time:
riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Workers: workers,
})
if err != nil {
panic(err)
}

// Run the client inline. All executed jobs will inherit from ctx:
if err := riverClient.Start(ctx); err != nil {
panic(err)
}
## Insert-only clients
It's often desirable to have a client that'll be used for inserting jobs, but
not working them. This is possible by omitting the `Queues` configuration, and
skipping the call to `Start`:
riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Workers: workers,
})
if err != nil {
panic(err)
}
`Workers` can also be omitted, but it's better to include it so River can check
that inserted job kinds have a worker that can run them.
## Stopping
The client should also be stopped on program shutdown:
// Stop fetching new work and wait for active jobs to finish.
if err := riverClient.Stop(ctx); err != nil {
panic(err)
}
There are some complexities around ensuring clients stop cleanly, but also in a
timely manner. See
graceful shutdown
for more details on River's stop modes.
Inserting jobs
¶
`Client.InsertTx`
is used in conjunction with an instance of job args to
insert a job to work on a transaction:
_, err = riverClient.InsertTx(ctx, tx, SortArgs{
Strings: []string{
"whale", "tiger", "bear",
},
}, nil)

if err != nil {
panic(err)
}
See the
`InsertAndWork` example
for complete code.
Other features
¶
Batch job insertion
for efficiently inserting many jobs at once using
Postgres `COPY FROM`.
Cancelling jobs
from inside a work function.
Error and panic handling
.
Multiple queues
to better guarantee job throughput, worker availability,
and isolation between components.
Periodic and cron jobs
.
Scheduled jobs
that run automatically at their scheduled time in the
future.
Snoozing jobs
from inside a work function.
Subscriptions
to queue activity and statistics, providing easy hooks for
telemetry like logging and metrics.
Test helpers
to verify that jobs are inserted as expected.
Transactional job completion
to guarantee job completion commits with
other changes in a transaction.
Unique jobs
by args, period, queue, and state.
Web UI
for inspecting and interacting with jobs and queues.
Work functions
for simplified worker implementation.
## Cross language enqueueing
River supports inserting jobs in some non-Go languages which are then worked by Go implementations. This may be desirable in performance sensitive cases so that jobs can take advantage of Go's fast runtime.
Inserting jobs from Python
.
Inserting jobs from Ruby
.
Development
¶
See
developing River
.
Example (BatchInsert)
¶
Example_batchInsert demonstrates how many jobs can be inserted for work as
part of a single operation.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type BatchInsertArgs struct{}

func (BatchInsertArgs) Kind() string { return "batch_insert" }

// BatchInsertWorker is a job worker demonstrating use of custom
// job-specific insertion options.
type BatchInsertWorker struct {
river.WorkerDefaults[BatchInsertArgs]
}

func (w *BatchInsertWorker) Work(ctx context.Context, job *river.Job[BatchInsertArgs]) error {
fmt.Printf("Worked a job\n")
return nil
}

// Example_batchInsert demonstrates how many jobs can be inserted for work as
// part of a single operation.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &BatchInsertWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

results, err := riverClient.InsertMany(ctx, []river.InsertManyParams{
{Args: BatchInsertArgs{}},
{Args: BatchInsertArgs{}},
{Args: BatchInsertArgs{}},
{Args: BatchInsertArgs{}, InsertOpts: &river.InsertOpts{Priority: 3}},
{Args: BatchInsertArgs{}, InsertOpts: &river.InsertOpts{Priority: 4}},
})
if err != nil {
panic(err)
}
fmt.Printf("Inserted %d jobs\n", len(results))

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 5)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
Inserted 5 jobs
Worked a job
Worked a job
Worked a job
Worked a job
Worked a job
Share
Format
Run
Example (CompleteJobWithinTx)
¶
Example_completeJobWithinTx demonstrates how to transactionally complete
a job alongside other database changes being made.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type TransactionalArgs struct{}

func (TransactionalArgs) Kind() string { return "transactional_worker" }

// TransactionalWorker is a job worker which runs an operation on the database
// and transactionally completes the current job.
//
// While this example is simplified, any operations could be performed within
// the transaction such as inserting additional jobs or manipulating other data.
type TransactionalWorker struct {
river.WorkerDefaults[TransactionalArgs]

dbPool *pgxpool.Pool
}

func (w *TransactionalWorker) Work(ctx context.Context, job *river.Job[TransactionalArgs]) error {
tx, err := w.dbPool.Begin(ctx)
if err != nil {
return err
}
defer tx.Rollback(ctx)

var result int
if err := tx.QueryRow(ctx, "SELECT 1").Scan(&result); err != nil {
return err
}

// The function needs to know the type of the database driver in use by the
// Client, but the other generic parameters can be inferred.
jobAfter, err := river.JobCompleteTx[*riverpgxv5.Driver](ctx, tx, job)
if err != nil {
return err
}
fmt.Printf("Transitioned TransactionalWorker job from %q to %q\n", job.State, jobAfter.State)

if err = tx.Commit(ctx); err != nil {
return err
}
return nil
}

// Example_completeJobWithinTx demonstrates how to transactionally complete
// a job alongside other database changes being made.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &TransactionalWorker{dbPool: dbPool})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Not strictly needed, but used to help this test wait until job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

if _, err = riverClient.Insert(ctx, TransactionalArgs{}, nil); err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
Transitioned TransactionalWorker job from "running" to "completed"
Share
Format
Run
Example (CronJob)
¶
Example_cronJob demonstrates how to create a cron job with a more complex
schedule using a third party cron package to parse more elaborate crontab
syntax.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"
"github.com/robfig/cron/v3"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type CronJobArgs struct{}

// Kind is the unique string name for this job.
func (CronJobArgs) Kind() string { return "cron" }

// CronJobWorker is a job worker for sorting strings.
type CronJobWorker struct {
river.WorkerDefaults[CronJobArgs]
}

func (w *CronJobWorker) Work(ctx context.Context, job *river.Job[CronJobArgs]) error {
fmt.Printf("This job will run once immediately then every hour on the half hour\n")
return nil
}

// Example_cronJob demonstrates how to create a cron job with a more complex
// schedule using a third party cron package to parse more elaborate crontab
// syntax.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &CronJobWorker{})

schedule, err := cron.ParseStandard("30 * * * *") // every hour on the half hour
if err != nil {
panic(err)
}

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
PeriodicJobs: []*river.PeriodicJob{
river.NewPeriodicJob(
schedule,
func() (river.JobArgs, *river.InsertOpts) {
return CronJobArgs{}, nil
},
&river.PeriodicJobOpts{RunOnStart: true},
),
},
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

// There's no need to explicitly insert a periodic job. One will be inserted
// (and worked soon after) as the client starts up.
if err := riverClient.Start(ctx); err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
This job will run once immediately then every hour on the half hour
Share
Format
Run
Example (CustomInsertOpts)
¶
Example_customInsertOpts demonstrates the use of a job with custom
job-specific insertion options.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type AlwaysHighPriorityArgs struct{}

func (AlwaysHighPriorityArgs) Kind() string { return "always_high_priority" }

// InsertOpts returns custom insert options that every job of this type will
// inherit by default.
func (AlwaysHighPriorityArgs) InsertOpts() river.InsertOpts {
return river.InsertOpts{
Queue: "high_priority",
}
}

// AlwaysHighPriorityWorker is a job worker demonstrating use of custom
// job-specific insertion options.
type AlwaysHighPriorityWorker struct {
river.WorkerDefaults[AlwaysHighPriorityArgs]
}

func (w *AlwaysHighPriorityWorker) Work(ctx context.Context, job *river.Job[AlwaysHighPriorityArgs]) error {
fmt.Printf("Ran in queue: %s\n", job.Queue)
return nil
}

type SometimesHighPriorityArgs struct{}

func (SometimesHighPriorityArgs) Kind() string { return "sometimes_high_priority" }

// SometimesHighPriorityWorker is a job worker that's made high-priority
// sometimes through the use of options at insertion time.
type SometimesHighPriorityWorker struct {
river.WorkerDefaults[SometimesHighPriorityArgs]
}

func (w *SometimesHighPriorityWorker) Work(ctx context.Context, job *river.Job[SometimesHighPriorityArgs]) error {
fmt.Printf("Ran in queue: %s\n", job.Queue)
return nil
}

// Example_customInsertOpts demonstrates the use of a job with custom
// job-specific insertion options.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &AlwaysHighPriorityWorker{})
river.AddWorker(workers, &SometimesHighPriorityWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
"high_priority":    {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

// This job always runs in the high-priority queue because its job-specific
// options on the struct above dictate that it will.
_, err = riverClient.Insert(ctx, AlwaysHighPriorityArgs{}, nil)
if err != nil {
panic(err)
}

// This job will run in the high-priority queue because of the options given
// at insertion time.
_, err = riverClient.Insert(ctx, SometimesHighPriorityArgs{}, &river.InsertOpts{
Queue: "high_priority",
})
if err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 2)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
Ran in queue: high_priority
Ran in queue: high_priority
Share
Format
Run
Example (ErrorHandler)
¶
Example_errorHandler demonstrates how to use the ErrorHandler interface for
custom application telemetry.
package main

import (
"context"
"errors"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
"github.com/riverqueue/river/rivertype"
)

type CustomErrorHandler struct{}

func (*CustomErrorHandler) HandleError(ctx context.Context, job *rivertype.JobRow, err error) *river.ErrorHandlerResult {
fmt.Printf("Job errored with: %s\n", err)
return nil
}

func (*CustomErrorHandler) HandlePanic(ctx context.Context, job *rivertype.JobRow, panicVal any, trace string) *river.ErrorHandlerResult {
fmt.Printf("Job panicked with: %v\n", panicVal)

// Either function can also set the job to be immediately cancelled.
return &river.ErrorHandlerResult{SetCancelled: true}
}

type ErroringArgs struct {
ShouldError bool
ShouldPanic bool
}

func (ErroringArgs) Kind() string { return "erroring" }

// Here to make sure our jobs are never accidentally retried which would add
// additional output and fail the example.
func (ErroringArgs) InsertOpts() river.InsertOpts {
return river.InsertOpts{MaxAttempts: 1}
}

type ErroringWorker struct {
river.WorkerDefaults[ErroringArgs]
}

func (w *ErroringWorker) Work(ctx context.Context, job *river.Job[ErroringArgs]) error {
switch {
case job.Args.ShouldError:
return errors.New("this job errored")
case job.Args.ShouldPanic:
panic("this job panicked")
}
return nil
}

// Example_errorHandler demonstrates how to use the ErrorHandler interface for
// custom application telemetry.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &ErroringWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
ErrorHandler: &CustomErrorHandler{},
Logger:       slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(9), ReplaceAttr: slogutil.NoLevelTime})), // Suppress logging so example output is cleaner (9 > slog.LevelError).
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 10},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Not strictly needed, but used to help this test wait until job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCancelled, river.EventKindJobFailed)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

if _, err = riverClient.Insert(ctx, ErroringArgs{ShouldError: true}, nil); err != nil {
panic(err)
}

// Wait for the first job before inserting another to guarantee test output
// is ordered correctly.
// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if _, err = riverClient.Insert(ctx, ErroringArgs{ShouldPanic: true}, nil); err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
Job errored with: this job errored
Job panicked with: this job panicked
Share
Format
Run
Example (GlobalHooks)
¶
Example_globalHooks demonstrates the use of hooks to modify River behavior
which are global to a River client.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
"github.com/riverqueue/river/rivertype"
)

type BothInsertAndWorkBeginHook struct{ river.HookDefaults }

func (BothInsertAndWorkBeginHook) InsertBegin(ctx context.Context, params *rivertype.JobInsertParams) error {
fmt.Printf("BothInsertAndWorkBeginHook.InsertBegin ran\n")
return nil
}

func (BothInsertAndWorkBeginHook) WorkBegin(ctx context.Context, job *rivertype.JobRow) error {
fmt.Printf("BothInsertAndWorkBeginHook.WorkBegin ran\n")
return nil
}

type InsertBeginHook struct{ river.HookDefaults }

func (InsertBeginHook) InsertBegin(ctx context.Context, params *rivertype.JobInsertParams) error {
fmt.Printf("InsertBeginHook.InsertBegin ran\n")
return nil
}

type WorkBeginHook struct{ river.HookDefaults }

func (WorkBeginHook) WorkBegin(ctx context.Context, job *rivertype.JobRow) error {
fmt.Printf("WorkBeginHook.WorkBegin ran\n")
return nil
}

// Verify interface compliance. It's recommended that these are included in your
// test suite to make sure that your hooks are complying to the specific
// interface hooks that you expected them to be.
var (
_ rivertype.HookInsertBegin = &BothInsertAndWorkBeginHook{}
_ rivertype.HookWorkBegin   = &BothInsertAndWorkBeginHook{}
_ rivertype.HookInsertBegin = &InsertBeginHook{}
_ rivertype.HookWorkBegin   = &WorkBeginHook{}
)

// Example_globalHooks demonstrates the use of hooks to modify River behavior
// which are global to a River client.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &NoOpWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
// Order is significant. See output below.
Hooks: []rivertype.Hook{
&BothInsertAndWorkBeginHook{},
&InsertBeginHook{},
&WorkBeginHook{},
},
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

_, err = riverClient.Insert(ctx, NoOpArgs{}, nil)
if err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
BothInsertAndWorkBeginHook.InsertBegin ran
InsertBeginHook.InsertBegin ran
BothInsertAndWorkBeginHook.WorkBegin ran
WorkBeginHook.WorkBegin ran
NoOpWorker.Work ran
Share
Format
Run
Example (GlobalMiddleware)
¶
Example_globalMiddleware demonstrates the use of middleware to modify River
behavior which are global to a River client.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
"github.com/riverqueue/river/rivertype"
)

type JobBothInsertAndWorkMiddleware struct{ river.MiddlewareDefaults }

func (JobBothInsertAndWorkMiddleware) InsertMany(ctx context.Context, manyParams []*rivertype.JobInsertParams, doInner func(ctx context.Context) ([]*rivertype.JobInsertResult, error)) ([]*rivertype.JobInsertResult, error) {
fmt.Printf("JobBothInsertAndWorkMiddleware.InsertMany ran\n")
return doInner(ctx)
}

func (JobBothInsertAndWorkMiddleware) Work(ctx context.Context, job *rivertype.JobRow, doInner func(ctx context.Context) error) error {
fmt.Printf("JobBothInsertAndWorkMiddleware.Work ran\n")
return doInner(ctx)
}

type JobInsertMiddleware struct{ river.MiddlewareDefaults }

func (JobInsertMiddleware) InsertMany(ctx context.Context, manyParams []*rivertype.JobInsertParams, doInner func(ctx context.Context) ([]*rivertype.JobInsertResult, error)) ([]*rivertype.JobInsertResult, error) {
fmt.Printf("JobInsertMiddleware.InsertMany ran\n")
return doInner(ctx)
}

type WorkerMiddleware struct{ river.MiddlewareDefaults }

func (WorkerMiddleware) Work(ctx context.Context, job *rivertype.JobRow, doInner func(ctx context.Context) error) error {
fmt.Printf("WorkerMiddleware.Work ran\n")
return doInner(ctx)
}

// Verify interface compliance. It's recommended that these are included in your
// test suite to make sure that your middlewares are complying to the specific
// interface middlewares that you expected them to be.
var (
_ rivertype.JobInsertMiddleware = &JobBothInsertAndWorkMiddleware{}
_ rivertype.WorkerMiddleware    = &JobBothInsertAndWorkMiddleware{}
_ rivertype.JobInsertMiddleware = &JobInsertMiddleware{}
_ rivertype.WorkerMiddleware    = &WorkerMiddleware{}
)

// Example_globalMiddleware demonstrates the use of middleware to modify River
// behavior which are global to a River client.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &NoOpWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
// Order is significant. See output below.
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Middleware: []rivertype.Middleware{
&JobBothInsertAndWorkMiddleware{},
&JobInsertMiddleware{},
&WorkerMiddleware{},
},
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

_, err = riverClient.Insert(ctx, NoOpArgs{}, nil)
if err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
JobBothInsertAndWorkMiddleware.InsertMany ran
JobInsertMiddleware.InsertMany ran
JobBothInsertAndWorkMiddleware.Work ran
WorkerMiddleware.Work ran
NoOpWorker.Work ran
Share
Format
Run
Example (GracefulShutdown)
¶
Example_gracefulShutdown demonstrates a realistic-looking stop loop for
River. It listens for SIGINT/SIGTERM (like might be received by a Ctrl+C
locally or on a platform like Heroku to stop a process) and when received,
tries a soft stop that waits for work to finish. If it doesn't finish in
time, a second SIGINT/SIGTERM will initiate a hard stop that cancels all jobs
using context cancellation. A third will give up on the stop procedure and
exit uncleanly.
package main

import (
"context"
"errors"
"fmt"
"log/slog"
"os"
"os/signal"
"syscall"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type WaitsForCancelOnlyArgs struct{}

func (WaitsForCancelOnlyArgs) Kind() string { return "waits_for_cancel_only" }

// WaitsForCancelOnlyWorker is a worker that will never finish jobs until its
// context is cancelled.
type WaitsForCancelOnlyWorker struct {
river.WorkerDefaults[WaitsForCancelOnlyArgs]

jobStarted chan struct{}
}

func (w *WaitsForCancelOnlyWorker) Work(ctx context.Context, job *river.Job[WaitsForCancelOnlyArgs]) error {
fmt.Printf("Working job that doesn't finish until cancelled\n")
close(w.jobStarted)

<-ctx.Done()
fmt.Printf("Job cancelled\n")

// In the event of cancellation, an error should be returned so that the job
// goes back in the retry queue.
return ctx.Err()
}

// Example_gracefulShutdown demonstrates a realistic-looking stop loop for
// River. It listens for SIGINT/SIGTERM (like might be received by a Ctrl+C
// locally or on a platform like Heroku to stop a process) and when received,
// tries a soft stop that waits for work to finish. If it doesn't finish in
// time, a second SIGINT/SIGTERM will initiate a hard stop that cancels all jobs
// using context cancellation. A third will give up on the stop procedure and
// exit uncleanly.
func main() {
ctx := context.Background()

jobStarted := make(chan struct{})

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &WaitsForCancelOnlyWorker{jobStarted: jobStarted})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTimeJobID})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

_, err = riverClient.Insert(ctx, WaitsForCancelOnlyArgs{}, nil)
if err != nil {
panic(err)
}

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

sigintOrTerm := make(chan os.Signal, 1)
signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

// This is meant to be a realistic-looking stop goroutine that might go in a
// real program. It waits for SIGINT/SIGTERM and when received, tries to stop
// gracefully by allowing a chance for jobs to finish. But if that isn't
// working, a second SIGINT/SIGTERM will tell it to terminate with prejudice and
// it'll issue a hard stop that cancels the context of all active jobs. In
// case that doesn't work, a third SIGINT/SIGTERM ignores River's stop procedure
// completely and exits uncleanly.
go func() {
<-sigintOrTerm
fmt.Printf("Received SIGINT/SIGTERM; initiating soft stop (try to wait for jobs to finish)\n")

softStopCtx, softStopCtxCancel := context.WithTimeout(ctx, 10*time.Second)
defer softStopCtxCancel()

go func() {
select {
case <-sigintOrTerm:
fmt.Printf("Received SIGINT/SIGTERM again; initiating hard stop (cancel everything)\n")
softStopCtxCancel()
case <-softStopCtx.Done():
fmt.Printf("Soft stop timeout; initiating hard stop (cancel everything)\n")
}
}()

err := riverClient.Stop(softStopCtx)
if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
panic(err)
}
if err == nil {
fmt.Printf("Soft stop succeeded\n")
return
}

hardStopCtx, hardStopCtxCancel := context.WithTimeout(ctx, 10*time.Second)
defer hardStopCtxCancel()

// As long as all jobs respect context cancellation, StopAndCancel will
// always work. However, in the case of a bug where a job blocks despite
// being cancelled, it may be necessary to either ignore River's stop
// result (what's shown here) or have a supervisor kill the process.
err = riverClient.StopAndCancel(hardStopCtx)
if err != nil && errors.Is(err, context.DeadlineExceeded) {
fmt.Printf("Hard stop timeout; ignoring stop procedure and exiting unsafely\n")
} else if err != nil {
panic(err)
}

// hard stop succeeded
}()

// Make sure our job starts being worked before doing anything else.
<-jobStarted

// Cheat a little by sending a SIGTERM manually for the purpose of this
// example (normally this will be sent by user or supervisory process). The
// first SIGTERM tries a soft stop in which jobs are given a chance to
// finish up.
sigintOrTerm <- syscall.SIGTERM

// The soft stop will never work in this example because our job only
// respects context cancellation, but wait a short amount of time to give it
// a chance. After it elapses, send another SIGTERM to initiate a hard stop.
select {
case <-riverClient.Stopped():
// Will never be reached in this example because our job will only ever
// finish on context cancellation.
fmt.Printf("Soft stop succeeded\n")

case <-time.After(100 * time.Millisecond):
sigintOrTerm <- syscall.SIGTERM
<-riverClient.Stopped()
}

}
Output:
Working job that doesn't finish until cancelled
Received SIGINT/SIGTERM; initiating soft stop (try to wait for jobs to finish)
Received SIGINT/SIGTERM again; initiating hard stop (cancel everything)
Job cancelled
msg="jobexecutor.JobExecutor: Job errored; retrying" error="context canceled" job_kind=waits_for_cancel_only
Share
Format
Run
Example (InsertAndWork)
¶
Example_insertAndWork demonstrates how to register job workers, start a
client, and insert a job on it to be worked.
package main

import (
"context"
"fmt"
"log/slog"
"os"
"sort"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type SortArgs struct {
// Strings is a slice of strings to sort.
Strings []string `json:"strings"`
}

func (SortArgs) Kind() string { return "sort" }

type SortWorker struct {
river.WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[SortArgs]) error {
sort.Strings(job.Args.Strings)
fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
return nil
}

// Example_insertAndWork demonstrates how to register job workers, start a
// client, and insert a job on it to be worked.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &SortWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

// Start a transaction to insert a job. It's also possible to insert a job
// outside a transaction, but this usage is recommended to ensure that all
// data a job needs to run is available by the time it starts. Because of
// snapshot visibility guarantees across transactions, the job will not be
// worked until the transaction has committed.
tx, err := dbPool.Begin(ctx)
if err != nil {
panic(err)
}
defer tx.Rollback(ctx)

_, err = riverClient.InsertTx(ctx, tx, SortArgs{
Strings: []string{
"whale", "tiger", "bear",
},
}, nil)
if err != nil {
panic(err)
}

if err := tx.Commit(ctx); err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
Sorted strings: [bear tiger whale]
Share
Format
Run
Example (JobArgsHooks)
¶
Example_jobArgsHooks demonstrates the use of hooks to modify River behavior.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
"github.com/riverqueue/river/rivertype"
)

type JobWithHooksArgs struct{}

func (JobWithHooksArgs) Kind() string { return "job_with_hooks" }

// Warning: Hooks is only called once per job insert or work and its return
// value is memoized. It should not vary based on the contents of any particular
// args because changes will be ignored.
func (JobWithHooksArgs) Hooks() []rivertype.Hook {
// Order is significant. See output below.
return []rivertype.Hook{
&JobWithHooksBothInsertAndWorkBeginHook{},
&JobWithHooksInsertBeginHook{},
&JobWithHooksWorkBeginHook{},
}
}

type JobWithHooksWorker struct {
river.WorkerDefaults[JobWithHooksArgs]
}

func (w *JobWithHooksWorker) Work(ctx context.Context, job *river.Job[JobWithHooksArgs]) error {
fmt.Printf("JobWithHooksWorker.Work ran\n")
return nil
}

type JobWithHooksBothInsertAndWorkBeginHook struct{ river.HookDefaults }

func (JobWithHooksBothInsertAndWorkBeginHook) InsertBegin(ctx context.Context, params *rivertype.JobInsertParams) error {
fmt.Printf("JobWithHooksInsertAndWorkBeginHook.InsertBegin ran\n")
return nil
}

func (JobWithHooksBothInsertAndWorkBeginHook) WorkBegin(ctx context.Context, job *rivertype.JobRow) error {
fmt.Printf("JobWithHooksInsertAndWorkBeginHook.WorkBegin ran\n")
return nil
}

type JobWithHooksInsertBeginHook struct{ river.HookDefaults }

func (JobWithHooksInsertBeginHook) InsertBegin(ctx context.Context, params *rivertype.JobInsertParams) error {
fmt.Printf("JobWithHooksInsertBeginHook.InsertBegin ran\n")
return nil
}

type JobWithHooksWorkBeginHook struct{ river.HookDefaults }

func (JobWithHooksWorkBeginHook) WorkBegin(ctx context.Context, job *rivertype.JobRow) error {
fmt.Printf("JobWithHooksWorkBeginHook.WorkBegin ran\n")
return nil
}

// Verify interface compliance. It's recommended that these are included in your
// test suite to make sure that your hooks are complying to the specific
// interface hooks that you expected them to be.
var (
_ rivertype.HookInsertBegin = &BothInsertAndWorkBeginHook{}
_ rivertype.HookWorkBegin   = &BothInsertAndWorkBeginHook{}
_ rivertype.HookInsertBegin = &InsertBeginHook{}
_ rivertype.HookWorkBegin   = &WorkBeginHook{}
)

// Example_jobArgsHooks demonstrates the use of hooks to modify River behavior.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &JobWithHooksWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

_, err = riverClient.Insert(ctx, JobWithHooksArgs{}, nil)
if err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
JobWithHooksInsertAndWorkBeginHook.InsertBegin ran
JobWithHooksInsertBeginHook.InsertBegin ran
JobWithHooksInsertAndWorkBeginHook.WorkBegin ran
JobWithHooksWorkBeginHook.WorkBegin ran
JobWithHooksWorker.Work ran
Share
Format
Run
Example (JobCancel)
¶
Example_jobCancel demonstrates how to permanently cancel a job from within
Work using JobCancel.
package main

import (
"context"
"errors"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type CancellingArgs struct {
ShouldCancel bool
}

func (args CancellingArgs) Kind() string { return "Cancelling" }

type CancellingWorker struct {
river.WorkerDefaults[CancellingArgs]
}

func (w *CancellingWorker) Work(ctx context.Context, job *river.Job[CancellingArgs]) error {
if job.Args.ShouldCancel {
fmt.Println("cancelling job")
return river.JobCancel(errors.New("this wrapped error message will be persisted to DB"))
}
return nil
}

// Example_jobCancel demonstrates how to permanently cancel a job from within
// Work using JobCancel.
func main() { //nolint:dupl
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &CancellingWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 10},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Not strictly needed, but used to help this test wait until job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCancelled)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}
if _, err = riverClient.Insert(ctx, CancellingArgs{ShouldCancel: true}, nil); err != nil {
panic(err)
}
// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
cancelling job
Share
Format
Run
Example (JobCancelFromClient)
¶
Example_jobCancelFromClient demonstrates how to permanently cancel a job from
any Client using JobCancel.
package main

import (
"context"
"errors"
"log/slog"
"os"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type SleepingArgs struct{}

func (args SleepingArgs) Kind() string { return "SleepingWorker" }

type SleepingWorker struct {
river.WorkerDefaults[CancellingArgs]

jobChan chan int64
}

func (w *SleepingWorker) Work(ctx context.Context, job *river.Job[CancellingArgs]) error {
w.jobChan <- job.ID
select {
case <-ctx.Done():
case <-time.After(5 * time.Second):
return errors.New("sleeping worker timed out")
}
return ctx.Err()
}

// Example_jobCancelFromClient demonstrates how to permanently cancel a job from
// any Client using JobCancel.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

jobChan := make(chan int64)

workers := river.NewWorkers()
river.AddWorker(workers, &SleepingWorker{jobChan: jobChan})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTimeJobID})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 10},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Not strictly needed, but used to help this test wait until job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCancelled)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}
insertRes, err := riverClient.Insert(ctx, CancellingArgs{ShouldCancel: true}, nil)
if err != nil {
panic(err)
}
select {
case <-jobChan:
case <-time.After(2 * time.Second):
panic("no jobChan signal received")
}

// There is presently no way to wait for the client to be 100% ready, so we
// sleep for a bit to give it time to start up. This is only needed in this
// example because we need the notifier to be ready for it to receive the
// cancellation signal.
time.Sleep(500 * time.Millisecond)

if _, err = riverClient.JobCancel(ctx, insertRes.Job.ID); err != nil {
panic(err)
}
// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
msg="jobexecutor.JobExecutor: job cancelled remotely"
Share
Format
Run
Example (JobSnooze)
¶
Example_jobSnooze demonstrates how to snooze a job from within Work using
JobSnooze. The job will be run again after 5 minutes and the snooze attempt
will decrement the job's attempt count, ensuring that one can snooze as many
times as desired without being impacted by the max attempts.
package main

import (
"context"
"fmt"
"log/slog"
"os"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type SnoozingArgs struct {
ShouldSnooze bool
}

func (args SnoozingArgs) Kind() string { return "Snoozing" }

type SnoozingWorker struct {
river.WorkerDefaults[SnoozingArgs]
}

func (w *SnoozingWorker) Work(ctx context.Context, job *river.Job[SnoozingArgs]) error {
if job.Args.ShouldSnooze {
fmt.Println("snoozing job for 5 minutes")
return river.JobSnooze(5 * time.Minute)
}
return nil
}

// Example_jobSnooze demonstrates how to snooze a job from within Work using
// JobSnooze. The job will be run again after 5 minutes and the snooze attempt
// will decrement the job's attempt count, ensuring that one can snooze as many
// times as desired without being impacted by the max attempts.
func main() { //nolint:dupl
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &SnoozingWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 10},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// The subscription bits are not needed in real usage, but are used to make
// sure the test waits until the job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobSnoozed)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}
if _, err = riverClient.Insert(ctx, SnoozingArgs{ShouldSnooze: true}, nil); err != nil {
panic(err)
}
// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
snoozing job for 5 minutes
Share
Format
Run
Example (PeriodicJob)
¶
Example_periodicJob demonstrates the use of a periodic job.
package main

import (
"context"
"fmt"
"log/slog"
"os"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type PeriodicJobArgs struct{}

// Kind is the unique string name for this job.
func (PeriodicJobArgs) Kind() string { return "periodic" }

// PeriodicJobWorker is a job worker for sorting strings.
type PeriodicJobWorker struct {
river.WorkerDefaults[PeriodicJobArgs]
}

func (w *PeriodicJobWorker) Work(ctx context.Context, job *river.Job[PeriodicJobArgs]) error {
fmt.Printf("This job will run once immediately then approximately once every 15 minutes\n")
return nil
}

// Example_periodicJob demonstrates the use of a periodic job.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &PeriodicJobWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
PeriodicJobs: []*river.PeriodicJob{
river.NewPeriodicJob(
river.PeriodicInterval(15*time.Minute),
func() (river.JobArgs, *river.InsertOpts) {
return PeriodicJobArgs{}, nil
},
&river.PeriodicJobOpts{RunOnStart: true},
),
},
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

// There's no need to explicitly insert a periodic job. One will be inserted
// (and worked soon after) as the client starts up.
if err := riverClient.Start(ctx); err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

// Periodic jobs can also be configured dynamically after a client has
// already started. Added jobs are scheduled for run immediately.
riverClient.PeriodicJobs().Clear()
riverClient.PeriodicJobs().Add(
river.NewPeriodicJob(
river.PeriodicInterval(15*time.Minute),
func() (river.JobArgs, *river.InsertOpts) {
return PeriodicJobArgs{}, nil
},
nil,
),
)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
This job will run once immediately then approximately once every 15 minutes
Share
Format
Run
Example (QueuePause)
¶
Example_queuePause demonstrates how to pause queues to prevent them from
working new jobs, and later resume them.
package main

import (
"context"
"fmt"
"log/slog"
"os"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type ReportingArgs struct{}

func (args ReportingArgs) Kind() string { return "Reporting" }

type ReportingWorker struct {
river.WorkerDefaults[ReportingArgs]

jobWorkedCh chan<- string
}

func (w *ReportingWorker) Work(ctx context.Context, job *river.Job[ReportingArgs]) error {
select {
case <-ctx.Done():
return ctx.Err()
case w.jobWorkedCh <- job.Queue:
return nil
}
}

// Example_queuePause demonstrates how to pause queues to prevent them from
// working new jobs, and later resume them.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

const (
unreliableQueue = "unreliable_external_service"
reliableQueue   = "reliable_jobs"
)

workers := river.NewWorkers()
jobWorkedCh := make(chan string)
river.AddWorker(workers, &ReportingWorker{jobWorkedCh: jobWorkedCh})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
unreliableQueue: {MaxWorkers: 10},
reliableQueue:   {MaxWorkers: 10},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

// Out of example scope, but used to wait until a queue is paused or unpaused.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindQueuePaused, river.EventKindQueueResumed)
defer subscribeCancel()

fmt.Printf("Pausing %s queue\n", unreliableQueue)
if err := riverClient.QueuePause(ctx, unreliableQueue, nil); err != nil {
panic(err)
}

// Wait for queue to be paused:
waitOrTimeout(subscribeChan)

fmt.Println("Inserting one job each into unreliable and reliable queues")
if _, err = riverClient.Insert(ctx, ReportingArgs{}, &river.InsertOpts{Queue: unreliableQueue}); err != nil {
panic(err)
}
if _, err = riverClient.Insert(ctx, ReportingArgs{}, &river.InsertOpts{Queue: reliableQueue}); err != nil {
panic(err)
}
// The unreliable queue is paused so its job should get worked yet, while
// reliable queue is not paused so its job should get worked immediately:
receivedQueue := waitOrTimeout(jobWorkedCh)
fmt.Printf("Job worked on %s queue\n", receivedQueue)

// Resume the unreliable queue so it can work the job:
fmt.Printf("Resuming %s queue\n", unreliableQueue)
if err := riverClient.QueueResume(ctx, unreliableQueue, nil); err != nil {
panic(err)
}
receivedQueue = waitOrTimeout(jobWorkedCh)
fmt.Printf("Job worked on %s queue\n", receivedQueue)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}

func waitOrTimeout[T any](ch <-chan T) T {
select {
case item := <-ch:
return item
case <-time.After(5 * time.Second):
panic("WaitOrTimeout timed out after waiting 5s")
}
}
Output:
Pausing unreliable_external_service queue
Inserting one job each into unreliable and reliable queues
Job worked on reliable_jobs queue
Resuming unreliable_external_service queue
Job worked on unreliable_external_service queue
Share
Format
Run
Example (ScheduledJob)
¶
Example_scheduledJob demonstrates how to schedule a job to be worked in the
future.
package main

import (
"context"
"fmt"
"log/slog"
"os"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type ScheduledArgs struct {
Message string `json:"message"`
}

func (ScheduledArgs) Kind() string { return "scheduled" }

type ScheduledWorker struct {
river.WorkerDefaults[ScheduledArgs]
}

func (w *ScheduledWorker) Work(ctx context.Context, job *river.Job[ScheduledArgs]) error {
fmt.Printf("Message: %s\n", job.Args.Message)
return nil
}

// Example_scheduledJob demonstrates how to schedule a job to be worked in the
// future.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &ScheduledWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

_, err = riverClient.Insert(ctx,
ScheduledArgs{
Message: "hello from the future",
},
&river.InsertOpts{
// Schedule the job to be worked in three hours.
ScheduledAt: time.Now().Add(3 * time.Hour),
})
if err != nil {
panic(err)
}

// Unlike most other examples, we don't wait for the job to be worked since
// doing so would require making the job's scheduled time contrived, and the
// example therefore less realistic/useful.

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Share
Format
Run
Example (Subscription)
¶
Example_subscription demonstrates the use of client subscriptions to receive
events containing information about worked jobs.
package main

import (
"context"
"errors"
"fmt"
"log/slog"
"os"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type SubscriptionArgs struct {
Cancel bool `json:"cancel"`
Fail   bool `json:"fail"`
}

func (SubscriptionArgs) Kind() string { return "subscription" }

type SubscriptionWorker struct {
river.WorkerDefaults[SubscriptionArgs]
}

func (w *SubscriptionWorker) Work(ctx context.Context, job *river.Job[SubscriptionArgs]) error {
switch {
case job.Args.Cancel:
return river.JobCancel(errors.New("cancelling job"))
case job.Args.Fail:
return errors.New("failing job")
}
return nil
}

// Example_subscription demonstrates the use of client subscriptions to receive
// events containing information about worked jobs.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &SubscriptionWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(9), ReplaceAttr: slogutil.NoLevelTime})), // Suppress logging so example output is cleaner (9 > slog.LevelError).
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Subscribers tell the River client the kinds of events they'd like to receive.
completedChan, completedSubscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer completedSubscribeCancel()

// Multiple simultaneous subscriptions are allowed.
failedChan, failedSubscribeCancel := riverClient.Subscribe(river.EventKindJobFailed)
defer failedSubscribeCancel()

otherChan, otherSubscribeCancel := riverClient.Subscribe(river.EventKindJobCancelled, river.EventKindJobSnoozed)
defer otherSubscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

// Insert one job for each subscription above: one to succeed, one to fail,
// and one that's cancelled that'll arrive on the "other" channel.
_, err = riverClient.Insert(ctx, SubscriptionArgs{}, nil)
if err != nil {
panic(err)
}
_, err = riverClient.Insert(ctx, SubscriptionArgs{Fail: true}, nil)
if err != nil {
panic(err)
}
_, err = riverClient.Insert(ctx, SubscriptionArgs{Cancel: true}, nil)
if err != nil {
panic(err)
}

waitForJob := func(subscribeChan <-chan *river.Event) {
select {
case event := <-subscribeChan:
if event == nil {
fmt.Printf("Channel is closed\n")
return
}

fmt.Printf("Got job with state: %s\n", event.Job.State)
case <-time.After(riversharedtest.WaitTimeout()):
panic("timed out waiting for job")
}
}

waitForJob(completedChan)
waitForJob(failedChan)
waitForJob(otherChan)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

fmt.Printf("Client stopped\n")

// Try waiting again, but none of these work because stopping the client
// closed all subscription channels automatically.
waitForJob(completedChan)
waitForJob(failedChan)
waitForJob(otherChan)

}
Output:
Got job with state: completed
Got job with state: available
Got job with state: cancelled
Client stopped
Channel is closed
Channel is closed
Channel is closed
Share
Format
Run
Example (UniqueJob)
¶
Example_uniqueJob demonstrates the use of a job with custom
job-specific insertion options.
package main

import (
"context"
"fmt"
"log/slog"
"os"
"time"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

// Account represents a minimal account including recent expenditures and a
// remaining total.
type Account struct {
RecentExpenditures int
AccountTotal       int
}

// Map of account ID -> account.
var allAccounts = map[int]Account{ //nolint:gochecknoglobals
1: {RecentExpenditures: 100, AccountTotal: 1_000},
2: {RecentExpenditures: 999, AccountTotal: 1_000},
}

type ReconcileAccountArgs struct {
AccountID int `json:"account_id"`
}

func (ReconcileAccountArgs) Kind() string { return "reconcile_account" }

// InsertOpts returns custom insert options that every job of this type will
// inherit, including unique options.
func (ReconcileAccountArgs) InsertOpts() river.InsertOpts {
return river.InsertOpts{
UniqueOpts: river.UniqueOpts{
ByArgs:   true,
ByPeriod: 24 * time.Hour,
},
}
}

type ReconcileAccountWorker struct {
river.WorkerDefaults[ReconcileAccountArgs]
}

func (w *ReconcileAccountWorker) Work(ctx context.Context, job *river.Job[ReconcileAccountArgs]) error {
account := allAccounts[job.Args.AccountID]

account.AccountTotal -= account.RecentExpenditures
account.RecentExpenditures = 0

fmt.Printf("Reconciled account %d; new total: %d\n", job.Args.AccountID, account.AccountTotal)

return nil
}

// Example_uniqueJob demonstrates the use of a job with custom
// job-specific insertion options.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &ReconcileAccountWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

// First job insertion for account 1.
_, err = riverClient.Insert(ctx, ReconcileAccountArgs{AccountID: 1}, nil)
if err != nil {
panic(err)
}

// Job is inserted a second time, but it doesn't matter because its unique
// args cause the insertion to be skipped because it's meant to only run
// once per account per 24 hour period.
_, err = riverClient.Insert(ctx, ReconcileAccountArgs{AccountID: 1}, nil)
if err != nil {
panic(err)
}

// Cheat a little by waiting for the first job to come back so we can
// guarantee that this example's output comes out in order.
// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

// Because the job is unique ByArgs, another job for account 2 is allowed.
_, err = riverClient.Insert(ctx, ReconcileAccountArgs{AccountID: 2}, nil)
if err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
Reconciled account 1; new total: 900
Reconciled account 2; new total: 1
Share
Format
Run
Example (WorkFunc)
¶
Example_workFunc demonstrates the use of river.WorkFunc, which can be used to
easily add a worker with only a function instead of having to implement a
full worker struct.
package main

import (
"context"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type WorkFuncArgs struct {
Message string `json:"message"`
}

func (WorkFuncArgs) Kind() string { return "work_func" }

// Example_workFunc demonstrates the use of river.WorkFunc, which can be used to
// easily add a worker with only a function instead of having to implement a
// full worker struct.
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, river.WorkFunc(func(ctx context.Context, job *river.Job[WorkFuncArgs]) error {
fmt.Printf("Message: %s", job.Args.Message)
return nil
}))

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 100},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Out of example scope, but used to wait until a job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}

_, err = riverClient.Insert(ctx, WorkFuncArgs{
Message: "hello from a function!",
}, nil)
if err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
Message: hello from a function!
Share
Format
Run
Index
¶
Constants
Variables
func AddWorker[T JobArgs](workers *Workers, worker Worker[T])
func AddWorkerArgs[T JobArgs](workers *Workers, jobArgs T, worker Worker[T])
func AddWorkerSafely[T JobArgs](workers *Workers, worker Worker[T]) error
func JobCancel(err error) error
func JobSnooze(duration time.Duration) error
func RecordOutput(ctx context.Context, output any) error
type Client
func ClientFromContext[TTx any](ctx context.Context) *Client[TTx]
func ClientFromContextSafely[TTx any](ctx context.Context) (*Client[TTx], error)
func NewClient[TTx any](driver riverdriver.Driver[TTx], config *Config) (*Client[TTx], error)
func (c *Client[TTx]) Driver() riverdriver.Driver[TTx]
func (c *Client[TTx]) ID() string
func (c *Client[TTx]) Insert(ctx context.Context, args JobArgs, opts *InsertOpts) (*rivertype.JobInsertResult, error)
func (c *Client[TTx]) InsertMany(ctx context.Context, params []InsertManyParams) ([]*rivertype.JobInsertResult, error)
func (c *Client[TTx]) InsertManyFast(ctx context.Context, params []InsertManyParams) (int, error)
func (c *Client[TTx]) InsertManyFastTx(ctx context.Context, tx TTx, params []InsertManyParams) (int, error)
func (c *Client[TTx]) InsertManyTx(ctx context.Context, tx TTx, params []InsertManyParams) ([]*rivertype.JobInsertResult, error)
func (c *Client[TTx]) InsertTx(ctx context.Context, tx TTx, args JobArgs, opts *InsertOpts) (*rivertype.JobInsertResult, error)
func (c *Client[TTx]) JobCancel(ctx context.Context, jobID int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobCancelTx(ctx context.Context, tx TTx, jobID int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobDelete(ctx context.Context, id int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobDeleteMany(ctx context.Context, params *JobDeleteManyParams) (*JobDeleteManyResult, error)
func (c *Client[TTx]) JobDeleteManyTx(ctx context.Context, tx TTx, params *JobDeleteManyParams) (*JobDeleteManyResult, error)
func (c *Client[TTx]) JobDeleteTx(ctx context.Context, tx TTx, id int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobGet(ctx context.Context, id int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobGetTx(ctx context.Context, tx TTx, id int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobList(ctx context.Context, params *JobListParams) (*JobListResult, error)
func (c *Client[TTx]) JobListTx(ctx context.Context, tx TTx, params *JobListParams) (*JobListResult, error)
func (c *Client[TTx]) JobRetry(ctx context.Context, id int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobRetryTx(ctx context.Context, tx TTx, id int64) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobUpdate(ctx context.Context, id int64, params *JobUpdateParams) (*rivertype.JobRow, error)
func (c *Client[TTx]) JobUpdateTx(ctx context.Context, tx TTx, id int64, params *JobUpdateParams) (*rivertype.JobRow, error)
func (c *Client[TTx]) Notify() *ClientNotifyBundle[TTx]
func (c *Client[TTx]) PeriodicJobs() *PeriodicJobBundle
func (c *Client[TTx]) Pilot() riverpilot.Pilot
func (c *Client[TTx]) QueueGet(ctx context.Context, name string) (*rivertype.Queue, error)
func (c *Client[TTx]) QueueGetTx(ctx context.Context, tx TTx, name string) (*rivertype.Queue, error)
func (c *Client[TTx]) QueueList(ctx context.Context, params *QueueListParams) (*QueueListResult, error)
func (c *Client[TTx]) QueueListTx(ctx context.Context, tx TTx, params *QueueListParams) (*QueueListResult, error)
func (c *Client[TTx]) QueuePause(ctx context.Context, name string, opts *QueuePauseOpts) error
func (c *Client[TTx]) QueuePauseTx(ctx context.Context, tx TTx, name string, opts *QueuePauseOpts) error
func (c *Client[TTx]) QueueResume(ctx context.Context, name string, opts *QueuePauseOpts) error
func (c *Client[TTx]) QueueResumeTx(ctx context.Context, tx TTx, name string, opts *QueuePauseOpts) error
func (c *Client[TTx]) QueueUpdate(ctx context.Context, name string, params *QueueUpdateParams) (*rivertype.Queue, error)
func (c *Client[TTx]) QueueUpdateTx(ctx context.Context, tx TTx, name string, params *QueueUpdateParams) (*rivertype.Queue, error)
func (c *Client[TTx]) Queues() *QueueBundle
func (c *Client[TTx]) Schema() string
func (c *Client[TTx]) Start(ctx context.Context) error
func (c *Client[TTx]) Stop(ctx context.Context) error
func (c *Client[TTx]) StopAndCancel(ctx context.Context) error
func (c *Client[TTx]) Stopped() <-chan struct{}
func (c *Client[TTx]) Subscribe(kinds ...EventKind) (<-chan *Event, func())
func (c *Client[TTx]) SubscribeConfig(config *SubscribeConfig) (<-chan *Event, func())
type ClientNotifyBundle
func (c *ClientNotifyBundle[TTx]) RequestResign(ctx context.Context) error
func (c *ClientNotifyBundle[TTx]) RequestResignTx(ctx context.Context, tx TTx) error
type ClientRetryPolicy
type Config
func (c *Config) WithDefaults() *Config
type DefaultClientRetryPolicy
func (p *DefaultClientRetryPolicy) NextRetry(job *rivertype.JobRow) time.Time
type ErrorHandler
type ErrorHandlerResult
type Event
type EventKind
type HookDefaults
func (d *HookDefaults) IsHook() bool
type HookInsertBeginFunc
func (f HookInsertBeginFunc) InsertBegin(ctx context.Context, params *rivertype.JobInsertParams) error
func (f HookInsertBeginFunc) IsHook() bool
type HookPeriodicJobsStartFunc
func (f HookPeriodicJobsStartFunc) IsHook() bool
func (f HookPeriodicJobsStartFunc) Start(ctx context.Context, params *rivertype.HookPeriodicJobsStartParams) error
type HookWorkBeginFunc
func (f HookWorkBeginFunc) IsHook() bool
func (f HookWorkBeginFunc) WorkBegin(ctx context.Context, job *rivertype.JobRow) error
type HookWorkEndFunc
func (f HookWorkEndFunc) IsHook() bool
func (f HookWorkEndFunc) WorkEnd(ctx context.Context, job *rivertype.JobRow, err error) error
type InsertManyParams
type InsertOpts
type Job
func JobCompleteTx[TDriver riverdriver.Driver[TTx], TTx any, TArgs JobArgs](ctx context.Context, tx TTx, job *Job[TArgs]) (*Job[TArgs], error)
type JobArgs
type JobArgsWithHooks
type JobArgsWithInsertOpts
type JobArgsWithKindAliases
type JobCancelError
type JobDeleteManyParams
func NewJobDeleteManyParams() *JobDeleteManyParams
func (p *JobDeleteManyParams) First(count int) *JobDeleteManyParams
func (p *JobDeleteManyParams) IDs(ids ...int64) *JobDeleteManyParams
func (p *JobDeleteManyParams) Kinds(kinds ...string) *JobDeleteManyParams
func (p *JobDeleteManyParams) Priorities(priorities ...int16) *JobDeleteManyParams
func (p *JobDeleteManyParams) Queues(queues ...string) *JobDeleteManyParams
func (p *JobDeleteManyParams) States(states ...rivertype.JobState) *JobDeleteManyParams
func (p *JobDeleteManyParams) UnsafeAll() *JobDeleteManyParams
type JobDeleteManyResult
type JobInsertMiddlewareDefaults
deprecated
func (d *JobInsertMiddlewareDefaults) InsertMany(ctx context.Context, manyParams []*rivertype.JobInsertParams, ...) ([]*rivertype.JobInsertResult, error)
type JobInsertMiddlewareFunc
func (f JobInsertMiddlewareFunc) InsertMany(ctx context.Context, manyParams []*rivertype.JobInsertParams, ...) ([]*rivertype.JobInsertResult, error)
func (f JobInsertMiddlewareFunc) IsMiddleware() bool
type JobListCursor
func JobListCursorFromJob(job *rivertype.JobRow) *JobListCursor
func (c JobListCursor) MarshalText() ([]byte, error)
func (c *JobListCursor) UnmarshalText(text []byte) error
type JobListOrderByField
type JobListParams
func NewJobListParams() *JobListParams
func (p *JobListParams) After(cursor *JobListCursor) *JobListParams
func (p *JobListParams) First(count int) *JobListParams
func (p *JobListParams) IDs(ids ...int64) *JobListParams
func (p *JobListParams) Kinds(kinds ...string) *JobListParams
func (p *JobListParams) Metadata(json string) *JobListParams
func (p *JobListParams) OrderBy(field JobListOrderByField, direction SortOrder) *JobListParams
func (p *JobListParams) Priorities(priorities ...int16) *JobListParams
func (p *JobListParams) Queues(queues ...string) *JobListParams
func (p *JobListParams) States(states ...rivertype.JobState) *JobListParams
func (p *JobListParams) Where(sql string, namedArgsMany ...NamedArgs) *JobListParams
type JobListResult
type JobSnoozeError
type JobStatistics
type JobUpdateParams
type MiddlewareDefaults
func (d *MiddlewareDefaults) IsMiddleware() bool
type NamedArgs
type PeriodicJob
func NewPeriodicJob(scheduleFunc PeriodicSchedule, constructorFunc PeriodicJobConstructor, ...) *PeriodicJob
type PeriodicJobBundle
func (b *PeriodicJobBundle) Add(periodicJob *PeriodicJob) rivertype.PeriodicJobHandle
func (b *PeriodicJobBundle) AddMany(periodicJobs []*PeriodicJob) []rivertype.PeriodicJobHandle
func (b *PeriodicJobBundle) AddManySafely(periodicJobs []*PeriodicJob) ([]rivertype.PeriodicJobHandle, error)
func (b *PeriodicJobBundle) AddSafely(periodicJob *PeriodicJob) (rivertype.PeriodicJobHandle, error)
func (b *PeriodicJobBundle) Clear()
func (b *PeriodicJobBundle) Remove(periodicJobHandle rivertype.PeriodicJobHandle)
func (b *PeriodicJobBundle) RemoveByID(id string) bool
func (b *PeriodicJobBundle) RemoveMany(periodicJobHandles []rivertype.PeriodicJobHandle)
func (b *PeriodicJobBundle) RemoveManyByID(ids []string)
type PeriodicJobConstructor
type PeriodicJobOpts
type PeriodicSchedule
func NeverSchedule() PeriodicSchedule
func PeriodicInterval(interval time.Duration) PeriodicSchedule
type QueueAlreadyAddedError
func (e *QueueAlreadyAddedError) Error() string
func (e *QueueAlreadyAddedError) Is(target error) bool
type QueueBundle
func (b *QueueBundle) Add(queueName string, queueConfig QueueConfig) error
type QueueConfig
type QueueListParams
func NewQueueListParams() *QueueListParams
func (p *QueueListParams) First(count int) *QueueListParams
type QueueListResult
type QueuePauseOpts
type QueueUpdateParams
type SortOrder
type SubscribeConfig
type TestConfig
type UniqueOpts
type UnknownJobKindError
type Worker
func WorkFunc[T JobArgs](f func(context.Context, *Job[T]) error) Worker[T]
type WorkerDefaults
func (w WorkerDefaults[T]) Middleware(*rivertype.JobRow) []rivertype.WorkerMiddleware
func (w WorkerDefaults[T]) NextRetry(*Job[T]) time.Time
func (w WorkerDefaults[T]) Timeout(*Job[T]) time.Duration
type WorkerMiddlewareDefaults
deprecated
func (d *WorkerMiddlewareDefaults) Work(ctx context.Context, job *rivertype.JobRow, ...) error
type WorkerMiddlewareFunc
func (f WorkerMiddlewareFunc) IsMiddleware() bool
func (f WorkerMiddlewareFunc) Work(ctx context.Context, job *rivertype.JobRow, ...) error
type Workers
func NewWorkers() *Workers
Examples
¶
Package (BatchInsert)
Package (CompleteJobWithinTx)
Package (CronJob)
Package (CustomInsertOpts)
Package (ErrorHandler)
Package (GlobalHooks)
Package (GlobalMiddleware)
Package (GracefulShutdown)
Package (InsertAndWork)
Package (JobArgsHooks)
Package (JobCancel)
Package (JobCancelFromClient)
Package (JobSnooze)
Package (PeriodicJob)
Package (QueuePause)
Package (ScheduledJob)
Package (Subscription)
Package (UniqueJob)
Package (WorkFunc)
ClientFromContext (Pgx)
Constants
¶
View Source
const (
FetchCooldownDefault = 100 *
time
.
Millisecond
FetchCooldownMin     = 1 *
time
.
Millisecond
FetchPollIntervalDefault = 1 *
time
.
Second
FetchPollIntervalMin     = 1 *
time
.
Millisecond
JobTimeoutDefault  = 1 *
time
.
Minute
MaxAttemptsDefault =
rivercommon
.
MaxAttemptsDefault
PriorityDefault    =
rivercommon
.
PriorityDefault
QueueDefault       =
rivercommon
.
QueueDefault
QueueNumWorkersMax = 10_000
)
Variables
¶
View Source
var ErrJobCancelledRemotely =
rivertype
.
ErrJobCancelledRemotely
ErrJobCancelledRemotely is a sentinel error indicating that the job was cancelled remotely.
View Source
var (
// ErrNotFound is returned when a query by ID does not match any existing
// rows. For example, attempting to cancel a job that doesn't exist will
// return this error.
ErrNotFound =
rivertype
.
ErrNotFound
)
Functions
¶
func
AddWorker
¶
func AddWorker[T
JobArgs
](workers *
Workers
, worker
Worker
[T])
AddWorker registers a Worker on the provided Workers bundle. Each Worker must
be registered so that the Client knows it should handle a specific kind of
job (as returned by its `Kind()` method).
Use by explicitly specifying a JobArgs type and then passing an instance of a
worker for the same type:
river.AddWorker(workers, &SortWorker{})
Note that AddWorker can panic in some situations, such as if the worker is
already registered or if its configuration is otherwise invalid. This default
probably makes sense for most applications because you wouldn't want to start
an application with invalid hardcoded runtime configuration. If you want to
avoid panics, use AddWorkerSafely instead.
func
AddWorkerArgs
¶
added in
v0.24.0
func AddWorkerArgs[T
JobArgs
](workers *
Workers
, jobArgs T, worker
Worker
[T])
AddWorkerArgs is the same as AddWorker except that it lets args be passed
explicitly rather than being instantiated implicitly. We don't know of any
use for this function beyond exercising some args-related edge cases in tests
are difficult/impossible to exercise otherwise, and its use should be
considered internal only.
func
AddWorkerSafely
¶
func AddWorkerSafely[T
JobArgs
](workers *
Workers
, worker
Worker
[T])
error
AddWorkerSafely registers a worker on the provided Workers bundle. Unlike AddWorker,
AddWorkerSafely does not panic and instead returns an error if the worker
is already registered or if its configuration is invalid.
Use by explicitly specifying a JobArgs type and then passing an instance of a
worker for the same type:
river.AddWorkerSafely[SortArgs](workers, &SortWorker{}).
func
JobCancel
¶
func JobCancel(err
error
)
error
JobCancel wraps err and can be returned from a Worker's Work method to cancel
the job at the end of execution. Regardless of whether or not the job has any
remaining attempts, this will ensure the job does not execute again.
func
JobSnooze
¶
func JobSnooze(duration
time
.
Duration
)
error
JobSnooze can be returned from a Worker's Work method to cause the job to be
tried again after the specified duration. This also has the effect of
incrementing the job's MaxAttempts by 1, meaning that jobs can be repeatedly
snoozed without ever being discarded.
A special duration of zero can be used to make the job immediately available
to be reworked. This may be useful in cases like where a long-running job is
being interrupted on shutdown. Instead of returning a context cancelled error
that'd schedule a retry for the future and count towards maximum attempts,
the work function can return JobSnooze(0) and the job will be retried
immediately the next time a client starts up.
Panics if duration is < 0.
func
RecordOutput
¶
added in
v0.18.0
func RecordOutput(ctx
context
.
Context
, output
any
)
error
RecordOutput records output JSON from a job. The "output" can be any
JSON-encodable value and will be stored in the database on the job row after
the current execution attempt completes. Output may be useful for debugging,
or for storing the result of a job temporarily without needing to create a
dedicated table to keep it in.
For example, with workflows, it's common for subsequent task to depend on
something done in an earlier dependency task. Consider the creation of an
external resource in another API or in an database—it will typically have a
unique ID that must be used to reference the resource later. A later step
may require that info in order to complete its work, and the output can be
a convenient way to store that info.
Output is stored in the job's metadata under the `"output"` key
(
github.com/riverqueue/river/rivertype.MetadataKeyOutput
).
This function must be called within an Worker's Work function. It returns an
error if called anywhere else. As with any stored value, care should be taken
to ensure that the payload size is not too large. Output is limited to 32MB
in size for safety, but should be kept much smaller than this.
Only one output can be stored per job. If this function is called more than
once, the output will be overwritten with the latest value. The output also
must be recorded _before_ the job finishes executing so that it can be stored
when the job's row is updated.
Once recorded, the output is stored regardless of the outcome of the
execution attempt (success, error, panic, etc.).
RecordOutput always stores output lazily as a job is being completed (whether
that's completion to success or failure). Client.JobUpdate and JobUpdateTx
are available to store output eagerly at any time, including from inside a
work function as the job is being executed.
The output is marshalled to JSON as part of this function and it will return
an error if the output is not JSON-encodable.
Types
¶
type
Client
¶
type Client[TTx
any
] struct {
// contains filtered or unexported fields
}
Client is a single isolated instance of River. Your application may use
multiple instances operating on different databases or Postgres schemas
within a single database.
func
ClientFromContext
¶
added in
v0.0.17
func ClientFromContext[TTx
any
](ctx
context
.
Context
) *
Client
[TTx]
ClientFromContext returns the Client from the context. This function can
only be used within a Worker's Work() method because that is the only place
River sets the Client on the context.
It panics if the context does not contain a Client, which will never happen
from the context provided to a Worker's Work() method.
When testing JobArgs.Work implementations, it might be useful to use
rivertest.WorkContext to initialize a context that has an available client.
The type parameter TTx is the transaction type used by the
Client
,
pgx.Tx for the pgx driver, and *sql.Tx for the
database/sql
driver.
Example (Pgx)
¶
ExampleClientFromContext_pgx demonstrates how to extract the River client
from the worker context when using the pgx/v5 driver.
(
github.com/riverqueue/river/riverdriver/riverpgxv5
).
package main

import (
"context"
"errors"
"fmt"
"log/slog"
"os"

"github.com/jackc/pgx/v5"
"github.com/jackc/pgx/v5/pgxpool"

"github.com/riverqueue/river"
"github.com/riverqueue/river/riverdbtest"
"github.com/riverqueue/river/riverdriver/riverpgxv5"
"github.com/riverqueue/river/rivershared/riversharedtest"
"github.com/riverqueue/river/rivershared/util/slogutil"
"github.com/riverqueue/river/rivershared/util/testutil"
)

type ContextClientArgs struct{}

func (args ContextClientArgs) Kind() string { return "ContextClientWorker" }

type ContextClientWorker struct {
river.WorkerDefaults[ContextClientArgs]
}

func (w *ContextClientWorker) Work(ctx context.Context, job *river.Job[ContextClientArgs]) error {
client := river.ClientFromContext[pgx.Tx](ctx)
if client == nil {
fmt.Println("client not found in context")
return errors.New("client not found in context")
}

fmt.Printf("client found in context, id=%s\n", client.ID())
return nil
}

// ExampleClientFromContext_pgx demonstrates how to extract the River client
// from the worker context when using the pgx/v5 driver.
// ([github.com/riverqueue/river/riverdriver/riverpgxv5]).
func main() {
ctx := context.Background()

dbPool, err := pgxpool.New(ctx, riversharedtest.TestDatabaseURL())
if err != nil {
panic(err)
}
defer dbPool.Close()

workers := river.NewWorkers()
river.AddWorker(workers, &ContextClientWorker{})

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
ID:     "ClientFromContextClient",
Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
Queues: map[string]river.QueueConfig{
river.QueueDefault: {MaxWorkers: 10},
},
Schema:   riverdbtest.TestSchema(ctx, testutil.PanicTB(), riverpgxv5.New(dbPool), nil), // only necessary for the example test
TestOnly: true,                                                                         // suitable only for use in tests; remove for live environments
Workers:  workers,
})
if err != nil {
panic(err)
}

// Not strictly needed, but used to help this test wait until job is worked.
subscribeChan, subscribeCancel := riverClient.Subscribe(river.EventKindJobCompleted)
defer subscribeCancel()

if err := riverClient.Start(ctx); err != nil {
panic(err)
}
if _, err = riverClient.Insert(ctx, ContextClientArgs{}, nil); err != nil {
panic(err)
}

// Wait for jobs to complete. Only needed for purposes of the example test.
riversharedtest.WaitOrTimeoutN(testutil.PanicTB(), subscribeChan, 1)

if err := riverClient.Stop(ctx); err != nil {
panic(err)
}

}
Output:
client found in context, id=ClientFromContextClient
Share
Format
Run
func
ClientFromContextSafely
¶
added in
v0.0.17
func ClientFromContextSafely[TTx
any
](ctx
context
.
Context
) (*
Client
[TTx],
error
)
ClientFromContextSafely returns the Client from the context. This function
can only be used within a Worker's Work() method because that is the only
place River sets the Client on the context.
It returns an error if the context does not contain a Client, which will
never happen from the context provided to a Worker's Work() method.
When testing JobArgs.Work implementations, it might be useful to use
rivertest.WorkContext to initialize a context that has an available client.
See the examples for
ClientFromContext
to understand how to use this
function.
func
NewClient
¶
func NewClient[TTx
any
](driver
riverdriver
.
Driver
[TTx], config *
Config
) (*
Client
[TTx],
error
)
NewClient creates a new Client with the given database driver and
configuration.
Currently only one driver is supported, which is Pgx v5. See package
riverpgxv5.
The function takes a generic parameter TTx representing a transaction type,
but it can be omitted because it'll generally always be inferred from the
driver. For example:
import "github.com/riverqueue/river"
import "github.com/riverqueue/river/riverdriver/riverpgxv5"

...

dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
if err != nil {
// handle error
}
defer dbPool.Close()

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
...
})
if err != nil {
// handle error
}
func (*Client[TTx])
Driver
¶
added in
v0.11.0
func (c *
Client
[TTx]) Driver()
riverdriver
.
Driver
[TTx]
Driver exposes the underlying driver used by the client.
API is not stable. DO NOT USE.
func (*Client[TTx])
ID
¶
added in
v0.0.21
func (c *
Client
[TTx]) ID()
string
ID returns the unique ID of this client as set in its config or
auto-generated if not specified.
func (*Client[TTx])
Insert
¶
func (c *
Client
[TTx]) Insert(ctx
context
.
Context
, args
JobArgs
, opts *
InsertOpts
) (*
rivertype
.
JobInsertResult
,
error
)
Insert inserts a new job with the provided args. Job opts can be used to
override any defaults that may have been provided by an implementation of
JobArgsWithInsertOpts.InsertOpts, as well as any global defaults. The
provided context is used for the underlying Postgres insert and can be used
to cancel the operation or apply a timeout.
jobRow, err := client.Insert(insertCtx, MyArgs{}, nil)
if err != nil {
// handle error
}
func (*Client[TTx])
InsertMany
¶
func (c *
Client
[TTx]) InsertMany(ctx
context
.
Context
, params []
InsertManyParams
) ([]*
rivertype
.
JobInsertResult
,
error
)
InsertMany inserts many jobs at once. Each job is inserted as an
InsertManyParams tuple, which takes job args along with an optional set of
insert options, which override insert options provided by an
JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults.
The provided context is used for the underlying Postgres inserts and can be
used to cancel the operation or apply a timeout.
count, err := client.InsertMany(ctx, []river.InsertManyParams{
{Args: BatchInsertArgs{}},
{Args: BatchInsertArgs{}, InsertOpts: &river.InsertOpts{Priority: 3}},
})
if err != nil {
// handle error
}
func (*Client[TTx])
InsertManyFast
¶
added in
v0.12.0
func (c *
Client
[TTx]) InsertManyFast(ctx
context
.
Context
, params []
InsertManyParams
) (
int
,
error
)
InsertManyFast inserts many jobs at once using Postgres' `COPY FROM` mechanism,
making the operation quite fast and memory efficient. Each job is inserted as
an InsertManyParams tuple, which takes job args along with an optional set of
insert options, which override insert options provided by an
JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults.
The provided context is used for the underlying Postgres inserts and can be
used to cancel the operation or apply a timeout.
count, err := client.InsertMany(ctx, []river.InsertManyParams{
{Args: BatchInsertArgs{}},
{Args: BatchInsertArgs{}, InsertOpts: &river.InsertOpts{Priority: 3}},
})
if err != nil {
// handle error
}
Unlike with `InsertMany`, unique conflicts cannot be handled gracefully. If a
unique constraint is violated, the operation will fail and no jobs will be inserted.
func (*Client[TTx])
InsertManyFastTx
¶
added in
v0.12.0
func (c *
Client
[TTx]) InsertManyFastTx(ctx
context
.
Context
, tx TTx, params []
InsertManyParams
) (
int
,
error
)
InsertManyFastTx inserts many jobs at once using Postgres' `COPY FROM`
mechanism, making the operation quite fast and memory efficient. Each job is
inserted as an InsertManyParams tuple, which takes job args along with an
optional set of insert options, which override insert options provided by an
JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults.
The provided context is used for the underlying Postgres inserts and can be
used to cancel the operation or apply a timeout.
count, err := client.InsertManyTx(ctx, tx, []river.InsertManyParams{
{Args: BatchInsertArgs{}},
{Args: BatchInsertArgs{}, InsertOpts: &river.InsertOpts{Priority: 3}},
})
if err != nil {
// handle error
}
This variant lets a caller insert jobs atomically alongside other database
changes. An inserted job isn't visible to be worked until the transaction
commits, and if the transaction rolls back, so too is the inserted job.
Unlike with `InsertManyTx`, unique conflicts cannot be handled gracefully. If
a unique constraint is violated, the operation will fail and no jobs will be
inserted.
func (*Client[TTx])
InsertManyTx
¶
func (c *
Client
[TTx]) InsertManyTx(ctx
context
.
Context
, tx TTx, params []
InsertManyParams
) ([]*
rivertype
.
JobInsertResult
,
error
)
InsertManyTx inserts many jobs at once. Each job is inserted as an
InsertManyParams tuple, which takes job args along with an optional set of
insert options, which override insert options provided by an
JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults.
The provided context is used for the underlying Postgres inserts and can be
used to cancel the operation or apply a timeout.
count, err := client.InsertManyTx(ctx, tx, []river.InsertManyParams{
{Args: BatchInsertArgs{}},
{Args: BatchInsertArgs{}, InsertOpts: &river.InsertOpts{Priority: 3}},
})
if err != nil {
// handle error
}
This variant lets a caller insert jobs atomically alongside other database
changes. An inserted job isn't visible to be worked until the transaction
commits, and if the transaction rolls back, so too is the inserted job.
func (*Client[TTx])
InsertTx
¶
func (c *
Client
[TTx]) InsertTx(ctx
context
.
Context
, tx TTx, args
JobArgs
, opts *
InsertOpts
) (*
rivertype
.
JobInsertResult
,
error
)
InsertTx inserts a new job with the provided args on the given transaction.
Job opts can be used to override any defaults that may have been provided by
an implementation of JobArgsWithInsertOpts.InsertOpts, as well as any global
defaults. The provided context is used for the underlying Postgres insert and
can be used to cancel the operation or apply a timeout.
jobRow, err := client.InsertTx(insertCtx, tx, MyArgs{}, nil)
if err != nil {
// handle error
}
This variant lets a caller insert jobs atomically alongside other database
changes. It's also possible to insert a job outside a transaction, but this
usage is recommended to ensure that all data a job needs to run is available
by the time it starts. Because of snapshot visibility guarantees across
transactions, the job will not be worked until the transaction has committed,
and if the transaction rolls back, so too is the inserted job.
func (*Client[TTx])
JobCancel
¶
added in
v0.0.17
func (c *
Client
[TTx]) JobCancel(ctx
context
.
Context
, jobID
int64
) (*
rivertype
.
JobRow
,
error
)
JobCancel cancels the job with the given ID. If possible, the job is
cancelled immediately and will not be retried. The provided context is used
for the underlying Postgres update and can be used to cancel the operation or
apply a timeout.
If the job is still in the queue (available, scheduled, or retryable), it is
immediately marked as cancelled and will not be retried.
If the job is already finalized (cancelled, completed, or discarded), no
changes are made.
If the job is currently running, it is not immediately cancelled, but is
instead marked for cancellation. The client running the job will also be
notified (via LISTEN/NOTIFY) to cancel the running job's context. Although
the job's context will be cancelled, since Go does not provide a mechanism to
interrupt a running goroutine the job will continue running until it returns.
As always, it is important for workers to respect context cancellation and
return promptly when the job context is done.
Once the cancellation signal is received by the client running the job, any
error returned by that job will result in it being cancelled permanently and
not retried. However if the job returns no error, it will be completed as
usual.
In the event the running job finishes executing _before_ the cancellation
signal is received but _after_ this update was made, the behavior depends on
which state the job is being transitioned into (based on its return error):
If the job completed successfully, was cancelled from within, or was
discarded due to exceeding its max attempts, the job will be updated as
usual.
If the job was snoozed to run again later or encountered a retryable error,
the job will be marked as cancelled and will not be attempted again.
Returns the up-to-date JobRow for the specified jobID if it exists. Returns
ErrNotFound if the job doesn't exist.
func (*Client[TTx])
JobCancelTx
¶
added in
v0.0.17
func (c *
Client
[TTx]) JobCancelTx(ctx
context
.
Context
, tx TTx, jobID
int64
) (*
rivertype
.
JobRow
,
error
)
JobCancelTx cancels the job with the given ID within the specified
transaction. This variant lets a caller cancel a job atomically alongside
other database changes. A cancelled job doesn't take effect until the
transaction commits, and if the transaction rolls back, so too is the
cancelled job.
If possible, the job is cancelled immediately and will not be retried. The
provided context is used for the underlying Postgres update and can be used
to cancel the operation or apply a timeout.
If the job is still in the queue (available, scheduled, or retryable), it is
immediately marked as cancelled and will not be retried.
If the job is already finalized (cancelled, completed, or discarded), no
changes are made.
If the job is currently running, it is not immediately cancelled, but is
instead marked for cancellation. The client running the job will also be
notified (via LISTEN/NOTIFY) to cancel the running job's context. Although
the job's context will be cancelled, since Go does not provide a mechanism to
interrupt a running goroutine the job will continue running until it returns.
As always, it is important for workers to respect context cancellation and
return promptly when the job context is done.
Once the cancellation signal is received by the client running the job, any
error returned by that job will result in it being cancelled permanently and
not retried. However if the job returns no error, it will be completed as
usual.
In the event the running job finishes executing _before_ the cancellation
signal is received but _after_ this update was made, the behavior depends on
which state the job is being transitioned into (based on its return error):
If the job completed successfully, was cancelled from within, or was
discarded due to exceeding its max attempts, the job will be updated as
usual.
If the job was snoozed to run again later or encountered a retryable error,
the job will be marked as cancelled and will not be attempted again.
Returns the up-to-date JobRow for the specified jobID if it exists. Returns
ErrNotFound if the job doesn't exist.
func (*Client[TTx])
JobDelete
¶
added in
v0.7.0
func (c *
Client
[TTx]) JobDelete(ctx
context
.
Context
, id
int64
) (*
rivertype
.
JobRow
,
error
)
JobDelete deletes the job with the given ID from the database, returning the
deleted row if it was deleted. Jobs in the running state are not deleted,
instead returning rivertype.ErrJobRunning.
func (*Client[TTx])
JobDeleteMany
¶
added in
v0.24.0
func (c *
Client
[TTx]) JobDeleteMany(ctx
context
.
Context
, params *
JobDeleteManyParams
) (*
JobDeleteManyResult
,
error
)
JobDeleteMany deletes many jobs at once based on the conditions defined by
JobDeleteManyParams. Running jobs are always ignored.
params := river.NewJobDeleteManyParams().First(10).State(rivertype.JobStateCompleted)
jobRows, err := client.JobDeleteMany(ctx, params)
if err != nil {
// handle error
}
func (*Client[TTx])
JobDeleteManyTx
¶
added in
v0.24.0
func (c *
Client
[TTx]) JobDeleteManyTx(ctx
context
.
Context
, tx TTx, params *
JobDeleteManyParams
) (*
JobDeleteManyResult
,
error
)
JobDeleteManyTx deletes many jobs at once based on the conditions defined by
JobDeleteManyParams. Running jobs are always ignored.
params := river.NewJobDeleteManyParams().First(10).States(river.JobStateCompleted)
jobRows, err := client.JobDeleteManyTx(ctx, tx, params)
if err != nil {
// handle error
}
func (*Client[TTx])
JobDeleteTx
¶
added in
v0.7.0
func (c *
Client
[TTx]) JobDeleteTx(ctx
context
.
Context
, tx TTx, id
int64
) (*
rivertype
.
JobRow
,
error
)
JobDeleteTx deletes the job with the given ID from the database, returning the
deleted row if it was deleted. Jobs in the running state are not deleted,
instead returning rivertype.ErrJobRunning. This variant lets a caller retry a
job atomically alongside other database changes. A deleted job isn't deleted
until the transaction commits, and if the transaction rolls back, so too is
the deleted job.
func (*Client[TTx])
JobGet
¶
added in
v0.0.19
func (c *
Client
[TTx]) JobGet(ctx
context
.
Context
, id
int64
) (*
rivertype
.
JobRow
,
error
)
JobGet fetches a single job by its ID. Returns the up-to-date JobRow for the
specified jobID if it exists. Returns ErrNotFound if the job doesn't exist.
func (*Client[TTx])
JobGetTx
¶
added in
v0.0.19
func (c *
Client
[TTx]) JobGetTx(ctx
context
.
Context
, tx TTx, id
int64
) (*
rivertype
.
JobRow
,
error
)
JobGetTx fetches a single job by its ID, within a transaction. Returns the
up-to-date JobRow for the specified jobID if it exists. Returns ErrNotFound
if the job doesn't exist.
func (*Client[TTx])
JobList
¶
added in
v0.0.17
func (c *
Client
[TTx]) JobList(ctx
context
.
Context
, params *
JobListParams
) (*
JobListResult
,
error
)
JobList returns a paginated list of jobs matching the provided filters. The
provided context is used for the underlying Postgres query and can be used to
cancel the operation or apply a timeout.
params := river.NewJobListParams().First(10).State(rivertype.JobStateCompleted)
jobRows, err := client.JobList(ctx, params)
if err != nil {
// handle error
}
func (*Client[TTx])
JobListTx
¶
added in
v0.0.17
func (c *
Client
[TTx]) JobListTx(ctx
context
.
Context
, tx TTx, params *
JobListParams
) (*
JobListResult
,
error
)
JobListTx returns a paginated list of jobs matching the provided filters. The
provided context is used for the underlying Postgres query and can be used to
cancel the operation or apply a timeout.
params := river.NewJobListParams().First(10).States(river.JobStateCompleted)
jobRows, err := client.JobListTx(ctx, tx, params)
if err != nil {
// handle error
}
func (*Client[TTx])
JobRetry
¶
added in
v0.0.19
func (c *
Client
[TTx]) JobRetry(ctx
context
.
Context
, id
int64
) (*
rivertype
.
JobRow
,
error
)
JobRetry updates the job with the given ID to make it immediately available
to be retried. Jobs in the running state are not touched, while jobs in any
other state are made available. To prevent jobs already waiting in the queue
from being set back in line, the job's scheduled_at field is set to the
current time only if it's not already in the past.
MaxAttempts is also incremented by one if the job has already exhausted its
max attempts.
func (*Client[TTx])
JobRetryTx
¶
added in
v0.0.19
func (c *
Client
[TTx]) JobRetryTx(ctx
context
.
Context
, tx TTx, id
int64
) (*
rivertype
.
JobRow
,
error
)
JobRetryTx updates the job with the given ID to make it immediately available
to be retried, within the specified transaction. This variant lets a caller
retry a job atomically alongside other database changes. A retried job isn't
visible to be worked until the transaction commits, and if the transaction
rolls back, so too is the retried job.
Jobs in the running state are not touched, while jobs in any other state are
made available. To prevent jobs already waiting in the queue from being set
back in line, the job's scheduled_at field is set to the current time only if
it's not already in the past.
MaxAttempts is also incremented by one if the job has already exhausted its
max attempts.
func (*Client[TTx])
JobUpdate
¶
added in
v0.29.0
func (c *
Client
[TTx]) JobUpdate(ctx
context
.
Context
, id
int64
, params *
JobUpdateParams
) (*
rivertype
.
JobRow
,
error
)
JobUpdate updates the job with the given ID.
If JobUpdateParams.Output is not set, this function may be used inside a job
work function to set a job's output based on output recorded so far using
RecordOutput.
func (*Client[TTx])
JobUpdateTx
¶
added in
v0.29.0
func (c *
Client
[TTx]) JobUpdateTx(ctx
context
.
Context
, tx TTx, id
int64
, params *
JobUpdateParams
) (*
rivertype
.
JobRow
,
error
)
JobUpdateTx updates the job with the given ID.
If JobUpdateParams.Output is not set, this function may be used inside a job
work function to set a job's output based on output recorded so far using
RecordOutput.
This variant updates the job inside of a transaction.
func (*Client[TTx])
Notify
¶
added in
v0.29.0
func (c *
Client
[TTx]) Notify() *
ClientNotifyBundle
[TTx]
Notify retrieves a notification bundle for the client (in the sense of
Postgres listen/notify) used to send notifications of various kinds.
func (*Client[TTx])
PeriodicJobs
¶
added in
v0.2.0
func (c *
Client
[TTx]) PeriodicJobs() *
PeriodicJobBundle
PeriodicJobs returns the currently configured set of periodic jobs for the
client, and can be used to add new or remove existing ones.
This function should only be invoked on clients capable of running perioidc
jobs. Running periodic jobs requires that the client be electable as leader
to run maintenance services, and being electable as leader requires that a
client be started. To be startable, a client must have Queues and Workers
configured. Invoking this function will panic if these conditions aren't met.
func (*Client[TTx])
Pilot
¶
added in
v0.13.0
func (c *
Client
[TTx]) Pilot()
riverpilot
.
Pilot
Pilot returns the pilot in use by the pilot. If not configured, this is often
simply StandardPilot.
API is not stable. DO NOT USE.
func (*Client[TTx])
QueueGet
¶
added in
v0.5.0
func (c *
Client
[TTx]) QueueGet(ctx
context
.
Context
, name
string
) (*
rivertype
.
Queue
,
error
)
QueueGet returns the queue with the given name. If the queue has not recently
been active or does not exist, returns ErrNotFound.
The provided context is used for the underlying Postgres query and can be
used to cancel the operation or apply a timeout.
func (*Client[TTx])
QueueGetTx
¶
added in
v0.8.0
func (c *
Client
[TTx]) QueueGetTx(ctx
context
.
Context
, tx TTx, name
string
) (*
rivertype
.
Queue
,
error
)
QueueGetTx returns the queue with the given name. If the queue has not recently
been active or does not exist, returns ErrNotFound.
The provided context is used for the underlying Postgres query and can be
used to cancel the operation or apply a timeout.
func (*Client[TTx])
QueueList
¶
added in
v0.5.0
func (c *
Client
[TTx]) QueueList(ctx
context
.
Context
, params *
QueueListParams
) (*
QueueListResult
,
error
)
QueueList returns a list of all queues that are currently active or were
recently active. Limit and offset can be used to paginate the results.
The provided context is used for the underlying Postgres query and can be
used to cancel the operation or apply a timeout.
params := river.NewQueueListParams().First(10)
queueRows, err := client.QueueListTx(ctx, tx, params)
if err != nil {
// handle error
}
func (*Client[TTx])
QueueListTx
¶
added in
v0.8.0
func (c *
Client
[TTx]) QueueListTx(ctx
context
.
Context
, tx TTx, params *
QueueListParams
) (*
QueueListResult
,
error
)
QueueListTx returns a list of all queues that are currently active or were
recently active. Limit and offset can be used to paginate the results.
The provided context is used for the underlying Postgres query and can be
used to cancel the operation or apply a timeout.
params := river.NewQueueListParams().First(10)
queueRows, err := client.QueueListTx(ctx, tx, params)
if err != nil {
// handle error
}
func (*Client[TTx])
QueuePause
¶
added in
v0.5.0
func (c *
Client
[TTx]) QueuePause(ctx
context
.
Context
, name
string
, opts *
QueuePauseOpts
)
error
QueuePause pauses the queue with the given name. When a queue is paused,
clients will not fetch any more jobs for that particular queue. To pause all
queues at once, use the special queue name "*".
Clients with a configured notifier should receive a notification about the
paused queue(s) within a few milliseconds of the transaction commit. Clients
in poll-only mode will pause after their next poll for queue configuration.
The provided context is used for the underlying Postgres update and can be
used to cancel the operation or apply a timeout. The opts are reserved for
future functionality.
func (*Client[TTx])
QueuePauseTx
¶
added in
v0.8.0
func (c *
Client
[TTx]) QueuePauseTx(ctx
context
.
Context
, tx TTx, name
string
, opts *
QueuePauseOpts
)
error
QueuePauseTx pauses the queue with the given name. When a queue is paused,
clients will not fetch any more jobs for that particular queue. To pause all
queues at once, use the special queue name "*".
Clients with a configured notifier should receive a notification about the
paused queue(s) within a few milliseconds of the transaction commit. Clients
in poll-only mode will pause after their next poll for queue configuration.
The provided context is used for the underlying Postgres update and can be
used to cancel the operation or apply a timeout. The opts are reserved for
future functionality.
func (*Client[TTx])
QueueResume
¶
added in
v0.5.0
func (c *
Client
[TTx]) QueueResume(ctx
context
.
Context
, name
string
, opts *
QueuePauseOpts
)
error
QueueResume resumes the queue with the given name. If the queue was
previously paused, any clients configured to work that queue will resume
fetching additional jobs. To resume all queues at once, use the special queue
name "*".
Clients with a configured notifier should receive a notification about the
resumed queue(s) within a few milliseconds of the transaction commit. Clients
in poll-only mode will resume after their next poll for queue configuration.
The provided context is used for the underlying Postgres update and can be
used to cancel the operation or apply a timeout. The opts are reserved for
future functionality.
func (*Client[TTx])
QueueResumeTx
¶
added in
v0.8.0
func (c *
Client
[TTx]) QueueResumeTx(ctx
context
.
Context
, tx TTx, name
string
, opts *
QueuePauseOpts
)
error
QueueResumeTx resumes the queue with the given name. If the queue was
previously paused, any clients configured to work that queue will resume
fetching additional jobs. To resume all queues at once, use the special queue
name "*".
Clients with a configured notifier should receive a notification about the
resumed queue(s) within a few milliseconds of the transaction commit. Clients
in poll-only mode will resume after their next poll for queue configuration.
The provided context is used for the underlying Postgres update and can be
used to cancel the operation or apply a timeout. The opts are reserved for
future functionality.
func (*Client[TTx])
QueueUpdate
¶
added in
v0.20.0
func (c *
Client
[TTx]) QueueUpdate(ctx
context
.
Context
, name
string
, params *
QueueUpdateParams
) (*
rivertype
.
Queue
,
error
)
QueueUpdate updates a queue's settings in the database. These settings
override the settings in the client (if applied).
func (*Client[TTx])
QueueUpdateTx
¶
added in
v0.20.2
func (c *
Client
[TTx]) QueueUpdateTx(ctx
context
.
Context
, tx TTx, name
string
, params *
QueueUpdateParams
) (*
rivertype
.
Queue
,
error
)
QueueUpdateTx updates a queue's settings in the database. These settings
override the settings in the client (if applied).
func (*Client[TTx])
Queues
¶
added in
v0.10.0
func (c *
Client
[TTx]) Queues() *
QueueBundle
Queues returns the currently configured set of queues for the client, and can
be used to add new ones.
func (*Client[TTx])
Schema
¶
added in
v0.24.0
func (c *
Client
[TTx]) Schema()
string
Schema returns the configured schema for the client.
func (*Client[TTx])
Start
¶
func (c *
Client
[TTx]) Start(ctx
context
.
Context
)
error
Start starts the client's job fetching and working loops. Once this is called,
the client will run in a background goroutine until stopped. All jobs are
run with a context inheriting from the provided context, but with a timeout
deadline applied based on the job's settings.
A graceful shutdown stops fetching new jobs but allows any previously fetched
jobs to complete. This can be initiated with the Stop method.
A more abrupt shutdown can be achieved by either cancelling the provided
context or by calling StopAndCancel. This will not only stop fetching new
jobs, but will also cancel the context for any currently-running jobs. If
using StopAndCancel, there's no need to also call Stop.
func (*Client[TTx])
Stop
¶
func (c *
Client
[TTx]) Stop(ctx
context
.
Context
)
error
Stop performs a graceful shutdown of the Client. It signals all producers
to stop fetching new jobs and waits for any fetched or in-progress jobs to
complete before exiting. If the provided context is done before shutdown has
completed, Stop will return immediately with the context's error.
There's no need to call this method if a hard stop has already been initiated
by cancelling the context passed to Start or by calling StopAndCancel.
func (*Client[TTx])
StopAndCancel
¶
func (c *
Client
[TTx]) StopAndCancel(ctx
context
.
Context
)
error
StopAndCancel shuts down the client and cancels all work in progress. It is a
more aggressive stop than Stop because the contexts for any in-progress jobs
are cancelled. However, it still waits for jobs to complete before returning,
even though their contexts are cancelled. If the provided context is done
before shutdown has completed, Stop will return immediately with the
context's error.
This can also be initiated by cancelling the context passed to Run. There is
no need to call this method if the context passed to Run is cancelled
instead.
func (*Client[TTx])
Stopped
¶
added in
v0.0.11
func (c *
Client
[TTx]) Stopped() <-chan struct{}
Stopped returns a channel that will be closed when the Client has stopped.
It can be used to wait for a graceful shutdown to complete.
It is not affected by any contexts passed to Stop or StopAndCancel.
func (*Client[TTx])
Subscribe
¶
func (c *
Client
[TTx]) Subscribe(kinds ...
EventKind
) (<-chan *
Event
, func())
Subscribe subscribes to the provided kinds of events that occur within the
client, like EventKindJobCompleted for when a job completes.
Returns a channel over which to receive events along with a cancel function
that can be used to cancel and tear down resources associated with the
subscription. It's recommended but not necessary to invoke the cancel
function. Resources will be freed when the client stops in case it's not.
The event channel is buffered and sends on it are non-blocking. Consumers
must process events in a timely manner or it's possible for events to be
dropped. Any slow operations performed in a response to a receipt (e.g.
persisting to a database) should be made asynchronous to avoid event loss.
Callers must specify the kinds of events they're interested in. This allows
for forward compatibility in case new kinds of events are added in future
versions. If new event kinds are added, callers will have to explicitly add
them to their requested list and ensure they can be handled correctly.
func (*Client[TTx])
SubscribeConfig
¶
added in
v0.1.0
func (c *
Client
[TTx]) SubscribeConfig(config *
SubscribeConfig
) (<-chan *
Event
, func())
SubscribeConfig is a special internal variant of Subscribe that lets us
inject an overridden channel size.
type
ClientNotifyBundle
¶
added in
v0.29.0
type ClientNotifyBundle[TTx
any
] struct {
// contains filtered or unexported fields
}
ClientNotifyBundle sends various notifications for a client (in the sense of
Postgres listen/notify). Functions are on this bundle struct instead of the
top-level client to keep them grouped together and better organized.
func (*ClientNotifyBundle[TTx])
RequestResign
¶
added in
v0.29.0
func (c *
ClientNotifyBundle
[TTx]) RequestResign(ctx
context
.
Context
)
error
RequestResign sends a notification requesting that the current leader resign.
This usually causes the resignation of the current leader, but may have no
effect if no leader is currently elected.
func (*ClientNotifyBundle[TTx])
RequestResignTx
¶
added in
v0.29.0
func (c *
ClientNotifyBundle
[TTx]) RequestResignTx(ctx
context
.
Context
, tx TTx)
error
RequestResignTx sends a notification requesting that the current leader
resign. This usually causes the resignation of the current leader, but may
have no effect if no leader is currently elected.
This variant sends a notification in a transaction, which means that no
notification is sent until the transaction commits.
type
ClientRetryPolicy
¶
type ClientRetryPolicy interface {
// NextRetry calculates when the next retry for a failed job should take place
// given when it was last attempted and its number of attempts, or any other
// of the job's properties a user-configured retry policy might want to
// consider.
NextRetry(job *
rivertype
.
JobRow
)
time
.
Time
}
ClientRetryPolicy is an interface that can be implemented to provide a retry
policy for how River deals with failed jobs at the client level (when a
worker does not define an override for `NextRetry`). Jobs are scheduled to be
retried in the future up until they've reached the job's max attempts, at
which pointed they're set as discarded.
The ClientRetryPolicy does not have access to generics and operates on the
raw JobRow struct with encoded args.
type
Config
¶
type Config struct {
// AdvisoryLockPrefix is a configurable 32-bit prefix that River will use
// when generating any key to acquire a Postgres advisory lock. All advisory
// locks share the same 64-bit number space, so this allows a calling
// application to guarantee that a River advisory lock will never conflict
// with one of its own by cordoning each type to its own prefix.
//
// If this value isn't set, River defaults to generating key hashes across
// the entire 64-bit advisory lock number space, which is large enough that
// conflicts are exceedingly unlikely. If callers don't strictly need this
// option then it's recommended to leave it unset because the prefix leaves
// only 32 bits of number space for advisory lock hashes, so it makes
// internally conflicting River-generated keys more likely.
//
// Advisory locks are currently only used for the deprecated fallback/slow
// path of unique job insertion when pending, scheduled, available, or running
// are omitted from a customized ByState configuration.
AdvisoryLockPrefix
int32
// CancelledJobRetentionPeriod is the amount of time to keep cancelled jobs
// around before they're removed permanently.
//
// The special value -1 disables deletion of cancelled jobs.
//
// Defaults to 24 hours.
CancelledJobRetentionPeriod
time
.
Duration
// CompletedJobRetentionPeriod is the amount of time to keep completed jobs
// around before they're removed permanently.
//
// The special value -1 disables deletion of completed jobs.
//
// Defaults to 24 hours.
CompletedJobRetentionPeriod
time
.
Duration
// DiscardedJobRetentionPeriod is the amount of time to keep discarded jobs
// around before they're removed permanently.
//
// The special value -1 disables deletion of discarded jobs.
//
// Defaults to 7 days.
DiscardedJobRetentionPeriod
time
.
Duration
// ErrorHandler can be configured to be invoked in case of an error or panic
// occurring in a job. This is often useful for logging and exception
// tracking, but can also be used to customize retry behavior.
ErrorHandler
ErrorHandler
// FetchCooldown is the minimum amount of time to wait between fetches of new
// jobs. Jobs will only be fetched *at most* this often, but if no new jobs
// are coming in via LISTEN/NOTIFY then fetches may be delayed as long as
// FetchPollInterval.
//
// Throughput is limited by this value.
//
// Individual QueueConfig structs may override this for a specific queue.
//
// Defaults to 100 ms.
FetchCooldown
time
.
Duration
// FetchPollInterval is the amount of time between periodic fetches for new
// jobs. Typically new jobs will be picked up ~immediately after insert via
// LISTEN/NOTIFY, but this provides a fallback.
//
// Individual QueueConfig structs may override this for a specific queue.
//
// Defaults to 1 second.
FetchPollInterval
time
.
Duration
// ID is the unique identifier for this client. If not set, a random
// identifier will be generated.
//
// This is used to identify the client in job attempts and for leader election.
// This value must be unique across all clients in the same database and
// schema and there must not be more than one process running with the same
// ID at the same time.
//
// A client ID should differ between different programs and must be unique
// across all clients in the same database and schema. There must not be
// more than one process running with the same ID at the same time.
// Duplicate IDs between processes will lead to facilities like leader
// election or client statistics to fail in novel ways. However, the client
// ID is shared by all executors within any given client. (i.e.  different
// Go processes have different IDs, but IDs are shared within any given
// process.)
//
// If in doubt, leave this property empty.
ID
string
// JobCleanerTimeout is the timeout of the individual queries within the job
// cleaner.
//
// Defaults to 30 seconds, which should be more than enough time for most
// deployments.
JobCleanerTimeout
time
.
Duration
// JobInsertMiddleware are optional functions that can be called around job
// insertion.
//
// Deprecated: Prefer the use of Middleware instead (which may contain
// instances of rivertype.JobInsertMiddleware).
JobInsertMiddleware []
rivertype
.
JobInsertMiddleware
// JobTimeout is the maximum amount of time a job is allowed to run before its
// context is cancelled. A timeout of zero means JobTimeoutDefault will be
// used, whereas a value of -1 means the job's context will not be cancelled
// unless the Client is shutting down.
//
// Defaults to 1 minute.
JobTimeout
time
.
Duration
// Hooks are functions that may activate at certain points during a job's
// lifecycle (see rivertype.Hook), installed globally.
//
// The effect of hooks in this list will depend on the specific hook
// interfaces they implement, so for example implementing
// rivertype.HookInsertBegin will cause the hook to be invoked before a job
// is inserted, or implementing rivertype.HookWorkBegin will cause it to be
// invoked before a job is worked. Hook structs may implement multiple hook
// interfaces.
//
// Order in this list is significant. A hook that appears first will be
// entered before a hook that appears later. For any particular phase, order
// is relevant only for hooks that will run for that phase. For example, if
// two rivertype.HookInsertBegin are separated by a rivertype.HookWorkBegin,
// during job insertion those two outer hooks will run one after another,
// and the work hook between them will not run. When a job is worked, the
// work hook runs and the insertion hooks on either side of it are skipped.
//
// Jobs may have their own specific hooks by implementing JobArgsWithHooks.
Hooks []
rivertype
.
Hook
// Logger is the structured logger to use for logging purposes. If none is
// specified, logs will be emitted to STDOUT with messages at warn level
// or higher.
Logger *
slog
.
Logger
// MaxAttempts is the default number of times a job will be retried before
// being discarded. This value is applied to all jobs by default, and can be
// overridden on individual job types on the JobArgs or on a per-job basis at
// insertion time.
//
// If not specified, defaults to 25 (MaxAttemptsDefault).
MaxAttempts
int
// Middleware contains middleware that may activate at certain points during
// a job's lifecycle (see rivertype.Middleware), installed globally.
//
// The effect of middleware in this list will depend on the specific
// middleware interfaces they implement, so for example implementing
// rivertype.JobInsertMiddleware will cause the middleware to be invoked
// when jobs are inserted, and implementing rivertype.WorkerMiddleware will
// cause it to be invoked when a job is worked. Middleware structs may
// implement multiple middleware interfaces.
//
// Order in this list is significant. Middleware that appears first will be
// entered before middleware that appears later. For any particular phase,
// order is relevant only for middlewares that will run for that phase. For
// example, if two rivertype.JobInsertMiddleware are separated by a
// rivertype.WorkerMiddleware, during job insertion those two outer
// middlewares will run one after another, and the work middleware between
// them will not run. When a job is worked, the work middleware runs and the
// insertion middlewares on either side of it are skipped.
Middleware []
rivertype
.
Middleware
// PeriodicJobs are a set of periodic jobs to run at the specified intervals
// in the client.
PeriodicJobs []*
PeriodicJob
// PollOnly starts the client in "poll only" mode, which avoids issuing
// `LISTEN` statements to wait for events like a leadership resignation or
// new job available. The program instead polls periodically to look for
// changes (checking for new jobs on the period in FetchPollInterval).
//
// The downside of this mode of operation is that events will usually be
// noticed less quickly. A new job in the queue may have to wait up to
// FetchPollInterval to be locked for work. When a leader resigns, it will
// be up to five seconds before a new one elects itself.
//
// The upside is that it makes River compatible with systems where
// listen/notify isn't available. For example, PgBouncer in transaction
// pooling mode.
PollOnly
bool
// Queues is a list of queue names for this client to operate on along with
// configuration for the queue like the maximum number of workers to run for
// each queue.
//
// This field may be omitted for a program that's only queueing jobs rather
// than working them. If it's specified, then Workers must also be given.
Queues map[
string
]
QueueConfig
// ReindexerSchedule is the schedule for running the reindexer. If nil, the
// reindexer will run at midnight UTC every day.
ReindexerSchedule
PeriodicSchedule
// ReindexerTimeout is the amount of time to wait for the reindexer to run a
// single reindex operation before cancelling it via context. Set to -1 to
// disable the timeout.
//
// Defaults to 1 minute.
ReindexerTimeout
time
.
Duration
// RescueStuckJobsAfter is the amount of time a job can be running before it
// is considered stuck. A stuck job which has not yet reached its max attempts
// will be scheduled for a retry, while one which has exhausted its attempts
// will be discarded.  This prevents jobs from being stuck forever if a worker
// crashes or is killed.
//
// Note that this can result in repeat or duplicate execution of a job that is
// not actually stuck but is still working. The value should be set higher
// than the maximum duration you expect your jobs to run. Setting a value too
// low will result in more duplicate executions, whereas too high of a value
// will result in jobs being stuck for longer than necessary before they are
// retried.
//
// RescueStuckJobsAfter must be greater than JobTimeout. Otherwise, jobs
// would become eligible for rescue while they're still running.
//
// Defaults to 1 hour, or in cases where JobTimeout has been configured and
// is greater than 1 hour, JobTimeout + 1 hour.
RescueStuckJobsAfter
time
.
Duration
// RetryPolicy is a configurable retry policy for the client.
//
// Defaults to DefaultRetryPolicy.
RetryPolicy
ClientRetryPolicy
// Schema is a non-standard Schema where River tables are located. All table
// references in database queries will use this value as a prefix.
//
// Defaults to empty, which causes the client to look for tables using the
// setting of Postgres `search_path`.
Schema
string
// SkipJobKindValidation causes the job kind format validation check to be
// skipped. This is available as an interim stopgap for users that have
// invalid job kind names, but would rather disable the check rather than
// fix them immediately.
//
// Deprecated: This option will be removed in a future versions so that job
// kinds will always have to have a valid format.
SkipJobKindValidation
bool
// SkipUnknownJobCheck is a flag to control whether the client should skip
// checking to see if a registered worker exists in the client's worker bundle
// for a job arg prior to insertion.
//
// This can be set to true to allow a client to insert jobs which are
// intended to be worked by a different client which effectively makes
// the client's insertion behavior mimic that of an insert-only client.
//
// Defaults to false.
SkipUnknownJobCheck
bool
// Test holds configuration specific to test environments.
Test
TestConfig
// TestOnly can be set to true to disable certain features that are useful
// in production, but which may be harmful to tests, in ways like having the
// effect of making them slower. It should not be used outside of test
// suites.
//
// For example, queue maintenance services normally stagger their startup
// with a random jittered sleep so they don't all try to work at the same
// time. This is nice in production, but makes starting and stopping the
// client in a test case slower.
TestOnly
bool
// Workers is a bundle of registered job workers.
//
// This field may be omitted for a program that's only enqueueing jobs
// rather than working them, but if it is configured the client can validate
// ahead of time that a worker is properly registered for an inserted job.
// (i.e.  That it wasn't forgotten by accident.)
Workers *
Workers
// WorkerMiddleware are optional functions that can be called around
// all job executions.
//
// Deprecated: Prefer the use of Middleware instead (which may contain
// instances of rivertype.WorkerMiddleware).
WorkerMiddleware []
rivertype
.
WorkerMiddleware
// contains filtered or unexported fields
}
Config is the configuration for a Client.
Both Queues and Workers are required for a client to work jobs, but an
insert-only client can be initialized by omitting Queues, and not calling
Start for the client. Workers can also be omitted, but it's better to include
it so River can check that inserted job kinds have a worker that can run
them.
func (*Config)
WithDefaults
¶
added in
v0.17.0
func (c *
Config
) WithDefaults() *
Config
WithDefaults returns a copy of the Config with all default values applied.
type
DefaultClientRetryPolicy
¶
type DefaultClientRetryPolicy struct {
// contains filtered or unexported fields
}
DefaultClientRetryPolicy is River's default retry policy.
func (*DefaultClientRetryPolicy)
NextRetry
¶
func (p *
DefaultClientRetryPolicy
) NextRetry(job *
rivertype
.
JobRow
)
time
.
Time
NextRetry gets the next retry given for the given job, accounting for when it
was last attempted and what attempt number that was. Reschedules using a
basic exponential backoff of `ATTEMPT^4`, so after the first failure a new
try will be scheduled in 1 seconds, 16 seconds after the second, 1 minute and
21 seconds after the third, etc.
Snoozes do not count as attempts and do not influence retry behavior.
Earlier versions of River would allow the attempt to increment each time a
job was snoozed. Although this has been changed and snoozes now decrement the
attempt count, we can maintain the same retry schedule even for pre-existing
jobs by using the number of errors instead of the attempt count. This ensures
consistent behavior across River versions.
At degenerately high retry counts (>= 310) the policy starts adding the
equivalent of the maximum of time.Duration to each retry, about 292 years.
The schedule is no longer exponential past this point.
type
ErrorHandler
¶
type ErrorHandler interface {
// HandleError is invoked in case of an error occurring in a job.
//
// Context is descended from the one used to start the River client that
// worked the job. Errors are handled above all middleware, so changes made
// to context by a middleware are not available in the context.
HandleError(ctx
context
.
Context
, job *
rivertype
.
JobRow
, err
error
) *
ErrorHandlerResult
// HandlePanic is invoked in case of a panic occurring in a job.
//
// Context is descended from the one used to start the River client that
// worked the job. Panics are handled above all middleware, so changes made
// to context by a middleware are not available in the context (however,
// panics can be recovered from in any middleware where middleware context
// is available).
HandlePanic(ctx
context
.
Context
, job *
rivertype
.
JobRow
, panicVal
any
, trace
string
) *
ErrorHandlerResult
}
ErrorHandler provides an interface that will be invoked in case of an error
or panic occurring in the job. This is often useful for logging and exception
tracking, but can also be used to customize retry behavior.
type
ErrorHandlerResult
¶
type ErrorHandlerResult struct {
// SetCancelled can be set to true to fail the job immediately and
// permanently. By default it'll continue to follow the configured retry
// schedule.
SetCancelled
bool
}
type
Event
¶
type Event struct {
// Kind is the kind of event. Receivers should read this field and respond
// accordingly. Subscriptions will only receive event kinds that they
// requested when creating a subscription with Subscribe.
Kind
EventKind
// Job contains job-related information.
Job *
rivertype
.
JobRow
// JobStats are statistics about the run of a job.
JobStats *
JobStatistics
// Queue contains queue-related information.
Queue *
rivertype
.
Queue
}
Event wraps an event that occurred within a River client, like a job being
completed.
type
EventKind
¶
type EventKind
string
EventKind is a kind of event to subscribe to from a client.
const (
// EventKindJobCancelled occurs when a job is cancelled.
EventKindJobCancelled
EventKind
= "job_cancelled"
// EventKindJobCompleted occurs when a job is completed.
EventKindJobCompleted
EventKind
= "job_completed"
// EventKindJobFailed occurs when a job fails. Occurs both when a job fails
// and will be retried and when a job fails for the last time and will be
// discarded. Callers can use job fields like `Attempt` and `State` to
// differentiate each type of occurrence.
EventKindJobFailed
EventKind
= "job_failed"
// EventKindJobSnoozed occurs when a job is snoozed.
EventKindJobSnoozed
EventKind
= "job_snoozed"
// EventKindQueuePaused occurs when a queue is paused.
EventKindQueuePaused
EventKind
= "queue_paused"
// EventKindQueueResumed occurs when a queue is resumed.
EventKindQueueResumed
EventKind
= "queue_resumed"
)
type
HookDefaults
¶
added in
v0.19.0
type HookDefaults struct{}
HookDefaults should be embedded on any hooks implementation. It helps
identify a struct as hooks, and guarantee forward compatibility in case
additions are necessary to the rivertype.Hook interface.
func (*HookDefaults)
IsHook
¶
added in
v0.19.0
func (d *
HookDefaults
) IsHook()
bool
type
HookInsertBeginFunc
¶
added in
v0.19.0
type HookInsertBeginFunc func(ctx
context
.
Context
, params *
rivertype
.
JobInsertParams
)
error
HookInsertBeginFunc is a convenience helper for implementing
rivertype.HookInsertBegin using a simple function instead of a struct.
func (HookInsertBeginFunc)
InsertBegin
¶
added in
v0.19.0
func (f
HookInsertBeginFunc
) InsertBegin(ctx
context
.
Context
, params *
rivertype
.
JobInsertParams
)
error
func (HookInsertBeginFunc)
IsHook
¶
added in
v0.19.0
func (f
HookInsertBeginFunc
) IsHook()
bool
type
HookPeriodicJobsStartFunc
¶
added in
v0.29.0
type HookPeriodicJobsStartFunc func(ctx
context
.
Context
, params *
rivertype
.
HookPeriodicJobsStartParams
)
error
HookPeriodicJobsStartFunc is a convenience helper for implementing
rivertype.HookPeriodicJobsStart using a simple function instead of a struct.
func (HookPeriodicJobsStartFunc)
IsHook
¶
added in
v0.29.0
func (f
HookPeriodicJobsStartFunc
) IsHook()
bool
func (HookPeriodicJobsStartFunc)
Start
¶
added in
v0.29.0
func (f
HookPeriodicJobsStartFunc
) Start(ctx
context
.
Context
, params *
rivertype
.
HookPeriodicJobsStartParams
)
error
type
HookWorkBeginFunc
¶
added in
v0.19.0
type HookWorkBeginFunc func(ctx
context
.
Context
, job *
rivertype
.
JobRow
)
error
HookWorkBeginFunc is a convenience helper for implementing
rivertype.HookWorkBegin using a simple function instead of a struct.
func (HookWorkBeginFunc)
IsHook
¶
added in
v0.19.0
func (f
HookWorkBeginFunc
) IsHook()
bool
func (HookWorkBeginFunc)
WorkBegin
¶
added in
v0.19.0
func (f
HookWorkBeginFunc
) WorkBegin(ctx
context
.
Context
, job *
rivertype
.
JobRow
)
error
type
HookWorkEndFunc
¶
added in
v0.21.0
type HookWorkEndFunc func(ctx
context
.
Context
, job *
rivertype
.
JobRow
, err
error
)
error
HookWorkEndFunc is a convenience helper for implementing
rivertype.HookWorkEnd using a simple function instead of a struct.
func (HookWorkEndFunc)
IsHook
¶
added in
v0.21.0
func (f
HookWorkEndFunc
) IsHook()
bool
func (HookWorkEndFunc)
WorkEnd
¶
added in
v0.21.0
func (f
HookWorkEndFunc
) WorkEnd(ctx
context
.
Context
, job *
rivertype
.
JobRow
, err
error
)
error
type
InsertManyParams
¶
type InsertManyParams struct {
// Args are the arguments of the job to insert.
Args
JobArgs
// InsertOpts are insertion options for this job.
InsertOpts *
InsertOpts
}
InsertManyParams encapsulates a single job combined with insert options for
use with batch insertion.
type
InsertOpts
¶
type InsertOpts struct {
// MaxAttempts is the maximum number of total attempts (including both the
// original run and all retries) before a job is abandoned and set as
// discarded.
MaxAttempts
int
// Metadata is a JSON object blob of arbitrary data that will be stored with
// the job. Users should not overwrite or remove anything stored in this
// field by River.
Metadata []
byte
// Pending indicates that the job should be inserted in the `pending` state.
// Pending jobs are not immediately available to be worked and are never
// deleted, but they can be used to indicate work which should be performed in
// the future once they are made available (or scheduled) by some external
// update.
Pending
bool
// Priority is the priority of the job, with 1 being the highest priority and
// 4 being the lowest. When fetching available jobs to work, the highest
// priority jobs will always be fetched before any lower priority jobs are
// fetched. Note that if your workers are swamped with more high-priority jobs
// then they can handle, lower priority jobs may not be fetched.
//
// Defaults to PriorityDefault.
Priority
int
// Queue is the name of the job queue in which to insert the job.
//
// Defaults to the job kind's default queue if set via
// `JobArgsWithInsertOpts`, or QueueDefault if not.
Queue
string
// ScheduledAt is a time in future at which to schedule the job (i.e. in
// cases where it shouldn't be run immediately). The job is guaranteed not
// to run before this time, but may run slightly after depending on the
// number of other scheduled jobs and how busy the queue is.
//
// Use of this option generally only makes sense when passing options into
// Insert rather than when a job args struct is implementing
// JobArgsWithInsertOpts, however, it will work in both cases.
ScheduledAt
time
.
Time
// Tags are an arbitrary list of keywords to add to the job. They have no
// functional behavior and are meant entirely as a user-specified construct
// to help group and categorize jobs.
//
// Tags should conform to the regex `\A[\w][\w\-]+[\w]\z` and be a maximum
// of 255 characters long. No special characters are allowed.
//
// If tags are specified from both a job args override and from options on
// Insert, the latter takes precedence. Tags are not merged.
Tags []
string
// UniqueOpts returns options relating to job uniqueness. An empty struct
// avoids setting any worker-level unique options.
UniqueOpts
UniqueOpts
}
InsertOpts are optional settings for a new job which can be provided at job
insertion time. These will override any default InsertOpts settings provided
by JobArgsWithInsertOpts, as well as any global defaults.
type
Job
¶
type Job[T
JobArgs
] struct {
*
rivertype
.
JobRow
// Args are the arguments for the job.
Args T
}
Job represents a single unit of work, holding both the arguments and
information for a job with args of type T.
func
JobCompleteTx
¶
func JobCompleteTx[TDriver
riverdriver
.
Driver
[TTx], TTx
any
, TArgs
JobArgs
](ctx
context
.
Context
, tx TTx, job *
Job
[TArgs]) (*
Job
[TArgs],
error
)
JobCompleteTx marks the job as completed as part of transaction tx. If tx is
rolled back, the completion will be as well.
The function needs to know the type of the River database driver, which is
the same as the one in use by Client, but the other generic parameters can be
inferred. An invocation should generally look like:
_, err := river.JobCompleteTx[*riverpgxv5.Driver](ctx, tx, job)
if err != nil {
// handle error
}
Returns the updated, completed job.
type
JobArgs
¶
type JobArgs interface {
// Kind is a string that uniquely identifies the type of job. This must be
// provided on your job arguments struct. Jobs are identified by a string
// instead of being based on type names so that previously inserted jobs
// can be worked across deploys even if job/worker types are renamed.
//
// Kinds should be formatted without spaces like `my_custom_job`,
// `mycustomjob`, or `my-custom-job`. Many special characters like colons,
// dots, hyphens, and underscores are allowed, but those like spaces and
// commas, which would interfere with UI functionality, are invalid.
//
// After initially deploying a job, it's generally not safe to rename its
// kind (unless the database is completely empty) because River won't know
// which worker should work the old kind. Job kinds can be renamed safely
// over multiple deploys using the JobArgsWithKindAliases interface.
Kind()
string
}
JobArgs is an interface that represents the arguments for a job of type T.
These arguments are serialized into JSON and stored in the database.
The struct is serialized using `encoding/json`. All exported fields are
serialized, unless skipped with a struct field tag.
type
JobArgsWithHooks
¶
added in
v0.19.0
type JobArgsWithHooks interface {
// Hooks returns specific hooks to run for this job type. These will run
// after the global hooks configured on the client.
//
// Warning: Hooks returned should be based on the job type only and be
// invariant of the specific contents of a job. Hooks are extracted by
// instantiating a generic instance of the job even when a specific instance
// is available, so any conditional logic within will be ignored. This is
// done because although specific job information may be available in some
// hook contexts like on InsertBegin, it won't be in others like WorkBegin.
Hooks() []
rivertype
.
Hook
}
JobArgsWithHooks is an interface that job args can implement to attach
specific hooks (i.e. other than those globally installed to a client) to
certain kinds of jobs.
type
JobArgsWithInsertOpts
¶
type JobArgsWithInsertOpts interface {
// InsertOpts returns options for all jobs of this job type, overriding any
// system defaults. These can also be overridden at insertion time.
InsertOpts()
InsertOpts
}
JobArgsWithInsertOpts is an extra interface that a job may implement on top
of JobArgs to provide insertion-time options for all jobs of this type.
type
JobArgsWithKindAliases
¶
added in
v0.22.0
type JobArgsWithKindAliases interface {
// KindAliases returns alias kinds that an associated job args worker will
// respond to.
KindAliases() []
string
}
JobArgsWithKindAliases  is an interface that jobs args can implement to
provide an alternate kind which a worker will be registered under in addition
to the primary kind. This is useful for renaming a job kind in a safe manner
so that any jobs already in the database aren't orphaned.
Renaming a job is a three part process. To begin, a job args with its
original name:
type jobArgsBeingRenamed struct{}

func (a jobArgsBeingRenamed) Kind() string { return "old_name" }
Rename by putting the new name in Kind and moving the old name to
KindAliases:
type jobArgsBeingRenamed struct{}

func (a jobArgsBeingRenamed) Kind() string          { return "new_name" }
func (a jobArgsBeingRenamed) KindAliases() []string { return []string{"old_name"} }
After all jobs inserted under the original name have finished working
(including all their possible retries, which notably might take up to three
weeks on the default retry policy), remove KindAliases:
type jobArgsBeingRenamed struct{}

func (a jobArgsBeingRenamed) Kind() string { return "new_name" }
type
JobCancelError
¶
added in
v0.14.0
type JobCancelError =
rivertype
.
JobCancelError
JobCancelError is the error type returned by JobCancel. It should not be
initialized directly, but is returned from the
JobCancel
function and can
be used for test assertions.
type
JobDeleteManyParams
¶
added in
v0.24.0
type JobDeleteManyParams struct {
// contains filtered or unexported fields
}
JobDeleteManyParams specifies the parameters for a JobDeleteMany query. It
must be initialized with NewJobDeleteManyParams. Params can be built by
chaining methods on the JobDeleteManyParams object:
params := NewJobDeleteManyParams().First(100).States(river.JobStateCompleted)
func
NewJobDeleteManyParams
¶
added in
v0.24.0
func NewJobDeleteManyParams() *
JobDeleteManyParams
NewJobDeleteManyParams creates a new JobDeleteManyParams to delete jobs
sorted by ID in ascending order, deleting 100 jobs at most.
func (*JobDeleteManyParams)
First
¶
added in
v0.24.0
func (p *
JobDeleteManyParams
) First(count
int
) *
JobDeleteManyParams
First returns an updated filter set that will only delete the first
count jobs.
Count must be between 1 and 10_000, inclusive, or this will panic.
func (*JobDeleteManyParams)
IDs
¶
added in
v0.24.0
func (p *
JobDeleteManyParams
) IDs(ids ...
int64
) *
JobDeleteManyParams
IDs returns an updated filter set that will only delete jobs with the given
IDs.
func (*JobDeleteManyParams)
Kinds
¶
added in
v0.24.0
func (p *
JobDeleteManyParams
) Kinds(kinds ...
string
) *
JobDeleteManyParams
Kinds returns an updated filter set that will only delete jobs of the given
kinds.
func (*JobDeleteManyParams)
Priorities
¶
added in
v0.24.0
func (p *
JobDeleteManyParams
) Priorities(priorities ...
int16
) *
JobDeleteManyParams
Priorities returns an updated filter set that will only delete jobs with the
given priorities.
func (*JobDeleteManyParams)
Queues
¶
added in
v0.24.0
func (p *
JobDeleteManyParams
) Queues(queues ...
string
) *
JobDeleteManyParams
Queues returns an updated filter set that will only delete jobs from the
given queues.
func (*JobDeleteManyParams)
States
¶
added in
v0.24.0
func (p *
JobDeleteManyParams
) States(states ...
rivertype
.
JobState
) *
JobDeleteManyParams
States returns an updated filter set that will only delete jobs in the given
states.
func (*JobDeleteManyParams)
UnsafeAll
¶
added in
v0.25.0
func (p *
JobDeleteManyParams
) UnsafeAll() *
JobDeleteManyParams
UnsafeAll is a special directive that allows unbounded job deletion without
any filters. Normally, filters like IDs or Kinds is required to scope down
the deletion so that the caller doesn't accidentally delete all non-running
jobs. Invoking UnsafeAll removes this safety guard so that all jobs can be
removed arbitrarily.
Example of use:
deleteRes, err = client.JobDeleteMany(ctx, NewJobDeleteManyParams().UnsafeAll())
if err != nil {
// handle error
}
It only makes sense to call this function if no filters have yet been applied
on the parameters object. If some have already, calling it will panic.
type
JobDeleteManyResult
¶
added in
v0.24.0
type JobDeleteManyResult struct {
// Jobs is a slice of job returned as part of the list operation.
Jobs []*
rivertype
.
JobRow
}
JobDeleteManyResult is the result of a job list operation. It contains a list of
jobs and a cursor for fetching the next page of results.
type
JobInsertMiddlewareDefaults
deprecated
added in
v0.13.0
type JobInsertMiddlewareDefaults struct{
MiddlewareDefaults
}
JobInsertMiddlewareDefaults is an embeddable struct that provides default
implementations for the rivertype.JobInsertMiddleware. Use of this struct is
recommended in case rivertype.JobInsertMiddleware is expanded in the future
so that existing code isn't unexpectedly broken during an upgrade.
Deprecated: Prefer embedding the more general MiddlewareDefaults instead.
func (*JobInsertMiddlewareDefaults)
InsertMany
¶
added in
v0.13.0
func (d *
JobInsertMiddlewareDefaults
) InsertMany(ctx
context
.
Context
, manyParams []*
rivertype
.
JobInsertParams
, doInner func(ctx
context
.
Context
) ([]*
rivertype
.
JobInsertResult
,
error
)) ([]*
rivertype
.
JobInsertResult
,
error
)
type
JobInsertMiddlewareFunc
¶
added in
v0.21.0
type JobInsertMiddlewareFunc func(ctx
context
.
Context
, manyParams []*
rivertype
.
JobInsertParams
, doInner func(ctx
context
.
Context
) ([]*
rivertype
.
JobInsertResult
,
error
)) ([]*
rivertype
.
JobInsertResult
,
error
)
JobInsertMiddlewareFunc is a convenience helper for implementing
rivertype.JobInsertMiddleware using a simple function instead of a struct.
func (JobInsertMiddlewareFunc)
InsertMany
¶
added in
v0.21.0
func (f
JobInsertMiddlewareFunc
) InsertMany(ctx
context
.
Context
, manyParams []*
rivertype
.
JobInsertParams
, doInner func(ctx
context
.
Context
) ([]*
rivertype
.
JobInsertResult
,
error
)) ([]*
rivertype
.
JobInsertResult
,
error
)
func (JobInsertMiddlewareFunc)
IsMiddleware
¶
added in
v0.21.0
func (f
JobInsertMiddlewareFunc
) IsMiddleware()
bool
type
JobListCursor
¶
added in
v0.0.17
type JobListCursor struct {
// contains filtered or unexported fields
}
JobListCursor is used to specify a starting point for a paginated
job list query.
func
JobListCursorFromJob
¶
added in
v0.0.17
func JobListCursorFromJob(job *
rivertype
.
JobRow
) *
JobListCursor
JobListCursorFromJob creates a JobListCursor from a JobRow.
func (JobListCursor)
MarshalText
¶
added in
v0.0.17
func (c
JobListCursor
) MarshalText() ([]
byte
,
error
)
MarshalText implements encoding.TextMarshaler to encode the cursor as an
opaque string.
func (*JobListCursor)
UnmarshalText
¶
added in
v0.0.17
func (c *
JobListCursor
) UnmarshalText(text []
byte
)
error
UnmarshalText implements encoding.TextUnmarshaler to decode the cursor from
a previously marshaled string.
type
JobListOrderByField
¶
added in
v0.0.17
type JobListOrderByField
string
JobListOrderByField specifies the field to sort by.
const (
// JobListOrderByID specifies that the sort should be by job ID.
JobListOrderByID
JobListOrderByField
= "id"
// JobListOrderByFinalizedAt specifies that the sort should be by
// `finalized_at`.
//
// This option must be used in conjunction with filtering by only finalized
// job states.
JobListOrderByFinalizedAt
JobListOrderByField
= "finalized_at"
// JobListOrderByScheduledAt specifies that the sort should be by
// `scheduled_at`.
JobListOrderByScheduledAt
JobListOrderByField
= "scheduled_at"
// JobListOrderByTime specifies that the sort should be by the "best fit"
// time field based on listed state. The best fit is determined by looking
// at the first value given to JobListParams.States. If multiple states are
// specified, the ones after the first will be ignored.
//
// The specific time field used for sorting depends on requested state:
//
// * States `available`, `retryable`, or `scheduled` use `scheduled_at`.
// * State `running` uses `attempted_at`.
// * States `cancelled`, `completed`, or `discarded` use `finalized_at`.
JobListOrderByTime
JobListOrderByField
= "time"
)
type
JobListParams
¶
added in
v0.0.17
type JobListParams struct {
// contains filtered or unexported fields
}
JobListParams specifies the parameters for a JobList query. It must be
initialized with NewJobListParams. Params can be built by chaining methods on
the JobListParams object:
params := NewJobListParams().OrderBy(JobListOrderByTime, SortOrderAsc).First(100)
func
NewJobListParams
¶
added in
v0.0.17
func NewJobListParams() *
JobListParams
NewJobListParams creates a new JobListParams to return available jobs sorted
by time in ascending order, returning 100 jobs at most.
func (*JobListParams)
After
¶
added in
v0.0.17
func (p *
JobListParams
) After(cursor *
JobListCursor
) *
JobListParams
After returns an updated filter set that will only return jobs
after the given cursor.
func (*JobListParams)
First
¶
added in
v0.0.17
func (p *
JobListParams
) First(count
int
) *
JobListParams
First returns an updated filter set that will only return the first
count jobs.
Count must be between 1 and 10_000, inclusive, or this will panic.
func (*JobListParams)
IDs
¶
added in
v0.21.0
func (p *
JobListParams
) IDs(ids ...
int64
) *
JobListParams
IDs returns an updated filter set that will only return jobs with the given
IDs.
func (*JobListParams)
Kinds
¶
added in
v0.0.23
func (p *
JobListParams
) Kinds(kinds ...
string
) *
JobListParams
Kinds returns an updated filter set that will only return jobs of the given
kinds.
func (*JobListParams)
Metadata
¶
added in
v0.0.17
func (p *
JobListParams
) Metadata(json
string
) *
JobListParams
Metadata returns an updated filter set that will return only jobs that has
metadata which contains the given JSON fragment at its top level. This is
equivalent to the `@>` operator in Postgres:
https://www.postgresql.org/docs/current/functions-json.html
This function isn't supported in SQLite due to SQLite not having an
equivalent operator to use, so there's no efficient way to implement it. We
recommend the use of Where using a condition with a comparison on the `->>`
operator instead.
func (*JobListParams)
OrderBy
¶
added in
v0.0.17
func (p *
JobListParams
) OrderBy(field
JobListOrderByField
, direction
SortOrder
) *
JobListParams
OrderBy returns an updated filter set that will sort the results using the
specified field and direction.
If ordering by FinalizedAt, the States filter will be set to only include
finalized job states unless it has already been overridden.
func (*JobListParams)
Priorities
¶
added in
v0.21.0
func (p *
JobListParams
) Priorities(priorities ...
int16
) *
JobListParams
Priorities returns an updated filter set that will only return jobs with the
given priorities.
func (*JobListParams)
Queues
¶
added in
v0.0.17
func (p *
JobListParams
) Queues(queues ...
string
) *
JobListParams
Queues returns an updated filter set that will only return jobs from the
given queues.
func (*JobListParams)
States
¶
added in
v0.4.0
func (p *
JobListParams
) States(states ...
rivertype
.
JobState
) *
JobListParams
States returns an updated filter set that will only return jobs in the given
states.
func (*JobListParams)
Where
¶
added in
v0.23.0
func (p *
JobListParams
) Where(sql
string
, namedArgsMany ...
NamedArgs
) *
JobListParams
Where is an all-encompassing query escape hatch that adds an arbitrary
predicate after a list query's `WHERE ...` clause. Use of other JobListParams
filters should be preferred where possible because they're safer and their
compatibility between drivers is better guaranteed, but in case none is
suitable, Where can be used as a last resort.
For example, using Where to query with `jsonb_path_query_first(...)` using a
JSON path, a function that's specific to Postgres:
listParams = listParams.Where("jsonb_path_query_first(metadata, @json_path) = @json_val", NamedArgs{"json_path": "$.foo", "json_val": `"bar"`})
A JSON path can be used in a query in SQLite as well, but there the `->` or
`->>` operators must be used instead:
listParams = listParams.Where("metadata ->> @json_path = @json_val", NamedArgs{"json_path": "$.foo", "json_val": "bar"})
Arguments beyond the first are interpreted as named parameters. Each one
should be present in the query SQL prefixed with a `@` symbol. Multiple sets
of named parameters will be merged together, with values in later sets
overwriting those in earlier ones.
Calling Where multiple times will add multiple conditions separate by `AND`.
Use `OR` instead by stuffing all conditions into a single Where invocation.
Consider use of this function possibly hazardous! Any time raw SQL is in
play, an application is opening itself up to SQL injection attacks. Never mix
unsanitized user input into a SQL string, and use named parameters to curb
the likelihood of injection.
type
JobListResult
¶
added in
v0.4.0
type JobListResult struct {
// Jobs is a slice of job returned as part of the list operation.
Jobs []*
rivertype
.
JobRow
// LastCursor is a cursor that can be used to list the next page of jobs.
LastCursor *
JobListCursor
}
JobListResult is the result of a job list operation. It contains a list of
jobs and a cursor for fetching the next page of results.
type
JobSnoozeError
¶
added in
v0.14.0
type JobSnoozeError =
rivertype
.
JobSnoozeError
JobSnoozeError is the error type returned by JobSnooze. It should not be
initialized directly, but is returned from the
JobSnooze
function and can
be used for test assertions.
type
JobStatistics
¶
type JobStatistics struct {
CompleteDuration
time
.
Duration
// Time it took to set the job completed, discarded, or errored.
QueueWaitDuration
time
.
Duration
// Time the job spent waiting in available state before starting execution.
RunDuration
time
.
Duration
// Time job spent running (measured around job worker.)
}
JobStatistics contains information about a single execution of a job.
type
JobUpdateParams
¶
added in
v0.29.0
type JobUpdateParams struct {
// Output is a new output value for a job.
//
// If not set, and a job is updated from inside a work function, the job's
// output is set based on output recorded so far using RecordOutput.
Output
any
}
JobUpdateParams contains parameters for Client.JobUpdate and Client.JobUpdateTx.
type
MiddlewareDefaults
¶
added in
v0.19.0
type MiddlewareDefaults struct{}
MiddlewareDefaults should be embedded on any middleware implementation. It
helps identify a struct as middleware, and guarantees forward compatibility in
case additions are necessary to the rivertype.Middleware interface.
func (*MiddlewareDefaults)
IsMiddleware
¶
added in
v0.19.0
func (d *
MiddlewareDefaults
) IsMiddleware()
bool
type
NamedArgs
¶
added in
v0.23.0
type NamedArgs map[
string
]
any
NamedArgs are named arguments for use with JobListParams.Where. Keys should
look like "my_param", and map to parameters like "@my_param" in SQL queries.
"@" are present in the SQL, but not in the keys of this map.
type
PeriodicJob
¶
type PeriodicJob struct {
// contains filtered or unexported fields
}
PeriodicJob is a configuration for a periodic job.
func
NewPeriodicJob
¶
func NewPeriodicJob(scheduleFunc
PeriodicSchedule
, constructorFunc
PeriodicJobConstructor
, opts *
PeriodicJobOpts
) *
PeriodicJob
NewPeriodicJob returns a new PeriodicJob given a schedule and a constructor
function.
The schedule returns a time until the next time the periodic job should run.
The helper PeriodicInterval is available for jobs that should run on simple,
fixed intervals (e.g. every 15 minutes), and a custom schedule or third party
cron package can be used for more complex scheduling (see the cron example).
The constructor function is invoked each time a periodic job's schedule
elapses, returning job arguments to insert along with optional insertion
options.
The periodic job scheduler is approximate and doesn't guarantee strong
durability. It's started by the elected leader in a River cluster, and each
periodic job is assigned an initial run time when that occurs. New run times
are scheduled each time a job's target run time is reached and a new job
inserted. However, each scheduler only retains in-memory state, so anytime a
process quits or a new leader is elected, the whole process starts over
without regard for the state of the last scheduler. The RunOnStart option
can be used as a hedge to make sure that jobs with long run durations are
guaranteed to occasionally run.
type
PeriodicJobBundle
¶
added in
v0.2.0
type PeriodicJobBundle struct {
// contains filtered or unexported fields
}
PeriodicJobBundle is a bundle of currently configured periodic jobs. It's
made accessible through Client, where new periodic jobs can be configured,
and old ones removed.
func (*PeriodicJobBundle)
Add
¶
added in
v0.2.0
func (b *
PeriodicJobBundle
) Add(periodicJob *
PeriodicJob
)
rivertype
.
PeriodicJobHandle
Add adds a new periodic job to the client. The job is queued immediately if
RunOnStart is enabled, and then scheduled normally.
Returns a periodic job handle which can be used to subsequently remove the
job if desired.
Adding or removing periodic jobs has no effect unless this client is elected
leader because only the leader enqueues periodic jobs. To make sure that a
new periodic job is fully enabled or disabled, it should be added or removed
from _every_ active River client across all processes.
func (*PeriodicJobBundle)
AddMany
¶
added in
v0.2.0
func (b *
PeriodicJobBundle
) AddMany(periodicJobs []*
PeriodicJob
) []
rivertype
.
PeriodicJobHandle
AddMany adds many new periodic jobs to the client. The jobs are queued
immediately if their RunOnStart is enabled, and then scheduled normally.
Returns a periodic job handle which can be used to subsequently remove the
job if desired.
Adding or removing periodic jobs has no effect unless this client is elected
leader because only the leader enqueues periodic jobs. To make sure that a
new periodic job is fully enabled or disabled, it should be added or removed
from _every_ active River client across all processes.
func (*PeriodicJobBundle)
AddManySafely
¶
added in
v0.23.0
func (b *
PeriodicJobBundle
) AddManySafely(periodicJobs []*
PeriodicJob
) ([]
rivertype
.
PeriodicJobHandle
,
error
)
AddManySafely is the same as AddMany, but it returns an error in the case of
a validation problem or duplicate ID instead of panicking.
func (*PeriodicJobBundle)
AddSafely
¶
added in
v0.23.0
func (b *
PeriodicJobBundle
) AddSafely(periodicJob *
PeriodicJob
) (
rivertype
.
PeriodicJobHandle
,
error
)
AddSafely is the same as Add, but it returns an error in the case of a
validation problem or duplicate ID instead of panicking.
func (*PeriodicJobBundle)
Clear
¶
added in
v0.2.0
func (b *
PeriodicJobBundle
) Clear()
Clear clears all periodic jobs, cancelling all scheduled runs.
Adding or removing periodic jobs has no effect unless this client is elected
leader because only the leader enqueues periodic jobs. To make sure that a
new periodic job is fully enabled or disabled, it should be added or removed
from _every_ active River client across all processes.
func (*PeriodicJobBundle)
Remove
¶
added in
v0.2.0
func (b *
PeriodicJobBundle
) Remove(periodicJobHandle
rivertype
.
PeriodicJobHandle
)
Remove removes a periodic job, cancelling all scheduled runs.
Requires the use of the periodic job handle that was returned when the job
was added.
Adding or removing periodic jobs has no effect unless this client is elected
leader because only the leader enqueues periodic jobs. To make sure that a
new periodic job is fully enabled or disabled, it should be added or removed
from _every_ active River client across all processes.
func (*PeriodicJobBundle)
RemoveByID
¶
added in
v0.27.0
func (b *
PeriodicJobBundle
) RemoveByID(id
string
)
bool
RemoveByID removes a periodic job by ID, cancelling all scheduled runs.
Adding or removing periodic jobs has no effect unless this client is elected
leader because only the leader enqueues periodic jobs. To make sure that a
new periodic job is fully enabled or disabled, it should be added or removed
from _every_ active River client across all processes.
Has no effect if no jobs with the given ID is configured.
Returns true if a job with the given ID existed (and was removed), and false
otherwise.
func (*PeriodicJobBundle)
RemoveMany
¶
added in
v0.2.0
func (b *
PeriodicJobBundle
) RemoveMany(periodicJobHandles []
rivertype
.
PeriodicJobHandle
)
RemoveMany removes many periodic jobs, cancelling all scheduled runs.
Requires the use of the periodic job handles that were returned when the jobs
were added.
Adding or removing periodic jobs has no effect unless this client is elected
leader because only the leader enqueues periodic jobs. To make sure that a
new periodic job is fully enabled or disabled, it should be added or removed
from _every_ active River client across all processes.
func (*PeriodicJobBundle)
RemoveManyByID
¶
added in
v0.27.0
func (b *
PeriodicJobBundle
) RemoveManyByID(ids []
string
)
RemoveManyByID removes many periodic jobs by ID, cancelling all scheduled
runs.
Adding or removing periodic jobs has no effect unless this client is elected
leader because only the leader enqueues periodic jobs. To make sure that a
new periodic job is fully enabled or disabled, it should be added or removed
from _every_ active River client across all processes.
Has no effect if no jobs with the given IDs are configured.
type
PeriodicJobConstructor
¶
type PeriodicJobConstructor func() (
JobArgs
, *
InsertOpts
)
PeriodicJobConstructor is a function that gets called each time the paired
PeriodicSchedule is triggered.
A constructor must never block. It may return nil to indicate that no job
should be inserted.
type
PeriodicJobOpts
¶
type PeriodicJobOpts struct {
// ID is an optional identifier for the job. Identifiers must be unique
// between all periodic jobs and adding a periodic job will error if they're
// not.
ID
string
// RunOnStart can be used to indicate that a periodic job should insert an
// initial job as a new scheduler is started. This can be used as a hedge
// for jobs with longer scheduled durations that may not get to expiry
// before a new scheduler is elected.
//
// RunOnStart also applies when a new periodic job is added dynamically with
// `PeriodicJobs().Add` or `PeriodicJobs().AddMany`. Jobs added this way
// with RunOnStart set to true are inserted once, then continue with their
// normal run schedule.
RunOnStart
bool
}
PeriodicJobOpts are options for a periodic job.
type
PeriodicSchedule
¶
type PeriodicSchedule interface {
// Next returns the next time at which the job should be run given the
// current time.
Next(current
time
.
Time
)
time
.
Time
}
PeriodicSchedule is a schedule for a periodic job. Periodic jobs should
generally have an interval of at least 1 minute, and never less than one
second.
func
NeverSchedule
¶
added in
v0.16.0
func NeverSchedule()
PeriodicSchedule
NeverSchedule returns a PeriodicSchedule that never runs.
func
PeriodicInterval
¶
func PeriodicInterval(interval
time
.
Duration
)
PeriodicSchedule
PeriodicInterval returns a simple PeriodicSchedule that runs at the given
interval.
type
QueueAlreadyAddedError
¶
added in
v0.23.0
type QueueAlreadyAddedError struct {
Name
string
}
QueueAlreadyAddedError is returned when attempting to add a queue that has
already been added to the Client.
func (*QueueAlreadyAddedError)
Error
¶
added in
v0.23.0
func (e *
QueueAlreadyAddedError
) Error()
string
func (*QueueAlreadyAddedError)
Is
¶
added in
v0.23.0
func (e *
QueueAlreadyAddedError
) Is(target
error
)
bool
type
QueueBundle
¶
added in
v0.10.0
type QueueBundle struct {
// contains filtered or unexported fields
}
QueueBundle is a bundle for adding additional queues. It's made accessible
through Client.Queues.
func (*QueueBundle)
Add
¶
added in
v0.10.0
func (b *
QueueBundle
) Add(queueName
string
, queueConfig
QueueConfig
)
error
Add adds a new queue to the client. If the client is already started, a
producer for the queue is started. Context is inherited from the one given to
Client.Start.
type
QueueConfig
¶
type QueueConfig struct {
// FetchCooldown is the minimum amount of time to wait between fetches of new
// jobs. Jobs will only be fetched *at most* this often, but if no new jobs
// are coming in via LISTEN/NOTIFY then fetches may be delayed as long as
// FetchPollInterval.
//
// Throughput is limited by this value.
//
// If non-zero, this overrides the FetchCooldown setting in the Client's
// Config.
FetchCooldown
time
.
Duration
// FetchPollInterval is the amount of time between periodic fetches for new
// jobs. Typically new jobs will be picked up ~immediately after insert via
// LISTEN/NOTIFY, but this provides a fallback.
//
// If non-zero, this overrides the FetchCooldown setting in the Client's
// Config.
FetchPollInterval
time
.
Duration
// MaxWorkers is the maximum number of workers to run for the queue, or put
// otherwise, the maximum parallelism to run.
//
// This is the maximum number of workers within this particular client
// instance, but note that it doesn't control the total number of workers
// across parallel processes. Installations will want to calculate their
// total number by multiplying this number by the number of parallel nodes
// running River clients configured to the same database and queue.
//
// Requires a minimum of 1, and a maximum of 10,000.
MaxWorkers
int
}
QueueConfig contains queue-specific configuration.
type
QueueListParams
¶
added in
v0.5.0
type QueueListParams struct {
// contains filtered or unexported fields
}
QueueListParams specifies the parameters for a QueueList query. It must be
initialized with NewQueueListParams. Params can be built by chaining methods
on the QueueListParams object:
params := NewQueueListParams().First(100)
func
NewQueueListParams
¶
added in
v0.5.0
func NewQueueListParams() *
QueueListParams
NewQueueListParams creates a new QueueListParams to return available queues
sorted by time in ascending order, returning 100 jobs at most.
func (*QueueListParams)
First
¶
added in
v0.5.0
func (p *
QueueListParams
) First(count
int
) *
QueueListParams
First returns an updated filter set that will only return the first count
queues.
Count must be between 1 and 10000, inclusive, or this will panic.
type
QueueListResult
¶
added in
v0.5.0
type QueueListResult struct {
// Queues is a slice of queues returned as part of the list operation.
Queues []*
rivertype
.
Queue
}
QueueListResult is the result of a job list operation. It contains a list of
jobs and leaves room for future cursor functionality.
type
QueuePauseOpts
¶
added in
v0.5.0
type QueuePauseOpts struct{}
QueuePauseOpts are optional settings for pausing or resuming a queue.
type
QueueUpdateParams
¶
added in
v0.20.0
type QueueUpdateParams struct {
// Metadata is the new metadata for the queue. If nil or empty, the metadata
// will not be changed.
Metadata []
byte
}
QueueUpdateParams are the parameters for a QueueUpdate operation.
type
SortOrder
¶
added in
v0.0.17
type SortOrder
int
SortOrder specifies the direction of a sort.
const (
// SortOrderAsc specifies that the sort should in ascending order.
SortOrderAsc
SortOrder
=
iota
// SortOrderDesc specifies that the sort should in descending order.
SortOrderDesc
)
type
SubscribeConfig
¶
added in
v0.1.0
type SubscribeConfig struct {
// ChanSize is the size of the buffered channel that will be created for the
// subscription. Incoming events that overall this number because a listener
// isn't reading from the channel in a timely manner will be dropped.
//
// Defaults to 1000.
ChanSize
int
// Kinds are the kinds of events that the subscription will receive.
// Requiring that kinds are specified explicitly allows for forward
// compatibility in case new kinds of events are added in future versions.
// If new event kinds are added, callers will have to explicitly add them to
// their requested list and ensure they can be handled correctly.
Kinds []
EventKind
}
SubscribeConfig is more thorough subscription configuration used for
Client.SubscribeConfig.
type
TestConfig
¶
added in
v0.17.0
type TestConfig struct {
// DisableUniqueEnforcement disables the application of unique job
// constraints. This is useful for testing scenarios when testing a worker
// that typically uses uniqueness, but where enforcing uniqueness would cause
// conflicts with parallel test execution.
//
// The [rivertest.Worker] type automatically disables uniqueness enforcement
// when creating jobs.
DisableUniqueEnforcement
bool
// Time is a time generator to make time stubbable in tests.
Time
rivertype
.
TimeGenerator
}
TestConfig contains configuration specific to test environments.
type
UniqueOpts
¶
type UniqueOpts struct {
// ByArgs indicates that uniqueness should be enforced for any specific
// instance of encoded args for a job.
//
// Default is false, meaning that as long as any other unique property is
// enabled, uniqueness will be enforced for a kind regardless of input args.
//
// When set to true, the entire encoded args field will be included in the
// uniqueness hash, which requires care to ensure that no irrelevant args are
// factored into the uniqueness check. It is also possible to use a subset of
// the args by indicating on the `JobArgs` struct which fields should be
// included in the uniqueness check using struct tags:
//
// 	type MyJobArgs struct {
// 		CustomerID string `json:"customer_id" river:"unique"`
// 		TraceID    string `json:"trace_id"
// 	}
//
// In this example, only the encoded `customer_id` key will be included in the
// uniqueness check and the `trace_id` key will be ignored.
//
// All keys are sorted alphabetically before hashing to ensure consistent
// results.
//
// River recurses into embedded structs and fields with struct values and
// looks for `river:"unique"` annotations on them as well:
//
// 	type MyJobArgs struct {
// 		Customer *Customer `json:"customer"`
// 		TraceID  string    `json:"trace_id"
// 	}
//
// 	type Customer struct {
// 		ID string `json:"id" river:"unique"`
// 	}
//
// In this example, the `id` value inside a `customer` subboject is used in
// the uniqueness check. It'd be the same story if Customer was embedded on
// MyJobArgs instead:
//
// 	type MyJobArgs struct {
// 		Customer
// 		TraceID string `json:"trace_id"
// 	}
//
// If the struct field itself has a `river:"unique"` annotation, but none on
// any fields in the substruct, then the entire JSON encoded value of the
// struct is used as a unique value:
//
// 	type MyJobArgs struct {
// 		Customer *Customer `json:"customer" river:"unique"`
// 		TraceID  string    `json:"trace_id"
// 	}
//
// 	type Customer struct {
// 		ID string `json:"id"`
// 	}
ByArgs
bool
// ByPeriod defines uniqueness within a given period. On an insert time is
// rounded down to the nearest multiple of the given period, and a job is
// only inserted if there isn't an existing job that will run between then
// and the next multiple of the period.
//
// Default is no unique period, meaning that as long as any other unique
// property is enabled, uniqueness will be enforced across all jobs of the
// kind in the database, regardless of when they were scheduled.
ByPeriod
time
.
Duration
// ByQueue indicates that uniqueness should be enforced within each queue.
//
// Default is false, meaning that as long as any other unique property is
// enabled, uniqueness will be enforced for a kind across all queues.
ByQueue
bool
// ByState indicates that uniqueness should be enforced across any of the
// states in the given set. Unlike other unique options, ByState gets a
// default when it's not set for user convenience. The default is equivalent
// to:
//
// 	ByState: []rivertype.JobState{rivertype.JobStateAvailable, rivertype.JobStateCompleted, rivertype.JobStatePending, rivertype.JobStateRunning, rivertype.JobStateRetryable, rivertype.JobStateScheduled}
//
// Or more succinctly:
//
// 	ByState: rivertype.UniqueOptsByStateDefault()
//
// With this setting, any jobs of the same kind that have been completed or
// discarded, but not yet cleaned out by the system, will still prevent a
// duplicate unique job from being inserted. For example, with the default
// states, if a unique job is actively `running`, a duplicate cannot be
// inserted. Likewise, if a unique job has `completed`, you still can't
// insert a duplicate, at least not until the job cleaner maintenance process
// eventually removes the completed job from the `river_job` table.
//
// The list may be safely customized to _add_ additional states (`cancelled`
// or `discarded`), though only `retryable` may be safely _removed_ from the
// list.
//
// Warning: Removing any states from the default list (other than `retryable`)
// forces a fallback to a slower insertion path that takes an advisory lock
// and performs a look up before insertion. This path is deprecated and should
// be avoided if possible.
ByState []
rivertype
.
JobState
// ExcludeKind indicates that the job kind should not be included in the
// uniqueness check. This is useful when you want to enforce uniqueness
// across all jobs regardless of kind.
ExcludeKind
bool
}
UniqueOpts contains parameters for uniqueness for a job.
When the options struct is uninitialized (its zero value) no uniqueness at is
enforced. As each property is initialized, it's added as a dimension on the
uniqueness matrix. When any property has a non-zero value specified, the
job's kind automatically counts toward uniqueness, but can be excluded by
setting ExcludeKind to true.
So for example, if only ByQueue is on, then for the given job kind, only a
single instance is allowed in any given queue, regardless of other properties
on the job. If both ByArgs and ByQueue are on, then for the given job kind, a
single instance is allowed for each combination of args and queues. If either
args or queue is changed on a new job, it's allowed to be inserted as a new
job.
Uniqueness relies on a hash of the job kind and any unique properties along
with a database unique constraint. See the note on ByState for more details
including about the fallback to a deprecated advisory lock method.
type
UnknownJobKindError
¶
type UnknownJobKindError =
rivertype
.
UnknownJobKindError
UnknownJobKindError is returned when a Client fetches and attempts to
work a job that has not been registered on the Client's Workers bundle (using AddWorker).
type
Worker
¶
type Worker[T
JobArgs
] interface {
// Middleware returns the type-specific middleware for this job.
Middleware(job *
rivertype
.
JobRow
) []
rivertype
.
WorkerMiddleware
// NextRetry calculates when the next retry for a failed job should take
// place given when it was last attempted and its number of attempts, or any
// other of the job's properties a user-configured retry policy might want
// to consider.
//
// Note that this method on a worker overrides any client-level retry policy.
// To use the client-level retry policy, return an empty `time.Time{}` or
// include WorkerDefaults to do this for you.
NextRetry(job *
Job
[T])
time
.
Time
// Timeout is the maximum amount of time the job is allowed to run before
// its context is cancelled. A timeout of zero (the default) means the job
// will inherit the Client-level timeout. A timeout of -1 means the job's
// context will never time out.
Timeout(job *
Job
[T])
time
.
Duration
// Work performs the job and returns an error if the job failed. The context
// will be configured with a timeout according to the worker settings and may
// be cancelled for other reasons.
//
// If no error is returned, the job is assumed to have succeeded and will be
// marked completed.
//
// It is important for any worker to respect context cancellation to enable
// the client to respond to shutdown requests; there is no way to cancel a
// running job that does not respect context cancellation, other than
// terminating the process.
Work(ctx
context
.
Context
, job *
Job
[T])
error
}
Worker is an interface that can perform a job with args of type T. A typical
implementation will be a JSON-serializable `JobArgs` struct that implements
`Kind()`, along with a Worker that embeds WorkerDefaults and implements `Work()`.
Workers may optionally override other methods to provide job-specific
configuration for all jobs of that type:
type SleepArgs struct {
Duration time.Duration `json:"duration"`
}

func (SleepArgs) Kind() string { return "sleep" }

type SleepWorker struct {
WorkerDefaults[SleepArgs]
}

func (w *SleepWorker) Work(ctx context.Context, job *Job[SleepArgs]) error {
select {
case <-ctx.Done():
return ctx.Err()
case <-time.After(job.Args.Duration):
return nil
}
}
In addition to fulfilling the Worker interface, workers must be registered
with the client using the AddWorker function.
func
WorkFunc
¶
func WorkFunc[T
JobArgs
](f func(
context
.
Context
, *
Job
[T])
error
)
Worker
[T]
WorkFunc wraps a function to implement the Worker interface. A job args
struct implementing JobArgs will still be required to specify a Kind.
For example:
river.AddWorker(workers, river.WorkFunc(func(ctx context.Context, job *river.Job[WorkFuncArgs]) error {
fmt.Printf("Message: %s", job.Args.Message)
return nil
}))
type
WorkerDefaults
¶
type WorkerDefaults[T
JobArgs
] struct{}
WorkerDefaults is an empty struct that can be embedded in your worker
struct to make it fulfill the Worker interface with default values.
func (WorkerDefaults[T])
Middleware
¶
added in
v0.13.0
func (w
WorkerDefaults
[T]) Middleware(*
rivertype
.
JobRow
) []
rivertype
.
WorkerMiddleware
func (WorkerDefaults[T])
NextRetry
¶
func (w
WorkerDefaults
[T]) NextRetry(*
Job
[T])
time
.
Time
NextRetry returns an empty time.Time{} to avoid setting any job or
Worker-specific overrides on the next retry time. This means that the
Client-level retry policy schedule will be used instead.
func (WorkerDefaults[T])
Timeout
¶
func (w
WorkerDefaults
[T]) Timeout(*
Job
[T])
time
.
Duration
Timeout returns the job-specific timeout. Override this method to set a
job-specific timeout, otherwise the Client-level timeout will be applied.
type
WorkerMiddlewareDefaults
deprecated
added in
v0.13.0
type WorkerMiddlewareDefaults struct{
MiddlewareDefaults
}
WorkerInsertMiddlewareDefaults is an embeddable struct that provides default
implementations for the rivertype.WorkerMiddleware. Use of this struct is
recommended in case rivertype.WorkerMiddleware is expanded in the future so
that existing code isn't unexpectedly broken during an upgrade.
Deprecated: Prefer embedding the more general MiddlewareDefaults instead.
func (*WorkerMiddlewareDefaults)
Work
¶
added in
v0.13.0
func (d *
WorkerMiddlewareDefaults
) Work(ctx
context
.
Context
, job *
rivertype
.
JobRow
, doInner func(ctx
context
.
Context
)
error
)
error
type
WorkerMiddlewareFunc
¶
added in
v0.21.0
type WorkerMiddlewareFunc func(ctx
context
.
Context
, job *
rivertype
.
JobRow
, doInner func(ctx
context
.
Context
)
error
)
error
WorkerMiddlewareFunc is a convenience helper for implementing
rivertype.WorkerMiddleware using a simple function instead of a struct.
func (WorkerMiddlewareFunc)
IsMiddleware
¶
added in
v0.21.0
func (f
WorkerMiddlewareFunc
) IsMiddleware()
bool
func (WorkerMiddlewareFunc)
Work
¶
added in
v0.21.0
func (f
WorkerMiddlewareFunc
) Work(ctx
context
.
Context
, job *
rivertype
.
JobRow
, doInner func(ctx
context
.
Context
)
error
)
error
type
Workers
¶
type Workers struct {
// contains filtered or unexported fields
}
Workers is a list of available job workers. A Worker must be registered for
each type of Job to be handled.
Use the top-level AddWorker function combined with a Workers to register a
worker.
func
NewWorkers
¶
func NewWorkers() *
Workers
NewWorkers initializes a new registry of available job workers.
Use the top-level AddWorker function combined with a Workers registry to
register each available worker.