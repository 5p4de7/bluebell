package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConf)

type AppConf struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`

	*LogConf   `mapstructure:"log"`
	*MysqlConf `mapstructure:"mysql"`
	*RedisConf `mapstructure:"redis"`
}

type LogConf struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConf struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Db   int    `mapstructure:"db"`
}

func Init() (err error) {
	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return err
	}

	if err = viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed")
		if err = viper.Unmarshal(&Conf); err != nil {
			fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
			return
		}
	})
	return
}
