package main

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
)

const ()

func joinStrings(strs ...string) string {
	joined := strings.Builder{}
	for _, s := range strs {
		joined.WriteString(s)
	}
	return joined.String()
}

func generateJavaRunCmd(serverPath string, iniHeapSize, maxHeapSize int) *exec.Cmd {
	iniHeapFlag := joinStrings("-", "Xms", strconv.Itoa(iniHeapSize), "M")
	maxHeapFlag := joinStrings("-", "Xmx", strconv.Itoa(maxHeapSize), "M")
	return exec.Command("java", iniHeapFlag, maxHeapFlag, "-jar", serverPath, "nogui")
}

func main() {
	execCmd := generateJavaRunCmd("server/server.jar", 1024, 2024)
	console := NewConsole(execCmd)
	msw := NewMSW(console)

	msw.Start(context.Background())

	for {
	}
}
