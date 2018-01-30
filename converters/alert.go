package converters

import (
	"github.com/ArthurHlt/gridana/model"
)

type AlertConverter struct {
	route         *model.Route
	probes        model.Probes
	colorByLabels model.ColorByLabels
	silenceColor  string
	defaultColor  string
}

func NewAlertConverter(
	route *model.Route,
	probes model.Probes,
	colorByLabels model.ColorByLabels,
	defaultColor string,
	silenceColor string,
) *AlertConverter {
	return &AlertConverter{
		route:         route,
		probes:        probes,
		colorByLabels: colorByLabels,
		silenceColor:  silenceColor,
		defaultColor:  defaultColor,
	}
}

func (c AlertConverter) Convert(alert model.Alert) (model.FormattedAlert, error) {
	route := c.route.FindRoute(alert)
	probe := c.probes.FindProbe(route.Probe)
	if probe == nil {
		return model.FormattedAlert{}, ErrProbeNotFound
	}
	message, err := GenAlertMessage(alert, probe.Template)
	if err != nil {
		return model.FormattedAlert{}, err
	}

	identifier, err := GenAlertIdentifier(alert, route.Template)
	if err != nil {
		return model.FormattedAlert{}, err
	}

	var color string
	var weight int
	if alert.Status == model.Ssilenced {
		color = c.silenceColor
		weight = -1
	} else {
		color, weight = c.colorByLabels.AlertColor(alert)
	}
	if color == "" {
		color = c.defaultColor
	}
	return model.FormattedAlert{
		Alert:      alert,
		Color:      color,
		Probe:      probe.Name,
		Message:    message,
		Identifier: identifier,
		Weight:     weight,
	}, nil
}
