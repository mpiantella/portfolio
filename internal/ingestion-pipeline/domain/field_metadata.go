package domain

import (
	"fmt"
	"time"
)

// FieldMetadata defines the structure and rules for a field
type FieldMetadata struct {
	FieldID         string                 `json:"field_id"`
	FieldName       string                 `json:"field_name"`
	DisplayName     string                 `json:"display_name"`
	FieldType       FieldType              `json:"field_type"`
	Description     string                 `json:"description"`
	IsRequired      bool                   `json:"is_required"`
	IsDrivingField  bool                   `json:"is_driving_field"`
	DefaultValue    *string                `json:"default_value,omitempty"`
	FormatPattern   string                 `json:"format_pattern,omitempty"` // regex
	MinValue        *float64               `json:"min_value,omitempty"`
	MaxValue        *float64               `json:"max_value,omitempty"`
	MaxLength       *int                   `json:"max_length,omitempty"`
	AllowedValues   []string               `json:"allowed_values,omitempty"` // enum
	BusinessRules   map[string]interface{} `json:"business_rules,omitempty"`
	QualityWeight   float64                `json:"quality_weight"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	CreatedBy       string                 `json:"created_by"`
	IsActive        bool                   `json:"is_active"`
}

// FieldType represents the data type of a field
type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInteger FieldType = "integer"
	FieldTypeDecimal FieldType = "decimal"
	FieldTypeDate    FieldType = "date"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeJSON    FieldType = "json"
)

// Validate checks if the field metadata is valid
func (fm *FieldMetadata) Validate() error {
	if fm.FieldName == "" {
		return NewValidationError("field name cannot be empty")
	}
	if fm.DisplayName == "" {
		return NewValidationError("display name cannot be empty")
	}
	if !fm.FieldType.IsValid() {
		return NewValidationError(fmt.Sprintf("invalid field type: %s", fm.FieldType))
	}
	if fm.QualityWeight < 0 || fm.QualityWeight > 100 {
		return NewValidationError("quality weight must be between 0 and 100")
	}
	return nil
}

// IsValid checks if the field type is valid
func (ft FieldType) IsValid() bool {
	switch ft {
	case FieldTypeString, FieldTypeInteger, FieldTypeDecimal, FieldTypeDate, FieldTypeBoolean, FieldTypeJSON:
		return true
	default:
		return false
	}
}

// ValidateValue validates a value against this field's rules
func (fm *FieldMetadata) ValidateValue(value interface{}) error {
	if value == nil {
		if fm.IsRequired {
			return NewValidationError(fmt.Sprintf("required field %s cannot be null", fm.FieldName))
		}
		return nil
	}

	// Type-specific validation would go here
	// In production, this would include proper type checking and conversion

	// Check allowed values (enum)
	if len(fm.AllowedValues) > 0 {
		strValue := fmt.Sprintf("%v", value)
		found := false
		for _, allowed := range fm.AllowedValues {
			if allowed == strValue {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError(fmt.Sprintf("value %v not in allowed values for field %s", value, fm.FieldName))
		}
	}

	// Check min/max for numeric types
	if fm.MinValue != nil || fm.MaxValue != nil {
		// Type assertion and comparison would happen here
		// Simplified for example
	}

	// Check max length for strings
	if fm.MaxLength != nil && fm.FieldType == FieldTypeString {
		if strValue, ok := value.(string); ok {
			if len(strValue) > *fm.MaxLength {
				return NewValidationError(fmt.Sprintf("value length %d exceeds max length %d for field %s", len(strValue), *fm.MaxLength, fm.FieldName))
			}
		}
	}

	return nil
}

// HasEnumValues checks if this field has enum constraints
func (fm *FieldMetadata) HasEnumValues() bool {
	return len(fm.AllowedValues) > 0
}

// HasRangeValidation checks if this field has min/max constraints
func (fm *FieldMetadata) HasRangeValidation() bool {
	return fm.MinValue != nil || fm.MaxValue != nil
}

// HasLengthValidation checks if this field has length constraints
func (fm *FieldMetadata) HasLengthValidation() bool {
	return fm.MaxLength != nil
}
