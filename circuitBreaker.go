package circuitbreaker

import (
	cbConfig "github.com/cbluebird/circuitbreaker/config"
)

var CB CircuitBreaker

type CircuitBreaker struct {
	LB       LoadBalance
	SnapShot map[string]*ApiSnapShot
}

func init() {
	lb := LoadBalance{
		zfLB:    &randomLB{},
		oauthLB: &randomLB{},
	}
	snapShot := make(map[string]*ApiSnapShot)

	for _, api := range cbConfig.GetLoadBalanceConfig().Apis {
		lb.Add(api, Oauth)
		lb.Add(api, ZF)
		snapShot[api+string(Oauth)] = NewApiSnapShot()
		snapShot[api+string(ZF)] = NewApiSnapShot()
	}

	CB = CircuitBreaker{
		LB:       lb,
		SnapShot: snapShot,
	}
}

func (c *CircuitBreaker) GetApi(zfFlag, oauthFlag bool) (string, LoginType) {
	return c.LB.Pick(zfFlag, oauthFlag)
}

func (c *CircuitBreaker) Fail(api string, loginType LoginType) {
	if c.SnapShot[api+string(loginType)].Fail() {
		c.LB.Remove(api, loginType)
		Probe.Add(api, loginType)
	}
}

func (c *CircuitBreaker) Success(api string, loginType LoginType) {
	c.SnapShot[api+string(loginType)].Success()
}
