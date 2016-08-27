package profilesvc

import (
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// errorer is implemented by all concrete response types that may contain
// errors.
type errorer interface {
	error() error
}

// encodeError encodes business-logic errors and return them as HTTP errors.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError called with nil error (programmer error)")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// decide status code
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// codeFrom decides which HTTP status code should be returned.
func codeFrom(err error) int {
	switch err {
	// business-logic errors
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidRequestBody:
		return http.StatusBadRequest

	// go-kit predefined errors
	default:
		if e, ok := err.(kithttp.Error); ok {
			switch e.Domain {
			case kithttp.DomainDecode:
				return http.StatusBadRequest
			case kithttp.DomainDo:
				return http.StatusServiceUnavailable
			default:
				return http.StatusInternalServerError
			}
		}
		return http.StatusInternalServerError
	}
}
