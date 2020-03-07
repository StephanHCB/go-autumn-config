# go-autumn-config

## About go-autumn

A collection of libraries for [enterprise microservices](https://github.com/StephanHCB/go-mailer-service/blob/master/README.md) in golang that
- is heavily inspired by Spring Boot / Spring Cloud
- is very opinionated
- names modules by what they do
- unlike Spring Boot avoids certain types of auto-magical behaviour
- is not a library monolith, that is every part only depends on the api parts of the other components
  (if at all), and the api parts do not add any dependencies.

Fall is my favourite season, so I'm calling it _go-autumn_.

## About go-autumn-config

A library that handles configuration for enterprise microservices.

## How to use

We recommend collecting all configuration related code in a package `internal/repository/configuration`.

You configure the configuration subsystem by a call to `auconfig.Setup(...)`. This function takes 3 arguments:
 - a list of `auconfigapi.ConfigItem` to specify what configuration items exist 
 - a failure handler of type `auconfigapi.ConfigFailFunc`, which is expected to somehow notify
   the user or environment that loading configuration has failed. It should terminate the application. 
   `panic` satisfies these requirements, but we hope you won't use that in an actual production ready
   enterprise service...
 - a warning message handler of type `auconfigapi.ConfigWarnFunc`, which should probably log a warning
   using your preferred logging framework. `log.Print` satisfies the type requirements, but again we
   hope this is not what you'll use in production...

See [go-autumn-config-api](https://github.com/StephanHCB/go-autumn-config-api/api.go) for the precise
type definitions.

When you request your configuration to be loaded, which you must do yourself with a call to 
`auconfig.Load()`, every key is assigned its value by going through the following precedence list:
 - command line flag
 - environment variable
 - configuration read from secrets.(yaml|json|properties)
 - configuration read from each config-profileName.(yaml|json|properties) in reverse order, so the last profile wins 
   (NOT IMPLEMENTED YET)
 - configuration read from config.(yaml|json|properties)
 - default value specified in ConfigItems

**Important:** avoid calling `Setup(...)` or `Load()` from inside an `init()` func, or you might get errors if another
library defines any command line parameters using flag or pflag. Setup calls `pflag.Parse()`!

The whole library is really just a thin wrapper around [spf13/viper](https://github.com/spf13/viper) and
[spf13/pflag](https://github.com/spf13/pflag). Once loaded, you access your configuration values simply
using the viper primitives such as `viper.GetString()`. 

The library will always add two extra GNU flags style command line parameters called `config-path` and `secrets-path`. 
If unset, these default to the current working directory, but they will log a warning, because you should set those
in production. 

```main --config-path=. --secrets-path=.```

Will get rid of the warnings. These parameters contain the path to two configuration files which must be
called `config.(yaml|json|properties)` and `secrets.(yaml|json|properties)`.

We have found that a good use pattern is to have a file called `access.go` inside your configuration
package where you can provide public accessor functions for all your configuration values.

## Example:

Using this can be as simple as:

```go
package configuration

import (
	"fmt"
	"github.com/StephanHCB/go-autumn-config-api"
	"github.com/StephanHCB/go-autumn-config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

// custom validation function example
func checkValidPortNumber(key string) error {
	port := viper.GetUint(key)
	if port < 1024 || port > 65535 {
		return fmt.Errorf("Fatal error: configuration value for key %s is not in range 1024..65535\n", key)
	}
	return nil
}

// define what configuration items you want
var configItems = []auconfigapi.ConfigItem{
	{
		Key:         "server.address",
		Default:     "",
		Description: "ip address or hostname to listen on, can be left blank for localhost",
		Validate:    auconfigapi.ConfigNeedsNoValidation,
	}, {
		Key:         "server.port",
		Default:     uint(8080),
		Description: "port to listen on, defaults to 8080 if not set",
		Validate:    checkValidPortNumber,
	},
}

// initialize the library.
func Setup() {
	auconfig.Setup(configItems, panic, log.Print)
	auconfig.Load()
}

// provide accessor functions, using viper.GetXYZ to read configuration values

func ServerAddress() string {
	return fmt.Sprintf("%v:%d", viper.GetString("server.address"), viper.GetUint("server.port"))
}
```

## TODOs

- add unit tests
- push mini example with a ContextLoads acceptance test to separate repo and link it here
