package main

/*
Structs are VALUE types
Slices and maps are REFERENCE types
so,
When you pass a struct to a function:
A copy is made
Changes do NOT affect the original

Note:to change the original use pointer


slice is a descriptor not an array it has pointer to array, length and capacity.
Maps is true refence type contains pointers to runtime hastable. no pointer needed always mutate original
*/

/*
Scenario	Pass Struct As
Read-only	Value
Store in slice	Value
Modify	Pointer
Large struct	Pointer
Returned from slice	Pointer




*/
import (
	"fmt"
)

type User struct {
	name string
	ID   int
}

/*
users := make([]User, 0, 100)-> preallocated slice like vector in c++
User → struct type
[]User → slice of User
*[]User → pointer to slice of User
*/

func addUser(users *[]User, user User) {
	*users = append(*users, user)
}

func getUserByID(users []User, id int) (*User, bool) {
	for i := range users {
		if users[i].ID == id {
			return &users[i], true
		}
	}
	return nil, false
}

func listUsers(users []User) {
	for _, user := range users {
		fmt.Printf("%d ID: %s Name\n", user.ID, user.name)
	}
}

func main() {
	var users []User
	addUser(&users, User{ID: 1, name: "Sudo"})
	addUser(&users, User{ID: 2, name: "Doe"})

	listUsers(users)

	result, found := getUserByID(users, 1)
	if found {
		fmt.Println("User Found", result.name)
	}

}
