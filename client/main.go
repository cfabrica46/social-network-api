package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

		fmt.Print("> ")

		fmt.Scan(&election)

		fmt.Println()

		switch election {

		case 0:

			exit = true

		case 1:

			u, err := login()

			if err != nil {
				log.Fatal(err)
			}

			c, err := profileGET(&u)

			if err != nil {
				log.Fatal(err)
			}

			if c {

				for !exit {

					loopIntoProfile(u, &exit)

				}
			}
		case 2:

			u, err := login()

			if err != nil {
				log.Fatal(err)
			}

			err = createUser(&u)

			if err != nil {
				log.Fatal(err)
			}

			c, err := profileGET(&u)

			if err != nil {
				log.Fatal(err)
			}

			if c {
				for !exit {

					loopIntoProfile(u, &exit)

				}
			}
		default:

			fmt.Println("Seleccione una opción válida")

		}

	}
}

func loopIntoProfile(u User, exit *bool) {

	var election int

	fmt.Println()
	fmt.Printf("Bienvenido %s\n", u.Username)
	fmt.Println("¿Qué deseas hacer?")
	fmt.Println()

	fmt.Println("1.Ver Tus Posts")
	fmt.Println("2.Ver Todos Los Posts De Tus Amigos")
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

		fmt.Println("Escribe lo que deseas publicar")

		fmt.Print("> ")

		reader := bufio.NewReader(os.Stdin)

		mensaje, err := reader.ReadString('\n')

		if err != nil {
			return
		}

		mensaje, err = reader.ReadString('\n')

		if err != nil {
			return
		}

		err = addPost(mensaje, u)

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
