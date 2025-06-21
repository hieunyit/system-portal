package handlers

import (
	"fmt"
	nethttp "net/http"
	"strings"
	dto "system-portal/internal/domains/openvpn/dto"
	"system-portal/internal/domains/openvpn/entities"
	"system-portal/internal/domains/openvpn/usecases"
	"system-portal/internal/shared/errors"
	"system-portal/internal/shared/infrastructure/xmlrpc"
	http "system-portal/internal/shared/response"
	"system-portal/pkg/logger"
	"system-portal/pkg/validator"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupUsecase  usecases.GroupUsecase
	configUsecase usecases.ConfigUsecase
	xmlrpcClient  *xmlrpc.Client
}

func NewGroupHandler(groupUsecase usecases.GroupUsecase, configUsecase usecases.ConfigUsecase, xmlrpcClient *xmlrpc.Client) *GroupHandler {
	return &GroupHandler{
		groupUsecase:  groupUsecase,
		configUsecase: configUsecase,
		xmlrpcClient:  xmlrpcClient,
	}
}

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new VPN user group
// @Tags Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.VpnCreateGroupRequest true "Group creation data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/openvpn/groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req dto.VpnCreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("Failed to bind create group request")
		http.RespondWithError(c, errors.BadRequest("Invalid request format", err))
		return
	}

	// Set default values
	if req.MFA == nil {
		defaultMFA := true
		req.MFA = &defaultMFA
	}
	if req.Role == "" {
		req.Role = entities.UserRoleUser
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		logger.Log.WithError(err).Error("Create group request validation failed")
		http.RespondWithValidationError(c, err)
		return
	}

	// BASIC FIX: Validate reserved group names
	if h.isReservedGroupName(req.GroupName) {
		http.RespondWithError(c, errors.BadRequest("Group name is reserved and cannot be used", nil))
		return
	}

	// Validate GroupSubnet and GroupRange with current state
	// Need to get effective values after merge to validate properly
	var effectiveGroupSubnet, effectiveGroupRange []string

	if req.GroupSubnet != nil {
		effectiveGroupSubnet = req.GroupSubnet
	} else {
		effectiveGroupSubnet = []string{}
	}

	if req.GroupRange != nil {
		effectiveGroupRange = req.GroupRange
	} else {
		effectiveGroupRange = []string{}
	}

	if err := h.validateGroupSubnetAndRange(effectiveGroupSubnet, effectiveGroupRange); err != nil {
		http.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}

	// Convert DTO to entity
	group := &entities.Group{
		GroupName:     req.GroupName,
		AuthMethod:    req.AuthMethod,
		AccessControl: req.AccessControl,
		Role:          req.Role,
		GroupSubnet:   req.GroupSubnet,
		GroupRange:    req.GroupRange,
	}

	// Set MFA
	group.SetMFA(*req.MFA)

	// Create group
	if err := h.groupUsecase.CreateGroup(c.Request.Context(), group); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.RespondWithError(c, appErr)
		} else {
			http.RespondWithError(c, errors.InternalServerError("Failed to create group", err))
		}
		return
	}

	http.RespondWithMessage(c, nethttp.StatusCreated, "Group created successfully")
}

// GetGroup godoc
// @Summary Get group by name
// @Description Get group information by name
// @Tags Groups
// @Security BearerAuth
// @Produce json
// @Param groupName path string true "Group name"
// @Success 200 {object} response.SuccessResponse{data=dto.VpnGroupResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /api/openvpn/groups/{groupName} [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	groupName := c.Param("groupName")
	if groupName == "" {
		http.RespondWithError(c, errors.BadRequest("Group name is required", nil))
		return
	}

	group, err := h.groupUsecase.GetGroup(c.Request.Context(), groupName)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.RespondWithError(c, appErr)
		} else {
			http.RespondWithError(c, errors.InternalServerError("Failed to get group", err))
		}
		return
	}

	// Convert entity to DTO
	response := dto.VpnGroupResponse{
		GroupName:     group.GroupName,
		AuthMethod:    group.AuthMethod,
		MFA:           group.MFA == "true",
		Role:          group.Role,
		DenyAccess:    group.DenyAccess == "true",
		AccessControl: group.AccessControl,
		GroupSubnet:   group.GroupSubnet,
		GroupRange:    group.GroupRange,
	}

	http.RespondWithSuccess(c, nethttp.StatusOK, response)
}

