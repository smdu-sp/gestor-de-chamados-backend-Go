// Package banner fornece a funcionalidade responsável por exibir o banner
// inicial da aplicação no momento de sua inicialização.
//
// Este pacote encapsula a lógica de carregamento e formatação do banner
// apresentado durante o boot da aplicação.
//
// O banner é carregado a partir do arquivo embedado `banner.txt`, utilizando
// a diretiva `go:embed`. O conteúdo inclui o logo da aplicação em ASCII art
// e um marcador {{VERSAO}}, que é substituído dinamicamente pela versão
// informada em runtime.
//
// # Função exposta
//
//   - CarregarBanner(version string) string
//       Substitui o marcador {{VERSAO}} pela versão informada e retorna
//       o banner formatado com códigos ANSI, pronto para exibição no terminal.
//
// # Uso típico
//
// O pacote é utilizado durante o processo de inicialização da aplicação,
// normalmente a partir do package main, como parte do fluxo de boot.
//
// Observações:
//   - Este pacote não possui dependências externas.
//   - Não contém regras de negócio.
//   - Sua responsabilidade é exclusivamente visual e informativa.
package banner
