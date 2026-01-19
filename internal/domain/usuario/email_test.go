package usuario

import "testing"

// TestEmail_Validar testa a validação de emails.
func TestEmail_Validar(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"email valido", "USER@EMAIL.COM", false},
		{"email vazio", "", true},
		{"email invalido", "email@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NovoEmail(tt.email)
			err := e.ValidarEmail()

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}
