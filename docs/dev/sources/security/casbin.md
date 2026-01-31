# Casbin

> Source: https://pkg.go.dev/github.com/casbin/casbin/v2
> Fetched: 2026-01-30T23:54:37.686834+00:00
> Content-Hash: fe2599881fd42305
> Type: html

---

Overview

¶

rbac_api_context.go

Index

¶

func CasbinJsGetPermissionForUser(e IEnforcer, user string) (string, error)

func CasbinJsGetPermissionForUserOld(e IEnforcer, user string) ([]byte, error)

func GetCacheKey(params ...interface{}) (string, bool)

type CacheableParam

type CachedEnforcer

func NewCachedEnforcer(params ...interface{}) (*CachedEnforcer, error)

func (e *CachedEnforcer) ClearPolicy()

func (e *CachedEnforcer) EnableCache(enableCache bool)

func (e *CachedEnforcer) Enforce(rvals ...interface{}) (bool, error)

func (e *CachedEnforcer) InvalidateCache() error

func (e *CachedEnforcer) LoadPolicy() error

func (e *CachedEnforcer) RemovePolicies(rules [][]string) (bool, error)

func (e *CachedEnforcer) RemovePolicy(params ...interface{}) (bool, error)

func (e *CachedEnforcer) SetCache(c cache.Cache)

func (e *CachedEnforcer) SetExpireTime(expireTime time.Duration)

type ConflictDetector

func NewConflictDetector(baseModel, currentModel model.Model, operations []persist.PolicyOperation) *ConflictDetector

func (cd *ConflictDetector) DetectConflicts() error

type ConflictError

func (e *ConflictError) Error() string

type ContextEnforcer

func (e *ContextEnforcer) AddGroupingPoliciesCtx(ctx context.Context, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddGroupingPoliciesExCtx(ctx context.Context, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddGroupingPolicyCtx(ctx context.Context, params ...interface{}) (bool, error)

func (e *ContextEnforcer) AddNamedGroupingPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddNamedGroupingPoliciesExCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddNamedGroupingPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)

func (e *ContextEnforcer) AddNamedPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddNamedPoliciesExCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddNamedPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)

func (e *ContextEnforcer) AddPermissionForUserCtx(ctx context.Context, user string, permission ...string) (bool, error)

func (e *ContextEnforcer) AddPermissionsForUserCtx(ctx context.Context, user string, permissions ...[]string) (bool, error)

func (e *ContextEnforcer) AddPoliciesCtx(ctx context.Context, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddPoliciesExCtx(ctx context.Context, rules [][]string) (bool, error)

func (e *ContextEnforcer) AddPolicyCtx(ctx context.Context, params ...interface{}) (bool, error)

func (e *ContextEnforcer) AddRoleForUserCtx(ctx context.Context, user string, role string, domain ...string) (bool, error)

func (e *ContextEnforcer) AddRoleForUserInDomainCtx(ctx context.Context, user string, role string, domain string) (bool, error)

func (e *ContextEnforcer) DeleteAllUsersByDomainCtx(ctx context.Context, domain string) (bool, error)

func (e *ContextEnforcer) DeleteDomainsCtx(ctx context.Context, domains ...string) (bool, error)

func (e *ContextEnforcer) DeletePermissionCtx(ctx context.Context, permission ...string) (bool, error)

func (e *ContextEnforcer) DeletePermissionForUserCtx(ctx context.Context, user string, permission ...string) (bool, error)

func (e *ContextEnforcer) DeletePermissionsForUserCtx(ctx context.Context, user string) (bool, error)

func (e *ContextEnforcer) DeleteRoleCtx(ctx context.Context, role string) (bool, error)

func (e *ContextEnforcer) DeleteRoleForUserCtx(ctx context.Context, user string, role string, domain ...string) (bool, error)

func (e *ContextEnforcer) DeleteRoleForUserInDomainCtx(ctx context.Context, user string, role string, domain string) (bool, error)

func (e *ContextEnforcer) DeleteRolesForUserCtx(ctx context.Context, user string, domain ...string) (bool, error)

func (e *ContextEnforcer) DeleteRolesForUserInDomainCtx(ctx context.Context, user string, domain string) (bool, error)

func (e *ContextEnforcer) DeleteUserCtx(ctx context.Context, user string) (bool, error)

func (e *ContextEnforcer) IsFilteredCtx(ctx context.Context) bool

func (e *ContextEnforcer) LoadPolicyCtx(ctx context.Context) error

func (e *ContextEnforcer) RemoveFilteredGroupingPolicyCtx(ctx context.Context, fieldIndex int, fieldValues ...string) (bool, error)

func (e *ContextEnforcer) RemoveFilteredNamedGroupingPolicyCtx(ctx context.Context, ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *ContextEnforcer) RemoveFilteredNamedPolicyCtx(ctx context.Context, ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *ContextEnforcer) RemoveFilteredPolicyCtx(ctx context.Context, fieldIndex int, fieldValues ...string) (bool, error)

func (e *ContextEnforcer) RemoveGroupingPoliciesCtx(ctx context.Context, rules [][]string) (bool, error)

func (e *ContextEnforcer) RemoveGroupingPolicyCtx(ctx context.Context, params ...interface{}) (bool, error)

func (e *ContextEnforcer) RemoveNamedGroupingPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) RemoveNamedGroupingPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)

func (e *ContextEnforcer) RemoveNamedPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) RemoveNamedPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)

func (e *ContextEnforcer) RemovePoliciesCtx(ctx context.Context, rules [][]string) (bool, error)

func (e *ContextEnforcer) RemovePolicyCtx(ctx context.Context, params ...interface{}) (bool, error)

func (e *ContextEnforcer) SavePolicyCtx(ctx context.Context) error

func (e *ContextEnforcer) SelfAddPoliciesCtx(ctx context.Context, sec string, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) SelfAddPoliciesExCtx(ctx context.Context, sec string, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) SelfAddPolicyCtx(ctx context.Context, sec string, ptype string, rule []string) (bool, error)

func (e *ContextEnforcer) SelfRemoveFilteredPolicyCtx(ctx context.Context, sec string, ptype string, fieldIndex int, ...) (bool, error)

func (e *ContextEnforcer) SelfRemovePoliciesCtx(ctx context.Context, sec string, ptype string, rules [][]string) (bool, error)

func (e *ContextEnforcer) SelfRemovePolicyCtx(ctx context.Context, sec string, ptype string, rule []string) (bool, error)

func (e *ContextEnforcer) SelfUpdatePoliciesCtx(ctx context.Context, sec string, ptype string, oldRules, newRules [][]string) (bool, error)

func (e *ContextEnforcer) SelfUpdatePolicyCtx(ctx context.Context, sec string, ptype string, oldRule, newRule []string) (bool, error)

func (e *ContextEnforcer) UpdateFilteredNamedPoliciesCtx(ctx context.Context, ptype string, newPolicies [][]string, fieldIndex int, ...) (bool, error)

func (e *ContextEnforcer) UpdateFilteredPoliciesCtx(ctx context.Context, newPolicies [][]string, fieldIndex int, ...) (bool, error)

func (e *ContextEnforcer) UpdateGroupingPoliciesCtx(ctx context.Context, oldRules [][]string, newRules [][]string) (bool, error)

func (e *ContextEnforcer) UpdateGroupingPolicyCtx(ctx context.Context, oldRule []string, newRule []string) (bool, error)

func (e *ContextEnforcer) UpdateNamedGroupingPoliciesCtx(ctx context.Context, ptype string, oldRules [][]string, newRules [][]string) (bool, error)

func (e *ContextEnforcer) UpdateNamedGroupingPolicyCtx(ctx context.Context, ptype string, oldRule []string, newRule []string) (bool, error)

func (e *ContextEnforcer) UpdateNamedPoliciesCtx(ctx context.Context, ptype string, p1 [][]string, p2 [][]string) (bool, error)

func (e *ContextEnforcer) UpdateNamedPolicyCtx(ctx context.Context, ptype string, p1 []string, p2 []string) (bool, error)

func (e *ContextEnforcer) UpdatePoliciesCtx(ctx context.Context, oldPolicies [][]string, newPolicies [][]string) (bool, error)

func (e *ContextEnforcer) UpdatePolicyCtx(ctx context.Context, oldPolicy []string, newPolicy []string) (bool, error)

type DistributedEnforcer

func NewDistributedEnforcer(params ...interface{}) (*DistributedEnforcer, error)

func (d *DistributedEnforcer) AddPoliciesSelf(shouldPersist func() bool, sec string, ptype string, rules [][]string) (affected [][]string, err error)

func (d *DistributedEnforcer) ClearPolicySelf(shouldPersist func() bool) error

func (d *DistributedEnforcer) RemoveFilteredPolicySelf(shouldPersist func() bool, sec string, ptype string, fieldIndex int, ...) (affected [][]string, err error)

func (d *DistributedEnforcer) RemovePoliciesSelf(shouldPersist func() bool, sec string, ptype string, rules [][]string) (affected [][]string, err error)

func (d *DistributedEnforcer) SetDispatcher(dispatcher persist.Dispatcher)

func (d *DistributedEnforcer) UpdateFilteredPoliciesSelf(shouldPersist func() bool, sec string, ptype string, newRules [][]string, ...) (bool, error)

func (d *DistributedEnforcer) UpdatePoliciesSelf(shouldPersist func() bool, sec string, ptype string, ...) (affected bool, err error)

func (d *DistributedEnforcer) UpdatePolicySelf(shouldPersist func() bool, sec string, ptype string, oldRule, newRule []string) (affected bool, err error)

type EnforceContext

func NewEnforceContext(suffix string) EnforceContext

func (e EnforceContext) GetCacheKey() string

type Enforcer

func NewEnforcer(params ...interface{}) (*Enforcer, error)

func (e *Enforcer) AddFunction(name string, function govaluate.ExpressionFunction)

func (e *Enforcer) AddGroupingPolicies(rules [][]string) (bool, error)

func (e *Enforcer) AddGroupingPoliciesEx(rules [][]string) (bool, error)

func (e *Enforcer) AddGroupingPolicy(params ...interface{}) (bool, error)

func (e *Enforcer) AddNamedDomainLinkConditionFunc(ptype, user, role string, domain string, fn rbac.LinkConditionFunc) bool

func (e *Enforcer) AddNamedDomainMatchingFunc(ptype, name string, fn rbac.MatchingFunc) bool

func (e *Enforcer) AddNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)

func (e *Enforcer) AddNamedGroupingPoliciesEx(ptype string, rules [][]string) (bool, error)

func (e *Enforcer) AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (e *Enforcer) AddNamedLinkConditionFunc(ptype, user, role string, fn rbac.LinkConditionFunc) bool

func (e *Enforcer) AddNamedMatchingFunc(ptype, name string, fn rbac.MatchingFunc) bool

func (e *Enforcer) AddNamedPolicies(ptype string, rules [][]string) (bool, error)

func (e *Enforcer) AddNamedPoliciesEx(ptype string, rules [][]string) (bool, error)

func (e *Enforcer) AddNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (e *Enforcer) AddPermissionForUser(user string, permission ...string) (bool, error)

func (e *Enforcer) AddPermissionsForUser(user string, permissions ...[]string) (bool, error)

func (e *Enforcer) AddPolicies(rules [][]string) (bool, error)

func (e *Enforcer) AddPoliciesEx(rules [][]string) (bool, error)

func (e *Enforcer) AddPolicy(params ...interface{}) (bool, error)

func (e *Enforcer) AddRoleForUser(user string, role string, domain ...string) (bool, error)

func (e *Enforcer) AddRoleForUserInDomain(user string, role string, domain string) (bool, error)

func (e *Enforcer) AddRolesForUser(user string, roles []string, domain ...string) (bool, error)

func (e *Enforcer) BatchEnforce(requests [][]interface{}) ([]bool, error)

func (e *Enforcer) BatchEnforceWithMatcher(matcher string, requests [][]interface{}) ([]bool, error)

func (e *Enforcer) BuildIncrementalConditionalRoleLinks(op model.PolicyOp, ptype string, rules [][]string) error

func (e *Enforcer) BuildIncrementalRoleLinks(op model.PolicyOp, ptype string, rules [][]string) error

func (e *Enforcer) BuildRoleLinks() error

func (e *Enforcer) ClearPolicy()

func (e *Enforcer) DeleteAllUsersByDomain(domain string) (bool, error)

func (e *Enforcer) DeleteDomains(domains ...string) (bool, error)

func (e *Enforcer) DeletePermission(permission ...string) (bool, error)

func (e *Enforcer) DeletePermissionForUser(user string, permission ...string) (bool, error)

func (e *Enforcer) DeletePermissionsForUser(user string) (bool, error)

func (e *Enforcer) DeleteRole(role string) (bool, error)

func (e *Enforcer) DeleteRoleForUser(user string, role string, domain ...string) (bool, error)

func (e *Enforcer) DeleteRoleForUserInDomain(user string, role string, domain string) (bool, error)

func (e *Enforcer) DeleteRolesForUser(user string, domain ...string) (bool, error)

func (e *Enforcer) DeleteRolesForUserInDomain(user string, domain string) (bool, error)

func (e *Enforcer) DeleteUser(user string) (bool, error)

func (e *Enforcer) EnableAcceptJsonRequest(acceptJsonRequest bool)

func (e *Enforcer) EnableAutoBuildRoleLinks(autoBuildRoleLinks bool)

func (e *Enforcer) EnableAutoNotifyDispatcher(enable bool)

func (e *Enforcer) EnableAutoNotifyWatcher(enable bool)

func (e *Enforcer) EnableAutoSave(autoSave bool)

func (e *Enforcer) EnableEnforce(enable bool)

func (e *Enforcer) EnableLog(enable bool)

func (e *Enforcer) Enforce(rvals ...interface{}) (bool, error)

func (e *Enforcer) EnforceEx(rvals ...interface{}) (bool, []string, error)

func (e *Enforcer) EnforceExWithMatcher(matcher string, rvals ...interface{}) (bool, []string, error)

func (e *Enforcer) EnforceWithMatcher(matcher string, rvals ...interface{}) (bool, error)

func (e *Enforcer) GetAdapter() persist.Adapter

func (e *Enforcer) GetAllActions() ([]string, error)

func (e *Enforcer) GetAllDomains() ([]string, error)

func (e *Enforcer) GetAllNamedActions(ptype string) ([]string, error)

func (e *Enforcer) GetAllNamedObjects(ptype string) ([]string, error)

func (e *Enforcer) GetAllNamedRoles(ptype string) ([]string, error)

func (e *Enforcer) GetAllNamedSubjects(ptype string) ([]string, error)

func (e *Enforcer) GetAllObjects() ([]string, error)

func (e *Enforcer) GetAllRoles() ([]string, error)

func (e *Enforcer) GetAllRolesByDomain(domain string) ([]string, error)

func (e *Enforcer) GetAllSubjects() ([]string, error)

func (e *Enforcer) GetAllUsersByDomain(domain string) ([]string, error)

func (e *Enforcer) GetAllowedObjectConditions(user string, action string, prefix string) ([]string, error)

func (e *Enforcer) GetDomainsForUser(user string) ([]string, error)

func (e *Enforcer) GetFieldIndex(ptype string, field string) (int, error)

func (e *Enforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *Enforcer) GetFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *Enforcer) GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *Enforcer) GetFilteredNamedPolicyWithMatcher(ptype string, matcher string) ([][]string, error)

func (e *Enforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *Enforcer) GetGroupingPolicy() ([][]string, error)

func (e *Enforcer) GetImplicitObjectPatternsForUser(user string, domain string, action string) ([]string, error)

func (e *Enforcer) GetImplicitPermissionsForUser(user string, domain ...string) ([][]string, error)

func (e *Enforcer) GetImplicitResourcesForUser(user string, domain ...string) ([][]string, error)

func (e *Enforcer) GetImplicitRolesForUser(name string, domain ...string) ([]string, error)

func (e *Enforcer) GetImplicitUsersForPermission(permission ...string) ([]string, error)

func (e *Enforcer) GetImplicitUsersForResource(resource string) ([][]string, error)

func (e *Enforcer) GetImplicitUsersForResourceByDomain(resource string, domain string) ([][]string, error)

func (e *Enforcer) GetImplicitUsersForRole(name string, domain ...string) ([]string, error)

func (e *Enforcer) GetModel() model.Model

func (e *Enforcer) GetNamedGroupingPolicy(ptype string) ([][]string, error)

func (e *Enforcer) GetNamedImplicitPermissionsForUser(ptype string, gtype string, user string, domain ...string) ([][]string, error)

func (e *Enforcer) GetNamedImplicitRolesForUser(ptype string, name string, domain ...string) ([]string, error)

func (e *Enforcer) GetNamedImplicitUsersForResource(ptype string, resource string) ([][]string, error)

func (e *Enforcer) GetNamedPermissionsForUser(ptype string, user string, domain ...string) ([][]string, error)

func (e *Enforcer) GetNamedPolicy(ptype string) ([][]string, error)

func (e *Enforcer) GetNamedRoleManager(ptype string) rbac.RoleManager

func (e *Enforcer) GetPermissionsForUser(user string, domain ...string) ([][]string, error)

func (e *Enforcer) GetPermissionsForUserInDomain(user string, domain string) [][]string

func (e *Enforcer) GetPolicy() ([][]string, error)

func (e *Enforcer) GetRoleManager() rbac.RoleManager

func (e *Enforcer) GetRolesForUser(name string, domain ...string) ([]string, error)

func (e *Enforcer) GetRolesForUserInDomain(name string, domain string) []string

func (e *Enforcer) GetUsersForRole(name string, domain ...string) ([]string, error)

func (e *Enforcer) GetUsersForRoleInDomain(name string, domain string) []string

func (e *Enforcer) HasGroupingPolicy(params ...interface{}) (bool, error)

func (e *Enforcer) HasNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (e *Enforcer) HasNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (e *Enforcer) HasPermissionForUser(user string, permission ...string) (bool, error)

func (e *Enforcer) HasPolicy(params ...interface{}) (bool, error)

func (e *Enforcer) HasRoleForUser(name string, role string, domain ...string) (bool, error)

func (e *Enforcer) InitWithAdapter(modelPath string, adapter persist.Adapter) error

func (e *Enforcer) InitWithFile(modelPath string, policyPath string) error

func (e *Enforcer) InitWithModelAndAdapter(m model.Model, adapter persist.Adapter) error

func (e *Enforcer) IsFiltered() bool

