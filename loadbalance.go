package circuitbreaker

import (
	"sync"

	"github.com/bytedance/gopkg/lang/fastrand"
)

type LoadBalanceType int

const (
	Random LoadBalanceType = iota
)

type LoadBalance struct {
	zfLB    *randomLB
	oauthLB *randomLB
}

func (lb *LoadBalance) Pick(zfFlag, oauthFlag bool) (string, LoginType) {
	var loginType LoginType

	if zfFlag {
		zfFlag = lb.zfLB.Len() > 0
	}
	if oauthFlag {
		oauthFlag = lb.oauthLB.Len() > 0
	}

	if oauthFlag && zfFlag {
		if fastrand.Intn(100) > 50 {
			loginType = Oauth
		} else {
			loginType = ZF
		}
	} else if oauthFlag {
		loginType = Oauth
	} else if zfFlag {
		loginType = ZF
	} else {
		return "", Unknown
	}

	if loginType == Oauth {
		return lb.oauthLB.Pick(), loginType
	}
	return lb.zfLB.Pick(), loginType
}

func (lb *LoadBalance) Add(api string, loginType LoginType) {
	if loginType == Oauth {
		lb.oauthLB.Add(api)
	} else {
		lb.zfLB.Add(api)
	}
}

func (lb *LoadBalance) Remove(api string, loginType LoginType) {
	if loginType == Oauth {
		lb.oauthLB.Remove(api)
	} else {
		lb.zfLB.Remove(api)
	}
}

type loadBalance interface {
	LoadBalance() LoadBalanceType
	Pick() (api string)
	ReBalance(apis []string)
	Add(api ...string)
	Remove(api string)
}

type randomLB struct {
	sync.Mutex
	Api  []string
	Size int
}

func newRandomLB(apis []string) loadBalance {
	return &randomLB{Api: apis, Size: len(apis)}
}

func (b *randomLB) LoadBalance() LoadBalanceType {
	return Random
}

func (b *randomLB) Pick() string {
	b.Lock()
	defer b.Unlock()
	idx := fastrand.Intn(b.Size)
	return b.Api[idx]
}

func (b *randomLB) ReBalance(apis []string) {
	b.Lock()
	defer b.Unlock()
	b.Api, b.Size = apis, len(apis)
}

func (b *randomLB) Add(api ...string) {
	b.Lock()
	defer b.Unlock()
	b.Api = append(b.Api, api...)
	b.Size = len(b.Api)
}

func (b *randomLB) Remove(api string) {
	b.Lock()
	defer b.Unlock()
	for i, s := range b.Api {
		if s == api {
			b.Api = append(b.Api[:i], b.Api[i+1:]...)
			break
		}
	}
	b.Size = len(b.Api)
}

func (b *randomLB) Len() int {
	b.Lock()
	defer b.Unlock()
	return b.Size
}
