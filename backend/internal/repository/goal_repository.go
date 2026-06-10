package repository

import (
	"fitnessapi/internal/model"
	"time"

	"gorm.io/gorm"
)

type GoalRepository struct{ db *gorm.DB }

func NewGoalRepository(db *gorm.DB) GoalRepository   { return GoalRepository{db: db} }
func (r GoalRepository) Save(goal *model.Goal) error { return r.db.Create(goal).Error }
func (r GoalRepository) Latest(userID uint) (model.Goal, error) {
	var goal model.Goal
	return goal, r.db.Where("user_id = ?", userID).Order("effective_monday desc").First(&goal).Error
}
func WeekStart(t time.Time) time.Time {
	day := int(t.Weekday()+6) % 7
	return time.Date(t.Year(), t.Month(), t.Day()-day, 0, 0, 0, 0, t.Location())
}
