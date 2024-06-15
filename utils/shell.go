package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

func shell(bin string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(bin, args...)
	cmd.Stdout = &out
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() != 0 {
				return "", err
			}
		}
	}

	return strings.TrimRight(strings.TrimSpace(out.String()), "\000"), nil
}
