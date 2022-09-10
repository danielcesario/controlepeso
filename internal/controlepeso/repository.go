package controlepeso

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PGRepository struct {
	DB *sql.DB
}

func NewPGRepository(conn *sql.DB) *PGRepository {
	return &PGRepository{
		DB: conn,
	}
}

func InitializeDB(user, password, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func (repository *PGRepository) Save(entry Entry) (*Entry, error) {
	err := repository.DB.QueryRow(
		"INSERT INTO t_entry(user_id, weight, date) VALUES ($1, $2, $3) RETURNING id",
		entry.UserId, entry.Weight, entry.Date).Scan(&entry.ID)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}
