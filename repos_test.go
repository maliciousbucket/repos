package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v66/github"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"testing"
)

func TestTopRepos(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	gitClient := github.NewClient(nil)
	mux := http.NewServeMux()
	AddRoutes(mux, gitClient)
	tests := []struct {
		language            string
		expectedHTTPStatus  int
		expectedResultCount int
		description         string
		expectError         bool
	}{
		{
			language:            "go",
			expectedHTTPStatus:  http.StatusOK,
			expectedResultCount: 30,
			description:         "Returns 200 ok with top Golang repositories",
			expectError:         false,
		},
		{
			language:            "python",
			expectedHTTPStatus:  http.StatusOK,
			expectedResultCount: 30,
			description:         "Returns 200OOK with top Python repositories",
			expectError:         false,
		},
		{
			language:            "SWE40006",
			expectedHTTPStatus:  http.StatusBadGateway,
			expectedResultCount: 30,
			description:         "For some reason returns results ....",
			expectError:         false,
		},
		{
			language:            "",
			expectedHTTPStatus:  http.StatusNotFound,
			expectedResultCount: 0,
			description:         "Returns bad request for empty language string",
			expectError:         true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := newLanguageRequest(test.language)
			log.Println(request)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)

			res := w.Result()
			if !test.expectError {
				defer res.Body.Close()
				data, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				var response []*RepoResponse
				err = json.Unmarshal(data, &response)
				if err != nil {
					t.Fatal(err)
				}
				assertResultCount(t, response, test.expectedResultCount)
			}

			assertStatus(t, res.StatusCode, test.expectedHTTPStatus)

		})
	}
	ctx.Done()

}

func TestUserRepos(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	gitClient := github.NewClient(nil)
	mux := http.NewServeMux()
	AddRoutes(mux, gitClient)

	tests := []struct {
		user               string
		expectedHTTPStatus int
		description        string
	}{
		{
			user:               "maliciousbucket",
			expectedHTTPStatus: http.StatusOK,
			description:        "Returns 200 ok",
		},
		{
			user:               "swe40006isagoodunit",
			expectedHTTPStatus: http.StatusNotFound,
			description:        "Returns bad request",
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := newUserRequest(test.user)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			res := w.Result()
			assertStatus(t, res.StatusCode, test.expectedHTTPStatus)
		})
	}
	ctx.Done()

}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResultCount(t testing.TB, result []*RepoResponse, want int) {
	t.Helper()
	if len(result) != want {
		t.Errorf("got %d results, want %d", len(result), want)
	}
}

func newLanguageRequest(language string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/top/%s", language), nil)
	return req
}

func newUserRequest(user string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/repos/%s", user), nil)
	return req
}

func newPingRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	return req
}
