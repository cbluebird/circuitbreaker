package circuitbreaker

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	cbConfig "github.com/cbluebird/circuitbreaker/config"
)

var Probe *LiveNessProbe

func init() {
	Probe = NewLiveNessProbe(cbConfig.GetLiveNessConfig())
}

type LiveNessProbe struct {
	sync.Mutex
	ApiMap   map[string]LoginType
	Duration time.Duration
	User     *User
}

func NewLiveNessProbe(config cbConfig.LiveNessProbeConfig) *LiveNessProbe {
	user := &User{
		StudentID:     config.StudentId,
		OauthPassword: config.OauthPassword,
		ZFPassword:    config.ZFPassword,
	}
	return &LiveNessProbe{
		ApiMap:   make(map[string]LoginType),
		Duration: config.Duration,
		User:     user,
	}
}

func (l *LiveNessProbe) Add(api string, loginType LoginType) {
	l.Lock()
	defer l.Unlock()
	l.ApiMap[api+string(loginType)] = loginType
}

func (l *LiveNessProbe) Remove(key string) {
	l.Lock()
	defer l.Unlock()
	delete(l.ApiMap, key)
}

func (l *LiveNessProbe) Start(ctx context.Context) {
	ticker := time.NewTicker(l.Duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for apiKey, loginType := range l.ApiMap {
				api := strings.TrimSuffix(apiKey, string(loginType))
				if err := liveNess(l.User, api, loginType); err == nil {
					CB.LB.Add(api, loginType)
					CB.Success(api, loginType)
					l.Remove(apiKey)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func liveNess(u *User, api string, loginType LoginType) error {
	var password string
	if loginType == Oauth {
		password = u.OauthPassword
	} else {
		password = u.ZFPassword
	}
	form := url.Values{}
	form.Add("username", u.StudentID)
	form.Add("password", password)
	form.Add("type", string(loginType))
	form.Add("year", strconv.Itoa(time.Now().Year()-1))
	form.Add("term", "上")

	f := Fetch{}
	f.Init()

	rc := struct {
		Code int `json:"code" binding:"required"`
	}{}
	for i := 0; i < 5; i++ {
		res, err := f.PostForm(api+string(ZFClassTable), form)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(res, &rc); err != nil {
			return err
		}
		if rc.Code != 413 {
			break
		}
	}
	if rc.Code == 200 || rc.Code == 412 || rc.Code == 416 {
		return nil
	}
	return errors.New("liveNessProbe failed")
}
