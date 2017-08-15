package config

import (
	"database/sql"
	"fmt"

	. "../models"

	"github.com/aws/aws-sdk-go/service/s3"
	//mysql libary
	_ "github.com/go-sql-driver/mysql"
)

//Env used for code injection
type Env struct {
	DB       *sql.DB
	Sessions map[string]Session
	S3       *s3.S3
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
