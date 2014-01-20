package gexpect

import (
	"bufio"
	"errors"
	shell "github.com/kballard/go-shellquote"
	"io"
	"os"
	"os/exec"
)

type ExpectSubprocess struct {
	cmd    *exec.Cmd
	output *bufio.Reader
	stdin  *io.WriteCloser
}

func spawnAtDirectory(command string, directory string) (*ExpectSubprocess, error) {
	expect, err := _spawn(command)
	if err != nil {
		return nil, err
	}
	expect.cmd.Dir = directory
	return _start(expect)
}

func spawn(command string) (*ExpectSubprocess, error) {
	expect, err := _spawn(command)
	if err != nil {
		return nil, err
	}
	return _start(expect)
}

func (expect *ExpectSubprocess) expect(searchString string) error {

	chunk := make([]byte, len(searchString))
	found := 0
	target := len(searchString)
	for {

		n, err := expect.output.Read(chunk)

		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			if chunk[i] == searchString[found] {
				found++
				if found == target {
					return nil
				}
			} else {
				found = 0
			}
		}
	}
}

func (expect *ExpectSubprocess) sendline(command string) error {
	_, err := io.WriteString(*expect.stdin, command+"\r\n")
	return err
}

func (expect *ExpectSubprocess) interact() {
	defer expect.cmd.Wait()
	go io.Copy(os.Stdout, expect.output)
	go io.Copy(*expect.stdin, os.Stdin)
}

func (expect *ExpectSubprocess) readLine() (string, error) {
	str, err := expect.output.ReadString('\n')
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func _start(expect *ExpectSubprocess) (*ExpectSubprocess, error) {

	err := expect.cmd.Start()
	if err != nil {
		return nil, err
	}

	return expect, nil
}

func _spawn(command string) (*ExpectSubprocess, error) {
	wrapper := new(ExpectSubprocess)

	splitArgs, err := shell.Split(command)
	if err != nil {
		return nil, err
	}
	numArguments := len(splitArgs) - 1
	if numArguments < 0 {
		return nil, errors.New("gexpect: No command given to spawn")
	}
	path, err := exec.LookPath(splitArgs[0])
	if err != nil {
		return nil, err
	}

	if numArguments >= 1 {
		wrapper.cmd = exec.Command(path, splitArgs[1:]...)
	} else {
		wrapper.cmd = exec.Command(path)
	}

	stdout, err := wrapper.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	wrapper.cmd.Stderr = wrapper.cmd.Stdout

	stdin, err := wrapper.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	wrapper.stdin = &stdin

	wrapper.output = bufio.NewReader(stdout)

	return wrapper, nil
}
