package jwtprovide

import (
	"time"

	"twitter-clone/config"

	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/samber/do"
)

type JWTProvider struct {
	Signer   *jwt.Signer
	Verifier *jwt.Verifier
}

func NewJWTProvider(i *do.Injector) (*JWTProvider, error) {
	cfg := do.MustInvoke[*config.Config](i)

	return &JWTProvider{
		Signer:   jwt.NewSigner(jwt.HS256, []byte(cfg.JWTSecretKey), 24*time.Hour),
		Verifier: jwt.NewVerifier(jwt.HS256, []byte(cfg.JWTSecretKey)),
	}, nil
}
