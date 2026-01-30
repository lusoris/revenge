// Package request provides QAR content request domain models.
package request

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the provision (content request) data access interface.
type Repository interface {
	// Provisions
	GetProvisionByID(ctx context.Context, id uuid.UUID) (*Provision, error)
	ListProvisions(ctx context.Context, params ListProvisionsParams) ([]Provision, error)
	CreateProvision(ctx context.Context, params CreateProvisionParams) (*Provision, error)
	UpdateProvisionStatus(ctx context.Context, id uuid.UUID, status ProvisionStatus, approvedBy *uuid.UUID, declinedReason string) (*Provision, error)
	UpdateProvisionPriority(ctx context.Context, id uuid.UUID, priority int) (*Provision, error)
	UpdateProvisionIntegration(ctx context.Context, id uuid.UUID, integrationID, integrationStatus string) (*Provision, error)
	SetProvisionAvailable(ctx context.Context, id uuid.UUID, actualCargoGB float64) (*Provision, error)
	SetProvisionAutoApproved(ctx context.Context, id uuid.UUID, articleID uuid.UUID) (*Provision, error)
	DeleteProvision(ctx context.Context, id uuid.UUID) error
	CountProvisionsByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	CountProvisionsByStatus(ctx context.Context, status ProvisionStatus) (int64, error)
	GetProvisionByExternalID(ctx context.Context, source ExternalSource, externalID string) (*Provision, error)
	SearchProvisions(ctx context.Context, query string, limit, offset int) ([]Provision, error)

	// Ayes (Votes)
	GetAye(ctx context.Context, provisionID, userID uuid.UUID) (*ProvisionAye, error)
	ListAyes(ctx context.Context, provisionID uuid.UUID) ([]ProvisionAye, error)
	CreateAye(ctx context.Context, provisionID, userID uuid.UUID) (*ProvisionAye, error)
	DeleteAye(ctx context.Context, provisionID, userID uuid.UUID) error
	HasUserVoted(ctx context.Context, provisionID, userID uuid.UUID) (bool, error)

	// Missives (Comments)
	GetMissive(ctx context.Context, id uuid.UUID) (*ProvisionMissive, error)
	ListMissives(ctx context.Context, provisionID uuid.UUID) ([]ProvisionMissive, error)
	CreateMissive(ctx context.Context, provisionID, userID uuid.UUID, message string, isCaptainOrder bool) (*ProvisionMissive, error)
	DeleteMissive(ctx context.Context, id uuid.UUID) error

	// Rations (User Quotas)
	GetRation(ctx context.Context, userID uuid.UUID) (*Ration, error)
	UpsertRation(ctx context.Context, userID uuid.UUID) (*Ration, error)
	UpdateRationLimits(ctx context.Context, userID uuid.UUID, params UpdateRationParams) (*Ration, error)
	IncrementRationUsage(ctx context.Context, userID uuid.UUID) (*Ration, error)
	AddRationCargoUsage(ctx context.Context, userID uuid.UUID, cargoGB float64) (*Ration, error)
	ResetDailyRations(ctx context.Context) error
	ResetWeeklyRations(ctx context.Context) error
	ResetMonthlyRations(ctx context.Context) error

	// Articles (Rules)
	GetArticle(ctx context.Context, id uuid.UUID) (*Article, error)
	ListArticles(ctx context.Context) ([]Article, error)
	ListEnabledArticles(ctx context.Context) ([]Article, error)
	ListArticlesByContentType(ctx context.Context, contentType ContentType) ([]Article, error)
	CreateArticle(ctx context.Context, params CreateArticleParams) (*Article, error)
	UpdateArticle(ctx context.Context, id uuid.UUID, params UpdateArticleParams) (*Article, error)
	DeleteArticle(ctx context.Context, id uuid.UUID) error
	SetArticleEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*Article, error)

	// Cargo Hold (Global Quotas)
	GetCargoHold(ctx context.Context) (*CargoHold, error)
	UpdateCargoHoldQuotas(ctx context.Context, totalQuota, expeditionQuota, voyageQuota float64) (*CargoHold, error)
	AddCargoHoldUsage(ctx context.Context, cargoGB float64, isExpedition bool) (*CargoHold, error)
}
