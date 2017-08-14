package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./models"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type session struct {
	username     string
	lastActivity time.Time
}

var dbUsers = map[string]models.User{} // user ID, user struct
var dbSessions = map[string]session{}  // session ID, session
var dbSessionsCleaned time.Time

type Env struct {
	db models.Datastore
}

const sessionLength int = 30

func init() {
	dbSessionsCleaned = time.Now()

	// for testing only
	bs, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	dbUsers["jimmy"] = models.User{"jimmy", bs, "jimmy", "xie"}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//login function handler
func (env *Env) login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Error(w, http.StatusText(http.StatusSeeOther), http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	userinfo, err := env.db.GetUserInfo(username)

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
		Name:   "session",
		Value:  sID.String(),
		MaxAge: sessionLength,
	}
	http.SetCookie(w, c)
	dbSessions[c.Value] = session{username, time.Now()}
	fmt.Println("cookie is ", c.Value)

	w.WriteHeader(http.StatusOK)
}

//reg function handler
func (env *Env) reg(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
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
		Name:   "session",
		Value:  sID.String(),
		MaxAge: sessionLength,
	}
	http.SetCookie(w, c)
	dbSessions[c.Value] = session{username, time.Now()}

	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// store user in dbUsers
	userinfo := models.User{username, bs, firstname, lastname}
	dbUsers[username] = userinfo
	fmt.Printf("userinfo %v", userinfo)
	fmt.Println("should return some status code and message")

}

//logout function handler
func (env *Env) logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
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

	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	w.WriteHeader(http.StatusOK)
}

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) >
			(time.Second * time.Duration(sessionLength)) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	session, ok := dbSessions[c.Value]
	if ok {
		session.lastActivity = time.Now()
		dbSessions[c.Value] = session
	}
	_, ok = dbUsers[session.username]
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	return ok
}

func main() {
	db, err := models.NewDB("root:yuki@tcp(localhost:3306)/userinfo?charset=utf8")
	check(err)
	defer db.Close()
	env := &Env{db: db}
	r := mux.NewRouter().StrictSlash(true)

	// Routes consist of a path and a handler function.
	r.HandleFunc("/login", env.login).Methods("POST")
	r.HandleFunc("/reg", env.reg).Methods("POST")
	r.HandleFunc("/logout", env.logout).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8001", r))
}
