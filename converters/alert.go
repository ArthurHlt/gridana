package converters

import (
	"github.com/ArthurHlt/gridana/model"
)

type AlertConverter struct {
	route         *model.Route
	probes        model.Probes
	colorByLabels model.ColorByLabels
	defaultColor  string
}

func NewAlertConverter(
	route *model.Route,
	probes model.Probes,
	colorByLabels model.ColorByLabels,
	defaultColor string,
) *AlertConverter {
	return &AlertConverter{
		route:         route,
		probes:        probes,
		colorByLabels: colorByLabels,
		defaultColor:  defaultColor,
	}
}

func (c AlertConverter) Convert(alert model.Alert) (model.FormattedAlert, error) {
	var err error
	route := c.route.FindRoute(alert)
	probe := c.probes.FindProbe(route.Probe)
	if probe == nil {
		return model.FormattedAlert{}, ErrProbeNotFound
	}

	color, weight := c.colorByLabels.AlertColor(alert)
	if alert.Status == model.Ssilenced {
		weight = -1
	}
	if color == "" {
		color = c.defaultColor
	}
	fmtAlert := model.FormattedAlert{
		Alert:  alert,
		Color:  color,
		Probe:  probe.Name,
		Weight: weight,
	}

	fmtAlert.Identifier, err = GenText(fmtAlert, route.Template)
	if err != nil {
		return model.FormattedAlert{}, err
	}
	fmtAlert.Message, err = GenHTML(fmtAlert, probe.Template.Message)
	if err != nil {
		return model.FormattedAlert{}, err
	}

	fmtAlert.Notification, err = GenText(fmtAlert, probe.Template.Notification)
	if err != nil {
		return model.FormattedAlert{}, err
	}

	return fmtAlert, nil
}
