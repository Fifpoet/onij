package prm

import "onij/infra/mysql"

type GetWeeklyTodParam struct {
}

type TodPrime struct {
	Tod  mysql.Tod
	Tags []mysql.Tag
}

type GetWeeklyTodResult struct {
	Tods [][]TodPrime `json:"tods"`
}
