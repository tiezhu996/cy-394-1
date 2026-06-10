package handler

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RecordHandler struct{ svc service.RecordService }

func NewRecordHandler(svc service.RecordService) RecordHandler { return RecordHandler{svc: svc} }
func (h RecordHandler) Create(c *gin.Context) {
	var record model.WorkoutRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.Error(err)
		return
	}
	if record.OccurredAt.IsZero() {
		record.OccurredAt = time.Now()
	}
	if err := h.svc.Create(&record); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, record)
}
func (h RecordHandler) List(c *gin.Context) {
	userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
	records, err := h.svc.List(uint(userID), nil, nil)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, records)
}
func (h RecordHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	record, err := h.svc.Find(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&record); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Update(&record); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, record)
}
func (h RecordHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
