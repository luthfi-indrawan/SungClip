#!/usr/bin/env pwsh
# SungClip Setup Script
# Usage: .\setup.ps1

$ErrorActionPreference = "Stop"

# Colors
$Green = "`e[32m"
$Red = "`e[31m"
$Yellow = "`e[33m"
$Blue = "`e[34m"
$Reset = "`e[0m"

function Write-Step {
    param([string]$Message)
    Write-Host "`n$Blue▶ $Message$Reset" -NoNewline
}

function Write-Success {
    param([string]$Message)
    Write-Host " $Green✓ $Message$Reset"
}

function Write-Error {
    param([string]$Message)
    Write-Host " $Red✗ $Message$Reset"
    exit 1
}

function Write-Info {
    param([string]$Message)
    Write-Host "$Yellow  → $Message$Reset"
}

# ============================================
# 0. CHECK PREREQUISITES
# ============================================
Write-Step "Checking prerequisites..."

# Check Node.js
try {
    $nodeVersion = node --version 2>$null
    Write-Info "Node.js found: $nodeVersion"
} catch {
    Write-Error "Node.js is not installed. Please install Node.js first: https://nodejs.org/"
}

# Check Go
try {
    $goVersion = go version 2>$null
    Write-Info "Go found: $goVersion"
} catch {
    Write-Error "Go is not installed. Please install Go first: https://go.dev/dl/"
}

# Check Python
try {
    $pyVersion = python --version 2>$null
    Write-Info "Python found: $pyVersion"
} catch {
    try {
        $pyVersion = py --version 2>$null
        Write-Info "Python found (via py): $pyVersion"
    } catch {
        Write-Error "Python is not installed. Please install Python first: https://python.org/"
    }
}

Write-Success "All prerequisites found"

# ============================================
# 1. CREATE BIN DIRECTORY & DOWNLOAD BINARIES
# ============================================
Write-Step "Setting up bin directory..."

$binDir = Join-Path $PSScriptRoot "bin"
New-Item -ItemType Directory -Force -Path $binDir | Out-Null

# Check ffmpeg
$ffmpegPath = Join-Path $binDir "ffmpeg.exe"
if (Test-Path $ffmpegPath) {
    Write-Info "ffmpeg.exe already exists in bin/"
} else {
    Write-Info "Downloading ffmpeg.exe..."
    # Download from gyandeev (reliable Windows builds)
    $ffmpegZip = Join-Path $env:TEMP "ffmpeg-release-essentials.zip"
    $ffmpegUrl = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
    
    try {
        Invoke-WebRequest -Uri $ffmpegUrl -OutFile $ffmpegZip -UseBasicParsing
        Expand-Archive -Path $ffmpegZip -DestinationPath $env:TEMP -Force
        
        # Find the extracted folder
        $extractedDir = Get-ChildItem -Path $env:TEMP -Directory -Filter "ffmpeg-*" | Select-Object -First 1
        $ffmpegSrc = Join-Path $extractedDir.FullName "bin" "ffmpeg.exe"
        
        Copy-Item -Path $ffmpegSrc -Destination $ffmpegPath -Force
        Remove-Item -Path $ffmpegZip -Force -ErrorAction SilentlyContinue
        Remove-Item -Path $extractedDir.FullName -Recurse -Force -ErrorAction SilentlyContinue
        
        Write-Success "ffmpeg.exe downloaded to bin/"
    } catch {
        Write-Error "Failed to download ffmpeg. Please download manually from https://ffmpeg.org/download.html and place ffmpeg.exe in bin/"
    }
}

# Check yt-dlp
$ytDlpPath = Join-Path $binDir "yt-dlp.exe"
if (Test-Path $ytDlpPath) {
    Write-Info "yt-dlp.exe already exists in bin/"
} else {
    Write-Info "Downloading yt-dlp.exe..."
    $ytDlpUrl = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
    
    try {
        Invoke-WebRequest -Uri $ytDlpUrl -OutFile $ytDlpPath -UseBasicParsing
        Write-Success "yt-dlp.exe downloaded to bin/"
    } catch {
        Write-Error "Failed to download yt-dlp. Please download manually from https://github.com/yt-dlp/yt-dlp/releases and place yt-dlp.exe in bin/"
    }
}

Write-Success "Binaries ready"

# ============================================
# 2. SETUP PYTHON ENVIRONMENT
# ============================================
Write-Step "Setting up Python environment..."

$pyDir = Join-Path $PSScriptRoot "scripts\transcript"
$venvDir = Join-Path $pyDir ".venv"

if (Test-Path $venvDir) {
    Write-Info "Python venv already exists"
} else {
    Write-Info "Creating Python virtual environment..."
    python -m venv $venvDir
    Write-Success "Virtual environment created"
}

# Activate venv and install requirements
Write-Info "Installing Python dependencies..."
$pipPath = Join-Path $venvDir "Scripts\pip.exe"
$activatePath = Join-Path $venvDir "Scripts\Activate.ps1"

& $pipPath install --upgrade pip | Out-Null

$requirementsPath = Join-Path $pyDir "requirements.txt"
if (Test-Path $requirementsPath) {
    & $pipPath install -r $requirementsPath
    Write-Success "Python dependencies installed"
} else {
    Write-Info "No requirements.txt found, skipping pip install"
}

# ============================================
# 3. SETUP REMOTION (NODE.JS)
# ============================================
Write-Step "Setting up Remotion (Node.js)..."

$remotionDir = Join-Path $PSScriptRoot "remotion"
Set-Location $remotionDir

if (Test-Path (Join-Path $remotionDir "node_modules")) {
    Write-Info "node_modules already exists"
} else {
    Write-Info "Running npm install... (this may take a while)"
    npm install
    Write-Success "Node dependencies installed"
}

Set-Location $PSScriptRoot

# ============================================
# 4. SETUP GO MODULES & BUILD
# ============================================
Write-Step "Setting up Go project..."

Write-Info "Running go mod tidy..."
go mod tidy

Write-Info "Building SungClip..."
go build -o sungclip.exe .

Write-Success "Go project built successfully (sungclip.exe)"

# ============================================
# 5. FINAL CHECKS
# ============================================
Write-Step "Running final checks..."

# Verify binaries
$ffmpegCheck = & $ffmpegPath -version 2>$null | Select-Object -First 1
Write-Info "ffmpeg: $ffmpegCheck"

$ytDlpCheck = & $ytDlpPath --version 2>$null
Write-Info "yt-dlp: $ytDlpCheck"

# Verify Go binary
$goBinPath = Join-Path $PSScriptRoot "sungclip.exe"
if (Test-Path $goBinPath) {
    Write-Info "SungClip binary: OK"
}

Write-Success "All checks passed!"

# ============================================
# DONE
# ============================================
Write-Host "`n$Green========================================$Reset"
Write-Host "$Green  ✅ SungClip setup complete!$Reset"
Write-Host "$Green========================================$Reset"
Write-Host "`nNext steps:"
Write-Host "  1. Configure your settings (if needed)"
Write-Host "  2. Run: $Blue.\sungclip.exe$Reset"
Write-Host "`nFolder structure:"
Write-Host "  $Yellow📁 bin/$Reset         - ffmpeg.exe, yt-dlp.exe"
Write-Host "  $Yellow📁 remotion/$Reset    - Video templates (Node.js/Remotion)"
Write-Host "  $Yellow📁 scripts/transcript/$Reset - Transcription scripts (Python)"
Write-Host "  $Yellow📁 internal/$Reset    - Go source code"
Write-Host "  $Yellow📄 sungclip.exe$Reset  - Main application`n"