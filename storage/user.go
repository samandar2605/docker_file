package storage

import (
	"github.com/jmoiron/sqlx"
)

type DBManager struct {
	db *sqlx.DB
}

func NewDBManager(db *sqlx.DB) *DBManager {
	return &DBManager{db: db}
}

type User struct {
	Id          int `json:"id"`
	FirstName   string `json:"first_name"` 
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

func (ur *DBManager) Create(u *User) (*User, error) {
	query := `
		INSERT INTO users(
			first_name,
			last_name,
			phone_number
		)values($1,$2,$3)
		RETURNING id
	`

	row := ur.db.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		u.PhoneNumber,
	)

	if err := row.Scan(
		&u.Id,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *DBManager) Get(id int) (*User, error) {
	var user User

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number
		from users
		where id=$1
	`
	row := ur.db.QueryRow(query, id)
	if err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
