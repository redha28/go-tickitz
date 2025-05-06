package handler

import (
	"gotickitz/internal/models"
	"gotickitz/internal/repositories"
	"gotickitz/pkg"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionsHandler struct {
	transactionsRepo *repositories.TransactionsRepository
}

func NewTransactionsHandler(transactionsRepo *repositories.TransactionsRepository) *TransactionsHandler {
	return &TransactionsHandler{transactionsRepo}
}

func (t *TransactionsHandler) CreateTransactionHandler(ctx *gin.Context) {
	var transactionReq models.TransactionReq
	payloads, _ := ctx.Get("payloads")
	userPayload := payloads.(*pkg.Payload)
	if err := ctx.ShouldBindJSON(&transactionReq); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	soldSeats, err := t.transactionsRepo.UseGetSeatsSold(ctx.Request.Context(), transactionReq.ShowingId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	for _, seat := range soldSeats {
		for _, reqSeat := range transactionReq.SeatId {
			if seat.ID == reqSeat {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Seat already sold"})
				return
			}
		}
	}

	if err := t.transactionsRepo.UseCreateTransaction(ctx.Request.Context(), &transactionReq, userPayload.Id); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := t.transactionsRepo.UseAddPoints(ctx.Request.Context(), userPayload.Id, len(transactionReq.SeatId)); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})

}

func (t *TransactionsHandler) GetTransactionsHandler(ctx *gin.Context) {
	payloads, _ := ctx.Get("payloads")
	userPayload := payloads.(*pkg.Payload)

	transactions, err := t.transactionsRepo.UseGetTransactions(ctx.Request.Context(), userPayload.Id)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(transactions) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No transactions found"})
		return
	}

	// Format ulang tanggal dan waktu
	type TransactionFormattedRes struct {
		MovieTitle string `json:"movie_title"`
		Schedule   string `json:"schedule"` // format yyyy-mm-dd
		Time       string `json:"time"`     // format hh:mm:ss
		CinemaName string `json:"cinema_name"`
		SeatNames  string `json:"seat_names"`
		GrandTotal int    `json:"grand_total"`
	}

	var formatted []TransactionFormattedRes
	for _, trx := range transactions {
		formatted = append(formatted, TransactionFormattedRes{
			MovieTitle: trx.MovieTitle,
			Schedule:   trx.Schedule.Format("2006-01-02"),
			Time:       formatTime(trx.Time),
			CinemaName: trx.CinemaName,
			SeatNames:  trx.SeatNames,
			GrandTotal: trx.GrandTotal,
		})
	}

	ctx.JSON(http.StatusOK, formatted)
}

func formatTime(t string) string {
	parsed, err := time.Parse("15:04:05.000000", t)
	if err != nil {
		return t
	}
	return parsed.Format("15:04:05")
}
