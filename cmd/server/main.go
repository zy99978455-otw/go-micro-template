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

	// 1. 引入业务层和数据层
	"github.com/zy99978455-otw/go-micro-template/internal/data"
	"github.com/zy99978455-otw/go-micro-template/internal/server"

	// 2. 引入基础设施层
	"github.com/zy99978455-otw/go-micro-template/pkg/bootstrap"
	"github.com/zy99978455-otw/go-micro-template/pkg/database"
	"github.com/zy99978455-otw/go-micro-template/pkg/global"
	"github.com/zy99978455-otw/go-micro-template/pkg/register"
)

func main() {
	// ================= 1. 初始化配置 =================
	// 以前是在 InitComponents 里做的，现在我们要显式做
	// 优先读取 config-local.yaml，没有则读取 config-debug.yaml
	configPath := "configs/config-local.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "configs/config-debug.yaml"
	}
	
	conf, err := bootstrap.NewConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// ================= 2. 初始化日志 =================
	bootstrap.InitLogger() // 这个暂时还没改，还是依赖 global.AppConfig，没问题

	// ================= 3. 初始化基础设施 (DB, Redis) =================
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
	
	// 先初始化 RPC Manager
	rpcMgr := data.NewRPCManager(conf)

	// 然后注入到 Data 层
	dataModule, cleanupData, err := data.NewData(db, rdb, rpcMgr)
	if err != nil {
		global.Log.Fatalf("Data 层初始化失败: %v", err)
	}
	defer cleanupData()

	fmt.Println("------------------------------------------------")

	// ================= 🔥 验证 RPC Manager 是否工作 =================
	targetChainID := int64(1)
	client, err := dataModule.GetRPCClient(targetChainID)

	if err != nil {
		global.Log.Errorf("❌ [验证失败] 无法获取 ChainID %d 的客户端: %v", targetChainID, err)
	} else {
		height, _ := client.BlockNumber(context.Background())
		global.Log.Infof("✅ [验证成功] 通过 RPCManager 拿到了客户端! ChainID: %d, 当前高度: %d", targetChainID, height)
	}

	fmt.Println("------------------------------------------------")

	httpPort := global.AppConfig.Server.Port
	fmt.Printf("\n🔥🔥🔥 HTTP服务启动！端口:%d 🔥🔥🔥\n\n", httpPort)

	// ================= 5. 启动 HTTP 服务 =================
	if global.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 调用 Server 层进行组装
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

	// ================= 6. 服务注册与优雅退出 =================
	registerToConsul(httpPort)

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

	global.Log.Info("👋 服务退出完成 (Bye!)")
}

// 注册函数保持不变...
func registerToConsul(httpPort int) {
    // ... (内容不变)
    consulReg, err := register.NewConsulRegister()
	if err == nil && consulReg != nil {
		serviceID := fmt.Sprintf("%s-%d", global.AppConfig.Server.Name, httpPort)
		registerErr := consulReg.RegisterService(
			global.AppConfig.Server.Name,
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
