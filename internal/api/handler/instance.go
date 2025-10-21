package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Unfield/Valdock/internal/response"
	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/namespaces"
	"github.com/Unfield/Valdock/utils"
	"github.com/gin-gonic/gin"
)

type GetInstancesResponse struct {
	Count     int                    `json:"count" xml:"count"`
	Instances []models.InstanceModel `json:"instances,omitempty" xml:"instances,omitempty"`
}

func (h *Handler) GetInstancesHandler(c *gin.Context) {
	instances, err := h.store.ListInstances()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to list instances")
		return
	}

	response.SendSuccess(c, GetInstancesResponse{Instances: instances, Count: len(instances)})
}

type usernamePasswordField struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
}

type CreateInstanceRequest struct {
	Name           string                  `json:"name" xml:"name"`
	ConfigTemplate string                  `json:"config_template" xml:"config_template"`
	Users          []usernamePasswordField `json:"users" xml:"users"`
}

func (h *Handler) CreateInstanceHandler(c *gin.Context) {
	var req CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, response.BadRequest, "invalid body")
		return
	}

	instanceID, err := utils.GenerateID()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to generate id")
	}

	port, err := h.pa.GetPort()
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to provision a port")
	}

	instance := models.InstanceModel{
		ID:             instanceID,
		ContainerID:    "",
		DataPath:       "",
		Name:           req.Name,
		Port:           port,
		ConfigTemplate: "default",
		Status:         "creating",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.store.SetJSON(fmt.Sprintf("%s:%s", namespaces.INSTANCES, instance.ID), instance); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to save instance to db")
		return
	}

	response.SendSuccess(c, instance)
}

func (h *Handler) GetInstanceHandler(c *gin.Context) {
	var instance models.InstanceModel
	err := h.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.INSTANCES, c.Param("id")), &instance)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to get instances")
		return
	}

	response.SendSuccess(c, instance)
}

func (h *Handler) DeleteInstanceHandler(c *gin.Context) {
	var instance models.InstanceModel
	err := h.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.INSTANCES, c.Param("id")), &instance)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to get instances")
		return
	}

	if err := h.pa.FreePort(instance.Port); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to free port")
		return
	}

	err = h.store.DeleteKey(fmt.Sprintf("%s:%s", namespaces.INSTANCES, c.Param("id")))
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to delete instances")
		return
	}

	response.SendSuccess(c, nil)
}
