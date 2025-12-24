package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zy99978455-otw/go-micro-template/internal/biz"
)

// 直接处理 HTTP 请求
type ChainHandler struct {
	uc *biz.ChainUsecase
}

// NewChainHandler 构造函数
func NewChainHandler(uc *biz.ChainUsecase) *ChainHandler {
	return &ChainHandler{uc: uc}
}

// GetBlock 处理 GET /api/v1/web3/block 请求
func (h *ChainHandler) GetBlock(c *gin.Context) {
	// 1. 解析参数
	chainIDStr := c.Query("chain_id")
	// 默认值为 1
	chainID, _ := strconv.ParseInt(chainIDStr, 10, 64)
	if chainID == 0 {
		chainID = 1 
	}

	// 2. 调用业务逻辑 (Biz)
	height, err := h.uc.GetCurrentHeight(c.Request.Context(), chainID)
	
	// 3. 处理错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"chain_id": chainID,
			"height":   height,
		},
	})
}
