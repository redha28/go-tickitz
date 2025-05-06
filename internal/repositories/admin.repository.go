package repositories

import (
	"context"
	"gotickitz/internal/models"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AdminRepository struct {
	*pgxpool.Pool
	*redis.Client
}

func NewAdminRepository(pg *pgxpool.Pool, rdb *redis.Client) *AdminRepository {
	return &AdminRepository{pg, rdb}
}

func (a *AdminRepository) UseCheckAdmin(ctx context.Context, uuid string) (bool, error) {
	query := "SELECT role FROM auth WHERE id = $1"

	var role string
	err := a.QueryRow(ctx, query, uuid).Scan(&role)
	if err != nil {
		return false, err
	}
	if role == "admin" {
		return true, nil
	}
	return false, nil
}

func (a *AdminRepository) UseCreateMovie(ctx context.Context, movieReq models.AdminCreateMovieReq, AdminId string) error {
	query := "INSERT INTO movies (title, image, banner, release, synopsis, duration, director, rating, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;"

	var movieID int
	err := a.QueryRow(ctx, query, movieReq.Title, movieReq.Image, movieReq.Banner, movieReq.Release, movieReq.Synopsis, movieReq.Duration, movieReq.Director, movieReq.Rating, AdminId).Scan(&movieID)
	if err != nil {
		return err
	}

	for _, genreID := range movieReq.Genres {
		query = "INSERT INTO movie_genres (movie_id, genre_id) VALUES ($1, $2);"
		_, err = a.Exec(ctx, query, movieID, genreID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AdminRepository) UseEditMovie(ctx context.Context, movieReq models.AdminEditMovieReq, id int) (string, error) {
	query := `UPDATE movies SET`
	var updates []string
	var values []any

	if movieReq.Title != "" {
		updates = append(updates, "title = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Title)
	}
	if movieReq.Image != "" {
		updates = append(updates, "image = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Image)
	}
	if movieReq.Banner != "" {
		updates = append(updates, "banner = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Banner)
	}
	if movieReq.Release != "" {
		updates = append(updates, "release = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Release)
	}
	if movieReq.Synopsis != "" {
		updates = append(updates, "synopsis = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Synopsis)
	}
	if movieReq.Duration != 0 {
		updates = append(updates, "duration = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Duration)
	}
	if movieReq.Director != "" {
		updates = append(updates, "director = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Director)
	}
	if movieReq.Rating != 0 {
		updates = append(updates, "rating = $"+strconv.Itoa(len(values)+1))
		values = append(values, movieReq.Rating)
	}

	if len(updates) == 0 {
		return "no fields to update", nil
	}

	query += " " + strings.Join(updates, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(values)+1)
	values = append(values, id)
	if _, err := a.Exec(ctx, query, values...); err != nil {
		return "", err
	}

	a.FlushDB(ctx)

	return "", nil
}

func (a *AdminRepository) UseEditMovieGenre(ctx context.Context, movieID int, genreIDs []int) error {

	query := `INSERT INTO movie_genres (movie_id, genre_id) VALUES`
	for i, genreID := range genreIDs {
		if i == len(genreIDs)-1 {
			query += " (" + strconv.Itoa(movieID) + ", " + strconv.Itoa(genreID) + ")"
		} else {
			query += " (" + strconv.Itoa(movieID) + ", " + strconv.Itoa(genreID) + "),"
		}
	}
	query += ` on conflict do nothing;`
	log.Println(query)
	if _, err := a.Exec(ctx, query); err != nil {
		return err
	}

	queryDelete := `DELETE from movie_genres WHERE movie_id = $1 AND genre_id NOT IN ( `

	var values []any

	values = append(values, movieID)
	for i, v := range genreIDs {
		values = append(values, v)
		queryDelete += "$" + strconv.Itoa(i+2)
		if i != len(genreIDs)-1 {
			queryDelete += `, `
		}
	}
	queryDelete += `)`

	log.Println(queryDelete)

	log.Println(values...)
	if _, err := a.Exec(ctx, queryDelete, values...); err != nil {
		return err
	}

	a.FlushDB(ctx)

	return nil
}

func (a *AdminRepository) UseDeleteMovie(ctx context.Context, movieReq models.AdminDeleteMovieReq) error {
	query := "DELETE FROM movies WHERE id = $1;"
	if _, err := a.Exec(ctx, query, movieReq.MoviesID); err != nil {
		return err
	}

	query = "DELETE FROM movie_genres WHERE movie_id = $1;"
	if _, err := a.Exec(ctx, query, movieReq.MoviesID); err != nil {
		return err
	}

	a.FlushDB(ctx)

	return nil
}
