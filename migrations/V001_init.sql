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