package mysql

import (
	"bluebell/models"
	"database/sql"
	"go.uber.org/zap"
)

// GetCommunityKind从数据库中查询出所有的社区分类name和id返回给前端
func GetCommunityKind() (data []*models.CommunityKind, err error) {
	strsql := "select community_id,community_name from community"
	err = db.Select(&data, strsql)
	if err != nil {
		//如果社区为空,不算错误，记录一个警告
		if err == sql.ErrNoRows {
			zap.L().Warn("there have not community kinds", zap.Error(err))
			err = nil
		}
	}
	return
}

// GetCommunityDetailu获取某个community_id的社区详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	details := new(models.CommunityDetail)
	var err error
	strsql := "select community_id,community_name,introduction,create_time " +
		"from community where community_id=?"
	err = db.Get(details, strsql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return details, err
}
