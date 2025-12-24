package global

import (
	"go.uber.org/zap"
)

// 只保留日志，其他的都去掉了
var (
	Log *zap.SugaredLogger
)
