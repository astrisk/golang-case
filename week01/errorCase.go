package week01

import (
	"database/sql"
	"errors"
)

type user struct {
	Age  int
	Name string
	Uid  int64
}

func getUser(db *sql.DB, uid int64) (*user, error) {
	user := new(user)
	sqlStmt := `SELECT name,age FROM user WHERE uid = ?`
	rows, err := db.Query(sqlStmt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&user.Name, &user.Age, &user.Uid)
	return user, nil
}
