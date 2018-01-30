package alertmanager

import (
	"encoding/json"
	"fmt"
	"github.com/ArthurHlt/gridana/drivers"
	"github.com/ArthurHlt/gridana/model"
	"github.com/ArthurHlt/gridana/utils"
	promtpl "github.com/prometheus/alertmanager/template"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func init() {
	drivers.Register("alertmanager", func() drivers.Driver {
		return New()
	})
}

const (
	alertRoute     = "/api/v1/alerts"
	alertNameLabel = "alertname"
	NS_UUID        = "531f431c-70bf-4275-af81-f6571eeff06a"
	stateActive    = "active"
	stateSilenced  = "suppressed"
)

type Driver struct {
	client *http.Client
	amUrl  string
}

func New() *Driver {
	return &Driver{}
}

type AmData struct {
	Status string    `json:"status"`
	Data   []AmAlert `json:"data"`
}
type AmAlert struct {
	Labels       model.KV
	Annotations  model.KV
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Status       struct {
		State       string        `json:"state"`
		SilencedBy  []interface{} `json:"silencedBy"`
		InhibitedBy []interface{} `json:"inhibitedBy"`
	} `json:"status"`
	Receivers   []string `json:"receivers"`
	Fingerprint string   `json:"fingerprint"`
}

func (d Driver) toAlertFromRetrieve(amAlert AmAlert) model.Alert {
	labels := amAlert.Labels
	alert := model.Alert{
		Name:         string(labels[alertNameLabel]),
		GeneratorURL: amAlert.GeneratorURL,
		NotifierURL:  d.amUrl,
		StartsAt:     amAlert.StartsAt,
		EndsAt:       amAlert.EndsAt,
		Annotations:  amAlert.Annotations,
	}
	delete(labels, alertNameLabel)
	alert.Labels = labels
	alert.ID = d.generateAlertId(alert)
	if amAlert.Status.State == stateActive {
		alert.Status = model.Sfiring
	} else {
		alert.Status = model.Ssilenced
	}
	return alert
}
func (d Driver) toAlertFromReceived(amAlert promtpl.Alert) model.Alert {
	labels := amAlert.Labels
	alert := model.Alert{
		Name:         string(labels[alertNameLabel]),
		GeneratorURL: amAlert.GeneratorURL,
		NotifierURL:  d.amUrl,
		StartsAt:     amAlert.StartsAt,
		EndsAt:       amAlert.EndsAt,
		Annotations:  model.KV(amAlert.Annotations),
	}
	alert.ID = d.generateAlertId(alert)
	delete(labels, alertNameLabel)
	alert.Labels = model.KV(labels)

	if amAlert.Status == "firing" {
		alert.Status = model.Sfiring
	} else {
		alert.Status = model.Sresolved
	}
	return alert
}
func (d Driver) generateAlertId(alert model.Alert) string {
	names := make([]string, len(alert.Labels))
	for i, pair := range alert.Labels.SortedPairs() {
		names[i] = fmt.Sprintf("%s=%s", pair.Name, pair.Value)
	}
	return uuid.NewV3(uuid.FromStringOrNil(NS_UUID), strings.Join(names, "-")).String()
}
func (d Driver) RetrieveAlerts() ([]model.Alert, error) {
	alerts := make([]model.Alert, 0)
	resp, err := d.client.Get(d.amUrl + alertRoute)
	if err != nil {
		return alerts, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return alerts, err
	}
	if resp.StatusCode >= 300 {
		return alerts, fmt.Errorf("Error when retrieving alerts (status code %d): %s", resp.StatusCode, string(b))
	}

	var amAlerts AmData
	err = json.Unmarshal(b, &amAlerts)
	if err != nil {
		return alerts, err
	}
	for _, amAlert := range amAlerts.Data {
		if amAlert.Status.State != stateActive && amAlert.Status.State != stateSilenced {
			continue
		}
		alerts = append(alerts, d.toAlertFromRetrieve(amAlert))
	}
	return alerts, nil
}

func (d Driver) ReceiveAlerts(data []byte) ([]model.Alert, error) {
	alerts := make([]model.Alert, 0)
	var wAlerts promtpl.Data
	err := json.Unmarshal(data, &wAlerts)
	if err != nil {
		return alerts, err
	}
	for _, wAlert := range wAlerts.Alerts {
		alerts = append(alerts, d.toAlertFromReceived(wAlert))
	}
	return alerts, nil
}
func (d *Driver) Config(config model.DriverConfig) error {
	d.amUrl = strings.TrimSuffix(config.URL, "/")
	d.client = utils.CreateClient(config)
	return nil
}
