// Package request provides QAR content request domain models.
package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service errors
var (
	ErrQuotaExceeded      = errors.New("request quota exceeded")
	ErrStorageQuotaExceeded = errors.New("storage quota exceeded")
	ErrAlreadyVoted       = errors.New("already voted on this provision")
	ErrNotVoted           = errors.New("not voted on this provision")
	ErrCannotModify       = errors.New("cannot modify provision in current status")
)

// Service provides provision (content request) business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new request service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "qar.request")),
	}
}

// ============================================================================
// PROVISIONS
// ============================================================================

// GetProvision retrieves a provision by ID.
func (s *Service) GetProvision(ctx context.Context, id uuid.UUID) (*Provision, error) {
	provision, err := s.repo.GetProvisionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get provision: %w", err)
	}
	return provision, nil
}

// ListProvisions retrieves provisions with filtering and pagination.
func (s *Service) ListProvisions(ctx context.Context, params ListProvisionsParams) ([]Provision, error) {
	provisions, err := s.repo.ListProvisions(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list provisions: %w", err)
	}
	return provisions, nil
}

// SearchProvisions searches provisions by title.
func (s *Service) SearchProvisions(ctx context.Context, query string, limit, offset int) ([]Provision, error) {
	provisions, err := s.repo.SearchProvisions(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search provisions: %w", err)
	}
	return provisions, nil
}

// CreateProvision creates a new content request.
func (s *Service) CreateProvision(ctx context.Context, params CreateProvisionParams) (*Provision, error) {
	// Check user quota
	ration, err := s.repo.UpsertRation(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user ration: %w", err)
	}

	// Reset quotas if needed
	if err := s.maybeResetRation(ctx, ration); err != nil {
		s.logger.Warn("failed to reset ration", slog.String("error", err.Error()))
	}

	// Check limits
	if ration.DailyUsed >= ration.DailyLimit {
		return nil, ErrQuotaExceeded
	}
	if ration.WeeklyUsed >= ration.WeeklyLimit {
		return nil, ErrQuotaExceeded
	}
	if ration.MonthlyUsed >= ration.MonthlyLimit {
		return nil, ErrQuotaExceeded
	}

	// Create provision
	provision, err := s.repo.CreateProvision(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create provision: %w", err)
	}

	// Increment usage
	if _, err := s.repo.IncrementRationUsage(ctx, params.UserID); err != nil {
		s.logger.Warn("failed to increment ration usage", slog.String("error", err.Error()))
	}

	// Check for auto-approval rules
	if article := s.matchArticle(ctx, provision); article != nil {
		if article.Action == ArticleActionAutoApprove {
			provision, err = s.repo.SetProvisionAutoApproved(ctx, provision.ID, article.ID)
			if err != nil {
				s.logger.Warn("failed to auto-approve provision", slog.String("error", err.Error()))
			} else {
				s.logger.Info("provision auto-approved",
					slog.String("provision_id", provision.ID.String()),
					slog.String("article_id", article.ID.String()),
					slog.String("article_name", article.Name),
				)
			}
		} else if article.Action == ArticleActionDecline {
			provision, err = s.repo.UpdateProvisionStatus(ctx, provision.ID, ProvisionStatusDeclined, nil, "Auto-declined by rule: "+article.Name)
			if err != nil {
				s.logger.Warn("failed to auto-decline provision", slog.String("error", err.Error()))
			}
		} else if article.Action == ArticleActionOnHold {
			provision, err = s.repo.UpdateProvisionStatus(ctx, provision.ID, ProvisionStatusOnHold, nil, "")
			if err != nil {
				s.logger.Warn("failed to put provision on hold", slog.String("error", err.Error()))
			}
		}
	}

	s.logger.Info("provision created",
		slog.String("id", provision.ID.String()),
		slog.String("title", provision.Title),
		slog.String("status", string(provision.Status)),
	)

	return provision, nil
}

