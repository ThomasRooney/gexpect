package gexpect

import (
	"testing"
)

func TestHelloWorld(*testing.T) {
	child, err := Spawn("echo \"Hello World\"")
	if err != nil {
		panic(err)
	}
	child.Expect("Hello World")
}

func TestExpectFtp(*testing.T) {
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
}

func TestInteractPing(*testing.T) {
	child, err := Spawn("ping -c8 8.8.8.8")
	if err != nil {
		panic(err)
	}
	child.Interact()
}
