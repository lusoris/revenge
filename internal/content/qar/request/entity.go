// Package request provides QAR content request domain models.
// QAR obfuscation: requests → provisions, votes → ayes, comments → missives,
// quotas → rations, rules → articles.
package request

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProvisionStatus represents the status of a content request.
type ProvisionStatus string

const (
	ProvisionStatusPending    ProvisionStatus = "pending"
	ProvisionStatusApproved   ProvisionStatus = "approved"
	ProvisionStatusProcessing ProvisionStatus = "processing"
	ProvisionStatusAvailable  ProvisionStatus = "available"
	ProvisionStatusDeclined   ProvisionStatus = "declined"
	ProvisionStatusOnHold     ProvisionStatus = "on_hold"
)

// ContentType represents the type of content being requested.
type ContentType string

const (
	ContentTypeExpedition ContentType = "expedition" // Movie
	ContentTypeVoyage     ContentType = "voyage"     // Scene
)

// RequestSubtype represents the specific request type.
type RequestSubtype string

const (
	RequestSubtypeScene    RequestSubtype = "scene"      // Single scene/movie
	RequestSubtypePort     RequestSubtype = "port"       // All content from studio
	RequestSubtypeCrew     RequestSubtype = "crew"       // All content with performer
	RequestSubtypeFlagCombo RequestSubtype = "flag_combo" // Tag combination
)

// ExternalSource represents the metadata source.
type ExternalSource string

const (
	ExternalSourceStashDB  ExternalSource = "stashdb"
	ExternalSourceTPDB     ExternalSource = "tpdb"
	ExternalSourceWhisparr ExternalSource = "whisparr"
)

// ArticleAction represents what happens when a rule matches.
type ArticleAction string

const (
	ArticleActionAutoApprove     ArticleAction = "auto_approve"
	ArticleActionRequireApproval ArticleAction = "require_approval"
	ArticleActionDecline         ArticleAction = "decline"
	ArticleActionOnHold          ArticleAction = "on_hold"
)

// ArticleConditionType represents the type of condition to evaluate.
type ArticleConditionType string

const (
	ArticleConditionUserRole         ArticleConditionType = "user_role"
	ArticleConditionTrustScore       ArticleConditionType = "trust_score"
	ArticleConditionCrewPreference   ArticleConditionType = "crew_preference"   // Performer
	ArticleConditionPortPreference   ArticleConditionType = "port_preference"   // Studio
	ArticleConditionFlagPreference   ArticleConditionType = "flag_preference"   // Tags
	ArticleConditionStorageAvailable ArticleConditionType = "storage_available"
	ArticleConditionReleaseYear      ArticleConditionType = "release_year"
)

