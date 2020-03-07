package auconfig

import auconfigapi "github.com/StephanHCB/go-autumn-config-api"

// config item for profiles supported by this package, use in your list of configItems that you pass to Setup()
var ConfigItemProfile = auconfigapi.ConfigItem{
	Key:         "profiles",
	Default:     []string{},
	Description: "list of profiles, separate by spaces in environment or command line parameters",
	Validate:    auconfigapi.ConfigNeedsNoValidation,
}
