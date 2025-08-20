package memory

import (
	"context"
	"sync"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/util"
)

type UserRepo struct {
	mu   sync.RWMutex
	byID map[string]*model.User
	byLG map[string]*model.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		byID: make(map[string]*model.User),
		byLG: make(map[string]*model.User),
	}
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, ok := r.byLG[login]; ok {
		copy := *u
		return &copy, nil
	}

	return nil, repository.ErrNotFound
}

func (r *UserRepo) Upsert(ctx context.Context, u *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()

	if u.ID == "" {
		u.ID = util.NewID()
		u.CriadoEm = now
	}

	u.AtualizadoEm = now
	r.byID[u.ID] = u
	r.byLG[u.Login] = r.byID[u.ID]
	return nil
}

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
