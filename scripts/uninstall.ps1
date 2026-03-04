# Bamboo Document MCP Windows Uninstall Script (PowerShell)
# Project: https://github.com/bamboo-services/bamboo-document-mcp

param(
    [string]$InstallDir = "$env:USERPROFILE\.local\bin"
)

$BinaryName = "bamboo-document"

function Write-Info { param($msg) Write-Host "[INFO] " -ForegroundColor Blue -NoNewline; Write-Host $msg }
function Write-Success { param($msg) Write-Host "[SUCCESS] " -ForegroundColor Green -NoNewline; Write-Host $msg }
function Write-Warning { param($msg) Write-Host "[WARN] " -ForegroundColor Yellow -NoNewline; Write-Host $msg }

Write-Host ""
Write-Host "============================================================"
Write-Host "       Bamboo Document MCP Uninstaller (Windows)            "
Write-Host "============================================================"
Write-Host ""

$BinaryPath = Join-Path $InstallDir "${BinaryName}.exe"

if (-not (Test-Path $BinaryPath)) {
    Write-Warning "$BinaryName is not installed in $InstallDir"

    # Try to find other locations
    $Found = Get-Command $BinaryName -ErrorAction SilentlyContinue
    if ($Found) {
        Write-Info "Found installation at: $($Found.Source)"
        $Confirm = Read-Host "Remove it? [y/N]"
        if ($Confirm -eq "y" -or $Confirm -eq "Y") {
            Remove-Item -Path $Found.Source -Force
            Write-Success "Deleted: $($Found.Source)"
        }
    }
    exit 0
}

Write-Info "Removing: $BinaryPath"
Remove-Item -Path $BinaryPath -Force

Write-Success "$BinaryName has been successfully uninstalled"
Write-Host ""
