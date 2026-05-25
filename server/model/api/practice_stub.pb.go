// 临时类型定义：运行 `makefile` 中的 hz update 后请删除本文件。
package api

type PracticeType int32

const (
	PracticeType_PCT_Unknown            PracticeType = 0
	PracticeType_PCT_AnaerobicExercise  PracticeType = 1
	PracticeType_PCT_AerobicExercise    PracticeType = 2
	PracticeType_PCT_Reading            PracticeType = 3
	PracticeType_PCT_Language           PracticeType = 4
	PracticeType_PCT_Game               PracticeType = 5
	PracticeType_PCT_Piano              PracticeType = 6
	PracticeType_PCT_Flute              PracticeType = 7
	PracticeType_PCT_Singing            PracticeType = 8
	PracticeType_PCT_Other              PracticeType = 99
)

type Practice struct {
	Id           int64        `json:"id,omitempty" form:"id" query:"id"`
	PracticeType PracticeType `json:"practice_type,omitempty" form:"practice_type" query:"practice_type"`
	Content      string       `json:"content,omitempty" form:"content" query:"content"`
	PracticeAt   int64        `json:"practice_at,omitempty" form:"practice_at" query:"practice_at"`
	Duration     int32        `json:"duration,omitempty" form:"duration" query:"duration"`
	CreatedAt    int64        `json:"created_at,omitempty" form:"created_at" query:"created_at"`
	UpdatedAt    int64        `json:"updated_at,omitempty" form:"updated_at" query:"updated_at"`
}

type UploadPracticeReq struct {
	PracticeType PracticeType `json:"practice_type,omitempty" form:"practice_type" query:"practice_type"`
	Content      string       `json:"content,omitempty" form:"content" query:"content"`
	PracticeAt   int64        `json:"practice_at,omitempty" form:"practice_at" query:"practice_at"`
	Duration     int32        `json:"duration,omitempty" form:"duration" query:"duration"`
	Id           *int64       `json:"id,omitempty" form:"id" query:"id"`
}

type UploadPracticeResp struct {
	Code       int32  `json:"code,omitempty" form:"code" query:"code"`
	Message    string `json:"message,omitempty" form:"message" query:"message"`
	PracticeId int64  `json:"practice_id,omitempty" form:"practice_id" query:"practice_id"`
}

type DeletePracticeReq struct {
	Id int64 `json:"id,omitempty" form:"id" query:"id"`
}

type DeletePracticeResp struct {
	Code    int32  `json:"code,omitempty" form:"code" query:"code"`
	Message string `json:"message,omitempty" form:"message" query:"message"`
}

type GetMonthlyPracticeReq struct {
	Year  int32 `json:"year,omitempty" form:"year" query:"year"`
	Month int32 `json:"month,omitempty" form:"month" query:"month"`
}

type PracticeDay struct {
	Day       int32       `json:"day,omitempty" form:"day" query:"day"`
	Practices []*Practice `json:"practices,omitempty" form:"practices" query:"practices"`
}

type GetMonthlyPracticeResp struct {
	Code    int32          `json:"code,omitempty" form:"code" query:"code"`
	Message string         `json:"message,omitempty" form:"message" query:"message"`
	Days    []*PracticeDay `json:"days,omitempty" form:"days" query:"days"`
}