// ApproveProvision approves a pending provision.
func (s *Service) ApproveProvision(ctx context.Context, provisionID, approverID uuid.UUID) (*Provision, error) {
	provision, err := s.repo.GetProvisionByID(ctx, provisionID)
	if err != nil {
		return nil, fmt.Errorf("get provision: %w", err)
	}

	if provision.Status != ProvisionStatusPending && provision.Status != ProvisionStatusOnHold {
		return nil, ErrCannotModify
	}

	provision, err = s.repo.UpdateProvisionStatus(ctx, provisionID, ProvisionStatusApproved, &approverID, "")
	if err != nil {
		return nil, fmt.Errorf("approve provision: %w", err)
	}

	s.logger.Info("provision approved",
		slog.String("provision_id", provisionID.String()),
		slog.String("approver_id", approverID.String()),
	)

	return provision, nil
}

// DeclineProvision declines a pending provision.
func (s *Service) DeclineProvision(ctx context.Context, provisionID, declinedByID uuid.UUID, reason string) (*Provision, error) {
	provision, err := s.repo.GetProvisionByID(ctx, provisionID)
	if err != nil {
		return nil, fmt.Errorf("get provision: %w", err)
	}

	if provision.Status != ProvisionStatusPending && provision.Status != ProvisionStatusOnHold {
		return nil, ErrCannotModify
	}

	provision, err = s.repo.UpdateProvisionStatus(ctx, provisionID, ProvisionStatusDeclined, &declinedByID, reason)
	if err != nil {
		return nil, fmt.Errorf("decline provision: %w", err)
	}

	s.logger.Info("provision declined",
		slog.String("provision_id", provisionID.String()),
		slog.String("reason", reason),
	)

	return provision, nil
}

// SetProvisionPriority updates the priority of a provision.
func (s *Service) SetProvisionPriority(ctx context.Context, provisionID uuid.UUID, priority int) (*Provision, error) {
	provision, err := s.repo.UpdateProvisionPriority(ctx, provisionID, priority)
	if err != nil {
		return nil, fmt.Errorf("set provision priority: %w", err)
	}

	s.logger.Info("provision priority updated",
		slog.String("provision_id", provisionID.String()),
		slog.Int("priority", priority),
	)

	return provision, nil
}

// SetProvisionAvailable marks a provision as available (content downloaded).
func (s *Service) SetProvisionAvailable(ctx context.Context, provisionID uuid.UUID, actualCargoGB float64) (*Provision, error) {
	provision, err := s.repo.SetProvisionAvailable(ctx, provisionID, actualCargoGB)
	if err != nil {
		return nil, fmt.Errorf("set provision available: %w", err)
	}

	// Update user cargo usage
	if _, err := s.repo.AddRationCargoUsage(ctx, provision.UserID, actualCargoGB); err != nil {
		s.logger.Warn("failed to add ration cargo usage", slog.String("error", err.Error()))
	}

	// Update global cargo usage
	isExpedition := provision.ContentType == ContentTypeExpedition
	if _, err := s.repo.AddCargoHoldUsage(ctx, actualCargoGB, isExpedition); err != nil {
		s.logger.Warn("failed to add cargo hold usage", slog.String("error", err.Error()))
	}

	s.logger.Info("provision available",
		slog.String("provision_id", provisionID.String()),
		slog.Float64("cargo_gb", actualCargoGB),
	)

	return provision, nil
}

// DeleteProvision deletes a provision.
func (s *Service) DeleteProvision(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteProvision(ctx, id); err != nil {
		return fmt.Errorf("delete provision: %w", err)
	}
	s.logger.Info("provision deleted", slog.String("id", id.String()))
	return nil
}

// ============================================================================
// AYES (Votes)
// ============================================================================

// VoteProvision adds a vote (aye) to a provision.
func (s *Service) VoteProvision(ctx context.Context, provisionID, userID uuid.UUID) error {
	// Check if already voted
	hasVoted, err := s.repo.HasUserVoted(ctx, provisionID, userID)
	if err != nil {
		return fmt.Errorf("check vote: %w", err)
	}
	if hasVoted {
		return ErrAlreadyVoted
	}

	// Create vote
	if _, err := s.repo.CreateAye(ctx, provisionID, userID); err != nil {
		return fmt.Errorf("create vote: %w", err)
	}

	s.logger.Info("provision voted",
		slog.String("provision_id", provisionID.String()),
		slog.String("user_id", userID.String()),
	)

	return nil
}

