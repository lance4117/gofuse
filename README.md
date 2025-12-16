# GoFuse

一个轻量、模块化、可扩展的 Go 项目启动框架，内置常用业务组件的封装，包括 Gin、Redis、数据库、配置管理、日志系统等，适合用于快速启动中小型后端项目。

> 🚀 快速开始： 快速开始： `go get github.com/lance4117/gofuse@latest`

---

## ✨ 功能特性

- 📜 模块化目录结构
- ✅ 单元测试/集成测试结构预留
- 🧱 常用工具初始化与封装
- ⚙️ 多环境配置加载
- 🔗 区块链客户端支持
- 📦 内置缓存、日志、时间等常用工具模块

---

## 📦 核心模块详解

### **🌐 server** - HTTP服务模块
基于 Gin 的 HTTP 服务封装，支持中间件、路由管理。提供了更简洁的API来处理HTTP请求和响应，并支持自定义上下文处理器。

### **⚙️ config** - 配置管理模块
基于 Viper 的配置管理，支持多种格式文件加载（如 YAML、JSON、TOML）。提供了类型安全的配置读取方法，支持默认值和类型转换。

### **🗄️ store** - 存储模块
基于 XORM、pebble、redis 的存储系统封装，提供统一的接口访问不同类型的存储系统，包括关系型数据库、嵌入式键值存储和分布式缓存。

### **📝 logger** - 日志模块
基于 Zap 的高性能日志系统，支持多级别日志记录（Debug、Info、Warn、Error、Panic、Fatal），提供结构化日志输出。

### **🔗 chain** - 区块链客户端模块
Cosmos 区块链客户端支持，提供与Cosmos生态链交互的能力，包括交易构建、签名和广播等功能。

### **📊 monitor** - 系统监控模块
进程监控（CPU/内存/IO/磁盘），支持 CSV 导出。可用于监控应用程序性能指标，并导出为CSV文件进行分析。

### **🧵 pool** - 工作池模块
基于 Ants 的高性能 goroutine 池，有效管理系统资源，避免频繁创建和销毁goroutine带来的开销。

### **📂 fileio** - 文件IO模块
通用文件 IO，内置 CSV 读写实现。提供统一的接口处理不同类型的文件操作。

### **🔄 codec** - 编解码模块
支持 JSON, Base64, MessagePack 等多种编解码格式，方便在不同数据格式之间进行转换。

### **💻 system** - 系统工具模块
系统工具函数（外部命令执行等），简化与操作系统交互的操作。

### **🎫 limiter** - 限流模块
通用限流封装包，按 key 区分，支持基于API路径或其他标识符的请求频率控制。

### **⚡ cache** - 缓存模块
基于 BigCache 的高性能缓存，提供内存级缓存功能，支持过期时间和自动清理机制。

### **🔑 gen** - 数据生成模块
分布式 ID（Sonyflake）和数据生成工具，提供全局唯一ID生成和其他数据生成功能。

### **🚨 errs** - 错误处理模块
统一错误定义与处理，集中管理项目中的各种错误类型。

### **🔒 once** - 单例模式模块
泛型单例模式支持，确保对象只被初始化一次，在并发环境下安全使用。

### **📡 eventbus** - 事件总线模块
轻量级发布/订阅事件总线，支持并发安全的事件分发、单个订阅者取消订阅、异步发布及 panic 恢复机制。

### **⏰ times** - 时间处理模块
时间工具函数（时间戳、格式化、计算等），简化时间相关的操作。

### **🔒 crypt** - 加密模块
提供现代加密算法支持，包括AES-GCM对称加密和Argon2id密钥派生函数。

### **🧮 conv** - 转换模块
提供不同类型之间的转换功能，包括数字和字符串之间的相互转换。

---

## 🛠️ 主要技术栈

- [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- [Viper](https://github.com/spf13/viper) - 配置解决方案
- [XORM](https://gitea.com/xorm/xorm) - ORM库
- [Redis Go](https://github.com/redis/go-redis) - Redis客户端
- [Pebble](https://github.com/cockroachdb/pebble) - 嵌入式键值数据库
- [Zap](https://github.com/uber-go/zap) - 日志库
- [BigCache](https://github.com/allegro/bigcache) - 高性能缓存
- [Sonyflake](https://github.com/sony/sonyflake) - 分布式ID生成器
- [Msgpack](https://github.com/vmihailenco/msgpack) - 高性能序列化库
- [Gofakeit](https://github.com/brianvoe/gofakeit) - 随机数据生成器
- [Ants](https://github.com/panjf2000/ants) - 高性能 goroutine 池
- [Gopsutil](https://github.com/shirou/gopsutil) - 系统监控库
- [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) - 区块链开发框架


---

## 📁 项目结构

```
GoFuse/
├── cache/           # 高性能缓存模块
├── chain/           # 区块链客户端模块
├── client/          # gRPC客户端封装
├── codec/           # 编解码模块
├── config/          # 配置管理模块
├── conv/            # 类型转换模块
├── crypt/           # 加密模块
├── errs/            # 错误处理模块
├── eventbus/        # 事件总线模块
├── fileio/          # 文件IO模块
├── gen/             # ID生成和数据生成模块
├── limiter/         # 限流模块
├── logger/          # 日志模块
├── monitor/         # 系统监控模块
├── once/            # 单例模式模块
├── pool/            # 工作池模块
├── server/          # HTTP服务模块
├── store/           # 存储模块
│   ├── dbs/         # 关系型数据库封装
│   └── kvs/         # 键值存储封装
├── system/          # 系统工具模块
└── times/           # 时间处理模块
```

---

## 🎯 设计理念

GoFuse 遵循以下设计理念：

1. **模块化设计**：每个功能都被封装在独立的模块中，可以根据需求选择性使用。
2. **易于集成**：各模块之间低耦合，可以轻松集成到现有项目中。
3. **高性能**：选用业界优秀的开源组件，确保系统的高性能表现。
4. **类型安全**：充分利用 Go 语言的类型系统，减少运行时错误。
5. **简洁易用**：提供简洁明了的API，降低学习和使用成本。

---

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。