package tvshow

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	db "github.com/lusoris/revenge/internal/content/tvshow/db"
)

// Networks

// GetNetworkByID retrieves a network by ID.
func (r *pgRepository) GetNetworkByID(ctx context.Context, id uuid.UUID) (*Network, error) {
	n, err := r.queries.GetNetworkByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNetworkNotFound
		}
		return nil, err
	}
	return fromDBNetwork(&n), nil
}

// GetNetworkByTmdbID retrieves a network by TMDb ID.
func (r *pgRepository) GetNetworkByTmdbID(ctx context.Context, tmdbID int) (*Network, error) {
	id := int32(tmdbID)
	n, err := r.queries.GetNetworkByTmdbID(ctx, &id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNetworkNotFound
		}
		return nil, err
	}
	return fromDBNetwork(&n), nil
}

// ListNetworks lists all networks.
func (r *pgRepository) ListNetworks(ctx context.Context, params ListParams) ([]*Network, error) {
	rows, err := r.queries.ListNetworks(ctx, db.ListNetworksParams{
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	networks := make([]*Network, len(rows))
	for i, row := range rows {
		networks[i] = fromDBNetwork(&row)
	}
	return networks, nil
}

// CreateNetwork creates a new network.
func (r *pgRepository) CreateNetwork(ctx context.Context, network *Network) error {
	var tmdbID *int32
	if network.TmdbID > 0 {
		id := int32(network.TmdbID)
		tmdbID = &id
	}

	n, err := r.queries.CreateNetwork(ctx, db.CreateNetworkParams{
		Name:          network.Name,
		LogoPath:      &network.LogoPath,
		OriginCountry: &network.OriginCountry,
		TmdbID:        tmdbID,
	})
	if err != nil {
		return err
	}
	network.ID = n.ID
	network.CreatedAt = n.CreatedAt
	return nil
}

// GetOrCreateNetwork gets or creates a network.
func (r *pgRepository) GetOrCreateNetwork(ctx context.Context, name string, logoPath, originCountry string, tmdbID int) (*Network, error) {
	var tid *int32
	if tmdbID > 0 {
		id := int32(tmdbID)
		tid = &id
	}

	n, err := r.queries.GetOrCreateNetwork(ctx, db.GetOrCreateNetworkParams{
		Name:          name,
		LogoPath:      &logoPath,
		OriginCountry: &originCountry,
		TmdbID:        tid,
	})
	if err != nil {
		return nil, err
	}
	return fromDBNetwork(&n), nil
}

// GetSeriesNetworks retrieves networks for a series.
func (r *pgRepository) GetSeriesNetworks(ctx context.Context, seriesID uuid.UUID) ([]Network, error) {
	rows, err := r.queries.GetSeriesNetworks(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	networks := make([]Network, len(rows))
	for i, row := range rows {
		networks[i] = *fromDBNetwork(&row)
	}
	return networks, nil
}

// LinkSeriesNetwork links a network to a series.
func (r *pgRepository) LinkSeriesNetwork(ctx context.Context, seriesID, networkID uuid.UUID, order int) error {
	return r.queries.LinkSeriesNetwork(ctx, db.LinkSeriesNetworkParams{
		SeriesID:     seriesID,
		NetworkID:    networkID,
		DisplayOrder: int32(order),
	})
}

// UnlinkSeriesNetworks removes all network links from a series.
func (r *pgRepository) UnlinkSeriesNetworks(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.UnlinkSeriesNetworks(ctx, seriesID)
}

// ListSeriesByNetwork lists series by network.
func (r *pgRepository) ListSeriesByNetwork(ctx context.Context, networkID uuid.UUID, params ListParams) ([]*Series, error) {
	rows, err := r.queries.ListSeriesByNetwork(ctx, db.ListSeriesByNetworkParams{
		NetworkID: networkID,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// CountSeriesByNetwork counts series by network.
func (r *pgRepository) CountSeriesByNetwork(ctx context.Context, networkID uuid.UUID) (int64, error) {
	return r.queries.CountSeriesByNetwork(ctx, networkID)
}

// Genres

// GetSeriesGenres retrieves genres for a series.
func (r *pgRepository) GetSeriesGenres(ctx context.Context, seriesID uuid.UUID) ([]Genre, error) {
	rows, err := r.queries.GetSeriesGenres(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	genres := make([]Genre, len(rows))
	for i, row := range rows {
		genres[i] = Genre{
			ID:   row.ID,
			Name: row.Name,
		}
	}
	return genres, nil
}

// LinkSeriesGenre links a genre to a series.
func (r *pgRepository) LinkSeriesGenre(ctx context.Context, seriesID, genreID uuid.UUID) error {
	return r.queries.LinkSeriesGenre(ctx, db.LinkSeriesGenreParams{
		SeriesID: seriesID,
		GenreID:  genreID,
	})
}

// UnlinkSeriesGenres removes all genre links from a series.
func (r *pgRepository) UnlinkSeriesGenres(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.UnlinkSeriesGenres(ctx, seriesID)
}

// ListSeriesByGenre lists series by genre.
func (r *pgRepository) ListSeriesByGenre(ctx context.Context, genreID uuid.UUID, params ListParams) ([]*Series, error) {
	rows, err := r.queries.ListSeriesByGenre(ctx, db.ListSeriesByGenreParams{
		GenreID: genreID,
		Limit:   int32(params.Limit),
		Offset:  int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	series := make([]*Series, len(rows))
	for i, row := range rows {
		series[i] = FromDBSeries(&row)
	}
	return series, nil
}

// CountSeriesByGenre counts series by genre.
func (r *pgRepository) CountSeriesByGenre(ctx context.Context, genreID uuid.UUID) (int64, error) {
	return r.queries.CountSeriesByGenre(ctx, genreID)
}

// GetOrCreateGenre gets or creates a genre.
func (r *pgRepository) GetOrCreateGenre(ctx context.Context, name string, tmdbID int) (*Genre, error) {
	var tid *int32
	if tmdbID > 0 {
		id := int32(tmdbID)
		tid = &id
	}

	g, err := r.queries.GetOrCreateTVShowGenre(ctx, db.GetOrCreateTVShowGenreParams{
		Name:   name,
		TmdbID: tid,
	})
	if err != nil {
		return nil, err
	}
	return &Genre{
		ID:   g.ID,
		Name: g.Name,
	}, nil
}

// Series Credits

// GetSeriesCast retrieves cast members for a series.
func (r *pgRepository) GetSeriesCast(ctx context.Context, seriesID uuid.UUID) ([]CastMember, error) {
	rows, err := r.queries.GetSeriesCast(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	cast := make([]CastMember, len(rows))
	for i, row := range rows {
		cast[i] = CastMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			BillingOrder: int(row.BillingOrder),
		}
		if row.CharacterName != nil {
			cast[i].CharacterName = *row.CharacterName
		}
		if row.PrimaryImageUrl != nil {
			cast[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			cast[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return cast, nil
}

// GetSeriesCrew retrieves crew members for a series.
func (r *pgRepository) GetSeriesCrew(ctx context.Context, seriesID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetSeriesCrew(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	crew := make([]CrewMember, len(rows))
	for i, row := range rows {
		crew[i] = CrewMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			Role:         string(row.Role),
			BillingOrder: int(row.BillingOrder),
		}
		if row.Department != nil {
			crew[i].Department = *row.Department
		}
		if row.Job != nil {
			crew[i].Job = *row.Job
		}
		if row.PrimaryImageUrl != nil {
			crew[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			crew[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return crew, nil
}

// GetSeriesCreators retrieves creators for a series.
func (r *pgRepository) GetSeriesCreators(ctx context.Context, seriesID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetSeriesCreators(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	creators := make([]CrewMember, len(rows))
	for i, row := range rows {
		creators[i] = CrewMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			Role:         string(row.Role),
			BillingOrder: int(row.BillingOrder),
		}
		if row.PrimaryImageUrl != nil {
			creators[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			creators[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return creators, nil
}

// CreateSeriesCredit creates a series credit.
func (r *pgRepository) CreateSeriesCredit(ctx context.Context, seriesID, personID uuid.UUID, role, character, department, job string, order int, tmdbCreditID string) error {
	params := db.CreateSeriesCreditParams{
		SeriesID:     seriesID,
		PersonID:     personID,
		Role:         role,
		BillingOrder: int32(order),
	}
	if character != "" {
		params.CharacterName = &character
	}
	if department != "" {
		params.Department = &department
	}
	if job != "" {
		params.Job = &job
	}
	if tmdbCreditID != "" {
		params.TmdbCreditID = &tmdbCreditID
	}

	_, err := r.queries.CreateSeriesCredit(ctx, params)
	return err
}

// DeleteSeriesCredits deletes all credits for a series.
func (r *pgRepository) DeleteSeriesCredits(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesCredits(ctx, seriesID)
}

// Episode Credits

// GetEpisodeCast retrieves cast members for an episode.
func (r *pgRepository) GetEpisodeCast(ctx context.Context, episodeID uuid.UUID) ([]CastMember, error) {
	rows, err := r.queries.GetEpisodeCast(ctx, episodeID)
	if err != nil {
		return nil, err
	}

	cast := make([]CastMember, len(rows))
	for i, row := range rows {
		cast[i] = CastMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			BillingOrder: int(row.BillingOrder),
			IsGuest:      row.IsGuest,
		}
		if row.CharacterName != nil {
			cast[i].CharacterName = *row.CharacterName
		}
		if row.PrimaryImageUrl != nil {
			cast[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			cast[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return cast, nil
}

// GetEpisodeGuestStars retrieves guest stars for an episode.
func (r *pgRepository) GetEpisodeGuestStars(ctx context.Context, episodeID uuid.UUID) ([]CastMember, error) {
	rows, err := r.queries.GetEpisodeGuestStars(ctx, episodeID)
	if err != nil {
		return nil, err
	}

	cast := make([]CastMember, len(rows))
	for i, row := range rows {
		cast[i] = CastMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			BillingOrder: int(row.BillingOrder),
			IsGuest:      true,
		}
		if row.CharacterName != nil {
			cast[i].CharacterName = *row.CharacterName
		}
		if row.PrimaryImageUrl != nil {
			cast[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			cast[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return cast, nil
}

// GetEpisodeCrew retrieves crew members for an episode.
func (r *pgRepository) GetEpisodeCrew(ctx context.Context, episodeID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetEpisodeCrew(ctx, episodeID)
	if err != nil {
		return nil, err
	}

	crew := make([]CrewMember, len(rows))
	for i, row := range rows {
		crew[i] = CrewMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			Role:         string(row.Role),
			BillingOrder: int(row.BillingOrder),
		}
		if row.Department != nil {
			crew[i].Department = *row.Department
		}
		if row.Job != nil {
			crew[i].Job = *row.Job
		}
		if row.PrimaryImageUrl != nil {
			crew[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			crew[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return crew, nil
}

// GetEpisodeDirectors retrieves directors for an episode.
func (r *pgRepository) GetEpisodeDirectors(ctx context.Context, episodeID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetEpisodeDirectors(ctx, episodeID)
	if err != nil {
		return nil, err
	}

	directors := make([]CrewMember, len(rows))
	for i, row := range rows {
		directors[i] = CrewMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			Role:         "director",
			BillingOrder: int(row.BillingOrder),
		}
		if row.PrimaryImageUrl != nil {
			directors[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			directors[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return directors, nil
}

// GetEpisodeWriters retrieves writers for an episode.
func (r *pgRepository) GetEpisodeWriters(ctx context.Context, episodeID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetEpisodeWriters(ctx, episodeID)
	if err != nil {
		return nil, err
	}

	writers := make([]CrewMember, len(rows))
	for i, row := range rows {
		writers[i] = CrewMember{
			PersonID:     row.PersonID,
			Name:         row.Name,
			Role:         "writer",
			BillingOrder: int(row.BillingOrder),
		}
		if row.PrimaryImageUrl != nil {
			writers[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			writers[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return writers, nil
}

// CreateEpisodeCredit creates an episode credit.
func (r *pgRepository) CreateEpisodeCredit(ctx context.Context, episodeID, personID uuid.UUID, role, character, department, job string, order int, isGuest bool, tmdbCreditID string) error {
	params := db.CreateEpisodeCreditParams{
		EpisodeID:    episodeID,
		PersonID:     personID,
		Role:         role,
		BillingOrder: int32(order),
		IsGuest:      isGuest,
	}
	if character != "" {
		params.CharacterName = &character
	}
	if department != "" {
		params.Department = &department
	}
	if job != "" {
		params.Job = &job
	}
	if tmdbCreditID != "" {
		params.TmdbCreditID = &tmdbCreditID
	}

	_, err := r.queries.CreateEpisodeCredit(ctx, params)
	return err
}

// DeleteEpisodeCredits deletes all credits for an episode.
func (r *pgRepository) DeleteEpisodeCredits(ctx context.Context, episodeID uuid.UUID) error {
	return r.queries.DeleteEpisodeCredits(ctx, episodeID)
}

// Images

// GetSeriesImages retrieves images for a series.
func (r *pgRepository) GetSeriesImages(ctx context.Context, seriesID uuid.UUID) ([]Image, error) {
	rows, err := r.queries.GetSeriesImages(ctx, seriesID)
	if err != nil {
		return nil, err
	}
	return fromDBSeriesImages(rows), nil
}

// GetSeriesImagesByType retrieves images of a specific type for a series.
func (r *pgRepository) GetSeriesImagesByType(ctx context.Context, seriesID uuid.UUID, imageType string) ([]Image, error) {
	rows, err := r.queries.GetSeriesImagesByType(ctx, db.GetSeriesImagesByTypeParams{
		SeriesID:  seriesID,
		ImageType: imageType,
	})
	if err != nil {
		return nil, err
	}
	return fromDBSeriesImages(rows), nil
}

// CreateSeriesImage creates a series image.
func (r *pgRepository) CreateSeriesImage(ctx context.Context, seriesID uuid.UUID, img *Image) error {
	params := db.CreateSeriesImageParams{
		SeriesID:  seriesID,
		ImageType: img.ImageType,
		Url:       img.URL,
		IsPrimary: img.IsPrimary,
	}
	if img.LocalPath != "" {
		params.LocalPath = &img.LocalPath
	}
	if img.Width > 0 {
		w := int32(img.Width)
		params.Width = &w
	}
	if img.Height > 0 {
		h := int32(img.Height)
		params.Height = &h
	}
	if img.AspectRatio > 0 {
		params.AspectRatio = numericFromFloat(img.AspectRatio)
	}
	if img.Language != "" {
		params.Language = &img.Language
	}
	if img.VoteAverage > 0 {
		params.VoteAverage = numericFromFloat(img.VoteAverage)
	}
	if img.VoteCount > 0 {
		vc := int32(img.VoteCount)
		params.VoteCount = &vc
	}
	if img.Blurhash != "" {
		params.Blurhash = &img.Blurhash
	}
	if img.Provider != "" {
		params.Provider = &img.Provider
	}
	if img.ProviderID != "" {
		params.ProviderID = &img.ProviderID
	}

	i, err := r.queries.CreateSeriesImage(ctx, params)
	if err != nil {
		return err
	}
	img.ID = i.ID
	img.CreatedAt = i.CreatedAt
	return nil
}

// DeleteSeriesImages deletes all images for a series.
func (r *pgRepository) DeleteSeriesImages(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesImages(ctx, seriesID)
}

// GetSeasonImages retrieves images for a season.
func (r *pgRepository) GetSeasonImages(ctx context.Context, seasonID uuid.UUID) ([]Image, error) {
	rows, err := r.queries.GetSeasonImages(ctx, seasonID)
	if err != nil {
		return nil, err
	}
	return fromDBSeasonImages(rows), nil
}

// CreateSeasonImage creates a season image.
func (r *pgRepository) CreateSeasonImage(ctx context.Context, seasonID uuid.UUID, img *Image) error {
	params := db.CreateSeasonImageParams{
		SeasonID:  seasonID,
		ImageType: img.ImageType,
		Url:       img.URL,
		IsPrimary: img.IsPrimary,
	}
	if img.LocalPath != "" {
		params.LocalPath = &img.LocalPath
	}
	if img.Width > 0 {
		w := int32(img.Width)
		params.Width = &w
	}
	if img.Height > 0 {
		h := int32(img.Height)
		params.Height = &h
	}
	if img.AspectRatio > 0 {
		params.AspectRatio = numericFromFloat(img.AspectRatio)
	}
	if img.Language != "" {
		params.Language = &img.Language
	}
	if img.VoteAverage > 0 {
		params.VoteAverage = numericFromFloat(img.VoteAverage)
	}
	if img.VoteCount > 0 {
		vc := int32(img.VoteCount)
		params.VoteCount = &vc
	}
	if img.Blurhash != "" {
		params.Blurhash = &img.Blurhash
	}
	if img.Provider != "" {
		params.Provider = &img.Provider
	}
	if img.ProviderID != "" {
		params.ProviderID = &img.ProviderID
	}

	i, err := r.queries.CreateSeasonImage(ctx, params)
	if err != nil {
		return err
	}
	img.ID = i.ID
	img.CreatedAt = i.CreatedAt
	return nil
}

// DeleteSeasonImages deletes all images for a season.
func (r *pgRepository) DeleteSeasonImages(ctx context.Context, seasonID uuid.UUID) error {
	return r.queries.DeleteSeasonImages(ctx, seasonID)
}

// GetEpisodeImages retrieves images for an episode.
func (r *pgRepository) GetEpisodeImages(ctx context.Context, episodeID uuid.UUID) ([]Image, error) {
	rows, err := r.queries.GetEpisodeImages(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	return fromDBEpisodeImages(rows), nil
}

// CreateEpisodeImage creates an episode image.
func (r *pgRepository) CreateEpisodeImage(ctx context.Context, episodeID uuid.UUID, img *Image) error {
	params := db.CreateEpisodeImageParams{
		EpisodeID: episodeID,
		ImageType: img.ImageType,
		Url:       img.URL,
		IsPrimary: img.IsPrimary,
	}
	if img.LocalPath != "" {
		params.LocalPath = &img.LocalPath
	}
	if img.Width > 0 {
		w := int32(img.Width)
		params.Width = &w
	}
	if img.Height > 0 {
		h := int32(img.Height)
		params.Height = &h
	}
	if img.AspectRatio > 0 {
		params.AspectRatio = numericFromFloat(img.AspectRatio)
	}
	if img.VoteAverage > 0 {
		params.VoteAverage = numericFromFloat(img.VoteAverage)
	}
	if img.VoteCount > 0 {
		vc := int32(img.VoteCount)
		params.VoteCount = &vc
	}
	if img.Blurhash != "" {
		params.Blurhash = &img.Blurhash
	}
	if img.Provider != "" {
		params.Provider = &img.Provider
	}
	if img.ProviderID != "" {
		params.ProviderID = &img.ProviderID
	}

	i, err := r.queries.CreateEpisodeImage(ctx, params)
	if err != nil {
		return err
	}
	img.ID = i.ID
	img.CreatedAt = i.CreatedAt
	return nil
}

// DeleteEpisodeImages deletes all images for an episode.
func (r *pgRepository) DeleteEpisodeImages(ctx context.Context, episodeID uuid.UUID) error {
	return r.queries.DeleteEpisodeImages(ctx, episodeID)
}

// Videos

// GetSeriesVideos retrieves videos for a series.
func (r *pgRepository) GetSeriesVideos(ctx context.Context, seriesID uuid.UUID) ([]Video, error) {
	rows, err := r.queries.GetSeriesVideos(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	videos := make([]Video, len(rows))
	for i, row := range rows {
		videos[i] = Video{
			ID:        row.ID,
			VideoType: row.VideoType,
			CreatedAt: row.CreatedAt,
		}
		if row.Name != nil {
			videos[i].Name = *row.Name
		}
		if row.Key != nil {
			videos[i].Key = *row.Key
		}
		if row.Site != nil {
			videos[i].Site = *row.Site
		}
		if row.Size != nil {
			videos[i].Size = int(*row.Size)
		}
		if row.Language != nil {
			videos[i].Language = *row.Language
		}
		videos[i].IsOfficial = row.IsOfficial
		if row.TmdbID != nil {
			videos[i].TmdbID = *row.TmdbID
		}
	}
	return videos, nil
}

// CreateSeriesVideo creates a series video.
func (r *pgRepository) CreateSeriesVideo(ctx context.Context, seriesID uuid.UUID, video *Video) error {
	params := db.CreateSeriesVideoParams{
		SeriesID:   seriesID,
		VideoType:  video.VideoType,
		IsOfficial: video.IsOfficial,
	}
	if video.Name != "" {
		params.Name = &video.Name
	}
	if video.Key != "" {
		params.Key = &video.Key
	}
	if video.Site != "" {
		params.Site = &video.Site
	}
	if video.Size > 0 {
		s := int32(video.Size)
		params.Size = &s
	}
	if video.Language != "" {
		params.Language = &video.Language
	}
	if video.TmdbID != "" {
		params.TmdbID = &video.TmdbID
	}

	v, err := r.queries.CreateSeriesVideo(ctx, params)
	if err != nil {
		return err
	}
	video.ID = v.ID
	video.CreatedAt = v.CreatedAt
	return nil
}

// DeleteSeriesVideos deletes all videos for a series.
func (r *pgRepository) DeleteSeriesVideos(ctx context.Context, seriesID uuid.UUID) error {
	return r.queries.DeleteSeriesVideos(ctx, seriesID)
}

// Helper functions

func fromDBNetwork(n *db.TvNetwork) *Network {
	if n == nil {
		return nil
	}
	network := &Network{
		ID:        n.ID,
		Name:      n.Name,
		CreatedAt: n.CreatedAt,
	}
	if n.LogoPath != nil {
		network.LogoPath = *n.LogoPath
	}
	if n.OriginCountry != nil {
		network.OriginCountry = *n.OriginCountry
	}
	if n.TmdbID != nil {
		network.TmdbID = int(*n.TmdbID)
	}
	return network
}

func fromDBSeriesImages(rows []db.SeriesImage) []Image {
	images := make([]Image, len(rows))
	for i, row := range rows {
		images[i] = Image{
			ID:        row.ID,
			ImageType: row.ImageType,
			URL:       row.Url,
			IsPrimary: row.IsPrimary,
			CreatedAt: row.CreatedAt,
		}
		if row.LocalPath != nil {
			images[i].LocalPath = *row.LocalPath
		}
		if row.Width != nil {
			images[i].Width = int(*row.Width)
		}
		if row.Height != nil {
			images[i].Height = int(*row.Height)
		}
		if row.AspectRatio.Valid {
			f, _ := row.AspectRatio.Float64Value()
			images[i].AspectRatio = f.Float64
		}
		if row.Language != nil {
			images[i].Language = *row.Language
		}
		if row.VoteAverage.Valid {
			f, _ := row.VoteAverage.Float64Value()
			images[i].VoteAverage = f.Float64
		}
		if row.VoteCount != nil {
			images[i].VoteCount = int(*row.VoteCount)
		}
		if row.Blurhash != nil {
			images[i].Blurhash = *row.Blurhash
		}
		if row.Provider != nil {
			images[i].Provider = *row.Provider
		}
		if row.ProviderID != nil {
			images[i].ProviderID = *row.ProviderID
		}
	}
	return images
}

func fromDBSeasonImages(rows []db.SeasonImage) []Image {
	images := make([]Image, len(rows))
	for i, row := range rows {
		images[i] = Image{
			ID:        row.ID,
			ImageType: row.ImageType,
			URL:       row.Url,
			IsPrimary: row.IsPrimary,
			CreatedAt: row.CreatedAt,
		}
		if row.LocalPath != nil {
			images[i].LocalPath = *row.LocalPath
		}
		if row.Width != nil {
			images[i].Width = int(*row.Width)
		}
		if row.Height != nil {
			images[i].Height = int(*row.Height)
		}
		if row.AspectRatio.Valid {
			f, _ := row.AspectRatio.Float64Value()
			images[i].AspectRatio = f.Float64
		}
		if row.Language != nil {
			images[i].Language = *row.Language
		}
		if row.VoteAverage.Valid {
			f, _ := row.VoteAverage.Float64Value()
			images[i].VoteAverage = f.Float64
		}
		if row.VoteCount != nil {
			images[i].VoteCount = int(*row.VoteCount)
		}
		if row.Blurhash != nil {
			images[i].Blurhash = *row.Blurhash
		}
		if row.Provider != nil {
			images[i].Provider = *row.Provider
		}
		if row.ProviderID != nil {
			images[i].ProviderID = *row.ProviderID
		}
	}
	return images
}

func fromDBEpisodeImages(rows []db.EpisodeImage) []Image {
	images := make([]Image, len(rows))
	for i, row := range rows {
		images[i] = Image{
			ID:        row.ID,
			ImageType: row.ImageType,
			URL:       row.Url,
			IsPrimary: row.IsPrimary,
			CreatedAt: row.CreatedAt,
		}
		if row.LocalPath != nil {
			images[i].LocalPath = *row.LocalPath
		}
		if row.Width != nil {
			images[i].Width = int(*row.Width)
		}
		if row.Height != nil {
			images[i].Height = int(*row.Height)
		}
		if row.AspectRatio.Valid {
			f, _ := row.AspectRatio.Float64Value()
			images[i].AspectRatio = f.Float64
		}
		if row.VoteAverage.Valid {
			f, _ := row.VoteAverage.Float64Value()
			images[i].VoteAverage = f.Float64
		}
		if row.VoteCount != nil {
			images[i].VoteCount = int(*row.VoteCount)
		}
		if row.Blurhash != nil {
			images[i].Blurhash = *row.Blurhash
		}
		if row.Provider != nil {
			images[i].Provider = *row.Provider
		}
		if row.ProviderID != nil {
			images[i].ProviderID = *row.ProviderID
		}
	}
	return images
}
