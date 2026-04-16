import {
  Component,
  OnInit,
  Output,
  EventEmitter,
  inject,
  signal,
} from "@angular/core";
import { CommonModule } from "@angular/common";
import {
  ReactiveFormsModule,
  FormBuilder,
  FormArray,
  Validators,
  AbstractControl,
  ValidationErrors,
} from "@angular/forms";
import { ProdutoService } from "../../../../core/services/produto.service";
import { NotaFiscalService } from "../../../../core/services/nota-fiscal.service";
import { NotaFiscal, Produto } from "../../../../core/models/models";
import { finalize } from "rxjs";

function inteiro(control: AbstractControl): ValidationErrors | null {
  const v = control.value;
  if (v === null || v === undefined || v === "") return null;
  return Number.isInteger(Number(v)) ? null : { inteiro: true };
}

function maxSaldo(getSaldo: () => number | null) {
  return (control: AbstractControl): ValidationErrors | null => {
    const valor = Number(control.value);
    const saldo = getSaldo();

    if (!saldo && saldo !== 0) return null;
    if (!valor) return null;

    return valor > saldo ? { saldoExcedido: true } : null;
  };
}

@Component({
  selector: "app-modal-criar-nota",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: "./modal-criar-nota.component.html",
  styleUrls: ["./modal-criar-nota.component.scss"],
})
export class ModalCriarNotaComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly produtoSvc = inject(ProdutoService);
  private readonly notaSvc = inject(NotaFiscalService);

  @Output() readonly notaCriada = new EventEmitter<NotaFiscal>();
  @Output() readonly fechado = new EventEmitter<void>();

  readonly produtos = signal<Produto[]>([]);
  readonly salvando = signal(false);
  readonly erro = signal<string | null>(null);

  readonly form = this.fb.group({ itens: this.fb.array([]) });

  get itens() {
    return this.form.get("itens") as FormArray;
  }

  getItemGroup(i: number) {
    return this.itens.at(i) as ReturnType<typeof this.fb.group>;
  }

  getProdutoSelecionado(i: number) {
    const produtoId = this.getItemGroup(i).get("produtoId")?.value;
    return this.produtos().find((p) => p.id === produtoId) || null;
  }

  ngOnInit(): void {
    this.produtoSvc.listarTodos().subscribe((ps) => this.produtos.set(ps));
    this.adicionarItem();
  }

  adicionarItem(): void {
    const group = this.fb.group({
      produtoId: ["", Validators.required],
      quantidade: [
        null as number | null,
        [Validators.required, Validators.min(1), inteiro],
      ],
    });

    const quantidadeCtrl = group.get("quantidade")!;
    const produtoCtrl = group.get("produtoId")!;

    quantidadeCtrl.addValidators(
      maxSaldo(() => {
        const produtoId = produtoCtrl.value;
        const p = this.produtos().find((x) => x.id === produtoId);
        return p?.saldo ?? null;
      }),
    );

    produtoCtrl.valueChanges.subscribe(() => {
      quantidadeCtrl.updateValueAndValidity();
      this.validarProdutosDuplicados();
    });

    this.itens.push(group);
    this.validarProdutosDuplicados();
  }

  removerItem(i: number): void {
    if (this.itens.length > 1) {
      this.itens.removeAt(i);
      this.validarProdutosDuplicados();
    }
  }

  fechar(): void {
    this.fechado.emit();
  }

  submit(): void {
    this.validarProdutosDuplicados();

    if (this.form.invalid || this.itens.length === 0) return;

    this.salvando.set(true);
    this.erro.set(null);

    const itens = this.itens.controls.map((group) => {
      const produtoId = group.get("produtoId")?.value as string;
      const produto = this.produtos().find((p) => p.id === produtoId);
      const produtoDescricao = produto?.descricao ?? "";
      const produtoCodigo = produto?.codigo ?? "";
      const quantidade = group.get("quantidade")?.value as number;

      return {
        produtoId,
        produtoCodigo,
        produtoDescricao,
        quantidade,
      };
    });

    this.notaSvc
      .criar({ itens })
      .pipe(finalize(() => this.salvando.set(false)))
      .subscribe({
        next: (nota) => this.notaCriada.emit(nota),
        error: (e: Error) => this.erro.set(e.message),
      });
  }

  isProdutoJaSelecionado(produtoId: string, indexAtual: number): boolean {
    return this.itens.controls.some((control, index) => {
      if (index === indexAtual) return false;
      return control.get("produtoId")?.value === produtoId;
    });
  }

  private validarProdutosDuplicados(): void {
    const idsSelecionados = this.itens.controls
      .map((control) => control.get("produtoId")?.value)
      .filter((id) => !!id);

    const duplicados = idsSelecionados.filter(
      (id, index) => idsSelecionados.indexOf(id) !== index,
    );

    this.itens.controls.forEach((control) => {
      const produtoCtrl = control.get("produtoId");
      if (!produtoCtrl) return;

      const valor = produtoCtrl.value;
      const errosAtuais = produtoCtrl.errors || {};

      if (valor && duplicados.includes(valor)) {
        produtoCtrl.setErrors({ ...errosAtuais, duplicado: true });
      } else {
        const { duplicado, ...restante } = errosAtuais;
        produtoCtrl.setErrors(Object.keys(restante).length ? restante : null);
      }
    });
  }
}
