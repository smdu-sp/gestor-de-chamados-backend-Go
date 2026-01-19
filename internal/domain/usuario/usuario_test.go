package usuario

import "testing"

// TestNovoUsuario testa a criação de novos usuários.
func TestNovoUsuario(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		nome      string
		login     string
		email     Email
		permissao Permissao
		wantErr   bool
	}{
		{
			name:      "cria usuario valido",
			id:        "user-123",
			nome:      "Rogério Silva",
			login:     "rogerio",
			email:     NovoEmail("rogerio@email.com"),
			permissao: PermADM,
			wantErr:   false,
		},
		{
			name:      "erro com nome curto",
			id:        "user-123",
			nome:      "Ro",
			login:     "rogerio",
			email:     NovoEmail("rogerio@email.com"),
			permissao: PermADM,
			wantErr:   true,
		},
		{
			name:      "erro com email invalido",
			id:        "user-123",
			nome:      "Rogério",
			login:     "rogerio",
			email:     NovoEmail("email-invalido"),
			permissao: PermADM,
			wantErr:   true,
		},
		{
			name:      "erro com permissao invalida",
			id:        "user-123",
			nome:      "Rogério",
			login:     "rogerio",
			email:     NovoEmail("rogerio@email.com"),
			permissao: Permissao("XXX"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := Novo(tt.id, tt.nome, tt.login, tt.email, tt.permissao, nil)

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if !tt.wantErr && u == nil {
				t.Fatal("usuario nao deveria ser nil")
			}
		})
	}
}
