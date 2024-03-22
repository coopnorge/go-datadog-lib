package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coopnorge/go-datadog-lib/v2/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestWrapClient(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	span, ctx := tracer.StartSpanFromContext(context.Background(), "http.request", tracer.ResourceName("/helloworld"))
	defer span.Finish()

	var traceparent string
	var ddParentID string
	var ddTraceID string
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceparent = r.Header.Get("Traceparent")
		ddTraceID = r.Header.Get("X-Datadog-Trace-Id")
		ddParentID = r.Header.Get("X-Datadog-Parent-Id")
		_, _ = w.Write([]byte("OK"))
	}))

	t.Cleanup(s.Close)

	client := WrapClient(&http.Client{Timeout: 500 * time.Millisecond})
	req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

	testTracer.Stop()

	spans := testTracer.FinishedSpans()
	require.Equal(t, 1, len(spans))
	finishedSpan := spans[0]
	assert.Equal(t, strconv.Itoa(int(finishedSpan.TraceID())), ddTraceID)
	assert.Equal(t, strconv.Itoa(int(finishedSpan.SpanID())), ddParentID)

	assert.Empty(t, traceparent, "Datadog's mocktracer does not propagate W3C-headers as of writing this test. If they start propagating it, we should remove the separate test below, and update this test to assert the correct W3C-header.")
}

func TestWrapClientW3C(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	// Start real tracer (not mocktracer), to propagate Traceparent.
	tracer.Start()

	span, ctx := tracer.StartSpanFromContext(context.Background(), "http.request", tracer.ResourceName("/helloworld"))
	defer span.Finish()

	var traceparent string
	var tracestate string
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceparent = r.Header.Get("Traceparent")
		tracestate = r.Header.Get("Tracestate")
		_, _ = w.Write([]byte("OK"))
	}))

	t.Cleanup(s.Close)

	client := WrapClient(&http.Client{Timeout: 500 * time.Millisecond})
	req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.NoError(t, err)

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
	assert.Equal(t, "01", parts[3], "w3c trace-flags not is not correct")

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
