# Gridana

Extensible and customizable grid system for any alerting system.

## Drivers

Drivers are extensions to retrieve and receive alert from an alerting system.

Builtin drivers:
- [alertmanager](https://github.com/prometheus/alertmanager)

## Run dev environment

**Requirements**:
- docker and docker-compose
- angular2 cli (see doc: https://angular.io/guide/quickstart#devenv )
- go(lang)

**Steps**:

**Retrieve with `go get github.com/ArthurHlt/gridana` and go to `$GOPATH/src/github.com/ArthurHlt/gridana`**

1. Run docker
2. Run alertmanager with docker-compose: `docker-compose up -d -f alertmanager-docker/docker-compose.yml`
3. Run gridana backend: `go run server/main.go` (You can configure the `config.yml` file for set your own routing)
4. Go to `front` directory and run `npm install`
5. Run from this folder: `ng serve --open`

You can now push alerts to alert manager with the helper `send-alerts` 
inside `alertmanager-docker`, e.g.: `./alertmanager-docker/send-alerts --timeout 60` (this will make alert expire after 1h)

**Note**: `config.yml` is not yet documented, see `model/config.go` file to see all available option for configuration.

## Roadmap

### Roadmap leading to mvp

- [x] See alerts in a grid scrapped from driver
- [x] Routing alerts by their labels to aggregate on a probe
- [x] Alert message can be templatized in markdown
- [x] Alert identifer can be templatized
- [x] Receive new firing/resolved alerts by websocket leading to update grid with new information
- [x] Provide mechanism to send all incoming alerts to a multi-instanced gridana (For now only amqp is available)
- [ ] Show all alerts aggregate by a probe and identifier as cards in an overlay
- [ ] Show all alerts on a side nav
- [ ] Options to see/hide silenced alert
- [ ] Js cron to remove expired alert which has been silenced on alerting system 
(alertmanager doesn't send that alert has been silenced, other driver could do the same)
- [ ] Push notification (on ui and OS/browser notification system) when receiving a firing alert 
(Potentially, have an option to see when they are resolved too)


### Roadmap for first release

- [ ] Create and manage dashboard outside of `config.yml` file (by passing files or from a database like grafana)
- [ ] Create and manage drivers from ui too (not only on `config.yml` file)
- [ ] Receive, scrap and susbscribe to websocket on a specific or multiple driver on ui
(actually, ui receive all alerts from all drivers)
- [ ] Filtering the grid items by identifier in ui

### Miscellaneous ideas (far from now)

- Provide authentication in a sidecar or directly through lib like [gobis](https://github.com/orange-cloudfoundry/gobis)
(a reverse proxy in front a multi-instanced gridana could be a real bottleneck cause of websocket)
- Wysiwyg to create dashboard (thinking about a graph system for routing)
- And all ideas we could have after mvp incubation

