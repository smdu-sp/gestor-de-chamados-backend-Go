package domain

// GeradorID define a interface para geração de IDs únicos.
// Implementações dessa interface devem fornecer um método NovoID que retorna um novo ID como string.
// Isso permite a flexibilidade de usar diferentes estratégias de geração de IDs, como UUIDs, IDs sequenciais, etc.
type GeradorID interface {
	NovoID() (string, error)
}