// Provision represents a content request (QAR: request → provision).
type Provision struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"userId"`
	ContentType     ContentType     `json:"contentType"`
	RequestSubtype  RequestSubtype  `json:"requestSubtype,omitempty"`
	ExternalID      string          `json:"externalId,omitempty"`
	ExternalSource  ExternalSource  `json:"externalSource,omitempty"`
	Title           string          `json:"title"`
	ReleaseYear     *int            `json:"releaseYear,omitempty"`
	Manifest        json.RawMessage `json:"manifest,omitempty"` // Type-specific metadata

	Status         ProvisionStatus `json:"status"`
	AutoApproved   bool            `json:"autoApproved"`
	AutoArticleID  *uuid.UUID      `json:"autoArticleId,omitempty"`

	ApprovedByUserID *uuid.UUID `json:"approvedByUserId,omitempty"`
	ApprovedAt       *time.Time `json:"approvedAt,omitempty"`
	DeclinedReason   string     `json:"declinedReason,omitempty"`

	Priority   int `json:"priority"`
	AyesCount  int `json:"ayesCount"`

	IntegrationID     string `json:"integrationId,omitempty"`
	IntegrationStatus string `json:"integrationStatus,omitempty"`

	EstimatedCargoGB *float64 `json:"estimatedCargoGb,omitempty"`
	ActualCargoGB    *float64 `json:"actualCargoGb,omitempty"`

	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	AvailableAt *time.Time `json:"availableAt,omitempty"`

	TriggeredByAutomation bool       `json:"triggeredByAutomation"`
	ParentProvisionID     *uuid.UUID `json:"parentProvisionId,omitempty"`
}

// ProvisionManifest contains type-specific metadata for a provision.
type ProvisionManifest struct {
	// For scene requests
	FlagIDs []uuid.UUID `json:"flagIds,omitempty"` // Tag IDs
	CrewIDs []uuid.UUID `json:"crewIds,omitempty"` // Performer IDs
	PortID  *uuid.UUID  `json:"portId,omitempty"`  // Studio ID

	// For port/crew requests (monitor mode)
	MonitorNew bool `json:"monitorNew,omitempty"` // Auto-add new releases

	// Quality preferences
	PreferredQuality string `json:"preferredQuality,omitempty"` // "4k", "1080p", etc.
}

// ProvisionAye represents a vote on a provision (QAR: vote → aye).
type ProvisionAye struct {
	ProvisionID uuid.UUID `json:"provisionId"`
	UserID      uuid.UUID `json:"userId"`
	VotedAt     time.Time `json:"votedAt"`
}

// ProvisionMissive represents a comment on a provision (QAR: comment → missive).
type ProvisionMissive struct {
	ID             uuid.UUID `json:"id"`
	ProvisionID    uuid.UUID `json:"provisionId"`
	UserID         uuid.UUID `json:"userId"`
	Message        string    `json:"message"`
	IsCaptainOrder bool      `json:"isCaptainOrder"` // Admin/mod comment
	CreatedAt      time.Time `json:"createdAt"`
}

// Ration represents per-user request quotas (QAR: quota → ration).
type Ration struct {
	UserID uuid.UUID `json:"userId"`

	// Request limits
	DailyLimit   int `json:"dailyLimit"`
	WeeklyLimit  int `json:"weeklyLimit"`
	MonthlyLimit int `json:"monthlyLimit"`
	DailyUsed    int `json:"dailyUsed"`
	WeeklyUsed   int `json:"weeklyUsed"`
	MonthlyUsed  int `json:"monthlyUsed"`

	// Storage quota
	CargoQuotaGB float64 `json:"cargoQuotaGb"`
	CargoUsedGB  float64 `json:"cargoUsedGb"`

	// Reset timestamps
	LastResetDaily   time.Time `json:"lastResetDaily"`
	LastResetWeekly  time.Time `json:"lastResetWeekly"`
	LastResetMonthly time.Time `json:"lastResetMonthly"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RemainingDaily returns remaining daily requests.
func (r *Ration) RemainingDaily() int {
	return r.DailyLimit - r.DailyUsed
}

// RemainingWeekly returns remaining weekly requests.
func (r *Ration) RemainingWeekly() int {
	return r.WeeklyLimit - r.WeeklyUsed
}

// RemainingMonthly returns remaining monthly requests.
func (r *Ration) RemainingMonthly() int {
	return r.MonthlyLimit - r.MonthlyUsed
}

// RemainingCargoGB returns remaining storage quota in GB.
func (r *Ration) RemainingCargoGB() float64 {
	return r.CargoQuotaGB - r.CargoUsedGB
}

// Article represents an auto-approval/decline rule (QAR: rule → article).
type Article struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`

	ContentType       *ContentType         `json:"contentType,omitempty"` // nil = all types
	ConditionType     ArticleConditionType `json:"conditionType"`
	ConditionValue    json.RawMessage      `json:"conditionValue"`
	Action            ArticleAction        `json:"action"`
	AutomationTrigger string               `json:"automationTrigger,omitempty"`

	Enabled  bool `json:"enabled"`
	Priority int  `json:"priority"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ArticleConditionValue represents common condition value structures.
