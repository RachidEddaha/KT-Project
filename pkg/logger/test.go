package logger

import (
	"KTOnlinePlatform/pkg/configuration"
)

func InitializeForTest() {
	if internalLogger != nil {
		return
	}
	Initialize(configuration.ConfigLogger{
		LogLevel:    "trace",
		LogSaveFile: false,
	})
}
