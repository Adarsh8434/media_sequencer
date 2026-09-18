package repository

import (
	"mediaSequencer/backend/internal/model"

	"gorm.io/gorm"
)

type PlaylistRepository struct {
	DB *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) *PlaylistRepository {
	return &PlaylistRepository{DB: db}
}

func (r *PlaylistRepository) GetByWindowID(windowID int) ([]model.Playlist, error) {
	var playlists []model.Playlist

	err := r.DB.
		Where("window_id = ?", windowID).
		Order("position ASC").
		Find(&playlists).Error

	return playlists, err
}

func (r *PlaylistRepository) Create(playlist *model.Playlist) error {
	return r.DB.Create(playlist).Error
}
