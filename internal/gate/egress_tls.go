package gate

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/quic-go/quic-go/http3"
)

type EgressTLSClient struct {
	mu          sync.RWMutex
	http2Client *http.Client
	http3Client *http.Client
}

func NewEgressTLSClient() *EgressTLSClient {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: false,
	}

	// Standard HTTP/2 Client with custom dialer
	h2Transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	// Native HTTP/3 (QUIC) Transport for modern CDN traversal
	h3Transport := &http3.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &EgressTLSClient{
		http2Client: &http.Client{Transport: h2Transport, Timeout: 30 * time.Second},
		http3Client: &http.Client{Transport: h3Transport, Timeout: 30 * time.Second},
	}
}

// DoExecute routes requests over HTTP/3 with automatic fallback to HTTP/2
func (e *EgressTLSClient) DoExecute(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	// Attempt HTTP/3 transport first
	resp, err := e.http3Client.Do(req)
	if err == nil {
		return resp, nil
	}

	// Fallback to fingerprint-matched HTTP/2
	return e.http2Client.Do(req)
}

// FetchAudioStream performs a secure GET request with exponential backoff retry.
func (e *EgressTLSClient) FetchAudioStream(ctx context.Context, targetURL string) (io.ReadCloser, error) {
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "audio/mpeg, audio/wav, audio/*;q=0.9, */*;q=0.8")

		resp, err := e.DoExecute(ctx, req)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				return resp.Body, nil
			}
			resp.Body.Close()
			lastErr = fmt.Errorf("egress: server returned status: %s", resp.Status)
		} else {
			lastErr = err
		}

		// Exponential backoff
		backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
		select {
		case <-time.After(backoff):
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("egress fetch failed after %d retries: %w", maxRetries, lastErr)
}
