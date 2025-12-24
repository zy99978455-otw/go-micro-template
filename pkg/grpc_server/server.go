package grpc_server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zy99978455-otw/go-micro-template/pkg/global"
)

// RegisterFn 是一个回调函数类型
// 业务层通过这个函数，把自己的服务注册到 grpcServer 上
type RegisterFn func(server *grpc.Server)

// Run 启动通用的 gRPC 服务
// port: 端口号
// register: 业务层的注册回调（把业务逻辑传进来）
func Run(port int, register RegisterFn) (*grpc.Server, error) {
	// 1. 监听端口
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("gRPC 监听端口失败: %w", err)
	}

	// 2. 创建 gRPC 服务器实例
	// 🔥 框架核心价值：在这里统一添加拦截器（中间件）
	// 比如：Recovery（防崩溃）、Logger（日志）、Tracer（链路追踪）
	// 暂时先裸奔，后面可以在这里加 opts
	server := grpc.NewServer()

	// 3. 调用回调函数，注册业务服务
	// 框架层根本不知道你在注册什么，只管执行这个函数
	if register != nil {
		register(server)
	}

	// 4. 开启 gRPC 反射 (Reflection)
	// 这样可以用 grpcui 等工具直接调试接口，非常方便
	reflection.Register(server)

	// 5. 启动服务 (在一个新的 goroutine 中启动，避免阻塞主线程)
	go func() {
		global.Log.Infof("🚀 gRPC Server is starting on %s", addr)
		if err := server.Serve(lis); err != nil {
			global.Log.Errorf("gRPC Server 异常退出: %v", err)
		}
	}()

	return server, nil
}