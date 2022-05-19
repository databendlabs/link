package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

func GoodFirstIssue(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

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

	// Take current unix nano as seed.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(len(issues))

	w.Header().Add("Location", *issues[index].HTMLURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write(nil)
}
