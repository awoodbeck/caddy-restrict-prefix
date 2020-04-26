package restrictprefix

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func TestRestrictPrefix(t *testing.T) {
	next := new(nextHandler)
	rp := new(RestrictPrefix)
	err := rp.Validate()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		path string
		code int
	}{
		{"http://test/sage.svg", http.StatusOK},
		{"http://test/.secret", http.StatusNotFound},
		{"http://test/.dir/secret", http.StatusNotFound},
	}

	for i, c := range testCases {
		r := httptest.NewRequest(http.MethodGet, c.path, nil)
		w := httptest.NewRecorder()
		err = rp.ServeHTTP(w, r, next)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}

		actual := w.Result().StatusCode
		if c.code != actual {
			t.Errorf("%d: expected %d; actual %d", i, c.code, actual)
		}
	}
}

type nextHandler struct{}

func (n nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

var _ caddyhttp.Handler = (*nextHandler)(nil)
