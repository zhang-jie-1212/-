package redis

const (
	KeyPrefix          = "bluebell:"  //项目所有key的前缀，需要手动和下面进行拼接，但是可以查询本项目所有的key
	KeyPostTimeZset    = "post:time"  //数据类型为后缀：zset,帖子及其时间戳信息，每当创建一个帖子，就添加进redis
	KeyPostScoreZset   = "post:score" //帖子的投票信息，本key中存放帖子及其分数
	KeyPostVoteZsetPre = "post:vote:" //帖子的投票信息，需要和post_id拼接，存放所有投票人user_id及其分数
	KeyCommunitySetPF  = "community:" //按社区分类的帖子信息：communityID:postID的set
)

// GetPostKey 获取post的key
func GetPostKey(key string) string {
	return KeyPrefix + key
}
