# GoFuse

一个轻量、模块化、可扩展的 Go 项目启动框架，内置常用业务组件的封装，包括 Gin、Redis、数据库、配置管理、日志系统等，适合用于快速启动中小型后端项目。

> 🚀 快速开始： go get gitee.com/lance4117/GoFuse@latest

---

## ✨ 功能特性

- 🔧 HTTP 框架封装
- 🧱 常用工具初始化与封装
- 💾 数据库连接管理
- ⚙️ 多环境配置加载
- 📜 模块化目录结构
- ✅ 单元测试/集成测试结构预留
- 🔗 区块链客户端支持
- 🎨 数据生成工具（文章,ID生成等）
- 📊 进程监控工具
- 🏊 工作池（协程池）支持
- 📝 通用写入器接口及CSV写入实现
- 📁 文件IO操作支持
- 📦 内置缓存、日志、时间等常用工具模块

---

## 📦 核心模块

### server - HTTP 服务封装

基于 Gin 框架封装，提供更简洁的 API 接口定义方式。支持中间件注册、GET/POST 请求处理等常用功能。

### config - 配置管理

基于 viper 实现，支持多种格式配置文件（JSON, YAML, TOML 等）的加载和读取。提供单例模式的配置访问方式。

### logger - 日志系统

基于 zap 日志库封装，提供高性能的日志记录功能。支持多种日志级别：Debug、Info、Warn、Error、Panic、Fatal。

### cache - 缓存管理

基于 bigcache 实现的高性能内存缓存，支持设置过期时间。提供键值对存储、获取、删除等操作。

### gen - 生成器

包含分布式 ID 生成器（Sonyflake算法）和基于gofakeit的文章生成器，可生成全局唯一的 ID 和英文文章。

### errs - 错误处理

统一的错误定义和处理模块，包含常见的错误类型定义。

### once - 单例模式

提供泛型支持的单例模式工具函数，支持带参数和错误返回的单例模式实现。

### times - 时间工具

提供时间处理相关的工具函数，如获取当前时间戳、时间格式化、时间计算等。

### system - 系统工具

提供系统相关工具函数，如执行外部命令、IO操作等。

### chain - 区块链客户端

提供cosmos区块链客户端支持。

### monitor - 进程监控

提供进程监控功能，可以监控指定进程的CPU、内存、IO、磁盘等指标，并将结果导出为CSV文件。支持插件化采集器设计。

### pool - 工作池

基于ants实现的工作池（协程池）封装，支持并发任务执行和结果收集。

### fileio - 文件IO操作

提供通用文件读写接口和CSV文件操作实现，用于处理各种格式的文件读写操作。

### writer - 数据写入器

提供通用写入器接口和CSV写入器实现，用于将数据写入不同格式的文件。

---

## 🛠️ 主要技术栈

- [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- [Viper](https://github.com/spf13/viper) - 配置解决方案
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