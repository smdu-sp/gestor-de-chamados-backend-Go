package service

// =====================================================================================================================
//  MOCKS E FAKES
// =====================================================================================================================

// --- Fake GeradorID --------------------------------------------------------------------------------------------------

// fakeGeradorID é uma implementação falsa de GeradorID para testes.
type fakeGeradorID struct {
	id  string
	err error
}

// NovoID retorna um ID falso ou um erro, dependendo da configuração do fakeGeradorID.
func (f *fakeGeradorID) NovoID() (string, error) {
	return f.id, f.err
}