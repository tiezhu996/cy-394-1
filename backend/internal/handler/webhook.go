package handler

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct{ recordSvc service.RecordService }

func NewWebhookHandler(recordSvc service.RecordService) WebhookHandler {
	return WebhookHandler{recordSvc: recordSvc}
}
func (h WebhookHandler) Device(c *gin.Context) {
	var record model.WorkoutRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.Error(err)
		return
	}
	if err := h.recordSvc.Create(&record); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"accepted": true, "record_id": record.ID})
}
