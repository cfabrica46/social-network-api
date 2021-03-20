package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID                 int
	Username, Password string
}

type Post struct {
	Propetary string
	ID        int
	Content   string
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

	http.HandleFunc("/user", db.user)

	http.HandleFunc("/user/posts", db.myPosts)

	http.HandleFunc("/user/friends", db.friends)

	http.HandleFunc("/user/friends/posts", db.friendsPosts)

	fmt.Println("Listening on 8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}
