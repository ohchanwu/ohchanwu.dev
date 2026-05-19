package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type config struct {
	tlsDomains   []string
	acmeEmail    string
	acmeCacheDir string
	acmeStaging  bool
	httpAddr     string
	httpsAddr    string
}

func parseConfig() (config, error) {
	cfg := config{
		acmeCacheDir: getenvDefault("ACME_CACHE_DIR", "/var/cache/autocert"),
		httpAddr:     getenvDefault("HTTP_ADDR", ":8080"),
		httpsAddr:    getenvDefault("HTTPS_ADDR", ":8443"),
		acmeEmail:    os.Getenv("ACME_EMAIL"),
		acmeStaging:  parseBool(os.Getenv("ACME_STAGING")),
	}

	if raw := strings.TrimSpace(os.Getenv("TLS_DOMAINS")); raw != "" {
		for _, d := range strings.Split(raw, ",") {
			if d = strings.TrimSpace(d); d != "" {
				cfg.tlsDomains = append(cfg.tlsDomains, d)
			}
		}
	}

	if len(cfg.tlsDomains) > 0 && cfg.acmeEmail == "" {
		return cfg, errors.New("ACME_EMAIL is required when TLS_DOMAINS is set")
	}
	return cfg, nil
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}

func newManager(cfg config) *autocert.Manager {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(cfg.acmeCacheDir),
		HostPolicy: autocert.HostWhitelist(cfg.tlsDomains...),
		Email:      cfg.acmeEmail,
	}
	if cfg.acmeStaging {
		m.Client = &acme.Client{DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory"}
	}
	return m
}

func run(ctx context.Context, cfg config) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))

	if len(cfg.tlsDomains) == 0 {
		log.Printf("TLS disabled (TLS_DOMAINS unset); serving plain HTTP on %s", cfg.httpAddr)
		srv := &http.Server{
			Addr:              cfg.httpAddr,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		}
		return serveUntilShutdown(ctx, srv, nil)
	}

	m := newManager(cfg)
	log.Printf("ACME enabled for %s (staging=%v, cache=%s); HTTP %s, HTTPS %s",
		strings.Join(cfg.tlsDomains, ","), cfg.acmeStaging, cfg.acmeCacheDir, cfg.httpAddr, cfg.httpsAddr)

	httpSrv := &http.Server{
		Addr:              cfg.httpAddr,
		Handler:           m.HTTPHandler(nil),
		ReadHeaderTimeout: 10 * time.Second,
	}
	httpsSrv := &http.Server{
		Addr:              cfg.httpsAddr,
		Handler:           mux,
		TLSConfig:         m.TLSConfig(),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return serveUntilShutdown(ctx, httpSrv, httpsSrv)
}

func serveUntilShutdown(ctx context.Context, httpSrv, httpsSrv *http.Server) error {
	errCh := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()

	if httpsSrv != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := httpsSrv.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("https server: %w", err)
			}
		}()
	}

	var runErr error
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case runErr = <-errCh:
		log.Printf("server error: %v", runErr)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	if httpsSrv != nil {
		_ = httpsSrv.Shutdown(shutdownCtx)
	}
	wg.Wait()
	return runErr
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
