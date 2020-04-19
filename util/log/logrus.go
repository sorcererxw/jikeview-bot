package log

import (
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
}

// 获取 logrus 实例
func Logger() *logrus.Logger {
	return logger
}
