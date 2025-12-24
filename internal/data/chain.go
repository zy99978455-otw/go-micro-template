package data

import (
	"context"
	"github.com/zy99978455-otw/go-micro-template/internal/biz" 
)

// 定义 Data 层的 ProviderSet


// chainRepo 是 biz.ChainRepo 的具体实现
type chainRepo struct {
	data *Data
}

// NewChainRepo 构造函数
func NewChainRepo(data *Data) biz.ChainRepo {
	return &chainRepo{
		data: data,
	}
}

// GetBlockHeight 实现接口方法
func (r *chainRepo) GetBlockHeight(ctx context.Context, chainID int64) (uint64, error) {
	// 1. 从 Data 层获取 RPC 客户端 (利用了我们的 RPC Manager)
	client, err := r.data.GetRPCClient(chainID)
	if err != nil {
		return 0, err
	}

	// 2. 调用 ethclient 的方法
	height, err := client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return height, nil
}
