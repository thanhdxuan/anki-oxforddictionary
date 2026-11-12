package oxford

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	HTTP      *http.Client
	BaseURL   string
	UserAgent string
}

const DefaultUserAgent = "oxford-lookup/0.1"

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}

func (c *Client) Lookup(ctx context.Context, word string) (Result, error) {
	formattedWord := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(word)), " ", "-")
	u := fmt.Sprintf("%s/definition/english/%s", c.BaseURL, url.PathEscape(formattedWord))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return Result{}, fmt.Errorf("build request: %w", err)
	}

	hc := c.HTTP
	ua := c.UserAgent

	if hc == nil {
		hc = defaultHTTPClient
	}

	if ua == "" {
		ua = DefaultUserAgent
	}

	req.Header.Set("User-Agent", ua)
	resp, err := hc.Do(req)

	if err != nil {
		return Result{}, fmt.Errorf("fetch %s: %w", u, err)
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	reqURL := resp.Request.URL.String()
	if code == http.StatusNotFound || strings.HasPrefix(resp.Request.URL.Path, "/spellcheck/english") {
		return Result{SourceURL: reqURL}, ErrNotFound
	}

	if code != http.StatusOK {
		return Result{SourceURL: reqURL}, fmt.Errorf("fetch %s: %d", u, code)
	}

	result, err := Parse(resp.Body)

	result.SourceURL = reqURL
	result.Word = word

	if err != nil {
		return result, err
	}

	return result, nil
}
