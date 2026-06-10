package handler

import (
	"fitnessapi/internal/constants"
	"fitnessapi/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct{ recordSvc service.RecordService }

func NewStatsHandler(recordSvc service.RecordService) StatsHandler {
	return StatsHandler{recordSvc: recordSvc}
}
func (h StatsHandler) Types(c *gin.Context) { c.JSON(http.StatusOK, constants.METValues) }
func (h StatsHandler) Stats(c *gin.Context) {
	userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
	records, err := h.recordSvc.List(uint(userID), nil, nil)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, service.BuildSummary(records))
}
func (h StatsHandler) PR(c *gin.Context) {
	userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
	records, err := h.recordSvc.List(uint(userID), nil, nil)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, service.PersonalRecords(records))
}
