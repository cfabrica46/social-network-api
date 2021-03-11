package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type user struct {
	ID                 int
	Username, Password string
}

func main() {

	u := createUserGet()

	createUserPost(u)

}

func createUserGet() (u user) {

	var username, password string

	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://192.168.1.2:8080/", nil)

	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", body)

	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)

	u = user{
		Username: username,
		Password: password,
	}

	return
}

func createUserPost(u user) {

	data, err := json.Marshal(u)

	if err != nil {
		log.Fatal(err)
	}

	dataReader := strings.NewReader(string(data))

	res, err := http.Post("http://192.168.1.2:8080/", "application/json", dataReader)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Printf("%s", body)

}
