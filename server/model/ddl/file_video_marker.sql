CREATE TABLE
  `file_video_marker` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `file_id` BIGINT NOT NULL COMMENT '文件id',
    `time_ms` BIGINT NOT NULL COMMENT '时间点(毫秒)',
    `label` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标记名称',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_file_time` (`file_id`, `time_ms`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '视频进度打点';
