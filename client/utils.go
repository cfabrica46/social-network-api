package main

import "fmt"

func printFriends(friends []User) {

	fmt.Println("Tus amigos:")

	for i := range friends {

		fmt.Printf("%d. %s\n", friends[i].ID, friends[i].Username)

	}
	fmt.Println()

}

func printMyPosts(posts []Post) {

	fmt.Printf("\nPosts:\n")

	for i := range posts {

		fmt.Printf("%d.%s: %s - %s\n", posts[i].ID, posts[i].Propetary, posts[i].Content, posts[i].Date)

	}
	fmt.Println()
}
