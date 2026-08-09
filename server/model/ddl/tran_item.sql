CREATE TABLE
  `tran_item` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `kind` TINYINT NOT NULL COMMENT '1=text 2=file',
    `content` TEXT NULL COMMENT '文本内容',
    `file_id` BIGINT NULL COMMENT '关联 file.id',
    `name` VARCHAR(255) NULL COMMENT '文件名',
    `size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小',
    `format` INT NULL COMMENT '文件格式',
    `pinned` TINYINT NOT NULL DEFAULT 0 COMMENT '是否置顶',
    `pinned_at` BIGINT NULL COMMENT '置顶时间戳(ms)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tran_pinned_created` (`pinned`, `pinned_at`, `created_at`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '中转站条目';
