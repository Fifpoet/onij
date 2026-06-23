package handler

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"onij/util"
)

// PreviewFile 预览/下载：proxy=1 时经七牛 SDK 回源（本地 dev 无 DNS）；否则 302 到 cloud.onij.fun 私有链
// @router /file/preview [GET]
func PreviewFile(ctx context.Context, c *app.RequestContext) {
	fileId, err := strconv.ParseInt(string(c.Query("file_id")), 10, 64)
	if err != nil || fileId <= 0 {
		c.String(consts.StatusBadRequest, "invalid file_id")
		return
	}

	attname := ""
	if string(c.Query("download")) == "1" {
		attname = string(c.Query("filename"))
	}

	if string(c.Query("proxy")) == "1" {
		name, data, err := ser.FileLogic.PreviewContent(ctx, fileId)
		if err != nil {
			c.String(consts.StatusNotFound, err.Error())
			return
		}
		if attname == "" {
			attname = name
		}
		contentType := util.MimeTypeByName(name)
		if string(c.Query("download")) == "1" && attname != "" {
			c.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.PathEscape(attname)))
		}
		c.Data(consts.StatusOK, contentType, data)
		return
	}

	redirectURL, err := ser.FileLogic.PreviewRedirect(ctx, fileId, attname)
	if err != nil {
		c.String(consts.StatusNotFound, err.Error())
		return
	}

	c.Redirect(consts.StatusFound, []byte(redirectURL))
}
