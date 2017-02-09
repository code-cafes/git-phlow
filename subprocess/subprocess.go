package subprocess

import (
	"os/exec"
	"bytes"
	"errors"
	"syscall"
	"log"
)

const (
	EmptyReturnString string = ""
)

type ExecError struct {
	err      error
	stdErr   string
	exitCode int
}

func (e ExecError) Error() string {
	return e.stdErr
}

//SimpleExec
//Executes a cmd on your operating system
func SimpleExec(name string, args ...string) (string, error) {

	cmd := exec.Command(name, args...)

	var stdOutBuffer, stdErrBuffer bytes.Buffer

	cmd.Stderr = &stdErrBuffer
	cmd.Stdout = &stdOutBuffer

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %d", err)

		//Return Error with stderr, error - and exit status 1
		return EmptyReturnString, ExecError{err, stdErrBuffer.String(), 1}
	}

	if err := cmd.Wait(); err != nil {

		if exitErr, ok := err.(*exec.ExitError); ok {
			//Program exited with exit code != 0

			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				//Trying to obtain exit error from failed command

				log.Printf("Exit Status: %d", status.ExitStatus())
				return EmptyReturnString, ExecError{err, stdErrBuffer.String(), status.ExitStatus()}
			}
		}

		//Return Error with stderr, error - and exit status 1
		return EmptyReturnString, ExecError{err, stdErrBuffer.String(), 1}
	}

	//If no errors are returned, return stdout
	return stdOutBuffer.String(), nil
}

//IsInPath
//Checks if an app has been added to $PATH
func IsInPath(application string) (error) {
	_, err := exec.LookPath(application)

	if err != nil {
		return errors.New(application + " is not in $PATH")

	}
	return nil
}
