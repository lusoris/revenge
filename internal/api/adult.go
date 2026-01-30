package api

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/content/qar/crew"
	"github.com/lusoris/revenge/internal/content/qar/expedition"
	"github.com/lusoris/revenge/internal/content/qar/flag"
	"github.com/lusoris/revenge/internal/content/qar/port"
	"github.com/lusoris/revenge/internal/content/qar/voyage"
	"github.com/lusoris/revenge/internal/content/shared"
)

// Adult content handlers.
// These implement the QAR (Queen Anne's Revenge) API endpoints.
// Internal services use QAR terminology (expedition, voyage, crew, port, flag).
// API uses user-facing terminology (movie, scene, performer, studio, tag).

// requireExpeditionService returns the expedition service or an error if nil.
func (h *Handler) requireExpeditionService() (*expedition.Service, error) {
	if h.expeditionService == nil {
		return nil, ErrModuleDisabled
	}
	return h.expeditionService, nil
}

// requireVoyageService returns the voyage service or an error if nil.
func (h *Handler) requireVoyageService() (*voyage.Service, error) {
	if h.voyageService == nil {
		return nil, ErrModuleDisabled
	}
	return h.voyageService, nil
}

// requireCrewService returns the crew service or an error if nil.
func (h *Handler) requireCrewService() (*crew.Service, error) {
	if h.crewService == nil {
		return nil, ErrModuleDisabled
	}
	return h.crewService, nil
}

// requirePortService returns the port service or an error if nil.
func (h *Handler) requirePortService() (*port.Service, error) {
	if h.portService == nil {
		return nil, ErrModuleDisabled
	}
	return h.portService, nil
}

// requireFlagService returns the flag service or an error if nil.
func (h *Handler) requireFlagService() (*flag.Service, error) {
	if h.flagService == nil {
		return nil, ErrModuleDisabled
	}
	return h.flagService, nil
}

// adultModuleDisabledError returns a generic error for adult endpoints.
func adultModuleDisabledError() *gen.Error {
	return &gen.Error{
		Code:    "module_disabled",
		Message: "Adult content module is not enabled",
	}
}

// =============================================================================
// Expedition (Adult Movie) Handlers
// =============================================================================

