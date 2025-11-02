package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/namespaces"
	"github.com/Unfield/Valdock/permissions"
	"github.com/Unfield/Valdock/response"
	"github.com/Unfield/Valdock/utils"
	"github.com/gin-gonic/gin"
)

type GetACLsResponse struct {
	Count int                   `json:"count" xml:"count"`
	ACLs  []models.ACLUserModel `json:"acls,omitempty" xml:"acls,omitempty"`
}

func (h *Handler) GetACLUsers(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.ACLRead, permissions.ACLAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	acls, err := h.store.ListACLs()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to list acl users")
		return
	}

	response.SendSuccess(c, GetACLsResponse{ACLs: acls, Count: len(acls)})
}

type CreateNewACLRequest struct {
	Username        string   `json:"username"`
	Enabled         bool     `json:"enabled"`
	NoPassword      bool     `json:"no_password"`
	Passwords       []string `json:"passwords"`
	KeyPatterns     []string `json:"key_patterns"`
	ChannelPatterns []string `json:"channel_patterns"`
	AllowedCommands []string `json:"allowed_commands"`
	DeniedCommands  []string `json:"denied_commands"`
}

func (h *Handler) CreateNewACLUser(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.ACLWrite, permissions.ACLAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	var req CreateNewACLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, response.BadRequest, "failed to add acl user: invalid body")
		return
	}

	instanceID := c.Param("id")

	if !h.store.KeyExists(fmt.Sprintf("%s:%s", namespaces.INSTANCES, instanceID)) {
		response.SendError(c, http.StatusBadRequest, response.BadRequest, "failed to add acl user: unknown instance")
		return
	}

	var existingACLUser models.ACLUserModel
	if err := h.store.GetJSON(fmt.Sprintf("%s:%s:%s", namespaces.ACLUSERS, instanceID, req.Username), &existingACLUser); err == nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to add acl user: user already exists")
		return
	}

	id, err := utils.GenerateID()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to add acl user: internal server error")
		return
	}

	var passwordHashes []string

	if len(req.Passwords) > 0 {
		for _, password := range req.Passwords {
			passwordHashes = append(passwordHashes, utils.GenerateValkeyACLHash(password))
		}
	}

	newUser := models.ACLUserModel{
		ID:              id,
		InstanceID:      instanceID,
		Username:        req.Username,
		Enabled:         req.Enabled,
		NoPassword:      req.NoPassword,
		PasswordHashes:  passwordHashes,
		KeyPatterns:     req.KeyPatterns,
		ChannelPatterns: req.ChannelPatterns,
		AllowedCommands: req.AllowedCommands,
		DeniedCommands:  req.DeniedCommands,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	if err := h.store.SetJSON(fmt.Sprintf("%s:%s:%s", namespaces.ACLUSERS, instanceID, req.Username), newUser); err == nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to add acl user: internal server error")
		return
	}

	response.SendSuccess(c, nil)
}

func (h *Handler) UpdateAclUserHandler(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.ACLWrite, permissions.ACLAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	var req CreateNewACLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, response.BadRequest, "failed to update acl user: invalid body")
		return
	}

	instanceID := c.Param("id")

	var aclUser models.ACLUserModel
	if err := h.store.GetJSON(fmt.Sprintf("%s:%s:%s", namespaces.ACLUSERS, instanceID, req.Username), &aclUser); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to update acl user: user does not exists")
		return
	}

	if req.Username != "" {
		aclUser.Username = req.Username
	}

	if req.Enabled != aclUser.Enabled {
		aclUser.Enabled = req.Enabled
	}

	if req.NoPassword != aclUser.NoPassword {
		aclUser.NoPassword = req.NoPassword
	}

	if len(req.Passwords) > 0 {
		for _, password := range req.Passwords {
			aclUser.PasswordHashes = append(aclUser.PasswordHashes, utils.GenerateValkeyACLHash(password))
		}
	}

	if len(req.KeyPatterns) > 0 {
		for _, keyPattern := range req.KeyPatterns {
			aclUser.KeyPatterns = append(aclUser.KeyPatterns, keyPattern)
		}
	}

	if len(req.ChannelPatterns) > 0 {
		for _, channelPattern := range req.ChannelPatterns {
			aclUser.ChannelPatterns = append(aclUser.ChannelPatterns, channelPattern)
		}
	}

	if len(req.AllowedCommands) > 0 {
		for _, allowedCommands := range req.AllowedCommands {
			aclUser.AllowedCommands = append(aclUser.AllowedCommands, allowedCommands)
		}
	}

	if len(req.DeniedCommands) > 0 {
		for _, deniedCommands := range req.DeniedCommands {
			aclUser.DeniedCommands = append(aclUser.DeniedCommands, deniedCommands)
		}
	}

	if err := h.store.SetJSON(fmt.Sprintf("%s:%s:%s", namespaces.ACLUSERS, instanceID, req.Username), aclUser); err == nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to update acl user: internal server error")
		return
	}

	response.SendSuccess(c, nil)
}

func (h *Handler) DeleteAclUserHandler(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.ACLDelete, permissions.ACLAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	var req CreateNewACLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, response.BadRequest, "failed to update acl user: invalid body")
		return
	}

	instanceID := c.Param("id")

	if err := h.store.DeleteKey(fmt.Sprintf("%s:%s:%s", namespaces.ACLUSERS, instanceID, req.Username)); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to delete acl user: internal server error")
		return
	}

	response.SendSuccess(c, nil)
}
