package wrapper

import (
	"fmt"
	"io"
	"os/exec"
)

type JavaExec interface {
	Stdout() io.ReadCloser
	Stdin() io.WriteCloser
	Start() error
	Kill() error
}

type defaultJavaExec struct {
	cmd *exec.Cmd
}

func (j *defaultJavaExec) Stdout() io.ReadCloser {
	r, _ := j.cmd.StdoutPipe()
	return r
}

func (j *defaultJavaExec) Stdin() io.WriteCloser {
	w, _ := j.cmd.StdinPipe()
	return w
}

func (j *defaultJavaExec) Start() error {
	return j.cmd.Start()
}

func (j *defaultJavaExec) Kill() error {
	return j.cmd.Process.Kill()
}

func javaExecCmd(serverPath string, initialHeapSize, maxHeapSize int) *defaultJavaExec {
	initialHeapFlag := fmt.Sprintf("-Xms%dM", initialHeapSize)
	maxHeapFlag := fmt.Sprintf("-Xmx%dM", maxHeapSize)
	cmd := exec.Command("java", initialHeapFlag, maxHeapFlag, "-jar", serverPath, "nogui")
	return &defaultJavaExec{cmd: cmd}
}
