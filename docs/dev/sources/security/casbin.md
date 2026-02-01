# Casbin

> Source: https://pkg.go.dev/github.com/casbin/casbin/v2
> Fetched: 2026-02-01T11:49:20.334614+00:00
> Content-Hash: 7a9c49c862d92e05
> Type: html

---

### Overview ¶

rbac_api_context.go

### Index ¶

- func CasbinJsGetPermissionForUser(e IEnforcer, user string) (string, error)
- func CasbinJsGetPermissionForUserOld(e IEnforcer, user string) ([]byte, error)
- func GetCacheKey(params ...interface{}) (string, bool)
- type CacheableParam
- type CachedEnforcer
-     * func NewCachedEnforcer(params ...interface{}) (*CachedEnforcer, error)
-     * func (e *CachedEnforcer) ClearPolicy()
  - func (e *CachedEnforcer) EnableCache(enableCache bool)
  - func (e *CachedEnforcer) Enforce(rvals ...interface{}) (bool, error)
  - func (e *CachedEnforcer) InvalidateCache() error
  - func (e *CachedEnforcer) LoadPolicy() error
  - func (e *CachedEnforcer) RemovePolicies(rules [][]string) (bool, error)
  - func (e *CachedEnforcer) RemovePolicy(params ...interface{}) (bool, error)
  - func (e *CachedEnforcer) SetCache(c cache.Cache)
  - func (e *CachedEnforcer) SetExpireTime(expireTime time.Duration)
- type ConflictDetector
-     * func NewConflictDetector(baseModel, currentModel model.Model, operations []persist.PolicyOperation) *ConflictDetector
-     * func (cd *ConflictDetector) DetectConflicts() error
- type ConflictError
-     * func (e *ConflictError) Error() string
- type ContextEnforcer
-     * func (e *ContextEnforcer) AddGroupingPoliciesCtx(ctx context.Context, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddGroupingPoliciesExCtx(ctx context.Context, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddGroupingPolicyCtx(ctx context.Context, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) AddNamedGroupingPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddNamedGroupingPoliciesExCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddNamedGroupingPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) AddNamedPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddNamedPoliciesExCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddNamedPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) AddPermissionForUserCtx(ctx context.Context, user string, permission ...string) (bool, error)
  - func (e *ContextEnforcer) AddPermissionsForUserCtx(ctx context.Context, user string, permissions ...[]string) (bool, error)
  - func (e *ContextEnforcer) AddPoliciesCtx(ctx context.Context, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddPoliciesExCtx(ctx context.Context, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) AddPolicyCtx(ctx context.Context, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) AddRoleForUserCtx(ctx context.Context, user string, role string, domain ...string) (bool, error)
  - func (e *ContextEnforcer) AddRoleForUserInDomainCtx(ctx context.Context, user string, role string, domain string) (bool, error)
  - func (e *ContextEnforcer) DeleteAllUsersByDomainCtx(ctx context.Context, domain string) (bool, error)
  - func (e *ContextEnforcer) DeleteDomainsCtx(ctx context.Context, domains ...string) (bool, error)
  - func (e *ContextEnforcer) DeletePermissionCtx(ctx context.Context, permission ...string) (bool, error)
  - func (e *ContextEnforcer) DeletePermissionForUserCtx(ctx context.Context, user string, permission ...string) (bool, error)
  - func (e *ContextEnforcer) DeletePermissionsForUserCtx(ctx context.Context, user string) (bool, error)
  - func (e *ContextEnforcer) DeleteRoleCtx(ctx context.Context, role string) (bool, error)
  - func (e *ContextEnforcer) DeleteRoleForUserCtx(ctx context.Context, user string, role string, domain ...string) (bool, error)
  - func (e *ContextEnforcer) DeleteRoleForUserInDomainCtx(ctx context.Context, user string, role string, domain string) (bool, error)
  - func (e *ContextEnforcer) DeleteRolesForUserCtx(ctx context.Context, user string, domain ...string) (bool, error)
  - func (e *ContextEnforcer) DeleteRolesForUserInDomainCtx(ctx context.Context, user string, domain string) (bool, error)
  - func (e *ContextEnforcer) DeleteUserCtx(ctx context.Context, user string) (bool, error)
  - func (e *ContextEnforcer) IsFilteredCtx(ctx context.Context) bool
  - func (e *ContextEnforcer) LoadPolicyCtx(ctx context.Context) error
  - func (e *ContextEnforcer) RemoveFilteredGroupingPolicyCtx(ctx context.Context, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *ContextEnforcer) RemoveFilteredNamedGroupingPolicyCtx(ctx context.Context, ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *ContextEnforcer) RemoveFilteredNamedPolicyCtx(ctx context.Context, ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *ContextEnforcer) RemoveFilteredPolicyCtx(ctx context.Context, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *ContextEnforcer) RemoveGroupingPoliciesCtx(ctx context.Context, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) RemoveGroupingPolicyCtx(ctx context.Context, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) RemoveNamedGroupingPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) RemoveNamedGroupingPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) RemoveNamedPoliciesCtx(ctx context.Context, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) RemoveNamedPolicyCtx(ctx context.Context, ptype string, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) RemovePoliciesCtx(ctx context.Context, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) RemovePolicyCtx(ctx context.Context, params ...interface{}) (bool, error)
  - func (e *ContextEnforcer) SavePolicyCtx(ctx context.Context) error
  - func (e *ContextEnforcer) SelfAddPoliciesCtx(ctx context.Context, sec string, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) SelfAddPoliciesExCtx(ctx context.Context, sec string, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) SelfAddPolicyCtx(ctx context.Context, sec string, ptype string, rule []string) (bool, error)
  - func (e *ContextEnforcer) SelfRemoveFilteredPolicyCtx(ctx context.Context, sec string, ptype string, fieldIndex int, ...) (bool, error)
  - func (e *ContextEnforcer) SelfRemovePoliciesCtx(ctx context.Context, sec string, ptype string, rules [][]string) (bool, error)
  - func (e *ContextEnforcer) SelfRemovePolicyCtx(ctx context.Context, sec string, ptype string, rule []string) (bool, error)
  - func (e *ContextEnforcer) SelfUpdatePoliciesCtx(ctx context.Context, sec string, ptype string, oldRules, newRules [][]string) (bool, error)
  - func (e *ContextEnforcer) SelfUpdatePolicyCtx(ctx context.Context, sec string, ptype string, oldRule, newRule []string) (bool, error)
  - func (e *ContextEnforcer) UpdateFilteredNamedPoliciesCtx(ctx context.Context, ptype string, newPolicies [][]string, fieldIndex int, ...) (bool, error)
  - func (e *ContextEnforcer) UpdateFilteredPoliciesCtx(ctx context.Context, newPolicies [][]string, fieldIndex int, ...) (bool, error)
  - func (e *ContextEnforcer) UpdateGroupingPoliciesCtx(ctx context.Context, oldRules [][]string, newRules [][]string) (bool, error)
  - func (e *ContextEnforcer) UpdateGroupingPolicyCtx(ctx context.Context, oldRule []string, newRule []string) (bool, error)
  - func (e *ContextEnforcer) UpdateNamedGroupingPoliciesCtx(ctx context.Context, ptype string, oldRules [][]string, newRules [][]string) (bool, error)
  - func (e *ContextEnforcer) UpdateNamedGroupingPolicyCtx(ctx context.Context, ptype string, oldRule []string, newRule []string) (bool, error)
  - func (e *ContextEnforcer) UpdateNamedPoliciesCtx(ctx context.Context, ptype string, p1 [][]string, p2 [][]string) (bool, error)
  - func (e *ContextEnforcer) UpdateNamedPolicyCtx(ctx context.Context, ptype string, p1 []string, p2 []string) (bool, error)
  - func (e *ContextEnforcer) UpdatePoliciesCtx(ctx context.Context, oldPolicies [][]string, newPolicies [][]string) (bool, error)
  - func (e *ContextEnforcer) UpdatePolicyCtx(ctx context.Context, oldPolicy []string, newPolicy []string) (bool, error)
- type DistributedEnforcer
-     * func NewDistributedEnforcer(params ...interface{}) (*DistributedEnforcer, error)
-     * func (d *DistributedEnforcer) AddPoliciesSelf(shouldPersist func() bool, sec string, ptype string, rules [][]string) (affected [][]string, err error)
  - func (d *DistributedEnforcer) ClearPolicySelf(shouldPersist func() bool) error
  - func (d *DistributedEnforcer) RemoveFilteredPolicySelf(shouldPersist func() bool, sec string, ptype string, fieldIndex int, ...) (affected [][]string, err error)
  - func (d *DistributedEnforcer) RemovePoliciesSelf(shouldPersist func() bool, sec string, ptype string, rules [][]string) (affected [][]string, err error)
  - func (d *DistributedEnforcer) SetDispatcher(dispatcher persist.Dispatcher)
  - func (d *DistributedEnforcer) UpdateFilteredPoliciesSelf(shouldPersist func() bool, sec string, ptype string, newRules [][]string, ...) (bool, error)
  - func (d *DistributedEnforcer) UpdatePoliciesSelf(shouldPersist func() bool, sec string, ptype string, ...) (affected bool, err error)
  - func (d *DistributedEnforcer) UpdatePolicySelf(shouldPersist func() bool, sec string, ptype string, oldRule, newRule []string) (affected bool, err error)
- type EnforceContext
-     * func NewEnforceContext(suffix string) EnforceContext
-     * func (e EnforceContext) GetCacheKey() string
- type Enforcer
-     * func NewEnforcer(params ...interface{}) (*Enforcer, error)
-     * func (e *Enforcer) AddFunction(name string, function govaluate.ExpressionFunction)
  - func (e *Enforcer) AddGroupingPolicies(rules [][]string) (bool, error)
  - func (e *Enforcer) AddGroupingPoliciesEx(rules [][]string) (bool, error)
  - func (e *Enforcer) AddGroupingPolicy(params ...interface{}) (bool, error)
  - func (e *Enforcer) AddNamedDomainLinkConditionFunc(ptype, user, role string, domain string, fn rbac.LinkConditionFunc) bool
  - func (e *Enforcer) AddNamedDomainMatchingFunc(ptype, name string, fn rbac.MatchingFunc) bool
  - func (e *Enforcer) AddNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) AddNamedGroupingPoliciesEx(ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *Enforcer) AddNamedLinkConditionFunc(ptype, user, role string, fn rbac.LinkConditionFunc) bool
  - func (e *Enforcer) AddNamedMatchingFunc(ptype, name string, fn rbac.MatchingFunc) bool
  - func (e *Enforcer) AddNamedPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) AddNamedPoliciesEx(ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) AddNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *Enforcer) AddPermissionForUser(user string, permission ...string) (bool, error)
  - func (e *Enforcer) AddPermissionsForUser(user string, permissions ...[]string) (bool, error)
  - func (e *Enforcer) AddPolicies(rules [][]string) (bool, error)
  - func (e *Enforcer) AddPoliciesEx(rules [][]string) (bool, error)
  - func (e *Enforcer) AddPolicy(params ...interface{}) (bool, error)
  - func (e *Enforcer) AddRoleForUser(user string, role string, domain ...string) (bool, error)
  - func (e *Enforcer) AddRoleForUserInDomain(user string, role string, domain string) (bool, error)
  - func (e *Enforcer) AddRolesForUser(user string, roles []string, domain ...string) (bool, error)
  - func (e *Enforcer) BatchEnforce(requests [][]interface{}) ([]bool, error)
  - func (e *Enforcer) BatchEnforceWithMatcher(matcher string, requests [][]interface{}) ([]bool, error)
  - func (e *Enforcer) BuildIncrementalConditionalRoleLinks(op model.PolicyOp, ptype string, rules [][]string) error
  - func (e *Enforcer) BuildIncrementalRoleLinks(op model.PolicyOp, ptype string, rules [][]string) error
  - func (e *Enforcer) BuildRoleLinks() error
  - func (e *Enforcer) ClearPolicy()
  - func (e *Enforcer) DeleteAllUsersByDomain(domain string) (bool, error)
  - func (e *Enforcer) DeleteDomains(domains ...string) (bool, error)
  - func (e *Enforcer) DeletePermission(permission ...string) (bool, error)
  - func (e *Enforcer) DeletePermissionForUser(user string, permission ...string) (bool, error)
  - func (e *Enforcer) DeletePermissionsForUser(user string) (bool, error)
  - func (e *Enforcer) DeleteRole(role string) (bool, error)
  - func (e *Enforcer) DeleteRoleForUser(user string, role string, domain ...string) (bool, error)
  - func (e *Enforcer) DeleteRoleForUserInDomain(user string, role string, domain string) (bool, error)
  - func (e *Enforcer) DeleteRolesForUser(user string, domain ...string) (bool, error)
  - func (e *Enforcer) DeleteRolesForUserInDomain(user string, domain string) (bool, error)
  - func (e *Enforcer) DeleteUser(user string) (bool, error)
  - func (e *Enforcer) EnableAcceptJsonRequest(acceptJsonRequest bool)
  - func (e *Enforcer) EnableAutoBuildRoleLinks(autoBuildRoleLinks bool)
  - func (e *Enforcer) EnableAutoNotifyDispatcher(enable bool)
  - func (e *Enforcer) EnableAutoNotifyWatcher(enable bool)
  - func (e *Enforcer) EnableAutoSave(autoSave bool)
  - func (e *Enforcer) EnableEnforce(enable bool)
  - func (e *Enforcer) EnableLog(enable bool)
  - func (e *Enforcer) Enforce(rvals ...interface{}) (bool, error)
  - func (e *Enforcer) EnforceEx(rvals ...interface{}) (bool, []string, error)
  - func (e *Enforcer) EnforceExWithMatcher(matcher string, rvals ...interface{}) (bool, []string, error)
  - func (e *Enforcer) EnforceWithMatcher(matcher string, rvals ...interface{}) (bool, error)
  - func (e *Enforcer) GetAdapter() persist.Adapter
  - func (e *Enforcer) GetAllActions() ([]string, error)
  - func (e *Enforcer) GetAllDomains() ([]string, error)
  - func (e *Enforcer) GetAllNamedActions(ptype string) ([]string, error)
  - func (e *Enforcer) GetAllNamedObjects(ptype string) ([]string, error)
  - func (e *Enforcer) GetAllNamedRoles(ptype string) ([]string, error)
  - func (e *Enforcer) GetAllNamedSubjects(ptype string) ([]string, error)
  - func (e *Enforcer) GetAllObjects() ([]string, error)
  - func (e *Enforcer) GetAllRoles() ([]string, error)
  - func (e *Enforcer) GetAllRolesByDomain(domain string) ([]string, error)
  - func (e *Enforcer) GetAllSubjects() ([]string, error)
  - func (e *Enforcer) GetAllUsersByDomain(domain string) ([]string, error)
  - func (e *Enforcer) GetAllowedObjectConditions(user string, action string, prefix string) ([]string, error)
  - func (e *Enforcer) GetDomainsForUser(user string) ([]string, error)
  - func (e *Enforcer) GetFieldIndex(ptype string, field string) (int, error)
  - func (e *Enforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *Enforcer) GetFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *Enforcer) GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *Enforcer) GetFilteredNamedPolicyWithMatcher(ptype string, matcher string) ([][]string, error)
  - func (e *Enforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *Enforcer) GetGroupingPolicy() ([][]string, error)
  - func (e *Enforcer) GetImplicitObjectPatternsForUser(user string, domain string, action string) ([]string, error)
  - func (e *Enforcer) GetImplicitPermissionsForUser(user string, domain ...string) ([][]string, error)
  - func (e *Enforcer) GetImplicitResourcesForUser(user string, domain ...string) ([][]string, error)
  - func (e *Enforcer) GetImplicitRolesForUser(name string, domain ...string) ([]string, error)
  - func (e *Enforcer) GetImplicitUsersForPermission(permission ...string) ([]string, error)
  - func (e *Enforcer) GetImplicitUsersForResource(resource string) ([][]string, error)
  - func (e *Enforcer) GetImplicitUsersForResourceByDomain(resource string, domain string) ([][]string, error)
  - func (e *Enforcer) GetImplicitUsersForRole(name string, domain ...string) ([]string, error)
  - func (e *Enforcer) GetModel() model.Model
  - func (e *Enforcer) GetNamedGroupingPolicy(ptype string) ([][]string, error)
  - func (e *Enforcer) GetNamedImplicitPermissionsForUser(ptype string, gtype string, user string, domain ...string) ([][]string, error)
  - func (e *Enforcer) GetNamedImplicitRolesForUser(ptype string, name string, domain ...string) ([]string, error)
  - func (e *Enforcer) GetNamedImplicitUsersForResource(ptype string, resource string) ([][]string, error)
  - func (e *Enforcer) GetNamedPermissionsForUser(ptype string, user string, domain ...string) ([][]string, error)
  - func (e *Enforcer) GetNamedPolicy(ptype string) ([][]string, error)
  - func (e *Enforcer) GetNamedRoleManager(ptype string) rbac.RoleManager
  - func (e *Enforcer) GetPermissionsForUser(user string, domain ...string) ([][]string, error)
  - func (e *Enforcer) GetPermissionsForUserInDomain(user string, domain string) [][]string
  - func (e *Enforcer) GetPolicy() ([][]string, error)
  - func (e *Enforcer) GetRoleManager() rbac.RoleManager
  - func (e *Enforcer) GetRolesForUser(name string, domain ...string) ([]string, error)
  - func (e *Enforcer) GetRolesForUserInDomain(name string, domain string) []string
  - func (e *Enforcer) GetUsersForRole(name string, domain ...string) ([]string, error)
  - func (e *Enforcer) GetUsersForRoleInDomain(name string, domain string) []string
  - func (e *Enforcer) HasGroupingPolicy(params ...interface{}) (bool, error)
  - func (e *Enforcer) HasNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *Enforcer) HasNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *Enforcer) HasPermissionForUser(user string, permission ...string) (bool, error)
  - func (e *Enforcer) HasPolicy(params ...interface{}) (bool, error)
  - func (e *Enforcer) HasRoleForUser(name string, role string, domain ...string) (bool, error)
  - func (e *Enforcer) InitWithAdapter(modelPath string, adapter persist.Adapter) error
  - func (e *Enforcer) InitWithFile(modelPath string, policyPath string) error
  - func (e *Enforcer) InitWithModelAndAdapter(m model.Model, adapter persist.Adapter) error
  - func (e *Enforcer) IsFiltered() bool
  - func (e *Enforcer) IsLogEnabled() bool
  - func (e *Enforcer) LoadFilteredPolicy(filter interface{}) error
  - func (e *Enforcer) LoadFilteredPolicyCtx(ctx context.Context, filter interface{}) error
  - func (e *Enforcer) LoadIncrementalFilteredPolicy(filter interface{}) error
  - func (e *Enforcer) LoadIncrementalFilteredPolicyCtx(ctx context.Context, filter interface{}) error
  - func (e *Enforcer) LoadModel() error
  - func (e *Enforcer) LoadPolicy() error
  - func (e *Enforcer) RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *Enforcer) RemoveFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *Enforcer) RemoveFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *Enforcer) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *Enforcer) RemoveGroupingPolicies(rules [][]string) (bool, error)
  - func (e *Enforcer) RemoveGroupingPolicy(params ...interface{}) (bool, error)
  - func (e *Enforcer) RemoveNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *Enforcer) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *Enforcer) RemovePolicies(rules [][]string) (bool, error)
  - func (e *Enforcer) RemovePolicy(params ...interface{}) (bool, error)
  - func (e *Enforcer) SavePolicy() error
  - func (e *Enforcer) SelfAddPolicies(sec string, ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) SelfAddPoliciesEx(sec string, ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) SelfAddPolicy(sec string, ptype string, rule []string) (bool, error)
  - func (e *Enforcer) SelfRemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *Enforcer) SelfRemovePolicies(sec string, ptype string, rules [][]string) (bool, error)
  - func (e *Enforcer) SelfRemovePolicy(sec string, ptype string, rule []string) (bool, error)
  - func (e *Enforcer) SelfUpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) (bool, error)
  - func (e *Enforcer) SelfUpdatePolicy(sec string, ptype string, oldRule, newRule []string) (bool, error)
  - func (e *Enforcer) SetAdapter(adapter persist.Adapter)
  - func (e *Enforcer) SetEffector(eft effector.Effector)
  - func (e *Enforcer) SetFieldIndex(ptype string, field string, index int)
  - func (e *Enforcer) SetLogger(logger log.Logger)
  - func (e *Enforcer) SetModel(m model.Model)
  - func (e *Enforcer) SetNamedDomainLinkConditionFuncParams(ptype, user, role, domain string, params ...string) bool
  - func (e *Enforcer) SetNamedLinkConditionFuncParams(ptype, user, role string, params ...string) bool
  - func (e *Enforcer) SetNamedRoleManager(ptype string, rm rbac.RoleManager)
  - func (e *Enforcer) SetRoleManager(rm rbac.RoleManager)
  - func (e *Enforcer) SetWatcher(watcher persist.Watcher) error
  - func (e *Enforcer) UpdateFilteredNamedPolicies(ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *Enforcer) UpdateFilteredPolicies(newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *Enforcer) UpdateGroupingPolicies(oldRules [][]string, newRules [][]string) (bool, error)
  - func (e *Enforcer) UpdateGroupingPolicy(oldRule []string, newRule []string) (bool, error)
  - func (e *Enforcer) UpdateNamedGroupingPolicies(ptype string, oldRules [][]string, newRules [][]string) (bool, error)
  - func (e *Enforcer) UpdateNamedGroupingPolicy(ptype string, oldRule []string, newRule []string) (bool, error)
  - func (e *Enforcer) UpdateNamedPolicies(ptype string, p1 [][]string, p2 [][]string) (bool, error)
  - func (e *Enforcer) UpdateNamedPolicy(ptype string, p1 []string, p2 []string) (bool, error)
  - func (e *Enforcer) UpdatePolicies(oldPolices [][]string, newPolicies [][]string) (bool, error)
  - func (e *Enforcer) UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error)
