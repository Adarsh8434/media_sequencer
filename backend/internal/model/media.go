package model

type Media struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	URL             string `json:"url"`
	DurationSeconds int    `json:"durationSeconds"`
}

func (Media) TableName() string {
	return "media"
}
