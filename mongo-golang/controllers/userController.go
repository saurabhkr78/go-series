package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	uc := UserController{
		session: s,
	}
	return &uc
}

/*
Note:
use methods when logic needs struct state, and normal functions for stateless utilities
If the code needs data stored inside a struct (state),
then write it as a METHOD with a receiver:
func (receiver) FunctionName(parameters) returnType {
}
If the code does NOT need any struct data (stateless),
then write it as a NORMAL function:
func functionName(parameters)returnType{
}
Receiver → who owns the function
Parameters → what it needs
Return → what it gives back

Multiple return values:func divide(a, b int) (int, error)
Named return values:func divide(a, b int) (result int, err error)
*/
/*
why pointer receiver method to avoid copying the struct and if needed can make modify

now question why we have response writer of not a pointer type whereas request is of pointer type

Q.why?

http.ResponseWriter is an interface whereas http.Request is a struct.

Q2.why they have been made like this?

Interfaces already contain a pointer internally, so a pointer to an interface is almost never needed.
Interface It does not hold data,It defines behavior,The actual response object is handled internally by Go.
so,You are already passing a reference to the real response object i.e w http.ResponseWriter

Q.why request is of pointer type then?
Request object is large,Copying it is expensive,Handlers may need to read body,Modify fields (context, headers) etc
so *http.Request provide shared access and avoid copying


#Ref
Long-lived objects → struct fields
Short-lived objects → request context

*/

/*
Get User steps:
| Step | Question to ask                      |
| ---- | ------------------------------------ |
| 1    | Where is input coming from?          |
| 2    | Is input valid?                      |
| 3    | What resource do I need (DB, cache)? |
| 4    | What work needs to be done?          |
| 5    | What response should client get?     |
*/
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //struct method:Method with pointer receiver on UserController
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // whosoever is making call we need to send them the header which wil contain context, statusCode etc
	}
	session := uc.session.Copy()
	defer session.Close()

}
func (uc *UserController) CreateUser() {}
func (uc *UserController) DeleteUser() {}
