package browser

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const fixtureDir = "testdata/pages"

// fixtureServer starts a test HTTP server that serves the static fixture
// pages from testdata/pages under /pages and adds the dynamic endpoints used
// by those pages:
//
//   - /redirect      -> 302 to /pages/nav.html (navigation fixture)
//   - /api/echo      -> JSON {"method","q"} (network / DOM fetch fixture)
func fixtureServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle("/pages/", http.StripPrefix("/pages", http.FileServer(http.Dir(fixtureDir))))

	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pages/nav.html", http.StatusFound)
	})

	mux.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"method": r.Method,
			"q":      r.URL.Query().Get("q"),
		}); err != nil {
			t.Errorf("encode echo: %v", err)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

func fetchFixture(t *testing.T, client *http.Client, url string) (int, http.Header, string) {
	t.Helper()

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read %s: %v", url, err)
	}

	return resp.StatusCode, resp.Header, string(body)
}

func TestFixtureNavigationPages(t *testing.T) {
	srv := fixtureServer(t)

	for _, tc := range []struct {
		name string
		path string
		want []string
	}{
		{
			name: "index",
			path: "/pages/index.html",
			want: []string{"go-bidi fixtures", "nav.html", "network.html", "search-form"},
		},
		{
			name: "nav",
			path: "/pages/nav.html",
			want: []string{"navigation page", "next.html", "/redirect"},
		},
		{
			name: "next",
			path: "/pages/next.html",
			want: []string{"next page", "index.html"},
		},
		{
			name: "network",
			path: "/pages/network.html",
			want: []string{"network page", "pixel.png", "styles.css", "app.js"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			status, _, body := fetchFixture(t, http.DefaultClient, srv.URL+tc.path)
			if status != http.StatusOK {
				t.Fatalf("status = %d, want 200", status)
			}

			for _, want := range tc.want {
				if !strings.Contains(body, want) {
					t.Errorf("body does not contain %q", want)
				}
			}
		})
	}
}

func TestFixtureNetworkResources(t *testing.T) {
	srv := fixtureServer(t)

	for _, tc := range []struct {
		path string
		want string
	}{
		{path: "/pages/styles.css", want: "text/css"},
		{path: "/pages/app.js", want: "text/javascript"},
		{path: "/pages/data.json", want: "application/json"},
		{path: "/pages/pixel.png", want: "image/png"},
	} {
		status, header, body := fetchFixture(t, http.DefaultClient, srv.URL+tc.path)
		if status != http.StatusOK {
			t.Fatalf("GET %s status = %d, want 200", tc.path, status)
		}

		if ctype := header.Get("Content-Type"); !strings.HasPrefix(ctype, tc.want) {
			t.Errorf("GET %s Content-Type = %q, want prefix %q", tc.path, ctype, tc.want)
		}

		if body == "" {
			t.Errorf("GET %s: empty body", tc.path)
		}
	}
}

func TestFixtureRedirect(t *testing.T) {
	srv := fixtureServer(t)

	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	status, header, _ := fetchFixture(t, client, srv.URL+"/redirect")
	if status != http.StatusFound {
		t.Fatalf("status = %d, want 302", status)
	}

	if location := header.Get("Location"); location != "/pages/nav.html" {
		t.Errorf("Location = %q, want /pages/nav.html", location)
	}
}

func TestFixtureEchoAPI(t *testing.T) {
	srv := fixtureServer(t)

	status, header, body := fetchFixture(t, http.DefaultClient, srv.URL+"/api/echo?q=hello")
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}

	if ctype := header.Get("Content-Type"); !strings.HasPrefix(ctype, "application/json") {
		t.Errorf("Content-Type = %q, want prefix application/json", ctype)
	}

	var got map[string]string
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("decode echo: %v", err)
	}

	if got["method"] != "GET" || got["q"] != "hello" {
		t.Errorf("echo = %+v, want method GET q hello", got)
	}
}
