# 📦 v5.3.0 发版指南

## 🎯 发版前准备

### ✅ 已完成的工作

- [x] 代码编译成功（wx_channel.exe）
- [x] 版本号更新（5.3.0）
- [x] 文档更新（README、CHANGELOG、RELEASE）
- [x] 文档整理（docs、dev-docs、assets）
- [x] .gitignore 配置

### ⚠️ 重要：文件重组

由于本次更新涉及大量文件移动，需要先清理 Git 仓库中的旧文件：

**移动的文件：**
- 图片文件：根目录 → `assets/`
- API 文档：根目录 → `docs/`
- 开发文档：`docs/` → `dev-docs/`

## 🧹 步骤 1：清理旧文件

运行清理脚本：

```powershell
.\cleanup_old_files.ps1
```

这个脚本会：
1. 删除根目录的旧图片文件（已移动到 assets/）
2. 删除根目录的旧 API 文档（已移动到 docs/）
3. 添加新位置的文件到 Git
4. 显示当前 Git 状态

**或者手动清理：**

```bash
# 删除旧位置的图片
git rm an.png jietu.png liang.png pinglun.png sous.png wxq.png zanshang.png

# 删除旧位置的文档（如果存在）
git rm API_README.md API_QUICK_START.md CHANGELOG.md

# 添加新位置的文件
git add assets/ docs/ .gitignore README.md DOCUMENTATION.md
```

## 📝 步骤 2：提交更改

```bash
# 查看状态
git status

# 提交更改
git commit -m "Release v5.3.0 - 通用批量下载组件

主要更新：
- 通用批量下载组件（减少400+行代码）
- Home页面分类视频批量下载
- 视频列表优化（完整信息、分页）
- 下载功能增强（强制重下、取消、进度）
- 搜索页面增强（直播数据、HTML清理）
- Bug修复（下载显示、复选框等）

文件重组：
- 图片移至 assets/ 目录
- API文档移至 docs/ 目录
- 开发文档移至 dev-docs/ 目录（不同步）
- 更新所有文档链接"

# 推送到远程
git push origin main
```

## 🏷️ 步骤 3：创建标签

```bash
# 创建标签
git tag -a v5.3.0 -m "Release v5.3.0 - 通用批量下载组件

主要更新：
- 通用批量下载组件
- Home页面分类视频批量下载
- 视频列表优化
- 下载功能增强
- 搜索页面增强
- Bug修复和文档优化"

# 推送标签
git push origin v5.3.0
```

## 🚀 步骤 4：创建 GitHub Release

1. **访问 GitHub Release 页面**
   - 打开：https://github.com/nobiyou/wx_channel/releases
   - 点击 "Draft a new release"

2. **填写 Release 信息**
   - **Tag**：选择 `v5.3.0`
   - **Title**：`v5.3.0 - 通用批量下载组件`
   - **Description**：复制 `RELEASE_v5.3.0.md` 的内容

3. **上传文件**
   - 上传 `wx_channel.exe`（12.01 MB）

4. **发布**
   - 点击 "Publish release"

## ✅ 步骤 5：验证

- [ ] 检查 Release 页面显示正常
- [ ] 下载链接可用
- [ ] 文档链接正确
- [ ] 图片显示正常

## 📊 预期的 Git 变更

### 删除的文件（旧位置）
```
D  an.png
D  jietu.png
D  liang.png
D  pinglun.png
D  sous.png
D  wxq.png
D  zanshang.png
```

### 新增的文件
```
A  assets/an.png
A  assets/jietu.png
A  assets/liang.png
A  assets/pinglun.png
A  assets/sous.png
A  assets/wxq.png
A  assets/zanshang.png
A  assets/README.md
A  docs/API_README.md
A  docs/API_QUICK_START.md
A  docs/BUILD.md
A  docs/CONFIGURATION.md
A  RELEASE_v5.3.0.md
A  PRE_RELEASE_CHECKLIST.md
A  RELEASE_GUIDE.md
```

### 修改的文件
```
M  .gitignore
M  README.md
M  DOCUMENTATION.md
M  docs/README.md
M  internal/config/config.go
M  main.go
```

## ⚠️ 注意事项

1. **dev-docs/ 不会同步**
   - 包含内部技术细节
   - 已在 .gitignore 中排除

2. **web/ 不会同步**
   - Web 控制台源代码
   - 已在 .gitignore 中排除

3. **编译产物不同步**
   - wx_channel.exe 只在 Release 中提供
   - 已在 .gitignore 中排除

4. **文件移动历史**
   - Git 会自动检测文件移动
   - 使用 `git log --follow` 可以查看文件历史

## 🔧 故障排查

### 问题：Git 显示大量删除和新增

**原因**：Git 可能没有识别出文件移动

**解决**：
```bash
# 使用 git mv 命令（如果还没移动）
git mv old_path new_path

# 或者让 Git 自动检测
git add -A
git status  # 应该显示 renamed
```

### 问题：图片在 GitHub 上不显示

**原因**：README 中的图片路径不正确

**解决**：
- 检查 README.md 中的图片路径
- 应该是 `assets/图片名.png`
- 不是 `./assets/图片名.png` 或 `/assets/图片名.png`

### 问题：文档链接失效

**原因**：文档移动后链接未更新

**解决**：
- 检查 DOCUMENTATION.md
- 检查 README.md
- 检查 docs/README.md
- 确保所有链接指向正确位置

## 📮 发版后

- [ ] 在社区/论坛发布更新公告
- [ ] 收集用户反馈
- [ ] 记录已知问题
- [ ] 规划下一版本

---

**发版日期**：2026-01-18  
**版本号**：v5.3.0  
**负责人**：_________
