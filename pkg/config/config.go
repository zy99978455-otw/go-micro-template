package config

// ServerConfig 服务端通用配置
type ServerConfig struct {
	Name        string      `mapstructure:"name" json:"name"`
	Mode        string      `mapstructure:"mode" json:"mode"`
	Port        int         `mapstructure:"port" json:"port"`
	Version     string      `mapstructure:"version" json:"version"`
	
	// 手动指定注册 IP (解决 Docker 网络隔离问题)
	RegisterIP  string      `mapstructure:"register_ip" json:"register_ip"` 

	JwtInfo     JwtConfig   `mapstructure:"jwt" json:"jwt"`
	ConsulInfo  ConsulConfig `mapstructure:"consul" json:"consul"`
}
// ================= Web2 基础设施 =================

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"name" json:"name"` // 数据库名
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	MaxIdle  int    `mapstructure:"max_idle" json:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open" json:"max_open"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"signing_key" json:"signing_key"`
	Expire     int64  `mapstructure:"expire" json:"expire"` // 过期时间(秒)
}

// ================= Web3 基础设施 (关键差异点) =================

// ChainConfig 通用链配置 (支持多链)
type ChainConfig struct {
	ChainID   int64  `mapstructure:"chain_id" json:"chain_id"`     // 链ID
	ChainName string `mapstructure:"chain_name" json:"chain_name"` // 链名称 (e.g. "eth_mainnet", "bsc_testnet")
	RpcUrl    string `mapstructure:"rpc_url" json:"rpc_url"`       // HTTP RPC 地址
	WssUrl    string `mapstructure:"wss_url" json:"wss_url"`       // WebSocket 地址 (监听事件用)
	ApiKey    string `mapstructure:"api_key" json:"api_key"`       // 如果用 Infura/Alchemy 需要 Key
}

// ContractConfig 智能合约配置 (DApp 常用)
type ContractConfig struct {
	Name    string `mapstructure:"name" json:"name"`       // 合约别名
	Address string `mapstructure:"address" json:"address"` // 合约地址
	AbiJson string `mapstructure:"abi_json" json:"abi_json"` // ABI 内容或路径
}

// ================= 总入口 =================

type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server" json:"server"`
	Mysql    MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	Redis    RedisConfig    `mapstructure:"redis" json:"redis"`
	
	// Web3 特有：支持配置多个链 (例如同时监听 ETH 和 BSC)
	Chains   []ChainConfig  `mapstructure:"chains" json:"chains"`
	Contracts []ContractConfig `mapstructure:"contracts" json:"contracts"`
}