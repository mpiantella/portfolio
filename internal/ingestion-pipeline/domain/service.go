package domain

import (
	"context"
	"io"
)

// FileParser defines the interface for parsing different file formats
// This interface is defined in the domain but implemented in infrastructure
type FileParser interface {
	// Parse parses a file and returns structured data
	Parse(ctx context.Context, file io.Reader, metadata *FileMetadata) (*ParsedData, error)

	// Validate validates file structure without parsing all data
	Validate(ctx context.Context, file io.Reader) error

	// GetMetadata extracts metadata from a file
	GetMetadata(ctx context.Context, file io.Reader) (*FileMetadata, error)

	// SupportsFormat returns true if this parser supports the given format
	SupportsFormat(format FileFormat) bool
}

// DataNormalizer defines the interface for data normalization
type DataNormalizer interface {
	// Normalize transforms raw parsed data into normalized entities
	Normalize(ctx context.Context, data *ParsedData, drivingFieldName string) ([]*Entity, error)

	// ValidateNormalization checks if normalized data is valid
	ValidateNormalization(ctx context.Context, entities []*Entity) error

	// DetectDrivingField attempts to detect the driving field from data
	DetectDrivingField(ctx context.Context, data *ParsedData) (string, error)
}

// DataValidator defines the interface for data validation
type DataValidator interface {
	// ValidateEntity validates a single entity against all rules
	ValidateEntity(ctx context.Context, entity *Entity, metadata []*FieldMetadata, rules []*ValidationRule) (*ValidationSummary, error)

	// ValidateBatch validates multiple entities
	ValidateBatch(ctx context.Context, entities []*Entity) ([]*ValidationSummary, error)

	// ValidateField validates a single field value
	ValidateField(ctx context.Context, fieldName string, value interface{}, metadata *FieldMetadata, rules []*ValidationRule) ([]*ValidationResult, error)
}

// QualityScorer defines the interface for data quality scoring
type QualityScorer interface {
	// CalculateScore calculates quality scores for an entity
	CalculateScore(ctx context.Context, entity *Entity, context *QualityContext) (*QualityScore, error)

	// CalculateBatchScores calculates quality scores for multiple entities
	CalculateBatchScores(ctx context.Context, entities []*Entity, context *QualityContext) ([]*QualityScore, error)

	// CalculateCompleteness calculates completeness score
	CalculateCompleteness(ctx context.Context, entity *Entity, metadata []*FieldMetadata) float64

	// CalculateAccuracy calculates accuracy score
	CalculateAccuracy(ctx context.Context, entity *Entity, validationSummary *ValidationSummary) float64
}

// QualityContext provides context for quality scoring
type QualityContext struct {
	Metadata         []*FieldMetadata
	ValidationRules  []*ValidationRule
	RelatedEntities  []*Entity
	ExistingEntities []*Entity
	Weights          QualityWeights
}

// Storage defines the interface for file storage operations
type Storage interface {
	// Upload uploads a file to storage
	Upload(ctx context.Context, path string, data io.Reader) error

	// Download downloads a file from storage
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete deletes a file from storage
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, path string) (bool, error)

	// GetURL returns a pre-signed URL for a file
	GetURL(ctx context.Context, path string, expirationMinutes int) (string, error)
}

// Logger defines the interface for logging
type Logger interface {
	// Debug logs a debug message
	Debug(ctx context.Context, message string, fields map[string]interface{})

	// Info logs an info message
	Info(ctx context.Context, message string, fields map[string]interface{})

	// Warn logs a warning message
	Warn(ctx context.Context, message string, fields map[string]interface{})

	// Error logs an error message
	Error(ctx context.Context, message string, err error, fields map[string]interface{})
}

// Notifier defines the interface for sending notifications
type Notifier interface {
	// Notify sends a notification
	Notify(ctx context.Context, notification *Notification) error
}

// Notification represents a notification to be sent
type Notification struct {
	Type      NotificationType       `json:"type"`
	Subject   string                 `json:"subject"`
	Message   string                 `json:"message"`
	Recipients []string               `json:"recipients"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypeSlack NotificationType = "slack"
)
