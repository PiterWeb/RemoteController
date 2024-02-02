package plugins

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

const file_extension string = ".exe"

type PluginAction struct {
	ActionName string   `json:"actionName"`
	Parameters []string `json:"parameters"`
}

func (pa PluginAction) GetName() string {
	return pa.ActionName
}

func (pa PluginAction) GetParameters() []string {
	return pa.Parameters
}

type Plugin struct {
	PluginName  string         `json:"pluginName"`
	PluginUrl   string         `json:"pluginUrl"`
	DownloadUrl string         `json:"downloadUrl"`
	Actions     []PluginAction `json:"actions"`
	downloaded  bool
}

func (p Plugin) GetName() string {
	return p.PluginName
}

func (p Plugin) GetActions() []PluginAction {
	return p.Actions
}

func (p Plugin) GetPluginUrl() string {
	return p.PluginUrl
}

func (p Plugin) GetDownloadUrl() string {
	return p.DownloadUrl
}

func (p *Plugin) New(pluginUrl string) error {
	p.PluginUrl = pluginUrl
	return p.init()
}

func (p *Plugin) init() error {

	res, err := http.Get(p.PluginUrl)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, p)

	if err != nil {
		return err
	}

	err = downloadPlugin(p.PluginName, p.DownloadUrl)

	if err != nil {
		return err
	}

	p.downloaded = true

	return nil

}

func (p Plugin) Call(values ...string) error {

	cmd := exec.Command("./"+p.PluginName+file_extension, values...)

	output, err := cmd.Output()

	if err != nil {
		return err
	}

	strOutput := string(output)

	if strings.Contains(strOutput, "ERROR =") {
		return errors.New(strings.Split(strOutput, "ERROR =")[1])
	}

	return nil
}
