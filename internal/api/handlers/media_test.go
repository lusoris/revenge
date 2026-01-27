package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/domain"
)

func TestItemToResponse(t *testing.T) {
	year := 2024
	overview := "A great movie about testing"
	runtime := int64(7200000000000) // 2 hours in ticks
	rating := 8.5
	sortName := "Movie, The"
	parentID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	path := "/media/movies/the_movie.mkv"
	premiereDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	container := "mkv"
	width := 1920
	height := 1080

	item := &domain.MediaItem{
		ID:              uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		LibraryID:       uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		ParentID:        &parentID,
		Type:            domain.MediaTypeMovie,
		Name:            "The Movie",
		SortName:        &sortName,
		Path:            &path,
		Overview:        &overview,
		Year:            &year,
		PremiereDate:    &premiereDate,
		RuntimeTicks:    &runtime,
		CommunityRating: &rating,
		Genres:          []string{"Action", "Comedy"},
		Tags:            []string{"featured", "new"},
		Studios:         []string{"Universal", "Paramount"},
		ProviderIDs:     map[string]string{"imdb": "tt1234567", "tmdb": "12345"},
		Container:       &container,
		Width:           &width,
		Height:          &height,
	}

	resp := itemToResponse(item)

	if resp.ID != item.ID.String() {
		t.Errorf("ID mismatch: got %s, want %s", resp.ID, item.ID.String())
	}
	if resp.Name != item.Name {
		t.Errorf("Name mismatch: got %s, want %s", resp.Name, item.Name)
	}
	if resp.Type != "Movie" {
		t.Errorf("Type mismatch: got %s, want Movie", resp.Type)
	}
	if resp.MediaType != "Video" {
		t.Errorf("MediaType mismatch: got %s, want Video", resp.MediaType)
	}
	if resp.SortName != sortName {
		t.Errorf("SortName mismatch: got %s, want %s", resp.SortName, sortName)
	}
	if resp.Overview == nil || *resp.Overview != overview {
		t.Error("Overview mismatch")
	}
	if resp.ProductionYear == nil || *resp.ProductionYear != year {
		t.Error("ProductionYear mismatch")
	}
	if resp.RunTimeTicks == nil || *resp.RunTimeTicks != runtime {
		t.Error("RunTimeTicks mismatch")
	}
	if resp.CommunityRating == nil || *resp.CommunityRating != rating {
		t.Error("CommunityRating mismatch")
	}
	if len(resp.Genres) != 2 {
		t.Errorf("Genres count mismatch: got %d, want 2", len(resp.Genres))
	}
	if len(resp.Studios) != 2 {
		t.Errorf("Studios count mismatch: got %d, want 2", len(resp.Studios))
	}
	if resp.IsFolder {
		t.Error("IsFolder should be false for Movie")
	}
	if resp.ParentID == nil || *resp.ParentID != parentID.String() {
		t.Error("ParentID mismatch")
	}
}

func TestItemToResponse_Episode(t *testing.T) {
	seasonNum := 2
	episodeNum := 5

	item := &domain.MediaItem{
		ID:            uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		LibraryID:     uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Type:          domain.MediaTypeEpisode,
		Name:          "Episode Title",
		SeasonNumber:  &seasonNum,
		EpisodeNumber: &episodeNum,
	}

	resp := itemToResponse(item)

	if resp.Type != "Episode" {
		t.Errorf("Type mismatch: got %s, want Episode", resp.Type)
	}
	if resp.ParentIndexNumber == nil || *resp.ParentIndexNumber != seasonNum {
		t.Errorf("ParentIndexNumber (season) mismatch: got %v, want %d", resp.ParentIndexNumber, seasonNum)
	}
	if resp.IndexNumber == nil || *resp.IndexNumber != episodeNum {
		t.Errorf("IndexNumber (episode) mismatch: got %v, want %d", resp.IndexNumber, episodeNum)
	}
}

