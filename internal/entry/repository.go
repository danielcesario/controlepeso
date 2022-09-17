package entry

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

func InitializeDB(user, password, host, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, host, dbname)
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

func (repository *PGRepository) ListAll(start, count int) ([]Entry, error) {
	rows, err := repository.DB.Query(
		"SELECT id, user_id, weight, date FROM t_entry LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	entries := []Entry{}
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.UserId, &entry.Weight, &entry.Date); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// create table t_entry (id serial PRIMARY KEY, user_id int not null, weight numeric(5,2) not null, date TIMESTAMP not null)