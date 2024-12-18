package gorm_test

import (
	"context"

	ddDatabase "github.com/coopnorge/go-datadog-lib/v2/middleware/database"
	ddGorm "github.com/coopnorge/go-datadog-lib/v2/middleware/gorm"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct{}

func ExampleNewORM_standard() {
	ctx := context.Background()

	dsn := "example.com/users"
	gormDB, err := ddGorm.NewORM(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := &User{}
	tx := gormDB.WithContext(ctx).Select("*").First(user)

	println(tx)
}

func ExampleNewORM_withDriver() {
	ctx := context.Background()

	dsn := "example.com/users"
	db, err := ddDatabase.RegisterDriverAndOpen("mysql", mysqlDriver.MySQLDriver{}, dsn)
	if err != nil {
		panic(err)
	}

	gormDB, err := ddGorm.NewORM(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := &User{}
	tx := gormDB.WithContext(ctx).Select("*").First(user)

	println(tx)
}
