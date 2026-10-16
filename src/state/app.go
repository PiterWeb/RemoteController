package state

import (
	"context"
	"io/fs"
	"sync"
)

type App struct {
	Ctx context.Context
	CancelCtx context.CancelFunc
	Assets fs.FS
	Peer *Peer
	Dev bool
}

type Peer struct {
	open bool
	mutex *sync.RWMutex
}

func (p *Peer) IsOpen() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.open
}

func (p *Peer) Open() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.open = true
}

func (p *Peer) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.open = false
}