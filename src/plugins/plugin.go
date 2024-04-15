package plugins

type plugin struct {
	name             string
	path             string
	init_client      func(...uintptr) (uintptr, error)
	init_client_args []plugin_arg
	init_server      func(...uintptr) (uintptr, error)
	init_server_args []plugin_arg
	background       func(...uintptr) (uintptr, error)
	background_args  []plugin_arg
	enabled          bool
}

func (p *plugin) Toogle() {
	p.enabled = !p.enabled
}

func (p plugin) IsEnabled() bool {
	return p.enabled
}

// Init client will get the input arguments from the struct populated by the json file
func (p plugin) Init_client() (uintptr, error) {

	// Reload the arguments from the json file
	p.ReloadArgs()

	args := []uintptr{}

	for _, arg := range p.init_client_args {
		args = append(args, uintptr(arg.Value.(uintptr)))
	}

	return p.init_client(args...)
}

// Init server will get the input arguments from the struct populated by the json file
func (p plugin) Init_server() (uintptr, error) {

	// Reload the arguments from the json file
	p.ReloadArgs()

	args := []uintptr{}

	for _, arg := range p.init_server_args {
		args = append(args, uintptr(arg.Value.(uintptr)))
	}

	return p.init_server(args...)
}

// Background will get the input arguments from the struct populated by the json file
func (p plugin) Background() (uintptr, error) {

	// Reload the arguments from the json file
	p.ReloadArgs()

	args := []uintptr{}

	for _, arg := range p.background_args {
		args = append(args, uintptr(arg.Value.(uintptr)))
	}

	return p.background(args...)

}

func (p *plugin) ReloadArgs() {

	spec := getPluginSpecification(p.path)

	p.init_client_args = spec.Init_client
	p.init_server_args = spec.Init_server
	p.background_args = spec.Background

}

type plugins_args struct {
	Init_client []plugin_arg `json:"init_client_args"`
	Init_server []plugin_arg `json:"init_server_args"`
	Background  []plugin_arg `json:"background_args"`
}

type plugin_arg struct {
	Name      string       `json:"name"`
	Value     any          `json:"value"`
	ValueList []plugin_arg `json:"value_list,omitempty"`
}
