package logger

import "go.uber.org/zap"

// New membuat logger zap.
// Development memakai format yang lebih ramah dibaca,
// sedangkan production memakai structured JSON log.
func New(appEnv string) (*zap.Logger, error) {
	if appEnv == "production" {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}
