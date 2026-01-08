package main

//Interfaces exist to REMOVE switch-case logic, not to wrap it.
/*
Interface decides what can be done
Struct decides how it is done
so,“All implementations must look the same from outside.”i.e All methods that implement an interface must have the EXACT same method signature.
Interface defines the shape of behavior, not the shape of data.


Go achieves OOP using structs, methods, interfaces, and composition.

struct replaces class
Methods are attached to structs
no inheritence using extend keyword Go uses composition i.e Composition means building objects by combining smaller objects, instead of inheriting from them.
Why Composition Exists (The Problem It Solves)
Inheritance problem (Java )
class Animal {
    void eat() {}
}

class Dog extends Animal {
    void bark() {}
}
Looks fine… until:

    Animal changes → all children break

    Deep inheritance trees form

    Behavior becomes tightly coupled

This is called the Fragile Base Class problem.

Composition in Simple Words

Instead of saying:

    “Dog is an Animal”

You say:

    “Dog has an Animal behavior”



Java: Composition Example


class Engine {
    void start() {}
}

class Car {
    private Engine engine = new Engine();

    void drive() {
        engine.start();
        System.out.println("Driving");
    }
}

👉 Car has an Engine
👉 Not extends Engine


Go: Composition Example (VERY IMPORTANT)
Without inheritance (Go way)


type Engine struct {}

func (e Engine) Start() {
    fmt.Println("Engine started")
}

type Car struct {
    Engine
}

func main() {
    c := Car{}
    c.Start() // promoted method
}

What happened?

    Car contains Engine

    Car automatically gets Start()

    No inheritance

    No hierarchy

👉 This is composition with method promotion


Composition + Interface (Most Powerful Pattern)

type Engine interface {
    Start()
}

type PetrolEngine struct {}

func (p PetrolEngine) Start() {
    fmt.Println("Petrol engine started")
}

type Car struct {
    Engine
}

Now:

    Car depends on behavior, not implementation

    You can swap engines easily

Benefits:

    Loose coupling

    Easier testing

    Better flexibility

    Fewer side effects

    Clear ownership

Go philosophy:

    “Don’t build type hierarchies. Build behavior.”


so,Composition is a design technique where a type is built by including other types, rather than inheriting from them.
Go uses composition with structs and interfaces instead of inheritance to achieve code reuse and flexibility.

Inheritance copies behavior.
Composition shares behavior.











If the method exists → interface is satisfied. no need to write extends in methods.This reduces coupling and improves flexibility.
Polymorphism via interfaces.Based on method sets, not class hierarchy
var n Notifier
n = &Email{}
n.Notify()

encapsulatio is simplerin go
Name string  // exported
name string  // unexported
Capitalization controls visibility — no getters/setters by default.


summary:
| OOP Concept    | How Go Achieves It                   |
| -------------- | ------------------------------------ |
| Encapsulation  | Structs + exported/unexported fields |
| Abstraction    | Interfaces                           |
| Polymorphism   | Interface-based method dispatch      |
| Inheritance    | Replaced by composition              |
| Method binding | Methods with receivers               |


Go is object-oriented, but it uses composition and interfaces instead of classes and inheritance.
This results in simpler, more maintainable designs.

Go supports object-oriented design, not object-oriented syntax.

Encapsulation:
type Email struct {
	Address string
}

func (e *Email) Notify(...) error   -> this is

    Data + behavior together

    Logic hidden inside methods

    This is classic encapsulation





Polymorphism :
var n Notifier
n = &Email{}
n = &SMS{}

Same call: n.Notify(message) Different behavior.




Abstraction

type Notifier interface {
	Notify(string) error
}

    Caller depends on what, not how

    Same idea as abstract base classes — but cleaner


*/

import (
	"errors"
	"fmt"
)

type Email struct {
	Address string
}
type SMS struct {
	MobNo int
}
type Whatsapp struct {
	WhNo int
}

type Notifier interface {
	Notify(message string) error
}

func (e *Email) Notify(message string) error {
	if e.Address == "" {
		return errors.New("Email address missing")
	}
	fmt.Printf("message sent to email address %s with message %s", e.Address, message)

	return nil
}
func (s *SMS) Notify(message string) error {
	fmt.Printf("message sent to mobile Number %d with message %s", s.MobNo, message)

	return nil
}
func (w *Whatsapp) Notify(message string) error {
	fmt.Printf("message sent to whatsapp Number %s with message %s", w.WhNo, message)

	return nil
}
func send(n Notifier, message string) error {

	if err := n.Notify(message); err != nil {
		return fmt.Errorf("send failed %w", err)
	}
	return nil

}

func main() {

	newMessage := "hello! this message is being sent on the behalf of simple fintech"
	/*read it: Create a slice where each element is a Notifier,
	  and store different concrete objects inside it.”
	*/
	notifiers := []Notifier{
		&Email{Address: "test@simpletech.com"},
		&SMS{MobNo: 917886666666},
		&Whatsapp{WhNo: 917886666666},
	}

	for _, n := range notifiers {
		if err := send(n, newMessage); err != nil {
			fmt.Println("notification sent failed", err)
		}
	}
}
