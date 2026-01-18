# Git 仓库清理脚本 - v5.3.0 发版前清理

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║  🧹 清理 Git 仓库中的旧文件          ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Cyan

# 检查是否是 Git 仓库
if (!(Test-Path ".git")) {
    Write-Host "❌ 错误：当前目录不是 Git 仓库" -ForegroundColor Red
    exit 1
}

Write-Host "📋 清理计划:" -ForegroundColor Yellow
Write-Host ""

# 1. 删除根目录的旧图片文件（已移动到 assets/）
Write-Host "1️⃣  删除根目录的旧图片文件（已移动到 assets/）" -ForegroundColor Cyan
$oldImages = @(
    "an.png",
    "jietu.png",
    "liang.png",
    "pinglun.png",
    "sous.png",
    "wxq.png",
    "zanshang.png"
)

foreach ($file in $oldImages) {
    $gitCheck = git ls-files $file 2>$null
    if ($gitCheck) {
        Write-Host "  • 删除: $file" -ForegroundColor White
        git rm $file 2>$null
    }
}

# 2. 检查是否有其他需要清理的文件
Write-Host "`n2️⃣  检查其他可能需要清理的文件" -ForegroundColor Cyan

# 检查根目录是否还有 API 文档（应该在 docs/ 或 dev-docs/）
$apiDocs = @(
    "API_README.md",
    "API_QUICK_START.md"
)

foreach ($file in $apiDocs) {
    $gitCheck = git ls-files $file 2>$null
    if ($gitCheck) {
        Write-Host "  • 删除: $file (已移动到 docs/)" -ForegroundColor White
        git rm $file 2>$null
    }
}

# 检查 CHANGELOG.md 是否在根目录（应该在 dev-docs/）
$gitCheck = git ls-files "CHANGELOG.md" 2>$null
if ($gitCheck) {
    Write-Host "  • 删除: CHANGELOG.md (已移动到 dev-docs/)" -ForegroundColor White
    git rm "CHANGELOG.md" 2>$null
}

# 3. 添加新位置的文件
Write-Host "`n3️⃣  添加新位置的文件" -ForegroundColor Cyan
Write-Host "  • 添加: assets/ 目录" -ForegroundColor White
git add assets/ 2>$null

Write-Host "  • 添加: docs/ 目录" -ForegroundColor White
git add docs/ 2>$null

Write-Host "  • 添加: .gitignore" -ForegroundColor White
git add .gitignore 2>$null

Write-Host "  • 添加: 其他更新的文件" -ForegroundColor White
git add README.md DOCUMENTATION.md RELEASE_v5.3.0.md PRE_RELEASE_CHECKLIST.md 2>$null

# 4. 显示当前状态
Write-Host "`n4️⃣  当前 Git 状态" -ForegroundColor Cyan
git status --short

Write-Host "`n✅ 清理完成！" -ForegroundColor Green
Write-Host ""
Write-Host "📝 下一步操作:" -ForegroundColor Yellow
Write-Host "  1. 检查上面的 Git 状态" -ForegroundColor White
Write-Host "  2. 如果正确，运行: git commit -m 'Release v5.3.0 - 文件重组和功能更新'" -ForegroundColor White
Write-Host "  3. 然后运行: git push origin main" -ForegroundColor White
Write-Host ""
