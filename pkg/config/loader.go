package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// NewConfig 加载配置并返回实例
func NewConfig(path string) (*AppConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

    // 🔥 暂时移除热重载，避免循环引用
	// v.WatchConfig()
	// v.OnConfigChange(...) 

	var conf AppConfig
	if err := v.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	// 🔥 删掉了 global.AppConfig = &conf

	fmt.Printf("✅ 配置加载成功! App Name: %s, Port: %d\n", conf.Server.Name, conf.Server.Port)
	if len(conf.Chains) > 0 {
		fmt.Printf(">>> 监测到 Web3 配置: 已加载 %d 条链信息 (ChainID: %d)\n", len(conf.Chains), conf.Chains[0].ChainID)
	}

	return &conf, nil
}
