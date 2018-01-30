package model_test

import (
	. "github.com/ArthurHlt/gridana/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ColorByLabels", func() {
	Context("AlertColor", func() {
		cByLabels := ColorByLabels{
			{
				Color: "yellow",
				LabelMatcher: LabelMatcher{
					Match: map[string]string{
						"severity": "warning",
					},
				},
			},
		}
		It("should give empty color if alert not matching", func() {
			color, weight := cByLabels.AlertColor(Alert{
				Labels: map[string]string{},
			})
			Expect(color).To(Equal(""))
			Expect(weight).To(Equal(10))
		})
		It("should give color if alert matching", func() {
			color, weight := cByLabels.AlertColor(Alert{
				Labels: map[string]string{"severity": "warning"},
			})
			Expect(color).To(Equal("yellow"))
			Expect(weight).To(Equal(0))
		})
	})
})
