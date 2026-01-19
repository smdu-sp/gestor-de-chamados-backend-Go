package domain

import "testing"

// TestNovoErrosValidacao_InicialmenteVazio verifica se uma nova instância de ErrosValidacao
func TestNovoErrosValidacao_InicialmenteVazio(t *testing.T) {
	v := NovoErrosValidacao()

	if v.HaErros() {
		t.Fatal("nao deveria haver erros inicialmente")
	}

	if v.Error() != "" {
		t.Fatalf("mensagem de erro deveria ser vazia, recebeu: %q", v.Error())
	}
}

// TestErrosValidacao_Add verifica se adicionar erros funciona corretamente.
func TestErrosValidacao_Add(t *testing.T) {
	v := NovoErrosValidacao()

	v.Add("Nome", "nome invalido")

	if !v.HaErros() {
		t.Fatal("deveria indicar que ha erros")
	}

	if len(v.Erros()) != 1 {
		t.Fatalf("esperava 1 erro, recebeu %d", len(v.Erros()))
	}
}

// TestErrosValidacao_ErrorOrdenado verifica se a mensagem de erro é gerada em ordem alfabética dos campos.
func TestErrosValidacao_ErrorOrdenado(t *testing.T) {
	v := NovoErrosValidacao()

	v.Add("Email", "email invalido")
	v.Add("Nome", "nome invalido")

	msg := v.Error()

	expected := "Email: email invalido; Nome: nome invalido"
	if msg != expected {
		t.Fatalf("mensagem inesperada\nesperado: %q\nrecebido: %q", expected, msg)
	}
}

// TestErrosValidacao_ErrosRetornaMapa verifica se o método Erros retorna o mapa correto.
func TestErrosValidacao_ErrosRetornaMapa(t *testing.T) {
	v := NovoErrosValidacao()

	v.Add("Login", "login obrigatorio")

	erros := v.Erros()

	if erros["Login"] != "login obrigatorio" {
		t.Fatalf("mensagem inesperada: %v", erros)
	}
}

// TestErrosValidacao_MarshalJSON verifica se a serialização JSON funciona corretamente.
func TestErrosValidacao_MarshalJSON_Vazio(t *testing.T) {
	v := NovoErrosValidacao()

	data, err := v.MarshalJSON()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if string(data) != "{}" {
		t.Fatalf("esperava {}, recebeu %s", string(data))
	}
}

// TestErrosValidacao_MarshalJSON_ComErros verifica se a serialização JSON com erros funciona corretamente.
func TestErrosValidacao_MarshalJSON_ComErros(t *testing.T) {
	v := NovoErrosValidacao()

	v.Add("Nome", "nome invalido")
	v.Add("Email", "email invalido")

	data, err := v.MarshalJSON()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	expected := `[{"campo":"Email","erro":"email invalido"},{"campo":"Nome","erro":"nome invalido"}]`

	if string(data) != expected {
		t.Fatalf("json inesperado\nesperado: %s\nrecebido: %s", expected, string(data))
	}
}
