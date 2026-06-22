@echo off
setlocal

echo === SungClip Transcript Setup ===

set "SCRIPT_DIR=%~dp0"
cd /d "%SCRIPT_DIR%"

set "VENV_DIR=.venv"
set "REQUIREMENTS=requirements.txt"

:: Cari Python
set "PYTHON_CMD="
for %%P in (python3.11 python3.12 python3.10 python3 python py) do (
    %%P --version >nul 2>&1 && (
        set "PYTHON_CMD=%%P"
        goto :found
    )
)

echo ❌ Python not found. Install Python 3.10+ from https://python.org
exit /b 1

:found
%PYTHON_CMD% --version

:: Buat venv
if not exist "%VENV_DIR%" (
    echo Creating .venv...
    %PYTHON_CMD% -m venv %VENV_DIR%
)

:: Install
call "%VENV_DIR%\Scripts\activate.bat"
python -m pip install --upgrade pip
pip install -r %REQUIREMENTS%

echo.
echo ✅ Done.
echo Python: %SCRIPT_DIR%%VENV_DIR%\Scripts\python.exe

pause