package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v66/github"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	client *github.Client
)

func main() {

	fmt.Println("Starting server....")
	err := run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func newServer(client *github.Client) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(mux, client)
	var handler http.Handler = mux
	return handler
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	registerMetrics()

	client = github.NewClient(nil)
	server := newServer(client)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutDownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutDownCtx, 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down the serrer: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func AddRoutes(mux *http.ServeMux, client *github.Client) {
	mux.Handle("/top/{language}", handleGetTopLanguageRepos(client))
	mux.Handle("/repos/{user}", handleUserRepos(client))
	mux.Handle("/ping", handlePing())
	mux.Handle("/livez", handleLiveCheck())
	mux.Handle("/metrics", promhttp.Handler())

}

func handleGetTopLanguageRepos(client *github.Client) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		topCounter.Inc()
		language := r.PathValue("language")
		if language == "" {
			err := fmt.Errorf("no language parameter provided")
			WriteErrorResponse(w, err, http.StatusBadRequest)
			return

		}
		getTopRepos(r.Context(), w, client, language)
	})
}

type RepoResponse struct {
	Name        string   `json:"name"`
	Link        string   `json:"link"`
	Language    string   `json:"language"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
}

func getTopRepos(ctx context.Context, w http.ResponseWriter, client *github.Client, language string) {
	repos, _, err := client.Search.Repositories(ctx, language, nil)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadGateway)
		return
	}

	if repos == nil || len(repos.Repositories) == 0 {
		WriteErrorResponse(w, fmt.Errorf("no repos found for language %s", language), http.StatusBadGateway)
		return
	}

	var result []*RepoResponse

	for _, repo := range repos.Repositories {
		result = append(result, &RepoResponse{
			Name:        repo.GetName(),
			Link:        repo.GetHTMLURL(),
			Language:    repo.GetLanguage(),
			Description: repo.GetDescription(),
			Topics:      repo.Topics,
		})
	}

	WriteJSONResponse(w, http.StatusOK, &result)
}

func handleUserRepos(client *github.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRepoCounter.Inc()
		user := r.PathValue("user")
		if user == "" {
			WriteErrorResponse(w, fmt.Errorf("no user parameter provided"), http.StatusBadRequest)
			return
		}
		getUserRepos(r.Context(), w, client, user)
	})
}

func getUserRepos(ctx context.Context, w http.ResponseWriter, client *github.Client, userName string) {
	user, _, err := client.Users.Get(ctx, userName)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusNotFound)
		return
	}
	if user == nil {
		WriteErrorResponse(w, fmt.Errorf("user %s not found", userName), http.StatusNotFound)
		return
	}

	repos, _, err := client.Search.Repositories(ctx, user.GetHTMLURL(), nil)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if repos == nil {
		WriteErrorResponse(w, fmt.Errorf("no repos for user %s found", userName), http.StatusNotFound)
		return
	}
	var result []*RepoResponse
	for _, repo := range repos.Repositories {
		result = append(result, &RepoResponse{
			Name:        repo.GetName(),
			Link:        repo.GetHTMLURL(),
			Language:    repo.GetLanguage(),
			Description: repo.GetDescription(),
			Topics:      repo.Topics,
		})
	}
	userRepoResultGauge.WithLabelValues(userName).Set(float64(len(result)))
	WriteJSONResponse(w, http.StatusOK, &result)
}

type PingResponse struct {
	Reply string `json:"reply"`
	Time  string `json:"time"`
}

func handlePing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pingCounter.Inc()
		resTime := time.Now().Format("2006-01-02 15:04:05")
		res := &PingResponse{
			Reply: "Pong",
			Time:  resTime,
		}
		WriteJSONResponse(w, http.StatusOK, res)
	})
}

func handleLiveCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func WriteJSONResponse(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)

	_, _ = w.Write(data)
}

func WriteErrorResponse(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(map[string]string{"error": err.Error()})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(data)

}
