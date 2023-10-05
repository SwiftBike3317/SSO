package storage

import (
	Types "SSO/types"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage interface {
	GetUser(email, password string) (Types.Account, bool)
	GetServices(userid string) ([]int, bool)
	AddToBlacklist(jwt string) error
	InBlacklist(jwt string) bool
}

type PostgresStore struct {
	db *sqlx.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "host=database port=5432 user=postgres dbname=postgres password=qwerty sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil

}

func (s *PostgresStore) GetUser(email, password string) (Types.Account, bool) {
	statement := `SELECT id, patronymic, first_name, last_name FROM users WHERE email = $1 AND password = $2 LIMIT 1;`
	rows, err := s.db.Query(statement, email, password)
	if err != nil {
		return Types.Account{}, false
	}

	var user Types.Account
	user.Email = email
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Patronymic, &user.FirstName, &user.LastName)

		if err != nil {
			return Types.Account{}, false
		}
	}

	return user, true

}

func (s *PostgresStore) GetServices(userid string) ([]int, bool) {
	statement := `SELECT service_id FROM access WHERE user_id = $1;`
	rows, err := s.db.Query(statement, userid)
	if err != nil {
		return []int{}, false
	}

	var id int
	var ids []int

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {

			return []int{}, false
		}
		ids = append(ids, id)
	}

	return ids, true

}

func (s *PostgresStore) AddToBlacklist(jwt string) error {
	schema := `INSERT INTO jwt_blacklist (token) VALUES ($1)`
	s.db.MustExec(schema, jwt)
	return nil
}

func (s *PostgresStore) InBlacklist(jwt string) bool {
	query := "SELECT id FROM jwt_blacklist WHERE token = $1"
	var id int
	err := s.db.QueryRow(query, jwt).Scan(&id)
	fmt.Println("test1", id)
	if err != nil && err != sql.ErrNoRows {
		return true
	}
	fmt.Println("test2")

	if id != 0 {
		return true
	}
	return false
}
