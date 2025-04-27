package repositories

import (
	"context"
	"gotickitz/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShowingsRepository struct {
	*pgxpool.Pool
}

func NewShowingsRepository(pg *pgxpool.Pool) *ShowingsRepository {
	return &ShowingsRepository{pg}
}

func (s *ShowingsRepository) UseGetShowings(ctx context.Context, movie_id int) ([]models.Showings, error) {
	query := `SELECT 
    s.id AS showing_id,
    s.schedule AS date,
    s.time AS time,
    c.name AS cinema_name,
    c.picture AS cinema_image,
    ci.name AS city_name,
    s.price AS price,
    m.title AS movie_title
		FROM showings s
		JOIN cinemas c ON s.cinema_id = c.id
		JOIN cinema_cities cc ON c.id = cc.cinema_id
		JOIN city ci ON cc.city_id = ci.id
		JOIN movies m ON s.movie_id = m.id
		WHERE m.id = $1
		ORDER BY s.schedule, s.time;`
	var showings []models.Showings
	rows, err := s.Query(ctx, query, movie_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var showing models.Showings
		if err := rows.Scan(&showing.ID, &showing.Date, &showing.Time, &showing.CinemaName, &showing.CinemaImage, &showing.CityName, &showing.Price, &showing.MovieTitle); err != nil {
			return nil, err
		}
		showings = append(showings, showing)
	}
	return showings, err
}

func (s *ShowingsRepository) UseGetSeats(ctx context.Context, showing_id int) ([]models.SeatRes, error) {
	query := `SELECT 
    s.id AS seat_id,
    s.name AS seat_name,
    CASE 
        WHEN ss.id IS NULL THEN TRUE
        ELSE FALSE
    END AS is_available
		FROM seats s
		LEFT JOIN showing_seats ss ON s.id = ss.seat_id AND ss.showing_id = $1
		ORDER BY s.id;`

	var seats []models.SeatRes
	rows, err := s.Query(ctx, query, showing_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var seat models.SeatRes
		if err := rows.Scan(&seat.ID, &seat.Name, &seat.IsAvailable); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, err
}
