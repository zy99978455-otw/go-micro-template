package bootstrap

import (
	"go.uber.org/zap"
	"github.com/zy99978455-otw/go-micro-template/pkg/global"
)

// InitLogger 初始化日志
// 修改点：函数名改为 InitLogger (为了匹配 init.go)
func InitLogger() {
	// 这里先用开发模式的 Logger，生产环境以后再配置 Zap Core
	logger, _ := zap.NewDevelopment()
	global.Log = logger.Sugar()
	global.Log.Info("✅ 日志组件初始化完成")
}