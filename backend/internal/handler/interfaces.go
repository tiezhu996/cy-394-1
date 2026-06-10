package handler

import (
	"fitnessapi/internal/model"
	"time"
)

type GoalService interface {
	Save(goal *model.Goal) error
	Latest(userID uint) (model.Goal, error)
	List(userID uint) ([]model.Goal, error)
}

type RecordService interface {
	Create(record *model.WorkoutRecord) error
	List(userID uint, start, end *time.Time) ([]model.WorkoutRecord, error)
	Update(record *model.WorkoutRecord) error
	Delete(id uint) error
	Find(id uint) (model.WorkoutRecord, error)
}

type UserRepository interface {
	List() ([]model.User, error)
	Find(id uint) (model.User, error)
}
