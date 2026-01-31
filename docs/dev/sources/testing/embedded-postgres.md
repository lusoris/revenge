# embedded-postgres

> Source: https://pkg.go.dev/github.com/fergusstrange/embedded-postgres
> Fetched: 2026-01-30T23:55:02.389570+00:00
> Content-Hash: 2fb32dcee8508970
> Type: html

---

Index

¶

Constants

Variables

func TestGetConnectionURL(t *testing.T)

type CacheLocator

type Config

func DefaultConfig() Config

func (c Config) BinariesPath(path string) Config

func (c Config) BinaryRepositoryURL(binaryRepositoryURL string) Config

func (c Config) CachePath(path string) Config

func (c Config) DataPath(path string) Config

func (c Config) Database(database string) Config

func (c Config) Encoding(encoding string) Config

func (c Config) GetConnectionURL() string

func (c Config) Locale(locale string) Config

func (c Config) Logger(logger io.Writer) Config

func (c Config) Password(password string) Config

func (c Config) Port(port uint32) Config

func (c Config) RuntimePath(path string) Config

func (c Config) StartParameters(parameters map[string]string) Config

func (c Config) StartTimeout(timeout time.Duration) Config

func (c Config) Username(username string) Config

func (c Config) Version(version PostgresVersion) Config

type EmbeddedPostgres

func NewDatabase(config ...Config) *EmbeddedPostgres

func (ep *EmbeddedPostgres) Start() error

func (ep *EmbeddedPostgres) Stop() error

type PostgresVersion

type RemoteFetchStrategy

type VersionStrategy

Constants

¶

View Source

const (

V18 =

PostgresVersion

("18.0.0")

V17 =

PostgresVersion

("17.5.0")

V16 =

PostgresVersion

("16.9.0")

V15 =

PostgresVersion

("15.13.0")

V14 =

PostgresVersion

("14.18.0")

V13 =

PostgresVersion

("13.21.0")

V12 =

PostgresVersion

("12.22.0")

V11 =

PostgresVersion

("11.22.0")

V10 =

PostgresVersion

("10.23.0")

V9  =

PostgresVersion

("9.6.24")

)

Predefined supported Postgres versions.

Variables

¶

View Source

var (

ErrServerNotStarted     =

errors

.

New

("server has not been started")

ErrServerAlreadyStarted =

errors

.

New

("server is already started")

)

Functions

¶

func

TestGetConnectionURL

¶

added in

v1.22.0

func TestGetConnectionURL(t *

testing

.

T

)

Types

¶

type

CacheLocator

¶

type CacheLocator func() (location

string

, exists

bool

)

CacheLocator retrieves the location of the Postgres binary cache returning it to location.
The result of whether this cache is present will be returned to exists.

type

Config

¶

type Config struct {

// contains filtered or unexported fields

}

Config maintains the runtime configuration for the Postgres process to be created.

func

DefaultConfig

¶

func DefaultConfig()

Config

DefaultConfig provides a default set of configuration to be used "as is" or modified using the provided builders.
The following can be assumed as defaults:
Version:      16
Port:         5432
Database:     postgres
Username:     postgres
Password:     postgres
StartTimeout: 15 Seconds

func (Config)

BinariesPath

¶

added in

v1.10.0

func (c

Config

) BinariesPath(path

string

)

Config

BinariesPath sets the path of the pre-downloaded postgres binaries.
If this option is left unset, the binaries will be downloaded.

func (Config)

BinaryRepositoryURL

¶

added in

v1.15.0

func (c

Config

) BinaryRepositoryURL(binaryRepositoryURL

string

)

Config

BinaryRepositoryURL set BinaryRepositoryURL to fetch PG Binary in case of Maven proxy

func (Config)

CachePath

¶

added in

v1.24.0

func (c

Config

) CachePath(path

string

)

Config

CachePath sets the path that will be used for storing Postgres binaries archive.
If this option is not set, ~/.go-embedded-postgres will be used.

