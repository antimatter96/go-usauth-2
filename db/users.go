package db

import (
	"database/sql"
	"fmt"

	// M
	_ "github.com/go-sql-driver/mysql"
)

// UserAuthService is the
type UserAuthService interface {
	Init() error

	CheckUser(string) (bool, error)
	GetPasswordHash(string) (bool, string, string, error)
	AddUser(string, string) error
}

type users struct {
	addUser       *sql.Stmt
	getPassword   *sql.Stmt
	checkUserStmt *sql.Stmt
	db            *sql.DB
}

// Init creates all the prepared statements
func (u *users) Init() error {
	var prepareStatementError error

	u.addUser, prepareStatementError = u.db.Prepare("insert into `users` (`email`,`password`) values (?,?)")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.getPassword, prepareStatementError = u.db.Prepare("select `id`,`password`  from `users` where `email` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}

	u.checkUserStmt, prepareStatementError = u.db.Prepare("select `id` from `users` where `email` = ?")
	if prepareStatementError != nil {
		return prepareStatementError
	}
	return nil
}

func (u *users) CheckUser(email string) (bool, error) {
	var userID int
	err := u.checkUserStmt.QueryRow(email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		fmt.Println(err)
		return false, err
	}
	if userID == 0 {
		return false, nil
	}
	return true, nil
}

func (u *users) AddUser(email, password string) error {
	_, err := u.addUser.Exec(email, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *users) GetPasswordHash(email string) (bool, string, string, error) {
	var password string
	var userID string
	err := u.getPassword.QueryRow(email).Scan(&userID, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", "", nil
		}
		return false, "", "", err
	}
	return true, userID, password, nil
}
