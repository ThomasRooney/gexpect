package gexpect

import (
	"fmt"
	"testing"
)

func ExampleHelloWorld() {

	fmt.Println("hello")
	// Output: hello
}

func TestSpawn(*testing.T) {
}

func TestInteract(*testing.T) {

}

func ExampleFTP() {
	child, err := Spawn("ftp ftp.openbsd.org")
	if err != nil {
		panic(err)
	}
	child.Expect("(?i)name .*: ")
	child.Sendline("anonymous")
	child.Expect("(?i)password")
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
	print("Escape character is '^]'.\n")
	// sys.stdout.write(child.after)
	// sys.stdout.flush()
	child.Interact()
}

func main() {
	ExampleFTP()
}
