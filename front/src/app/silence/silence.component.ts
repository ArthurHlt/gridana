import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Alert} from "../alert";
import {AddPipe, DateFormatPipe} from "angular2-moment";
import {AlertService} from "../alert.service";

@Component({
  selector: 'app-silence',
  templateUrl: './silence.component.html',
  styleUrls: ['./silence.component.css']
})
export class SilenceComponent implements OnInit {
  @Input() alert: Alert;
  @Output() close = new EventEmitter<string>();
  startAt: any;
  endsAt: any;
  durationNumber: number;
  durationType: any;
  durationRaw: string;
  createdBy: string;
  reason: string;
  error: string;
  success: boolean;

  constructor(private dateFormatPipe: DateFormatPipe,
              private addPipe: AddPipe,
              private alertService: AlertService) {
  }

  ngOnInit() {
    this.success = false;
    this.durationNumber = 2;
    this.durationType = 'h';
    this.durationRaw = this.durationNumber + this.durationType;
    this.startAt = this.dateFormatPipe.transform(Date.now());
    this.endsAt = this.dateFormatPipe.transform(this.addPipe.transform(this.startAt, this.durationNumber, this.durationType));
  }

  parseDuration() {
    let reg = new RegExp('^([1-9]+)([A-Za-z]+)$');
    if (!reg.test(this.durationRaw)) {
      return
    }
    let match = reg.exec(this.durationRaw);
    this.durationNumber = Number(match[1]);
    this.durationType = match[2];
    this.endsAt = this.dateFormatPipe.transform(this.addPipe.transform(this.startAt, this.durationNumber, this.durationType));
  }

  createSilence() {
    this.error = null;
    this.alert.silence.createdBy = this.createdBy;
    this.alert.silence.reason = this.reason;
    this.alert.silence.startsAt = this.startAt;
    this.alert.silence.endsAt = this.endsAt;
    this.alert.status = "silenced";
    this.alertService.setSilence(this.alert).subscribe(
      data => {
        this.success = true;
        setTimeout(() => {
          this.close.emit("closing");
        }, 500)
      },
      err => {
        console.log(err.error);
        this.error = err.error.title + ": " + err.error.details;
      }
    );
  }

  cancel() {
    this.close.emit("closing");
  }
}
