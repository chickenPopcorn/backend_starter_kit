package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"./config"
	. "./models"
	"./services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	svc := s3.New(session.New(), cfg)

	file, err := os.Open("test.jpg")
	if err != nil {
		fmt.Printf("err opening file: %s", err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := "/media/" + file.Name()

	params := &s3.PutObjectInput{
		Bucket:        aws.String("go-backend-s3-jimmy"),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String("public-read"),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}
	fmt.Printf("response %s", awsutil.StringValue(resp))

	db, err := config.NewDB("root:yuki@tcp(localhost:3306)/userinfo?charset=utf8")
	check(err)
	defer db.Close()
	var dbSessions = map[string]Session{} // session ID, session
	env := &config.Env{db, dbSessions}
	r := mux.NewRouter().StrictSlash(true)

	// Routes consist of a path and a handler function.
	r.HandleFunc("/login", services.Login(env)).Methods("POST")
	r.HandleFunc("/reg", services.Reg(env)).Methods("POST")
	r.HandleFunc("/logout", services.Logout(env)).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8001", r))
}
