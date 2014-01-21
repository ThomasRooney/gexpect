package main

import "github.com/ThomasRooney/gexpect"
import "fmt"

func main() {
	fmt.Printf("Starting screen.. \n")
	child, err := gexpect.Spawn("screen")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Expecting @.. \n")
	child.Expect("@")

	fmt.Printf("Interacting.. \n")
	child.Interact()
	fmt.Printf("Done \n")
}
