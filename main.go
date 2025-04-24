package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := gin.Default()

	// CONNECT TO DB
	dbEnv := []any{}
	dbEnv = append(dbEnv, os.Getenv("DBUSER"))
	dbEnv = append(dbEnv, os.Getenv("DBPASS"))
	dbEnv = append(dbEnv, os.Getenv("DBHOST"))
	dbEnv = append(dbEnv, os.Getenv("DBPORT"))
	dbEnv = append(dbEnv, os.Getenv("DBNAME"))

	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbEnv...)
	dbClient, err := pgxpool.New(context.Background(), dbString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		log.Println("DB connection closed")
		dbClient.Close()
	}()

	// var c *gin.Context

	// ROUTES
	v1 := router.Group("/api/v1")

	users := v1.Group("/users")
	// AUTH

	// type RegisterHandler struct {
	// 	c        *gin.Context
	// 	dbClient *pgxpool.Pool
	// }

	auth := users.Group("/auth")
	{
		auth.GET("", func(ctx *gin.Context) {
			loginHandler(ctx, dbClient)
		})
		auth.POST("", func(ctx *gin.Context) {
			registerHandler(ctx, dbClient)
		})
	}
	// PROFILE
	profile := users.Group("/profile")
	{
		profile.GET("", getProfileHandler)
		profile.PATCH("", updateProfileHandler)
	}
	// MOVIE
	movies := v1.Group("/movies")
	{
		movies.GET("", getMoviesHandler)
		movies.GET("/:id", getMovieDetailHandler)
		movies.GET("/popular", getPopularMoviesHandler)
		movies.GET("/upcoming", getUpcomingMoviesHandler)
	}

	// SHOWING MOVIE
	showings := v1.Group("/showings")
	{
		showings.GET("", getShowingsHandler)
		showings.GET("/:showing_id", getShowingDetailHandler)
	}

	// TRANSACTIONS
	transactions := v1.Group("/transactions")
	{
		transactions.GET("/:user_id", getTransactionsHandler)
		transactions.POST("", createTransactionHandler)
	}

	// ADMIN
	admin := v1.Group("/admin")
	{
		admin.GET("/movies", getMoviesAdminHandler)
		admin.POST("/movies", createMovieHandler)
		admin.PATCH("/movies/:id", updateMovieHandler)
		admin.DELETE("/movies/:id", deleteMovieHandler)
	}

	router.Run("localhost:8080")
}

type Response struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type UserReq struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required"`
	Role     string `json:"role" db:"role"`
}

type UserRes struct {
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}

func registerHandler(c *gin.Context, dbClient *pgxpool.Pool) {
	var userReq UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "Invalid input",
			"error": err.Error(),
		})
		return
	}

	queryCheckMail := `SELECT email FROM auth WHERE email = $1`
	valuesMail := []any{userReq.Email}
	var findUser UserRes
	if err := dbClient.QueryRow(c.Request.Context(), queryCheckMail, valuesMail...).Scan(&findUser.Email); err != nil && err != pgx.ErrNoRows {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan pada server",
		})
		return
	}
	if findUser.Email == userReq.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Email sudah terdaftar",
		})
		return
	}

	query := `INSERT INTO auth (email, "password") VALUES ($1, $2) RETURNING email, "role";`
	values := []any{userReq.Email, userReq.Password}
	var result UserRes
	if err := dbClient.QueryRow(c.Request.Context(), query, values...).Scan(&result.Email, &result.Role); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg":  "success",
		"data": result,
	})
}

func loginHandler(c *gin.Context, dbClient *pgxpool.Pool) {
	var userReq UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "Invalid input",
			"error": err.Error(),
		})
		return
	}
	query := `SELECT email, "role" FROM auth WHERE email = $1 AND "password" = $2`
	values := []any{userReq.Email, userReq.Password}
	var result UserRes
	if err := dbClient.QueryRow(c.Request.Context(), query, values...).Scan(&result.Email, &result.Role); err != nil && err != pgx.ErrNoRows {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Terjadi kesalahan pada server",
		})
		return
	}
	if result.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Email atau password salah",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "login success",
		"data": result,
	})
}
func getProfileHandler(c *gin.Context)        {}
func updateProfileHandler(c *gin.Context)     {}
func getMoviesHandler(c *gin.Context)         {}
func getMovieDetailHandler(c *gin.Context)    {}
func getPopularMoviesHandler(c *gin.Context)  {}
func getUpcomingMoviesHandler(c *gin.Context) {}
func getShowingsHandler(c *gin.Context)       {}
func getShowingDetailHandler(c *gin.Context)  {}
func getTransactionsHandler(c *gin.Context)   {}
func createTransactionHandler(c *gin.Context) {}
func getMoviesAdminHandler(c *gin.Context)    {}
func createMovieHandler(c *gin.Context)       {}
func updateMovieHandler(c *gin.Context)       {}
func deleteMovieHandler(c *gin.Context)       {}
