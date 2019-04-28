package log

import (
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger

func init() {
	lg, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Logger init failed, %v", err)
	}
	Logger = lg
}
