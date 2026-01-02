// Package usecase contains application business logic
// This layer orchestrates the domain entities and services
package usecase

import (
	"context"
	"fmt"

	"portfolio/internal/ingestion-pipeline/domain"
)

// ProcessFileUseCase handles the complete file processing workflow
type ProcessFileUseCase struct {
	fileRepo     domain.FileRepository
	parser       domain.FileParser
	storage      domain.Storage
	logger       domain.Logger
}

// NewProcessFileUseCase creates a new ProcessFileUseCase instance
func NewProcessFileUseCase(
	fileRepo domain.FileRepository,
	parser domain.FileParser,
	storage domain.Storage,
	logger domain.Logger,
) *ProcessFileUseCase {
	return &ProcessFileUseCase{
		fileRepo: fileRepo,
		parser:   parser,
		storage:  storage,
		logger:   logger,
	}
}

// ProcessFileRequest represents the input for processing a file
type ProcessFileRequest struct {
	FileID   string
	FilePath string
}

// ProcessFileResponse represents the output of file processing
type ProcessFileResponse struct {
	FileID      string
	ParsedData  *domain.ParsedData
	RecordCount int
	Success     bool
	Error       string
}

// Execute processes a file through validation, parsing, and metadata extraction
func (uc *ProcessFileUseCase) Execute(ctx context.Context, req *ProcessFileRequest) (*ProcessFileResponse, error) {
	uc.logger.Info(ctx, "Starting file processing", map[string]interface{}{
		"file_id":   req.FileID,
		"file_path": req.FilePath,
	})

	// Get file metadata
	fileMetadata, err := uc.fileRepo.GetFileMetadata(ctx, req.FileID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to get file metadata", err, map[string]interface{}{
			"file_id": req.FileID,
		})
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	// Update status to processing
	fileMetadata.MarkProcessing()
	if err := uc.fileRepo.SaveFileMetadata(ctx, fileMetadata); err != nil {
		return nil, fmt.Errorf("failed to update file status: %w", err)
	}

	// Download file from storage
	file, err := uc.storage.Download(ctx, req.FilePath)
	if err != nil {
		uc.logger.Error(ctx, "Failed to download file", err, map[string]interface{}{
			"file_path": req.FilePath,
		})
		fileMetadata.MarkFailed(fmt.Sprintf("Failed to download file: %v", err))
		_ = uc.fileRepo.SaveFileMetadata(ctx, fileMetadata)
		return &ProcessFileResponse{
			FileID:  req.FileID,
			Success: false,
			Error:   err.Error(),
		}, err
	}
	defer file.Close()

	// Validate file format
	if err := uc.parser.Validate(ctx, file); err != nil {
		uc.logger.Error(ctx, "File validation failed", err, map[string]interface{}{
			"file_id": req.FileID,
		})
		fileMetadata.MarkFailed(fmt.Sprintf("Validation failed: %v", err))
		_ = uc.fileRepo.SaveFileMetadata(ctx, fileMetadata)
		return &ProcessFileResponse{
			FileID:  req.FileID,
			Success: false,
			Error:   fmt.Sprintf("validation failed: %v", err),
		}, nil
	}

	// Reopen file for parsing (after validation consumed it)
	file.Close()
	file, err = uc.storage.Download(ctx, req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to reopen file: %w", err)
	}
	defer file.Close()

	// Parse file
	parsedData, err := uc.parser.Parse(ctx, file, fileMetadata)
	if err != nil {
		uc.logger.Error(ctx, "File parsing failed", err, map[string]interface{}{
			"file_id": req.FileID,
		})
		fileMetadata.MarkFailed(fmt.Sprintf("Parsing failed: %v", err))
		_ = uc.fileRepo.SaveFileMetadata(ctx, fileMetadata)
		return &ProcessFileResponse{
			FileID:  req.FileID,
			Success: false,
			Error:   fmt.Sprintf("parsing failed: %v", err),
		}, nil
	}

	// Check for parsing errors
	if parsedData.HasErrors() {
		uc.logger.Warn(ctx, "File parsed with errors", map[string]interface{}{
			"file_id":     req.FileID,
			"error_count": parsedData.GetErrorCount(),
		})
	}

	// Update file metadata with results
	fileMetadata.MarkCompleted(parsedData.RecordCount)
	if err := uc.fileRepo.SaveFileMetadata(ctx, fileMetadata); err != nil {
		uc.logger.Error(ctx, "Failed to update file metadata", err, map[string]interface{}{
			"file_id": req.FileID,
		})
		return nil, fmt.Errorf("failed to update file metadata: %w", err)
	}

	uc.logger.Info(ctx, "File processing completed", map[string]interface{}{
		"file_id":      req.FileID,
		"record_count": parsedData.RecordCount,
	})

	return &ProcessFileResponse{
		FileID:      req.FileID,
		ParsedData:  parsedData,
		RecordCount: parsedData.RecordCount,
		Success:     true,
	}, nil
}

// ValidateFile validates a file without processing it completely
func (uc *ProcessFileUseCase) ValidateFile(ctx context.Context, req *ProcessFileRequest) error {
	uc.logger.Info(ctx, "Validating file", map[string]interface{}{
		"file_id":   req.FileID,
		"file_path": req.FilePath,
	})

	// Download file from storage
	file, err := uc.storage.Download(ctx, req.FilePath)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer file.Close()

	// Validate file format
	if err := uc.parser.Validate(ctx, file); err != nil {
		uc.logger.Error(ctx, "File validation failed", err, map[string]interface{}{
			"file_id": req.FileID,
		})
		return fmt.Errorf("validation failed: %w", err)
	}

	uc.logger.Info(ctx, "File validation passed", map[string]interface{}{
		"file_id": req.FileID,
	})

	return nil
}
