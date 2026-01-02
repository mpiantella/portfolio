package domain

import "time"

// FileMetadata represents metadata about an uploaded file
type FileMetadata struct {
	FileID            string            `json:"file_id"`
	FileName          string            `json:"file_name"`
	FilePath          string            `json:"file_path"`
	FileSizeBytes     int64             `json:"file_size_bytes"`
	FileFormat        FileFormat        `json:"file_format"`
	UploadTimestamp   time.Time         `json:"upload_timestamp"`
	ProcessingStatus  ProcessingStatus  `json:"processing_status"`
	ProcessedAt       *time.Time        `json:"processed_at,omitempty"`
	RecordCount       int               `json:"record_count"`
	ErrorMessage      string            `json:"error_message,omitempty"`
	Checksum          string            `json:"checksum"` // SHA-256
	AdditionalMetadata map[string]string `json:"additional_metadata,omitempty"`
}

// FileFormat represents the format of an uploaded file
type FileFormat string

const (
	FileFormatExcelXLSX FileFormat = "xlsx"
	FileFormatExcelXLS  FileFormat = "xls"
	FileFormatCSV       FileFormat = "csv"
	FileFormatJSON      FileFormat = "json"
	FileFormatXML       FileFormat = "xml"
	FileFormatParquet   FileFormat = "parquet"
)

// ProcessingStatus represents the processing status of a file
type ProcessingStatus string

const (
	StatusPending    ProcessingStatus = "pending"
	StatusProcessing ProcessingStatus = "processing"
	StatusCompleted  ProcessingStatus = "completed"
	StatusFailed     ProcessingStatus = "failed"
)

// Validate checks if the file metadata is valid
func (fm *FileMetadata) Validate() error {
	if fm.FileName == "" {
		return NewValidationError("file name cannot be empty")
	}
	if fm.FilePath == "" {
		return NewValidationError("file path cannot be empty")
	}
	if fm.FileSizeBytes <= 0 {
		return NewValidationError("file size must be positive")
	}
	if !fm.FileFormat.IsValid() {
		return NewValidationError("invalid file format")
	}
	if !fm.ProcessingStatus.IsValid() {
		return NewValidationError("invalid processing status")
	}
	return nil
}

// IsValid checks if the file format is valid
func (ff FileFormat) IsValid() bool {
	switch ff {
	case FileFormatExcelXLSX, FileFormatExcelXLS, FileFormatCSV,
		FileFormatJSON, FileFormatXML, FileFormatParquet:
		return true
	default:
		return false
	}
}

// IsValid checks if the processing status is valid
func (ps ProcessingStatus) IsValid() bool {
	switch ps {
	case StatusPending, StatusProcessing, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

// MarkProcessing updates the status to processing
func (fm *FileMetadata) MarkProcessing() {
	fm.ProcessingStatus = StatusProcessing
}

// MarkCompleted updates the status to completed
func (fm *FileMetadata) MarkCompleted(recordCount int) {
	fm.ProcessingStatus = StatusCompleted
	fm.RecordCount = recordCount
	now := time.Now()
	fm.ProcessedAt = &now
}

// MarkFailed updates the status to failed with an error message
func (fm *FileMetadata) MarkFailed(errorMessage string) {
	fm.ProcessingStatus = StatusFailed
	fm.ErrorMessage = errorMessage
	now := time.Now()
	fm.ProcessedAt = &now
}

// IsCompleted returns true if the file has been successfully processed
func (fm *FileMetadata) IsCompleted() bool {
	return fm.ProcessingStatus == StatusCompleted
}

// IsFailed returns true if the file processing failed
func (fm *FileMetadata) IsFailed() bool {
	return fm.ProcessingStatus == StatusFailed
}

// IsProcessing returns true if the file is currently being processed
func (fm *FileMetadata) IsProcessing() bool {
	return fm.ProcessingStatus == StatusProcessing
}

// GetFileExtension returns the file extension
func (fm *FileMetadata) GetFileExtension() string {
	return string(fm.FileFormat)
}

// ParsedData represents the parsed data from a file
type ParsedData struct {
	FileMetadata *FileMetadata              `json:"file_metadata"`
	Headers      []string                   `json:"headers"`
	Records      []map[string]interface{}   `json:"records"`
	RecordCount  int                        `json:"record_count"`
	ParsedAt     time.Time                  `json:"parsed_at"`
	Errors       []ParseError               `json:"errors,omitempty"`
}

// ParseError represents an error that occurred during parsing
type ParseError struct {
	Row     int    `json:"row"`
	Column  string `json:"column"`
	Message string `json:"message"`
}

// NewParsedData creates a new ParsedData instance
func NewParsedData(metadata *FileMetadata, headers []string) *ParsedData {
	return &ParsedData{
		FileMetadata: metadata,
		Headers:      headers,
		Records:      make([]map[string]interface{}, 0),
		ParsedAt:     time.Now(),
		Errors:       make([]ParseError, 0),
	}
}

// AddRecord adds a record to the parsed data
func (pd *ParsedData) AddRecord(record map[string]interface{}) {
	pd.Records = append(pd.Records, record)
	pd.RecordCount++
}

// AddError adds a parse error
func (pd *ParsedData) AddError(row int, column, message string) {
	pd.Errors = append(pd.Errors, ParseError{
		Row:     row,
		Column:  column,
		Message: message,
	})
}

// HasErrors returns true if there are any parse errors
func (pd *ParsedData) HasErrors() bool {
	return len(pd.Errors) > 0
}

// GetErrorCount returns the number of parse errors
func (pd *ParsedData) GetErrorCount() int {
	return len(pd.Errors)
}
