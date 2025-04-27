package models

type UserReq struct {
	Email    string `json:"email" db:"email" binding:"required,email"`
	Password string `json:"password" db:"password" binding:"required"`
	Role     string `json:"role" db:"role" binding:"omitempty"`
}

type UserRes struct {
	Email  string `json:"email" db:"email"`
	Role   string `json:"role" db:"role"`
	AuthID string `json:"authId" db:"auth_id"`
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
	Firstname string `json:"firstname" binding:"omitempty"`
	Lastname  string `json:"lastname" binding:"omitempty"`
	Phone     string `json:"phone" binding:"omitempty"`
	Picture   string `json:"picture" binding:"omitempty"`
}

type IdParams struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}
