package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	UserName string
	Password []byte
	First    string
	Last     string
}

var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID

func init() {
	bs, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	dbUsers["jimmy"] = user{"jimmy", bs, "jimmy", "xie"}
}

//login function handler
func login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Error(w, http.StatusText(http.StatusSeeOther), http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	userinfo, ok := dbUsers[username]
	if !ok {
		http.Error(w, "Username and/or password do not match", http.StatusForbidden)
		return
	}

	err := bcrypt.CompareHashAndPassword(userinfo.Password, []byte(password))
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
	dbSessions[c.Value] = username
	fmt.Println("cookie is ", c.Value)

	w.WriteHeader(http.StatusOK)
}

//reg function handler
func reg(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		//TODO change to properate message
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// get form values
	username := r.FormValue("username")
	password := r.FormValue("password")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")

	// username taken?
	if _, ok := dbUsers[username]; ok {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// create session
	sID := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	http.SetCookie(w, c)
	dbSessions[c.Value] = username

	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// store user in dbUsers
	userinfo := user{username, bs, firstname, lastname}
	dbUsers[username] = userinfo
	fmt.Printf("userinfo %v", userinfo)
	fmt.Println("should return some status code and message")

}

//logout function handler
func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		fmt.Println("i'm hrere")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	c, _ := r.Cookie("session")
	delete(dbSessions, c.Value)
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, c)
	w.WriteHeader(http.StatusOK)
}

func getUser(w http.ResponseWriter, r *http.Request) user {
	// get cookie
	c, err := r.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	http.SetCookie(w, c)

	// if the user exists already, get user
	var userinfo user
	if username, ok := dbSessions[c.Value]; ok {
		userinfo = dbUsers[username]
	}
	return userinfo
}

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	username := dbSessions[c.Value]
	_, ok := dbUsers[username]
	fmt.Println("already logged in? ", ok)
	return ok
}

func main() {

	r := mux.NewRouter().StrictSlash(true)

	// Routes consist of a path and a handler function.
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/reg", reg).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
