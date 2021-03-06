import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProbeComponent } from './probe.component';

describe('ProbeComponent', () => {
  let component: ProbeComponent;
  let fixture: ComponentFixture<ProbeComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ ProbeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProbeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
