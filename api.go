package circuitbreaker

type LoginType string

const (
	Oauth   LoginType = "OAUTH"
	ZF      LoginType = "ZF"
	Unknown LoginType = "Unknown"
)

const (
	ZFClassTable string = "student/zf/table"
)
