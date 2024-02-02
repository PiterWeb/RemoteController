package plugins

import (
	"io"
	"net/http"
	"os"
)

func downloadPlugin(pluginName, url string) error {

	if _, err := os.ReadFile("./plugins/" + pluginName + file_extension); err == nil {
		return nil
	}

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create("./plugins/" + pluginName + file_extension)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err

}
