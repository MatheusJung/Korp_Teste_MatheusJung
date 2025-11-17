import { Component , Input } from '@angular/core';
import { MenuOptions } from "../menu-options/menu-options";

@Component({
  selector: 'app-sidebar',
  imports: [MenuOptions],
  templateUrl: './sidebar.html',
  styleUrl: './sidebar.scss',
})
export class Sidebar {
  @Input() isOpen = false;
}

