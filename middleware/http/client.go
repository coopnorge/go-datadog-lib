package http

import (
	"net/http"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

// WrapClient wraps the net/http.Client to automatically create child-spans, and append to HTTP Headers.
func WrapClient(client *http.Client) *http.Client {
	if internal.IsDatadogConfigured() {
		client = httptrace.WrapClient(client)
	}
	return client
}
