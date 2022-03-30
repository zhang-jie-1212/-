package models

//validator进行参数校验
import (
	"github.com/go-playground/validator/v10"
	"time"
)

//注册参数
type ParamSignUp struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required" `
	//校验不出来Repassword字段不能为空
	Repassword string `json:"repassword" binding:"required,eqfield=Password"`
	//Date       string `json:"date" binding:"required,datetime=2016-01-02,checkDate"`
}

//登录参数
type LoginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//按规则排序的帖子参数
type ParamSortPost struct {
	Page  int    `json:"page" form:"page"`
	Size  int    `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

//某个社区下按规则排序的帖子参数
type ParamSortPostCommunity struct {
	*ParamSortPost
	CommunityID int64 `json:"community_id" form:"community_id"`
}

// SignUpParamStructLevelValidation 自定义SignUpParam结构体校验函数,使eqfield=Password"出错输出的是password
func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(ParamSignUp)
	if su.Password != su.Repassword {
		// 输出错误提示信息，最后一个参数就是传递的param
		sl.ReportError(su.Repassword, "repassword", "Repassword", "eqfield", "password")
	}
}

// customFunc 自定义字段级别校验方法:检验日期类型（1）是否为2006-01-02类型；（2）是否为未来日期（必须为未来日期）
//fl.Field()，validator.FiledLevel.Filed()，validator中存储的request中的参数（具体哪个字段，在validator初始化时进行了注册
//本例注册为checkDate字段
func lblCustomFunc(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	//校验是否是未来时间
	if date.Before(time.Now()) {
		return false
	}
	return true
}
