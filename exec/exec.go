package exec

import (
	"bytes"
	"os/exec"
)

func Execute(command string) (stdout, stderr string, retCode int, err error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	return out.String(), errOut.String(), cmd.ProcessState.ExitCode(), err
}
