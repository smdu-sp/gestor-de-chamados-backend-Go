package usuario

import (
	"testing"
	"time"
)

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

// TestCarregarDoBD testa a função CarregarDoBD.
func TestCarregarDoBD(t *testing.T) {
	now := time.Now()
	usuarioDB := UsuarioDB{
		ID:           "user-123",
		Nome:         "Rogério Silva",
		Login:        "rogerio",
		Email:        "rogerio@email.com",
		Permissao:    "ADM",
		Status:       true,
		Avatar:       nil,
		UltimoLogin:  now,
		CriadoEm:     now,
		AtualizadoEm: now,
	}

	u := CarregarDoBD(usuarioDB)

	if u.ID() != usuarioDB.ID {
		t.Errorf("esperava ID %s, mas recebeu %s", usuarioDB.ID, u.ID())
	}

	if u.Nome() != usuarioDB.Nome {
		t.Errorf("esperava Nome %s, mas recebeu %s", usuarioDB.Nome, u.Nome())
	}

	if u.Login() != usuarioDB.Login {
		t.Errorf("esperava Login %s, mas recebeu %s", usuarioDB.Login, u.Login())
	}

	if u.Email() != usuarioDB.Email {
		t.Errorf("esperava Email %s, mas recebeu %s", usuarioDB.Email, u.Email())
	}

	if u.Permissao().String() != usuarioDB.Permissao {
		t.Errorf("esperava Permissao %s, mas recebeu %s", usuarioDB.Permissao, u.Permissao().String())
	}

	if u.Status() != usuarioDB.Status {
		t.Errorf("esperava Status %v, mas recebeu %v", usuarioDB.Status, u.Status())
	}

	if u.UltimoLogin() != usuarioDB.UltimoLogin {
		t.Errorf("esperava UltimoLogin %v, mas recebeu %v", usuarioDB.UltimoLogin, u.UltimoLogin())
	}

	if !u.CriadoEm().Equal(usuarioDB.CriadoEm) {
		t.Errorf("esperava CriadoEm %v, mas recebeu %v", usuarioDB.CriadoEm, u.CriadoEm())
	}

	if !u.AtualizadoEm().Equal(usuarioDB.AtualizadoEm) {
		t.Errorf("esperava AtualizadoEm %v, mas recebeu %v", usuarioDB.AtualizadoEm, u.AtualizadoEm())
	}
}

