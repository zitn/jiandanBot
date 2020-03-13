package types

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigName("config") //配置文件名
	viper.AddConfigPath("config") //配置文件所在的路径
	viper.SetConfigType("json")   //配置文件类型
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
}
