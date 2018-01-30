import {Component, OnInit} from '@angular/core';
import {Alert} from '../alert';
import {AlertService} from '../alert.service';
import {OrderedAlerts} from "../orderedAlerts";
import {WsAlertService} from "../ws-alert.service";
import {Probe} from "../probe";

@Component({
  selector: 'app-alerts',
  templateUrl: './alerts.component.html',
  styleUrls: ['./alerts.component.css']
})
export class AlertsComponent implements OnInit {
  selectedAlert: Alert;
  alerts: Alert[];
  probes: Probe[];
  orderedAlerts: OrderedAlerts;


  constructor(private alertService: AlertService, private wsAlert: WsAlertService) {
  }

  ngOnInit() {
    this.getAlerts();
    this.getOrderedAlerts();
    this.getProbes();
    this.wsAlert.observable.subscribe((event: MessageEvent) => {
      let alert = JSON.parse(event.data);
      console.log(alert);
      this.addAlert(alert);
      this.updateOrderedAlert(alert);
    });
  }

  getColor(identifier: string, probe: string): string {
    if (!this.orderedAlerts.alerts[identifier]) {
      return ""
    }
    if (this.orderedAlerts.alerts[identifier][probe].length == 0) {
      return ""
    }
    let color = this.orderedAlerts.alerts[identifier][probe][0].color
    if (color == "yellow") {
      color = "#ffeb3b"
    }
    return color
  }


  addAlert(alert: Alert) {
    if (this.alertExists(alert)) {
      return
    }
    this.alerts.unshift(alert);
  }

  alertExists(alert: Alert): boolean {
    for (let cAlert of this.alerts) {
      if (alert.id == cAlert.id) {
        return true
      }
    }
    return false
  }

  orderedAlertPos(alert: Alert): number {
    if (!this.orderedAlerts.alerts[alert.identifier]) {
      return -1
    }
    let alerts = this.orderedAlerts.alerts[alert.identifier][alert.probe]
    if (alerts.length == 0) {
      return -1
    }
    for (let i = 0; i < alerts.length; i++) {
      let cAlert = alerts[i];
      if (cAlert.id == alert.id) {
        return i
      }
    }
    return -1
  }

  updateOrderedAlert(alert: Alert) {
    if (alert.status == "resolved") {
      this.deleteOrderedAlert(alert);
      return
    }
    this.addOrderedAlert(alert);
  }

  addOrderedAlert(alert: Alert) {
    let alertPos = this.orderedAlertPos(alert);
    if (alertPos >= 0) {
      this.orderedAlerts.alerts[alert.identifier][alert.probe][alertPos] = alert;
      return
    }
    if (this.orderedAlerts.alerts[alert.identifier]) {
      let alerts = this.orderedAlerts.alerts[alert.identifier][alert.probe];
      for (let i = 0; i < alerts; i++) {
        if (alert.weight >= alerts[i].weight) {
          this.orderedAlerts.alerts[alert.identifier][alert.probe].splice(i, 0, alert);
          return
        }
      }
      this.orderedAlerts.alerts[alert.identifier][alert.probe].push(alert);
      return
    }

    this.orderedAlerts.identifiers.unshift(alert.identifier);
    let alertsToAdd = [alert];
    let probes = {};
    for (let probe of this.probes) {
      if (probe.name == alert.probe) {
        probes[probe.name] = alertsToAdd;
        continue
      }
      probes[probe.name] = []
    }
    this.orderedAlerts.alerts[alert.identifier] = probes
  }

  deleteOrderedAlert(alert: Alert) {
    if (!this.orderedAlerts.alerts[alert.identifier]) {
      return
    }
    let alerts = this.orderedAlerts.alerts[alert.identifier][alert.probe];
    if (alerts.length == 0) {
      return
    }
    for (let i = 0; i < this.alerts.length; i++) {
      if (alerts[i].id == alert.id) {
        this.orderedAlerts.alerts[alert.identifier][alert.probe].splice(i, 1);
        break;
      }
    }
    if (this.orderedAlerts.alerts[alert.identifier][alert.probe].length > 0) {
      return
    }
    for (let i = 0; i < this.orderedAlerts.identifiers.length; i++) {
      if (this.orderedAlerts.identifiers[i] == alert.identifier) {
        this.orderedAlerts.identifiers.splice(i, 1);
        break;
      }
    }
  }

  onSelect(alert: Alert): void {
    this.selectedAlert = alert;
  }

  getAlerts(): void {
    this.alertService.getAlerts()
      .subscribe(alerts => this.alerts = alerts);
  }

  getProbes(): void {
    this.alertService.getProbes()
      .subscribe(probes => this.probes = probes);
  }

  getOrderedAlerts(): void {
    this.alertService.getOrderedAlerts()
      .subscribe(orderedAlerts => this.orderedAlerts = orderedAlerts);
  }
}

