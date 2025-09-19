CREATE TABLE
  `album_music` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name` VARCHAR(255) NOT NULL COMMENT '歌曲名称',
    `time_length` INT NULL COMMENT '时长(s)',
    `music_id` BIGINT NOT NULL COMMENT '音乐id',
    `album_id` BIGINT NOT NULL COMMENT '专辑id',
    `artist_names` VARCHAR(1024) NULL COMMENT '艺术家名称',
    `third_id` BIGINT NULL COMMENT '三方id',
    `is_available` TINYINT NOT NULL COMMENT '是否可用',
    `created_at` DATETIME NULL COMMENT '创建时间',
    `updated_at` DATETIME NULL COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '专辑音乐';
