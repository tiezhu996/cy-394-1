package service

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
)

type GoalService struct{ repo repository.GoalRepository }

func NewGoalService(repo repository.GoalRepository) GoalService { return GoalService{repo: repo} }
func (s GoalService) Save(goal *model.Goal) error               { return s.repo.Save(goal) }
func GoalProgress(goal model.Goal, summary Summary) map[string]float64 {
	return map[string]float64{
		"sessions_percent": percent(float64(summary.Count), float64(goal.WeeklySessions)),
		"minutes_percent":  percent(float64(summary.TotalDuration), float64(goal.WeeklyMinutes)),
		"calories_percent": percent(summary.TotalCalories, goal.WeeklyCalories),
	}
}
func percent(done, target float64) float64 {
	if target <= 0 {
		return 0
	}
	if done/target > 1 {
		return 100
	}
	return done / target * 100
}
