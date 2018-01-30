package model_test

import (
	. "github.com/ArthurHlt/gridana/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"regexp"
)

var _ = Describe("Type", func() {
	Context("LabelMatcher", func() {
		Context("MatchAlert", func() {
			It("should return false by default if there is no matches filled", func() {
				lm := LabelMatcher{
					MatchRE: map[string]Regexps{},
				}
				alert := Alert{
					Labels: KV{
						"label": "value",
					},
				}
				Expect(lm.MatchAlert(alert)).To(BeFalse())
			})
			It("should return the default value set by user if there is no matches filled", func() {
				lm := LabelMatcher{
					MatchRE: map[string]Regexps{},
				}
				alert := Alert{
					Labels: KV{
						"label": "value",
					},
				}
				Expect(lm.MatchAlert(alert, true)).To(BeTrue())
			})
			It("should return true when one label match and require_all is false", func() {
				lm := LabelMatcher{
					MatchRE: map[string]Regexps{
						"label": Regexps{
							{regexp.MustCompile("^(?:avalue)$")},
						},
						"label2": Regexps{
							{regexp.MustCompile("^(?:avalue)$")},
						},
					},
				}
				alert := Alert{
					Labels: KV{
						"label2": "avalue",
						"label":  "value",
					},
				}
				Expect(lm.MatchAlert(alert)).To(BeTrue())
			})
			It("should return false when just one label match and require_all is true", func() {
				lm := LabelMatcher{
					RequireAll: true,
					MatchRE: map[string]Regexps{
						"label": Regexps{
							{regexp.MustCompile("^(?:avalue)$")},
						},
						"label2": Regexps{
							{regexp.MustCompile("^(?:avalue)$")},
						},
					},
				}
				alert := Alert{
					Labels: KV{
						"label2": "avalue",
						"label":  "value",
					},
				}
				Expect(lm.MatchAlert(alert)).To(BeFalse())
			})
			It("should return true when all label match and require_all is true", func() {
				lm := LabelMatcher{
					RequireAll: true,
					MatchRE: map[string]Regexps{
						"label": Regexps{
							{regexp.MustCompile("^(?:avalue)$")},
						},
						"label2": Regexps{
							{regexp.MustCompile("^(?:avalue)$")},
						},
					},
				}
				alert := Alert{
					Labels: KV{
						"label2": "avalue",
						"label":  "avalue",
					},
				}
				Expect(lm.MatchAlert(alert)).To(BeTrue())
			})
			It("should return true when alert name match and require_all is false", func() {
				lm := LabelMatcher{
					MatchRE: map[string]Regexps{
						"__name__": Regexps{
							{regexp.MustCompile("^(?:myname)$")},
						},
						"label2": Regexps{
							{regexp.MustCompile("^(?:avalue)$")},
						},
					},
				}
				alert := Alert{
					Name:   "myname",
					Labels: KV{},
				}
				Expect(lm.MatchAlert(alert)).To(BeTrue())
			})
		})
	})
})
