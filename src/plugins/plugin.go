package plugins

import (
	"github.com/PiterWeb/RemoteController/src/plugins/game_share"
	"github.com/PiterWeb/RemoteController/src/plugins/shared"
)

func GetPlugins() []shared.Plugin {
	return []shared.Plugin{
		game_share.Plugin(),
	}
}