func (e *Enforcer) IsLogEnabled() bool

func (e *Enforcer) LoadFilteredPolicy(filter interface{}) error

func (e *Enforcer) LoadFilteredPolicyCtx(ctx context.Context, filter interface{}) error

func (e *Enforcer) LoadIncrementalFilteredPolicy(filter interface{}) error

func (e *Enforcer) LoadIncrementalFilteredPolicyCtx(ctx context.Context, filter interface{}) error

func (e *Enforcer) LoadModel() error

func (e *Enforcer) LoadPolicy() error

func (e *Enforcer) RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) (bool, error)

func (e *Enforcer) RemoveFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *Enforcer) RemoveFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *Enforcer) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error)

func (e *Enforcer) RemoveGroupingPolicies(rules [][]string) (bool, error)

func (e *Enforcer) RemoveGroupingPolicy(params ...interface{}) (bool, error)

func (e *Enforcer) RemoveNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)

func (e *Enforcer) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (e *Enforcer) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error)

func (e *Enforcer) RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (e *Enforcer) RemovePolicies(rules [][]string) (bool, error)

func (e *Enforcer) RemovePolicy(params ...interface{}) (bool, error)

func (e *Enforcer) SavePolicy() error

func (e *Enforcer) SelfAddPolicies(sec string, ptype string, rules [][]string) (bool, error)

func (e *Enforcer) SelfAddPoliciesEx(sec string, ptype string, rules [][]string) (bool, error)

func (e *Enforcer) SelfAddPolicy(sec string, ptype string, rule []string) (bool, error)

func (e *Enforcer) SelfRemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *Enforcer) SelfRemovePolicies(sec string, ptype string, rules [][]string) (bool, error)

func (e *Enforcer) SelfRemovePolicy(sec string, ptype string, rule []string) (bool, error)

func (e *Enforcer) SelfUpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) (bool, error)

func (e *Enforcer) SelfUpdatePolicy(sec string, ptype string, oldRule, newRule []string) (bool, error)

func (e *Enforcer) SetAdapter(adapter persist.Adapter)

func (e *Enforcer) SetEffector(eft effector.Effector)

func (e *Enforcer) SetFieldIndex(ptype string, field string, index int)

func (e *Enforcer) SetLogger(logger log.Logger)

func (e *Enforcer) SetModel(m model.Model)

func (e *Enforcer) SetNamedDomainLinkConditionFuncParams(ptype, user, role, domain string, params ...string) bool

func (e *Enforcer) SetNamedLinkConditionFuncParams(ptype, user, role string, params ...string) bool

func (e *Enforcer) SetNamedRoleManager(ptype string, rm rbac.RoleManager)

func (e *Enforcer) SetRoleManager(rm rbac.RoleManager)

func (e *Enforcer) SetWatcher(watcher persist.Watcher) error

func (e *Enforcer) UpdateFilteredNamedPolicies(ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *Enforcer) UpdateFilteredPolicies(newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *Enforcer) UpdateGroupingPolicies(oldRules [][]string, newRules [][]string) (bool, error)

func (e *Enforcer) UpdateGroupingPolicy(oldRule []string, newRule []string) (bool, error)

func (e *Enforcer) UpdateNamedGroupingPolicies(ptype string, oldRules [][]string, newRules [][]string) (bool, error)

func (e *Enforcer) UpdateNamedGroupingPolicy(ptype string, oldRule []string, newRule []string) (bool, error)

func (e *Enforcer) UpdateNamedPolicies(ptype string, p1 [][]string, p2 [][]string) (bool, error)

func (e *Enforcer) UpdateNamedPolicy(ptype string, p1 []string, p2 []string) (bool, error)

func (e *Enforcer) UpdatePolicies(oldPolices [][]string, newPolicies [][]string) (bool, error)

func (e *Enforcer) UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error)

type IDistributedEnforcer

type IEnforcer

type IEnforcerContext

func NewContextEnforcer(params ...interface{}) (IEnforcerContext, error)

type SyncedCachedEnforcer

func NewSyncedCachedEnforcer(params ...interface{}) (*SyncedCachedEnforcer, error)

func (e *SyncedCachedEnforcer) AddPolicies(rules [][]string) (bool, error)

func (e *SyncedCachedEnforcer) AddPolicy(params ...interface{}) (bool, error)

func (e *SyncedCachedEnforcer) EnableCache(enableCache bool)

func (e *SyncedCachedEnforcer) Enforce(rvals ...interface{}) (bool, error)

func (e *SyncedCachedEnforcer) InvalidateCache() error

func (e *SyncedCachedEnforcer) LoadPolicy() error

func (e *SyncedCachedEnforcer) RemovePolicies(rules [][]string) (bool, error)

func (e *SyncedCachedEnforcer) RemovePolicy(params ...interface{}) (bool, error)

func (e *SyncedCachedEnforcer) SetCache(c cache.Cache)

func (e *SyncedCachedEnforcer) SetExpireTime(expireTime time.Duration)

type SyncedEnforcer

func NewSyncedEnforcer(params ...interface{}) (*SyncedEnforcer, error)

func (e *SyncedEnforcer) AddFunction(name string, function govaluate.ExpressionFunction)

func (e *SyncedEnforcer) AddGroupingPolicies(rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddGroupingPoliciesEx(rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddGroupingPolicy(params ...interface{}) (bool, error)

func (e *SyncedEnforcer) AddNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddNamedGroupingPoliciesEx(ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (e *SyncedEnforcer) AddNamedPolicies(ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddNamedPoliciesEx(ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (e *SyncedEnforcer) AddPermissionForUser(user string, permission ...string) (bool, error)

func (e *SyncedEnforcer) AddPermissionsForUser(user string, permissions ...[]string) (bool, error)

func (e *SyncedEnforcer) AddPolicies(rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddPoliciesEx(rules [][]string) (bool, error)

func (e *SyncedEnforcer) AddPolicy(params ...interface{}) (bool, error)

func (e *SyncedEnforcer) AddRoleForUser(user string, role string, domain ...string) (bool, error)

func (e *SyncedEnforcer) AddRoleForUserInDomain(user string, role string, domain string) (bool, error)

func (e *SyncedEnforcer) AddRolesForUser(user string, roles []string, domain ...string) (bool, error)

func (e *SyncedEnforcer) BatchEnforce(requests [][]interface{}) ([]bool, error)

func (e *SyncedEnforcer) BatchEnforceWithMatcher(matcher string, requests [][]interface{}) ([]bool, error)

func (e *SyncedEnforcer) BuildRoleLinks() error

func (e *SyncedEnforcer) ClearPolicy()

func (e *SyncedEnforcer) DeleteDomains(domains ...string) (bool, error)

func (e *SyncedEnforcer) DeletePermission(permission ...string) (bool, error)

func (e *SyncedEnforcer) DeletePermissionForUser(user string, permission ...string) (bool, error)

func (e *SyncedEnforcer) DeletePermissionsForUser(user string) (bool, error)

func (e *SyncedEnforcer) DeleteRole(role string) (bool, error)

func (e *SyncedEnforcer) DeleteRoleForUser(user string, role string, domain ...string) (bool, error)

func (e *SyncedEnforcer) DeleteRoleForUserInDomain(user string, role string, domain string) (bool, error)

func (e *SyncedEnforcer) DeleteRolesForUser(user string, domain ...string) (bool, error)

func (e *SyncedEnforcer) DeleteRolesForUserInDomain(user string, domain string) (bool, error)

func (e *SyncedEnforcer) DeleteUser(user string) (bool, error)

func (e *SyncedEnforcer) Enforce(rvals ...interface{}) (bool, error)

func (e *SyncedEnforcer) EnforceEx(rvals ...interface{}) (bool, []string, error)

func (e *SyncedEnforcer) EnforceExWithMatcher(matcher string, rvals ...interface{}) (bool, []string, error)

func (e *SyncedEnforcer) EnforceWithMatcher(matcher string, rvals ...interface{}) (bool, error)

func (e *SyncedEnforcer) GetAllActions() ([]string, error)

func (e *SyncedEnforcer) GetAllNamedActions(ptype string) ([]string, error)

func (e *SyncedEnforcer) GetAllNamedObjects(ptype string) ([]string, error)

func (e *SyncedEnforcer) GetAllNamedRoles(ptype string) ([]string, error)

func (e *SyncedEnforcer) GetAllNamedSubjects(ptype string) ([]string, error)

func (e *SyncedEnforcer) GetAllObjects() ([]string, error)

func (e *SyncedEnforcer) GetAllRoles() ([]string, error)

func (e *SyncedEnforcer) GetAllSubjects() ([]string, error)

func (e *SyncedEnforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *SyncedEnforcer) GetFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *SyncedEnforcer) GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *SyncedEnforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)

func (e *SyncedEnforcer) GetGroupingPolicy() ([][]string, error)

func (e *SyncedEnforcer) GetImplicitObjectPatternsForUser(user string, domain string, action string) ([]string, error)

func (e *SyncedEnforcer) GetImplicitPermissionsForUser(user string, domain ...string) ([][]string, error)

func (e *SyncedEnforcer) GetImplicitRolesForUser(name string, domain ...string) ([]string, error)

func (e *SyncedEnforcer) GetImplicitUsersForPermission(permission ...string) ([]string, error)

func (e *SyncedEnforcer) GetLock() *sync.RWMutex

func (e *SyncedEnforcer) GetNamedGroupingPolicy(ptype string) ([][]string, error)

func (e *SyncedEnforcer) GetNamedImplicitPermissionsForUser(ptype string, gtype string, user string, domain ...string) ([][]string, error)

func (e *SyncedEnforcer) GetNamedPermissionsForUser(ptype string, user string, domain ...string) ([][]string, error)

func (e *SyncedEnforcer) GetNamedPolicy(ptype string) ([][]string, error)

func (e *SyncedEnforcer) GetNamedRoleManager(ptype string) rbac.RoleManager

func (e *SyncedEnforcer) GetPermissionsForUser(user string, domain ...string) ([][]string, error)

func (e *SyncedEnforcer) GetPermissionsForUserInDomain(user string, domain string) [][]string

func (e *SyncedEnforcer) GetPolicy() ([][]string, error)

func (e *SyncedEnforcer) GetRoleManager() rbac.RoleManager

func (e *SyncedEnforcer) GetRolesForUser(name string, domain ...string) ([]string, error)

func (e *SyncedEnforcer) GetRolesForUserInDomain(name string, domain string) []string

func (e *SyncedEnforcer) GetUsersForRole(name string, domain ...string) ([]string, error)

func (e *SyncedEnforcer) GetUsersForRoleInDomain(name string, domain string) []string

func (e *SyncedEnforcer) HasGroupingPolicy(params ...interface{}) (bool, error)

func (e *SyncedEnforcer) HasNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (e *SyncedEnforcer) HasNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (e *SyncedEnforcer) HasPermissionForUser(user string, permission ...string) (bool, error)

func (e *SyncedEnforcer) HasPolicy(params ...interface{}) (bool, error)

func (e *SyncedEnforcer) HasRoleForUser(name string, role string, domain ...string) (bool, error)

func (e *SyncedEnforcer) IsAutoLoadingRunning() bool

func (e *SyncedEnforcer) LoadFilteredPolicy(filter interface{}) error

func (e *SyncedEnforcer) LoadIncrementalFilteredPolicy(filter interface{}) error

func (e *SyncedEnforcer) LoadModel() error

func (e *SyncedEnforcer) LoadPolicy() error

func (e *SyncedEnforcer) RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) (bool, error)

func (e *SyncedEnforcer) RemoveFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *SyncedEnforcer) RemoveFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *SyncedEnforcer) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error)

func (e *SyncedEnforcer) RemoveGroupingPolicies(rules [][]string) (bool, error)

func (e *SyncedEnforcer) RemoveGroupingPolicy(params ...interface{}) (bool, error)

func (e *SyncedEnforcer) RemoveNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (e *SyncedEnforcer) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (e *SyncedEnforcer) RemovePolicies(rules [][]string) (bool, error)

func (e *SyncedEnforcer) RemovePolicy(params ...interface{}) (bool, error)

func (e *SyncedEnforcer) SavePolicy() error

func (e *SyncedEnforcer) SelfAddPolicies(sec string, ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) SelfAddPoliciesEx(sec string, ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) SelfAddPolicy(sec string, ptype string, rule []string) (bool, error)

func (e *SyncedEnforcer) SelfRemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *SyncedEnforcer) SelfRemovePolicies(sec string, ptype string, rules [][]string) (bool, error)

func (e *SyncedEnforcer) SelfRemovePolicy(sec string, ptype string, rule []string) (bool, error)

func (e *SyncedEnforcer) SelfUpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) (bool, error)

func (e *SyncedEnforcer) SelfUpdatePolicy(sec string, ptype string, oldRule, newRule []string) (bool, error)

func (e *SyncedEnforcer) SetNamedRoleManager(ptype string, rm rbac.RoleManager)

func (e *SyncedEnforcer) SetRoleManager(rm rbac.RoleManager)

func (e *SyncedEnforcer) SetWatcher(watcher persist.Watcher) error

func (e *SyncedEnforcer) StartAutoLoadPolicy(d time.Duration)

func (e *SyncedEnforcer) StopAutoLoadPolicy()

func (e *SyncedEnforcer) UpdateFilteredNamedPolicies(ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *SyncedEnforcer) UpdateFilteredPolicies(newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)

func (e *SyncedEnforcer) UpdateGroupingPolicies(oldRules [][]string, newRules [][]string) (bool, error)

func (e *SyncedEnforcer) UpdateGroupingPolicy(oldRule []string, newRule []string) (bool, error)

func (e *SyncedEnforcer) UpdateNamedGroupingPolicies(ptype string, oldRules [][]string, newRules [][]string) (bool, error)

func (e *SyncedEnforcer) UpdateNamedGroupingPolicy(ptype string, oldRule []string, newRule []string) (bool, error)

func (e *SyncedEnforcer) UpdateNamedPolicies(ptype string, p1 [][]string, p2 [][]string) (bool, error)

func (e *SyncedEnforcer) UpdateNamedPolicy(ptype string, p1 []string, p2 []string) (bool, error)

func (e *SyncedEnforcer) UpdatePolicies(oldPolices [][]string, newPolicies [][]string) (bool, error)

func (e *SyncedEnforcer) UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error)

type Transaction

func (tx *Transaction) AddGroupingPolicy(params ...interface{}) (bool, error)

func (tx *Transaction) AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (tx *Transaction) AddNamedPolicies(ptype string, rules [][]string) (bool, error)

func (tx *Transaction) AddNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (tx *Transaction) AddPolicies(rules [][]string) (bool, error)

func (tx *Transaction) AddPolicy(params ...interface{}) (bool, error)

func (tx *Transaction) Commit() error

func (tx *Transaction) GetBufferedModel() (model.Model, error)

func (tx *Transaction) HasOperations() bool

func (tx *Transaction) IsActive() bool

func (tx *Transaction) IsCommitted() bool

func (tx *Transaction) IsRolledBack() bool

func (tx *Transaction) OperationCount() int

func (tx *Transaction) RemoveGroupingPolicy(params ...interface{}) (bool, error)

func (tx *Transaction) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)

func (tx *Transaction) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error)

func (tx *Transaction) RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error)

func (tx *Transaction) RemovePolicies(rules [][]string) (bool, error)

func (tx *Transaction) RemovePolicy(params ...interface{}) (bool, error)

func (tx *Transaction) Rollback() error

func (tx *Transaction) UpdateNamedPolicy(ptype string, oldPolicy []string, newPolicy []string) (bool, error)

func (tx *Transaction) UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error)

type TransactionBuffer

func NewTransactionBuffer(baseModel model.Model) *TransactionBuffer

func (tb *TransactionBuffer) AddOperation(op persist.PolicyOperation)

func (tb *TransactionBuffer) ApplyOperationsToModel(baseModel model.Model) (model.Model, error)

func (tb *TransactionBuffer) Clear()

func (tb *TransactionBuffer) GetModelSnapshot() model.Model

func (tb *TransactionBuffer) GetOperations() []persist.PolicyOperation

func (tb *TransactionBuffer) HasOperations() bool

func (tb *TransactionBuffer) OperationCount() int

type TransactionalEnforcer

func NewTransactionalEnforcer(params ...interface{}) (*TransactionalEnforcer, error)

func (te *TransactionalEnforcer) BeginTransaction(ctx context.Context) (*Transaction, error)

func (te *TransactionalEnforcer) GetTransaction(id string) *Transaction

func (te *TransactionalEnforcer) IsTransactionActive(id string) bool

func (te *TransactionalEnforcer) WithTransaction(ctx context.Context, fn func(*Transaction) error) error

Constants

¶

This section is empty.

Variables

¶

This section is empty.

Functions

¶

func

CasbinJsGetPermissionForUser

¶

added in

v2.9.0

func CasbinJsGetPermissionForUser(e

IEnforcer

, user

string

) (

string

,

error

)

func

CasbinJsGetPermissionForUserOld

¶

added in

v2.31.1

func CasbinJsGetPermissionForUserOld(e

IEnforcer

, user

string

) ([]

byte

,

error

)

func

GetCacheKey

¶

added in

v2.66.0

func GetCacheKey(params ...interface{}) (

string

,

bool

)

Types

¶

type

CacheableParam

¶

added in

v2.40.0

type CacheableParam interface {

GetCacheKey()

string

}

type

CachedEnforcer

¶

type CachedEnforcer struct {

*

Enforcer

// contains filtered or unexported fields

}

CachedEnforcer wraps Enforcer and provides decision cache.

func

NewCachedEnforcer

¶

func NewCachedEnforcer(params ...interface{}) (*

CachedEnforcer

,

error

)

NewCachedEnforcer creates a cached enforcer via file or DB.

func (*CachedEnforcer)

ClearPolicy

¶

added in

v2.97.0

func (e *

CachedEnforcer

) ClearPolicy()

ClearPolicy clears all policy.

func (*CachedEnforcer)

EnableCache

¶

func (e *

CachedEnforcer

) EnableCache(enableCache

bool

)

EnableCache determines whether to enable cache on Enforce(). When enableCache is enabled, cached result (true | false) will be returned for previous decisions.

func (*CachedEnforcer)

Enforce

¶

func (e *

CachedEnforcer

) Enforce(rvals ...interface{}) (

bool

,

error

)

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).
if rvals is not string , ignore the cache.

func (*CachedEnforcer)

InvalidateCache

¶

