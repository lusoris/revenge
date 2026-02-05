package metadata

import (
	"strconv"
	"strings"
	"time"
)

// ParseReleaseDate parses a date string in ISO format (YYYY-MM-DD).
// Returns nil if the string is empty or invalid.
func ParseReleaseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}

	return &t
}

// ExtractYearFromDate extracts the year from an ISO date string (YYYY-MM-DD).
// Returns nil if the string is empty or invalid.
func ExtractYearFromDate(dateStr string) *int {
	if dateStr == "" {
		return nil
	}

	parts := strings.Split(dateStr, "-")
	if len(parts) == 0 {
		return nil
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}

	return &year
}

// ParseOptionalString returns nil if the string pointer is nil or empty.
func ParseOptionalString(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

// ParseOptionalStringValue returns nil if the string is empty, otherwise returns a pointer.
func ParseOptionalStringValue(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// SafeIntToInt32 safely converts an int to int32, clamping if necessary.
func SafeIntToInt32(i int) int32 {
	if i > 2147483647 {
		return 2147483647
	}
	if i < -2147483648 {
		return -2147483648
	}
	return int32(i)
}

// ParseOptionalInt32 returns nil if i is 0, otherwise returns a pointer to int32.
func ParseOptionalInt32(i int) *int32 {
	if i == 0 {
		return nil
	}
	val := SafeIntToInt32(i)
	return &val
}

// ParseOptionalInt32Ptr returns nil if i is nil or 0.
func ParseOptionalInt32Ptr(i *int) *int32 {
	if i == nil || *i == 0 {
		return nil
	}
	val := SafeIntToInt32(*i)
	return &val
}

// ParseOptionalInt64Ptr returns nil if i is nil or 0.
func ParseOptionalInt64Ptr(i *int64) *int64 {
	if i == nil || *i == 0 {
		return nil
	}
	return i
}

// AgeRatingSystem represents a regional age rating system.
type AgeRatingSystem string

const (
	AgeRatingMPAA  AgeRatingSystem = "MPAA"  // US
	AgeRatingFSK   AgeRatingSystem = "FSK"   // Germany
	AgeRatingBBFC  AgeRatingSystem = "BBFC"  // UK
	AgeRatingCNC   AgeRatingSystem = "CNC"   // France
	AgeRatingEirin AgeRatingSystem = "Eirin" // Japan
	AgeRatingKMRB  AgeRatingSystem = "KMRB"  // South Korea
	AgeRatingDJCTQ AgeRatingSystem = "DJCTQ" // Brazil
	AgeRatingACB   AgeRatingSystem = "ACB"   // Australia
)

// GetAgeRatingSystem returns the rating system name for a given country code.
func GetAgeRatingSystem(countryISO string) AgeRatingSystem {
	switch countryISO {
	case "US":
		return AgeRatingMPAA
	case "DE":
		return AgeRatingFSK
	case "GB":
		return AgeRatingBBFC
	case "FR":
		return AgeRatingCNC
	case "JP":
		return AgeRatingEirin
	case "KR":
		return AgeRatingKMRB
	case "BR":
		return AgeRatingDJCTQ
	case "AU":
		return AgeRatingACB
	default:
		// Use country code as fallback
		return AgeRatingSystem(countryISO)
	}
}

// LanguageToISO converts a TMDb language code (en-US) to ISO 639-1 (en).
func LanguageToISO(lang string) string {
	if len(lang) > 2 && lang[2] == '-' {
		return lang[:2]
	}
	return lang
}

// ISOToLanguage converts an ISO 639-1 code (en) to a common TMDb format (en-US).
// This is a best-effort mapping for common languages.
func ISOToLanguage(iso string) string {
	switch iso {
	case "en":
		return "en-US"
	case "de":
		return "de-DE"
	case "fr":
		return "fr-FR"
	case "es":
		return "es-ES"
	case "it":
		return "it-IT"
	case "pt":
		return "pt-BR"
	case "ja":
		return "ja-JP"
	case "ko":
		return "ko-KR"
	case "zh":
		return "zh-CN"
	case "ru":
		return "ru-RU"
	default:
		return iso
	}
}
