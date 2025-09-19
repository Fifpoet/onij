CREATE TABLE
  `music` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name` VARCHAR(255) NOT NULL COMMENT '名称',
    `full_name` VARCHAR(255) NULL COMMENT '全名称',
    `artist_ids` VARCHAR(255) NOT NULL COMMENT '艺术家id列表',
    `composer_ids` VARCHAR(255) NULL COMMENT '作曲家id列表',
    `writer_ids` VARCHAR(255) NULL COMMENT '作词家id列表',
    `issue_time` DATETIME NULL COMMENT '发行时间',
    `perform_type` INT NULL COMMENT '表演类型',
    `time_length` INT NULL COMMENT '时长(s)',
    `mv_url` VARCHAR(1024) NULL COMMENT 'mv链接',
    `audio_file_id` BIGINT NULL COMMENT '音频文件id',
    `audio_quality` INT NULL COMMENT '音频质量类型',
    `lyric_file_id` BIGINT NULL COMMENT '歌词文件id',
    `lyric_content` TEXT NULL COMMENT '歌词内容',
    `root_id` BIGINT NULL COMMENT '根id',
    `priority` INT NULL COMMENT '优先级',
    `third_id` BIGINT NULL COMMENT '三方id',
    `created_at` DATETIME NULL COMMENT '创建时间',
    `updated_at` DATETIME NULL COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '音乐';
