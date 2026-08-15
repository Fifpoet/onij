CREATE TABLE
  `ai_session` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `title` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '会话标题',
    `kind` TINYINT NOT NULL DEFAULT 0 COMMENT '0普通 1语音助手',
    `token_total` BIGINT NOT NULL DEFAULT 0 COMMENT '累计token',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_ai_session_updated` (`updated_at`),
    KEY `idx_ai_session_kind` (`kind`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'AI会话';
