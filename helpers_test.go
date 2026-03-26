package gitforge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCall struct {
	Method  string
	Path    string
	Query   string
	Headers http.Header
	Body    json.RawMessage
}

type mockResponse struct {
	Status int
	Body   any
}

func setupMock(t *testing.T, responses ...mockResponse) (*httptest.Server, *[]mockCall) {
	t.Helper()
	calls := &[]mockCall{}
	idx := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body json.RawMessage
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&body)
		}
		*calls = append(*calls, mockCall{
			Method:  r.Method,
			Path:    r.URL.Path,
			Query:   r.URL.RawQuery,
			Headers: r.Header,
			Body:    body,
		})
		if idx >= len(responses) {
			w.WriteHeader(500)
			return
		}
		resp := responses[idx]
		idx++
		w.WriteHeader(resp.Status)
		if resp.Body != nil {
			_ = json.NewEncoder(w).Encode(resp.Body)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, calls
}
