# 📋 v5.3.0 发版前检查清单

## ✅ 代码和编译

- [x] 代码编译成功（无错误）
- [x] 版本号已更新（5.3.0）
  - [x] `internal/config/config.go`
  - [x] `README.md`
  - [x] `dev-docs/CHANGELOG.md`
- [x] 生成 `wx_channel.exe`（约 12 MB）
- [x] 使用 `-ldflags="-s -w"` 优化体积

## ✅ 文档更新

- [x] `README.md` - 版本说明更新
- [x] `dev-docs/CHANGELOG.md` - 完整更新日志
- [x] `RELEASE_v5.3.0.md` - 发布说明
- [x] `docs/README.md` - 用户文档索引
- [x] `dev-docs/README.md` - 开发文档索引
- [x] `DOCUMENTATION.md` - 文档导航

## ✅ 文档整理

- [x] 用户文档（docs/）- 11 个文件
  - [x] 基础文档（安装、介绍、故障排查）
  - [x] 功能文档（批量下载、Web 控制台）
  - [x] API 文档（API_README、API_QUICK_START）
  - [x] 开发文档（BUILD、CONFIGURATION）
- [x] 开发文档（dev-docs/）- 11 个文件
  - [x] 技术细节（API 实现、评论采集等）
  - [x] 版本管理（CHANGELOG、RELEASE_NOTES）
- [x] 资源文件（assets/）- 7 个图片
  - [x] 界面截图
  - [x] 功能演示图

## ✅ .gitignore 配置

- [x] `dev-docs/` - 已排除（内部技术细节）
- [x] `downloads/` - 已排除（用户下载的视频）
- [x] `logs/` - 已排除（运行日志）
- [x] `web/` - **已排除**（Web 控制台源代码，暂不公开）
- [x] `*.exe` - 已排除（编译产物）
- [x] `*.zip` - 已排除（压缩文件）

## ✅ 链接检查

- [x] README.md 中的所有链接
- [x] docs/README.md 中的链接
- [x] dev-docs/README.md 中的链接
- [x] DOCUMENTATION.md 中的链接
- [x] 图片资源路径（assets/）

## 🔍 需要同步到 GitHub 的内容

### ✅ 源代码
- [x] `cmd/` - 命令行入口
- [x] `internal/` - 内部包
- [x] `inject/` - 前端注入脚本
- [x] `pkg/` - 公共包
- [x] `lib/` - 第三方库
- [x] `scripts/` - 脚本文件
- [x] `winres/` - Windows 资源

### ✅ 文档
- [x] `docs/` - 用户文档
- [x] `assets/` - 图片资源
- [x] `README.md` - 项目主页
- [x] `DOCUMENTATION.md` - 文档导航
- [x] `LICENSE` - 许可证

### ✅ 配置文件
- [x] `go.mod` - Go 模块
- [x] `go.sum` - 依赖锁定
- [x] `.gitignore` - Git 忽略配置

### ❌ 不同步的内容
- [x] `dev-docs/` - 开发文档（内部技术细节）
- [x] `web/` - **Web 控制台源代码（暂不公开）**
- [x] `downloads/` - 下载的视频文件
- [x] `logs/` - 运行日志
- [x] `*.exe` - 编译产物
- [x] `*.zip` - 压缩文件
- [x] `.vscode/` - IDE 配置

## 📦 发版文件

需要上传到 GitHub Release 的文件：
- [x] `wx_channel.exe` - Windows 可执行文件
- [x] `RELEASE_v5.3.0.md` - 发布说明（作为 Release 描述）

## 🚀 发版步骤

### 1. Git 提交
```bash
git add .
git commit -m "Release v5.3.0 - 通用批量下载组件"
git push origin main
```

### 2. 创建 Tag
```bash
git tag -a v5.3.0 -m "Release v5.3.0 - 通用批量下载组件

主要更新：
- 通用批量下载组件
- Home 页面分类视频下载
- 视频列表优化
- 下载功能增强
- 搜索页面增强
- Bug 修复和文档优化"

git push origin v5.3.0
```

### 3. 创建 GitHub Release
1. 访问 GitHub 仓库的 Releases 页面
2. 点击 "Draft a new release"
3. 选择 Tag: `v5.3.0`
4. 标题：`v5.3.0 - 通用批量下载组件`
5. 描述：复制 `RELEASE_v5.3.0.md` 的内容
6. 上传文件：`wx_channel.exe`
7. 点击 "Publish release"

### 4. 验证
- [ ] 检查 Release 页面显示正常
- [ ] 下载链接可用
- [ ] 文档链接正确

## 📝 发版后

- [ ] 在社区/论坛发布更新公告
- [ ] 收集用户反馈
- [ ] 记录已知问题
- [ ] 规划下一版本

## ⚠️ 重要提醒

1. **web/ 目录不同步**：包含 Web 控制台源代码，暂不公开
2. **dev-docs/ 不同步**：包含内部技术细节，不对外公开
3. **编译产物不同步**：.exe 文件只在 Release 中提供
4. **测试后再发布**：确保所有功能正常工作

---

**检查人员：** _________  
**检查日期：** 2026-01-18  
**版本号：** v5.3.0  
**检查状态：** ✅ 已完成
