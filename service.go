// This file defines our application/package logics for implementations.

package profilesvc

import "golang.org/x/net/context"

type Service interface {
	PostProfile(ctx context.Context, p Profile) error
	GetProfile(ctx context.Context, id string) (Profile, error)
	DeleteProfile(ctx context.Context, id string) error
}

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}
