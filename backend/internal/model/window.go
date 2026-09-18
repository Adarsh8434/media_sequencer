package model

type Window struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Window) TableName() string {
	return "windows"
}