func TestItemToResponse_Audio(t *testing.T) {
	albumArtist := "Test Artist"
	trackNum := 5

	item := &domain.MediaItem{
		ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		LibraryID:   uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Type:        domain.MediaTypeAudio,
		Name:        "Track Name",
		AlbumArtist: &albumArtist,
		TrackNumber: &trackNum,
	}

	resp := itemToResponse(item)

	if resp.Type != "Audio" {
		t.Errorf("Type mismatch: got %s, want Audio", resp.Type)
	}
	if resp.MediaType != "Audio" {
		t.Errorf("MediaType mismatch: got %s, want Audio", resp.MediaType)
	}
	if resp.AlbumArtist == nil || *resp.AlbumArtist != albumArtist {
		t.Error("AlbumArtist mismatch")
	}
	if resp.IndexNumber == nil || *resp.IndexNumber != trackNum {
		t.Errorf("IndexNumber (track) mismatch: got %v, want %d", resp.IndexNumber, trackNum)
	}
}

func TestItemToResponse_Container(t *testing.T) {
	tests := []struct {
		mediaType domain.MediaType
		isFolder  bool
	}{
		{domain.MediaTypeMovie, false},
		{domain.MediaTypeEpisode, false},
		{domain.MediaTypeSeries, true},
		{domain.MediaTypeSeason, true},
		{domain.MediaTypeAlbum, true},
		{domain.MediaTypePlaylist, true},
		{domain.MediaTypeBoxSet, true},
		{domain.MediaTypeFolder, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.mediaType), func(t *testing.T) {
			item := &domain.MediaItem{
				ID:        uuid.New(),
				LibraryID: uuid.New(),
				Type:      tt.mediaType,
				Name:      "Test",
			}
			resp := itemToResponse(item)
			if resp.IsFolder != tt.isFolder {
				t.Errorf("IsFolder mismatch for %s: got %v, want %v", tt.mediaType, resp.IsFolder, tt.isFolder)
			}
		})
	}
}

