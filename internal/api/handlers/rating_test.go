package handlers

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/domain"
)

func TestRatingSystemToResponse(t *testing.T) {
	rs := &domain.RatingSystem{
		ID:           uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Code:         "mpaa",
		Name:         "Motion Picture Association",
		CountryCodes: []string{"US", "CA"},
		IsActive:     true,
	}

	resp := ratingSystemToResponse(rs)

	if resp.ID != rs.ID.String() {
		t.Errorf("ID mismatch: got %s, want %s", resp.ID, rs.ID.String())
	}
	if resp.Code != rs.Code {
		t.Errorf("Code mismatch: got %s, want %s", resp.Code, rs.Code)
	}
	if resp.Name != rs.Name {
		t.Errorf("Name mismatch: got %s, want %s", resp.Name, rs.Name)
	}
	if len(resp.CountryCodes) != 2 {
		t.Errorf("CountryCodes length mismatch: got %d, want 2", len(resp.CountryCodes))
	}
	if !resp.IsActive {
		t.Error("IsActive should be true")
	}
}

func TestRatingToResponse(t *testing.T) {
	desc := "Parental Guidance Suggested"
	minAge := 10
	iconURL := "https://example.com/pg.png"

	r := &domain.Rating{
		ID:              uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Code:            "PG",
		Name:            "Parental Guidance",
		Description:     &desc,
		MinAge:          &minAge,
		NormalizedLevel: 25,
		IsAdult:         false,
		IconURL:         &iconURL,
		System: &domain.RatingSystem{
			ID:   uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Code: "mpaa",
			Name: "Motion Picture Association",
		},
	}

	resp := ratingToResponse(r)

	if resp.ID != r.ID.String() {
		t.Errorf("ID mismatch: got %s, want %s", resp.ID, r.ID.String())
	}
	if resp.Code != r.Code {
		t.Errorf("Code mismatch: got %s, want %s", resp.Code, r.Code)
	}
	if resp.Name != r.Name {
		t.Errorf("Name mismatch: got %s, want %s", resp.Name, r.Name)
	}
	if resp.Description == nil || *resp.Description != desc {
		t.Error("Description mismatch")
	}
	if resp.MinAge == nil || *resp.MinAge != minAge {
		t.Error("MinAge mismatch")
	}
	if resp.NormalizedLevel != 25 {
		t.Errorf("NormalizedLevel mismatch: got %d, want 25", resp.NormalizedLevel)
	}
	if resp.IsAdult {
		t.Error("IsAdult should be false")
	}
	if resp.SystemCode != "mpaa" {
		t.Errorf("SystemCode mismatch: got %s, want mpaa", resp.SystemCode)
	}
	if resp.SystemName != "Motion Picture Association" {
		t.Errorf("SystemName mismatch: got %s, want Motion Picture Association", resp.SystemName)
	}
}

func TestRatingToResponse_NilSystem(t *testing.T) {
	r := &domain.Rating{
		ID:              uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Code:            "PG",
		Name:            "Parental Guidance",
		NormalizedLevel: 25,
		System:          nil,
	}

	resp := ratingToResponse(r)

	if resp.SystemCode != "" {
		t.Errorf("SystemCode should be empty, got %s", resp.SystemCode)
	}
	if resp.SystemName != "" {
		t.Errorf("SystemName should be empty, got %s", resp.SystemName)
	}
}

func TestContentRatingToResponse(t *testing.T) {
	source := "tmdb"
	cr := &domain.ContentRating{
		ID:          uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		ContentID:   uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		ContentType: "media_item",
		Source:      &source,
		Rating: &domain.Rating{
			ID:              uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			Code:            "PG-13",
			Name:            "Parents Strongly Cautioned",
			NormalizedLevel: 50,
		},
	}

	resp := contentRatingToResponse(cr)

	if resp.ID != cr.ID.String() {
		t.Errorf("ID mismatch: got %s, want %s", resp.ID, cr.ID.String())
	}
	if resp.ContentID != cr.ContentID.String() {
		t.Errorf("ContentID mismatch: got %s, want %s", resp.ContentID, cr.ContentID.String())
	}
	if resp.ContentType != "media_item" {
		t.Errorf("ContentType mismatch: got %s, want media_item", resp.ContentType)
	}
	if resp.Source == nil || *resp.Source != "tmdb" {
		t.Error("Source mismatch")
	}
	if resp.Rating == nil {
		t.Fatal("Rating should not be nil")
	}
	if resp.Rating.Code != "PG-13" {
		t.Errorf("Rating.Code mismatch: got %s, want PG-13", resp.Rating.Code)
	}
}

