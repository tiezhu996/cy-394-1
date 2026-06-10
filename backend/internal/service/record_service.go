package service

import (
	"fitnessapi/internal/constants"
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
	"time"
)

type RecordService struct{ repo repository.RecordRepository }

func NewRecordService(repo repository.RecordRepository) RecordService {
	return RecordService{repo: repo}
}

func (s RecordService) Create(record *model.WorkoutRecord) error {
	if record.Calories <= 0 {
		record.Calories = float64(record.DurationMin) / 60 * constants.METValues[constants.SportType(record.SportType)] * 70
	}
	return s.repo.Create(record)
}
func (s RecordService) List(userID uint, start, end *time.Time) ([]model.WorkoutRecord, error) {
	return s.repo.List(userID, start, end)
}
func (s RecordService) Update(record *model.WorkoutRecord) error  { return s.repo.Update(record) }
func (s RecordService) Delete(id uint) error                      { return s.repo.Delete(id) }
func (s RecordService) Find(id uint) (model.WorkoutRecord, error) { return s.repo.Find(id) }
