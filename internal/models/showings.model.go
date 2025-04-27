package models

import "time"

type Showings struct {
	ID          int       `json:"id" db:"id"`
	Date        time.Time `json:"date" db:"date"`
	Time        time.Time `json:"time" db:"time"`
	CinemaName  string    `json:"cinemaName" db:"cinema_name"`
	CinemaImage string    `json:"cinemaImage" db:"cinema_image"`
	CityName    string    `json:"cityName" db:"city_name"`
	Price       int       `json:"price" db:"price"`
	MovieTitle  string    `json:"movieTitle" db:"movie_title"`
}

type SeatRes struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	IsAvailable bool   `json:"isAvailable" db:"is_available"`
}
