//雪花算法生成ID
package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node //创建一个节点，节点的Generate()方法实现生成ID
//初始化节点
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		fmt.Println("parse time failed...")
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

//生成唯一ID
func GenID() int64 {
	return node.Generate().Int64()
}

//测试
func main() {
	if err := Init("2020-07-01", 1); err != nil {
		fmt.Printf("Init failed,err:%v\n", err)
		return
	}
	id := GenID()
	fmt.Println(id)
}
