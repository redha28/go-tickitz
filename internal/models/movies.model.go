package models

import "time"

type Movies struct {
	ID       int       `json:"id" db:"id"`
	Title    string    `json:"title" db:"title"`
	Image    string    `json:"image" db:"image"`
	Banner   string    `json:"banner" db:"banner"`
	Release  time.Time `json:"release" db:"release"`
	Genre    string    `json:"genre" db:"genre"`
	Synopsis string    `json:"synopsis" db:"synopsis"`
	Duration int       `json:"duration" db:"duration"`
	Director string    `json:"director" db:"director"`
	Rating   float64   `json:"rating" db:"rating"`
}

type ReqMetaMovies struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
