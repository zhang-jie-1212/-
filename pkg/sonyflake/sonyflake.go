package sonyflake

//使用索尼的算法生成ID
import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var (
	node          *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}
func Init(startTime string, machineId uint16) (err error) {
	sonyMachineID = machineId
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	settings := sonyflake.Settings{
		StartTime: st,
		MachineID: getMachineID,
	}
	node = sonyflake.NewSonyflake(settings)
	return
}

//genId 生成ID
func GenID() (id uint64, err error) {
	if node == nil {
		err = fmt.Errorf("sony flake not inited")
		return
	}
	id, err = node.NextID()
	return
}

//func main() {
//	if err := Init("2020-07-01", 1); err != nil {
//		fmt.Printf("Init failed,err:%v\n", err)
//		return
//	}
//	id, _ := GenID()
//	fmt.Println(id)
//}