// UnvoteProvision removes a vote from a provision.
func (s *Service) UnvoteProvision(ctx context.Context, provisionID, userID uuid.UUID) error {
	// Check if voted
	hasVoted, err := s.repo.HasUserVoted(ctx, provisionID, userID)
	if err != nil {
		return fmt.Errorf("check vote: %w", err)
	}
	if !hasVoted {
		return ErrNotVoted
	}

	// Delete vote
	if err := s.repo.DeleteAye(ctx, provisionID, userID); err != nil {
		return fmt.Errorf("delete vote: %w", err)
	}

	s.logger.Info("provision unvoted",
		slog.String("provision_id", provisionID.String()),
		slog.String("user_id", userID.String()),
	)

	return nil
}

// HasUserVoted checks if a user has voted on a provision.
func (s *Service) HasUserVoted(ctx context.Context, provisionID, userID uuid.UUID) (bool, error) {
	return s.repo.HasUserVoted(ctx, provisionID, userID)
}

// ============================================================================
// MISSIVES (Comments)
// ============================================================================

// ListMissives retrieves comments for a provision.
func (s *Service) ListMissives(ctx context.Context, provisionID uuid.UUID) ([]ProvisionMissive, error) {
	missives, err := s.repo.ListMissives(ctx, provisionID)
	if err != nil {
		return nil, fmt.Errorf("list missives: %w", err)
	}
	return missives, nil
}

// CreateMissive creates a comment on a provision.
func (s *Service) CreateMissive(ctx context.Context, provisionID, userID uuid.UUID, message string, isCaptainOrder bool) (*ProvisionMissive, error) {
	missive, err := s.repo.CreateMissive(ctx, provisionID, userID, message, isCaptainOrder)
	if err != nil {
		return nil, fmt.Errorf("create missive: %w", err)
	}

	s.logger.Info("missive created",
		slog.String("provision_id", provisionID.String()),
		slog.String("user_id", userID.String()),
	)

	return missive, nil
}

// DeleteMissive deletes a comment.
func (s *Service) DeleteMissive(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteMissive(ctx, id); err != nil {
		return fmt.Errorf("delete missive: %w", err)
	}
	return nil
}

// ============================================================================
// RATIONS (User Quotas)
// ============================================================================

// GetRation retrieves a user's quota/ration.
func (s *Service) GetRation(ctx context.Context, userID uuid.UUID) (*Ration, error) {
	ration, err := s.repo.GetRation(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrRationNotFound) {
			// Create default ration
			return s.repo.UpsertRation(ctx, userID)
		}
		return nil, fmt.Errorf("get ration: %w", err)
	}
	return ration, nil
}

// UpdateRation updates a user's quota limits.
func (s *Service) UpdateRation(ctx context.Context, userID uuid.UUID, params UpdateRationParams) (*Ration, error) {
	// Ensure ration exists
	if _, err := s.repo.UpsertRation(ctx, userID); err != nil {
		return nil, fmt.Errorf("upsert ration: %w", err)
	}

	ration, err := s.repo.UpdateRationLimits(ctx, userID, params)
	if err != nil {
		return nil, fmt.Errorf("update ration: %w", err)
	}

	s.logger.Info("ration updated",
		slog.String("user_id", userID.String()),
	)

	return ration, nil
}

// ============================================================================
// ARTICLES (Rules)
// ============================================================================

// GetArticle retrieves a rule by ID.
func (s *Service) GetArticle(ctx context.Context, id uuid.UUID) (*Article, error) {
	article, err := s.repo.GetArticle(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get article: %w", err)
	}
	return article, nil
}

// ListArticles retrieves all rules.
func (s *Service) ListArticles(ctx context.Context) ([]Article, error) {
	articles, err := s.repo.ListArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	return articles, nil
}

// CreateArticle creates a new auto-approval rule.
func (s *Service) CreateArticle(ctx context.Context, params CreateArticleParams) (*Article, error) {
	article, err := s.repo.CreateArticle(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create article: %w", err)
	}

	s.logger.Info("article created",
		slog.String("id", article.ID.String()),
		slog.String("name", article.Name),
	)

	return article, nil
}

