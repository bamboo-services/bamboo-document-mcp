# Bamboo Document MCP Windows Installation Script (PowerShell)
# Project: https://github.com/bamboo-services/bamboo-document-mcp
#
# Usage:
#   iwr -useb https://raw.githubusercontent.com/bamboo-services/bamboo-document-mcp/master/scripts/install.ps1 | iex

param(
    [string]$Version = "",
    [string]$InstallDir = "$env:USERPROFILE\.local\bin"
)

$Repo = "bamboo-services/bamboo-document-mcp"
$BinaryName = "bamboo-document"

# Print functions
function Write-Info { param($msg) Write-Host "[INFO] " -ForegroundColor Blue -NoNewline; Write-Host $msg }
function Write-Success { param($msg) Write-Host "[SUCCESS] " -ForegroundColor Green -NoNewline; Write-Host $msg }
function Write-Err { param($msg) Write-Host "[ERROR] " -ForegroundColor Red -NoNewline; Write-Host $msg; exit 1 }

# Detect architecture
function Get-Arch {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64") { return "amd64" }
    elseif ($arch -eq "ARM64") { return "arm64" }
    else { Write-Err "Unsupported architecture: $arch" }
}

# Get latest version
function Get-LatestVersion {
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $release.tag_name
    } catch {
        Write-Err "Failed to get latest version"
    }
}

# Main process
Write-Host ""
Write-Host "============================================================"
Write-Host "       Bamboo Document MCP Installer (Windows)              "
Write-Host "============================================================"
Write-Host ""

# Get version
if ([string]::IsNullOrEmpty($Version)) {
    Write-Info "Fetching latest version..."
    $Version = Get-LatestVersion
}
Write-Info "Installing version: $Version"

# Detect architecture
$Arch = Get-Arch
Write-Info "Detected architecture: windows/$Arch"

# Set filename
$Binary = "${BinaryName}-windows-${Arch}.exe"
$DownloadUrl = "https://github.com/${Repo}/releases/download/${Version}/${Binary}"

# Create install directory
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# Download
Write-Info "Downloading: $DownloadUrl"
$TempFile = Join-Path $env:TEMP $Binary
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -UseBasicParsing
} catch {
    Write-Err "Download failed. Please check if the version is correct."
}

# Install
$DestPath = Join-Path $InstallDir "${BinaryName}.exe"
Write-Info "Installing to: $DestPath"
Move-Item -Path $TempFile -Destination $DestPath -Force

# Check PATH
$PathDirs = $env:PATH -split ";"
if ($PathDirs -notcontains $InstallDir) {
    Write-Host ""
    Write-Warning "Install directory $InstallDir is not in PATH"
    Write-Host "Please add the following to your PowerShell profile:"
    Write-Host ""
    Write-Host "    `$env:PATH += `";$InstallDir`"" -ForegroundColor Yellow
    Write-Host ""
}

# Success
Write-Host ""
Write-Success "$BinaryName $Version installed successfully!"
Write-Host ""
Write-Host "Usage:"
Write-Host "    $BinaryName --help"
Write-Host ""
