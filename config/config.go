package config

import (
	"github.com/spf13/viper"
	"gofuse/single"
	"log"
)

type Config struct {
	Flags *viper.Viper
}

var GetConfig = single.DoSingleWithParam(func(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
		return nil
	}
	return &Config{v}
})

func (c *Config) All() map[string]any {
	return c.Flags.AllSettings()
}

func (c *Config) AllKeys() []string {
	return c.Flags.AllKeys()
}
