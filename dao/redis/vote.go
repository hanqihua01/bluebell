package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("voting period has expired")
)

func CreatePost(postID, communityID int64) error {
	pipeline := rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()), // 初试分数就是当前时间
		Member: postID,
	})
	// 把帖子id加到社区的set
	ckey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ckey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 判断投票限制
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val() // 获取帖子发布时间
	if float64(time.Now().Unix()-int64(postTime)) > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 更新帖子分数
	// 先查之前的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	diff := value - ov // 计算两次投票的差值

	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), diff*scorePerVote, postID)
	// 记录用户
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