func TestRevengeItemType(t *testing.T) {
	tests := []struct {
		input    domain.MediaType
		expected string
	}{
		{domain.MediaTypeMovie, "Movie"},
		{domain.MediaTypeEpisode, "Episode"},
		{domain.MediaTypeSeries, "Series"},
		{domain.MediaTypeSeason, "Season"},
		{domain.MediaTypeAudio, "Audio"},
		{domain.MediaTypeAlbum, "MusicAlbum"},
		{domain.MediaTypeArtist, "MusicArtist"},
		{domain.MediaTypeMusicVideo, "MusicVideo"},
		{domain.MediaTypePhoto, "Photo"},
		{domain.MediaTypeBook, "Book"},
		{domain.MediaTypeAudiobook, "AudioBook"},
		{domain.MediaTypePodcast, "Podcast"},
		{domain.MediaTypeBoxSet, "BoxSet"},
		{domain.MediaTypePlaylist, "Playlist"},
		{domain.MediaTypeFolder, "Folder"},
		{domain.MediaTypeTrailer, "Trailer"},
		{domain.MediaTypeChannel, "Channel"},
		{domain.MediaTypeProgram, "Program"},
		{domain.MediaTypeRecording, "Recording"},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got := jellfinItemType(tt.input)
			if got != tt.expected {
				t.Errorf("jellfinItemType(%s) = %s, want %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRevengeMediaType(t *testing.T) {
	tests := []struct {
		input    domain.MediaType
		expected string
	}{
		{domain.MediaTypeMovie, "Video"},
		{domain.MediaTypeEpisode, "Video"},
		{domain.MediaTypeMusicVideo, "Video"},
		{domain.MediaTypeTrailer, "Video"},
		{domain.MediaTypeRecording, "Video"},
		{domain.MediaTypeAudio, "Audio"},
		{domain.MediaTypePodcastEpisode, "Audio"},
		{domain.MediaTypePhoto, "Photo"},
		{domain.MediaTypeBook, "Book"},
		{domain.MediaTypeSeries, ""},  // Containers don't have media type
		{domain.MediaTypeAlbum, ""},
		{domain.MediaTypePlaylist, ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got := revengeMediaType(tt.input)
			if got != tt.expected {
				t.Errorf("revengeMediaType(%s) = %s, want %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRevengeTypeToMediaType(t *testing.T) {
	tests := []struct {
		input    string
		expected domain.MediaType
	}{
		{"Movie", domain.MediaTypeMovie},
		{"movie", domain.MediaTypeMovie}, // Case insensitive
		{"Episode", domain.MediaTypeEpisode},
		{"Series", domain.MediaTypeSeries},
		{"Season", domain.MediaTypeSeason},
		{"Audio", domain.MediaTypeAudio},
		{"MusicAlbum", domain.MediaTypeAlbum},
		{"MusicArtist", domain.MediaTypeArtist},
		{"MusicVideo", domain.MediaTypeMusicVideo},
		{"Photo", domain.MediaTypePhoto},
		{"PhotoAlbum", domain.MediaTypePhotoAlbum},
		{"Book", domain.MediaTypeBook},
		{"AudioBook", domain.MediaTypeAudiobook},
		{"Podcast", domain.MediaTypePodcast},
		{"BoxSet", domain.MediaTypeBoxSet},
		{"Playlist", domain.MediaTypePlaylist},
		{"Folder", domain.MediaTypeFolder},
		{"Trailer", domain.MediaTypeTrailer},
		{"Channel", domain.MediaTypeChannel},
		{"Program", domain.MediaTypeProgram},
		{"Recording", domain.MediaTypeRecording},
		{"Unknown", domain.MediaTypeMovie}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := revengeTypeToMediaType(tt.input)
			if got != tt.expected {
				t.Errorf("revengeTypeToMediaType(%s) = %s, want %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestItemResponse_JSONFormat(t *testing.T) {
	resp := ItemResponse{
		ID:             "11111111-1111-1111-1111-111111111111",
		Name:           "Test Movie",
		Type:           "Movie",
		MediaType:      "Video",
		Genres:         []string{"Action"},
		Tags:           []string{"featured"},
		IsFolder:       false,
		ProviderIDs:    map[string]string{"imdb": "tt1234567"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check required fields match Revenge API format
	requiredFields := []string{"Id", "Name", "Type", "IsFolder"}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Missing expected JSON field: %s", field)
		}
	}

	// Verify ProviderIds is correct (Revenge uses "ProviderIds" not "ProviderIDs")
	if _, ok := raw["ProviderIds"]; !ok {
		t.Error("Missing ProviderIds field")
	}
}

func TestItemsResponse_JSONFormat(t *testing.T) {
	resp := ItemsResponse{
		Items:            []ItemResponse{{ID: "1", Name: "Test", Type: "Movie"}},
		TotalRecordCount: 100,
		StartIndex:       0,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	requiredFields := []string{"Items", "TotalRecordCount", "StartIndex"}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Missing expected JSON field: %s", field)
		}
	}
}

func TestSearchHintResponse_JSONFormat(t *testing.T) {
	resp := SearchHintResponse{
		ItemID:    "11111111-1111-1111-1111-111111111111",
		ID:        "11111111-1111-1111-1111-111111111111",
		Name:      "Test Movie",
		Type:      "Movie",
		MediaType: "Video",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Revenge search hints have both ItemId and Id
	if _, ok := raw["ItemId"]; !ok {
		t.Error("Missing ItemId field")
	}
	if _, ok := raw["Id"]; !ok {
		t.Error("Missing Id field")
	}
}

func TestUserDataResponse_JSONFormat(t *testing.T) {
	played := 75.5
	lastPlayed := "2024-01-15T10:30:00Z"
	unplayed := 5

	resp := UserDataResponse{
		PlaybackPositionTicks: 12345678900,
		PlayCount:             3,
		IsFavorite:            true,
		Played:                false,
		UnplayedItemCount:     &unplayed,
		PlayedPercentage:      &played,
		LastPlayedDate:        &lastPlayed,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	requiredFields := []string{
		"PlaybackPositionTicks", "PlayCount", "IsFavorite", "Played",
		"UnplayedItemCount", "PlayedPercentage", "LastPlayedDate",
	}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Missing expected JSON field: %s", field)
		}
	}
}
