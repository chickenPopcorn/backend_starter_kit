package main

import (
	"fmt"
	"log"
	"net/http"

	"./config"
	. "./models"
	"./services"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	// "github.com/aws/aws-sdk-go/aws/session"
)

var dbUsers = map[string]User{} // user ID, user struct

func init() {

	// for testing only
	bs, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	dbUsers["jimmy"] = User{"jimmy", bs, "jimmy", "xie"}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s3 := config.NewS3()

	db, err := config.NewDB("root:yuki@tcp(localhost:3306)/userinfo?charset=utf8")
	check(err)
	defer db.Close()

	var dbSessions = map[string]Session{} // session ID, session

	env := &config.Env{db, dbSessions, s3}

	r := mux.NewRouter().StrictSlash(true)

	// Routes consist of a path and a handler function.
	r.HandleFunc("/login", services.Login(env)).Methods("POST")
	r.HandleFunc("/reg", services.Reg(env)).Methods("POST")
	r.HandleFunc("/logout", services.Logout(env)).Methods("GET")
	r.HandleFunc("/upload", services.Upload(env)).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8001", r))
}
