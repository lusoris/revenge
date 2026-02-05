package matcher

import (
	"strings"
	"unicode"
)

// LevenshteinDistance calculates the edit distance between two strings.
// The edit distance is the minimum number of single-character edits
// (insertions, deletions, or substitutions) required to transform one
// string into the other.
func LevenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Convert to runes for proper Unicode handling
	r1 := []rune(s1)
	r2 := []rune(s2)

	// Create distance matrix
	matrix := make([][]int, len(r1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(r2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill in the matrix
	for i := 1; i <= len(r1); i++ {
		for j := 1; j <= len(r2); j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(r1)][len(r2)]
}

// NormalizedSimilarity calculates the similarity between two strings
// as a value between 0.0 (completely different) and 1.0 (identical).
// It uses Levenshtein distance normalized by the maximum string length.
func NormalizedSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	distance := LevenshteinDistance(s1, s2)
	maxLen := max(len(s1), len(s2))

	if maxLen == 0 {
		return 1.0
	}

	return 1.0 - float64(distance)/float64(maxLen)
}

// TitleSimilarity calculates similarity between two titles,
// normalizing them first (lowercase, remove articles, punctuation).
func TitleSimilarity(title1, title2 string) float64 {
	// Normalize both titles
	norm1 := normalizeForComparison(title1)
	norm2 := normalizeForComparison(title2)

	return NormalizedSimilarity(norm1, norm2)
}

// normalizeForComparison prepares a title for comparison
func normalizeForComparison(title string) string {
	// Lowercase
	title = strings.ToLower(title)

	// Remove leading articles
	articles := []string{"the ", "a ", "an "}
	for _, article := range articles {
		if strings.HasPrefix(title, article) {
			title = strings.TrimPrefix(title, article)
			break
		}
	}

	// Remove non-alphanumeric characters except spaces
	var normalized strings.Builder
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			normalized.WriteRune(r)
		}
	}

	// Collapse multiple spaces
	result := normalized.String()
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}

	return strings.TrimSpace(result)
}

// YearMatch returns a score for how well two years match.
// Returns 1.0 for exact match, 0.5 for ±1 year, 0.0 otherwise.
func YearMatch(year1, year2 *int) float64 {
	if year1 == nil || year2 == nil {
		return 0.0
	}

	diff := abs(*year1 - *year2)
	switch diff {
	case 0:
		return 1.0
	case 1:
		return 0.5
	default:
		return 0.0
	}
}

// YearMatchInt is like YearMatch but takes int values directly
func YearMatchInt(year1, year2 int) float64 {
	diff := abs(year1 - year2)
	switch diff {
	case 0:
		return 1.0
	case 1:
		return 0.5
	default:
		return 0.0
	}
}

// abs returns the absolute value of an integer
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ConfidenceScore combines multiple weighted scores into a single confidence value.
// Each score should be between 0.0 and 1.0.
type ConfidenceScore struct {
	scores []weightedScore
}

type weightedScore struct {
	score  float64
	weight float64
}

// NewConfidenceScore creates a new confidence score calculator
func NewConfidenceScore() *ConfidenceScore {
	return &ConfidenceScore{}
}

// Add adds a weighted score component
func (c *ConfidenceScore) Add(score, weight float64) *ConfidenceScore {
	c.scores = append(c.scores, weightedScore{score: score, weight: weight})
	return c
}

// AddBonus adds a bonus (or penalty if negative) to the final score
func (c *ConfidenceScore) AddBonus(bonus float64) *ConfidenceScore {
	// Store as score with weight 1.0 but mark as bonus
	c.scores = append(c.scores, weightedScore{score: bonus, weight: 0})
	return c
}

// Calculate returns the weighted average confidence score, clamped to [0, 1]
func (c *ConfidenceScore) Calculate() float64 {
	if len(c.scores) == 0 {
		return 0.0
	}

	var totalWeight float64
	var weightedSum float64
	var bonuses float64

	for _, ws := range c.scores {
		if ws.weight == 0 {
			// This is a bonus
			bonuses += ws.score
		} else {
			totalWeight += ws.weight
			weightedSum += ws.score * ws.weight
		}
	}

	var result float64
	if totalWeight > 0 {
		result = weightedSum / totalWeight
	}

	// Add bonuses
	result += bonuses

	// Clamp to [0, 1]
	if result < 0 {
		return 0
	}
	if result > 1 {
		return 1
	}
	return result
}
