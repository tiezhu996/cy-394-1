package handler

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
	"fitnessapi/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GoalHandler struct {
	goalSvc   service.GoalService
	recordSvc service.RecordService
	goalRepo  repository.GoalRepository
}

func NewGoalHandler(goalSvc service.GoalService, recordSvc service.RecordService, goalRepo repository.GoalRepository) GoalHandler {
	return GoalHandler{goalSvc: goalSvc, recordSvc: recordSvc, goalRepo: goalRepo}
}
func (h GoalHandler) Save(c *gin.Context) {
	var goal model.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.Error(err)
		return
	}
	if goal.EffectiveMonday.IsZero() {
		goal.EffectiveMonday = repository.WeekStart(time.Now())
	}
	if err := h.goalSvc.Save(&goal); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, goal)
}
func (h GoalHandler) Progress(c *gin.Context) {
	userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
	goal, err := h.goalSvc.Latest(uint(userID))
	if err != nil {
		c.Error(err)
		return
	}
	records, err := h.recordSvc.List(uint(userID), nil, nil)
	if err != nil {
		c.Error(err)
		return
	}
	summary := service.BuildSummary(records)
	progress := service.GoalProgress(goal, summary)

	allGoals, err := h.goalSvc.List(uint(userID))
	if err != nil {
		c.Error(err)
		return
	}
	streak := service.BuildWeeklyStreak(allGoals, records, time.Now(), 4)

	c.JSON(http.StatusOK, gin.H{
		"goal":     goal,
		"progress": progress,
		"streak":   streak,
	})
}
