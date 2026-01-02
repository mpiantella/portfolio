package domain

import "time"

// ValidationRule represents a business rule for data validation
type ValidationRule struct {
	RuleID         string            `json:"rule_id"`
	RuleName       string            `json:"rule_name"`
	RuleType       ValidationRuleType `json:"rule_type"`
	FieldID        string            `json:"field_id,omitempty"`
	RuleExpression string            `json:"rule_expression"`
	ErrorMessage   string            `json:"error_message"`
	Severity       ValidationSeverity `json:"severity"`
	IsActive       bool              `json:"is_active"`
	CreatedAt      time.Time         `json:"created_at"`
}

// ValidationRuleType represents the type of validation rule
type ValidationRuleType string

const (
	RuleTypeFormat     ValidationRuleType = "format"
	RuleTypeRange      ValidationRuleType = "range"
	RuleTypeReference  ValidationRuleType = "reference"
	RuleTypeCustom     ValidationRuleType = "custom"
	RuleTypeConsistency ValidationRuleType = "consistency"
)

// ValidationSeverity represents the severity level of a validation rule
type ValidationSeverity string

const (
	SeverityError   ValidationSeverity = "error"
	SeverityWarning ValidationSeverity = "warning"
	SeverityInfo    ValidationSeverity = "info"
)

// Validate checks if the validation rule is valid
func (vr *ValidationRule) Validate() error {
	if vr.RuleName == "" {
		return NewValidationError("rule name cannot be empty")
	}
	if vr.RuleExpression == "" {
		return NewValidationError("rule expression cannot be empty")
	}
	if vr.ErrorMessage == "" {
		return NewValidationError("error message cannot be empty")
	}
	if !vr.RuleType.IsValid() {
		return NewValidationError("invalid rule type")
	}
	if !vr.Severity.IsValid() {
		return NewValidationError("invalid severity")
	}
	return nil
}

// IsValid checks if the rule type is valid
func (rt ValidationRuleType) IsValid() bool {
	switch rt {
	case RuleTypeFormat, RuleTypeRange, RuleTypeReference, RuleTypeCustom, RuleTypeConsistency:
		return true
	default:
		return false
	}
}

// IsValid checks if the severity is valid
func (vs ValidationSeverity) IsValid() bool {
	switch vs {
	case SeverityError, SeverityWarning, SeverityInfo:
		return true
	default:
		return false
	}
}

// IsBlockingError returns true if the validation failure should block processing
func (vr *ValidationRule) IsBlockingError() bool {
	return vr.Severity == SeverityError
}

// ValidationResult represents the result of applying a validation rule
type ValidationResult struct {
	RuleID        string             `json:"rule_id"`
	RuleName      string             `json:"rule_name"`
	Passed        bool               `json:"passed"`
	ActualValue   interface{}        `json:"actual_value,omitempty"`
	ExpectedValue interface{}        `json:"expected_value,omitempty"`
	ErrorDetails  string             `json:"error_details,omitempty"`
	Severity      ValidationSeverity `json:"severity"`
	CheckedAt     time.Time          `json:"checked_at"`
}

// NewValidationResult creates a new validation result
func NewValidationResult(ruleID, ruleName string, passed bool, severity ValidationSeverity) *ValidationResult {
	return &ValidationResult{
		RuleID:    ruleID,
		RuleName:  ruleName,
		Passed:    passed,
		Severity:  severity,
		CheckedAt: time.Now(),
	}
}

// WithError adds error details to the validation result
func (vr *ValidationResult) WithError(message string) *ValidationResult {
	vr.ErrorDetails = message
	return vr
}

// WithValues adds actual and expected values to the validation result
func (vr *ValidationResult) WithValues(actual, expected interface{}) *ValidationResult {
	vr.ActualValue = actual
	vr.ExpectedValue = expected
	return vr
}

// ValidationSummary aggregates multiple validation results
type ValidationSummary struct {
	EntityID       string              `json:"entity_id"`
	TotalRules     int                 `json:"total_rules"`
	PassedRules    int                 `json:"passed_rules"`
	FailedRules    int                 `json:"failed_rules"`
	WarningCount   int                 `json:"warning_count"`
	ErrorCount     int                 `json:"error_count"`
	ValidationRate float64             `json:"validation_rate"` // % passed
	Results        []*ValidationResult `json:"results"`
	ValidatedAt    time.Time           `json:"validated_at"`
}

// NewValidationSummary creates a new validation summary from results
func NewValidationSummary(entityID string, results []*ValidationResult) *ValidationSummary {
	summary := &ValidationSummary{
		EntityID:    entityID,
		TotalRules:  len(results),
		Results:     results,
		ValidatedAt: time.Now(),
	}

	for _, result := range results {
		if result.Passed {
			summary.PassedRules++
		} else {
			summary.FailedRules++
			if result.Severity == SeverityError {
				summary.ErrorCount++
			} else if result.Severity == SeverityWarning {
				summary.WarningCount++
			}
		}
	}

	if summary.TotalRules > 0 {
		summary.ValidationRate = (float64(summary.PassedRules) / float64(summary.TotalRules)) * 100.0
	}

	return summary
}

// IsValid returns true if there are no blocking errors
func (vs *ValidationSummary) IsValid() bool {
	return vs.ErrorCount == 0
}

// HasWarnings returns true if there are any warnings
func (vs *ValidationSummary) HasWarnings() bool {
	return vs.WarningCount > 0
}
