package gridana

import (
	"encoding/json"
	"fmt"
	"github.com/ArthurHlt/gridana/emitter"
	"github.com/ArthurHlt/gridana/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	resp "github.com/nicklaw5/go-respond"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type API struct {
	processor   *Processor
	probes      model.Probes
	driversName []string
	upgrader    websocket.Upgrader
}

func NewApi(processor *Processor, upgrader websocket.Upgrader, probes model.Probes, driversName []string) *API {
	sort.Sort(probes)
	return &API{
		processor:   processor,
		probes:      probes,
		driversName: driversName,
		upgrader:    upgrader,
	}
}

func (api API) Register(r *mux.Router) {
	corsHandler := cors.AllowAll()
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/probes", api.listProbes).Methods("GET")
	s.HandleFunc("/alerts", api.listAlerts).Methods("GET")
	s.HandleFunc("/alerts/ordered", api.listOrderedAlerts).Methods("GET")
	s.HandleFunc("/silence", api.silence).Methods("POST", "PUT", "OPTIONS")
	s.Use(handlers.CompressHandler)
	s.Use(corsHandler.Handler)
	s.Use(panicHandler)

	r.HandleFunc("/notify", api.notify)
	r.HandleFunc("/webhook/{driver}", api.webhook).Methods("POST")
}
func (api API) notify(w http.ResponseWriter, req *http.Request) {
	fmt.Println(fmt.Sprintf("Listeners: %d", len(emitter.Listeners())))
	entry := log.WithField("verb", "notify")
	c, err := api.upgrader.Upgrade(w, req, nil)
	if err != nil {
		entry.Error(err)
		return
	}
	entry = entry.WithField("user_addr", c.RemoteAddr())
	entry.Debug("User connected.")
	closeWrite := make(chan bool)
	defer func() {
		entry.Debug("User disconnected.")
		c.Close()
		closeWrite <- true
	}()
	go notifyWriter(c, closeWrite)
	notifyReader(c)
}
func (api API) webhook(w http.ResponseWriter, req *http.Request) {
	driver := mux.Vars(req)["driver"]
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(HttpError{
			Verb:    "webhook",
			Details: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}
	api.processor.ReceiveAlerts(driver, data)
}

func (api API) silence(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(HttpError{
			Verb:    "silence",
			Details: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}
	var fmtAlert model.FormattedAlert
	err = json.Unmarshal(data, &fmtAlert)
	if err != nil {
		panic(HttpError{
			Verb:    "silence",
			Details: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}
	if fmtAlert.Silence.CreatedBy == "" {
		panic(HttpError{
			Verb:    "silence",
			Details: "User must be set",
			Status:  http.StatusBadRequest,
		})
	}
	if fmtAlert.Silence.Reason == "" {
		panic(HttpError{
			Verb:    "silence",
			Details: "Reason must be set",
			Status:  http.StatusBadRequest,
		})
	}
	if fmtAlert.Silence.EndsAt.Before(fmtAlert.Silence.StartsAt) {
		panic(HttpError{
			Verb:    "silence",
			Details: "end time can't be before start at",
			Status:  http.StatusBadRequest,
		})
	}
	if fmtAlert.Silence.EndsAt.Before(time.Now()) {
		panic(HttpError{
			Verb:    "silence",
			Details: "end time can't be in the past",
			Status:  http.StatusBadRequest,
		})
	}

	for _, driverName := range api.driversName {
		err = api.processor.SilenceAlert(driverName, fmtAlert.Alert)
		if err != nil {
			panic(HttpError{
				Verb:    "silence",
				Details: err.Error(),
				Status:  http.StatusInternalServerError,
			})
		}
	}
	resp.NewResponse(w).Ok(fmtAlert)
}

func (api API) listOrderedAlerts(w http.ResponseWriter, req *http.Request) {
	alerts := model.FormattedAlerts(api.processor.RetrieveAlerts()).Ordered(api.probes)
	identifiers := make([]string, 0)
	for identifier, _ := range alerts {
		identifiers = append(identifiers, identifier)
	}
	sort.Strings(identifiers)
	resp.NewResponse(w).Ok(struct {
		Alerts      model.OrderedAlerts `json:"alerts"`
		Identifiers []string            `json:"identifiers"`
	}{alerts, identifiers})
}

func (api API) listAlerts(w http.ResponseWriter, req *http.Request) {
	alerts := api.processor.RetrieveAlerts()
	sort.Sort(model.ByStartAt{alerts})
	resp.NewResponse(w).Ok(alerts)
}

func (api API) listProbes(w http.ResponseWriter, req *http.Request) {
	probes := api.probes
	resp.NewResponse(w).Ok(probes)
}
