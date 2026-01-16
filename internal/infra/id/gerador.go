package id

import dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"

// Gerador é um gerador de IDs.
type Gerador struct{}

// Asserção de interface para garantir que Gerador implementa dmn.GeradorID.
var _ dmn.GeradorID = (*Gerador)(nil)

// NovoID retorna um novo identificador único.
// Neste caso, um UUIDv7 em formato string, porém a implementação pode ser alterada
// conforme necessário, sem impactar o restante do sistema, mudando apenas o retorno desta função.
func (Gerador) NovoID() (string, error) {
	return GerarUUIDv7String()
}