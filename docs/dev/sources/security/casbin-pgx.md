# Casbin pgx Adapter

> Source: https://pkg.go.dev/github.com/pckhoi/casbin-pgx-adapter/v3
> Fetched: 2026-01-30T23:54:42.911232+00:00
> Content-Hash: e2da7930d66e03d2
> Type: html

---

Index

¶

Constants

type Adapter

func NewAdapter(conn interface{}, opts ...Option) (*Adapter, error)

func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error

func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error

func (a *Adapter) Close()

func (a *Adapter) IsFiltered() bool

func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error

func (a *Adapter) LoadPolicy(model model.Model) error

func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error

func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error

func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error

func (a *Adapter) SavePolicy(model model.Model) error

func (a *Adapter) UpdateFilteredPolicies(sec string, ptype string, newPolicies [][]string, fieldIndex int, ...) ([][]string, error)

func (a *Adapter) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error

func (a *Adapter) UpdatePolicy(sec string, ptype string, oldRule, newPolicy []string) error

type Filter

type Option

func WithConnectionPool(pool *pgxpool.Pool) Option

func WithDatabase(dbname string) Option

func WithSchema(s string) Option

func WithSkipTableCreate() Option

func WithTableName(tableName string) Option

func WithTimeout(timeout time.Duration) Option

Constants

¶

View Source

const (

DefaultTableName    = "casbin_rule"

DefaultDatabaseName = "casbin"

DefaultTimeout      =

time

.

Second

* 10

)

Variables

¶

This section is empty.

Functions

¶

This section is empty.

Types

¶

type

Adapter

¶

type Adapter struct {

// contains filtered or unexported fields

}

Adapter represents the github.com/jackc/pgx/v5 adapter for policy storage.

func

NewAdapter

¶

func NewAdapter(conn interface{}, opts ...

Option

) (*

Adapter

,

error

)

NewAdapter creates a new adapter with connection conn which must either be a PostgreSQL
connection string or an instance of *pgx.ConnConfig from package github.com/jackc/pgx/v5.

func (*Adapter)

AddPolicies

¶

func (a *

Adapter

) AddPolicies(sec

string

, ptype

string

, rules [][]

string

)

error

AddPolicies adds policy rules to the storage.

func (*Adapter)

AddPolicy

¶

func (a *

Adapter

) AddPolicy(sec

string

, ptype

string

, rule []

string

)

error

AddPolicy adds a policy rule to the storage.

func (*Adapter)

Close

¶

func (a *

Adapter

) Close()

func (*Adapter)

IsFiltered

¶

func (a *

Adapter

) IsFiltered()

bool

func (*Adapter)

LoadFilteredPolicy

¶

func (a *

Adapter

) LoadFilteredPolicy(model

model

.

Model

, filter interface{})

error

LoadFilteredPolicy can query policies with a filter. Make sure that filter is of type *pgxadapter.Filter

func (*Adapter)

LoadPolicy

¶

func (a *

Adapter

) LoadPolicy(model

model

.

Model

)

error

LoadPolicy loads policy from database.

func (*Adapter)

RemoveFilteredPolicy

¶

func (a *

Adapter

) RemoveFilteredPolicy(sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

)

error

RemoveFilteredPolicy removes policy rules that match the filter from the storage.

func (*Adapter)

RemovePolicies

¶

func (a *

Adapter

) RemovePolicies(sec

string

, ptype

string

, rules [][]

string

)

error

RemovePolicies removes policy rules from the storage.

func (*Adapter)

RemovePolicy

¶

func (a *

Adapter

) RemovePolicy(sec

string

, ptype

string

, rule []

string

)

error

RemovePolicy removes a policy rule from the storage.

func (*Adapter)

SavePolicy

¶

func (a *

Adapter

) SavePolicy(model

model

.

Model

)

error

SavePolicy saves policy to database.

func (*Adapter)

UpdateFilteredPolicies

¶

func (a *

Adapter

) UpdateFilteredPolicies(sec

string

, ptype

string

, newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) ([][]

string

,

error

)

UpdateFilteredPolicies deletes old rules and adds new rules.

func (*Adapter)

UpdatePolicies

¶

func (a *

Adapter

) UpdatePolicies(sec

string

, ptype

string

, oldRules, newRules [][]

string

)

error

UpdatePolicies updates some policy rules to storage, like db, redis.

func (*Adapter)

UpdatePolicy

¶

func (a *

Adapter

) UpdatePolicy(sec

string

, ptype

string

, oldRule, newPolicy []

string

)

error

UpdatePolicy updates a policy rule from storage.
This is part of the Auto-Save feature.

type

Filter

¶

type Filter struct {

P [][]

string

G [][]

string

}

type

Option

¶

type Option func(a *

Adapter

)

func

WithConnectionPool

¶

added in

v3.1.0

func WithConnectionPool(pool *

pgxpool

.

Pool

)

Option

WithConnectionPool can be used to pass an existing *pgxpool.Pool instance

func

WithDatabase

¶

func WithDatabase(dbname

string

)

Option

WithDatabase can be used to pass custom database name for Casbin rules

func

WithSchema

¶

func WithSchema(s

string

)

Option

WithSchema can be used to pass a custom schema name. Note that the schema
name is case-sensitive. If you don't create the schema before hand, the
schema will be created for you.

func

WithSkipTableCreate

¶

func WithSkipTableCreate()

Option

WithSkipTableCreate skips the table creation step when the adapter starts
If the Casbin rules table does not exist, it will lead to issues when using the adapter

func

WithTableName

¶

func WithTableName(tableName

string

)

Option

WithTableName can be used to pass custom table name for Casbin rules

func

WithTimeout

¶

func WithTimeout(timeout

time

.

Duration

)

Option

WithTimeout can be used to pass a different timeout than DefaultTimeout
for each request to Postgres