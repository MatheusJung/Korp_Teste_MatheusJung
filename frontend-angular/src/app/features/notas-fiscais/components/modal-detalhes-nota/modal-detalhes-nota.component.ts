import {
  Component,
  inject,
  input,
  output,
  OnInit,
  signal,
} from "@angular/core";
import { CommonModule } from "@angular/common";
import { NotaFiscal } from "../../../../core/models/models";
import { ProdutoService } from "../../../../core/services/produto.service";
import { NotaFiscalService } from "../../../../core/services/nota-fiscal.service";
import { Produto } from "../../../../core/models/models";
import { finalize } from "rxjs";
import { HealthService } from "../../../../core/services/health.service";

@Component({
  selector: "app-modal-detalhes-nota",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./modal-detalhes-nota.component.html",
  styleUrl: "./modal-detalhes-nota.component.scss",
})
export class ModalDetalhesNotaComponent implements OnInit {
  private readonly produtoSvc = inject(ProdutoService);
  private readonly notaSvc = inject(NotaFiscalService);

  private readonly healthService = inject(HealthService);
  readonly health = this.healthService.health;

  // Input obrigatório — nota a exibir
  readonly nota = input.required<NotaFiscal>();

  // Eventos para o pai
  readonly fechado = output<void>();
  readonly impressaoIniciada = output<NotaFiscal>();

  readonly produtos = signal<Produto[]>([]);
  readonly imprimindo = signal(false);
  readonly erro = signal<string | null>(null);

  ngOnInit(): void {
    this.produtoSvc.listarTodos().subscribe((ps) => this.produtos.set(ps));
  }

  nomeProduto(produtoId: string): string {
    const p = this.produtos().find((x) => x.id === produtoId);
    return p ? `${p.codigo} — ${p.descricao}` : "Carregando...";
  }

  fechar(): void {
    this.fechado.emit();
  }

  imprimir(): void {
    this.imprimindo.set(true);
    this.erro.set(null);

    this.notaSvc
      .imprimir(this.nota().id)
      .pipe(finalize(() => this.imprimindo.set(false)))
      .subscribe({
        next: () => {
          const notaAtualizada: NotaFiscal = {
            ...this.nota(),
            status: "Processando",
          };

          this.impressaoIniciada.emit(notaAtualizada);
          this.fechar();
        },
        error: (e: Error) => this.erro.set(e.message),
      });
  }
}
