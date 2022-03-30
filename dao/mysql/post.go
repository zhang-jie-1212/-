package mysql

import (
	"bluebell/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

// CreatePost插入一条帖子
func CreatePost(p *models.Post) (err error) {
	strsql := "insert into post(post_id,title,content,community_id,user_id)" +
		" values(?,?,?,?,?)"
	_, err = db.Exec(strsql, p.ID, p.Title, p.Content, p.CommunityID, p.UserID)
	return
}

// GetPostDetail获取帖子详情
func GetPostDetail(postId int64) (postDetail *models.Post, err error) {
	strsql := "select post_id,title,content,user_id,community_id,create_time from post where post_id=?"
	postDetail = new(models.Post)
	err = db.Get(postDetail, strsql, postId)
	if err == sql.ErrNoRows {
		err = ErrorPostNotExist
	}
	return

}

// GetPostList获取全部贴子信息列表
func GetPostList(page int64, size int64) (postList []*models.Post, err error) {
	//limit后面跟着数据库的偏移量和查询数据量
	strsql := "select post_id,title,content,user_id,community_id from post limit ?,?"
	postList = make([]*models.Post, 0, size)
	err = db.Select(&postList, strsql, (page-1)*size, size)
	if err == sql.ErrNoRows {
		err = ErrorPostNotExist
	}
	return
}

// GetSortPostDetail获取排序的帖子信息列表
func GetSortPostDetail(postIDList []string) (sortPostDetail []*models.Post, err error) {
	//1.sqlx IN在切片中的帖子信息
	strsql := "select post_id,title,content,community_id,user_id,create_time " +
		"from post " +
		"where post_id in (?) " +
		"order by FIND_IN_SET(post_id,?)"
	query, args, err := sqlx.In(strsql, postIDList, strings.Join(postIDList, ","))
	if err != nil {
		zap.L().Error("mysql.GetSortPostDetail sqlx.In failed", zap.Error(err))
		return
	}
	query = db.Rebind(query)
	sortPostDetail = make([]*models.Post, 0, len(postIDList))
	err = db.Select(&sortPostDetail, query, args...) //注意这里sortPostDeatil还要取址，防止对切片进行appendbl扩容
	//zap.L().Info("post detail", zap.Any("post details", sortPostDetail))
	if err != nil {
		zap.L().Error("mysql.GetSortPostDetail db.Select failed", zap.Error(err))
	}
	return
}
