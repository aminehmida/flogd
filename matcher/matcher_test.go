package matcher

import (
	"sync"
	"testing"
	"time"
)

func TestMatcherSingle(t *testing.T) {
	input := make(chan string)
	result := make(chan string)

	go Monitor(`Accepted publickey for root from ((?:[0-9]{1,3}\.){3}[0-9]{1,3})`, 1, 1, input, result, nil)
	input <- "Jul 22 15:36:11 debian-1 sshd[1871963]: Accepted publickey for root from 51.148.183.198 port 51558 ssh2: RSA SHA256:Nh9OTrMgt7pLwn80MMBuuPZEPdN8Ie3JHoo/zZ9kSeo"
	// input <- "Jul 22 15:36:11 debian-1 sshd[1871963]: Accepted publickey for root from 51.148.183.198 port 51558 ssh2: RSA SHA256:Nh9OTrMgt7pLwn80MMBuuPZEPdN8Ie3JHoo/zZ9kSeo"

	if r := <-result; r != "51.148.183.198" {
		t.Error("Matcher failed. Got: ", r)
	}

}

func TestMatcherSingle2(t *testing.T) {
	input := make(chan string)
	result := make(chan string)

	go Monitor("^(hello) world$", 2, 5, input, result, nil)
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

	go Monitor(`^(\d\d) world$`, 2, 10, input, result, &wg)
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

	go Monitor(`^\d\d world$`, 2, 10, input, result, nil)
	input <- "01 world"
	input <- "02 world"
	input <- "03 world"
	input <- "04 world"
	input <- "05 world"

	close(input)
	close(result)
	wg.Wait()

}
