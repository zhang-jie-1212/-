//将配置文件config.yaml中的信息读取到viper中，再读取到结构体中
package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(multipleConfig)

type multipleConfig struct {
	*AppConfig   `mapstructure:"app"`
	*AutoConfig  `mapstructure:"auto"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}
type AutoConfig struct {
	AutoExpire int `mapstructure:"auto_expire"`
}
type AppConfig struct {
	//无论config是什么类型（后缀名称），统一用mapstructure来表示tag
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapdtructure:"port"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineId int64  `mapstructure:"machine_id"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backup"`
}
type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conn"`
	MaxIdleConns int    `mapstructure:"max_idle_conn"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

//初始化配置
//配置文件定位两种方式：
//1.直接指定配置文件路径（相对路径或绝对路径）
//相对路径：相对于执行程序的目录（即终端目录） （一般用这个）
//viper.SetConfigFile("./config.yaml")
//绝对路径：系统实际文件路径
//viper.SetConfigFile("E:/GoCode/GinBubble/webTemplate2/config.yaml")
//2.指定配置文件名和配置文件位置，viper执行查找可用的配置文件
//配置文件名不需要带后缀，需要注意指定配置目录内不要有同名的配置文件
//配置文件位置可配置多个,一次从文件路径查找
//viper.SetConfigName("config")
//viper.AddConfigPath(".")
//viper.AddConfigPath("./conf")
func Init() (err error) {
	//设置文件名称,若为本地config文件，不用指定文件类型（要求config文件只能有一个），SetConfigType语句不会起作用
	//当本地config有多个，需要加后缀：”config.json"/"config.yaml"
	viper.SetConfigName("config")
	//设置文件类型，专门用于从远程获取配置文件信息的时候指定配置文件类型
	viper.SetConfigType("yaml")
	//设置文件路径（相对路径）
	viper.AddConfigPath(".")
	//读取配置信息到viper
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("viper ReadInConfig failed,err:%v\n", err)
		return
	}
	//viper.SetDefault("start_time", "2020-07-01")
	//将viper中的信息反序列化到结构体变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Println("unmarshal config failed")
		return
	}
	fmt.Println(Conf)
	//t := viper.Get("app.start_time")
	//fmt.Printf("starttime is:%v,type is:%T\n", t, t)
	//监视配置信息变化
	viper.WatchConfig()
	//若变化，回调此函数重新写入文件信息到viper,并反序列化到Conf结构体中
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		//此处还可以加一些功能，通知项目组成员等
		//配置文件修改了，viper中的值重新读取到Conf中
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Println("unmarshal config failed")
			return
		}
	})
	return
}
