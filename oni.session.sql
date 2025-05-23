
INSERT INTO onij_tag (
    id ,
    resource_id,
    resource_type,
    tag_biz,
    tag_group,
    tag_type,
    target_id,
    target_type,
    list_show,
    extra,
    created_at,
    updated_at,
    deleted_at
  )
VALUES (
    111,
    'resource_id:bigint',
    resource_type:int,
    tag_biz:int,
    tag_group:int,
    tag_type:int,
    'target_id:bigint',
    target_type:int,
    'list_show:tinyint',
    'extra:varchar',
    'created_at:datetime',
    'updated_at:datetime',
    'deleted_at:datetime'
  );

  SELECT * FROM 