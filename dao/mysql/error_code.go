package mysql

import "errors"

var (
	ErrorUserExist    = errors.New("用户已存在！")
	ErrorUserNotExist = errors.New("用户不存在！")
	ErrorPassword     = errors.New("密码错误！")
	ErrorInvalidID    = errors.New("无效的ID！")
	ErrorPostNotExist = errors.New("帖子不存在！")
)
