package controller

import "C"
import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CommunityHandler(c *gin.Context) {
	//1.返回给前端所有社区分类name及ID列表（切片），进行展示
	data, err := logic.GetCommunityKind()
	if err != nil {
		zap.L().Error("logic.GetCommunity failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //没有必要将具体的错误信息返回给前端，自己记录日志即可
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler根据路径参数ID提取社区详细信息返回给前端
func CommunityDetailHandler(c *gin.Context) {
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	details, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, details)
}
