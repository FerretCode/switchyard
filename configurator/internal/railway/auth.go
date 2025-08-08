package railway

import "net/http"

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.token)
	r.Header.Set("Content-Type", "application/json")

	return t.wrapped.RoundTrip(r)
}
