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

func TestDoubleHelloWorld(*testing.T) {
	log.Printf("Testing Double Hello World... ")
	child, err := Spawn(`sh -c "echo Hello World ; echo Hello ; echo Hi"`)
	if err != nil {
		panic(err)
	}
	err = child.Expect("Hello World")
	if err != nil {
		panic(err)
	}
	err = child.Expect("Hello")
	if err != nil {
		panic(err)
	}
	err = child.Expect("Hi")
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
	sender <- "\n"
	sender <- "echo Hello World\n"
	wait("Hello World")
	sender <- "times\n"
	wait("s")
	sender <- "^D\n"
	log.Printf("Success\n")

}

func TestCommandStart(*testing.T) {
	log.Printf("Testing Command... ")

	// Doing this allows you to modify the cmd struct prior to execution, for example to add environment variables
	child, err := Command("echo 'Hello World'")
	if err != nil {
		panic(err)
	}
	child.Start()
	child.Expect("Hello World")
	log.Printf("Success\n")
}

var regexMatchTests = []struct {
	re   string
	good string
	bad  string
}{
	{`a`, `a`, `b`},
	{`.b`, `ab`, `ac`},
	{`a+hello`, `aaaahello`, `bhello`},
	{`(hello|world)`, `hello`, `unknown`},
	{`(hello|world)`, `world`, `unknown`},
}

func TestRegex(t *testing.T) {
	log.Printf("Testing Regular Expression Matching... ")
	for _, tt := range regexMatchTests {
		runTest := func(input string) bool {
			var match bool
			child, err := Spawn("echo \"\"; echo \"" + input + "\"")
			if err != nil {
				t.Fatal(err)
			}
			match, err = child.ExpectRegex(tt.re)
			if err != nil {
				t.Fatal(err)
			}
			return match
		}
		if !runTest(tt.good) {
			t.Errorf("Regex Not matching [%#q] with pattern [%#q]", tt.good, tt.re)
		}
		if runTest(tt.bad) {
			t.Errorf("Regex Matching [%#q] with pattern [%#q]", tt.bad, tt.re)
		}
	}
	log.Printf("Success\n")
}
