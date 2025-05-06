package models

type AdminCreateMovieReq struct {
	Title    string  `json:"title" db:"title" binding:"required"`
	Image    string  `json:"image" db:"image" binding:"required"`
	Banner   string  `json:"banner" db:"banner" binding:"required"`
	Release  string  `json:"release" db:"release" binding:"required"`
	Synopsis string  `json:"synopsis" db:"synopsis" binding:"required"`
	Duration int     `json:"duration" db:"duration" binding:"required"`
	Director string  `json:"director" db:"director" binding:"required"`
	Genres   []int   `json:"genres" db:"-" binding:"required"`
	Rating   float64 `json:"rating" db:"rating" binding:"required"`
	// CreatedBy string  `json:"createdBy" db:"created_by" binding:"required"`
}

type GenresReq struct {
	MoviesID int `json:"movies_id" db:"movies_id" binding:"required"`
	GenreId  int `json:"genre_id" db:"genre_id" binding:"required"`
}

type AdminEditMovieReq struct {
	Title    string  `json:"title" db:"title" binding:"omitempty"`
	Image    string  `json:"image" db:"image" binding:"omitempty"`
	Banner   string  `json:"banner" db:"banner" binding:"omitempty"`
	Release  string  `json:"release" db:"release" binding:"omitempty"`
	Synopsis string  `json:"synopsis" db:"synopsis" binding:"omitempty"`
	Duration int     `json:"duration" db:"duration" binding:"omitempty"`
	Director string  `json:"director" db:"director" binding:"omitempty"`
	Genres   []int   `json:"genres" db:"-" binding:"omitempty"`
	Rating   float64 `json:"rating" db:"rating" binding:"omitempty"`
	// CreatedBy string  `json:"createdBy" db:"created_by" binding:"required"`
}

type AdminDeleteMovieReq struct {
	MoviesID int `json:"moviesId" db:"movies_id" binding:"required"`
}
