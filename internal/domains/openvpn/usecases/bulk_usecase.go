package usecases

import (
	"context"
	openvpndto "system-portal/internal/domains/openvpn/dto"
)

// BulkUsecase defines operations for bulk processing of users and groups
type BulkUsecase interface {
	// =================== BULK USER OPERATIONS ===================

	// BulkCreateUsers creates multiple users in a single operation
	// Returns detailed results for each user creation attempt
	BulkCreateUsers(ctx context.Context, req *openvpndto.BulkCreateUsersRequest) (*openvpndto.BulkCreateUsersResponse, error)

	// BulkUserActions performs the same action on multiple users
	// Supported actions: enable, disable, reset-otp
	BulkUserActions(ctx context.Context, req *openvpndto.BulkUserActionsRequest) (*openvpndto.BulkActionResponse, error)

	// BulkExtendUsers extends expiration date for multiple users
	BulkExtendUsers(ctx context.Context, req *openvpndto.BulkUserExtendRequest) (*openvpndto.BulkActionResponse, error)

	// ImportUsers imports users from uploaded file (CSV, JSON, XLSX)
	// Supports dry-run mode for validation only
	ImportUsers(ctx context.Context, req *openvpndto.ImportUsersRequest) (*openvpndto.ImportResponse, error)

	// GenerateUserTemplate generates template file for user import
	// Returns filename and file content for download
	GenerateUserTemplate(format string) (filename string, content []byte, error error)

	// =================== BULK GROUP OPERATIONS ===================

	// BulkCreateGroups creates multiple groups in a single operation
	// Returns detailed results for each group creation attempt
	BulkCreateGroups(ctx context.Context, req *openvpndto.BulkCreateGroupsRequest) (*openvpndto.BulkCreateGroupsResponse, error)

	// BulkGroupActions performs the same action on multiple groups
	// Supported actions: enable, disable
	BulkGroupActions(ctx context.Context, req *openvpndto.BulkGroupActionsRequest) (*openvpndto.BulkGroupActionResponse, error)

	// ImportGroups imports groups from uploaded file (CSV, JSON, XLSX)
	// Supports dry-run mode for validation only
	ImportGroups(ctx context.Context, req *openvpndto.ImportGroupsRequest) (*openvpndto.ImportResponse, error)

	// GenerateGroupTemplate generates template file for group import
	// Returns filename and file content for download
	GenerateGroupTemplate(format string) (filename string, content []byte, error error)

	// =================== BULK VALIDATION & UTILITIES ===================

	// ValidateUserBatch validates a batch of users before processing
	// Returns validation errors and cleaned data
	ValidateUserBatch(users []openvpndto.CreateUserRequest) (valid []openvpndto.CreateUserRequest, errors []openvpndto.ImportValidationError, err error)

	// ValidateGroupBatch validates a batch of groups before processing
	// Returns validation errors and cleaned data
	ValidateGroupBatch(groups []openvpndto.CreateGroupRequest) (valid []openvpndto.CreateGroupRequest, errors []openvpndto.ImportValidationError, err error)

	// ParseImportFile parses uploaded file and returns structured data
	// Supports multiple file formats and validates content
	ParseImportFile(filename string, content []byte, format string, entityType string) (interface{}, []openvpndto.ImportValidationError, error)

	// =================== BATCH REPORTING ===================

	// GetBulkOperationStatus gets status of ongoing bulk operations
	// Returns progress information for long-running operations
	GetBulkOperationStatus(ctx context.Context, operationId string) (interface{}, error)

	// GetBulkOperationHistory gets history of bulk operations
	// Returns list of previous bulk operations with results
	GetBulkOperationHistory(ctx context.Context, entityType string, limit int) ([]interface{}, error)
}
