#Requires -Version 5.1
$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

$VenvDir = ".\.venv"
$Requirements = ".\requirements.txt"

Write-Host "=== SungClip Transcript Setup ===" -ForegroundColor Cyan

# Cari Python (prioritas 3.11, 3.12, 3.10, lalu python3)
$PythonCmd = $null
$candidates = @("python3.11", "python3.12", "python3.10", "python3", "python")

foreach ($py in $candidates) {
    $found = Get-Command $py -ErrorAction SilentlyContinue
    if ($found) {
        # Cek version minimal 3.10
        $verStr = & $py --version 2>&1
        if ($verStr -match "Python (\d+)\.(\d+)") {
            $major = [int]$Matches[1]
            $minor = [int]$Matches[2]
            if ($major -gt 3 -or ($major -eq 3 -and $minor -ge 10)) {
                $PythonCmd = $py
                break
            }
        }
    }
}

if (-not $PythonCmd) {
    Write-Host "❌ Python 3.10+ not found. Install from https://python.org" -ForegroundColor Red
    exit 1
}

$ver = & $PythonCmd --version 2>&1
Write-Host "Using: $ver" -ForegroundColor Green

# Buat venv kalau belum ada
if (-not (Test-Path $VenvDir)) {
    Write-Host "Creating .venv..." -ForegroundColor Yellow
    & $PythonCmd -m venv $VenvDir
}

# Activate & install
$VenvPython = Join-Path $VenvDir "Scripts\python.exe"
$VenvPip = Join-Path $VenvDir "Scripts\pip.exe"

Write-Host "Upgrading pip..." -ForegroundColor Yellow
& $VenvPython -m pip install --upgrade pip

Write-Host "Installing dependencies..." -ForegroundColor Yellow
& $VenvPip install -r $Requirements

Write-Host ""
Write-Host "✅ Setup complete!" -ForegroundColor Green
Write-Host "Python executable:" -ForegroundColor Cyan
Write-Host "   $VenvPython"
Write-Host ""
Write-Host "Set in your config:" -ForegroundColor Cyan
Write-Host "   PYEXE=$VenvPython"
Write-Host "   PYTranscribe=$ScriptDir\main.py"