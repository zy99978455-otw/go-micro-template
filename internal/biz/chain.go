package biz

import (
	"context"
	"github.com/google/wire"
)

// 定义 Biz 层的 ProviderSet
var ProviderSet = wire.NewSet(NewChainUsecase)


// ChainUsecase 定义了与链交互的业务逻辑接口
type ChainUsecase struct {
	repo ChainRepo
}

// ChainRepo 定义了数据层必须实现的方法 (依赖倒置)
type ChainRepo interface {
	GetBlockHeight(ctx context.Context, chainID int64) (uint64, error)
	// 未来可以在这里加: GetBalance, SendTransaction ...
}

// NewChainUsecase 构造函数
func NewChainUsecase(repo ChainRepo) *ChainUsecase {
	return &ChainUsecase{repo: repo}
}

// GetCurrentHeight 业务方法：获取当前高度
func (uc *ChainUsecase) GetCurrentHeight(ctx context.Context, chainID int64) (uint64, error) {
	return uc.repo.GetBlockHeight(ctx, chainID)
}
