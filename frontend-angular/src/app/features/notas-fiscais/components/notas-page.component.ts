import { Component, OnInit, OnDestroy, inject, signal } from "@angular/core";
import { CommonModule } from "@angular/common";
import { NotaFiscalService } from "../../../core/services/nota-fiscal.service";
import { HealthService } from "../../../core/services/health.service";
import { NotaFiscal } from "../../../core/models/models";
import { finalize, Subscription } from "rxjs";
import { ModalCriarNotaComponent } from "./modal-criar-nota/modal-criar-nota.component";
import { ModalDetalhesNotaComponent } from "./modal-detalhes-nota/modal-detalhes-nota.component";

@Component({
  selector: "app-notas-page",
  standalone: true,
  imports: [CommonModule, ModalCriarNotaComponent, ModalDetalhesNotaComponent],
  templateUrl: "./notas-page.component.html",
  styleUrls: ["./notas-page.component.scss"],
})
export class NotasPageComponent implements OnInit, OnDestroy {
  private readonly svc = inject(NotaFiscalService);
  private readonly healthService = inject(HealthService);
  private pollSub?: Subscription;

  readonly health = this.healthService.health;
  readonly notas = signal<NotaFiscal[]>([]);
  readonly carregando = signal(false);
  readonly notaProcessando = signal<NotaFiscal | null>(null);
  readonly erro = signal<string | null>(null);

  readonly modalCriar = signal(false);
  readonly modalDetalhe = signal(false);
  readonly notaDetalhe = signal<NotaFiscal | null>(null);

  readonly page = signal(1);
  readonly pageSize = signal(10);
  readonly totalItems = signal(0);
  readonly totalPages = signal(0);

  readonly searchInput = signal("");
  readonly search = signal("");
  readonly sortBy = signal("codigo");
  readonly sortDirection = signal<"asc" | "desc">("asc");

  formatarNumeroNota(numero: number): string {
    const padded = numero.toString().padStart(6, "0");
    return padded.replace(/\B(?=(\d{3})+(?!\d))/g, ".");
  }

  ngOnInit(): void {
    this.carregar();
  }

  ngOnDestroy(): void {
    this.pollSub?.unsubscribe();
  }

  carregar(): void {
    this.carregando.set(true);
    this.erro.set(null);

    this.svc
      .listarPaginado({
        page: this.page(),
        pageSize: this.pageSize(),
        search: this.search(),
        sortBy: this.sortBy(),
        sortDirection: this.sortDirection(),
      })
      .pipe(finalize(() => this.carregando.set(false)))
      .subscribe({
        next: (res) => {
          this.notas.set(res.items);
          this.totalPages.set(res.totalPages);
          this.totalItems.set(res.totalItems);
        },
        error: (e: Error) => this.erro.set(e.message),
      });
  }

  pesquisar(): void {
    this.search.set(this.searchInput().trim());
    this.page.set(1);
    this.carregar();
  }

  alterarBusca(valor: string): void {
    this.search.set(valor);
    this.page.set(1);
    this.carregar();
  }

  ordenarPor(coluna: string): void {
    if (this.sortBy() === coluna) {
      this.sortDirection.set(this.sortDirection() === "asc" ? "desc" : "asc");
    } else {
      this.sortBy.set(coluna);
      this.sortDirection.set("asc");
    }

    this.page.set(1);
    this.carregar();
  }

  alterarOrdenacao(sortBy: string): void {
    this.sortBy.set(sortBy);
    this.page.set(1);
    this.carregar();
  }

  alterarDirecao(sortDirection: "asc" | "desc"): void {
    this.sortDirection.set(sortDirection);
    this.page.set(1);
    this.carregar();
  }

  proximaPagina(): void {
    if (this.page() < this.totalPages()) {
      this.page.update((p) => p + 1);
      this.carregar();
    }
  }

  paginaAnterior(): void {
    if (this.page() > 1) {
      this.page.update((p) => p - 1);
      this.carregar();
    }
  }

  verDetalhe(n: NotaFiscal): void {
    this.notaDetalhe.set(n);
    this.modalDetalhe.set(true);
  }

  fecharTudo(): void {
    this.modalCriar.set(false);
    this.modalDetalhe.set(false);
    this.notaDetalhe.set(null);
  }

  // Chamado pelo ModalCriarNotaComponent via (notaCriada)
  onNotaCriada(nota: NotaFiscal): void {
    this.notas.update((ns) => [nota, ...ns]);
    this.fecharTudo();
  }

  // Chamado pelo ModalDetalhesNotaComponent via (impressaoIniciada)
  onImpressaoIniciada(nota: NotaFiscal): void {
    this.notas.update((ns) => ns.map((n) => (n.id === nota.id ? nota : n)));
    this.notaProcessando.set(nota);
    this.fecharTudo();
    this.iniciarPolling(nota.id);
  }

  private iniciarPolling(id: string): void {
    this.pollSub?.unsubscribe();

    let pdfAberto = false;

    this.pollSub = this.svc.aguardarProcessamento(id).subscribe({
      next: (nota) => {
        this.notas.update((ns) => ns.map((n) => (n.id === nota.id ? nota : n)));

        if (nota.status === "Fechada" && !pdfAberto) {
          pdfAberto = true;
          this.notaProcessando.set(null);
          window.open(
            `http://localhost:5002/api/notas/${nota.id}/pdf`,
            "_blank",
          );
          this.pollSub?.unsubscribe();
          return;
        }

        if (nota.status !== "Processando") {
          this.notaProcessando.set(null);
        }
      },
      complete: () => this.notaProcessando.set(null),
    });
  }
}
