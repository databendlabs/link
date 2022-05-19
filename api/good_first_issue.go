package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/google/go-github/v44/github"
)

func GoodFirstIssue(w http.ResponseWriter, r *http.Request) {
	client := github.NewClient(nil)

	ctx := context.Background()

	issues, _, err := client.Issues.ListByOrg(ctx, "datafuselabs", &github.IssueListOptions{Labels: []string{"good first issue"}})
	if err != nil {
		fmt.Fprintf(w, "ListByOrg: %s", err)
		return
	}

	index := rand.Intn(len(issues))
	url := *issues[index].URL

	w.Header().Add("Location", url)
	w.WriteHeader(302)
	w.Write(nil)
}
