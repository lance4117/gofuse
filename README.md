# GoFuse

一个轻量、模块化、可扩展的 Go 项目启动框架，内置常用业务组件的封装，包括 Gin、Redis、数据库、配置管理、日志系统等，适合用于快速启动中小型后端项目。

> 🚀 快速开始： go get gitee.com/lance4117/GoFuse@latest

---

## ✨ 功能特性

- 🔧 Gin HTTP 框架封装
- 🧱 常用工具初始化与封装
- 💾 数据库连接管理
- ⚙️ 多环境配置加载（基于 viper）
- 📜 模块化目录结构
- ✅ 单元测试/集成测试结构预留
- 🔗 区块链客户端支持
- 📦 内置缓存、日志、配置、时间、ID生成等常用工具模块
- 🎨 数据生成工具（文章生成等）

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

### errors - 错误处理

统一的错误定义和处理模块，包含常见的错误类型定义。

### once - 单例模式

提供泛型支持的单例模式工具函数，支持带参数和错误返回的单例模式实现。

### times - 时间工具

提供时间处理相关的工具函数，如获取当前时间戳、时间格式化、时间计算等。

### utils - 通用工具

提供常用的工具函数，如基于 msgpack 的序列化/反序列化功能。

### chain - 区块链客户端

提供cosmos区块链客户端支持。

---

## 🛠️ 主要技术栈

- [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- [Viper](https://github.com/spf13/viper) - 配置解决方案
- [Zap](https://github.com/uber-go/zap) - 日志库
- [BigCache](https://github.com/allegro/bigcache) - 高性能缓存
- [Sonyflake](https://github.com/sony/sonyflake) - 分布式ID生成器
- [Msgpack](https://github.com/vmihailenco/msgpack) - 高性能序列化库
- [Gofakeit](https://github.com/brianvoe/gofakeit) - 随机数据生成器

---

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。