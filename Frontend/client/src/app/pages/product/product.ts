import { Component} from '@angular/core';
import { Topbar } from "../../shared/topbar/topbar";
import { Sidebar } from "../../shared/sidebar/sidebar";
import { Footer } from "../../shared/footer/footer";
import { ProductTable } from "./product-table/product-table";

@Component({
  selector: 'app-product',
  imports: [Topbar, Sidebar, Footer, ProductTable],
  templateUrl: './product.html',
  styleUrl: './product.scss',
})
export class ProductComponent{
  showModal = false;
  sidebarOpen = false;

  openModal() {
    this.showModal = true;
  }

  closeModal() {
    this.showModal = false;
  }

  handleToggleSidebar(){
    this.sidebarOpen = !this.sidebarOpen;
  }
}
