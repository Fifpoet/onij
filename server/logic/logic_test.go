package logic

import (
	"fmt"
	"onij/infra"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化逻辑
	fmt.Println("Running init logic...")

	infra.NewAllInfra()
}

func TestGetWeeklyTodList(t *testing.T) {

}
