CREATE TABLE
  `practice` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `practice_type` INT NOT NULL COMMENT '练习类型',
    `content` TEXT NOT NULL COMMENT '练习内容',
    `practice_at` BIGINT NOT NULL COMMENT '练习时间戳',
    `duration` INT NOT NULL DEFAULT 0 COMMENT '时长(分钟)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_practice_at` (`practice_at`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '练习';
