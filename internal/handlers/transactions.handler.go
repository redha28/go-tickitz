package handler

import (
	"gotickitz/internal/models"
	"gotickitz/internal/repositories"
	"log"
	"net/http"

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

	if err := t.transactionsRepo.UseCreateTransaction(ctx.Request.Context(), &transactionReq); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := t.transactionsRepo.UseAddPoints(ctx.Request.Context(), transactionReq.AuthID, len(transactionReq.SeatId)); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})

}
