package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func (d *dataBases) profile(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "CONNECT":

		p := page{
			Title:   "Ingrese sus Datos",
			Options: []string{userName, passWord},
		}

		err := json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

	case "GET":

		var userBeta user
		var p page

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		u, err := getUser(userBeta)

		if err != nil {
			if err == sql.ErrNoRows {
				p.Err = http.StatusText(http.StatusNotFound)
				json.NewEncoder(w).Encode(p)
				return
			}
			log.Fatal(err)

			p.Err = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(p)
			return
		}

		p = page{
			Title: "Mi Perfil",
			User:  u,
		}

		err = json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

	}

}

func (d *dataBases) createUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "CONNECT":

		p := page{
			Title:   "Ingrese sus Datos",
			Options: []string{userName, passWord},
		}

		err := json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

	case "POST":

		var userBeta user
		var p page

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		_, err = getUser(userBeta)

		if err != nil {
			if err != sql.ErrNoRows {
				p.Err = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(p)
				return
			}
		}

		err = insertIntoDatabase(&userBeta)

		if err != nil {
			p.Err = http.StatusText(http.StatusLocked)
			json.NewEncoder(w).Encode(p)
			return
		}

		p = page{
			Title: "Nuevo Usuario",
			User:  userBeta,
		}

		err = json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(p)
			return
		}

	}

}

func (d *dataBases) deleteUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "DELETE":

		var p page
		var userBeta user

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		u, err := getUser(userBeta)

		if err != nil {
			if err != sql.ErrNoRows {
				p.Err = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(p)
				return
			}
		}

		err = deleteUserIntoDatabases(u)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(p)
		}

		p = page{
			Title: "Se eliminó la cuenta con éxito",
		}

		err = json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

	}

}

func (d dataBases) getMyPosts(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		var p page
		var userBeta user

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		_, err = getUser(userBeta)

		if err != nil {
			if err != sql.ErrNoRows {
				p.Err = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(p)
				return
			}
		}

		posts, err := obtainMyPosts(userBeta)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(p)
			return
		}

		p = page{
			Posts: posts,
		}

		err = json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

	}
}
