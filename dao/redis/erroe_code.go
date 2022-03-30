package redis

import "errors"

var (
	ErrorPostNotExist   = errors.New("帖子不存在！")
	ErrorPostHasExpire  = errors.New("帖子已过期！")
	ErrorPostIDNotExist = errors.New("帖子ID不存在")
)
