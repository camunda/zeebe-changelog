package github

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/v28/github"
)

func TestAddLabel_404Error(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Create a client with the test server
	client := github.NewClient(nil)
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	ghc := &Client{
		client: client,
		ctx:    context.Background(),
	}

	// Call AddLabel with a non-existent issue
	ghc.AddLabel("testorg", "testrepo", 12345, "test-label")

	// Verify that a warning was logged instead of a fatal error
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Warning: Issue #12345 not found") {
		t.Errorf("Expected warning message in log output, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "testorg/testrepo") {
		t.Errorf("Expected org/repo in log output, got: %s", logOutput)
	}
}
