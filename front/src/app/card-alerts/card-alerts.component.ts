import {Component, OnInit} from '@angular/core';
import {Alert} from "../alert";


@Component({
  selector: 'app-card-alerts',
  templateUrl: './card-alerts.component.html',
  styleUrls: ['./card-alerts.component.css'],
})
export class CardAlertsComponent implements OnInit {

  alerts: Alert[];
  show: boolean;
  scale: boolean;

  constructor() {
  }

  static stopPropagation(event) {
    event.stopPropagation();
  }

  toggle() {

    let willShow = !this.show;
    if (!willShow) {
      this.scale = false;
      setTimeout(() => {
        this.show = willShow;
      }, 200);
      return
    }
    this.show = willShow;
    setTimeout(() => {
      this.scale = true;
    }, 50);

  }

  ngOnInit() {
    this.show = false;
    this.scale = false;
  }

}
