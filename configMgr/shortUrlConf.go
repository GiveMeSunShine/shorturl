/**
 * @Author : ysh
 * @Description :
 * @File : shortUrlConf
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/6 上午8:36
 */

package configMgr

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

type ShortUrlConf struct {
	Viper *viper.Viper
}

func NewShortUrlConf() (*ShortUrlConf, error) {
	cfg := viper.New()
	cfg.AddConfigPath("./configFile/")
	cfg.AddConfigPath("../configFile/")
	cfg.AddConfigPath("../../configFile/")
	cfg.SetConfigType("yaml")
	cfg.SetConfigName("shorturl") // name of config file (without extension)
	cfg.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	cfg.SetEnvKeyReplacer(replacer)
	err := cfg.ReadInConfig()
	if err != nil {
		log.Fatal("read config Err : ",err)
		panic(err)
	}
	shortConf := ShortUrlConf{
		Viper: cfg,
	}

	return &shortConf,nil
}

func (short *ShortUrlConf) GetShortConfigViper() *viper.Viper {
	return short.Viper
}

