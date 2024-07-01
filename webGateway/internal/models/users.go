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
	UserName  string
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
