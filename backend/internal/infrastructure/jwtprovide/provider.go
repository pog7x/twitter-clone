package jwtprovide

import (
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
		Signer:   jwt.NewSigner(jwt.HS256, []byte(cfg.JWTSecretKey), cfg.SessionExpiredAt),
		Verifier: jwt.NewVerifier(jwt.HS256, []byte(cfg.JWTSecretKey)),
	}, nil
}
