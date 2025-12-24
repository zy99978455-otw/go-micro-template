package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zy99978455-otw/go-micro-template/internal/biz"
	"github.com/zy99978455-otw/go-micro-template/internal/data"
)

// NewHTTPServer 初始化 HTTP 服务器
func NewHTTPServer(dataModule *data.Data) *gin.Engine {
	
	// 1. 组装 (Wiring)
	// data -> biz
	chainRepo := data.NewChainRepo(dataModule)
	// biz -> handler (直接在 server 包内实例化 handler)
	chainUseCase := biz.NewChainUsecase(chainRepo)
	chainHandler := NewChainHandler(chainUseCase)

	// 2. 路由
	r := gin.Default()

	// 🔥健康检查接口 
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "UP"})
    })
	
	v1 := r.Group("/api/v1")
	{
		web3 := v1.Group("/web3")
		{
			web3.GET("/block", chainHandler.GetBlock)
		}
	}
	
	return r
}
