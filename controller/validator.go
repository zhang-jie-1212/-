package controller

import (
	"bluebell/models"
	"fmt"
	"github.com/gin-gonic/gin/binding"    //gin框架中binding库，将校验器和翻译器连接到一起
	"github.com/go-playground/locales/en" //中文和英文包
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"                      //transfer的一个实例
	"github.com/go-playground/validator/v10"                                //validator校验器
	enTranslations "github.com/go-playground/validator/v10/translations/en" //validator中中英文的两个翻译包
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

//定义一个全局翻译器
var trans ut.Translator

//初始化一个全局翻译器，并绑定到validator校验器中。此部分在main。go中初始化即可
func InitTrans(locale string) (err error) {
	//拿到gin框架中的validator(翻译器）引擎，将其改为Validator.validate类型
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法；使validator校验错误的输出的参数字段是json中的字段而不是后端定义的字段
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		// 为SignUpParam注册自定义校验方法
		v.RegisterStructValidation(models.SignUpParamStructLevelValidation, models.ParamSignUp{})
		// 在校验器注册自定义的校验方法
		//if err := v.RegisterValidation("checkDate", models.CustomFunc); err != nil {
		//	return err
		//}

		//成功，注册两个（中英文）翻译器
		zhT := zh.New()
		enT := en.New()
		//创建一个全局翻译器
		// 第一个参数是备用（fallback）的语言环境，后面的参数是应该支持的语言环境（支持多个）
		//uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)
		//将要翻译的类型传给uni翻译器（locale）,实现一个实例化的翻译器
		//locale一般取决于http请求头中的”Accep-Language"
		var ok bool
		trans, ok = uni.GetTranslator(locale)
		//翻译器实例不成功，返回信息
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		//注册翻译器到校验器,让那个validator将娇艳的错误翻译成中文发送给前端
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default: //默认中文翻译器
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}
