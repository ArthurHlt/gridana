import {Component, OnInit} from '@angular/core';
import {ConfigService} from '../config.service';

@Component({
  selector: 'app-config',
  templateUrl: './config.component.html',
  styleUrls: ['./config.component.css']
})
export class ConfigComponent implements OnInit {

  pushOnResolved: boolean;
  showSilenced: boolean;

  constructor(private configService: ConfigService) {
    configService.defineShowSilenced(this, 'showSilenced');
    configService.defineReceiveResolved(this, 'pushOnResolved');
    if (!this.showSilenced) {
      this.showSilenced = false;
    }
    if (!this.pushOnResolved) {
      this.pushOnResolved = false;
    }
  }

  ngOnInit() {
  }

}
