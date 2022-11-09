# go-logger

This Go package is used to offer a unified logging interface among projects.

## Import

```
import "github.com/coopnorge/go-logger"
```

## Default behavior

By default, all logs will include:
- Log level
- Full path to file which called the logger, and line number
- Signature of the function that called the logger
- Timestamp of the log entry

Example log entry with default settings:

```go
logger.Warn("something went wrong")
// Output:
// {"file":"/Users/anonymous/Projects/my-project/main.go:7","function":"main.main","level":"warning","msg":"something went wrong","time":"2022-02-17T15:04:06+01:00"}
```

## Example usage

See [logger_examples_test.go](logger_examples_test.go) for more examples.

<details>
<summary>Using global logger</summary>

```go
package main

import "github.com/coopnorge/go-logger"

func main() {
    logger.Info("this won't be logged because the default log level is higher than info")
    logger.Warn("but this will be logged")
    // Output:
    // {"level":"warning","msg":"but this will be logged","time":"2022-02-17T11:01:28+01:00"}
}
```

</details>

<details>
<summary>Setting log level</summary>

```go
package main

import "github.com/coopnorge/go-logger"

func main() {
    // global logger
    logger.Info("this won't be logged because the default log level is higher than info")
    logger.ConfigureGlobalLogger(logger.WithLevel(logger.LevelInfo))
    logger.Info("now this will be logged")
    // Output:
    // {"level":"info","msg":"now this will be logged","time":"2022-02-17T10:54:54+01:00"}

    // logger instance
    prodLogger := logger.New(logger.WithLevel(logger.LevelWarn))
    prodLogger.Info("this won't be logged because prodLogger's level is set to Warn...")
    prodLogger.Error("...but this will, because Error >= Warn")
    // Output:
    // {"level":"error","msg":"...but this will, because Error \u003e= Warn","time":"2022-02-17T10:54:54+01:00"}

    debugLogger := logger.New(logger.WithLevel(logger.LevelDebug))
    debugLogger.Debug("this logger will log anything as Debug is the lowest available level")
    debugLogger.Warn("and this will be logged too")
    // Output:
    // {"level":"debug","msg":"this logger will log anything as Debug is the lowest available level","time":"2022-02-17T10:54:54+01:00"}
    // {"level":"warning","msg":"and this will be logged too","time":"2022-02-17T10:54:54+01:00"}
}
```

</details>