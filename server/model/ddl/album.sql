CREATE TABLE
  `album` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name` VARCHAR(255) NOT NULL COMMENT '名称',
    `profile` VARCHAR(1024) NOT NULL COMMENT '简介',
    `album_type` INT NULL COMMENT '专辑类型',
    `artist_ids` VARCHAR(255) NOT NULL COMMENT '艺术家id列表',
    `issue_time` DATETIME NULL COMMENT '发行时间',
    `cover_file_id` BIGINT NULL COMMENT '封面文件id',
    `created_at` DATETIME NULL COMMENT '创建时间',
    `updated_at` DATETIME NULL COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '专辑';
