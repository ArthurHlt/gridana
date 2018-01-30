package converters

import "github.com/ArthurHlt/gridana/model"

type Converter interface {
	Convert(alert model.Alert) (model.FormattedAlert, error)
}
