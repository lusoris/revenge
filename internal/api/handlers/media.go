package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jellyfin/jellyfin-go/internal/api/middleware"
	"github.com/jellyfin/jellyfin-go/internal/domain"
)

// MediaService defines the interface for media item operations.
type MediaService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.MediaItem, error)
	List(ctx context.Context, params domain.ListMediaItemsParams) ([]*domain.MediaItem, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.MediaItem, error)
	Search(ctx context.Context, query string, limit int32) ([]*domain.MediaItem, error)
}

// MediaHandler handles media item HTTP requests.
type MediaHandler struct {
	mediaService MediaService
}

// NewMediaHandler creates a new MediaHandler.
func NewMediaHandler(mediaService MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

// RegisterRoutes registers media item routes on the given mux.
func (h *MediaHandler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
	// Item endpoints (Jellyfin API compatible)
	mux.Handle("GET /Items", auth.Required(http.HandlerFunc(h.ListItems)))
	mux.Handle("GET /Items/{itemId}", auth.Required(http.HandlerFunc(h.GetItem)))
	mux.Handle("GET /Items/{itemId}/Similar", auth.Required(http.HandlerFunc(h.GetSimilarItems)))
	mux.Handle("GET /Items/{itemId}/ThemeSongs", auth.Required(http.HandlerFunc(h.GetThemeSongs)))
	mux.Handle("GET /Items/{itemId}/ThemeVideos", auth.Required(http.HandlerFunc(h.GetThemeVideos)))

	// User item endpoints
	mux.Handle("GET /Users/{userId}/Items", auth.Required(http.HandlerFunc(h.ListUserItems)))
	mux.Handle("GET /Users/{userId}/Items/{itemId}", auth.Required(http.HandlerFunc(h.GetUserItem)))
	mux.Handle("GET /Users/{userId}/Items/Latest", auth.Required(http.HandlerFunc(h.GetLatestItems)))
	mux.Handle("GET /Users/{userId}/Items/Resume", auth.Required(http.HandlerFunc(h.GetResumeItems)))

	// Search
	mux.Handle("GET /Search/Hints", auth.Required(http.HandlerFunc(h.SearchHints)))

	// Special collections
	mux.Handle("GET /Movies/Recommendations", auth.Required(http.HandlerFunc(h.GetMovieRecommendations)))
	mux.Handle("GET /Shows/NextUp", auth.Required(http.HandlerFunc(h.GetNextUp)))
}

// ItemResponse represents a media item in API responses.
// Matches Jellyfin API BaseItemDto format.
type ItemResponse struct {
	ID                    string            `json:"Id"`
	Name                  string            `json:"Name"`
	SortName              string            `json:"SortName,omitempty"`
	Type                  string            `json:"Type"`
	MediaType             string            `json:"MediaType,omitempty"`
	Overview              *string           `json:"Overview,omitempty"`
	Tagline               *string           `json:"Taglines,omitempty"`
	ParentID              *string           `json:"ParentId,omitempty"`
	LibraryID             *string           `json:"ParentLogoItemId,omitempty"`
	Path                  *string           `json:"Path,omitempty"`

	// Dates
	PremiereDate          *string           `json:"PremiereDate,omitempty"`
	EndDate               *string           `json:"EndDate,omitempty"`
	ProductionYear        *int              `json:"ProductionYear,omitempty"`
	DateCreated           *string           `json:"DateCreated,omitempty"`

	// Runtime
	RunTimeTicks          *int64            `json:"RunTimeTicks,omitempty"`

	// Series/Episode info
	SeriesName            *string           `json:"SeriesName,omitempty"`
	SeriesID              *string           `json:"SeriesId,omitempty"`
	SeasonID              *string           `json:"SeasonId,omitempty"`
	SeasonName            *string           `json:"SeasonName,omitempty"`
	IndexNumber           *int              `json:"IndexNumber,omitempty"`
	ParentIndexNumber     *int              `json:"ParentIndexNumber,omitempty"`

	// Music info
	AlbumArtist           *string           `json:"AlbumArtist,omitempty"`
	Album                 *string           `json:"Album,omitempty"`
	Artists               []string          `json:"Artists,omitempty"`

	// Ratings
	CommunityRating       *float64          `json:"CommunityRating,omitempty"`
	CriticRating          *float64          `json:"CriticRating,omitempty"`
	OfficialRating        *string           `json:"OfficialRating,omitempty"`

	// Metadata
	Genres                []string          `json:"Genres,omitempty"`
	Tags                  []string          `json:"Tags,omitempty"`
	Studios               []StudioInfo      `json:"Studios,omitempty"`
	ProviderIDs           map[string]string `json:"ProviderIds,omitempty"`

	// Media info
	Container             *string           `json:"Container,omitempty"`
	Width                 *int              `json:"Width,omitempty"`
	Height                *int              `json:"Height,omitempty"`

	// Collection info
	ChildCount            *int              `json:"ChildCount,omitempty"`
	RecursiveItemCount    *int              `json:"RecursiveItemCount,omitempty"`

	// User data (if requested)
	UserData              *UserDataResponse `json:"UserData,omitempty"`

	// Image info
	ImageTags             map[string]string `json:"ImageTags,omitempty"`
	BackdropImageTags     []string          `json:"BackdropImageTags,omitempty"`

	// Status
	IsFolder              bool              `json:"IsFolder"`
	CanDelete             bool              `json:"CanDelete,omitempty"`
	CanDownload           bool              `json:"CanDownload,omitempty"`
}

// StudioInfo represents a studio.
type StudioInfo struct {
	Name string `json:"Name"`
	ID   string `json:"Id,omitempty"`
}

// UserDataResponse represents user-specific item data.
type UserDataResponse struct {
	PlaybackPositionTicks int64   `json:"PlaybackPositionTicks"`
	PlayCount             int     `json:"PlayCount"`
	IsFavorite            bool    `json:"IsFavorite"`
	Played                bool    `json:"Played"`
	UnplayedItemCount     *int    `json:"UnplayedItemCount,omitempty"`
	PlayedPercentage      *float64 `json:"PlayedPercentage,omitempty"`
	LastPlayedDate        *string `json:"LastPlayedDate,omitempty"`
}

// ItemsResponse represents a paginated list of items.
type ItemsResponse struct {
	Items            []ItemResponse `json:"Items"`
	TotalRecordCount int            `json:"TotalRecordCount"`
	StartIndex       int            `json:"StartIndex"`
}

// SearchHintResponse represents a search hint.
type SearchHintResponse struct {
	ItemID            string   `json:"ItemId"`
	ID                string   `json:"Id"`
	Name              string   `json:"Name"`
	Type              string   `json:"Type"`
	MediaType         string   `json:"MediaType,omitempty"`
	RunTimeTicks      *int64   `json:"RunTimeTicks,omitempty"`
	ProductionYear    *int     `json:"ProductionYear,omitempty"`
	PrimaryImageTag   *string  `json:"PrimaryImageTag,omitempty"`
	ThumbImageTag     *string  `json:"ThumbImageTag,omitempty"`
	PrimaryImageAspectRatio *float64 `json:"PrimaryImageAspectRatio,omitempty"`
	Album             *string  `json:"Album,omitempty"`
	AlbumArtist       *string  `json:"AlbumArtist,omitempty"`
	Artists           []string `json:"Artists,omitempty"`
	SeriesName        *string  `json:"SeriesName,omitempty"`
}

// SearchHintsResponse represents search results.
type SearchHintsResponse struct {
	SearchHints      []SearchHintResponse `json:"SearchHints"`
	TotalRecordCount int                  `json:"TotalRecordCount"`
}

// itemToResponse converts a domain.MediaItem to ItemResponse.
func itemToResponse(item *domain.MediaItem) ItemResponse {
	resp := ItemResponse{
		ID:              item.ID.String(),
		Name:            item.Name,
		Type:            jellfinItemType(item.Type),
		MediaType:       jellyfinMediaType(item.Type),
		Overview:        item.Overview,
		Tagline:         item.Tagline,
		ProductionYear:  item.Year,
		RunTimeTicks:    item.RuntimeTicks,
		CommunityRating: item.CommunityRating,
		CriticRating:    item.CriticRating,
		Genres:          item.Genres,
		Tags:            item.Tags,
		ProviderIDs:     item.ProviderIDs,
		IsFolder:        item.Type.IsContainer(),
		Container:       item.Container,
		Width:           item.Width,
		Height:          item.Height,
	}

	if item.SortName != nil {
		resp.SortName = *item.SortName
	}

	if item.ParentID != nil {
		s := item.ParentID.String()
		resp.ParentID = &s
	}

	s := item.LibraryID.String()
	resp.LibraryID = &s

	if item.Path != nil {
		resp.Path = item.Path
	}

	if item.PremiereDate != nil {
		s := item.PremiereDate.Format(time.RFC3339)
		resp.PremiereDate = &s
	}

	if item.EndDate != nil {
		s := item.EndDate.Format(time.RFC3339)
		resp.EndDate = &s
	}

	if item.DateCreated != nil {
		s := item.DateCreated.Format(time.RFC3339)
		resp.DateCreated = &s
	}

	// Series/Episode info
	if item.SeasonNumber != nil {
		resp.ParentIndexNumber = item.SeasonNumber
	}
	if item.EpisodeNumber != nil {
		resp.IndexNumber = item.EpisodeNumber
	}

	// Music info
	if item.AlbumArtist != nil {
		resp.AlbumArtist = item.AlbumArtist
	}
	if item.TrackNumber != nil {
		resp.IndexNumber = item.TrackNumber
	}

	// Studios
	if len(item.Studios) > 0 {
		resp.Studios = make([]StudioInfo, len(item.Studios))
		for i, s := range item.Studios {
			resp.Studios[i] = StudioInfo{Name: s}
		}
	}

	return resp
}

// jellfinItemType converts our MediaType to Jellyfin's Type field.
func jellfinItemType(t domain.MediaType) string {
	switch t {
	case domain.MediaTypeMovie:
		return "Movie"
	case domain.MediaTypeEpisode:
		return "Episode"
	case domain.MediaTypeSeries:
		return "Series"
	case domain.MediaTypeSeason:
		return "Season"
	case domain.MediaTypeAudio:
		return "Audio"
	case domain.MediaTypeAlbum:
		return "MusicAlbum"
	case domain.MediaTypeArtist:
		return "MusicArtist"
	case domain.MediaTypeMusicVideo:
		return "MusicVideo"
	case domain.MediaTypePhoto:
		return "Photo"
	case domain.MediaTypePhotoAlbum:
		return "PhotoAlbum"
	case domain.MediaTypeBook:
		return "Book"
	case domain.MediaTypeAudiobook:
		return "AudioBook"
	case domain.MediaTypePodcast:
		return "Podcast"
	case domain.MediaTypePodcastEpisode:
		return "PodcastEpisode"
	case domain.MediaTypeBoxSet:
		return "BoxSet"
	case domain.MediaTypePlaylist:
		return "Playlist"
	case domain.MediaTypeFolder:
		return "Folder"
	case domain.MediaTypeTrailer:
		return "Trailer"
	case domain.MediaTypeChannel:
		return "Channel"
	case domain.MediaTypeProgram:
		return "Program"
	case domain.MediaTypeRecording:
		return "Recording"
	default:
		return "Unknown"
	}
}

// jellyfinMediaType converts our MediaType to Jellyfin's MediaType field.
func jellyfinMediaType(t domain.MediaType) string {
	switch t {
	case domain.MediaTypeMovie, domain.MediaTypeEpisode, domain.MediaTypeMusicVideo,
		domain.MediaTypeTrailer, domain.MediaTypeHomeVideo, domain.MediaTypeRecording:
		return "Video"
	case domain.MediaTypeAudio, domain.MediaTypeAudiobookChapter, domain.MediaTypePodcastEpisode:
		return "Audio"
	case domain.MediaTypePhoto:
		return "Photo"
	case domain.MediaTypeBook:
		return "Book"
	default:
		return ""
	}
}

// ListItems handles GET /Items
func (h *MediaHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	params := parseItemQueryParams(r)

	items, err := h.mediaService.List(r.Context(), params)
	if err != nil {
		InternalError(w, err)
		return
	}

	response := ItemsResponse{
		Items:            make([]ItemResponse, len(items)),
		TotalRecordCount: len(items), // TODO: Get actual count
		StartIndex:       int(params.Offset),
	}

	for i, item := range items {
		response.Items[i] = itemToResponse(item)
	}

	OK(w, response)
}

// GetItem handles GET /Items/{itemId}
func (h *MediaHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	item, err := h.mediaService.GetByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFound(w, "Item not found")
			return
		}
		InternalError(w, err)
		return
	}

	OK(w, itemToResponse(item))
}

