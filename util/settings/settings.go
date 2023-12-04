package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf是全局变量，用来保存程序的所有配置信息，被其他包读取
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(confFile string) (err error) {
	viper.SetConfigFile(confFile) // 指定配置文件
	err = viper.ReadInConfig()    // 读取配置文件
	if err != nil {               // 读取配置文件失败
		fmt.Printf("viper.ReadInConfig() failed, err: %v\n", err)
		return
	}

	// 反序列化配置信息到Conf中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
	}

	// 监视配置文件变化并做出处理（回调函数）
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
		}
	})

	return
}
