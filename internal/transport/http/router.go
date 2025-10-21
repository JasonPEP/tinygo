package http

import (
	stdhttp "net/http"
)

type Middleware func(stdhttp.Handler) stdhttp.Handler

// Chain composes middlewares left-to-right.
func Chain(h stdhttp.Handler, mws ...Middleware) stdhttp.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
