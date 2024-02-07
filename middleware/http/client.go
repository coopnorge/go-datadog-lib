package http

import (
	"net/http"

	"github.com/coopnorge/go-datadog-lib/v2/internal"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

// WrapClient wraps the net/http.Client to automatically create child-spans, and append to HTTP Headers.
//
// Deprecated: Use AddTracingToClient instead, and set a proper ResourceNamer. This function will be removed in a later version.
func WrapClient(client *http.Client) *http.Client {
	// Note: Explicitly setting ResourceNamer to `nil`, to prevent leaking paths and keeping backwards-compatibility.
	return AddTracingToClient(client, WithResourceNamer(nil))
}

// AddTracingToClient wraps the net/http.Client to automatically create child-spans, and append to HTTP Headers.
func AddTracingToClient(client *http.Client, options ...Option) *http.Client {
	if !internal.IsDatadogConfigured() {
		return client
	}
	opts := convertClientOptions(options...)
	return httptrace.WrapClient(client, opts...)
}
