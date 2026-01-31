# pgx PostgreSQL Driver

> Source: https://pkg.go.dev/github.com/jackc/pgx/v5
> Fetched: 2026-01-31T10:57:30.079719+00:00
> Content-Hash: 77eef8486c31f787
> Type: html

---

### Overview ¶

  * Establishing a Connection
  * Connection Pool
  * Query Interface
  * PostgreSQL Data Types
  * Transactions
  * Prepared Statements
  * Copy Protocol
  * Listen and Notify
  * Tracing and Logging
  * Lower Level PostgreSQL Functionality
  * PgBouncer



Package pgx is a PostgreSQL database driver. 

pgx provides a native PostgreSQL driver and can act as a database/sql driver. The native PostgreSQL interface is similar to the database/sql interface while providing better speed and access to PostgreSQL specific features. Use github.com/jackc/pgx/v5/stdlib to use pgx as a database/sql compatible driver. See that package's documentation for details. 

#### Establishing a Connection ¶

The primary way of establishing a connection is with pgx.Connect: 
    
    
    conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    

The database connection string can be in URL or key/value format. Both PostgreSQL settings and pgx settings can be specified here. In addition, a config struct can be created by ParseConfig and modified before establishing the connection with ConnectConfig to configure settings such as tracing that cannot be configured with a connection string. 

#### Connection Pool ¶

*pgx.Conn represents a single connection to the database and is not concurrency safe. Use package github.com/jackc/pgx/v5/pgxpool for a concurrency safe connection pool. 

#### Query Interface ¶

pgx implements Query in the familiar database/sql style. However, pgx provides generic functions such as CollectRows and ForEachRow that are a simpler and safer way of processing rows than manually calling defer rows.Close(), rows.Next(), rows.Scan, and rows.Err(). 

CollectRows can be used collect all returned rows into a slice. 
    
    
    rows, _ := conn.Query(context.Background(), "select generate_series(1,$1)", 5)
    numbers, err := pgx.CollectRows(rows, pgx.RowTo[int32])
    if err != nil {
      return err
    }
    // numbers => [1 2 3 4 5]
    

ForEachRow can be used to execute a callback function for every row. This is often easier than iterating over rows directly. 
    
    
    var sum, n int32
    rows, _ := conn.Query(context.Background(), "select generate_series(1,$1)", 10)
    _, err := pgx.ForEachRow(rows, []any{&n}, func() error {
      sum += n
      return nil
    })
    if err != nil {
      return err
    }
    

pgx also implements QueryRow in the same style as database/sql. 
    
    
    var name string
    var weight int64
    err := conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
    if err != nil {
        return err
    }
    

Use Exec to execute a query that does not return a result set. 
    
    
    commandTag, err := conn.Exec(context.Background(), "delete from widgets where id=$1", 42)
    if err != nil {
        return err
    }
    if commandTag.RowsAffected() != 1 {
        return errors.New("No row found to delete")
    }
    

#### PostgreSQL Data Types ¶

pgx uses the pgtype package to converting Go values to and from PostgreSQL values. It supports many PostgreSQL types directly and is customizable and extendable. User defined data types such as enums, domains, and composite types may require type registration. See that package's documentation for details. 

#### Transactions ¶

Transactions are started by calling Begin. 
    
    
    tx, err := conn.Begin(context.Background())
    if err != nil {
        return err
    }
    // Rollback is safe to call even if the tx is already closed, so if
    // the tx commits successfully, this is a no-op
    defer tx.Rollback(context.Background())
    
    _, err = tx.Exec(context.Background(), "insert into foo(id) values (1)")
    if err != nil {
        return err
    }
    
    err = tx.Commit(context.Background())
    if err != nil {
        return err
    }
    

The Tx returned from Begin also implements the Begin method. This can be used to implement pseudo nested transactions. These are internally implemented with savepoints. 

Use BeginTx to control the transaction mode. BeginTx also can be used to ensure a new transaction is created instead of a pseudo nested transaction. 

BeginFunc and BeginTxFunc are functions that begin a transaction, execute a function, and commit or rollback the transaction depending on the return value of the function. These can be simpler and less error prone to use. 
    
    
    err = pgx.BeginFunc(context.Background(), conn, func(tx pgx.Tx) error {
        _, err := tx.Exec(context.Background(), "insert into foo(id) values (1)")
        return err
    })
    if err != nil {
        return err
    }
    

#### Prepared Statements ¶

Prepared statements can be manually created with the Prepare method. However, this is rarely necessary because pgx includes an automatic statement cache by default. Queries run through the normal Query, QueryRow, and Exec functions are automatically prepared on first execution and the prepared statement is reused on subsequent executions. See ParseConfig for information on how to customize or disable the statement cache. 

#### Copy Protocol ¶

Use CopyFrom to efficiently insert multiple rows at a time using the PostgreSQL copy protocol. CopyFrom accepts a CopyFromSource interface. If the data is already in a [][]any use CopyFromRows to wrap it in a CopyFromSource interface. Or implement CopyFromSource to avoid buffering the entire data set in memory. 
    
    
    rows := [][]any{
        {"John", "Smith", int32(36)},
        {"Jane", "Doe", int32(29)},
    }
    
    copyCount, err := conn.CopyFrom(
        context.Background(),
        pgx.Identifier{"people"},
        []string{"first_name", "last_name", "age"},
        pgx.CopyFromRows(rows),
    )
    

When you already have a typed array using CopyFromSlice can be more convenient. 
    
    
    rows := []User{
        {"John", "Smith", 36},
        {"Jane", "Doe", 29},
    }
    
    copyCount, err := conn.CopyFrom(
        context.Background(),
        pgx.Identifier{"people"},
        []string{"first_name", "last_name", "age"},
        pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
            return []any{rows[i].FirstName, rows[i].LastName, rows[i].Age}, nil
        }),
    )
    

CopyFrom can be faster than an insert with as few as 5 rows. 

#### Listen and Notify ¶

pgx can listen to the PostgreSQL notification system with the `Conn.WaitForNotification` method. It blocks until a notification is received or the context is canceled. 
    
    
    _, err := conn.Exec(context.Background(), "listen channelname")
    if err != nil {
        return err
    }
    
    notification, err := conn.WaitForNotification(context.Background())
    if err != nil {
        return err
    }
    // do something with notification
    

#### Tracing and Logging ¶

pgx supports tracing by setting ConnConfig.Tracer. To combine several tracers you can use the multitracer.Tracer. 

In addition, the tracelog package provides the TraceLog type which lets a traditional logger act as a Tracer. 

For debug tracing of the actual PostgreSQL wire protocol messages see github.com/jackc/pgx/v5/pgproto3. 

#### Lower Level PostgreSQL Functionality ¶

github.com/jackc/pgx/v5/pgconn contains a lower level PostgreSQL driver roughly at the level of libpq. pgx.Conn is implemented on top of pgconn. The Conn.PgConn() method can be used to access this lower layer. 

#### PgBouncer ¶

By default pgx automatically uses prepared statements. Prepared statements are incompatible with PgBouncer. This can be disabled by setting a different QueryExecMode in ConnConfig.DefaultQueryExecMode. 

