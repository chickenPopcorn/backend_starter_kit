package services

import (
	"database/sql"
	"fmt"
	"net/http"

	"../config"
	. "../models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//Env used for code injection
type Env struct {
	Db       *sql.DB
	Sessions map[string]Session
}

//regUserInfo fetch user info from database
func regUserInfo(userinfo User, db *sql.DB) error {
	// insert
	stmt, err := db.Prepare("INSERT userlogin SET username=?,password=?,firstname=?,lastname=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(userinfo.Username, userinfo.Password, userinfo.First, userinfo.Last)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Println("inserted id is ", id)

	return nil
}

//getUserInfo fetch user info from database
func getUserInfo(username string, db *sql.DB) (*User, error) {
	rows, err := db.Query("SELECT * FROM userlogin WHERE username='" + username + "'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userInfo := new(User)
	for rows.Next() {
		err := rows.Scan(&userInfo.Username, &userInfo.Password, &userInfo.First, &userInfo.Last)
		if err != nil {
			fmt.Println(46)
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println(46)
		return nil, err
	}
	return userInfo, nil
}

//Login function handler
func Login(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,
		r *http.Request) {
		if alreadyLoggedIn(w, r, env.Sessions) {
			http.Error(w, http.StatusText(http.StatusSeeOther), http.StatusSeeOther)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		userinfo, err := getUserInfo(username, env.DB)

		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		err = bcrypt.CompareHashAndPassword(userinfo.Password, []byte(password))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		fmt.Println("cookie is ", c.Value)

		w.WriteHeader(http.StatusOK)
	}
}

//Logout function handler
func Logout(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,
		r *http.Request) {
		if !alreadyLoggedIn(w, r, env.Sessions) {
			fmt.Println("i'm hrere")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		c, _ := r.Cookie("session")
		c = &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, c)

		w.WriteHeader(http.StatusOK)
	}
}

//Reg function handler
func Reg(env *config.Env) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,
		r *http.Request) {
		if alreadyLoggedIn(w, r, env.Sessions) {
			//TODO change to properate message
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// get form values
		username := r.FormValue("username")
		password := r.FormValue("password")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		fmt.Printf("%v, %v, %v, %v", username, password, firstname, lastname)

		// username taken?
		if _, err := getUserInfo(username, env.DB); err != nil {
			fmt.Println(151)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		// store user in dbUsers
		userinfo := User{username, bs, firstname, lastname}
		// username taken?
		if err := regUserInfo(userinfo, env.DB); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		fmt.Printf("userinfo %v", userinfo)
		fmt.Println("should return some status code and message")

		// create session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		env.Sessions[c.Value] = Session{username}
	}
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request, dbSessions map[string]Session) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	session, ok := dbSessions[c.Value]
	if ok {
		dbSessions[c.Value] = session
	}

	// c.MaxAge = sessionLength
	// http.SetCookie(w, c)
	return ok
}
