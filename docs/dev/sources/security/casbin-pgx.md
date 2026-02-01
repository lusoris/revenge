# Casbin pgx Adapter

> Source: https://pkg.go.dev/github.com/pckhoi/casbin-pgx-adapter/v3
> Fetched: 2026-02-01T11:49:25.841351+00:00
> Content-Hash: 136ff6db0df064d9
> Type: html

---

### Index ¶

- Constants
- type Adapter
-     * func NewAdapter(conn interface{}, opts ...Option) (*Adapter, error)
-     * func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error
  - func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error
  - func (a *Adapter) Close()
  - func (a *Adapter) IsFiltered() bool
  - func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error
  - func (a *Adapter) LoadPolicy(model model.Model) error
  - func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error
  - func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error
  - func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error
  - func (a *Adapter) SavePolicy(model model.Model) error
  - func (a *Adapter) UpdateFilteredPolicies(sec string, ptype string, newPolicies [][]string, fieldIndex int, ...) ([][]string, error)
  - func (a *Adapter) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error
  - func (a *Adapter) UpdatePolicy(sec string, ptype string, oldRule, newPolicy []string) error
- type Filter
- type Option
-     * func WithConnectionPool(pool *pgxpool.Pool) Option
  - func WithDatabase(dbname string) Option
  - func WithSchema(s string) Option
  - func WithSkipTableCreate() Option
  - func WithTableName(tableName string) Option
  - func WithTimeout(timeout time.Duration) Option

### Constants ¶

[View Source](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L19)

    const (
     DefaultTableName    = "casbin_rule"
     DefaultDatabaseName = "casbin"
     DefaultTimeout      = [time](/time).[Second](/time#Second) * 10
    )

### Variables ¶

This section is empty.

### Functions ¶

This section is empty.

### Types ¶

#### type [Adapter](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L26) ¶

    type Adapter struct {
     // contains filtered or unexported fields
    }

Adapter represents the github.com/jackc/pgx/v5 adapter for policy storage.

#### func [NewAdapter](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L45) ¶

    func NewAdapter(conn interface{}, opts ...Option) (*Adapter, [error](/builtin#error))

NewAdapter creates a new adapter with connection conn which must either be a PostgreSQL connection string or an instance of *pgx.ConnConfig from package github.com/jackc/pgx/v5.

#### func (*Adapter) [AddPolicies](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L247) ¶

    func (a *Adapter) AddPolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) [error](/builtin#error)

AddPolicies adds policy rules to the storage.

#### func (*Adapter) [AddPolicy](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L236) ¶

    func (a *Adapter) AddPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) [error](/builtin#error)

AddPolicy adds a policy rule to the storage.

#### func (*Adapter) [Close](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L436) ¶

    func (a *Adapter) Close()

#### func (*Adapter) [IsFiltered](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L397) ¶

    func (a *Adapter) IsFiltered() [bool](/builtin#bool)

#### func (*Adapter) [LoadFilteredPolicy](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L380) ¶

    func (a *Adapter) LoadFilteredPolicy(model [model](/github.com/casbin/casbin/v2/model).[Model](/github.com/casbin/casbin/v2/model#Model), filter interface{}) [error](/builtin#error)

LoadFilteredPolicy can query policies with a filter. Make sure that filter is of type *pgxadapter.Filter

#### func (*Adapter) [LoadPolicy](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L144) ¶

    func (a *Adapter) LoadPolicy(model [model](/github.com/casbin/casbin/v2/model).[Model](/github.com/casbin/casbin/v2/model#Model)) [error](/builtin#error)

LoadPolicy loads policy from database.

#### func (*Adapter) [RemoveFilteredPolicy](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L301) ¶

    func (a *Adapter) RemoveFilteredPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) [error](/builtin#error)

RemoveFilteredPolicy removes policy rules that match the filter from the storage.

#### func (*Adapter) [RemovePolicies](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L280) ¶

    func (a *Adapter) RemovePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) [error](/builtin#error)

RemovePolicies removes policy rules from the storage.

#### func (*Adapter) [RemovePolicy](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L268) ¶

    func (a *Adapter) RemovePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) [error](/builtin#error)

RemovePolicy removes a policy rule from the storage.

#### func (*Adapter) [SavePolicy](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L198) ¶

    func (a *Adapter) SavePolicy(model [model](/github.com/casbin/casbin/v2/model).[Model](/github.com/casbin/casbin/v2/model#Model)) [error](/builtin#error)

SavePolicy saves policy to database.

#### func (*Adapter) [UpdateFilteredPolicies](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L432) ¶

    func (a *Adapter) UpdateFilteredPolicies(sec [string](/builtin#string), ptype [string](/builtin#string), newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

UpdateFilteredPolicies deletes old rules and adds new rules.

#### func (*Adapter) [UpdatePolicies](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L408) ¶

    func (a *Adapter) UpdatePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) [error](/builtin#error)

UpdatePolicies updates some policy rules to storage, like db, redis.

#### func (*Adapter) [UpdatePolicy](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L403) ¶

    func (a *Adapter) UpdatePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newPolicy [][string](/builtin#string)) [error](/builtin#error)

UpdatePolicy updates a policy rule from storage. This is part of the Auto-Save feature.

#### type [Filter](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L36) ¶

    type Filter struct {
     P [][][string](/builtin#string)
     G [][][string](/builtin#string)
    }

#### type [Option](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L41) ¶

    type Option func(a *Adapter)

#### func [WithConnectionPool](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L103) ¶ added in v3.1.0

    func WithConnectionPool(pool *[pgxpool](/github.com/jackc/pgx/v5/pgxpool).[Pool](/github.com/jackc/pgx/v5/pgxpool#Pool)) Option

WithConnectionPool can be used to pass an existing *pgxpool.Pool instance

#### func [WithDatabase](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L88) ¶

    func WithDatabase(dbname [string](/builtin#string)) Option

WithDatabase can be used to pass custom database name for Casbin rules

#### func [WithSchema](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L112) ¶

    func WithSchema(s [string](/builtin#string)) Option

WithSchema can be used to pass a custom schema name. Note that the schema name is case-sensitive. If you don't create the schema before hand, the schema will be created for you.

#### func [WithSkipTableCreate](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L81) ¶

    func WithSkipTableCreate() Option

WithSkipTableCreate skips the table creation step when the adapter starts If the Casbin rules table does not exist, it will lead to issues when using the adapter

#### func [WithTableName](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L73) ¶

    func WithTableName(tableName [string](/builtin#string)) Option

WithTableName can be used to pass custom table name for Casbin rules

#### func [WithTimeout](https://github.com/pckhoi/casbin-pgx-adapter/blob/v3.2.0/adapter.go#L96) ¶

    func WithTimeout(timeout [time](/time).[Duration](/time#Duration)) Option

WithTimeout can be used to pass a different timeout than DefaultTimeout for each request to Postgres
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
