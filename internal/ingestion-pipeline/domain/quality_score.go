package domain

import "time"

// QualityScore represents data quality metrics for an entity
type QualityScore struct {
	ScoreID          string    `json:"score_id"`
	EntityID         string    `json:"entity_id"`
	CompletenessScore float64   `json:"completeness_score"`
	AccuracyScore    float64   `json:"accuracy_score"`
	ConsistencyScore float64   `json:"consistency_score"`
	TimelinessScore  float64   `json:"timeliness_score"`
	UniquenessScore  float64   `json:"uniqueness_score"`
	ValidityScore    float64   `json:"validity_score"`
	OverallScore     float64   `json:"overall_score"`
	CalculatedAt     time.Time `json:"calculated_at"`
}

// QualityWeights defines the weight of each quality dimension
type QualityWeights struct {
	Completeness float64 `json:"completeness"`
	Accuracy     float64 `json:"accuracy"`
	Consistency  float64 `json:"consistency"`
	Timeliness   float64 `json:"timeliness"`
	Uniqueness   float64 `json:"uniqueness"`
	Validity     float64 `json:"validity"`
}

// DefaultQualityWeights returns the default weights for quality scoring
func DefaultQualityWeights() QualityWeights {
	return QualityWeights{
		Completeness: 20.0,
		Accuracy:     25.0,
		Consistency:  20.0,
		Timeliness:   10.0,
		Uniqueness:   15.0,
		Validity:     10.0,
	}
}

// Validate checks if the weights sum to 100
func (qw *QualityWeights) Validate() error {
	sum := qw.Completeness + qw.Accuracy + qw.Consistency +
		qw.Timeliness + qw.Uniqueness + qw.Validity

	if sum != 100.0 {
		return NewValidationError("quality weights must sum to 100")
	}

	return nil
}

// CalculateOverall computes the overall score using the weights
func (qs *QualityScore) CalculateOverall(weights QualityWeights) float64 {
	qs.OverallScore = (weights.Completeness*qs.CompletenessScore +
		weights.Accuracy*qs.AccuracyScore +
		weights.Consistency*qs.ConsistencyScore +
		weights.Timeliness*qs.TimelinessScore +
		weights.Uniqueness*qs.UniquenessScore +
		weights.Validity*qs.ValidityScore) / 100.0

	return qs.OverallScore
}

// QualityLevel returns the quality level based on the overall score
func (qs *QualityScore) QualityLevel() QualityLevel {
	switch {
	case qs.OverallScore >= 90:
		return QualityLevelExcellent
	case qs.OverallScore >= 70:
		return QualityLevelGood
	case qs.OverallScore >= 50:
		return QualityLevelPoor
	default:
		return QualityLevelCritical
	}
}

// QualityLevel represents the quality classification
type QualityLevel string

const (
	QualityLevelExcellent QualityLevel = "excellent"
	QualityLevelGood      QualityLevel = "good"
	QualityLevelPoor      QualityLevel = "poor"
	QualityLevelCritical  QualityLevel = "critical"
)

// RequiresReview returns true if the quality score requires manual review
func (qs *QualityScore) RequiresReview() bool {
	return qs.OverallScore < 70.0
}

// IsAcceptable returns true if the quality score is acceptable for production
func (qs *QualityScore) IsAcceptable() bool {
	return qs.OverallScore >= 70.0
}

// GetLowestDimension returns the quality dimension with the lowest score
func (qs *QualityScore) GetLowestDimension() (string, float64) {
	dimensions := map[string]float64{
		"completeness": qs.CompletenessScore,
		"accuracy":     qs.AccuracyScore,
		"consistency":  qs.ConsistencyScore,
		"timeliness":   qs.TimelinessScore,
		"uniqueness":   qs.UniquenessScore,
		"validity":     qs.ValidityScore,
	}

	minDimension := "completeness"
	minScore := qs.CompletenessScore

	for dimension, score := range dimensions {
		if score < minScore {
			minDimension = dimension
			minScore = score
		}
	}

	return minDimension, minScore
}

// Validate checks if all scores are within valid range
func (qs *QualityScore) Validate() error {
	scores := []struct {
		name  string
		value float64
	}{
		{"completeness", qs.CompletenessScore},
		{"accuracy", qs.AccuracyScore},
		{"consistency", qs.ConsistencyScore},
		{"timeliness", qs.TimelinessScore},
		{"uniqueness", qs.UniquenessScore},
		{"validity", qs.ValidityScore},
		{"overall", qs.OverallScore},
	}

	for _, score := range scores {
		if score.value < 0 || score.value > 100 {
			return NewValidationError("quality score " + score.name + " must be between 0 and 100")
		}
	}

	return nil
}
