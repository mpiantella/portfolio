package usecase

import (
	"context"
	"fmt"

	"portfolio/internal/ingestion-pipeline/domain"
)

// ScoreQualityUseCase handles data quality scoring
type ScoreQualityUseCase struct {
	scorer       domain.QualityScorer
	metadataRepo domain.MetadataRepository
	qualityRepo  domain.QualityRepository
	logger       domain.Logger
	weights      domain.QualityWeights
}

// NewScoreQualityUseCase creates a new ScoreQualityUseCase instance
func NewScoreQualityUseCase(
	scorer domain.QualityScorer,
	metadataRepo domain.MetadataRepository,
	qualityRepo domain.QualityRepository,
	logger domain.Logger,
) *ScoreQualityUseCase {
	return &ScoreQualityUseCase{
		scorer:       scorer,
		metadataRepo: metadataRepo,
		qualityRepo:  qualityRepo,
		logger:       logger,
		weights:      domain.DefaultQualityWeights(),
	}
}

// ScoreQualityRequest represents the input for quality scoring
type ScoreQualityRequest struct {
	Entity           *domain.Entity
	ValidationSummary *domain.ValidationSummary
	CustomWeights    *domain.QualityWeights // Optional custom weights
}

// ScoreQualityResponse represents the output of quality scoring
type ScoreQualityResponse struct {
	Score          *domain.QualityScore
	QualityLevel   domain.QualityLevel
	RequiresReview bool
	Success        bool
	Error          string
}

// Execute calculates quality scores for an entity
func (uc *ScoreQualityUseCase) Execute(ctx context.Context, req *ScoreQualityRequest) (*ScoreQualityResponse, error) {
	uc.logger.Info(ctx, "Starting quality scoring", map[string]interface{}{
		"entity_id": req.Entity.EntityID,
	})

	// Get field metadata
	metadata, err := uc.metadataRepo.GetAllFieldMetadata(ctx)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get field metadata", err, nil)
		return &ScoreQualityResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get metadata: %v", err),
		}, err
	}

	// Get validation rules
	rules, err := uc.metadataRepo.GetAllValidationRules(ctx)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get validation rules", err, nil)
		return &ScoreQualityResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get rules: %v", err),
		}, err
	}

	// Determine weights
	weights := uc.weights
	if req.CustomWeights != nil {
		if err := req.CustomWeights.Validate(); err != nil {
			uc.logger.Warn(ctx, "Invalid custom weights, using defaults", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			weights = *req.CustomWeights
		}
	}

	// Prepare quality context
	qualityContext := &domain.QualityContext{
		Metadata:        metadata,
		ValidationRules: rules,
		Weights:         weights,
	}

	// Calculate quality score
	score, err := uc.scorer.CalculateScore(ctx, req.Entity, qualityContext)
	if err != nil {
		uc.logger.Error(ctx, "Quality scoring failed", err, map[string]interface{}{
			"entity_id": req.Entity.EntityID,
		})
		return &ScoreQualityResponse{
			Success: false,
			Error:   fmt.Sprintf("scoring failed: %v", err),
		}, err
	}

	// Validate scores
	if err := score.Validate(); err != nil {
		uc.logger.Error(ctx, "Invalid quality scores", err, nil)
		return &ScoreQualityResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid scores: %v", err),
		}, err
	}

	// Save quality score
	if err := uc.qualityRepo.SaveQualityScore(ctx, score); err != nil {
		uc.logger.Error(ctx, "Failed to save quality score", err, map[string]interface{}{
			"entity_id": req.Entity.EntityID,
		})
		return &ScoreQualityResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to save score: %v", err),
		}, err
	}

	// Determine quality level and review requirement
	qualityLevel := score.QualityLevel()
	requiresReview := score.RequiresReview()

	if requiresReview {
		uc.logger.Warn(ctx, "Entity requires quality review", map[string]interface{}{
			"entity_id":     req.Entity.EntityID,
			"overall_score": score.OverallScore,
		})
	}

	// Get lowest dimension for logging
	lowestDim, lowestScore := score.GetLowestDimension()
	uc.logger.Info(ctx, "Quality scoring completed", map[string]interface{}{
		"entity_id":        req.Entity.EntityID,
		"overall_score":    score.OverallScore,
		"quality_level":    qualityLevel,
		"requires_review":  requiresReview,
		"lowest_dimension": lowestDim,
		"lowest_score":     lowestScore,
	})

	return &ScoreQualityResponse{
		Score:          score,
		QualityLevel:   qualityLevel,
		RequiresReview: requiresReview,
		Success:        true,
	}, nil
}

// ScoreBatch calculates quality scores for multiple entities
func (uc *ScoreQualityUseCase) ScoreBatch(ctx context.Context, entities []*domain.Entity) ([]*ScoreQualityResponse, error) {
	uc.logger.Info(ctx, "Starting batch quality scoring", map[string]interface{}{
		"entity_count": len(entities),
	})

	responses := make([]*ScoreQualityResponse, 0, len(entities))
	successCount := 0
	failCount := 0

	for _, entity := range entities {
		req := &ScoreQualityRequest{
			Entity: entity,
		}

		response, err := uc.Execute(ctx, req)
		if err != nil || !response.Success {
			failCount++
			uc.logger.Error(ctx, "Failed to score entity", err, map[string]interface{}{
				"entity_id": entity.EntityID,
			})
		} else {
			successCount++
		}

		responses = append(responses, response)
	}

	uc.logger.Info(ctx, "Batch quality scoring completed", map[string]interface{}{
		"total_count":   len(entities),
		"success_count": successCount,
		"fail_count":    failCount,
	})

	return responses, nil
}

// GetQualityReport generates a quality report for an entity
func (uc *ScoreQualityUseCase) GetQualityReport(ctx context.Context, entityID string) (map[string]interface{}, error) {
	// Get quality score
	score, err := uc.qualityRepo.GetQualityScore(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quality score: %w", err)
	}

	// Get validation results
	validationResults, err := uc.qualityRepo.GetValidationResults(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get validation results: %w", err)
	}

	// Build report
	report := map[string]interface{}{
		"entity_id":        entityID,
		"overall_score":    score.OverallScore,
		"quality_level":    score.QualityLevel(),
		"requires_review":  score.RequiresReview(),
		"calculated_at":    score.CalculatedAt,
		"dimensions": map[string]float64{
			"completeness": score.CompletenessScore,
			"accuracy":     score.AccuracyScore,
			"consistency":  score.ConsistencyScore,
			"timeliness":   score.TimelinessScore,
			"uniqueness":   score.UniquenessScore,
			"validity":     score.ValidityScore,
		},
		"validation_results": validationResults,
	}

	// Identify weakest dimension
	lowestDim, lowestScore := score.GetLowestDimension()
	report["weakest_dimension"] = map[string]interface{}{
		"name":  lowestDim,
		"score": lowestScore,
	}

	return report, nil
}
