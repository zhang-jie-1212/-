package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

//CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.参数绑定和参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create post with invalid params", zap.Error(err))
		//不是validator校验错误
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.创建帖子：将details插入数据库
	//从token中获取的user_id
	userid, err := GetCurentUserID(c)
	fmt.Println(userid)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.UserID = userid
	zap.L().Info("get useid success ", zap.Int64("userid", userid))
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// PostDetailHnadler获取帖子及作者等详情信息
func PostDetailHandler(c *gin.Context) {
	//1.获取帖子id
	postId := c.Param("id")
	id, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		zap.L().Error("invalid param with post id", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	zap.L().Info("post_id", zap.Int64("post_id", id))
	//2.获取帖子信息，社区信息和作者信息
	all_detail, err := logic.GetAllDetail(id)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, all_detail)
}

//PostListHandler得到全部的帖子详情信息
func PostListHandler(c *gin.Context) {
	//1.获取页面和长度参数
	page, size := GetPostLimit(c)
	//2.获取全部帖子详细信息
	allDetail, err := logic.GetAllDetailList(page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
	}
	//3.返回响应
	ResponseSuccess(c, allDetail)
}

// GetSortPostListd得到按照一定规则排序的帖子的详细信息
func GetSortPostList(c *gin.Context) {
	//1.参数绑定 page size order communityID
	queryList := &models.ParamSortPostCommunity{
		ParamSortPost: &models.ParamSortPost{
			Page:  1,  //设置page默认值
			Size:  10, //设置size默认值
			Order: "time",
		},
	}
	if err := c.ShouldBindQuery(queryList); err != nil {
		zap.L().Error("c.ShouldBindQuery failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.得到按照参数排序的post详细信息
	if queryList.CommunityID == 0 {
		allDetailList, err := logic.GetSortPosts(queryList)
		if err != nil {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseSuccess(c, allDetailList)
		return
	} else {
		allDetailList, err := logic.GetCommunitySortPosts(queryList)
		if err != nil {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseSuccess(c, allDetailList)
		return
	}

}
