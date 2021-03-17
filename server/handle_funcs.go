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

func insertIntoDatabase(u *user) (err error) {

	stmt, err := db.d.Prepare("INSERT INTO users(username,password) VALUES (?,?)")

	if err != nil {
		return
	}
	res, err := stmt.Exec(u.Username, u.Password)

	if err != nil {
		return
	}

	id, err := res.LastInsertId()

	if err != nil {
		return
	}

	u.ID = int(id)

	return
}
