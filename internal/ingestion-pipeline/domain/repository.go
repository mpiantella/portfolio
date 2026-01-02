package domain

import "context"

// EntityRepository defines the interface for entity persistence
// This interface is defined in the domain but implemented in infrastructure
type EntityRepository interface {
	// Save persists an entity
	Save(ctx context.Context, entity *Entity) error

	// FindByID retrieves an entity by ID
	FindByID(ctx context.Context, entityID string) (*Entity, error)

	// FindByDrivingField retrieves an entity by driving field value
	FindByDrivingField(ctx context.Context, drivingFieldValue string) (*Entity, error)

	// FindByType retrieves entities by type with pagination
	FindByType(ctx context.Context, entityType string, limit, offset int) ([]*Entity, error)

	// Update updates an existing entity
	Update(ctx context.Context, entity *Entity) error

	// Delete soft deletes an entity
	Delete(ctx context.Context, entityID string) error

	// Exists checks if an entity exists
	Exists(ctx context.Context, drivingFieldValue string) (bool, error)
}

// MetadataRepository defines the interface for metadata persistence
type MetadataRepository interface {
	// GetFieldMetadata retrieves field metadata by field name
	GetFieldMetadata(ctx context.Context, fieldName string) (*FieldMetadata, error)

	// GetAllFieldMetadata retrieves all active field metadata
	GetAllFieldMetadata(ctx context.Context) ([]*FieldMetadata, error)

	// GetDrivingField retrieves the driving field metadata
	GetDrivingField(ctx context.Context) (*FieldMetadata, error)

	// SaveFieldMetadata persists field metadata
	SaveFieldMetadata(ctx context.Context, metadata *FieldMetadata) error

	// GetValidationRules retrieves validation rules for a field
	GetValidationRules(ctx context.Context, fieldID string) ([]*ValidationRule, error)

	// GetAllValidationRules retrieves all active validation rules
	GetAllValidationRules(ctx context.Context) ([]*ValidationRule, error)

	// SaveValidationRule persists a validation rule
	SaveValidationRule(ctx context.Context, rule *ValidationRule) error
}

// QualityRepository defines the interface for quality score persistence
type QualityRepository interface {
	// SaveQualityScore persists a quality score
	SaveQualityScore(ctx context.Context, score *QualityScore) error

	// GetQualityScore retrieves a quality score for an entity
	GetQualityScore(ctx context.Context, entityID string) (*QualityScore, error)

	// SaveValidationResults persists validation results
	SaveValidationResults(ctx context.Context, entityID string, results []*ValidationResult) error

	// GetValidationResults retrieves validation results for an entity
	GetValidationResults(ctx context.Context, entityID string) ([]*ValidationResult, error)

	// GetQualityMetrics retrieves aggregated quality metrics
	GetQualityMetrics(ctx context.Context, aggregationLevel, aggregationKey string) (map[string]float64, error)
}

// FileRepository defines the interface for file metadata persistence
type FileRepository interface {
	// SaveFileMetadata persists file metadata
	SaveFileMetadata(ctx context.Context, metadata *FileMetadata) error

	// GetFileMetadata retrieves file metadata by ID
	GetFileMetadata(ctx context.Context, fileID string) (*FileMetadata, error)

	// UpdateFileStatus updates the processing status of a file
	UpdateFileStatus(ctx context.Context, fileID string, status ProcessingStatus, errorMessage string) error

	// GetPendingFiles retrieves files pending processing
	GetPendingFiles(ctx context.Context, limit int) ([]*FileMetadata, error)
}

// LineageRepository defines the interface for data lineage tracking
type LineageRepository interface {
	// RecordLineage records a data lineage entry
	RecordLineage(ctx context.Context, lineage *LineageRecord) error

	// GetLineage retrieves lineage for an entity
	GetLineage(ctx context.Context, entityID string) ([]*LineageRecord, error)

	// GetSourceLineage retrieves lineage for a source file
	GetSourceLineage(ctx context.Context, sourceFileID string) ([]*LineageRecord, error)
}

// LineageRecord represents a data lineage entry
type LineageRecord struct {
	LineageID              string                 `json:"lineage_id"`
	EntityID               string                 `json:"entity_id"`
	SourceFileID           string                 `json:"source_file_id"`
	TransformationStep     string                 `json:"transformation_step"`
	TransformationTimestamp string                 `json:"transformation_timestamp"`
	TransformationDetails  map[string]interface{} `json:"transformation_details"`
	PerformedBy            string                 `json:"performed_by"`
}
