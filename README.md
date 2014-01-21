## Gexpect

Gexpect is a pure golang expect-like module.

It makes it simple and safe to control other terminal applications.  

pexpect-like syntax for golang

	child, err := gexpect.Spawn("python")
	if err != nil {
		panic(err)
	}
	child.Expect(">>>")
	child.Sendline("print 'Hello World'")
	child.Interact()
	child.Close()

Free,  MIT open source licenced, etc etc.

Check gexpect_test.go for full examples