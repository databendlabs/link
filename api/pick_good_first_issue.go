package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v78/github"
)

func PickGoodFirstIssue(w http.ResponseWriter, r *http.Request) {
	log.Printf("start PickGoodFirstIssue")

	resp, err := http.Get(fmt.Sprintf("https://%s/api/fetch_good_first_issue", os.Getenv("VERCEL_URL")))
	if err != nil {
		log.Fatalf("Fetch good first issues: %s", err)
	}

	log.Printf("X-Vercel-Cache: %s", resp.Header.Get("X-Vercel-Cache"))

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read good first issues: %s", err)
	}
	log.Printf("content: %s", string(content))

	var issues []*github.Issue
	err = json.Unmarshal(content, &issues)
	if err != nil {
		log.Fatalf("Unmarshal good first issues: %s", err)
	}

	// Take current unix nano as seed.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(len(issues))

	w.Header().Add("Location", *issues[index].HTMLURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write(nil)

	log.Printf("finish PickGoodFirstIssue")
}