// UpdateArticle updates an existing rule.
func (s *Service) UpdateArticle(ctx context.Context, id uuid.UUID, params UpdateArticleParams) (*Article, error) {
	article, err := s.repo.UpdateArticle(ctx, id, params)
	if err != nil {
		return nil, fmt.Errorf("update article: %w", err)
	}

	s.logger.Info("article updated", slog.String("id", id.String()))

	return article, nil
}

// DeleteArticle deletes a rule.
func (s *Service) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteArticle(ctx, id); err != nil {
		return fmt.Errorf("delete article: %w", err)
	}

	s.logger.Info("article deleted", slog.String("id", id.String()))

	return nil
}

// ============================================================================
// CARGO HOLD (Global Quotas)
// ============================================================================

// GetCargoHold retrieves global storage quotas.
func (s *Service) GetCargoHold(ctx context.Context) (*CargoHold, error) {
	cargoHold, err := s.repo.GetCargoHold(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cargo hold: %w", err)
	}
	return cargoHold, nil
}

// UpdateCargoHoldQuotas updates global storage quotas.
func (s *Service) UpdateCargoHoldQuotas(ctx context.Context, totalQuota, expeditionQuota, voyageQuota float64) (*CargoHold, error) {
	cargoHold, err := s.repo.UpdateCargoHoldQuotas(ctx, totalQuota, expeditionQuota, voyageQuota)
	if err != nil {
		return nil, fmt.Errorf("update cargo hold quotas: %w", err)
	}

	s.logger.Info("cargo hold quotas updated",
		slog.Float64("total_quota_gb", totalQuota),
	)

	return cargoHold, nil
}

// ============================================================================
// PRIVATE HELPERS
// ============================================================================

// maybeResetRation checks if the ration counters need to be reset.
func (s *Service) maybeResetRation(ctx context.Context, ration *Ration) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Daily reset
	if ration.LastResetDaily.Before(today) {
		if err := s.repo.ResetDailyRations(ctx); err != nil {
			return err
		}
	}

	// Weekly reset (7 days)
	weekAgo := today.AddDate(0, 0, -7)
	if ration.LastResetWeekly.Before(weekAgo) {
		if err := s.repo.ResetWeeklyRations(ctx); err != nil {
			return err
		}
	}

	// Monthly reset (30 days)
	monthAgo := today.AddDate(0, 0, -30)
	if ration.LastResetMonthly.Before(monthAgo) {
		if err := s.repo.ResetMonthlyRations(ctx); err != nil {
			return err
		}
	}

	return nil
}

// matchArticle finds the first matching auto-approval rule for a provision.
func (s *Service) matchArticle(ctx context.Context, provision *Provision) *Article {
	articles, err := s.repo.ListEnabledArticles(ctx)
	if err != nil {
		s.logger.Warn("failed to list articles", slog.String("error", err.Error()))
		return nil
	}

	for _, article := range articles {
		// Check content type match
		if article.ContentType != nil && *article.ContentType != provision.ContentType {
			continue
		}

		// Evaluate condition
		if s.evaluateCondition(ctx, &article, provision) {
			return &article
		}
	}

	return nil
}

// evaluateCondition checks if a provision matches an article's condition.
func (s *Service) evaluateCondition(ctx context.Context, article *Article, provision *Provision) bool {
	var condValue ArticleConditionValue
	if err := json.Unmarshal(article.ConditionValue, &condValue); err != nil {
		s.logger.Warn("failed to unmarshal condition value", slog.String("error", err.Error()))
		return false
	}

	switch article.ConditionType {
	case ArticleConditionReleaseYear:
		if provision.ReleaseYear == nil {
			return false
		}
		year := *provision.ReleaseYear
		if condValue.MinYear != nil && year < *condValue.MinYear {
			return false
		}
		if condValue.MaxYear != nil && year > *condValue.MaxYear {
			return false
		}
		return true

	case ArticleConditionStorageAvailable:
		if condValue.MinFreeGB == nil {
			return false
		}
		cargoHold, err := s.repo.GetCargoHold(ctx)
		if err != nil {
			return false
		}
		return cargoHold.RemainingTotalGB() >= *condValue.MinFreeGB

	// Add more condition evaluators as needed
	default:
		return false
	}
}
