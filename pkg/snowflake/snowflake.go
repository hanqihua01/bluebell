package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime) // 解析startTime
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano()          // 设置snowflake起始时间
	node, err = sf.NewNode(machineID) // 创建snowflake节点
	return
}

// 生成一个新的ID
func GenID() int64 {
	return node.Generate().Int64()
}
