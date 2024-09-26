package config

import "github.com/phillippbertram/xc-strings/internal/constants"

type Config struct {
	BaseStringsFile string `mapstructure:"baseStringsFile"`
	StringsPath     string `mapstructure:"stringsPath"`
	SwiftPath       string `mapstructure:"swiftPath"`
}

var Cfg Config = Config{
	StringsPath: constants.DefaultStringsGlob,
}
