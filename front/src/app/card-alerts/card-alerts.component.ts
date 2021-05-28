import {Component, OnInit} from '@angular/core';
import {Alert} from '../alert';


@Component({
  selector: 'app-card-alerts',
  templateUrl: './card-alerts.component.html',
  styleUrls: ['./card-alerts.component.css'],
})
export class CardAlertsComponent implements OnInit {

  alerts: Alert[];
  selectAlert: Alert;
  show: boolean;
  scale: boolean;
  silenceSection: boolean;

  constructor() {
    this.alerts = [];
  }

  stopPropagation(event) {
    event.stopPropagation();
  }

  showAlerts(alerts: Alert[]) {
    this.selectAlert = null;
    this.alerts = alerts;
    this.toggle();

  }

  toggle() {
    if (this.selectAlert && this.alerts.length > 0) {
      this.scale = false;
      setTimeout(() => {
        this.selectAlert = null;
        this.scale = true;
      }, 200);
      return;
    }
    const willShow = !this.show;
    if (!willShow) {
      this.scale = false;
      setTimeout(() => {
        this.show = willShow;
      }, 200);
      return;
    }
    this.show = willShow;
    setTimeout(() => {
      this.scale = true;
    }, 50);

  }

  showDetailsFromCards(alert: Alert, event) {
    event.stopPropagation();
    this.silenceSection = false;
    this.scale = false;
    setTimeout(() => {
      this.selectAlert = alert;
      this.scale = true;
    }, 200);
  }

  showSilenceFromCards(alert: Alert, event) {
    event.stopPropagation();
    this.scale = false;
    setTimeout(() => {
      this.selectAlert = alert;
      this.silenceSection = true;
      this.scale = true;
    }, 200);
  }

  showSilenceFromDetails(alert: Alert) {
    this.scale = false;
    setTimeout(() => {
      this.silenceSection = true;
      this.scale = true;
    }, 200);
  }

  showDetails(alert: Alert) {
    this.silenceSection = false;
    this.alerts = [];
    this.selectAlert = alert;
    this.toggle();
  }

  showSilence(alert: Alert) {
    this.silenceSection = true;
    this.alerts = [];
    this.selectAlert = alert;
    this.toggle();
  }

  ngOnInit() {
    this.show = false;
    this.scale = false;
  }

}
