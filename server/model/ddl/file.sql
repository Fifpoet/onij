CREATE TABLE
  `file` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name` VARCHAR(255) NOT NULL COMMENT '文件名称',
    `format` INT NOT NULL COMMENT '文件格式',
    `size` BIGINT NOT NULL COMMENT '文件大小',
    `store_key` VARCHAR(255) NOT NULL COMMENT '文件存储key',
    `hash` VARCHAR(255) NOT NULL COMMENT '文件hash',
    `origin_at` BIGINT NOT NULL COMMENT '原始时间',
    `created_at` DATETIME NULL COMMENT '创建时间',
    `updated_at` DATETIME NULL COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文件';