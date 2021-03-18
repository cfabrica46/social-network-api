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
	Friends []user
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
	Content   string
	Date      string
}

var (
	errNotAccept           = errors.New("username o password incorrectas")
	errUsernameAlreadyUsed = errors.New("username ya en uso")
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
	fmt.Println("3.Añadir Un Post")
	fmt.Println("4.Eliminar Un Post")
	fmt.Println("5.Mostrar Amigos")
	fmt.Println("6.Añadir Amigo")
	fmt.Println("7.Eliminar Amigo")
	fmt.Println("8.Eliminar Cuenta")
	fmt.Println("0.Salir")
	fmt.Println()

	fmt.Print("> ")

	fmt.Scan(&election)

	fmt.Println()

	switch election {
	case 0:
		*exit = true
	case 1:

		posts, err := getPosts(u, "http://localhost:8080/user/post/one")

		if err != nil {
			log.Fatal(err)
		}

		if len(posts) == 0 {
			fmt.Printf("No tienes ningun Post aun\n")
			return
		}

		printMyPosts(posts)

	case 2:

		posts, err := getPosts(u, "http://localhost:8080/user/post/all")

		if err != nil {
			log.Fatal(err)
		}

		if len(posts) == 0 {
			fmt.Printf("Tus amigos no han publicado ningun post\n")
			return
		}

		printMyPosts(posts)

	case 3:

		var post string

		fmt.Println("Escribe lo que deseas publicar")

		fmt.Print("> ")

		fmt.Scan(&post)

		err := addPost(post, u)

		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}

	case 4:

		var postID int

		fmt.Println("Escribe el ID del post que deseas eliminar")

		fmt.Print("> ")

		fmt.Scan(&postID)

		err := deletePost(u, postID)

		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}

	case 5:

		friends, err := getFriends(u)

		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}

		if len(friends) == 0 {
			fmt.Printf("No tienes amigos\n")
			return
		}

		printFriends(friends)

	case 6:

		var friendID int

		fmt.Println("Escribe el ID del amigo a agregar")
		fmt.Print("> ")
		fmt.Scan(&friendID)

		err := addFriend(u, friendID)

		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}

	case 7:
		var friendID int

		fmt.Println("Escribe el ID del amigo que deseas eliminar")
		fmt.Print("> ")
		fmt.Scan(&friendID)

		err := deleteFriend(u, friendID)

		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}

	case 8:

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

func getPosts(u user, stringURL string) (posts []post, err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("GET", stringURL, buf)

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

	posts = p.Posts

	return
}

func printMyPosts(posts []post) {

	fmt.Printf("\nPosts:\n")

	for i := range posts {

		fmt.Printf("%d.%s: %s - %s\n", posts[i].ID, posts[i].Propetary, posts[i].Content, posts[i].Date)

	}
	fmt.Println()
}

func addPost(po string, u user) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	p = page{
		User:  u,
		Posts: []post{{Content: po}},
	}

	dataJSON, err := json.Marshal(p)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("POST", "http://localhost:8080/user/post/add", buf)

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

	fmt.Printf("\nSe publico tu post con exito\n")

	return
}

func deletePost(u user, postID int) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	p = page{
		User:  u,
		Posts: []post{{ID: postID}},
	}

	dataJSON, err := json.Marshal(p)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/user/post/delete", buf)

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

	fmt.Printf("\nSe eliminó tu post con exito\n")

	return

}

func getFriends(u user) (friends []user, err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("GET", "http://localhost:8080/user/friends/show", buf)

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

	friends = p.Friends

	return
}

func printFriends(friends []user) {

	fmt.Println("Tus amigos:")

	for i := range friends {

		fmt.Printf("%d. %s\n", friends[i].ID, friends[i].Username)

	}
	fmt.Println()

}

func addFriend(u user, friendID int) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	p = page{
		User:    u,
		Friends: []user{{ID: friendID}},
	}

	dataJSON, err := json.Marshal(p)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("POST", "http://localhost:8080/user/friends/add", buf)

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

	fmt.Printf("\nSe añadio un nuevo amigo con éxito\n")

	return

}

func deleteFriend(u user, friendID int) (err error) {

	var p page

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	p = page{
		User:    u,
		Friends: []user{{ID: friendID}},
	}

	dataJSON, err := json.Marshal(p)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/user/friends/delete", buf)

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

	fmt.Printf("\nSe eliminó tu amistad con exito\n")

	return

}
