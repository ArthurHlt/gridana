export const environment = {
  production: true,
  alertsUrl: "/v1/alerts",
  probesUrl: "/v1/probes",
  alertsOrderedUrl: "/v1/alerts/ordered",
  wsUrl: function () {
    let loc = window.location;
    let scheme = "ws:";
    if (loc.protocol === "https:") {
      scheme = "wss:";
    }
    return scheme + "//" + loc.host + loc.pathname + "/notify";
  }()
};
