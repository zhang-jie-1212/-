package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"
)

//投票功能
//
func PostVote(userID int64, p *models.ParamPostVote) error {
	return redis.PostVote(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
