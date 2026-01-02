// Package domain contains the core business entities and rules
// This layer has NO external dependencies - only standard library
package domain

import (
	"time"
)

// Entity represents a core business entity identified by a driving field
type Entity struct {
	EntityID          string                 `json:"entity_id"`
	DrivingFieldValue string                 `json:"driving_field_value"`
	EntityType        string                 `json:"entity_type"`
	Attributes        map[string]interface{} `json:"attributes"`
	SourceFileID      string                 `json:"source_file_id"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	IsActive          bool                   `json:"is_active"`
}

// Validate checks if the entity is valid according to business rules
func (e *Entity) Validate() error {
	if e.DrivingFieldValue == "" {
		return NewValidationError("driving field value cannot be empty")
	}
	if e.EntityType == "" {
		return NewValidationError("entity type cannot be empty")
	}
	return nil
}

// IsComplete checks if all required fields are populated
func (e *Entity) IsComplete(requiredFields []string) bool {
	for _, field := range requiredFields {
		if val, exists := e.Attributes[field]; !exists || val == nil {
			return false
		}
	}
	return true
}

// GetAttribute safely retrieves an attribute value
func (e *Entity) GetAttribute(key string) (interface{}, bool) {
	val, exists := e.Attributes[key]
	return val, exists
}

// SetAttribute safely sets an attribute value
func (e *Entity) SetAttribute(key string, value interface{}) {
	if e.Attributes == nil {
		e.Attributes = make(map[string]interface{})
	}
	e.Attributes[key] = value
	e.UpdatedAt = time.Now()
}

// AttributeCount returns the number of attributes
func (e *Entity) AttributeCount() int {
	return len(e.Attributes)
}

// MarkUpdated updates the UpdatedAt timestamp
func (e *Entity) MarkUpdated() {
	e.UpdatedAt = time.Now()
}

// Deactivate marks the entity as inactive
func (e *Entity) Deactivate() {
	e.IsActive = false
	e.UpdatedAt = time.Now()
}

// Activate marks the entity as active
func (e *Entity) Activate() {
	e.IsActive = true
	e.UpdatedAt = time.Now()
}
