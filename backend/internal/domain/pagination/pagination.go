package pagination

import "github.com/kataras/iris/v12"

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// Page holds normalized limit/offset values taken from the request query.
type Page struct {
	Limit  int
	Offset int
}

// FromContext reads "limit" and "offset" query params and clamps them to sane
// bounds, so a client can never ask for an unbounded result set.
func FromContext(ctx iris.Context) Page {
	limit := ctx.URLParamIntDefault("limit", DefaultLimit)
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset := ctx.URLParamIntDefault("offset", 0)
	if offset < 0 {
		offset = 0
	}

	return Page{Limit: limit, Offset: offset}
}
