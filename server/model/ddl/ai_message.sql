CREATE TABLE
  `ai_message` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `session_id` BIGINT NOT NULL COMMENT '会话id',
    `role` VARCHAR(32) NOT NULL COMMENT 'user/assistant/tool/system',
    `content` MEDIUMTEXT NOT NULL COMMENT '内容',
    `tool_name` VARCHAR(128) NULL COMMENT 'tool名',
    `tool_call_id` VARCHAR(128) NULL COMMENT 'tool调用id',
    `token_prompt` INT NOT NULL DEFAULT 0 COMMENT 'prompt tokens',
    `token_completion` INT NOT NULL DEFAULT 0 COMMENT 'completion tokens',
    `token_total` INT NOT NULL DEFAULT 0 COMMENT 'total tokens',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_ai_msg_session` (`session_id`, `id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'AI消息';
