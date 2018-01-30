package model

import (
	"fmt"
)

const defaultWeight = 10

type ColorByLabel struct {
	Color        string                 `yaml:"color"`
	Weight       int                    `yaml:"weight"`
	LabelMatcher LabelMatcher           `yaml:"label_matcher"`
	XXX          map[string]interface{} `yaml:",inline"`
}

func (c *ColorByLabel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain ColorByLabel
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.Color == "" {
		return fmt.Errorf("color_by_labels config: You must set a color")
	}
	if c.Weight < 0 {
		return fmt.Errorf("color_by_labels config: Weight must be a positive number")
	}
	return checkOverflow(c.XXX, "color_by_labels config")
}

type ColorByLabels []ColorByLabel

func (c ColorByLabels) AlertColor(alert Alert) (string, int) {
	for _, colorByLabel := range c {
		if colorByLabel.LabelMatcher.MatchAlert(alert) {
			return colorByLabel.Color, colorByLabel.Weight
		}
	}
	return "", defaultWeight
}