type ArticleConditionValue struct {
	// For user_role
	Roles []string `json:"roles,omitempty"`

	// For trust_score
	MinTrustScore *int `json:"minTrustScore,omitempty"`

	// For storage_available
	MinFreeGB *float64 `json:"minFreeGb,omitempty"`

	// For release_year
	MinYear *int `json:"minYear,omitempty"`
	MaxYear *int `json:"maxYear,omitempty"`

	// For crew/port/flag preference
	MinWatchCount  *int     `json:"minWatchCount,omitempty"`
	CompletionRate *float64 `json:"completionRate,omitempty"`
}

// CargoHold represents global storage quotas (QAR: global storage → cargo hold).
type CargoHold struct {
	ID                 int       `json:"id"`
	TotalQuotaGB       float64   `json:"totalQuotaGb"`
	TotalUsedGB        float64   `json:"totalUsedGb"`
	ExpeditionQuotaGB  float64   `json:"expeditionQuotaGb"`
	ExpeditionUsedGB   float64   `json:"expeditionUsedGb"`
	VoyageQuotaGB      float64   `json:"voyageQuotaGb"`
	VoyageUsedGB       float64   `json:"voyageUsedGb"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// RemainingTotalGB returns remaining total storage.
func (c *CargoHold) RemainingTotalGB() float64 {
	return c.TotalQuotaGB - c.TotalUsedGB
}

// RemainingExpeditionGB returns remaining expedition storage.
func (c *CargoHold) RemainingExpeditionGB() float64 {
	return c.ExpeditionQuotaGB - c.ExpeditionUsedGB
}

// RemainingVoyageGB returns remaining voyage storage.
func (c *CargoHold) RemainingVoyageGB() float64 {
	return c.VoyageQuotaGB - c.VoyageUsedGB
}

// ProvisionWithDetails includes the provision and related data.
type ProvisionWithDetails struct {
	Provision
	UserName       string             `json:"userName,omitempty"`
	ApproverName   string             `json:"approverName,omitempty"`
	Missives       []ProvisionMissive `json:"missives,omitempty"`
	HasVoted       bool               `json:"hasVoted,omitempty"` // Current user voted
}

// CreateProvisionParams contains parameters for creating a provision.
type CreateProvisionParams struct {
	UserID         uuid.UUID
	ContentType    ContentType
	RequestSubtype RequestSubtype
	ExternalID     string
	ExternalSource ExternalSource
	Title          string
	ReleaseYear    *int
	Manifest       json.RawMessage
}

// ListProvisionsParams contains parameters for listing provisions.
type ListProvisionsParams struct {
	UserID      *uuid.UUID
	Status      *ProvisionStatus
	ContentType *ContentType
	Query       string
	Limit       int
	Offset      int
	SortBy      string // "created_at", "priority", "ayes_count"
	SortDesc    bool
}

// UpdateProvisionParams contains parameters for updating a provision.
type UpdateProvisionParams struct {
	Status            *ProvisionStatus
	Priority          *int
	DeclinedReason    *string
	IntegrationID     *string
	IntegrationStatus *string
	ActualCargoGB     *float64
	AvailableAt       *time.Time
}

// CreateArticleParams contains parameters for creating an article.
type CreateArticleParams struct {
	Name              string
	Description       string
	ContentType       *ContentType
	ConditionType     ArticleConditionType
	ConditionValue    json.RawMessage
	Action            ArticleAction
	AutomationTrigger string
	Enabled           bool
	Priority          int
}

// UpdateArticleParams contains parameters for updating an article.
type UpdateArticleParams struct {
	Name              *string
	Description       *string
	ContentType       *ContentType
	ConditionType     *ArticleConditionType
	ConditionValue    json.RawMessage
	Action            *ArticleAction
	AutomationTrigger *string
	Enabled           *bool
	Priority          *int
}

// UpdateRationParams contains parameters for updating a ration.
type UpdateRationParams struct {
	DailyLimit   *int
	WeeklyLimit  *int
	MonthlyLimit *int
	CargoQuotaGB *float64
}
