package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/justinndidit/stringAnalyzer/internal/dto"
)

type Envelope map[string]any

// func ToDTO(filters *dto.QueryParams) dto.QueryParams {
// 	return dto.QueryParams{
// 		IsPalindrome:      filters.IsPalindrome,
// 		MinLength:         filters.MinLength,
// 		MaxLength:         filters.MaxLength,
// 		WordCount:         filters.WordCount,
// 		ContainsCharacter: filters.ContainsCharacter,
// 	}
// }

func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "")

	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// func WriteResponse(data any) *Envelope {
// 	return &Envelope{
// 		"": data,
// 	}
// }

func CountWords(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	inWord := false
	count := 0

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				inWord = true
				count++
			}
		} else {
			inWord = false
		}
	}

	return count
}

func IsPalindrome(s string) bool {
	// Convert once to lowercase
	s = strings.ToLower(s)

	// Filter only alphanumeric runes
	var filtered []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			filtered = append(filtered, r)
		}
	}

	// Check palindrome using two-pointer technique
	i, j := 0, len(filtered)-1
	for i < j {
		if filtered[i] != filtered[j] {
			return false
		}
		i++
		j--
	}

	return true
}

func CountUniqueCharacters(s string) int {
	s = strings.ToLower(s)
	seen := make(map[rune]bool)

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			seen[r] = true
		}
	}

	return len(seen)
}

func CharacterFrequencyMap(s string) map[string]int {
	frequencies := make(map[string]int)

	// Convert the string to lowercase for consistency
	s = strings.ToLower(s)

	for _, char := range s {
		// Only count letters and digits
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			frequencies[string(char)]++
		}
	}

	return frequencies
}
func Hash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func CharacterCount(s string) int {
	return len([]rune(s))
}

func ParseNaturalLanguageQuery(query string) (*dto.FilterParams, map[string]interface{}, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	filters := &dto.FilterParams{}
	parsedFilters := make(map[string]interface{})

	// Check for palindrome
	if containsAny(query, []string{"palindrome", "palindromic"}) {
		t := true
		filters.IsPalindrome = &t
		parsedFilters["is_palindrome"] = true
	}

	// Check for word count
	if wordCount := extractWordCount(query); wordCount != nil {
		filters.WordCount = wordCount
		parsedFilters["word_count"] = *wordCount
	}

	// Check for length constraints
	if minLen := extractMinLength(query); minLen != nil {
		filters.MinLength = minLen
		parsedFilters["min_length"] = *minLen
	}

	if maxLen := extractMaxLength(query); maxLen != nil {
		filters.MaxLength = maxLen
		parsedFilters["max_length"] = *maxLen
	}

	// Check for contains character
	if char := extractContainsCharacter(query); char != nil {
		filters.ContainsCharacter = char
		parsedFilters["contains_character"] = *char
	}

	// If no filters were parsed, return error
	if len(parsedFilters) == 0 {
		return nil, nil, fmt.Errorf("could not extract any filters from query")
	}

	return filters, parsedFilters, nil
}

func extractWordCount(query string) *int {
	patterns := map[string]int{
		"single word": 1,
		"one word":    1,
		"two word":    2,
		"three word":  3,
		"four word":   4,
		"five word":   5,
	}

	for pattern, count := range patterns {
		if strings.Contains(query, pattern) {
			return &count
		}
	}

	// Check for numeric patterns like "3 words"
	re := regexp.MustCompile(`(\d+)\s*words?`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return &num
		}
	}

	return nil
}

func extractMinLength(query string) *int {
	// "longer than 10" → min_length = 11
	re := regexp.MustCompile(`longer than (\d+)`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			minLen := num + 1
			return &minLen
		}
	}

	// "at least 10 characters"
	re = regexp.MustCompile(`at least (\d+) characters?`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return &num
		}
	}

	// "minimum length 10"
	re = regexp.MustCompile(`minimum length (\d+)`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return &num
		}
	}

	return nil
}

// extractMaxLength looks for patterns like "shorter than X", "at most X characters"
func extractMaxLength(query string) *int {
	// "shorter than 10" → max_length = 9
	re := regexp.MustCompile(`shorter than (\d+)`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			maxLen := num - 1
			return &maxLen
		}
	}

	// "at most 10 characters"
	re = regexp.MustCompile(`at most (\d+) characters?`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return &num
		}
	}

	// "maximum length 10"
	re = regexp.MustCompile(`maximum length (\d+)`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return &num
		}
	}

	return nil
}

// extractContainsCharacter looks for patterns like "containing letter x", "with character a"
func extractContainsCharacter(query string) *string {
	// "containing the letter z"
	re := regexp.MustCompile(`containing (?:the )?letter ([a-z])`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		return &matches[1]
	}

	// "strings containing z"
	re = regexp.MustCompile(`containing ([a-z])(?:\s|$)`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		return &matches[1]
	}

	// "with character a"
	re = regexp.MustCompile(`with (?:the )?character ([a-z])`)
	if matches := re.FindStringSubmatch(query); len(matches) > 1 {
		return &matches[1]
	}

	// "first vowel" → 'a'
	if strings.Contains(query, "first vowel") {
		a := "a"
		return &a
	}

	// "last vowel" → 'u'
	if strings.Contains(query, "last vowel") {
		u := "u"
		return &u
	}

	return nil
}

// validateFilters checks for conflicting filter combinations
func ValidateFilters(filters *dto.FilterParams) error {
	// Check if min_length > max_length
	if filters.MinLength != nil && filters.MaxLength != nil {
		if *filters.MinLength > *filters.MaxLength {
			return fmt.Errorf("min_length (%d) cannot be greater than max_length (%d)",
				*filters.MinLength, *filters.MaxLength)
		}
	}

	// Check if word_count and length constraints are impossible
	// (e.g., single word but min_length = 100 might be unrealistic)
	// Add your domain-specific validation here

	return nil
}

// Helper function to check if string contains any of the given substrings
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
