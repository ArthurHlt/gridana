import {Component} from '@angular/core';
import {WsAlertService} from "./ws-alert.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  title = 'Gridana';

  constructor(private wsAlert: WsAlertService) {
  }

  stopPropagation(event) {
    event.stopPropagation();
  }

  ngOnInit() {
  }
}
