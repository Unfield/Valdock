package handler

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/namespaces"
	"github.com/Unfield/Valdock/response"
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
	Hostname       string                  `json:"hostname" xml:"hostname"`
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

	if req.Hostname == "" {
		req.Hostname = h.cfg.Docker.Instance.DefaultHostname
	}

	instance := models.InstanceModel{
		ID:              instanceID,
		ContainerID:     "",
		DataPath:        path.Join(h.cfg.Docker.Instance.DataPath, instanceID),
		Name:            req.Name,
		Port:            port,
		PrimaryHostname: req.Hostname,
		ConfigTemplate:  "default",
		Status:          "creating",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.store.SetJSON(fmt.Sprintf("%s:%s", namespaces.INSTANCES, instance.ID), instance); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to save instance to db")
		return
	}

	if err := h.jobClient.EnqueueCreateInstance(&instance); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to queue instance creation job")
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
	instanceID := c.Param("id")

	var instance models.InstanceModel
	err := h.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.INSTANCES, instanceID), &instance)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to get instances")
		return
	}

	if err := h.pa.FreePort(instance.Port); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to free port")
		return
	}

	/*
		err = h.store.DeleteKey(fmt.Sprintf("%s:%s", namespaces.INSTANCES, c.Param("id")))
		if err != nil {
			response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to delete instances")
			return
		}
	*/

	if err := h.jobClient.EnqueueDeleteInstance(instanceID); err != nil {
		response.SendError(c, http.StatusInternalServerError, response.InternalServerError, "failed to queue instance deletion job")
		return
	}

	response.SendSuccess(c, nil)
}