// GetSimilarItems handles GET /Items/{itemId}/Similar
func (h *MediaHandler) GetSimilarItems(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	// TODO: Implement similar items algorithm
	_ = itemID

	OK(w, ItemsResponse{
		Items:            []ItemResponse{},
		TotalRecordCount: 0,
		StartIndex:       0,
	})
}

// GetThemeSongs handles GET /Items/{itemId}/ThemeSongs
func (h *MediaHandler) GetThemeSongs(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	// TODO: Implement theme songs
	OK(w, ItemsResponse{
		Items:            []ItemResponse{},
		TotalRecordCount: 0,
		StartIndex:       0,
	})
}

// GetThemeVideos handles GET /Items/{itemId}/ThemeVideos
func (h *MediaHandler) GetThemeVideos(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	// TODO: Implement theme videos
	OK(w, ItemsResponse{
		Items:            []ItemResponse{},
		TotalRecordCount: 0,
		StartIndex:       0,
	})
}

// ListUserItems handles GET /Users/{userId}/Items
func (h *MediaHandler) ListUserItems(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	params := parseItemQueryParams(r)
	params.UserID = &userID

	items, err := h.mediaService.List(r.Context(), params)
	if err != nil {
		InternalError(w, err)
		return
	}

	response := ItemsResponse{
		Items:            make([]ItemResponse, len(items)),
		TotalRecordCount: len(items),
		StartIndex:       int(params.Offset),
	}

	for i, item := range items {
		response.Items[i] = itemToResponse(item)
	}

	OK(w, response)
}

