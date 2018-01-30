import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {SimpleNotificationsModule} from 'angular2-notifications';
import {MaterializeModule} from 'angular2-materialize';
import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AlertService} from './alert.service';


import {AppComponent} from './app.component';
import {WsAlertService} from './ws-alert.service';
import {AlertsComponent} from './alerts/alerts.component';
import {CardAlertsComponent} from './card-alerts/card-alerts.component';


@NgModule({
  declarations: [
    AppComponent,
    AlertsComponent,
    CardAlertsComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MaterializeModule,
    SimpleNotificationsModule.forRoot()
  ],
  providers: [AlertService, WsAlertService],
  bootstrap: [AppComponent]
})
export class AppModule {
}
