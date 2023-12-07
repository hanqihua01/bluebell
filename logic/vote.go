package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"
)

// 本项目使用简化版的投票规则
// 投一票就加432分    86400/200  ->  200张赞成票可以给帖子续一天
/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票
	2. 之前投过反对票，现在改投赞成票
direction=0时，有两种情况：
	1. 之前投过赞成票，现在要取消投票
	2. 之前投过反对票，现在要取消投票
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票
	2. 之前投过赞成票，现在改投反对票
投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了
	1. 到期之后将redis中保存的赞成票和反对票存储到mysql表中
	2. 到期之后删除KeyPostVotedZSetPrefix
*/
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
