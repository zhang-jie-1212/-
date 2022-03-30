package controller

//错误码和返回信息的定义
type TypeCode int64

const (
	CodeSuccess TypeCode = 1000 + iota
	CodeInvalidParam
	CodeInvalidPassword
	CodeExistUser
	CodeNotExistUser
	CodeServerBusy
	CodeInvaildToken
	CodeNeedLogin
)

var CodeMsgMap = map[TypeCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeInvalidPassword: "密码错误",
	CodeExistUser:       "用户已存在",
	CodeNotExistUser:    "用户不存在",
	CodeServerBusy:      "服务器忙",
	CodeInvaildToken:    "无效的token",
	CodeNeedLogin:       "用户未登录",
}

//定义通过code获取codemsgmap的方法
func (c TypeCode) GetMsg() string {
	msg, ok := CodeMsgMap[c]
	if !ok {
		msg = CodeMsgMap[CodeServerBusy]
	}
	return msg
}
