package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"fmt"
)

// GetCommunityKind获取社区所有分类列表
func GetCommunityKind() (data []*models.CommunityKind, err error) {
	return mysql.GetCommunityKind()
}

// GetCommunityDeatil获取特定社区的详细信息
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	fmt.Println(id)
	return mysql.GetCommunityDetail(id)
}
