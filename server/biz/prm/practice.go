package prm

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
	"sort"
	"time"
)

type UploadPracticeParam struct {
	Id           *int64
	PracticeType int32
	Content      string
	PracticeAt   int64
	Duration     int32
}

func NewUploadPracticeParam(req *api.UploadPracticeReq) *UploadPracticeParam {
	return &UploadPracticeParam{
		Id:           req.Id,
		PracticeType: int32(req.PracticeType),
		Content:      req.Content,
		PracticeAt:   req.PracticeAt,
		Duration:     req.Duration,
	}
}

type UploadPracticeResult struct {
	PracticeId int64
}

func (r *UploadPracticeResult) Resp() *api.UploadPracticeResp {
	return &api.UploadPracticeResp{
		Code:       util.BaseCodeOK,
		Message:    util.BaseMsgOK,
		PracticeId: r.PracticeId,
	}
}

type DeletePracticeParam struct {
	Id int64
}

func NewDeletePracticeParam(req *api.DeletePracticeReq) *DeletePracticeParam {
	return &DeletePracticeParam{Id: req.Id}
}

type DeletePracticeResult struct{}

func (r *DeletePracticeResult) Resp() *api.DeletePracticeResp {
	return &api.DeletePracticeResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
	}
}

type GetMonthlyPracticeParam struct {
	Year  int32
	Month int32
}

func NewGetMonthlyPracticeParam(req *api.GetMonthlyPracticeReq) *GetMonthlyPracticeParam {
	return &GetMonthlyPracticeParam{
		Year:  req.Year,
		Month: req.Month,
	}
}

type GetMonthlyPracticeResult struct {
	Days []*api.PracticeDay
}

func toAPIPractice(p *model.Practice) *api.Practice {
	if p == nil {
		return nil
	}
	return &api.Practice{
		Id:           p.Id,
		PracticeType: api.PracticeType(p.PracticeType),
		Content:      p.Content,
		PracticeAt:   p.PracticeAt,
		Duration:     p.Duration,
		CreatedAt:    p.CreatedAt.Unix(),
		UpdatedAt:    p.UpdatedAt.Unix(),
	}
}

func (r *GetMonthlyPracticeResult) Resp() *api.GetMonthlyPracticeResp {
	days := r.Days
	if days == nil {
		days = []*api.PracticeDay{}
	}
	return &api.GetMonthlyPracticeResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Days:    days,
	}
}

func BuildMonthlyPracticeDays(practices []*model.Practice, loc *time.Location) []*api.PracticeDay {
	if loc == nil {
		loc = time.Local
	}
	dayMap := make(map[int32][]*api.Practice)
	for _, p := range practices {
		if p == nil {
			continue
		}
		day := int32(time.Unix(p.PracticeAt, 0).In(loc).Day())
		dayMap[day] = append(dayMap[day], toAPIPractice(p))
	}

	days := make([]*api.PracticeDay, 0, len(dayMap))
	for day, list := range dayMap {
		sort.Slice(list, func(i, j int) bool {
			return list[i].PracticeAt < list[j].PracticeAt
		})
		days = append(days, &api.PracticeDay{
			Day:       day,
			Practices: list,
		})
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Day < days[j].Day
	})
	if len(days) == 0 {
		return []*api.PracticeDay{}
	}
	return days
}

func MonthRange(year, month int32, loc *time.Location) (start, end int64) {
	if loc == nil {
		loc = time.Local
	}
	startTime := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, loc)
	endTime := startTime.AddDate(0, 1, 0)
	return startTime.Unix(), endTime.Unix()
}