// GetUserItem handles GET /Users/{userId}/Items/{itemId}
func (h *MediaHandler) GetUserItem(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	itemID, err := uuid.Parse(r.PathValue("itemId"))
	if err != nil {
		BadRequest(w, "Invalid item ID")
		return
	}

	item, err := h.mediaService.GetByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFound(w, "Item not found")
			return
		}
		InternalError(w, err)
		return
	}

	// TODO: Add user-specific data (played status, favorites, etc.)
	OK(w, itemToResponse(item))
}

// GetLatestItems handles GET /Users/{userId}/Items/Latest
func (h *MediaHandler) GetLatestItems(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	params := parseItemQueryParams(r)
	params.UserID = &userID
	params.SortBy = "DateCreated"
	params.SortOrder = "desc"
	if params.Limit == 0 {
		params.Limit = 20
	}

	items, err := h.mediaService.List(r.Context(), params)
	if err != nil {
		InternalError(w, err)
		return
	}

	response := make([]ItemResponse, len(items))
	for i, item := range items {
		response[i] = itemToResponse(item)
	}

	OK(w, response)
}

// GetResumeItems handles GET /Users/{userId}/Items/Resume
func (h *MediaHandler) GetResumeItems(w http.ResponseWriter, r *http.Request) {
	_, err := uuid.Parse(r.PathValue("userId"))
	if err != nil {
		BadRequest(w, "Invalid user ID")
		return
	}

	// TODO: Implement resume items (items with playback position)
	OK(w, ItemsResponse{
		Items:            []ItemResponse{},
		TotalRecordCount: 0,
		StartIndex:       0,
	})
}

