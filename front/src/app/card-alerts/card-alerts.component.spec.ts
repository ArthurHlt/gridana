import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { CardAlertsComponent } from './card-alerts.component';

describe('CardAlertsComponent', () => {
  let component: CardAlertsComponent;
  let fixture: ComponentFixture<CardAlertsComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ CardAlertsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CardAlertsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
