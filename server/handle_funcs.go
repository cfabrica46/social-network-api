package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func (d *dataBases) user(w http.ResponseWriter, r *http.Request) {

	errMensaje := struct {
		Mensaje string
	}{}

	switch r.Method {
	case "GET":

		var userBeta User
		var err error

		idBeta := r.Header.Get("id")

		userBeta.ID, err = strconv.Atoi(idBeta)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		userBeta.Username = r.Header.Get("username")
		userBeta.Password = r.Header.Get("password")

		u := getUser(userBeta)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)

			return

		}

		err = json.NewEncoder(w).Encode(*u)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

	case "POST":

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		u := getUser(userBeta)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		err = insertIntoDatabase(u)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusLocked)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(*u)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

	case "DELETE":

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		u := getUser(userBeta)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		err = deleteUserIntoDatabases(*u)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}
	}

}

func (d dataBases) myPosts(w http.ResponseWriter, r *http.Request) {

	errMensaje := struct {
		Mensaje string
	}{}

	switch r.Method {

	case "GET":

		var userBeta User
		var err error

		idBeta := r.Header.Get("id")

		userBeta.ID, err = strconv.Atoi(idBeta)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		userBeta.Username = r.Header.Get("username")
		userBeta.Password = r.Header.Get("password")

		u := getUser(userBeta)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		posts, err := obtainMyPosts(*u)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(posts)
				return
			}
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(posts)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

	case "POST":

		data := struct {
			User
			Post
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		u := getUser(data.User)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		err = insertPostIntoDatabases(data.Post.Content, u.ID)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(data.Post)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

	case "DELETE":

		data := struct {
			User
			Post
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		u := getUser(data.User)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		check, err := checkIfMyPostExist(data.Post.ID, data.User.ID)

		if err != nil {
			if err != sql.ErrNoRows {
				errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errMensaje)
				return
			}
		}

		if !check {
			errMensaje.Mensaje = http.StatusText(http.StatusSeeOther)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = deletePostFromDatabases(data.Post.ID)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(data.Post)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}
	}
}

func (d dataBases) friendsPosts(w http.ResponseWriter, r *http.Request) {

	errMensaje := struct {
		Mensaje string
	}{}

	if r.Method == "GET" {

		var userBeta User
		var err error

		idBeta := r.Header.Get("id")

		userBeta.ID, err = strconv.Atoi(idBeta)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		userBeta.Username = r.Header.Get("username")
		userBeta.Password = r.Header.Get("password")

		u := getUser(userBeta)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		posts, err := obtainAllFriendsPosts(*u)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(posts)
				return
			}
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(posts)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

	}
}

func (d dataBases) friends(w http.ResponseWriter, r *http.Request) {

	errMensaje := struct {
		Mensaje string
	}{}

	switch r.Method {

	case "GET":

		var userBeta User
		var err error

		idBeta := r.Header.Get("id")

		userBeta.ID, err = strconv.Atoi(idBeta)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errMensaje)
		}

		userBeta.Username = r.Header.Get("username")
		userBeta.Password = r.Header.Get("password")

		u := getUser(userBeta)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		friends, err := obtainAllFriends(*u)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(friends)
				return
			}
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(friends)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

	case "POST":

		data := struct {
			User
			Friend User
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		u := getUser(data.User)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		check := checkIfMyFriendAlreadyIsMyFriend(data.Friend.ID, data.User.ID)

		if check {
			errMensaje.Mensaje = http.StatusText(http.StatusSeeOther)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = addFriendIntoDatabases(data.User.ID, data.Friend.ID)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(data.Friend)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

	case "DELETE":

		data := struct {
			User
			Friend User
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		u := getUser(data.User)

		if u == nil {

			errMensaje.Mensaje = http.StatusText(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(errMensaje)
			return

		}

		check := checkIfMyFriendAlreadyIsMyFriend(data.Friend.ID, data.User.ID)

		if !check {
			errMensaje.Mensaje = http.StatusText(http.StatusSeeOther)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = deleteFriendFromDatabases(data.Friend.ID, data.User.ID)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}

		err = json.NewEncoder(w).Encode(data.Friend)

		if err != nil {
			errMensaje.Mensaje = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errMensaje)
			return
		}
	}

}
