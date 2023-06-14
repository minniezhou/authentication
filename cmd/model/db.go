package model

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

type UserInterface interface {
	GetInfoByEmail(email string) (*User, error)
	MatchPassword(password string) bool
}

type Model struct {
	user *User
}

func NewModel(database *sql.DB) *Model {
	db = database
	return &Model{
		user: &User{},
	}
}

func (m *Model) GetInfoByEmail(email string) (*User, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.QueryRowContext(queryCtx, "SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&m.user.ID, &m.user.Email, &m.user.FirstName, &m.user.LastName, &m.user.Password,
		&m.user.Active, &m.user.CreatedAt, &m.user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("did not find user with this email")
	}
	if err != nil {
		return nil, err
	}
	return m.user, nil
}

func (m *Model) MatchPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.user.Password), []byte(password))
	return err == nil
}
