package memory

import (
	"context"
	"sync"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/repository"
)

type RefreshRepo struct {
	mu   sync.RWMutex
	data map[string]string // token -> userID
}

func NewRefreshRepo() *RefreshRepo {
	return &RefreshRepo{data: make(map[string]string)}
}

func (r *RefreshRepo) Save(ctx context.Context, rt repository.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[rt.Token] = rt.UserID
	return nil
}

func (r *RefreshRepo) Delete(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, token)
	return nil
}

func (r *RefreshRepo) DeleteByUser(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for t, uid := range r.data {
		if uid == userID {
			delete(r.data, t)
		}
	}
	return nil
}

func (r *RefreshRepo) Exists(ctx context.Context, token string) (bool,
	error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.data[token]
	return ok, nil
}
