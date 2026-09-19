import { HttpErrorResponse } from '@angular/common/http';

export function getApiErrorMessage(error: unknown, fallback: string): string {
  if (!(error instanceof HttpErrorResponse)) {
    return fallback;
  }

  if (error.status === 0) {
    return 'Não foi possível conectar ao servidor. Confirme se os serviços estão em execução.';
  }

  if (error.status === 503) {
    return 'O serviço de estoque está indisponível. Aguarde a recuperação e tente a mesma nota novamente.';
  }

  if (error.status === 504) {
    return 'O serviço de estoque demorou para responder. Tente a mesma nota novamente.';
  }

  if (error.status === 502) {
    return 'O serviço de estoque retornou uma resposta inválida. Tente novamente.';
  }

  const message = error.error?.message;
  return typeof message === 'string' && message.trim() ? message : fallback;
}
