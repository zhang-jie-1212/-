package redis

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"math"
	"strconv"
	"time"
)

const (
	PostTimeExpire = 7 * 24 * 3600
	scorePerVote   = 432
)

/*
1.帖子时间限制：
	每个帖子仅能存活一周的时间(秒），超过一周不能投票
2.帖子的投票规则
direction=1,有两种情况：
	1.之前从未投过票，现在投赞成    -->更新记录和投票分数 差值的绝对值：1   +432
	2.之前投反对票，现在投赞成票    -->更新记录和投票分数  差值的绝对值：2  +432*2
direction=0，有两种情况：
	1.以前投赞成票，现在取消投票    -->更新记录和投票分数 差值的绝对值：1   +432
	2.之前投反对票，现在取消投票    -->更新记录和投票分数  差值的绝对值：1  -432
direction=-1，有两种情况：
	1.以前从未投过票，现在投反对    -->更新记录和投票分数 差值的绝对值：1   -432
	2.以前投赞成票，现在投反对票    -->更新记录和投票分数  差值的绝对值：2  -432*2



*/
func PostVote(userID string, postID string, direction float64) error {
	//1.根据post_id查看帖子是否过期（帖子及其日期信息存放在redis的KeyPostTime中)
	postTime, err := rdb.ZScore(GetPostKey(KeyPostTimeZset), postID).Result()
	if err != nil {
		zap.L().Error("rdb.ZScore GetPostKey(KeyPostTimeZset) failed", zap.Error(err))
		err = ErrorPostNotExist
		return err
	}
	//计算是否过期
	if (float64(time.Now().Unix()) - postTime) > PostTimeExpire {
		err = ErrorPostHasExpire
		return err
	}
	info := "帖子没有过期"
	zap.L().Info("post time expire?", zap.Any("result", info))
	//2.根据规则，改变帖子分数
	//获取之前的投票分数，找不到则返回类型零值
	ov := rdb.ZScore((GetPostKey(KeyPostVoteZsetPre + postID)), userID).Val()
	zap.L().Info("post oringinal score", zap.Float64("postScore", ov))
	//计算差值绝对值
	weight := math.Abs(direction - ov)
	//计算加号还是减号
	var op float64
	if direction > ov {
		op = 1
	} else {
		op = -1
	}
	//将投票情况和分数改变更新数据库（事物）
	//pipeline := rdb.TxPipeline()
	//更新分数
	_, err = rdb.ZIncrBy(GetPostKey(KeyPostScoreZset), float64(scorePerVote)*weight*op, postID).Result()
	if err != nil {
		zap.L().Error("update redis post score failed", zap.Error(err))
		return err
	}
	//更新用户选择
	_, err = rdb.ZAdd(GetPostKey(KeyPostVoteZsetPre+postID), redis.Z{
		Score:  direction,
		Member: userID,
	}).Result()
	if err != nil {
		zap.L().Error("update redis post vote user failed", zap.Error(err))
	}
	//_, err = pipeline.Exec()
	//zap.L().Error("redis.vote.pipeline failed", zap.Error(err))
	return nil
}

//向redis.KeyPostTimeZset中插入当前时间
func CreatePost(postID, communityID int64) (err error) {
	//1.两个添加操作为一个事物
	//pipeline := rdb.TxPipeline()
	//2.添加到KeyPostTimeZset
	_, err = rdb.ZAdd(GetPostKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()
	if err != nil {
		zap.L().Error("rdb ZAdd post create time failed", zap.Int64("postid", postID), zap.Float64("createTim", float64(time.Now().Unix())), zap.Error(err))
		return
	}
	//3.添加到KeyPostScoreZset
	_, err = rdb.ZAdd(GetPostKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()
	if err != nil {
		zap.L().Error("rdb ZAdd post create time failed", zap.Int64("postid", postID), zap.Float64("postScore", float64(time.Now().Unix())), zap.Error(err))
		return
	}
	//4.添加到KeyCommunitySetPF
	_, err = rdb.SAdd(GetPostKey(KeyCommunitySetPF+strconv.Itoa(int(communityID))), postID).Result()
	if err != nil {
		zap.L().Error("rdb.SAdd KeyCommunitySetPF failed", zap.Error(err))
	}
	return
	//_, err := pipeline.Exec()
	//if err != nil {
	//	zap.L().Error("redis.CreatPost pipeline failed", zap.Error(err))
	//}
	//return err
}
