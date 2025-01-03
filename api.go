package circuitbreaker

type LoginType string

const (
	Oauth   LoginType = "OAUTH"
	ZF                = "ZF"
	Unknown           = "Unknown"
)
