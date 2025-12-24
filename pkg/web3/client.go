package web3

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zy99978455-otw/go-micro-template/pkg/global"
)

// InitWeb3Clients 初始化所有链的 RPC 连接
func InitWeb3Clients() {
	chains := global.AppConfig.Chains
	global.EthClients = make(map[int64]*ethclient.Client)

	// 🔥 智能开关：如果没有配置任何链，直接跳过
	if len(chains) == 0 {
		global.Log.Info(">>> [Web3] Chain 配置为空，跳过初始化 (当前可能为纯 Web2 模式)")
		return
	}

	for _, chain := range chains {
		// 建立 RPC 连接 (设置 10秒超时，防止卡死)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := ethclient.DialContext(ctx, chain.RpcUrl)
		cancel()

		if err != nil {
			global.Log.Errorf("链 [%s] 连接失败: %v", chain.ChainName, err)
			continue
		}

		// 简单的连通性测试 (获取 ChainID)
		cid, err := client.ChainID(context.Background())
		if err != nil {
			global.Log.Errorf("链 [%s] 通信失败 (ChainID获取失败): %v", chain.ChainName, err)
			continue
		}

		// 存入全局 Map
		global.EthClients[chain.ChainID] = client
		global.Log.Infof(">>> [Web3] 节点连接成功: %s (ChainID: %d)", chain.ChainName, cid)
	}
}