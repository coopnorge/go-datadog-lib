package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"sort"
	"strconv"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestRegisterAndOpen(t *testing.T) {
	// Ensure valid datadog config, even if we don't have a datadog agent running, to fully instrument the application.
	t.Setenv("DD_ENV", "unittest")
	t.Setenv("DD_SERVICE", "unittest-service")
	require.True(t, internal.IsDatadogConfigured())

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	db, err := RegisterDriverAndOpen(newMockDriver(), "")
	require.NoError(t, err)

	span, ctx := tracer.StartSpanFromContext(context.Background(), "http.request", tracer.ResourceName("/helloworld"))

	dbString, numRows, err := readFromDB(ctx, db)
	require.NoError(t, err)
	assert.Equal(t, 1, numRows)
	assert.Equal(t, "hello", dbString)

	err = db.Close()
	require.NoError(t, err)

	// Finish the span manually
	span.Finish()

	testTracer.Stop()

	require.Equal(t, 0, len(testTracer.OpenSpans()))
	spans := testTracer.FinishedSpans()
	require.Equal(t, 2, len(spans))
	sort.Slice(spans, func(i, j int) bool { return spans[i].FinishTime().Before(spans[j].FinishTime()) })
	{
		// SQL span
		finishedSpan := spans[0]
		assert.Equal(t, "mysql.query", finishedSpan.OperationName())
		assert.Equal(t, "client", finishedSpan.Tag("span.kind"))
		assert.Equal(t, "Query", finishedSpan.Tag("sql.query_type"))
		assert.Equal(t, "unittest-service", finishedSpan.Tag("service.name"))
	}
	{
		// HTTP span
		finishedSpan := spans[1]
		assert.Equal(t, "http.request", finishedSpan.OperationName())
		assert.Equal(t, "/helloworld", finishedSpan.Tag("resource.name"))
		assert.Equal(t, strconv.Itoa(int(finishedSpan.TraceID())), strconv.Itoa(int(span.Context().TraceID())))
		assert.Equal(t, strconv.Itoa(int(finishedSpan.SpanID())), strconv.Itoa(int(span.Context().SpanID())))
	}
}

func TestRegisterAndOpenNoTrace(t *testing.T) {
	// Ensure valid datadog config, even if we don't have a datadog agent running, to fully instrument the application.
	t.Setenv("DD_ENV", "unittest")
	t.Setenv("DD_SERVICE", "unittest-service")
	require.True(t, internal.IsDatadogConfigured())

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	db, err := RegisterDriverAndOpen(newMockDriver(), "")
	require.NoError(t, err)

	ctx := context.Background() // Note: We are not creating a span in the context.

	dbString, numRows, err := readFromDB(ctx, db)
	require.NoError(t, err)
	assert.Equal(t, 1, numRows)
	assert.Equal(t, "hello", dbString)

	err = db.Close()
	require.NoError(t, err)

	testTracer.Stop()

	require.Equal(t, 0, len(testTracer.OpenSpans()))
	require.Equal(t, 0, len(testTracer.FinishedSpans()))
}

func readFromDB(ctx context.Context, db *sql.DB) (string, int, error) {
	rows, err := db.QueryContext(ctx, "SELECT 'hello' AS value FROM TRACETEST")
	if err != nil {
		return "", 0, err
	}
	var count int
	var myStr string
	for rows.Next() {
		err = rows.Scan(&myStr)
		if err != nil {
			return "", 0, err
		}
		count++
	}
	err = rows.Err()
	if err != nil {
		return "", 0, err
	}
	err = rows.Close()
	if err != nil {
		return "", 0, err
	}
	return myStr, count, nil
}

// Create mock driver and friends, to avoid having to import a specific mysql driver to test.

type mockDriver struct {
	conn driver.Conn
}

func (d *mockDriver) Open(_ string) (driver.Conn, error) {
	return d.conn, nil
}

func newMockDriver() *mockDriver {
	conn := newMockDriverConn()
	return &mockDriver{conn: conn}
}

type mockDriverConn struct{}

func newMockDriverConn() *mockDriverConn {
	return &mockDriverConn{}
}

func (c *mockDriverConn) Prepare(query string) (driver.Stmt, error) {
	return &mockStmt{query: query}, nil
}

func (c *mockDriverConn) Close() error {
	return nil
}

func (c *mockDriverConn) Begin() (driver.Tx, error) {
	panic("Begin not implemented")
}

type mockStmt struct {
	query string
}

func (s *mockStmt) Close() error {
	return nil
}

func (s *mockStmt) NumInput() int {
	return 0
}

func (s *mockStmt) Exec(_ []driver.Value) (driver.Result, error) {
	panic("exec not implemented")
}

func (s *mockStmt) Query(_ []driver.Value) (driver.Rows, error) {
	return &mockRows{}, nil
}

type mockRows struct {
	doneReading bool
}

func (r *mockRows) Close() error {
	return nil
}

func (r *mockRows) Columns() []string {
	return []string{"value"}
}

func (r *mockRows) Next(dest []driver.Value) error {
	if r.doneReading {
		return io.EOF
	}

	dest[0] = "hello"
	r.doneReading = true
	return nil
}
