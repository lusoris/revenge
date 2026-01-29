//go:build ogen

package api

import (
	"context"
	"errors"

	gen "github.com/lusoris/revenge/api/generated"
)

func adultNotImplemented() *gen.Error {
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Adult content API not implemented yet",
	}
}

func adultModuleDisabled() *gen.Error {
	return &gen.Error{
		Code:    "module_disabled",
		Message: "Adult module disabled",
	}
}

func adultUnauthorized() *gen.Error {
	return &gen.Error{
		Code:    "unauthorized",
		Message: "Not authenticated",
	}
}

func adultForbidden() *gen.Error {
	return &gen.Error{
		Code:    "forbidden",
		Message: "Admin access required",
	}
}

func (h *Handler) CreateAdultMovie(ctx context.Context, req *gen.AdultMovieCreate) (gen.CreateAdultMovieRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) CreateAdultPerformer(ctx context.Context, req *gen.AdultPerformerCreate) (gen.CreateAdultPerformerRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) CreateAdultScene(ctx context.Context, req *gen.AdultSceneCreate) (gen.CreateAdultSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) CreateAdultStudio(ctx context.Context, req *gen.AdultStudioCreate) (gen.CreateAdultStudioRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) DeleteAdultMovie(ctx context.Context, params gen.DeleteAdultMovieParams) (gen.DeleteAdultMovieRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) DeleteAdultPerformer(ctx context.Context, params gen.DeleteAdultPerformerParams) (gen.DeleteAdultPerformerRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) DeleteAdultScene(ctx context.Context, params gen.DeleteAdultSceneParams) (gen.DeleteAdultSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) DeleteAdultStudio(ctx context.Context, params gen.DeleteAdultStudioParams) (gen.DeleteAdultStudioRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultMovie(ctx context.Context, params gen.GetAdultMovieParams) (gen.GetAdultMovieRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultPerformer(ctx context.Context, params gen.GetAdultPerformerParams) (gen.GetAdultPerformerRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultScene(ctx context.Context, params gen.GetAdultSceneParams) (gen.GetAdultSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultStudio(ctx context.Context, params gen.GetAdultStudioParams) (gen.GetAdultStudioRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultTag(ctx context.Context, params gen.GetAdultTagParams) (gen.GetAdultTagRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) IdentifyAdultScene(ctx context.Context, req *gen.AdultIdentifyRequest) (gen.IdentifyAdultSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultMovieMarkers(ctx context.Context, params gen.ListAdultMovieMarkersParams) (gen.ListAdultMovieMarkersRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultMoviePerformers(ctx context.Context, params gen.ListAdultMoviePerformersParams) (gen.ListAdultMoviePerformersRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultMovies(ctx context.Context, params gen.ListAdultMoviesParams) (gen.ListAdultMoviesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultPerformerMovies(ctx context.Context, params gen.ListAdultPerformerMoviesParams) (gen.ListAdultPerformerMoviesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultPerformers(ctx context.Context, params gen.ListAdultPerformersParams) (gen.ListAdultPerformersRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultScenes(ctx context.Context, params gen.ListAdultScenesParams) (gen.ListAdultScenesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultSimilarMovies(ctx context.Context, params gen.ListAdultSimilarMoviesParams) (gen.ListAdultSimilarMoviesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultStudioMovies(ctx context.Context, params gen.ListAdultStudioMoviesParams) (gen.ListAdultStudioMoviesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultStudios(ctx context.Context, params gen.ListAdultStudiosParams) (gen.ListAdultStudiosRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultTagMovies(ctx context.Context, params gen.ListAdultTagMoviesParams) (gen.ListAdultTagMoviesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultTags(ctx context.Context, params gen.ListAdultTagsParams) (gen.ListAdultTagsRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) MatchAdultScene(ctx context.Context, req *gen.AdultMatchRequest) (gen.MatchAdultSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) UpdateAdultMovie(ctx context.Context, req *gen.AdultMovieUpdate, params gen.UpdateAdultMovieParams) (gen.UpdateAdultMovieRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) UpdateAdultPerformer(ctx context.Context, req *gen.AdultPerformerUpdate, params gen.UpdateAdultPerformerParams) (gen.UpdateAdultPerformerRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) UpdateAdultScene(ctx context.Context, req *gen.AdultSceneUpdate, params gen.UpdateAdultSceneParams) (gen.UpdateAdultSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) UpdateAdultStudio(ctx context.Context, req *gen.AdultStudioUpdate, params gen.UpdateAdultStudioParams) (gen.UpdateAdultStudioRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) SearchAdultRequests(ctx context.Context, params gen.SearchAdultRequestsParams) (gen.SearchAdultRequestsRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultRequests(ctx context.Context, params gen.ListAdultRequestsParams) (gen.ListAdultRequestsRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) CreateAdultRequest(ctx context.Context, req *gen.AdultRequestCreate) (gen.CreateAdultRequestRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultRequest(ctx context.Context, params gen.GetAdultRequestParams) (gen.GetAdultRequestRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) VoteAdultRequest(ctx context.Context, req *gen.AdultRequestVoteCreate, params gen.VoteAdultRequestParams) (gen.VoteAdultRequestRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) CommentAdultRequest(ctx context.Context, req *gen.AdultRequestCommentCreate, params gen.CommentAdultRequestParams) (gen.CommentAdultRequestRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ListAdultAdminRequests(ctx context.Context, params gen.ListAdultAdminRequestsParams) (gen.ListAdultAdminRequestsRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.ListAdultAdminRequestsUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.ListAdultAdminRequestsForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.ListAdultAdminRequestsForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.ListAdultAdminRequestsForbidden)(adultNotImplemented()), nil
}

func (h *Handler) ApproveAdultRequest(ctx context.Context, params gen.ApproveAdultRequestParams) (gen.ApproveAdultRequestRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.ApproveAdultRequestUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.ApproveAdultRequestForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.ApproveAdultRequestForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.ApproveAdultRequestForbidden)(adultNotImplemented()), nil
}

func (h *Handler) DeclineAdultRequest(ctx context.Context, params gen.DeclineAdultRequestParams) (gen.DeclineAdultRequestRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.DeclineAdultRequestUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.DeclineAdultRequestForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.DeclineAdultRequestForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.DeclineAdultRequestForbidden)(adultNotImplemented()), nil
}

func (h *Handler) UpdateAdultRequestQuota(ctx context.Context, req *gen.AdultRequestQuotaUpdate, params gen.UpdateAdultRequestQuotaParams) (gen.UpdateAdultRequestQuotaRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.UpdateAdultRequestQuotaUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.UpdateAdultRequestQuotaForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.UpdateAdultRequestQuotaForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.UpdateAdultRequestQuotaForbidden)(adultNotImplemented()), nil
}

func (h *Handler) ListAdultRequestRules(ctx context.Context) (gen.ListAdultRequestRulesRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.ListAdultRequestRulesUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.ListAdultRequestRulesForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.ListAdultRequestRulesForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.ListAdultRequestRulesForbidden)(adultNotImplemented()), nil
}

func (h *Handler) CreateAdultRequestRule(ctx context.Context, req *gen.AdultRequestRuleCreate) (gen.CreateAdultRequestRuleRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.CreateAdultRequestRuleUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.CreateAdultRequestRuleForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.CreateAdultRequestRuleForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.CreateAdultRequestRuleForbidden)(adultNotImplemented()), nil
}

func (h *Handler) UpdateAdultRequestRule(ctx context.Context, req *gen.AdultRequestRuleUpdate, params gen.UpdateAdultRequestRuleParams) (gen.UpdateAdultRequestRuleRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.UpdateAdultRequestRuleUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.UpdateAdultRequestRuleForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.UpdateAdultRequestRuleForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.UpdateAdultRequestRuleForbidden)(adultNotImplemented()), nil
}

func (h *Handler) DeleteAdultRequestRule(ctx context.Context, params gen.DeleteAdultRequestRuleParams) (gen.DeleteAdultRequestRuleRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return (*gen.DeleteAdultRequestRuleUnauthorized)(adultUnauthorized()), nil
		}
		return (*gen.DeleteAdultRequestRuleForbidden)(adultForbidden()), nil
	}
	if !h.adultEnabled {
		return (*gen.DeleteAdultRequestRuleForbidden)(adultModuleDisabled()), nil
	}
	return (*gen.DeleteAdultRequestRuleForbidden)(adultNotImplemented()), nil
}

func (h *Handler) SearchAdultStashDBScenes(ctx context.Context, params gen.SearchAdultStashDBScenesParams) (gen.SearchAdultStashDBScenesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultStashDBScene(ctx context.Context, params gen.GetAdultStashDBSceneParams) (gen.GetAdultStashDBSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) SearchAdultStashDBPerformers(ctx context.Context, params gen.SearchAdultStashDBPerformersParams) (gen.SearchAdultStashDBPerformersRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultStashDBPerformer(ctx context.Context, params gen.GetAdultStashDBPerformerParams) (gen.GetAdultStashDBPerformerRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) IdentifyAdultStashDBScene(ctx context.Context, req *gen.AdultMatchRequest) (gen.IdentifyAdultStashDBSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) SearchAdultTPDBScenes(ctx context.Context, params gen.SearchAdultTPDBScenesParams) (gen.SearchAdultTPDBScenesRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultTPDBScene(ctx context.Context, params gen.GetAdultTPDBSceneParams) (gen.GetAdultTPDBSceneRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultTPDBPerformer(ctx context.Context, params gen.GetAdultTPDBPerformerParams) (gen.GetAdultTPDBPerformerRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) SyncAdultStash(ctx context.Context) (gen.SyncAdultStashRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) ImportAdultStash(ctx context.Context) (gen.ImportAdultStashRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}

func (h *Handler) GetAdultStashStatus(ctx context.Context) (gen.GetAdultStashStatusRes, error) {
	if _, err := requireUser(ctx); err != nil {
		return adultUnauthorized(), nil
	}
	if !h.adultEnabled {
		return adultModuleDisabled(), nil
	}
	return adultNotImplemented(), nil
}
