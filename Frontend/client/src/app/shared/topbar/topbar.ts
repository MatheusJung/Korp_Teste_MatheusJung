import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-topbar',
  imports: [],
  templateUrl: './topbar.html',
  styleUrl: './topbar.scss',
})
export class Topbar {
    @Output() toggleSidebar = new EventEmitter<void>();

  onToggleClick() {
    this.toggleSidebar.emit(); // envia o evento para o pai
  }
}
