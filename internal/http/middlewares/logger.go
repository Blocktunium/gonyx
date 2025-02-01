package middlewares

import (
	"github.com/Blocktunium/gonyx/internal/logger"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func ZapLogger() gin.HandlerFunc {
	logge, err := logger.GetManager().GetLogger()
	if err != nil {
		return nil
	}
	logger1 := logge.(*logger.ZapWrapper).Instance()

	return ginzap.Ginzap(logger1, time.RFC3339, true)
}

func ZapRecoveryLogger() gin.HandlerFunc {
	logge, err := logger.GetManager().GetLogger()
	if err != nil {
		return nil
	}
	logger1 := logge.(*logger.ZapWrapper).Instance()

	return ginzap.RecoveryWithZap(logger1, true)
}

func LogMeLogger() gin.HandlerFunc {
	logge, err := logger.GetManager().GetLogger()
	if err != nil {
		return nil
	}
	logger1 := logge.(*logger.LogMeWrapper).Instance()

	return gin.LoggerWithWriter(logger1.Writer())
}

func LogMeRecoveryLogger() gin.HandlerFunc {
	logge, err := logger.GetManager().GetLogger()
	if err != nil {
		return nil
	}
	logger1 := logge.(*logger.LogMeWrapper).Instance()

	return gin.RecoveryWithWriter(logger1.Writer())
}
