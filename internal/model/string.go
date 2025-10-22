package model

import (
	"time"
)

type String struct {
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	StringValue      string    `json:"value" db:"string_value"`
	IsPalindrome     bool      `json:"is_palindrome" db:"is_palindrome"`
	UniqueCharacters int       `json:"unique_characters" db:"unique_characters"`
	WordCount        int       `json:"word_count" db:"word_count"`
	Hash             string    `json:"sha256_hash" db:"sha256_hash"`
	Length           int       `json:"length" db:"length"`
}
