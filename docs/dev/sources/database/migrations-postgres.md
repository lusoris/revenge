# golang-migrate PostgreSQL

> Source: https://pkg.go.dev/github.com/golang-migrate/migrate/v4/database/postgres
> Fetched: 2026-01-30T23:50:06.687032+00:00
> Content-Hash: 44811f45ac85ae98
> Type: html

---

Index

¶

Variables

func WithInstance(instance *sql.DB, config *Config) (database.Driver, error)

type Config

type Postgres

func WithConnection(ctx context.Context, conn *sql.Conn, config *Config) (*Postgres, error)

func (p *Postgres) Close() error

func (p *Postgres) Drop() (err error)

func (p *Postgres) Lock() error

func (p *Postgres) Open(url string) (database.Driver, error)

func (p *Postgres) Run(migration io.Reader) error

func (p *Postgres) SetVersion(version int, dirty bool) error

func (p *Postgres) Unlock() error

func (p *Postgres) Version() (version int, dirty bool, err error)

Constants

¶

This section is empty.

Variables

¶

View Source

var (

DefaultMigrationsTable       = "schema_migrations"

DefaultMultiStatementMaxSize = 10 * 1 << 20

// 10 MB

)

View Source

var (

ErrNilConfig      =

fmt

.

Errorf

("no config")

ErrNoDatabaseName =

fmt

.

Errorf

("no database name")

ErrNoSchema       =

fmt

.

Errorf

("no schema")

ErrDatabaseDirty  =

fmt

.

Errorf

("database is dirty")

)

Functions

¶

func

WithInstance

¶

func WithInstance(instance *

sql

.

DB

, config *

Config

) (

database

.

Driver

,

error

)

Types

¶

type

Config

¶

type Config struct {

MigrationsTable

string

MigrationsTableQuoted

bool

MultiStatementEnabled

bool

DatabaseName

string

SchemaName

string

StatementTimeout

time

.

Duration

MultiStatementMaxSize

int

// contains filtered or unexported fields

}

type

Postgres

¶

type Postgres struct {

// contains filtered or unexported fields

}

func

WithConnection

¶

added in

v4.15.2

func WithConnection(ctx

context

.

Context

, conn *

sql

.

Conn

, config *

Config

) (*

Postgres

,

error

)

func (*Postgres)

Close

¶

func (p *

Postgres

) Close()

error

func (*Postgres)

Drop

¶

func (p *

Postgres

) Drop() (err

error

)

func (*Postgres)

Lock

¶

func (p *

Postgres

) Lock()

error

https://www.postgresql.org/docs/9.6/static/explicit-locking.html#ADVISORY-LOCKS

func (*Postgres)

Open

¶

func (p *

Postgres

) Open(url

string

) (

database

.

Driver

,

error

)

func (*Postgres)

Run

¶

func (p *

Postgres

) Run(migration

io

.

Reader

)

error

func (*Postgres)

SetVersion

¶

func (p *

Postgres

) SetVersion(version

int

, dirty

bool

)

error

func (*Postgres)

Unlock

¶

func (p *

Postgres

) Unlock()

error

func (*Postgres)

Version

¶

func (p *

Postgres

) Version() (version

int

, dirty

bool

, err

error

)