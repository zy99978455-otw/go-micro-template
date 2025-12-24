package global

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/zy99978455-otw/go-micro-template/pkg/config"
)

var (
	AppConfig *config.AppConfig
	Log       *zap.SugaredLogger

	// Web2: MySQL 全局连接
	DB *gorm.DB

	Redis *redis.Client

	// Web3: 多链客户端容器
	// Key: ChainID (e.g., 1, 56) -> Value: Client
	EthClients map[int64]*ethclient.Client
)

// GetEthClient 辅助函数：快速获取指定链的客户端
func GetEthClient(chainID int64) *ethclient.Client {
	if client, ok := EthClients[chainID]; ok {
		return client
	}
	return nil
}