// TestUsuario_Validar testa o método Validar do Usuario.
func TestUsuario_Validar(t *testing.T) {
	tests := []struct {
		name    string
		usuario Usuario
		wantErr bool
	}{
		{
			name: "usuario valido",
			usuario: Usuario{
				id:        "user-123",
				nome:      "Rogério Silva",
				login:     "rogerio",
				email:     NovoEmail("rogerio@email.com"),
				permissao: PermTEC,
				status:    true,
			},
			wantErr: false,
		},
		{
			name: "erro com nome vazio",
			usuario: Usuario{
				id:        "user-123",
				nome:      "",
				login:     "rogerio",
				email:     NovoEmail("rogerio@email.com"),
				permissao: PermTEC,
				status:    true,
			},
			wantErr: true,
		},
		{
			name: "erro com login vazio",
			usuario: Usuario{
				id:        "user-123",
				nome:      "Rogério Silva",
				login:     "",
				email:     NovoEmail("rogerio@email.com"),
				permissao: PermTEC,
				status:    true,
			},
			wantErr: true,
		},
		{
			name: "erro com email invalido",
			usuario: Usuario{
				id:        "user-123",
				nome:      "Rogério Silva",
				login:     "rogerio",
				email:     NovoEmail("email-invalido"),
				permissao: PermTEC,
				status:    true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.usuario.Validar()

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestAtualizarDados testa o método AtualizarDados do Usuario.
func TestAtualizarDados(t *testing.T) {
	usuario, _ := Novo(
		"user-123",
		"Rogério Silva",
		"rogerio",
		NovoEmail("rogerio@email.com"),
		PermTEC,
		nil,
	)

	time.Sleep(1 * time.Second) // garante que o tempo de atualizadoEm será diferente

	params := AtualizarParams{
		Nome:  ptrString("Rogério S."),
		Login: ptrString("rogerios"),
		Email: ptrEmail(NovoEmail("rogerio.s@email.com")),
	}
	err := usuario.AtualizarDados(params)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar dados: %v", err)
	}

	if usuario.Nome() != "Rogério S." {
		t.Errorf("esperava Nome atualizado para 'Rogério S.', mas recebeu %s", usuario.Nome())
	}

	if usuario.Login() != "rogerios" {
		t.Errorf("esperava Login atualizado para 'rogerios', mas recebeu %s", usuario.Login())
	}

	if usuario.Email() != "rogerio.s@email.com" {
		t.Errorf("esperava Email atualizado para 'rogerio.s@email.com', mas recebeu %s", usuario.Email())
	}

	if !usuario.AtualizadoEm().After(usuario.CriadoEm()) {
		t.Errorf("esperava AtualizadoEm após CriadoEm")
	}
}

// TestAtivarDesativar testa os métodos Ativar e Desativar do Usuario.
func TestAtivarDesativar(t *testing.T) {
	usuario, _ := Novo(
		"user-123",
		"Rogério Silva",
		"rogerio",
		NovoEmail("rogerio@email.com"),
		PermTEC,
		nil,
	)

	usuario.Desativar()
	if usuario.Status() != false {
		t.Errorf("esperava Status false após desativar, mas recebeu %v", usuario.Status())
	}

	usuario.Ativar()
	if usuario.Status() != true {
		t.Errorf("esperava Status true após ativar, mas recebeu %v", usuario.Status())
	}
}

// TestAtualizarPermissao testa o método AtualizarPermissao do Usuario.
func TestAtualizarPermissao(t *testing.T) {
	usuario, _ := Novo(
		"user-123",
		"Rogério Silva",
		"rogerio",
		NovoEmail("rogerio@email.com"),
		PermTEC,
		nil,
	)
	err := usuario.AtualizarPermissao(PermADM)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar permissao: %v", err)
	}
	if usuario.Permissao() != PermADM {
		t.Errorf("esperava Permissao ADM após atualizar, mas recebeu %s", usuario.Permissao().String())
	}
}

// TestAtualizarPermissao_Invalida testa o método AtualizarPermissao do Usuario com uma permissão inválida.
func TestAtualizarPermissao_Invalida(t *testing.T) {
	usuario, _ := Novo(
		"user-123",
		"Rogério Silva",
		"rogerio",
		NovoEmail("rogerio@email.com"),
		PermTEC,
		nil,
	)

	err := usuario.AtualizarPermissao("PermissaoInvalida")
	if err == nil {
		t.Fatalf("esperava erro ao atualizar permissao inválida, mas recebeu nil")
	}
}

// TestAtualizarUltimoLogin testa o método AtualizarUltimoLogin do Usuario.
func TestAtualizarUltimoLogin(t *testing.T) {
	usuario, _ := Novo(
		"user-123",
		"Rogério Silva",
		"rogerio",
		NovoEmail("rogerio@email.com"),
		PermTEC,
		nil,
	)

	antes := time.Now()
	time.Sleep(1 * time.Second) // garante que o tempo será diferente

	usuario.AtualizarUltimoLogin()
	depois := time.Now()

	if usuario.UltimoLogin().Before(antes) || usuario.UltimoLogin().After(depois) {
		t.Errorf("UltimoLogin não foi atualizado corretamente, recebeu %v", usuario.UltimoLogin())
	}
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// ptrString retorna um ponteiro para a string fornecida.
func ptrString(s string) *string {
	return &s
}

// ptrEmail retorna um ponteiro para o Email fornecido.
func ptrEmail(e Email) *Email {
	return &e
}