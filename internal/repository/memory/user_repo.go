package memory

import (
	"context"
	"sync"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/util"
)

// UserRepo é um repositório em memória para armazenar usuários.
// Usa mutex para concorrência segura e mapas para acesso rápido por ID e login.
type UserRepo struct {
	mu   sync.RWMutex           // Mutex para leitura/escrita segura em múltiplas goroutines
	byID map[string]*model.User // Mapeia usuário por ID
	byLG map[string]*model.User // Mapeia usuário por login
}

// Cria uma instância nova de UserRepo
func NewUserRepo() *UserRepo {
	return &UserRepo{
		byID: make(map[string]*model.User),
		byLG: make(map[string]*model.User),
	}
}

// Retorna uma cópia do usuário para evitar alterações externas diretas
func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	r.mu.RLock()         // Bloqueio apenas para leitura
	defer r.mu.RUnlock() // Garante desbloqueio ao final

	if u, ok := r.byLG[login]; ok {
		copy := *u       // cria cópia para segurança
		return &copy, nil
	}

	return nil, repository.ErrNotFound
}

// Upsert cria ou atualiza um usuário no repositório
func (r *UserRepo) Upsert(ctx context.Context, u *model.User) error {
	r.mu.Lock()          // Bloqueio para escrita
	defer r.mu.Unlock()  // Garante desbloqueio ao final
	now := time.Now()

	if u.ID == "" {      // Se não tiver ID, cria um novo
		u.ID = util.NewID()
		u.CriadoEm = now
	}

	u.AtualizadoEm = now
	r.byID[u.ID] = u     // Atualiza ou adiciona no mapa por ID
	r.byLG[u.Login] = r.byID[u.ID] // Atualiza ou adiciona no mapa por login
	return nil
}

// Atualiza o último login do usuário e timestamp de atualização
func (r *UserRepo) TouchLogin(ctx context.Context, login string, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if u, ok := r.byLG[login]; ok {
		novoLogin := t
		u.UltimoLogin = &novoLogin
		u.AtualizadoEm = t
		return nil
	}

	return repository.ErrNotFound
}
