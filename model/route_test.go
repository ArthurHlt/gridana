package model_test

import (
	. "github.com/ArthurHlt/gridana/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Route", func() {
	Context("FindRoute", func() {
		It("should return nil when no route match", func() {
			r := &Route{
				Probe: "a receiver",
				LabelMatcher: LabelMatcher{
					Match: map[string]string{
						"severity": "warning",
					},
				},
			}

			Expect(r.FindRoute(Alert{
				Labels: map[string]string{},
			})).To(BeNil())
		})
		It("should return the route when route have no matcher", func() {
			r := &Route{
				Probe: "a receiver",
			}

			Expect(r.FindRoute(Alert{
				Labels: map[string]string{},
			})).To(Equal(r))
		})

		It("should return the first route matching", func() {
			r := &Route{
				Probe: "a receiver",
				Routes: []*Route{
					&Route{
						LabelMatcher: LabelMatcher{
							Match: map[string]string{
								"label": "value",
							},
						},
					},
					&Route{
						LabelMatcher: LabelMatcher{
							Match: map[string]string{
								"label": "value2",
							},
						},
					},
				},
			}

			Expect(r.FindRoute(Alert{
				Labels: map[string]string{"label": "value2"},
			})).To(Equal(r.Routes[1]))
		})

		Context("With sub routes", func() {
			It("should return the first route matching", func() {
				r := &Route{
					Probe: "a receiver",
					Routes: []*Route{
						&Route{
							LabelMatcher: LabelMatcher{
								Match: map[string]string{
									"label": "value",
								},
							},
						},
						&Route{
							LabelMatcher: LabelMatcher{
								Match: map[string]string{
									"label": "value2",
								},
							},
							Routes: []*Route{
								&Route{
									LabelMatcher: LabelMatcher{
										Match: map[string]string{
											"sub_label": "value",
										},
									},
								},
								&Route{
									LabelMatcher: LabelMatcher{
										Match: map[string]string{
											"sub_label": "value2",
										},
									},
								},
							},
						},
					},
				}

				Expect(r.FindRoute(Alert{
					Labels: map[string]string{
						"label":     "value2",
						"sub_label": "value",
					},
				})).To(Equal(r.Routes[1].Routes[0]))
			})
		})
	})
})
