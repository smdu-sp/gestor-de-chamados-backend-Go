package migrations

import (
	"context"
	"embed"
)

//go:embed *.sql
var FS embed.FS // FS contém os arquivos de migração embutidos.

// Migrador define o contrato para execução de migrações.
type Migrador interface {
	Migrar(ctx context.Context) error
}
