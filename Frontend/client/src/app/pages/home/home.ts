import { Component } from '@angular/core';
import { Footer } from "../../shared/footer/footer";
import { Topbar } from "../../shared/topbar/topbar";
import { Sidebar } from "../../shared/sidebar/sidebar";
import { MenuOptions } from "../../shared/menu-options/menu-options";

@Component({
  selector: 'app-home',
  imports: [Footer, Topbar, Sidebar, MenuOptions],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home {
  sidebarOpen = false;

  handleToggleSidebar() {
    this.sidebarOpen = !this.sidebarOpen;
  }
}
