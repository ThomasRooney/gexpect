package gexpect

import (
	"log"
	"strings"
	"testing"
)

func TestHelloWorld(*testing.T) {
	log.Printf("Testing Hello World... ")
	child, err := Spawn("echo \"Hello World\"")
	if err != nil {
		panic(err)
	}
	err = child.Expect("Hello World")
	if err != nil {
		panic(err)
	}
	log.Printf("Success\n")
}

func TestHelloWorldFailureCase(*testing.T) {
	log.Printf("Testing Hello World Failure case... ")
	child, err := Spawn("echo \"Hello World\"")
	if err != nil {
		panic(err)
	}
	err = child.Expect("YOU WILL NEVER FIND ME")
	if err != nil {
		log.Printf("Success\n")
		return
	}
	panic("Expected an error for TestHelloWorldFailureCase")
}

func TestBiChannel(*testing.T) {
	log.Printf("Testing BiChannel screen... ")
	child, err := Spawn("screen")
	if err != nil {
		panic(err)
	}
	sender, reciever := child.AsyncInteractChannels()
	wait := func(str string) {
		for {
			msg, open := <-reciever
			if !open {
				return
			}
			if strings.Contains(msg, str) {
				return
			}
		}
	}
	sender <- ""
	sender <- "echo Hello World"
	wait("Hello World")
	sender <- "times"
	wait("s")
	sender <- "^D"
	log.Printf("Success\n")

}

func TestExpectRegex(*testing.T) {
	log.Printf("Testing ExpectRegex... ")

	child, err := Spawn("/bin/sh times")
	if err != nil {
		panic(err)
	}
	child.ExpectRegex("Name")
	log.Printf("Success\n")

}

func TestExpectFtp(*testing.T) {
	log.Printf("Testing Ftp... ")

	child, err := Spawn("ftp ftp.openbsd.org")
	if err != nil {
		panic(err)
	}
	child.Expect("Name")
	child.Sendline("anonymous")
	child.Expect("Password")
	child.Sendline("pexpect@sourceforge.net")
	child.Expect("ftp> ")
	child.Sendline("cd /pub/OpenBSD/3.7/packages/i386")
	child.Expect("ftp> ")
	child.Sendline("bin")
	child.Expect("ftp> ")
	child.Sendline("prompt")
	child.Expect("ftp> ")
	child.Sendline("pwd")
	child.Expect("ftp> ")
	log.Printf("Success\n")

}

func TestInteractPing(*testing.T) {
	log.Printf("Testing Ping interact... \n")

	child, err := Spawn("ping -c8 8.8.8.8")
	if err != nil {
		panic(err)
	}
	child.Interact()
	log.Printf("Success\n")

}
