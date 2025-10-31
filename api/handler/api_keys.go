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

type GetApiKeysResponse struct {
	Count   int             `json:"count" xml:"count"`
	ApiKeys []models.APIKey `json:"api_keys,omitempty" xml:"instances,omitempty"`
}

func (h *Handler) GetApiKeysHandler(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.APIKeyRead, permissions.APIKeyAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	apikeys, err := h.store.ListApiKeys()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to list api keys")
		return
	}

	response.SendSuccess(c, GetApiKeysResponse{ApiKeys: apikeys, Count: len(apikeys)})
}

type CreateApiKeyRequest struct {
	Name        string        `json:"name"`
	Permissions string        `json:"permissions"`
	Enabled     bool          `json:"enabled"`
	ExpiresIn   time.Duration `json:"expires_in"`
}

func (h *Handler) CreateApiKeyHandler(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.APIKeyCreate, permissions.APIKeyAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	var req CreateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, response.BadRequest, "invalid body")
		return
	}

	apiKeyStr, err := utils.NewApiKey()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to create api key")
		return
	}

	apiKeyID, err := utils.GenerateID()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to create api key")
		return
	}

	permissions := permissions.ParsePermissionString(req.Permissions)

	apiKey := models.APIKey{
		ID:          apiKeyID,
		Value:       apiKeyStr,
		Name:        req.Name,
		Enabled:     req.Enabled,
		Permissions: permissions,
		ExpiresIn:   req.ExpiresIn,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	if err := h.store.SetJSON(fmt.Sprintf("%s:%s", namespaces.API_KEYS, apiKey.Value), apiKey); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to create api key")
		return
	}
}

func (h *Handler) GetApiKeyHandler(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.APIKeyRead, permissions.APIKeyAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	apiKeyValue := c.Param("key")

	var apiKey models.APIKey
	if err := h.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.API_KEYS, apiKeyValue), &apiKey); err != nil {
		response.SendError(c, http.StatusBadRequest, response.BadRequest, "api key invalid")
		return
	}

	response.SendSuccess(c, apiKey)
}

func (h *Handler) DeleteApiKeyHandler(c *gin.Context) {
	perms, ok := c.Get("permissions")
	if !ok {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}
	userPerms := perms.([]permissions.Permission)
	if !permissions.HasOnePermission(userPerms, permissions.APIKeyDelete, permissions.APIKeyAdmin) {
		response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
		return
	}

	apiKeyValue := c.Param("key")

	if err := h.store.DeleteKey(fmt.Sprintf("%s:%s", namespaces.API_KEYS, apiKeyValue)); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to delete api key")
		return
	}

	response.SendSuccess(c, nil)
}

func (h *Handler) UpdateApiKeyStateHandler(action string) func(*gin.Context) {
	return func(c *gin.Context) {
		perms, ok := c.Get("permissions")
		if !ok {
			response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
			return
		}
		userPerms := perms.([]permissions.Permission)
		if !permissions.HasOnePermission(userPerms, permissions.APIKeyWrite, permissions.APIKeyAdmin) {
			response.SendError(c, http.StatusForbidden, response.Forbidden, "forbidden")
			return
		}

		apiKeyValue := c.Param("key")

		var apiKey models.APIKey
		if err := h.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.API_KEYS, apiKeyValue), &apiKey); err != nil {
			response.SendError(c, http.StatusBadRequest, response.BadRequest, "api key invalid")
			return
		}

		if action == "disable" {
			if apiKey.Enabled {
				apiKey.Enabled = false
				if err := h.store.SetJSON(fmt.Sprintf("%s:%s", namespaces.API_KEYS, apiKeyValue), apiKey); err != nil {
					response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to disable api key")
					return
				}
			}
			response.SendSuccess(c, nil)
			return
		}
		if action == "enable" {
			if !apiKey.Enabled {
				apiKey.Enabled = true
				if err := h.store.SetJSON(fmt.Sprintf("%s:%s", namespaces.API_KEYS, apiKeyValue), apiKey); err != nil {
					response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to enable api key")
					return
				}
			}
			response.SendSuccess(c, nil)
			return
		}
		if action == "expire" {
			if time.Now().After(apiKey.CreatedAt.Add(apiKey.ExpiresIn)) {
				response.SendSuccess(c, nil)
				return
			}
			apiKey.ExpiresIn = time.Since(apiKey.CreatedAt)
			if err := h.store.SetJSON(fmt.Sprintf("%s:%s", namespaces.API_KEYS, apiKeyValue), apiKey); err != nil {
				response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to expire api key")
				return
			}
		}

		response.SendError(c, http.StatusBadRequest, response.BadRequest, "unknown action")
	}
}
