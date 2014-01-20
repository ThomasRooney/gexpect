package main

import "github.com/ThomasRooney/gexpect"

func main() {
	child, err := gexpect.Spawn("/usr/bin/env python")
	if err != nil {
		panic(err)
	}

	// child.expect(">>>")

	child.Interact()
}
