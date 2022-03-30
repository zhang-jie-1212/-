package models

type ParamPostVote struct {
	//接收的时候postID直接接受为string类型，后面操作redis数据库用的是string类型
	PostID    string `json:"post_id" db:"post_id" binding:"required"`
	Direction int8   `json:"direction" db:"direction" binding:"oneof=-1 0 1"`
}
