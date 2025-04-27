package repositories

import (
	"context"
	"gotickitz/internal/models"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionsRepository struct {
	*pgxpool.Pool
}

func NewTransactionsRepository(pg *pgxpool.Pool) *TransactionsRepository {
	return &TransactionsRepository{pg}
}

func (t *TransactionsRepository) UseCreateTransaction(ctx context.Context, datas *models.TransactionReq) error {
	query := "INSERT INTO transactions (auth_id, payment_id, status, grand_total) values ($1, $2, $3, $4) RETURNING id;"

	var transactions_id int
	if err := t.QueryRow(ctx, query, datas.AuthID, datas.PaymentId, datas.Status, datas.GrandTotal).Scan(&transactions_id); err != nil {
		return err
	}
	log.Println(transactions_id)
	querySeat := "INSERT INTO showing_seats (showing_id, seat_id, transactions_id) values ($1, $2, $3);"
	for _, seat_id := range datas.SeatId {
		if _, err := t.Exec(ctx, querySeat, datas.ShowingId, seat_id, transactions_id); err != nil {
			return err
		}
	}
	return nil
}

func (s *TransactionsRepository) UseGetSeatsSold(ctx context.Context, showing_id int) ([]models.SeatRes, error) {
	query := `SELECT 
		s.id AS seat_id,
		s.name AS seat_name
		FROM showing_seats ss
		JOIN seats s ON ss.seat_id = s.id
		WHERE ss.showing_id = $1;`

	var seats []models.SeatRes
	rows, err := s.Query(ctx, query, showing_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var seat models.SeatRes
		if err := rows.Scan(&seat.ID, &seat.Name); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, err
}

func (s *TransactionsRepository) UseAddPoints(ctx context.Context, auth_id string, points int) error {
	query := "UPDATE profile SET point = point + $1 WHERE auth_id = $2;"
	totalPoint := points * 10
	if _, err := s.Exec(ctx, query, totalPoint, auth_id); err != nil {
		return err
	}
	return nil
}
