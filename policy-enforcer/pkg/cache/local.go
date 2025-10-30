package cache

import (
	"context"
	"sync"
	"time"
)

// Local implements Cache and is only intended for local tests.
type Local struct {
	sync.RWMutex

	data map[string][]byte
}

func NewLocal() *Local {
	return &Local{
		data: map[string][]byte{},
	}
}

func (l *Local) Get(_ context.Context, key string, val any) (bool, error) {
	l.RLock()
	defer l.RUnlock()

	item, ok := l.data[key]
	if !ok {
		return false, nil
	}

	return true, unmarshal(item, val)
}

func (l *Local) Set(_ context.Context, key string, val any, _ time.Duration) error {
	l.Lock()
	defer l.Unlock()

	item, err := marshal(val)
	if err != nil {
		return err
	}

	l.data[key] = item

	return nil
}

func (l *Local) Del(_ context.Context, key string) error {
	l.Lock()
	defer l.Unlock()

	delete(l.data, key)

	return nil
}
