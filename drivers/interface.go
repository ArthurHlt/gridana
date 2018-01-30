package drivers

import "github.com/ArthurHlt/gridana/model"

type Driver interface {
	RetrieveAlerts() ([]model.Alert, error)
	Config(config model.DriverConfig) error
}

type DriverNotifier interface {
	ReceiveAlerts(data []byte) ([]model.Alert, error)
}
