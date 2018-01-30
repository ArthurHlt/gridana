package model

const defaultIdentifierTemplate = `{{ (.Labels.SortedPairs.Values) | join "/" }}`

type Route struct {
	Probe        string       `yaml:"probe,omitempty"`
	LabelMatcher LabelMatcher `yaml:"label_matcher"`
	Routes       []*Route     `yaml:"routes,omitempty"`
	Template     string       `yaml:"template"`
	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline" json:"-"`
}

func (r *Route) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Route
	if err := unmarshal((*plain)(r)); err != nil {
		return err
	}
	FillIdentifierTemplate(r, defaultIdentifierTemplate)
	return checkOverflow(r.XXX, "route")
}
func (r *Route) MatchAlert(alert Alert) bool {
	return r.LabelMatcher.MatchAlert(alert, true)
}
func (r *Route) FindRoute(alert Alert) *Route {
	finalRoute := r
	if !r.MatchAlert(alert) {
		return nil
	}
	if r.Routes == nil || len(r.Routes) == 0 {
		return finalRoute
	}
	for _, subRoute := range r.Routes {
		if !subRoute.MatchAlert(alert) {
			continue
		}
		finalRoute = subRoute.FindRoute(alert)
	}
	return finalRoute
}

func FillIdentifierTemplate(r *Route, currentTpl string) {
	if r.Template == "" {
		r.Template = currentTpl
	}
	for _, route := range r.Routes {
		FillIdentifierTemplate(route, r.Template)
	}
}
