package chain

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Config 客户端的配置参数
type Config struct {
	// 基础设置
	ChainID        string                // 链 ID, node id
	NodeRPC        string                // 节点 RPC 地址 e.g. "http://127.0.0.1:26657"
	NodeGRPC       string                // 节点 gRPC 地址 e.g. "127.0.0.1:9090"
	HomeDir        string                // 数据目录
	KeyringBackend string                // "test" | "file" | "os"
	BroadcastMode  txtypes.BroadcastMode // 广播模式

	// 发送交易相关设置
	Fee      string
	GasLimit uint64

	// gRPC 限制/保活
	GRPCMaxRecvBytes int
	GRPCMaxSendBytes int
	KeepaliveTime    time.Duration
	KeepaliveTimeout time.Duration

	// 调试 / 其他
	SkipConfirmation bool

	// 自定义模块注册（可选，覆盖默认：auth + bank）
	// 用于自定义 encoding config（注册更多模块）
	RegisterModules func() moduletestutil.TestEncodingConfig
}

// DefaultConfig 获取默认配置
// modules yourmodule.AppModule{}
// (通过go.mod引入yourmodule "module/x/module/foo")
func DefaultConfig(chainid, homedir string, modules ...module.AppModuleBasic) Config {
	// 默认带auth和bank模块
	var m []module.AppModuleBasic
	m = append(modules,
		bank.AppModule{},
		auth.AppModule{},
	)
	return Config{
		ChainID:          chainid,
		NodeRPC:          "http://127.0.0.1:26657",
		NodeGRPC:         "127.0.0.1:9090",
		HomeDir:          homedir,
		KeyringBackend:   keyring.BackendTest,
		BroadcastMode:    txtypes.BroadcastMode_BROADCAST_MODE_ASYNC,
		Fee:              "0stake",
		GasLimit:         200000,
		GRPCMaxRecvBytes: 32 << 20, // 32 MiB
		GRPCMaxSendBytes: 32 << 20, // 32 MiB
		KeepaliveTime:    30 * time.Second,
		KeepaliveTimeout: 10 * time.Second,
		SkipConfirmation: false,
		RegisterModules: func() moduletestutil.TestEncodingConfig {
			return moduletestutil.MakeTestEncodingConfig(
				m...,
			)
		},
	}
}
