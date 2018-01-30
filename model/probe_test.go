package model_test

import (
	. "github.com/ArthurHlt/gridana/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Probes", func() {
	Context("FindProbeByAlert", func() {
		It("should return nil when no route match", func() {
			r := &Route{
				Probe: "a receiver",
				LabelMatcher: LabelMatcher{
					Match: map[string]string{
						"severity": "warning",
					},
				},
			}
			rec := Probes{
				{
					Name: "a receiver",
				},
			}

			Expect(rec.FindProbeByAlert(Alert{
				Labels: map[string]string{},
			}, r)).To(BeNil())
		})
		It("should return nil when no receiver not exists", func() {
			r := &Route{
				Probe: "wrong receiver",
				LabelMatcher: LabelMatcher{
					Match: map[string]string{
						"severity": "warning",
					},
				},
			}
			rec := Probes{
				{
					Name: "a receiver",
				},
			}

			Expect(rec.FindProbeByAlert(Alert{
				Labels: map[string]string{"severity": "warning"},
			}, r)).To(BeNil())
		})
		It("should return receiver found in route when route match", func() {
			r := &Route{
				Probe: "a receiver",
				LabelMatcher: LabelMatcher{
					Match: map[string]string{
						"severity": "warning",
					},
				},
			}
			rec := Probes{
				{
					Name: "a receiver",
				},
			}

			Expect(rec.FindProbeByAlert(Alert{
				Labels: map[string]string{"severity": "warning"},
			}, r)).To(Equal(rec[0]))
		})
	})
})
