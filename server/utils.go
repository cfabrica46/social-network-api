package main

import (
	"database/sql"
	"io"
	"log"
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

			databases, err = migracion()

			if err != nil {

				archivo.Close()
				return
			}

			return
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

func getUser(userBeta User) (u *User) {

	row := db.d.QueryRow("SELECT id,username,password FROM users WHERE username = ? AND password = ?", userBeta.Username, userBeta.Password)

	err := row.Scan(&userBeta.ID, &userBeta.Username, &userBeta.Password)

	if err != nil {
		log.Fatal(err)
		return
	}

	u = &userBeta

	return
}

func insertIntoDatabase(u *User) (err error) {

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

func deleteUserIntoDatabases(u User) (err error) {

	stmt, err := db.d.Prepare("DELETE FROM users WHERE id = ? AND username = ? AND password = ?")

	if err != nil {
		return
	}

	_, err = stmt.Exec(u.ID, u.Username, u.Password)

	if err != nil {
		return
	}

	return
}

func obtainMyPosts(u User) (posts []Post, err error) {

	rows, err := db.d.Query("SELECT posts.id,users.username,posts.content,posts.date FROM users_posts INNER JOIN users ON users_posts.user_ID = users.id INNER JOIN posts ON users_posts.post_ID = posts.id WHERE users.id = ? ORDER BY posts.date DESC", u.ID)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var postAux Post

		err = rows.Scan(&postAux.ID, &postAux.Propetary, &postAux.Content, &postAux.Date)

		if err != nil {
			return
		}

		posts = append(posts, postAux)
	}

	return
}

func obtainAllFriendsPosts(u User) (posts []Post, err error) {

	idFriends, err := obtainIDFriends(u)

	if err != nil {
		return
	}

	for i := range idFriends {

		var rows *sql.Rows

		rows, err = db.d.Query("SELECT posts.id,users.username,posts.content,posts.date FROM users_posts INNER JOIN users ON users_posts.user_ID = users.id INNER JOIN posts ON users_posts.post_ID = posts.id WHERE users.id = ? ORDER BY posts.date DESC", idFriends[i])

		if err != nil {
			return
		}

		defer rows.Close()

		for rows.Next() {

			var postAux Post

			err = rows.Scan(&postAux.ID, &postAux.Propetary, &postAux.Content, &postAux.Date)

			if err != nil {
				return
			}

			posts = append(posts, postAux)
		}

	}

	return
}

func obtainIDFriends(u User) (idFriends []int, err error) {

	rows, err := db.d.Query("SELECT id1,id2 FROM friends WHERE id1=? OR  id2=?", u.ID, u.ID)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var id1, id2 int

		err = rows.Scan(&id1, &id2)

		if err != nil {
			return
		}

		if id1 != u.ID {
			idFriends = append(idFriends, id1)
		} else {
			idFriends = append(idFriends, id2)
		}

	}

	return
}

func insertPostIntoDatabases(post string, userID int) (err error) {

	stmt, err := db.d.Prepare("INSERT INTO posts(content,date) VALUES (?,datetime('now','localtime'))")

	if err != nil {
		return
	}

	res, err := stmt.Exec(post)

	if err != nil {
		return
	}

	id, err := res.LastInsertId()

	if err != nil {
		return
	}

	postID := int(id)

	err = insertPostPivote(postID, userID)

	if err != nil {
		return
	}

	return

}

func insertPostPivote(postID, userID int) (err error) {

	stmt, err := db.d.Prepare("INSERT INTO users_posts(user_ID,post_ID) VALUES (?,?)")

	if err != nil {
		return
	}

	_, err = stmt.Exec(userID, postID)

	if err != nil {
		return
	}

	return
}

func deletePostFromDatabases(postID int) (err error) {

	stmt, err := db.d.Prepare("DELETE FROM posts WHERE id = ? ")

	if err != nil {
		return
	}

	_, err = stmt.Exec(postID)

	if err != nil {
		return
	}

	return

}

func obtainAllFriends(u User) (friends []User, err error) {

	idFriends, err := obtainIDFriends(u)

	if err != nil {
		return
	}

	for i := range idFriends {

		var userAux User

		row := db.d.QueryRow("SELECT id,username FROM users WHERE id = ? ", idFriends[i])

		err = row.Scan(&userAux.ID, &userAux.Username)

		if err != nil {
			return
		}

		friends = append(friends, userAux)

	}

	return
}

func checkIfMyFriendExist(id int, userID int) (check bool, err error) {

	var id1, id2 int

	row := db.d.QueryRow("SELECT id1,id2 FROM friends WHERE (id1 = ? AND id2 = ?) OR (id1 = ? AND id2 = ?) ", id, userID, userID, id)

	err = row.Scan(&id1, &id2)

	if err != nil {
		if err == sql.ErrNoRows {
			check = false
		}
		check = false
		return
	}

	check = true

	return

}

func checkIfMyPostExist(id int, userID int) (check bool, err error) {

	var idPost, idUser int

	row := db.d.QueryRow("SELECT user_ID,post_ID FROM users_posts WHERE post_ID = ? AND user_ID = ?", id, userID, userID, id)

	err = row.Scan(&idUser, &idPost)

	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		return
	}

	check = true

	return

}

func addFriendIntoDatabases(userID, friendID int) (err error) {

	stmt, err := db.d.Prepare("INSERT INTO friends(id1,id2) VALUES (?,?)")

	if err != nil {
		return
	}

	_, err = stmt.Exec(userID, friendID)

	if err != nil {
		return
	}
	return
}

func deleteFriendFromDatabases(friendID, userID int) (err error) {

	stmt, err := db.d.Prepare("DELETE FROM friends WHERE id1 = ? AND id2 = ? OR id1 = ? AND id2 = ?")

	if err != nil {
		return
	}

	_, err = stmt.Exec(friendID, userID, userID, friendID)

	if err != nil {
		return
	}

	return
}
