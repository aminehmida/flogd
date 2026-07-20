package tailer

import (
	"bufio"
	"fmt"
	"os/exec"

	shellwords "github.com/mattn/go-shellwords"
)

func ProcessTailer(command string, out chan<- string, errors chan<- error) {
	args, err := shellwords.Parse(command)
	if err != nil {
		close(out)
		errors <- fmt.Errorf("error while parsing command \"%s\": %v", command, err)
		close(errors)
		return
	}
	if len(args) == 0 {
		close(out)
		errors <- fmt.Errorf("empty command")
		close(errors)
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	r, err := cmd.StdoutPipe()
	if err != nil {
		close(out)
		errors <- fmt.Errorf("error while creating stdout pipe for \"%s\": %v", command, err)
		close(errors)
		return
	}
	cmd.Stderr = cmd.Stdout
	done := make(chan struct{})
	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
		}
		done <- struct{}{}
	}()

	// Start the command and check for errors
	err = cmd.Start()
	if err != nil {
		close(out)
		errors <- fmt.Errorf("error while starting command \"%s\": %v", command, err)
		close(errors)
		return
	}

	<-done

	err = cmd.Wait()
	if err != nil {
		close(out)
		errors <- fmt.Errorf("error while running \"%s\": %v", command, err)
		close(errors)
		return
	}

	errors <- nil
	close(out)
	close(errors)
}
