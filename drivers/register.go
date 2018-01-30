package drivers

import "fmt"

var drivers map[string]Factory = make(map[string]Factory)

type Factory func() Driver

func Register(driverType string, factory Factory) {
	drivers[driverType] = factory
}
func GenerateDriver(driverType string) (Driver, error) {
	if _, ok := drivers[driverType]; !ok {
		return nil, fmt.Errorf("Can't found driver")
	}
	return drivers[driverType](), nil
}
