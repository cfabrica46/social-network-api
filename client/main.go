package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

var (
	errNotAccept           = errors.New("username o password incorrectas")
	errUsernameAlreadyUsed = errors.New("username ya en uso")
	errNotPost             = errors.New("no hay posts")
)

func main() {

	log.SetFlags(log.Lshortfile)

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

			err = profileGET(&u)

			if err != nil {
				log.Fatal(err)
			}

			for !exit {

				loopIntoProfile(u, &exit)

			}

		case 2:

			u, err := login("http://localhost:8080/user/create")

			if err != nil {
				log.Fatal(err)
			}

			err = createUserPOST(&u)

			if err != nil {
				log.Fatal(err)
			}

			err = profileGET(&u)

			if err != nil {
				log.Fatal(err)
			}

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

	fmt.Println("1.Ver Tus Posts")
	fmt.Println("2.Ver Todos Los Posts")
	fmt.Println("3.Eliminar Cuenta")
	fmt.Println("0.Salir")
	fmt.Println()

	fmt.Print(">")

	fmt.Scan(&election)

	fmt.Println()

	switch election {
	case 0:
		*exit = true
	case 1:

		posts, err := getMyPosts(u)

		if err != nil {
			log.Fatal(err)
		}

		printMyPosts(posts)

	case 2:

	case 3:

		var security string

		fmt.Println("¿Esta seguro?[S/N]")
		fmt.Scan(&security)
		security = strings.ToLower(security)

		if security != "s" {
			return
		}

		err := deleteUser(u)

		if err != nil {
			log.Fatal(err)
		}

		*exit = true

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
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&p)

	if err != nil {
		if err != io.EOF {
			return
		}
	}

	if p.Err != "" {

		fmt.Printf("\nERROR: %s\n", p.Err)
		err = errNotAccept
		return

	}

	*user = p.User

	fmt.Printf("\n%s\n", p.Title)

	fmt.Printf("Bienvenido %s tu ID es: %d\n", user.Username, user.ID)

	return
}

func createUserPOST(u *user) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(*u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("POST", "http://localhost:8080/user/create", buf)

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
		if err != io.EOF {
			return
		}
	}

	if p.Err != "" {

		fmt.Printf("\nERROR: %s\n", p.Err)
		err = errUsernameAlreadyUsed
		return

	}

	*u = p.User

	fmt.Printf("\n%s\n", p.Title)

	fmt.Printf("Se creo usuario: %s con contraseña: %s y ID: %d\n", u.Username, u.Password, u.ID)

	return
}

func deleteUser(u user) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/user/delete", buf)

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
		if err != io.EOF {
			return
		}
	}

	if p.Err != "" {
		fmt.Printf("\nERROR: %s\n", p.Err)
		return
	}

	fmt.Printf("\n%s\n", p.Title)

	return

}

func getMyPosts(u user) (posts []post, err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("GET", "http://localhost:8080/user/post/one", buf)

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
		if err != io.EOF {
			return
		}
	}

	if p.Err != "" {
		fmt.Printf("\nERROR: %s\n", p.Err)
		err = nil
		return
	}

	if len(p.Posts) == 0 {
		fmt.Printf("No tienes ningun Post aun\n")
		return
	}

	posts = p.Posts

	return
}

func printMyPosts(posts []post) {

	fmt.Printf("\nTus Posts:\n")

	for i := range posts {

		fmt.Printf("%s: %s - %s\n", posts[i].Propetary, posts[i].Contet, posts[i].Date)

	}
	fmt.Println()
}
