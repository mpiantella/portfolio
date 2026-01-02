package domain

import "fmt"

// DomainError is the base error type for all domain errors
type DomainError struct {
	message string
	code    string
}

func (e *DomainError) Error() string {
	if e.code != "" {
		return fmt.Sprintf("[%s] %s", e.code, e.message)
	}
	return e.message
}

// ValidationError represents a validation error in the domain
type ValidationError struct {
	DomainError
	Field string
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		DomainError: DomainError{
			message: message,
			code:    "VALIDATION_ERROR",
		},
	}
}

// WithField adds field information to the validation error
func (e *ValidationError) WithField(field string) *ValidationError {
	e.Field = field
	return e
}

// ParsingError represents an error that occurred during file parsing
type ParsingError struct {
	DomainError
	Row    int
	Column string
}

// NewParsingError creates a new parsing error
func NewParsingError(message string, row int, column string) *ParsingError {
	return &ParsingError{
		DomainError: DomainError{
			message: message,
			code:    "PARSING_ERROR",
		},
		Row:    row,
		Column: column,
	}
}

// BusinessRuleError represents a business rule violation
type BusinessRuleError struct {
	DomainError
	Rule string
}

// NewBusinessRuleError creates a new business rule error
func NewBusinessRuleError(message, rule string) *BusinessRuleError {
	return &BusinessRuleError{
		DomainError: DomainError{
			message: message,
			code:    "BUSINESS_RULE_ERROR",
		},
		Rule: rule,
	}
}

// NotFoundError represents an entity not found error
type NotFoundError struct {
	DomainError
	EntityType string
	EntityID   string
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(entityType, entityID string) *NotFoundError {
	return &NotFoundError{
		DomainError: DomainError{
			message: fmt.Sprintf("%s with ID %s not found", entityType, entityID),
			code:    "NOT_FOUND",
		},
		EntityType: entityType,
		EntityID:   entityID,
	}
}

// DuplicateError represents a duplicate entity error
type DuplicateError struct {
	DomainError
	EntityType string
	Field      string
	Value      string
}

// NewDuplicateError creates a new duplicate error
func NewDuplicateError(entityType, field, value string) *DuplicateError {
	return &DuplicateError{
		DomainError: DomainError{
			message: fmt.Sprintf("duplicate %s found: %s = %s", entityType, field, value),
			code:    "DUPLICATE",
		},
		EntityType: entityType,
		Field:      field,
		Value:      value,
	}
}
