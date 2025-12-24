package logger

import (
	"go.uber.org/zap"
	"github.com/zy99978455-otw/go-micro-template/pkg/global"
)

// InitLogger 初始化全局日志
func InitLogger() {
	// 暂时还是硬编码 Development 模式
	logger, _ := zap.NewDevelopment()
	
	// 赋值给全局变量，供全项目使用
	global.Log = logger.Sugar()
	global.Log.Info("✅ 日志组件初始化完成")
}
