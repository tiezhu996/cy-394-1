package handler

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
	"fitnessapi/internal/service"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RankHandler struct {
	recordSvc service.RecordService
	goalSvc   service.GoalService
	userRepo  repository.UserRepository
}

func NewRankHandler(recordSvc service.RecordService, goalSvc service.GoalService, userRepo repository.UserRepository) RankHandler {
	return RankHandler{recordSvc: recordSvc, goalSvc: goalSvc, userRepo: userRepo}
}
func (h RankHandler) Rankings(c *gin.Context) {
	scope := c.DefaultQuery("scope", "week")
	friendCircle := c.Query("friend_circle")
	now := time.Now()

	users, err := h.userRepo.List()
	if err != nil {
		c.Error(err)
		return
	}

	type rankItem struct {
		UserID           uint    `json:"user_id"`
		Name             string  `json:"name"`
		Duration         int     `json:"duration"`
		Distance         float64 `json:"distance"`
		Calories         float64 `json:"calories"`
		ConsecutiveWeeks int     `json:"consecutive_weeks,omitempty"`
	}

	var rankings []rankItem

	for _, u := range users {
		records, err := h.recordSvc.List(u.ID, nil, nil)
		if err != nil {
			continue
		}

		var filtered []model.WorkoutRecord
		if scope == "week" {
			weekStart := repository.WeekStart(now)
			weekEnd := weekStart.AddDate(0, 0, 7)
			filtered, _ = h.recordSvc.List(u.ID, &weekStart, &weekEnd)
		} else if scope == "month" {
			year, month, _ := now.Date()
			monthStart := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
			nextMonth := monthStart.AddDate(0, 1, 0)
			filtered, _ = h.recordSvc.List(u.ID, &monthStart, &nextMonth)
		} else {
			filtered = records
		}
		scopeSummary := service.BuildSummary(filtered)

		item := rankItem{
			UserID:   u.ID,
			Name:     u.Name,
			Duration: scopeSummary.TotalDuration,
			Distance: scopeSummary.TotalDistance,
			Calories: scopeSummary.TotalCalories,
		}

		if scope == "streak" {
			goals, err := h.goalSvc.List(u.ID)
			if err == nil {
				item.ConsecutiveWeeks = service.ConsecutiveWeeks(goals, records, now)
			}
		}

		rankings = append(rankings, item)
	}

	if scope == "streak" {
		sort.Slice(rankings, func(i, j int) bool {
			return rankings[i].ConsecutiveWeeks > rankings[j].ConsecutiveWeeks
		})
	} else {
		sort.Slice(rankings, func(i, j int) bool {
			if rankings[i].Duration != rankings[j].Duration {
				return rankings[i].Duration > rankings[j].Duration
			}
			if rankings[i].Calories != rankings[j].Calories {
				return rankings[i].Calories > rankings[j].Calories
			}
			return rankings[i].Distance > rankings[j].Distance
		})
	}

	result := make([]gin.H, 0, len(rankings))
	for _, r := range rankings {
		entry := gin.H{
			"user_id":  r.UserID,
			"name":     r.Name,
			"duration": r.Duration,
			"distance": r.Distance,
			"calories": r.Calories,
		}
		if scope == "streak" {
			entry["consecutive_weeks"] = r.ConsecutiveWeeks
		}
		result = append(result, entry)
	}

	myID, _ := strconv.Atoi(c.DefaultQuery("user_id", "1"))
	myRank := 0
	for i, r := range result {
		if r["user_id"] == uint(myID) {
			myRank = i + 1
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"scope":         scope,
		"friend_circle": friendCircle,
		"my_rank":       myRank,
		"ranking":       result,
	})
}
