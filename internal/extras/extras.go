package extras

import (
	"fmt"
	"sync"
	"time"
)

var tc *tokenCache

func ClearTokenCache() {
	tc = nil
}

func GetAuthToken() (string, error) {
	if tc == nil {
		tc = &tokenCache{}
	}
	token := tc.getToken()
	if token == "" {
		token, err := tc.setToken()
		if err != nil {
			return "", err
		}
		return token, nil
	}
	return token, nil
}

type tokenCache struct {
	token   string
	expires time.Time
	lock    sync.RWMutex
}

func (tc *tokenCache) getToken() string {
	tc.lock.RLock()
	defer tc.lock.RUnlock()
	now := time.Now()
	if tc.expires.Before(now) {
		return ""
	}
	return tc.token
}

func (tc *tokenCache) setToken() (string, error) {
	tc.lock.Lock()
	defer tc.lock.Unlock()
	token, expiresIn := fmt.Sprintf("token-%d", time.Now().Unix()), 2*time.Minute
	tc.token = token
	tc.expires = time.Now().Add(time.Duration(expiresIn/2) * time.Second)
	return token, nil
}
