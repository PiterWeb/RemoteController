package plugins

import (
	"encoding/json"
	"os"
	"strings"
)

func getPluginSpecification(path string) plugins_args {

	// Load the .json file with the same name as the plugin
	// Return the init_client, init_server and background arguments

	path = strings.Replace(path, ".dll", ".json", 1)
	path = strings.Replace(path, ".so", ".json", 1)

	specsFile, err := os.ReadFile(path)

	specs := plugins_args{}

	if err != nil {
		return specs
	}

	json.Unmarshal(specsFile, &specs)

	return specs

}
