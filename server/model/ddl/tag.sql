CREATE TABLE
  `tag` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `resource_id` BIGINT NULL COMMENT '资源id',
    `resource_type` INT NULL COMMENT '资源类型',
    `tag_biz` INT NULL COMMENT '标签业务',
    `tag_group` INT NULL COMMENT '标签组',
    `tag_type` INT NULL COMMENT '标签类型',
    `target_id` BIGINT NULL COMMENT '目标id',
    `target_type` INT NULL COMMENT '目标类型',
    `extra` VARCHAR(256) NULL COMMENT '额外信息',
    `created_at` DATETIME NULL COMMENT '创建时间',
    `updated_at` DATETIME NULL COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_resource_biz_group_type_target` (`resource_id`, `tag_biz`, `tag_group`, `tag_type`, `target_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '标签';
