package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"errors"
	"go.uber.org/zap"
)

// CreatePost创建帖子
func CreatePost(p *models.Post) (err error) {
	//1.生成post_id
	id := snowflake.GenID()
	zap.L().Info(" set postid", zap.Int64("id", id))
	post := models.Post{
		ID:          id,
		UserID:      p.UserID,
		CommunityID: p.CommunityID,
		Title:       p.Title,
		Content:     p.Content,
	}
	zap.L().Info("set postid success ", zap.Int64("postid", post.ID))
	//2.插入mysql数据库
	err = mysql.CreatePost(&post)
	if err != nil {
		zap.L().Error("mysql.CreatePost failed", zap.Error(err))
		return
	}
	//3.插入redis:将创建的postID和其create_time分别保存到redis的KeyPostTime、KeyPostScore、KeyCommunitySetPF中
	err = redis.CreatePost(post.ID, post.CommunityID)
	return
}

// GetPostDetail查找帖子、用户和社区详细信息
func GetAllDetail(post_id int64) (all_detail *models.AllDetail, err error) {
	//1.从数据库中查找帖子信息
	postDetail, err := mysql.GetPostDetail(post_id)
	if err != nil {
		zap.L().Error(" mysql.GetPostDetail failed", zap.Error(err))
		return
	}
	//2.从数据库中查找user信息
	userDetail, err := mysql.GetUserDeatil(postDetail.UserID)
	if err != nil {
		zap.L().Error("mysql.GetUserDeatil failed", zap.Error(err))
		return
	}
	//3.从数据库中查找社区信息
	communityDetail, err := mysql.GetCommunityDetail(postDetail.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
		return
	}
	//4.信息合并并返回
	all_detail = &models.AllDetail{
		UserName:        userDetail.Username,
		Post:            postDetail,
		CommunityDetail: communityDetail,
	}
	return
}

// GetAllDetailList获取全部帖子详情列表信息
func GetAllDetailList(page int64, size int64) (allDeatil []*models.AllDetail, err error) {
	//1.获取帖子列表
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
	}
	allDeatil = make([]*models.AllDetail, 0, size)
	//2.遍历列表拿到每个数据
	for _, post := range postList {
		//1.从数据库中查找user信息
		userDetail, err := mysql.GetUserDeatil(post.UserID)
		if err != nil {
			zap.L().Error("mysql.GetUserDeatil failed", zap.Error(err))
			continue
		}
		//2.从数据库中查找社区信息
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
			continue
		}
		aDetail := &models.AllDetail{
			UserName:        userDetail.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		allDeatil = append(allDeatil, aDetail)
	}
	return
}

// GetAllDetailbyPostID 通过postID获得帖子相关信息的调用函数
func GetAllDetailbyPostID(postIDList []string) (allDetailList []*models.AllSortDetail, err error) {
	//2.按照帖子ID从mysql数据库中取出相关详细信息
	//根据postIDList取出帖子信息列表
	sortPostDetail, err := mysql.GetSortPostDetail(postIDList)
	if err != nil {
		return
	}
	//3.从redis中拿到每个postID相应的赞成票票数（拼成切片）
	scores, err := redis.GetPostScore(postIDList)
	if err != nil {
		return
	}
	allDetailList = make([]*models.AllSortDetail, 0, len(postIDList))
	//4.遍历帖子信息列表，分别取出帖子分数，用户详细信息和社区详细信息
	for idx, post := range sortPostDetail {
		//1.从数据库中查找user信息
		userDetail, err := mysql.GetUserDeatil(post.UserID)
		if err != nil {
			zap.L().Error("mysql.GetUserDeatil failed", zap.Error(err))
			continue
		}
		//2.从数据库中查找社区信息
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail failed", zap.Error(err))
			continue
		}
		aDetail := &models.AllSortDetail{
			UserName:        userDetail.Username,
			Post:            post,
			Score:           scores[idx],
			CommunityDetail: communityDetail,
		}
		allDetailList = append(allDetailList, aDetail)
	}
	//5.返回排序帖子相关信息列表和err
	return
}

// GetSortPosts 获取安规则排序的帖子的相关信息
func GetSortPosts(querList *models.ParamSortPostCommunity) (allDetailList []*models.AllSortDetail, err error) {
	//1.从redis中按排序规则取出page 和size规定的postID排列list
	postIDList, err := redis.GetPostIDList(querList)
	if err != nil {
		zap.L().Error("redis.GetSortPosts failed", zap.Error(err))
		return
	}
	return GetAllDetailbyPostID(postIDList)

}

// GetCommunitySortPosts 获取某个社区下按规则分类的帖子的详细信息
func GetCommunitySortPosts(queryList *models.ParamSortPostCommunity) (allDetailList []*models.AllSortDetail, err error) {
	//1.从redis中取出按社区存储的postid
	postIDList, err := redis.GetCommunityPostList(queryList)
	if err != nil {
		return
	}
	if len(postIDList) == 0 {
		err = errors.New("没有任何帖子")
		zap.L().Error("redis.GetCommunityPostList have no post", zap.Error(err))
	}
	//2.调用中间件函数，根据postID获取帖子相应信息
	return GetAllDetailbyPostID(postIDList)
}
