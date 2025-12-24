package data

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zy99978455-otw/go-micro-template/pkg/config" // 引入 config 包
	"github.com/zy99978455-otw/go-micro-template/pkg/global" // 仅用于日志
)

// Node 代表一个具体的 RPC 节点
type Node struct {
	URL         string
	ChainID     int64
	Client      *ethclient.Client
	
	IsHealthy   bool
	Latency     time.Duration
	BlockHeight uint64
	ErrorCount  int
	
	mu          sync.RWMutex
}

// RPCManager 管理多链的所有节点
type RPCManager struct {
	chainNodes map[int64][]*Node
	mu sync.RWMutex
}

// ================= 2. 初始化逻辑 =================

// NewRPCManager 根据传入的配置初始化管理器
func NewRPCManager(cfg *config.AppConfig) *RPCManager {
	mgr := &RPCManager{
		chainNodes: make(map[int64][]*Node),
	}

	// 1. 遍历配置，初始化连接
	// 🔥 使用传入的 cfg，不再使用 global.AppConfig
	if cfg != nil && len(cfg.Chains) > 0 {
		for _, chainConf := range cfg.Chains {
			
			// 尝试初始连接
			client, err := ethclient.Dial(chainConf.RpcUrl)
			isHealthy := false
			if err == nil {
				isHealthy = true 
			} else {
				if global.Log != nil {
					global.Log.Warnf("⚠️ [RPC] Init failed for chain %d (%s): %v", chainConf.ChainID, chainConf.RpcUrl, err)
				} else {
					fmt.Printf("⚠️ [RPC] Init failed for chain %d (%s): %v\n", chainConf.ChainID, chainConf.RpcUrl, err)
				}
			}

			node := &Node{
				URL:       chainConf.RpcUrl,
				ChainID:   chainConf.ChainID,
				Client:    client,
				IsHealthy: isHealthy,
			}

			mgr.chainNodes[chainConf.ChainID] = append(mgr.chainNodes[chainConf.ChainID], node)
            
            if global.Log != nil {
                global.Log.Infof("✅ [RPC] Added node for chain %d: %s", chainConf.ChainID, chainConf.RpcUrl)
            }
		}
	}

	// 2. 启动后台健康检查
	go mgr.startHealthCheckLoop()

	return mgr
}

// ================= 3. 健康检查核心逻辑 =================

func (m *RPCManager) startHealthCheckLoop() {
	// 每 30 秒检查一次
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.checkAllNodes()
	}
}

func (m *RPCManager) checkAllNodes() {
	m.mu.RLock()
	// 复制一份节点列表，避免在检查时长时间持有锁
	// 这里其实可以直接遍历，因为 map 只有初始化时才写，后面基本只读。
	// 但为了严谨，我们还是在锁内只做简单的遍历。
	
    var allNodes []*Node
    for _, nodes := range m.chainNodes {
        allNodes = append(allNodes, nodes...)
    }
	m.mu.RUnlock()

	var wg sync.WaitGroup
	for _, node := range allNodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			m.checkOneNode(n)
		}(node)
	}
	wg.Wait()
}

func (m *RPCManager) checkOneNode(n *Node) {
	// 设置 5 秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	
	// 如果 client 为空（初始化失败），尝试重连
	if n.Client == nil {
		client, err := ethclient.Dial(n.URL)
		if err != nil {
			m.markUnhealthy(n, err)
			return
		}
		n.Client = client
	}

	// 核心检查：获取区块高度
	// BlockNumber 返回的是 uint64
	height, err := n.Client.BlockNumber(ctx)
	latency := time.Since(start)

	if err != nil {
		m.markUnhealthy(n, err)
		return
	}

	// 标记为健康
	n.mu.Lock()
	n.IsHealthy = true
	n.Latency = latency
	n.BlockHeight = height
	n.ErrorCount = 0
	n.mu.Unlock()
}

func (m *RPCManager) markUnhealthy(n *Node, err error) {
	n.mu.Lock()
	n.IsHealthy = false
	n.ErrorCount++
	currentErrCount := n.ErrorCount
	n.mu.Unlock()
	
	// 只有连续错误多次才打印 Error 日志，避免刷屏
	if currentErrCount <= 3 && global.Log != nil {
		global.Log.Warnf("⚠️ [RPC] Node unhealthy: %s, Err: %v", n.URL, err)
	}
}

// ================= 4. 对外接口 =================

// GetClient 获取指定链的一个最佳节点
func (m *RPCManager) GetClient(chainID int64) (*ethclient.Client, error) {
	m.mu.RLock()
	nodes, ok := m.chainNodes[chainID]
	m.mu.RUnlock()

	if !ok || len(nodes) == 0 {
		return nil, fmt.Errorf("chain %d not configured", chainID)
	}

	// 简单的负载均衡策略：选第一个健康的
	// 进阶策略：可以遍历 nodes，找 latency 最小的
	for _, node := range nodes {
		node.mu.RLock()
		isHealthy := node.IsHealthy
        client := node.Client
		node.mu.RUnlock()

		if isHealthy && client != nil {
			return client, nil
		}
	}

	return nil, fmt.Errorf("no healthy node available for chain %d", chainID)
}