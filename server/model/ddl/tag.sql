create table onij_tag
(
    id            bigint auto_increment
        primary key,
    resource_id   bigint       null,
    resource_type int          null,
    tag_biz       int          null,
    tag_group     int          null,
    tag_type      int          null,
    target_id     bigint       null,
    target_type   int          null,
    list_show     tinyint(1)   null,
    extra         varchar(256) null,
    created_at    datetime     null,
    updated_at    datetime     null,
    deleted_at    datetime     null,
    constraint uk_resource_biz_group_type_target
        unique (resource_id, tag_biz, tag_group, tag_type, target_id)
);