- type IDistributedEnforcer
- type IEnforcer
- type IEnforcerContext
-     * func NewContextEnforcer(params ...interface{}) (IEnforcerContext, error)
- type SyncedCachedEnforcer
-     * func NewSyncedCachedEnforcer(params ...interface{}) (*SyncedCachedEnforcer, error)
-     * func (e *SyncedCachedEnforcer) AddPolicies(rules [][]string) (bool, error)
  - func (e *SyncedCachedEnforcer) AddPolicy(params ...interface{}) (bool, error)
  - func (e *SyncedCachedEnforcer) EnableCache(enableCache bool)
  - func (e *SyncedCachedEnforcer) Enforce(rvals ...interface{}) (bool, error)
  - func (e *SyncedCachedEnforcer) InvalidateCache() error
  - func (e *SyncedCachedEnforcer) LoadPolicy() error
  - func (e *SyncedCachedEnforcer) RemovePolicies(rules [][]string) (bool, error)
  - func (e *SyncedCachedEnforcer) RemovePolicy(params ...interface{}) (bool, error)
  - func (e *SyncedCachedEnforcer) SetCache(c cache.Cache)
  - func (e *SyncedCachedEnforcer) SetExpireTime(expireTime time.Duration)
- type SyncedEnforcer
-     * func NewSyncedEnforcer(params ...interface{}) (*SyncedEnforcer, error)
-     * func (e *SyncedEnforcer) AddFunction(name string, function govaluate.ExpressionFunction)
  - func (e *SyncedEnforcer) AddGroupingPolicies(rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddGroupingPoliciesEx(rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddGroupingPolicy(params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) AddNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddNamedGroupingPoliciesEx(ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) AddNamedPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddNamedPoliciesEx(ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) AddPermissionForUser(user string, permission ...string) (bool, error)
  - func (e *SyncedEnforcer) AddPermissionsForUser(user string, permissions ...[]string) (bool, error)
  - func (e *SyncedEnforcer) AddPolicies(rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddPoliciesEx(rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) AddPolicy(params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) AddRoleForUser(user string, role string, domain ...string) (bool, error)
  - func (e *SyncedEnforcer) AddRoleForUserInDomain(user string, role string, domain string) (bool, error)
  - func (e *SyncedEnforcer) AddRolesForUser(user string, roles []string, domain ...string) (bool, error)
  - func (e *SyncedEnforcer) BatchEnforce(requests [][]interface{}) ([]bool, error)
  - func (e *SyncedEnforcer) BatchEnforceWithMatcher(matcher string, requests [][]interface{}) ([]bool, error)
  - func (e *SyncedEnforcer) BuildRoleLinks() error
  - func (e *SyncedEnforcer) ClearPolicy()
  - func (e *SyncedEnforcer) DeleteDomains(domains ...string) (bool, error)
  - func (e *SyncedEnforcer) DeletePermission(permission ...string) (bool, error)
  - func (e *SyncedEnforcer) DeletePermissionForUser(user string, permission ...string) (bool, error)
  - func (e *SyncedEnforcer) DeletePermissionsForUser(user string) (bool, error)
  - func (e *SyncedEnforcer) DeleteRole(role string) (bool, error)
  - func (e *SyncedEnforcer) DeleteRoleForUser(user string, role string, domain ...string) (bool, error)
  - func (e *SyncedEnforcer) DeleteRoleForUserInDomain(user string, role string, domain string) (bool, error)
  - func (e *SyncedEnforcer) DeleteRolesForUser(user string, domain ...string) (bool, error)
  - func (e *SyncedEnforcer) DeleteRolesForUserInDomain(user string, domain string) (bool, error)
  - func (e *SyncedEnforcer) DeleteUser(user string) (bool, error)
  - func (e *SyncedEnforcer) Enforce(rvals ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) EnforceEx(rvals ...interface{}) (bool, []string, error)
  - func (e *SyncedEnforcer) EnforceExWithMatcher(matcher string, rvals ...interface{}) (bool, []string, error)
  - func (e *SyncedEnforcer) EnforceWithMatcher(matcher string, rvals ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) GetAllActions() ([]string, error)
  - func (e *SyncedEnforcer) GetAllNamedActions(ptype string) ([]string, error)
  - func (e *SyncedEnforcer) GetAllNamedObjects(ptype string) ([]string, error)
  - func (e *SyncedEnforcer) GetAllNamedRoles(ptype string) ([]string, error)
  - func (e *SyncedEnforcer) GetAllNamedSubjects(ptype string) ([]string, error)
  - func (e *SyncedEnforcer) GetAllObjects() ([]string, error)
  - func (e *SyncedEnforcer) GetAllRoles() ([]string, error)
  - func (e *SyncedEnforcer) GetAllSubjects() ([]string, error)
  - func (e *SyncedEnforcer) GetFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetGroupingPolicy() ([][]string, error)
  - func (e *SyncedEnforcer) GetImplicitObjectPatternsForUser(user string, domain string, action string) ([]string, error)
  - func (e *SyncedEnforcer) GetImplicitPermissionsForUser(user string, domain ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetImplicitRolesForUser(name string, domain ...string) ([]string, error)
  - func (e *SyncedEnforcer) GetImplicitUsersForPermission(permission ...string) ([]string, error)
  - func (e *SyncedEnforcer) GetLock()*sync.RWMutex
  - func (e *SyncedEnforcer) GetNamedGroupingPolicy(ptype string) ([][]string, error)
  - func (e *SyncedEnforcer) GetNamedImplicitPermissionsForUser(ptype string, gtype string, user string, domain ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetNamedPermissionsForUser(ptype string, user string, domain ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetNamedPolicy(ptype string) ([][]string, error)
  - func (e *SyncedEnforcer) GetNamedRoleManager(ptype string) rbac.RoleManager
  - func (e *SyncedEnforcer) GetPermissionsForUser(user string, domain ...string) ([][]string, error)
  - func (e *SyncedEnforcer) GetPermissionsForUserInDomain(user string, domain string) [][]string
  - func (e *SyncedEnforcer) GetPolicy() ([][]string, error)
  - func (e *SyncedEnforcer) GetRoleManager() rbac.RoleManager
  - func (e *SyncedEnforcer) GetRolesForUser(name string, domain ...string) ([]string, error)
  - func (e *SyncedEnforcer) GetRolesForUserInDomain(name string, domain string) []string
  - func (e *SyncedEnforcer) GetUsersForRole(name string, domain ...string) ([]string, error)
  - func (e *SyncedEnforcer) GetUsersForRoleInDomain(name string, domain string) []string
  - func (e *SyncedEnforcer) HasGroupingPolicy(params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) HasNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) HasNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) HasPermissionForUser(user string, permission ...string) (bool, error)
  - func (e *SyncedEnforcer) HasPolicy(params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) HasRoleForUser(name string, role string, domain ...string) (bool, error)
  - func (e *SyncedEnforcer) IsAutoLoadingRunning() bool
  - func (e *SyncedEnforcer) LoadFilteredPolicy(filter interface{}) error
  - func (e *SyncedEnforcer) LoadIncrementalFilteredPolicy(filter interface{}) error
  - func (e *SyncedEnforcer) LoadModel() error
  - func (e *SyncedEnforcer) LoadPolicy() error
  - func (e *SyncedEnforcer) RemoveFilteredGroupingPolicy(fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *SyncedEnforcer) RemoveFilteredNamedGroupingPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *SyncedEnforcer) RemoveFilteredNamedPolicy(ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *SyncedEnforcer) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *SyncedEnforcer) RemoveGroupingPolicies(rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) RemoveGroupingPolicy(params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) RemoveNamedGroupingPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) RemovePolicies(rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) RemovePolicy(params ...interface{}) (bool, error)
  - func (e *SyncedEnforcer) SavePolicy() error
  - func (e *SyncedEnforcer) SelfAddPolicies(sec string, ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) SelfAddPoliciesEx(sec string, ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) SelfAddPolicy(sec string, ptype string, rule []string) (bool, error)
  - func (e *SyncedEnforcer) SelfRemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *SyncedEnforcer) SelfRemovePolicies(sec string, ptype string, rules [][]string) (bool, error)
  - func (e *SyncedEnforcer) SelfRemovePolicy(sec string, ptype string, rule []string) (bool, error)
  - func (e *SyncedEnforcer) SelfUpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) (bool, error)
  - func (e *SyncedEnforcer) SelfUpdatePolicy(sec string, ptype string, oldRule, newRule []string) (bool, error)
  - func (e *SyncedEnforcer) SetNamedRoleManager(ptype string, rm rbac.RoleManager)
  - func (e *SyncedEnforcer) SetRoleManager(rm rbac.RoleManager)
  - func (e *SyncedEnforcer) SetWatcher(watcher persist.Watcher) error
  - func (e *SyncedEnforcer) StartAutoLoadPolicy(d time.Duration)
  - func (e *SyncedEnforcer) StopAutoLoadPolicy()
  - func (e *SyncedEnforcer) UpdateFilteredNamedPolicies(ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *SyncedEnforcer) UpdateFilteredPolicies(newPolicies [][]string, fieldIndex int, fieldValues ...string) (bool, error)
  - func (e *SyncedEnforcer) UpdateGroupingPolicies(oldRules [][]string, newRules [][]string) (bool, error)
  - func (e *SyncedEnforcer) UpdateGroupingPolicy(oldRule []string, newRule []string) (bool, error)
  - func (e *SyncedEnforcer) UpdateNamedGroupingPolicies(ptype string, oldRules [][]string, newRules [][]string) (bool, error)
  - func (e *SyncedEnforcer) UpdateNamedGroupingPolicy(ptype string, oldRule []string, newRule []string) (bool, error)
  - func (e *SyncedEnforcer) UpdateNamedPolicies(ptype string, p1 [][]string, p2 [][]string) (bool, error)
  - func (e *SyncedEnforcer) UpdateNamedPolicy(ptype string, p1 []string, p2 []string) (bool, error)
  - func (e *SyncedEnforcer) UpdatePolicies(oldPolices [][]string, newPolicies [][]string) (bool, error)
  - func (e *SyncedEnforcer) UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error)
- type Transaction
-     * func (tx *Transaction) AddGroupingPolicy(params ...interface{}) (bool, error)
  - func (tx *Transaction) AddNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (tx *Transaction) AddNamedPolicies(ptype string, rules [][]string) (bool, error)
  - func (tx *Transaction) AddNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (tx *Transaction) AddPolicies(rules [][]string) (bool, error)
  - func (tx *Transaction) AddPolicy(params ...interface{}) (bool, error)
  - func (tx *Transaction) Commit() error
  - func (tx *Transaction) GetBufferedModel() (model.Model, error)
  - func (tx *Transaction) HasOperations() bool
  - func (tx *Transaction) IsActive() bool
  - func (tx *Transaction) IsCommitted() bool
  - func (tx *Transaction) IsRolledBack() bool
  - func (tx *Transaction) OperationCount() int
  - func (tx *Transaction) RemoveGroupingPolicy(params ...interface{}) (bool, error)
  - func (tx *Transaction) RemoveNamedGroupingPolicy(ptype string, params ...interface{}) (bool, error)
  - func (tx *Transaction) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error)
  - func (tx *Transaction) RemoveNamedPolicy(ptype string, params ...interface{}) (bool, error)
  - func (tx *Transaction) RemovePolicies(rules [][]string) (bool, error)
  - func (tx *Transaction) RemovePolicy(params ...interface{}) (bool, error)
  - func (tx *Transaction) Rollback() error
  - func (tx *Transaction) UpdateNamedPolicy(ptype string, oldPolicy []string, newPolicy []string) (bool, error)
  - func (tx *Transaction) UpdatePolicy(oldPolicy []string, newPolicy []string) (bool, error)
- type TransactionBuffer
-     * func NewTransactionBuffer(baseModel model.Model) *TransactionBuffer
-     * func (tb *TransactionBuffer) AddOperation(op persist.PolicyOperation)
  - func (tb *TransactionBuffer) ApplyOperationsToModel(baseModel model.Model) (model.Model, error)
  - func (tb *TransactionBuffer) Clear()
  - func (tb *TransactionBuffer) GetModelSnapshot() model.Model
  - func (tb *TransactionBuffer) GetOperations() []persist.PolicyOperation
  - func (tb *TransactionBuffer) HasOperations() bool
  - func (tb *TransactionBuffer) OperationCount() int
- type TransactionalEnforcer
-     * func NewTransactionalEnforcer(params ...interface{}) (*TransactionalEnforcer, error)
-     * func (te *TransactionalEnforcer) BeginTransaction(ctx context.Context) (*Transaction, error)
  - func (te *TransactionalEnforcer) GetTransaction(id string)*Transaction
  - func (te *TransactionalEnforcer) IsTransactionActive(id string) bool
  - func (te *TransactionalEnforcer) WithTransaction(ctx context.Context, fn func(*Transaction) error) error

### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

