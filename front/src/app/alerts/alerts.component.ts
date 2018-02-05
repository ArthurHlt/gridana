import {Component, NgZone, OnInit, ViewChild} from '@angular/core';
import {Alert} from '../alert';
import {AlertService} from '../alert.service';
import {OrderedAlerts} from "../orderedAlerts";
import {WsAlertService} from "../ws-alert.service";
import {Probe} from "../probe";
import {NotificationsService} from "angular2-notifications";
import {CardAlertsComponent} from "../card-alerts/card-alerts.component";
import {PushNotificationsService} from "ng-push";
import {capitalizeFirstLetter} from "../utils";
import {ConfigService} from "../config.service";

@Component({
  selector: 'app-alerts',
  templateUrl: './alerts.component.html',
  styleUrls: ['./alerts.component.css']
})
export class AlertsComponent implements OnInit {
  @ViewChild(CardAlertsComponent) cardAlerts: CardAlertsComponent;
  alerts: Alert[];
  probes: Probe[];
  pushOnResolved: boolean;
  orderedAlerts: OrderedAlerts;
  showSilenced: boolean;
  notifOptions = {
    position: ["top", "right"],
    timeOut: 5000,
    lastOnBottom: true,
    showProgressBar: true,
    pauseOnHover: true,
    clickToClose: false,
    clickIconToClose: true
  };

  constructor(private alertService: AlertService,
              private wsAlert: WsAlertService,
              private notifService: NotificationsService,
              private browserNotifSvc: PushNotificationsService,
              private zone: NgZone,
              private configService: ConfigService) {
    this.alerts = [];
    configService.defineShowSilenced(this, 'showSilenced');
    configService.defineReceiveResolved(this, 'pushOnResolved');
  }


  ngOnInit() {
    this.browserNotifSvc.requestPermission();
    this.getOrderedAlerts();
    this.getProbes();
    this.wsAlert.observable.subscribe((event: MessageEvent) => {
        let alert = JSON.parse(event.data);
        this.pushAlertNotification(alert);
        this.updateAlert(alert);
        this.updateOrderedAlert(alert);
      },
      () => {
        this.pushNotification("Websocket lost connection", "Connection has been lost, please refresh page.")
      });
  }

  getFirstAlert(identifier: string, probe: string): Alert {
    if (!this.orderedAlerts.alerts[identifier]) {
      return null
    }
    if (this.orderedAlerts.alerts[identifier][probe].length == 0) {
      return null
    }

    return this.orderedAlerts.alerts[identifier][probe][0];
  }

  getGroupAlerts(identifier: string, probe: string): Alert[] {
    if (!this.orderedAlerts.alerts[identifier]) {
      return []
    }
    if (this.orderedAlerts.alerts[identifier][probe].length == 0) {
      return []
    }
    return this.orderedAlerts.alerts[identifier][probe]
  }

  addAlert(alert: Alert) {
    let alertPos = this.alertPos(alert);
    if (alertPos < 0) {
      this.alerts.unshift(alert);
      return
    }
    this.alerts[alertPos] = alert;
  }

  deleteAlert(alert: Alert) {
    if (!this.alertExists(alert)) {
      return
    }
    for (let i = 0; i < this.alerts.length; i++) {
      if (this.alerts[i].id == alert.id) {
        this.alerts.splice(i, 1);
        break;
      }
    }
  }

  alertPos(alert: Alert): number {
    for (let i = 0; i < this.alerts.length; i++) {
      if (alert.id == this.alerts[i].id) {
        return i
      }
    }
    return -1
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
    let alerts = this.orderedAlerts.alerts[alert.identifier][alert.probe];
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

  updateAlert(alert: Alert) {
    if (alert.status == "resolved") {
      this.deleteAlert(alert);
      return
    }
    this.addAlert(alert);
  }

  updateOrderedAlert(alert: Alert) {
    if (alert.status == "resolved") {
      this.deleteOrderedAlert(alert);
      return
    }
    this.addOrderedAlert(alert);
  }

  pushSiteNotification(alert: Alert) {
    let notif = this.notifService.alert(capitalizeFirstLetter(alert.status), alert.notification);
    notif.click.subscribe((event) => {
      this.cardAlerts.showDetails(alert);
    });
  }

  pushNotification(title: string, body: string) {
    let notif = this.browserNotifSvc.create(
      title,
      {body: body}
    );
    notif.subscribe(
      res => {
      },
      err => {
        this.notifService.alert(title, body);
      }
    );
  }

  pushAlertNotification(alert: Alert) {
    let alertExists = this.alertExists(alert);
    if (alertExists && !this.pushOnResolved) {
      return
    }
    if (alertExists && alert.status != "resolved") {
      return
    }
    let notif = this.browserNotifSvc.create(
      capitalizeFirstLetter(alert.status),
      {body: alert.notification}
    );
    notif.subscribe(
      res => {
        if (res.event.type === 'click') {
          this.zone.run(() => {
            res.notification.close();
            this.cardAlerts.showDetails(alert);
          });
        }
      },
      err => {
        this.pushSiteNotification(alert);
      }
    );
  }

  addOrderedAlert(alert: Alert) {
    let alertPos = this.orderedAlertPos(alert);
    if (alertPos >= 0) {
      this.orderedAlerts.alerts[alert.identifier][alert.probe][alertPos] = alert;
      return
    }
    if (this.orderedAlerts.alerts[alert.identifier]) {
      let alerts = this.orderedAlerts.alerts[alert.identifier][alert.probe];
      for (let i = 0; i < alerts.length; i++) {
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

  hasShowableAlerts(identifier: string): boolean {
    if (!this.orderedAlerts.alerts[identifier]) {
      return false;
    }
    let count = 0;
    let probes = this.orderedAlerts.alerts[identifier];
    for (let probe in probes) {
      let alerts = this.orderedAlerts.alerts[identifier][probe];
      for (let alert of alerts) {
        if (alert.status == "silenced" && !this.showSilenced) {
          continue;
        }
        count++;
      }
    }
    return count > 0;
  }

  gridShowable(): boolean {
    if (!this.orderedAlerts) {
      return false;
    }
    let ids = this.orderedAlerts.identifiers;
    for (let id of ids) {
      if (this.hasShowableAlerts(id)) {
        return true;
      }
    }
    return false;
  }

  getProbes(): void {
    this.alertService.getProbes()
      .subscribe(probes => this.probes = probes);
  }

  getOrderedAlerts(): void {
    this.alertService.getOrderedAlerts()
      .subscribe(orderedAlerts => {
        this.orderedAlerts = orderedAlerts;
        let identifiers = this.orderedAlerts.alerts;
        for (let identifier in identifiers) {
          let probes = this.orderedAlerts.alerts[identifier];
          for (let probe in probes) {
            let alerts = this.orderedAlerts.alerts[identifier][probe];
            for (let alert of alerts) {
              this.addAlert(alert);
            }
          }
        }
      });
  }
}

