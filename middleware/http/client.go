package http //nolint:revive

import (
	"net/http"

	httptrace "github.com/DataDog/dd-trace-go/contrib/net/http/v2"
	"github.com/coopnorge/go-datadog-lib/v2/internal"
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
	if internal.IsDatadogDisabled() {
		return client
	}
	opts := convertClientOptions(options...)
	return httptrace.WrapClient(client, opts...)
}
