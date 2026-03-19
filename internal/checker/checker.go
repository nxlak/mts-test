package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/nxlak/mts-test/internal/models"
)

type latestResponse struct {
	Version string `json:"Version"`
}

// search available updates for dependencies
type Checker interface {
	FindUpdates(ctx context.Context, deps []models.Dependency) ([]models.Update, error)
}

type checker struct {
	httpClient *http.Client
	proxyBase  string
	workers    int
}

type checkerOption func(*checker)

// sets custom proxy url
func WithProxyBase(base string) checkerOption {
	return func(c *checker) { c.proxyBase = strings.TrimRight(base, "/") }
}

// sets max number of concurrent requests
func WithWorkers(n int) checkerOption {
	return func(c *checker) { c.workers = n }
}

// sets custom http client
func WithHTTPClient(client *http.Client) checkerOption {
	return func(c *checker) { c.httpClient = client }
}

func NewChecker(opts ...checkerOption) Checker {
	const (
		defaultProxyBase   = "https://proxy.golang.org"
		defaultWorkers     = 8
		defaultHTTPTimeout = 15 * time.Second
	)

	c := &checker{
		httpClient: &http.Client{Timeout: defaultHTTPTimeout},
		proxyBase:  defaultProxyBase,
		workers:    defaultWorkers,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *checker) FindUpdates(ctx context.Context, deps []models.Dependency) ([]models.Update, error) {
	type job struct {
		idx int
		dep models.Dependency
	}

	type result struct {
		update models.Update
		hasNew bool
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]result, len(deps))
	jobs := make(chan job, c.workers)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	for w := 0; w < c.workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case j, ok := <-jobs:
					if !ok {
						return
					}

					latest, err := c.fetchLatest(ctx, j.dep.Path)
					if err != nil {
						select {
						case errChan <- err:
						default:
						}
						cancel()
						return
					}

					if isNewer(latest, j.dep.CurrentVersion) {
						results[j.idx] = result{
							update: models.Update{
								Dependency:    j.dep,
								LatestVersion: latest,
							},
							hasNew: true,
						}
					}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)

		for i, dep := range deps {
			select {
			case <-ctx.Done():
				return
			case jobs <- job{idx: i, dep: dep}:
			}
		}
	}()

	wg.Wait()
	close(errChan)

	if err, ok := <-errChan; ok {
		return nil, fmt.Errorf("proxy: check updates: %w", err)
	}

	updates := make([]models.Update, 0, len(deps))
	for _, r := range results {
		if r.hasNew {
			updates = append(updates, r.update)
		}
	}

	return updates, nil
}

func (c *checker) fetchLatest(ctx context.Context, modulePath string) (string, error) {
	encoded, err := module.EscapePath(modulePath)
	if err != nil {
		return "", fmt.Errorf("proxy: encode path for %q: %w", modulePath, err)
	}

	reqURL := fmt.Sprintf("%s/%s/@latest", c.proxyBase, encoded)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("proxy: build request for %q: %w", modulePath, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("proxy: GET %q: %w", reqURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return "", fmt.Errorf("proxy: module %q not found (status %d)", modulePath, resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("proxy: unexpected status %d for %q", resp.StatusCode, modulePath)
	}

	var payload latestResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("proxy: decode response for %q: %w", modulePath, err)
	}
	if payload.Version == "" {
		return "", fmt.Errorf("proxy: empty version for %q", modulePath)
	}

	return payload.Version, nil
}

func isNewer(latestVer, currentVer string) bool {
	if !semver.IsValid(latestVer) || !semver.IsValid(currentVer) {
		return false
	}
	return semver.Compare(latestVer, currentVer) > 0
}
