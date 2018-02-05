package gridana

import (
	"github.com/ArthurHlt/gridana/converters"
	"github.com/ArthurHlt/gridana/drivers"
	"github.com/ArthurHlt/gridana/emitter"
	"github.com/ArthurHlt/gridana/model"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	drivers    map[string]drivers.Driver
	converter  converters.Converter
	dropLabels model.LabelMatcher
}

func NewProcessor(drivers map[string]drivers.Driver, converter converters.Converter, dropLabels model.LabelMatcher) *Processor {
	return &Processor{
		drivers:    drivers,
		converter:  converter,
		dropLabels: dropLabels,
	}
}

func (p Processor) RetrieveAlerts() []model.FormattedAlert {
	log.Debug("Retrieving alerts ...")
	alerts := make([]model.Alert, 0)
	for name, driver := range p.drivers {
		entry := log.WithField("driver_name", name)
		entry.Debug("Retrieving alerts for driver ...")
		newAlerts, err := driver.RetrieveAlerts()
		if err != nil {
			entry.Error(err)
			continue
		}
		alerts = append(alerts, newAlerts...)
		entry.Debug("Finished retrieving alerts for driver.")
	}
	log.Debug("Finished retrieving alerts.")

	log.Debug("Processing alerts ...")
	fmtAlerts := make([]model.FormattedAlert, 0)
	for _, alert := range alerts {
		if p.dropLabels.MatchAlert(alert) {
			continue
		}
		entry := log.WithField("alert_id", alert.ID)
		fmtAlert, err := p.converter.Convert(alert)
		if err != nil {
			entry.Error(err)
			continue
		}
		fmtAlerts = append(fmtAlerts, fmtAlert)
	}
	log.Debug("Finished processing alerts.")
	return fmtAlerts
}
func (p Processor) SilenceAlert(driverName string, alert model.Alert) error {
	entry := log.WithField("driver_name", driverName)
	entry.Debug("Silencing alert ...")
	if _, ok := p.drivers[driverName]; !ok {
		entry.Warning("Can't find this driver, skipping.")
		return nil
	}
	if _, ok := p.drivers[driverName].(drivers.DriverSilencer); !ok {
		entry.Warning("This driver doesn't implement silencer, skipping.")
		return nil
	}
	driver := p.drivers[driverName].(drivers.DriverSilencer)
	err := driver.Silence(alert)
	if err != nil {
		return err
	}
	entry.Debug("Finished silencing alert.")
	return nil
}
func (p Processor) ReceiveAlerts(driverName string, data []byte) {
	entry := log.WithField("driver_name", driverName)
	entry.Debug("Receiving alerts ...")
	if _, ok := p.drivers[driverName]; !ok {
		entry.Warning("Can't find this driver, skipping.")
		return
	}
	if _, ok := p.drivers[driverName].(drivers.DriverNotifier); !ok {
		entry.Warning("This driver doesn't implement notifier, skipping.")
		return
	}
	driver := p.drivers[driverName].(drivers.DriverNotifier)
	alerts, err := driver.ReceiveAlerts(data)
	if err != nil {
		entry.Error(err)
		return
	}
	entry.Debug("Finished receiving alerts.")

	entry.Debug("Processing alerts ...")
	for _, alert := range alerts {
		if p.dropLabels.MatchAlert(alert) {
			continue
		}
		entry := log.WithField("alert_id", alert.ID)
		fmtAlert, err := p.converter.Convert(alert)
		if err != nil {
			entry.Error(err)
			continue
		}
		emitter.Emit(fmtAlert)
	}
	entry.Debug("Finished processing alerts.")
}
