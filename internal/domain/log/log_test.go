package log

import (
	"testing"
	"time"
)

// TestNovoLog_Sucesso testa a criação bem-sucedida de um novo log.
func TestNovoLog_Sucesso(t *testing.T) {
	id := "log-123"
	usuarioID := "user-456"
	acao := Criar
	entidade := "Subcategoria"
	detalhes := "Subcategoria criada com sucesso"

	log, err := Novo(id, usuarioID, acao, entidade, detalhes)

	if err != nil {
		t.Fatalf("não esperava erro, recebeu: %v", err)
	}

	if log.ID() != id {
		t.Errorf("ID esperado %s, recebeu %s", id, log.ID())
	}

	if log.UsuarioID() != usuarioID {
		t.Errorf("UsuarioID esperado %s, recebeu %s", usuarioID, log.UsuarioID())
	}

	if log.Acao() != acao.String() {
		t.Errorf("Acao esperada %s, recebeu %s", acao, log.Acao())
	}

	if log.Entidade() != entidade {
		t.Errorf("Entidade esperada %s, recebeu %s", entidade, log.Entidade())
	}

	if log.Detalhes() != detalhes {
		t.Errorf("Detalhes esperados %s, recebeu %s", detalhes, log.Detalhes())
	}

	if log.CriadoEm().IsZero() {
		t.Errorf("CriadoEm não deveria ser zero")
	}
}

// TestCarregarLogDoBD testa a função CarregarDoBD.
func TestCarregarLogDoBD(t *testing.T) {
	now := time.Now()

	db := LogDB{
		ID:        "log-1",
		UsuarioID: "user-1",
		Acao:      "CRIAR",
		Entidade:  "Categoria",
		Detalhes:  "Categoria criada",
		CriadoEm:  now,
	}

	log := CarregarDoBD(db)

	if log.ID() != db.ID {
		t.Errorf("ID esperado %s, recebeu %s", db.ID, log.ID())
	}

	if log.UsuarioID() != db.UsuarioID {
		t.Errorf("UsuarioID esperado %s, recebeu %s", db.UsuarioID, log.UsuarioID())
	}

	if log.Acao() != db.Acao {
		t.Errorf("Acao esperada %s, recebeu %s", db.Acao, log.Acao())
	}

	if log.CriadoEm() != now {
		t.Errorf("CriadoEm esperado %v, recebeu %v", now, log.CriadoEm())
	}
}

// TestLogValidar testa o método Validar do Log.
func TestLogValidar(t *testing.T) {
	tests := []struct {
		name    string
		log     Log
		wantErr bool
	}{
		{
			name: "Log válido",
			log: Log{
				id:        "log-1",
				usuarioID: "user-1",
				acao:      Criar,
				entidade:  "Entidade",
				detalhes:  "Detalhes do log",
			},
			wantErr: false,
		},
		{
			name: "ID vazio",
			log: Log{
				id:        "",
				usuarioID: "user-1",
				acao:      Criar,
				entidade:  "Entidade",
				detalhes:  "Detalhes do log",
			},
			wantErr: true,
		},
		{
			name: "UsuarioID vazio",
			log: Log{
				id:				"log-1",
				usuarioID: "",
				acao:      Criar,
				entidade:  "Entidade",
				detalhes:  "Detalhes do log",
			},
			wantErr: true,
		},
		{
			name: "Ação inválida",
			log: Log{
				id:        "log-1",
				usuarioID: "user-1",
				acao:      Acao("INVALIDA"),
				entidade:  "Entidade",
				detalhes:  "Detalhes do log",
			},
			wantErr: true,
		},
		{
			name: "Entidade vazia",
			log: Log{
				id:        "log-1",
				usuarioID: "user-1",
				acao:      Criar,
				entidade:  "",
				detalhes:  "Detalhes do log",
			},
			wantErr: true,
		},
		{
			name: "Detalhes vazios",
			log: Log{
				id:        "log-1",
				usuarioID: "user-1",
				acao:      Criar,
				entidade:  "Entidade",
				detalhes:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validar()
			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("não esperava erro, mas recebeu: %v", err)
			}
		})
	}
}