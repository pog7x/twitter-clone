package corsheaders

import (
	"strings"

	"github.com/kataras/iris/v12"

	"twitter-clone/config"
)

const allowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"

const allowedHeaders = "Content-Type, Authorization"

// defaultOrigins are used when CORS_ALLOWED_ORIGINS is not configured, so a
// local frontend keeps working out of the box.
var defaultOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000"}

// New builds the CORS middleware for the configured origins. The request
// origin is echoed back only when it is on the list, which keeps credentialed
// requests limited to known frontends.
func New(cfg *config.Config) iris.Handler {
	allowed := parseOrigins(cfg.CORSAllowedOrigins)

	return func(ctx iris.Context) {
		origin := ctx.GetHeader("Origin")

		// Vary is set unconditionally: the response body depends on the Origin
		// header even when the origin is rejected.
		ctx.Header("Vary", "Origin")

		if origin != "" && allowed[origin] {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}

		if ctx.Method() == iris.MethodOptions {
			ctx.Header("Access-Control-Allow-Methods", allowedMethods)
			ctx.Header("Access-Control-Allow-Headers", allowedHeaders)
			ctx.Header("Access-Control-Max-Age", "86400")

			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func parseOrigins(raw string) map[string]bool {
	origins := defaultOrigins

	if trimmed := strings.TrimSpace(raw); trimmed != "" {
		origins = strings.Split(trimmed, ",")
	}

	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = true
		}
	}

	return allowed
}
