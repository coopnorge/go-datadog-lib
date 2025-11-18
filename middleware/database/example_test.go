package database_test

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func Example() {
	ctx := context.Background()

	dsn := "example.com/users"
	db, err := ddDatabase.RegisterDriverAndOpen("mysql", mysqlDriver.MySQLDriver{}, dsn)
	if err != nil {
		panic(err)
	}

	span, ctx := tracer.StartSpanFromContext(ctx, "http.request")
	defer span.Finish()
	rows, err := db.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		span.Finish(tracer.WithError(err))
		panic(err)
	}
	println(rows)
}
