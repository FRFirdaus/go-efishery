package requestidmw

import (
	"net/http"

	"github.com/google/uuid"
)

// HttpMiddleware is inject request id and response id
func HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// set request id header if not exist
		if r.Header.Get(HeaderKey) == "" {
			r.Header.Set(HeaderKey, uuid.New().String())
		}

		// set request id header response if not exist
		if rw.Header().Get(HeaderKey) == "" {
			rw.Header().Set(HeaderKey, r.Header.Get(HeaderKey))
		}

		if next != nil {
			next.ServeHTTP(rw, r.WithContext(ctx))
		}
	})
}

// GetRequestId is get request id from request header
func GetRequestId(r *http.Request) string {
	if r == nil {
		return ""
	}
	return r.Header.Get(HeaderKey)
}
