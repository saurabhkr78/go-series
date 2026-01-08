package main

/*
Errors are for callers.
Panics are for developers.


Creating errors
errors.New("something went wrong")

or with context:
fmt.Errorf("read config failed: %w", err)

Sentinel errors
var ErrEmptyName = errors.New("name cannot be empty") Why?

    Can be compared

    Can be handled specifically
if err == ErrEmptyName {

}



*/

import ()

func main() {

}
