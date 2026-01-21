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

// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// ptr é um helper que retorna um ponteiro para o valor fornecido.
func ptr[T any](v T) *T {
	return &v
}

// ptrString é uma função auxiliar para obter um ponteiro para uma string.
func ptrString(s string) *string {
	return &s
}

// testLogWriter é um writer personalizado para capturar logs durante os testes.
type testLogWriter struct {
	writeFn func(msg string)
}

// Write implementa a interface io.Writer.
func (w *testLogWriter) Write(p []byte) (n int, err error) {
	w.writeFn(string(p))
	return len(p), nil
}
