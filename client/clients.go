package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
