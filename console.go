package wrapper

import (
	"bufio"
	"fmt"
)

type Console interface {
	Start() error
	Kill() error
	WriteCmd(string) error
	ReadLine() (string, error)
}

type defaultConsole struct {
	cmd    JavaExec
	stdout *bufio.Reader
	stdin  *bufio.Writer
}

func newConsole(cmd JavaExec) *defaultConsole {
	c := &defaultConsole{
		cmd: cmd,
	}

	c.stdout = bufio.NewReader(cmd.Stdout())
	c.stdin = bufio.NewWriter(cmd.Stdin())
	return c
}

func (c *defaultConsole) Start() error {
	return c.cmd.Start()
}

func (c *defaultConsole) Kill() error {
	return c.cmd.Kill()
}

func (c *defaultConsole) WriteCmd(cmd string) error {
	wrappedCmd := fmt.Sprintf("%s\r\n", cmd)
	_, err := c.stdin.WriteString(wrappedCmd)
	if err != nil {
		return err
	}
	return c.stdin.Flush()
}

func (c *defaultConsole) ReadLine() (string, error) {
	return c.stdout.ReadString('\n')
}
