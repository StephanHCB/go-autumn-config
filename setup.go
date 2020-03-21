package auconfig

import (
	"fmt"
	"github.com/StephanHCB/go-autumn-config-api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"regexp"
)

var configPath string
var secretsPath string

var failFunction auconfigapi.ConfigFailFunc = fail
var warnFunction auconfigapi.ConfigWarnFunc = warn

var configItems []auconfigapi.ConfigItem

var configItemKeysWithNoFlags = map[string]bool{}

// initialize configuration with full setup - you need to call this from your code
func Setup(items []auconfigapi.ConfigItem, failFunc auconfigapi.ConfigFailFunc, warnFunc auconfigapi.ConfigWarnFunc) {
	SetupWithOverriddenConfigPath(items, failFunc, warnFunc, "", "")
}

// load any configuration files - you need to call this from your code after calling Setup()
func Load() {
	performLoad()
	validate()
}

// use this for unit tests.
//
// This just sets all configuration settings to their default values. No need to call Load() after this.
func SetupDefaultsOnly(items []auconfigapi.ConfigItem, failFunc auconfigapi.ConfigFailFunc, warnFunc auconfigapi.ConfigWarnFunc) {
	configItems = items
	failFunction = failFunc
	warnFunction = warnFunc

	setupDefaults()
}

// use this for integration tests instead of Setup().
//
// This allows you to specify a default path for both config and secrets files, avoiding the need for command line parameters in integration tests.
// You still need to call Load(). Set defaultSecretsPath to "" to disable loading it.
func SetupWithOverriddenConfigPath(items []auconfigapi.ConfigItem, failFunc auconfigapi.ConfigFailFunc, warnFunc auconfigapi.ConfigWarnFunc, defaultConfigPath string, defaultSecretsPath string) {
	configItems = items
	failFunction = failFunc
	warnFunction = warnFunc

	initializeFlags(defaultConfigPath, defaultSecretsPath)
	pflag.Parse()

	setupDefaults()
	setupEnv()
	setupFlags()
}

func initializeFlags(defaultConfigPath string, defaultSecretsPath string) {
	pflag.StringVar(&configPath, "config-path", defaultConfigPath, "config file path without file name")
	pflag.StringVar(&secretsPath, "secrets-path", defaultSecretsPath, "secrets file path without file name")

	for _, item := range configItems {
		flagname := item.FlagName
		if flagname == "" {
			flagname = item.Key
		}

		// must set null default values here, or else this value will overwrite config values from config file
		if _, ok := item.Default.(string); ok {
			pflag.String(flagname, "", item.Description)
		} else if _, ok := item.Default.(int); ok {
			pflag.Int(flagname, 0, item.Description)
		} else if _, ok := item.Default.(int8); ok {
			pflag.Int8(flagname, 0, item.Description)
		} else if _, ok := item.Default.(int16); ok {
			pflag.Int16(flagname, 0, item.Description)
		} else if _, ok := item.Default.(int32); ok {
			pflag.Int32(flagname, 0, item.Description)
		} else if _, ok := item.Default.(uint); ok {
			pflag.Uint(flagname, 0, item.Description)
		} else if _, ok := item.Default.(uint8); ok {
			pflag.Uint8(flagname, 0, item.Description)
		} else if _, ok := item.Default.(uint16); ok {
			pflag.Uint16(flagname, 0, item.Description)
		} else if _, ok := item.Default.(uint32); ok {
			pflag.Uint32(flagname, 0, item.Description)
		} else if _, ok := item.Default.(bool); ok {
			pflag.Bool(flagname, false, item.Description)
		} else if _, ok := item.Default.([]string); ok {
			pflag.StringSlice(flagname, []string{}, item.Description)
		} else {
			configItemKeysWithNoFlags[item.Key] = true
			warnFunction(fmt.Sprintf("unsupported data type for config item %v, cannot initialize command line argument, skipping this key", item.Key))
		}
	}
}

func setupFlags() {
	for _, item := range configItems {
		if _, ok := configItemKeysWithNoFlags[item.Key]; !ok {
			flagname := item.FlagName
			if flagname == "" {
				flagname = item.Key
			}

			bindFlag(item.Key, flagname)
		}
	}
}

func bindFlag(key string, flagname string) {
	err := viper.BindPFlag(key, pflag.Lookup(flagname))
	if err != nil {
		failFunction(fmt.Errorf("Fatal error could not bind configuration flag %s: %s\n", key, err))
	}
}

func setupDefaults() {
	for _, item := range configItems {
		viper.SetDefault(item.Key, item.Default)
	}
}

func setupEnv() {
	re := regexp.MustCompile(`[^a-z0-9]`)
	for _, item := range configItems {
		// simply fill in EnvName if unset
		if item.EnvName == "" {
			item.EnvName = "CONFIG_" + re.ReplaceAllString(item.Key, "_")
		}

		// the only error that can occur is when the Key is empty
		_ = viper.BindEnv(item.Key, item.EnvName)
	}
}

func validate() {
	for _, item := range configItems {
		err := item.Validate(item.Key)
		if err != nil {
			fail(err)
		}
	}
}

func fail(err error) {
	panic(err)
}

func warn(message string) {
	log.Print(message)
}
