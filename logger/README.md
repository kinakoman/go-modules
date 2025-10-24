Logger package
==============

The `logger` package provides a small wrapper around `zap` with optional file output and rotation (via `lumberjack`). It aims to be simple to use for console logging, writing to a file, or creating a rotating log file.

Key types and functions
-----------------------

- `Logger`: thin wrapper around `*zap.Logger` with convenience methods:
  - `Info(msg string)`
  - `Error(msg string)`
  - `Warn(msg string)`
  - `Debug(msg string)`
- `LoggerConfig`: configuration used for file rotation (used by `NewLogger` when passed)
  - `Filename string`
  - `MaxSize int` (MB)
  - `MaxBackups int`
  - `MaxAge int` (days)
  - `Compress bool`
- `NewLogger(args ...interface{}) (*Logger, error)` – constructor with flexible arguments:
  - `NewLogger()` — console logger only
  - `NewLogger("path/to/file.log")` — console + file logger (append/create)
  - `NewLogger(LoggerConfig{...})` — console + lumberjack-based rotating file logger

Usage
-----

Import the package (module path):

    import "github.com/kinakoman/go-modules/logger"

Example — console logger:

    l, err := logger.NewLogger()
    if err != nil {
        // handle error
    }
    l.Info("starting application")

Example — simple file logger:

    l, err := logger.NewLogger("logs/app.log")
    if err != nil {
        // handle error
    }
    l.Warn("something to watch")

Example — rotating log file using LoggerConfig:

    cfg := logger.LoggerConfig{
        Filename:   "logs/app_rot.log",
        MaxSize:    10, // MB
        MaxBackups: 3,
        MaxAge:     30, // days
        Compress:   true,
    }
    l, err := logger.NewLogger(cfg)
    if err != nil {
        // handle error
    }
    l.Error("an error occurred")

Notes
-----

- `NewLogger` accepts either no arguments, a single file path string, or a `LoggerConfig` value. Passing both a string and a `LoggerConfig` will prefer the `LoggerConfig` behavior.
- The logger is configured to include short caller information (file:line) and timestamps. Internally it uses `zap` and `lumberjack` for rotation.
- The package sets a sensible default encoder and Info-level output. Adjust the code if you need different levels or formats.

Installing
----------

If you're using modules, add the import to your code and run `go get`:

    go get github.com/kinakoman/go-modules/logger@latest

Dependencies
------------

- go.uber.org/zap
- gopkg.in/natefinch/lumberjack.v2

License
-------

This project is **not under an open-source license**,  
but you are welcome to **use, modify, and share** the code for any purpose.  
Please keep the author credit somewhere in your project.  
No warranty or liability is provided.