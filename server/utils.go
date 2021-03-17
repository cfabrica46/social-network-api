package main

import (
	"database/sql"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func migracion() (databases *sql.DB, err error) {

	archivoDB, err := os.Create("databases.db")

	if err != nil {
		return
	}
	archivoDB.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	archivoSQL, err := os.Open("databases.sql")

	if err != nil {
		return
	}

	defer archivoSQL.Close()

	contendio, err := io.ReadAll(archivoSQL)

	if err != nil {
		return
	}

	_, err = databases.Exec(string(contendio))
	if err != nil {
		return
	}

	return

}

func open() (databases *sql.DB, err error) {

	archivo, err := os.Open("databases.db")

	if err != nil {
		if os.IsNotExist(err) {

			databases, err := migracion()

			if err != nil {

				archivo.Close()
				return databases, err
			}

			return databases, err
		}
		return
	}
	defer archivo.Close()

	databases, err = sql.Open("sqlite3", "./databases.db?_foreign_keys=on")

	if err != nil {
		return
	}

	return
}

func getUser(userBeta user) (u user, err error) {

	var userAux user

	row := db.d.QueryRow("SELECT id,username,password FROM users WHERE username = ? AND password = ?", userBeta.Username, userBeta.Password)

	err = row.Scan(&userAux.ID, &userAux.Username, &userAux.Password)

	if err != nil {
		return
	}

	u = userAux

	return
}
