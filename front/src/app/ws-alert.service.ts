import {Injectable} from '@angular/core'
import 'rxjs/add/operator/share'
import {Observable} from 'rxjs/Observable';
import {Observer} from "rxjs/Observer";

import {environment} from '../environments/environment';

@Injectable()
export class WsAlertService {

  public observable: Observable<MessageEvent>;

  constructor() {
    let url = environment.wsUrl;
    this.observable = this.create(url);
    console.log("Successfully connected: " + url);
  }

  private create(url): Observable<MessageEvent> {
    let ws = new WebSocket(url);

    return Observable.create(
      (obs: Observer<MessageEvent>) => {
        ws.onmessage = obs.next.bind(obs);
        ws.onerror = obs.error.bind(obs);
        ws.onclose = obs.complete.bind(obs);
        return ws.close.bind(ws);
      });
  }

}

