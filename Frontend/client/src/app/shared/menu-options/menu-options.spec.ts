import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MenuOptions } from './menu-options';

describe('MenuOptions', () => {
  let component: MenuOptions;
  let fixture: ComponentFixture<MenuOptions>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MenuOptions]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MenuOptions);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
