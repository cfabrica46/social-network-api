package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type user struct {
	ID                 int
	Username, Password string
}

func main() {

	var election int
	var exit bool

	for exit == false {

		fmt.Println("1.Get User")
		fmt.Println("2.Get All Users")
		fmt.Println("3.Create User")
		fmt.Println("4.Update User")
		fmt.Println("5.Delete User")
		fmt.Println("0.Salir")

		fmt.Scan(&election)

		switch election {

		case 0:

			exit = true

		case 1:

			findUser()

		case 2:

			findUsers()

		case 3:

			u := createUserGet()

			createUserPost(u)

		case 4:

			u := updateUserGet()

			updateUserPost(u)

		case 5:

			deleteUser()

		default:

			fmt.Println("Seleccione una opción válida")

		}

	}
}

func findUsers() {

	res, err := http.Get("http://192.168.1.2:8080/users/all")

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotFound))
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusInternalServerError))
		return
	}

	fmt.Printf("%s\n", body)

}

func findUser() {

	var id string

	fmt.Println("Escribe el ID del usuario que deseas ver")
	fmt.Scan(&id)

	s := fmt.Sprintf("http://192.168.1.2:8080/users/one?id=%s", id)

	res, err := http.Get(s)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotFound))
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusInternalServerError))
		return
	}

	fmt.Printf("%s\n", body)

}

func createUserGet() (u user) {

	var username, password string

	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://192.168.1.2:8080/user/create", nil)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotFound))
		return
	}

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusConflict))
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotAcceptable))
		return
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
		fmt.Println(http.StatusText(http.StatusNotAcceptable))
		return
	}

	dataReader := strings.NewReader(string(data))

	res, err := http.Post("http://192.168.1.2:8080/user/create", "application/json", dataReader)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotFound))
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusInternalServerError))
		return
	}

	fmt.Printf("%s\n", body)

}

func updateUserGet() (u user) {

	var id string

	fmt.Println("Escribe el ID del usuario que deseas editar")
	fmt.Print("> ")
	fmt.Scan(&id)

	s := fmt.Sprintf("http://192.168.1.2:8080/user/update?id=%s", id)

	res, err := http.Get(s)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotFound))
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusInternalServerError))
		return
	}
	err = json.Unmarshal(body, &u)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusInternalServerError))
		return
	}

	return

}

func updateUserPost(u user) {

	var username, password string

	fmt.Println("Introduzca los nuevos valores")
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)

	u.Username = username
	u.Password = password

	data, err := json.Marshal(u)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotAcceptable))
		return
	}

	dataReader := strings.NewReader(string(data))

	res, err := http.Post("http://192.168.1.2:8080/user/update", "application/json", dataReader)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusNotFound))
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(http.StatusText(http.StatusInternalServerError))
		return
	}

	fmt.Printf("%s\n", body)

}

func deleteUser() {

}
