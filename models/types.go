package models

import (
	"database/sql"
	"net/http"
)

//Session in server
type Session struct {
	Username string
}

//Datastore interface for DB
type Datastore interface {
	GetUserInfo(string) (*User, error)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Reg(w http.ResponseWriter, r *http.Request)
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
