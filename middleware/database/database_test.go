package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"sort"
	"strconv"
	"testing"

	"github.com/coopnorge/go-datadog-lib/v2/internal/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestRegisterAndOpen(t *testing.T) {
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	db, err := RegisterDriverAndOpen("mysql", &fakeDriver{}, "")
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
	testhelpers.ConfigureDatadog(t)

	// Start Datadog tracer, so that we don't create NoopSpans.
	testTracer := mocktracer.Start()

	db, err := RegisterDriverAndOpen("mysql", &fakeDriver{}, "")
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

// Create fake driver++, to avoid having to import a specific database driver to test.

type fakeDriver struct{}

func (d *fakeDriver) Open(_ string) (driver.Conn, error) {
	return &fakeConn{}, nil
}

type fakeConn struct{}

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) {
	return &fakeStmt{query: query}, nil
}

func (c *fakeConn) Close() error {
	return nil
}

func (c *fakeConn) Begin() (driver.Tx, error) {
	panic("Begin not implemented")
}

type fakeStmt struct {
	query string
}

func (s *fakeStmt) Close() error {
	return nil
}

func (s *fakeStmt) NumInput() int {
	return 0
}

func (s *fakeStmt) Exec(_ []driver.Value) (driver.Result, error) {
	panic("exec not implemented")
}

func (s *fakeStmt) Query(_ []driver.Value) (driver.Rows, error) {
	return &fakeRows{}, nil
}

type fakeRows struct {
	doneReading bool
}

func (r *fakeRows) Close() error {
	return nil
}

func (r *fakeRows) Columns() []string {
	return []string{"value"}
}

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.doneReading {
		return io.EOF
	}

	dest[0] = "hello"
	r.doneReading = true
	return nil
}
