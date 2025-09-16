-- Usuários
CREATE TABLE IF NOT EXISTS usuarios (
  id            CHAR(36)      NOT NULL PRIMARY KEY, -- UUID
  nome          VARCHAR(255)  NOT NULL,
  login         VARCHAR(255)  NOT NULL UNIQUE,
  email         VARCHAR(255)  NOT NULL UNIQUE,
  permissao     ENUM('ADM','TEC','SUP','INF','VOIP','IMP','CAD','USR','DEV') NOT NULL DEFAULT 'USR',
  status        BOOLEAN       NOT NULL DEFAULT TRUE, -- ativo/inativo
  avatar        TEXT NULL,
  ultimo_login  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  criado_em     DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_usuarios_login (login),
  INDEX idx_usuarios_email (email),
  INDEX idx_usuarios_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Categorias
CREATE TABLE IF NOT EXISTS categorias (
  id            CHAR(36) NOT NULL PRIMARY KEY,
  nome          VARCHAR(255) NOT NULL UNIQUE,
  status        BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Subcategorias
CREATE TABLE IF NOT EXISTS subcategorias (
  id           CHAR(36) NOT NULL PRIMARY KEY,
  nome         VARCHAR(255) NOT NULL,
  categoria_id CHAR(36) NOT NULL,
  criado_em    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (categoria_id) REFERENCES categorias(id) ON DELETE CASCADE ON UPDATE CASCADE,
  INDEX idx_subcategorias_categoria_id (categoria_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Permissões específicas para categorias
CREATE TABLE IF NOT EXISTS categoria_permissoes (
  categoria_id CHAR(36) NOT NULL,
  permissao    ENUM('ADM','TEC','SUP','INF','VOIP','IMP','CAD','USR','DEV') NOT NULL,
  criado_em    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(categoria_id, permissao),
  FOREIGN KEY (categoria_id) REFERENCES categorias(id) ON DELETE CASCADE ON UPDATE CASCADE,
  INDEX idx_categoria_permissoes_categoria_id (categoria_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Chamados
CREATE TABLE IF NOT EXISTS chamados (
  id               CHAR(36)     NOT NULL PRIMARY KEY,
  titulo           VARCHAR(255) NOT NULL,
  descricao        TEXT         NOT NULL,
  status           ENUM('NOVO','ATRIBUIDO','RESOLVIDO','REJEITADO','FECHADO','ARQUIVADO') NOT NULL DEFAULT 'NOVO',
  criado_em        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  solucionado_em   DATETIME NULL,
  fechado_em       DATETIME NULL,
  categoria_id     CHAR(36) NOT NULL,
  subcategoria_id  CHAR(36) NOT NULL,
  criador_id       CHAR(36) NOT NULL,
  FOREIGN KEY (categoria_id) REFERENCES categorias(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (subcategoria_id) REFERENCES subcategorias(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (criador_id) REFERENCES usuarios(id) ON UPDATE CASCADE,
  INDEX idx_chamados_categoria_id (categoria_id),
  INDEX idx_chamados_subcategoria_id (subcategoria_id),
  INDEX idx_chamados_criador_id (criador_id),
  INDEX idx_chamados_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Acompanhamentos (chat entre técnico e usuário nos chamados)
CREATE TABLE IF NOT EXISTS acompanhamentos (
  id            CHAR(36) NOT NULL PRIMARY KEY,
  conteudo      TEXT     NOT NULL,
  criado_em     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  chamado_id    CHAR(36) NOT NULL,
  usuario_id    CHAR(36) NOT NULL,
  FOREIGN KEY (chamado_id) REFERENCES chamados(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON UPDATE CASCADE,
  INDEX idx_acompanhamentos_chamado_id (chamado_id),
  INDEX idx_acompanhamentos_usuario_id (usuario_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Relacionamentos auxiliares

-- Chamados abertos
CREATE TABLE IF NOT EXISTS chamados_abertos (
  usuario_id CHAR(36) NOT NULL,
  chamado_id CHAR(36) NOT NULL,
  criado_em  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(usuario_id, chamado_id),
  FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON UPDATE CASCADE,
  FOREIGN KEY (chamado_id) REFERENCES chamados(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Chamados atendidos
CREATE TABLE IF NOT EXISTS chamados_atendidos (
  usuario_id CHAR(36) NOT NULL,
  chamado_id CHAR(36) NOT NULL,
  criado_em  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(usuario_id, chamado_id),
  FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON UPDATE CASCADE,
  FOREIGN KEY (chamado_id) REFERENCES chamados(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Relacionamento muitos-para-muitos entre chamados e técnicos
CREATE TABLE IF NOT EXISTS chamado_tecnicos (
  chamado_id CHAR(36) NOT NULL,
  tecnico_id CHAR(36) NOT NULL,
  criado_em  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(chamado_id, tecnico_id),
  FOREIGN KEY (chamado_id) REFERENCES chamados(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (tecnico_id) REFERENCES usuarios(id) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

