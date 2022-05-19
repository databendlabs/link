package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
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

	repos := []string{"databend", "openraft", "opendal"}

	wg := &sync.WaitGroup{}
	var issues []*github.Issue
	lock := sync.Mutex{}

	for _, v := range repos {
		wg.Add(1)

		go func(repo string) {
			defer wg.Done()

			is, resp, err := client.Issues.ListByRepo(ctx, "datafuselabs", repo, &github.IssueListByRepoOptions{Labels: []string{"good first issue"},
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

	// Take current unix nano as seed.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Shuffle the whole issues and choosing the first one.
	rnd.Shuffle(len(issues), func(i, j int) {
		issues[i], issues[j] = issues[j], issues[i]
	})

	w.Header().Add("Location", *issues[0].HTMLURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write(nil)
}
