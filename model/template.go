package model

type Template struct {
	Message      string `yaml:"message,omitempty"`
	Notification string `yaml:"notification"`
	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline" json:"-"`
}
