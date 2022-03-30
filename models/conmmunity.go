package models

import "time"

//所有社区分类
type CommunityKind struct {
	ID   int64  `json:"id,string" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

//社区详细信息 omitempty:如果字段为空，则返回给前端的json响应中不展示
//CreatTime为time.Time。而数据库中是timestamp类型，在链接数据库时要添加ParseTime=true&loc=Local，会自动进行类型转换
type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"creat_time" db:"create_time"`
}
