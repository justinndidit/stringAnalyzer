package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/justinndidit/stringAnalyzer/internal/database"
	"github.com/justinndidit/stringAnalyzer/internal/dto"
	"github.com/justinndidit/stringAnalyzer/internal/errs"

	"github.com/justinndidit/stringAnalyzer/internal/model"
	"github.com/rs/zerolog"
)

type StringRepository struct {
	logger *zerolog.Logger
	db     *database.Database
}

func NewStringRepository(logger *zerolog.Logger, db *database.Database) *StringRepository {
	return &StringRepository{
		logger: logger,
		db:     db,
	}
}

func (r *StringRepository) GetFilteredStrings(ctx context.Context, params dto.QueryParams) ([]model.String, error) {
	stmt := `
		SELECT
			*
		FROM
			strings
		WHERE
			(@is_palindrome::boolean IS NULL OR is_palindrome = @is_palindrome::boolean)
			AND (@min_length::int IS NULL OR length >= @min_length::int)
			AND (@max_length::int IS NULL OR length <= @max_length::int)
			AND (@word_count::int IS NULL OR word_count = @word_count::int)
			AND (@contains_character::text IS NULL OR string_value ILIKE '%' || @contains_character || '%');
	`

	rows, err := r.db.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"is_palindrome": func() any {
			if params.IsPalindrome == nil {
				return nil
			}
			return *params.IsPalindrome
		}(),
		"min_length": func() any {
			if params.MinLength == nil {
				return nil
			}
			return *params.MinLength
		}(),
		"max_length": func() any {
			if params.MaxLength == nil {
				return nil
			}
			return *params.MaxLength
		}(),
		"word_count": func() any {
			if params.WordCount == nil {
				return nil
			}
			return *params.WordCount
		}(),
		"contains_character": params.ContainsCharacter,
	})

	if err != nil {
		r.logger.Error().Err(err).Msg("Query Failed!")
		return nil, fmt.Errorf("failed to execute string query: %w", err)
	}

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.String])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:strings: %w", err)
	}

	return records, nil
}

func (r *StringRepository) GetStringByValue(ctx context.Context, value string) (*model.String, error) {
	stmt := `
		SELECT
			*
		FROM
			strings
		WHERE
			string_value = @value
	`

	args := pgx.NamedArgs{
		"value": value,
	}

	rows, err := r.db.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.String])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info().Msg("No rows matching query")
			return nil, nil // no match found, return nil without error
		}

		r.logger.Error().Err(err).Msg("No")
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &record, nil
}

func (r *StringRepository) CreateString(ctx context.Context, payload *dto.CreateString) (*model.String, error) {

	stmt := `
		INSERT INTO strings (
			string_value,
			is_palindrome,
			unique_characters,
			word_count,
			sha256_hash,
			length
		)
		VALUES (
			@string_value,
			@is_palindrome,
			@unique_characters,
			@word_count,
			@sha256_hash,
			@length
		)
		RETURNING *
	`

	rows, err := r.db.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"string_value":      payload.StringValue,
		"is_palindrome":     payload.IsPalindrome,
		"unique_characters": payload.UniqueCharacters,
		"word_count":        payload.WordCount,
		"sha256_hash":       payload.Hash,
		"length":            payload.Length,
	})
	if err != nil {
		r.logger.Error().Err(err).Msg("Query Failed!")
		return nil, fmt.Errorf("failed to execute create string query: %w", err)
	}

	newString, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.String])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:strings: %w", err)
	}

	return &newString, nil
}

func (r *StringRepository) DeleteString(ctx context.Context, value string) error {
	stmt := `DELETE FROM strings WHERE string_value = @string_value`

	cmdTag, err := r.db.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"string_value": value,
	})

	if err != nil {
		r.logger.Error().Err(err).Msg("Delete query failed!")
		return fmt.Errorf("failed to execute delete string query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *StringRepository) GetFilteredStringsByNaturalLanguage(ctx context.Context, params *dto.FilterParams) ([]model.String, error) {
	stmt := `
		SELECT
			*
		FROM
			strings
		WHERE
			(@is_palindrome::boolean IS NULL OR is_palindrome = @is_palindrome::boolean)
			AND (@min_length::int IS NULL OR length >= @min_length::int)
			AND (@max_length::int IS NULL OR length <= @max_length::int)
			AND (@word_count::int IS NULL OR word_count = @word_count::int)
			AND (@contains_character::text IS NULL OR string_value ILIKE '%' || @contains_character || '%')
		ORDER BY
			created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"is_palindrome": func() any {
			if params.IsPalindrome == nil {
				return nil
			}
			return *params.IsPalindrome
		}(),
		"min_length": func() any {
			if params.MinLength == nil {
				return nil
			}
			return *params.MinLength
		}(),
		"max_length": func() any {
			if params.MaxLength == nil {
				return nil
			}
			return *params.MaxLength
		}(),
		"word_count": func() any {
			if params.WordCount == nil {
				return nil
			}
			return *params.WordCount
		}(),
		"contains_character": func() any {
			if params.ContainsCharacter == nil {
				return nil
			}
			return *params.ContainsCharacter
		}(),
	})

	if err != nil {
		r.logger.Error().
			Err(err).
			Interface("params", params).
			Msg("Natural language filter query failed")
		return nil, fmt.Errorf("failed to execute natural language filter query: %w", err)
	}

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.String])
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to collect rows from natural language query")
		return nil, fmt.Errorf("failed to collect rows from table:strings: %w", err)
	}

	r.logger.Info().
		Int("count", len(records)).
		Interface("filters", params).
		Msg("Natural language query executed successfully")

	return records, nil
}
