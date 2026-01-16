// Package usuario gerencia os usuários do sistema.
//
// Um Usuario representa uma pessoa que interage com o sistema e possui permissões, status e informações de login.
//
// O pacote fornece:
//
//   - A entidade principal `Usuario` com validação de campos, campos de auditoria (criadoEm, atualizadoEm, ultimoLogin) e métodos de acesso;
//   - Objeto de valor `Email` com validação de formato e serialização JSON;
//   - Enum `Permissao` para definir níveis de acesso (ADM, TEC, USR, DEV);
//   - Estruturas de filtro (`Filtro`) para listagem de usuários, com suporte a paginação, busca textual, status e permissão;
//   - Interface de repositório (`Repository`) para persistência e consulta no banco de dados;
//   - Interfaces de serviço (`Service`, `Leitor`, `Escritor`) para casos de uso.
//
// # Validações de Domínio
//
// Principais validações da entidade Usuario:
//
//   - ID do usuário não pode estar vazio;
//   - Nome deve ter entre 3 e 100 caracteres;
//   - Login não pode estar vazio;
//   - Email deve ter formato válido (user@example.com);
//   - Permissão deve ser uma das válidas (ADM, TEC, USR, DEV);
//
// Todas as validações utilizam erros estruturados do pacote `domain` (`ErrosValidacao`).
//
// # Detalhes da Implementação
//
// A entidade Usuario oferece:
//
//   - Novo():
//       Cria um novo usuário aplicando validações de domínio;
//
//   - CarregarDoBD():
//       Cria instâncias a partir de dados do banco;
//
//   - AtualizarDados():
//       Permite atualizar nome, login, email, permissão, status e avatar;
//
//   - Ativar() / Desativar():
//       Alteram o status do usuário;
//
//   - AtualizarPermissao():
//       Altera a permissão do usuário;
//
//   - AtualizarUltimoLogin():
//       Registra a data do último login do usuário;
//
//   - Métodos de acesso:
//       ID(), Nome(), Login(), Email(), Permissao(), Status(), CriadoEm(), AtualizadoEm(), UltimoLogin();
//
// # Filtro de Listagem
//
// A estrutura `Filtro` permite consultas flexíveis:
//
//   - Paginação (Pagina, Limite);
//   - Busca textual por nome ou login;
//   - Filtragem por status e permissão;
//
// Métodos auxiliares:
//
//   - Pagina(), Limite(), Offset();
//   - Normalizar();
//   - SemLimite().
//
// # Repositório e Serviços
//
// Interfaces definidas:
//
//   - Leitor:
//       BuscarPorID, BuscarPorLogin, Listar;
//
//   - Escritor:
//       Criar, Atualizar, Ativar, Desativar, AtualizarPermissao, AtualizarUltimoLogin;
//
//   - Service:
//       Agrega Leitor e Escritor.
//
// # Estrutura do Banco de Dados
//
// A estrutura UsuarioDB representa o usuário no banco de dados.
//
// # Exemplo de Uso
//
//	u, err := usuario.Novo(
//	    "u123",
//	    "João Silva",
//	    "joaosilva",
//	    usuario.NovoEmail("joao@example.com"),
//	    usuario.PermUSR,
//	    nil,
//	)
//	if err != nil {
//	    // tratar erro de validação
//	}
//
//	if err := repo.Criar(ctx, *u); err != nil {
//	    // tratar erro de persistência
//	}
package usuario
