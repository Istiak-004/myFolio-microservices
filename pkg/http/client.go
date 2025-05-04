package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/istiak-004/myFolio-microservices/pkg/logger"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// ClientConfig holds configuration for the HTTP client
type ClientConfig struct {
	Timeout         time.Duration
	MaxConnections  int
	IdleConnTimeout time.Duration
}

// Client represents an HTTP client with circuit breaker
type Client struct {
	http.Client
	circuitBreaker *gobreaker.CircuitBreaker // Circuit breaker for fault tolerance
	logger         *logger.Logger
}

// NewClient creates a new HTTP client with circuit breaker
func NewClient(cfg ClientConfig, logger *logger.Logger) *Client {
	transport := &http.Transport{
		MaxIdleConns:        cfg.MaxConnections,
		IdleConnTimeout:     cfg.IdleConnTimeout,
		DisableCompression:  false,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "HTTPClient",
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Info("Circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", fmt.Sprintf("%d", from)),
				zap.String("to", fmt.Sprintf("%d", to)))
		},
	})

	return &Client{
		Client: http.Client{
			Transport: transport,
			Timeout:   cfg.Timeout,
		},
		circuitBreaker: cb,
		logger:         logger,
	}
}

// Do executes an HTTP request with circuit breaker protection
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	_, cbErr := c.circuitBreaker.Execute(func() (interface{}, error) {
		resp, err = c.Client.Do(req.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		// Consider 5xx errors as failures for circuit breaker
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, fmt.Errorf("server error: %d", resp.StatusCode)
		}

		return resp, nil
	})

	if cbErr != nil {
		return nil, cbErr
	}

	return resp, err
}
