package exec

import (
	"testing"
)

func TestExecuteSuccess(t *testing.T) {
	stdout, stderr, retCode, err := Execute("echo hello")
	if err != nil || stderr != "" || retCode != 0 {
		t.Error("Execute failed. Error:", err, "Return code:", retCode)
	}
	if stdout != "hello\n" {
		t.Error("Execute failed", stdout)
	}
}

func TestExecuteComposedOutput(t *testing.T) {
	stdout, stderr, retCode, err := Execute("printf hello; printf world")
	if err != nil || stderr != "" || retCode != 0 {
		t.Error("Execute failed. Error:", err, "Return code:", retCode, "Stderr:", stderr)
	}
	if stdout != "helloworld" {
		t.Error("Expected stdout to be 'helloworld', got", stdout)
	}
}

func TestExecuteStderr(t *testing.T) {
	stdout, stderr, retCode, err := Execute("printf hello >&2")
	if err != nil || stdout != "" || retCode != 0 {
		t.Error("Execute failed. Error:", err, "Return code:", retCode, "Stdout:", stdout)
	}
	if stderr != "hello" {
		t.Error("Expected stderr to be 'hello', got", stdout)
	}
}

func TestExecuteFail(t *testing.T) {
	_, _, retCode, _ := Execute("false")
	if retCode != 1 {
		t.Error("Expected return code 1, got", retCode)
	}
}
