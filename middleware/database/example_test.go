package database_test

import (
	"context"

	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
