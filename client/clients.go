package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func profileGET(u *User) (check bool, err error) {

	var userAux User
	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(*u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("GET", "http://localhost:8080/user", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &userAux)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}

			fmt.Println()
			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	check = true
	*u = userAux
	return

}

func createUser(u *User) (err error) {

	var userAux User
	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(*u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("POST", "http://localhost:8080/user", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &userAux)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	*u = userAux

	fmt.Printf("Se creo usuario: %s con contraseña: %s y ID: %d\n", u.Username, u.Password, u.ID)
	return

}

func deleteUser(u User) (err error) {

	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/user", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &u)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	fmt.Printf("\nSe eliminó la cuenta con éxito,\n")

	return

}

func getPosts(u User, stringURL string) (posts []Post, err error) {

	var errString string

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

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &posts)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	return
}

func addPost(postContent string, u User) (err error) {

	var postAux Post
	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	post := Post{Content: postContent}

	dataJSON, err := json.Marshal(struct {
		User
		Post
	}{
		u,
		post,
	})

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("POST", "http://localhost:8080/user/posts", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &postAux)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	fmt.Printf("\nSe publico tu post con exito\n")

	return
}

func deletePost(u User, postID int) (err error) {

	var postAux Post
	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	post := Post{ID: postID}

	dataJSON, err := json.Marshal(struct {
		User
		Post
	}{
		u,
		post,
	})

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/user/posts", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &postAux)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	fmt.Printf("\nSe eliminó tu post con exito\n")

	return

}

func getFriends(u User) (friends []User, err error) {

	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	dataJSON, err := json.Marshal(u)

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("GET", "http://localhost:8080/user/friends", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &friends)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	return
}

func addFriend(u User, friendID int) (err error) {

	var friendAux User
	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	friend := User{ID: friendID}

	dataJSON, err := json.Marshal(struct {
		User
		Friend User
	}{
		u,
		friend,
	})

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("POST", "http://localhost:8080/user/friends", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &friendAux)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	fmt.Printf("\nSe añadio un nuevo amigo con éxito\n")

	return

}

func deleteFriend(u User, friendID int) (err error) {

	var friendAux User
	var errString string

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	friend := User{ID: friendID}

	dataJSON, err := json.Marshal(struct {
		User
		Friend User
	}{
		u,
		friend,
	})

	if err != nil {
		return
	}

	buf := bytes.NewBuffer(dataJSON)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/user/friends", buf)

	if err != nil {
		return
	}

	res, err := client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &friendAux)

	if err != nil {
		if err != io.EOF {
			err = json.Unmarshal(data, &errString)
			if err != nil {
				if err != io.EOF {
					return
				}
				err = nil
			}
			fmt.Println()

			fmt.Printf("ERROR: %s\n", errString)
			err = nil
			return

		} else {
			err = nil
		}
	}

	fmt.Printf("\nSe eliminó tu amistad con exito\n")

	return

}
