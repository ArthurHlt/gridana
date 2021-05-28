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
import {LabelsComponent} from './labels/labels.component';
import {PushNotificationsModule} from 'ng-push-ivy';
import {mainRoutingProviders, routing} from './main.route';
import {ProbeComponent} from './probe/probe.component';
import {PersistenceModule} from 'angular-persistence';
import {ConfigComponent} from './config/config.component';
import {ConfigService} from './config.service';
import {SilenceComponent} from './silence/silence.component';
import {AddPipe, DateFormatPipe, MomentModule} from 'ngx-moment';
import {FormsModule} from '@angular/forms';
import {RootComponent} from './root.component';


@NgModule({
  declarations: [
    AppComponent,
    AlertsComponent,
    CardAlertsComponent,
    LabelsComponent,
    ProbeComponent,
    ConfigComponent,
    SilenceComponent,
    RootComponent
  ],
  imports: [
    BrowserModule,
    PersistenceModule,
    FormsModule,
    MomentModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MaterializeModule,
    mainRoutingProviders,
    SimpleNotificationsModule.forRoot(),
    PushNotificationsModule,
    routing
  ],
  providers: [AlertService, WsAlertService, ConfigService, DateFormatPipe, AddPipe],
  bootstrap: [AppComponent]
})
export class AppModule {
}
