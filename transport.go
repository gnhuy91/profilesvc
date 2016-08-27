package profilesvc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(ctx context.Context, s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	r.Methods("POST").Path("/profile").Handler(kithttp.NewServer(
		ctx,
		MakePostProfileEndpoint(s),
		decodePostProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/profile/{id}").Handler(kithttp.NewServer(
		ctx,
		MakeGetProfileEndpoint(s),
		decodeGetProfileRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/profile/{id}").Handler(kithttp.NewServer(
		ctx,
		MakeDeleteProfileEndpoint(s),
		decodeDeleteProfileRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request postProfileRequest
	var buf bytes.Buffer

	// capture r.Body contents to buf
	body := io.TeeReader(r.Body, &buf)
	err := json.NewDecoder(body).Decode(&request.Profile)
	if err != nil {
		// inspect buf...
		return nil, err
	}
	return request, nil
}

func decodeGetProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getProfileRequest{ID: id}, nil
}

func decodeDeleteProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteProfileRequest{ID: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a go-kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
