package route

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	router := New()
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://localhost:9090/foo", nil)
	if err != nil {
		t.Fatalf("Error building test request: %s", err)
	}

	router.Redirect(w, r, "/some/endpoint", http.StatusFound)
	if w.Code != http.StatusFound {
		t.Fatalf("Unexpected redirect status code: got %d, want %d", w.Code, http.StatusFound)
	}

	want := "/some/endpoint"
	got := w.Header()["Location"][0]
	if want != got {
		t.Fatalf("Unexpected redirect location: got %s, want %s", got, want)
	}
}

func TestInstrumentation(t *testing.T) {
	var got string
	cases := []struct {
		router *Router
		want   string
	}{
		{
			router: New(),
			want:   "",
		}, {
			router: New().WithInstrumentation(func(handlerName string, handler http.Handler) http.Handler {
				got = handlerName
				return handler
			}),
			want: "/foo",
		},
	}

	for _, c := range cases {
		c.router.Get("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		r, err := http.NewRequest("GET", "http://localhost:9090/foo", nil)
		if err != nil {
			t.Fatalf("Error building test request: %s", err)
		}
		c.router.ServeHTTP(nil, r)
		if c.want != got {
			t.Fatalf("Unexpected value: want %q, got %q", c.want, got)
		}
	}
}
