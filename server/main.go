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

	http.HandleFunc("/user/profile", db.profile)
	http.HandleFunc("/user/create", db.createUser)
	http.HandleFunc("/user/delete", db.deleteUser)
	http.HandleFunc("/user/post/one", db.getMyPosts)
	http.HandleFunc("/user/post/all", db.getAllFriendsPosts)
	http.HandleFunc("/user/post/add", db.addPost)
	http.HandleFunc("/user/post/delete", db.deletePost)
	http.HandleFunc("/user/friends/show", db.showFriends)
	http.HandleFunc("/user/friends/add", db.addFriend)
	http.HandleFunc("/user/friends/delete", db.deleteFriend)

	fmt.Println("Listening on 8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}
