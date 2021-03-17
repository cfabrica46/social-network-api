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

		fmt.Println("1.Ir a tu Perfil")
		fmt.Println("0.Salir")
		fmt.Println()

		fmt.Print(">")

		fmt.Scan(&election)

		fmt.Println()

		switch election {

		case 0:

			exit = true

		case 1:

			userBeta, err := profileCONNECT()

			if err != nil {
				log.Fatal(err)
			}

			profileGET(userBeta)

		default:

			fmt.Println("Seleccione una opción válida")

		}

	}
}

func profileCONNECT() (userBeta user, err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("CONNECT", "http://localhost:8080/user/profile", nil)

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
	fmt.Scan(&userBeta.Username)
	fmt.Printf("%s: ", p.Options[1])
	fmt.Scan(&userBeta.Password)

	return
}

func profileGET(userBeta user) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(userBeta)

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

	fmt.Printf("\n%s\n", p.Title)

	fmt.Printf("Bienvenido %s tu ID es: %d\n", p.User.Username, p.User.ID)

	return
}
