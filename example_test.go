package coopdatadog_test

import (
	"context"

	coopdatadog "github.com/coopnorge/go-datadog-lib/v2"
)

func ExampleStart() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	stop, err := coopdatadog.Start(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		err := stop()
		if err != nil {
			panic(err)
		}
	}()

	// ...

	return nil
}