// SearchHints handles GET /Search/Hints
func (h *MediaHandler) SearchHints(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("searchTerm")
	if searchTerm == "" {
		BadRequest(w, "searchTerm is required")
		return
	}

	limit := int32(20)
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 32); err == nil {
			limit = int32(parsed)
		}
	}

	items, err := h.mediaService.Search(r.Context(), searchTerm, limit)
	if err != nil {
		InternalError(w, err)
		return
	}

	hints := make([]SearchHintResponse, len(items))
	for i, item := range items {
		hints[i] = SearchHintResponse{
			ItemID:         item.ID.String(),
			ID:             item.ID.String(),
			Name:           item.Name,
			Type:           jellfinItemType(item.Type),
			MediaType:      jellyfinMediaType(item.Type),
			RunTimeTicks:   item.RuntimeTicks,
			ProductionYear: item.Year,
		}
	}

	OK(w, SearchHintsResponse{
		SearchHints:      hints,
		TotalRecordCount: len(hints),
	})
}

// GetMovieRecommendations handles GET /Movies/Recommendations
func (h *MediaHandler) GetMovieRecommendations(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement movie recommendations
	OK(w, []any{})
}

// GetNextUp handles GET /Shows/NextUp
func (h *MediaHandler) GetNextUp(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement NextUp for TV shows
	OK(w, ItemsResponse{
		Items:            []ItemResponse{},
		TotalRecordCount: 0,
		StartIndex:       0,
	})
}

