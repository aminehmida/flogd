package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncodeConfig(t *testing.T) {
	c := Config{
		Name:        "test",
		SType:       "process",
		Description: "test",
		Regex:       "test",
		Do:          "test",
		Count:       1,
		Interval:    1,
		Command:     "test",
	}
	r := `- name: test
  type: process
  description: test
  regex: test
  do: test
  count: 1
  interval: 1
  command: test
`

	config := Configs{c}
	data, err := config.Encode()
	if err != nil {
		t.Errorf("Error encoding config: %v", err)
	}
	// os.WriteFile("test.yaml", data, 0644)
	if string(data) != r {
		t.Errorf("Expected: \n%s\nGot: \n%s", r, string(data))
	}
	// t.Logf("Encoded config: %s", data)
}

func TestDecodeConfig(t *testing.T) {
	c := Config{
		Name:        "test",
		SType:       "process",
		Description: "test",
		Regex:       "test",
		Do:          "test",
		Count:       1,
		Interval:    1,
		Command:     "test",
	}
	r := `- name: test
  type: process
  description: test
  regex: test
  do: test
  count: 1
  interval: 1
  command: test
`

	config := Configs{c}
	result := Configs{}
	err := result.Decode([]byte(r))
	if err != nil {
		t.Errorf("Error encoding config: %v", err)
	}
	t.Logf("Decoded config: %v", result)
	if !cmp.Equal(config, result) {
		t.Errorf("Config different from result: \n%v", cmp.Diff(config, result))
	}
}

func TestEncodeInvalidConfig(t *testing.T) {
	c := Config{
		Name:        "test",
		SType:       "process",
		Description: "test",
		Regex:       "test",
		Do:          "test",
		Count:       1,
		Interval:    1,
	}

	config := Configs{c}
	_, err := config.Encode()
	if err == nil || err.Error() != "config type set to \"process\" but \"command\" not defined" {
		t.Errorf("Expected error: config type set to \"process\" but \"command\" not defined, got: %v", err)
	}
}

func TestDecodeInvalidConfig(t *testing.T) {
	r := `- name: test
  type: process
  description: test
  regex: test
  do: test
  count: 1
  interval: 1
`
	result := Configs{}
	err := result.Decode([]byte(r))
	if err == nil || err.Error() != "config type set to \"process\" but \"command\" not defined" {
		t.Errorf("Expected error: config type set to \"process\" but \"command\" not defined, got: %v", err)
	}
}
