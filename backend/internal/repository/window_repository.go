package repository

import (
	"mediaSequencer/backend/internal/model"

	"gorm.io/gorm"
)

type WindowRepository struct {
	DB *gorm.DB
}

func NewWindowRepository(db *gorm.DB) *WindowRepository {
	return &WindowRepository{DB: db}
}

func (r *WindowRepository) GetAll() ([]model.Window, error) {
	var windows []model.Window

	err := r.DB.Find(&windows).Error

	return windows, err
}
