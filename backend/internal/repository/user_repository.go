package repository

import (
	"fitnessapi/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return UserRepository{db: db} }
func (r UserRepository) List() ([]model.User, error) {
	var users []model.User
	return users, r.db.Order("id asc").Find(&users).Error
}
func (r UserRepository) Find(id uint) (model.User, error) {
	var user model.User
	return user, r.db.First(&user, id).Error
}
