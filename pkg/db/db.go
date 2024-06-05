package db

import (
	"context"
	"fmt"
	"html/template"
	"log"

	"github.com/DuvanAlbarracin/movies_auth/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Conn *pgxpool.Pool
}

func Init(url string) Handler {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalln("Error creating connection with the database:", err)
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalln("Error while acquiring connection from the database pool:", err)
	}
	defer conn.Release()

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalln("Error pinging the database:", err)
	}

	log.Println("Database conection success!")

	createUserTable(pool)

	return Handler{Conn: pool}
}

func createUserTable(pool *pgxpool.Pool) (err error) {
	_, err = pool.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username VARCHAR(10) UNIQUE NOT NULL,email VARCHAR(30) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL);")
	if err != nil {
		log.Fatalln("Error creating the Users table")
		return
	}

	return nil
}

func FindUserByEmail(pool *pgxpool.Pool, email string) (models.User, error) {
	var user models.User
	err := pool.QueryRow(context.Background(),
		"SELECT * FROM users WHERE email = $1", template.HTMLEscapeString(email)).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	fmt.Println("USER:", user)

	return user, err
}

func CreateUser(pool *pgxpool.Pool, user *models.User) (err error) {
	_, err = pool.Exec(context.Background(),
		"insert into users(username, email, password) values($1, $2, $3)",
		user.Username,
		user.Email,
		user.Password,
	)

	return
}
