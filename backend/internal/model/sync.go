package model

type SyncRequest struct {
	MediaID         int `json:"mediaId"`
	DurationSeconds int `json:"durationSeconds"`
}

type SyncState struct {
	MediaID         int   `json:"mediaId"`
	DurationSeconds int   `json:"durationSeconds"`
	StartedAt       int64 `json:"startedAt"`
}
