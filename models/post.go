package models

import "time"

//帖子详情
type Post struct {
	ID          int64     `json:"id" db:"post_id"`      //插入数据库时雪花算法生成，不用在参数中
	UserID      int64     `json:"user_id" db:"user_id"` //在校验token时获取
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreatTime   time.Time `json:"create_time" db:"create_time"`
}

//帖子详情页展示的总详情：帖子详情，社区详情，用户姓名
type AllDetail struct {
	UserName         string `json:"user_name" db:"username"`
	*Post            `json:"post"`
	*CommunityDetail `json:"community"`
}

//按一定规则查询帖子相关信息
type AllSortDetail struct {
	UserName         string `json:"user_info"`
	*Post            `json:"post_info"`
	Score            int64 `json:"score"`
	*CommunityDetail `json:"community_info"`
}
