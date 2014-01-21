package gexpect

import (
	"errors"
	"fmt"
	shell "github.com/kballard/go-shellquote"
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
	"regexp"
)

type ExpectSubprocess struct {
	cmd *exec.Cmd
	f   *os.File
}

func SpawnAtDirectory(command string, directory string) (*ExpectSubprocess, error) {
	expect, err := _spawn(command)
	if err != nil {
		return nil, err
	}
	expect.cmd.Dir = directory
	return _start(expect)
}

func Spawn(command string) (*ExpectSubprocess, error) {
	expect, err := _spawn(command)
	if err != nil {
		return nil, err
	}
	return _start(expect)
}

func (expect *ExpectSubprocess) Close() error {
	return expect.cmd.Process.Kill()
}

func (expect *ExpectSubprocess) AsyncInteractBiChannel() chan string {

	ch := make(chan string)
	readChan := make(chan string)

	go func() {
		defer func() {
			r := recover()
			if r != nil {
				ch <- fmt.Sprintf("Unknown error: %+#v", r)
			}
		}()
		for {
			str, err := expect.ReadLine()
			if err != nil {
				close(readChan)
				return
			}
			readChan <- str
		}
	}()

	go func() {
		defer func() {
			r := recover()
			if r != nil {
				ch <- fmt.Sprintf("Unknown error: %+#v", r)
			}
		}()
		for {
			select {
			case sendCommand, exists := <-ch:
				{
					if !exists {
						return
					}
					err := expect.Sendline(sendCommand)
					if err != nil {
						close(ch)
						return
					}
				}
			case output, exists := <-readChan:
				{
					if !exists {
						close(ch)
						return
					}
					ch <- output
				}
			}
		}
	}()

	return ch
}

func (expect *ExpectSubprocess) ExpectRegex(regexSearchString string) (e error) {
	defer func() {
		r := recover()
		if r != nil {
			e = errors.New(fmt.Sprintf("Unknown error: %+#v", r))
		}
	}()
	var size = len(regexSearchString)

	if size < 255 {
		size = 255
	}

	chunk := make([]byte, size)

	for {
		n, err := expect.f.Read(chunk)

		if err != nil {
			return err
		}

		// fmt.Printf("%d: %s\n", n, string(chunk))
		success, err := regexp.Match(regexSearchString, chunk[:n])
		if err != nil {
			return err
		}
		if success {
			return nil
		}
	}
}

func (expect *ExpectSubprocess) Expect(searchString string) (e error) {
	defer func() {
		r := recover()
		if r != nil {
			e = errors.New(fmt.Sprintf("Unknown error: %+#v", r))
		}
	}()
	chunk := make([]byte, len(searchString))
	found := 0
	target := len(searchString)
	for {

		n, err := expect.f.Read(chunk)

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
				if found > 0 {
					// TODO: Calculate jump value
					i = i - found + 1
				}
				found = 0
			}
		}
	}
}

func (expect *ExpectSubprocess) Sendline(command string) error {
	_, err := io.WriteString(expect.f, command+"\r\n")
	return err
}

func (expect *ExpectSubprocess) Interact() {
	defer expect.cmd.Wait()
	// go io.Copy(os.Stdout, os.Stdin)
	go io.Copy(os.Stdout, expect.f)
	go io.Copy(os.Stderr, expect.f)
	go io.Copy(expect.f, os.Stdin)
}

func (expect *ExpectSubprocess) ReadUntil(delim byte) ([]byte, error) {
	join := make([]byte, 1, 512)
	chunk := make([]byte, 255)

	for {

		n, err := expect.f.Read(chunk)

		if err != nil {
			return join, err
		}

		for i := 0; i < n; i++ {
			if chunk[i] == delim {
				return join, nil
			} else {
				join = append(join, chunk[i])
			}
		}
	}
}

func (expect *ExpectSubprocess) ReadLine() (string, error) {
	str, err := expect.ReadUntil('\n')
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func _start(expect *ExpectSubprocess) (*ExpectSubprocess, error) {
	f, err := pty.Start(expect.cmd)
	if err != nil {
		return nil, err
	}
	expect.f = f

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

	// wrapper.cmd.SysProcAttr.
	// wrapper.cmd.Stdout = wrapper.cmd.Stderr
	// go io.Copy(os.Stdout, stdout)

	return wrapper, nil
}
