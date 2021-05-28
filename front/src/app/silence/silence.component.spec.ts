import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SilenceComponent } from './silence.component';

describe('SilenceComponent', () => {
  let component: SilenceComponent;
  let fixture: ComponentFixture<SilenceComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ SilenceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SilenceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
