package tailer

import (
	"testing"
)

func TestProcessTailerSingleLine(t *testing.T) {
	out := make(chan string)
	errors := make(chan error)
	go ProcessTailer("echo hello", out, errors)
	if <-out != "hello" {
		t.Error("ProcessTailer failed")
	}

	if err := <-errors; err != nil {
		t.Error("ProcessTailer failed:", err)
	}

}
func TestProcessTailerMultipleLinesFor(t *testing.T) {
	out := make(chan string)
	errors := make(chan error)
	expected := []string{"hello", "world"}

	go ProcessTailer("printf 'hello\nworld'", out, errors)

	for i := 0; i < len(expected); i++ {
		if <-out != expected[i] {
			t.Error("ProcessTailer failed")
		}
	}

	if err := <-errors; err != nil {
		t.Error("ProcessTailer failed:", err)
	}
}

func TestProcessTailerMultipleLinesRange(t *testing.T) {
	out := make(chan string)
	errors := make(chan error)
	go func() {
		for l := range out {
			if l != "test" {
				t.Error("ProcessTailer failed")
			}
		}
	}()

	go ProcessTailer("printf 'test\ntest\ntest'", out, errors)

	if err := <-errors; err != nil {
		t.Error("ProcessTailer failed:", err)
	}

}

func TestProcessTailerError(t *testing.T) {
	out := make(chan string)
	errors := make(chan error)
	go ProcessTailer("command-not-found", out, errors)
	for l := range out {
		t.Log("ProcessTailer error output:", l)
	}
	if err := <-errors; err == nil {
		t.Error("ProcessTailer failed")
	}
}
