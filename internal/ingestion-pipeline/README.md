# Data Ingestion Pipeline - Go Implementation

This directory contains a production-ready implementation of the data ingestion pipeline following Clean Architecture principles in Go.

## Overview

This implementation demonstrates:
- **Clean Architecture** with proper separation of concerns
- **Domain-Driven Design** with rich domain models
- **SOLID Principles** throughout the codebase
- **Interface-driven design** for testability and flexibility
- **Go best practices** and idiomatic code

## Architecture Layers

### 1. Domain Layer (`domain/`)

The core business logic with **zero external dependencies**. Contains:

- **Entities**: Core business objects (`entity.go`, `field_metadata.go`, `file_metadata.go`)
- **Value Objects**: Immutable objects (`quality_score.go`, `validation_rule.go`)
- **Repository Interfaces**: Data access contracts (`repository.go`)
- **Service Interfaces**: External service contracts (`service.go`)
- **Domain Errors**: Business error types (`errors.go`)

**Key Files**:
- `entity.go` - Core business entity with driving field
- `quality_score.go` - Data quality scoring models
- `validation_rule.go` - Business validation rules
- `field_metadata.go` - Field definitions and constraints
- `file_metadata.go` - File processing metadata
- `repository.go` - Repository interfaces (implemented in infrastructure)
- `service.go` - Service interfaces (FileParser, DataNormalizer, etc.)
- `errors.go` - Domain-specific error types

### 2. Use Case Layer (`usecase/`)

Application business logic that orchestrates domain entities and services.

**Key Files**:
- `process_file.go` - File processing workflow (validate, parse, store metadata)
- `normalize_data.go` - Data normalization and deduplication
- `score_quality.go` - Data quality scoring and reporting

**Example Use Case Flow**:
```go
// Process File Use Case
ProcessFileUseCase:
  1. Retrieve file metadata from repository
  2. Download file from storage
  3. Validate file format
  4. Parse file content
  5. Update file status
  6. Return parsed data
```

### 3. Infrastructure Layer (To be implemented)

Contains implementations of domain interfaces:

- **PostgreSQL Repository** - Implements `EntityRepository`, `MetadataRepository`, etc.
- **S3 Storage** - Implements `Storage` interface
- **Excel/CSV/JSON Parsers** - Implements `FileParser` interface
- **CloudWatch Logger** - Implements `Logger` interface
- **SNS Notifier** - Implements `Notifier` interface

### 4. Interface Layer (To be implemented)

Entry points to the application:

- **Lambda Handlers** - AWS Lambda function handlers
- **API Gateway** - REST API endpoints
- **Step Functions** - State machine task handlers
- **Event Processors** - S3 event listeners

## Domain Model

### Entity
The core business entity identified by a driving field (e.g., customer ID, order ID).

```go
type Entity struct {
    EntityID          string                 
    DrivingFieldValue string                 // Unique business identifier
    EntityType        string                 
    Attributes        map[string]interface{} // Flexible key-value storage
    SourceFileID      string                 
    CreatedAt         time.Time              
    UpdatedAt         time.Time              
    IsActive          bool                   
}
```

### FieldMetadata
Defines structure, validation rules, and quality weights for fields.

```go
type FieldMetadata struct {
    FieldID         string
    FieldName       string
    FieldType       FieldType  // string, integer, decimal, date, boolean, json
    IsRequired      bool
    IsDrivingField  bool
    FormatPattern   string     // regex validation
    MinValue        *float64
    MaxValue        *float64
    AllowedValues   []string   // enum values
    QualityWeight   float64    // weight in quality calculation
}
```

### QualityScore
Multi-dimensional data quality scoring:

```go
type QualityScore struct {
    EntityID          string
    CompletenessScore float64  // % of required fields populated
    AccuracyScore     float64  // % conforming to patterns/ranges
    ConsistencyScore  float64  // cross-field consistency
    TimelinessScore   float64  // data freshness
    UniquenessScore   float64  // duplicate detection
    ValidityScore     float64  // business rule compliance
    OverallScore      float64  // weighted average
}
```

**Quality Levels**:
- Excellent: 90-100
- Good: 70-89
- Poor: 50-69
- Critical: 0-49

### ValidationRule
Flexible business rule validation:

```go
type ValidationRule struct {
    RuleID         string
    RuleName       string
    RuleType       ValidationRuleType  // format, range, reference, custom
    RuleExpression string              // SQL or custom DSL
    ErrorMessage   string
    Severity       ValidationSeverity  // error, warning, info
}
```

## Use Cases

### 1. Process File Use Case

**Purpose**: Validate, parse, and extract metadata from uploaded files.

**Dependencies**:
- `FileRepository` - File metadata persistence
- `FileParser` - File parsing implementation
- `Storage` - File storage (S3)
- `Logger` - Structured logging

**Workflow**:
```
1. Retrieve file metadata
2. Update status to "processing"
3. Download file from storage
4. Validate file format
5. Parse file content
6. Extract records and metadata
7. Update status to "completed" or "failed"
8. Return parsed data
```

**Usage Example**:
```go
useCase := NewProcessFileUseCase(fileRepo, parser, storage, logger)

request := &ProcessFileRequest{
    FileID:   "file-123",
    FilePath: "landing/upload-2024-01-01.xlsx",
}

response, err := useCase.Execute(ctx, request)
if err != nil {
    // Handle error
}

fmt.Printf("Parsed %d records\n", response.RecordCount)
```

### 2. Normalize Data Use Case

**Purpose**: Transform raw parsed data into normalized domain entities based on a driving field.

