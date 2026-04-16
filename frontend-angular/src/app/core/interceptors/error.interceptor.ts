import { HttpInterceptorFn, HttpErrorResponse } from "@angular/common/http";
import { catchError, throwError } from "rxjs";

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  return next(req).pipe(
    catchError((err: HttpErrorResponse) => {
      if (err.status === 0) {
        console.warn("Backend offline ou inacessível");
      }
      const mensagem =
        err.error?.erro ||
        err.error?.detalhe ||
        err.error?.message ||
        mensagemPorStatus(err.status);
      return throwError(() => new Error(mensagem));
    }),
  );
};

function mensagemPorStatus(status: number): string {
  switch (status) {
    case 0:
      return "Serviço indisponível. Verifique sua conexão.";
    case 400:
      return "Dados inválidos. Verifique os campos.";
    case 404:
      return "Recurso não encontrado.";
    case 409:
      return "Conflito de concorrência. Tente novamente.";
    case 422:
      return "Operação não permitida.";
    case 500:
      return "Erro interno do servidor.";
    case 503:
      return "Serviço temporariamente indisponível.";
    default:
      return `Erro inesperado (${status}).`;
  }
}