func (h *Handler) GetAdultMovie(ctx context.Context, params gen.GetAdultMovieParams) (gen.GetAdultMovieRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.GetAdultMovieUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.GetAdultMovieUnauthorized{
			Code:    "forbidden",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return &gen.GetAdultMovieNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	exp, err := svc.GetByID(ctx, params.MovieId)
	if err != nil {
		return &gen.GetAdultMovieNotFound{
			Code:    "not_found",
			Message: "Adult movie not found",
		}, nil
	}

	result := expeditionToAPI(exp)
	return &result, nil
}

func (h *Handler) ListAdultMovies(ctx context.Context, params gen.ListAdultMoviesParams) (gen.ListAdultMoviesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	limit := params.Limit.Or(20)
	offset := params.Offset.Or(0)

	var expeditions []expedition.Expedition
	if params.Query.IsSet() && params.Query.Value != "" {
		// Search by query
		expeditions, err = svc.Search(ctx, params.Query.Value, limit, offset)
	} else if params.LibraryId.IsSet() {
		// Filter by library
		expeditions, err = svc.ListByFleet(ctx, params.LibraryId.Value, limit, offset)
	} else {
		// List all
		expeditions, err = svc.List(ctx, limit, offset)
	}

	if err != nil {
		h.logger.Error("List adult movies failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list adult movies",
		}, nil
	}

	result := gen.AdultMovieListResponse{
		Movies:     make([]gen.AdultMovie, len(expeditions)),
		Pagination: paginationMeta(int64(len(expeditions)), limit, offset),
	}

	for i, e := range expeditions {
		result.Movies[i] = expeditionToAPI(&e)
	}

	return &result, nil
}

func (h *Handler) CreateAdultMovie(ctx context.Context, req *gen.AdultMovieCreate) (gen.CreateAdultMovieRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		return &gen.Error{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	exp := &expedition.Expedition{
		ContentEntity: shared.ContentEntity{
			BaseEntity: shared.BaseEntity{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Path:  req.Path,
			Title: req.Title,
		},
		FleetID: req.LibraryId,
	}

	if req.ReleaseDate.IsSet() {
		t := time.Time(req.ReleaseDate.Value)
		exp.LaunchDate = &t
	}
	if req.StudioId.IsSet() {
		exp.PortID = &req.StudioId.Value
	}

	if err := svc.Create(ctx, exp); err != nil {
		h.logger.Error("Create adult movie failed", "error", err)
		return &gen.Error{
			Code:    "create_failed",
			Message: "Failed to create adult movie",
		}, nil
	}

	result := expeditionToAPI(exp)
	return &result, nil
}

func (h *Handler) UpdateAdultMovie(ctx context.Context, req *gen.AdultMovieUpdate, params gen.UpdateAdultMovieParams) (gen.UpdateAdultMovieRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.UpdateAdultMovieUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.UpdateAdultMovieUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return &gen.UpdateAdultMovieNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	exp, err := svc.GetByID(ctx, params.MovieId)
	if err != nil {
		return &gen.UpdateAdultMovieNotFound{
			Code:    "not_found",
			Message: "Adult movie not found",
		}, nil
	}

	// Apply updates
	if req.Title.IsSet() {
		exp.Title = req.Title.Value
	}
	if req.SortTitle.IsSet() {
		exp.SortTitle = req.SortTitle.Value
	}
	if req.Overview.IsSet() {
		exp.Overview = req.Overview.Value
	}
	if req.ReleaseDate.IsSet() {
		t := time.Time(req.ReleaseDate.Value)
		exp.LaunchDate = &t
	}
	if req.StudioId.IsSet() {
		exp.PortID = &req.StudioId.Value
	}
	exp.UpdatedAt = time.Now()

	if err := svc.Update(ctx, exp); err != nil {
		h.logger.Error("Update adult movie failed", "error", err, "movie_id", params.MovieId)
		return &gen.UpdateAdultMovieNotFound{
			Code:    "update_failed",
			Message: "Failed to update adult movie",
		}, nil
	}

	result := expeditionToAPI(exp)
	return &result, nil
}

func (h *Handler) DeleteAdultMovie(ctx context.Context, params gen.DeleteAdultMovieParams) (gen.DeleteAdultMovieRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteAdultMovieUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteAdultMovieUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return &gen.DeleteAdultMovieNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	if err := svc.Delete(ctx, params.MovieId); err != nil {
		return &gen.DeleteAdultMovieNotFound{
			Code:    "not_found",
			Message: "Adult movie not found",
		}, nil
	}

	return &gen.DeleteAdultMovieNoContent{}, nil
}

// =============================================================================
// Voyage (Adult Scene) Handlers
// =============================================================================

func (h *Handler) GetAdultScene(ctx context.Context, params gen.GetAdultSceneParams) (gen.GetAdultSceneRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.GetAdultSceneUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.GetAdultSceneUnauthorized{
			Code:    "forbidden",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireVoyageService()
	if err != nil {
		return &gen.GetAdultSceneNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	voy, err := svc.GetByID(ctx, params.SceneId)
	if err != nil {
		return &gen.GetAdultSceneNotFound{
			Code:    "not_found",
			Message: "Adult scene not found",
		}, nil
	}

	result := voyageToAPI(voy)
	return &result, nil
}

func (h *Handler) ListAdultScenes(ctx context.Context, params gen.ListAdultScenesParams) (gen.ListAdultScenesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireVoyageService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	limit := params.Limit.Or(20)
	offset := params.Offset.Or(0)

	var voyages []voyage.Voyage
	if params.Query.IsSet() && params.Query.Value != "" {
		// Search by query
		voyages, err = svc.Search(ctx, params.Query.Value, limit, offset)
	} else if params.LibraryId.IsSet() {
		// Filter by library
		voyages, err = svc.ListByFleet(ctx, params.LibraryId.Value, limit, offset)
	} else {
		// List all
		voyages, err = svc.List(ctx, limit, offset)
	}

	if err != nil {
		h.logger.Error("List adult scenes failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list adult scenes",
		}, nil
	}

	result := gen.AdultSceneListResponse{
		Scenes:     make([]gen.AdultScene, len(voyages)),
		Pagination: paginationMeta(int64(len(voyages)), limit, offset),
	}

	for i, v := range voyages {
		result.Scenes[i] = voyageToAPI(&v)
	}

	return &result, nil
}

func (h *Handler) CreateAdultScene(ctx context.Context, req *gen.AdultSceneCreate) (gen.CreateAdultSceneRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		return &gen.Error{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireVoyageService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	voy := &voyage.Voyage{
		ContentEntity: shared.ContentEntity{
			BaseEntity: shared.BaseEntity{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Path:  req.Path,
			Title: req.Title,
		},
		FleetID: req.LibraryId,
	}

	if req.ReleaseDate.IsSet() {
		t := time.Time(req.ReleaseDate.Value)
		voy.LaunchDate = &t
	}
	if req.StudioId.IsSet() {
		voy.PortID = &req.StudioId.Value
	}

	if err := svc.Create(ctx, voy); err != nil {
		h.logger.Error("Create adult scene failed", "error", err)
		return &gen.Error{
			Code:    "create_failed",
			Message: "Failed to create adult scene",
		}, nil
	}

	result := voyageToAPI(voy)
	return &result, nil
}

func (h *Handler) UpdateAdultScene(ctx context.Context, req *gen.AdultSceneUpdate, params gen.UpdateAdultSceneParams) (gen.UpdateAdultSceneRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.UpdateAdultSceneUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.UpdateAdultSceneUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireVoyageService()
	if err != nil {
		return &gen.UpdateAdultSceneNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	voy, err := svc.GetByID(ctx, params.SceneId)
	if err != nil {
		return &gen.UpdateAdultSceneNotFound{
			Code:    "not_found",
			Message: "Adult scene not found",
		}, nil
	}

	// Apply updates
	if req.Title.IsSet() {
		voy.Title = req.Title.Value
	}
	if req.SortTitle.IsSet() {
		voy.SortTitle = req.SortTitle.Value
	}
	if req.Overview.IsSet() {
		voy.Overview = req.Overview.Value
	}
	if req.ReleaseDate.IsSet() {
		t := time.Time(req.ReleaseDate.Value)
		voy.LaunchDate = &t
	}
	if req.StudioId.IsSet() {
		voy.PortID = &req.StudioId.Value
	}
	voy.UpdatedAt = time.Now()

	if err := svc.Update(ctx, voy); err != nil {
		h.logger.Error("Update adult scene failed", "error", err, "scene_id", params.SceneId)
		return &gen.UpdateAdultSceneNotFound{
			Code:    "update_failed",
			Message: "Failed to update adult scene",
		}, nil
	}

	result := voyageToAPI(voy)
	return &result, nil
}

func (h *Handler) DeleteAdultScene(ctx context.Context, params gen.DeleteAdultSceneParams) (gen.DeleteAdultSceneRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteAdultSceneUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteAdultSceneUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireVoyageService()
	if err != nil {
		return &gen.DeleteAdultSceneNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	if err := svc.Delete(ctx, params.SceneId); err != nil {
		return &gen.DeleteAdultSceneNotFound{
			Code:    "not_found",
			Message: "Adult scene not found",
		}, nil
	}

	return &gen.DeleteAdultSceneNoContent{}, nil
}

// =============================================================================
// Crew (Adult Performer) Handlers
// =============================================================================

func (h *Handler) GetAdultPerformer(ctx context.Context, params gen.GetAdultPerformerParams) (gen.GetAdultPerformerRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.GetAdultPerformerUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.GetAdultPerformerUnauthorized{
			Code:    "forbidden",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireCrewService()
	if err != nil {
		return &gen.GetAdultPerformerNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	c, err := svc.GetByID(ctx, params.PerformerId)
	if err != nil {
		return &gen.GetAdultPerformerNotFound{
			Code:    "not_found",
			Message: "Adult performer not found",
		}, nil
	}

	result := crewToAPI(c)
	return &result, nil
}

func (h *Handler) ListAdultPerformers(ctx context.Context, params gen.ListAdultPerformersParams) (gen.ListAdultPerformersRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireCrewService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	limit := params.Limit.Or(20)
	offset := params.Offset.Or(0)

	var members []crew.Crew
	if params.Query.IsSet() && params.Query.Value != "" {
		// Search by query
		members, err = svc.Search(ctx, params.Query.Value, limit, offset)
	} else {
		// List all
		members, err = svc.List(ctx, limit, offset)
	}

	if err != nil {
		h.logger.Error("List adult performers failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list adult performers",
		}, nil
	}

	result := gen.AdultPerformerListResponse{
		Performers: make([]gen.AdultPerformer, len(members)),
		Pagination: paginationMeta(int64(len(members)), limit, offset),
	}

	for i, c := range members {
		result.Performers[i] = crewToAPI(&c)
	}

	return &result, nil
}

func (h *Handler) CreateAdultPerformer(ctx context.Context, req *gen.AdultPerformerCreate) (gen.CreateAdultPerformerRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		return &gen.Error{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireCrewService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	c := &crew.Crew{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Gender.IsSet() {
		c.Gender = req.Gender.Value
	}
	if req.Disambiguation.IsSet() {
		c.Disambiguation = req.Disambiguation.Value
	}

	if err := svc.Create(ctx, c); err != nil {
		h.logger.Error("Create adult performer failed", "error", err)
		return &gen.Error{
			Code:    "create_failed",
			Message: "Failed to create adult performer",
		}, nil
	}

	result := crewToAPI(c)
	return &result, nil
}

func (h *Handler) UpdateAdultPerformer(ctx context.Context, req *gen.AdultPerformerUpdate, params gen.UpdateAdultPerformerParams) (gen.UpdateAdultPerformerRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.UpdateAdultPerformerUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.UpdateAdultPerformerUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireCrewService()
	if err != nil {
		return &gen.UpdateAdultPerformerNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	c, err := svc.GetByID(ctx, params.PerformerId)
	if err != nil {
		return &gen.UpdateAdultPerformerNotFound{
			Code:    "not_found",
			Message: "Adult performer not found",
		}, nil
	}

	// Apply updates
	if req.Name.IsSet() {
		c.Name = req.Name.Value
	}
	if req.Gender.IsSet() {
		c.Gender = req.Gender.Value
	}
	if req.Disambiguation.IsSet() {
		c.Disambiguation = req.Disambiguation.Value
	}
	c.UpdatedAt = time.Now()

	if err := svc.Update(ctx, c); err != nil {
		h.logger.Error("Update adult performer failed", "error", err, "performer_id", params.PerformerId)
		return &gen.UpdateAdultPerformerNotFound{
			Code:    "update_failed",
			Message: "Failed to update adult performer",
		}, nil
	}

	result := crewToAPI(c)
	return &result, nil
}

func (h *Handler) DeleteAdultPerformer(ctx context.Context, params gen.DeleteAdultPerformerParams) (gen.DeleteAdultPerformerRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteAdultPerformerUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteAdultPerformerUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requireCrewService()
	if err != nil {
		return &gen.DeleteAdultPerformerNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	if err := svc.Delete(ctx, params.PerformerId); err != nil {
		return &gen.DeleteAdultPerformerNotFound{
			Code:    "not_found",
			Message: "Adult performer not found",
		}, nil
	}

	return &gen.DeleteAdultPerformerNoContent{}, nil
}

// =============================================================================
// Port (Adult Studio) Handlers
// =============================================================================

func (h *Handler) GetAdultStudio(ctx context.Context, params gen.GetAdultStudioParams) (gen.GetAdultStudioRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.GetAdultStudioUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.GetAdultStudioUnauthorized{
			Code:    "forbidden",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requirePortService()
	if err != nil {
		return &gen.GetAdultStudioNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	p, err := svc.GetByID(ctx, params.StudioId)
	if err != nil {
		return &gen.GetAdultStudioNotFound{
			Code:    "not_found",
			Message: "Adult studio not found",
		}, nil
	}

	result := portToAPI(p)
	return &result, nil
}

func (h *Handler) ListAdultStudios(ctx context.Context, params gen.ListAdultStudiosParams) (gen.ListAdultStudiosRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requirePortService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	limit := params.Limit.Or(20)
	offset := params.Offset.Or(0)

	var ports []port.Port
	if params.Query.IsSet() && params.Query.Value != "" {
		// Search by query
		ports, err = svc.Search(ctx, params.Query.Value, limit, offset)
	} else {
		// List all
		ports, err = svc.List(ctx, limit, offset)
	}

	if err != nil {
		h.logger.Error("List adult studios failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list adult studios",
		}, nil
	}

	result := gen.AdultStudioListResponse{
		Studios:    make([]gen.AdultStudio, len(ports)),
		Pagination: paginationMeta(int64(len(ports)), limit, offset),
	}

	for i, p := range ports {
		result.Studios[i] = portToAPI(&p)
	}

	return &result, nil
}

func (h *Handler) CreateAdultStudio(ctx context.Context, req *gen.AdultStudioCreate) (gen.CreateAdultStudioRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		return &gen.Error{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requirePortService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	p := &port.Port{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.ParentId.IsSet() {
		p.ParentID = &req.ParentId.Value
	}
	if req.URL.IsSet() {
		p.URL = req.URL.Value
	}

	if err := svc.Create(ctx, p); err != nil {
		h.logger.Error("Create adult studio failed", "error", err)
		return &gen.Error{
			Code:    "create_failed",
			Message: "Failed to create adult studio",
		}, nil
	}

	result := portToAPI(p)
	return &result, nil
}

func (h *Handler) UpdateAdultStudio(ctx context.Context, req *gen.AdultStudioUpdate, params gen.UpdateAdultStudioParams) (gen.UpdateAdultStudioRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.UpdateAdultStudioUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.UpdateAdultStudioUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requirePortService()
	if err != nil {
		return &gen.UpdateAdultStudioNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	p, err := svc.GetByID(ctx, params.StudioId)
	if err != nil {
		return &gen.UpdateAdultStudioNotFound{
			Code:    "not_found",
			Message: "Adult studio not found",
		}, nil
	}

	// Apply updates
	if req.Name.IsSet() {
		p.Name = req.Name.Value
	}
	if req.ParentId.IsSet() {
		p.ParentID = &req.ParentId.Value
	}
	if req.URL.IsSet() {
		p.URL = req.URL.Value
	}
	p.UpdatedAt = time.Now()

	if err := svc.Update(ctx, p); err != nil {
		h.logger.Error("Update adult studio failed", "error", err, "studio_id", params.StudioId)
		return &gen.UpdateAdultStudioNotFound{
			Code:    "update_failed",
			Message: "Failed to update adult studio",
		}, nil
	}

	result := portToAPI(p)
	return &result, nil
}

func (h *Handler) DeleteAdultStudio(ctx context.Context, params gen.DeleteAdultStudioParams) (gen.DeleteAdultStudioRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.DeleteAdultStudioUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.DeleteAdultStudioUnauthorized{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}

	svc, err := h.requirePortService()
	if err != nil {
		return &gen.DeleteAdultStudioNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	if err := svc.Delete(ctx, params.StudioId); err != nil {
		return &gen.DeleteAdultStudioNotFound{
			Code:    "not_found",
			Message: "Adult studio not found",
		}, nil
	}

	return &gen.DeleteAdultStudioNoContent{}, nil
}

// =============================================================================
// Flag (Adult Tag) Handlers
// =============================================================================

func (h *Handler) GetAdultTag(ctx context.Context, params gen.GetAdultTagParams) (gen.GetAdultTagRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return &gen.GetAdultTagUnauthorized{
				Code:    "unauthorized",
				Message: "Not authenticated",
			}, nil
		}
		return &gen.GetAdultTagUnauthorized{
			Code:    "forbidden",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireFlagService()
	if err != nil {
		return &gen.GetAdultTagNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	f, err := svc.GetByID(ctx, params.TagId)
	if err != nil {
		return &gen.GetAdultTagNotFound{
			Code:    "not_found",
			Message: "Adult tag not found",
		}, nil
	}

	result := flagToAPI(f)
	return &result, nil
}

func (h *Handler) ListAdultTags(ctx context.Context, params gen.ListAdultTagsParams) (gen.ListAdultTagsRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireFlagService()
	if err != nil {
		return adultModuleDisabledError(), nil
	}

	limit := params.Limit.Or(20)
	offset := params.Offset.Or(0)

	var flags []flag.Flag
	if params.Query.IsSet() && params.Query.Value != "" {
		// Search by query
		flags, err = svc.Search(ctx, params.Query.Value, limit, offset)
	} else {
		// List all
		flags, err = svc.List(ctx, limit, offset)
	}

	if err != nil {
		h.logger.Error("List adult tags failed", "error", err)
		return &gen.Error{
			Code:    "list_failed",
			Message: "Failed to list adult tags",
		}, nil
	}

	result := gen.AdultTagListResponse{
		Tags:       make([]gen.AdultTag, len(flags)),
		Pagination: paginationMeta(int64(len(flags)), limit, offset),
	}

	for i, f := range flags {
		result.Tags[i] = flagToAPI(&f)
	}

	return &result, nil
}

// =============================================================================
// Relationship Handlers
// =============================================================================

func (h *Handler) ListAdultMoviePerformers(ctx context.Context, params gen.ListAdultMoviePerformersParams) (gen.ListAdultMoviePerformersRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.ListAdultMoviePerformersUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireCrewService()
	if err != nil {
		return &gen.ListAdultMoviePerformersNotFound{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	members, err := svc.ListExpeditionCrew(ctx, params.MovieId)
	if err != nil {
		h.logger.Error("List movie performers failed", "error", err, "movie_id", params.MovieId)
		return &gen.ListAdultMoviePerformersNotFound{
			Code:    "list_failed",
			Message: "Failed to list performers",
		}, nil
	}

	result := make(gen.ListAdultMoviePerformersOKApplicationJSON, len(members))
	for i, c := range members {
		result[i] = crewToAPI(&c)
	}

	return &result, nil
}

func (h *Handler) ListAdultPerformerMovies(ctx context.Context, params gen.ListAdultPerformerMoviesParams) (gen.ListAdultPerformerMoviesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.ListAdultPerformerMoviesUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return &gen.ListAdultPerformerMoviesNotFound{
			Code:    "module_disabled",
			Message: "Adult module is not enabled",
		}, nil
	}

	limit := 100
	offset := 0

	expeditions, total, err := svc.ListByPerformer(ctx, params.PerformerId, limit, offset)
	if err != nil {
		h.logger.Error("List performer movies failed", "error", err, "performerId", params.PerformerId)
		return &gen.ListAdultPerformerMoviesNotFound{
			Code:    "list_failed",
			Message: "Failed to list performer movies",
		}, nil
	}

	result := gen.AdultMovieListResponse{
		Movies:     make([]gen.AdultMovie, len(expeditions)),
		Pagination: paginationMeta(total, limit, offset),
	}

	for i, e := range expeditions {
		result.Movies[i] = expeditionToAPI(&e)
	}

	return &result, nil
}

func (h *Handler) ListAdultStudioMovies(ctx context.Context, params gen.ListAdultStudioMoviesParams) (gen.ListAdultStudioMoviesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.ListAdultStudioMoviesUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return &gen.ListAdultStudioMoviesNotFound{
			Code:    "module_disabled",
			Message: "Adult module is not enabled",
		}, nil
	}

	limit := 100
	offset := 0

	expeditions, total, err := svc.ListByStudio(ctx, params.StudioId, limit, offset)
	if err != nil {
		h.logger.Error("List studio movies failed", "error", err, "studioId", params.StudioId)
		return &gen.ListAdultStudioMoviesNotFound{
			Code:    "list_failed",
			Message: "Failed to list studio movies",
		}, nil
	}

	result := gen.AdultMovieListResponse{
		Movies:     make([]gen.AdultMovie, len(expeditions)),
		Pagination: paginationMeta(total, limit, offset),
	}

	for i, e := range expeditions {
		result.Movies[i] = expeditionToAPI(&e)
	}

	return &result, nil
}

func (h *Handler) ListAdultTagMovies(ctx context.Context, params gen.ListAdultTagMoviesParams) (gen.ListAdultTagMoviesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.ListAdultTagMoviesUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireExpeditionService()
	if err != nil {
		return &gen.ListAdultTagMoviesNotFound{
			Code:    "module_disabled",
			Message: "Adult module is not enabled",
		}, nil
	}

	limit := 100
	offset := 0

	expeditions, total, err := svc.ListByTag(ctx, params.TagId, limit, offset)
	if err != nil {
		h.logger.Error("List tag movies failed", "error", err, "tagId", params.TagId)
		return &gen.ListAdultTagMoviesNotFound{
			Code:    "list_failed",
			Message: "Failed to list tag movies",
		}, nil
	}

	result := gen.AdultMovieListResponse{
		Movies:     make([]gen.AdultMovie, len(expeditions)),
		Pagination: paginationMeta(total, limit, offset),
	}

	for i, e := range expeditions {
		result.Movies[i] = expeditionToAPI(&e)
	}

	return &result, nil
}

func (h *Handler) ListAdultSimilarMovies(ctx context.Context, params gen.ListAdultSimilarMoviesParams) (gen.ListAdultSimilarMoviesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.ListAdultSimilarMoviesUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	// TODO: Implement similar movie recommendations
	return &gen.ListAdultSimilarMoviesNotFound{
		Code:    "not_implemented",
		Message: "Similar movies not yet implemented",
	}, nil
}

func (h *Handler) ListAdultMovieMarkers(ctx context.Context, params gen.ListAdultMovieMarkersParams) (gen.ListAdultMovieMarkersRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.ListAdultMovieMarkersUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	// TODO: Implement markers/chapters for adult movies
	return &gen.ListAdultMovieMarkersNotFound{
		Code:    "not_implemented",
		Message: "Movie markers not yet implemented",
	}, nil
}

// =============================================================================
// Fingerprinting & Matching Handlers
// =============================================================================

func (h *Handler) IdentifyAdultScene(ctx context.Context, req *gen.AdultIdentifyRequest) (gen.IdentifyAdultSceneRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireVoyageService()
	if err != nil {
		return &gen.Error{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	oshash := ""
	phash := ""
	if req.Oshash.IsSet() {
		oshash = req.Oshash.Value
	}
	if req.Phash.IsSet() {
		phash = req.Phash.Value
	}

	voy, err := svc.MatchByFingerprint(ctx, oshash, phash)
	if err != nil {
		// Return empty candidates when no match
		return &gen.AdultIdentifyResponse{
			Candidates: []gen.AdultScene{},
		}, nil
	}

	scene := voyageToAPI(voy)
	return &gen.AdultIdentifyResponse{
		Candidates: []gen.AdultScene{scene},
	}, nil
}

func (h *Handler) MatchAdultScene(ctx context.Context, req *gen.AdultMatchRequest) (gen.MatchAdultSceneRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}

	svc, err := h.requireVoyageService()
	if err != nil {
		return &gen.Error{
			Code:    "module_disabled",
			Message: "Adult content module is not enabled",
		}, nil
	}

	// Convert algorithm-based request to oshash/phash parameters
	var oshash, phash string
	switch req.Algorithm {
	case gen.AdultMatchRequestAlgorithmOshash:
		oshash = req.Hash
	case gen.AdultMatchRequestAlgorithmPhash:
		phash = req.Hash
	}

	voy, err := svc.MatchByFingerprint(ctx, oshash, phash)
	if err != nil {
		return &gen.AdultMatchResponse{
			Matched:    gen.NewOptBool(false),
			Confidence: gen.NewOptFloat64(0),
		}, nil
	}

	scene := voyageToAPI(voy)
	return &gen.AdultMatchResponse{
		Matched:    gen.NewOptBool(true),
		Scene:      gen.NewOptAdultScene(scene),
		Confidence: gen.NewOptFloat64(1.0),
	}, nil
}

// =============================================================================
// Request System Handlers (Stub implementations)
// =============================================================================

func (h *Handler) SearchAdultRequests(ctx context.Context, params gen.SearchAdultRequestsParams) (gen.SearchAdultRequestsRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	// TODO: Implement request system
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) ListAdultRequests(ctx context.Context, params gen.ListAdultRequestsParams) (gen.ListAdultRequestsRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) CreateAdultRequest(ctx context.Context, req *gen.AdultRequestCreate) (gen.CreateAdultRequestRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) GetAdultRequest(ctx context.Context, params gen.GetAdultRequestParams) (gen.GetAdultRequestRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.GetAdultRequestUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.GetAdultRequestNotFound{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) VoteAdultRequest(ctx context.Context, req *gen.AdultRequestVoteCreate, params gen.VoteAdultRequestParams) (gen.VoteAdultRequestRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) CommentAdultRequest(ctx context.Context, req *gen.AdultRequestCommentCreate, params gen.CommentAdultRequestParams) (gen.CommentAdultRequestRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) ListAdultAdminRequests(ctx context.Context, params gen.ListAdultAdminRequestsParams) (gen.ListAdultAdminRequestsRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.ListAdultAdminRequestsForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	return &gen.ListAdultAdminRequestsUnauthorized{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) ApproveAdultRequest(ctx context.Context, params gen.ApproveAdultRequestParams) (gen.ApproveAdultRequestRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.ApproveAdultRequestForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	return &gen.ApproveAdultRequestNotFound{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) DeclineAdultRequest(ctx context.Context, params gen.DeclineAdultRequestParams) (gen.DeclineAdultRequestRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.DeclineAdultRequestForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	return &gen.DeclineAdultRequestNotFound{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) UpdateAdultRequestQuota(ctx context.Context, req *gen.AdultRequestQuotaUpdate, params gen.UpdateAdultRequestQuotaParams) (gen.UpdateAdultRequestQuotaRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.UpdateAdultRequestQuotaForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	// TODO: Implement request quota system
	return &gen.UpdateAdultRequestQuotaNoContent{}, nil
}

func (h *Handler) ListAdultRequestRules(ctx context.Context) (gen.ListAdultRequestRulesRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.ListAdultRequestRulesForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	// TODO: Implement request rules system
	return &gen.AdultRequestRuleListResponse{
		Rules: []gen.AdultRequestRule{},
	}, nil
}

func (h *Handler) CreateAdultRequestRule(ctx context.Context, req *gen.AdultRequestRuleCreate) (gen.CreateAdultRequestRuleRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.CreateAdultRequestRuleForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	return &gen.CreateAdultRequestRuleUnauthorized{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) UpdateAdultRequestRule(ctx context.Context, req *gen.AdultRequestRuleUpdate, params gen.UpdateAdultRequestRuleParams) (gen.UpdateAdultRequestRuleRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.UpdateAdultRequestRuleNotFound{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	return &gen.UpdateAdultRequestRuleNotFound{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

func (h *Handler) DeleteAdultRequestRule(ctx context.Context, params gen.DeleteAdultRequestRuleParams) (gen.DeleteAdultRequestRuleRes, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return &gen.DeleteAdultRequestRuleForbidden{
			Code:    "forbidden",
			Message: "Admin access required",
		}, nil
	}
	return &gen.DeleteAdultRequestRuleNotFound{
		Code:    "not_implemented",
		Message: "Request system not yet implemented",
	}, nil
}

// =============================================================================
// External Metadata Handlers (StashDB, TPDB, Stash-App)
// =============================================================================

func (h *Handler) SearchAdultStashDBScenes(ctx context.Context, params gen.SearchAdultStashDBScenesParams) (gen.SearchAdultStashDBScenesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	// TODO: Implement StashDB integration
	return &gen.Error{
		Code:    "not_implemented",
		Message: "StashDB integration not yet implemented",
	}, nil
}

func (h *Handler) GetAdultStashDBScene(ctx context.Context, params gen.GetAdultStashDBSceneParams) (gen.GetAdultStashDBSceneRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.GetAdultStashDBSceneUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.GetAdultStashDBSceneNotFound{
		Code:    "not_implemented",
		Message: "StashDB integration not yet implemented",
	}, nil
}

func (h *Handler) SearchAdultStashDBPerformers(ctx context.Context, params gen.SearchAdultStashDBPerformersParams) (gen.SearchAdultStashDBPerformersRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "StashDB integration not yet implemented",
	}, nil
}

func (h *Handler) GetAdultStashDBPerformer(ctx context.Context, params gen.GetAdultStashDBPerformerParams) (gen.GetAdultStashDBPerformerRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.GetAdultStashDBPerformerUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.GetAdultStashDBPerformerNotFound{
		Code:    "not_implemented",
		Message: "StashDB integration not yet implemented",
	}, nil
}

func (h *Handler) IdentifyAdultStashDBScene(ctx context.Context, req *gen.AdultMatchRequest) (gen.IdentifyAdultStashDBSceneRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "StashDB integration not yet implemented",
	}, nil
}

func (h *Handler) SearchAdultTPDBScenes(ctx context.Context, params gen.SearchAdultTPDBScenesParams) (gen.SearchAdultTPDBScenesRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "TPDB integration not yet implemented",
	}, nil
}

func (h *Handler) GetAdultTPDBScene(ctx context.Context, params gen.GetAdultTPDBSceneParams) (gen.GetAdultTPDBSceneRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.GetAdultTPDBSceneUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.GetAdultTPDBSceneNotFound{
		Code:    "not_implemented",
		Message: "TPDB integration not yet implemented",
	}, nil
}

func (h *Handler) GetAdultTPDBPerformer(ctx context.Context, params gen.GetAdultTPDBPerformerParams) (gen.GetAdultTPDBPerformerRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.GetAdultTPDBPerformerUnauthorized{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.GetAdultTPDBPerformerNotFound{
		Code:    "not_implemented",
		Message: "TPDB integration not yet implemented",
	}, nil
}

func (h *Handler) SyncAdultStash(ctx context.Context) (gen.SyncAdultStashRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		return &gen.Error{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Stash-App sync not yet implemented",
	}, nil
}

func (h *Handler) ImportAdultStash(ctx context.Context) (gen.ImportAdultStashRes, error) {
	if _, err := h.requireAdultMetadataWrite(ctx); err != nil {
		return &gen.Error{
			Code:    "forbidden",
			Message: "Adult content write access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Stash-App import not yet implemented",
	}, nil
}

func (h *Handler) GetAdultStashStatus(ctx context.Context) (gen.GetAdultStashStatusRes, error) {
	if _, err := h.requireAdultBrowse(ctx); err != nil {
		return &gen.Error{
			Code:    "unauthorized",
			Message: "Adult content access denied",
		}, nil
	}
	return &gen.Error{
		Code:    "not_implemented",
		Message: "Stash-App integration not yet implemented",
	}, nil
}
