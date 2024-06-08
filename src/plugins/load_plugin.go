package plugins

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/Binject/universal"
)

func loadPlugins() []Plugin {

	pluginPaths := getPluginPaths()

	plugins := []Plugin{}

	loader, err := universal.NewLoader()

	if err != nil {
		return []Plugin{}
	}

	for _, path := range pluginPaths {
		plugin := Plugin{}

		func_specs := getPluginSpecification(path)

		image, err := os.ReadFile(path)

		if err != nil {
			continue
		}

		lib, err := loader.LoadLibrary("main", &image)

		if err != nil {
			continue
		}

		plugin.Init_client_args = func_specs.Init_client
		plugin.Init_host_args = func_specs.Init_host
		plugin.Background_args = func_specs.Background

		plugin.init_client = func(u ...uintptr) (uintptr, error) {
			return lib.Call("init_client", u...)
		}

		plugin.init_server = func(u ...uintptr) (uintptr, error) {
			return lib.Call("init_server", u...)
		}

		plugin.background = func(u ...uintptr) (uintptr, error) {
			return lib.Call("background", u...)
		}

		plugin.Name = strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
		fmt.Println(plugin.Name)
		plugin.Path = path

		plugins = append(plugins, plugin)

	}

	return plugins

}

func getPluginPaths() []string {
	wd, err := os.Getwd()

	if err != nil {
		return []string{}
	}

	wd += "/plugins/"

	pluginFiles, err := os.ReadDir(wd)

	if err != nil {
		os.Mkdir(wd, fs.FileMode(os.O_CREATE))
		return []string{}
	}

	pluginPaths := []string{}

	for _, file := range pluginFiles {
		if file.IsDir() || !strings.Contains(file.Name(), ".dll") || !strings.Contains(file.Name(), ".so") {
			continue
		}

		pluginPaths = append(pluginPaths, wd+file.Name())
	}

	return pluginPaths
}
