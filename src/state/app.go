package state

import (
	"context"
	"embed"
)

type App struct {
	Ctx context.Context
	CancelCtx context.CancelFunc
	Assets embed.FS
}