### Index ¶

  * Constants
  * Variables
  * func AppendRows[T any, S ~[]T](slice S, rows Rows, fn RowToFunc[T]) (S, error)
  * func BeginFunc(ctx context.Context, db interface{ ... }, fn func(Tx) error) (err error)
  * func BeginTxFunc(ctx context.Context, db interface{ ... }, txOptions TxOptions, ...) (err error)
  * func CollectExactlyOneRow[T any](rows Rows, fn RowToFunc[T]) (T, error)
  * func CollectOneRow[T any](rows Rows, fn RowToFunc[T]) (T, error)
  * func CollectRows[T any](rows Rows, fn RowToFunc[T]) ([]T, error)
  * func ForEachRow(rows Rows, scans []any, fn func() error) (pgconn.CommandTag, error)
  * func RowTo[T any](row CollectableRow) (T, error)
  * func RowToAddrOf[T any](row CollectableRow) (*T, error)
  * func RowToAddrOfStructByName[T any](row CollectableRow) (*T, error)
  * func RowToAddrOfStructByNameLax[T any](row CollectableRow) (*T, error)
  * func RowToAddrOfStructByPos[T any](row CollectableRow) (*T, error)
  * func RowToMap(row CollectableRow) (map[string]any, error)
  * func RowToStructByName[T any](row CollectableRow) (T, error)
  * func RowToStructByNameLax[T any](row CollectableRow) (T, error)
  * func RowToStructByPos[T any](row CollectableRow) (T, error)
  * func ScanRow(typeMap *pgtype.Map, fieldDescriptions []pgconn.FieldDescription, ...) error
  * type Batch
  *     * func (b *Batch) Len() int
    * func (b *Batch) Queue(query string, arguments ...any) *QueuedQuery
  * type BatchResults
  * type BatchTracer
  * type CollectableRow
  * type Conn
  *     * func Connect(ctx context.Context, connString string) (*Conn, error)
    * func ConnectConfig(ctx context.Context, connConfig *ConnConfig) (*Conn, error)
    * func ConnectWithOptions(ctx context.Context, connString string, options ParseConfigOptions) (*Conn, error)
  *     * func (c *Conn) Begin(ctx context.Context) (Tx, error)
    * func (c *Conn) BeginTx(ctx context.Context, txOptions TxOptions) (Tx, error)
    * func (c *Conn) Close(ctx context.Context) error
    * func (c *Conn) Config() *ConnConfig
    * func (c *Conn) CopyFrom(ctx context.Context, tableName Identifier, columnNames []string, ...) (int64, error)
    * func (c *Conn) Deallocate(ctx context.Context, name string) error
    * func (c *Conn) DeallocateAll(ctx context.Context) error
    * func (c *Conn) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
    * func (c *Conn) IsClosed() bool
    * func (c *Conn) LoadType(ctx context.Context, typeName string) (*pgtype.Type, error)
    * func (c *Conn) LoadTypes(ctx context.Context, typeNames []string) ([]*pgtype.Type, error)
    * func (c *Conn) PgConn() *pgconn.PgConn
    * func (c *Conn) Ping(ctx context.Context) error
    * func (c *Conn) Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
    * func (c *Conn) Query(ctx context.Context, sql string, args ...any) (Rows, error)
    * func (c *Conn) QueryRow(ctx context.Context, sql string, args ...any) Row
    * func (c *Conn) SendBatch(ctx context.Context, b *Batch) (br BatchResults)
    * func (c *Conn) TypeMap() *pgtype.Map
    * func (c *Conn) WaitForNotification(ctx context.Context) (*pgconn.Notification, error)
  * type ConnConfig
  *     * func ParseConfig(connString string) (*ConnConfig, error)
    * func ParseConfigWithOptions(connString string, options ParseConfigOptions) (*ConnConfig, error)
  *     * func (cc *ConnConfig) ConnString() string
    * func (cc *ConnConfig) Copy() *ConnConfig
  * type ConnectTracer
  * type CopyFromSource
  *     * func CopyFromFunc(nxtf func() (row []any, err error)) CopyFromSource
    * func CopyFromRows(rows [][]any) CopyFromSource
    * func CopyFromSlice(length int, next func(int) ([]any, error)) CopyFromSource
  * type CopyFromTracer
  * type ExtendedQueryBuilder
  *     * func (eqb *ExtendedQueryBuilder) Build(m *pgtype.Map, sd *pgconn.StatementDescription, args []any) error
  * type Identifier
  *     * func (ident Identifier) Sanitize() string
  * type LargeObject
  *     * func (o *LargeObject) Close() error
    * func (o *LargeObject) Read(p []byte) (int, error)
    * func (o *LargeObject) Seek(offset int64, whence int) (n int64, err error)
    * func (o *LargeObject) Tell() (n int64, err error)
    * func (o *LargeObject) Truncate(size int64) (err error)
    * func (o *LargeObject) Write(p []byte) (int, error)
  * type LargeObjectMode
  * type LargeObjects
  *     * func (o *LargeObjects) Create(ctx context.Context, oid uint32) (uint32, error)
    * func (o *LargeObjects) Open(ctx context.Context, oid uint32, mode LargeObjectMode) (*LargeObject, error)
    * func (o *LargeObjects) Unlink(ctx context.Context, oid uint32) error
  * type NamedArgs
  *     * func (na NamedArgs) RewriteQuery(ctx context.Context, conn *Conn, sql string, args []any) (newSQL string, newArgs []any, err error)
  * type ParseConfigOptions
  * type PrepareTracer
  * type QueryExecMode
  *     * func (m QueryExecMode) String() string
  * type QueryResultFormats
  * type QueryResultFormatsByOID
  * type QueryRewriter
  * type QueryTracer
  * type QueuedQuery
  *     * func (qq *QueuedQuery) Exec(fn func(ct pgconn.CommandTag) error)
    * func (qq *QueuedQuery) Query(fn func(rows Rows) error)
    * func (qq *QueuedQuery) QueryRow(fn func(row Row) error)
  * type Row
  * type RowScanner
  * type RowToFunc
  * type Rows
  *     * func RowsFromResultReader(typeMap *pgtype.Map, resultReader *pgconn.ResultReader) Rows
  * type ScanArgError
  *     * func (e ScanArgError) Error() string
    * func (e ScanArgError) Unwrap() error
  * type StrictNamedArgs
  *     * func (sna StrictNamedArgs) RewriteQuery(ctx context.Context, conn *Conn, sql string, args []any) (newSQL string, newArgs []any, err error)
  * type TraceBatchEndData
  * type TraceBatchQueryData
  * type TraceBatchStartData
  * type TraceConnectEndData
  * type TraceConnectStartData
  * type TraceCopyFromEndData
  * type TraceCopyFromStartData
  * type TracePrepareEndData
  * type TracePrepareStartData
  * type TraceQueryEndData
  * type TraceQueryStartData
  * type Tx
  * type TxAccessMode
  * type TxDeferrableMode
  * type TxIsoLevel
  * type TxOptions



### Examples ¶

  * CollectRows
  * Conn.Query
  * Conn.SendBatch
  * ForEachRow
  * RowTo
  * RowToAddrOf
  * RowToStructByName
  * RowToStructByNameLax
  * RowToStructByPos



### Constants ¶

