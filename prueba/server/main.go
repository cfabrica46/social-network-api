package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type user struct {
	ID                 int
	Username, Password string
}

func main() {

	http.HandleFunc("/", home)

	http.ListenAndServe(":8080", nil)

}

func home(w http.ResponseWriter, r *http.Request) {

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

		err = json.NewEncoder(w).Encode(u)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}

}
