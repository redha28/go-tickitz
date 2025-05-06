package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"gotickitz/internal/models"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MoviesRepository struct {
	*pgxpool.Pool
	*redis.Client
}

func NewMoviesRepository(pg *pgxpool.Pool, rdb *redis.Client) *MoviesRepository {
	return &MoviesRepository{pg, rdb}
}

func (m *MoviesRepository) UseGetMovies(ctx context.Context, name, options, genre string, limit, offset int) ([]models.Movies, models.Meta, error) {
	// Redis key untuk caching
	redisKey := fmt.Sprintf("movies:%s:%d:%d", options, limit, offset)

	if options != "" && name == "" && genre == "" {
		cache, err := m.Get(ctx, redisKey).Result()
		if err == nil {
			var cachedData models.RedistMovies
			if err := json.Unmarshal([]byte(cache), &cachedData); err == nil {
				return cachedData.Movies, cachedData.Meta, nil
			}
		} else if err != redis.Nil {
			log.Println("Redis error:", err)
		}
	}

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

	var conditions []string
	var args []any
	argIndex := 1

	switch options {
	case "upcoming":
		conditions = append(conditions, fmt.Sprintf("m.release IS NOT NULL AND m.release > NOW()"))
	case "popular":
		conditions = append(conditions, fmt.Sprintf("m.rating IS NOT NULL AND m.rating > 6"))
	case "allrelease":
		conditions = append(conditions, fmt.Sprintf("m.release IS NOT NULL AND m.release < NOW()"))
	default:
		conditions = append(conditions, "1=1") // default semua movie
	}

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(m.title) LIKE LOWER($%d)", argIndex))
		args = append(args, "%"+name+"%")
		argIndex++
	}
	if genre != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(g.name) LIKE LOWER($%d)", argIndex))
		args = append(args, "%"+genre+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY m.id
		ORDER BY m.id
		LIMIT $` + fmt.Sprint(argIndex) + ` OFFSET $` + fmt.Sprint(argIndex+1) + `;`

	args = append(args, limit, offset)

	// Query main movies
	rows, err := m.Query(ctx, query, args...)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, models.Meta{}, err
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
			return nil, models.Meta{}, err
		}
		movies = append(movies, movie)
	}

	// Count query
	countQuery := `
		SELECT COUNT(DISTINCT m.id)
		FROM movies m
		JOIN movie_genres mg ON m.id = mg.movie_id
		JOIN genres g ON mg.genre_id = g.id
	`

	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Hanya kirim parameter filtering (tanpa limit/offset)
	argsForCount := args[:len(args)-2]

	var totalRecords int
	err = m.QueryRow(ctx, countQuery, argsForCount...).Scan(&totalRecords)
	if err != nil {
		log.Println("Failed to count total records:", err)
		return nil, models.Meta{}, err
	}

	totalPages := (totalRecords + limit - 1) / limit
	currentPage := (offset / limit) + 1

	meta := models.Meta{
		Page:       currentPage,
		TotalPages: totalPages,
	}

	if options != "" && name == "" && genre == "" {
		cachedData := models.RedistMovies{
			Movies: movies,
			Meta:   meta,
		}
		cachedJSON, err := json.Marshal(cachedData)
		if err == nil {
			err = m.Set(ctx, redisKey, cachedJSON, 24*time.Hour).Err()
			if err != nil {
				log.Println("Failed to save movies to Redis:", err)
			}
		}
	}

	return movies, meta, nil
}

func (m *MoviesRepository) UseGetMovie(ctx context.Context, id int) (models.Movies, error) {
	query := `SELECT m.id, m.title, m.image, m.banner, m.release, m.synopsis, m.duration, m.director, m.rating, STRING_AGG(g.name, ', ') AS genres FROM movies m JOIN movie_genres mg ON m.id = mg.movie_id JOIN genres g ON mg.genre_id = g.id WHERE m.id = $1 GROUP BY m.id;`
	var movie models.Movies
	err := m.QueryRow(ctx, query, id).Scan(&movie.ID, &movie.Title, &movie.Image, &movie.Banner, &movie.Release, &movie.Synopsis, &movie.Duration, &movie.Director, &movie.Rating, &movie.Genre)
	return movie, err
}
