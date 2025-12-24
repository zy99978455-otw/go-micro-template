package bootstrap

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/zy99978455-otw/go-micro-template/pkg/config" 
	"github.com/zy99978455-otw/go-micro-template/pkg/global"
)

// NewConfig 加载配置并返回实例 (供 Wire 使用)
// path: 配置文件路径
func NewConfig(path string) (*config.AppConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// 1. 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	// 2. 监听配置文件变化 (热加载)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件被修改:", e.Name)
		// 重新解析到全局变量
		if err := v.Unmarshal(&global.AppConfig); err != nil {
			fmt.Println("配置重载失败:", err)
		}
	})

	var conf config.AppConfig
	// 3. 解析配置
	if err := v.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	// 🔥 兼容旧代码：赋值给全局变量
	global.AppConfig = &conf

	fmt.Printf("✅ 配置加载成功! App Name: %s, Port: %d\n", conf.Server.Name, conf.Server.Port)
	if len(conf.Chains) > 0 {
		fmt.Printf(">>> 监测到 Web3 配置: 已加载 %d 条链信息 (ChainID: %d)\n", len(conf.Chains), conf.Chains[0].ChainID)
	}

	return &conf, nil
}

// InitConfig 旧的初始化函数 (为了保持兼容性，让它调用 NewConfig)
func InitConfig() {
    // 这里硬编码路径，和之前一样
	_, err := NewConfig("configs/config-local.yaml")
    if err != nil {
        panic(err)
    }
}
