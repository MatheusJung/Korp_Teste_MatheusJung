import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ServiceOfflineComponent } from './service-offline';

describe('ServiceOffline', () => {
  let component: ServiceOfflineComponent;
  let fixture: ComponentFixture<ServiceOfflineComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ServiceOfflineComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ServiceOfflineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
