package http_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"
	datadogMiddleware "github.com/coopnorge/go-datadog-lib/v2/middleware/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapClient(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	var traceparent string
	var ddParentID string
	var ddTraceID string
	var tracestate string
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceparent = r.Header.Get("Traceparent")
		tracestate = r.Header.Get("Tracestate")
		ddTraceID = r.Header.Get("X-Datadog-Trace-Id")
		ddParentID = r.Header.Get("X-Datadog-Parent-Id")
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	}))

	t.Cleanup(s.Close)

	client := datadogMiddleware.WrapClient(&http.Client{Timeout: 500 * time.Millisecond})

	span, ctx := tracer.StartSpanFromContext(context.Background(), "http.request", tracer.ResourceName("/helloworld"))
	req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	span.Finish()

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 2, len(spans))
	finishedSpan := spans[0]
	assert.Equal(t, strconv.FormatUint(finishedSpan.TraceID(), 10), ddTraceID)
	assert.Equal(t, strconv.Itoa(int(finishedSpan.SpanID())), ddParentID)

	// Assert TraceParent
	require.NotEmpty(t, traceparent)
	parts := strings.Split(traceparent, "-")
	require.Equal(t, 4, len(parts))
	// version
	assert.Equal(t, "00", parts[0], "w3c version is not correct")
	// trace-id
	assert.Equal(t, 32, len(parts[1]), "w3c trace-id has invalid length")
	assert.NotEqual(t, "00000000000000000000000000000000", parts[1], "w3c trace-id is zero")
	// parent-id
	assert.Equal(t, 16, len(parts[2]), "w3c parent-id has invalid length")
	assert.NotEqual(t, "0000000000000000", parts[2], "w3c parent-id is zero")
	// trace-flags
	assert.Equal(t, "00", parts[3], "w3c trace-flags not is not correct")

	// Assert TraceState
	parts = strings.Split(tracestate, ",")
	require.True(t, len(parts) >= 1)
	found := false
	for _, listMember := range parts {
		if strings.HasPrefix(listMember, "dd=") {
			assert.NotEmpty(t, strings.TrimPrefix(listMember, "dd="))
			found = true
		}
	}
	assert.True(t, found, "Did not find Datadog's list-member in w3c tracestate")
}

func TestURLIsNotInTags(t *testing.T) {
	// This test is a regression-test for this issue: https://github.com/coopnorge/go-datadog-lib/issues/495
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	span, ctx := tracer.StartSpanFromContext(context.Background(), "http.request", tracer.ResourceName("/helloworld"))

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("OK"))
		require.NoError(t, err)
	}))

	t.Cleanup(s.Close)

	url := fmt.Sprintf("%s/some-path-with-pii?some-query-with-pii=true", s.URL)

	// Adding tracing to client, with a static resource-name, as we want to make sure that no tags automatically add the full URL, which might contain PII (Personally Identifiable Information).
	client := datadogMiddleware.AddTracingToClient(&http.Client{Timeout: 500 * time.Millisecond}, datadogMiddleware.WithResourceNamer(datadogMiddleware.StaticResourceNamer("my-resource-name")))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	span.Finish()
	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 2, len(spans))
	finishedSpan := spans[0]
	assert.Empty(t, finishedSpan.Tag(ext.HTTPURL))
	assert.NotContains(t, finishedSpan.OperationName(), "pii")
	assert.NotContains(t, finishedSpan.String(), "pii")
	// Iterate over every tag, to make sure that none of the tags contain the full URL or the word 'pii' from the URL.
	for tag, tagValue := range finishedSpan.Tags() {
		if str, ok := tagValue.(string); ok {
			assert.NotContains(t, str, "pii", "Tag %q contained the word 'pii', but it should not, as it is somehow picked up from the path or query! Full value: %s", tag, str)
			assert.NotContains(t, str, url, "Tag %q contained the full URL, but it should not, as it might contain PII! Full value: %s", tag, str)
		}
	}
}
