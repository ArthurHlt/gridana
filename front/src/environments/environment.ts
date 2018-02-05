// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.

export const environment = {
  production: false,
  alertsUrl: "http://localhost:8080/v1/alerts",
  silenceUrl: "http://localhost:8080/v1/silence",
  probesUrl: "http://localhost:8080/v1/probes",
  alertsOrderedUrl: "http://localhost:8080/v1/alerts/ordered",
  wsUrl: "ws://localhost:8080/notify"
};
