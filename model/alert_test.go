package model_test

import (
	. "github.com/ArthurHlt/gridana/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
	"time"
)

func timeMustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

var _ = Describe("FormattedAlerts", func() {
	Context("OrderByStartAt", func() {
		It("should order alerts by their time start at", func() {
			alerts := FormattedAlerts{
				{
					Alert: Alert{
						Name:     "alert1",
						StartsAt: timeMustParse("15", "10"),
					},
				},
				{
					Alert: Alert{
						Name:     "alert2",
						StartsAt: timeMustParse("15", "5"),
					},
				},
				{
					Alert: Alert{
						Name:     "alert3",
						StartsAt: timeMustParse("15", "1"),
					},
				},
				{
					Alert: Alert{
						Name:     "alert4",
						StartsAt: timeMustParse("15", "11"),
					},
				},
			}
			sort.Sort(ByStartAt{alerts})
			Expect(alerts[0].Name).To(Equal("alert3"))
			Expect(alerts[1].Name).To(Equal("alert2"))
			Expect(alerts[2].Name).To(Equal("alert1"))
			Expect(alerts[3].Name).To(Equal("alert4"))
		})
	})
	Context("Ordered", func() {
		It("should give a map of alerts ordered by identifiers and probes sorted by weight", func() {
			alerts := FormattedAlerts{
				{
					Alert: Alert{
						Name: "alert1",
					},
					Identifier: "identifier1",
					Probe:      "probe1",
					Weight:     2,
				},
				{
					Alert: Alert{
						Name: "alert1.1",
					},
					Identifier: "identifier1",
					Probe:      "probe1",
					Weight:     5,
				},
				{
					Alert: Alert{
						Name: "alert1.2",
					},
					Identifier: "identifier1",
					Probe:      "probe1",
					Weight:     3,
				},
				{
					Alert: Alert{
						Name: "alert1.3",
					},
					Identifier: "identifier1",
					Probe:      "probe2",
					Weight:     1,
				},
				{
					Alert: Alert{
						Name: "alert2",
					},
					Identifier: "identifier2",
					Probe:      "probe1",
					Weight:     1,
				},
				{
					Alert: Alert{
						Name: "alert2.1",
					},
					Identifier: "identifier2",
					Probe:      "probe2",
					Weight:     4,
				},
			}
			aOrdered := alerts.Ordered(
				Probes{
					{
						Name: "probe1",
					},
					{
						Name: "probe2",
					},
					{
						Name: "probe3",
					},
				},
			)
			Expect(aOrdered).To(Equal(OrderedAlerts{
				"identifier1": AlertsByProbe{
					"probe1": []FormattedAlert{
						{
							Alert: Alert{
								Name: "alert1.1",
							},
							Identifier: "identifier1",
							Probe:      "probe1",
							Weight:     5,
						},
						{
							Alert: Alert{
								Name: "alert1.2",
							},
							Identifier: "identifier1",
							Probe:      "probe1",
							Weight:     3,
						},
						{
							Alert: Alert{
								Name: "alert1",
							},
							Identifier: "identifier1",
							Probe:      "probe1",
							Weight:     2,
						},
					},
					"probe2": []FormattedAlert{
						{
							Alert: Alert{
								Name: "alert1.3",
							},
							Identifier: "identifier1",
							Probe:      "probe2",
							Weight:     1,
						},
					},
					"probe3": []FormattedAlert{},
				},
				"identifier2": AlertsByProbe{
					"probe1": []FormattedAlert{
						{
							Alert: Alert{
								Name: "alert2",
							},
							Identifier: "identifier2",
							Probe:      "probe1",
							Weight:     1,
						},
					},
					"probe2": []FormattedAlert{
						{
							Alert: Alert{
								Name: "alert2.1",
							},
							Identifier: "identifier2",
							Probe:      "probe2",
							Weight:     4,
						},
					},
					"probe3": []FormattedAlert{},
				},
			}))
			Expect(aOrdered["identifier1"]["probe1"][0].Name).To(Equal("alert1.1"))
			Expect(aOrdered["identifier1"]["probe2"][0].Name).To(Equal("alert1.3"))
			Expect(aOrdered["identifier2"]["probe1"][0].Name).To(Equal("alert2"))
			Expect(aOrdered["identifier2"]["probe2"][0].Name).To(Equal("alert2.1"))
		})

	})
})
