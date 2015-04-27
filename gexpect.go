package gexpect

import (
	"bytes"
	"errors"
	shell "github.com/kballard/go-shellquote"
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type ExpectSubprocess struct {
	Cmd *exec.Cmd
	buf *buffer
}

type buffer struct {
	f *os.File
	b bytes.Buffer
}

func (buf *buffer) Read(chunk []byte) (int, error) {
	nread := 0
	if buf.b.Len() > 0 {
		n, err := buf.b.Read(chunk)
		if err != nil {
			return n, err
		}
		if n == len(chunk) {
			return n, nil
		}
		nread = n
	}
	fn, err := buf.f.Read(chunk[nread:])
	return fn + nread, err
}

func (buf *buffer) PutBack(chunk []byte) {
	if len(chunk) == 0 {
		return
	}
	if buf.b.Len() == 0 {
		buf.b.Write(chunk)
		return
	}
	d := make([]byte, 0, len(chunk)+buf.b.Len())
	d = append(d, chunk...)
	d = append(d, buf.b.Bytes()...)
	buf.b.Reset()
	buf.b.Write(d)
}

func SpawnAtDirectory(command string, directory string) (*ExpectSubprocess, error) {
	expect, err := _spawn(command)
	if err != nil {
		return nil, err
	}
	expect.Cmd.Dir = directory
	return _start(expect)
}

func Command(command string) (*ExpectSubprocess, error) {
	expect, err := _spawn(command)
	if err != nil {
		return nil, err
	}
	return expect, nil
}

func (expect *ExpectSubprocess) Start() error {
	_, err := _start(expect)
	return err
}

func Spawn(command string) (*ExpectSubprocess, error) {
	expect, err := _spawn(command)
	if err != nil {
		return nil, err
	}
	return _start(expect)
}

func (expect *ExpectSubprocess) Close() error {
	return expect.Cmd.Process.Kill()
}

func (expect *ExpectSubprocess) AsyncInteractChannels() (send chan string, receive chan string) {
	receive = make(chan string)
	send = make(chan string)

	go func() {
		for {
			str, err := expect.ReadLine()
			if err != nil {
				close(receive)
				return
			}
			receive <- str
		}
	}()

	go func() {
		for {
			select {
			case sendCommand, exists := <-send:
				{
					if !exists {
						return
					}
					err := expect.Send(sendCommand)
					if err != nil {
						receive <- "gexpect Error: " + err.Error()
						return
					}
				}
			}
		}
	}()

	return
}

// This is an unsound function. It shouldn't be trusted, as we're not using a stream based regex library.
// TODO: Find a regex stream library, plug it in, or develop my own for fun.
func (expect *ExpectSubprocess) ExpectRegex(regexSearchString string) (e error) {
	var size = len(regexSearchString)

	if size < 255 {
		size = 255
	}

	chunk := make([]byte, size)

	for {
		n, err := expect.buf.Read(chunk)

		if err != nil {
			return err
		}
		success, err := regexp.Match(regexSearchString, chunk[:n])
		if err != nil {
			return err
		}
		if success {
			return nil
		}
	}
}

func buildKMPTable(searchString string) []int {
	pos := 2
	cnd := 0
	length := len(searchString)

	var table []int
	if length < 2 {
		length = 2
	}

	table = make([]int, length)
	table[0] = -1
	table[1] = 0

	for pos < len(searchString) {
		if searchString[pos-1] == searchString[cnd] {
			cnd += 1
			table[pos] = cnd
			pos += 1
		} else if cnd > 0 {
			cnd = table[cnd]
		} else {
			table[pos] = 0
			pos += 1
		}
	}
	return table
}

func (expect *ExpectSubprocess) ExpectTimeout(searchString string, timeout time.Duration) (e error) {
	result := make(chan error)
	go func() {
		result <- expect.Expect(searchString)
	}()
	select {
	case e = <-result:
	case <-time.After(timeout):
		e = errors.New("Expect timed out.")
	}
	return e
}

func (expect *ExpectSubprocess) Expect(searchString string) (e error) {
	chunk := make([]byte, len(searchString)*2)
	target := len(searchString)
	m := 0
	i := 0
	// Build KMP Table
	table := buildKMPTable(searchString)

	for {
		n, err := expect.buf.Read(chunk)

		if err != nil {
			return err
		}
		offset := m + i
		for m+i-offset < n {
			if searchString[i] == chunk[m+i-offset] {
				i += 1
				if i == target {
					unreadIndex := m + i - offset
					if len(chunk) > unreadIndex {
						expect.buf.PutBack(chunk[unreadIndex:])
					}
					return nil
				}
			} else {
				m += i - table[i]
				if table[i] > -1 {
					i = table[i]
				} else {
					i = 0
				}
			}
		}
	}
}

func (expect *ExpectSubprocess) ExpectAny(searchStrings []string) (foundString string, e error) {
	// Get number of elements
	numElem := len(searchStrings)

	// Find the longest string
	// Create targets for each string
	// Initialize vars
	maxLen := 0
	table := make([][]int, numElem)
	target := make([]int, numElem)
	m := make([]int, numElem)
	i := make([]int, numElem)
	offset := make([]int, numElem)
	for elem, searchString := range searchStrings {
		// Find longest
		target[elem] = len(searchString)
		if target[elem] > maxLen {
			maxLen = target[elem]
		}
		// Create target
		table[elem] = buildKMPTable(searchString)
		// Initialize vars
		m[elem] = 0
		i[elem] = 0
	}

	// Make chunk depending on the lenght of the longest string
	chunk := make([]byte, maxLen*2)

	for {
		// Read chunk, check for erros
		n, err := expect.f.Read(chunk)
		if err != nil {
			return "", err
		}

		// Using same algorithm as normal Expect, only modified to work with an array of strings
		for elem, searchString := range searchStrings {
			offset[elem] = m[elem] + i[elem]
			for m[elem]+i[elem]-offset[elem] < n {
				if searchString[i[elem]] == chunk[m[elem]+i[elem]-offset[elem]] {
					i[elem] += 1
					if i[elem] == target[elem] {
						return searchString, nil
					}
				} else {
					m[elem] += i[elem] - table[elem][i[elem]]
					if table[elem][i[elem]] > -1 {
						i[elem] = table[elem][i[elem]]
					} else {
						i[elem] = 0
					}
				}
			}
		}
	}
}

func (expect *ExpectSubprocess) Send(command string) error {
	_, err := io.WriteString(expect.buf.f, command)
	return err
}

func (expect *ExpectSubprocess) SendLine(command string) error {
	_, err := io.WriteString(expect.buf.f, command+"\r\n")
	return err
}

func (expect *ExpectSubprocess) Interact() {
	defer expect.Cmd.Wait()
	io.Copy(os.Stdout, &expect.buf.b)
	go io.Copy(os.Stdout, expect.buf.f)
	go io.Copy(expect.buf.f, os.Stdin)
}

func (expect *ExpectSubprocess) ReadUntil(delim byte) ([]byte, error) {
	join := make([]byte, 1, 512)
	chunk := make([]byte, 255)

	for {

		n, err := expect.buf.Read(chunk)

		if err != nil {
			return join, err
		}

		for i := 0; i < n; i++ {
			if chunk[i] == delim {
				if len(chunk) > i+1 {
					expect.buf.PutBack(chunk[i+1:])
				}
				return join, nil
			} else {
				join = append(join, chunk[i])
			}
		}
	}
}

func (expect *ExpectSubprocess) Wait() error {
	return expect.Cmd.Wait()
}

func (expect *ExpectSubprocess) ReadLine() (string, error) {
	str, err := expect.ReadUntil('\n')
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func _start(expect *ExpectSubprocess) (*ExpectSubprocess, error) {
	f, err := pty.Start(expect.Cmd)
	if err != nil {
		return nil, err
	}
	expect.buf.f = f

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
		wrapper.Cmd = exec.Command(path, splitArgs[1:]...)
	} else {
		wrapper.Cmd = exec.Command(path)
	}
	wrapper.buf = new(buffer)

	return wrapper, nil
}
