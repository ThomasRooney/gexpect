package main

import "github.com/ThomasRooney/gexpect"
import "fmt"

func main() {
	fmt.Printf("Starting screen.. \n")
	child, err := gexpect.Spawn("screen")
	if err != nil {
		panic(err)
	}
	ch := child.AsyncInteractBiChannel()
	ch <- "echo Hello World"
	ch <- "echo Hello World"
	ch <- "echo Hello World"
}
