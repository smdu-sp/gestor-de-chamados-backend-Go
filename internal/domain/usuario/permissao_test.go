package usuario

import "testing"

// TestValidarPermissao testa a funcao ValidarPermissao
func TestValidarPermissao(t *testing.T) {
	tests := []struct {
		name      string
		permissao Permissao
		wantErr   bool
	}{
		{"ADM valida", PermADM, false},
		{"TEC valida", PermTEC, false},
		{"permissao vazia", "", true},
		{"permissao invalida", Permissao("XXX"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarPermissao(tt.permissao)

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}
