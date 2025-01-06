package circuitbreaker

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type LiveNessProbeConfig struct {
	StudentId     string
	OauthPassword string
	ZFPassword    string
	Duration      time.Duration
}

func GetLiveNessConfig() LiveNessProbeConfig {
	cfg := LiveNessProbeConfig{}
	if Config.IsSet("liveNess.studentId") {
		cfg.StudentId = Config.GetString("liveNess.studentId")
	}
	if Config.IsSet("liveNess.oauthPassword") {
		cfg.OauthPassword = Config.GetString("liveNess.oauthPassword")
	}
	if Config.IsSet("liveNess.zfPassword") {
		cfg.ZFPassword = Config.GetString("liveNess.zfPassword")
	}
	if Config.IsSet("liveNess.duration") {
		cfg.Duration = Config.GetDuration("liveNess.duration")
	}
	return cfg
}

type LoadBalanceConfig struct {
	Name string
	Url  string
	Type LoginType
}

func GetLoadBalanceConfig() []LoadBalanceConfig {
	var configs []LoadBalanceConfig
	err := viper.UnmarshalKey("loadBalance", &configs)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return configs
}
