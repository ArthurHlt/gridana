package model

type Logs struct {
	Level      string `yaml:"level"`
	InJson     bool   `yaml:"in_json"`
	NoColor    bool   `yaml:"no_color"`
	SyslogHost string `yaml:"syslog_host"`
	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

func (c *Logs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Logs
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return checkOverflow(c.XXX, "logs config")
}
