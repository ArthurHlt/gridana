import {Component, Input} from '@angular/core';

@Component({
  selector: 'app-labels',
  templateUrl: './labels.component.html',
  styleUrls: ['./labels.component.css']
})
export class LabelsComponent {
  @Input() labels: Map<string, string>;

  constructor() {
  }

  labelKeys() {
    return Object.keys(this.labels);
  }

}
