package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (d *dataBases) user(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		u, err := getUser(userBeta)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

	case "POST":

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = getUser(userBeta)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = insertIntoDatabase(&userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusLocked))
			return
		}

		err = json.NewEncoder(w).Encode(userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

	case "DELETE":

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		u, err := getUser(userBeta)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = deleteUserIntoDatabases(u)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}
	}

}

func (d dataBases) myPosts(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = getUser(userBeta)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		posts, err := obtainMyPosts(userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(posts)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

	case "POST":

		data := struct {
			User
			Post
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		u, err := getUser(data.User)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = insertPostIntoDatabases(data.Post.Content, u.ID)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(data.Post)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

	case "DELETE":

		data := struct {
			User
			Post
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = getUser(data.User)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		check, err := checkIfMyPostExist(data.Post.ID, data.User.ID)

		if err != nil {
			if err != sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		if !check {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusSeeOther))
			return
		}

		err = deletePostFromDatabases(data.Post.ID)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(data.Post)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}
	}
}

func (d dataBases) friendsPosts(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = getUser(userBeta)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		posts, err := obtainAllFriendsPosts(userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(posts)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

	}
}

func (d dataBases) friends(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		var userBeta User

		err := json.NewDecoder(r.Body).Decode(&userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = getUser(userBeta)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		friends, err := obtainAllFriends(userBeta)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(friends)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

	case "POST":

		data := struct {
			User
			Friend User
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = getUser(data.User)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		check, err := checkIfMyFriendExist(data.Friend.ID, data.User.ID)

		if err != nil {
			if err != sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		if check {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusSeeOther))
			return
		}

		err = addFriendIntoDatabases(data.User.ID, data.Friend.ID)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(data.Friend)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

	case "DELETE":

		data := struct {
			User
			Friend User
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = getUser(data.User)

		if err != nil {
			if err == sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusNetworkAuthenticationRequired))
				return
			}
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		check, err := checkIfMyFriendExist(data.Friend.ID, data.User.ID)

		if err != nil {
			if err != sql.ErrNoRows {
				json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		if !check {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusSeeOther))
			return
		}

		err = deleteFriendFromDatabases(data.Friend.ID, data.User.ID)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}

		err = json.NewEncoder(w).Encode(data.Friend)

		if err != nil {
			json.NewEncoder(w).Encode(http.StatusText(http.StatusInternalServerError))
			return
		}
	}

}