// UpdateGroup godoc
// @Summary Update group
// @Description Update group information
// @Tags Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param groupName path string true "Group name"
// @Param request body dto.VpnUpdateGroupRequest true "Group update data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/openvpn/groups/{groupName} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	groupName := c.Param("groupName")
	if groupName == "" {
		http.RespondWithError(c, errors.BadRequest("Group name is required", nil))
		return
	}

	// BASIC FIX: Check if group is system group
	if h.isSystemGroup(groupName) {
		http.RespondWithError(c, errors.BadRequest("Cannot modify system group", nil))
		return
	}

	var req dto.VpnUpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("Failed to bind update group request")
		http.RespondWithError(c, errors.BadRequest("Invalid request format", err))
		return
	}

	// Set default Role if empty
	if req.Role == "" {
		req.Role = entities.UserRoleUser
	}

	// Validate request
	if err := validator.Validate(&req); err != nil {
		logger.Log.WithError(err).Error("Update group request validation failed")
		http.RespondWithValidationError(c, err)
		return
	}

	// Validate GroupSubnet and GroupRange - Note: Full validation with conflict check is done in usecase

	// Convert DTO to entity
	group := &entities.Group{
		GroupName: groupName,
		Role:      req.Role,
	}

	// Handle arrays: nil vs [] semantics
	// nil = preserve existing, [] = clear, [values] = replace
	if req.AccessControl != nil {
		group.AccessControl = req.AccessControl
	}
	if req.GroupSubnet != nil {
		group.GroupSubnet = req.GroupSubnet
	}
	if req.GroupRange != nil {
		group.GroupRange = req.GroupRange
	}

	if req.DenyAccess != nil {
		group.SetDenyAccess(*req.DenyAccess)
	}

	if req.MFA != nil {
		group.SetMFA(*req.MFA)
	}

	// Update group
	if err := h.groupUsecase.UpdateGroup(c.Request.Context(), group); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.RespondWithError(c, appErr)
		} else {
			http.RespondWithError(c, errors.InternalServerError("Failed to update group", err))
		}
		return
	}

	http.RespondWithMessage(c, nethttp.StatusOK, "Group updated successfully")
}

// DeleteGroup godoc
// @Summary Delete group
// @Description Delete group by name
// @Tags Groups
// @Security BearerAuth
// @Produce json
// @Param groupName path string true "Group name"
// @Success 200 {object} response.SuccessResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/openvpn/groups/{groupName} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	groupName := c.Param("groupName")
	if groupName == "" {
		http.RespondWithError(c, errors.BadRequest("Group name is required", nil))
		return
	}

	// BASIC FIX: Check if group is system group
	if h.isSystemGroup(groupName) {
		http.RespondWithError(c, errors.BadRequest("Cannot delete system group", nil))
		return
	}

	if err := h.groupUsecase.DeleteGroup(c.Request.Context(), groupName); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.RespondWithError(c, appErr)
		} else {
			http.RespondWithError(c, errors.InternalServerError("Failed to delete group", err))
		}
		return
	}

	http.RespondWithMessage(c, nethttp.StatusOK, "Group deleted successfully")
}

// ListGroups godoc
// @Summary List groups
// @Description List groups with pagination and filtering
// @Tags Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query dto.VpnGroupFilter false "Filter parameters"
// @Success 200 {object} response.SuccessResponse{data=dto.VpnGroupListResponse}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/openvpn/groups [get]
func (h *GroupHandler) ListGroups(c *gin.Context) {
	var filter dto.VpnGroupFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Log.WithError(err).Error("Failed to bind group filter")
		http.RespondWithError(c, errors.BadRequest("Invalid filter parameters", err))
		return
	}

	// Validate filter
	if err := validator.Validate(&filter); err != nil {
		logger.Log.WithError(err).Error("Group filter validation failed")
		http.RespondWithValidationError(c, err)
		return
	}

	// Convert DTO filter to entity filter
	entityFilter := &entities.GroupFilter{
		GroupName:  filter.GroupName,
		AuthMethod: filter.AuthMethod,
		Role:       filter.Role,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Offset:     (filter.Page - 1) * filter.Limit,
	}

	groups, totalCount, err := h.groupUsecase.ListGroupsWithTotal(c.Request.Context(), entityFilter)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.RespondWithError(c, appErr)
		} else {
			http.RespondWithError(c, errors.InternalServerError("Failed to list groups", err))
		}
		return
	}

	// Convert entities to DTOs
	var groupResponses []dto.VpnGroupResponse
	for _, group := range groups {
		groupResponses = append(groupResponses, dto.VpnGroupResponse{
			GroupName:     group.GroupName,
			AuthMethod:    group.AuthMethod,
			MFA:           group.MFA == "true",
			Role:          group.Role,
			DenyAccess:    group.DenyAccess == "true",
			AccessControl: group.AccessControl,
			GroupSubnet:   group.GroupSubnet,
			GroupRange:    group.GroupRange,
		})
	}

	response := dto.VpnGroupListResponse{
		Groups: groupResponses,
		Total:  totalCount,
		Page:   filter.Page,
		Limit:  filter.Limit,
	}

	http.RespondWithSuccess(c, nethttp.StatusOK, response)
}

