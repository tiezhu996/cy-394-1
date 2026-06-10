package handler

import (
	"fitnessapi/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RankHandler struct{ recordSvc service.RecordService }

func NewRankHandler(recordSvc service.RecordService) RankHandler {
	return RankHandler{recordSvc: recordSvc}
}
func (h RankHandler) Rankings(c *gin.Context) {
	userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
	records, err := h.recordSvc.List(uint(userID), nil, nil)
	if err != nil {
		c.Error(err)
		return
	}
	summary := service.BuildSummary(records)
	c.JSON(http.StatusOK, gin.H{"scope": c.DefaultQuery("scope", "week"), "friend_circle": c.Query("friend_circle"), "ranking": []gin.H{{"user_id": userID, "duration": summary.TotalDuration, "distance": summary.TotalDistance, "calories": summary.TotalCalories}}})
}
