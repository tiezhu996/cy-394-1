package handler

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
	"fitnessapi/internal/service"
	"net/http"
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
	goal, err := h.goalRepo.Latest(1)
	if err != nil {
		c.Error(err)
		return
	}
	records, err := h.recordSvc.List(goal.UserID, nil, nil)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, service.GoalProgress(goal, service.BuildSummary(records)))
}
