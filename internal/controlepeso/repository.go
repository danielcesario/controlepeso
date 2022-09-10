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

func (repository *PGRepository) Save(entry Entry) error {
	fmt.Println("Save")

	_, err := repository.DB.Exec("INSERT INTO t_entry(user_id, weight, date) VALUES ($1, $2, $3)",
		entry.UserId, entry.Weight, entry.Date)

	if err != nil {
		return err
	}
	return nil
}
