package redis

import (
	"bluebell/models"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

//封装通过key获得postID的公用部分
func GetPostIDbyKey(key string, page, size int) (postIDList []string, err error) {
	//2.根据排序规则以及page 和size从redis中取出相应的值
	start := (page - 1) * size
	stop := start + size - 1
	postIDList, err = rdb.ZRevRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		zap.L().Error("redis.GetPostIDbyKey failed", zap.Error(err))
		return
	}
	if len(postIDList) == 0 {
		err = errors.New("postID为空")
		zap.L().Error("redis.GetPostIDList get invalid PostIDList", zap.Error(err))
	}
	fmt.Println(postIDList)
	return
}

// GetSortPosts 从redis中取出按一定规则排序的PostID
func GetPostIDList(queryList *models.ParamSortPostCommunity) (postIDList []string, err error) {
	//1.通过排序规则获取redis的键,m默认为通过时间排序
	key := GetPostKey(KeyPostTimeZset)
	if queryList.ParamSortPost.Order == "score" {
		key = GetPostKey(KeyPostScoreZset)
	}
	postIDList, err = GetPostIDbyKey(key, queryList.Page, queryList.Size)
	return
}

// GetPostScore获取帖子投赞成票的人数
func GetPostScore(postIDList []string) (scores []int64, err error) {
	//1.pipeline得到帖子分数;不用pipeline,for range每次遍历都要请求链接数据库，pipeline只用链接释放一次
	pipeline := rdb.Pipeline()
	for _, postID := range postIDList {
		key := GetPostKey(KeyPostVoteZsetPre + postID)
		pipeline.ZCount(key, "1", "1")
	}
	cmder, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("redis.GetPostScore pipeline Exec failed", zap.Error(err))
		return
	}
	scores = make([]int64, 0, len(postIDList))
	for _, cmd := range cmder {
		val := cmd.(*redis.IntCmd).Val()
		scores = append(scores, val)
	}
	return
}

//GetCommunityPostList 取到按社区存储的postID
func GetCommunityPostList(queryList *models.ParamSortPostCommunity) (postIDList []string, err error) {
	//1.使用ZInterStore将comunityID:postID set和key:postID:order(time/score)zset联合存储
	//针对新的zset按之间的逻辑取PostIList
	//利用缓存key(新创建的zset的key)减少zinterstore执行的次数
	//社区的key
	cKey := GetPostKey(KeyCommunitySetPF + strconv.Itoa(int(queryList.CommunityID)))
	//排序的key
	orderKey := GetPostKey(KeyPostTimeZset)
	if queryList.Order == "score" {
		orderKey = GetPostKey(KeyPostScoreZset)
	}
	newKey := orderKey + strconv.Itoa(int(queryList.CommunityID))
	//如果newKey不存在，则先进行联合存储
	if rdb.Exists(newKey).Val() < 1 {
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(newKey, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(newKey, 60*time.Second) //联合存储形成的新的zset只有60秒的存在事假
		_, err = pipeline.Exec()
		if err != nil {
			zap.L().Error("redis.GetCommunityPostList pipleline.ZInterStore failed", zap.Error(err))
			return
		}
	}
	//否则直接从key中查找postID的值
	postIDList, err = GetPostIDbyKey(newKey, queryList.Page, queryList.Size)
	return
}
