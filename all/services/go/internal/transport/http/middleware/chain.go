package middleware

import "net/http"

// Chain applies middleware in the order provided:
// Chain(h, A, B) => A(B(h))
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
