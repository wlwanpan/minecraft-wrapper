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

type JavaExecImpl struct {
	cmd *exec.Cmd
}

func (j *JavaExecImpl) Stdout() io.ReadCloser {
	r, _ := j.cmd.StdoutPipe()
	return r
}

func (j *JavaExecImpl) Stdin() io.WriteCloser {
	w, _ := j.cmd.StdinPipe()
	return w
}

func (j *JavaExecImpl) Start() error {
	return j.cmd.Start()
}

func (j *JavaExecImpl) Kill() error {
	return j.cmd.Process.Kill()
}

func JavaExecCmd(serverPath string, initialHeapSize, maxHeapSize int) *JavaExecImpl {
	initialHeapFlag := fmt.Sprintf("-Xms%dM", initialHeapSize)
	maxHeapFlag := fmt.Sprintf("-Xmx%dM", maxHeapSize)
	cmd := exec.Command("java", initialHeapFlag, maxHeapFlag, "-jar", serverPath, "nogui")
	return &JavaExecImpl{cmd: cmd}
}
