import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InvoiceTable } from './invoice-table';

describe('InvoiceTable', () => {
  let component: InvoiceTable;
  let fixture: ComponentFixture<InvoiceTable>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [InvoiceTable]
    })
    .compileComponents();

    fixture = TestBed.createComponent(InvoiceTable);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
