package models

type UserReq struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required"`
	Role     string `json:"role" db:"role"`
}

type UserRes struct {
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}

type ProfileRes struct {
	Firstname string `json:"firstname" db:"firstname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Picture   string `json:"picture" db:"picture"`
	Phone     string `json:"phone" db:"phone"`
	Point     int    `json:"point" db:"point"`
	Email     string `json:"email" db:"email"`
}

type UpdateProfileReq struct {
	Firstname string `json:"firstname" binding:"required"`
}

type IdParams struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}
