package config

import (
	"database/sql"
	"fmt"

	. "../models"
	//mysql libary
	_ "github.com/go-sql-driver/mysql"
)

//Env used for code injection
type Env struct {
	DB       *sql.DB
	Sessions map[string]Session
}

//NewDB opens database
func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
