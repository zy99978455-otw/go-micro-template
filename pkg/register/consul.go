package register

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/zy99978455-otw/go-micro-template/pkg/common"
	"github.com/zy99978455-otw/go-micro-template/pkg/global"
)

type ConsulRegister struct {
	Client *api.Client
}

// NewConsulRegister 创建 Consul 客户端
func NewConsulRegister() (*ConsulRegister, error) {
    // 1. 获取配置
    consulInfo := global.AppConfig.Server.ConsulInfo

    // ❌ 删掉: if consulInfo == nil { ... } (这是导致报错的原因)
    
    // ✅ 保留: 检查 Host 是否为空即可判断配置是否存在
    if consulInfo.Host == "" {
        return nil, fmt.Errorf("consul host 未配置 (请检查 config.yaml 中 consul 是否缩进在 server 下面)")
    }

    cfg := api.DefaultConfig()
    cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

    // 防御性代码：设置超时
    if cfg.HttpClient == nil {
        cfg.HttpClient = api.DefaultConfig().HttpClient
    }
    // 注意：如果 api.DefaultConfig().HttpClient 也是 nil，这里可能会崩
    // 更加稳妥的写法：
    if cfg.HttpClient != nil {
        cfg.HttpClient.Timeout = 10 * time.Second
    }

    client, err := api.NewClient(cfg)
    if err != nil {
        return nil, fmt.Errorf("创建 Consul 客户端失败: %w", err)
    }

    return &ConsulRegister{Client: client}, nil
}

// RegisterService 注册服务（增加可选重试次数参数）
func (r *ConsulRegister) RegisterService(name, id string, port int, tags []string, retryTimes ...int) error {
	// 可变参数支持：不传或传 0 时默认重试 5 次
	maxRetry := 5
	if len(retryTimes) > 0 {
		maxRetry = retryTimes[0]
		if maxRetry <= 0 {
			maxRetry = 1
		}
	}

	registerAddr := getRegisterIP(port)

	var err error
	for attempt := 1; attempt <= maxRetry; attempt++ {
		if attempt > 1 {
			sleepTime := time.Duration(attempt*2) * time.Second // 指数退避
			global.Log.Warnf("Consul 注册重试第 %d/%d 次，%v 后重试...", attempt, maxRetry, sleepTime)
			time.Sleep(sleepTime)
		}

		registration := &api.AgentServiceRegistration{
			Name:    name,
			ID:      id,
			Port:    port,
			Tags:    tags,
			Address: registerAddr,
		}

		// 🔥 推荐使用 HTTP 检查（因为你有 /health 接口）
		registration.Check = &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", registerAddr, port),
			Method:                         "GET",
			Timeout:                        "5s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "60s",
		}

		// 如果你暂时不想用 HTTP 检查，想用 TCP 检查，改成下面这块即可：
		// registration.Check = &api.AgentServiceCheck{
		// 	TCP:                            fmt.Sprintf("%s:%d", registerAddr, port),
		// 	Timeout:                        "5s",
		// 	Interval:                       "10s",
		// 	DeregisterCriticalServiceAfter: "60s",
		// }

		err = r.Client.Agent().ServiceRegister(registration)
		if err == nil {
			global.Log.Info("✅ Consul 服务注册成功")
			return nil
		}
		global.Log.Warnf("Consul 服务注册失败（第 %d 次）: %v", attempt, err)
	}

	return fmt.Errorf("Consul 服务注册最终失败（已重试 %d 次）: %w", maxRetry, err)
}

// 提取 IP 获取逻辑
func getRegisterIP(port int) string {
	registerAddr := global.AppConfig.Server.RegisterIP
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