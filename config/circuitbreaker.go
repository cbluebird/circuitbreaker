package config

import (
	"time"
)

type LiveNessProbeConfig struct {
	StudentId     string
	OauthPassword string
	ZFPassword    string
	Duration      time.Duration
}

func GetLiveNessConfig() LiveNessProbeConfig {
	return LiveNessProbeConfig{
		StudentId:     Config.GetString("zfCircuit.studentId"),
		OauthPassword: Config.GetString("zfCircuit.oauthPassword"),
		ZFPassword:    Config.GetString("zfCircuit.zfPassword"),
		Duration:      Config.GetDuration("zfCircuit.duration"),
	}
}

type LoadBalanceConfig struct {
	Apis []string
}

func GetLoadBalanceConfig() LoadBalanceConfig {
	return LoadBalanceConfig{
		Apis: Config.GetStringSlice("zfCircuit.apis"),
	}
}
