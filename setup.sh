#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   SungClip Setup Script${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check prerequisites
echo -e "${BLUE}▶ Checking prerequisites...${NC}"

check_cmd() {
    if ! command -v "$1" &> /dev/null; then
        echo -e "${RED}✗ $1 is not installed${NC}"
        exit 1
    fi
    echo -e "${YELLOW}  → $1: $(command "$2" 2>/dev/null | head -1)${NC}"
}

check_cmd node "node --version"
check_cmd go "go version"
check_cmd python3 "python3 --version"

echo -e "${GREEN}✓ All prerequisites found${NC}"

# Setup bin directory
echo -e "\n${BLUE}▶ Setting up bin directory...${NC}"
mkdir -p bin

# ffmpeg
if [ -f "bin/ffmpeg" ]; then
    echo -e "${YELLOW}  → ffmpeg already exists${NC}"
else
    echo -e "${YELLOW}  → Please install ffmpeg via your package manager:${NC}"
    echo -e "${YELLOW}     Ubuntu/Debian: sudo apt install ffmpeg${NC}"
    echo -e "${YELLOW}     macOS: brew install ffmpeg${NC}"
    exit 1
fi

# yt-dlp
if [ -f "bin/yt-dlp" ]; then
    echo -e "${YELLOW}  → yt-dlp already exists${NC}"
else
    echo -e "${YELLOW}  → Downloading yt-dlp...${NC}"
    curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o bin/yt-dlp
    chmod +x bin/yt-dlp
fi
echo -e "${GREEN}✓ Binaries ready${NC}"

# Python setup
echo -e "\n${BLUE}▶ Setting up Python environment...${NC}"
cd scripts/transcript
if [ -d ".venv" ]; then
    echo -e "${YELLOW}  → Virtual environment exists${NC}"
else
    python3 -m venv .venv
fi
.venv/bin/pip install --upgrade pip > /dev/null
[ -f "requirements.txt" ] && .venv/bin/pip install -r requirements.txt
cd ../..

# Remotion setup
echo -e "\n${BLUE}▶ Setting up Remotion...${NC}"
cd remotion
[ -d "node_modules" ] || npm install
cd ..

# Go setup
echo -e "\n${BLUE}▶ Setting up Go project...${NC}"
go mod tidy
go build -o sungclip .

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}  ✅ SungClip setup complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Run: ./sungclip"