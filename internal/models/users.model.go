package models

import "mime/multipart"

type UserReq struct {
	Email    string `json:"email" db:"email" binding:"required"`
	Password string `json:"password" db:"password" binding:"required,min=8"`
	Role     string `json:"role" db:"role" binding:"omitempty"`
}

type UserRes struct {
	Email  string `json:"email" db:"email"`
	Role   string `json:"role" db:"role"`
	AuthID string `json:"id" db:"auth_id"`
	Pass   string `json:"password" db:"password"`
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
	Firstname string                `form:"firstname" binding:"omitempty"`
	Lastname  string                `form:"lastname" binding:"omitempty"`
	Phone     string                `from:"phone" binding:"omitempty"`
	Picture   *multipart.FileHeader `form:"picture" binding:"omitempty"`
}

type IdParams struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}
