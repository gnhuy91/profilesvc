package profilesvc

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

func MakePostProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postProfileRequest)
		e := s.PostProfile(ctx, req.Profile)
		if e != nil {
			return postProfileResponse{Err: e.Error()}, nil
		}
		return postProfileResponse{Err: ""}, nil
	}
}

func MakeGetProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProfileRequest)
		p, e := s.GetProfile(ctx, req.ID)
		if e != nil {
			return getProfileResponse{Profile: p, Err: e.Error()}, nil
		}
		return getProfileResponse{Profile: p, Err: ""}, nil
	}
}

func MakeDeleteProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteProfileRequest)
		e := s.DeleteProfile(ctx, req.ID)
		if e != nil {
			return deleteProfileResponse{Err: e.Error()}, nil
		}
		return deleteProfileResponse{Err: ""}, nil
	}
}

type postProfileRequest struct {
	Profile Profile
}

type postProfileResponse struct {
	Err string `json:"err,omitempty"`
}

type getProfileRequest struct {
	ID string
}

type getProfileResponse struct {
	Profile Profile `json:"profile,omitempty"`
	Err     string  `json:"err,omitempty"`
}

type deleteProfileRequest struct {
	ID string
}

type deleteProfileResponse struct {
	Err string `json:"err,omitempty"`
}
