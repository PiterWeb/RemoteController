package plugins

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

type Plugin struct {
	Name             string                            `json:"name"`
	Path             string                            `json:"path"`
	init_client      func(...uintptr) (uintptr, error) `json:"-"` // ignore this field when marshalling
	Init_client_args []Plugin_arg                      `json:"init_client_args"`
	init_server      func(...uintptr) (uintptr, error) `json:"-"` // ignore this field when marshalling
	Init_host_args   []Plugin_arg                      `json:"init_host_args"`
	background       func(...uintptr) (uintptr, error) `json:"-"` // ignore this field when marshalling
	Background_args  []Plugin_arg                      `json:"background_args"`
	Enabled          bool                              `json:"enabled"`
}

var pluginsLock = &sync.Mutex{}
var pluginsInstance []Plugin

func GetPlugins() []Plugin {
	if pluginsInstance != nil {
		pluginsLock.Lock()
		defer pluginsLock.Unlock()
		if pluginsInstance == nil {
			pluginsInstance = loadPlugins()
		}
	}

	return pluginsInstance
}

// ReloadPlugins will reload the plugins from the json file
func ReloadPlugins() {
	pluginsLock.Lock()
	defer pluginsLock.Unlock()
	pluginsInstance = loadPlugins()
}

// ModifyArgs will modify the arguments of a plugin
// mode can be "init_client", "init_host" or "background"
func ModifyArgs(pluginName string, args []Plugin_arg, mode string) {
	plugins := GetPlugins()

	mode = strings.TrimSpace(mode)

	for _, plugin := range plugins {
		if plugin.Name == pluginName {

			if mode == "init_client" {
				plugin.Init_client_args = args
			} else if mode == "init_host" {
				plugin.Init_host_args = args
			} else if mode == "background" {
				plugin.Background_args = args
			}

			plugin.PersistPlugin()
			break

		}
	}
}

func (p Plugin) PersistPlugin() {

	args := plugins_args{
		Init_client: p.Init_client_args,
		Init_host:   p.Init_host_args,
		Background:  p.Background_args,
	}

	argsData, err := json.Marshal(args)

	if err != nil {
		return
	}

	// Write the arguments to the json file
	os.WriteFile(p.Path, argsData, 0644)

}

func (p *Plugin) Toogle() {
	p.Enabled = !p.Enabled
}

func (p Plugin) IsEnabled() bool {
	return p.Enabled
}

// Init client will get the input arguments from the struct populated by the json file and will return the result of the function
func (p Plugin) Init_client(comms_port uint16) (uintptr, error) {

	if !p.IsEnabled() {
		return 0, nil
	}

	// Reload the arguments from the json file
	p.ReloadArgs()

	args := []uintptr{}

	args = append(args, uintptr(comms_port))

	for _, arg := range p.Init_client_args {
		args = append(args, uintptr(arg.Value.(uintptr)))
	}

	return p.init_client(args...)
}

// Init server will get the input arguments from the struct populated by the json file and will return the result of the function
func (p Plugin) Init_host(comms_port uint16) (uintptr, error) {

	if !p.IsEnabled() {
		return 0, nil
	}
	// Reload the arguments from the json file
	p.ReloadArgs()

	args := []uintptr{}

	args = append(args, uintptr(comms_port))

	for _, arg := range p.Init_host_args {
		args = append(args, uintptr(arg.Value.(uintptr)))
	}

	return p.init_server(args...)
}

// Background will get the input arguments from the struct populated by the json file and will return the result of the function
func (p Plugin) Background(comms_port uint16) (uintptr, error) {

	if !p.IsEnabled() {
		return 0, nil
	}
	// Reload the arguments from the json file
	p.ReloadArgs()

	args := []uintptr{}

	args = append(args, uintptr(comms_port))

	for _, arg := range p.Background_args {
		args = append(args, uintptr(arg.Value.(uintptr)))
	}

	return p.background(args...)

}

func (p *Plugin) ReloadArgs() {

	spec := getPluginSpecification(p.Path)

	p.Init_client_args = spec.Init_client
	p.Init_host_args = spec.Init_host
	p.Background_args = spec.Background

}

type plugins_args struct {
	Init_client []Plugin_arg `json:"init_client_args"`
	Init_host   []Plugin_arg `json:"init_host_args"`
	Background  []Plugin_arg `json:"background_args"`
}

type Plugin_arg struct {
	Name      string       `json:"name"`
	Value     any          `json:"value"`
	ValueList []Plugin_arg `json:"value_list,omitempty"`
}
