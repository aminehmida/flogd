package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name        string `yaml:"name"`
	SType       string `yaml:"type"`
	Description string `yaml:"description"`
	Regex       string `yaml:"regex"`
	Do          string `yaml:"do"`
	Count       int    `yaml:"count"`
	Interval    int    `yaml:"interval"`
	Command     string `yaml:"command"`
}

type Configs []Config

func (c *Configs) Decode(data []byte) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return c.validate()
}

func (c *Configs) Encode() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return yaml.Marshal(c)
}

func (c *Configs) validate() error {
	for _, config := range *c {
		if config.Name == "" {
			return fmt.Errorf("config \"name\" not defined")
		}
		if config.SType == "" {
			return fmt.Errorf("config \"type\" not defined")
		}
		if config.Regex == "" {
			return fmt.Errorf("config \"regex\" not defined")
		}
		if config.Do == "" {
			return fmt.Errorf("config \"do\" not defined")
		}
		if config.Count == 0 {
			return fmt.Errorf("config \"count\" not defined or set to 0")
		}
		if config.SType == "process" && config.Command == "" {
			return fmt.Errorf("config type set to \"process\" but \"command\" not defined")
		}
	}
	return nil
}
