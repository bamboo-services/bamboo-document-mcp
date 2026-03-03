# Bamboo Document MCP Windows 安装脚本 (PowerShell)
# 项目地址：https://github.com/bamboo-services/bamboo-document-mcp
#
# 用法：
#   iwr -useb https://raw.githubusercontent.com/bamboo-services/bamboo-document-mcp/master/scripts/install.ps1 | iex

param(
    [string]$Version = "",
    [string]$InstallDir = "$env:USERPROFILE\.local\bin"
)

$Repo = "bamboo-services/bamboo-document-mcp"
$BinaryName = "bamboo-document"

# 打印函数
function Write-Info { param($msg) Write-Host "[INFO] " -ForegroundColor Blue -NoNewline; Write-Host $msg }
function Write-Success { param($msg) Write-Host "[SUCCESS] " -ForegroundColor Green -NoNewline; Write-Host $msg }
function Write-Err { param($msg) Write-Host "[ERROR] " -ForegroundColor Red -NoNewline; Write-Host $msg; exit 1 }

# 检测架构
function Get-Arch {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64") { return "amd64" }
    elseif ($arch -eq "ARM64") { return "arm64" }
    else { Write-Err "不支持的架构: $arch" }
}

# 获取最新版本
function Get-LatestVersion {
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $release.tag_name
    } catch {
        Write-Err "无法获取最新版本号"
    }
}

# 主流程
Write-Host ""
Write-Host "╔══════════════════════════════════════════════════════╗"
Write-Host "║       Bamboo Document MCP 安装脚本 (Windows)         ║"
Write-Host "╚══════════════════════════════════════════════════════╝"
Write-Host ""

# 获取版本
if ([string]::IsNullOrEmpty($Version)) {
    Write-Info "正在获取最新版本..."
    $Version = Get-LatestVersion
}
Write-Info "安装版本: $Version"

# 检测架构
$Arch = Get-Arch
Write-Info "检测到架构: windows/$Arch"

# 设置文件名
$Binary = "${BinaryName}-windows-${Arch}.exe"
$DownloadUrl = "https://github.com/${Repo}/releases/download/${Version}/${Binary}"

# 创建安装目录
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# 下载
Write-Info "正在下载: $DownloadUrl"
$TempFile = Join-Path $env:TEMP $Binary
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -UseBasicParsing
} catch {
    Write-Err "下载失败，请检查版本号是否正确"
}

# 安装
$DestPath = Join-Path $InstallDir "${BinaryName}.exe"
Write-Info "正在安装到: $DestPath"
Move-Item -Path $TempFile -Destination $DestPath -Force

# 检查 PATH
$PathDirs = $env:PATH -split ";"
if ($PathDirs -notcontains $InstallDir) {
    Write-Host ""
    Write-Warning "安装目录 $InstallDir 不在 PATH 中"
    Write-Host "请将以下内容添加到你的 PowerShell 配置文件:"
    Write-Host ""
    Write-Host "    `$env:PATH += `";$InstallDir`"" -ForegroundColor Yellow
    Write-Host ""
}

# 成功
Write-Host ""
Write-Success "✅ $BinaryName $Version 安装成功！"
Write-Host ""
Write-Host "使用方法："
Write-Host "    $BinaryName --help"
Write-Host ""
