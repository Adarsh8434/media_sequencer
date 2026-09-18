package model

type Playlist struct {
	ID       int `json:"id"`
	WindowID int `json:"windowId"`
	MediaID  int `json:"mediaId"`
	Position int `json:"position"`
}

func (Playlist) TableName() string {
	return "playlists"
}
