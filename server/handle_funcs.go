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

func (d dataBases) getAllFriendsPosts(w http.ResponseWriter, r *http.Request) {

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

		posts, err := obtainAllFriendsPosts(userBeta)

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

func (d dataBases) addPost(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		var p page

		err := json.NewDecoder(r.Body).Decode(&p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		_, err = getUser(p.User)

		if err != nil {
			if err != sql.ErrNoRows {
				p.Err = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(p)
				return
			}
		}

		err = insertPostIntoDatabases(p.Posts[0].Content, p.User.ID)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

	}

}

func (d dataBases) deletePost(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {

		var p page

		err := json.NewDecoder(r.Body).Decode(&p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		_, err = getUser(p.User)

		if err != nil {
			if err != sql.ErrNoRows {
				p.Err = http.StatusText(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(p)
				return
			}
		}

		err = deletePostFromDatabases(p.Posts[0].ID, p.User.ID)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(p)

		if err != nil {
			p.Err = http.StatusText(http.StatusInternalServerError)
			return
		}

	}

}

func deletePostFromDatabases(postID int, userID int) (err error) {

	stmt, err := db.d.Prepare("DELETE FROM posts WHERE id = ?")

	if err != nil {
		return
	}

	_, err = stmt.Exec(postID)

	if err != nil {
		return
	}

	return
}
