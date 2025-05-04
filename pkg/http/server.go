package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/istiak-004/myFolio-microservices/pkg/logger"
	"golang.org/x/crypto/acme/autocert"
)

// ServerConfig holds the configuration for the HTTP server
type ServerConfig struct {
	Host            string
	Mode            string // "debug", "release", or "test"
	TLSCertPath     string
	TLSKeyPath      string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	EnableHTTPS     bool
	DomainWhitelist []string
}

// Server represents an HTTP server
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	logger     *logger.Logger
	config     ServerConfig
}

// NewServer creates a new HTTP server with sensible defaults
func NewServer(cfg ServerConfig, logger *logger.Logger) *Server {
	// Set Gin mode
	gin.SetMode(cfg.Mode)

	// Create router with default middleware stack
	router := gin.New()
	router.Use(
		gin.Recovery(), // Handle panics
	)

	// Configure server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: srv,
		router:     router,
		logger:     logger,
		config:     cfg,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	if s.config.EnableHTTPS {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
			PreferServerCipherSuites: true,
		}

		if s.config.TLSCertPath != "" && s.config.TLSKeyPath != "" {
			// Use provided certs
			tlsConfig.Certificates = make([]tls.Certificate, 1)
			cert, err := tls.LoadX509KeyPair(s.config.TLSCertPath, s.config.TLSKeyPath)
			if err != nil {
				return fmt.Errorf("failed to load TLS cert: %w", err)
			}
			tlsConfig.Certificates[0] = cert
			s.httpServer.TLSConfig = tlsConfig
			return s.httpServer.ListenAndServeTLS("", "")
		} else {
			// Use Let's Encrypt
			certManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(s.config.DomainWhitelist...),
				Cache:      autocert.DirCache("/var/www/.cache"),
			}
			tlsConfig.GetCertificate = certManager.GetCertificate
			s.httpServer.TLSConfig = tlsConfig

			// Redirect HTTP to HTTPS
			go func() {
				redirectServer := &http.Server{
					Addr:         ":80",
					Handler:      certManager.HTTPHandler(nil),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 5 * time.Second,
					IdleTimeout:  120 * time.Second,
				}
				redirectServer.ListenAndServe()
			}()

			return s.httpServer.ListenAndServeTLS("", "")
		}
	}
	s.logger.Info("Starting HTTP server ", s.logger.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server without interrupting any active connections.
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	return nil
}

// Router returns the underlying Gin router
func (s *Server) Router() *gin.Engine {
	return s.router
}

// RegisterMiddleware registers global middleware
func (s *Server) RegisterMiddleware(middleware ...gin.HandlerFunc) {
	s.router.Use(middleware...)
}

// RegisterRouteGroup creates a new route group with common middleware
func (s *Server) RegisterRouteGroup(prefix string, middleware ...gin.HandlerFunc) *gin.RouterGroup {
	group := s.router.Group(prefix)
	group.Use(middleware...)
	return group
}