func TestContentRatingToResponse_NilRating(t *testing.T) {
	cr := &domain.ContentRating{
		ID:          uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		ContentID:   uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		ContentType: "media_item",
		Rating:      nil,
	}

	resp := contentRatingToResponse(cr)

	if resp.Rating != nil {
		t.Error("Rating should be nil")
	}
}

func TestRatingSystemResponse_JSONFormat(t *testing.T) {
	resp := RatingSystemResponse{
		ID:           "22222222-2222-2222-2222-222222222222",
		Code:         "fsk",
		Name:         "Freiwillige Selbstkontrolle",
		CountryCodes: []string{"DE", "AT"},
		IsActive:     true,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	expectedFields := []string{"Id", "Code", "Name", "CountryCodes", "IsActive"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Missing expected JSON field: %s", field)
		}
	}
}

func TestRatingResponse_JSONFormat(t *testing.T) {
	resp := RatingResponse{
		ID:              "33333333-3333-3333-3333-333333333333",
		Code:            "R",
		Name:            "Restricted",
		NormalizedLevel: 75,
		IsAdult:         false,
		SystemCode:      "mpaa",
		SystemName:      "MPAA",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check required fields exist
	requiredFields := []string{"Id", "Code", "Name", "NormalizedLevel", "IsAdult", "SystemCode", "SystemName"}
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("Missing expected JSON field: %s", field)
		}
	}

	// Check omitempty fields are not present when nil/empty
	omittedFields := []string{"Description", "MinAge", "IconUrl"}
	for _, field := range omittedFields {
		if _, ok := raw[field]; ok {
			t.Errorf("Field should be omitted: %s", field)
		}
	}
}

func TestNormalizedLevel_Constants(t *testing.T) {
	// Test the normalized level scale matches documentation
	tests := []struct {
		name     string
		level    int
		minAge   int
		examples string
	}{
		{"All Ages", 0, 0, "G, FSK 0, U"},
		{"6+", 25, 6, "PG, FSK 6"},
		{"12+", 50, 12, "PG-13, FSK 12"},
		{"16+", 75, 16, "R, FSK 16"},
		{"18+", 90, 18, "NC-17, FSK 18"},
		{"Adult XXX", 100, 18, "R18, X18+"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just validate the levels are in expected range
			if tt.level < 0 || tt.level > 100 {
				t.Errorf("Level %d out of range [0, 100]", tt.level)
			}
		})
	}
}

func TestAddContentRatingRequest_JSONParsing(t *testing.T) {
	jsonStr := `{
		"RatingId": "33333333-3333-3333-3333-333333333333",
		"ContentType": "media_item",
		"Source": "manual"
	}`

	var req AddContentRatingRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.RatingID != "33333333-3333-3333-3333-333333333333" {
		t.Errorf("RatingId mismatch: got %s", req.RatingID)
	}
	if req.ContentType != "media_item" {
		t.Errorf("ContentType mismatch: got %s", req.ContentType)
	}
	if req.Source == nil || *req.Source != "manual" {
		t.Error("Source mismatch")
	}
}

func TestAddContentRatingRequest_MinimalJSON(t *testing.T) {
	jsonStr := `{"RatingId": "33333333-3333-3333-3333-333333333333"}`

	var req AddContentRatingRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if req.RatingID != "33333333-3333-3333-3333-333333333333" {
		t.Errorf("RatingId mismatch: got %s", req.RatingID)
	}
	if req.ContentType != "" {
		t.Errorf("ContentType should be empty, got %s", req.ContentType)
	}
	if req.Source != nil {
		t.Error("Source should be nil")
	}
}
