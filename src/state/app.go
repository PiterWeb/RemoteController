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
	Dev bool
	Conn *Conn
}

type Conn struct {
	Ctx context.Context
	CancelCtx context.CancelFunc
	Peer *Peer
}

type Peer struct {
	open bool
	Mutex *sync.RWMutex
}

func (p *Peer) IsOpen() bool {
	p.Mutex.RLock()
	defer p.Mutex.RUnlock()
	return p.open
}

func (p *Peer) Open() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.open = true
}

func (p *Peer) Close() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.open = false
}