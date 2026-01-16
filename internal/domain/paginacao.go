package domain

// Paginacao contém os campos comuns de paginação que todas as entidades usam.
type Paginacao struct {
	pagina int
	limite int
}

// NovoPaginacao cria uma nova instância de Paginacao com os dados fornecidos.
func NovoPaginacao(pagina, limite int) Paginacao {
	return Paginacao{
		pagina: pagina,
		limite: limite,
	}
}

// Normalizar aplica valores padrão para paginação.
func (p *Paginacao) Normalizar() {
	if p.pagina <= 0 {
		p.pagina = 1
	}
	if p.limite <= 0 || p.limite > 100 {
		p.limite = 10
	}
}

// Offset calcula o offset para queries SQL
func (p *Paginacao) Offset() int {
	return (p.pagina - 1) * p.limite
}

// SemLimite define paginação com limite alto para retornar todos os registros
func (p *Paginacao) SemLimite() {
	p.pagina = 1
	p.limite = 100000
}

// Pagina retorna a página atual
func (p *Paginacao) Pagina() int {
	return p.pagina
}

// Limite retorna o limite atual
func (p *Paginacao) Limite() int {
	return p.limite
}