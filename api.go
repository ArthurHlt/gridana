package gridana

import (
	"github.com/ArthurHlt/gridana/emitter"
	"github.com/ArthurHlt/gridana/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	resp "github.com/nicklaw5/go-respond"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
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
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/probes", api.listProbes).Methods("GET")
	s.HandleFunc("/alerts", api.listAlerts).Methods("GET")
	s.HandleFunc("/alerts/ordered", api.listOrderedAlerts).Methods("GET")

	r.HandleFunc("/notify", api.notify)
	r.HandleFunc("/webhook/{driver}", api.webhook).Methods("POST")
}
func (api API) notify(w http.ResponseWriter, req *http.Request) {
	entry := log.WithField("verb", "notify")
	c, err := api.upgrader.Upgrade(w, req, nil)
	if err != nil {
		entry.Error(err)
		return
	}

	entry = entry.WithField("user_addr", c.RemoteAddr())
	entry.Debug("User connected.")
	defer func() {
		entry.Debug("User disconnected.")
		c.Close()
	}()
	// TODO: add mechanism to release disconnected client from socket
	// for now they will be released only when events have been received which can be long
	for event := range emitter.On() {
		alert := emitter.ToAlert(event)
		err := c.WriteJSON(alert)
		if err != nil {
			entry.Debug(err)
			break
		}
		entry.Debug("Event sent to client")
	}

}
func (api API) webhook(w http.ResponseWriter, req *http.Request) {
	driver := mux.Vars(req)["driver"]
	entry := log.WithField("verb", "webhook").WithField("driver", driver)
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		entry.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	api.processor.ReceiveAlerts(driver, data)
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
