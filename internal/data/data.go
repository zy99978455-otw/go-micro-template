package data

import (
	"github.com/google/wire" // 引入 wire
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zy99978455-otw/go-micro-template/pkg/global" // 仅用于日志
)

// 🔥 定义 ProviderSet，告诉 Wire data 层有哪些组件
var ProviderSet = wire.NewSet(NewData, NewRPCManager, NewChainRepo)

type Data struct {
	db         *gorm.DB
	redis      *redis.Client
	rpcManager *RPCManager
}

// NewData 显式接收依赖
// 参数 db, redis, rpcMgr 都会由 Wire 自动注入
func NewData(db *gorm.DB, rdb *redis.Client, rpcMgr *RPCManager) (*Data, func(), error) {
	d := &Data{
		db:         db,
		redis:      rdb,
		rpcManager: rpcMgr,
	}

	cleanup := func() {
		global.Log.Info("正在关闭 Data 层资源...")
	}

	return d, cleanup, nil
}

func (d *Data) GetRPCClient(chainID int64) (*ethclient.Client, error) {
	return d.rpcManager.GetClient(chainID)
}

func (d *Data) GetDB() *gorm.DB {
	return d.db
}
