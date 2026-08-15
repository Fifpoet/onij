CREATE TABLE
  `ai_confirm_pending` (
    `id` BIGINT NOT NULL COMMENT '主键 confirm_id',
    `session_id` BIGINT NOT NULL COMMENT '会话id',
    `tool_call_id` VARCHAR(128) NOT NULL COMMENT '对应 assistant.tool_calls.id',
    `tool_name` VARCHAR(128) NOT NULL COMMENT 'tool名',
    `draft_args` MEDIUMTEXT NOT NULL COMMENT '待执行参数JSON',
    `ui_json` MEDIUMTEXT NOT NULL COMMENT '确认UI JSON',
    `status` VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT 'pending/confirmed/cancelled/expired',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_ai_confirm_session` (`session_id`, `status`),
    KEY `idx_ai_confirm_tool_call` (`tool_call_id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'AI HITL待确认';
