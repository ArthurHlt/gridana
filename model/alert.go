package model

import (
	"sort"
	"time"
)

const (
	Sfiring   alertStatus = "firing"
	Sresolved alertStatus = "resolved"
	Ssilenced alertStatus = "silenced"
)

type alertStatus string

type FormattedAlerts []FormattedAlert

type ByStartAt struct{ FormattedAlerts }

func (a ByStartAt) Less(i, j int) bool {
	return a.FormattedAlerts[i].StartsAt.Before(a.FormattedAlerts[j].StartsAt)
}

type Alert struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
	NotifierURL  string      `json:"notifierURL"`
	Status       alertStatus `json:"status"`
	Labels       KV          `json:"labels"`
	Annotations  KV          `json:"annotations"`
	Silence      Silence     `json:"silence"`
}

type Silence struct {
	ID        string    `json:"id"`
	CreatedBy string    `json:"createdBy"`
	Reason    string    `json:"reason"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
}

type AlertsByProbe map[string][]FormattedAlert

type OrderedAlerts map[string]AlertsByProbe

type FormattedAlert struct {
	*Alert
	Color        string `json:"color"`
	Message      string `json:"message"`
	Probe        string `json:"probe"`
	Identifier   string `json:"identifier"`
	Notification string `json:"notification"`
	Weight       int    `json:"weight"`
}

func (a FormattedAlerts) Len() int {
	return len(a)
}
func (a FormattedAlerts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a FormattedAlerts) Less(i, j int) bool {
	if a[i].Weight == a[j].Weight {
		return a[i].StartsAt.Before(a[j].StartsAt)
	}
	return a[i].Weight > a[j].Weight
}

func (a FormattedAlerts) Ordered(probes Probes) OrderedAlerts {
	orderedAlerts := make(OrderedAlerts)

	orderTemp := make(map[string]FormattedAlerts)
	for _, alert := range a {
		if _, ok := orderTemp[alert.Identifier]; !ok {
			orderTemp[alert.Identifier] = FormattedAlerts{alert}
			continue
		}
		tmpAlerts := orderTemp[alert.Identifier]
		tmpAlerts = append(tmpAlerts, alert)
		orderTemp[alert.Identifier] = tmpAlerts
	}

	for key, alerts := range orderTemp {
		sort.Sort(alerts)
		aByProbes := make(AlertsByProbe)
		for _, probe := range probes {
			aByProbes[probe.Name] = make([]FormattedAlert, 0)
		}
		for _, alert := range alerts {
			tmpAlerts := aByProbes[alert.Probe]
			tmpAlerts = append(tmpAlerts, alert)
			aByProbes[alert.Probe] = tmpAlerts
		}
		orderedAlerts[key] = aByProbes
	}
	return orderedAlerts
}
