package gexpect

import (
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	t.Logf("Testing Hello World... ")
	child, err := Spawn("echo \"Hello World\"")
	if err != nil {
		t.Fatal(err)
	}
	err = child.Expect("Hello World")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Success\n")
}

func TestDoubleHelloWorld(t *testing.T) {
	t.Logf("Testing Double Hello World... ")
	child, err := Spawn(`sh -c "echo Hello World ; echo Hello ; echo Hi"`)
	if err != nil {
		t.Fatal(err)
	}
	err = child.Expect("Hello World")
	if err != nil {
		t.Fatal(err)
	}
	err = child.Expect("Hello")
	if err != nil {
		t.Fatal(err)
	}
	err = child.Expect("Hi")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Success\n")
}

func TestHelloWorldFailureCase(t *testing.T) {
	t.Logf("Testing Hello World Failure case... ")
	child, err := Spawn("echo \"Hello World\"")
	if err != nil {
		t.Fatal(err)
	}
	err = child.Expect("YOU WILL NEVER FIND ME")
	if err != nil {
		t.Logf("Success\n")
		return
	}
	t.Fatal("Expected an error for TestHelloWorldFailureCase")
}

func TestBiChannel(t *testing.T) {
	t.Logf("Testing BiChannel screen... ")
	child, err := Spawn("screen")
	if err != nil {
		t.Fatal(err)
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
	t.Logf("Success\n")

}

func TestCommandStart(t *testing.T) {
	t.Logf("Testing Command... ")

	// Doing this allows you to modify the cmd struct prior to execution, for example to add environment variables
	child, err := Command("echo 'Hello World'")
	if err != nil {
		t.Fatal(err)
	}
	child.Start()
	child.Expect("Hello World")
	t.Logf("Success\n")
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
	t.Logf("Testing Regular Expression Matching... ")
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
	t.Logf("Success\n")
}