func (e *

CachedEnforcer

) InvalidateCache()

error

InvalidateCache deletes all the existing cached decisions.

func (*CachedEnforcer)

LoadPolicy

¶

added in

v2.32.0

func (e *

CachedEnforcer

) LoadPolicy()

error

func (*CachedEnforcer)

RemovePolicies

¶

added in

v2.32.0

func (e *

CachedEnforcer

) RemovePolicies(rules [][]

string

) (

bool

,

error

)

func (*CachedEnforcer)

RemovePolicy

¶

added in

v2.32.0

func (e *

CachedEnforcer

) RemovePolicy(params ...interface{}) (

bool

,

error

)

func (*CachedEnforcer)

SetCache

¶

added in

v2.32.0

func (e *

CachedEnforcer

) SetCache(c

cache

.

Cache

)

func (*CachedEnforcer)

SetExpireTime

¶

added in

v2.32.0

func (e *

CachedEnforcer

) SetExpireTime(expireTime

time

.

Duration

)

type

ConflictDetector

¶

added in

v2.128.0

type ConflictDetector struct {

// contains filtered or unexported fields

}

ConflictDetector detects conflicts between transaction operations and current model state.

func

NewConflictDetector

¶

added in

v2.128.0

func NewConflictDetector(baseModel, currentModel

model

.

Model

, operations []

persist

.

PolicyOperation

) *

ConflictDetector

NewConflictDetector creates a new conflict detector instance.

func (*ConflictDetector)

DetectConflicts

¶

added in

v2.128.0

func (cd *

ConflictDetector

) DetectConflicts()

error

DetectConflicts checks for conflicts between the transaction operations and current model state.
Returns nil if no conflicts are found, otherwise returns a ConflictError describing the conflict.

type

ConflictError

¶

added in

v2.128.0

type ConflictError struct {

Operation

persist

.

PolicyOperation

Reason

string

}

ConflictError represents a transaction conflict error.

func (*ConflictError)

Error

¶

added in

v2.128.0

func (e *

ConflictError

) Error()

string

type

ContextEnforcer

¶

added in

v2.128.0

type ContextEnforcer struct {

*

Enforcer

// contains filtered or unexported fields

}

ContextEnforcer wraps Enforcer and provides context-aware operations.

func (*ContextEnforcer)

AddGroupingPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddGroupingPoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

AddGroupingPoliciesCtx adds grouping policy rules to the storage with context.

func (*ContextEnforcer)

AddGroupingPoliciesExCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddGroupingPoliciesExCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

func (*ContextEnforcer)

AddGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddGroupingPolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

AddGroupingPolicyCtx adds a grouping policy rule to the storage with context.

func (*ContextEnforcer)

AddNamedGroupingPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddNamedGroupingPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPoliciesCtx adds named grouping policy rules to the storage with context.

func (*ContextEnforcer)

AddNamedGroupingPoliciesExCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddNamedGroupingPoliciesExCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*ContextEnforcer)

AddNamedGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedGroupingPolicyCtx adds a named grouping policy rule to the storage with context.

func (*ContextEnforcer)

AddNamedPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddNamedPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedPoliciesCtx adds named policy rules to the storage with context.

func (*ContextEnforcer)

AddNamedPoliciesExCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddNamedPoliciesExCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*ContextEnforcer)

AddNamedPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddNamedPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedPolicyCtx adds a named policy rule to the storage with context.

func (*ContextEnforcer)

AddPermissionForUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddPermissionForUserCtx(ctx

context

.

Context

, user

string

, permission ...

string

) (

bool

,

error

)

AddPermissionForUserCtx adds a permission for a user or role with context support.
Returns false if the user or role already has the permission (aka not affected).

func (*ContextEnforcer)

AddPermissionsForUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddPermissionsForUserCtx(ctx

context

.

Context

, user

string

, permissions ...[]

string

) (

bool

,

error

)

AddPermissionsForUserCtx adds multiple permissions for a user or role with context support.
Returns false if the user or role already has one of the permissions (aka not affected).

func (*ContextEnforcer)

AddPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddPoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

AddPoliciesCtx adds policy rules to the storage with context.

func (*ContextEnforcer)

AddPoliciesExCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddPoliciesExCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

func (*ContextEnforcer)

AddPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddPolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

AddPolicyCtx adds a policy rule to the storage with context.

func (*ContextEnforcer)

AddRoleForUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddRoleForUserCtx(ctx

context

.

Context

, user

string

, role

string

, domain ...

string

) (

bool

,

error

)

AddRoleForUserCtx adds a role for a user with context support.
Returns false if the user already has the role (aka not affected).

func (*ContextEnforcer)

AddRoleForUserInDomainCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) AddRoleForUserInDomainCtx(ctx

context

.

Context

, user

string

, role

string

, domain

string

) (

bool

,

error

)

AddRoleForUserInDomainCtx adds a role for a user inside a domain with context support.
Returns false if the user already has the role (aka not affected).

func (*ContextEnforcer)

DeleteAllUsersByDomainCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteAllUsersByDomainCtx(ctx

context

.

Context

, domain

string

) (

bool

,

error

)

DeleteAllUsersByDomainCtx deletes all users associated with the domain with context support.

func (*ContextEnforcer)

DeleteDomainsCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteDomainsCtx(ctx

context

.

Context

, domains ...

string

) (

bool

,

error

)

DeleteDomainsCtx deletes all associated policies for domains with context support.
It would delete all domains if parameter is not provided.

func (*ContextEnforcer)

DeletePermissionCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeletePermissionCtx(ctx

context

.

Context

, permission ...

string

) (

bool

,

error

)

DeletePermissionCtx deletes a permission with context support.
Returns false if the permission does not exist (aka not affected).

func (*ContextEnforcer)

DeletePermissionForUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeletePermissionForUserCtx(ctx

context

.

Context

, user

string

, permission ...

string

) (

bool

,

error

)

DeletePermissionForUserCtx deletes a permission for a user or role with context support.
Returns false if the user or role does not have the permission (aka not affected).

func (*ContextEnforcer)

DeletePermissionsForUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeletePermissionsForUserCtx(ctx

context

.

Context

, user

string

) (

bool

,

error

)

DeletePermissionsForUserCtx deletes permissions for a user or role with context support.
Returns false if the user or role does not have any permissions (aka not affected).

func (*ContextEnforcer)

DeleteRoleCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteRoleCtx(ctx

context

.

Context

, role

string

) (

bool

,

error

)

DeleteRoleCtx deletes a role with context support.
Returns false if the role does not exist (aka not affected).

func (*ContextEnforcer)

DeleteRoleForUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteRoleForUserCtx(ctx

context

.

Context

, user

string

, role

string

, domain ...

string

) (

bool

,

error

)

DeleteRoleForUserCtx deletes a role for a user with context support.
Returns false if the user does not have the role (aka not affected).

func (*ContextEnforcer)

DeleteRoleForUserInDomainCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteRoleForUserInDomainCtx(ctx

context

.

Context

, user

string

, role

string

, domain

string

) (

bool

,

error

)

DeleteRoleForUserInDomainCtx deletes a role for a user inside a domain with context support.
Returns false if the user does not have the role (aka not affected).

func (*ContextEnforcer)

DeleteRolesForUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteRolesForUserCtx(ctx

context

.

Context

, user

string

, domain ...

string

) (

bool

,

error

)

DeleteRolesForUserCtx deletes all roles for a user with context support.
Returns false if the user does not have any roles (aka not affected).

func (*ContextEnforcer)

DeleteRolesForUserInDomainCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteRolesForUserInDomainCtx(ctx

context

.

Context

, user

string

, domain

string

) (

bool

,

error

)

DeleteRolesForUserInDomainCtx deletes all roles for a user inside a domain with context support.
Returns false if the user does not have any roles (aka not affected).

func (*ContextEnforcer)

DeleteUserCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) DeleteUserCtx(ctx

context

.

Context

, user

string

) (

bool

,

error

)

DeleteUserCtx deletes a user with context support.
Returns false if the user does not exist (aka not affected).

func (*ContextEnforcer)

IsFilteredCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) IsFilteredCtx(ctx

context

.

Context

)

bool

IsFilteredCtx returns true if the loaded policy has been filtered with context.

func (*ContextEnforcer)

LoadPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) LoadPolicyCtx(ctx

context

.

Context

)

error

LoadPolicyCtx loads all policy rules from the storage with context.

func (*ContextEnforcer)

RemoveFilteredGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveFilteredGroupingPolicyCtx(ctx

context

.

Context

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredGroupingPolicyCtx removes grouping policy rules that match the filter from the storage with context.

func (*ContextEnforcer)

RemoveFilteredNamedGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveFilteredNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredNamedGroupingPolicyCtx removes named grouping policy rules that match the filter from the storage with context.

func (*ContextEnforcer)

RemoveFilteredNamedPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveFilteredNamedPolicyCtx(ctx

context

.

Context

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredNamedPolicyCtx removes named policy rules that match the filter from the storage with context.

func (*ContextEnforcer)

RemoveFilteredPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveFilteredPolicyCtx(ctx

context

.

Context

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredPolicyCtx removes policy rules that match the filter from the storage with context.

func (*ContextEnforcer)

RemoveGroupingPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveGroupingPoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

RemoveGroupingPoliciesCtx removes grouping policy rules from the storage with context.

func (*ContextEnforcer)

RemoveGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveGroupingPolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

RemoveGroupingPolicyCtx removes a grouping policy rule from the storage with context.

func (*ContextEnforcer)

RemoveNamedGroupingPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveNamedGroupingPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveNamedGroupingPoliciesCtx removes named grouping policy rules from the storage with context.

func (*ContextEnforcer)

RemoveNamedGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedGroupingPolicyCtx removes a named grouping policy rule from the storage with context.

func (*ContextEnforcer)

RemoveNamedPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveNamedPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveNamedPoliciesCtx removes named policy rules from the storage with context.

func (*ContextEnforcer)

RemoveNamedPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemoveNamedPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedPolicyCtx removes a named policy rule from the storage with context.

func (*ContextEnforcer)

RemovePoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemovePoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

RemovePoliciesCtx removes policy rules from the storage with context.

func (*ContextEnforcer)

RemovePolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) RemovePolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

RemovePolicyCtx removes a policy rule from the storage with context.

func (*ContextEnforcer)

SavePolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SavePolicyCtx(ctx

context

.

Context

)

error

func (*ContextEnforcer)

SelfAddPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfAddPoliciesCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfAddPoliciesCtx adds policy rules to the current policy with context.

func (*ContextEnforcer)

SelfAddPoliciesExCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfAddPoliciesExCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*ContextEnforcer)

SelfAddPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfAddPolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

SelfAddPolicyCtx adds a policy rule to the current policy with context.

func (*ContextEnforcer)

SelfRemoveFilteredPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfRemoveFilteredPolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

SelfRemoveFilteredPolicyCtx removes policy rules that match the filter from the current policy with context.

func (*ContextEnforcer)

SelfRemovePoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfRemovePoliciesCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfRemovePoliciesCtx removes policy rules from the current policy with context.

func (*ContextEnforcer)

SelfRemovePolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfRemovePolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

SelfRemovePolicyCtx removes a policy rule from the current policy with context.

func (*ContextEnforcer)

SelfUpdatePoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfUpdatePoliciesCtx(ctx

context

.

Context

, sec

string

, ptype

string

, oldRules, newRules [][]

string

) (

bool

,

error

)

SelfUpdatePoliciesCtx updates policy rules in the current policy with context.

func (*ContextEnforcer)

SelfUpdatePolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) SelfUpdatePolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, oldRule, newRule []

string

) (

bool

,

error

)

SelfUpdatePolicyCtx updates a policy rule in the current policy with context.

func (*ContextEnforcer)

UpdateFilteredNamedPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateFilteredNamedPoliciesCtx(ctx

context

.

Context

, ptype

string

, newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

UpdateFilteredNamedPoliciesCtx updates named policy rules that match the filter in the storage with context.

func (*ContextEnforcer)

UpdateFilteredPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateFilteredPoliciesCtx(ctx

context

.

Context

, newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

UpdateFilteredPoliciesCtx updates policy rules that match the filter in the storage with context.

func (*ContextEnforcer)

UpdateGroupingPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateGroupingPoliciesCtx(ctx

context

.

Context

, oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

UpdateGroupingPoliciesCtx updates grouping policy rules in the storage with context.

func (*ContextEnforcer)

UpdateGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateGroupingPolicyCtx(ctx

context

.

Context

, oldRule []

string

, newRule []

string

) (

bool

,

error

)

UpdateGroupingPolicyCtx updates a grouping policy rule in the storage with context.

func (*ContextEnforcer)

UpdateNamedGroupingPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateNamedGroupingPoliciesCtx(ctx

context

.

Context

, ptype

string

, oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

UpdateNamedGroupingPoliciesCtx updates named grouping policy rules in the storage with context.

func (*ContextEnforcer)

UpdateNamedGroupingPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, oldRule []

string

, newRule []

string

) (

bool

,

error

)

UpdateNamedGroupingPolicyCtx updates a named grouping policy rule in the storage with context.

func (*ContextEnforcer)

UpdateNamedPoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateNamedPoliciesCtx(ctx

context

.

Context

, ptype

string

, p1 [][]

string

, p2 [][]

string

) (

bool

,

error

)

UpdateNamedPoliciesCtx updates named policy rules in the storage with context.

func (*ContextEnforcer)

UpdateNamedPolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdateNamedPolicyCtx(ctx

context

.

Context

, ptype

string

, p1 []

string

, p2 []

string

) (

bool

,

error

)

UpdateNamedPolicyCtx updates a named policy rule in the storage with context.

func (*ContextEnforcer)

UpdatePoliciesCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdatePoliciesCtx(ctx

context

.

Context

, oldPolicies [][]

string

, newPolicies [][]

string

) (

bool

,

error

)

UpdatePoliciesCtx updates policy rules in the storage with context.

func (*ContextEnforcer)

UpdatePolicyCtx

¶

added in

v2.128.0

func (e *

ContextEnforcer

) UpdatePolicyCtx(ctx

context

.

Context

, oldPolicy []

string

, newPolicy []

string

) (

bool

,

error

)

UpdatePolicyCtx updates a policy rule in the storage with context.

type

DistributedEnforcer

¶

added in

v2.19.0

type DistributedEnforcer struct {

*

SyncedEnforcer

}

DistributedEnforcer wraps SyncedEnforcer for dispatcher.

func

NewDistributedEnforcer

¶

added in

v2.19.0

func NewDistributedEnforcer(params ...interface{}) (*

DistributedEnforcer

,

error

)

func (*DistributedEnforcer)

AddPoliciesSelf

¶

added in

v2.23.3

func (d *

DistributedEnforcer

) AddPoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, rules [][]

string

) (affected [][]

string

, err

error

)

AddPoliciesSelf provides a method for dispatcher to add authorization rules to the current policy.
The function returns the rules affected and error.

func (*DistributedEnforcer)

ClearPolicySelf

¶

added in

v2.19.0

func (d *

DistributedEnforcer

) ClearPolicySelf(shouldPersist func()

bool

)

error

ClearPolicySelf provides a method for dispatcher to clear all rules from the current policy.

func (*DistributedEnforcer)

RemoveFilteredPolicySelf

¶

added in

v2.19.0

func (d *

DistributedEnforcer

) RemoveFilteredPolicySelf(shouldPersist func()

bool

, sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (affected [][]

string

, err

error

)

RemoveFilteredPolicySelf provides a method for dispatcher to remove an authorization rule from the current policy, field filters can be specified.
The function returns the rules affected and error.

func (*DistributedEnforcer)

RemovePoliciesSelf

¶

added in

v2.23.3

func (d *

DistributedEnforcer

) RemovePoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, rules [][]

string

) (affected [][]

string

, err

error

)

RemovePoliciesSelf provides a method for dispatcher to remove a set of rules from current policy.
The function returns the rules affected and error.

func (*DistributedEnforcer)

SetDispatcher

¶

added in

v2.20.1

func (d *

DistributedEnforcer

) SetDispatcher(dispatcher

persist

.

Dispatcher

)

SetDispatcher sets the current dispatcher.

func (*DistributedEnforcer)

UpdateFilteredPoliciesSelf

¶

added in

v2.28.0

func (d *

DistributedEnforcer

) UpdateFilteredPoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, newRules [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

UpdateFilteredPoliciesSelf provides a method for dispatcher to update a set of authorization rules from the current policy.

func (*DistributedEnforcer)

UpdatePoliciesSelf

¶

added in

v2.23.3

func (d *

DistributedEnforcer

) UpdatePoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, oldRules, newRules [][]

string

) (affected

bool

, err

error

)

UpdatePoliciesSelf provides a method for dispatcher to update a set of authorization rules from the current policy.

func (*DistributedEnforcer)

UpdatePolicySelf

¶

added in

v2.19.0

func (d *

DistributedEnforcer

) UpdatePolicySelf(shouldPersist func()

bool

, sec

string

, ptype

string

, oldRule, newRule []

string

) (affected

bool

, err

error

)

UpdatePolicySelf provides a method for dispatcher to update an authorization rule from the current policy.

type

EnforceContext

¶

added in

v2.36.0

type EnforceContext struct {

RType

string

PType

string

EType

string

MType

string

}

EnforceContext is used as the first element of the parameter "rvals" in method "enforce".

func

NewEnforceContext

¶

added in

v2.36.0

func NewEnforceContext(suffix

string

)

EnforceContext

NewEnforceContext Create a default structure based on the suffix.

func (EnforceContext)

GetCacheKey

¶

added in

v2.71.1

func (e

EnforceContext

) GetCacheKey()

string

type

Enforcer

¶

type Enforcer struct {

// contains filtered or unexported fields

}

Enforcer is the main interface for authorization enforcement and policy management.

func

NewEnforcer

¶

func NewEnforcer(params ...interface{}) (*

Enforcer

,

error

)

NewEnforcer creates an enforcer via file or DB.

File:

e := casbin.NewEnforcer("path/to/basic_model.conf", "path/to/basic_policy.csv")

MySQL DB:

a := mysqladapter.NewDBAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
e := casbin.NewEnforcer("path/to/basic_model.conf", a)

func (*Enforcer)

AddFunction

¶

func (e *

Enforcer

) AddFunction(name

string

, function

govaluate

.

ExpressionFunction

)

AddFunction adds a customized function.

func (*Enforcer)

AddGroupingPolicies

¶

added in

v2.2.2

func (e *

Enforcer

) AddGroupingPolicies(rules [][]

string

) (

bool

,

error

)

AddGroupingPolicies adds role inheritance rules to the current policy.
If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added.
Otherwise the function returns true for the corresponding policy rule by adding the new rule.

func (*Enforcer)

AddGroupingPoliciesEx

¶

added in

v2.63.0

func (e *

Enforcer

) AddGroupingPoliciesEx(rules [][]

string

) (

bool

,

error

)

AddGroupingPoliciesEx adds role inheritance rules to the current policy.
If the rule already exists, the rule will not be added.
But unlike AddGroupingPolicies, other non-existent rules are added instead of returning false directly.

func (*Enforcer)

AddGroupingPolicy

¶

func (e *

Enforcer

) AddGroupingPolicy(params ...interface{}) (

bool

,

error

)

AddGroupingPolicy adds a role inheritance rule to the current policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*Enforcer)

AddNamedDomainLinkConditionFunc

¶

added in

v2.75.0

func (e *

Enforcer

) AddNamedDomainLinkConditionFunc(ptype, user, role

string

, domain

string

, fn

rbac

.

LinkConditionFunc

)

bool

AddNamedDomainLinkConditionFunc Add condition function fn for Link userName-> {roleName, domain},
when fn returns true, Link is valid, otherwise invalid.

func (*Enforcer)

AddNamedDomainMatchingFunc

¶

added in

v2.21.0

func (e *

Enforcer

) AddNamedDomainMatchingFunc(ptype, name

string

, fn

rbac

.

MatchingFunc

)

bool

AddNamedDomainMatchingFunc add MatchingFunc by ptype to RoleManager.

func (*Enforcer)

AddNamedGroupingPolicies

¶

added in

v2.2.2

func (e *

Enforcer

) AddNamedGroupingPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPolicies adds named role inheritance rules to the current policy.
If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added.
Otherwise the function returns true for the corresponding policy rule by adding the new rule.

func (*Enforcer)

AddNamedGroupingPoliciesEx

¶

added in

v2.63.0

func (e *

Enforcer

) AddNamedGroupingPoliciesEx(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPoliciesEx adds named role inheritance rules to the current policy.
If the rule already exists, the rule will not be added.
But unlike AddNamedGroupingPolicies, other non-existent rules are added instead of returning false directly.

func (*Enforcer)

AddNamedGroupingPolicy

¶

func (e *

Enforcer

) AddNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedGroupingPolicy adds a named role inheritance rule to the current policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*Enforcer)

AddNamedLinkConditionFunc

¶

added in

v2.75.0

func (e *

Enforcer

) AddNamedLinkConditionFunc(ptype, user, role

string

, fn

rbac

.

LinkConditionFunc

)

bool

AddNamedLinkConditionFunc Add condition function fn for Link userName->roleName,
when fn returns true, Link is valid, otherwise invalid.

func (*Enforcer)

AddNamedMatchingFunc

¶

added in

v2.21.0

func (e *

Enforcer

) AddNamedMatchingFunc(ptype, name

string

, fn

rbac

.

MatchingFunc

)

bool

AddNamedMatchingFunc add MatchingFunc by ptype RoleManager.

func (*Enforcer)

AddNamedPolicies

¶

added in

v2.2.2

func (e *

Enforcer

) AddNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedPolicies adds authorization rules to the current named policy.
If the rule already exists, the function returns false for the corresponding rule and the rule will not be added.
Otherwise the function returns true for the corresponding by adding the new rule.

func (*Enforcer)

AddNamedPoliciesEx

¶

added in

v2.63.0

func (e *

Enforcer

) AddNamedPoliciesEx(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedPoliciesEx adds authorization rules to the current named policy.
If the rule already exists, the rule will not be added.
But unlike AddNamedPolicies, other non-existent rules are added instead of returning false directly.

func (*Enforcer)

AddNamedPolicy

¶

func (e *

Enforcer

) AddNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedPolicy adds an authorization rule to the current named policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*Enforcer)

AddPermissionForUser

¶

func (e *

Enforcer

) AddPermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

AddPermissionForUser adds a permission for a user or role.
Returns false if the user or role already has the permission (aka not affected).

func (*Enforcer)

AddPermissionsForUser

¶

added in

v2.38.0

func (e *

Enforcer

) AddPermissionsForUser(user

string

, permissions ...[]

string

) (

bool

,

error

)

AddPermissionsForUser adds multiple permissions for a user or role.
Returns false if the user or role already has one of the permissions (aka not affected).

func (*Enforcer)

AddPolicies

¶

added in

v2.2.2

func (e *

Enforcer

) AddPolicies(rules [][]

string

) (

bool

,

error

)

AddPolicies adds authorization rules to the current policy.
If the rule already exists, the function returns false for the corresponding rule and the rule will not be added.
Otherwise the function returns true for the corresponding rule by adding the new rule.

func (*Enforcer)

AddPoliciesEx

¶

added in

v2.63.0

func (e *

Enforcer

) AddPoliciesEx(rules [][]

string

) (

bool

,

error

)

AddPoliciesEx adds authorization rules to the current policy.
If the rule already exists, the rule will not be added.
But unlike AddPolicies, other non-existent rules are added instead of returning false directly.

func (*Enforcer)

AddPolicy

¶

func (e *

Enforcer

) AddPolicy(params ...interface{}) (

bool

,

error

)

AddPolicy adds an authorization rule to the current policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*Enforcer)

