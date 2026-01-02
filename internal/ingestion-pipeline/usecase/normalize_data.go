package usecase

import (
	"context"
	"fmt"

	"portfolio/internal/ingestion-pipeline/domain"
)

// NormalizeDataUseCase handles data normalization logic
type NormalizeDataUseCase struct {
	normalizer   domain.DataNormalizer
	metadataRepo domain.MetadataRepository
	entityRepo   domain.EntityRepository
	logger       domain.Logger
}

// NewNormalizeDataUseCase creates a new NormalizeDataUseCase instance
func NewNormalizeDataUseCase(
	normalizer domain.DataNormalizer,
	metadataRepo domain.MetadataRepository,
	entityRepo domain.EntityRepository,
	logger domain.Logger,
) *NormalizeDataUseCase {
	return &NormalizeDataUseCase{
		normalizer:   normalizer,
		metadataRepo: metadataRepo,
		entityRepo:   entityRepo,
		logger:       logger,
	}
}

// NormalizeDataRequest represents the input for data normalization
type NormalizeDataRequest struct {
	ParsedData       *domain.ParsedData
	DrivingFieldName string // Optional - will be detected if not provided
}

// NormalizeDataResponse represents the output of data normalization
type NormalizeDataResponse struct {
	Entities         []*domain.Entity
	DrivingFieldName string
	NormalizedCount  int
	Success          bool
	Error            string
}

// Execute normalizes parsed data into domain entities
func (uc *NormalizeDataUseCase) Execute(ctx context.Context, req *NormalizeDataRequest) (*NormalizeDataResponse, error) {
	uc.logger.Info(ctx, "Starting data normalization", map[string]interface{}{
		"record_count":       req.ParsedData.RecordCount,
		"driving_field_name": req.DrivingFieldName,
	})

	// Determine driving field
	drivingFieldName := req.DrivingFieldName
	if drivingFieldName == "" {
		// Try to get from metadata
		drivingField, err := uc.metadataRepo.GetDrivingField(ctx)
		if err != nil {
			// Attempt auto-detection
			detected, detectErr := uc.normalizer.DetectDrivingField(ctx, req.ParsedData)
			if detectErr != nil {
				uc.logger.Error(ctx, "Failed to detect driving field", detectErr, nil)
				return &NormalizeDataResponse{
					Success: false,
					Error:   fmt.Sprintf("no driving field specified and detection failed: %v", detectErr),
				}, fmt.Errorf("failed to determine driving field: %w", detectErr)
			}
			drivingFieldName = detected
			uc.logger.Info(ctx, "Auto-detected driving field", map[string]interface{}{
				"field_name": drivingFieldName,
			})
		} else {
			drivingFieldName = drivingField.FieldName
		}
	}

	// Normalize data
	entities, err := uc.normalizer.Normalize(ctx, req.ParsedData, drivingFieldName)
	if err != nil {
		uc.logger.Error(ctx, "Data normalization failed", err, map[string]interface{}{
			"driving_field": drivingFieldName,
		})
		return &NormalizeDataResponse{
			Success: false,
			Error:   fmt.Sprintf("normalization failed: %v", err),
		}, err
	}

	// Validate normalized data
	if err := uc.normalizer.ValidateNormalization(ctx, entities); err != nil {
		uc.logger.Error(ctx, "Normalization validation failed", err, nil)
		return &NormalizeDataResponse{
			Success: false,
			Error:   fmt.Sprintf("normalization validation failed: %v", err),
		}, err
	}

	uc.logger.Info(ctx, "Data normalization completed", map[string]interface{}{
		"entity_count":  len(entities),
		"driving_field": drivingFieldName,
	})

	return &NormalizeDataResponse{
		Entities:         entities,
		DrivingFieldName: drivingFieldName,
		NormalizedCount:  len(entities),
		Success:          true,
	}, nil
}

// NormalizeAndMerge normalizes data and merges with existing entities
func (uc *NormalizeDataUseCase) NormalizeAndMerge(ctx context.Context, req *NormalizeDataRequest) (*NormalizeDataResponse, error) {
	// First normalize the data
	response, err := uc.Execute(ctx, req)
	if err != nil || !response.Success {
		return response, err
	}

	// Check for existing entities and merge
	mergedCount := 0
	newCount := 0

	for _, entity := range response.Entities {
		// Check if entity already exists
		exists, err := uc.entityRepo.Exists(ctx, entity.DrivingFieldValue)
		if err != nil {
			uc.logger.Error(ctx, "Failed to check entity existence", err, map[string]interface{}{
				"driving_field_value": entity.DrivingFieldValue,
			})
			continue
		}

		if exists {
			// Entity exists - merge attributes
			existingEntity, err := uc.entityRepo.FindByDrivingField(ctx, entity.DrivingFieldValue)
			if err != nil {
				uc.logger.Error(ctx, "Failed to retrieve existing entity", err, map[string]interface{}{
					"driving_field_value": entity.DrivingFieldValue,
				})
				continue
			}

			// Merge attributes (new values overwrite existing)
			for key, value := range entity.Attributes {
				existingEntity.SetAttribute(key, value)
			}

			if err := uc.entityRepo.Update(ctx, existingEntity); err != nil {
				uc.logger.Error(ctx, "Failed to update entity", err, map[string]interface{}{
					"entity_id": existingEntity.EntityID,
				})
				continue
			}

			mergedCount++
		} else {
			// New entity - save
			if err := uc.entityRepo.Save(ctx, entity); err != nil {
				uc.logger.Error(ctx, "Failed to save entity", err, map[string]interface{}{
					"driving_field_value": entity.DrivingFieldValue,
				})
				continue
			}

			newCount++
		}
	}

	uc.logger.Info(ctx, "Normalization and merge completed", map[string]interface{}{
		"merged_count": mergedCount,
		"new_count":    newCount,
	})

	return response, nil
}
