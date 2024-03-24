package model_test

import (
	"github.com/Waelson/go-o11y/internal/model"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		errorVar       error
		expectedErrMsg string
	}{
		{
			name:           "InvalidCepError",
			errorVar:       model.InvalidCepError,
			expectedErrMsg: "cep vazio",
		},
		{
			name:           "CepNotFoundError",
			errorVar:       model.CepNotFoundError,
			expectedErrMsg: "cep nao encontrado",
		},
		{
			name:           "InternalError",
			errorVar:       model.InternalError,
			expectedErrMsg: "ocorreu um erro interno",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errorVar.Error() != tc.expectedErrMsg {
				t.Errorf("erro inesperado, recebido: %v, esperado: %v", tc.errorVar, tc.expectedErrMsg)
			}
		})
	}
}
