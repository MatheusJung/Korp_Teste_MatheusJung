import {
  Component,
  EventEmitter,
  Input,
  Output,
  inject,
  OnChanges,
  SimpleChanges,
} from "@angular/core";
import { CommonModule } from "@angular/common";
import {
  ReactiveFormsModule,
  FormBuilder,
  Validators,
  AbstractControl,
  ValidationErrors,
} from "@angular/forms";
import { Produto } from "../../../../core/models/models";

function inteiro(control: AbstractControl): ValidationErrors | null {
  const v = control.value;
  if (v === null || v === undefined || v === "") return null;
  return Number.isInteger(Number(v)) ? null : { inteiro: true };
}

@Component({
  selector: "app-modal-entrada-produto",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: "./modal-entrada-produto.component.html",
  styleUrls: ["./modal-entrada-produto.component.scss"],
})
export class ModalEntradaProdutoComponent implements OnChanges {
  private fb = inject(FormBuilder);

  @Input() produto: Produto | null = null;
  @Input() salvando = false;
  @Input() erro: string | null = null;

  @Output() fechar = new EventEmitter<void>();
  @Output() salvar = new EventEmitter<{ quantidade: number }>();

  form = this.fb.group({
    quantidade: [
      null as number | null,
      [Validators.required, Validators.min(1), inteiro],
    ],
  });

  ngOnChanges(changes: SimpleChanges): void {
    // Quando trocar o produto → limpa formulário
    if (changes["produto"] && !changes["produto"].firstChange) {
      this.resetForm();
    }

    // Se veio erro do backend → reseta input
    if (changes["erro"] && this.erro) {
      this.form.get("quantidade")?.reset();
      this.form.markAsPristine();
      this.form.markAsUntouched();
    }
  }

  submit(): void {
    if (this.form.invalid || this.salvando) return;

    const quantidade = Number(this.form.value.quantidade);

    this.salvar.emit({ quantidade });
  }

  fecharModal(): void {
    this.resetForm();
    this.fechar.emit();
  }

  private resetForm(): void {
    this.form.reset();
  }
}
