package model

import "fmt"

type DriverConfig struct {
	Name               string `yaml:"name" json:"name"`
	URL                string `yaml:"url" json:"url"`
	User               string `yaml:"user" json:"user"`
	Password           string `yaml:"password" json:"password"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" json:"insecure_skip_verify"`
	Type               string `yaml:"type" json:"type"`
	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

func (c *DriverConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain DriverConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.Name == "" {
		return fmt.Errorf("driver config: you must set a name to your driver")
	}
	if c.Type == "" {
		return fmt.Errorf("driver config: you must set a type to your driver")
	}
	return checkOverflow(c.XXX, "driver config")
}
