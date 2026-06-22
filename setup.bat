@echo off
setlocal EnableDelayedExpansion

echo ========================================
echo    SungClip Setup Script
echo ========================================
echo.

:: Colors (limited in batch)
set "GREEN=[92m"
set "RED=[91m"
set "YELLOW=[93m"
set "RESET=[0m"

:: Check prerequisites
echo [CHECK] Checking prerequisites...

where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Node.js not found. Install from https://nodejs.org/
    exit /b 1
)
for /f "tokens=*" %%a in ('node --version') do echo   Node.js: %%a

where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go not found. Install from https://go.dev/dl/
    exit /b 1
)
for /f "tokens=*" %%a in ('go version') do echo   Go: %%a

where python >nul 2>&1
if %errorlevel% neq 0 (
    where py >nul 2>&1
    if %errorlevel% neq 0 (
        echo [ERROR] Python not found. Install from https://python.org/
        exit /b 1
    )
    for /f "tokens=*" %%a in ('py --version') do echo   Python: %%a
) else (
    for /f "tokens=*" %%a in ('python --version') do echo   Python: %%a
)
echo [OK] All prerequisites found
echo.

:: Setup bin directory
echo [SETUP] Setting up bin directory...
if not exist "bin" mkdir bin

:: Download ffmpeg if not exists
if exist "bin\ffmpeg.exe" (
    echo   ffmpeg.exe already exists
) else (
    echo   Downloading ffmpeg.exe...
    echo   [INFO] Please manually download ffmpeg from https://ffmpeg.org/download.html
    echo   [INFO] Extract and copy ffmpeg.exe to bin\ folder
    echo   [WARN] Auto-download not available in batch. Use setup.ps1 for auto-download.
    pause
)

:: Download yt-dlp if not exists
if exist "bin\yt-dlp.exe" (
    echo   yt-dlp.exe already exists
) else (
    echo   Downloading yt-dlp.exe...
    powershell -Command "Invoke-WebRequest -Uri 'https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe' -OutFile 'bin\yt-dlp.exe' -UseBasicParsing"
    if %errorlevel% neq 0 (
        echo   [ERROR] Failed to download yt-dlp.exe
        exit /b 1
    )
    echo   yt-dlp.exe downloaded
)
echo [OK] Binaries ready
echo.

:: Setup Python
echo [SETUP] Setting up Python environment...
cd scripts\transcript

if exist ".venv" (
    echo   Virtual environment exists
) else (
    echo   Creating virtual environment...
    python -m venv .venv
    if %errorlevel% neq 0 py -m venv .venv
)

echo   Installing Python dependencies...
.venv\Scripts\pip.exe install --upgrade pip >nul
if exist "requirements.txt" (
    .venv\Scripts\pip.exe install -r requirements.txt
)
cd ..\..
echo [OK] Python environment ready
echo.

:: Setup Remotion
echo [SETUP] Setting up Remotion...
cd remotion
if exist "node_modules" (
    echo   node_modules already exists
) else (
    echo   Running npm install...
    call npm install
)
cd ..
echo [OK] Remotion ready
echo.

:: Setup Go
echo [SETUP] Setting up Go project...
go mod tidy
echo   Building SungClip...
go build -o sungclip.exe .
echo [OK] Go project built
echo.

:: Final message
echo ========================================
echo    ✅ SungClip Setup Complete!
echo ========================================
echo.
echo Run: sungclip.exe
echo.
pause