import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';

import {Observable} from 'rxjs/Observable';

import {Alert} from './alert';
import {OrderedAlerts} from './orderedAlerts';
import {Probe} from "./probe";

import {environment} from '../environments/environment';

const httpOptions = {
  headers: new HttpHeaders({'Content-Type': 'application/json'})
};

@Injectable()
export class AlertService {
  private alertsUrl = environment.alertsUrl;
  private probesUrl = environment.probesUrl;
  private alertsOrderedUrl = environment.alertsOrderedUrl;


  constructor(private http: HttpClient) {
  }

  getProbes(): Observable<Probe[]> {
    return this.http.get<Probe[]>(this.probesUrl);
  }

  getAlerts(): Observable<Alert[]> {
    return this.http.get<Alert[]>(this.alertsUrl);
  }

  getOrderedAlerts(): Observable<OrderedAlerts> {
    return this.http.get<OrderedAlerts>(this.alertsOrderedUrl);
  }
}
