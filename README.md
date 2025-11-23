# uslogs

`uslogs` is a high-performance, unstructured logging library for Go, built on top of the standard library's `log/slog`. 
It is designed to be a drop-in replacement for `slog.TextHandler` with more place for customization, keeping performance in mind.

## Features

*   **High Performance**: Optimized for speed and low memory allocations.
*   **Drop-in Replacement**: Implements `slog.Handler` interface, making it compatible with existing `slog` setups.
*   **Plain Text Output**: Produces unstructured text logs by default.
*   **Flexible Configuration**: Supports customizable time formats, source code location logging, and log levels.
*   **Efficient Writers**: Includes specialized writers (currently `AsyncWriter`) to handle output efficiently.

## Installation

```bash
go get github.com/Drathveloper/uslogs
```

## Usage
Using `uslogs` is similar to using the standard `slog` package.
### Basic Example

``` go
package main

import (
    "log/slog"
    "os"
    
    "github.com/Drathveloper/uslogs"
)

func main() {
    // Create a new handler
    handler := uslogs.NewUnstructuredHandler(uslog.WithWriter(os.Stdout))

    // Create a logger with the handler
    logger := slog.New(handler)

    // Set as default logger (optional)
    slog.SetDefault(logger)

    // Log something
    logger.Info("Hello from uslogs!", "status", "ok", "count", 42)
}
```

### Configuration Options
You can customize the handler using built-in uslogs.LogWriterOption that includes:
*   `WithWriter`: Sets the writer to be used for output. Defaults to `os.Stdout`.
*   `WithTimestamp`: Sets if timestamps should be included in the output. Defaults to `false`.
*   `WithLevel`: Sets the minimum log level to be logged. Defaults to `slog.InfoLevel`.
*   `WithSeparator`: Sets the separator between fields. Defaults to `' '`
*   `WithMaskedFields`: Sets the attribute fields that should be masked in the output. Defaults to not masked fields.
*   `WithResponsivePool`: Allows the usage of multiple buffer pools to reduce memory allocations. Defaults to `false`.

## Benchmarks
uslogs is designed to be as fast as the standard library's text handler but more configurable and with support for asynchrony.
(You can run the included benchmark tests to verify performance on your machine)

``` bash
go test -bench=. -benchmem
```

## License
This project is licensed under the MIT License.