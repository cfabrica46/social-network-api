package main

import (
	"fmt"
	"log"
	"net/http"
)

type user struct {
	ID                 int
	Username, Password string
}

var users = []user{
	{
		ID:       1,
		Username: "cfabrica46",
		Password: "01234",
	},
	{
		ID:       2,
		Username: "arthuronavah",
		Password: "456456",
	},
}

func main() {

	http.HandleFunc("/users/one", findUser)

	http.HandleFunc("/users/all", findUsers)

	http.HandleFunc("/user/create", createUser)

	http.HandleFunc("/user/delete", deleteUser)

	http.HandleFunc("/user/update", updateUser)

	fmt.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}
