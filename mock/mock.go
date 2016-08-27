// This package provides mock implementations of profilesvc.Service for testing.

package mock

import (
	"golang.org/x/net/context"

	"github.com/gnhuy91/profilesvc"
)

// Service is a mock implementation of profilesvc.Service.
type Service struct {
	PostProfileFunc   func(ctx context.Context, p profilesvc.Profile) error
	GetProfileFunc    func(ctx context.Context, id string) (profilesvc.Profile, error)
	DeleteProfileFunc func(ctx context.Context, id string) error
}

// PostProfile invokes the mock implementation.
func (s *Service) PostProfile(ctx context.Context, p profilesvc.Profile) error {
	return s.PostProfileFunc(ctx, p)
}

// GetProfile invokes the mock implementation.
func (s *Service) GetProfile(ctx context.Context, id string) (profilesvc.Profile, error) {
	return s.GetProfileFunc(ctx, id)
}

// DeleteProfile invokes the mock implementation.
func (s *Service) DeleteProfile(ctx context.Context, id string) error {
	return s.DeleteProfileFunc(ctx, id)
}
