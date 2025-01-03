package circuitbreaker

import (
	"time"
)

type LiveNessProbeConfig struct {
	StudentId     string
	OauthPassword string
	ZFPassword    string
	Duration      time.Duration
}

func GetConfig() LiveNessProbeConfig {
	cfg := LiveNessProbeConfig{}
	return cfg
}
