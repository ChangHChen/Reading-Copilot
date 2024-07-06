package models

import (
	"database/sql"
	"errors"
	"fmt"
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
				if strings.Contains(mySQLError.Message, "users.username") {
					return ErrDuplicateUserName
				} else if strings.Contains(mySQLError.Message, "users.email") {
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

func (m *UserModel) UpdatePWD(id int, curPWD, newPWD string) error {
	var hashedPassword []byte
	stmt := `SELECT hashed_password FROM users WHERE id=?`

	err := m.DB.QueryRow(stmt, id).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(curPWD))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPWD), 14)
	if err != nil {
		return err
	}
	stmt = `UPDATE users SET hashed_password=? WHERE id=?`
	_, err = m.DB.Exec(stmt, newHashedPassword, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetReadingProgress(userID, bookID int) (int, error) {
	stmt := `SELECT page FROM reading_progress WHERE user_id=? AND book_id=?`
	var pageNum int
	err := m.DB.QueryRow(stmt, userID, bookID).Scan(&pageNum)
	if err != nil {
		if err == sql.ErrNoRows {
			pageNum = 1
			updateErr := m.UpdateReadingProgress(userID, bookID, pageNum)
			if updateErr != nil {
				return 0, fmt.Errorf("failed to initialize reading progress for user %d on book %d: %w", userID, bookID, updateErr)
			}
		} else {
			return 0, err
		}
	}
	return pageNum, nil
}

func (m *UserModel) UpdateReadingProgress(userID, bookID, pageNum int) error {
	stmt := `INSERT INTO reading_progress (user_id, book_id, page)
	VALUES (?, ?, ?)
	ON DUPLICATE KEY UPDATE page = VALUES(page)`

	_, err := m.DB.Exec(stmt, userID, bookID, pageNum)
	if err != nil {
		return err
	}
	return nil
}
