package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ExperienceData represents a single experience data record
type ExperienceData struct {
	ID             uuid.UUID       `json:"id"`
	CollectedAt    time.Time       `json:"collected_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	SourceType     string          `json:"source_type"`
	SourceID       *string         `json:"source_id,omitempty"`
	SourceName     *string         `json:"source_name,omitempty"`
	FieldID        string          `json:"field_id"`
	FieldLabel     *string         `json:"field_label,omitempty"`
	FieldType      string          `json:"field_type"`
	ValueText      *string         `json:"value_text,omitempty"`
	ValueNumber    *float64        `json:"value_number,omitempty"`
	ValueBoolean   *bool           `json:"value_boolean,omitempty"`
	ValueDate      *time.Time      `json:"value_date,omitempty"`
	ValueJSON      json.RawMessage `json:"value_json,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	Language       *string         `json:"language,omitempty"`
	UserIdentifier *string         `json:"user_identifier,omitempty"`
}

// CreateExperienceRequest represents the request to create experience data
type CreateExperienceRequest struct {
	CollectedAt    *time.Time      `json:"collected_at,omitempty"`
	SourceType     string          `json:"source_type"`
	SourceID       *string         `json:"source_id,omitempty"`
	SourceName     *string         `json:"source_name,omitempty"`
	FieldID        string          `json:"field_id"`
	FieldLabel     *string         `json:"field_label,omitempty"`
	FieldType      string          `json:"field_type"`
	ValueText      *string         `json:"value_text,omitempty"`
	ValueNumber    *float64        `json:"value_number,omitempty"`
	ValueBoolean   *bool           `json:"value_boolean,omitempty"`
	ValueDate      *time.Time      `json:"value_date,omitempty"`
	ValueJSON      json.RawMessage `json:"value_json,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	Language       *string         `json:"language,omitempty"`
	UserIdentifier *string         `json:"user_identifier,omitempty"`
}

// UpdateExperienceRequest represents the request to update experience data
type UpdateExperienceRequest struct {
	SourceType     *string         `json:"source_type,omitempty"`
	SourceID       *string         `json:"source_id,omitempty"`
	SourceName     *string         `json:"source_name,omitempty"`
	FieldID        *string         `json:"field_id,omitempty"`
	FieldLabel     *string         `json:"field_label,omitempty"`
	FieldType      *string         `json:"field_type,omitempty"`
	ValueText      *string         `json:"value_text,omitempty"`
	ValueNumber    *float64        `json:"value_number,omitempty"`
	ValueBoolean   *bool           `json:"value_boolean,omitempty"`
	ValueDate      *time.Time      `json:"value_date,omitempty"`
	ValueJSON      json.RawMessage `json:"value_json,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	Language       *string         `json:"language,omitempty"`
	UserIdentifier *string         `json:"user_identifier,omitempty"`
}

// ListExperiencesFilters represents filters for listing experiences
type ListExperiencesFilters struct {
	SourceType     *string
	SourceID       *string
	FieldID        *string
	UserIdentifier *string
	Limit          int
	Offset         int
}

// SearchExperiencesRequest represents search parameters for experiences
type SearchExperiencesRequest struct {
	Query          *string    `json:"query,omitempty"`           // Full-text search query
	SourceType     *string    `json:"source_type,omitempty"`     // Filter by source type
	SourceID       *string    `json:"source_id,omitempty"`       // Filter by source ID
	FieldID        *string    `json:"field_id,omitempty"`        // Filter by field ID
	FieldType      *string    `json:"field_type,omitempty"`      // Filter by field type
	UserIdentifier *string    `json:"user_identifier,omitempty"` // Filter by user identifier
	StartDate      *time.Time `json:"start_date,omitempty"`      // Filter by collected_at >= start_date
	EndDate        *time.Time `json:"end_date,omitempty"`        // Filter by collected_at <= end_date
	Limit          int        `json:"limit,omitempty"`           // Maximum number of results
	Offset         int        `json:"offset,omitempty"`          // Number of results to skip
}
