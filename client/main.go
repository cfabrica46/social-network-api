package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type page struct {
	Title   string
	Options []string
	User    user
	Err     string
}

type user struct {
	ID                 int
	Username, Password string
}

func main() {

	log.SetFlags(log.Llongfile)

	var election int
	var exit bool

	for !exit {

		fmt.Println()
		fmt.Println("Bienvenido")
		fmt.Println("¿Qué deseas hacer?")
		fmt.Println()

		fmt.Println("1.Iniciar Sesión")
		fmt.Println("2.Registrarse")
		fmt.Println("0.Salir")
		fmt.Println()

		fmt.Print(">")

		fmt.Scan(&election)

		fmt.Println()

		switch election {

		case 0:

			exit = true

		case 1:

			u, err := login("http://localhost:8080/user/profile")

			if err != nil {
				log.Fatal(err)
			}

			profileGET(&u)

			for !exit {

				loopIntoProfile(u, &exit)

			}

		case 2:

			u, err := login("http://localhost:8080/user/create")

			if err != nil {
				log.Fatal(err)
			}

			createUserPOST(&u)

			profileGET(&u)

			for !exit {

				loopIntoProfile(u, &exit)

			}

		default:

			fmt.Println("Seleccione una opción válida")

		}

	}
}

func loopIntoProfile(u user, exit *bool) {

	var election int

	fmt.Println("¿Qué deseas hacer?")
	fmt.Println()

	fmt.Println("1.Eliminar Cuenta")
	fmt.Println("0.Salir")
	fmt.Println()

	fmt.Print(">")

	fmt.Scan(&election)

	fmt.Println()

	switch election {
	case 0:
		*exit = true
	case 1:
		deleteUser(u)
	}

}

func login(stringURL string) (user user, err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("CONNECT", stringURL, nil)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&p)

	if err != nil {
		return
	}

	if p.Err != "" {

		fmt.Printf("\nERROR: %s\n", p.Err)
		return

	}

	fmt.Printf("\n%s\n", p.Title)
	fmt.Printf("%s: ", p.Options[0])
	fmt.Scan(&user.Username)
	fmt.Printf("%s: ", p.Options[1])
	fmt.Scan(&user.Password)

	return
}

func profileGET(user *user) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(*user)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("GET", "http://localhost:8080/user/profile", buf)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&p)

	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}

	if p.Err != "" {

		fmt.Printf("\nERROR: %s\n", p.Err)
		return

	}

	*user = p.User

	fmt.Printf("\n%s\n", p.Title)

	fmt.Printf("Bienvenido %s tu ID es: %d\n", user.Username, user.ID)

	return
}

func createUserPOST(user *user) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(*user)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("POST", "http://localhost:8080/user/create", buf)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&p)

	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}

	if p.Err != "" {

		fmt.Printf("\nERROR: %s\n", p.Err)
		return

	}

	*user = p.User

	fmt.Printf("\n%s\n", p.Title)

	fmt.Printf("Se creo usuario: %s con contraseña: %s y ID: %d\n", user.Username, user.Password, user.ID)

	return
}

func deleteUser(u user) {

}
