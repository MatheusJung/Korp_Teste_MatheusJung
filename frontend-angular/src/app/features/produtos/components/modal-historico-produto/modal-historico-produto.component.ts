import { Component, Input, Output, EventEmitter } from "@angular/core";
import { CommonModule } from "@angular/common";
import { Movimentacao, Produto } from "../../../../core/models/models";

@Component({
  selector: "app-modal-historico-produto",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./modal-historico-produto.component.html",
  styleUrls: ["./modal-historico-produto.component.scss"],
})
export class ModalHistoricoProdutoComponent {
  @Input() produto: Produto | null = null;
  @Input() movimentacoes: Movimentacao[] = [];
  @Input() carregando: boolean = false;

  @Output() fecharModal = new EventEmitter<void>();

  fechar() {
    this.fecharModal.emit();
  }
}
