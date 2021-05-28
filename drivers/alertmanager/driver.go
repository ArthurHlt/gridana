package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ArthurHlt/gridana/drivers"
	"github.com/ArthurHlt/gridana/model"
	"github.com/ArthurHlt/gridana/utils"
	promtpl "github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/alertmanager/types"
	promodel "github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

func init() {
	drivers.Register("alertmanager", func() drivers.Driver {
		return New()
	})
}

const (
	alertRoute     = "/api/v1/alerts"
	silenceRoute   = "/api/v1/silence/%s"
	silencesRoute  = "/api/v1/silences"
	alertNameLabel = "alertname"
	stateActive    = "active"
	stateSilenced  = "suppressed"
)

type Driver struct {
	client   *http.Client
	amUrl    string
	cache    *sync.Map
	periodic *promodel.Duration
}

func New() *Driver {
	return &Driver{
		cache: &sync.Map{},
	}
}

type AmData struct {
	Status string    `json:"status"`
	Data   []AmAlert `json:"data"`
}

type AmSilence struct {
	ID        string    `json:"id"`
	CreatedBy string    `json:"createdBy"`
	Comment   string    `json:"comment"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
}

type AmAlert struct {
	Labels       model.KV
	Annotations  model.KV
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Status       struct {
		State       string        `json:"state"`
		SilencedBy  []string      `json:"silencedBy"`
		InhibitedBy []interface{} `json:"inhibitedBy"`
	} `json:"status"`
	Receivers   []string `json:"receivers"`
	Fingerprint string   `json:"fingerprint"`
}

func (d Driver) toAlertFromRetrieve(amAlert AmAlert) (*model.Alert, error) {
	labels := amAlert.Labels
	alert := &model.Alert{
		Name:         labels[alertNameLabel],
		GeneratorURL: amAlert.GeneratorURL,
		NotifierURL:  d.amUrl,
		StartsAt:     amAlert.StartsAt,
		EndsAt:       amAlert.EndsAt,
		Annotations:  amAlert.Annotations,
	}

	alert.Labels = labels
	alert.ID = amAlert.Fingerprint
	delete(labels, alertNameLabel)
	alert.Labels = labels
	if amAlert.Status.State == stateActive {
		alert.Status = model.Sfiring
		return alert, nil
	}
	alert.Status = model.Ssilenced
	amSilence, err := d.retrieveSilence(amAlert.Status.SilencedBy[0])
	if err != nil {
		return alert, err
	}
	alert.Silence = model.Silence{
		ID:        amSilence.ID,
		CreatedBy: amSilence.CreatedBy,
		Reason:    amSilence.Comment,
		StartsAt:  amSilence.StartsAt,
		EndsAt:    amSilence.EndsAt,
	}
	return alert, nil
}

func (d Driver) retrieveSilence(silenceID string) (AmSilence, error) {
	silenceFRoute := fmt.Sprintf(silenceRoute, silenceID)
	req, err := http.NewRequest("GET", d.amUrl+silenceFRoute, nil)
	if err != nil {
		return AmSilence{}, err
	}
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := d.client.Do(req)
	if err != nil {
		return AmSilence{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AmSilence{}, err
	}
	if resp.StatusCode >= 300 {
		return AmSilence{}, fmt.Errorf("Error when retrieving alerts (status code %d): %s", resp.StatusCode, string(b))
	}
	data := struct {
		Data AmSilence `json:"data"`
	}{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return AmSilence{}, err
	}
	return data.Data, nil
}

func (d Driver) toAlertFromReceived(amAlert promtpl.Alert) *model.Alert {
	labels := amAlert.Labels
	alert := &model.Alert{
		Name:         string(labels[alertNameLabel]),
		GeneratorURL: amAlert.GeneratorURL,
		NotifierURL:  d.amUrl,
		StartsAt:     amAlert.StartsAt,
		EndsAt:       amAlert.EndsAt,
		Annotations:  model.KV(amAlert.Annotations),
		ID:           amAlert.Fingerprint,
	}
	alert.Labels = model.KV(labels)
	delete(labels, alertNameLabel)
	alert.Labels = model.KV(labels)

	if amAlert.Status == "firing" {
		alert.Status = model.Sfiring
	} else {
		alert.Status = model.Sresolved
	}
	return alert
}

func (d Driver) RetrieveAlerts() ([]*model.Alert, error) {
	alerts := make([]*model.Alert, 0)
	req, err := http.NewRequest("GET", d.amUrl+alertRoute, nil)
	if err != nil {
		return alerts, err
	}
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := d.client.Do(req)
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
		alert, err := d.toAlertFromRetrieve(amAlert)
		if err != nil {
			return alerts, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

func (d Driver) ReceiveAlerts(data []byte) ([]*model.Alert, error) {

	alerts := make([]*model.Alert, 0)
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
	d.periodic = config.PeriodicRetrieve
	return nil
}

func (d Driver) Silence(alert *model.Alert) error {
	silence := types.Silence{
		ID:        alert.Silence.ID,
		Comment:   alert.Silence.Reason,
		CreatedBy: alert.Silence.CreatedBy,
		StartsAt:  alert.Silence.StartsAt,
		EndsAt:    alert.Silence.EndsAt,
	}
	matchers := make(labels.Matchers, 0)
	matchers = append(matchers, &labels.Matcher{
		Name:  alertNameLabel,
		Value: alert.Name,
	})
	for key, value := range alert.Labels {
		matchers = append(matchers, &labels.Matcher{
			Name:  key,
			Value: value,
		})
	}
	silence.Matchers = matchers
	b, err := json.Marshal(silence)
	if err != nil {
		return err
	}
	method := "POST"
	if silence.ID != "" {
		method = "PUT"
	}
	req, err := http.NewRequest(method, d.amUrl+silencesRoute, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("Error when retrieving alerts (status code %d): %s", resp.StatusCode, string(b))
	}
	return nil
}

func (d Driver) StartPeriodicAlerts(alertsChan chan<- []*model.Alert) error {
	if d.periodic == nil {
		return nil
	}

	go func() {
		for {
			toAdd := make([]*model.Alert, 0)
			toDelete := make([]string, 0)

			alerts, err := d.RetrieveAlerts()
			if err != nil {
				logrus.Error(err.Error())
			}
			for _, alert := range alerts {
				_, ok := d.cache.Load(alert.ID)
				if ok {
					continue
				}
				toAdd = append(toAdd, alert)
				d.cache.Store(alert.ID, alert)
			}
			d.cache.Range(func(key, value interface{}) bool {
				exists := false
				for _, alert := range alerts {
					if alert.ID == key.(string) {
						exists = true
						break
					}
				}
				if !exists {
					toDelete = append(toDelete, key.(string))
				}
				return true
			})

			for _, key := range toDelete {
				alertRaw, _ := d.cache.Load(key)
				alert := alertRaw.(*model.Alert)
				alert.Status = model.Sresolved
				toAdd = append(toAdd, alert)
				d.cache.Delete(key)
			}
			alertsChan <- toAdd
			time.Sleep(time.Duration(*d.periodic))
		}
	}()
	return nil
}
