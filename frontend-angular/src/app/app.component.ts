import { Component, OnInit, OnDestroy, inject, signal } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HealthService } from "./core/services/health.service";
import { ProdutosPageComponent } from "./features/produtos/components/produtos-page.component";
import { NotasPageComponent } from "./features/notas-fiscais/components/notas-page.component";

type Tab = "produtos" | "notas";

@Component({
  selector: "app-root",
  standalone: true,
  imports: [CommonModule, ProdutosPageComponent, NotasPageComponent],
  template: `
    <!-- Health banner -->
    @if (hasIssue()) {
      <div class="health-banner">
        <span class="health-banner__icon">⚠</span>
        <span>
          @if (health().estoque !== "ok") {
            Serviço de estoque {{ health().estoque }}.
          }
          @if (health().faturamento !== "ok") {
            Serviço de faturamento {{ health().faturamento }}.
          }
          Algumas funcionalidades podem estar indisponíveis.
        </span>
      </div>
    }

    <!-- Header -->
    <header class="header">
      <div class="header__inner">
        <div class="header__brand">
          <div class="header__logo">NF</div>
          <div>
            <div class="header__title">Sistema de Notas Fiscais</div>
            <div class="header__sub">Gestão de estoque e faturamento</div>
          </div>
        </div>

        <div class="header__status">
          <div
            class="status-dot"
            [class]="'status-dot--' + health().estoque"
            title="Estoque"
          >
            <span class="badge badge--{{ health().estoque }}">Estoque</span>
          </div>
          <div
            class="status-dot"
            [class]="'status-dot--' + health().faturamento"
            title="Faturamento"
          >
            <span class="badge badge--{{ health().faturamento }}"
              >Faturamento</span
            >
          </div>
        </div>
      </div>
    </header>

    <!-- Tabs -->
    <div class="tabs-bar">
      <div class="tabs-bar__inner">
        <button
          class="tab"
          [class.tab--active]="tab() === 'produtos'"
          (click)="tab.set('produtos')"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="m7.5 4.27 9 5.15" />
            <path
              d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"
            />
            <path d="m3.3 7 8.7 5 8.7-5" />
            <path d="M12 22V12" />
          </svg>

          Produtos
        </button>
        <button
          class="tab"
          [class.tab--active]="tab() === 'notas'"
          (click)="tab.set('notas')"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path
              d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"
            />
            <polyline points="14 2 14 8 20 8" />
          </svg>
          Notas Fiscais
        </button>
      </div>
    </div>

    <!-- Conteúdo -->
    <main class="main">
      @if (tab() === "produtos") {
        <app-produtos-page />
      } @else {
        <app-notas-page />
      }
    </main>
  `,
  styleUrls: ["./app.component.scss"],
})
export class AppComponent implements OnInit, OnDestroy {
  private readonly healthService = inject(HealthService);

  readonly tab = signal<Tab>("produtos");
  readonly health = this.healthService.health;

  hasIssue() {
    const h = this.health();
    return h.estoque !== "ok" || h.faturamento !== "ok";
  }

  ngOnInit(): void {
    this.healthService.startMonitoring();
  }
  ngOnDestroy(): void {
    this.healthService.stopMonitoring();
  }
}
