package repository

import (
	"fitnessapi/internal/model"
	"time"

	"gorm.io/gorm"
)

type RecordRepository struct{ db *gorm.DB }

func NewRecordRepository(db *gorm.DB) RecordRepository { return RecordRepository{db: db} }

func (r RecordRepository) Create(record *model.WorkoutRecord) error { return r.db.Create(record).Error }
func (r RecordRepository) Update(record *model.WorkoutRecord) error { return r.db.Save(record).Error }
func (r RecordRepository) Delete(id uint) error                     { return r.db.Delete(&model.WorkoutRecord{}, id).Error }
func (r RecordRepository) Find(id uint) (model.WorkoutRecord, error) {
	var record model.WorkoutRecord
	return record, r.db.First(&record, id).Error
}
func (r RecordRepository) List(userID uint, start, end *time.Time) ([]model.WorkoutRecord, error) {
	query := r.db.Where("user_id = ?", userID).Order("occurred_at desc")
	if start != nil {
		query = query.Where("occurred_at >= ?", *start)
	}
	if end != nil {
		query = query.Where("occurred_at <= ?", *end)
	}
	var records []model.WorkoutRecord
	return records, query.Find(&records).Error
}
