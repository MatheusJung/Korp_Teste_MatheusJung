import { Component, OnInit, inject, signal } from "@angular/core";
import { CommonModule } from "@angular/common";
import { finalize } from "rxjs";

import { ProdutoService } from "../../../core/services/produto.service";
import { HealthService } from "../../../core/services/health.service";
import { Produto, Movimentacao } from "../../../core/models/models";

import { ModalCadastroProdutoComponent } from "./modal-cadastro-produto/modal-cadastro-produto.component";
import { ModalEntradaProdutoComponent } from "./modal-entrada-produto/modal-entrada-produto.component";
import { ModalHistoricoProdutoComponent } from "./modal-historico-produto/modal-historico-produto.component";

@Component({
  selector: "app-produtos-page",
  standalone: true,
  imports: [
    CommonModule,
    ModalCadastroProdutoComponent,
    ModalEntradaProdutoComponent,
    ModalHistoricoProdutoComponent,
  ],
  templateUrl: "./produtos-page.component.html",
  styleUrl: "./produtos-page.component.scss",
})
export class ProdutosPageComponent implements OnInit {
  private readonly svc = inject(ProdutoService);
  private readonly healthService = inject(HealthService);

  readonly produtos = signal<Produto[]>([]);
  readonly movimentacoes = signal<Movimentacao[]>([]);
  readonly produtoSelecionado = signal<Produto | null>(null);

  readonly carregando = signal(false);
  readonly carregandoHistorico = signal(false);
  readonly salvando = signal(false);
  readonly erro = signal<string | null>(null);

  readonly modalCadastro = signal(false);
  readonly modalEntrada = signal(false);
  readonly modalHistorico = signal(false);

  readonly page = signal(1);
  readonly pageSize = signal(10);
  readonly totalItems = signal(0);
  readonly totalPages = signal(0);

  readonly searchInput = signal("");
  readonly search = signal("");
  readonly sortBy = signal("codigo");
  readonly sortDirection = signal<"asc" | "desc">("asc");

  readonly health = this.healthService.health;

  ngOnInit(): void {
    this.carregar();
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
          this.produtos.set(res.items);
          this.totalItems.set(res.totalItems);
          this.totalPages.set(res.totalPages);
        },
        error: (e: Error) => this.erro.set(e.message),
      });
  }

  pesquisar(): void {
    this.search.set(this.searchInput().trim());
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

  alterarBusca(valor: string): void {
    this.search.set(valor);
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

  abrirModalCadastro(): void {
    this.erro.set(null);
    this.modalCadastro.set(true);
  }

  abrirEntrada(produto: Produto): void {
    this.produtoSelecionado.set(produto);
    this.erro.set(null);
    this.modalEntrada.set(true);
  }

  abrirMovimentacoes(produto: Produto): void {
    this.produtoSelecionado.set(produto);
    this.movimentacoes.set([]);
    this.erro.set(null);
    this.modalHistorico.set(true);
    this.carregarMovimentacoes(produto.id);
  }

  fecharTudo(): void {
    this.modalCadastro.set(false);
    this.modalEntrada.set(false);
    this.modalHistorico.set(false);
    this.erro.set(null);
  }

  salvarProdutoFilho(data: {
    codigo: string;
    descricao: string;
    saldoInicial: number;
  }): void {
    this.salvando.set(true);
    this.erro.set(null);

    this.svc
      .criar({
        codigo: data.codigo.toUpperCase(),
        descricao: data.descricao,
        saldoInicial: data.saldoInicial ?? 0,
      })
      .pipe(finalize(() => this.salvando.set(false)))
      .subscribe({
        next: () => {
          this.fecharTudo();
          this.carregar();
        },
        error: (e: Error) => this.erro.set(e.message),
      });
  }

  salvarEntradaFilho(data: { quantidade: number }): void {
    const produto = this.produtoSelecionado();
    if (!produto) return;

    this.salvando.set(true);
    this.erro.set(null);

    this.svc
      .adicionarEntrada({
        produtoId: produto.id,
        quantidade: data.quantidade,
      })
      .pipe(finalize(() => this.salvando.set(false)))
      .subscribe({
        next: (atualizado: Produto) => {
          this.produtos.update((ps) =>
            ps.map((p) => (p.id === atualizado.id ? atualizado : p)),
          );

          this.produtoSelecionado.set(atualizado);

          if (this.modalHistorico()) {
            this.carregarMovimentacoes(atualizado.id);
          }

          this.fecharTudo();
        },
        error: (e: Error) => this.erro.set(e.message),
      });
  }

  private carregarMovimentacoes(produtoId: string): void {
    this.carregandoHistorico.set(true);

    this.svc
      .listarMovimentacoes(produtoId)
      .pipe(finalize(() => this.carregandoHistorico.set(false)))
      .subscribe({
        next: (ms) => this.movimentacoes.set(ms),
        error: (e: Error) => this.erro.set(e.message),
      });
  }
}
