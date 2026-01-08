package main

import (
	"errors"
	"fmt"
)

/*
Notes:
If it needs to be accessed outside the package — capitalize it.


    In Go, capitalization controls VISIBILITY (exporting).

    Capital letter → Exported (public)

    Small letter → Unexported (private)

This applies to:

    Structs

    Struct fields

    Functions

    Methods

    Variables

    Constants
 Why JSON Needs Capital Fields?
    JSON package uses reflection

    Reflection sees only exported fields

why go avoids exceptions?
In exception-based languages, errors can jump across many function calls, which makes code harder to read and debug.
Go returns errors as normal values, so failure handling is visible at the exact place it happens.
This makes programs easier to understand, test, and maintain—especially in concurrent systems.


Rules:
If a method mutates state → return nothing.
Behavior that belongs to a type should be a METHOD, not a free function.Free function breaks encapsulation.
If a function logically belongs to a type → make it a method.
Methods should report problems, not decide how to handle them
Functions return errors.main (or the caller) decides what to do with them.
Errors are values in Go.

*/
/*
| Scenario          | Capitalization  |
| ----------------- | --------------- |
| Library / API     | Capital         |
| Models / DTOs     | Capital         |
| Internal helpers  | lowercase       |
| Main package only | lowercase is OK |
*/

type User struct {
	Name string
	ID   int
}

var ErrEmptyName = errors.New("User name cannot be empty")

func (u *User) Rename(newName string) error {
	if newName == "" {

		return ErrEmptyName
	}
	u.Name = newName
	return nil
}
func (u *User) Display() {
	fmt.Printf("User Name is :%s\n", u.Name)
}

func main() {

	newUser := User{
		Name: "Saurabh",
		ID:   1,
	}
	newUser.Display()
	err := newUser.Rename("")
	if err != nil {
		fmt.Println("Rename failed")
	}

	newUser.Display()

}
