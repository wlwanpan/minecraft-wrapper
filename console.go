package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

const (
	NewlineByte byte = '\n'
)

type Console struct {
	cmd    *exec.Cmd
	stdout *bufio.Reader
	stdin  *bufio.Writer
}

func NewConsole(cmd *exec.Cmd) *Console {
	return &Console{
		cmd: cmd,
	}
}

func (c *Console) Start() error {
	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	c.stdout = bufio.NewReader(stdout)

	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}
	c.stdin = bufio.NewWriter(stdin)

	return c.cmd.Start()
}

func (c *Console) Kill() error {
	return c.cmd.Process.Kill()
}

func (c *Console) Write(cmd string) error {
	wrappedCmd := fmt.Sprintf("%s\r\n", cmd)
	_, err := c.stdin.WriteString(wrappedCmd)
	if err != nil {
		return err
	}
	return c.stdin.Flush()
}

func (c *Console) Read() (string, error) {
	return c.stdout.ReadString(NewlineByte)
}
