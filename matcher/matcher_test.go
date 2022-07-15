package matcher

import (
	"sync"
	"testing"
	"time"
)

func TestMatcherSingle(t *testing.T) {
	input := make(chan string)
	result := make(chan string)

	go Monitor("^(hello) world$", 2, 5, input, result)
	input <- "hello world"
	time.Sleep(1 * time.Second)
	input <- "hello world"

	if <-result != "hello" {
		t.Error("Matcher failed")
	}

}

func TestMatcherMultiple(t *testing.T) {
	input := make(chan string)
	result := make(chan string)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		if <-result != "01" {
			t.Error("Matcher failed")
		}
		if <-result != "02" {
			t.Error("Matcher failed")
		}
	}()

	go Monitor(`^(\d\d) world$`, 2, 10, input, result)
	input <- "01 world"
	input <- "02 world"
	input <- "01 world"
	input <- "02 world"
	input <- "03 world"

	close(input)
	close(result)
	wg.Wait()

}

func TestMatcherNoGroup(t *testing.T) {
	input := make(chan string)
	result := make(chan string)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		if <-result != "02 world" {
			t.Error("Matcher failed")
		}
		if <-result != "04 world" {
			t.Error("Matcher failed")
		}
	}()

	go Monitor(`^\d\d world$`, 2, 10, input, result)
	input <- "01 world"
	input <- "02 world"
	input <- "03 world"
	input <- "04 world"
	input <- "05 world"

	close(input)
	close(result)
	wg.Wait()

}
