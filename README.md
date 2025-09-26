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

## 📦 核心模块

| 模块             | 功能简介                          |
|----------------|-------------------------------|
| **🌐 server**  | 基于 Gin 的 HTTP 服务封装，支持中间件、路由管理 |
| **⚙️ config**  | 基于 Viper 的配置管理，支持多格式文件加载      |
| **🗄️ store**  | 基于 XORM 的数据库操作封装，支持事务和常见 CRUD |
| **📝 logger**  | 基于 Zap 的高性能日志系统，支持多级别日志       |
| **🔗 chain**   | Cosmos 区块链客户端支持               |
| **📊 monitor** | 进程监控（CPU/内存/IO/磁盘），支持 CSV 导出  |
| **🧵 pool**    | 基于 Ants 的高性能 goroutine 池      |
| **📂 fileio**  | 通用文件 IO，内置 CSV 读写实现           |
| **🔄 codec**   | JSON / MessagePack 编解码支持      |
| **💻 system**  | 系统工具函数（外部命令执行等）               |
| **🎫 limiter** | 通用限流封装包,按 key 区分              |
| **⚡ cache**    | 基于 BigCache 的高性能缓存            |
| **🔑 gen**     | 分布式 ID（Sonyflake）和数据生成工具      |
| **🚨 errs**    | 统一错误定义与处理                     |
| **🔒 once**    | 泛型单例模式支持                      |
| **💱 conv**    | 数据类型与字符串之间的相互转换               |
| **⏰ times**    | 时间工具函数（时间戳、格式化、计算等）           |

---

## 🛠️ 主要技术栈

- [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- [Viper](https://github.com/spf13/viper) - 配置解决方案
- [XORM](https://gitea.com/xorm/xorm) - ORM库
- [Zap](https://github.com/uber-go/zap) - 日志库
- [BigCache](https://github.com/allegro/bigcache) - 高性能缓存
- [Sonyflake](https://github.com/sony/sonyflake) - 分布式ID生成器
- [Msgpack](https://github.com/vmihailenco/msgpack) - 高性能序列化库
- [Gofakeit](https://github.com/brianvoe/gofakeit) - 随机数据生成器
- [Ants](https://github.com/panjf2000/ants) - 高性能 goroutine 池
- [Gopsutil](https://github.com/shirou/gopsutil) - 系统监控库

---

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。