AddRoleForUser

¶

func (e *

Enforcer

) AddRoleForUser(user

string

, role

string

, domain ...

string

) (

bool

,

error

)

AddRoleForUser adds a role for a user.
Returns false if the user already has the role (aka not affected).

func (*Enforcer)

AddRoleForUserInDomain

¶

func (e *

Enforcer

) AddRoleForUserInDomain(user

string

, role

string

, domain

string

) (

bool

,

error

)

AddRoleForUserInDomain adds a role for a user inside a domain.
Returns false if the user already has the role (aka not affected).

func (*Enforcer)

AddRolesForUser

¶

added in

v2.5.0

func (e *

Enforcer

) AddRolesForUser(user

string

, roles []

string

, domain ...

string

) (

bool

,

error

)

AddRolesForUser adds roles for a user.
Returns false if the user already has the roles (aka not affected).

func (*Enforcer)

BatchEnforce

¶

added in

v2.23.0

func (e *

Enforcer

) BatchEnforce(requests [][]interface{}) ([]

bool

,

error

)

BatchEnforce enforce in batches.

func (*Enforcer)

BatchEnforceWithMatcher

¶

added in

v2.23.0

func (e *

Enforcer

) BatchEnforceWithMatcher(matcher

string

, requests [][]interface{}) ([]

bool

,

error

)

BatchEnforceWithMatcher enforce with matcher in batches.

func (*Enforcer)

BuildIncrementalConditionalRoleLinks

¶

added in

v2.75.0

func (e *

Enforcer

) BuildIncrementalConditionalRoleLinks(op

model

.

PolicyOp

, ptype

string

, rules [][]

string

)

error

BuildIncrementalConditionalRoleLinks provides incremental build the role inheritance relations with conditions.

func (*Enforcer)

BuildIncrementalRoleLinks

¶

added in

v2.6.0

func (e *

Enforcer

) BuildIncrementalRoleLinks(op

model

.

PolicyOp

, ptype

string

, rules [][]

string

)

error

BuildIncrementalRoleLinks provides incremental build the role inheritance relations.

func (*Enforcer)

BuildRoleLinks

¶

func (e *

Enforcer

) BuildRoleLinks()

error

BuildRoleLinks manually rebuild the role inheritance relations.

func (*Enforcer)

ClearPolicy

¶

func (e *

Enforcer

) ClearPolicy()

ClearPolicy clears all policy.

func (*Enforcer)

DeleteAllUsersByDomain

¶

added in

v2.29.0

func (e *

Enforcer

) DeleteAllUsersByDomain(domain

string

) (

bool

,

error

)

DeleteAllUsersByDomain would delete all users associated with the domain.

func (*Enforcer)

DeleteDomains

¶

added in

v2.29.0

func (e *

Enforcer

) DeleteDomains(domains ...

string

) (

bool

,

error

)

DeleteDomains would delete all associated policies.
It would delete all domains if parameter is not provided.

func (*Enforcer)

DeletePermission

¶

func (e *

Enforcer

) DeletePermission(permission ...

string

) (

bool

,

error

)

DeletePermission deletes a permission.
Returns false if the permission does not exist (aka not affected).

func (*Enforcer)

DeletePermissionForUser

¶

func (e *

Enforcer

) DeletePermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

DeletePermissionForUser deletes a permission for a user or role.
Returns false if the user or role does not have the permission (aka not affected).

func (*Enforcer)

DeletePermissionsForUser

¶

func (e *

Enforcer

) DeletePermissionsForUser(user

string

) (

bool

,

error

)

DeletePermissionsForUser deletes permissions for a user or role.
Returns false if the user or role does not have any permissions (aka not affected).

func (*Enforcer)

DeleteRole

¶

func (e *

Enforcer

) DeleteRole(role

string

) (

bool

,

error

)

DeleteRole deletes a role.
Returns false if the role does not exist (aka not affected).

func (*Enforcer)

DeleteRoleForUser

¶

func (e *

Enforcer

) DeleteRoleForUser(user

string

, role

string

, domain ...

string

) (

bool

,

error

)

DeleteRoleForUser deletes a role for a user.
Returns false if the user does not have the role (aka not affected).

func (*Enforcer)

DeleteRoleForUserInDomain

¶

func (e *

Enforcer

) DeleteRoleForUserInDomain(user

string

, role

string

, domain

string

) (

bool

,

error

)

DeleteRoleForUserInDomain deletes a role for a user inside a domain.
Returns false if the user does not have the role (aka not affected).

func (*Enforcer)

DeleteRolesForUser

¶

func (e *

Enforcer

) DeleteRolesForUser(user

string

, domain ...

string

) (

bool

,

error

)

DeleteRolesForUser deletes all roles for a user.
Returns false if the user does not have any roles (aka not affected).

func (*Enforcer)

DeleteRolesForUserInDomain

¶

added in

v2.8.4

func (e *

Enforcer

) DeleteRolesForUserInDomain(user

string

, domain

string

) (

bool

,

error

)

DeleteRolesForUserInDomain deletes all roles for a user inside a domain.
Returns false if the user does not have any roles (aka not affected).

func (*Enforcer)

DeleteUser

¶

func (e *

Enforcer

) DeleteUser(user

string

) (

bool

,

error

)

DeleteUser deletes a user.
Returns false if the user does not exist (aka not affected).

func (*Enforcer)

EnableAcceptJsonRequest

¶

added in

v2.72.0

func (e *

Enforcer

) EnableAcceptJsonRequest(acceptJsonRequest

bool

)

EnableAcceptJsonRequest controls whether to accept json as a request parameter.

func (*Enforcer)

EnableAutoBuildRoleLinks

¶

func (e *

Enforcer

) EnableAutoBuildRoleLinks(autoBuildRoleLinks

bool

)

EnableAutoBuildRoleLinks controls whether to rebuild the role inheritance relations when a role is added or deleted.

func (*Enforcer)

EnableAutoNotifyDispatcher

¶

added in

v2.18.0

func (e *

Enforcer

) EnableAutoNotifyDispatcher(enable

bool

)

EnableAutoNotifyDispatcher controls whether to save a policy rule automatically notify the Dispatcher when it is added or removed.

func (*Enforcer)

EnableAutoNotifyWatcher

¶

added in

v2.2.1

func (e *

Enforcer

) EnableAutoNotifyWatcher(enable

bool

)

EnableAutoNotifyWatcher controls whether to save a policy rule automatically notify the Watcher when it is added or removed.

func (*Enforcer)

EnableAutoSave

¶

func (e *

Enforcer

) EnableAutoSave(autoSave

bool

)

EnableAutoSave controls whether to save a policy rule automatically to the adapter when it is added or removed.

func (*Enforcer)

EnableEnforce

¶

func (e *

Enforcer

) EnableEnforce(enable

bool

)

EnableEnforce changes the enforcing state of Casbin, when Casbin is disabled, all access will be allowed by the Enforce() function.

func (*Enforcer)

EnableLog

¶

func (e *

Enforcer

) EnableLog(enable

bool

)

EnableLog changes whether Casbin will log messages to the Logger.

func (*Enforcer)

Enforce

¶

func (e *

Enforcer

) Enforce(rvals ...interface{}) (

bool

,

error

)

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).

func (*Enforcer)

EnforceEx

¶

added in

v2.4.1

func (e *

Enforcer

) EnforceEx(rvals ...interface{}) (

bool

, []

string

,

error

)

EnforceEx explain enforcement by informing matched rules.

func (*Enforcer)

EnforceExWithMatcher

¶

added in

v2.4.1

func (e *

Enforcer

) EnforceExWithMatcher(matcher

string

, rvals ...interface{}) (

bool

, []

string

,

error

)

EnforceExWithMatcher use a custom matcher and explain enforcement by informing matched rules.

func (*Enforcer)

EnforceWithMatcher

¶

added in

v2.0.2

func (e *

Enforcer

) EnforceWithMatcher(matcher

string

, rvals ...interface{}) (

bool

,

error

)

EnforceWithMatcher use a custom matcher to decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (matcher, sub, obj, act), use model matcher by default when matcher is "".

func (*Enforcer)

GetAdapter

¶

func (e *

Enforcer

) GetAdapter()

persist

.

Adapter

GetAdapter gets the current adapter.

func (*Enforcer)

GetAllActions

¶

func (e *

Enforcer

) GetAllActions() ([]

string

,

error

)

GetAllActions gets the list of actions that show up in the current policy.

func (*Enforcer)

GetAllDomains

¶

added in

v2.43.0

func (e *

Enforcer

) GetAllDomains() ([]

string

,

error

)

GetAllDomains would get all domains.

func (*Enforcer)

GetAllNamedActions

¶

func (e *

Enforcer

) GetAllNamedActions(ptype

string

) ([]

string

,

error

)

GetAllNamedActions gets the list of actions that show up in the current named policy.

func (*Enforcer)

GetAllNamedObjects

¶

func (e *

Enforcer

) GetAllNamedObjects(ptype

string

) ([]

string

,

error

)

GetAllNamedObjects gets the list of objects that show up in the current named policy.

func (*Enforcer)

GetAllNamedRoles

¶

func (e *

Enforcer

) GetAllNamedRoles(ptype

string

) ([]

string

,

error

)

GetAllNamedRoles gets the list of roles that show up in the current named policy.

func (*Enforcer)

GetAllNamedSubjects

¶

func (e *

Enforcer

) GetAllNamedSubjects(ptype

string

) ([]

string

,

error

)

GetAllNamedSubjects gets the list of subjects that show up in the current named policy.

func (*Enforcer)

GetAllObjects

¶

func (e *

Enforcer

) GetAllObjects() ([]

string

,

error

)

GetAllObjects gets the list of objects that show up in the current policy.

func (*Enforcer)

GetAllRoles

¶

func (e *

Enforcer

) GetAllRoles() ([]

string

,

error

)

GetAllRoles gets the list of roles that show up in the current policy.

func (*Enforcer)

GetAllRolesByDomain

¶

added in

v2.61.0

func (e *

Enforcer

) GetAllRolesByDomain(domain

string

) ([]

string

,

error

)

GetAllRolesByDomain would get all roles associated with the domain.
note: Not applicable to Domains with inheritance relationship  (implicit roles)

func (*Enforcer)

GetAllSubjects

¶

func (e *

Enforcer

) GetAllSubjects() ([]

string

,

error

)

GetAllSubjects gets the list of subjects that show up in the current policy.

func (*Enforcer)

GetAllUsersByDomain

¶

added in

v2.29.0

func (e *

Enforcer

) GetAllUsersByDomain(domain

string

) ([]

string

,

error

)

GetAllUsersByDomain would get all users associated with the domain.

func (*Enforcer)

GetAllowedObjectConditions

¶

added in

v2.68.0

func (e *

Enforcer

) GetAllowedObjectConditions(user

string

, action

string

, prefix

string

) ([]

string

,

error

)

GetAllowedObjectConditions returns a string array of object conditions that the user can access.
For example: conditions, err := e.GetAllowedObjectConditions("alice", "read", "r.obj.")
Note:

0. prefix: You can customize the prefix of the object conditions, and "r.obj." is commonly used as a prefix.
After removing the prefix, the remaining part is the condition of the object.
If there is an obj policy that does not meet the prefix requirement, an errors.ERR_OBJ_CONDITION will be returned.

1. If the 'objectConditions' array is empty, return errors.ERR_EMPTY_CONDITION
This error is returned because some data adapters' ORM return full table data by default
when they receive an empty condition, which tends to behave contrary to expectations.(e.g. GORM)
If you are using an adapter that does not behave like this, you can choose to ignore this error.

func (*Enforcer)

GetDomainsForUser

¶

added in

v2.26.0

func (e *

Enforcer

) GetDomainsForUser(user

string

) ([]

string

,

error

)

GetDomainsForUser gets all domains.

func (*Enforcer)

GetFieldIndex

¶

added in

v2.48.0

func (e *

Enforcer

) GetFieldIndex(ptype

string

, field

string

) (

int

,

error

)

func (*Enforcer)

GetFilteredGroupingPolicy

¶

func (e *

Enforcer

) GetFilteredGroupingPolicy(fieldIndex

int

, fieldValues ...

string

) ([][]

string

,

error

)

GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

func (*Enforcer)

GetFilteredNamedGroupingPolicy

¶

