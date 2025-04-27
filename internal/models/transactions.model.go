package models

type TransactionRes struct{}

type TransactionReq struct {
	AuthID     string `json:"authId" db:"auth_id"`
	PaymentId  int    `json:"paymentId" db:"payment_id"`
	Status     string `json:"status" db:"status"`
	GrandTotal int    `json:"grandTotal" db:"grand_total"`
	ShowingId  int    `json:"showingId" db:"showing_id"`
	SeatId     []int  `json:"seatId" db:"seat_id"`
}
