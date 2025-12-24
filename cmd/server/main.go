package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	// 引入各层
	"github.com/zy99978455-otw/go-micro-template/internal/data"
	"github.com/zy99978455-otw/go-micro-template/internal/server"
	"github.com/zy99978455-otw/go-micro-template/pkg/config"
	"github.com/zy99978455-otw/go-micro-template/pkg/database"
	"github.com/zy99978455-otw/go-micro-template/pkg/global"
	"github.com/zy99978455-otw/go-micro-template/pkg/logger"
	"github.com/zy99978455-otw/go-micro-template/pkg/register"
)

func main() {
	// ================= 1. 初始化配置 (不再依赖 global) =================
	configPath := "configs/config-local.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "configs/config-debug.yaml"
	}
	
	// 注意：如果你的 NewConfig 在 pkg/config/loader.go 里，这里包名可能是 config
	// 如果在 pkg/bootstrap/config.go 里，包名可能是 bootstrap
	// 请根据你实际的包名修改调用
	conf, err := config.NewConfig(configPath) 
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// ================= 2. 初始化日志 =================
	logger.InitLogger() // 暂时保持原样

	// ================= 3. 初始化基础设施 (显式传参 conf) =================
	// MySQL
	db, cleanupDB, err := database.NewMySQLClient(conf)
	if err != nil {
		global.Log.Errorf("MySQL Init Failed: %v", err)
	}
	defer cleanupDB()

	// Redis
	rdb, cleanupRedis, err := database.NewRedisClient(conf)
	if err != nil {
		global.Log.Errorf("Redis Init Failed: %v", err)
	}
	defer cleanupRedis()
	// ================= 4. 初始化 Data 层 (依赖注入) =================
	
	// 4.1 先初始化 RPC Manager (传入 conf)
	rpcMgr := data.NewRPCManager(conf)

	// 4.2 然后注入到 Data 层
	dataModule, cleanupData, err := data.NewData(db, rdb, rpcMgr)
	if err != nil {
		global.Log.Fatalf("Data 层初始化失败: %v", err)
	}
	defer cleanupData()

	// 验证 RPC
	fmt.Println("------------------------------------------------")
	targetChainID := int64(1)
	client, err := dataModule.GetRPCClient(targetChainID)
	if err != nil {
		global.Log.Errorf("❌ [验证失败] 无法获取 ChainID %d: %v", targetChainID, err)
	} else {
		height, _ := client.BlockNumber(context.Background())
		global.Log.Infof("✅ [验证成功] RPC工作正常! ChainID: %d, Height: %d", targetChainID, height)
	}
	fmt.Println("------------------------------------------------")

	// ================= 5. 启动 HTTP 服务 =================
	httpPort := conf.Server.Port // 使用 conf，不用 global
	fmt.Printf("\n🔥🔥🔥 HTTP服务启动！端口:%d 🔥🔥🔥\n\n", httpPort)

	if conf.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 组装 Server
	r := server.NewHTTPServer(dataModule)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: r,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Fatalf("HTTP Server 启动失败: %v", err)
		}
	}()

	// ================= 6. 服务注册 (Consul) =================
	// 🔥 传入 conf
	registerToConsul(httpPort, conf)

	// ================= 7. 优雅停机 =================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit 

	global.Log.Info("正在关闭服务 (Shutting down)...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := httpSrv.Shutdown(ctx); err != nil {
		global.Log.Error("HTTP 强制关闭:", err)
	} else {
		global.Log.Info("✅ [HTTP] 服务已停止")
	}
	global.Log.Info("👋 服务退出完成")
}
// registerToConsul 辅助函数
// 🔥 修改：接收 conf *config.AppConfig 参数
func registerToConsul(httpPort int, conf *config.AppConfig) {
	
	// 传入 conf 初始化
	consulReg, err := register.NewConsulRegister(conf)

	if err == nil && consulReg != nil {
		// 使用 conf 获取服务名
		serviceID := fmt.Sprintf("%s-%d", conf.Server.Name, httpPort)

		registerErr := consulReg.RegisterService(
			conf.Server.Name,
			serviceID,
			httpPort,
			[]string{"http", "web3"},
		)

		if registerErr != nil {
			global.Log.Warnf("Consul 注册失败: %v", registerErr)
		} else {
			global.Log.Infof("✅ 服务已注册到 Consul (ID: %s)", serviceID)
		}
	}
}
