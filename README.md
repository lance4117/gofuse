# GoFuse

一个轻量、模块化、可扩展的 Go 项目启动框架，内置常用业务组件的封装，包括 Gin、Redis、数据库、配置管理、日志系统等，适合用于快速启动中小型后端项目。

> 🛠️ 本项目旨在构建一套贴近真实业务的项目脚手架。

---

## ✨ 功能特性

- 🔧 Gin HTTP 框架封装
- 🧱 常用工具初始化与封装
- 💾 数据库连接管理
- ⚙️ 多环境配置加载（基于 viper）
- 📜 模块化目录结构
- ✅ 单元测试/集成测试结构预留
- 📦 内置缓存、日志、配置、时间、ID生成等常用工具模块

---

## 🗂️ 目录结构

```
GoFuse/
├── fscache/        # 缓存操作封装（基于 bigcache）
├── fsconfig/       # 配置管理模块（基于 viper）
├── fserror/        # 错误码与错误处理
├── fshttp/         # HTTP 服务封装（基于 Gin）
├── fsid/           # ID 生成器（基于 Sonyflake）
├── fslogger/       # 日志系统封装（基于 zap）
├── fsonce/         # 单例模式实现
├── fstime/         # 时间工具类封装
├── test/           # 测试配置和示例
├── go.mod          # 项目依赖管理
└── README.md       # 项目说明文档
```

---

## 📦 核心模块

### fshttp - HTTP 服务封装
基于 Gin 框架封装，提供更简洁的 API 接口定义方式。

### fsconfig - 配置管理
基于 viper 实现，支持多种格式配置文件（JSON, YAML, TOML 等）的加载和读取。

### fslogger - 日志系统
基于 zap 日志库封装，提供高性能的日志记录功能。

### fscache - 缓存管理
基于 bigcache 实现的高性能内存缓存，支持设置过期时间。

### fsonce - 单例模式
提供单例模式工具函数，确保各组件只初始化一次。

### fstime - 时间工具
提供常用的时间处理函数，如获取当前时间戳、时间格式化等。

### fsid - ID 生成器
基于 Sonyflake 算法实现分布式 ID 生成器。

### fserror - 错误处理
统一的错误码定义和错误处理工具。

---

## 🚀 快速开始

### 克隆项目

```bash
git clone git@gitee.com:lance4117/GoFuse.git
cd GoFuse
```

---

## 🛠️ 技术栈

- [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- [Viper](https://github.com/spf13/viper) - 配置解决方案
- [Zap](https://github.com/uber-go/zap) - 日志库
- [BigCache](https://github.com/allegro/bigcache) - 高性能缓存
- [Sonyflake](https://github.com/sony/sonyflake) - 分布式ID生成器

---
