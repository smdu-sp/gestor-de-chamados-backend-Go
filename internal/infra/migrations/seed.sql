-- ======================================================================
-- SEED IDEMPOTENTE - DESENVOLVIMENTO / TESTES
-- ======================================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ======================================================================
-- USUÁRIOS
-- ======================================================================
INSERT INTO usuarios (id, nome, login, email, permissao, status)
VALUES
('11111111-1111-1111-1111-111111111111', 'Administrador', 'admin', 'admin@sistema.local', 'ADM', TRUE),
('22222222-2222-2222-2222-222222222222', 'Técnico João', 'joao.tec', 'joao@sistema.local', 'TEC', TRUE),
('33333333-3333-3333-3333-333333333333', 'Usuário Maria', 'maria.usr', 'maria@sistema.local', 'USR', TRUE),
('44444444-4444-4444-4444-444444444444', 'Dev Carlos', 'carlos.dev', 'carlos@sistema.local', 'DEV', TRUE)
ON DUPLICATE KEY UPDATE
  login = login;

-- ======================================================================
-- CATEGORIAS
-- ======================================================================
INSERT INTO categorias (id, nome, status)
VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Infraestrutura', TRUE),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Software', TRUE)
ON DUPLICATE KEY UPDATE
  nome = nome;

-- ======================================================================
-- SUBCATEGORIAS
-- ======================================================================
INSERT INTO subcategorias (id, nome, status, categoria_id)
VALUES
('aaaaaaaa-1111-aaaa-aaaa-aaaaaaaaaaaa', 'Rede', TRUE, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
('aaaaaaaa-2222-aaaa-aaaa-aaaaaaaaaaaa', 'Servidores', TRUE, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
('bbbbbbbb-1111-bbbb-bbbb-bbbbbbbbbbbb', 'Sistema Interno', TRUE, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'),
('bbbbbbbb-2222-bbbb-bbbb-bbbbbbbbbbbb', 'Aplicação Web', TRUE, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb')
ON DUPLICATE KEY UPDATE
  nome = nome;

-- ======================================================================
-- CHAMADOS
-- ======================================================================
INSERT INTO chamados (
  id, titulo, descricao, status,
  categoria_id, subcategoria_id, criador_id
)
VALUES
(
  'cccccccc-1111-cccc-cccc-cccccccccccc',
  'Internet instável',
  'A conexão cai diversas vezes durante o dia.',
  'ABERTO',
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'aaaaaaaa-1111-aaaa-aaaa-aaaaaaaaaaaa',
  '33333333-3333-3333-3333-333333333333'
),
(
  'cccccccc-2222-cccc-cccc-cccccccccccc',
  'Erro no sistema interno',
  'Sistema não permite salvar registros.',
  'ATRIBUIDO',
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'bbbbbbbb-1111-bbbb-bbbb-bbbbbbbbbbbb',
  '33333333-3333-3333-3333-333333333333'
)
ON DUPLICATE KEY UPDATE
  id = id;

-- ======================================================================
-- ATENDIMENTOS
-- ======================================================================
INSERT INTO atendimentos (id, atribuido_id, chamado_id)
VALUES
(
  'dddddddd-1111-dddd-dddd-dddddddddddd',
  '22222222-2222-2222-2222-222222222222',
  'cccccccc-2222-cccc-cccc-cccccccccccc'
)
ON DUPLICATE KEY UPDATE
  id = id;

-- ======================================================================
-- ACOMPANHAMENTOS
-- ======================================================================
INSERT INTO acompanhamentos (id, conteudo, chamado_id, usuario_id, remetente)
VALUES
(
  'eeeeeeee-1111-eeee-eeee-eeeeeeeeeeee',
  'Chamado aberto, aguardando análise.',
  'cccccccc-1111-cccc-cccc-cccccccccccc',
  '33333333-3333-3333-3333-333333333333',
  'USR'
),
(
  'eeeeeeee-2222-eeee-eeee-eeeeeeeeeeee',
  'Estou analisando o problema.',
  'cccccccc-2222-cccc-cccc-cccccccccccc',
  '22222222-2222-2222-2222-222222222222',
  'TEC'
)
ON DUPLICATE KEY UPDATE
  id = id;

-- ======================================================================
-- PERMISSÕES POR CATEGORIA
-- ======================================================================
INSERT INTO categoria_permissoes (categoria_id, permissao, usuario_id)
VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'TEC', '22222222-2222-2222-2222-222222222222'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'TEC', '22222222-2222-2222-2222-222222222222'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'ADM', '11111111-1111-1111-1111-111111111111')
ON DUPLICATE KEY UPDATE
  permissao = permissao;

-- ======================================================================
-- LOGS
-- ======================================================================
-- Logs são eventos históricos.
-- Normalmente NÃO entram em seed.
-- Mantidos aqui apenas como exemplo técnico.

INSERT INTO logs (id, usuario_id, acao, entidade, detalhes)
VALUES
(
  'ffffffff-1111-ffff-ffff-ffffffffffff',
  '33333333-3333-3333-3333-333333333333',
  'CRIAR',
  'chamados',
  'Chamado criado pelo usuário'
),
(
  'ffffffff-2222-ffff-ffff-ffffffffffff',
  '22222222-2222-2222-2222-222222222222',
  'ATUALIZAR',
  'chamados',
  'Chamado atribuído ao técnico'
)
ON DUPLICATE KEY UPDATE
  id = id;

SET FOREIGN_KEY_CHECKS = 1;
