package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockRepository é um stub manual para RepositoryInterface
type mockRepository struct {
	findByIDFunc        func(ctx context.Context, id string) (*Usuario, error)
	findByLoginFunc     func(ctx context.Context, login string) (*Usuario, error)
	insertFunc          func(ctx context.Context, u *Usuario) error
	updateFunc          func(ctx context.Context, id string, u *Usuario) error
	deleteFunc          func(ctx context.Context, id string) error
	updateLastLoginFunc func(ctx context.Context, id string) error
	listFunc            func(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]Usuario, int, error)
}

func (m *mockRepository) FindByID(ctx context.Context, id string) (*Usuario, error) {
	return m.findByIDFunc(ctx, id)
}
func (m *mockRepository) FindByLogin(ctx context.Context, login string) (*Usuario, error) {
	return m.findByLoginFunc(ctx, login)
}
func (m *mockRepository) Insert(ctx context.Context, u *Usuario) error {
	return m.insertFunc(ctx, u)
}
func (m *mockRepository) Update(ctx context.Context, id string, u *Usuario) error {
	return m.updateFunc(ctx, id, u)
}
func (m *mockRepository) Delete(ctx context.Context, id string) error {
	return m.deleteFunc(ctx, id)
}
func (m *mockRepository) UpdateLastLogin(ctx context.Context, id string) error {
	return m.updateLastLoginFunc(ctx, id)
}
func (m *mockRepository) List(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]Usuario, int, error) {
	return m.listFunc(ctx, pagina, limite, busca, status, permissao)
}

// --- TESTES ---

// TestBuscarPorIDSucesso verifica se a busca por ID retorna o usuário correto.
func TestBuscarPorIDSucesso(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*Usuario, error) {
			return &Usuario{ID: id, Nome: "João"}, nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	u, err := svc.BuscarPorID(ctx, "123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "123", u.ID)
	assert.Equal(t, "João", u.Nome)
}

// TestBuscarPorIDNotFound verifica se a busca por ID retorna um erro quando o usuário não é encontrado.
func TestBuscarPorIDNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*Usuario, error) {
			return nil, nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	u, err := svc.BuscarPorID(ctx, "123")

	// Assert
	assert.ErrorIs(t, err, ErrUsuarioNaoEncontrado)
	assert.Nil(t, u)
}

// TestCriarCamposObrigatorios verifica se a criação de um usuário sem campos obrigatórios resulta em um erro.
func TestCriarCamposObrigatorios(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{}
	svc := NewService(mockRepo)

	// Act
	err := svc.Criar(ctx, &Usuario{})

	// Assert
	assert.ErrorIs(t, err, ErrCamposObrigatorios)
}

// TestCriarSucesso verifica se a criação de um usuário com dados válidos resulta em sucesso.
func TestCriarSucesso(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var saved *Usuario
	mockRepo := &mockRepository{
		insertFunc: func(ctx context.Context, u *Usuario) error {
			saved = u
			return nil
		},
	}
	svc := NewService(mockRepo)

	u := &Usuario{Nome: "Maria", Login: "maria", Email: "maria@test.com"}

	// Act
	err := svc.Criar(ctx, u)

	// Assert
	assert.NoError(t, err)
	assert.True(t, saved.Status)
	assert.Equal(t, PermUSR, saved.Permissao) // default
}

// TestAtualizarSucesso verifica se a atualização de um usuário com dados válidos resulta em sucesso.
func TestAtualizarSucesso(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		updateFunc: func(ctx context.Context, id string, u *Usuario) error {
			return nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	err := svc.Atualizar(ctx, "123", &Usuario{Nome: "Novo"})

	// Assert
	assert.NoError(t, err)
}

// TestRemoverSucesso verifica se a remoção de um usuário com ID válido resulta em sucesso.
func TestRemoverSucesso(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		deleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	err := svc.Remover(ctx, "123")

	// Assert
	assert.NoError(t, err)
}

// TestAtualizarUltimoLoginSucesso verifica se a atualização do último login de um usuário com ID válido resulta em sucesso.
func TestAtualizarUltimoLoginSucesso(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		updateLastLoginFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	err := svc.AtualizarUltimoLogin(ctx, "123")

	// Assert
	assert.NoError(t, err)
}

// TestPermitidoSucesso verifica se a permissão é concedida corretamente.
func TestPermitidoSucesso(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*Usuario, error) {
			return &Usuario{ID: id, Permissao: PermADM}, nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	ok, err := svc.Permitido(ctx, "123", []string{"USR", "ADM"})

	// Assert
	assert.NoError(t, err)
	assert.True(t, ok)
}

// TestPermitidoNaoPermitido verifica se a permissão é negada corretamente.
func TestPermitidoNaoPermitido(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*Usuario, error) {
			return &Usuario{ID: id, Permissao: PermUSR}, nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	ok, err := svc.Permitido(ctx, "123", []string{"ADM"})

	// Assert
	assert.NoError(t, err)
	assert.False(t, ok)
}

// TestPermitidoNotFound verifica se a permissão retorna um erro quando o usuário não é encontrado.
func TestPermitidoNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*Usuario, error) {
			return nil, nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	ok, err := svc.Permitido(ctx, "123", []string{"ADM"})

	// Assert
	assert.ErrorIs(t, err, ErrUsuarioNaoEncontrado)
	assert.False(t, ok)
}

// TestListarSucesso verifica se a listagem de usuários retorna os usuários corretos.
func TestListarSucesso(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mockRepository{
		listFunc: func(ctx context.Context, pagina, limite int, busca, status, permissao *string) ([]Usuario, int, error) {
			return []Usuario{
				{ID: "1", Nome: "João", CriadoEm: time.Now()},
				{ID: "2", Nome: "Maria", CriadoEm: time.Now()},
			}, 2, nil
		},
	}
	svc := NewService(mockRepo)

	// Act
	users, total, err := svc.Listar(ctx, 1, 10, nil, nil, nil)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, 2, total)
}
