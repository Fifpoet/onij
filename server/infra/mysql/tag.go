package mysql

// Tag设计:
// 1. 需要联动 tod 图片 card poem, 一级业务分类

type Tag struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Biz        int    `json:"biz"`
	Group      int    `json:"group"`
	TagType    int    `json:"tag_type"`
	Target     int    `json:"target"`
	TargetType int    `json:"target_type"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
}
