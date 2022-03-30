package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//函数必须以Test开头
func TestCreatePostHandler(t *testing.T) {
	//1.设置gin的测试模式
	gin.SetMode(gin.TestMode)
	//2.创建一个测试的路由
	//r:=routers.SetupRouter()这里存在routers包和controller包循环引用的问题，所以在开发的时候包分的太细就会有问题。
	//自己创建路由
	r := gin.Default()
	url := "/air/v1/post"
	r.POST(url, CreatePostHandler)
	//3.进行测试（t对象放置测试结果/内容
	//创建一个请求参数和请求
	body := `{
	 "community_id":2,
		"title":"好吃的水果",
		"content":"香蕉"
}`
	//创建新的请求req,参数为请求方法，请求路径，以及路径参数/body
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	//创建存放响应的测试对象
	w := httptest.NewRecorder()
	//进行请求
	r.ServeHTTP(w, req)
	//4.查看响应内容是否符合预期,预期结果存入t testify/assert(断言)
	//判断响应状态码是否为200
	assert.Equal(t, 200, w.Code)
	//方法1：判断返回错误是是否为代码中的预期错误
	//assert.Contains(t, w.Body.String(), "用户未登录")
	//方法2 将w中body的内容解析到响应中进行判断
	//创建一个响应
	res := new(ResponseData)
	//将w.Body()内容反序列化到res中
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		//解析失败，向t中记录错误日志
		t.Errorf("json.Unmashal w.body failed,err:%v", err)
	}
	//判断响应码是否正确
	assert.Equal(t, res.Code, CodeNeedLogin)

}
