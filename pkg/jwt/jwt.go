package jwt

//定义了JWT Tocken的生成和解密算法
import (
	"bluebell/settings"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//定义过期时间:2小时

//定义加密字段：将此字段一起加密生成Tocken,没有此字段就没有办法真正解密
var MySecret = []byte("肖战王一博好帅")

//定义JWT的负载信息，由自定义字段和官方字段构成
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

//生成JWT Tocken
func GetTocken(userid int64, username string) (string, error) {
	//定义负载信息实例
	c := MyClaims{
		UserID:   userid,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(settings.Conf.AutoExpire)).Unix(), // 过期时间
			Issuer:    "bluebell",                                                                 // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

//解析Tocken
func ParseTocken(tokenString string) (*MyClaims, error) {
	mc := new(MyClaims)
	//将tockenstring解析到mc（不加言）中，返回解析的token（加言），第三个参数是获取加言（加一个语句）MySecret的参数
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	//若token有效
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid tocken")

}
