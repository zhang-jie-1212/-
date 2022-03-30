package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "liwenzhou.com"

//定义全局变量，方便controller层验证是哪个出错信息，并和Code对应，返回响应

// CheckUserExist判断已知用户名的用户是否存在
func CheckUserExist(username string) error {
	strsql := "select count(username) from user where username=?"
	var c int
	err := db.Get(&c, strsql, username)
	if err != nil {
		return err
	}
	if c > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser,向数据库user表中插入一条密码加密的用户数据
func InsertUser(user *models.User) error {
	strsql := "insert into user(user_id,username,password) values(?,?,?)"
	user.Password = encryptPassword(user.Password)
	_, err := db.Exec(strsql, user.UserID, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

//md5加密算法加密password
func encryptPassword(opassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(opassword)))
}

//CheckLogin：校验登录。用户不存在/用户名或密码错误
func CheckPasswordEq(user *models.User) (err error) {
	opassword := user.Password
	//查询数据库的密码
	strsql := "select user_id,username,password from user where username=?"
	if err = db.Get(user, strsql, user.Username); err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotExist
		}
		//查询数据库失败
		return err
	}
	//密码校验
	opassword = encryptPassword(opassword)
	if opassword != user.Password {
		return ErrorPassword
	}
	return nil
}

//获取用户信息
func GetUserDeatil(user_id int64) (user_detail *models.User, err error) {
	strsql := "select user_id,username from user where user_id=?"
	user_detail = new(models.User)
	err = db.Get(user_detail, strsql, user_id)
	return
}
