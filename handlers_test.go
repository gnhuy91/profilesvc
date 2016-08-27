package profilesvc_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"github.com/gnhuy91/profilesvc"
	"github.com/gnhuy91/profilesvc/mock"
)

func TestPostProfileHandler(t *testing.T) {
	cases := []struct {
		path, method, body string
		expectedRespBody   string
		expectedStatusCode int
	}{
		{
			path:               "/profile",
			method:             "POST",
			body:               `{ "name": "Blah" }`,
			expectedRespBody:   `{}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			path:               "/profile",
			method:             "POST",
			body:               `{ "invalid": "field" }`,
			expectedRespBody:   `{"error":"invalid request body"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	svc := &mock.Service{}
	ctx := context.Background()
	h := profilesvc.MakeHTTPHandler(ctx, svc, kitlog.NewNopLogger())
	// mock func call
	svc.PostProfileFunc = func(_ context.Context, p profilesvc.Profile) error {
		if (p == profilesvc.Profile{}) {
			return profilesvc.ErrInvalidRequestBody
		}
		return nil
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			req, _ := http.NewRequest(c.method, c.path, strings.NewReader(c.body))
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != c.expectedStatusCode {
				t.Errorf("unexpected status code: %v, want %v", rec.Code, c.expectedStatusCode)
			}
			respBody := strings.TrimSpace(rec.Body.String())
			if respBody != c.expectedRespBody {
				t.Errorf("unexpected response body: %s, want %s", respBody, c.expectedRespBody)
			}
		})
	}
}

func TestGetProfileHandler(t *testing.T) {
	cases := []struct {
		path, method, body string
		expectedRespBody   string
		expectedStatusCode int
		getProfileFunc     func(ctx context.Context, id string) (profilesvc.Profile, error)
	}{
		{
			path:               "/profile/1",
			method:             "GET",
			expectedRespBody:   `{"profile":{"id":"1","name":"Blah"}}`,
			expectedStatusCode: http.StatusOK,
			getProfileFunc: func(_ context.Context, id string) (profilesvc.Profile, error) {
				return profilesvc.Profile{ID: id, Name: "Blah"}, nil
			},
		},
		{
			path:               "/profile/9999",
			method:             "GET",
			expectedRespBody:   `{"error":"not found"}`,
			expectedStatusCode: http.StatusNotFound,
			getProfileFunc: func(_ context.Context, id string) (profilesvc.Profile, error) {
				return profilesvc.Profile{}, profilesvc.ErrNotFound
			},
		},
	}

	svc := &mock.Service{}
	ctx := context.Background()
	h := profilesvc.MakeHTTPHandler(ctx, svc, kitlog.NewNopLogger())

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			// mock func call
			svc.GetProfileFunc = c.getProfileFunc

			req, _ := http.NewRequest(c.method, c.path, strings.NewReader(c.body))
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != c.expectedStatusCode {
				t.Errorf("unexpected status code: %v, want %v", rec.Code, c.expectedStatusCode)
			}
			respBody := strings.TrimSpace(rec.Body.String())
			if respBody != c.expectedRespBody {
				t.Errorf("unexpected response body: %s, want %s", respBody, c.expectedRespBody)
			}
		})
	}
}

func TestDeleteProfileHandler(t *testing.T) {
	cases := []struct {
		path, method, body string
		expectedRespBody   string
		expectedStatusCode int
	}{
		{
			path:               "/profile/1",
			method:             "DELETE",
			expectedRespBody:   `{}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	svc := &mock.Service{}
	ctx := context.Background()
	h := profilesvc.MakeHTTPHandler(ctx, svc, kitlog.NewNopLogger())
	// mock func call
	svc.DeleteProfileFunc = func(_ context.Context, id string) error {
		return nil
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			req, _ := http.NewRequest(c.method, c.path, strings.NewReader(c.body))
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != c.expectedStatusCode {
				t.Errorf("unexpected status code: %v, want %v", rec.Code, c.expectedStatusCode)
			}
			respBody := strings.TrimSpace(rec.Body.String())
			if respBody != c.expectedRespBody {
				t.Errorf("unexpected response body: %s, want %s", respBody, c.expectedRespBody)
			}
		})
	}
}