func (e *

Enforcer

) GetFilteredNamedGroupingPolicy(ptype

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

GetFilteredNamedGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

func (*Enforcer)

GetFilteredNamedPolicy

¶

func (e *

Enforcer

) GetFilteredNamedPolicy(ptype

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

GetFilteredNamedPolicy gets all the authorization rules in the named policy, field filters can be specified.

func (*Enforcer)

GetFilteredNamedPolicyWithMatcher

¶

added in

v2.47.0

func (e *

Enforcer

) GetFilteredNamedPolicyWithMatcher(ptype

string

, matcher

string

) ([][]

string

,

error

)

GetFilteredNamedPolicyWithMatcher gets rules based on matcher from the policy.

func (*Enforcer)

GetFilteredPolicy

¶

func (e *

Enforcer

) GetFilteredPolicy(fieldIndex

int

, fieldValues ...

string

) ([][]

string

,

error

)

GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.

func (*Enforcer)

GetGroupingPolicy

¶

func (e *

Enforcer

) GetGroupingPolicy() ([][]

string

,

error

)

GetGroupingPolicy gets all the role inheritance rules in the policy.

func (*Enforcer)

GetImplicitObjectPatternsForUser

¶

added in

v2.121.0

func (e *

Enforcer

) GetImplicitObjectPatternsForUser(user

string

, domain

string

, action

string

) ([]

string

,

error

)

GetImplicitObjectPatternsForUser returns all object patterns (with wildcards) that a user has for a given domain and action.
For example:
p, admin, chronicle/123, location/*, read
p, user, chronicle/456, location/789, read
g, alice, admin
g, bob, user

GetImplicitObjectPatternsForUser("alice", "chronicle/123", "read") will return ["location/*"].
GetImplicitObjectPatternsForUser("bob", "chronicle/456", "read") will return ["location/789"].

func (*Enforcer)

GetImplicitPermissionsForUser

¶

func (e *

Enforcer

) GetImplicitPermissionsForUser(user

string

, domain ...

string

) ([][]

string

,

error

)

GetImplicitPermissionsForUser gets implicit permissions for a user or role.
Compared to GetPermissionsForUser(), this function retrieves permissions for inherited roles.
For example:
p, admin, data1, read
p, alice, data2, read
g, alice, admin

GetPermissionsForUser("alice") can only get: [["alice", "data2", "read"]].
But GetImplicitPermissionsForUser("alice") will get: [["admin", "data1", "read"], ["alice", "data2", "read"]].

func (*Enforcer)

GetImplicitResourcesForUser

¶

added in

v2.31.0

func (e *

Enforcer

) GetImplicitResourcesForUser(user

string

, domain ...

string

) ([][]

string

,

error

)

GetImplicitResourcesForUser returns all policies that user obtaining in domain.

func (*Enforcer)

GetImplicitRolesForUser

¶

func (e *

Enforcer

) GetImplicitRolesForUser(name

string

, domain ...

string

) ([]

string

,

error

)

GetImplicitRolesForUser gets implicit roles that a user has.
Compared to GetRolesForUser(), this function retrieves indirect roles besides direct roles.
For example:
g, alice, role:admin
g, role:admin, role:user

GetRolesForUser("alice") can only get: ["role:admin"].
But GetImplicitRolesForUser("alice") will get: ["role:admin", "role:user"].

func (*Enforcer)

GetImplicitUsersForPermission

¶

func (e *

Enforcer

) GetImplicitUsersForPermission(permission ...

string

) ([]

string

,

error

)

GetImplicitUsersForPermission gets implicit users for a permission.
For example:
p, admin, data1, read
p, bob, data1, read
g, alice, admin

GetImplicitUsersForPermission("data1", "read") will get: ["alice", "bob"].
Note: only users will be returned, roles (2nd arg in "g") will be excluded.

func (*Enforcer)

GetImplicitUsersForResource

¶

added in

v2.69.0

func (e *

Enforcer

) GetImplicitUsersForResource(resource

string

) ([][]

string

,

error

)

GetImplicitUsersForResource return implicit user based on resource.
for example:
p, alice, data1, read
p, bob, data2, write
p, data2_admin, data2, read
p, data2_admin, data2, write
g, alice, data2_admin
GetImplicitUsersForResource("data2") will return [[bob data2 write] [alice data2 read] [alice data2 write]]
GetImplicitUsersForResource("data1") will return [[alice data1 read]]
Note: only users will be returned, roles (2nd arg in "g") will be excluded.

func (*Enforcer)

GetImplicitUsersForResourceByDomain

¶

added in

v2.69.0

func (e *

Enforcer

) GetImplicitUsersForResourceByDomain(resource

string

, domain

string

) ([][]

string

,

error

)

GetImplicitUsersForResourceByDomain return implicit user based on resource and domain.
Compared to GetImplicitUsersForResource, domain is supported.

func (*Enforcer)

GetImplicitUsersForRole

¶

added in

v2.31.0

func (e *

Enforcer

) GetImplicitUsersForRole(name

string

, domain ...

string

) ([]

string

,

error

)

GetImplicitUsersForRole gets implicit users for a role.

func (*Enforcer)

GetModel

¶

func (e *

Enforcer

) GetModel()

model

.

Model

GetModel gets the current model.

func (*Enforcer)

GetNamedGroupingPolicy

¶

func (e *

Enforcer

) GetNamedGroupingPolicy(ptype

string

) ([][]

string

,

error

)

GetNamedGroupingPolicy gets all the role inheritance rules in the policy.

func (*Enforcer)

GetNamedImplicitPermissionsForUser

¶

added in

v2.45.0

func (e *

Enforcer

) GetNamedImplicitPermissionsForUser(ptype

string

, gtype

string

, user

string

, domain ...

string

) ([][]

string

,

error

)

GetNamedImplicitPermissionsForUser gets implicit permissions for a user or role by named policy.
Compared to GetNamedPermissionsForUser(), this function retrieves permissions for inherited roles.
For example:
p, admin, data1, read
p2, admin, create
g, alice, admin

GetImplicitPermissionsForUser("alice") can only get: [["admin", "data1", "read"]], whose policy is default policy "p"
But you can specify the named policy "p2" to get: [["admin", "create"]] by    GetNamedImplicitPermissionsForUser("p2","alice").

func (*Enforcer)

GetNamedImplicitRolesForUser

¶

added in

v2.95.0

func (e *

Enforcer

) GetNamedImplicitRolesForUser(ptype

string

, name

string

, domain ...

string

) ([]

string

,

error

)

GetNamedImplicitRolesForUser gets implicit roles that a user has by named role definition.
Compared to GetImplicitRolesForUser(), this function retrieves indirect roles besides direct roles.
For example:
g, alice, role:admin
g, role:admin, role:user
g2, alice, role:admin2

GetImplicitRolesForUser("alice") can only get: ["role:admin", "role:user"].
But GetNamedImplicitRolesForUser("g2", "alice") will get: ["role:admin2"].

func (*Enforcer)

GetNamedImplicitUsersForResource

¶

added in

v2.120.0

func (e *

Enforcer

) GetNamedImplicitUsersForResource(ptype

string

, resource

string

) ([][]

string

,

error

)

GetNamedImplicitUsersForResource return implicit user based on resource with named policy support.
This function handles resource role relationships through named policies (e.g., g2, g3, etc.).
for example:
p, admin_group, admin_data, *
g, admin, admin_group
g2, app, admin_data
GetNamedImplicitUsersForResource("g2", "app") will return users who have access to admin_data through g2 relationship.

func (*Enforcer)

GetNamedPermissionsForUser

¶

added in

v2.45.0

func (e *

Enforcer

) GetNamedPermissionsForUser(ptype

string

, user

string

, domain ...

string

) ([][]

string

,

error

)

GetNamedPermissionsForUser gets permissions for a user or role by named policy.

func (*Enforcer)

GetNamedPolicy

¶

func (e *

Enforcer

) GetNamedPolicy(ptype

string

) ([][]

string

,

error

)

GetNamedPolicy gets all the authorization rules in the named policy.

func (*Enforcer)

GetNamedRoleManager

¶

added in

v2.52.0

func (e *

Enforcer

) GetNamedRoleManager(ptype

string

)

rbac

.

RoleManager

GetNamedRoleManager gets the role manager for the named policy.

func (*Enforcer)

GetPermissionsForUser

¶

func (e *

Enforcer

) GetPermissionsForUser(user

string

, domain ...

string

) ([][]

string

,

error

)

GetPermissionsForUser gets permissions for a user or role.

func (*Enforcer)

GetPermissionsForUserInDomain

¶

func (e *

Enforcer

) GetPermissionsForUserInDomain(user

string

, domain

string

) [][]

string

GetPermissionsForUserInDomain gets permissions for a user or role inside a domain.

func (*Enforcer)

GetPolicy

¶

func (e *

Enforcer

) GetPolicy() ([][]

string

,

error

)

GetPolicy gets all the authorization rules in the policy.

func (*Enforcer)

GetRoleManager

¶

added in

v2.1.0

func (e *

Enforcer

) GetRoleManager()

rbac

.

RoleManager

GetRoleManager gets the current role manager.

func (*Enforcer)

GetRolesForUser

¶

func (e *

Enforcer

) GetRolesForUser(name

string

, domain ...

string

) ([]

string

,

error

)

GetRolesForUser gets the roles that a user has.

func (*Enforcer)

GetRolesForUserInDomain

¶

func (e *

Enforcer

) GetRolesForUserInDomain(name

string

, domain

string

) []

string

GetRolesForUserInDomain gets the roles that a user has inside a domain.

func (*Enforcer)

GetUsersForRole

¶

func (e *

Enforcer

) GetUsersForRole(name

string

, domain ...

string

) ([]

string

,

error

)

GetUsersForRole gets the users that has a role.

func (*Enforcer)

GetUsersForRoleInDomain

¶

func (e *

Enforcer

) GetUsersForRoleInDomain(name

string

, domain

string

) []

string

GetUsersForRoleInDomain gets the users that has a role inside a domain. Add by Gordon.

func (*Enforcer)

HasGroupingPolicy

¶

func (e *

Enforcer

) HasGroupingPolicy(params ...interface{}) (

bool

,

error

)

HasGroupingPolicy determines whether a role inheritance rule exists.

func (*Enforcer)

HasNamedGroupingPolicy

¶

func (e *

Enforcer

) HasNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

HasNamedGroupingPolicy determines whether a named role inheritance rule exists.

func (*Enforcer)

HasNamedPolicy

¶

func (e *

Enforcer

) HasNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

HasNamedPolicy determines whether a named authorization rule exists.

func (*Enforcer)

HasPermissionForUser

¶

func (e *

Enforcer

) HasPermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

HasPermissionForUser determines whether a user has a permission.

func (*Enforcer)

HasPolicy

¶

func (e *

Enforcer

) HasPolicy(params ...interface{}) (

bool

,

error

)

HasPolicy determines whether an authorization rule exists.

func (*Enforcer)

HasRoleForUser

¶

func (e *

Enforcer

) HasRoleForUser(name

string

, role

string

, domain ...

string

) (

bool

,

error

)

HasRoleForUser determines whether a user has a role.

func (*Enforcer)

InitWithAdapter

¶

func (e *

Enforcer

) InitWithAdapter(modelPath

string

, adapter

persist

.

Adapter

)

error

InitWithAdapter initializes an enforcer with a database adapter.

func (*Enforcer)

InitWithFile

¶

func (e *

Enforcer

) InitWithFile(modelPath

string

, policyPath

string

)

error

InitWithFile initializes an enforcer with a model file and a policy file.

func (*Enforcer)

InitWithModelAndAdapter

¶

func (e *

Enforcer

) InitWithModelAndAdapter(m

model

.

Model

, adapter

persist

.

Adapter

)

error

InitWithModelAndAdapter initializes an enforcer with a model and a database adapter.

func (*Enforcer)

IsFiltered

¶

func (e *

Enforcer

) IsFiltered()

bool

IsFiltered returns true if the loaded policy has been filtered.

func (*Enforcer)

IsLogEnabled

¶

added in

v2.16.0

func (e *

Enforcer

) IsLogEnabled()

bool

IsLogEnabled returns the current logger's enabled status.

func (*Enforcer)

LoadFilteredPolicy

¶

func (e *

Enforcer

) LoadFilteredPolicy(filter interface{})

error

LoadFilteredPolicy reloads a filtered policy from file/database.

func (*Enforcer)

LoadFilteredPolicyCtx

¶

added in

v2.128.0

func (e *

Enforcer

) LoadFilteredPolicyCtx(ctx

context

.

Context

, filter interface{})

error

LoadFilteredPolicyCtx loads all policy rules from the storage with context and filter.

func (*Enforcer)

LoadIncrementalFilteredPolicy

¶

added in

v2.11.0

func (e *

Enforcer

) LoadIncrementalFilteredPolicy(filter interface{})

error

LoadIncrementalFilteredPolicy append a filtered policy from file/database.

func (*Enforcer)

LoadIncrementalFilteredPolicyCtx

¶

added in

v2.128.0

func (e *

Enforcer

) LoadIncrementalFilteredPolicyCtx(ctx

context

.

Context

, filter interface{})

error

LoadIncrementalFilteredPolicyCtx append a filtered policy from file/database with context.

func (*Enforcer)

LoadModel

¶

func (e *

Enforcer

) LoadModel()

error

LoadModel reloads the model from the model CONF file.
Because the policy is attached to a model, so the policy is invalidated and needs to be reloaded by calling LoadPolicy().

func (*Enforcer)

LoadPolicy

¶

func (e *

Enforcer

) LoadPolicy()

error

LoadPolicy reloads the policy from file/database.

func (*Enforcer)

RemoveFilteredGroupingPolicy

¶

func (e *

Enforcer

) RemoveFilteredGroupingPolicy(fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredGroupingPolicy removes a role inheritance rule from the current policy, field filters can be specified.

func (*Enforcer)

RemoveFilteredNamedGroupingPolicy

¶

func (e *

Enforcer

) RemoveFilteredNamedGroupingPolicy(ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredNamedGroupingPolicy removes a role inheritance rule from the current named policy, field filters can be specified.

func (*Enforcer)

RemoveFilteredNamedPolicy

¶

func (e *

Enforcer

) RemoveFilteredNamedPolicy(ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredNamedPolicy removes an authorization rule from the current named policy, field filters can be specified.

func (*Enforcer)

RemoveFilteredPolicy

¶

func (e *

Enforcer

) RemoveFilteredPolicy(fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredPolicy removes an authorization rule from the current policy, field filters can be specified.

func (*Enforcer)

RemoveGroupingPolicies

¶

added in

v2.2.2

func (e *

Enforcer

) RemoveGroupingPolicies(rules [][]

string

) (

bool

,

error

)

RemoveGroupingPolicies removes role inheritance rules from the current policy.

func (*Enforcer)

RemoveGroupingPolicy

¶

func (e *

Enforcer

) RemoveGroupingPolicy(params ...interface{}) (

bool

,

error

)

RemoveGroupingPolicy removes a role inheritance rule from the current policy.

func (*Enforcer)

RemoveNamedGroupingPolicies

¶

added in

v2.2.2

func (e *

Enforcer

) RemoveNamedGroupingPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveNamedGroupingPolicies removes role inheritance rules from the current named policy.

func (*Enforcer)

RemoveNamedGroupingPolicy

¶

func (e *

Enforcer

) RemoveNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedGroupingPolicy removes a role inheritance rule from the current named policy.

func (*Enforcer)

RemoveNamedPolicies

¶

added in

v2.2.2

func (e *

Enforcer

) RemoveNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveNamedPolicies removes authorization rules from the current named policy.

func (*Enforcer)

RemoveNamedPolicy

¶

func (e *

Enforcer

) RemoveNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedPolicy removes an authorization rule from the current named policy.

func (*Enforcer)

RemovePolicies

¶

added in

v2.2.2

func (e *

Enforcer

) RemovePolicies(rules [][]

string

) (

bool

,

error

)

RemovePolicies removes authorization rules from the current policy.

func (*Enforcer)

RemovePolicy

¶

func (e *

Enforcer

) RemovePolicy(params ...interface{}) (

bool

,

error

)

RemovePolicy removes an authorization rule from the current policy.

func (*Enforcer)

SavePolicy

¶

func (e *

Enforcer

) SavePolicy()

error

SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.

func (*Enforcer)

SelfAddPolicies

¶

added in

v2.53.0

func (e *

Enforcer

) SelfAddPolicies(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*Enforcer)

SelfAddPoliciesEx

¶

added in

v2.63.0

func (e *

Enforcer

) SelfAddPoliciesEx(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*Enforcer)

SelfAddPolicy

¶

added in

v2.53.0

func (e *

Enforcer

) SelfAddPolicy(sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

func (*Enforcer)

SelfRemoveFilteredPolicy

¶

added in

v2.53.0

func (e *

Enforcer

) SelfRemoveFilteredPolicy(sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

func (*Enforcer)

SelfRemovePolicies

¶

added in

v2.53.1

func (e *

Enforcer

) SelfRemovePolicies(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*Enforcer)

SelfRemovePolicy

¶

added in

v2.53.0

func (e *

Enforcer

) SelfRemovePolicy(sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

func (*Enforcer)

SelfUpdatePolicies

¶

added in

v2.53.0

func (e *

Enforcer

) SelfUpdatePolicies(sec

string

, ptype

string

, oldRules, newRules [][]

string

) (

bool

,

error

)

func (*Enforcer)

SelfUpdatePolicy

¶

added in

v2.53.0

func (e *

Enforcer

) SelfUpdatePolicy(sec

string

, ptype

string

, oldRule, newRule []

string

) (

bool

,

error

)

func (*Enforcer)

SetAdapter

¶

func (e *

Enforcer

) SetAdapter(adapter

persist

.

Adapter

)

SetAdapter sets the current adapter.

func (*Enforcer)

SetEffector

¶

func (e *

Enforcer

) SetEffector(eft

effector

.

Effector

)

SetEffector sets the current effector.

func (*Enforcer)

SetFieldIndex

¶

added in

v2.48.0

func (e *

Enforcer

) SetFieldIndex(ptype

string

, field

string

, index

int

)

func (*Enforcer)

SetLogger

¶

added in

v2.16.0

func (e *

Enforcer

) SetLogger(logger

log

.

Logger

)

SetLogger changes the current enforcer's logger.

func (*Enforcer)

SetModel

¶

func (e *

Enforcer

) SetModel(m

model

.

Model

)

SetModel sets the current model.

func (*Enforcer)

SetNamedDomainLinkConditionFuncParams

¶

added in

v2.75.0

func (e *

Enforcer

) SetNamedDomainLinkConditionFuncParams(ptype, user, role, domain

string

, params ...

string

)

bool

SetNamedDomainLinkConditionFuncParams Sets the parameters of the condition function fn
for Link userName->{roleName, domain}.

func (*Enforcer)

SetNamedLinkConditionFuncParams

¶

added in

v2.75.0

func (e *

Enforcer

) SetNamedLinkConditionFuncParams(ptype, user, role

string

, params ...

string

)

bool

SetNamedLinkConditionFuncParams Sets the parameters of the condition function fn for Link userName->roleName.

func (*Enforcer)

SetNamedRoleManager

¶

added in

v2.52.0

func (e *

Enforcer

) SetNamedRoleManager(ptype

string

, rm

rbac

.

RoleManager

)

SetNamedRoleManager sets the role manager for the named policy.

func (*Enforcer)

SetRoleManager

¶

func (e *

Enforcer

) SetRoleManager(rm

rbac

.

RoleManager

)

SetRoleManager sets the current role manager.

func (*Enforcer)

SetWatcher

¶

func (e *

Enforcer

) SetWatcher(watcher

persist

.

Watcher

)

error

SetWatcher sets the current watcher.

func (*Enforcer)

UpdateFilteredNamedPolicies

¶

added in

v2.28.0

func (e *

Enforcer

) UpdateFilteredNamedPolicies(ptype

string

, newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

func (*Enforcer)

UpdateFilteredPolicies

¶

added in

v2.28.0

func (e *

Enforcer

) UpdateFilteredPolicies(newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

func (*Enforcer)

UpdateGroupingPolicies

¶

added in

v2.41.0

func (e *

Enforcer

) UpdateGroupingPolicies(oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

UpdateGroupingPolicies updates authorization rules from the current policies.

func (*Enforcer)

UpdateGroupingPolicy

¶

added in

v2.19.0

func (e *

Enforcer

) UpdateGroupingPolicy(oldRule []

string

, newRule []

string

) (

bool

,

error

)

func (*Enforcer)

UpdateNamedGroupingPolicies

¶

added in

v2.41.0

func (e *

Enforcer

) UpdateNamedGroupingPolicies(ptype

string

, oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

func (*Enforcer)

UpdateNamedGroupingPolicy

¶

added in

v2.19.0

func (e *

Enforcer

) UpdateNamedGroupingPolicy(ptype

string

, oldRule []

string

, newRule []

string

) (

bool

,

error

)

func (*Enforcer)

UpdateNamedPolicies

¶

added in

v2.22.0

func (e *

Enforcer

) UpdateNamedPolicies(ptype

string

, p1 [][]

string

, p2 [][]

string

) (

bool

,

error

)

func (*Enforcer)

UpdateNamedPolicy

¶

added in

v2.14.0

func (e *

Enforcer

) UpdateNamedPolicy(ptype

string

, p1 []

string

, p2 []

string

) (

bool

,

error

)

func (*Enforcer)

UpdatePolicies

¶

added in

v2.22.0

func (e *

Enforcer

) UpdatePolicies(oldPolices [][]

string

, newPolicies [][]

string

) (

bool

,

error

)

UpdatePolicies updates authorization rules from the current policies.

func (*Enforcer)

UpdatePolicy

¶

added in

v2.14.0

func (e *

Enforcer

) UpdatePolicy(oldPolicy []

string

, newPolicy []

string

) (

bool

,

error

)

UpdatePolicy updates an authorization rule from the current policy.

type

IDistributedEnforcer

¶

added in

v2.19.0

type IDistributedEnforcer interface {

IEnforcer

SetDispatcher(dispatcher

persist

.

Dispatcher

)

/* Management API for DistributedEnforcer*/

AddPoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, rules [][]

string

) (affected [][]

string

, err

error

)

RemovePoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, rules [][]

string

) (affected [][]

string

, err

error

)

