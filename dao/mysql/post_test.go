package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

//初始化db。在CreatePost函数中调用了db.Exec，但是db需要先初始化链接，而本次text仅测试的是CreatePost这个函数，所以会出现空指针引用
//另外定义一个初始化db的函数
//定义为init(),调用此模块时会自动调用init这个函数
func init() {
	dbcfg := settings.MysqlConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "root1234",
		DbName:       "bubble",
		Port:         13306,
		MaxOpenConns: 120,
		MaxIdleConns: 10,
	}
	err := Init(&dbcfg)
	if err != nil {
		panic(err)
	}
}

//数据库测试，创建帖子
func TestCreatePost(t *testing.T) {
	//调用CreatePost函数，传递一个数据，看一下返回的err
	post := models.Post{
		ID:          10,
		UserID:      123,
		CommunityID: 2,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("createpost insert into mysql failed,err:%v", err)
	}
	//成功的话，打印日志
	t.Logf("createpost insert record into mysql success")

}
