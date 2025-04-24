package repositories

import (
	"context"
	"gotickitz/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	*pgxpool.Pool
}

func NewUsersRepository(pg *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{pg}
}

func (auth *UsersRepository) UseLogin(c context.Context, email, password string) (models.UserRes, error) {
	query := `SELECT email, "role" FROM auth WHERE email = $1 AND "password" = $2`
	values := []any{email, password}
	var result models.UserRes
	if err := auth.QueryRow(c, query, values...).Scan(&result.Email, &result.Role); err != nil && err != pgx.ErrNoRows {
		return models.UserRes{}, err
	}
	return result, nil
}

func (auth *UsersRepository) UseRegister(c context.Context, userReq models.UserReq) (models.UserRes, models.UserRes, error) {
	queryCheckMail := `SELECT email FROM auth WHERE email = $1`
	valuesMail := []any{userReq.Email}
	var findUser models.UserRes
	if err := auth.QueryRow(c, queryCheckMail, valuesMail...).Scan(&findUser.Email); err != nil && err != pgx.ErrNoRows {
		return models.UserRes{}, models.UserRes{}, err
	}
	if findUser.Email == userReq.Email {
		return models.UserRes{}, findUser, nil
	}

	query := `INSERT INTO auth (email, "password") VALUES ($1, $2) RETURNING email, "role";`
	values := []any{userReq.Email, userReq.Password}
	var result models.UserRes
	if err := auth.QueryRow(c, query, values...).Scan(&result.Email, &result.Role); err != nil {
		return models.UserRes{}, models.UserRes{}, err
	}
	return result, models.UserRes{}, nil
}