// GroupAction godoc
// @Summary Perform action on group
// @Description Enable or disable a group
// @Tags Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param groupName path string true "Group name"
// @Param action path string true "Action (enable/disable)"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/openvpn/groups/{groupName}/{action} [put]
func (h *GroupHandler) GroupAction(c *gin.Context) {
	groupName := c.Param("groupName")
	action := c.Param("action")

	if groupName == "" {
		http.RespondWithError(c, errors.BadRequest("Group name is required", nil))
		return
	}

	// BASIC FIX: Check if group is system group
	if h.isSystemGroup(groupName) {
		http.RespondWithError(c, errors.BadRequest("Cannot modify system group", nil))
		return
	}

	if action != "enable" && action != "disable" {
		http.RespondWithError(c, errors.BadRequest("Invalid action. Allowed actions: enable, disable", nil))
		return
	}

	// Get existing group to check current state
	existingGroup, err := h.groupUsecase.GetGroup(c.Request.Context(), groupName)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.RespondWithError(c, appErr)
		} else {
			http.RespondWithError(c, errors.InternalServerError("Failed to get group", err))
		}
		return
	}

	switch action {
	case "enable":
		if existingGroup.DenyAccess != "true" {
			http.RespondWithError(c, errors.BadRequest("Group is already enabled", nil))
			return
		}

		group := &entities.Group{
			GroupName: groupName,
		}
		group.SetDenyAccess(false)

		if err := h.groupUsecase.UpdateGroup(c.Request.Context(), group); err != nil {
			if appErr, ok := err.(*errors.AppError); ok {
				http.RespondWithError(c, appErr)
			} else {
				http.RespondWithError(c, errors.InternalServerError("Failed to enable group", err))
			}
			return
		}
		http.RespondWithMessage(c, nethttp.StatusOK, "Group enabled successfully")

	case "disable":
		if existingGroup.DenyAccess == "true" {
			http.RespondWithError(c, errors.BadRequest("Group is already disabled", nil))
			return
		}

		group := &entities.Group{
			GroupName: groupName,
		}
		group.SetDenyAccess(true)

		if err := h.groupUsecase.UpdateGroup(c.Request.Context(), group); err != nil {
			if appErr, ok := err.(*errors.AppError); ok {
				http.RespondWithError(c, appErr)
			} else {
				http.RespondWithError(c, errors.InternalServerError("Failed to disable group", err))
			}
			return
		}
		http.RespondWithMessage(c, nethttp.StatusOK, "Group disabled successfully")
	}
}

// validateGroupSubnetAndRange validates GroupSubnet and GroupRange according to business rules
// NOTE: This function is now moved to usecase layer for better separation of concerns
// and to include comprehensive conflict checking with existing groups
func (h *GroupHandler) validateGroupSubnetAndRange(groupSubnets, groupRanges []string) error {
	// Basic validation can be done here, but comprehensive validation
	// including conflict checking is now handled in usecase layer
	if len(groupSubnets) == 0 && len(groupRanges) > 0 {
		return fmt.Errorf("GroupRange requires GroupSubnet to be specified")
	}

	// For detailed validation including overlap checking with existing groups,
	// see validateGroupSubnetAndRangeWithConflictCheck in group_usecase_impl.go
	return nil
}

// BASIC FIX: Helper functions to validate group names
func (h *GroupHandler) isReservedGroupName(groupName string) bool {
	reservedNames := []string{"__DEFAULT__", "admin", "root", "system", "default"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(groupName, reserved) {
			return true
		}
	}
	return false
}

func (h *GroupHandler) isSystemGroup(groupName string) bool {
	systemGroups := []string{"__DEFAULT__", "admin", "system"}
	for _, systemGroup := range systemGroups {
		if strings.EqualFold(groupName, systemGroup) {
			return true
		}
	}
	return false
}
