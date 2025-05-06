package models

import "time"

type TransactionRes struct {
	MovieTitle string    `json:"movie_title"`
	Schedule   time.Time `json:"schedule"`
	Time       string    `json:"time"`
	CinemaName string    `json:"cinema_name"`
	SeatNames  string    `json:"seat_names"`
	GrandTotal int       `json:"grand_total"`
}

type TransactionReq struct {
	PaymentId  int    `json:"paymentId" db:"payment_id"`
	Status     string `json:"status" db:"status"`
	GrandTotal int    `json:"grandTotal" db:"grand_total"`
	ShowingId  int    `json:"showingId" db:"showing_id"`
	SeatId     []int  `json:"seatId" db:"seat_id"`
}
