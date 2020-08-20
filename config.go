package gallop

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

func initConfig(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml") // 设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	viper.SetEnvPrefix("ENV")   // 读取环境变量的前缀为EG
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		log.Fatal(err)
	}
}
