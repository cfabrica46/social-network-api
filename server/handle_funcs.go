package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func findUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		u, err := getUser(r.URL.Query())

		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), 400)
			return
		}

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			return
		}
	}
}

func findUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode(users)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			return
		}
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		fmt.Fprint(w, "Introduce tus Datos\n")

	case "POST":

		var u user

		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		idInt64 := time.Now().UnixNano()

		id := int(idInt64)

		u.ID = id

		users = append(users, u)

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}

}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		u, err := getUser(r.URL.Query())

		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), 400)
			return
		}

		if len(users) <= 1 {

			users = []user{}

		} else {

			for i := range users {
				if users[i] == u {
					users = append(users[:i], users[i+1:]...)
					break
				}
			}
		}

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			return
		}
	}

}

func updateUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		u, err := getUser(r.URL.Query())

		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), 400)
			return
		}

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			return
		}

	case "POST":

		var u user

		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		for i := range users {

			if users[i].ID == u.ID {

				users[i] = u
				break

			}

		}

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}

}
