package repository

import (
	"context"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/model"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	Upsert(ctx context.Context, u *model.User) error
	TouchLogin(ctx context.Context, login string, t time.Time) error
}

// ErrNotFound é um erro genérico retornado quando o usuário não é encontrado
var (
	ErrNotFound = fmtError("not found")
)

// fmtError é um tipo customizado que implementa a interface error
// Permite criar erros simples com string
type fmtError string

// Error implementa a interface error
func (e fmtError) Error() string {
	return string(e)
}