RemoveFilteredPolicySelf(shouldPersist func()

bool

, sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (affected [][]

string

, err

error

)

ClearPolicySelf(shouldPersist func()

bool

)

error

UpdatePolicySelf(shouldPersist func()

bool

, sec

string

, ptype

string

, oldRule, newRule []

string

) (affected

bool

, err

error

)

UpdatePoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, oldRules, newRules [][]

string

) (affected

bool

, err

error

)

UpdateFilteredPoliciesSelf(shouldPersist func()

bool

, sec

string

, ptype

string

, newRules [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

}

IDistributedEnforcer defines dispatcher enforcer.

type

IEnforcer

¶

added in

v2.1.1

type IEnforcer interface {

/* Enforcer API */

InitWithFile(modelPath

string

, policyPath

string

)

error

InitWithAdapter(modelPath

string

, adapter

persist

.

Adapter

)

error

InitWithModelAndAdapter(m

model

.

Model

, adapter

persist

.

Adapter

)

error

LoadModel()

error

GetModel()

model

.

Model

SetModel(m

model

.

Model

)

GetAdapter()

persist

.

Adapter

SetAdapter(adapter

persist

.

Adapter

)

SetWatcher(watcher

persist

.

Watcher

)

error

GetRoleManager()

rbac

.

RoleManager

SetRoleManager(rm

rbac

.

RoleManager

)

SetEffector(eft

effector

.

Effector

)

ClearPolicy()

LoadPolicy()

error

LoadFilteredPolicy(filter interface{})

error

LoadIncrementalFilteredPolicy(filter interface{})

error

IsFiltered()

bool

SavePolicy()

error

EnableEnforce(enable

bool

)

EnableLog(enable

bool

)

EnableAutoNotifyWatcher(enable

bool

)

EnableAutoSave(autoSave

bool

)

EnableAutoBuildRoleLinks(autoBuildRoleLinks

bool

)

BuildRoleLinks()

error

Enforce(rvals ...interface{}) (

bool

,

error

)

EnforceWithMatcher(matcher

string

, rvals ...interface{}) (

bool

,

error

)

EnforceEx(rvals ...interface{}) (

bool

, []

string

,

error

)

EnforceExWithMatcher(matcher

string

, rvals ...interface{}) (

bool

, []

string

,

error

)

BatchEnforce(requests [][]interface{}) ([]

bool

,

error

)

BatchEnforceWithMatcher(matcher

string

, requests [][]interface{}) ([]

bool

,

error

)

/* RBAC API */

GetRolesForUser(name

string

, domain ...

string

) ([]

string

,

error

)

GetUsersForRole(name

string

, domain ...

string

) ([]

string

,

error

)

HasRoleForUser(name

string

, role

string

, domain ...

string

) (

bool

,

error

)

AddRoleForUser(user

string

, role

string

, domain ...

string

) (

bool

,

error

)

AddPermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

AddPermissionsForUser(user

string

, permissions ...[]

string

) (

bool

,

error

)

DeletePermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

DeletePermissionsForUser(user

string

) (

bool

,

error

)

GetPermissionsForUser(user

string

, domain ...

string

) ([][]

string

,

error

)

HasPermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

GetImplicitRolesForUser(name

string

, domain ...

string

) ([]

string

,

error

)

GetImplicitPermissionsForUser(user

string

, domain ...

string

) ([][]

string

,

error

)

GetImplicitUsersForPermission(permission ...

string

) ([]

string

,

error

)

DeleteRoleForUser(user

string

, role

string

, domain ...

string

) (

bool

,

error

)

DeleteRolesForUser(user

string

, domain ...

string

) (

bool

,

error

)

DeleteUser(user

string

) (

bool

,

error

)

DeleteRole(role

string

) (

bool

,

error

)

DeletePermission(permission ...

string

) (

bool

,

error

)

/* RBAC API with domains*/

GetUsersForRoleInDomain(name

string

, domain

string

) []

string

GetRolesForUserInDomain(name

string

, domain

string

) []

string

GetPermissionsForUserInDomain(user

string

, domain

string

) [][]

string

AddRoleForUserInDomain(user

string

, role

string

, domain

string

) (

bool

,

error

)

DeleteRoleForUserInDomain(user

string

, role

string

, domain

string

) (

bool

,

error

)

GetAllUsersByDomain(domain

string

) ([]

string

,

error

)

DeleteRolesForUserInDomain(user

string

, domain

string

) (

bool

,

error

)

DeleteAllUsersByDomain(domain

string

) (

bool

,

error

)

DeleteDomains(domains ...

string

) (

bool

,

error

)

GetAllDomains() ([]

string

,

error

)

GetAllRolesByDomain(domain

string

) ([]

string

,

error

)

/* Management API */

GetAllSubjects() ([]

string

,

error

)

GetAllNamedSubjects(ptype

string

) ([]

string

,

error

)

GetAllObjects() ([]

string

,

error

)

GetAllNamedObjects(ptype

string

) ([]

string

,

error

)

GetAllActions() ([]

string

,

error

)

GetAllNamedActions(ptype

string

) ([]

string

,

error

)

GetAllRoles() ([]

string

,

error

)

GetAllNamedRoles(ptype

string

) ([]

string

,

error

)

GetPolicy() ([][]

string

,

error

)

GetFilteredPolicy(fieldIndex

int

, fieldValues ...

string

) ([][]

string

,

error

)

GetNamedPolicy(ptype

string

) ([][]

string

,

error

)

GetFilteredNamedPolicy(ptype

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

GetGroupingPolicy() ([][]

string

,

error

)

GetFilteredGroupingPolicy(fieldIndex

int

, fieldValues ...

string

) ([][]

string

,

error

)

GetNamedGroupingPolicy(ptype

string

) ([][]

string

,

error

)

GetFilteredNamedGroupingPolicy(ptype

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

HasPolicy(params ...interface{}) (

bool

,

error

)

HasNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddPolicy(params ...interface{}) (

bool

,

error

)

AddPolicies(rules [][]

string

) (

bool

,

error

)

AddNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddPoliciesEx(rules [][]

string

) (

bool

,

error

)

AddNamedPoliciesEx(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemovePolicy(params ...interface{}) (

bool

,

error

)

RemovePolicies(rules [][]

string

) (

bool

,

error

)

RemoveFilteredPolicy(fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveFilteredNamedPolicy(ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

HasGroupingPolicy(params ...interface{}) (

bool

,

error

)

HasNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddGroupingPolicy(params ...interface{}) (

bool

,

error

)

AddGroupingPolicies(rules [][]

string

) (

bool

,

error

)

AddGroupingPoliciesEx(rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedGroupingPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPoliciesEx(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveGroupingPolicy(params ...interface{}) (

bool

,

error

)

RemoveGroupingPolicies(rules [][]

string

) (

bool

,

error

)

RemoveFilteredGroupingPolicy(fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedGroupingPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveFilteredNamedGroupingPolicy(ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

AddFunction(name

string

, function

govaluate

.

ExpressionFunction

)

UpdatePolicy(oldPolicy []

string

, newPolicy []

string

) (

bool

,

error

)

UpdatePolicies(oldPolicies [][]

string

, newPolicies [][]

string

) (

bool

,

error

)

UpdateFilteredPolicies(newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

UpdateGroupingPolicy(oldRule []

string

, newRule []

string

) (

bool

,

error

)

UpdateGroupingPolicies(oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

UpdateNamedGroupingPolicy(ptype

string

, oldRule []

string

, newRule []

string

) (

bool

,

error

)

UpdateNamedGroupingPolicies(ptype

string

, oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

/* Management API with autoNotifyWatcher disabled */

SelfAddPolicy(sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

SelfAddPolicies(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfAddPoliciesEx(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfRemovePolicy(sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

SelfRemovePolicies(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfRemoveFilteredPolicy(sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

SelfUpdatePolicy(sec

string

, ptype

string

, oldRule, newRule []

string

) (

bool

,

error

)

SelfUpdatePolicies(sec

string

, ptype

string

, oldRules, newRules [][]

string

) (

bool

,

error

)

}

IEnforcer is the API interface of Enforcer.

type

IEnforcerContext

¶

added in

v2.128.0

type IEnforcerContext interface {

IEnforcer

/* Enforcer API */

LoadPolicyCtx(ctx

context

.

Context

)

error

LoadFilteredPolicyCtx(ctx

context

.

Context

, filter interface{})

error

LoadIncrementalFilteredPolicyCtx(ctx

context

.

Context

, filter interface{})

error

IsFilteredCtx(ctx

context

.

Context

)

bool

SavePolicyCtx(ctx

context

.

Context

)

error

/* RBAC API */

AddRoleForUserCtx(ctx

context

.

Context

, user

string

, role

string

, domain ...

string

) (

bool

,

error

)

AddPermissionForUserCtx(ctx

context

.

Context

, user

string

, permission ...

string

) (

bool

,

error

)

AddPermissionsForUserCtx(ctx

context

.

Context

, user

string

, permissions ...[]

string

) (

bool

,

error

)

DeletePermissionForUserCtx(ctx

context

.

Context

, user

string

, permission ...

string

) (

bool

,

error

)

DeletePermissionsForUserCtx(ctx

context

.

Context

, user

string

) (

bool

,

error

)

DeleteRoleForUserCtx(ctx

context

.

Context

, user

string

, role

string

, domain ...

string

) (

bool

,

error

)

DeleteRolesForUserCtx(ctx

context

.

Context

, user

string

, domain ...

string

) (

bool

,

error

)

DeleteUserCtx(ctx

context

.

Context

, user

string

) (

bool

,

error

)

DeleteRoleCtx(ctx

context

.

Context

, role

string

) (

bool

,

error

)

DeletePermissionCtx(ctx

context

.

Context

, permission ...

string

) (

bool

,

error

)

/* RBAC API with domains*/

AddRoleForUserInDomainCtx(ctx

context

.

Context

, user

string

, role

string

, domain

string

) (

bool

,

error

)

DeleteRoleForUserInDomainCtx(ctx

context

.

Context

, user

string

, role

string

, domain

string

) (

bool

,

error

)

DeleteRolesForUserInDomainCtx(ctx

context

.

Context

, user

string

, domain

string

) (

bool

,

error

)

DeleteAllUsersByDomainCtx(ctx

context

.

Context

, domain

string

) (

bool

,

error

)

DeleteDomainsCtx(ctx

context

.

Context

, domains ...

string

) (

bool

,

error

)

/* Management API */

AddPolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

AddPoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

AddNamedPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

AddPoliciesExCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

AddNamedPoliciesExCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

RemovePolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

RemovePoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

RemoveFilteredPolicyCtx(ctx

context

.

Context

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveNamedPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveFilteredNamedPolicyCtx(ctx

context

.

Context

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

AddGroupingPolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

AddGroupingPoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

AddGroupingPoliciesExCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedGroupingPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPoliciesExCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveGroupingPolicyCtx(ctx

context

.

Context

, params ...interface{}) (

bool

,

error

)

RemoveGroupingPoliciesCtx(ctx

context

.

Context

, rules [][]

string

) (

bool

,

error

)

RemoveFilteredGroupingPolicyCtx(ctx

context

.

Context

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedGroupingPoliciesCtx(ctx

context

.

Context

, ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveFilteredNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

UpdatePolicyCtx(ctx

context

.

Context

, oldPolicy []

string

, newPolicy []

string

) (

bool

,

error

)

UpdatePoliciesCtx(ctx

context

.

Context

, oldPolicies [][]

string

, newPolicies [][]

string

) (

bool

,

error

)

UpdateFilteredPoliciesCtx(ctx

context

.

Context

, newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

UpdateGroupingPolicyCtx(ctx

context

.

Context

, oldRule []

string

, newRule []

string

) (

bool

,

error

)

UpdateGroupingPoliciesCtx(ctx

context

.

Context

, oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

UpdateNamedGroupingPolicyCtx(ctx

context

.

Context

, ptype

string

, oldRule []

string

, newRule []

string

) (

bool

,

error

)

UpdateNamedGroupingPoliciesCtx(ctx

context

.

Context

, ptype

string

, oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

/* Management API with autoNotifyWatcher disabled */

SelfAddPolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

SelfAddPoliciesCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfAddPoliciesExCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfRemovePolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

SelfRemovePoliciesCtx(ctx

context

.

Context

, sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

SelfRemoveFilteredPolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

SelfUpdatePolicyCtx(ctx

context

.

Context

, sec

string

, ptype

string

, oldRule, newRule []

string

) (

bool

,

error

)

SelfUpdatePoliciesCtx(ctx

context

.

Context

, sec

string

, ptype

string

, oldRules, newRules [][]

string

) (

bool

,

error

)

}

func

NewContextEnforcer

¶

added in

v2.128.0

func NewContextEnforcer(params ...interface{}) (

IEnforcerContext

,

error

)

NewContextEnforcer creates a context-aware enforcer via file or DB.

type

SyncedCachedEnforcer

¶

added in

v2.66.0

type SyncedCachedEnforcer struct {

*

SyncedEnforcer

// contains filtered or unexported fields

}

SyncedCachedEnforcer wraps Enforcer and provides decision sync cache.

func

NewSyncedCachedEnforcer

¶

added in

v2.66.0

func NewSyncedCachedEnforcer(params ...interface{}) (*

SyncedCachedEnforcer

,

error

)

NewSyncedCachedEnforcer creates a sync cached enforcer via file or DB.

func (*SyncedCachedEnforcer)

AddPolicies

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) AddPolicies(rules [][]

string

) (

bool

,

error

)

func (*SyncedCachedEnforcer)

AddPolicy

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) AddPolicy(params ...interface{}) (

bool

,

error

)

func (*SyncedCachedEnforcer)

EnableCache

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) EnableCache(enableCache

bool

)

EnableCache determines whether to enable cache on Enforce(). When enableCache is enabled, cached result (true | false) will be returned for previous decisions.

func (*SyncedCachedEnforcer)

Enforce

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) Enforce(rvals ...interface{}) (

bool

,

error

)

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).
if rvals is not string , ignore the cache.

func (*SyncedCachedEnforcer)

InvalidateCache

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) InvalidateCache()

error

InvalidateCache deletes all the existing cached decisions.

func (*SyncedCachedEnforcer)

LoadPolicy

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) LoadPolicy()

error

func (*SyncedCachedEnforcer)

RemovePolicies

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) RemovePolicies(rules [][]

string

) (

bool

,

error

)

func (*SyncedCachedEnforcer)

RemovePolicy

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) RemovePolicy(params ...interface{}) (

bool

,

error

)

func (*SyncedCachedEnforcer)

SetCache

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) SetCache(c

cache

.

Cache

)

SetCache need to be sync cache.

func (*SyncedCachedEnforcer)

SetExpireTime

¶

added in

v2.66.0

func (e *

SyncedCachedEnforcer

) SetExpireTime(expireTime

time

.

Duration

)

type

SyncedEnforcer

¶

type SyncedEnforcer struct {

*

Enforcer

// contains filtered or unexported fields

}

SyncedEnforcer wraps Enforcer and provides synchronized access.

func

NewSyncedEnforcer

¶

func NewSyncedEnforcer(params ...interface{}) (*

SyncedEnforcer

,

error

)

NewSyncedEnforcer creates a synchronized enforcer via file or DB.

func (*SyncedEnforcer)

AddFunction

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) AddFunction(name

string

, function

govaluate

.

ExpressionFunction

)

AddFunction adds a customized function.

func (*SyncedEnforcer)

AddGroupingPolicies

¶

added in

v2.8.1

func (e *

SyncedEnforcer

) AddGroupingPolicies(rules [][]

string

) (

bool

,

error

)

AddGroupingPolicies adds role inheritance rulea to the current policy.
If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added.
Otherwise the function returns true for the corresponding policy rule by adding the new rule.

func (*SyncedEnforcer)

AddGroupingPoliciesEx

¶

added in

v2.63.0

func (e *

SyncedEnforcer

) AddGroupingPoliciesEx(rules [][]

string

) (

bool

,

error

)

AddGroupingPoliciesEx adds role inheritance rules to the current policy.
If the rule already exists, the rule will not be added.
But unlike AddGroupingPolicies, other non-existent rules are added instead of returning false directly.

func (*SyncedEnforcer)

AddGroupingPolicy

¶

func (e *

SyncedEnforcer

) AddGroupingPolicy(params ...interface{}) (

bool

,

error

)

AddGroupingPolicy adds a role inheritance rule to the current policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*SyncedEnforcer)

AddNamedGroupingPolicies

¶

added in

v2.8.1

func (e *

SyncedEnforcer

) AddNamedGroupingPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPolicies adds named role inheritance rules to the current policy.
If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added.
Otherwise the function returns true for the corresponding policy rule by adding the new rule.

func (*SyncedEnforcer)

AddNamedGroupingPoliciesEx

¶

added in

v2.63.0

func (e *

SyncedEnforcer

) AddNamedGroupingPoliciesEx(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedGroupingPoliciesEx adds named role inheritance rules to the current policy.
If the rule already exists, the rule will not be added.
But unlike AddNamedGroupingPolicies, other non-existent rules are added instead of returning false directly.

func (*SyncedEnforcer)

AddNamedGroupingPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) AddNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedGroupingPolicy adds a named role inheritance rule to the current policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*SyncedEnforcer)

AddNamedPolicies

¶

added in

v2.8.1

func (e *

SyncedEnforcer

) AddNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedPolicies adds authorization rules to the current named policy.
If the rule already exists, the function returns false for the corresponding rule and the rule will not be added.
Otherwise the function returns true for the corresponding by adding the new rule.

func (*SyncedEnforcer)

AddNamedPoliciesEx

¶

added in

v2.63.0

func (e *

SyncedEnforcer

) AddNamedPoliciesEx(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedPoliciesEx adds authorization rules to the current named policy.
If the rule already exists, the rule will not be added.
But unlike AddNamedPolicies, other non-existent rules are added instead of returning false directly.

func (*SyncedEnforcer)

AddNamedPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) AddNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedPolicy adds an authorization rule to the current named policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*SyncedEnforcer)

AddPermissionForUser

¶

func (e *

SyncedEnforcer

) AddPermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

AddPermissionForUser adds a permission for a user or role.
Returns false if the user or role already has the permission (aka not affected).

func (*SyncedEnforcer)

AddPermissionsForUser

¶

added in

v2.73.0

func (e *

SyncedEnforcer

) AddPermissionsForUser(user

string

, permissions ...[]

string

) (

bool

,

error

)

AddPermissionsForUser adds permissions for a user or role.
Returns false if the user or role already has the permissions (aka not affected).

func (*SyncedEnforcer)

AddPolicies

¶

added in

v2.8.1

func (e *

SyncedEnforcer

) AddPolicies(rules [][]

string

) (

bool

,

error

)

AddPolicies adds authorization rules to the current policy.
If the rule already exists, the function returns false for the corresponding rule and the rule will not be added.
Otherwise the function returns true for the corresponding rule by adding the new rule.

func (*SyncedEnforcer)

AddPoliciesEx

¶

added in

v2.63.0

func (e *

SyncedEnforcer

) AddPoliciesEx(rules [][]

string

) (

bool

,

error

)

AddPoliciesEx adds authorization rules to the current policy.
If the rule already exists, the rule will not be added.
But unlike AddPolicies, other non-existent rules are added instead of returning false directly.

func (*SyncedEnforcer)

AddPolicy

¶

func (e *

SyncedEnforcer

) AddPolicy(params ...interface{}) (

bool

,

error

)

AddPolicy adds an authorization rule to the current policy.
If the rule already exists, the function returns false and the rule will not be added.
Otherwise the function returns true by adding the new rule.

func (*SyncedEnforcer)

AddRoleForUser

¶

func (e *

SyncedEnforcer

) AddRoleForUser(user

string

, role

string

, domain ...

string

) (

bool

,

error

)

AddRoleForUser adds a role for a user.
Returns false if the user already has the role (aka not affected).

func (*SyncedEnforcer)

AddRoleForUserInDomain

¶

func (e *

SyncedEnforcer

) AddRoleForUserInDomain(user

string

, role

string

, domain

string

) (

bool

,

error

)

AddRoleForUserInDomain adds a role for a user inside a domain.
Returns false if the user already has the role (aka not affected).

func (*SyncedEnforcer)

AddRolesForUser

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) AddRolesForUser(user

string

, roles []

string

, domain ...

string

) (

bool

,

error

)

AddRolesForUser adds roles for a user.
Returns false if the user already has the roles (aka not affected).

func (*SyncedEnforcer)

BatchEnforce

¶

added in

v2.25.0

func (e *

SyncedEnforcer

) BatchEnforce(requests [][]interface{}) ([]

bool

,

error

)

BatchEnforce enforce in batches.

func (*SyncedEnforcer)

BatchEnforceWithMatcher

¶

added in

v2.25.0

func (e *

SyncedEnforcer

) BatchEnforceWithMatcher(matcher

string

, requests [][]interface{}) ([]

bool

,

error

)

BatchEnforceWithMatcher enforce with matcher in batches.

func (*SyncedEnforcer)

BuildRoleLinks

¶

func (e *

SyncedEnforcer

) BuildRoleLinks()

error

BuildRoleLinks manually rebuild the role inheritance relations.

func (*SyncedEnforcer)

ClearPolicy

¶

func (e *

SyncedEnforcer

) ClearPolicy()

ClearPolicy clears all policy.

func (*SyncedEnforcer)

DeleteDomains

¶

added in

v2.109.0

func (e *

SyncedEnforcer

) DeleteDomains(domains ...

string

) (

bool

,

error

)

DeleteDomains deletes domains from the model.
Returns false if the domain does not exist (aka not affected).

func (*SyncedEnforcer)

DeletePermission

¶

func (e *

SyncedEnforcer

) DeletePermission(permission ...

string

) (

bool

,

error

)

DeletePermission deletes a permission.
Returns false if the permission does not exist (aka not affected).

func (*SyncedEnforcer)

DeletePermissionForUser

¶

func (e *

SyncedEnforcer

) DeletePermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

DeletePermissionForUser deletes a permission for a user or role.
Returns false if the user or role does not have the permission (aka not affected).

func (*SyncedEnforcer)

DeletePermissionsForUser

¶

func (e *

SyncedEnforcer

) DeletePermissionsForUser(user

string

) (

bool

,

error

)

DeletePermissionsForUser deletes permissions for a user or role.
Returns false if the user or role does not have any permissions (aka not affected).

func (*SyncedEnforcer)

DeleteRole

¶

func (e *

SyncedEnforcer

) DeleteRole(role

string

) (

bool

,

error

)

DeleteRole deletes a role.
Returns false if the role does not exist (aka not affected).

func (*SyncedEnforcer)

DeleteRoleForUser

¶

func (e *

SyncedEnforcer

) DeleteRoleForUser(user

string

, role

string

, domain ...

string

) (

bool

,

error

)

DeleteRoleForUser deletes a role for a user.
Returns false if the user does not have the role (aka not affected).

func (*SyncedEnforcer)

DeleteRoleForUserInDomain

¶

func (e *

SyncedEnforcer

) DeleteRoleForUserInDomain(user

string

, role

string

, domain

string

) (

bool

,

error

)

DeleteRoleForUserInDomain deletes a role for a user inside a domain.
Returns false if the user does not have the role (aka not affected).

func (*SyncedEnforcer)

DeleteRolesForUser

¶

func (e *

SyncedEnforcer

) DeleteRolesForUser(user

string

, domain ...

string

) (

bool

,

error

)

DeleteRolesForUser deletes all roles for a user.
Returns false if the user does not have any roles (aka not affected).

func (*SyncedEnforcer)

DeleteRolesForUserInDomain

¶

added in

v2.8.4

func (e *

SyncedEnforcer

) DeleteRolesForUserInDomain(user

string

, domain

string

) (

bool

,

error

)

DeleteRolesForUserInDomain deletes all roles for a user inside a domain.
Returns false if the user does not have any roles (aka not affected).

func (*SyncedEnforcer)

DeleteUser

¶

func (e *

SyncedEnforcer

) DeleteUser(user

string

) (

bool

,

error

)

DeleteUser deletes a user.
Returns false if the user does not exist (aka not affected).

func (*SyncedEnforcer)

Enforce

¶

func (e *

SyncedEnforcer

) Enforce(rvals ...interface{}) (

bool

,

error

)

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).

func (*SyncedEnforcer)

EnforceEx

¶

added in

v2.29.1

func (e *

SyncedEnforcer

) EnforceEx(rvals ...interface{}) (

bool

, []

string

,

error

)

EnforceEx explain enforcement by informing matched rules.

func (*SyncedEnforcer)

EnforceExWithMatcher

¶

added in

v2.29.1

func (e *

SyncedEnforcer

) EnforceExWithMatcher(matcher

string

, rvals ...interface{}) (

bool

, []

string

,

error

)

EnforceExWithMatcher use a custom matcher and explain enforcement by informing matched rules.

func (*SyncedEnforcer)

EnforceWithMatcher

¶

added in

v2.29.1

func (e *

SyncedEnforcer

) EnforceWithMatcher(matcher

string

, rvals ...interface{}) (

bool

,

error

)

EnforceWithMatcher use a custom matcher to decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (matcher, sub, obj, act), use model matcher by default when matcher is "".

func (*SyncedEnforcer)

GetAllActions

¶

func (e *

SyncedEnforcer

) GetAllActions() ([]

string

,

error

)

GetAllActions gets the list of actions that show up in the current policy.

func (*SyncedEnforcer)

GetAllNamedActions

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetAllNamedActions(ptype

string

) ([]

string

,

error

)

GetAllNamedActions gets the list of actions that show up in the current named policy.

func (*SyncedEnforcer)

GetAllNamedObjects

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetAllNamedObjects(ptype

string

) ([]

string

,

error

)

GetAllNamedObjects gets the list of objects that show up in the current named policy.

func (*SyncedEnforcer)

GetAllNamedRoles

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetAllNamedRoles(ptype

string

) ([]

string

,

error

)

GetAllNamedRoles gets the list of roles that show up in the current named policy.

func (*SyncedEnforcer)

GetAllNamedSubjects

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetAllNamedSubjects(ptype

string

) ([]

string

,

error

)

GetAllNamedSubjects gets the list of subjects that show up in the current named policy.

func (*SyncedEnforcer)

GetAllObjects

¶

func (e *

SyncedEnforcer

) GetAllObjects() ([]

string

,

error

)

GetAllObjects gets the list of objects that show up in the current policy.

func (*SyncedEnforcer)

GetAllRoles

¶

func (e *

SyncedEnforcer

) GetAllRoles() ([]

string

,

error

)

GetAllRoles gets the list of roles that show up in the current policy.

func (*SyncedEnforcer)

GetAllSubjects

¶

func (e *

SyncedEnforcer

) GetAllSubjects() ([]

string

,

error

)

GetAllSubjects gets the list of subjects that show up in the current policy.

func (*SyncedEnforcer)

GetFilteredGroupingPolicy

¶

func (e *

SyncedEnforcer

) GetFilteredGroupingPolicy(fieldIndex

int

, fieldValues ...

string

) ([][]

string

,

error

)

GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

func (*SyncedEnforcer)

GetFilteredNamedGroupingPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetFilteredNamedGroupingPolicy(ptype

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

GetFilteredNamedGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

func (*SyncedEnforcer)

GetFilteredNamedPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetFilteredNamedPolicy(ptype

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

GetFilteredNamedPolicy gets all the authorization rules in the named policy, field filters can be specified.

func (*SyncedEnforcer)

GetFilteredPolicy

¶

func (e *

SyncedEnforcer

) GetFilteredPolicy(fieldIndex

int

, fieldValues ...

string

) ([][]

string

,

error

)

GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.

func (*SyncedEnforcer)

GetGroupingPolicy

¶

func (e *

SyncedEnforcer

) GetGroupingPolicy() ([][]

string

,

error

)

GetGroupingPolicy gets all the role inheritance rules in the policy.

func (*SyncedEnforcer)

GetImplicitObjectPatternsForUser

¶

added in

v2.121.0

func (e *

SyncedEnforcer

) GetImplicitObjectPatternsForUser(user

string

, domain

string

, action

string

) ([]

string

,

error

)

GetImplicitObjectPatternsForUser returns all object patterns (with wildcards) that a user has for a given domain and action.
For example:
p, admin, chronicle/123, location/*, read
p, user, chronicle/456, location/789, read
g, alice, admin
g, bob, user

GetImplicitObjectPatternsForUser("alice", "chronicle/123", "read") will return ["location/*"].
GetImplicitObjectPatternsForUser("bob", "chronicle/456", "read") will return ["location/789"].

func (*SyncedEnforcer)

GetImplicitPermissionsForUser

¶

added in

v2.13.0

func (e *

SyncedEnforcer

) GetImplicitPermissionsForUser(user

string

, domain ...

string

) ([][]

string

,

error

)

GetImplicitPermissionsForUser gets implicit permissions for a user or role.
Compared to GetPermissionsForUser(), this function retrieves permissions for inherited roles.
For example:
p, admin, data1, read
p, alice, data2, read
g, alice, admin

GetPermissionsForUser("alice") can only get: [["alice", "data2", "read"]].
But GetImplicitPermissionsForUser("alice") will get: [["admin", "data1", "read"], ["alice", "data2", "read"]].

func (*SyncedEnforcer)

GetImplicitRolesForUser

¶

added in

v2.13.0

func (e *

SyncedEnforcer

) GetImplicitRolesForUser(name

string

, domain ...

string

) ([]

string

,

error

)

GetImplicitRolesForUser gets implicit roles that a user has.
Compared to GetRolesForUser(), this function retrieves indirect roles besides direct roles.
For example:
g, alice, role:admin
g, role:admin, role:user

GetRolesForUser("alice") can only get: ["role:admin"].
But GetImplicitRolesForUser("alice") will get: ["role:admin", "role:user"].

func (*SyncedEnforcer)

GetImplicitUsersForPermission

¶

added in

v2.13.0

func (e *

SyncedEnforcer

) GetImplicitUsersForPermission(permission ...

string

) ([]

string

,

error

)

GetImplicitUsersForPermission gets implicit users for a permission.
For example:
p, admin, data1, read
p, bob, data1, read
g, alice, admin

GetImplicitUsersForPermission("data1", "read") will get: ["alice", "bob"].
Note: only users will be returned, roles (2nd arg in "g") will be excluded.

func (*SyncedEnforcer)

GetLock

¶

added in

v2.52.1

func (e *

SyncedEnforcer

) GetLock() *

sync

.

RWMutex

GetLock return the private RWMutex lock.

func (*SyncedEnforcer)

GetNamedGroupingPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetNamedGroupingPolicy(ptype

string

) ([][]

string

,

error

)

GetNamedGroupingPolicy gets all the role inheritance rules in the policy.

func (*SyncedEnforcer)

GetNamedImplicitPermissionsForUser

¶

added in

v2.45.0

func (e *

SyncedEnforcer

) GetNamedImplicitPermissionsForUser(ptype

string

, gtype

string

, user

string

, domain ...

string

) ([][]

string

,

error

)

GetNamedImplicitPermissionsForUser gets implicit permissions for a user or role by named policy.
Compared to GetNamedPermissionsForUser(), this function retrieves permissions for inherited roles.
For example:
p, admin, data1, read
p2, admin, create
g, alice, admin

GetImplicitPermissionsForUser("alice") can only get: [["admin", "data1", "read"]], whose policy is default policy "p"
But you can specify the named policy "p2" to get: [["admin", "create"]] by    GetNamedImplicitPermissionsForUser("p2","alice").

func (*SyncedEnforcer)

GetNamedPermissionsForUser

¶

added in

v2.45.0

func (e *

SyncedEnforcer

) GetNamedPermissionsForUser(ptype

string

, user

string

, domain ...

string

) ([][]

string

,

error

)

GetNamedPermissionsForUser gets permissions for a user or role by named policy.

func (*SyncedEnforcer)

GetNamedPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) GetNamedPolicy(ptype

string

) ([][]

string

,

error

)

GetNamedPolicy gets all the authorization rules in the named policy.

func (*SyncedEnforcer)

GetNamedRoleManager

¶

added in

v2.125.0

func (e *

SyncedEnforcer

) GetNamedRoleManager(ptype

string

)

rbac

.

RoleManager

GetNamedRoleManager gets the role manager for the named policy with synchronization.

func (*SyncedEnforcer)

GetPermissionsForUser

¶

func (e *

SyncedEnforcer

) GetPermissionsForUser(user

string

, domain ...

string

) ([][]

string

,

error

)

GetPermissionsForUser gets permissions for a user or role.

func (*SyncedEnforcer)

GetPermissionsForUserInDomain

¶

func (e *

SyncedEnforcer

) GetPermissionsForUserInDomain(user

string

, domain

string

) [][]

string

GetPermissionsForUserInDomain gets permissions for a user or role inside a domain.

func (*SyncedEnforcer)

GetPolicy

¶

func (e *

SyncedEnforcer

) GetPolicy() ([][]

string

,

error

)

GetPolicy gets all the authorization rules in the policy.

func (*SyncedEnforcer)

GetRoleManager

¶

added in

v2.125.0

func (e *

SyncedEnforcer

) GetRoleManager()

rbac

.

RoleManager

GetRoleManager gets the current role manager with synchronization.

func (*SyncedEnforcer)

GetRolesForUser

¶

func (e *

SyncedEnforcer

) GetRolesForUser(name

string

, domain ...

string

) ([]

string

,

error

)

GetRolesForUser gets the roles that a user has.

func (*SyncedEnforcer)

GetRolesForUserInDomain

¶

func (e *

SyncedEnforcer

) GetRolesForUserInDomain(name

string

, domain

string

) []

string

GetRolesForUserInDomain gets the roles that a user has inside a domain.

func (*SyncedEnforcer)

GetUsersForRole

¶

func (e *

SyncedEnforcer

) GetUsersForRole(name

string

, domain ...

string

) ([]

string

,

error

)

GetUsersForRole gets the users that has a role.

func (*SyncedEnforcer)

GetUsersForRoleInDomain

¶

func (e *

SyncedEnforcer

) GetUsersForRoleInDomain(name

string

, domain

string

) []

string

GetUsersForRoleInDomain gets the users that has a role inside a domain. Add by Gordon.

func (*SyncedEnforcer)

HasGroupingPolicy

¶

func (e *

SyncedEnforcer

) HasGroupingPolicy(params ...interface{}) (

bool

,

error

)

HasGroupingPolicy determines whether a role inheritance rule exists.

func (*SyncedEnforcer)

HasNamedGroupingPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) HasNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

HasNamedGroupingPolicy determines whether a named role inheritance rule exists.

func (*SyncedEnforcer)

HasNamedPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) HasNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

HasNamedPolicy determines whether a named authorization rule exists.

func (*SyncedEnforcer)

HasPermissionForUser

¶

func (e *

SyncedEnforcer

) HasPermissionForUser(user

string

, permission ...

string

) (

bool

,

error

)

HasPermissionForUser determines whether a user has a permission.

func (*SyncedEnforcer)

HasPolicy

¶

func (e *

SyncedEnforcer

) HasPolicy(params ...interface{}) (

bool

,

error

)

HasPolicy determines whether an authorization rule exists.

func (*SyncedEnforcer)

HasRoleForUser

¶

func (e *

SyncedEnforcer

) HasRoleForUser(name

string

, role

string

, domain ...

string

) (

bool

,

error

)

HasRoleForUser determines whether a user has a role.

func (*SyncedEnforcer)

IsAutoLoadingRunning

¶

added in

v2.11.3

func (e *

SyncedEnforcer

) IsAutoLoadingRunning()

bool

IsAutoLoadingRunning check if SyncedEnforcer is auto loading policies.

func (*SyncedEnforcer)

LoadFilteredPolicy

¶

added in

v2.2.2

func (e *

SyncedEnforcer

) LoadFilteredPolicy(filter interface{})

error

LoadFilteredPolicy reloads a filtered policy from file/database.

func (*SyncedEnforcer)

LoadIncrementalFilteredPolicy

¶

added in

v2.11.0

func (e *

SyncedEnforcer

) LoadIncrementalFilteredPolicy(filter interface{})

error

LoadIncrementalFilteredPolicy reloads a filtered policy from file/database.

func (*SyncedEnforcer)

LoadModel

¶

added in

v2.19.6

func (e *

SyncedEnforcer

) LoadModel()

error

LoadModel reloads the model from the model CONF file.

func (*SyncedEnforcer)

LoadPolicy

¶

func (e *

SyncedEnforcer

) LoadPolicy()

error

LoadPolicy reloads the policy from file/database.

func (*SyncedEnforcer)

RemoveFilteredGroupingPolicy

¶

func (e *

SyncedEnforcer

) RemoveFilteredGroupingPolicy(fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredGroupingPolicy removes a role inheritance rule from the current policy, field filters can be specified.

func (*SyncedEnforcer)

RemoveFilteredNamedGroupingPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) RemoveFilteredNamedGroupingPolicy(ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredNamedGroupingPolicy removes a role inheritance rule from the current named policy, field filters can be specified.

func (*SyncedEnforcer)

RemoveFilteredNamedPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) RemoveFilteredNamedPolicy(ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredNamedPolicy removes an authorization rule from the current named policy, field filters can be specified.

func (*SyncedEnforcer)

RemoveFilteredPolicy

¶

func (e *

SyncedEnforcer

) RemoveFilteredPolicy(fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

RemoveFilteredPolicy removes an authorization rule from the current policy, field filters can be specified.

func (*SyncedEnforcer)

RemoveGroupingPolicies

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) RemoveGroupingPolicies(rules [][]

string

) (

bool

,

error

)

RemoveGroupingPolicies removes role inheritance rules from the current policy.

func (*SyncedEnforcer)

RemoveGroupingPolicy

¶

func (e *

SyncedEnforcer

) RemoveGroupingPolicy(params ...interface{}) (

bool

,

error

)

RemoveGroupingPolicy removes a role inheritance rule from the current policy.

func (*SyncedEnforcer)

RemoveNamedGroupingPolicies

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) RemoveNamedGroupingPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveNamedGroupingPolicies removes role inheritance rules from the current named policy.

func (*SyncedEnforcer)

RemoveNamedGroupingPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) RemoveNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedGroupingPolicy removes a role inheritance rule from the current named policy.

func (*SyncedEnforcer)

RemoveNamedPolicies

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) RemoveNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveNamedPolicies removes authorization rules from the current named policy.

func (*SyncedEnforcer)

RemoveNamedPolicy

¶

added in

v2.0.2

func (e *

SyncedEnforcer

) RemoveNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedPolicy removes an authorization rule from the current named policy.

func (*SyncedEnforcer)

RemovePolicies

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) RemovePolicies(rules [][]

string

) (

bool

,

error

)

RemovePolicies removes authorization rules from the current policy.

func (*SyncedEnforcer)

RemovePolicy

¶

func (e *

SyncedEnforcer

) RemovePolicy(params ...interface{}) (

bool

,

error

)

RemovePolicy removes an authorization rule from the current policy.

func (*SyncedEnforcer)

SavePolicy

¶

func (e *

SyncedEnforcer

) SavePolicy()

error

SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.

func (*SyncedEnforcer)

SelfAddPolicies

¶

added in

v2.62.0

func (e *

SyncedEnforcer

) SelfAddPolicies(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SelfAddPoliciesEx

¶

added in

v2.63.0

func (e *

SyncedEnforcer

) SelfAddPoliciesEx(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SelfAddPolicy

¶

added in

v2.62.0

func (e *

SyncedEnforcer

) SelfAddPolicy(sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SelfRemoveFilteredPolicy

¶

added in

v2.62.0

func (e *

SyncedEnforcer

) SelfRemoveFilteredPolicy(sec

string

, ptype

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SelfRemovePolicies

¶

added in

v2.62.0

func (e *

SyncedEnforcer

) SelfRemovePolicies(sec

string

, ptype

string

, rules [][]

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SelfRemovePolicy

¶

added in

v2.62.0

func (e *

SyncedEnforcer

) SelfRemovePolicy(sec

string

, ptype

string

, rule []

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SelfUpdatePolicies

¶

added in

v2.62.0

func (e *

SyncedEnforcer

) SelfUpdatePolicies(sec

string

, ptype

string

, oldRules, newRules [][]

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SelfUpdatePolicy

¶

added in

v2.62.0

func (e *

SyncedEnforcer

) SelfUpdatePolicy(sec

string

, ptype

string

, oldRule, newRule []

string

) (

bool

,

error

)

func (*SyncedEnforcer)

SetNamedRoleManager

¶

added in

v2.125.0

func (e *

SyncedEnforcer

) SetNamedRoleManager(ptype

string

, rm

rbac

.

RoleManager

)

SetNamedRoleManager sets the role manager for the named policy with synchronization.

func (*SyncedEnforcer)

SetRoleManager

¶

added in

v2.125.0

func (e *

SyncedEnforcer

) SetRoleManager(rm

rbac

.

RoleManager

)

SetRoleManager sets the current role manager with synchronization.

func (*SyncedEnforcer)

SetWatcher

¶

func (e *

SyncedEnforcer

) SetWatcher(watcher

persist

.

Watcher

)

error

SetWatcher sets the current watcher.

func (*SyncedEnforcer)

StartAutoLoadPolicy

¶

func (e *

SyncedEnforcer

) StartAutoLoadPolicy(d

time

.

Duration

)

StartAutoLoadPolicy starts a go routine that will every specified duration call LoadPolicy.

func (*SyncedEnforcer)

StopAutoLoadPolicy

¶

func (e *

SyncedEnforcer

) StopAutoLoadPolicy()

StopAutoLoadPolicy causes the go routine to exit.

func (*SyncedEnforcer)

UpdateFilteredNamedPolicies

¶

added in

v2.28.0

func (e *

SyncedEnforcer

) UpdateFilteredNamedPolicies(ptype

string

, newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdateFilteredPolicies

¶

added in

v2.28.0

func (e *

SyncedEnforcer

) UpdateFilteredPolicies(newPolicies [][]

string

, fieldIndex

int

, fieldValues ...

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdateGroupingPolicies

¶

added in

v2.41.0

func (e *

SyncedEnforcer

) UpdateGroupingPolicies(oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdateGroupingPolicy

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) UpdateGroupingPolicy(oldRule []

string

, newRule []

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdateNamedGroupingPolicies

¶

added in

v2.41.0

func (e *

SyncedEnforcer

) UpdateNamedGroupingPolicies(ptype

string

, oldRules [][]

string

, newRules [][]

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdateNamedGroupingPolicy

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) UpdateNamedGroupingPolicy(ptype

string

, oldRule []

string

, newRule []

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdateNamedPolicies

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) UpdateNamedPolicies(ptype

string

, p1 [][]

string

, p2 [][]

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdateNamedPolicy

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) UpdateNamedPolicy(ptype

string

, p1 []

string

, p2 []

string

) (

bool

,

error

)

func (*SyncedEnforcer)

UpdatePolicies

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) UpdatePolicies(oldPolices [][]

string

, newPolicies [][]

string

) (

bool

,

error

)

UpdatePolicies updates authorization rules from the current policies.

func (*SyncedEnforcer)

UpdatePolicy

¶

added in

v2.25.1

func (e *

SyncedEnforcer

) UpdatePolicy(oldPolicy []

string

, newPolicy []

string

) (

bool

,

error

)

UpdatePolicy updates an authorization rule from the current policy.

type

Transaction

¶

added in

v2.123.0

type Transaction struct {

// contains filtered or unexported fields

}

Transaction represents a Casbin transaction.
It provides methods to perform policy operations within a transaction.
and commit or rollback all changes atomically.

func (*Transaction)

AddGroupingPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) AddGroupingPolicy(params ...interface{}) (

bool

,

error

)

AddGroupingPolicy adds a grouping policy within the transaction.

func (*Transaction)

AddNamedGroupingPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) AddNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedGroupingPolicy adds a named grouping policy within the transaction.

func (*Transaction)

AddNamedPolicies

¶

added in

v2.123.0

func (tx *

Transaction

) AddNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

AddNamedPolicies adds multiple named policies within the transaction.

func (*Transaction)

AddNamedPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) AddNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

AddNamedPolicy adds a named policy within the transaction.
The policy is buffered and will be applied when the transaction is committed.

func (*Transaction)

AddPolicies

¶

added in

v2.123.0

func (tx *

Transaction

) AddPolicies(rules [][]

string

) (

bool

,

error

)

AddPolicies adds multiple policies within the transaction.

func (*Transaction)

AddPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) AddPolicy(params ...interface{}) (

bool

,

error

)

AddPolicy adds a policy within the transaction.
The policy is buffered and will be applied when the transaction is committed.

func (*Transaction)

Commit

¶

added in

v2.123.0

func (tx *

Transaction

) Commit()

error

Commit commits the transaction using a two-phase commit protocol.
Phase 1: Apply all operations to the database
Phase 2: Apply changes to the in-memory model and rebuild role links.

func (*Transaction)

GetBufferedModel

¶

added in

v2.123.0

func (tx *

Transaction

) GetBufferedModel() (

model

.

Model

,

error

)

GetBufferedModel returns the model as it would look after applying all buffered operations.
This is useful for preview or validation purposes within the transaction.

func (*Transaction)

HasOperations

¶

added in

v2.123.0

func (tx *

Transaction

) HasOperations()

bool

HasOperations returns true if the transaction has any buffered operations.

func (*Transaction)

IsActive

¶

added in

v2.123.0

func (tx *

Transaction

) IsActive()

bool

IsActive returns true if the transaction is still active (not committed or rolled back).

func (*Transaction)

IsCommitted

¶

added in

v2.123.0

func (tx *

Transaction

) IsCommitted()

bool

IsCommitted returns true if the transaction has been committed.

func (*Transaction)

IsRolledBack

¶

added in

v2.123.0

func (tx *

Transaction

) IsRolledBack()

bool

IsRolledBack returns true if the transaction has been rolled back.

func (*Transaction)

OperationCount

¶

added in

v2.123.0

func (tx *

Transaction

) OperationCount()

int

OperationCount returns the number of buffered operations in the transaction.

func (*Transaction)

RemoveGroupingPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) RemoveGroupingPolicy(params ...interface{}) (

bool

,

error

)

RemoveGroupingPolicy removes a grouping policy within the transaction.

func (*Transaction)

RemoveNamedGroupingPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) RemoveNamedGroupingPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedGroupingPolicy removes a named grouping policy within the transaction.

func (*Transaction)

RemoveNamedPolicies

¶

added in

v2.123.0

func (tx *

Transaction

) RemoveNamedPolicies(ptype

string

, rules [][]

string

) (

bool

,

error

)

RemoveNamedPolicies removes multiple named policies within the transaction.

func (*Transaction)

RemoveNamedPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) RemoveNamedPolicy(ptype

string

, params ...interface{}) (

bool

,

error

)

RemoveNamedPolicy removes a named policy within the transaction.

func (*Transaction)

RemovePolicies

¶

added in

v2.123.0

func (tx *

Transaction

) RemovePolicies(rules [][]

string

) (

bool

,

error

)

RemovePolicies removes multiple policies within the transaction.

func (*Transaction)

RemovePolicy

¶

added in

v2.123.0

func (tx *

Transaction

) RemovePolicy(params ...interface{}) (

bool

,

error

)

RemovePolicy removes a policy within the transaction.

func (*Transaction)

Rollback

¶

added in

v2.123.0

func (tx *

Transaction

) Rollback()

error

Rollback rolls back the transaction.
This will rollback the database transaction and clear the transaction state.

func (*Transaction)

UpdateNamedPolicy

¶

added in

v2.123.0

func (tx *

Transaction

) UpdateNamedPolicy(ptype

string

, oldPolicy []

string

, newPolicy []

string

) (

bool

,

error

)

UpdateNamedPolicy updates a named policy within the transaction.

func (*Transaction)

UpdatePolicy

¶

added in

v2.123.0

func (tx *

Transaction

) UpdatePolicy(oldPolicy []

string

, newPolicy []

string

) (

bool

,

error

)

UpdatePolicy updates a policy within the transaction.

type

TransactionBuffer

¶

added in

v2.123.0

type TransactionBuffer struct {

// contains filtered or unexported fields

}

TransactionBuffer holds all policy changes made within a transaction.
It maintains a list of operations and a snapshot of the model state
at the beginning of the transaction.

func

NewTransactionBuffer

¶

added in

v2.123.0

func NewTransactionBuffer(baseModel

model

.

Model

) *

TransactionBuffer

NewTransactionBuffer creates a new transaction buffer with a model snapshot.
The snapshot represents the state of the model at the beginning of the transaction.

func (*TransactionBuffer)

AddOperation

¶

added in

v2.123.0

func (tb *

TransactionBuffer

) AddOperation(op

persist

.

PolicyOperation

)

AddOperation adds a policy operation to the buffer.
This operation will be applied when the transaction is committed.

func (*TransactionBuffer)

ApplyOperationsToModel

¶

added in

v2.123.0

func (tb *

TransactionBuffer

) ApplyOperationsToModel(baseModel

model

.

Model

) (

model

.

Model

,

error

)

ApplyOperationsToModel applies all buffered operations to a model and returns the result.
This simulates what the model would look like after all operations are applied.
It's used for validation and preview purposes within the transaction.

func (*TransactionBuffer)

Clear

¶

added in

v2.123.0

func (tb *

TransactionBuffer

) Clear()

Clear removes all buffered operations.
This is typically called after a successful commit or rollback.

func (*TransactionBuffer)

GetModelSnapshot

¶

added in

v2.123.0

func (tb *

TransactionBuffer

) GetModelSnapshot()

model

.

Model

GetModelSnapshot returns the model snapshot taken at transaction start.
This represents the original state before any transaction operations.

func (*TransactionBuffer)

GetOperations

¶

added in

v2.123.0

func (tb *

TransactionBuffer

) GetOperations() []

persist

.

PolicyOperation

GetOperations returns all buffered operations.
Returns a copy to prevent external modifications.

func (*TransactionBuffer)

HasOperations

¶

added in

v2.123.0

func (tb *

TransactionBuffer

) HasOperations()

bool

HasOperations returns true if there are any buffered operations.

func (*TransactionBuffer)

OperationCount

¶

added in

v2.123.0

func (tb *

TransactionBuffer

) OperationCount()

int

OperationCount returns the number of buffered operations.

type

TransactionalEnforcer

¶

added in

v2.123.0

type TransactionalEnforcer struct {

*

Enforcer

// Embedded enforcer for all standard functionality

// contains filtered or unexported fields

}

TransactionalEnforcer extends Enforcer with transaction support.
It provides atomic policy operations through transactions.

func

NewTransactionalEnforcer

¶

added in

v2.123.0

func NewTransactionalEnforcer(params ...interface{}) (*

TransactionalEnforcer

,

error

)

NewTransactionalEnforcer creates a new TransactionalEnforcer.
It accepts the same parameters as NewEnforcer.

func (*TransactionalEnforcer)

BeginTransaction

¶

added in

v2.123.0

func (te *

TransactionalEnforcer

) BeginTransaction(ctx

context

.

Context

) (*

Transaction

,

error

)

BeginTransaction starts a new transaction.
Returns an error if a transaction is already in progress or if the adapter doesn't support transactions.

func (*TransactionalEnforcer)

GetTransaction

¶

added in

v2.128.0

func (te *

TransactionalEnforcer

) GetTransaction(id

string

) *

Transaction

GetTransaction returns a transaction by its ID, or nil if not found.

func (*TransactionalEnforcer)

IsTransactionActive

¶

added in

v2.128.0

func (te *

TransactionalEnforcer

) IsTransactionActive(id

string

)

bool

IsTransactionActive returns true if the transaction with the given ID is active.

func (*TransactionalEnforcer)

WithTransaction

¶

added in

v2.123.0

func (te *

TransactionalEnforcer

) WithTransaction(ctx

context

.

Context

, fn func(*

Transaction

)

error

)

error

WithTransaction executes a function within a transaction.
If the function returns an error, the transaction is rolled back.
Otherwise, it's committed automatically.