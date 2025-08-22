package repository

import "context"

// RefreshToken representa um refresh armazenado
// Pode incluir metadados (IP, UserAgent, etc)
type RefreshToken struct {
	Token  string
	UserID string
}

type RefreshRepository interface {
	Save(ctx context.Context, rt RefreshToken) error
	Delete(ctx context.Context, token string) error
	DeleteByUser(ctx context.Context, userID string) error
	Exists(ctx context.Context, token string) (bool, error)
}
