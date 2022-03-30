package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(q *models.ParamSignUp) (err error) {
	//1.判断数据库中该用户是否存在
	if err = mysql.CheckUserExist(q.Username); err != nil {
		return err
	}
	//2.生成唯一ID和数据库用户实例
	ID := snowflake.GenID()
	//生成数据实例
	user := models.User{
		UserID:   ID,
		Username: q.Username,
		Password: q.Password,
	}
	//3.插入数据库
	err = mysql.InsertUser(&user)
	return
}

// Login进行登录用户:用户名和密码与数据库的校验,返回UserID
func Login(q *models.LoginParams) (string, error) {
	//数据库中取出该用户的密码和username,(username不存在，取出失败）
	u := models.User{
		Username: q.Username,
		Password: q.Password,
	}
	//判断密码是否正确，直接更改u为数据库信息，之后生成token用
	if err := mysql.CheckPasswordEq(&u); err != nil {
		return "", err
	}
	//正确表示登录成功，创建Tocken
	return jwt.GetTocken(u.UserID, u.Username)
}
