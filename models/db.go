package models

import (
	"database/sql"
	"fmt"
	//mysql libary
	_ "github.com/go-sql-driver/mysql"
)

//Datastore interface for DB
type Datastore interface {
	GetUserInfo(string) (*User, error)
}

//DB user defined wrapper
type DB struct {
	*sql.DB
}

//NewDB opens database
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//GetUserInfo fetch user info from database
func (db *DB) GetUserInfo(username string) (*User, error) {
	rows, err := db.Query("SELECT * FROM userlogin WHERE username=" + username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userInfo := new(User)
	for rows.Next() {
		err := rows.Scan(&userInfo.UserName, &userInfo.Password, &userInfo.First, &userInfo.Last)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return userInfo, nil
}
