package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

const (
	login    = "login"
	register = "register"
	userName = "Username"
	passWord = "Password"
)

type page struct {
	Title   string
	Options []string
	User    user
	Posts   []post
	Err     string
}

type user struct {
	ID                 int
	Username, Password string
}

type post struct {
	Propetary string
	ID        int
	Contet    string
	Date      string
}

type dataBases struct {
	d *sql.DB
}

var db dataBases

func main() {

	log.SetFlags(log.Lshortfile)

	databases, err := open()

	if err != nil {
		log.Fatal(err)
	}

	db = dataBases{
		d: databases,
	}

	http.HandleFunc("/user/profile", db.profile)
	http.HandleFunc("/user/create", db.createUser)
	http.HandleFunc("/user/delete", db.deleteUser)
	http.HandleFunc("/user/post/one", db.getMyPosts)

	fmt.Println("Listening on 8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}
