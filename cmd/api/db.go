package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	ID        string    `json:"ID"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"_"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type Model struct {
	user User
}

func NewModel(database *sql.DB) *Model {
	db = database
	return &Model{
		user: User{},
	}
}

func (u *User) GetInfoByEmail(email string) (*User, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.QueryRowContext(queryCtx, "SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Password,
		&u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("did not find user with this email")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) MatchPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