[View Source](https://github.com/jackc/pgx/blob/v5.8.0/values.go#L11)
    
    
    const (
    	TextFormatCode   = 0
    	BinaryFormatCode = 1
    )

PostgreSQL format codes 

### Variables ¶

[View Source](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L105)
    
    
    var (
    	// ErrNoRows occurs when rows are expected but none are returned.
    	ErrNoRows = newProxyErr([sql](/database/sql).[ErrNoRows](/database/sql#ErrNoRows), "no rows in result set")
    	// ErrTooManyRows occurs when more rows than expected are returned.
    	ErrTooManyRows = [errors](/errors).[New](/errors#New)("too many rows in result set")
    )

[View Source](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L85)
    
    
    var ErrTxClosed = [errors](/errors).[New](/errors#New)("tx is closed")

[View Source](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L90)
    
    
    var ErrTxCommitRollback = [errors](/errors).[New](/errors#New)("commit unexpectedly resulted in rollback")

ErrTxCommitRollback occurs when an error has occurred in a transaction and Commit() is called. PostgreSQL accepts COMMIT on aborted transactions, but it is treated as ROLLBACK. 

### Functions ¶

####  func [AppendRows](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L437) ¶ added in v5.5.3
    
    
    func AppendRows[T [any](/builtin#any), S ~[]T](slice S, rows Rows, fn RowToFunc[T]) (S, [error](/builtin#error))

AppendRows iterates through rows, calling fn for each row, and appending the results into a slice of T. 

This function closes the rows automatically on return. 

####  func [BeginFunc](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L391) ¶
    
    
    func BeginFunc(
    	ctx [context](/context).[Context](/context#Context),
    	db interface {
    		Begin(ctx [context](/context).[Context](/context#Context)) (Tx, [error](/builtin#error))
    	},
    	fn func(Tx) [error](/builtin#error),
    ) (err [error](/builtin#error))

BeginFunc calls Begin on db and then calls fn. If fn does not return an error then it calls Commit on db. If fn returns an error it calls Rollback on db. The context will be used when executing the transaction control statements (BEGIN, ROLLBACK, and COMMIT) but does not otherwise affect the execution of fn. 

####  func [BeginTxFunc](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L410) ¶
    
    
    func BeginTxFunc(
    	ctx [context](/context).[Context](/context#Context),
    	db interface {
    		BeginTx(ctx [context](/context).[Context](/context#Context), txOptions TxOptions) (Tx, [error](/builtin#error))
    	},
    	txOptions TxOptions,
    	fn func(Tx) [error](/builtin#error),
    ) (err [error](/builtin#error))

BeginTxFunc calls BeginTx on db and then calls fn. If fn does not return an error then it calls Commit on db. If fn returns an error it calls Rollback on db. The context will be used when executing the transaction control statements (BEGIN, ROLLBACK, and COMMIT) but does not otherwise affect the execution of fn. 

####  func [CollectExactlyOneRow](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L495) ¶ added in v5.5.0
    
    
    func CollectExactlyOneRow[T [any](/builtin#any)](rows Rows, fn RowToFunc[T]) (T, [error](/builtin#error))

CollectExactlyOneRow calls fn for the first row in rows and returns the result. 

  * If no rows are found returns an error where errors.Is(ErrNoRows) is true.
  * If more than 1 row is found returns an error where errors.Is(ErrTooManyRows) is true.



This function closes the rows automatically on return. 

####  func [CollectOneRow](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L466) ¶
    
    
    func CollectOneRow[T [any](/builtin#any)](rows Rows, fn RowToFunc[T]) (T, [error](/builtin#error))

CollectOneRow calls fn for the first row in rows and returns the result. If no rows are found returns an error where errors.Is(ErrNoRows) is true. CollectOneRow is to CollectRows as QueryRow is to Query. 

This function closes the rows automatically on return. 

####  func [CollectRows](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L458) ¶
    
    
    func CollectRows[T [any](/builtin#any)](rows Rows, fn RowToFunc[T]) ([]T, [error](/builtin#error))

CollectRows iterates through rows, calling fn for each row, and collecting the results into a slice of T. 

This function closes the rows automatically on return. 

Example ¶

This example uses CollectRows with a manually written collector function. In most cases RowTo, RowToAddrOf, RowToStructByPos, RowToAddrOfStructByPos, or another generic function would be used. 
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    rows, _ := conn.Query(ctx, `select n from generate_series(1, 5) n`)
    numbers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (int32, error) {
    	var n int32
    	err := row.Scan(&n)
    	return n, err
    })
    if err != nil {
    	fmt.Printf("CollectRows error: %v", err)
    	return
    }
    
    fmt.Println(numbers)
    
    
    
    Output:
    
    [1 2 3 4 5]
    

####  func [ForEachRow](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L401) ¶
    
    
    func ForEachRow(rows Rows, scans [][any](/builtin#any), fn func() [error](/builtin#error)) ([pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag), [error](/builtin#error))

ForEachRow iterates through rows. For each row it scans into the elements of scans and calls fn. If any row fails to scan or fn returns an error the query will be aborted and the error will be returned. Rows will be closed when ForEachRow returns. 

Example ¶
    
    
    conn, err := pgx.Connect(context.Background(), os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    rows, _ := conn.Query(
    	context.Background(),
    	"select n, n * 2 from generate_series(1, $1) n",
    	3,
    )
    var a, b int
    _, err = pgx.ForEachRow(rows, []any{&a, &b}, func() error {
    	fmt.Printf("%v, %v\n", a, b)
    	return nil
    })
    if err != nil {
    	fmt.Printf("ForEachRow error: %v", err)
    	return
    }
    
    
    
    Output:
    
    1, 2
    2, 4
    3, 6
    

####  func [RowTo](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L526) ¶
    
    
    func RowTo[T [any](/builtin#any)](row CollectableRow) (T, [error](/builtin#error))

RowTo returns a T scanned from row. 

Example ¶
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    rows, _ := conn.Query(ctx, `select n from generate_series(1, 5) n`)
    numbers, err := pgx.CollectRows(rows, pgx.RowTo[int32])
    if err != nil {
    	fmt.Printf("CollectRows error: %v", err)
    	return
    }
    
    fmt.Println(numbers)
    
    
    
    Output:
    
    [1 2 3 4 5]
    

####  func [RowToAddrOf](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L533) ¶
    
    
    func RowToAddrOf[T [any](/builtin#any)](row CollectableRow) (*T, [error](/builtin#error))

RowTo returns a the address of a T scanned from row. 

Example ¶
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    rows, _ := conn.Query(ctx, `select n from generate_series(1, 5) n`)
    pNumbers, err := pgx.CollectRows(rows, pgx.RowToAddrOf[int32])
    if err != nil {
    	fmt.Printf("CollectRows error: %v", err)
    	return
    }
    
    for _, p := range pNumbers {
    	fmt.Println(*p)
    }
    
    
    
    Output:
    
    1
    2
    3
    4
    5
    

####  func [RowToAddrOfStructByName](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L654) ¶ added in v5.1.0
    
    
    func RowToAddrOfStructByName[T [any](/builtin#any)](row CollectableRow) (*T, [error](/builtin#error))

RowToAddrOfStructByName returns the address of a T scanned from row. T must be a struct. T must have the same number of named public fields as row has fields. The row and T fields will be matched by name. The match is case-insensitive. The database column name can be overridden with a "db" struct tag. If the "db" struct tag is "-" then the field will be ignored. 

####  func [RowToAddrOfStructByNameLax](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L673) ¶ added in v5.4.0
    
    
    func RowToAddrOfStructByNameLax[T [any](/builtin#any)](row CollectableRow) (*T, [error](/builtin#error))

RowToAddrOfStructByNameLax returns the address of a T scanned from row. T must be a struct. T must have greater than or equal number of named public fields as row has fields. The row and T fields will be matched by name. The match is case-insensitive. The database column name can be overridden with a "db" struct tag. If the "db" struct tag is "-" then the field will be ignored. 

####  func [RowToAddrOfStructByPos](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L575) ¶
    
    
    func RowToAddrOfStructByPos[T [any](/builtin#any)](row CollectableRow) (*T, [error](/builtin#error))

RowToAddrOfStructByPos returns the address of a T scanned from row. T must be a struct. T must have the same number a public fields as row has fields. The row and T fields will be matched by position. If the "db" struct tag is "-" then the field will be ignored. 

####  func [RowToMap](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L540) ¶
    
    
    func RowToMap(row CollectableRow) (map[[string](/builtin#string)][any](/builtin#any), [error](/builtin#error))

RowToMap returns a map scanned from row. 

####  func [RowToStructByName](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L644) ¶ added in v5.1.0
    
    
    func RowToStructByName[T [any](/builtin#any)](row CollectableRow) (T, [error](/builtin#error))

RowToStructByName returns a T scanned from row. T must be a struct. T must have the same number of named public fields as row has fields. The row and T fields will be matched by name. The match is case-insensitive. The database column name can be overridden with a "db" struct tag. If the "db" struct tag is "-" then the field will be ignored. 

Example ¶
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    if conn.PgConn().ParameterStatus("crdb_version") != "" {
    	// Skip test / example when running on CockroachDB. Since an example can't be skipped fake success instead.
    	fmt.Println(`Cheeseburger: $10
    Fries: $5
    Soft Drink: $3`)
    	return
    }
    
    // Setup example schema and data.
    _, err = conn.Exec(ctx, `
    create temporary table products (
    	id int primary key generated by default as identity,
    	name varchar(100) not null,
    	price int not null
    );
    
    insert into products (name, price) values
    	('Cheeseburger', 10),
    	('Double Cheeseburger', 14),
    	('Fries', 5),
    	('Soft Drink', 3);
    `)
    if err != nil {
    	fmt.Printf("Unable to setup example schema and data: %v", err)
    	return
    }
    
    type product struct {
    	ID    int32
    	Name  string
    	Price int32
    }
    
    rows, _ := conn.Query(ctx, "select * from products where price < $1 order by price desc", 12)
    products, err := pgx.CollectRows(rows, pgx.RowToStructByName[product])
    if err != nil {
    	fmt.Printf("CollectRows error: %v", err)
    	return
    }
    
    for _, p := range products {
    	fmt.Printf("%s: $%d\n", p.Name, p.Price)
    }
    
    
    
    Output:
    
    Cheeseburger: $10
    Fries: $5
    Soft Drink: $3
    

####  func [RowToStructByNameLax](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L663) ¶ added in v5.4.0
    
    
    func RowToStructByNameLax[T [any](/builtin#any)](row CollectableRow) (T, [error](/builtin#error))

RowToStructByNameLax returns a T scanned from row. T must be a struct. T must have greater than or equal number of named public fields as row has fields. The row and T fields will be matched by name. The match is case-insensitive. The database column name can be overridden with a "db" struct tag. If the "db" struct tag is "-" then the field will be ignored. 

Example ¶
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    if conn.PgConn().ParameterStatus("crdb_version") != "" {
    	// Skip test / example when running on CockroachDB. Since an example can't be skipped fake success instead.
    	fmt.Println(`Cheeseburger: $10
    Fries: $5
    Soft Drink: $3`)
    	return
    }
    
    // Setup example schema and data.
    _, err = conn.Exec(ctx, `
    create temporary table products (
    	id int primary key generated by default as identity,
    	name varchar(100) not null,
    	price int not null
    );
    
    insert into products (name, price) values
    	('Cheeseburger', 10),
    	('Double Cheeseburger', 14),
    	('Fries', 5),
    	('Soft Drink', 3);
    `)
    if err != nil {
    	fmt.Printf("Unable to setup example schema and data: %v", err)
    	return
    }
    
    type product struct {
    	ID    int32
    	Name  string
    	Type  string
    	Price int32
    }
    
    rows, _ := conn.Query(ctx, "select * from products where price < $1 order by price desc", 12)
    products, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[product])
    if err != nil {
    	fmt.Printf("CollectRows error: %v", err)
    	return
    }
    
    for _, p := range products {
    	fmt.Printf("%s: $%d\n", p.Name, p.Price)
    }
    
    
    
    Output:
    
    Cheeseburger: $10
    Fries: $5
    Soft Drink: $3
    

####  func [RowToStructByPos](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L566) ¶
    
    
    func RowToStructByPos[T [any](/builtin#any)](row CollectableRow) (T, [error](/builtin#error))

RowToStructByPos returns a T scanned from row. T must be a struct. T must have the same number of public fields as row has fields. The row and T fields will be matched by position. If the "db" struct tag is "-" then the field will be ignored. 

Example ¶
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    if conn.PgConn().ParameterStatus("crdb_version") != "" {
    	// Skip test / example when running on CockroachDB. Since an example can't be skipped fake success instead.
    	fmt.Println(`Cheeseburger: $10
    Fries: $5
    Soft Drink: $3`)
    	return
    }
    
    // Setup example schema and data.
    _, err = conn.Exec(ctx, `
    create temporary table products (
    	id int primary key generated by default as identity,
    	name varchar(100) not null,
    	price int not null
    );
    
    insert into products (name, price) values
    	('Cheeseburger', 10),
    	('Double Cheeseburger', 14),
    	('Fries', 5),
    	('Soft Drink', 3);
    `)
    if err != nil {
    	fmt.Printf("Unable to setup example schema and data: %v", err)
    	return
    }
    
    type product struct {
    	ID    int32
    	Name  string
    	Price int32
    }
    
    rows, _ := conn.Query(ctx, "select * from products where price < $1 order by price desc", 12)
    products, err := pgx.CollectRows(rows, pgx.RowToStructByPos[product])
    if err != nil {
    	fmt.Printf("CollectRows error: %v", err)
    	return
    }
    
    for _, p := range products {
    	fmt.Printf("%s: $%d\n", p.Name, p.Price)
    }
    
    
    
    Output:
    
    Cheeseburger: $10
    Fries: $5
    Soft Drink: $3
    

####  func [ScanRow](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L367) ¶
    
    
    func ScanRow(typeMap *[pgtype](/github.com/jackc/pgx/v5@v5.8.0/pgtype).[Map](/github.com/jackc/pgx/v5@v5.8.0/pgtype#Map), fieldDescriptions [][pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[FieldDescription](/github.com/jackc/pgx/v5@v5.8.0/pgconn#FieldDescription), values [][][byte](/builtin#byte), dest ...[any](/builtin#any)) [error](/builtin#error)

ScanRow decodes raw row data into dest. It can be used to scan rows read from the lower level pgconn interface. 

typeMap - OID to Go type mapping. fieldDescriptions - OID and format of values values - the raw data as returned from the PostgreSQL server dest - the destination that values will be decoded into 

### Types ¶

####  type [Batch](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L63) ¶
    
    
    type Batch struct {
    	QueuedQueries []*QueuedQuery
    }

Batch queries are a way of bundling multiple queries together to avoid unnecessary network round trips. A Batch must only be sent once. 

####  func (*Batch) [Len](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L84) ¶
    
    
    func (b *Batch) Len() [int](/builtin#int)

Len returns number of queries that have been queued so far. 

####  func (*Batch) [Queue](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L74) ¶
    
    
    func (b *Batch) Queue(query [string](/builtin#string), arguments ...[any](/builtin#any)) *QueuedQuery

Queue queues a query to batch b. query can be an SQL query or the name of a prepared statement. The only pgx option argument that is supported is QueryRewriter. Queries are executed using the connection's DefaultQueryExecMode. 

While query can contain multiple statements if the connection's DefaultQueryExecMode is QueryModeSimple, this should be avoided. QueuedQuery.Fn must not be set as it will only be called for the first query. That is, QueuedQuery.Query, QueuedQuery.QueryRow, and QueuedQuery.Exec must not be called. In addition, any error messages or tracing that include the current query may reference the wrong query. 

####  type [BatchResults](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L88) ¶
    
    
    type BatchResults interface {
    	// Exec reads the results from the next query in the batch as if the query has been sent with Conn.Exec. Prefer
    	// calling Exec on the QueuedQuery, or just calling Close.
    	Exec() ([pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag), [error](/builtin#error))
    
    	// Query reads the results from the next query in the batch as if the query has been sent with Conn.Query. Prefer
    	// calling Query on the QueuedQuery.
    	Query() (Rows, [error](/builtin#error))
    
    	// QueryRow reads the results from the next query in the batch as if the query has been sent with Conn.QueryRow.
    	// Prefer calling QueryRow on the QueuedQuery.
    	QueryRow() Row
    
    	// Close closes the batch operation. All unread results are read and any callback functions registered with
    	// QueuedQuery.Query, QueuedQuery.QueryRow, or QueuedQuery.Exec will be called. If a callback function returns an
    	// error or the batch encounters an error subsequent callback functions will not be called.
    	//
    	// For simple batch inserts inside a transaction or similar queries, it's sufficient to not set any callbacks,
    	// and just handle the return value of Close.
    	//
    	// Close must be called before the underlying connection can be used again. Any error that occurred during a batch
    	// operation may have made it impossible to resyncronize the connection with the server. In this case the underlying
    	// connection will have been closed.
    	//
    	// Close is safe to call multiple times. If it returns an error subsequent calls will return the same error. Callback
    	// functions will not be rerun.
    	Close() [error](/builtin#error)
    }

####  type [BatchTracer](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L29) ¶
    
    
    type BatchTracer interface {
    	// TraceBatchStart is called at the beginning of SendBatch calls. The returned context is used for the
    	// rest of the call and will be passed to TraceBatchQuery and TraceBatchEnd.
    	TraceBatchStart(ctx [context](/context).[Context](/context#Context), conn *Conn, data TraceBatchStartData) [context](/context).[Context](/context#Context)
    
    	TraceBatchQuery(ctx [context](/context).[Context](/context#Context), conn *Conn, data TraceBatchQueryData)
    	TraceBatchEnd(ctx [context](/context).[Context](/context#Context), conn *Conn, data TraceBatchEndData)
    }

BatchTracer traces SendBatch. 

####  type [CollectableRow](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L424) ¶
    
    
    type CollectableRow interface {
    	FieldDescriptions() [][pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[FieldDescription](/github.com/jackc/pgx/v5@v5.8.0/pgconn#FieldDescription)
    	Scan(dest ...[any](/builtin#any)) [error](/builtin#error)
    	Values() ([][any](/builtin#any), [error](/builtin#error))
    	RawValues() [][][byte](/builtin#byte)
    }

CollectableRow is the subset of Rows methods that a RowToFunc is allowed to call. 

####  type [Conn](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L67) ¶
    
    
    type Conn struct {
    	// contains filtered or unexported fields
    }

Conn is a PostgreSQL connection handle. It is not safe for concurrent usage. Use a connection pool to manage access to multiple database connections from multiple goroutines. 

####  func [Connect](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L135) ¶
    
    
    func Connect(ctx [context](/context).[Context](/context#Context), connString [string](/builtin#string)) (*Conn, [error](/builtin#error))

Connect establishes a connection with a PostgreSQL server with a connection string. See pgconn.Connect for details. 

####  func [ConnectConfig](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L155) ¶
    
    
    func ConnectConfig(ctx [context](/context).[Context](/context#Context), connConfig *ConnConfig) (*Conn, [error](/builtin#error))

ConnectConfig establishes a connection with a PostgreSQL server with a configuration struct. connConfig must have been created by ParseConfig. 

####  func [ConnectWithOptions](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L145) ¶ added in v5.1.0
    
    
    func ConnectWithOptions(ctx [context](/context).[Context](/context#Context), connString [string](/builtin#string), options ParseConfigOptions) (*Conn, [error](/builtin#error))

ConnectWithOptions behaves exactly like Connect with the addition of options. At the present options is only used to provide a GetSSLPassword function. 

####  func (*Conn) [Begin](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L94) ¶
    
    
    func (c *Conn) Begin(ctx [context](/context).[Context](/context#Context)) (Tx, [error](/builtin#error))

Begin starts a transaction. Unlike database/sql, the context only affects the begin command. i.e. there is no auto-rollback on context cancellation. 

####  func (*Conn) [BeginTx](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L100) ¶
    
    
    func (c *Conn) BeginTx(ctx [context](/context).[Context](/context#Context), txOptions TxOptions) (Tx, [error](/builtin#error))

BeginTx starts a transaction with txOptions determining the transaction mode. Unlike database/sql, the context only affects the begin command. i.e. there is no auto-rollback on context cancellation. 

####  func (*Conn) [Close](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L299) ¶
    
    
    func (c *Conn) Close(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)

Close closes a connection. It is safe to call Close on an already closed connection. 

####  func (*Conn) [Config](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L466) ¶
    
    
    func (c *Conn) Config() *ConnConfig

Config returns a copy of config that was used to establish this connection. 

####  func (*Conn) [CopyFrom](https://github.com/jackc/pgx/blob/v5.8.0/copy_from.go#L265) ¶
    
    
    func (c *Conn) CopyFrom(ctx [context](/context).[Context](/context#Context), tableName Identifier, columnNames [][string](/builtin#string), rowSrc CopyFromSource) ([int64](/builtin#int64), [error](/builtin#error))

CopyFrom uses the PostgreSQL copy protocol to perform bulk data insertion. It returns the number of rows copied and an error. 

CopyFrom requires all values use the binary format. A pgtype.Type that supports the binary format must be registered for the type of each column. Almost all types implemented by pgx support the binary format. 

Even though enum types appear to be strings they still must be registered to use with CopyFrom. This can be done with Conn.LoadType and pgtype.Map.RegisterType. 

####  func (*Conn) [Deallocate](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L373) ¶
    
    
    func (c *Conn) Deallocate(ctx [context](/context).[Context](/context#Context), name [string](/builtin#string)) [error](/builtin#error)

Deallocate releases a prepared statement. Calling Deallocate on a non-existent prepared statement will succeed. 

####  func (*Conn) [DeallocateAll](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L395) ¶ added in v5.1.0
    
    
    func (c *Conn) DeallocateAll(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)

DeallocateAll releases all previously prepared statements from the server and client, where it also resets the statement and description cache. 

####  func (*Conn) [Exec](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L470) ¶
    
    
    func (c *Conn) Exec(ctx [context](/context).[Context](/context#Context), sql [string](/builtin#string), arguments ...[any](/builtin#any)) ([pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag), [error](/builtin#error))

Exec executes sql. sql can be either a prepared statement name or an SQL string. arguments should be referenced positionally from the sql string as $1, $2, etc. 

####  func (*Conn) [IsClosed](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L432) ¶
    
    
    func (c *Conn) IsClosed() [bool](/builtin#bool)

IsClosed reports if the connection has been closed. 

####  func (*Conn) [LoadType](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L1284) ¶
    
    
    func (c *Conn) LoadType(ctx [context](/context).[Context](/context#Context), typeName [string](/builtin#string)) (*[pgtype](/github.com/jackc/pgx/v5@v5.8.0/pgtype).[Type](/github.com/jackc/pgx/v5@v5.8.0/pgtype#Type), [error](/builtin#error))

LoadType inspects the database for typeName and produces a pgtype.Type suitable for registration. typeName must be the name of a type where the underlying type(s) is already understood by pgx. It is for derived types. In particular, typeName must be one of the following: 

  * An array type name of a type that is already registered. e.g. "_foo" when "foo" is registered.
  * A composite type name where all field types are already registered.
  * A domain type name where the base type is already registered.
  * An enum type name.
  * A range type name where the element type is already registered.
  * A multirange type name where the element type is already registered.



####  func (*Conn) [LoadTypes](https://github.com/jackc/pgx/blob/v5.8.0/derived_types.go#L162) ¶ added in v5.7.0
    
    
    func (c *Conn) LoadTypes(ctx [context](/context).[Context](/context#Context), typeNames [][string](/builtin#string)) ([]*[pgtype](/github.com/jackc/pgx/v5@v5.8.0/pgtype).[Type](/github.com/jackc/pgx/v5@v5.8.0/pgtype#Type), [error](/builtin#error))

LoadTypes performs a single (complex) query, returning all the required information to register the named types, as well as any other types directly or indirectly required to complete the registration. The result of this call can be passed into RegisterTypes to complete the process. 

####  func (*Conn) [PgConn](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L460) ¶
    
    
    func (c *Conn) PgConn() *[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[PgConn](/github.com/jackc/pgx/v5@v5.8.0/pgconn#PgConn)

PgConn returns the underlying *pgconn.PgConn. This is an escape hatch method that allows lower level access to the PostgreSQL connection than pgx exposes. 

It is strongly recommended that the connection be idle (no in-progress queries) before the underlying *pgconn.PgConn is used and the connection must be returned to the same state before any *pgx.Conn methods are again used. 

####  func (*Conn) [Ping](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L451) ¶
    
    
    func (c *Conn) Ping(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)

Ping delegates to the underlying *pgconn.PgConn.Ping. 

####  func (*Conn) [Prepare](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L317) ¶
    
    
    func (c *Conn) Prepare(ctx [context](/context).[Context](/context#Context), name, sql [string](/builtin#string)) (sd *[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[StatementDescription](/github.com/jackc/pgx/v5@v5.8.0/pgconn#StatementDescription), err [error](/builtin#error))

Prepare creates a prepared statement with name and sql. sql can contain placeholders for bound parameters. These placeholders are referenced positionally as $1, $2, etc. name can be used instead of sql with Query, QueryRow, and Exec to execute the statement. It can also be used with Batch.Queue. 

The underlying PostgreSQL identifier for the prepared statement will be name if name != sql or a digest of sql if name == sql. 

Prepare is idempotent; i.e. it is safe to call Prepare multiple times with the same name and sql arguments. This allows a code path to Prepare and Query/Exec without concern for if the statement has already been prepared. 

####  func (*Conn) [Query](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L749) ¶
    
    
    func (c *Conn) Query(ctx [context](/context).[Context](/context#Context), sql [string](/builtin#string), args ...[any](/builtin#any)) (Rows, [error](/builtin#error))

Query sends a query to the server and returns a Rows to read the results. Only errors encountered sending the query and initializing Rows will be returned. Err() on the returned Rows must be checked after the Rows is closed to determine if the query executed successfully. 

The returned Rows must be closed before the connection can be used again. It is safe to attempt to read from the returned Rows even if an error is returned. The error will be the available in rows.Err() after rows are closed. It is allowed to ignore the error returned from Query and handle it in Rows. 

It is possible for a call of FieldDescriptions on the returned Rows to return nil even if the Query call did not return an error. 

It is possible for a query to return one or more rows before encountering an error. In most cases the rows should be collected before processing rather than processed while receiving each row. This avoids the possibility of the application processing rows from a query that the server rejected. The CollectRows function is useful here. 

An implementor of QueryRewriter may be passed as the first element of args. It can rewrite the sql and change or replace args. For example, NamedArgs is QueryRewriter that implements named arguments. 

For extra control over how the query is executed, the types QueryExecMode, QueryResultFormats, and QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely needed. See the documentation for those types for details. 

Example ¶

This example uses Query without using any helpers to read the results. Normally CollectRows, ForEachRow, or another helper function should be used. 
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    if conn.PgConn().ParameterStatus("crdb_version") != "" {
    	// Skip test / example when running on CockroachDB. Since an example can't be skipped fake success instead.
    	fmt.Println(`Cheeseburger: $10
    Fries: $5
    Soft Drink: $3`)
    	return
    }
    
    // Setup example schema and data.
    _, err = conn.Exec(ctx, `
    create temporary table products (
    	id int primary key generated by default as identity,
    	name varchar(100) not null,
    	price int not null
    );
    
    insert into products (name, price) values
    	('Cheeseburger', 10),
    	('Double Cheeseburger', 14),
    	('Fries', 5),
    	('Soft Drink', 3);
    `)
    if err != nil {
    	fmt.Printf("Unable to setup example schema and data: %v", err)
    	return
    }
    
    rows, err := conn.Query(ctx, "select name, price from products where price < $1 order by price desc", 12)
    // It is unnecessary to check err. If an error occurred it will be returned by rows.Err() later. But in rare
    // cases it may be useful to detect the error as early as possible.
    if err != nil {
    	fmt.Printf("Query error: %v", err)
    	return
    }
    
    // Ensure rows is closed. It is safe to close rows multiple times.
    defer rows.Close()
    
    // Iterate through the result set
    for rows.Next() {
    	var name string
    	var price int32
    
    	err = rows.Scan(&name, &price)
    	if err != nil {
    		fmt.Printf("Scan error: %v", err)
    		return
    	}
    
    	fmt.Printf("%s: $%d\n", name, price)
    }
    
    // rows is closed automatically when rows.Next() returns false so it is not necessary to manually close rows.
    
    // The first error encountered by the original Query call, rows.Next or rows.Scan will be returned here.
    if rows.Err() != nil {
    	fmt.Printf("rows error: %v", rows.Err())
    	return
    }
    
    
    
    Output:
    
    Cheeseburger: $10
    Fries: $5
    Soft Drink: $3
    

####  func (*Conn) [QueryRow](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L928) ¶
    
    
    func (c *Conn) QueryRow(ctx [context](/context).[Context](/context#Context), sql [string](/builtin#string), args ...[any](/builtin#any)) Row

QueryRow is a convenience wrapper over Query. Any error that occurs while querying is deferred until calling Scan on the returned Row. That Row will error with ErrNoRows if no rows are returned. 

####  func (*Conn) [SendBatch](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L939) ¶
    
    
    func (c *Conn) SendBatch(ctx [context](/context).[Context](/context#Context), b *Batch) (br BatchResults)

SendBatch sends all queued queries to the server at once. All queries are run in an implicit transaction unless explicit transaction control statements are executed. The returned BatchResults must be closed before the connection is used again. 

Depending on the QueryExecMode, all queries may be prepared before any are executed. This means that creating a table and using it in a subsequent query in the same batch can fail. 

Example ¶
    
    
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()
    
    conn, err := pgx.Connect(ctx, os.Getenv("PGX_TEST_DATABASE"))
    if err != nil {
    	fmt.Printf("Unable to establish connection: %v", err)
    	return
    }
    
    batch := &pgx.Batch{}
    batch.Queue("select 1 + 1").QueryRow(func(row pgx.Row) error {
    	var n int32
    	err := row.Scan(&n)
    	if err != nil {
    		return err
    	}
    
    	fmt.Println(n)
    
    	return err
    })
    
    batch.Queue("select 1 + 2").QueryRow(func(row pgx.Row) error {
    	var n int32
    	err := row.Scan(&n)
    	if err != nil {
    		return err
    	}
    
    	fmt.Println(n)
    
    	return err
    })
    
    batch.Queue("select 2 + 3").QueryRow(func(row pgx.Row) error {
    	var n int32
    	err := row.Scan(&n)
    	if err != nil {
    		return err
    	}
    
    	fmt.Println(n)
    
    	return err
    })
    
    err = conn.SendBatch(ctx, batch).Close()
    if err != nil {
    	fmt.Printf("SendBatch error: %v", err)
    	return
    }
    
    
    
    Output:
    
    2
    3
    5
    

####  func (*Conn) [TypeMap](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L463) ¶
    
    
    func (c *Conn) TypeMap() *[pgtype](/github.com/jackc/pgx/v5@v5.8.0/pgtype).[Map](/github.com/jackc/pgx/v5@v5.8.0/pgtype#Map)

TypeMap returns the connection info used for this connection. 

####  func (*Conn) [WaitForNotification](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L413) ¶
    
    
    func (c *Conn) WaitForNotification(ctx [context](/context).[Context](/context#Context)) (*[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[Notification](/github.com/jackc/pgx/v5@v5.8.0/pgconn#Notification), [error](/builtin#error))

WaitForNotification waits for a PostgreSQL notification. It wraps the underlying pgconn notification system in a slightly more convenient form. 

####  type [ConnConfig](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L22) ¶
    
    
    type ConnConfig struct {
    	[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[Config](/github.com/jackc/pgx/v5@v5.8.0/pgconn#Config)
    
    	Tracer QueryTracer
    
    	// StatementCacheCapacity is maximum size of the statement cache used when executing a query with "cache_statement"
    	// query exec mode.
    	StatementCacheCapacity [int](/builtin#int)
    
    	// DescriptionCacheCapacity is the maximum size of the description cache used when executing a query with
    	// "cache_describe" query exec mode.
    	DescriptionCacheCapacity [int](/builtin#int)
    
    	// DefaultQueryExecMode controls the default mode for executing queries. By default pgx uses the extended protocol
    	// and automatically prepares and caches prepared statements. However, this may be incompatible with proxies such as
    	// PGBouncer. In this case it may be preferable to use QueryExecModeExec or QueryExecModeSimpleProtocol. The same
    	// functionality can be controlled on a per query basis by passing a QueryExecMode as the first query argument.
    	DefaultQueryExecMode QueryExecMode
    	// contains filtered or unexported fields
    }

ConnConfig contains all the options used to establish a connection. It must be created by ParseConfig and then it can be modified. A manually initialized ConnConfig will cause ConnectConfig to panic. 

####  func [ParseConfig](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L236) ¶
    
    
    func ParseConfig(connString [string](/builtin#string)) (*ConnConfig, [error](/builtin#error))

ParseConfig creates a ConnConfig from a connection string. ParseConfig handles all options that [pgconn.ParseConfig](/github.com/jackc/pgx/v5@v5.8.0/pgconn#ParseConfig) does. In addition, it accepts the following options: 

  * default_query_exec_mode. Possible values: "cache_statement", "cache_describe", "describe_exec", "exec", and "simple_protocol". See QueryExecMode constant documentation for the meaning of these values. Default: "cache_statement". 

  * statement_cache_capacity. The maximum size of the statement cache used when executing a query with "cache_statement" query exec mode. Default: 512. 

  * description_cache_capacity. The maximum size of the description cache used when executing a query with "cache_describe" query exec mode. Default: 512. 




####  func [ParseConfigWithOptions](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L165) ¶ added in v5.1.0
    
    
    func ParseConfigWithOptions(connString [string](/builtin#string), options ParseConfigOptions) (*ConnConfig, [error](/builtin#error))

ParseConfigWithOptions behaves exactly as ParseConfig does with the addition of options. At the present options is only used to provide a GetSSLPassword function. 

####  func (*ConnConfig) [ConnString](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L63) ¶
    
    
    func (cc *ConnConfig) ConnString() [string](/builtin#string)

ConnString returns the connection string as parsed by pgx.ParseConfig into pgx.ConnConfig. 

####  func (*ConnConfig) [Copy](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L55) ¶
    
    
    func (cc *ConnConfig) Copy() *ConnConfig

Copy returns a deep copy of the config that is safe to use and modify. The only exception is the tls.Config: according to the tls.Config docs it must not be modified after creation. 

####  type [ConnectTracer](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L92) ¶
    
    
    type ConnectTracer interface {
    	// TraceConnectStart is called at the beginning of Connect and ConnectConfig calls. The returned context is used for
    	// the rest of the call and will be passed to TraceConnectEnd.
    	TraceConnectStart(ctx [context](/context).[Context](/context#Context), data TraceConnectStartData) [context](/context).[Context](/context#Context)
    
    	TraceConnectEnd(ctx [context](/context).[Context](/context#Context), data TraceConnectEndData)
    }

ConnectTracer traces Connect and ConnectConfig. 

####  type [CopyFromSource](https://github.com/jackc/pgx/blob/v5.8.0/copy_from.go#L95) ¶
    
    
    type CopyFromSource interface {
    	// Next returns true if there is another row and makes the next row data
    	// available to Values(). When there are no more rows available or an error
    	// has occurred it returns false.
    	Next() [bool](/builtin#bool)
    
    	// Values returns the values for the current row.
    	Values() ([][any](/builtin#any), [error](/builtin#error))
    
    	// Err returns any error that has been encountered by the CopyFromSource. If
    	// this is not nil *Conn.CopyFrom will abort the copy.
    	Err() [error](/builtin#error)
    }

CopyFromSource is the interface used by *Conn.CopyFrom as the source for copy data. 

####  func [CopyFromFunc](https://github.com/jackc/pgx/blob/v5.8.0/copy_from.go#L70) ¶ added in v5.5.1
    
    
    func CopyFromFunc(nxtf func() (row [][any](/builtin#any), err [error](/builtin#error))) CopyFromSource

CopyFromFunc returns a CopyFromSource interface that relies on nxtf for values. nxtf returns rows until it either signals an 'end of data' by returning row=nil and err=nil, or it returns an error. If nxtf returns an error, the copy is aborted. 

####  func [CopyFromRows](https://github.com/jackc/pgx/blob/v5.8.0/copy_from.go#L15) ¶
    
    
    func CopyFromRows(rows [][][any](/builtin#any)) CopyFromSource

CopyFromRows returns a CopyFromSource interface over the provided rows slice making it usable by *Conn.CopyFrom. 

####  func [CopyFromSlice](https://github.com/jackc/pgx/blob/v5.8.0/copy_from.go#L39) ¶
    
    
    func CopyFromSlice(length [int](/builtin#int), next func([int](/builtin#int)) ([][any](/builtin#any), [error](/builtin#error))) CopyFromSource

CopyFromSlice returns a CopyFromSource interface over a dynamic func making it usable by *Conn.CopyFrom. 

####  type [CopyFromTracer](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L54) ¶
    
    
    type CopyFromTracer interface {
    	// TraceCopyFromStart is called at the beginning of CopyFrom calls. The returned context is used for the
    	// rest of the call and will be passed to TraceCopyFromEnd.
    	TraceCopyFromStart(ctx [context](/context).[Context](/context#Context), conn *Conn, data TraceCopyFromStartData) [context](/context).[Context](/context#Context)
    
    	TraceCopyFromEnd(ctx [context](/context).[Context](/context#Context), conn *Conn, data TraceCopyFromEndData)
    }

CopyFromTracer traces CopyFrom. 

####  type [ExtendedQueryBuilder](https://github.com/jackc/pgx/blob/v5.8.0/extended_query_builder.go#L12) ¶
    
    
    type ExtendedQueryBuilder struct {
    	ParamValues [][][byte](/builtin#byte)
    
    	ParamFormats  [][int16](/builtin#int16)
    	ResultFormats [][int16](/builtin#int16)
    	// contains filtered or unexported fields
    }

ExtendedQueryBuilder is used to choose the parameter formats, to format the parameters and to choose the result formats for an extended query. 

####  func (*ExtendedQueryBuilder) [Build](https://github.com/jackc/pgx/blob/v5.8.0/extended_query_builder.go#L21) ¶
    
    
    func (eqb *ExtendedQueryBuilder) Build(m *[pgtype](/github.com/jackc/pgx/v5@v5.8.0/pgtype).[Map](/github.com/jackc/pgx/v5@v5.8.0/pgtype#Map), sd *[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[StatementDescription](/github.com/jackc/pgx/v5@v5.8.0/pgconn#StatementDescription), args [][any](/builtin#any)) [error](/builtin#error)

Build sets ParamValues, ParamFormats, and ResultFormats for use with *PgConn.ExecParams or *PgConn.ExecPrepared. If sd is nil then QueryExecModeExec behavior will be used. 

####  type [Identifier](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L93) ¶
    
    
    type Identifier [][string](/builtin#string)

Identifier a PostgreSQL identifier or name. Identifiers can be composed of multiple parts such as ["schema", "table"] or ["table", "column"]. 

####  func (Identifier) [Sanitize](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L96) ¶
    
    
    func (ident Identifier) Sanitize() [string](/builtin#string)

Sanitize returns a sanitized string safe for SQL interpolation. 

####  type [LargeObject](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L70) ¶
    
    
    type LargeObject struct {
    	// contains filtered or unexported fields
    }

A LargeObject is a large object stored on the server. It is only valid within the transaction that it was initialized in. It uses the context it was initialized with for all operations. It implements these interfaces: 
    
    
    io.Writer
    io.Reader
    io.Seeker
    io.Closer
    

####  func (*LargeObject) [Close](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L158) ¶
    
    
    func (o *LargeObject) Close() [error](/builtin#error)

Close the large object descriptor. 

####  func (*LargeObject) [Read](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L110) ¶
    
    
    func (o *LargeObject) Read(p [][byte](/builtin#byte)) ([int](/builtin#int), [error](/builtin#error))

Read reads up to len(p) bytes into p returning the number of bytes read. 

####  func (*LargeObject) [Seek](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L140) ¶
    
    
    func (o *LargeObject) Seek(offset [int64](/builtin#int64), whence [int](/builtin#int)) (n [int64](/builtin#int64), err [error](/builtin#error))

Seek moves the current location pointer to the new location specified by offset. 

####  func (*LargeObject) [Tell](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L146) ¶
    
    
    func (o *LargeObject) Tell() (n [int64](/builtin#int64), err [error](/builtin#error))

Tell returns the current read or write location of the large object descriptor. 

####  func (*LargeObject) [Truncate](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L152) ¶
    
    
    func (o *LargeObject) Truncate(size [int64](/builtin#int64)) (err [error](/builtin#error))

Truncate the large object to size. 

####  func (*LargeObject) [Write](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L77) ¶
    
    
    func (o *LargeObject) Write(p [][byte](/builtin#byte)) ([int](/builtin#int), [error](/builtin#error))

Write writes p to the large object and returns the number of bytes written and an error if not all of p was written. 

####  type [LargeObjectMode](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L24) ¶
    
    
    type LargeObjectMode [int32](/builtin#int32)
    
    
    const (
    	LargeObjectModeWrite LargeObjectMode = 0x20000
    	LargeObjectModeRead  LargeObjectMode = 0x40000
    )

####  type [LargeObjects](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L20) ¶
    
    
    type LargeObjects struct {
    	// contains filtered or unexported fields
    }

LargeObjects is a structure used to access the large objects API. It is only valid within the transaction where it was created. 

For more details see: <http://www.postgresql.org/docs/current/static/largeobjects.html>

####  func (*LargeObjects) [Create](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L32) ¶
    
    
    func (o *LargeObjects) Create(ctx [context](/context).[Context](/context#Context), oid [uint32](/builtin#uint32)) ([uint32](/builtin#uint32), [error](/builtin#error))

Create creates a new large object. If oid is zero, the server assigns an unused OID. 

####  func (*LargeObjects) [Open](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L39) ¶
    
    
    func (o *LargeObjects) Open(ctx [context](/context).[Context](/context#Context), oid [uint32](/builtin#uint32), mode LargeObjectMode) (*LargeObject, [error](/builtin#error))

Open opens an existing large object with the given mode. ctx will also be used for all operations on the opened large object. 

####  func (*LargeObjects) [Unlink](https://github.com/jackc/pgx/blob/v5.8.0/large_objects.go#L49) ¶
    
    
    func (o *LargeObjects) Unlink(ctx [context](/context).[Context](/context#Context), oid [uint32](/builtin#uint32)) [error](/builtin#error)

Unlink removes a large object from the database. 

####  type [NamedArgs](https://github.com/jackc/pgx/blob/v5.8.0/named_args.go#L21) ¶
    
    
    type NamedArgs map[[string](/builtin#string)][any](/builtin#any)

NamedArgs can be used as the first argument to a query method. It will replace every '@' named placeholder with a '$' ordinal placeholder and construct the appropriate arguments. 

For example, the following two queries are equivalent: 
    
    
    conn.Query(ctx, "select * from widgets where foo = @foo and bar = @bar", pgx.NamedArgs{"foo": 1, "bar": 2})
    conn.Query(ctx, "select * from widgets where foo = $1 and bar = $2", 1, 2)
    

Named placeholders are case sensitive and must start with a letter or underscore. Subsequent characters can be letters, numbers, or underscores. 

####  func (NamedArgs) [RewriteQuery](https://github.com/jackc/pgx/blob/v5.8.0/named_args.go#L24) ¶
    
    
    func (na NamedArgs) RewriteQuery(ctx [context](/context).[Context](/context#Context), conn *Conn, sql [string](/builtin#string), args [][any](/builtin#any)) (newSQL [string](/builtin#string), newArgs [][any](/builtin#any), err [error](/builtin#error))

RewriteQuery implements the QueryRewriter interface. 

####  type [ParseConfigOptions](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L48) ¶ added in v5.1.0
    
    
    type ParseConfigOptions struct {
    	[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[ParseConfigOptions](/github.com/jackc/pgx/v5@v5.8.0/pgconn#ParseConfigOptions)
    }

ParseConfigOptions contains options that control how a config is built such as getsslpassword. 

####  type [PrepareTracer](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L73) ¶
    
    
    type PrepareTracer interface {
    	// TracePrepareStart is called at the beginning of Prepare calls. The returned context is used for the
    	// rest of the call and will be passed to TracePrepareEnd.
    	TracePrepareStart(ctx [context](/context).[Context](/context#Context), conn *Conn, data TracePrepareStartData) [context](/context).[Context](/context#Context)
    
    	TracePrepareEnd(ctx [context](/context).[Context](/context#Context), conn *Conn, data TracePrepareEndData)
    }

PrepareTracer traces Prepare. 

####  type [QueryExecMode](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L641) ¶
    
    
    type QueryExecMode [int32](/builtin#int32)
    
    
    const (
    
    	// Automatically prepare and cache statements. This uses the extended protocol. Queries are executed in a single round
    	// trip after the statement is cached. This is the default. If the database schema is modified or the search_path is
    	// changed after a statement is cached then the first execution of a previously cached query may fail. e.g. If the
    	// number of columns returned by a "SELECT *" changes or the type of a column is changed.
    	QueryExecModeCacheStatement QueryExecMode
    
    	// Cache statement descriptions (i.e. argument and result types) and assume they do not change. This uses the extended
    	// protocol. Queries are executed in a single round trip after the description is cached. If the database schema is
    	// modified or the search_path is changed after a statement is cached then the first execution of a previously cached
    	// query may fail. e.g. If the number of columns returned by a "SELECT *" changes or the type of a column is changed.
    	QueryExecModeCacheDescribe
    
    	// Get the statement description on every execution. This uses the extended protocol. Queries require two round trips
    	// to execute. It does not use named prepared statements. But it does use the unnamed prepared statement to get the
    	// statement description on the first round trip and then uses it to execute the query on the second round trip. This
    	// may cause problems with connection poolers that switch the underlying connection between round trips. It is safe
    	// even when the database schema is modified concurrently.
    	QueryExecModeDescribeExec
    
    	// Assume the PostgreSQL query parameter types based on the Go type of the arguments. This uses the extended protocol
    	// with text formatted parameters and results. Queries are executed in a single round trip. Type mappings can be
    	// registered with pgtype.Map.RegisterDefaultPgType. Queries will be rejected that have arguments that are
    	// unregistered or ambiguous. e.g. A map[string]string may have the PostgreSQL type json or hstore. Modes that know
    	// the PostgreSQL type can use a map[string]string directly as an argument. This mode cannot.
    	//
    	// On rare occasions user defined types may behave differently when encoded in the text format instead of the binary
    	// format. For example, this could happen if a "type RomanNumeral int32" implements fmt.Stringer to format integers as
    	// Roman numerals (e.g. 7 is VII). The binary format would properly encode the integer 7 as the binary value for 7.
    	// But the text format would encode the integer 7 as the string "VII". As QueryExecModeExec uses the text format, it
    	// is possible that changing query mode from another mode to QueryExecModeExec could change the behavior of the query.
    	// This should not occur with types pgx supports directly and can be avoided by registering the types with
    	// pgtype.Map.RegisterDefaultPgType and implementing the appropriate type interfaces. In the cas of RomanNumeral, it
    	// should implement pgtype.Int64Valuer.
    	QueryExecModeExec
    
    	// Use the simple protocol. Assume the PostgreSQL query parameter types based on the Go type of the arguments. This is
    	// especially significant for []byte values. []byte values are encoded as PostgreSQL bytea. string must be used
    	// instead for text type values including json and jsonb. Type mappings can be registered with
    	// pgtype.Map.RegisterDefaultPgType. Queries will be rejected that have arguments that are unregistered or ambiguous.
    	// e.g. A map[string]string may have the PostgreSQL type json or hstore. Modes that know the PostgreSQL type can use a
    	// map[string]string directly as an argument. This mode cannot. Queries are executed in a single round trip.
    	//
    	// QueryExecModeSimpleProtocol should have the user application visible behavior as QueryExecModeExec. This includes
    	// the warning regarding differences in text format and binary format encoding with user defined types. There may be
    	// other minor exceptions such as behavior when multiple result returning queries are erroneously sent in a single
    	// string.
    	//
    	// QueryExecModeSimpleProtocol uses client side parameter interpolation. All values are quoted and escaped. Prefer
    	// QueryExecModeExec over QueryExecModeSimpleProtocol whenever possible. In general QueryExecModeSimpleProtocol should
    	// only be used if connecting to a proxy server, connection pool server, or non-PostgreSQL server that does not
    	// support the extended protocol.
    	QueryExecModeSimpleProtocol
    )

####  func (QueryExecMode) [String](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L700) ¶
    
    
    func (m QueryExecMode) String() [string](/builtin#string)

####  type [QueryResultFormats](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L718) ¶
    
    
    type QueryResultFormats [][int16](/builtin#int16)

QueryResultFormats controls the result format (text=0, binary=1) of a query by result column position. 

####  type [QueryResultFormatsByOID](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L721) ¶
    
    
    type QueryResultFormatsByOID map[[uint32](/builtin#uint32)][int16](/builtin#int16)

QueryResultFormatsByOID controls the result format (text=0, binary=1) of a query by the result column OID. 

####  type [QueryRewriter](https://github.com/jackc/pgx/blob/v5.8.0/conn.go#L724) ¶
    
    
    type QueryRewriter interface {
    	RewriteQuery(ctx [context](/context).[Context](/context#Context), conn *Conn, sql [string](/builtin#string), args [][any](/builtin#any)) (newSQL [string](/builtin#string), newArgs [][any](/builtin#any), err [error](/builtin#error))
    }

QueryRewriter rewrites a query when used as the first arguments to a query method. 

####  type [QueryTracer](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L10) ¶
    
    
    type QueryTracer interface {
    	// TraceQueryStart is called at the beginning of Query, QueryRow, and Exec calls. The returned context is used for the
    	// rest of the call and will be passed to TraceQueryEnd.
    	TraceQueryStart(ctx [context](/context).[Context](/context#Context), conn *Conn, data TraceQueryStartData) [context](/context).[Context](/context#Context)
    
    	TraceQueryEnd(ctx [context](/context).[Context](/context#Context), conn *Conn, data TraceQueryEndData)
    }

QueryTracer traces Query, QueryRow, and Exec. 

####  type [QueuedQuery](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L12) ¶
    
    
    type QueuedQuery struct {
    	SQL       [string](/builtin#string)
    	Arguments [][any](/builtin#any)
    	Fn        batchItemFunc
    	// contains filtered or unexported fields
    }

QueuedQuery is a query that has been queued for execution via a Batch. 

####  func (*QueuedQuery) [Exec](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L50) ¶
    
    
    func (qq *QueuedQuery) Exec(fn func(ct [pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag)) [error](/builtin#error))

Exec sets fn to be called when the response to qq is received. 

Note: for simple batch insert uses where it is not required to handle each potential error individually, it's sufficient to not set any callbacks, and just handle the return value of BatchResults.Close. 

####  func (*QueuedQuery) [Query](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L22) ¶
    
    
    func (qq *QueuedQuery) Query(fn func(rows Rows) [error](/builtin#error))

Query sets fn to be called when the response to qq is received. 

####  func (*QueuedQuery) [QueryRow](https://github.com/jackc/pgx/blob/v5.8.0/batch.go#L38) ¶
    
    
    func (qq *QueuedQuery) QueryRow(fn func(row Row) [error](/builtin#error))

Query sets fn to be called when the response to qq is received. 

####  type [Row](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L79) ¶
    
    
    type Row interface {
    	// Scan works the same as Rows. with the following exceptions. If no
    	// rows were found it returns ErrNoRows. If multiple rows are returned it
    	// ignores all but the first.
    	Scan(dest ...[any](/builtin#any)) [error](/builtin#error)
    }

Row is a convenience wrapper over Rows that is returned by QueryRow. 

Row is an interface instead of a struct to allow tests to mock QueryRow. However, adding a method to an interface is technically a breaking change. Because of this the Row interface is partially excluded from semantic version requirements. Methods will not be removed or changed, but new methods may be added. 

####  type [RowScanner](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L87) ¶
    
    
    type RowScanner interface {
    	// ScanRows scans the row.
    	ScanRow(rows Rows) [error](/builtin#error)
    }

RowScanner scans an entire row at a time into the RowScanner. 

####  type [RowToFunc](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L432) ¶
    
    
    type RowToFunc[T [any](/builtin#any)] func(row CollectableRow) (T, [error](/builtin#error))

RowToFunc is a function that scans or otherwise converts row to a T. 

####  type [Rows](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L27) ¶
    
    
    type Rows interface {
    	// Close closes the rows, making the connection ready for use again. It is safe
    	// to call Close after rows is already closed.
    	Close()
    
    	// Err returns any error that occurred while executing a query or reading its results. Err must be called after the
    	// Rows is closed (either by calling Close or by Next returning false) to check if the query was successful. If it is
    	// called before the Rows is closed it may return nil even if the query failed on the server.
    	Err() [error](/builtin#error)
    
    	// CommandTag returns the command tag from this query. It is only available after Rows is closed.
    	CommandTag() [pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag)
    
    	// FieldDescriptions returns the field descriptions of the columns. It may return nil. In particular this can occur
    	// when there was an error executing the query.
    	FieldDescriptions() [][pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[FieldDescription](/github.com/jackc/pgx/v5@v5.8.0/pgconn#FieldDescription)
    
    	// Next prepares the next row for reading. It returns true if there is another row and false if no more rows are
    	// available or a fatal error has occurred. It automatically closes rows upon returning false (whether due to all rows
    	// having been read or due to an error).
    	//
    	// Callers should check rows.Err() after rows.Next() returns false to detect whether result-set reading ended
    	// prematurely due to an error. See Conn.Query for details.
    	//
    	// For simpler error handling, consider using the higher-level pgx v5 CollectRows() and ForEachRow() helpers instead.
    	Next() [bool](/builtin#bool)
    
    	// Scan reads the values from the current row into dest values positionally. dest can include pointers to core types,
    	// values implementing the Scanner interface, and nil. nil will skip the value entirely. It is an error to call Scan
    	// without first calling Next() and checking that it returned true. Rows is automatically closed upon error.
    	Scan(dest ...[any](/builtin#any)) [error](/builtin#error)
    
    	// Values returns the decoded row values. As with Scan(), it is an error to
    	// call Values without first calling Next() and checking that it returned
    	// true.
    	Values() ([][any](/builtin#any), [error](/builtin#error))
    
    	// RawValues returns the unparsed bytes of the row values. The returned data is only valid until the next Next
    	// call or the Rows is closed.
    	RawValues() [][][byte](/builtin#byte)
    
    	// Conn returns the underlying *Conn on which the query was executed. This may return nil if Rows did not come from a
    	// *Conn (e.g. if it was created by RowsFromResultReader)
    	Conn() *Conn
    }

Rows is the result set returned from *Conn.Query. Rows must be closed before the *Conn can be used again. Rows are closed by explicitly calling Close(), calling Next() until it returns false, or when a fatal error occurs. 

Once a Rows is closed the only methods that may be called are Close(), Err(), and CommandTag(). 

Rows is an interface instead of a struct to allow tests to mock Query. However, adding a method to an interface is technically a breaking change. Because of this the Rows interface is partially excluded from semantic version requirements. Methods will not be removed or changed, but new methods may be added. 

####  func [RowsFromResultReader](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L391) ¶
    
    
    func RowsFromResultReader(typeMap *[pgtype](/github.com/jackc/pgx/v5@v5.8.0/pgtype).[Map](/github.com/jackc/pgx/v5@v5.8.0/pgtype#Map), resultReader *[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[ResultReader](/github.com/jackc/pgx/v5@v5.8.0/pgconn#ResultReader)) Rows

RowsFromResultReader returns a Rows that will read from values resultReader and decode with typeMap. It can be used to read from the lower level pgconn interface. 

####  type [ScanArgError](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L343) ¶
    
    
    type ScanArgError struct {
    	ColumnIndex [int](/builtin#int)
    	FieldName   [string](/builtin#string)
    	Err         [error](/builtin#error)
    }

####  func (ScanArgError) [Error](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L349) ¶
    
    
    func (e ScanArgError) Error() [string](/builtin#string)

####  func (ScanArgError) [Unwrap](https://github.com/jackc/pgx/blob/v5.8.0/rows.go#L357) ¶
    
    
    func (e ScanArgError) Unwrap() [error](/builtin#error)

####  type [StrictNamedArgs](https://github.com/jackc/pgx/blob/v5.8.0/named_args.go#L30) ¶ added in v5.6.0
    
    
    type StrictNamedArgs map[[string](/builtin#string)][any](/builtin#any)

StrictNamedArgs can be used in the same way as NamedArgs, but provided arguments are also checked to include all named arguments that the sql query uses, and no extra arguments. 

####  func (StrictNamedArgs) [RewriteQuery](https://github.com/jackc/pgx/blob/v5.8.0/named_args.go#L33) ¶ added in v5.6.0
    
    
    func (sna StrictNamedArgs) RewriteQuery(ctx [context](/context).[Context](/context#Context), conn *Conn, sql [string](/builtin#string), args [][any](/builtin#any)) (newSQL [string](/builtin#string), newArgs [][any](/builtin#any), err [error](/builtin#error))

RewriteQuery implements the QueryRewriter interface. 

####  type [TraceBatchEndData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L49) ¶
    
    
    type TraceBatchEndData struct {
    	Err [error](/builtin#error)
    }

####  type [TraceBatchQueryData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L42) ¶
    
    
    type TraceBatchQueryData struct {
    	SQL        [string](/builtin#string)
    	Args       [][any](/builtin#any)
    	CommandTag [pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag)
    	Err        [error](/builtin#error)
    }

####  type [TraceBatchStartData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L38) ¶
    
    
    type TraceBatchStartData struct {
    	Batch *Batch
    }

####  type [TraceConnectEndData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L104) ¶
    
    
    type TraceConnectEndData struct {
    	Conn *Conn
    	Err  [error](/builtin#error)
    }

####  type [TraceConnectStartData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L100) ¶
    
    
    type TraceConnectStartData struct {
    	ConnConfig *ConnConfig
    }

####  type [TraceCopyFromEndData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L67) ¶
    
    
    type TraceCopyFromEndData struct {
    	CommandTag [pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag)
    	Err        [error](/builtin#error)
    }

####  type [TraceCopyFromStartData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L62) ¶
    
    
    type TraceCopyFromStartData struct {
    	TableName   Identifier
    	ColumnNames [][string](/builtin#string)
    }

####  type [TracePrepareEndData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L86) ¶
    
    
    type TracePrepareEndData struct {
    	AlreadyPrepared [bool](/builtin#bool)
    	Err             [error](/builtin#error)
    }

####  type [TracePrepareStartData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L81) ¶
    
    
    type TracePrepareStartData struct {
    	Name [string](/builtin#string)
    	SQL  [string](/builtin#string)
    }

####  type [TraceQueryEndData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L23) ¶
    
    
    type TraceQueryEndData struct {
    	CommandTag [pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag)
    	Err        [error](/builtin#error)
    }

####  type [TraceQueryStartData](https://github.com/jackc/pgx/blob/v5.8.0/tracer.go#L18) ¶
    
    
    type TraceQueryStartData struct {
    	SQL  [string](/builtin#string)
    	Args [][any](/builtin#any)
    }

####  type [Tx](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L122) ¶
    
    
    type Tx interface {
    	// Begin starts a pseudo nested transaction.
    	Begin(ctx [context](/context).[Context](/context#Context)) (Tx, [error](/builtin#error))
    
    	// Commit commits the transaction if this is a real transaction or releases the savepoint if this is a pseudo nested
    	// transaction. Commit will return an error where errors.Is(ErrTxClosed) is true if the Tx is already closed, but is
    	// otherwise safe to call multiple times. If the commit fails with a rollback status (e.g. the transaction was already
    	// in a broken state) then an error where errors.Is(ErrTxCommitRollback) is true will be returned.
    	Commit(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)
    
    	// Rollback rolls back the transaction if this is a real transaction or rolls back to the savepoint if this is a
    	// pseudo nested transaction. Rollback will return an error where errors.Is(ErrTxClosed) is true if the Tx is already
    	// closed, but is otherwise safe to call multiple times. Hence, a defer tx.Rollback() is safe even if tx.Commit() will
    	// be called first in a non-error condition. Any other failure of a real transaction will result in the connection
    	// being closed.
    	Rollback(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)
    
    	CopyFrom(ctx [context](/context).[Context](/context#Context), tableName Identifier, columnNames [][string](/builtin#string), rowSrc CopyFromSource) ([int64](/builtin#int64), [error](/builtin#error))
    	SendBatch(ctx [context](/context).[Context](/context#Context), b *Batch) BatchResults
    	LargeObjects() LargeObjects
    
    	Prepare(ctx [context](/context).[Context](/context#Context), name, sql [string](/builtin#string)) (*[pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[StatementDescription](/github.com/jackc/pgx/v5@v5.8.0/pgconn#StatementDescription), [error](/builtin#error))
    
    	Exec(ctx [context](/context).[Context](/context#Context), sql [string](/builtin#string), arguments ...[any](/builtin#any)) (commandTag [pgconn](/github.com/jackc/pgx/v5@v5.8.0/pgconn).[CommandTag](/github.com/jackc/pgx/v5@v5.8.0/pgconn#CommandTag), err [error](/builtin#error))
    	Query(ctx [context](/context).[Context](/context#Context), sql [string](/builtin#string), args ...[any](/builtin#any)) (Rows, [error](/builtin#error))
    	QueryRow(ctx [context](/context).[Context](/context#Context), sql [string](/builtin#string), args ...[any](/builtin#any)) Row
    
    	// Conn returns the underlying *Conn that on which this transaction is executing.
    	Conn() *Conn
    }

Tx represents a database transaction. 

Tx is an interface instead of a struct to enable connection pools to be implemented without relying on internal pgx state, to support pseudo-nested transactions with savepoints, and to allow tests to mock transactions. However, adding a method to an interface is technically a breaking change. If new methods are added to Conn it may be desirable to add them to Tx as well. Because of this the Tx interface is partially excluded from semantic version requirements. Methods will not be removed or changed, but new methods may be added. 

####  type [TxAccessMode](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L24) ¶
    
    
    type TxAccessMode [string](/builtin#string)

TxAccessMode is the transaction access mode (read write or read only) 
    
    
    const (
    	ReadWrite TxAccessMode = "read write"
    	ReadOnly  TxAccessMode = "read only"
    )

Transaction access modes 

####  type [TxDeferrableMode](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L33) ¶
    
    
    type TxDeferrableMode [string](/builtin#string)

TxDeferrableMode is the transaction deferrable mode (deferrable or not deferrable) 
    
    
    const (
    	Deferrable    TxDeferrableMode = "deferrable"
    	NotDeferrable TxDeferrableMode = "not deferrable"
    )

Transaction deferrable modes 

####  type [TxIsoLevel](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L13) ¶
    
    
    type TxIsoLevel [string](/builtin#string)

TxIsoLevel is the transaction isolation level (serializable, repeatable read, read committed or read uncommitted) 
    
    
    const (
    	Serializable    TxIsoLevel = "serializable"
    	RepeatableRead  TxIsoLevel = "repeatable read"
    	ReadCommitted   TxIsoLevel = "read committed"
    	ReadUncommitted TxIsoLevel = "read uncommitted"
    )

Transaction isolation levels 

####  type [TxOptions](https://github.com/jackc/pgx/blob/v5.8.0/tx.go#L42) ¶
    
    
    type TxOptions struct {
    	IsoLevel       TxIsoLevel
    	AccessMode     TxAccessMode
    	DeferrableMode TxDeferrableMode
    
    	// BeginQuery is the SQL query that will be executed to begin the transaction. This allows using non-standard syntax
    	// such as BEGIN PRIORITY HIGH with CockroachDB. If set this will override the other settings.
    	BeginQuery [string](/builtin#string)
    	// CommitQuery is the SQL query that will be executed to commit the transaction.
    	CommitQuery [string](/builtin#string)
    }

TxOptions are transaction modes within a transaction block 
