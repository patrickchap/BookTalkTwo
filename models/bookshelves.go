package models

type Bookshelf struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Access      string   `json:"access"`
	VolumeCount int      `json:"volumeCount"`
	Volumes     []Volume `json:"volumes"`
}
