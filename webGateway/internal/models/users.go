package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Username  string
	Email     string
	HashedPWD []byte
	Created   time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(userName, email, pwd string) error {
	hashedPWD, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (username, email, hashed_password, created)
	VALUES(?, ?, ?, NOW())`

	_, err = m.DB.Exec(stmt, userName, email, hashedPWD)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 {
				if strings.Contains(mySQLError.Message, "users.users_uc_username") {
					return ErrDuplicateUserName
				} else if strings.Contains(mySQLError.Message, "users.users_uc_email") {
					return ErrDuplicateEmail
				}
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, string, error) {
	var id int
	var username string
	var hashedPassword []byte
	stmt := "SELECT id, username, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &username, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, "", ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}
	return id, username, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) Get(id int) (User, error) {
	var user User
	stmt := `SELECT username, email, created FROM users where id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&user.Username, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		} else {
			return User{}, err
		}
	}
	return user, err
}
