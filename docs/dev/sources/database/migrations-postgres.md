# golang-migrate PostgreSQL

> Source: https://pkg.go.dev/github.com/golang-migrate/migrate/v4/database/postgres
> Fetched: 2026-02-01T11:43:28.392621+00:00
> Content-Hash: cdb024435ad8d454
> Type: html

---

### Index ¶

- Variables
- func WithInstance(instance *sql.DB, config*Config) (database.Driver, error)
- type Config
- type Postgres
-     * func WithConnection(ctx context.Context, conn *sql.Conn, config *Config) (*Postgres, error)
-     * func (p *Postgres) Close() error
  - func (p *Postgres) Drop() (err error)
  - func (p *Postgres) Lock() error
  - func (p *Postgres) Open(url string) (database.Driver, error)
  - func (p *Postgres) Run(migration io.Reader) error
  - func (p *Postgres) SetVersion(version int, dirty bool) error
  - func (p *Postgres) Unlock() error
  - func (p *Postgres) Version() (version int, dirty bool, err error)

### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L30)

    var (
     DefaultMigrationsTable       = "schema_migrations"
     DefaultMultiStatementMaxSize = 10 * 1 << 20 // 10 MB
    )

[View Source](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L37)

    var (
     ErrNilConfig      = [fmt](/fmt).[Errorf](/fmt#Errorf)("no config")
     ErrNoDatabaseName = [fmt](/fmt).[Errorf](/fmt#Errorf)("no database name")
     ErrNoSchema       = [fmt](/fmt).[Errorf](/fmt#Errorf)("no schema")
     ErrDatabaseDirty  = [fmt](/fmt).[Errorf](/fmt#Errorf)("database is dirty")
    )

### Functions ¶

#### func [WithInstance](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L132) ¶

    func WithInstance(instance *[sql](/database/sql).[DB](/database/sql#DB), config *Config) ([database](/github.com/golang-migrate/migrate/v4@v4.19.1/database).[Driver](/github.com/golang-migrate/migrate/v4@v4.19.1/database#Driver), [error](/builtin#error))

### Types ¶

#### type [Config](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L44) ¶

    type Config struct {
     MigrationsTable       [string](/builtin#string)
     MigrationsTableQuoted [bool](/builtin#bool)
     MultiStatementEnabled [bool](/builtin#bool)
     DatabaseName          [string](/builtin#string)
     SchemaName            [string](/builtin#string)
    
     StatementTimeout      [time](/time).[Duration](/time#Duration)
     MultiStatementMaxSize [int](/builtin#int)
     // contains filtered or unexported fields
    }

#### type [Postgres](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L56) ¶

    type Postgres struct {
     // contains filtered or unexported fields
    }

#### func [WithConnection](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L66) ¶ added in v4.15.2

    func WithConnection(ctx [context](/context).[Context](/context#Context), conn *[sql](/database/sql).[Conn](/database/sql#Conn), config *Config) (*Postgres, [error](/builtin#error))

#### func (*Postgres) [Close](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L219) ¶

    func (p *Postgres) Close() [error](/builtin#error)

#### func (*Postgres) [Drop](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L409) ¶

    func (p *Postgres) Drop() (err [error](/builtin#error))

#### func (*Postgres) [Lock](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L233) ¶

    func (p *Postgres) Lock() [error](/builtin#error)

<https://www.postgresql.org/docs/9.6/static/explicit-locking.html#ADVISORY-LOCKS>

#### func (*Postgres) [Open](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L152) ¶

    func (p *Postgres) Open(url [string](/builtin#string)) ([database](/github.com/golang-migrate/migrate/v4@v4.19.1/database).[Driver](/github.com/golang-migrate/migrate/v4@v4.19.1/database#Driver), [error](/builtin#error))

#### func (*Postgres) [Run](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L265) ¶

    func (p *Postgres) Run(migration [io](/io).[Reader](/io#Reader)) [error](/builtin#error)

#### func (*Postgres) [SetVersion](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L355) ¶

    func (p *Postgres) SetVersion(version [int](/builtin#int), dirty [bool](/builtin#bool)) [error](/builtin#error)

#### func (*Postgres) [Unlock](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L250) ¶

    func (p *Postgres) Unlock() [error](/builtin#error)

#### func (*Postgres) [Version](https://github.com/golang-migrate/migrate/blob/v4.19.1/database/postgres/postgres.go#L389) ¶

    func (p *Postgres) Version() (version [int](/builtin#int), dirty [bool](/builtin#bool), err [error](/builtin#error))
