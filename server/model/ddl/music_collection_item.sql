CREATE TABLE
  `music_collection_item` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    `collection_id` BIGINT NOT NULL COMMENT '合集id',
    `song_id` BIGINT NOT NULL COMMENT '歌曲三方id(网易云)',
    `song_name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '歌名(搜索冗余)',
    `artist_names` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '歌手名(搜索冗余)',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_coll_song` (`collection_id`, `song_id`),
    KEY `idx_mci_coll` (`collection_id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '音乐合集曲目';
