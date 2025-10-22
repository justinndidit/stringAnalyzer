package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/justinndidit/stringAnalyzer/internal/database"
	"github.com/justinndidit/stringAnalyzer/internal/dto"
	"github.com/justinndidit/stringAnalyzer/internal/errs"
	"github.com/justinndidit/stringAnalyzer/internal/repository"

	"github.com/justinndidit/stringAnalyzer/internal/util"
	"github.com/rs/zerolog"
)

type StringAnalyzerHandler struct {
	logger *zerolog.Logger
	db     *database.Database
	repo   *repository.StringRepository
}

func NewStringAnalyzerHandler(logger *zerolog.Logger, db *database.Database, repo *repository.StringRepository) *StringAnalyzerHandler {
	return &StringAnalyzerHandler{
		logger: logger,
		db:     db,
		repo:   repo,
	}
}

func (s *StringAnalyzerHandler) UploadString(w http.ResponseWriter, r *http.Request) {
	var body dto.UploadString
	defer r.Body.Close()

	if r.ContentLength == 0 {
		s.logger.Error().Msg("request body is empty!")

		rb := &util.Envelope{"message": "Invalid request body or missing \"value\" field"}
		util.WriteJson(w, http.StatusBadRequest, *rb)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.logger.Error().Msg(fmt.Sprintf("Error decoding request body: %v", err))

		var typeErr *errs.InvalidTypeError
		switch {
		case errors.As(err, &typeErr):
			rb := &util.Envelope{"message": "Invalid data type for \"value\" (must be string)"}
			util.WriteJson(w, http.StatusBadRequest, *rb)

		default:
			rb := &util.Envelope{"meesage": "Invalid request body or missing \"value\" field"}
			util.WriteJson(w, http.StatusBadRequest, *rb)
		}
		return
	}

	record, err := s.repo.GetStringByValue(r.Context(), body.Value)

	if err != nil {
		s.logger.Error().Msg(fmt.Sprintf("error getting string from database: %v", err))
		rb := &util.Envelope{"message": "Something went wrong"}
		util.WriteJson(w, http.StatusInternalServerError, *rb)
	}

	if record != nil {
		s.logger.Info().Msg("string already exists!")
		rb := &util.Envelope{"message": "String already exists in the system"}

		util.WriteJson(w, http.StatusConflict, *rb)
		return
	}

	payload := &dto.CreateString{
		StringValue:      body.Value,
		IsPalindrome:     util.IsPalindrome(body.Value),
		UniqueCharacters: util.CountUniqueCharacters(body.Value),
		WordCount:        util.CountWords(body.Value),
		Hash:             util.Hash(body.Value),
		Length:           util.CharacterCount(body.Value),
	}

	newString, err := s.repo.CreateString(r.Context(), payload)

	if err != nil {
		s.logger.Error().Err(err).Msg("error creating new string")
		rb := &util.Envelope{"message": "Something went wrong"}
		util.WriteJson(w, http.StatusInternalServerError, *rb)
	}

	properties := map[string]any{
		"length":                  newString.Length,
		"is_palindrome":           newString.IsPalindrome,
		"unique_characters":       newString.UniqueCharacters,
		"word_count":              newString.WordCount,
		"sha256_hash":             newString.Hash,
		"character_frequency_map": util.CharacterFrequencyMap(body.Value),
	}
	rb := &util.Envelope{
		"id":         newString.Hash,
		"value":      newString.StringValue,
		"properties": properties,
		"created_at": newString.CreatedAt,
	}
	util.WriteJson(w, http.StatusCreated, *rb)
}

func (s *StringAnalyzerHandler) GetString(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "string_value")

	record, err := s.repo.GetStringByValue(r.Context(), param)

	if err != nil {
		s.logger.Error().Err(err).Msg("Invalid query param")
		rb := &util.Envelope{"message": "Something went wrong!"}
		util.WriteJson(w, http.StatusInternalServerError, *rb)
	}

	if record == nil {
		rb := &util.Envelope{"message": "String does not exist in the system"}
		util.WriteJson(w, http.StatusNotFound, *rb)
		return
	}

	properties := map[string]any{
		"length":                  record.Length,
		"is_palindrome":           record.IsPalindrome,
		"unique_characters":       record.UniqueCharacters,
		"word_count":              record.WordCount,
		"sha256_hash":             record.Hash,
		"character_frequency_map": util.CharacterFrequencyMap(param),
	}

	rb := &util.Envelope{
		"id":         record.Hash,
		"value":      record.StringValue,
		"properties": properties,
		"created_at": record.CreatedAt,
	}

	util.WriteJson(w, http.StatusOK, *rb)

}

