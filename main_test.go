package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    config
		wantErr bool
	}{
		{
			name: "all unset → TLS off with defaults",
			env:  map[string]string{},
			want: config{
				acmeCacheDir: "/var/cache/autocert",
				httpAddr:     ":8080",
				httpsAddr:    ":8443",
			},
		},
		{
			name: "TLS_DOMAINS without ACME_EMAIL → error",
			env: map[string]string{
				"TLS_DOMAINS": "ohchanwu.dev",
			},
			wantErr: true,
		},
		{
			name: "TLS_DOMAINS + ACME_EMAIL → valid",
			env: map[string]string{
				"TLS_DOMAINS": "ohchanwu.dev",
				"ACME_EMAIL":  "ohchanwu@gmail.com",
			},
			want: config{
				tlsDomains:   []string{"ohchanwu.dev"},
				acmeEmail:    "ohchanwu@gmail.com",
				acmeCacheDir: "/var/cache/autocert",
				httpAddr:     ":8080",
				httpsAddr:    ":8443",
			},
		},
		{
			name: "ACME_STAGING=true and =1 both parse",
			env: map[string]string{
				"TLS_DOMAINS":  "ohchanwu.dev",
				"ACME_EMAIL":   "ohchanwu@gmail.com",
				"ACME_STAGING": "1",
			},
			want: config{
				tlsDomains:   []string{"ohchanwu.dev"},
				acmeEmail:    "ohchanwu@gmail.com",
				acmeCacheDir: "/var/cache/autocert",
				acmeStaging:  true,
				httpAddr:     ":8080",
				httpsAddr:    ":8443",
			},
		},
		{
			name: "multiple domains with whitespace are trimmed",
			env: map[string]string{
				"TLS_DOMAINS": "ohchanwu.dev, www.ohchanwu.dev",
				"ACME_EMAIL":  "ohchanwu@gmail.com",
			},
			want: config{
				tlsDomains:   []string{"ohchanwu.dev", "www.ohchanwu.dev"},
				acmeEmail:    "ohchanwu@gmail.com",
				acmeCacheDir: "/var/cache/autocert",
				httpAddr:     ":8080",
				httpsAddr:    ":8443",
			},
		},
		{
			name: "custom addrs and cache dir override defaults",
			env: map[string]string{
				"HTTP_ADDR":      ":7080",
				"HTTPS_ADDR":     ":7443",
				"ACME_CACHE_DIR": "/tmp/certs",
			},
			want: config{
				acmeCacheDir: "/tmp/certs",
				httpAddr:     ":7080",
				httpsAddr:    ":7443",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, k := range []string{"TLS_DOMAINS", "ACME_EMAIL", "ACME_CACHE_DIR", "ACME_STAGING", "HTTP_ADDR", "HTTPS_ADDR"} {
				t.Setenv(k, "")
			}
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			got, err := parseConfig()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (config=%+v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("config mismatch\n got:  %+v\n want: %+v", got, tc.want)
			}
		})
	}
}

func TestHostPolicyRejectsUnknown(t *testing.T) {
	cfg := config{
		tlsDomains:   []string{"ohchanwu.dev"},
		acmeEmail:    "ohchanwu@gmail.com",
		acmeCacheDir: t.TempDir(),
	}
	m := newManager(cfg)

	if err := m.HostPolicy(context.Background(), "ohchanwu.dev"); err != nil {
		t.Errorf("expected whitelisted host to pass, got %v", err)
	}
	if err := m.HostPolicy(context.Background(), "evil.example.com"); err == nil {
		t.Errorf("expected non-whitelisted host to be rejected, got nil")
	}
}

func TestInlinePDFsSetsDispositionOnlyForPDFs(t *testing.T) {
	handler := inlinePDFs(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	pdfReq := httptest.NewRequest(http.MethodGet, "/resume_kr.pdf", nil)
	pdfRes := httptest.NewRecorder()
	handler.ServeHTTP(pdfRes, pdfReq)
	if got := pdfRes.Header().Get("Content-Disposition"); got != "inline" {
		t.Fatalf("PDF Content-Disposition = %q, want inline", got)
	}

	htmlReq := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	htmlRes := httptest.NewRecorder()
	handler.ServeHTTP(htmlRes, htmlReq)
	if got := htmlRes.Header().Get("Content-Disposition"); got != "" {
		t.Fatalf("HTML Content-Disposition = %q, want empty", got)
	}
}
