// Package conf handles configurations, read configs, also generate default configs
package conf

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	// total retry when fetch ipaddres.com
	RETRY int
	// update hours
	CRON_SPEC string
	// lists of urls
	LISTS []string
	// database filename
	DB_NAME string
)

func Initconf(conf *viper.Viper) {
	conf.SetConfigType("toml")

	for k, v := range defaultconf {
		conf.SetDefault(k, v)
	}

}

// ReadConfig read config file, panic when failed
func ReadConfig(fpath string) *viper.Viper {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatalf("open config file failed: %s, if dosn't have one, generate it by `ghip -gen -conf file_name`", err)
	}
	conf := viper.New()
	Initconf(conf)
	if err := conf.ReadConfig(f); err != nil {
		log.Fatalf("read config file failed: %s", err)
	}
	return conf
}

// ParseConfigs will parse config file and save to conf vairalbes, MUSE RUN AT EACH INSTANSE
func ParseConfigs(conf *viper.Viper) {
	RETRY = conf.GetInt("net.retry")
	CRON_SPEC = conf.GetString("cron.spec")
	LISTS = conf.GetStringSlice("urls.list")
	DB_NAME = conf.GetString("database.file_path")
}

// GenDefaultConfig will generate default config file to given path
func GenDefaultConfig(fpath string) error {
	conf := viper.New()
	Initconf(conf)
	if err := conf.WriteConfigAs(fpath); err != nil {
		return fmt.Errorf("write default config failed: %s", err)
	}

	return nil
}
