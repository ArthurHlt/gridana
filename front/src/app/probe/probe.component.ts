import {Component, Input, OnInit} from '@angular/core';
import {Alert} from "../alert";

@Component({
  selector: 'app-probe',
  templateUrl: './probe.component.html',
  styleUrls: ['./probe.component.css']
})
export class ProbeComponent implements OnInit {
  @Input() alert: Alert;
  @Input() zindex: string;

  constructor() {
    this.zindex = "auto";
  }

  ngOnInit() {
  }

  isSilenced(): boolean {
    return this.alert.status == "silenced";
  }
}
