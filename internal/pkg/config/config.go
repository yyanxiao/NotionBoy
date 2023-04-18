package config

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var globalConfig = &Config{}

const (
	SETTINGS_FOLDER        = "."
	SETTINGS_NAME          = "settings"
	SETTINGS_OVERRIDE_NAME = "settings_local"
	SETTINGS_ENV_PREFIX    = "app"
)

func init() {
	Load(globalConfig)
}

func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = &Config{}
	}
	return globalConfig
}

// Load use viper to load config
// read settings.{yaml/json/toml} and unmarshal to interface cfg
// settings_local will override config in settings
// support auto reload when file change
// auto read config from env with prefix 'APP_' and will auto convert "_" to "."
// for example APP_TELEGRAM_TOKEN will map to telegram.token
// avoid to use "_" in config name for better env reading
func Load(cfg interface{}) {
	// read configs from seetings in root directory

	viper.AddConfigPath(getConfigDir())
	viper.AddConfigPath(SETTINGS_FOLDER)
	viper.SetConfigName(SETTINGS_NAME)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// read configs from settings_local
	// this will override configs from settings
	viper.SetConfigName(SETTINGS_OVERRIDE_NAME)
	// ignore when settings_local not exists
	_ = viper.MergeInConfig()

	// env should start with app
	// should use '_' instead of '.'
	// for example APP_TELEGRAM_TOKEN map to telegram.token
	viper.AutomaticEnv()
	viper.SetEnvPrefix(SETTINGS_ENV_PREFIX)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	loadConfig := func() {
		if err := viper.Unmarshal(cfg); err != nil {
			panic(fmt.Errorf("unable to decode into struct: %v", err))
		}
	}
	// autoreload watch config change
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Reload config since file changed:", e.Name)
		loadConfig()
	})
	loadConfig()
	viper.WatchConfig()
}

func getConfigDir() string {
	_, filename, _, _ := runtime.Caller(1)
	filepath := path.Join(path.Dir(filename), "../../../")
	return filepath
}