#### func [CasbinJsGetPermissionForUser](https://github.com/casbin/casbin/blob/v2.135.0/frontend.go#L22) ¶ added in v2.9.0

    func CasbinJsGetPermissionForUser(e IEnforcer, user [string](/builtin#string)) ([string](/builtin#string), [error](/builtin#error))

#### func [CasbinJsGetPermissionForUserOld](https://github.com/casbin/casbin/blob/v2.135.0/frontend_old.go#L19) ¶ added in v2.31.1

    func CasbinJsGetPermissionForUserOld(e IEnforcer, user [string](/builtin#string)) ([][byte](/builtin#byte), [error](/builtin#error))

#### func [GetCacheKey](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L160) ¶ added in v2.66.0

    func GetCacheKey(params ...interface{}) ([string](/builtin#string), [bool](/builtin#bool))

### Types ¶

#### type [CacheableParam](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L35) ¶ added in v2.40.0

    type CacheableParam interface {
     GetCacheKey() [string](/builtin#string)
    }

#### type [CachedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L27) ¶

    type CachedEnforcer struct {
     *Enforcer
     // contains filtered or unexported fields
    }

CachedEnforcer wraps Enforcer and provides decision cache.

#### func [NewCachedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L40) ¶

    func NewCachedEnforcer(params ...interface{}) (*CachedEnforcer, [error](/builtin#error))

NewCachedEnforcer creates a cached enforcer via file or DB.

#### func (*CachedEnforcer) [ClearPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L177) ¶ added in v2.97.0

    func (e *CachedEnforcer) ClearPolicy()

ClearPolicy clears all policy.

#### func (*CachedEnforcer) [EnableCache](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L55) ¶

    func (e *CachedEnforcer) EnableCache(enableCache [bool](/builtin#bool))

EnableCache determines whether to enable cache on Enforce(). When enableCache is enabled, cached result (true | false) will be returned for previous decisions.

#### func (*CachedEnforcer) [Enforce](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L65) ¶

    func (e *CachedEnforcer) Enforce(rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act). if rvals is not string , ignore the cache.

#### func (*CachedEnforcer) [InvalidateCache](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L154) ¶

    func (e *CachedEnforcer) InvalidateCache() [error](/builtin#error)

InvalidateCache deletes all the existing cached decisions.

#### func (*CachedEnforcer) [LoadPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L90) ¶ added in v2.32.0

    func (e *CachedEnforcer) LoadPolicy() [error](/builtin#error)

#### func (*CachedEnforcer) [RemovePolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L111) ¶ added in v2.32.0

    func (e *CachedEnforcer) RemovePolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*CachedEnforcer) [RemovePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L99) ¶ added in v2.32.0

    func (e *CachedEnforcer) RemovePolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*CachedEnforcer) [SetCache](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L139) ¶ added in v2.32.0

    func (e *CachedEnforcer) SetCache(c [cache](/github.com/casbin/casbin/v2@v2.135.0/persist/cache).[Cache](/github.com/casbin/casbin/v2@v2.135.0/persist/cache#Cache))

#### func (*CachedEnforcer) [SetExpireTime](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached.go#L135) ¶ added in v2.32.0

    func (e *CachedEnforcer) SetExpireTime(expireTime [time](/time).[Duration](/time#Duration))

#### type [ConflictDetector](https://github.com/casbin/casbin/blob/v2.135.0/transaction_conflict.go#L35) ¶ added in v2.128.0

    type ConflictDetector struct {
     // contains filtered or unexported fields
    }

ConflictDetector detects conflicts between transaction operations and current model state.

#### func [NewConflictDetector](https://github.com/casbin/casbin/blob/v2.135.0/transaction_conflict.go#L42) ¶ added in v2.128.0

    func NewConflictDetector(baseModel, currentModel [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model), operations [][persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[PolicyOperation](/github.com/casbin/casbin/v2@v2.135.0/persist#PolicyOperation)) *ConflictDetector

NewConflictDetector creates a new conflict detector instance.

#### func (*ConflictDetector) [DetectConflicts](https://github.com/casbin/casbin/blob/v2.135.0/transaction_conflict.go#L52) ¶ added in v2.128.0

    func (cd *ConflictDetector) DetectConflicts() [error](/builtin#error)

DetectConflicts checks for conflicts between the transaction operations and current model state. Returns nil if no conflicts are found, otherwise returns a ConflictError describing the conflict.

#### type [ConflictError](https://github.com/casbin/casbin/blob/v2.135.0/transaction_conflict.go#L25) ¶ added in v2.128.0

    type ConflictError struct {
     Operation [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[PolicyOperation](/github.com/casbin/casbin/v2@v2.135.0/persist#PolicyOperation)
     Reason    [string](/builtin#string)
    }

ConflictError represents a transaction conflict error.

#### func (*ConflictError) [Error](https://github.com/casbin/casbin/blob/v2.135.0/transaction_conflict.go#L30) ¶ added in v2.128.0

    func (e *ConflictError) Error() [string](/builtin#string)

#### type [ContextEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L28) ¶ added in v2.128.0

    type ContextEnforcer struct {
     *Enforcer
     // contains filtered or unexported fields
    }

ContextEnforcer wraps Enforcer and provides context-aware operations.

#### func (*ContextEnforcer) [AddGroupingPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L274) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPoliciesCtx adds grouping policy rules to the storage with context.

#### func (*ContextEnforcer) [AddGroupingPoliciesExCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L278) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddGroupingPoliciesExCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*ContextEnforcer) [AddGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L269) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPolicyCtx adds a grouping policy rule to the storage with context.

#### func (*ContextEnforcer) [AddNamedGroupingPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L300) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddNamedGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPoliciesCtx adds named grouping policy rules to the storage with context.

#### func (*ContextEnforcer) [AddNamedGroupingPoliciesExCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L304) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddNamedGroupingPoliciesExCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*ContextEnforcer) [AddNamedGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L283) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPolicyCtx adds a named grouping policy rule to the storage with context.

#### func (*ContextEnforcer) [AddNamedPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L186) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddNamedPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPoliciesCtx adds named policy rules to the storage with context.

#### func (*ContextEnforcer) [AddNamedPoliciesExCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L194) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddNamedPoliciesExCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*ContextEnforcer) [AddNamedPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L172) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddNamedPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPolicyCtx adds a named policy rule to the storage with context.

#### func (*ContextEnforcer) [AddPermissionForUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L103) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddPermissionForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPermissionForUserCtx adds a permission for a user or role with context support. Returns false if the user or role already has the permission (aka not affected).

#### func (*ContextEnforcer) [AddPermissionsForUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L109) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddPermissionsForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), permissions ...[][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPermissionsForUserCtx adds multiple permissions for a user or role with context support. Returns false if the user or role already has one of the permissions (aka not affected).

#### func (*ContextEnforcer) [AddPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L167) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddPoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPoliciesCtx adds policy rules to the storage with context.

#### func (*ContextEnforcer) [AddPoliciesExCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L190) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddPoliciesExCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*ContextEnforcer) [AddPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L162) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddPolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddPolicyCtx adds a policy rule to the storage with context.

#### func (*ContextEnforcer) [AddRoleForUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L28) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddRoleForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRoleForUserCtx adds a role for a user with context support. Returns false if the user already has the role (aka not affected).

#### func (*ContextEnforcer) [AddRoleForUserInDomainCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_context.go#L26) ¶ added in v2.128.0

    func (e *ContextEnforcer) AddRoleForUserInDomainCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRoleForUserInDomainCtx adds a role for a user inside a domain with context support. Returns false if the user already has the role (aka not affected).

#### func (*ContextEnforcer) [DeleteAllUsersByDomainCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_context.go#L56) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteAllUsersByDomainCtx(ctx [context](/context).[Context](/context#Context), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteAllUsersByDomainCtx deletes all users associated with the domain with context support.

#### func (*ContextEnforcer) [DeleteDomainsCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_context.go#L93) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteDomainsCtx(ctx [context](/context).[Context](/context#Context), domains ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteDomainsCtx deletes all associated policies for domains with context support. It would delete all domains if parameter is not provided.

#### func (*ContextEnforcer) [DeletePermissionCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L97) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeletePermissionCtx(ctx [context](/context).[Context](/context#Context), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermissionCtx deletes a permission with context support. Returns false if the permission does not exist (aka not affected).

#### func (*ContextEnforcer) [DeletePermissionForUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L119) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeletePermissionForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermissionForUserCtx deletes a permission for a user or role with context support. Returns false if the user or role does not have the permission (aka not affected).

#### func (*ContextEnforcer) [DeletePermissionsForUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L125) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeletePermissionsForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermissionsForUserCtx deletes permissions for a user or role with context support. Returns false if the user or role does not have any permissions (aka not affected).

#### func (*ContextEnforcer) [DeleteRoleCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L75) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteRoleCtx(ctx [context](/context).[Context](/context#Context), role [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRoleCtx deletes a role with context support. Returns false if the role does not exist (aka not affected).

#### func (*ContextEnforcer) [DeleteRoleForUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L36) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteRoleForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRoleForUserCtx deletes a role for a user with context support. Returns false if the user does not have the role (aka not affected).

#### func (*ContextEnforcer) [DeleteRoleForUserInDomainCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_context.go#L32) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteRoleForUserInDomainCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRoleForUserInDomainCtx deletes a role for a user inside a domain with context support. Returns false if the user does not have the role (aka not affected).

#### func (*ContextEnforcer) [DeleteRolesForUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L44) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteRolesForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRolesForUserCtx deletes all roles for a user with context support. Returns false if the user does not have any roles (aka not affected).

#### func (*ContextEnforcer) [DeleteRolesForUserInDomainCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_context.go#L38) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteRolesForUserInDomainCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRolesForUserInDomainCtx deletes all roles for a user inside a domain with context support. Returns false if the user does not have any roles (aka not affected).

#### func (*ContextEnforcer) [DeleteUserCtx](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_context.go#L58) ¶ added in v2.128.0

    func (e *ContextEnforcer) DeleteUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteUserCtx deletes a user with context support. Returns false if the user does not exist (aka not affected).

#### func (*ContextEnforcer) [IsFilteredCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L134) ¶ added in v2.128.0

    func (e *ContextEnforcer) IsFilteredCtx(ctx [context](/context).[Context](/context#Context)) [bool](/builtin#bool)

IsFilteredCtx returns true if the loaded policy has been filtered with context.

#### func (*ContextEnforcer) [LoadPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L56) ¶ added in v2.128.0

    func (e *ContextEnforcer) LoadPolicyCtx(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)

LoadPolicyCtx loads all policy rules from the storage with context.

#### func (*ContextEnforcer) [RemoveFilteredGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L342) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveFilteredGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredGroupingPolicyCtx removes grouping policy rules that match the filter from the storage with context.

#### func (*ContextEnforcer) [RemoveFilteredNamedGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L347) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveFilteredNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredNamedGroupingPolicyCtx removes named grouping policy rules that match the filter from the storage with context.

#### func (*ContextEnforcer) [RemoveFilteredNamedPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L232) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveFilteredNamedPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredNamedPolicyCtx removes named policy rules that match the filter from the storage with context.

#### func (*ContextEnforcer) [RemoveFilteredPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L227) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredPolicyCtx removes policy rules that match the filter from the storage with context.

#### func (*ContextEnforcer) [RemoveGroupingPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L332) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveGroupingPoliciesCtx removes grouping policy rules from the storage with context.

#### func (*ContextEnforcer) [RemoveGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L309) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveGroupingPolicyCtx removes a grouping policy rule from the storage with context.

#### func (*ContextEnforcer) [RemoveNamedGroupingPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L337) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveNamedGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedGroupingPoliciesCtx removes named grouping policy rules from the storage with context.

#### func (*ContextEnforcer) [RemoveNamedGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L314) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedGroupingPolicyCtx removes a named grouping policy rule from the storage with context.

#### func (*ContextEnforcer) [RemoveNamedPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L222) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveNamedPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPoliciesCtx removes named policy rules from the storage with context.

#### func (*ContextEnforcer) [RemoveNamedPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L204) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemoveNamedPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPolicyCtx removes a named policy rule from the storage with context.

#### func (*ContextEnforcer) [RemovePoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L217) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemovePoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemovePoliciesCtx removes policy rules from the storage with context.

#### func (*ContextEnforcer) [RemovePolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L199) ¶ added in v2.128.0

    func (e *ContextEnforcer) RemovePolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemovePolicyCtx removes a policy rule from the storage with context.

#### func (*ContextEnforcer) [SavePolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L142) ¶ added in v2.128.0

    func (e *ContextEnforcer) SavePolicyCtx(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)

#### func (*ContextEnforcer) [SelfAddPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L379) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfAddPoliciesCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

SelfAddPoliciesCtx adds policy rules to the current policy with context.

#### func (*ContextEnforcer) [SelfAddPoliciesExCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L383) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfAddPoliciesExCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*ContextEnforcer) [SelfAddPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L374) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfAddPolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

SelfAddPolicyCtx adds a policy rule to the current policy with context.

#### func (*ContextEnforcer) [SelfRemoveFilteredPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L398) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfRemoveFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

SelfRemoveFilteredPolicyCtx removes policy rules that match the filter from the current policy with context.

#### func (*ContextEnforcer) [SelfRemovePoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L393) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfRemovePoliciesCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

SelfRemovePoliciesCtx removes policy rules from the current policy with context.

#### func (*ContextEnforcer) [SelfRemovePolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L388) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfRemovePolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

SelfRemovePolicyCtx removes a policy rule from the current policy with context.

#### func (*ContextEnforcer) [SelfUpdatePoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L408) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfUpdatePoliciesCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

SelfUpdatePoliciesCtx updates policy rules in the current policy with context.

#### func (*ContextEnforcer) [SelfUpdatePolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L403) ¶ added in v2.128.0

    func (e *ContextEnforcer) SelfUpdatePolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

SelfUpdatePolicyCtx updates a policy rule in the current policy with context.

#### func (*ContextEnforcer) [UpdateFilteredNamedPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L262) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateFilteredNamedPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateFilteredNamedPoliciesCtx updates named policy rules that match the filter in the storage with context.

#### func (*ContextEnforcer) [UpdateFilteredPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L257) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateFilteredPoliciesCtx(ctx [context](/context).[Context](/context#Context), newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateFilteredPoliciesCtx updates policy rules that match the filter in the storage with context.

#### func (*ContextEnforcer) [UpdateGroupingPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L362) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateGroupingPoliciesCtx updates grouping policy rules in the storage with context.

#### func (*ContextEnforcer) [UpdateGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L352) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateGroupingPolicyCtx updates a grouping policy rule in the storage with context.

#### func (*ContextEnforcer) [UpdateNamedGroupingPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L367) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateNamedGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateNamedGroupingPoliciesCtx updates named grouping policy rules in the storage with context.

#### func (*ContextEnforcer) [UpdateNamedGroupingPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L357) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateNamedGroupingPolicyCtx updates a named grouping policy rule in the storage with context.

#### func (*ContextEnforcer) [UpdateNamedPoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L252) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateNamedPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), p1 [][][string](/builtin#string), p2 [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateNamedPoliciesCtx updates named policy rules in the storage with context.

#### func (*ContextEnforcer) [UpdateNamedPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L242) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdateNamedPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), p1 [][string](/builtin#string), p2 [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateNamedPolicyCtx updates a named policy rule in the storage with context.

#### func (*ContextEnforcer) [UpdatePoliciesCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L247) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdatePoliciesCtx(ctx [context](/context).[Context](/context#Context), oldPolicies [][][string](/builtin#string), newPolicies [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdatePoliciesCtx updates policy rules in the storage with context.

#### func (*ContextEnforcer) [UpdatePolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L237) ¶ added in v2.128.0

    func (e *ContextEnforcer) UpdatePolicyCtx(ctx [context](/context).[Context](/context#Context), oldPolicy [][string](/builtin#string), newPolicy [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdatePolicyCtx updates a policy rule in the storage with context.

#### type [DistributedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L9) ¶ added in v2.19.0

    type DistributedEnforcer struct {
     *SyncedEnforcer
    }

DistributedEnforcer wraps SyncedEnforcer for dispatcher.

#### func [NewDistributedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L13) ¶ added in v2.19.0

    func NewDistributedEnforcer(params ...interface{}) (*DistributedEnforcer, [error](/builtin#error))

#### func (*DistributedEnforcer) [AddPoliciesSelf](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L31) ¶ added in v2.23.3

    func (d *DistributedEnforcer) AddPoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) (affected [][][string](/builtin#string), err [error](/builtin#error))

AddPoliciesSelf provides a method for dispatcher to add authorization rules to the current policy. The function returns the rules affected and error.

#### func (*DistributedEnforcer) [ClearPolicySelf](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L124) ¶ added in v2.19.0

    func (d *DistributedEnforcer) ClearPolicySelf(shouldPersist func() [bool](/builtin#bool)) [error](/builtin#error)

ClearPolicySelf provides a method for dispatcher to clear all rules from the current policy.

#### func (*DistributedEnforcer) [RemoveFilteredPolicySelf](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L97) ¶ added in v2.19.0

    func (d *DistributedEnforcer) RemoveFilteredPolicySelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) (affected [][][string](/builtin#string), err [error](/builtin#error))

RemoveFilteredPolicySelf provides a method for dispatcher to remove an authorization rule from the current policy, field filters can be specified. The function returns the rules affected and error.

#### func (*DistributedEnforcer) [RemovePoliciesSelf](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L69) ¶ added in v2.23.3

    func (d *DistributedEnforcer) RemovePoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) (affected [][][string](/builtin#string), err [error](/builtin#error))

RemovePoliciesSelf provides a method for dispatcher to remove a set of rules from current policy. The function returns the rules affected and error.

#### func (*DistributedEnforcer) [SetDispatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L25) ¶ added in v2.20.1

    func (d *DistributedEnforcer) SetDispatcher(dispatcher [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Dispatcher](/github.com/casbin/casbin/v2@v2.135.0/persist#Dispatcher))

SetDispatcher sets the current dispatcher.

#### func (*DistributedEnforcer) [UpdateFilteredPoliciesSelf](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L200) ¶ added in v2.28.0

    func (d *DistributedEnforcer) UpdateFilteredPoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), newRules [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateFilteredPoliciesSelf provides a method for dispatcher to update a set of authorization rules from the current policy.

#### func (*DistributedEnforcer) [UpdatePoliciesSelf](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L170) ¶ added in v2.23.3

    func (d *DistributedEnforcer) UpdatePoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) (affected [bool](/builtin#bool), err [error](/builtin#error))

UpdatePoliciesSelf provides a method for dispatcher to update a set of authorization rules from the current policy.

#### func (*DistributedEnforcer) [UpdatePolicySelf](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_distributed.go#L140) ¶ added in v2.19.0

    func (d *DistributedEnforcer) UpdatePolicySelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newRule [][string](/builtin#string)) (affected [bool](/builtin#bool), err [error](/builtin#error))

UpdatePolicySelf provides a method for dispatcher to update an authorization rule from the current policy.

#### type [EnforceContext](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L61) ¶ added in v2.36.0

    type EnforceContext struct {
     RType [string](/builtin#string)
     PType [string](/builtin#string)
     EType [string](/builtin#string)
     MType [string](/builtin#string)
    }

EnforceContext is used as the first element of the parameter "rvals" in method "enforce".

#### func [NewEnforceContext](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L598) ¶ added in v2.36.0

    func NewEnforceContext(suffix [string](/builtin#string)) EnforceContext

NewEnforceContext Create a default structure based on the suffix.

#### func (EnforceContext) [GetCacheKey](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L68) ¶ added in v2.71.1

    func (e EnforceContext) GetCacheKey() [string](/builtin#string)

#### type [Enforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L37) ¶

    type Enforcer struct {
     // contains filtered or unexported fields
    }

Enforcer is the main interface for authorization enforcement and policy management.

#### func [NewEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L82) ¶

    func NewEnforcer(params ...interface{}) (*Enforcer, [error](/builtin#error))

NewEnforcer creates an enforcer via file or DB.

File:

    e := casbin.NewEnforcer("path/to/basic_model.conf", "path/to/basic_policy.csv")
    

MySQL DB:

    a := mysqladapter.NewDBAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
    e := casbin.NewEnforcer("path/to/basic_model.conf", a)
    

#### func (*Enforcer) [AddFunction](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L466) ¶

    func (e *Enforcer) AddFunction(name [string](/builtin#string), function [govaluate](/github.com/casbin/govaluate).[ExpressionFunction](/github.com/casbin/govaluate#ExpressionFunction))

AddFunction adds a customized function.

#### func (*Enforcer) [AddGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L360) ¶ added in v2.2.2

    func (e *Enforcer) AddGroupingPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPolicies adds role inheritance rules to the current policy. If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added. Otherwise the function returns true for the corresponding policy rule by adding the new rule.

#### func (*Enforcer) [AddGroupingPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L367) ¶ added in v2.63.0

    func (e *Enforcer) AddGroupingPoliciesEx(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPoliciesEx adds role inheritance rules to the current policy. If the rule already exists, the rule will not be added. But unlike AddGroupingPolicies, other non-existent rules are added instead of returning false directly.

#### func (*Enforcer) [AddGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L353) ¶

    func (e *Enforcer) AddGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPolicy adds a role inheritance rule to the current policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*Enforcer) [AddNamedDomainLinkConditionFunc](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L937) ¶ added in v2.75.0

    func (e *Enforcer) AddNamedDomainLinkConditionFunc(ptype, user, role [string](/builtin#string), domain [string](/builtin#string), fn [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[LinkConditionFunc](/github.com/casbin/casbin/v2@v2.135.0/rbac#LinkConditionFunc)) [bool](/builtin#bool)

AddNamedDomainLinkConditionFunc Add condition function fn for Link userName-> {roleName, domain}, when fn returns true, Link is valid, otherwise invalid.

#### func (*Enforcer) [AddNamedDomainMatchingFunc](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L913) ¶ added in v2.21.0

    func (e *Enforcer) AddNamedDomainMatchingFunc(ptype, name [string](/builtin#string), fn [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[MatchingFunc](/github.com/casbin/casbin/v2@v2.135.0/rbac#MatchingFunc)) [bool](/builtin#bool)

AddNamedDomainMatchingFunc add MatchingFunc by ptype to RoleManager.

#### func (*Enforcer) [AddNamedGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L394) ¶ added in v2.2.2

    func (e *Enforcer) AddNamedGroupingPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPolicies adds named role inheritance rules to the current policy. If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added. Otherwise the function returns true for the corresponding policy rule by adding the new rule.

#### func (*Enforcer) [AddNamedGroupingPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L401) ¶ added in v2.63.0

    func (e *Enforcer) AddNamedGroupingPoliciesEx(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPoliciesEx adds named role inheritance rules to the current policy. If the rule already exists, the rule will not be added. But unlike AddNamedGroupingPolicies, other non-existent rules are added instead of returning false directly.

#### func (*Enforcer) [AddNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L374) ¶

    func (e *Enforcer) AddNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPolicy adds a named role inheritance rule to the current policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*Enforcer) [AddNamedLinkConditionFunc](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L927) ¶ added in v2.75.0

    func (e *Enforcer) AddNamedLinkConditionFunc(ptype, user, role [string](/builtin#string), fn [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[LinkConditionFunc](/github.com/casbin/casbin/v2@v2.135.0/rbac#LinkConditionFunc)) [bool](/builtin#bool)

AddNamedLinkConditionFunc Add condition function fn for Link userName->roleName, when fn returns true, Link is valid, otherwise invalid.

#### func (*Enforcer) [AddNamedMatchingFunc](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L904) ¶ added in v2.21.0

    func (e *Enforcer) AddNamedMatchingFunc(ptype, name [string](/builtin#string), fn [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[MatchingFunc](/github.com/casbin/casbin/v2@v2.135.0/rbac#MatchingFunc)) [bool](/builtin#bool)

AddNamedMatchingFunc add MatchingFunc by ptype RoleManager.

#### func (*Enforcer) [AddNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L256) ¶ added in v2.2.2

    func (e *Enforcer) AddNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPolicies adds authorization rules to the current named policy. If the rule already exists, the function returns false for the corresponding rule and the rule will not be added. Otherwise the function returns true for the corresponding by adding the new rule.

#### func (*Enforcer) [AddNamedPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L263) ¶ added in v2.63.0

    func (e *Enforcer) AddNamedPoliciesEx(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPoliciesEx adds authorization rules to the current named policy. If the rule already exists, the rule will not be added. But unlike AddNamedPolicies, other non-existent rules are added instead of returning false directly.

#### func (*Enforcer) [AddNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L240) ¶

    func (e *Enforcer) AddNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPolicy adds an authorization rule to the current named policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*Enforcer) [AddPermissionForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L154) ¶

    func (e *Enforcer) AddPermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPermissionForUser adds a permission for a user or role. Returns false if the user or role already has the permission (aka not affected).

#### func (*Enforcer) [AddPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L160) ¶ added in v2.38.0

    func (e *Enforcer) AddPermissionsForUser(user [string](/builtin#string), permissions ...[][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPermissionsForUser adds multiple permissions for a user or role. Returns false if the user or role already has one of the permissions (aka not affected).

#### func (*Enforcer) [AddPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L226) ¶ added in v2.2.2

    func (e *Enforcer) AddPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPolicies adds authorization rules to the current policy. If the rule already exists, the function returns false for the corresponding rule and the rule will not be added. Otherwise the function returns true for the corresponding rule by adding the new rule.

#### func (*Enforcer) [AddPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L233) ¶ added in v2.63.0

    func (e *Enforcer) AddPoliciesEx(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPoliciesEx adds authorization rules to the current policy. If the rule already exists, the rule will not be added. But unlike AddPolicies, other non-existent rules are added instead of returning false directly.

#### func (*Enforcer) [AddPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L219) ¶

    func (e *Enforcer) AddPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddPolicy adds an authorization rule to the current policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*Enforcer) [AddRoleForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L67) ¶

    func (e *Enforcer) AddRoleForUser(user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRoleForUser adds a role for a user. Returns false if the user already has the role (aka not affected).

#### func (*Enforcer) [AddRoleForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L49) ¶

    func (e *Enforcer) AddRoleForUserInDomain(user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRoleForUserInDomain adds a role for a user inside a domain. Returns false if the user already has the role (aka not affected).

#### func (*Enforcer) [AddRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L75) ¶ added in v2.5.0

    func (e *Enforcer) AddRolesForUser(user [string](/builtin#string), roles [][string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRolesForUser adds roles for a user. Returns false if the user already has the roles (aka not affected).

#### func (*Enforcer) [BatchEnforce](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L878) ¶ added in v2.23.0

    func (e *Enforcer) BatchEnforce(requests [][]interface{}) ([][bool](/builtin#bool), [error](/builtin#error))

BatchEnforce enforce in batches.

#### func (*Enforcer) [BatchEnforceWithMatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L891) ¶ added in v2.23.0

    func (e *Enforcer) BatchEnforceWithMatcher(matcher [string](/builtin#string), requests [][]interface{}) ([][bool](/builtin#bool), [error](/builtin#error))

BatchEnforceWithMatcher enforce with matcher in batches.

#### func (*Enforcer) [BuildIncrementalConditionalRoleLinks](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L592) ¶ added in v2.75.0

    func (e *Enforcer) BuildIncrementalConditionalRoleLinks(op [model](/github.com/casbin/casbin/v2@v2.135.0/model).[PolicyOp](/github.com/casbin/casbin/v2@v2.135.0/model#PolicyOp), ptype [string](/builtin#string), rules [][][string](/builtin#string)) [error](/builtin#error)

BuildIncrementalConditionalRoleLinks provides incremental build the role inheritance relations with conditions.

#### func (*Enforcer) [BuildIncrementalRoleLinks](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L586) ¶ added in v2.6.0

    func (e *Enforcer) BuildIncrementalRoleLinks(op [model](/github.com/casbin/casbin/v2@v2.135.0/model).[PolicyOp](/github.com/casbin/casbin/v2@v2.135.0/model#PolicyOp), ptype [string](/builtin#string), rules [][][string](/builtin#string)) [error](/builtin#error)

BuildIncrementalRoleLinks provides incremental build the role inheritance relations.

#### func (*Enforcer) [BuildRoleLinks](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L570) ¶

    func (e *Enforcer) BuildRoleLinks() [error](/builtin#error)

BuildRoleLinks manually rebuild the role inheritance relations.

#### func (*Enforcer) [ClearPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L318) ¶

    func (e *Enforcer) ClearPolicy()

ClearPolicy clears all policy.

#### func (*Enforcer) [DeleteAllUsersByDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L112) ¶ added in v2.29.0

    func (e *Enforcer) DeleteAllUsersByDomain(domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteAllUsersByDomain would delete all users associated with the domain.

#### func (*Enforcer) [DeleteDomains](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L149) ¶ added in v2.29.0

    func (e *Enforcer) DeleteDomains(domains ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteDomains would delete all associated policies. It would delete all domains if parameter is not provided.

#### func (*Enforcer) [DeletePermission](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L148) ¶

    func (e *Enforcer) DeletePermission(permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermission deletes a permission. Returns false if the permission does not exist (aka not affected).

#### func (*Enforcer) [DeletePermissionForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L170) ¶

    func (e *Enforcer) DeletePermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermissionForUser deletes a permission for a user or role. Returns false if the user or role does not have the permission (aka not affected).

#### func (*Enforcer) [DeletePermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L176) ¶

    func (e *Enforcer) DeletePermissionsForUser(user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermissionsForUser deletes permissions for a user or role. Returns false if the user or role does not have any permissions (aka not affected).

#### func (*Enforcer) [DeleteRole](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L126) ¶

    func (e *Enforcer) DeleteRole(role [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRole deletes a role. Returns false if the role does not exist (aka not affected).

#### func (*Enforcer) [DeleteRoleForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L87) ¶

    func (e *Enforcer) DeleteRoleForUser(user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRoleForUser deletes a role for a user. Returns false if the user does not have the role (aka not affected).

#### func (*Enforcer) [DeleteRoleForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L55) ¶

    func (e *Enforcer) DeleteRoleForUserInDomain(user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRoleForUserInDomain deletes a role for a user inside a domain. Returns false if the user does not have the role (aka not affected).

#### func (*Enforcer) [DeleteRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L95) ¶

    func (e *Enforcer) DeleteRolesForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRolesForUser deletes all roles for a user. Returns false if the user does not have any roles (aka not affected).

#### func (*Enforcer) [DeleteRolesForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L61) ¶ added in v2.8.4

    func (e *Enforcer) DeleteRolesForUserInDomain(user [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRolesForUserInDomain deletes all roles for a user inside a domain. Returns false if the user does not have any roles (aka not affected).

#### func (*Enforcer) [DeleteUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L109) ¶

    func (e *Enforcer) DeleteUser(user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteUser deletes a user. Returns false if the user does not exist (aka not affected).

#### func (*Enforcer) [EnableAcceptJsonRequest](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L565) ¶ added in v2.72.0

    func (e *Enforcer) EnableAcceptJsonRequest(acceptJsonRequest [bool](/builtin#bool))

EnableAcceptJsonRequest controls whether to accept json as a request parameter.

#### func (*Enforcer) [EnableAutoBuildRoleLinks](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L560) ¶

    func (e *Enforcer) EnableAutoBuildRoleLinks(autoBuildRoleLinks [bool](/builtin#bool))

EnableAutoBuildRoleLinks controls whether to rebuild the role inheritance relations when a role is added or deleted.

#### func (*Enforcer) [EnableAutoNotifyDispatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L550) ¶ added in v2.18.0

    func (e *Enforcer) EnableAutoNotifyDispatcher(enable [bool](/builtin#bool))

EnableAutoNotifyDispatcher controls whether to save a policy rule automatically notify the Dispatcher when it is added or removed.

#### func (*Enforcer) [EnableAutoNotifyWatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L545) ¶ added in v2.2.1

    func (e *Enforcer) EnableAutoNotifyWatcher(enable [bool](/builtin#bool))

EnableAutoNotifyWatcher controls whether to save a policy rule automatically notify the Watcher when it is added or removed.

#### func (*Enforcer) [EnableAutoSave](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L555) ¶

    func (e *Enforcer) EnableAutoSave(autoSave [bool](/builtin#bool))

EnableAutoSave controls whether to save a policy rule automatically to the adapter when it is added or removed.

#### func (*Enforcer) [EnableEnforce](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L530) ¶

    func (e *Enforcer) EnableEnforce(enable [bool](/builtin#bool))

EnableEnforce changes the enforcing state of Casbin, when Casbin is disabled, all access will be allowed by the Enforce() function.

#### func (*Enforcer) [EnableLog](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L535) ¶

    func (e *Enforcer) EnableLog(enable [bool](/builtin#bool))

EnableLog changes whether Casbin will log messages to the Logger.

#### func (*Enforcer) [Enforce](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L854) ¶

    func (e *Enforcer) Enforce(rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).

#### func (*Enforcer) [EnforceEx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L864) ¶ added in v2.4.1

    func (e *Enforcer) EnforceEx(rvals ...interface{}) ([bool](/builtin#bool), [][string](/builtin#string), [error](/builtin#error))

EnforceEx explain enforcement by informing matched rules.

#### func (*Enforcer) [EnforceExWithMatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L871) ¶ added in v2.4.1

    func (e *Enforcer) EnforceExWithMatcher(matcher [string](/builtin#string), rvals ...interface{}) ([bool](/builtin#bool), [][string](/builtin#string), [error](/builtin#error))

EnforceExWithMatcher use a custom matcher and explain enforcement by informing matched rules.

#### func (*Enforcer) [EnforceWithMatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L859) ¶ added in v2.0.2

    func (e *Enforcer) EnforceWithMatcher(matcher [string](/builtin#string), rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

EnforceWithMatcher use a custom matcher to decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (matcher, sub, obj, act), use model matcher by default when matcher is "".

#### func (*Enforcer) [GetAdapter](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L257) ¶

    func (e *Enforcer) GetAdapter() [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter)

GetAdapter gets the current adapter.

#### func (*Enforcer) [GetAllActions](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L56) ¶

    func (e *Enforcer) GetAllActions() ([][string](/builtin#string), [error](/builtin#error))

GetAllActions gets the list of actions that show up in the current policy.

#### func (*Enforcer) [GetAllDomains](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L172) ¶ added in v2.43.0

    func (e *Enforcer) GetAllDomains() ([][string](/builtin#string), [error](/builtin#error))

GetAllDomains would get all domains.

#### func (*Enforcer) [GetAllNamedActions](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L61) ¶

    func (e *Enforcer) GetAllNamedActions(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedActions gets the list of actions that show up in the current named policy.

#### func (*Enforcer) [GetAllNamedObjects](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L47) ¶

    func (e *Enforcer) GetAllNamedObjects(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedObjects gets the list of objects that show up in the current named policy.

#### func (*Enforcer) [GetAllNamedRoles](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L75) ¶

    func (e *Enforcer) GetAllNamedRoles(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedRoles gets the list of roles that show up in the current named policy.

#### func (*Enforcer) [GetAllNamedSubjects](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L33) ¶

    func (e *Enforcer) GetAllNamedSubjects(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedSubjects gets the list of subjects that show up in the current named policy.

#### func (*Enforcer) [GetAllObjects](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L42) ¶

    func (e *Enforcer) GetAllObjects() ([][string](/builtin#string), [error](/builtin#error))

GetAllObjects gets the list of objects that show up in the current policy.

#### func (*Enforcer) [GetAllRoles](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L70) ¶

    func (e *Enforcer) GetAllRoles() ([][string](/builtin#string), [error](/builtin#error))

GetAllRoles gets the list of roles that show up in the current policy.

#### func (*Enforcer) [GetAllRolesByDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L181) ¶ added in v2.61.0

    func (e *Enforcer) GetAllRolesByDomain(domain [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllRolesByDomain would get all roles associated with the domain. note: Not applicable to Domains with inheritance relationship (implicit roles)

#### func (*Enforcer) [GetAllSubjects](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L28) ¶

    func (e *Enforcer) GetAllSubjects() ([][string](/builtin#string), [error](/builtin#error))

GetAllSubjects gets the list of subjects that show up in the current policy.

#### func (*Enforcer) [GetAllUsersByDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L79) ¶ added in v2.29.0

    func (e *Enforcer) GetAllUsersByDomain(domain [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllUsersByDomain would get all users associated with the domain.

#### func (*Enforcer) [GetAllowedObjectConditions](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L486) ¶ added in v2.68.0

    func (e *Enforcer) GetAllowedObjectConditions(user [string](/builtin#string), action [string](/builtin#string), prefix [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllowedObjectConditions returns a string array of object conditions that the user can access. For example: conditions, err := e.GetAllowedObjectConditions("alice", "read", "r.obj.") Note:

0. prefix: You can customize the prefix of the object conditions, and "r.obj." is commonly used as a prefix. After removing the prefix, the remaining part is the condition of the object. If there is an obj policy that does not meet the prefix requirement, an errors.ERR_OBJ_CONDITION will be returned.

1. If the 'objectConditions' array is empty, return errors.ERR_EMPTY_CONDITION This error is returned because some data adapters' ORM return full table data by default when they receive an empty condition, which tends to behave contrary to expectations.(e.g. GORM) If you are using an adapter that does not behave like this, you can choose to ignore this error.

#### func (*Enforcer) [GetDomainsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L409) ¶ added in v2.26.0

    func (e *Enforcer) GetDomainsForUser(user [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetDomainsForUser gets all domains.

#### func (*Enforcer) [GetFieldIndex](https://github.com/casbin/casbin/blob/v2.135.0/internal_api.go#L490) ¶ added in v2.48.0

    func (e *Enforcer) GetFieldIndex(ptype [string](/builtin#string), field [string](/builtin#string)) ([int](/builtin#int), [error](/builtin#error))

#### func (*Enforcer) [GetFilteredGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L105) ¶

    func (e *Enforcer) GetFilteredGroupingPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

#### func (*Enforcer) [GetFilteredNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L115) ¶

    func (e *Enforcer) GetFilteredNamedGroupingPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredNamedGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

#### func (*Enforcer) [GetFilteredNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L95) ¶

    func (e *Enforcer) GetFilteredNamedPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredNamedPolicy gets all the authorization rules in the named policy, field filters can be specified.

#### func (*Enforcer) [GetFilteredNamedPolicyWithMatcher](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L120) ¶ added in v2.47.0

    func (e *Enforcer) GetFilteredNamedPolicyWithMatcher(ptype [string](/builtin#string), matcher [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredNamedPolicyWithMatcher gets rules based on matcher from the policy.

#### func (*Enforcer) [GetFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L85) ¶

    func (e *Enforcer) GetFilteredPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.

#### func (*Enforcer) [GetGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L100) ¶

    func (e *Enforcer) GetGroupingPolicy() ([][][string](/builtin#string), [error](/builtin#error))

GetGroupingPolicy gets all the role inheritance rules in the policy.

#### func (*Enforcer) [GetImplicitObjectPatternsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L663) ¶ added in v2.121.0

    func (e *Enforcer) GetImplicitObjectPatternsForUser(user [string](/builtin#string), domain [string](/builtin#string), action [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetImplicitObjectPatternsForUser returns all object patterns (with wildcards) that a user has for a given domain and action. For example: p, admin, chronicle/123, location/*, read p, user, chronicle/456, location/789, read g, alice, admin g, bob, user

GetImplicitObjectPatternsForUser("alice", "chronicle/123", "read") will return ["location/*"]. GetImplicitObjectPatternsForUser("bob", "chronicle/456", "read") will return ["location/789"].

#### func (*Enforcer) [GetImplicitPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L307) ¶

    func (e *Enforcer) GetImplicitPermissionsForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetImplicitPermissionsForUser gets implicit permissions for a user or role. Compared to GetPermissionsForUser(), this function retrieves permissions for inherited roles. For example: p, admin, data1, read p, alice, data2, read g, alice, admin

GetPermissionsForUser("alice") can only get: [["alice", "data2", "read"]]. But GetImplicitPermissionsForUser("alice") will get: [["admin", "data1", "read"], ["alice", "data2", "read"]].

#### func (*Enforcer) [GetImplicitResourcesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L429) ¶ added in v2.31.0

    func (e *Enforcer) GetImplicitResourcesForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetImplicitResourcesForUser returns all policies that user obtaining in domain.

#### func (*Enforcer) [GetImplicitRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L233) ¶

    func (e *Enforcer) GetImplicitRolesForUser(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetImplicitRolesForUser gets implicit roles that a user has. Compared to GetRolesForUser(), this function retrieves indirect roles besides direct roles. For example: g, alice, role:admin g, role:admin, role:user

GetRolesForUser("alice") can only get: ["role:admin"]. But GetImplicitRolesForUser("alice") will get: ["role:admin", "role:user"].

#### func (*Enforcer) [GetImplicitUsersForPermission](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L373) ¶

    func (e *Enforcer) GetImplicitUsersForPermission(permission ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetImplicitUsersForPermission gets implicit users for a permission. For example: p, admin, data1, read p, bob, data1, read g, alice, admin

GetImplicitUsersForPermission("data1", "read") will get: ["alice", "bob"]. Note: only users will be returned, roles (2nd arg in "g") will be excluded.

#### func (*Enforcer) [GetImplicitUsersForResource](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L535) ¶ added in v2.69.0

    func (e *Enforcer) GetImplicitUsersForResource(resource [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetImplicitUsersForResource return implicit user based on resource. for example: p, alice, data1, read p, bob, data2, write p, data2_admin, data2, read p, data2_admin, data2, write g, alice, data2_admin GetImplicitUsersForResource("data2") will return [[bob data2 write] [alice data2 read] [alice data2 write]] GetImplicitUsersForResource("data1") will return [[alice data1 read]] Note: only users will be returned, roles (2nd arg in "g") will be excluded.

#### func (*Enforcer) [GetImplicitUsersForResourceByDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L603) ¶ added in v2.69.0

    func (e *Enforcer) GetImplicitUsersForResourceByDomain(resource [string](/builtin#string), domain [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetImplicitUsersForResourceByDomain return implicit user based on resource and domain. Compared to GetImplicitUsersForResource, domain is supported.

#### func (*Enforcer) [GetImplicitUsersForRole](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L275) ¶ added in v2.31.0

    func (e *Enforcer) GetImplicitUsersForRole(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetImplicitUsersForRole gets implicit users for a role.

#### func (*Enforcer) [GetModel](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L243) ¶

    func (e *Enforcer) GetModel() [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model)

GetModel gets the current model.

#### func (*Enforcer) [GetNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L110) ¶

    func (e *Enforcer) GetNamedGroupingPolicy(ptype [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedGroupingPolicy gets all the role inheritance rules in the policy.

#### func (*Enforcer) [GetNamedImplicitPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L320) ¶ added in v2.45.0

    func (e *Enforcer) GetNamedImplicitPermissionsForUser(ptype [string](/builtin#string), gtype [string](/builtin#string), user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedImplicitPermissionsForUser gets implicit permissions for a user or role by named policy. Compared to GetNamedPermissionsForUser(), this function retrieves permissions for inherited roles. For example: p, admin, data1, read p2, admin, create g, alice, admin

GetImplicitPermissionsForUser("alice") can only get: [["admin", "data1", "read"]], whose policy is default policy "p" But you can specify the named policy "p2" to get: [["admin", "create"]] by GetNamedImplicitPermissionsForUser("p2","alice").

#### func (*Enforcer) [GetNamedImplicitRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L264) ¶ added in v2.95.0

    func (e *Enforcer) GetNamedImplicitRolesForUser(ptype [string](/builtin#string), name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetNamedImplicitRolesForUser gets implicit roles that a user has by named role definition. Compared to GetImplicitRolesForUser(), this function retrieves indirect roles besides direct roles. For example: g, alice, role:admin g, role:admin, role:user g2, alice, role:admin2

GetImplicitRolesForUser("alice") can only get: ["role:admin", "role:user"]. But GetNamedImplicitRolesForUser("g2", "alice") will get: ["role:admin2"].

#### func (*Enforcer) [GetNamedImplicitUsersForResource](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L546) ¶ added in v2.120.0

    func (e *Enforcer) GetNamedImplicitUsersForResource(ptype [string](/builtin#string), resource [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedImplicitUsersForResource return implicit user based on resource with named policy support. This function handles resource role relationships through named policies (e.g., g2, g3, etc.). for example: p, admin_group, admin_data, * g, admin, admin_group g2, app, admin_data GetNamedImplicitUsersForResource("g2", "app") will return users who have access to admin_data through g2 relationship.

#### func (*Enforcer) [GetNamedPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L190) ¶ added in v2.45.0

    func (e *Enforcer) GetNamedPermissionsForUser(ptype [string](/builtin#string), user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedPermissionsForUser gets permissions for a user or role by named policy.

#### func (*Enforcer) [GetNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L90) ¶

    func (e *Enforcer) GetNamedPolicy(ptype [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedPolicy gets all the authorization rules in the named policy.

#### func (*Enforcer) [GetNamedRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L290) ¶ added in v2.52.0

    func (e *Enforcer) GetNamedRoleManager(ptype [string](/builtin#string)) [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager)

GetNamedRoleManager gets the role manager for the named policy.

#### func (*Enforcer) [GetPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L185) ¶

    func (e *Enforcer) GetPermissionsForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetPermissionsForUser gets permissions for a user or role.

#### func (*Enforcer) [GetPermissionsForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L42) ¶

    func (e *Enforcer) GetPermissionsForUserInDomain(user [string](/builtin#string), domain [string](/builtin#string)) [][][string](/builtin#string)

GetPermissionsForUserInDomain gets permissions for a user or role inside a domain.

#### func (*Enforcer) [GetPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L80) ¶

    func (e *Enforcer) GetPolicy() ([][][string](/builtin#string), [error](/builtin#error))

GetPolicy gets all the authorization rules in the policy.

#### func (*Enforcer) [GetRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L279) ¶ added in v2.1.0

    func (e *Enforcer) GetRoleManager() [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager)

GetRoleManager gets the current role manager.

#### func (*Enforcer) [GetRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L29) ¶

    func (e *Enforcer) GetRolesForUser(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetRolesForUser gets the roles that a user has.

#### func (*Enforcer) [GetRolesForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L33) ¶

    func (e *Enforcer) GetRolesForUserInDomain(name [string](/builtin#string), domain [string](/builtin#string)) [][string](/builtin#string)

GetRolesForUserInDomain gets the roles that a user has inside a domain.

#### func (*Enforcer) [GetUsersForRole](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L39) ¶

    func (e *Enforcer) GetUsersForRole(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetUsersForRole gets the users that has a role.

#### func (*Enforcer) [GetUsersForRoleInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains.go#L24) ¶

    func (e *Enforcer) GetUsersForRoleInDomain(name [string](/builtin#string), domain [string](/builtin#string)) [][string](/builtin#string)

GetUsersForRoleInDomain gets the users that has a role inside a domain. Add by Gordon.

#### func (*Enforcer) [HasGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L332) ¶

    func (e *Enforcer) HasGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasGroupingPolicy determines whether a role inheritance rule exists.

#### func (*Enforcer) [HasNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L337) ¶

    func (e *Enforcer) HasNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasNamedGroupingPolicy determines whether a named role inheritance rule exists.

#### func (*Enforcer) [HasNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L203) ¶

    func (e *Enforcer) HasNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasNamedPolicy determines whether a named authorization rule exists.

#### func (*Enforcer) [HasPermissionForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L221) ¶

    func (e *Enforcer) HasPermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

HasPermissionForUser determines whether a user has a permission.

#### func (*Enforcer) [HasPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L198) ¶

    func (e *Enforcer) HasPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasPolicy determines whether an authorization rule exists.

#### func (*Enforcer) [HasRoleForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api.go#L49) ¶

    func (e *Enforcer) HasRoleForUser(name [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

HasRoleForUser determines whether a user has a role.

#### func (*Enforcer) [InitWithAdapter](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L159) ¶

    func (e *Enforcer) InitWithAdapter(modelPath [string](/builtin#string), adapter [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter)) [error](/builtin#error)

InitWithAdapter initializes an enforcer with a database adapter.

#### func (*Enforcer) [InitWithFile](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L153) ¶

    func (e *Enforcer) InitWithFile(modelPath [string](/builtin#string), policyPath [string](/builtin#string)) [error](/builtin#error)

InitWithFile initializes an enforcer with a model file and a policy file.

#### func (*Enforcer) [InitWithModelAndAdapter](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L175) ¶

    func (e *Enforcer) InitWithModelAndAdapter(m [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model), adapter [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter)) [error](/builtin#error)

InitWithModelAndAdapter initializes an enforcer with a model and a database adapter.

#### func (*Enforcer) [IsFiltered](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L471) ¶

    func (e *Enforcer) IsFiltered() [bool](/builtin#bool)

IsFiltered returns true if the loaded policy has been filtered.

#### func (*Enforcer) [IsLogEnabled](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L540) ¶ added in v2.16.0

    func (e *Enforcer) IsLogEnabled() [bool](/builtin#bool)

IsLogEnabled returns the current logger's enabled status.

#### func (*Enforcer) [LoadFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L459) ¶

    func (e *Enforcer) LoadFilteredPolicy(filter interface{}) [error](/builtin#error)

LoadFilteredPolicy reloads a filtered policy from file/database.

#### func (*Enforcer) [LoadFilteredPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L88) ¶ added in v2.128.0

    func (e *Enforcer) LoadFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), filter interface{}) [error](/builtin#error)

LoadFilteredPolicyCtx loads all policy rules from the storage with context and filter.

#### func (*Enforcer) [LoadIncrementalFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L466) ¶ added in v2.11.0

    func (e *Enforcer) LoadIncrementalFilteredPolicy(filter interface{}) [error](/builtin#error)

LoadIncrementalFilteredPolicy append a filtered policy from file/database.

#### func (*Enforcer) [LoadIncrementalFilteredPolicyCtx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L94) ¶ added in v2.128.0

    func (e *Enforcer) LoadIncrementalFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), filter interface{}) [error](/builtin#error)

LoadIncrementalFilteredPolicyCtx append a filtered policy from file/database with context.

#### func (*Enforcer) [LoadModel](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L226) ¶

    func (e *Enforcer) LoadModel() [error](/builtin#error)

LoadModel reloads the model from the model CONF file. Because the policy is attached to a model, so the policy is invalidated and needs to be reloaded by calling LoadPolicy().

#### func (*Enforcer) [LoadPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L329) ¶

    func (e *Enforcer) LoadPolicy() [error](/builtin#error)

LoadPolicy reloads the policy from file/database.

#### func (*Enforcer) [RemoveFilteredGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L416) ¶

    func (e *Enforcer) RemoveFilteredGroupingPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredGroupingPolicy removes a role inheritance rule from the current policy, field filters can be specified.

#### func (*Enforcer) [RemoveFilteredNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L461) ¶

    func (e *Enforcer) RemoveFilteredNamedGroupingPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredNamedGroupingPolicy removes a role inheritance rule from the current named policy, field filters can be specified.

#### func (*Enforcer) [RemoveFilteredNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L327) ¶

    func (e *Enforcer) RemoveFilteredNamedPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredNamedPolicy removes an authorization rule from the current named policy, field filters can be specified.

#### func (*Enforcer) [RemoveFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L304) ¶

    func (e *Enforcer) RemoveFilteredPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredPolicy removes an authorization rule from the current policy, field filters can be specified.

#### func (*Enforcer) [RemoveGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L411) ¶ added in v2.2.2

    func (e *Enforcer) RemoveGroupingPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveGroupingPolicies removes role inheritance rules from the current policy.

#### func (*Enforcer) [RemoveGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L406) ¶

    func (e *Enforcer) RemoveGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveGroupingPolicy removes a role inheritance rule from the current policy.

#### func (*Enforcer) [RemoveNamedGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L439) ¶ added in v2.2.2

    func (e *Enforcer) RemoveNamedGroupingPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedGroupingPolicies removes role inheritance rules from the current named policy.

#### func (*Enforcer) [RemoveNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L421) ¶

    func (e *Enforcer) RemoveNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedGroupingPolicy removes a role inheritance rule from the current named policy.

#### func (*Enforcer) [RemoveNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L322) ¶ added in v2.2.2

    func (e *Enforcer) RemoveNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPolicies removes authorization rules from the current named policy.

#### func (*Enforcer) [RemoveNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L309) ¶

    func (e *Enforcer) RemoveNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPolicy removes an authorization rule from the current named policy.

#### func (*Enforcer) [RemovePolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L299) ¶ added in v2.2.2

    func (e *Enforcer) RemovePolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemovePolicies removes authorization rules from the current policy.

#### func (*Enforcer) [RemovePolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L268) ¶

    func (e *Enforcer) RemovePolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemovePolicy removes an authorization rule from the current policy.

#### func (*Enforcer) [SavePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L480) ¶

    func (e *Enforcer) SavePolicy() [error](/builtin#error)

SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.

#### func (*Enforcer) [SelfAddPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L474) ¶ added in v2.53.0

    func (e *Enforcer) SelfAddPolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SelfAddPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L478) ¶ added in v2.63.0

    func (e *Enforcer) SelfAddPoliciesEx(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SelfAddPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L470) ¶ added in v2.53.0

    func (e *Enforcer) SelfAddPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SelfRemoveFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L490) ¶ added in v2.53.0

    func (e *Enforcer) SelfRemoveFilteredPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SelfRemovePolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L486) ¶ added in v2.53.1

    func (e *Enforcer) SelfRemovePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SelfRemovePolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L482) ¶ added in v2.53.0

    func (e *Enforcer) SelfRemovePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SelfUpdatePolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L498) ¶ added in v2.53.0

    func (e *Enforcer) SelfUpdatePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SelfUpdatePolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L494) ¶ added in v2.53.0

    func (e *Enforcer) SelfUpdatePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [SetAdapter](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L262) ¶

    func (e *Enforcer) SetAdapter(adapter [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter))

SetAdapter sets the current adapter.

#### func (*Enforcer) [SetEffector](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L313) ¶

    func (e *Enforcer) SetEffector(eft [effector](/github.com/casbin/casbin/v2@v2.135.0/effector).[Effector](/github.com/casbin/casbin/v2@v2.135.0/effector#Effector))

SetEffector sets the current effector.

#### func (*Enforcer) [SetFieldIndex](https://github.com/casbin/casbin/blob/v2.135.0/internal_api.go#L494) ¶ added in v2.48.0

    func (e *Enforcer) SetFieldIndex(ptype [string](/builtin#string), field [string](/builtin#string), index [int](/builtin#int))

#### func (*Enforcer) [SetLogger](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L198) ¶ added in v2.16.0

    func (e *Enforcer) SetLogger(logger [log](/github.com/casbin/casbin/v2@v2.135.0/log).[Logger](/github.com/casbin/casbin/v2@v2.135.0/log#Logger))

SetLogger changes the current enforcer's logger.

#### func (*Enforcer) [SetModel](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L248) ¶

    func (e *Enforcer) SetModel(m [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model))

SetModel sets the current model.

#### func (*Enforcer) [SetNamedDomainLinkConditionFuncParams](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L956) ¶ added in v2.75.0

    func (e *Enforcer) SetNamedDomainLinkConditionFuncParams(ptype, user, role, domain [string](/builtin#string), params ...[string](/builtin#string)) [bool](/builtin#bool)

SetNamedDomainLinkConditionFuncParams Sets the parameters of the condition function fn for Link userName->{roleName, domain}.

#### func (*Enforcer) [SetNamedLinkConditionFuncParams](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L946) ¶ added in v2.75.0

    func (e *Enforcer) SetNamedLinkConditionFuncParams(ptype, user, role [string](/builtin#string), params ...[string](/builtin#string)) [bool](/builtin#bool)

SetNamedLinkConditionFuncParams Sets the parameters of the condition function fn for Link userName->roleName.

#### func (*Enforcer) [SetNamedRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L307) ¶ added in v2.52.0

    func (e *Enforcer) SetNamedRoleManager(ptype [string](/builtin#string), rm [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager))

SetNamedRoleManager sets the role manager for the named policy.

#### func (*Enforcer) [SetRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L301) ¶

    func (e *Enforcer) SetRoleManager(rm [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager))

SetRoleManager sets the current role manager.

#### func (*Enforcer) [SetWatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer.go#L267) ¶

    func (e *Enforcer) SetWatcher(watcher [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Watcher](/github.com/casbin/casbin/v2@v2.135.0/persist#Watcher)) [error](/builtin#error)

SetWatcher sets the current watcher.

#### func (*Enforcer) [UpdateFilteredNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L294) ¶ added in v2.28.0

    func (e *Enforcer) UpdateFilteredNamedPolicies(ptype [string](/builtin#string), newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [UpdateFilteredPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L290) ¶ added in v2.28.0

    func (e *Enforcer) UpdateFilteredPolicies(newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [UpdateGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L448) ¶ added in v2.41.0

    func (e *Enforcer) UpdateGroupingPolicies(oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateGroupingPolicies updates authorization rules from the current policies.

#### func (*Enforcer) [UpdateGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L443) ¶ added in v2.19.0

    func (e *Enforcer) UpdateGroupingPolicy(oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [UpdateNamedGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L456) ¶ added in v2.41.0

    func (e *Enforcer) UpdateNamedGroupingPolicies(ptype [string](/builtin#string), oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [UpdateNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L452) ¶ added in v2.19.0

    func (e *Enforcer) UpdateNamedGroupingPolicy(ptype [string](/builtin#string), oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [UpdateNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L286) ¶ added in v2.22.0

    func (e *Enforcer) UpdateNamedPolicies(ptype [string](/builtin#string), p1 [][][string](/builtin#string), p2 [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [UpdateNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L277) ¶ added in v2.14.0

    func (e *Enforcer) UpdateNamedPolicy(ptype [string](/builtin#string), p1 [][string](/builtin#string), p2 [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*Enforcer) [UpdatePolicies](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L282) ¶ added in v2.22.0

    func (e *Enforcer) UpdatePolicies(oldPolices [][][string](/builtin#string), newPolicies [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdatePolicies updates authorization rules from the current policies.

#### func (*Enforcer) [UpdatePolicy](https://github.com/casbin/casbin/blob/v2.135.0/management_api.go#L273) ¶ added in v2.14.0

    func (e *Enforcer) UpdatePolicy(oldPolicy [][string](/builtin#string), newPolicy [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdatePolicy updates an authorization rule from the current policy.

#### type [IDistributedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_interface.go#L166) ¶ added in v2.19.0

    type IDistributedEnforcer interface {
     IEnforcer
     SetDispatcher(dispatcher [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Dispatcher](/github.com/casbin/casbin/v2@v2.135.0/persist#Dispatcher))
     /* Management API for DistributedEnforcer*/
     AddPoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) (affected [][][string](/builtin#string), err [error](/builtin#error))
     RemovePoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) (affected [][][string](/builtin#string), err [error](/builtin#error))
     RemoveFilteredPolicySelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) (affected [][][string](/builtin#string), err [error](/builtin#error))
     ClearPolicySelf(shouldPersist func() [bool](/builtin#bool)) [error](/builtin#error)
     UpdatePolicySelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newRule [][string](/builtin#string)) (affected [bool](/builtin#bool), err [error](/builtin#error))
     UpdatePoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) (affected [bool](/builtin#bool), err [error](/builtin#error))
     UpdateFilteredPoliciesSelf(shouldPersist func() [bool](/builtin#bool), sec [string](/builtin#string), ptype [string](/builtin#string), newRules [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    }

IDistributedEnforcer defines dispatcher enforcer.

#### type [IEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_interface.go#L30) ¶ added in v2.1.1

    type IEnforcer interface {
     /* Enforcer API */
     InitWithFile(modelPath [string](/builtin#string), policyPath [string](/builtin#string)) [error](/builtin#error)
     InitWithAdapter(modelPath [string](/builtin#string), adapter [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter)) [error](/builtin#error)
     InitWithModelAndAdapter(m [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model), adapter [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter)) [error](/builtin#error)
     LoadModel() [error](/builtin#error)
     GetModel() [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model)
     SetModel(m [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model))
     GetAdapter() [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter)
     SetAdapter(adapter [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Adapter](/github.com/casbin/casbin/v2@v2.135.0/persist#Adapter))
     SetWatcher(watcher [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Watcher](/github.com/casbin/casbin/v2@v2.135.0/persist#Watcher)) [error](/builtin#error)
     GetRoleManager() [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager)
     SetRoleManager(rm [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager))
     SetEffector(eft [effector](/github.com/casbin/casbin/v2@v2.135.0/effector).[Effector](/github.com/casbin/casbin/v2@v2.135.0/effector#Effector))
     ClearPolicy()
     LoadPolicy() [error](/builtin#error)
     LoadFilteredPolicy(filter interface{}) [error](/builtin#error)
     LoadIncrementalFilteredPolicy(filter interface{}) [error](/builtin#error)
     IsFiltered() [bool](/builtin#bool)
     SavePolicy() [error](/builtin#error)
     EnableEnforce(enable [bool](/builtin#bool))
     EnableLog(enable [bool](/builtin#bool))
     EnableAutoNotifyWatcher(enable [bool](/builtin#bool))
     EnableAutoSave(autoSave [bool](/builtin#bool))
     EnableAutoBuildRoleLinks(autoBuildRoleLinks [bool](/builtin#bool))
     BuildRoleLinks() [error](/builtin#error)
     Enforce(rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     EnforceWithMatcher(matcher [string](/builtin#string), rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     EnforceEx(rvals ...interface{}) ([bool](/builtin#bool), [][string](/builtin#string), [error](/builtin#error))
     EnforceExWithMatcher(matcher [string](/builtin#string), rvals ...interface{}) ([bool](/builtin#bool), [][string](/builtin#string), [error](/builtin#error))
     BatchEnforce(requests [][]interface{}) ([][bool](/builtin#bool), [error](/builtin#error))
     BatchEnforceWithMatcher(matcher [string](/builtin#string), requests [][]interface{}) ([][bool](/builtin#bool), [error](/builtin#error))
    
     /* RBAC API */
     GetRolesForUser(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     GetUsersForRole(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     HasRoleForUser(name [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddRoleForUser(user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddPermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddPermissionsForUser(user [string](/builtin#string), permissions ...[][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeletePermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeletePermissionsForUser(user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     GetPermissionsForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     HasPermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     GetImplicitRolesForUser(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     GetImplicitPermissionsForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     GetImplicitUsersForPermission(permission ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     DeleteRoleForUser(user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteRolesForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteUser(user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteRole(role [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeletePermission(permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     /* RBAC API with domains*/
     GetUsersForRoleInDomain(name [string](/builtin#string), domain [string](/builtin#string)) [][string](/builtin#string)
     GetRolesForUserInDomain(name [string](/builtin#string), domain [string](/builtin#string)) [][string](/builtin#string)
     GetPermissionsForUserInDomain(user [string](/builtin#string), domain [string](/builtin#string)) [][][string](/builtin#string)
     AddRoleForUserInDomain(user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteRoleForUserInDomain(user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     GetAllUsersByDomain(domain [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     DeleteRolesForUserInDomain(user [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteAllUsersByDomain(domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteDomains(domains ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     GetAllDomains() ([][string](/builtin#string), [error](/builtin#error))
     GetAllRolesByDomain(domain [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
    
     /* Management API */
     GetAllSubjects() ([][string](/builtin#string), [error](/builtin#error))
     GetAllNamedSubjects(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     GetAllObjects() ([][string](/builtin#string), [error](/builtin#error))
     GetAllNamedObjects(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     GetAllActions() ([][string](/builtin#string), [error](/builtin#error))
     GetAllNamedActions(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     GetAllRoles() ([][string](/builtin#string), [error](/builtin#error))
     GetAllNamedRoles(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))
     GetPolicy() ([][][string](/builtin#string), [error](/builtin#error))
     GetFilteredPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     GetNamedPolicy(ptype [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     GetFilteredNamedPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     GetGroupingPolicy() ([][][string](/builtin#string), [error](/builtin#error))
     GetFilteredGroupingPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     GetNamedGroupingPolicy(ptype [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     GetFilteredNamedGroupingPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))
     HasPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     HasNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddPoliciesEx(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedPoliciesEx(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemovePolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemovePolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredNamedPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     HasGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     HasNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddGroupingPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddGroupingPoliciesEx(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedGroupingPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedGroupingPoliciesEx(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveGroupingPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredGroupingPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedGroupingPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredNamedGroupingPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddFunction(name [string](/builtin#string), function [govaluate](/github.com/casbin/govaluate).[ExpressionFunction](/github.com/casbin/govaluate#ExpressionFunction))
    
     UpdatePolicy(oldPolicy [][string](/builtin#string), newPolicy [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdatePolicies(oldPolicies [][][string](/builtin#string), newPolicies [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateFilteredPolicies(newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     UpdateGroupingPolicy(oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateGroupingPolicies(oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateNamedGroupingPolicy(ptype [string](/builtin#string), oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateNamedGroupingPolicies(ptype [string](/builtin#string), oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     /* Management API with autoNotifyWatcher disabled */
     SelfAddPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfAddPolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfAddPoliciesEx(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfRemovePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfRemovePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfRemoveFilteredPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfUpdatePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfUpdatePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    }

IEnforcer is the API interface of Enforcer.

#### type [IEnforcerContext](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context_interface.go#L19) ¶ added in v2.128.0

    type IEnforcerContext interface {
     IEnforcer
    
     /* Enforcer API */
     LoadPolicyCtx(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)
     LoadFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), filter interface{}) [error](/builtin#error)
     LoadIncrementalFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), filter interface{}) [error](/builtin#error)
     IsFilteredCtx(ctx [context](/context).[Context](/context#Context)) [bool](/builtin#bool)
     SavePolicyCtx(ctx [context](/context).[Context](/context#Context)) [error](/builtin#error)
    
     /* RBAC API */
     AddRoleForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddPermissionForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddPermissionsForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), permissions ...[][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeletePermissionForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeletePermissionsForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     DeleteRoleForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteRolesForUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteUserCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteRoleCtx(ctx [context](/context).[Context](/context#Context), role [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeletePermissionCtx(ctx [context](/context).[Context](/context#Context), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     /* RBAC API with domains*/
     AddRoleForUserInDomainCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteRoleForUserInDomainCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteRolesForUserInDomainCtx(ctx [context](/context).[Context](/context#Context), user [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteAllUsersByDomainCtx(ctx [context](/context).[Context](/context#Context), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     DeleteDomainsCtx(ctx [context](/context).[Context](/context#Context), domains ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     /* Management API */
     AddPolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddPoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddPoliciesExCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedPoliciesExCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     RemovePolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemovePoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredNamedPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     AddGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddGroupingPoliciesExCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     AddNamedGroupingPoliciesExCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     RemoveGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveNamedGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     RemoveFilteredNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     UpdatePolicyCtx(ctx [context](/context).[Context](/context#Context), oldPolicy [][string](/builtin#string), newPolicy [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdatePoliciesCtx(ctx [context](/context).[Context](/context#Context), oldPolicies [][][string](/builtin#string), newPolicies [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateFilteredPoliciesCtx(ctx [context](/context).[Context](/context#Context), newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     UpdateGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateNamedGroupingPolicyCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     UpdateNamedGroupingPoliciesCtx(ctx [context](/context).[Context](/context#Context), ptype [string](/builtin#string), oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    
     /* Management API with autoNotifyWatcher disabled */
     SelfAddPolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfAddPoliciesCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfAddPoliciesExCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfRemovePolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfRemovePoliciesCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfRemoveFilteredPolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfUpdatePolicyCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
     SelfUpdatePoliciesCtx(ctx [context](/context).[Context](/context#Context), sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))
    }

#### func [NewContextEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_context.go#L34) ¶ added in v2.128.0

    func NewContextEnforcer(params ...interface{}) (IEnforcerContext, [error](/builtin#error))

NewContextEnforcer creates a context-aware enforcer via file or DB.

#### type [SyncedCachedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L26) ¶ added in v2.66.0

    type SyncedCachedEnforcer struct {
     *SyncedEnforcer
     // contains filtered or unexported fields
    }

SyncedCachedEnforcer wraps Enforcer and provides decision sync cache.

#### func [NewSyncedCachedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L35) ¶ added in v2.66.0

    func NewSyncedCachedEnforcer(params ...interface{}) (*SyncedCachedEnforcer, [error](/builtin#error))

NewSyncedCachedEnforcer creates a sync cached enforcer via file or DB.

#### func (*SyncedCachedEnforcer) [AddPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L101) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) AddPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedCachedEnforcer) [AddPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L94) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) AddPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedCachedEnforcer) [EnableCache](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L50) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) EnableCache(enableCache [bool](/builtin#bool))

EnableCache determines whether to enable cache on Enforce(). When enableCache is enabled, cached result (true | false) will be returned for previous decisions.

#### func (*SyncedCachedEnforcer) [Enforce](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L60) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) Enforce(rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act). if rvals is not string , ignore the cache.

#### func (*SyncedCachedEnforcer) [InvalidateCache](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L148) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) InvalidateCache() [error](/builtin#error)

InvalidateCache deletes all the existing cached decisions.

#### func (*SyncedCachedEnforcer) [LoadPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L85) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) LoadPolicy() [error](/builtin#error)

#### func (*SyncedCachedEnforcer) [RemovePolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L115) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) RemovePolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedCachedEnforcer) [RemovePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L108) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) RemovePolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedCachedEnforcer) [SetCache](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L133) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) SetCache(c [cache](/github.com/casbin/casbin/v2@v2.135.0/persist/cache).[Cache](/github.com/casbin/casbin/v2@v2.135.0/persist/cache#Cache))

SetCache need to be sync cache.

#### func (*SyncedCachedEnforcer) [SetExpireTime](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_cached_synced.go#L126) ¶ added in v2.66.0

    func (e *SyncedCachedEnforcer) SetExpireTime(expireTime [time](/time).[Duration](/time#Duration))

#### type [SyncedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L29) ¶

    type SyncedEnforcer struct {
     *Enforcer
     // contains filtered or unexported fields
    }

SyncedEnforcer wraps Enforcer and provides synchronized access.

#### func [NewSyncedEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L37) ¶

    func NewSyncedEnforcer(params ...interface{}) (*SyncedEnforcer, [error](/builtin#error))

NewSyncedEnforcer creates a synchronized enforcer via file or DB.

#### func (*SyncedEnforcer) [AddFunction](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L627) ¶ added in v2.0.2

    func (e *SyncedEnforcer) AddFunction(name [string](/builtin#string), function [govaluate](/github.com/casbin/govaluate).[ExpressionFunction](/github.com/casbin/govaluate#ExpressionFunction))

AddFunction adds a customized function.

#### func (*SyncedEnforcer) [AddGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L518) ¶ added in v2.8.1

    func (e *SyncedEnforcer) AddGroupingPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPolicies adds role inheritance rulea to the current policy. If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added. Otherwise the function returns true for the corresponding policy rule by adding the new rule.

#### func (*SyncedEnforcer) [AddGroupingPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L527) ¶ added in v2.63.0

    func (e *SyncedEnforcer) AddGroupingPoliciesEx(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPoliciesEx adds role inheritance rules to the current policy. If the rule already exists, the rule will not be added. But unlike AddGroupingPolicies, other non-existent rules are added instead of returning false directly.

#### func (*SyncedEnforcer) [AddGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L509) ¶

    func (e *SyncedEnforcer) AddGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPolicy adds a role inheritance rule to the current policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*SyncedEnforcer) [AddNamedGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L545) ¶ added in v2.8.1

    func (e *SyncedEnforcer) AddNamedGroupingPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPolicies adds named role inheritance rules to the current policy. If the rule already exists, the function returns false for the corresponding policy rule and the rule will not be added. Otherwise the function returns true for the corresponding policy rule by adding the new rule.

#### func (*SyncedEnforcer) [AddNamedGroupingPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L554) ¶ added in v2.63.0

    func (e *SyncedEnforcer) AddNamedGroupingPoliciesEx(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPoliciesEx adds named role inheritance rules to the current policy. If the rule already exists, the rule will not be added. But unlike AddNamedGroupingPolicies, other non-existent rules are added instead of returning false directly.

#### func (*SyncedEnforcer) [AddNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L536) ¶ added in v2.0.2

    func (e *SyncedEnforcer) AddNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPolicy adds a named role inheritance rule to the current policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*SyncedEnforcer) [AddNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L397) ¶ added in v2.8.1

    func (e *SyncedEnforcer) AddNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPolicies adds authorization rules to the current named policy. If the rule already exists, the function returns false for the corresponding rule and the rule will not be added. Otherwise the function returns true for the corresponding by adding the new rule.

#### func (*SyncedEnforcer) [AddNamedPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L406) ¶ added in v2.63.0

    func (e *SyncedEnforcer) AddNamedPoliciesEx(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPoliciesEx adds authorization rules to the current named policy. If the rule already exists, the rule will not be added. But unlike AddNamedPolicies, other non-existent rules are added instead of returning false directly.

#### func (*SyncedEnforcer) [AddNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L388) ¶ added in v2.0.2

    func (e *SyncedEnforcer) AddNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPolicy adds an authorization rule to the current named policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*SyncedEnforcer) [AddPermissionForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L96) ¶

    func (e *SyncedEnforcer) AddPermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPermissionForUser adds a permission for a user or role. Returns false if the user or role already has the permission (aka not affected).

#### func (*SyncedEnforcer) [AddPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L104) ¶ added in v2.73.0

    func (e *SyncedEnforcer) AddPermissionsForUser(user [string](/builtin#string), permissions ...[][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPermissionsForUser adds permissions for a user or role. Returns false if the user or role already has the permissions (aka not affected).

#### func (*SyncedEnforcer) [AddPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L370) ¶ added in v2.8.1

    func (e *SyncedEnforcer) AddPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPolicies adds authorization rules to the current policy. If the rule already exists, the function returns false for the corresponding rule and the rule will not be added. Otherwise the function returns true for the corresponding rule by adding the new rule.

#### func (*SyncedEnforcer) [AddPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L379) ¶ added in v2.63.0

    func (e *SyncedEnforcer) AddPoliciesEx(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPoliciesEx adds authorization rules to the current policy. If the rule already exists, the rule will not be added. But unlike AddPolicies, other non-existent rules are added instead of returning false directly.

#### func (*SyncedEnforcer) [AddPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L361) ¶

    func (e *SyncedEnforcer) AddPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddPolicy adds an authorization rule to the current policy. If the rule already exists, the function returns false and the rule will not be added. Otherwise the function returns true by adding the new rule.

#### func (*SyncedEnforcer) [AddRoleForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L40) ¶

    func (e *SyncedEnforcer) AddRoleForUser(user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRoleForUser adds a role for a user. Returns false if the user already has the role (aka not affected).

#### func (*SyncedEnforcer) [AddRoleForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_synced.go#L40) ¶

    func (e *SyncedEnforcer) AddRoleForUserInDomain(user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRoleForUserInDomain adds a role for a user inside a domain. Returns false if the user already has the role (aka not affected).

#### func (*SyncedEnforcer) [AddRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L48) ¶ added in v2.25.1

    func (e *SyncedEnforcer) AddRolesForUser(user [string](/builtin#string), roles [][string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddRolesForUser adds roles for a user. Returns false if the user already has the roles (aka not affected).

#### func (*SyncedEnforcer) [BatchEnforce](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L219) ¶ added in v2.25.0

    func (e *SyncedEnforcer) BatchEnforce(requests [][]interface{}) ([][bool](/builtin#bool), [error](/builtin#error))

BatchEnforce enforce in batches.

#### func (*SyncedEnforcer) [BatchEnforceWithMatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L226) ¶ added in v2.25.0

    func (e *SyncedEnforcer) BatchEnforceWithMatcher(matcher [string](/builtin#string), requests [][]interface{}) ([][bool](/builtin#bool), [error](/builtin#error))

BatchEnforceWithMatcher enforce with matcher in batches.

#### func (*SyncedEnforcer) [BuildRoleLinks](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L184) ¶

    func (e *SyncedEnforcer) BuildRoleLinks() [error](/builtin#error)

BuildRoleLinks manually rebuild the role inheritance relations.

#### func (*SyncedEnforcer) [ClearPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L139) ¶

    func (e *SyncedEnforcer) ClearPolicy()

ClearPolicy clears all policy.

#### func (*SyncedEnforcer) [DeleteDomains](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_synced.go#L64) ¶ added in v2.109.0

    func (e *SyncedEnforcer) DeleteDomains(domains ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteDomains deletes domains from the model. Returns false if the domain does not exist (aka not affected).

#### func (*SyncedEnforcer) [DeletePermission](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L88) ¶

    func (e *SyncedEnforcer) DeletePermission(permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermission deletes a permission. Returns false if the permission does not exist (aka not affected).

#### func (*SyncedEnforcer) [DeletePermissionForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L112) ¶

    func (e *SyncedEnforcer) DeletePermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermissionForUser deletes a permission for a user or role. Returns false if the user or role does not have the permission (aka not affected).

#### func (*SyncedEnforcer) [DeletePermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L120) ¶

    func (e *SyncedEnforcer) DeletePermissionsForUser(user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeletePermissionsForUser deletes permissions for a user or role. Returns false if the user or role does not have any permissions (aka not affected).

#### func (*SyncedEnforcer) [DeleteRole](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L80) ¶

    func (e *SyncedEnforcer) DeleteRole(role [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRole deletes a role. Returns false if the role does not exist (aka not affected).

#### func (*SyncedEnforcer) [DeleteRoleForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L56) ¶

    func (e *SyncedEnforcer) DeleteRoleForUser(user [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRoleForUser deletes a role for a user. Returns false if the user does not have the role (aka not affected).

#### func (*SyncedEnforcer) [DeleteRoleForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_synced.go#L48) ¶

    func (e *SyncedEnforcer) DeleteRoleForUserInDomain(user [string](/builtin#string), role [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRoleForUserInDomain deletes a role for a user inside a domain. Returns false if the user does not have the role (aka not affected).

#### func (*SyncedEnforcer) [DeleteRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L64) ¶

    func (e *SyncedEnforcer) DeleteRolesForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRolesForUser deletes all roles for a user. Returns false if the user does not have any roles (aka not affected).

#### func (*SyncedEnforcer) [DeleteRolesForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_synced.go#L56) ¶ added in v2.8.4

    func (e *SyncedEnforcer) DeleteRolesForUserInDomain(user [string](/builtin#string), domain [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteRolesForUserInDomain deletes all roles for a user inside a domain. Returns false if the user does not have any roles (aka not affected).

#### func (*SyncedEnforcer) [DeleteUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L72) ¶

    func (e *SyncedEnforcer) DeleteUser(user [string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

DeleteUser deletes a user. Returns false if the user does not exist (aka not affected).

#### func (*SyncedEnforcer) [Enforce](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L191) ¶

    func (e *SyncedEnforcer) Enforce(rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).

#### func (*SyncedEnforcer) [EnforceEx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L205) ¶ added in v2.29.1

    func (e *SyncedEnforcer) EnforceEx(rvals ...interface{}) ([bool](/builtin#bool), [][string](/builtin#string), [error](/builtin#error))

EnforceEx explain enforcement by informing matched rules.

#### func (*SyncedEnforcer) [EnforceExWithMatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L212) ¶ added in v2.29.1

    func (e *SyncedEnforcer) EnforceExWithMatcher(matcher [string](/builtin#string), rvals ...interface{}) ([bool](/builtin#bool), [][string](/builtin#string), [error](/builtin#error))

EnforceExWithMatcher use a custom matcher and explain enforcement by informing matched rules.

#### func (*SyncedEnforcer) [EnforceWithMatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L198) ¶ added in v2.29.1

    func (e *SyncedEnforcer) EnforceWithMatcher(matcher [string](/builtin#string), rvals ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

EnforceWithMatcher use a custom matcher to decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (matcher, sub, obj, act), use model matcher by default when matcher is "".

#### func (*SyncedEnforcer) [GetAllActions](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L261) ¶

    func (e *SyncedEnforcer) GetAllActions() ([][string](/builtin#string), [error](/builtin#error))

GetAllActions gets the list of actions that show up in the current policy.

#### func (*SyncedEnforcer) [GetAllNamedActions](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L268) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetAllNamedActions(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedActions gets the list of actions that show up in the current named policy.

#### func (*SyncedEnforcer) [GetAllNamedObjects](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L254) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetAllNamedObjects(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedObjects gets the list of objects that show up in the current named policy.

#### func (*SyncedEnforcer) [GetAllNamedRoles](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L282) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetAllNamedRoles(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedRoles gets the list of roles that show up in the current named policy.

#### func (*SyncedEnforcer) [GetAllNamedSubjects](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L240) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetAllNamedSubjects(ptype [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetAllNamedSubjects gets the list of subjects that show up in the current named policy.

#### func (*SyncedEnforcer) [GetAllObjects](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L247) ¶

    func (e *SyncedEnforcer) GetAllObjects() ([][string](/builtin#string), [error](/builtin#error))

GetAllObjects gets the list of objects that show up in the current policy.

#### func (*SyncedEnforcer) [GetAllRoles](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L275) ¶

    func (e *SyncedEnforcer) GetAllRoles() ([][string](/builtin#string), [error](/builtin#error))

GetAllRoles gets the list of roles that show up in the current policy.

#### func (*SyncedEnforcer) [GetAllSubjects](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L233) ¶

    func (e *SyncedEnforcer) GetAllSubjects() ([][string](/builtin#string), [error](/builtin#error))

GetAllSubjects gets the list of subjects that show up in the current policy.

#### func (*SyncedEnforcer) [GetFilteredGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L324) ¶

    func (e *SyncedEnforcer) GetFilteredGroupingPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

#### func (*SyncedEnforcer) [GetFilteredNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L338) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetFilteredNamedGroupingPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredNamedGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.

#### func (*SyncedEnforcer) [GetFilteredNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L310) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetFilteredNamedPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredNamedPolicy gets all the authorization rules in the named policy, field filters can be specified.

#### func (*SyncedEnforcer) [GetFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L296) ¶

    func (e *SyncedEnforcer) GetFilteredPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.

#### func (*SyncedEnforcer) [GetGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L317) ¶

    func (e *SyncedEnforcer) GetGroupingPolicy() ([][][string](/builtin#string), [error](/builtin#error))

GetGroupingPolicy gets all the role inheritance rules in the policy.

#### func (*SyncedEnforcer) [GetImplicitObjectPatternsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L214) ¶ added in v2.121.0

    func (e *SyncedEnforcer) GetImplicitObjectPatternsForUser(user [string](/builtin#string), domain [string](/builtin#string), action [string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetImplicitObjectPatternsForUser returns all object patterns (with wildcards) that a user has for a given domain and action. For example: p, admin, chronicle/123, location/*, read p, user, chronicle/456, location/789, read g, alice, admin g, bob, user

GetImplicitObjectPatternsForUser("alice", "chronicle/123", "read") will return ["location/*"]. GetImplicitObjectPatternsForUser("bob", "chronicle/456", "read") will return ["location/789"].

#### func (*SyncedEnforcer) [GetImplicitPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L170) ¶ added in v2.13.0

    func (e *SyncedEnforcer) GetImplicitPermissionsForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetImplicitPermissionsForUser gets implicit permissions for a user or role. Compared to GetPermissionsForUser(), this function retrieves permissions for inherited roles. For example: p, admin, data1, read p, alice, data2, read g, alice, admin

GetPermissionsForUser("alice") can only get: [["alice", "data2", "read"]]. But GetImplicitPermissionsForUser("alice") will get: [["admin", "data1", "read"], ["alice", "data2", "read"]].

#### func (*SyncedEnforcer) [GetImplicitRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L155) ¶ added in v2.13.0

    func (e *SyncedEnforcer) GetImplicitRolesForUser(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetImplicitRolesForUser gets implicit roles that a user has. Compared to GetRolesForUser(), this function retrieves indirect roles besides direct roles. For example: g, alice, role:admin g, role:admin, role:user

GetRolesForUser("alice") can only get: ["role:admin"]. But GetImplicitRolesForUser("alice") will get: ["role:admin", "role:user"].

#### func (*SyncedEnforcer) [GetImplicitUsersForPermission](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L199) ¶ added in v2.13.0

    func (e *SyncedEnforcer) GetImplicitUsersForPermission(permission ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetImplicitUsersForPermission gets implicit users for a permission. For example: p, admin, data1, read p, bob, data1, read g, alice, admin

GetImplicitUsersForPermission("data1", "read") will get: ["alice", "bob"]. Note: only users will be returned, roles (2nd arg in "g") will be excluded.

#### func (*SyncedEnforcer) [GetLock](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L51) ¶ added in v2.52.1

    func (e *SyncedEnforcer) GetLock() *[sync](/sync).[RWMutex](/sync#RWMutex)

GetLock return the private RWMutex lock.

#### func (*SyncedEnforcer) [GetNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L331) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetNamedGroupingPolicy(ptype [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedGroupingPolicy gets all the role inheritance rules in the policy.

#### func (*SyncedEnforcer) [GetNamedImplicitPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L185) ¶ added in v2.45.0

    func (e *SyncedEnforcer) GetNamedImplicitPermissionsForUser(ptype [string](/builtin#string), gtype [string](/builtin#string), user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedImplicitPermissionsForUser gets implicit permissions for a user or role by named policy. Compared to GetNamedPermissionsForUser(), this function retrieves permissions for inherited roles. For example: p, admin, data1, read p2, admin, create g, alice, admin

GetImplicitPermissionsForUser("alice") can only get: [["admin", "data1", "read"]], whose policy is default policy "p" But you can specify the named policy "p2" to get: [["admin", "create"]] by GetNamedImplicitPermissionsForUser("p2","alice").

#### func (*SyncedEnforcer) [GetNamedPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L134) ¶ added in v2.45.0

    func (e *SyncedEnforcer) GetNamedPermissionsForUser(ptype [string](/builtin#string), user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedPermissionsForUser gets permissions for a user or role by named policy.

#### func (*SyncedEnforcer) [GetNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L303) ¶ added in v2.0.2

    func (e *SyncedEnforcer) GetNamedPolicy(ptype [string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetNamedPolicy gets all the authorization rules in the named policy.

#### func (*SyncedEnforcer) [GetNamedRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L63) ¶ added in v2.125.0

    func (e *SyncedEnforcer) GetNamedRoleManager(ptype [string](/builtin#string)) [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager)

GetNamedRoleManager gets the role manager for the named policy with synchronization.

#### func (*SyncedEnforcer) [GetPermissionsForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L127) ¶

    func (e *SyncedEnforcer) GetPermissionsForUser(user [string](/builtin#string), domain ...[string](/builtin#string)) ([][][string](/builtin#string), [error](/builtin#error))

GetPermissionsForUser gets permissions for a user or role.

#### func (*SyncedEnforcer) [GetPermissionsForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_synced.go#L32) ¶

    func (e *SyncedEnforcer) GetPermissionsForUserInDomain(user [string](/builtin#string), domain [string](/builtin#string)) [][][string](/builtin#string)

GetPermissionsForUserInDomain gets permissions for a user or role inside a domain.

#### func (*SyncedEnforcer) [GetPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L289) ¶

    func (e *SyncedEnforcer) GetPolicy() ([][][string](/builtin#string), [error](/builtin#error))

GetPolicy gets all the authorization rules in the policy.

#### func (*SyncedEnforcer) [GetRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L56) ¶ added in v2.125.0

    func (e *SyncedEnforcer) GetRoleManager() [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager)

GetRoleManager gets the current role manager with synchronization.

#### func (*SyncedEnforcer) [GetRolesForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L18) ¶

    func (e *SyncedEnforcer) GetRolesForUser(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetRolesForUser gets the roles that a user has.

#### func (*SyncedEnforcer) [GetRolesForUserInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_synced.go#L25) ¶

    func (e *SyncedEnforcer) GetRolesForUserInDomain(name [string](/builtin#string), domain [string](/builtin#string)) [][string](/builtin#string)

GetRolesForUserInDomain gets the roles that a user has inside a domain.

#### func (*SyncedEnforcer) [GetUsersForRole](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L25) ¶

    func (e *SyncedEnforcer) GetUsersForRole(name [string](/builtin#string), domain ...[string](/builtin#string)) ([][string](/builtin#string), [error](/builtin#error))

GetUsersForRole gets the users that has a role.

#### func (*SyncedEnforcer) [GetUsersForRoleInDomain](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_with_domains_synced.go#L18) ¶

    func (e *SyncedEnforcer) GetUsersForRoleInDomain(name [string](/builtin#string), domain [string](/builtin#string)) [][string](/builtin#string)

GetUsersForRoleInDomain gets the users that has a role inside a domain. Add by Gordon.

#### func (*SyncedEnforcer) [HasGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L493) ¶

    func (e *SyncedEnforcer) HasGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasGroupingPolicy determines whether a role inheritance rule exists.

#### func (*SyncedEnforcer) [HasNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L500) ¶ added in v2.0.2

    func (e *SyncedEnforcer) HasNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasNamedGroupingPolicy determines whether a named role inheritance rule exists.

#### func (*SyncedEnforcer) [HasNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L352) ¶ added in v2.0.2

    func (e *SyncedEnforcer) HasNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasNamedPolicy determines whether a named authorization rule exists.

#### func (*SyncedEnforcer) [HasPermissionForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L141) ¶

    func (e *SyncedEnforcer) HasPermissionForUser(user [string](/builtin#string), permission ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

HasPermissionForUser determines whether a user has a permission.

#### func (*SyncedEnforcer) [HasPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L345) ¶

    func (e *SyncedEnforcer) HasPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

HasPolicy determines whether an authorization rule exists.

#### func (*SyncedEnforcer) [HasRoleForUser](https://github.com/casbin/casbin/blob/v2.135.0/rbac_api_synced.go#L32) ¶

    func (e *SyncedEnforcer) HasRoleForUser(name [string](/builtin#string), role [string](/builtin#string), domain ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

HasRoleForUser determines whether a user has a role.

#### func (*SyncedEnforcer) [IsAutoLoadingRunning](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L84) ¶ added in v2.11.3

    func (e *SyncedEnforcer) IsAutoLoadingRunning() [bool](/builtin#bool)

IsAutoLoadingRunning check if SyncedEnforcer is auto loading policies.

#### func (*SyncedEnforcer) [LoadFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L163) ¶ added in v2.2.2

    func (e *SyncedEnforcer) LoadFilteredPolicy(filter interface{}) [error](/builtin#error)

LoadFilteredPolicy reloads a filtered policy from file/database.

#### func (*SyncedEnforcer) [LoadIncrementalFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L170) ¶ added in v2.11.0

    func (e *SyncedEnforcer) LoadIncrementalFilteredPolicy(filter interface{}) [error](/builtin#error)

LoadIncrementalFilteredPolicy reloads a filtered policy from file/database.

#### func (*SyncedEnforcer) [LoadModel](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L132) ¶ added in v2.19.6

    func (e *SyncedEnforcer) LoadModel() [error](/builtin#error)

LoadModel reloads the model from the model CONF file.

#### func (*SyncedEnforcer) [LoadPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L146) ¶

    func (e *SyncedEnforcer) LoadPolicy() [error](/builtin#error)

LoadPolicy reloads the policy from file/database.

#### func (*SyncedEnforcer) [RemoveFilteredGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L575) ¶

    func (e *SyncedEnforcer) RemoveFilteredGroupingPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredGroupingPolicy removes a role inheritance rule from the current policy, field filters can be specified.

#### func (*SyncedEnforcer) [RemoveFilteredNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L620) ¶ added in v2.0.2

    func (e *SyncedEnforcer) RemoveFilteredNamedGroupingPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredNamedGroupingPolicy removes a role inheritance rule from the current named policy, field filters can be specified.

#### func (*SyncedEnforcer) [RemoveFilteredNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L486) ¶ added in v2.0.2

    func (e *SyncedEnforcer) RemoveFilteredNamedPolicy(ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredNamedPolicy removes an authorization rule from the current named policy, field filters can be specified.

#### func (*SyncedEnforcer) [RemoveFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L465) ¶

    func (e *SyncedEnforcer) RemoveFilteredPolicy(fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveFilteredPolicy removes an authorization rule from the current policy, field filters can be specified.

#### func (*SyncedEnforcer) [RemoveGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L568) ¶ added in v2.25.1

    func (e *SyncedEnforcer) RemoveGroupingPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveGroupingPolicies removes role inheritance rules from the current policy.

#### func (*SyncedEnforcer) [RemoveGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L561) ¶

    func (e *SyncedEnforcer) RemoveGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveGroupingPolicy removes a role inheritance rule from the current policy.

#### func (*SyncedEnforcer) [RemoveNamedGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L589) ¶ added in v2.25.1

    func (e *SyncedEnforcer) RemoveNamedGroupingPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedGroupingPolicies removes role inheritance rules from the current named policy.

#### func (*SyncedEnforcer) [RemoveNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L582) ¶ added in v2.0.2

    func (e *SyncedEnforcer) RemoveNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedGroupingPolicy removes a role inheritance rule from the current named policy.

#### func (*SyncedEnforcer) [RemoveNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L479) ¶ added in v2.25.1

    func (e *SyncedEnforcer) RemoveNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPolicies removes authorization rules from the current named policy.

#### func (*SyncedEnforcer) [RemoveNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L472) ¶ added in v2.0.2

    func (e *SyncedEnforcer) RemoveNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPolicy removes an authorization rule from the current named policy.

#### func (*SyncedEnforcer) [RemovePolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L458) ¶ added in v2.25.1

    func (e *SyncedEnforcer) RemovePolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemovePolicies removes authorization rules from the current policy.

#### func (*SyncedEnforcer) [RemovePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L413) ¶

    func (e *SyncedEnforcer) RemovePolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemovePolicy removes an authorization rule from the current policy.

#### func (*SyncedEnforcer) [SavePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L177) ¶

    func (e *SyncedEnforcer) SavePolicy() [error](/builtin#error)

SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.

#### func (*SyncedEnforcer) [SelfAddPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L639) ¶ added in v2.62.0

    func (e *SyncedEnforcer) SelfAddPolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SelfAddPoliciesEx](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L645) ¶ added in v2.63.0

    func (e *SyncedEnforcer) SelfAddPoliciesEx(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SelfAddPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L633) ¶ added in v2.62.0

    func (e *SyncedEnforcer) SelfAddPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SelfRemoveFilteredPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L663) ¶ added in v2.62.0

    func (e *SyncedEnforcer) SelfRemoveFilteredPolicy(sec [string](/builtin#string), ptype [string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SelfRemovePolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L657) ¶ added in v2.62.0

    func (e *SyncedEnforcer) SelfRemovePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SelfRemovePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L651) ¶ added in v2.62.0

    func (e *SyncedEnforcer) SelfRemovePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), rule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SelfUpdatePolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L675) ¶ added in v2.62.0

    func (e *SyncedEnforcer) SelfUpdatePolicies(sec [string](/builtin#string), ptype [string](/builtin#string), oldRules, newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SelfUpdatePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L669) ¶ added in v2.62.0

    func (e *SyncedEnforcer) SelfUpdatePolicy(sec [string](/builtin#string), ptype [string](/builtin#string), oldRule, newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [SetNamedRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L77) ¶ added in v2.125.0

    func (e *SyncedEnforcer) SetNamedRoleManager(ptype [string](/builtin#string), rm [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager))

SetNamedRoleManager sets the role manager for the named policy with synchronization.

#### func (*SyncedEnforcer) [SetRoleManager](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L70) ¶ added in v2.125.0

    func (e *SyncedEnforcer) SetRoleManager(rm [rbac](/github.com/casbin/casbin/v2@v2.135.0/rbac).[RoleManager](/github.com/casbin/casbin/v2@v2.135.0/rbac#RoleManager))

SetRoleManager sets the current role manager with synchronization.

#### func (*SyncedEnforcer) [SetWatcher](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L125) ¶

    func (e *SyncedEnforcer) SetWatcher(watcher [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[Watcher](/github.com/casbin/casbin/v2@v2.135.0/persist#Watcher)) [error](/builtin#error)

SetWatcher sets the current watcher.

#### func (*SyncedEnforcer) [StartAutoLoadPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L89) ¶

    func (e *SyncedEnforcer) StartAutoLoadPolicy(d [time](/time).[Duration](/time#Duration))

StartAutoLoadPolicy starts a go routine that will every specified duration call LoadPolicy.

#### func (*SyncedEnforcer) [StopAutoLoadPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L118) ¶

    func (e *SyncedEnforcer) StopAutoLoadPolicy()

StopAutoLoadPolicy causes the go routine to exit.

#### func (*SyncedEnforcer) [UpdateFilteredNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L451) ¶ added in v2.28.0

    func (e *SyncedEnforcer) UpdateFilteredNamedPolicies(ptype [string](/builtin#string), newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdateFilteredPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L445) ¶ added in v2.28.0

    func (e *SyncedEnforcer) UpdateFilteredPolicies(newPolicies [][][string](/builtin#string), fieldIndex [int](/builtin#int), fieldValues ...[string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdateGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L601) ¶ added in v2.41.0

    func (e *SyncedEnforcer) UpdateGroupingPolicies(oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdateGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L595) ¶ added in v2.25.1

    func (e *SyncedEnforcer) UpdateGroupingPolicy(oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdateNamedGroupingPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L613) ¶ added in v2.41.0

    func (e *SyncedEnforcer) UpdateNamedGroupingPolicies(ptype [string](/builtin#string), oldRules [][][string](/builtin#string), newRules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdateNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L607) ¶ added in v2.25.1

    func (e *SyncedEnforcer) UpdateNamedGroupingPolicy(ptype [string](/builtin#string), oldRule [][string](/builtin#string), newRule [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdateNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L439) ¶ added in v2.25.1

    func (e *SyncedEnforcer) UpdateNamedPolicies(ptype [string](/builtin#string), p1 [][][string](/builtin#string), p2 [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdateNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L426) ¶ added in v2.25.1

    func (e *SyncedEnforcer) UpdateNamedPolicy(ptype [string](/builtin#string), p1 [][string](/builtin#string), p2 [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

#### func (*SyncedEnforcer) [UpdatePolicies](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L433) ¶ added in v2.25.1

    func (e *SyncedEnforcer) UpdatePolicies(oldPolices [][][string](/builtin#string), newPolicies [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdatePolicies updates authorization rules from the current policies.

#### func (*SyncedEnforcer) [UpdatePolicy](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_synced.go#L420) ¶ added in v2.25.1

    func (e *SyncedEnforcer) UpdatePolicy(oldPolicy [][string](/builtin#string), newPolicy [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdatePolicy updates an authorization rule from the current policy.

#### type [Transaction](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L35) ¶ added in v2.123.0

    type Transaction struct {
     // contains filtered or unexported fields
    }

Transaction represents a Casbin transaction. It provides methods to perform policy operations within a transaction. and commit or rollback all changes atomically.

#### func (*Transaction) [AddGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L305) ¶ added in v2.123.0

    func (tx *Transaction) AddGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddGroupingPolicy adds a grouping policy within the transaction.

#### func (*Transaction) [AddNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L310) ¶ added in v2.123.0

    func (tx *Transaction) AddNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedGroupingPolicy adds a named grouping policy within the transaction.

#### func (*Transaction) [AddNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L120) ¶ added in v2.123.0

    func (tx *Transaction) AddNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPolicies adds multiple named policies within the transaction.

#### func (*Transaction) [AddNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L81) ¶ added in v2.123.0

    func (tx *Transaction) AddNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddNamedPolicy adds a named policy within the transaction. The policy is buffered and will be applied when the transaction is committed.

#### func (*Transaction) [AddPolicies](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L115) ¶ added in v2.123.0

    func (tx *Transaction) AddPolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

AddPolicies adds multiple policies within the transaction.

#### func (*Transaction) [AddPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L50) ¶ added in v2.123.0

    func (tx *Transaction) AddPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

AddPolicy adds a policy within the transaction. The policy is buffered and will be applied when the transaction is committed.

#### func (*Transaction) [Commit](https://github.com/casbin/casbin/blob/v2.135.0/transaction_commit.go#L27) ¶ added in v2.123.0

    func (tx *Transaction) Commit() [error](/builtin#error)

Commit commits the transaction using a two-phase commit protocol. Phase 1: Apply all operations to the database Phase 2: Apply changes to the in-memory model and rebuild role links.

#### func (*Transaction) [GetBufferedModel](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L384) ¶ added in v2.123.0

    func (tx *Transaction) GetBufferedModel() ([model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model), [error](/builtin#error))

GetBufferedModel returns the model as it would look after applying all buffered operations. This is useful for preview or validation purposes within the transaction.

#### func (*Transaction) [HasOperations](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L396) ¶ added in v2.123.0

    func (tx *Transaction) HasOperations() [bool](/builtin#bool)

HasOperations returns true if the transaction has any buffered operations.

#### func (*Transaction) [IsActive](https://github.com/casbin/casbin/blob/v2.135.0/transaction_commit.go#L260) ¶ added in v2.123.0

    func (tx *Transaction) IsActive() [bool](/builtin#bool)

IsActive returns true if the transaction is still active (not committed or rolled back).

#### func (*Transaction) [IsCommitted](https://github.com/casbin/casbin/blob/v2.135.0/transaction_commit.go#L246) ¶ added in v2.123.0

    func (tx *Transaction) IsCommitted() [bool](/builtin#bool)

IsCommitted returns true if the transaction has been committed.

#### func (*Transaction) [IsRolledBack](https://github.com/casbin/casbin/blob/v2.135.0/transaction_commit.go#L253) ¶ added in v2.123.0

    func (tx *Transaction) IsRolledBack() [bool](/builtin#bool)

IsRolledBack returns true if the transaction has been rolled back.

#### func (*Transaction) [OperationCount](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L403) ¶ added in v2.123.0

    func (tx *Transaction) OperationCount() [int](/builtin#int)

OperationCount returns the number of buffered operations in the transaction.

#### func (*Transaction) [RemoveGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L344) ¶ added in v2.123.0

    func (tx *Transaction) RemoveGroupingPolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveGroupingPolicy removes a grouping policy within the transaction.

#### func (*Transaction) [RemoveNamedGroupingPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L349) ¶ added in v2.123.0

    func (tx *Transaction) RemoveNamedGroupingPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedGroupingPolicy removes a named grouping policy within the transaction.

#### func (*Transaction) [RemoveNamedPolicies](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L210) ¶ added in v2.123.0

    func (tx *Transaction) RemoveNamedPolicies(ptype [string](/builtin#string), rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPolicies removes multiple named policies within the transaction.

#### func (*Transaction) [RemoveNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L171) ¶ added in v2.123.0

    func (tx *Transaction) RemoveNamedPolicy(ptype [string](/builtin#string), params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemoveNamedPolicy removes a named policy within the transaction.

#### func (*Transaction) [RemovePolicies](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L205) ¶ added in v2.123.0

    func (tx *Transaction) RemovePolicies(rules [][][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

RemovePolicies removes multiple policies within the transaction.

#### func (*Transaction) [RemovePolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L166) ¶ added in v2.123.0

    func (tx *Transaction) RemovePolicy(params ...interface{}) ([bool](/builtin#bool), [error](/builtin#error))

RemovePolicy removes a policy within the transaction.

#### func (*Transaction) [Rollback](https://github.com/casbin/casbin/blob/v2.135.0/transaction_commit.go#L105) ¶ added in v2.123.0

    func (tx *Transaction) Rollback() [error](/builtin#error)

Rollback rolls back the transaction. This will rollback the database transaction and clear the transaction state.

#### func (*Transaction) [UpdateNamedPolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L261) ¶ added in v2.123.0

    func (tx *Transaction) UpdateNamedPolicy(ptype [string](/builtin#string), oldPolicy [][string](/builtin#string), newPolicy [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdateNamedPolicy updates a named policy within the transaction.

#### func (*Transaction) [UpdatePolicy](https://github.com/casbin/casbin/blob/v2.135.0/transaction.go#L256) ¶ added in v2.123.0

    func (tx *Transaction) UpdatePolicy(oldPolicy [][string](/builtin#string), newPolicy [][string](/builtin#string)) ([bool](/builtin#bool), [error](/builtin#error))

UpdatePolicy updates a policy within the transaction.

#### type [TransactionBuffer](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L27) ¶ added in v2.123.0

    type TransactionBuffer struct {
     // contains filtered or unexported fields
    }

TransactionBuffer holds all policy changes made within a transaction. It maintains a list of operations and a snapshot of the model state at the beginning of the transaction.

#### func [NewTransactionBuffer](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L35) ¶ added in v2.123.0

    func NewTransactionBuffer(baseModel [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model)) *TransactionBuffer

NewTransactionBuffer creates a new transaction buffer with a model snapshot. The snapshot represents the state of the model at the beginning of the transaction.

#### func (*TransactionBuffer) [AddOperation](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L44) ¶ added in v2.123.0

    func (tb *TransactionBuffer) AddOperation(op [persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[PolicyOperation](/github.com/casbin/casbin/v2@v2.135.0/persist#PolicyOperation))

AddOperation adds a policy operation to the buffer. This operation will be applied when the transaction is committed.

#### func (*TransactionBuffer) [ApplyOperationsToModel](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L81) ¶ added in v2.123.0

    func (tb *TransactionBuffer) ApplyOperationsToModel(baseModel [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model)) ([model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model), [error](/builtin#error))

ApplyOperationsToModel applies all buffered operations to a model and returns the result. This simulates what the model would look like after all operations are applied. It's used for validation and preview purposes within the transaction.

#### func (*TransactionBuffer) [Clear](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L64) ¶ added in v2.123.0

    func (tb *TransactionBuffer) Clear()

Clear removes all buffered operations. This is typically called after a successful commit or rollback.

#### func (*TransactionBuffer) [GetModelSnapshot](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L72) ¶ added in v2.123.0

    func (tb *TransactionBuffer) GetModelSnapshot() [model](/github.com/casbin/casbin/v2@v2.135.0/model).[Model](/github.com/casbin/casbin/v2@v2.135.0/model#Model)

GetModelSnapshot returns the model snapshot taken at transaction start. This represents the original state before any transaction operations.

#### func (*TransactionBuffer) [GetOperations](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L52) ¶ added in v2.123.0

    func (tb *TransactionBuffer) GetOperations() [][persist](/github.com/casbin/casbin/v2@v2.135.0/persist).[PolicyOperation](/github.com/casbin/casbin/v2@v2.135.0/persist#PolicyOperation)

GetOperations returns all buffered operations. Returns a copy to prevent external modifications.

#### func (*TransactionBuffer) [HasOperations](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L120) ¶ added in v2.123.0

    func (tb *TransactionBuffer) HasOperations() [bool](/builtin#bool)

HasOperations returns true if there are any buffered operations.

#### func (*TransactionBuffer) [OperationCount](https://github.com/casbin/casbin/blob/v2.135.0/transaction_buffer.go#L127) ¶ added in v2.123.0

    func (tb *TransactionBuffer) OperationCount() [int](/builtin#int)

OperationCount returns the number of buffered operations.

#### type [TransactionalEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_transactional.go#L30) ¶ added in v2.123.0

    type TransactionalEnforcer struct {
     *Enforcer // Embedded enforcer for all standard functionality
     // contains filtered or unexported fields
    }

TransactionalEnforcer extends Enforcer with transaction support. It provides atomic policy operations through transactions.

#### func [NewTransactionalEnforcer](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_transactional.go#L39) ¶ added in v2.123.0

    func NewTransactionalEnforcer(params ...interface{}) (*TransactionalEnforcer, [error](/builtin#error))

NewTransactionalEnforcer creates a new TransactionalEnforcer. It accepts the same parameters as NewEnforcer.

#### func (*TransactionalEnforcer) [BeginTransaction](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_transactional.go#L52) ¶ added in v2.123.0

    func (te *TransactionalEnforcer) BeginTransaction(ctx [context](/context).[Context](/context#Context)) (*Transaction, [error](/builtin#error))

BeginTransaction starts a new transaction. Returns an error if a transaction is already in progress or if the adapter doesn't support transactions.

#### func (*TransactionalEnforcer) [GetTransaction](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_transactional.go#L83) ¶ added in v2.128.0

    func (te *TransactionalEnforcer) GetTransaction(id [string](/builtin#string)) *Transaction

GetTransaction returns a transaction by its ID, or nil if not found.

#### func (*TransactionalEnforcer) [IsTransactionActive](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_transactional.go#L91) ¶ added in v2.128.0

    func (te *TransactionalEnforcer) IsTransactionActive(id [string](/builtin#string)) [bool](/builtin#bool)

IsTransactionActive returns true if the transaction with the given ID is active.

#### func (*TransactionalEnforcer) [WithTransaction](https://github.com/casbin/casbin/blob/v2.135.0/enforcer_transactional.go#L101) ¶ added in v2.123.0

    func (te *TransactionalEnforcer) WithTransaction(ctx [context](/context).[Context](/context#Context), fn func(*Transaction) [error](/builtin#error)) [error](/builtin#error)

WithTransaction executes a function within a transaction. If the function returns an error, the transaction is rolled back. Otherwise, it's committed automatically.
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
