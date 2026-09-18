package service

import (
	"sync"
	"time"
)

type SyncState struct {
	MediaID         int
	DurationSeconds int
	StartedAt       time.Time
}

type SyncService struct {
	mu    sync.RWMutex
	state *SyncState
}

func NewSyncService() *SyncService {
	return &SyncService{}
}

func (s *SyncService) Start(mediaID, durationSeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = &SyncState{
		MediaID:         mediaID,
		DurationSeconds: durationSeconds,
		StartedAt:       time.Now(),
	}
}

func (s *SyncService) Get() *SyncState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.state == nil {
		return nil
	}

	if time.Since(s.state.StartedAt).Seconds() >= float64(s.state.DurationSeconds) {
		return nil
	}

	return s.state
}
