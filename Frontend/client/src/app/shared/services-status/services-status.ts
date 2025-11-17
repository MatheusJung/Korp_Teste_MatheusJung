import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { Subscription, timer, switchMap } from 'rxjs';
import { HealthService } from '../../core/healthService';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-service-status',
  imports: [CommonModule],
  template: `
    <div [ngClass]="statusClass">
      <span class="dot"></span> {{ statusText }}
    </div>
  `,
  styles: [`
    .dot {
      height: 10px;
      width: 10px;
      border-radius: 50%;
      display: inline-block;
      margin-right: 5px;
    }
    .online .dot { background-color: green; }
    .offline .dot { background-color: red; }
  `]
})
export class ServiceStatus implements OnInit, OnDestroy {
  @Input() healthUrl!: string;      // URL do health check do serviço
  @Input() refreshMs: number = 5000; // intervalo de atualização

  statusText = 'Carregando...';
  statusClass = '';
  private subscription!: Subscription;

  constructor(private healthService: HealthService) {}

  ngOnInit() {
    if (!this.healthUrl) throw new Error('O atributo healthUrl é obrigatório');

    // Dispara imediatamente e depois a cada refreshMs
    this.subscription = timer(0, this.refreshMs).pipe(
      switchMap(() => this.healthService.getServiceStatus(this.healthUrl))
    ).subscribe(isOnline => this.updateStatus(isOnline));
  }

  private updateStatus(isOnline: boolean) {
    this.statusText = isOnline ? 'ONLINE' : 'OFFLINE';
    this.statusClass = isOnline ? 'online' : 'offline';
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
