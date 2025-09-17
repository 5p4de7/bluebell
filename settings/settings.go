package settings

import "github.com/spf13/viper"

func Init() {
	viper.SetConfigFile("config.yaml")

}
