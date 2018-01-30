import { TestBed, inject } from '@angular/core/testing';

import { WsAlertService } from './ws-alert.service';

describe('WsAlertService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [WsAlertService]
    });
  });

  it('should be created', inject([WsAlertService], (service: WsAlertService) => {
    expect(service).toBeTruthy();
  }));
});
