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
	"sync"
	"testing"
	"time"

	"github.com/google/go-github/v83/github"
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

func TestAddLabel_Concurrent(t *testing.T) {
	// Track the number of concurrent requests
	concurrentRequests := 0
	maxConcurrent := 0
	var mutex sync.Mutex

	// Create a test server that simulates API responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		concurrentRequests++
		if concurrentRequests > maxConcurrent {
			maxConcurrent = concurrentRequests
		}
		mutex.Unlock()

		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)

		mutex.Lock()
		concurrentRequests--
		mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"name":"test-label"}]`))
	}))
	defer server.Close()

	// Create a client with the test server
	client := github.NewClient(nil)
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	ghc := &Client{
		client: client,
		ctx:    context.Background(),
	}

	// Test concurrent labeling with multiple goroutines
	const numIssues = 20
	var wg sync.WaitGroup

	for i := 1; i <= numIssues; i++ {
		wg.Add(1)
		go func(issueID int) {
			defer wg.Done()
			ghc.AddLabel("testorg", "testrepo", issueID, "test-label")
		}(i)
	}

	wg.Wait()

	// Verify that multiple requests were processed concurrently
	if maxConcurrent <= 1 {
		t.Errorf("Expected concurrent requests, but maxConcurrent was %d", maxConcurrent)
	}
	t.Logf("Max concurrent requests: %d", maxConcurrent)
}

func TestEnsureLabelExists_LabelAlreadyExists(t *testing.T) {
	getCalls := 0
	postCalls := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/testorg/testrepo/labels/test-label":
			getCalls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"test-label","color":"8e8e8e"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/testorg/testrepo/labels":
			postCalls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"name":"test-label","color":"8e8e8e"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := github.NewClient(nil)
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	ghc := &Client{
		client: client,
		ctx:    context.Background(),
		sleep:  func(_ time.Duration) {},
	}

	ghc.EnsureLabelExists("testorg", "testrepo", "test-label")

	if getCalls != 1 {
		t.Errorf("Expected 1 GET call, got %d", getCalls)
	}
	if postCalls != 0 {
		t.Errorf("Expected 0 POST calls, got %d", postCalls)
	}
}

func TestEnsureLabelExists_CreatesAndVerifiesLabel(t *testing.T) {
	getCalls := 0
	postCalls := 0
	var slept time.Duration

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/testorg/testrepo/labels/test-label":
			getCalls++
			w.Header().Set("Content-Type", "application/json")
			if getCalls == 1 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message":"Not Found"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"test-label","color":"8e8e8e"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/testorg/testrepo/labels":
			postCalls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"name":"test-label","color":"8e8e8e"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := github.NewClient(nil)
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	ghc := &Client{
		client: client,
		ctx:    context.Background(),
		sleep:  func(d time.Duration) { slept = d },
	}

	ghc.EnsureLabelExists("testorg", "testrepo", "test-label")

	if getCalls != 2 {
		t.Errorf("Expected 2 GET calls, got %d", getCalls)
	}
	if postCalls != 1 {
		t.Errorf("Expected 1 POST call, got %d", postCalls)
	}
	if slept != labelVerificationRetryDelay {
		t.Errorf("Expected sleep duration %v, got %v", labelVerificationRetryDelay, slept)
	}
}

func TestLabelExists_LabelMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/repos/testorg/testrepo/labels/missing-label" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message":"Not Found"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := github.NewClient(nil)
	baseURL, _ := url.Parse(server.URL + "/")
	client.BaseURL = baseURL

	ghc := &Client{
		client: client,
		ctx:    context.Background(),
		sleep:  func(_ time.Duration) {},
	}

	exists, err := ghc.LabelExists("testorg", "testrepo", "missing-label")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected label to be missing")
	}
}