func (Config)

DataPath

¶

added in

v1.4.0

func (c

Config

) DataPath(path

string

)

Config

DataPath sets the path that will be used for the Postgres data directory.
If this option is set, a previously initialized data directory will be reused if possible.

func (Config)

Database

¶

func (c

Config

) Database(database

string

)

Config

Database sets the database name that will be created.

func (Config)

Encoding

¶

added in

v1.27.0

func (c

Config

) Encoding(encoding

string

)

Config

Encoding sets the default character set for initdb

func (Config)

GetConnectionURL

¶

added in

v1.22.0

func (c

Config

) GetConnectionURL()

string

func (Config)

Locale

¶

added in

v1.1.0

func (c

Config

) Locale(locale

string

)

Config

Locale sets the default locale for initdb

func (Config)

Logger

¶

added in

v1.3.0

func (c

Config

) Logger(logger

io

.

Writer

)

Config

Logger sets the logger for postgres output

func (Config)

Password

¶

func (c

Config

) Password(password

string

)

Config

Password sets the password that will be used to connect.

func (Config)

Port

¶

func (c

Config

) Port(port

uint32

)

Config

Port sets the runtime port that Postgres can be accessed on.

func (Config)

RuntimePath

¶

func (c

Config

) RuntimePath(path

string

)

Config

RuntimePath sets the path that will be used for the extracted Postgres runtime directory.
If Postgres data directory is not set with DataPath(), this directory is also used as data directory.

func (Config)

StartParameters

¶

added in

v1.23.0

func (c

Config

) StartParameters(parameters map[

string

]

string

)

Config

StartParameters sets run-time parameters when starting Postgres (passed to Postgres via "-c").

These parameters can be used to override the default configuration values in postgres.conf such
as max_connections=100. See

https://www.postgresql.org/docs/current/runtime-config.html

func (Config)

StartTimeout

¶

func (c

Config

) StartTimeout(timeout

time

.

Duration

)

Config

StartTimeout sets the max timeout that will be used when starting the Postgres process and creating the initial database.

func (Config)

Username

¶

func (c

Config

) Username(username

string

)

Config

Username sets the username that will be used to connect.

func (Config)

Version

¶

func (c

Config

) Version(version

PostgresVersion

)

Config

Version will set the Postgres binary version.

type

EmbeddedPostgres

¶

type EmbeddedPostgres struct {

// contains filtered or unexported fields

}

EmbeddedPostgres maintains all configuration and runtime functions for maintaining the lifecycle of one Postgres process.

func

NewDatabase

¶

func NewDatabase(config ...

Config

) *

EmbeddedPostgres

NewDatabase creates a new EmbeddedPostgres struct that can be used to start and stop a Postgres process.
When called with no parameters it will assume a default configuration state provided by the DefaultConfig method.
When called with parameters the first Config parameter will be used for configuration.

func (*EmbeddedPostgres)

Start

¶

func (ep *

EmbeddedPostgres

) Start()

error

Start will try to start the configured Postgres process returning an error when there were any problems with invocation.
If any error occurs Start will try to also Stop the Postgres process in order to not leave any sub-process running.

func (*EmbeddedPostgres)

Stop

¶

func (ep *

EmbeddedPostgres

) Stop()

error

Stop will try to stop the Postgres process gracefully returning an error when there were any problems.

type

PostgresVersion

¶

type PostgresVersion

string

PostgresVersion represents the semantic version used to fetch and run the Postgres process.

type

RemoteFetchStrategy

¶

type RemoteFetchStrategy func()

error

RemoteFetchStrategy provides a strategy to fetch a Postgres binary so that it is available for use.

type

VersionStrategy

¶

type VersionStrategy func() (operatingSystem

string

, architecture

string

, postgresVersion

PostgresVersion

)

VersionStrategy provides a strategy that can be used to determine which version of Postgres should be used based on
the operating system, architecture and desired Postgres version.