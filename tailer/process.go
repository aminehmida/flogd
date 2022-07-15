package tailer

import (
	"bufio"
	"os/exec"

	shellwords "github.com/mattn/go-shellwords"
)

func ProcessTailer(command string, out chan<- string, errors chan<- error) {
	args, err := shellwords.Parse(command)
	if err != nil {
		errors <- err
		// close(out)
	}

	// fmt.Printf("args: %v\n", args)

	cmd := exec.Command(args[0], args[1:]...)
	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	done := make(chan struct{})
	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			out <- line
		}
		done <- struct{}{}
		close(out)
	}()

	// Start the command and check for errors
	err = cmd.Start()
	if err != nil {
		close(out)
		errors <- err
		return
	}

	<-done

	err = cmd.Wait()
	if err != nil {
		close(out)
		errors <- err
		return
	}

	errors <- nil

}
