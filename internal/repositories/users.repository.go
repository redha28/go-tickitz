package repositories

import (
	"context"
	"fmt"
	"gotickitz/internal/models"
	"strconv"
	"strings"

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

	query := `INSERT INTO auth (email, "password") VALUES ($1, $2) RETURNING email, "role", id;`
	values := []any{userReq.Email, userReq.Password}
	var result models.UserRes
	if err := auth.QueryRow(c, query, values...).Scan(&result.Email, &result.Role, &result.AuthID); err != nil {
		return models.UserRes{}, models.UserRes{}, err
	}

	// log.Println(result.AuthID)
	queryProfile := `INSERT INTO profile (auth_id, firstname, lastname, phone, point, picture) VALUES ($1, $2, $3, $4, $5, $6);`
	valuesProfile := []any{result.AuthID, "enjoyer", "tickitz", "", 0, "https://static.vecteezy.com/system/resources/thumbnails/020/765/399/small_2x/default-profile-account-unknown-icon-black-silhouette-free-vector.jpg"}
	if _, err := auth.Exec(c, queryProfile, valuesProfile...); err != nil {
		return models.UserRes{}, models.UserRes{}, err
	}

	return result, models.UserRes{}, nil
}

func (u *UsersRepository) UseUpdateProfile(c context.Context, authID string, updateReq models.UpdateProfileReq) error {
	query := `UPDATE profile SET`
	var updates []string
	var values []any

	if updateReq.Firstname != "" {
		updates = append(updates, "firstname = $"+strconv.Itoa(len(values)+1))
		values = append(values, updateReq.Firstname)
	}
	if updateReq.Lastname != "" {
		updates = append(updates, "lastname = $"+strconv.Itoa(len(values)+1))
		values = append(values, updateReq.Lastname)
	}
	if updateReq.Phone != "" {
		updates = append(updates, "phone = $"+strconv.Itoa(len(values)+1))
		values = append(values, updateReq.Phone)
	}
	if updateReq.Picture != "" {
		updates = append(updates, "picture = $"+strconv.Itoa(len(values)+1))
		values = append(values, updateReq.Picture)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query += " " + strings.Join(updates, ", ")
	query += " WHERE auth_id = $" + strconv.Itoa(len(values)+1)
	values = append(values, authID)
	if _, err := u.Exec(c, query, values...); err != nil {
		return err
	}

	return nil
}
