# Bamboo Document MCP Windows 卸载脚本 (PowerShell)
# 项目地址：https://github.com/bamboo-services/bamboo-document-mcp

param(
    [string]$InstallDir = "$env:USERPROFILE\.local\bin"
)

$BinaryName = "bamboo-document"

function Write-Info { param($msg) Write-Host "[INFO] " -ForegroundColor Blue -NoNewline; Write-Host $msg }
function Write-Success { param($msg) Write-Host "[SUCCESS] " -ForegroundColor Green -NoNewline; Write-Host $msg }
function Write-Warning { param($msg) Write-Host "[WARN] " -ForegroundColor Yellow -NoNewline; Write-Host $msg }

Write-Host ""
Write-Host "╔══════════════════════════════════════════════════════╗"
Write-Host "║       Bamboo Document MCP 卸载脚本 (Windows)         ║"
Write-Host "╚══════════════════════════════════════════════════════╝"
Write-Host ""

$BinaryPath = Join-Path $InstallDir "${BinaryName}.exe"

if (-not (Test-Path $BinaryPath)) {
    Write-Warning "$BinaryName 未安装在 $InstallDir"

    # 尝试查找其他位置
    $Found = Get-Command $BinaryName -ErrorAction SilentlyContinue
    if ($Found) {
        Write-Info "发现安装位置: $($Found.Source)"
        $Confirm = Read-Host "是否删除? [y/N]"
        if ($Confirm -eq "y" -or $Confirm -eq "Y") {
            Remove-Item -Path $Found.Source -Force
            Write-Success "已删除: $($Found.Source)"
        }
    }
    exit 0
}

Write-Info "正在删除: $BinaryPath"
Remove-Item -Path $BinaryPath -Force

Write-Success "✅ $BinaryName 已成功卸载"
Write-Host ""
