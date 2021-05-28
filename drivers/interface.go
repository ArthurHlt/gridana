package drivers

import "github.com/ArthurHlt/gridana/model"

type Driver interface {
	RetrieveAlerts() ([]*model.Alert, error)
	Config(config model.DriverConfig) error
}

type DriverNotifier interface {
	ReceiveAlerts(data []byte) ([]*model.Alert, error)
}

type DriverSilencer interface {
	Silence(alert *model.Alert) error
}

type DriverPeriodic interface {
	StartPeriodicAlerts(alerts chan<- []*model.Alert) error
}
