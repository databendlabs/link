package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGoodFirstIssue(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	os.Setenv("GITHUB_REPOS", "datafuselabs/databend,datafuselabs/opendal")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GoodFirstIssue)

	handler.ServeHTTP(rr, req)

	t.Logf("%s", rr.Header())
	t.Logf("%s", rr.Body)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
