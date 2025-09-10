CREATE TABLE
  `artist` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name` VARCHAR(255) NOT NULL COMMENT '艺术家名字',
    `avatar_file_id` BIGINT NOT NULL COMMENT '艺术家头像',
    `created_at` DATETIME NULL COMMENT '创建时间',
    `updated_at` DATETIME NULL COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '艺术家';