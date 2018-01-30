package model

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultSilenceColor = "blue"
	defaultColor        = "red"
)

type GridanaConfig struct {
	Route         *Route         `yaml:"route,omitempty"`
	Probes        Probes         `yaml:"probes,omitempty"`
	DropLabels    LabelMatcher   `yaml:"drop_labels"`
	ColorByLabels ColorByLabels  `yaml:"color_by_labels"`
	Logs          Logs           `yaml:"logs"`
	SilenceColor  string         `yaml:"silence_color"`
	DefaultColor  string         `yaml:"default_color"`
	ListenAddr    string         `yaml:"listen_addr"`
	Drivers       []DriverConfig `yaml:"drivers"`
	NoCheckOrigin bool           `yaml:"no_check_origin"`
	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline" json:"-"`
}

func (c *GridanaConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain GridanaConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if c.ListenAddr == "" {
		c.ListenAddr = "0.0.0.0:" + port
	}
	if len(strings.Split(c.ListenAddr, ":")) < 2 {
		c.ListenAddr += ":" + port
	}
	if c.SilenceColor == "" {
		c.SilenceColor = defaultSilenceColor
	}
	if c.DefaultColor == "" {
		c.DefaultColor = defaultColor
	}
	names := map[string]struct{}{}
	for _, rcv := range c.Probes {
		if _, ok := names[rcv.Name]; ok {
			return fmt.Errorf("config: probe name %q is not unique", rcv.Name)
		}
		names[rcv.Name] = struct{}{}
	}
	if len(c.Drivers) == 0 {
		return fmt.Errorf("config: you must set at least one driver")
	}
	names = map[string]struct{}{}
	for _, driver := range c.Drivers {
		if _, ok := names[driver.Name]; ok {
			return fmt.Errorf("config: driver name %q is not unique", driver.Name)
		}
		names[driver.Name] = struct{}{}
	}

	return checkOverflow(c.XXX, "config")
}
