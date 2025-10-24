package logger

import "testing"

func Test_Logger(t *testing.T) {
	// if no arguments, create console logger
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("Failed to initialize : %v", err)
	}

	logger.Info("This is an info message")
	logger.Error("This is an error message")
	logger.Warn("This is a warning message")

	// if a file path is provided, create file logger
	loggerWithFile, err := NewLogger(".test_log.log")
	if err != nil {
		t.Fatalf("Failed to initialize : %v", err)
	}

	loggerWithFile.Info("This is an info message")
	loggerWithFile.Error("This is an error message")
	loggerWithFile.Warn("This is a warning message")

	// logger with rotation
	loggerWithRoatation, err := NewLogger(LoggerConfig{
		Filename:   ".test_rotation.log",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     28,
		Compress:   false,
	})
	if err != nil {
		t.Fatalf("Logger with rotation creation failed: %v", err)
	}

	for i := 0; i < 50000; i++ {
		loggerWithRoatation.Info("test log output")
	}

	t.Log("Logger test completed successfully")
}
