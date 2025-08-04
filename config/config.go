package config

import (
	"github.com/spf13/viper"
	"gofuse/errs"
	"gofuse/single"
)

type Config struct {
	*viper.Viper
}

var GetConfig = single.DoSingleWithParam(func(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		errs.FatalErr(errs.ErrConfigRead, err)
		return nil
	}
	return &Config{v}
})
