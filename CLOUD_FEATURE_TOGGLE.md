# 云端管理功能开关

## 📅 日期
2026-02-08

## ✨ 新增功能

### 云端管理功能开关

现在可以通过配置文件控制是否启用云端管理功能（Hub Server 集中管理）。

## 🎯 配置方式

### 方法 1: 配置文件

在 `config.yaml` 中添加：

```yaml
# 是否启用云端管理功能
cloud_enabled: true   # 启用（默认）
# cloud_enabled: false  # 禁用
```

### 方法 2: 环境变量

```bash
# Windows PowerShell
$env:WX_CHANNEL_CLOUD_ENABLED="false"

# Linux/Mac
export WX_CHANNEL_CLOUD_ENABLED=false
```

## 📊 两种模式对比

### 启用云端管理 (cloud_enabled: true)

**特点**:
- ✅ 连接到 Hub Server
- ✅ 支持远程管理
- ✅ 支持任务分发
- ✅ 支持集中监控
- ✅ 支持多设备管理

**适用场景**:
- 需要集中管理多个客户端
- 需要远程控制和监控
- 需要任务分发和负载均衡

**日志输出**:
```
✓ 云端管理功能已启用
正在连接到云端服务器: ws://wx.dujulaoren.com/ws/client
```

### 禁用云端管理 (cloud_enabled: false)

**特点**:
- ✅ 客户端独立运行
- ✅ 不连接 Hub Server
- ✅ 节省网络资源
- ✅ 更简单的部署
- ✅ 更高的隐私性

**适用场景**:
- 单机使用
- 不需要远程管理
- 网络环境受限
- 注重隐私

**日志输出**:
```
云端管理功能已禁用 (cloud_enabled: false)
```

## 📝 配置示例

### 示例 1: 启用云端管理（默认）

```yaml
# config.yaml
cloud_enabled: true
cloud_hub_url: "ws://wx.dujulaoren.com/ws/client"
machine_id: "DEV-UTc_bGO8vT6cIRun"
bind_token: "fcb7a3"
```

### 示例 2: 禁用云端管理

```yaml
# config.yaml
cloud_enabled: false

# 以下配置可以省略（不会使用）
# cloud_hub_url: ""
# machine_id: ""
# bind_token: ""
```

### 示例 3: 最小配置（独立运行）

```yaml
# config.yaml
port: 2025
download_dir: downloads
cloud_enabled: false
```

## 🔧 实现细节

### 代码变更

1. **配置结构** (`internal/config/config.go`)
   ```go
   type Config struct {
       // ...
       CloudEnabled bool   `mapstructure:"cloud_enabled"`
       CloudHubURL  string `mapstructure:"cloud_hub_url"`
       // ...
   }
   ```

2. **默认值**
   ```go
   viper.SetDefault("cloud_enabled", true) // 默认启用
   ```

3. **启动逻辑** (`internal/app/app.go`)
   ```go
   if app.Cfg.CloudEnabled {
       app.CloudConnector = cloud.NewConnector(app.Cfg, app.WSHub)
       app.CloudConnector.Start()
       utils.Info("✓ 云端管理功能已启用")
   } else {
       utils.Info("云端管理功能已禁用 (cloud_enabled: false)")
   }
   ```

## 📋 配置文件

### config.yaml.example（简化版）
包含最常用的配置项，适合快速开始。

### config.yaml.full（完整版）
包含所有可配置选项及详细说明，适合高级用户。

## ⚠️ 注意事项

### 1. 默认行为
- 如果不配置 `cloud_enabled`，默认为 `true`（启用）
- 保持向后兼容，现有配置不受影响

### 2. 配置优先级
```
数据库配置 > 环境变量 > 配置文件 > 默认值
```

### 3. 禁用后的影响
- 不会连接到 Hub Server
- 不会发送心跳和监控数据
- 不会接收远程任务
- 本地功能不受影响（下载、代理等）

### 4. 动态切换
- 修改配置后需要重启程序
- 不支持运行时动态切换

## 🎉 优势

### 1. 灵活性
- 用户可以根据需求选择模式
- 支持不同的使用场景

### 2. 资源优化
- 独立模式节省网络资源
- 减少不必要的连接

### 3. 隐私保护
- 独立模式不连接外部服务器
- 数据完全本地化

### 4. 简化部署
- 独立模式部署更简单
- 不需要配置 Hub Server

## 📖 使用建议

### 个人用户
建议设置 `cloud_enabled: false`，独立运行即可。

### 企业用户
建议设置 `cloud_enabled: true`，便于集中管理。

### 开发测试
可以根据测试场景灵活切换。

## 🔗 相关文档

- [config.yaml.example](config.yaml.example) - 简化配置示例
- [config.yaml.full](config.yaml.full) - 完整配置示例
- [README.md](README.md) - 项目主文档

---

**创建日期**: 2026-02-08  
**版本**: v1.0.17
