#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

# ANSI Colors
$Green = "$([char]27)[32m"
$Red = "$([char]27)[31m"
$Yellow = "$([char]27)[33m"
$Blue = "$([char]27)[34m"
$Reset = "$([char]27)[0m"

function Write-Step {
    param([string]$Message)

    Write-Host ""
    Write-Host "$Blue[STEP] $Message$Reset"
}

function Write-Success {
    param([string]$Message)

    Write-Host "$Green[OK] $Message$Reset"
}

function Exit-WithError {
    param([string]$Message)

    Write-Host "$Red[FAIL] $Message$Reset"
    exit 1
}

function Write-Info {
    param([string]$Message)

    Write-Host "$Yellow -> $Message$Reset"
}

# ============================================
# 0. CHECK PREREQUISITES
# ============================================

Write-Step "Checking prerequisites..."

# Node.js
try {
    $nodeVersion = node --version 2>$null
    Write-Info "Node.js found: $nodeVersion"
}
catch {
    Exit-WithError "Node.js is not installed. Install from https://nodejs.org/"
}

# Go
try {
    $goVersion = go version 2>$null
    Write-Info "Go found: $goVersion"
}
catch {
    Exit-WithError "Go is not installed. Install from https://go.dev/dl/"
}

# Python
$pythonCmd = $null

try {
    $pyVersion = python --version 2>$null
    $pythonCmd = "python"
    Write-Info "Python found: $pyVersion"
}
catch {
    try {
        $pyVersion = py --version 2>$null
        $pythonCmd = "py"
        Write-Info "Python found via py launcher: $pyVersion"
    }
    catch {
        Exit-WithError "Python is not installed. Install from https://python.org/"
    }
}

Write-Success "All prerequisites found"

# ============================================
# 1. SETUP BIN DIRECTORY
# ============================================

Write-Step "Setting up binaries..."

$binDir = Join-Path $PSScriptRoot "bin"

if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir -Force | Out-Null
}

# --------------------------------------------
# FFMPEG
# --------------------------------------------

$ffmpegPath = Join-Path $binDir "ffmpeg.exe"

if (Test-Path $ffmpegPath) {

    Write-Info "ffmpeg.exe already exists"

}
else {

    Write-Info "Downloading ffmpeg..."

    $ffmpegZip = Join-Path $env:TEMP "ffmpeg-release-essentials.zip"
    $extractDir = Join-Path $env:TEMP "sungclip-ffmpeg"

    try {

        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue

        Invoke-WebRequest `
            -Uri "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" `
            -OutFile $ffmpegZip

        Expand-Archive `
            -Path $ffmpegZip `
            -DestinationPath $extractDir `
            -Force

        $ffmpegExe = Get-ChildItem `
            -Path $extractDir `
            -Recurse `
            -Filter "ffmpeg.exe" |
            Select-Object -First 1

        if (-not $ffmpegExe) {
            throw "ffmpeg.exe not found"
        }

        Copy-Item `
            -Path $ffmpegExe.FullName `
            -Destination $ffmpegPath `
            -Force

        Remove-Item $ffmpegZip -Force -ErrorAction SilentlyContinue
        Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue

        Write-Success "ffmpeg installed"
    }
    catch {
        Exit-WithError "Failed downloading ffmpeg. $_"
    }
}

# --------------------------------------------
# yt-dlp
# --------------------------------------------

$ytDlpPath = Join-Path $binDir "yt-dlp.exe"

if (Test-Path $ytDlpPath) {

    Write-Info "yt-dlp.exe already exists"

}
else {

    Write-Info "Downloading yt-dlp..."

    try {

        Invoke-WebRequest `
            -Uri "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" `
            -OutFile $ytDlpPath

        Write-Success "yt-dlp installed"
    }
    catch {
        Exit-WithError "Failed downloading yt-dlp. $_"
    }
}

Write-Success "Binaries ready"

# ============================================
# 2. PYTHON ENVIRONMENT
# ============================================

Write-Step "Setting up Python environment..."

$pyDir = Join-Path $PSScriptRoot "scripts\transcript"
$venvDir = Join-Path $pyDir ".venv"

if (-not (Test-Path $pyDir)) {
    Write-Info "Python script directory not found, skipping"
}
else {

    if (-not (Test-Path $venvDir)) {

        Write-Info "Creating virtual environment..."

        & $pythonCmd -m venv $venvDir

        Write-Success "Virtual environment created"
    }
    else {

        Write-Info "Virtual environment already exists"
    }

    $pipPath = Join-Path $venvDir "Scripts\pip.exe"

    if (Test-Path $pipPath) {

        Write-Info "Updating pip..."

        & $pipPath install --upgrade pip

        $requirementsPath = Join-Path $pyDir "requirements.txt"

        if (Test-Path $requirementsPath) {

            Write-Info "Installing Python dependencies..."

            & $pipPath install -r $requirementsPath

            Write-Success "Python dependencies installed"
        }
        else {

            Write-Info "requirements.txt not found"
        }
    }
}

# ============================================
# 3. REMOTION
# ============================================

Write-Step "Setting up Remotion..."

$remotionDir = Join-Path $PSScriptRoot "remotion"

if (Test-Path $remotionDir) {

    Push-Location $remotionDir

    try {

        if (Test-Path "node_modules") {

            Write-Info "node_modules already exists"
        }
        else {

            Write-Info "Running npm install..."

            npm install

            Write-Success "Node dependencies installed"
        }
    }
    finally {

        Pop-Location
    }
}
else {

    Write-Info "remotion directory not found, skipping"
}

# ============================================
# 4. GO BUILD
# ============================================

Write-Step "Setting up Go project..."

Write-Info "Running go mod tidy..."

go mod tidy

if ($LASTEXITCODE -ne 0) {
    Exit-WithError "go mod tidy failed"
}

Write-Info "Building SungClip..."

go build -o sungclip.exe .

if ($LASTEXITCODE -ne 0) {
    Exit-WithError "go build failed"
}

Write-Success "SungClip built successfully"

# ============================================
# 5. VERIFY
# ============================================

Write-Step "Running final checks..."

try {
    $ffmpegVersion = & $ffmpegPath -version 2>$null | Select-Object -First 1
    Write-Info "ffmpeg: $ffmpegVersion"
}
catch {
    Write-Info "ffmpeg verification skipped"
}

try {
    $ytVersion = & $ytDlpPath --version
    Write-Info "yt-dlp: $ytVersion"
}
catch {
    Write-Info "yt-dlp verification skipped"
}

$sungclipBinary = Join-Path $PSScriptRoot "sungclip.exe"

if (Test-Path $sungclipBinary) {
    Write-Info "sungclip.exe found"
}
else {
    Exit-WithError "sungclip.exe missing"
}

Write-Success "All checks passed"

# ============================================
# DONE
# ============================================

Write-Host ""
Write-Host "$Green========================================$Reset"
Write-Host "$Green[DONE] SungClip setup complete!$Reset"
Write-Host "$Green========================================$Reset"
Write-Host ""

Write-Host "Next steps:"
Write-Host "  1. Configure your settings"
Write-Host "  2. Run: .\sungclip.exe"
Write-Host ""

Write-Host "Folder structure:"
Write-Host "  bin\                  - ffmpeg.exe, yt-dlp.exe"
Write-Host "  remotion\             - Remotion project"
Write-Host "  scripts\transcript\   - Python scripts"
Write-Host "  internal\             - Go source"
Write-Host "  sungclip.exe          - Main binary"
Write-Host ""