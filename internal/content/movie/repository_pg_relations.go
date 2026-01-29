package movie

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/lusoris/revenge/internal/content/movie/db"
)

// Genres

// GetMovieGenres retrieves genres for a movie.
func (r *pgRepository) GetMovieGenres(ctx context.Context, movieID uuid.UUID) ([]Genre, error) {
	rows, err := r.queries.GetMovieGenres(ctx, movieID)
	if err != nil {
		return nil, err
	}

	genres := make([]Genre, len(rows))
	for i, row := range rows {
		genres[i] = Genre{
			ID:   row.ID,
			Name: row.Name,
			Slug: row.Slug,
		}
		if row.Description != nil {
			genres[i].Description = *row.Description
		}
	}
	return genres, nil
}

// LinkMovieGenre links a genre to a movie.
func (r *pgRepository) LinkMovieGenre(ctx context.Context, movieID, genreID uuid.UUID) error {
	return r.queries.LinkMovieGenre(ctx, db.LinkMovieGenreParams{
		MovieID: movieID,
		GenreID: genreID,
	})
}

// UnlinkMovieGenres removes all genre links from a movie.
func (r *pgRepository) UnlinkMovieGenres(ctx context.Context, movieID uuid.UUID) error {
	return r.queries.UnlinkMovieGenres(ctx, movieID)
}

// ListMoviesByGenre retrieves movies with a specific genre.
func (r *pgRepository) ListMoviesByGenre(ctx context.Context, genreID uuid.UUID, params ListParams) ([]*Movie, error) {
	rows, err := r.queries.ListMoviesByGenre(ctx, db.ListMoviesByGenreParams{
		GenreID: genreID,
		Limit:   int32(params.Limit),
		Offset:  int32(params.Offset),
	})
	if err != nil {
		return nil, err
	}

	movies := make([]*Movie, len(rows))
	for i, row := range rows {
		movies[i] = FromDBMovie(&row)
	}
	return movies, nil
}

// CountMoviesByGenre returns the number of movies with a specific genre.
func (r *pgRepository) CountMoviesByGenre(ctx context.Context, genreID uuid.UUID) (int64, error) {
	return r.queries.CountMoviesByGenre(ctx, genreID)
}

// Credits

