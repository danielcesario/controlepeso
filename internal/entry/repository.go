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

func (repository *PGRepository) FindById(id int) (*Entry, error) {
	var entry Entry
	err := repository.DB.QueryRow("SELECT id, user_id, weight, date FROM t_entry WHERE id = $1", id).Scan(&entry.ID, &entry.UserId, &entry.Weight, &entry.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

func (repository *PGRepository) DeleteById(id int) error {
	_, err := repository.DB.Exec("DELETE FROM t_entry WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PGRepository) Update(entry Entry) (*Entry, error) {
	sqlStatement := `
		UPDATE t_entry SET
			weight = $2, date = $3
		WHERE id = $1
		RETURNING weight, date
	`

	err := repository.DB.QueryRow(sqlStatement, entry.ID, entry.Weight, entry.Date).Scan(&entry.Weight, &entry.Date)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// create table t_entry (id serial PRIMARY KEY, user_id int not null, weight numeric(5,2) not null, date TIMESTAMP not null)
