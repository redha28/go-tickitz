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

func (t *TransactionsRepository) UseCreateTransaction(ctx context.Context, datas *models.TransactionReq, UUID string) error {
	query := "INSERT INTO transactions (auth_id, payment_id, status, grand_total) values ($1, $2, $3, $4) RETURNING id;"

	var transactions_id int
	if err := t.QueryRow(ctx, query, UUID, datas.PaymentId, datas.Status, datas.GrandTotal).Scan(&transactions_id); err != nil {
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

func (s *TransactionsRepository) UseGetTransactions(ctx context.Context, authID string) ([]models.TransactionRes, error) {
	query := `
	WITH seat_info AS (
		SELECT 
			ss.transactions_id,
			STRING_AGG(st.name, ', ') AS seat_names
		FROM showing_seats ss
		JOIN seats st ON ss.seat_id = st.id
		GROUP BY ss.transactions_id
	)
	SELECT 
		m.title AS movie_title,
		s.schedule AS date,
		s.time AS time,
		c.name AS cinema_name,
		si.seat_names,
		t.grand_total
	FROM transactions t
	JOIN showing_seats ss ON t.id = ss.transactions_id
	JOIN showings s ON ss.showing_id = s.id
	JOIN movies m ON s.movie_id = m.id
	JOIN cinemas c ON s.cinema_id = c.id
	JOIN seat_info si ON si.transactions_id = t.id
	WHERE t.auth_id = $1
	GROUP BY 
		m.title, s.schedule, s.time, c.name, si.seat_names, t.grand_total
	ORDER BY s.schedule DESC, s.time DESC;
	`

	rows, err := s.Query(ctx, query, authID)
	if err != nil {
		log.Println("QueryContext error:", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.TransactionRes

	for rows.Next() {
		var trx models.TransactionRes
		err := rows.Scan(
			&trx.MovieTitle,
			&trx.Schedule,
			&trx.Time,
			&trx.CinemaName,
			&trx.SeatNames,
			&trx.GrandTotal,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			return nil, err
		}
		transactions = append(transactions, trx)
	}

	return transactions, nil
}
