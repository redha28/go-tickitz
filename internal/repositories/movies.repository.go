package repositories

import (
	"context"
	"fmt"
	"gotickitz/internal/models"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MoviesRepository struct {
	*pgxpool.Pool
}

func NewMoviesRepository(pg *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{pg}
}

func (m *MoviesRepository) UseGetMovies(ctx context.Context, name, options, genre string, limit, offset int) ([]models.Movies, map[string]int, error) {
	query := `
		SELECT 
		m.id,
		m.title,
		m.image,
		m.banner,
		m.release,
		m.synopsis,
		m.duration,
		m.director,
		m.rating,
		STRING_AGG(g.name, ', ') AS genres
	FROM movie_genres mg
	JOIN movies m ON mg.movie_id = m.id
	JOIN genres g ON mg.genre_id = g.id
	`
	// log.Println(options)
	if options == "upcoming" {
		query += `
			WHERE m.release IS NOT NULL AND m.release > NOW()
		`
	}

	if options == "popular" {
		query += `
			WHERE m.rating IS NOT NULL AND m.rating > 6
		`
	}

	if options == "allrelease" {
		query += `
			WHERE m.release IS NOT NULL AND m.release < NOW()
		`
	}

	if name != "" {
		query += fmt.Sprintf(" AND LOWER(m.title) LIKE '%%%s%%'", strings.ToLower(name))
	}

	if genre != "" {
		query += fmt.Sprintf(" AND LOWER(g.name) LIKE '%%%s%%'", strings.ToLower(genre))
	}

	query += `
		GROUP BY m.id
		ORDER BY m.id
		LIMIT $1 OFFSET $2;
	`

	rows, err := m.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, nil, err
	}
	defer rows.Close()

	var movies []models.Movies
	for rows.Next() {
		var movie models.Movies
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Image,
			&movie.Banner,
			&movie.Release,
			&movie.Synopsis,
			&movie.Duration,
			&movie.Director,
			&movie.Rating,
			&movie.Genre,
		)
		if err != nil {
			log.Println("Failed to scan row:", err)
			return nil, nil, err
		}
		movies = append(movies, movie)
	}

	countQuery := `
	SELECT COUNT(DISTINCT m.id)
	FROM movies m
	JOIN movie_genres mg ON m.id = mg.movie_id
	JOIN genres g ON mg.genre_id = g.id
`
	var conditions []string

	if options == "upcoming" {
		conditions = append(conditions, "m.release IS NOT NULL AND m.release > NOW()")
	}

	if options == "popular" {
		conditions = append(conditions, "m.rating IS NOT NULL AND m.rating > 6")
	}

	if options == "allrelease" {
		conditions = append(conditions, "m.release IS NOT NULL AND m.release < NOW()")
	}

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(m.title) LIKE '%%%s%%'", strings.ToLower(name)))
	}

	if genre != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(g.name) LIKE '%%%s%%'", strings.ToLower(genre)))
	}

	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var totalRecords int
	err = m.QueryRow(ctx, countQuery).Scan(&totalRecords)
	if err != nil {
		log.Println("Failed to count total records:", err)
		return nil, nil, err
	}

	totalPages := (totalRecords + limit - 1) / limit
	currentPage := (offset / limit) + 1

	meta := map[string]int{
		"page":        currentPage,
		"total_pages": totalPages,
	}

	return movies, meta, nil
}

func (m *MoviesRepository) UseGetMovie(ctx context.Context, id int) (models.Movies, error) {
	query := `SELECT m.id, m.title, m.image, m.banner, m.release, m.synopsis, m.duration, m.director, m.rating, STRING_AGG(g.name, ', ') AS genres FROM movies m JOIN movie_genres mg ON m.id = mg.movie_id JOIN genres g ON mg.genre_id = g.id WHERE m.id = $1 GROUP BY m.id;`
	var movie models.Movies
	err := m.QueryRow(ctx, query, id).Scan(&movie.ID, &movie.Title, &movie.Image, &movie.Banner, &movie.Release, &movie.Synopsis, &movie.Duration, &movie.Director, &movie.Rating, &movie.Genre)
	return movie, err
}
