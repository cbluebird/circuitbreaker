package circuitbreaker

import (
	"sync"
	"time"
)

var Probe LiveNessProbe

func init() {
	Probe = LiveNessProbe{
		Config: GetLiveNessConfig(),
		ApiMap: make(map[string]LoginType),
	}
	go Probe.Start()
}

type LiveNessProbe struct {
	sync.Mutex
	Config LiveNessProbeConfig
	ApiMap map[string]LoginType
}

func (l *LiveNessProbe) Add(api string, loginType LoginType) {
	l.Lock()
	defer l.Unlock()
	l.ApiMap[api] = loginType
}

func (l *LiveNessProbe) Remove(key string) {
	l.Lock()
	defer l.Unlock()
	delete(l.ApiMap, key)
}

func (l *LiveNessProbe) Start() {
	ticker := time.NewTicker(l.Config.Duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for api, loginType := range l.ApiMap {
				// do some healthy check
				_ = api
				_ = loginType
			}
		}
	}
}
