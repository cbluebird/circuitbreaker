package circuitbreaker

import (
	cbConfig "github.com/cbluebird/circuitbreaker/config"
	"sync"
)

var CB CircuitBreaker

type CircuitBreaker struct {
	LB       LoadBalance
	SnapShot sync.Map
}

func init() {
	lb := LoadBalance{
		zfLB:    &randomLB{},
		oauthLB: &randomLB{},
	}
	snapShot := sync.Map{}

	for _, api := range cbConfig.GetLoadBalanceConfig().Apis {
		lb.Add(api, Oauth)
		lb.Add(api, ZF)
		snapShot.Store(api+string(Oauth), NewApiSnapShot())
		snapShot.Store(api+string(ZF), NewApiSnapShot())
	}

	CB = CircuitBreaker{
		LB:       lb,
		SnapShot: snapShot,
	}
}

func (c *CircuitBreaker) GetApi(zfFlag, oauthFlag bool) (string, LoginType, error) {
	return c.LB.Pick(zfFlag, oauthFlag)
}

func (c *CircuitBreaker) Fail(api string, loginType LoginType) {
	if value, ok := c.SnapShot.Load(api + string(loginType)); ok {
		if snapshot, ok := value.(*ApiSnapShot); ok {
			if snapshot.Fail() {
				c.LB.Remove(api, loginType)
				Probe.Add(api, loginType)
			}
		}
	}
}

func (c *CircuitBreaker) Success(api string, loginType LoginType) {
	if value, ok := c.SnapShot.Load(api + string(loginType)); ok {
		if snapshot, ok := value.(*ApiSnapShot); ok {
			snapshot.Success()
		}
	}
}
