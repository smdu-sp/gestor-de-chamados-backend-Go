package categoria_permissao

import (
	"testing"
	"time"

	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// TestNovoCategoriaPermissao_Valido testa a criação de uma CategoriaPermissao válida.
func TestNovoCategoriaPermissao_Valido(t *testing.T) {
	permissao := usr.PermADM
	cp, err := Novo("cat-123", "user-456", permissao)
	if err != nil {
		t.Fatalf("Esperava nenhuma erro, mas recebeu: %v", err)
	}

	if cp.categoriaID != "cat-123" {
		t.Errorf("Esperava CategoriaID %s, mas recebeu %s", "cat-123", cp.categoriaID)
	}
	if cp.usuarioID != "user-456" {
		t.Errorf("Esperava UsuarioID %s, mas recebeu %s", "user-456", cp.usuarioID)
	}
	if cp.permissao != permissao {
		t.Errorf("Esperava Permissao %s, mas recebeu %s", permissao, cp.permissao)
	}
}

// TestCarregarDoDB testa a função CarregarDoDB.
func TestCarregarDoDB(t *testing.T) {
	now := time.Now()
	dbStruct := CategoriaPermissaoDB{
		CategoriaID:  "cat-123",
		UsuarioID:    "user-456",
		Permissao:    string(usr.PermADM),
		CriadoEm:     now,
		AtualizadoEm: now,
	}

	cp := CarregarDoDB(dbStruct)

	if cp.categoriaID != dbStruct.CategoriaID {
		t.Errorf("Esperava CategoriaID %s, mas recebeu %s", dbStruct.CategoriaID, cp.categoriaID)
	}
	if cp.usuarioID != dbStruct.UsuarioID {
		t.Errorf("Esperava UsuarioID %s, mas recebeu %s", dbStruct.UsuarioID, cp.usuarioID)
	}
	if cp.permissao != usr.Permissao(dbStruct.Permissao) {
		t.Errorf("Esperava Permissao %s, mas recebeu %s", dbStruct.Permissao, cp.permissao)
	}
	if !cp.criadoEm.Equal(dbStruct.CriadoEm) {
		t.Errorf("Esperava CriadoEm %v, mas recebeu %v", dbStruct.CriadoEm, cp.criadoEm)
	}
	if !cp.atualizadoEm.Equal(dbStruct.AtualizadoEm) {
		t.Errorf("Esperava AtualizadoEm %v, mas recebeu %v", dbStruct.AtualizadoEm, cp.atualizadoEm)
	}
}

// TestCategoriaPermissao_Validar testa a validação da CategoriaPermissao.
func TestCategoriaPermissao_Validar(t *testing.T) {
	tests := []struct {
		name        string
		categoriaID string
		usuarioID   string
		permissao   usr.Permissao
		wantErr     bool
	}{
		{
			name:        "válido",
			categoriaID: "cat-123",
			usuarioID:   "user-456",
			permissao:   usr.PermADM,
			wantErr:     false,
		},
		{
			name:        "categoriaID vazio",
			categoriaID: "",
			usuarioID:   "user-456",
			permissao:   usr.PermADM,
			wantErr:     true,
		},
		{
			name:        "usuarioID vazio",
			categoriaID: "cat-123",
			usuarioID:   "",
			permissao:   usr.PermADM,
			wantErr:     true,
		},
		{
			name:        "permissao inválida",
			categoriaID: "cat-123",
			usuarioID:   "user-456",
			permissao:   usr.Permissao("INVALIDA"),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &CategoriaPermissao{
				categoriaID:  tt.categoriaID,
				usuarioID:    tt.usuarioID,
				permissao:    tt.permissao,
				criadoEm:     time.Now(),
				atualizadoEm: time.Now(),
			}
			err := cp.Validar()
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoriaPermissao.Validar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCategoriaPermissao_AtualizarDados testa a atualização dos dados da CategoriaPermissao.
func TestCategoriaPermissao_AtualizarDados(t *testing.T) {
	cpm, _ := Novo("cat-123", "user-456", usr.PermTEC)

	time.Sleep(1 * time.Second) // garante que o tempo de atualizadoEm será diferente

	err := cpm.AtualizarDados(AtualizarParams{
		Permissao: usr.PermADM,
	})
	if err != nil {
		t.Fatalf("Esperava nenhuma erro, mas recebeu: %v", err)
	}

	if cpm.permissao != usr.PermADM {
		t.Errorf("Esperava Permissao %s, mas recebeu %s", usr.PermADM, cpm.permissao)
	}

	if !cpm.atualizadoEm.After(cpm.CriadoEm()) {
		t.Errorf("Esperava AtualizadoEm atualizado, mas não foi")
	}
}
