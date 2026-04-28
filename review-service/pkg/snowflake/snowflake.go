package snowflake

import (
	"errors"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var (
	InvalidInitParamErr = errors.New("sonwflake初始化失败, 无效的startTime或者machineID")
	InvalidTimeFormErr  = errors.New("sonwflake初始化失败, 无效的startTime格式")
	node                *sf.Node
)

// InitSnowflake 初始化雪花算法节点
// nodeID: 当前服务的节点 ID (0-1023)，在分布式系统中每个实例应不同
func InitSnowflake(startTime string, machineID int64) error {
	if len(startTime) == 0 || machineID <= 0 {
		return InvalidInitParamErr
	}
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return InvalidTimeFormErr
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return err
}

// GenerateUserID 生成一个 64 位的唯一 ID
func GenerateReviewID() int64 {
	return node.Generate().Int64()
}
