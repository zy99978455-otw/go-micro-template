package register

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/zy99978455-otw/go-micro-template/pkg/common"
	"github.com/zy99978455-otw/go-micro-template/pkg/config" // 引入 config 包
	"github.com/zy99978455-otw/go-micro-template/pkg/global" // 仅用于日志
)

type ConsulRegister struct {
	Client *api.Client
	Config *config.AppConfig // 保存配置，供后续方法使用
}

// NewConsulRegister 创建 Consul 客户端
// 接收 cfg 参数，不再读 global
func NewConsulRegister(cfg *config.AppConfig) (*ConsulRegister, error) {
	
	consulInfo := cfg.Server.ConsulInfo // 使用传入的 cfg

	// 检查配置
	if consulInfo.Host == "" {
		return nil, fmt.Errorf("consul host 未配置")
	}

	apiConfig := api.DefaultConfig()
	apiConfig.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	// 设置超时
	if apiConfig.HttpClient != nil {
		apiConfig.HttpClient.Timeout = 10 * time.Second
	}

	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, fmt.Errorf("创建 Consul 客户端失败: %w", err)
	}

	return &ConsulRegister{
		Client: client,
		Config: cfg, // 保存起来
	}, nil
}

// RegisterService 注册服务
func (r *ConsulRegister) RegisterService(name, id string, port int, tags []string, retryTimes ...int) error {
	
	maxRetry := 5
	if len(retryTimes) > 0 && retryTimes[0] > 0 {
		maxRetry = retryTimes[0]
	}

	// 使用 r.Config 获取注册 IP
	registerAddr := r.getRegisterIP(port) 

	var err error
	for attempt := 1; attempt <= maxRetry; attempt++ {
		if attempt > 1 {
			sleepTime := time.Duration(attempt*2) * time.Second
			global.Log.Warnf("Consul 注册重试第 %d/%d 次，%v 后重试...", attempt, maxRetry, sleepTime)
			time.Sleep(sleepTime)
		}

		registration := &api.AgentServiceRegistration{
			Name:    name,
			ID:      id,
			Port:    port,
			Tags:    tags,
			Address: registerAddr,
			Check: &api.AgentServiceCheck{
				HTTP:                           fmt.Sprintf("http://%s:%d/health", registerAddr, port),
				Method:                         "GET",
				Timeout:                        "5s",
				Interval:                       "10s",
				DeregisterCriticalServiceAfter: "60s",
			},
		}

		err = r.Client.Agent().ServiceRegister(registration)
		if err == nil {
			global.Log.Info("✅ Consul 服务注册成功")
			return nil
		}
		global.Log.Warnf("Consul 服务注册失败（第 %d 次）: %v", attempt, err)
	}

	return fmt.Errorf("Consul 服务注册最终失败: %w", err)
}

// getRegisterIP 变成内部方法，使用 r.Config
func (r *ConsulRegister) getRegisterIP(port int) string {
	registerAddr := r.Config.Server.RegisterIP // 从保存的配置里拿

	if registerAddr == "" {
		ip, err := common.GetOutboundIP()
		if err != nil {
			global.Log.Errorf("自动获取本机 IP 失败: %v，使用 127.0.0.1 兜底", err)
			registerAddr = "127.0.0.1"
		} else {
			registerAddr = ip
		}
		global.Log.Infof(">>> 自动探测本机 IP: %s", registerAddr)
	} else {
		global.Log.Infof(">>> 使用配置文件指定的 IP: %s", registerAddr)
	}

	global.Log.Infof(">>> 准备注册服务到 Consul, 地址: %s:%d", registerAddr, port)
	return registerAddr
}

// DeregisterService 注销服务
func (r *ConsulRegister) DeregisterService(serviceID string) error {
	err := r.Client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		global.Log.Errorf("Consul 服务注销失败: %v", err)
		return err
	}
	global.Log.Info("✅ 服务已从 Consul 注销")
	return nil
}