func (s *StringAnalyzerHandler) GetFilteredStrings(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var (
		isPalindrome      *bool
		minLength         *int
		maxLength         *int
		wordCount         *int
		containsCharacter string
	)

	// Parse is_palindrome
	if v := query.Get("is_palindrome"); v != "" {
		b, _ := strconv.ParseBool(v)
		isPalindrome = &b
	}

	// Parse min_length
	if v := query.Get("min_length"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			minLength = &i
		}
	}

	// Parse max_length
	if v := query.Get("max_length"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			maxLength = &i
		}
	}

	// Parse word_count
	if v := query.Get("word_count"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			wordCount = &i
		}
	}

	// Parse contains_character
	containsCharacter = query.Get("contains_character")

	params := dto.QueryParams{
		IsPalindrome:      isPalindrome,
		MinLength:         minLength,
		MaxLength:         maxLength,
		WordCount:         wordCount,
		ContainsCharacter: containsCharacter,
	}

	// ✅ Validate inputs
	if err := validator.New().Struct(params); err != nil {
		s.logger.Error().Err(err).Msg("error validating params")
		rb := &util.Envelope{"message": "Invalid query parameters"}
		util.WriteJson(w, http.StatusBadRequest, *rb)
		return
	}

	// 🔍 Fetch filtered records
	records, err := s.repo.GetFilteredStrings(r.Context(), params)
	if err != nil {
		s.logger.Error().Err(err).Msg("error fetching records")
		rb := &util.Envelope{"message": "Something went wrong!"}
		util.WriteJson(w, http.StatusInternalServerError, *rb)
		return
	}

	// Format and return
	data := []map[string]any{}
	for _, record := range records {
		props := map[string]any{
			"length":                  record.Length,
			"is_palindrome":           record.IsPalindrome,
			"unique_characters":       record.UniqueCharacters,
			"word_count":              record.WordCount,
			"sha256_hash":             record.Hash,
			"character_frequency_map": util.CharacterFrequencyMap(record.StringValue),
		}

		data = append(data, map[string]any{
			"id":         record.Hash,
			"value":      record.StringValue,
			"properties": props,
			"created_at": record.CreatedAt,
		})
	}

	respBody := map[string]any{
		"count":           len(data),
		"data":            data,
		"filters_applied": params,
	}

	util.WriteJson(w, http.StatusOK, respBody)
}

func (s *StringAnalyzerHandler) DeleteString(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "string_value")

	err := s.repo.DeleteString(r.Context(), param)

	if err != nil {

		switch {
		case errors.Is(err, errs.ErrNotFound):
			rb := &util.Envelope{"message": "String does not exist in the system"}
			util.WriteJson(w, http.StatusNotFound, *rb)

		default:
			s.logger.Error().Err(err).Msg("error deleting string")
			rb := &util.Envelope{"message": "Something went wrong!"}
			util.WriteJson(w, http.StatusInternalServerError, *rb)
		}
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *StringAnalyzerHandler) FilterByNaturalLanguage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		s.logger.Error().Msg("error deleting string")
		rb := &util.Envelope{"message": "Query string is required"}
		util.WriteJson(w, http.StatusBadRequest, *rb)
	}

	// Parse natural language query into filters
	filters, parsedFilters, err := util.ParseNaturalLanguageQuery(query)
	if err != nil {
		s.logger.Error().Msg("unable to parse natural language")
		rb := &util.Envelope{"message": "Something went wrong"}
		util.WriteJson(w, http.StatusBadRequest, *rb)
		return
	}

	// Validate filters for conflicts
	if err = util.ValidateFilters(filters); err != nil {
		s.logger.Error().Err(err).Msg("conflicted filters")
		rb := &util.Envelope{"message": "Query parsed but resulted in conflicting filters"}
		util.WriteJson(w, http.StatusUnprocessableEntity, *rb)
		return
	}

	results, err := s.repo.GetFilteredStringsByNaturalLanguage(r.Context(), filters)

	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to query natural language")
		rb := &util.Envelope{"message": "Something went wrong"}
		util.WriteJson(w, http.StatusInternalServerError, *rb)
		return
	}

	response := util.Envelope{
		"data":  results,
		"count": len(results),
		"interpreted_query": map[string]any{
			"original":       query,
			"parsed_filters": parsedFilters,
		},
	}

	util.WriteJson(w, http.StatusOK, response)
}
