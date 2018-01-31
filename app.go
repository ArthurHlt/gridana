package gridana

import (
	"github.com/ArthurHlt/gridana/converters"
	"github.com/ArthurHlt/gridana/drivers"
	"github.com/ArthurHlt/gridana/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net/http"
	"strings"
)

type App struct {
	config *model.GridanaConfig
}

func NewApp() *App {
	return &App{}
}
func (a App) Run(args []string) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-path, c",
			Value: "./config.yml",
			Usage: "Path to configuration file",
		},
	}
	app.Name = "gridana"
	app.Action = a.start

	app.Run(args)
}
func (a App) start(c *cli.Context) error {
	err := a.loadConfig(c.GlobalString("config-path"))
	if err != nil {
		return err
	}
	a.loadLogConfig()
	mDrivers, err := a.bootstrapDrivers()
	if err != nil {
		return err
	}
	api := a.bootstrapApi(mDrivers)
	r := mux.NewRouter()
	api.Register(r)
	r.NewRoute().Handler(http.FileServer(assetFS()))
	log.WithField("listen_addr", a.config.ListenAddr).
		Info("Server started and listening.")
	corsHandler := cors.AllowAll()

	return http.ListenAndServe(a.config.ListenAddr, corsHandler.Handler(r))
}
func (a App) bootstrapApi(d map[string]drivers.Driver) *API {
	driversName := make([]string, 0)
	for driverName, _ := range d {
		driversName = append(driversName, driverName)
	}
	upgrader := websocket.Upgrader{}
	if a.config.NoCheckOrigin {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	return NewApi(a.bootstrapProcessor(d), upgrader, a.config.Probes, driversName)
}

func (a App) bootstrapProcessor(d map[string]drivers.Driver) *Processor {
	return NewProcessor(d, a.bootstrapConverter(), a.config.DropLabels)
}

func (a App) bootstrapConverter() converters.Converter {
	return converters.NewAlertConverter(a.config.Route, a.config.Probes, a.config.ColorByLabels, a.config.DefaultColor, a.config.SilenceColor)
}

func (a App) bootstrapDrivers() (map[string]drivers.Driver, error) {
	mDrivers := make(map[string]drivers.Driver)
	for _, driverConfig := range a.config.Drivers {
		entry := log.WithField("driver_name", driverConfig.Name).WithField("driver_type", driverConfig.Type)
		driver, err := drivers.GenerateDriver(driverConfig.Type)
		if err != nil {
			return mDrivers, err
		}
		err = driver.Config(driverConfig)
		if err != nil {
			return mDrivers, err
		}
		mDrivers[driverConfig.Name] = driver
		entry.Debug("Driver loaded")
	}
	return mDrivers, nil
}
func (a *App) loadConfig(configPath string) error {
	c, err := model.LoadFile(configPath)
	if err != nil {
		return err
	}
	a.config = c
	return nil
}

func (a App) loadLogConfig() {
	c := a.config
	if c.Logs.InJson {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: c.Logs.NoColor,
		})
	}
	switch strings.ToUpper(c.Logs.Level) {
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
		break
	case "WARN":
		log.SetLevel(log.WarnLevel)
		break
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		break
	case "PANIC":
		log.SetLevel(log.PanicLevel)
		break
	case "FATAL":
		log.SetLevel(log.FatalLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
		break
	}
}
