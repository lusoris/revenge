// Package request provides QAR content request domain models.
package request

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	adultdb "github.com/lusoris/revenge/internal/content/qar/db"
)

// Errors
var (
	ErrProvisionNotFound = errors.New("provision not found")
	ErrAyeNotFound       = errors.New("vote not found")
	ErrMissiveNotFound   = errors.New("comment not found")
	ErrRationNotFound    = errors.New("ration not found")
	ErrArticleNotFound   = errors.New("article not found")
	ErrCargoHoldNotFound = errors.New("cargo hold not found")
)

// SQLCRepository implements Repository using sqlc-generated queries.
type SQLCRepository struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewSQLCRepository creates a new SQLC-backed request repository.
func NewSQLCRepository(pool *pgxpool.Pool, logger *slog.Logger) Repository {
	return &SQLCRepository{
		queries: adultdb.New(pool),
		logger:  logger.With(slog.String("repository", "qar.request")),
	}
}

// ============================================================================
// PROVISIONS
// ============================================================================

func (r *SQLCRepository) GetProvisionByID(ctx context.Context, id uuid.UUID) (*Provision, error) {
	row, err := r.queries.GetProvisionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProvisionNotFound
		}
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) ListProvisions(ctx context.Context, params ListProvisionsParams) ([]Provision, error) {
	sortBy := "created_at"
	if params.SortBy != "" {
		sortBy = params.SortBy
	}

	limit := int32(params.Limit)
	if limit == 0 {
		limit = 20
	}

	rows, err := r.queries.ListProvisions(ctx, adultdb.ListProvisionsParams{
		Limit:  limit,
		Offset: int32(params.Offset),
		SortBy: sortBy,
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToProvisions(rows), nil
}

func (r *SQLCRepository) CreateProvision(ctx context.Context, params CreateProvisionParams) (*Provision, error) {
	row, err := r.queries.CreateProvision(ctx, adultdb.CreateProvisionParams{
		UserID:           params.UserID,
		ContentType:      string(params.ContentType),
		RequestSubtype:   ptrString(string(params.RequestSubtype)),
		ExternalID:       ptrString(params.ExternalID),
		ExternalSource:   ptrString(string(params.ExternalSource)),
		Title:            params.Title,
		ReleaseYear:      ptrInt32(params.ReleaseYear),
		Manifest:         params.Manifest,
		EstimatedCargoGb: pgtype.Numeric{},
	})
	if err != nil {
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) UpdateProvisionStatus(ctx context.Context, id uuid.UUID, status ProvisionStatus, approvedBy *uuid.UUID, declinedReason string) (*Provision, error) {
	var approvedByPG pgtype.UUID
	if approvedBy != nil {
		approvedByPG = pgtype.UUID{Bytes: *approvedBy, Valid: true}
	}

	var approvedAt pgtype.Timestamptz
	if status == ProvisionStatusApproved {
		approvedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	}

	row, err := r.queries.UpdateProvisionStatus(ctx, adultdb.UpdateProvisionStatusParams{
		ID:               id,
		Status:           string(status),
		ApprovedByUserID: approvedByPG,
		ApprovedAt:       approvedAt,
		DeclinedReason:   ptrString(declinedReason),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProvisionNotFound
		}
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) UpdateProvisionPriority(ctx context.Context, id uuid.UUID, priority int) (*Provision, error) {
	row, err := r.queries.UpdateProvisionPriority(ctx, adultdb.UpdateProvisionPriorityParams{
		ID:       id,
		Priority: int32(priority),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProvisionNotFound
		}
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) UpdateProvisionIntegration(ctx context.Context, id uuid.UUID, integrationID, integrationStatus string) (*Provision, error) {
	row, err := r.queries.UpdateProvisionIntegration(ctx, adultdb.UpdateProvisionIntegrationParams{
		ID:                id,
		IntegrationID:     ptrString(integrationID),
		IntegrationStatus: ptrString(integrationStatus),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProvisionNotFound
		}
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) SetProvisionAvailable(ctx context.Context, id uuid.UUID, actualCargoGB float64) (*Provision, error) {
	numericVal := pgtype.Numeric{}
	_ = numericVal.Scan(actualCargoGB)

	row, err := r.queries.UpdateProvisionAvailable(ctx, adultdb.UpdateProvisionAvailableParams{
		ID:            id,
		ActualCargoGb: numericVal,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProvisionNotFound
		}
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) SetProvisionAutoApproved(ctx context.Context, id uuid.UUID, articleID uuid.UUID) (*Provision, error) {
	row, err := r.queries.SetProvisionAutoApproved(ctx, adultdb.SetProvisionAutoApprovedParams{
		ID:            id,
		AutoArticleID: pgtype.UUID{Bytes: articleID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProvisionNotFound
		}
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) DeleteProvision(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteProvision(ctx, id)
}

func (r *SQLCRepository) CountProvisionsByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountProvisionsByUser(ctx, userID)
}

func (r *SQLCRepository) CountProvisionsByStatus(ctx context.Context, status ProvisionStatus) (int64, error) {
	return r.queries.CountProvisionsByStatus(ctx, string(status))
}

func (r *SQLCRepository) GetProvisionByExternalID(ctx context.Context, source ExternalSource, externalID string) (*Provision, error) {
	row, err := r.queries.GetProvisionByExternalID(ctx, adultdb.GetProvisionByExternalIDParams{
		ExternalSource: ptrString(string(source)),
		ExternalID:     ptrString(externalID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProvisionNotFound
		}
		return nil, err
	}
	return r.rowToProvision(&row), nil
}

func (r *SQLCRepository) SearchProvisions(ctx context.Context, query string, limit, offset int) ([]Provision, error) {
	rows, err := r.queries.SearchProvisions(ctx, adultdb.SearchProvisionsParams{
		Column1: &query,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}
	return r.rowsToProvisions(rows), nil
}

// ============================================================================
// AYES (Votes)
// ============================================================================

func (r *SQLCRepository) GetAye(ctx context.Context, provisionID, userID uuid.UUID) (*ProvisionAye, error) {
	row, err := r.queries.GetProvisionAye(ctx, adultdb.GetProvisionAyeParams{
		ProvisionID: provisionID,
		UserID:      userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAyeNotFound
		}
		return nil, err
	}
	return &ProvisionAye{
		ProvisionID: row.ProvisionID,
		UserID:      row.UserID,
		VotedAt:     row.VotedAt,
	}, nil
}

func (r *SQLCRepository) ListAyes(ctx context.Context, provisionID uuid.UUID) ([]ProvisionAye, error) {
	rows, err := r.queries.ListProvisionAyes(ctx, provisionID)
	if err != nil {
		return nil, err
	}
	result := make([]ProvisionAye, len(rows))
	for i, row := range rows {
		result[i] = ProvisionAye{
			ProvisionID: row.ProvisionID,
			UserID:      row.UserID,
			VotedAt:     row.VotedAt,
		}
	}
	return result, nil
}

func (r *SQLCRepository) CreateAye(ctx context.Context, provisionID, userID uuid.UUID) (*ProvisionAye, error) {
	row, err := r.queries.CreateProvisionAye(ctx, adultdb.CreateProvisionAyeParams{
		ProvisionID: provisionID,
		UserID:      userID,
	})
	if err != nil {
		return nil, err
	}
	return &ProvisionAye{
		ProvisionID: row.ProvisionID,
		UserID:      row.UserID,
		VotedAt:     row.VotedAt,
	}, nil
}

func (r *SQLCRepository) DeleteAye(ctx context.Context, provisionID, userID uuid.UUID) error {
	return r.queries.DeleteProvisionAye(ctx, adultdb.DeleteProvisionAyeParams{
		ProvisionID: provisionID,
		UserID:      userID,
	})
}

func (r *SQLCRepository) HasUserVoted(ctx context.Context, provisionID, userID uuid.UUID) (bool, error) {
	return r.queries.HasUserVoted(ctx, adultdb.HasUserVotedParams{
		ProvisionID: provisionID,
		UserID:      userID,
	})
}

// ============================================================================
// MISSIVES (Comments)
// ============================================================================

func (r *SQLCRepository) GetMissive(ctx context.Context, id uuid.UUID) (*ProvisionMissive, error) {
	row, err := r.queries.GetProvisionMissive(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMissiveNotFound
		}
		return nil, err
	}
	return r.rowToMissive(&row), nil
}

func (r *SQLCRepository) ListMissives(ctx context.Context, provisionID uuid.UUID) ([]ProvisionMissive, error) {
	rows, err := r.queries.ListProvisionMissives(ctx, provisionID)
	if err != nil {
		return nil, err
	}
	result := make([]ProvisionMissive, len(rows))
	for i, row := range rows {
		result[i] = *r.rowToMissive(&row)
	}
	return result, nil
}

func (r *SQLCRepository) CreateMissive(ctx context.Context, provisionID, userID uuid.UUID, message string, isCaptainOrder bool) (*ProvisionMissive, error) {
	row, err := r.queries.CreateProvisionMissive(ctx, adultdb.CreateProvisionMissiveParams{
		ProvisionID:    provisionID,
		UserID:         pgtype.UUID{Bytes: userID, Valid: true},
		Message:        message,
		IsCaptainOrder: isCaptainOrder,
	})
	if err != nil {
		return nil, err
	}
	return r.rowToMissive(&row), nil
}

func (r *SQLCRepository) DeleteMissive(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteProvisionMissive(ctx, id)
}

// ============================================================================
// RATIONS (User Quotas)
// ============================================================================

func (r *SQLCRepository) GetRation(ctx context.Context, userID uuid.UUID) (*Ration, error) {
	row, err := r.queries.GetRation(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRationNotFound
		}
		return nil, err
	}
	return r.rowToRation(&row), nil
}

func (r *SQLCRepository) UpsertRation(ctx context.Context, userID uuid.UUID) (*Ration, error) {
	row, err := r.queries.UpsertRation(ctx, userID)
	if err != nil {
		return nil, err
	}
	return r.rowToRation(&row), nil
}

func (r *SQLCRepository) UpdateRationLimits(ctx context.Context, userID uuid.UUID, params UpdateRationParams) (*Ration, error) {
	row, err := r.queries.UpdateRationLimits(ctx, adultdb.UpdateRationLimitsParams{
		UserID:       userID,
		DailyLimit:   derefInt32(params.DailyLimit),
		WeeklyLimit:  derefInt32(params.WeeklyLimit),
		MonthlyLimit: derefInt32(params.MonthlyLimit),
		CargoQuotaGb: numericFromFloat(params.CargoQuotaGB),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRationNotFound
		}
		return nil, err
	}
	return r.rowToRation(&row), nil
}

func (r *SQLCRepository) IncrementRationUsage(ctx context.Context, userID uuid.UUID) (*Ration, error) {
	row, err := r.queries.IncrementRationUsage(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRationNotFound
		}
		return nil, err
	}
	return r.rowToRation(&row), nil
}

func (r *SQLCRepository) AddRationCargoUsage(ctx context.Context, userID uuid.UUID, cargoGB float64) (*Ration, error) {
	row, err := r.queries.AddRationCargoUsage(ctx, adultdb.AddRationCargoUsageParams{
		UserID:      userID,
		CargoUsedGb: numericFromFloat(&cargoGB),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRationNotFound
		}
		return nil, err
	}
	return r.rowToRation(&row), nil
}

func (r *SQLCRepository) ResetDailyRations(ctx context.Context) error {
	return r.queries.ResetDailyRations(ctx)
}

func (r *SQLCRepository) ResetWeeklyRations(ctx context.Context) error {
	return r.queries.ResetWeeklyRations(ctx)
}

func (r *SQLCRepository) ResetMonthlyRations(ctx context.Context) error {
	return r.queries.ResetMonthlyRations(ctx)
}

// ============================================================================
// ARTICLES (Rules)
// ============================================================================

func (r *SQLCRepository) GetArticle(ctx context.Context, id uuid.UUID) (*Article, error) {
	row, err := r.queries.GetArticle(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return r.rowToArticle(&row), nil
}

func (r *SQLCRepository) ListArticles(ctx context.Context) ([]Article, error) {
	rows, err := r.queries.ListArticles(ctx)
	if err != nil {
		return nil, err
	}
	return r.rowsToArticles(rows), nil
}

func (r *SQLCRepository) ListEnabledArticles(ctx context.Context) ([]Article, error) {
	rows, err := r.queries.ListEnabledArticles(ctx)
	if err != nil {
		return nil, err
	}
	return r.rowsToArticles(rows), nil
}

func (r *SQLCRepository) ListArticlesByContentType(ctx context.Context, contentType ContentType) ([]Article, error) {
	rows, err := r.queries.ListArticlesByContentType(ctx, ptrString(string(contentType)))
	if err != nil {
		return nil, err
	}
	return r.rowsToArticles(rows), nil
}

func (r *SQLCRepository) CreateArticle(ctx context.Context, params CreateArticleParams) (*Article, error) {
	var contentType *string
	if params.ContentType != nil {
		s := string(*params.ContentType)
		contentType = &s
	}

	row, err := r.queries.CreateArticle(ctx, adultdb.CreateArticleParams{
		Name:              params.Name,
		Description:       ptrString(params.Description),
		ContentType:       contentType,
		ConditionType:     string(params.ConditionType),
		ConditionValue:    params.ConditionValue,
		Action:            string(params.Action),
		AutomationTrigger: ptrString(params.AutomationTrigger),
		Enabled:           params.Enabled,
		Priority:          int32(params.Priority),
	})
	if err != nil {
		return nil, err
	}
	return r.rowToArticle(&row), nil
}

func (r *SQLCRepository) UpdateArticle(ctx context.Context, id uuid.UUID, params UpdateArticleParams) (*Article, error) {
	var contentType *string
	if params.ContentType != nil {
		s := string(*params.ContentType)
		contentType = &s
	}

	row, err := r.queries.UpdateArticle(ctx, adultdb.UpdateArticleParams{
		ID:                id,
		Name:              derefString(params.Name),
		Description:       params.Description,
		ContentType:       contentType,
		ConditionType:     derefConditionType(params.ConditionType),
		ConditionValue:    params.ConditionValue,
		Action:            derefAction(params.Action),
		AutomationTrigger: params.AutomationTrigger,
		Enabled:           derefBool(params.Enabled),
		Priority:          derefInt32(params.Priority),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return r.rowToArticle(&row), nil
}

func (r *SQLCRepository) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteArticle(ctx, id)
}

func (r *SQLCRepository) SetArticleEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*Article, error) {
	row, err := r.queries.SetArticleEnabled(ctx, adultdb.SetArticleEnabledParams{
		ID:      id,
		Enabled: enabled,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return r.rowToArticle(&row), nil
}

// ============================================================================
// CARGO HOLD (Global Quotas)
// ============================================================================

func (r *SQLCRepository) GetCargoHold(ctx context.Context) (*CargoHold, error) {
	row, err := r.queries.GetCargoHold(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCargoHoldNotFound
		}
		return nil, err
	}
	return r.rowToCargoHold(&row), nil
}

func (r *SQLCRepository) UpdateCargoHoldQuotas(ctx context.Context, totalQuota, expeditionQuota, voyageQuota float64) (*CargoHold, error) {
	row, err := r.queries.UpdateCargoHoldQuotas(ctx, adultdb.UpdateCargoHoldQuotasParams{
		TotalQuotaGb:      numericFromFloat(&totalQuota),
		ExpeditionQuotaGb: numericFromFloat(&expeditionQuota),
		VoyageQuotaGb:     numericFromFloat(&voyageQuota),
	})
	if err != nil {
		return nil, err
	}
	return r.rowToCargoHold(&row), nil
}

func (r *SQLCRepository) AddCargoHoldUsage(ctx context.Context, cargoGB float64, isExpedition bool) (*CargoHold, error) {
	// ExpeditionUsedGb is actually the boolean flag for whether it's an expedition
	// The SQL uses it in a CASE expression
	row, err := r.queries.AddCargoHoldUsage(ctx, adultdb.AddCargoHoldUsageParams{
		TotalUsedGb:      numericFromFloat(&cargoGB),
		ExpeditionUsedGb: boolToNumeric(isExpedition),
	})
	if err != nil {
		return nil, err
	}
	return r.rowToCargoHold(&row), nil
}

// ============================================================================
// CONVERTERS
// ============================================================================

func (r *SQLCRepository) rowToProvision(row *adultdb.QarProvision) *Provision {
	if row == nil {
		return nil
	}

	p := &Provision{
		ID:                    row.ID,
		UserID:                row.UserID,
		ContentType:           ContentType(row.ContentType),
		Status:                ProvisionStatus(row.Status),
		AutoApproved:          row.AutoApproved,
		Title:                 row.Title,
		Priority:              int(row.Priority),
		AyesCount:             int(row.AyesCount),
		TriggeredByAutomation: row.TriggeredByAutomation,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}

	if row.RequestSubtype != nil {
		p.RequestSubtype = RequestSubtype(*row.RequestSubtype)
	}
	if row.ExternalID != nil {
		p.ExternalID = *row.ExternalID
	}
	if row.ExternalSource != nil {
		p.ExternalSource = ExternalSource(*row.ExternalSource)
	}
	if row.ReleaseYear != nil {
		year := int(*row.ReleaseYear)
		p.ReleaseYear = &year
	}
	if row.Manifest != nil {
		p.Manifest = row.Manifest
	}
	if row.AutoArticleID.Valid {
		id := uuid.UUID(row.AutoArticleID.Bytes)
		p.AutoArticleID = &id
	}
	if row.ApprovedByUserID.Valid {
		id := uuid.UUID(row.ApprovedByUserID.Bytes)
		p.ApprovedByUserID = &id
	}
	if row.ApprovedAt.Valid {
		p.ApprovedAt = &row.ApprovedAt.Time
	}
	if row.DeclinedReason != nil {
		p.DeclinedReason = *row.DeclinedReason
	}
	if row.IntegrationID != nil {
		p.IntegrationID = *row.IntegrationID
	}
	if row.IntegrationStatus != nil {
		p.IntegrationStatus = *row.IntegrationStatus
	}
	if row.EstimatedCargoGb.Valid {
		f, _ := row.EstimatedCargoGb.Float64Value()
		p.EstimatedCargoGB = &f.Float64
	}
	if row.ActualCargoGb.Valid {
		f, _ := row.ActualCargoGb.Float64Value()
		p.ActualCargoGB = &f.Float64
	}
	if row.AvailableAt.Valid {
		p.AvailableAt = &row.AvailableAt.Time
	}
	if row.ParentProvisionID.Valid {
		id := uuid.UUID(row.ParentProvisionID.Bytes)
		p.ParentProvisionID = &id
	}

	return p
}

func (r *SQLCRepository) rowsToProvisions(rows []adultdb.QarProvision) []Provision {
	result := make([]Provision, 0, len(rows))
	for i := range rows {
		if p := r.rowToProvision(&rows[i]); p != nil {
			result = append(result, *p)
		}
	}
	return result
}

func (r *SQLCRepository) rowToMissive(row *adultdb.QarProvisionMissive) *ProvisionMissive {
	if row == nil {
		return nil
	}
	m := &ProvisionMissive{
		ID:             row.ID,
		ProvisionID:    row.ProvisionID,
		Message:        row.Message,
		IsCaptainOrder: row.IsCaptainOrder,
		CreatedAt:      row.CreatedAt,
	}
	if row.UserID.Valid {
		m.UserID = uuid.UUID(row.UserID.Bytes)
	}
	return m
}

func (r *SQLCRepository) rowToRation(row *adultdb.QarRation) *Ration {
	if row == nil {
		return nil
	}

	quotaGB, _ := row.CargoQuotaGb.Float64Value()
	usedGB, _ := row.CargoUsedGb.Float64Value()

	return &Ration{
		UserID:           row.UserID,
		DailyLimit:       int(row.DailyLimit),
		WeeklyLimit:      int(row.WeeklyLimit),
		MonthlyLimit:     int(row.MonthlyLimit),
		DailyUsed:        int(row.DailyUsed),
		WeeklyUsed:       int(row.WeeklyUsed),
		MonthlyUsed:      int(row.MonthlyUsed),
		CargoQuotaGB:     quotaGB.Float64,
		CargoUsedGB:      usedGB.Float64,
		LastResetDaily:   row.LastResetDaily.Time,
		LastResetWeekly:  row.LastResetWeekly.Time,
		LastResetMonthly: row.LastResetMonthly.Time,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

func (r *SQLCRepository) rowToArticle(row *adultdb.QarArticle) *Article {
	if row == nil {
		return nil
	}

	a := &Article{
		ID:            row.ID,
		Name:          row.Name,
		ConditionType: ArticleConditionType(row.ConditionType),
		Action:        ArticleAction(row.Action),
		Enabled:       row.Enabled,
		Priority:      int(row.Priority),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}

	if row.Description != nil {
		a.Description = *row.Description
	}
	if row.ContentType != nil {
		ct := ContentType(*row.ContentType)
		a.ContentType = &ct
	}
	if row.ConditionValue != nil {
		a.ConditionValue = row.ConditionValue
	}
	if row.AutomationTrigger != nil {
		a.AutomationTrigger = *row.AutomationTrigger
	}

	return a
}

func (r *SQLCRepository) rowsToArticles(rows []adultdb.QarArticle) []Article {
	result := make([]Article, 0, len(rows))
	for i := range rows {
		if a := r.rowToArticle(&rows[i]); a != nil {
			result = append(result, *a)
		}
	}
	return result
}

func (r *SQLCRepository) rowToCargoHold(row *adultdb.QarCargoHold) *CargoHold {
	if row == nil {
		return nil
	}

	totalQuota, _ := row.TotalQuotaGb.Float64Value()
	totalUsed, _ := row.TotalUsedGb.Float64Value()
	expQuota, _ := row.ExpeditionQuotaGb.Float64Value()
	expUsed, _ := row.ExpeditionUsedGb.Float64Value()
	voyQuota, _ := row.VoyageQuotaGb.Float64Value()
	voyUsed, _ := row.VoyageUsedGb.Float64Value()

	return &CargoHold{
		ID:                int(row.ID),
		TotalQuotaGB:      totalQuota.Float64,
		TotalUsedGB:       totalUsed.Float64,
		ExpeditionQuotaGB: expQuota.Float64,
		ExpeditionUsedGB:  expUsed.Float64,
		VoyageQuotaGB:     voyQuota.Float64,
		VoyageUsedGB:      voyUsed.Float64,
		UpdatedAt:         row.UpdatedAt,
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrInt32(i *int) *int32 {
	if i == nil {
		return nil
	}
	v := int32(*i)
	return &v
}

func derefInt32(i *int) int32 {
	if i == nil {
		return 0
	}
	return int32(*i)
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func derefConditionType(ct *ArticleConditionType) string {
	if ct == nil {
		return ""
	}
	return string(*ct)
}

func derefAction(a *ArticleAction) string {
	if a == nil {
		return ""
	}
	return string(*a)
}

// boolToNumeric converts a boolean to pgtype.Numeric (1 for true, 0 for false).
// Used for SQL queries that use boolean in CASE WHEN expressions.
func boolToNumeric(b bool) pgtype.Numeric {
	n := pgtype.Numeric{}
	if b {
		_ = n.Scan(1)
	} else {
		_ = n.Scan(0)
	}
	return n
}

func numericFromFloat(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{}
	}
	n := pgtype.Numeric{}
	_ = n.Scan(*f)
	return n
}

// marshalJSON is a helper for encoding JSON.
func marshalJSON(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
