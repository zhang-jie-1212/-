package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//帖子投票
func PostVoteHandler(c *gin.Context) {
	//1.参数校验（post_id direction)user_id在校验token的时候可以获取到
	postvote := new(models.ParamPostVote)
	if err := c.ShouldBindJSON(postvote); err != nil {
		zap.L().Error("c.ShouldBindJSON failed", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseSuccess(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	zap.L().Info("bind JSON success", zap.Any("paramPostVote", postvote))
	//2.实现投票功能(更新redis中的分数）
	//根据token得到userID
	userID, err := GetCurentUserID(c)
	if err != nil {
		zap.L().Error("GetCurentUserID failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	zap.L().Info("Get UserID success", zap.Int64("userID", userID))
	err = logic.PostVote(userID, postvote)
	//3.返回响应
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
