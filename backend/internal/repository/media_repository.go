package repository

import (
	"mediaSequencer/backend/internal/model"

	"gorm.io/gorm"
)

type MediaRepository struct {
	DB *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{DB: db}
}

func (r *MediaRepository) GetAll() ([]model.Media, error) {
	var media []model.Media

	err := r.DB.Find(&media).Error

	return media, err
}

func (r *MediaRepository) Create(media *model.Media) error {
	return r.DB.Create(media).Error
}
