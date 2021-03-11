package main

import (
	"net/http"
	"net/url"
	"strconv"
)

func parseID(r url.Values) (key int, err error) {

	keys, ok := r["id"]

	if !ok || len(keys) != 1 {

		err = http.ErrAbortHandler

		return

	}

	key, err = strconv.Atoi(keys[0])

	if err != nil {

		return
	}

	return
}

func getUser(r url.Values) (u user, err error) {

	var check bool

	key, err := parseID(r)

	for i := range users {

		if users[i].ID == key {

			check = true
			u = users[i]
			break
		}

	}

	if check == false {
		err = http.ErrAbortHandler
		return
	}
	return
}