// parseItemQueryParams parses query parameters into ListMediaItemsParams.
func parseItemQueryParams(r *http.Request) domain.ListMediaItemsParams {
	params := domain.ListMediaItemsParams{
		Limit:  50,
		Offset: 0,
	}

	if parentID := r.URL.Query().Get("parentId"); parentID != "" {
		if id, err := uuid.Parse(parentID); err == nil {
			params.ParentID = &id
		}
	}

	if libraryID := r.URL.Query().Get("libraryId"); libraryID != "" {
		if id, err := uuid.Parse(libraryID); err == nil {
			params.LibraryID = &id
		}
	}

	// Parse includeItemTypes
	if types := r.URL.Query().Get("includeItemTypes"); types != "" {
		typeList := strings.Split(types, ",")
		if len(typeList) == 1 {
			mt := jellyfinTypeToMediaType(typeList[0])
			params.MediaType = &mt
		}
	}

	// Parse genres
	if genres := r.URL.Query().Get("genres"); genres != "" {
		params.Genres = strings.Split(genres, ",")
	}

	// Parse years
	if years := r.URL.Query().Get("years"); years != "" {
		yearStrs := strings.Split(years, ",")
		for _, ys := range yearStrs {
			if y, err := strconv.Atoi(ys); err == nil {
				params.Years = append(params.Years, y)
			}
		}
	}

	// Parse tags
	if tags := r.URL.Query().Get("tags"); tags != "" {
		params.Tags = strings.Split(tags, ",")
	}

	// Parse sorting
	if sortBy := r.URL.Query().Get("sortBy"); sortBy != "" {
		params.SortBy = sortBy
	}
	if sortOrder := r.URL.Query().Get("sortOrder"); sortOrder != "" {
		params.SortOrder = strings.ToLower(sortOrder)
	}

	// Parse pagination
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.ParseInt(limit, 10, 32); err == nil {
			params.Limit = int32(l)
		}
	}
	if startIndex := r.URL.Query().Get("startIndex"); startIndex != "" {
		if s, err := strconv.ParseInt(startIndex, 10, 32); err == nil {
			params.Offset = int32(s)
		}
	}

	return params
}

// jellyfinTypeToMediaType converts Jellyfin type to our MediaType.
func jellyfinTypeToMediaType(jType string) domain.MediaType {
	switch strings.ToLower(jType) {
	case "movie":
		return domain.MediaTypeMovie
	case "episode":
		return domain.MediaTypeEpisode
	case "series":
		return domain.MediaTypeSeries
	case "season":
		return domain.MediaTypeSeason
	case "audio":
		return domain.MediaTypeAudio
	case "musicalbum":
		return domain.MediaTypeAlbum
	case "musicartist":
		return domain.MediaTypeArtist
	case "musicvideo":
		return domain.MediaTypeMusicVideo
	case "photo":
		return domain.MediaTypePhoto
	case "photoalbum":
		return domain.MediaTypePhotoAlbum
	case "book":
		return domain.MediaTypeBook
	case "audiobook":
		return domain.MediaTypeAudiobook
	case "podcast":
		return domain.MediaTypePodcast
	case "boxset":
		return domain.MediaTypeBoxSet
	case "playlist":
		return domain.MediaTypePlaylist
	case "folder":
		return domain.MediaTypeFolder
	case "trailer":
		return domain.MediaTypeTrailer
	case "channel":
		return domain.MediaTypeChannel
	case "program":
		return domain.MediaTypeProgram
	case "recording":
		return domain.MediaTypeRecording
	default:
		return domain.MediaTypeMovie
	}
}
