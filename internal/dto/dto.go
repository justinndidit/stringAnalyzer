package dto

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/justinndidit/stringAnalyzer/internal/errs"
)

type UploadString struct {
	Value string `json:"value"`
}

func (u *UploadString) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	val, ok := raw["value"]
	if !ok {
		return fmt.Errorf("missing 'value' field")
	}

	strVal, ok := val.(string)
	if !ok {
		return &errs.InvalidTypeError{
			Field:    "value",
			Expected: "string",
			Got:      fmt.Sprintf("%T", val),
		}
	}

	u.Value = strVal
	return nil
}

type CreateString struct {
	StringValue      string
	IsPalindrome     bool
	UniqueCharacters int
	WordCount        int
	Hash             string
	Length           int
}

type QueryParams struct {
	IsPalindrome      *bool  `validate:"omitempty"` // optional, pointer differentiates false vs not provided
	MinLength         *int   `validate:"omitempty,gte=0"`
	MaxLength         *int   `validate:"omitempty,gte=0"`
	WordCount         *int   `validate:"omitempty,gte=0"`
	ContainsCharacter string `validate:"omitempty,len=1"` // optional, must be 1 char if provided
}

func (q *QueryParams) Validate() error {
	validate := validator.New()
	return validate.Struct(q)
}

// NLP
type FilterParams struct {
	IsPalindrome      *bool
	MinLength         *int
	MaxLength         *int
	WordCount         *int
	ContainsCharacter *string
}

type InterpretedQuery struct {
	Original      string                 `json:"original"`
	ParsedFilters map[string]interface{} `json:"parsed_filters"`
}

type NLQueryResponse struct {
	Data             []interface{}    `json:"data"`
	Count            int              `json:"count"`
	InterpretedQuery InterpretedQuery `json:"interpreted_query"`
}
