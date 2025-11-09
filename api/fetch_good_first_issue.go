package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/google/go-github/v78/github"
)

func FetchGoodFirstIssue(w http.ResponseWriter, r *http.Request) {
	log.Printf("start FetchGoodFirstIssue")

	ctx := r.Context()
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	repos := strings.Split(os.Getenv("GITHUB_REPOS"), ",")

	wg := &sync.WaitGroup{}
	var issues []*github.Issue
	lock := sync.Mutex{}

	for _, v := range repos {
		wg.Add(1)

		go func(repo string) {
			defer wg.Done()
			ownerRepo := strings.Split(repo, "/")
			if len(ownerRepo) < 2 {
				log.Fatalf("Invalid repo: %s", repo)
			}

			is, resp, err := client.Issues.ListByRepo(ctx, ownerRepo[0], ownerRepo[1], &github.IssueListByRepoOptions{Labels: []string{"good first issue"},
				ListOptions: github.ListOptions{
					Page:    0,
					PerPage: 500,
				}})
			if err != nil {
				log.Fatalf("ListByOrg: %s", err)
			}

			log.Printf("rate limit: %v", resp.Header.Get("X-RateLimit-Limit"))
			log.Printf("rate limit used: %v", resp.Header.Get("X-RateLimit-Used"))
			log.Printf("rate limit remaining: %v", resp.Header.Get("X-RateLimit-Remaining"))

			lock.Lock()
			defer lock.Unlock()
			issues = append(issues, is...)
		}(v)
	}
	wg.Wait()

	log.Printf("Got %d issues", len(issues))
	if len(issues) == 0 {
		// Return a Not Found error
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No issues found"))
		return
	}

	content, err := json.Marshal(issues)
	if err != nil {
		log.Fatalf("Invalid issues: %s", err)
	}

	// Cache content up to 600 (10 mins)
	w.Header().Add("Cache-Control", "s-maxage=600")
	w.WriteHeader(http.StatusOK)
	w.Write(content)

	log.Printf("finish FetchGoodFirstIssue")
}
