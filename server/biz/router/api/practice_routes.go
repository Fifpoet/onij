package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	handler "onij/handler"
)

// RegisterPractice 练习相关路由（hz 生成后可合并进 all.go，届时可删除本文件）
func RegisterPractice(r *server.Hertz) {
	root := r.Group("/")
	{
		_practice := root.Group("/practice")
		_practice.POST("/upload", handler.UploadPractice)
		_practice.POST("/delete", handler.DeletePractice)
		_practice.POST("/monthly", handler.GetMonthlyPractice)
	}
}
