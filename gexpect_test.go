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
	child, err := spawn("/usr/bin/env python")
	if err != nil {
		panic(err)
	}

	// child.expect(">>>")

	child.interact()
}

func ExampleFTP() {
	child, err := spawn("ftp ftp.openbsd.org")
	if err != nil {
		panic(err)
	}
	child.expect("(?i)name .*: ")
	child.sendline("anonymous")
	child.expect("(?i)password")
	child.sendline("pexpect@sourceforge.net")
	child.expect("ftp> ")
	child.sendline("cd /pub/OpenBSD/3.7/packages/i386")
	child.expect("ftp> ")
	child.sendline("bin")
	child.expect("ftp> ")
	child.sendline("prompt")
	child.expect("ftp> ")
	child.sendline("pwd")
	child.expect("ftp> ")
	print("Escape character is '^]'.\n")
	// sys.stdout.write(child.after)
	// sys.stdout.flush()
	child.interact()
}

func main() {
	ExampleFTP()
}