**Dependencies**:
- `DataNormalizer` - Normalization logic
- `MetadataRepository` - Field metadata
- `EntityRepository` - Entity persistence
- `Logger` - Structured logging

**Features**:
- Auto-detection of driving field
- Deduplication based on driving field
- Merge with existing entities
- Type conversion and validation

**Usage Example**:
```go
useCase := NewNormalizeDataUseCase(normalizer, metadataRepo, entityRepo, logger)

request := &NormalizeDataRequest{
    ParsedData:       parsedData,
    DrivingFieldName: "customer_id", // or auto-detect
}

response, err := useCase.Execute(ctx, request)
if err != nil {
    // Handle error
}

fmt.Printf("Normalized %d entities\n", response.NormalizedCount)
```

### 3. Score Quality Use Case

**Purpose**: Calculate multi-dimensional quality scores for entities.

**Dependencies**:
- `QualityScorer` - Scoring algorithms
- `MetadataRepository` - Field definitions
- `QualityRepository` - Score persistence
- `Logger` - Structured logging

**Quality Dimensions**:
1. **Completeness** - Required fields populated
2. **Accuracy** - Values conform to patterns
3. **Consistency** - Cross-field coherence
4. **Timeliness** - Data freshness
5. **Uniqueness** - No duplicates
6. **Validity** - Business rule compliance

**Usage Example**:
```go
useCase := NewScoreQualityUseCase(scorer, metadataRepo, qualityRepo, logger)

request := &ScoreQualityRequest{
    Entity: entity,
}

response, err := useCase.Execute(ctx, request)
if err != nil {
    // Handle error
}

fmt.Printf("Overall Quality Score: %.2f%%\n", response.Score.OverallScore)
fmt.Printf("Quality Level: %s\n", response.QualityLevel)
fmt.Printf("Requires Review: %v\n", response.RequiresReview)
```

## Dependency Injection

The architecture uses constructor injection for clean dependency management:

```go
// Wire up dependencies
logger := NewCloudWatchLogger()
storage := NewS3Storage()
parser := NewExcelParser()
fileRepo := NewPostgresFileRepository(db)
metadataRepo := NewPostgresMetadataRepository(db)

// Create use case with injected dependencies
processFileUseCase := NewProcessFileUseCase(
    fileRepo,
    parser,
    storage,
    logger,
)
```

## Testing Strategy

### Unit Tests
Test domain logic in isolation:

```go
func TestEntity_Validate(t *testing.T) {
    entity := &domain.Entity{
        DrivingFieldValue: "",
        EntityType:        "customer",
    }
    
    err := entity.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "driving field value cannot be empty")
}
```

### Integration Tests
Test use cases with mocked dependencies:

```go
func TestProcessFileUseCase_Execute(t *testing.T) {
    // Mock dependencies
    mockFileRepo := &MockFileRepository{}
    mockParser := &MockFileParser{}
    mockStorage := &MockStorage{}
    mockLogger := &MockLogger{}
    
    // Create use case
    useCase := NewProcessFileUseCase(
        mockFileRepo,
        mockParser,
        mockStorage,
        mockLogger,
    )
    
    // Test execution
    request := &ProcessFileRequest{
        FileID:   "test-file",
        FilePath: "test/file.xlsx",
    }
    
    response, err := useCase.Execute(context.Background(), request)
    assert.NoError(t, err)
    assert.True(t, response.Success)
}
```

## Benefits of This Architecture

### 1. **Testability**
- Domain logic testable without external dependencies
- Use cases testable with mocks
- Integration tests for infrastructure layer

### 2. **Flexibility**
- Easy to swap AWS for GCP or Azure
- Can change database without touching business logic
- Add new file formats by implementing FileParser interface

### 3. **Maintainability**
- Business logic centralized in domain and use cases
- Infrastructure changes don't affect core logic
- Clear separation of concerns

### 4. **Scalability**
- Stateless use cases enable horizontal scaling
- Repository pattern abstracts data access
- Easy to add caching, connection pooling

### 5. **Security**
- Input validation in domain layer
- Clear boundaries prevent injection attacks
- Dependency injection prevents tight coupling

## Best Practices Demonstrated

✅ **Clean Architecture** - Dependency rule strictly followed  
✅ **SOLID Principles** - Single responsibility, interface segregation  
✅ **DDD** - Rich domain models, ubiquitous language  
✅ **Go Idioms** - Error handling, interfaces, composition  
✅ **Testability** - Constructor injection, interface-based design  
✅ **Documentation** - Comprehensive comments and examples  
✅ **Error Handling** - Custom error types, error wrapping  
✅ **Logging** - Structured logging with context  
✅ **Type Safety** - Strong typing, enums, validation  

## Next Steps

To complete the implementation:

1. **Infrastructure Layer**
   - PostgreSQL repository implementations
   - S3 storage implementation
   - File parsers (Excel, CSV, JSON, XML, Parquet)
   - CloudWatch logger implementation

2. **Interface Layer**
   - AWS Lambda handlers
   - API Gateway integrations
   - Step Functions task handlers
   - Event processors

3. **Testing**
   - Unit tests for domain entities
   - Integration tests for use cases
   - End-to-end tests for complete workflows

4. **Deployment**
   - Terraform/CDK infrastructure code
   - CI/CD pipeline configuration
   - Docker containers for local development

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html) by Eric Evans
- [Effective Go](https://golang.org/doc/effective_go.html) - Go best practices
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)

## Author

Software Architect & AWS Solutions Architect  
Demonstrated expertise in: Go, Clean Architecture, AWS, Domain-Driven Design
