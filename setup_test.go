package auconfig

import (
	"fmt"
	"testing"

	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// tstValidateNonEmptyString is an example for a validation function.
//
// If a configuration item does not pass it, loading the configuration fails.
func tstValidateNonEmptyString(key string) error {
	value := viper.GetString(key)
	if value == "" {
		return fmt.Errorf("%s is empty", key)
	}
	return nil
}

// tstConfigItems contains examples for all available types of configItems.
var tstConfigItems = []auconfigapi.ConfigItem{
	{
		Key:         "string-key-no-default",
		Default:     "",
		Description: "required config item of string type with no default, using default env and flag names, not allowing empty value",
		Validate:    tstValidateNonEmptyString,
	},
	{
		Key:         "string-key-with-default",
		Default:     "default",
		Description: "optional config item of string type with default, explicitly specified flag and env name, no validation",
		EnvName:     "CONFIG_CUSTOM_string_key_with_default",
		FlagName:    "custom-string-key-with-default",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "secret-key",
		Default:     "",
		Description: "required secret of string type",
		Validate:    tstValidateNonEmptyString,
	},
	{
		Key:         "int-key",
		Default:     24,
		Description: "optional config item of int type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "int8-key",
		Default:     int8(3),
		Description: "optional config item of int8 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "int16-key",
		Default:     int16(30000),
		Description: "optional config item of int16 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "int32-key",
		Default:     int32(2000000000),
		Description: "optional config item of int32 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "int64-key",
		Default:     int64(8000000000),
		Description: "optional config item of int64 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "uint-key",
		Default:     uint(1234),
		Description: "optional config item of uint type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "uint8-key",
		Default:     uint8(12),
		Description: "optional config item of uint8 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "uint16-key",
		Default:     uint16(65535),
		Description: "optional config item of uint16 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "uint32-key",
		Default:     uint32(65535000),
		Description: "optional config item of uint16 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "uint64-key",
		Default:     uint64(65535000000),
		Description: "optional config item of uint64 type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "bool-key",
		Default:     false,
		Description: "optional config item of boolean type",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key:         "string-slice-key",
		Default:     []string{}, // empty default, but must set so type is recognized
		Description: "optional config item of string slice type",
		EnvName:     "CONFIG_string_slice_key", // the default env name, set explicitly for documentation purposes
		FlagName:    "string-slice-key",        // the default flag name, set explicitly for documentation purposes
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	},
	{
		Key: "map-key",
		Default: map[string]string{
			"key": "value",
		},
		Description: "optional config item of map (flat structure) type",
		EnvName:     "CONFIG_map_key", // value must be JSON
		// cannot set a structured key with command line flags
		Validate: auconfigapi.ConfigNeedsNoValidation,
	},
}

func TestLoadAllowMissingPaths_Success_MostlyDefaults(t *testing.T) {
	ResetForTesting()

	expectNoFailure, verifyFailures := tstExpectFailFunc(t, "")
	expectNoWarning, verifyWarnings := tstExpectWarnFunc(t)
	Setup(tstConfigItems, expectNoFailure, expectNoWarning)

	require.NoError(t, pflag.CommandLine.Set("string-key-no-default", "flag-value"))
	t.Setenv("CONFIG_secret_key", "secret-value")

	LoadAllowMissingPaths()

	require.Equal(t, "flag-value", viper.GetString("string-key-no-default"))
	require.Equal(t, "default", viper.GetString("string-key-with-default"))
	require.Equal(t, "secret-value", viper.GetString("secret-key"))
	require.Equal(t, 24, viper.GetInt("int-key"))
	require.Equal(t, 3, viper.GetInt("int8-key"))
	require.Equal(t, 30000, viper.GetInt("int16-key"))
	require.Equal(t, int32(2000000000), viper.GetInt32("int32-key"))
	require.Equal(t, int64(8000000000), viper.GetInt64("int64-key"))
	require.Equal(t, uint(1234), viper.GetUint("uint-key"))
	require.Equal(t, uint(12), viper.GetUint("uint8-key"))
	require.Equal(t, uint16(65535), viper.GetUint16("uint16-key"))
	require.Equal(t, uint32(65535000), viper.GetUint32("uint32-key"))
	require.Equal(t, uint64(65535000000), viper.GetUint64("uint64-key"))
	require.Equal(t, false, viper.GetBool("bool-key"))
	require.EqualValues(t, []string{}, viper.GetStringSlice("string-slice-key"))
	require.EqualValues(t, map[string]string{"key": "value"}, viper.GetStringMapString("map-key"))
	verifyFailures()
	verifyWarnings()
}

func TestLoadAllowMissingPaths_Success_FlagValues(t *testing.T) {
	ResetForTesting()

	expectNoFailure, verifyFailures := tstExpectFailFunc(t, "")
	expectNoWarning, verifyWarnings := tstExpectWarnFunc(t)
	Setup(tstConfigItems, expectNoFailure, expectNoWarning)

	// in a real command line, you need to prefix the flag names with "-" or "--", such as
	//
	// -flag value (classic style)
	// --flag=value (GNU style)
	//
	require.NoError(t, pflag.CommandLine.Set("string-key-no-default", "flag-value"))
	require.NoError(t, pflag.CommandLine.Set("custom-string-key-with-default", "flag-value-2"))
	require.NoError(t, pflag.CommandLine.Set("int-key", "18"))
	require.NoError(t, pflag.CommandLine.Set("int8-key", "13"))
	require.NoError(t, pflag.CommandLine.Set("int16-key", "4444"))
	require.NoError(t, pflag.CommandLine.Set("int32-key", "44444444"))
	require.NoError(t, pflag.CommandLine.Set("int64-key", "888888888888"))
	require.NoError(t, pflag.CommandLine.Set("uint-key", "19900"))
	require.NoError(t, pflag.CommandLine.Set("uint8-key", "199"))
	require.NoError(t, pflag.CommandLine.Set("uint16-key", "55555"))
	require.NoError(t, pflag.CommandLine.Set("uint32-key", "55555555"))
	require.NoError(t, pflag.CommandLine.Set("uint64-key", "555555555555"))
	require.NoError(t, pflag.CommandLine.Set("bool-key", "true"))

	// multiple values for a slice parameter are set by passing command line flag multiple times
	require.NoError(t, pflag.CommandLine.Set("string-slice-key", "value1"))
	require.NoError(t, pflag.CommandLine.Set("string-slice-key", "value2"))

	// cannot set a map or structured key with command line flags, so leaving at default value

	// secrets should not be set by flags, or they will be visible in command line, so we use environment variable
	t.Setenv("CONFIG_secret_key", "secret-value")

	LoadAllowMissingPaths()

	require.Equal(t, "flag-value", viper.GetString("string-key-no-default"))
	require.Equal(t, "flag-value-2", viper.GetString("string-key-with-default"))
	require.Equal(t, 18, viper.GetInt("int-key"))
	require.Equal(t, 13, viper.GetInt("int8-key"))
	require.Equal(t, 4444, viper.GetInt("int16-key"))
	require.Equal(t, int32(44444444), viper.GetInt32("int32-key"))
	require.Equal(t, int64(888888888888), viper.GetInt64("int64-key"))
	require.Equal(t, uint(19900), viper.GetUint("uint-key"))
	require.Equal(t, uint(199), viper.GetUint("uint8-key"))
	require.Equal(t, uint16(55555), viper.GetUint16("uint16-key"))
	require.Equal(t, uint32(55555555), viper.GetUint32("uint32-key"))
	require.Equal(t, uint64(555555555555), viper.GetUint64("uint64-key"))
	require.Equal(t, true, viper.GetBool("bool-key"))
	require.EqualValues(t, []string{"value1", "value2"}, viper.GetStringSlice("string-slice-key"))
	require.EqualValues(t, map[string]string{"key": "value"}, viper.GetStringMapString("map-key"))
	verifyFailures()
	verifyWarnings()
}

func TestLoadAllowMissingPaths_ValidationFailure(t *testing.T) {
	ResetForTesting()

	// first validation failure strikes
	expectedFailure, verifyFailures := tstExpectFailFunc(t, "is empty")
	expectNoWarning, verifyWarnings := tstExpectWarnFunc(t)
	Setup(tstConfigItems, expectedFailure, expectNoWarning)

	LoadAllowMissingPaths()
	verifyFailures()
	verifyWarnings()
}

func TestLoad_Success(t *testing.T) {
	ResetForTesting()

	expectNoFailure, verifyFailures := tstExpectFailFunc(t, "")
	expectNoWarning, verifyWarnings := tstExpectWarnFunc(t)
	SetupWithOverriddenConfigPath(tstConfigItems, expectNoFailure, expectNoWarning,
		"test/resources/valid", "test/resources/valid")

	Load()

	// some of these values come from test/resources/config.yaml and test/resources/secrets.yaml
	require.Equal(t, "hello world", viper.GetString("string-key-no-default"))
	require.Equal(t, "default", viper.GetString("string-key-with-default"))
	require.Equal(t, "some secret value", viper.GetString("secret-key"))
	require.Equal(t, 24, viper.GetInt("int-key"))
	require.Equal(t, 3, viper.GetInt("int8-key"))
	require.Equal(t, 30000, viper.GetInt("int16-key"))
	require.Equal(t, int32(2000000000), viper.GetInt32("int32-key"))
	require.Equal(t, int64(8000000000), viper.GetInt64("int64-key"))
	require.Equal(t, uint(1234), viper.GetUint("uint-key"))
	require.Equal(t, uint(12), viper.GetUint("uint8-key"))
	require.Equal(t, uint16(65535), viper.GetUint16("uint16-key"))
	require.Equal(t, uint32(65535000), viper.GetUint32("uint32-key"))
	require.Equal(t, uint64(65535000000), viper.GetUint64("uint64-key"))
	require.Equal(t, false, viper.GetBool("bool-key"))
	require.EqualValues(t, []string{"a", "b", "c", "yes", "no"}, viper.GetStringSlice("string-slice-key"))
	require.EqualValues(t, map[string]string{"lion": "roar", "cat": "meow"}, viper.GetStringMapString("map-key"))
	verifyFailures()
	verifyWarnings()
}

func TestLoad_Default_Filenames_Missing(t *testing.T) {
	ResetForTesting()

	expectNoFailure, verifyFailures := tstExpectFailFunc(t,
		"Fatal error: configuration file config.(yaml|json) not found in .")
	expectWarnings, verifyWarnings := tstExpectWarnFunc(t,
		"you did not provide the config-path command line flag",
		"you did not provide the secrets-path command line flag",
	)
	SetupWithOverriddenConfigPath(tstConfigItems, expectNoFailure, expectWarnings, "", "")

	Load()

	verifyFailures()
	verifyWarnings()
}

func TestLoad_Invalid_Config_File(t *testing.T) {
	ResetForTesting()

	expectNoFailure, verifyFailures := tstExpectFailFunc(t,
		"Fatal error: configuration file config.(yaml|json) found but failed to load")
	expectWarnings, verifyWarnings := tstExpectWarnFunc(t)
	SetupWithOverriddenConfigPath(tstConfigItems, expectNoFailure, expectWarnings,
		"test/resources/invalid", "test/resources/invalid")

	Load()

	verifyFailures()
	verifyWarnings()
}

func TestSetupDefaultsOnly(t *testing.T) {
	ResetForTesting()

	expectNoFailure, verifyFailures := tstExpectFailFunc(t, "")
	expectNoWarning, verifyWarnings := tstExpectWarnFunc(t)
	SetupDefaultsOnly(tstConfigItems, expectNoFailure, expectNoWarning)

	require.Equal(t, "", viper.GetString("string-key-no-default"))
	require.Equal(t, "default", viper.GetString("string-key-with-default"))
	require.Equal(t, "", viper.GetString("secret-key"))
	require.Equal(t, 24, viper.GetInt("int-key"))
	require.Equal(t, 3, viper.GetInt("int8-key"))
	require.Equal(t, 30000, viper.GetInt("int16-key"))
	require.Equal(t, int32(2000000000), viper.GetInt32("int32-key"))
	require.Equal(t, int64(8000000000), viper.GetInt64("int64-key"))
	require.Equal(t, uint(1234), viper.GetUint("uint-key"))
	require.Equal(t, uint(12), viper.GetUint("uint8-key"))
	require.Equal(t, uint16(65535), viper.GetUint16("uint16-key"))
	require.Equal(t, uint32(65535000), viper.GetUint32("uint32-key"))
	require.Equal(t, uint64(65535000000), viper.GetUint64("uint64-key"))
	require.Equal(t, false, viper.GetBool("bool-key"))
	require.EqualValues(t, []string{}, viper.GetStringSlice("string-slice-key"))
	require.EqualValues(t, map[string]string{"key": "value"}, viper.GetStringMapString("map-key"))
	verifyFailures()
	verifyWarnings()
}

// --- helpers ---

type tstVerifyFunc func()

func tstExpectFailFunc(t *testing.T, expectContains string) (auconfigapi.ConfigFailFunc, tstVerifyFunc) {
	t.Helper()

	shouldBeCalled := expectContains != ""
	called := false

	return func(err error) {
			if called {
				// normally, failFunc would exit, so we would never have got here - ignore subsequent calls
				return
			}
			if expectContains == "" {
				t.Error("unexpected call to error function")
			} else {
				called = true
				require.ErrorContains(t, err, expectContains)
			}
		}, func() {
			if called != shouldBeCalled {
				t.Error("call expectation for failFunc not met")
			}
		}
}

func tstExpectWarnFunc(t *testing.T, expectContains ...string) (auconfigapi.ConfigWarnFunc, tstVerifyFunc) {
	t.Helper()

	position := -1

	return func(message string) {
			position++
			require.Less(t, position, len(expectContains), "unexpected additional call to warnFunc")
			require.Contains(t, message, expectContains[position], fmt.Sprintf("expected message to contain %q in call %d to warnFunc", expectContains[position], position+1))
		}, func() {
			require.Equal(t, len(expectContains), position+1, "unexpected number of calls to warnFunc")
		}
}
