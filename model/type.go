package model

import (
	"regexp"
	"sort"
)

type Regexp struct {
	*regexp.Regexp
}

func (re *Regexp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	regex, err := regexp.Compile("^(?:" + s + ")$")
	if err != nil {
		return err
	}
	re.Regexp = regex
	return nil
}

type Regexps []Regexp

func (r Regexps) MatchString(s string) bool {
	for _, match := range r {
		if match.MatchString(s) {
			return true
		}
	}
	return false
}

type LabelMatcher struct {
	RequireAll bool               `yaml:"require_all"`
	MatchRE    map[string]Regexps `yaml:"match_re"`
	Match      map[string]string  `yaml:"match"`
}

func (lm LabelMatcher) MatchAlert(alert Alert, defaultMatch ...bool) bool {
	labels := alert.Labels
	labels["__name__"] = alert.Name
	defer func() {
		delete(labels, "__name__")
	}()
	if len(defaultMatch) > 0 && len(lm.Match) == 0 && len(lm.MatchRE) == 0 {
		return defaultMatch[0]
	}
	match := false
	for labelName, matcher := range lm.MatchRE {
		if _, ok := labels[labelName]; ok {
			match = matcher.MatchString(labels[labelName])
		}
		if !match && lm.RequireAll {
			return false
		}
		if match && !lm.RequireAll {
			return true
		}
	}
	for labelName, value := range lm.Match {
		if _, ok := labels[labelName]; ok {
			match = (labels[labelName] == value)
		}
		if !match && lm.RequireAll {
			return false
		}
		if match && !lm.RequireAll {
			return true
		}
	}

	return match
}

// KV is a set of key/value string pairs.
type KV map[string]string

// SortedPairs returns a sorted list of key/value pairs.
func (kv KV) SortedPairs() Pairs {
	var (
		pairs     = make([]Pair, 0, len(kv))
		keys      = make([]string, 0, len(kv))
		sortStart = 0
	)
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys[sortStart:])

	for _, k := range keys {
		pairs = append(pairs, Pair{k, kv[k]})
	}
	return pairs
}

// Remove returns a copy of the key/value set without the given keys.
func (kv KV) Remove(keys []string) KV {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	res := KV{}
	for k, v := range kv {
		if _, ok := keySet[k]; !ok {
			res[k] = v
		}
	}
	return res
}

// Pair is a key/value string pair.
type Pair struct {
	Name, Value string
}

// Pairs is a list of key/value string pairs.
type Pairs []Pair

// Names returns a list of names of the pairs.
func (ps Pairs) Names() []string {
	ns := make([]string, 0, len(ps))
	for _, p := range ps {
		ns = append(ns, p.Name)
	}
	return ns
}

// Values returns a list of values of the pairs.
func (ps Pairs) Values() []string {
	vs := make([]string, 0, len(ps))
	for _, p := range ps {
		vs = append(vs, p.Value)
	}
	return vs
}