// GetMovieCast retrieves cast members for a movie.
func (r *pgRepository) GetMovieCast(ctx context.Context, movieID uuid.UUID) ([]CastMember, error) {
	rows, err := r.queries.GetMovieCast(ctx, movieID)
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

// GetMovieCrew retrieves crew members for a movie.
func (r *pgRepository) GetMovieCrew(ctx context.Context, movieID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetMovieCrew(ctx, movieID)
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

// GetMovieDirectors retrieves directors for a movie.
func (r *pgRepository) GetMovieDirectors(ctx context.Context, movieID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetMovieDirectors(ctx, movieID)
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
		if row.PrimaryImageUrl != nil {
			crew[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			crew[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return crew, nil
}

// GetMovieWriters retrieves writers for a movie.
func (r *pgRepository) GetMovieWriters(ctx context.Context, movieID uuid.UUID) ([]CrewMember, error) {
	rows, err := r.queries.GetMovieWriters(ctx, movieID)
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
		if row.PrimaryImageUrl != nil {
			crew[i].PrimaryImageURL = *row.PrimaryImageUrl
		}
		if row.PrimaryImageBlurhash != nil {
			crew[i].PrimaryImageBlurhash = *row.PrimaryImageBlurhash
		}
	}
	return crew, nil
}

// CreateMovieCredit creates a credit for a movie.
func (r *pgRepository) CreateMovieCredit(ctx context.Context, movieID, personID uuid.UUID, role, character, department, job string, order int, isGuest bool, tmdbCreditID string) error {
	_, err := r.queries.CreateMovieCredit(ctx, db.CreateMovieCreditParams{
		MovieID:       movieID,
		PersonID:      personID,
		Role:          role,
		CharacterName: ptrString(character),
		Department:    ptrString(department),
		Job:           ptrString(job),
		BillingOrder:  int32(order),
		IsGuest:       isGuest,
		TmdbCreditID:  ptrString(tmdbCreditID),
	})
	return err
}

// DeleteMovieCredits deletes all credits for a movie.
func (r *pgRepository) DeleteMovieCredits(ctx context.Context, movieID uuid.UUID) error {
	return r.queries.DeleteMovieCredits(ctx, movieID)
}

// Images

// GetMovieImages retrieves images for a movie.
func (r *pgRepository) GetMovieImages(ctx context.Context, movieID uuid.UUID) ([]Image, error) {
	rows, err := r.queries.GetMovieImages(ctx, movieID)
	if err != nil {
		return nil, err
	}

	return fromDBImages(rows), nil
}

// GetMovieImagesByType retrieves images of a specific type for a movie.
func (r *pgRepository) GetMovieImagesByType(ctx context.Context, movieID uuid.UUID, imageType string) ([]Image, error) {
	rows, err := r.queries.GetMovieImagesByType(ctx, db.GetMovieImagesByTypeParams{
		MovieID:   movieID,
		ImageType: imageType,
	})
	if err != nil {
		return nil, err
	}

	return fromDBImages(rows), nil
}

// CreateMovieImage creates an image for a movie.
func (r *pgRepository) CreateMovieImage(ctx context.Context, movieID uuid.UUID, img *Image) error {
	var aspectRatio pgtype.Numeric
	if img.AspectRatio > 0 {
		aspectRatio = numericFromFloat(img.AspectRatio)
	}
	var voteAverage pgtype.Numeric
	if img.VoteAverage > 0 {
		voteAverage = numericFromFloat(img.VoteAverage)
	}

	row, err := r.queries.CreateMovieImage(ctx, db.CreateMovieImageParams{
		MovieID:     movieID,
		ImageType:   img.ImageType,
		Path:        img.Path,
		Blurhash:    ptrString(img.Blurhash),
		Width:       ptrInt32(int32(img.Width)),
		Height:      ptrInt32(int32(img.Height)),
		AspectRatio: aspectRatio,
		Language:    ptrString(img.Language),
		VoteAverage: voteAverage,
		VoteCount:   ptrInt32(int32(img.VoteCount)),
		IsPrimary:   img.IsPrimary,
		Source:      img.Source,
	})
	if err != nil {
		return err
	}
	img.ID = row.ID
	img.CreatedAt = row.CreatedAt
	return nil
}

// DeleteMovieImages deletes all images for a movie.
func (r *pgRepository) DeleteMovieImages(ctx context.Context, movieID uuid.UUID) error {
	return r.queries.DeleteMovieImages(ctx, movieID)
}

// Videos

// GetMovieVideos retrieves videos for a movie.
func (r *pgRepository) GetMovieVideos(ctx context.Context, movieID uuid.UUID) ([]Video, error) {
	rows, err := r.queries.GetMovieVideos(ctx, movieID)
	if err != nil {
		return nil, err
	}

	videos := make([]Video, len(rows))
	for i, row := range rows {
		videos[i] = Video{
			ID:        row.ID,
			VideoType: row.VideoType,
			Site:      row.Site,
			Key:       row.Key,
			CreatedAt: row.CreatedAt,
		}
		if row.Name != nil {
			videos[i].Name = *row.Name
		}
		if row.Language != nil {
			videos[i].Language = *row.Language
		}
		if row.Size != nil {
			videos[i].Size = int(*row.Size)
		}
	}
	return videos, nil
}

// CreateMovieVideo creates a video for a movie.
func (r *pgRepository) CreateMovieVideo(ctx context.Context, movieID uuid.UUID, video *Video) error {
	row, err := r.queries.CreateMovieVideo(ctx, db.CreateMovieVideoParams{
		MovieID:   movieID,
		VideoType: video.VideoType,
		Site:      video.Site,
		Key:       video.Key,
		Name:      ptrString(video.Name),
		Language:  ptrString(video.Language),
		Size:      ptrInt32(int32(video.Size)),
	})
	if err != nil {
		return err
	}
	video.ID = row.ID
	video.CreatedAt = row.CreatedAt
	return nil
}

// DeleteMovieVideos deletes all videos for a movie.
func (r *pgRepository) DeleteMovieVideos(ctx context.Context, movieID uuid.UUID) error {
	return r.queries.DeleteMovieVideos(ctx, movieID)
}

// Helper functions

func fromDBImages(rows []db.MovieImage) []Image {
	images := make([]Image, len(rows))
	for i, row := range rows {
		images[i] = Image{
			ID:        row.ID,
			ImageType: row.ImageType,
			Path:      row.Path,
			IsPrimary: row.IsPrimary,
			Source:    row.Source,
			CreatedAt: row.CreatedAt,
		}
		if row.Blurhash != nil {
			images[i].Blurhash = *row.Blurhash
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
	}
	return images
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrInt32(i int32) *int32 {
	if i == 0 {
		return nil
	}
	return &i
}

func numericFromFloat(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(f)
	return n
}
