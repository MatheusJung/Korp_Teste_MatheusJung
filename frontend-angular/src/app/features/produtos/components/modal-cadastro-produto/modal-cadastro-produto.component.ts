import { Component, EventEmitter, Output, Input, inject } from "@angular/core";
import { CommonModule } from "@angular/common";
import { ReactiveFormsModule, FormBuilder, Validators } from "@angular/forms";
import { AbstractControl, ValidationErrors } from "@angular/forms";

function inteiro(control: AbstractControl): ValidationErrors | null {
  const v = control.value;
  if (v === null || v === undefined || v === "") return null;
  return Number.isInteger(Number(v)) ? null : { inteiro: true };
}

@Component({
  selector: "app-modal-cadastro-produto",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: "./modal-cadastro-produto.component.html",
})
export class ModalCadastroProdutoComponent {
  private fb = inject(FormBuilder);

  @Input() salvando = false;
  @Input() erro: string | null = null;

  @Output() fechar = new EventEmitter<void>();
  @Output() salvar = new EventEmitter<any>();

  form = this.fb.group({
    codigo: [
      "",
      [
        Validators.required,
        Validators.pattern(/^\d{5}$/), // exatamente 5 dígitos
      ],
    ],
    descricao: [
      "",
      [
        Validators.required,
        Validators.pattern(/^[a-zA-ZÀ-ÿ0-9 .,\-()/]+$/),
        Validators.minLength(3),
        Validators.maxLength(30),
      ],
    ],
    saldoInicial: [0, [Validators.required, Validators.min(0), inteiro]],
  });

  submit() {
    if (this.form.invalid) return;
    this.salvar.emit(this.form.value);
  }
}
