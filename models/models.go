package models

import (
	"database/sql"
)

//Session in server
type Session struct {
	Username string
}

//DB user defined wrapper
type DB struct {
	*sql.DB
}

//User Data Model
type User struct {
	Username string
	Password []byte
	First    string
	Last     string
}
