package authcookie

import (
	"encoding/base64"
	"fmt"

	"twitter-clone/config"

	"github.com/gorilla/securecookie"
	"github.com/samber/do"
)

func NewCookieCodec(i *do.Injector) (*securecookie.SecureCookie, error) {
	cfg := do.MustInvoke[*config.Config](i)

	hashKey, err := base64.StdEncoding.DecodeString(cfg.EncodedSessionHashKey)
	if err != nil {
		return nil, fmt.Errorf("decode auth session block hash error %v", err)
	}

	blockKey, err := base64.StdEncoding.DecodeString(cfg.EncodedSessionBlockKey)
	if err != nil {
		return nil, fmt.Errorf("decode auth session block key error %v", err)
	}

	return securecookie.New(hashKey, blockKey), nil
}
