#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO_OWNER="FangTianwd"  # 替换为你的GitHub用户名
REPO_NAME="mdns-reflector-go"
FORMULA_FILE="mdns-reflector-go.rb"

# Check if tag is provided
if [ $# -eq 0 ]; then
    echo -e "${RED}Error: Please provide a version tag (e.g., v1.0.0)${NC}"
    echo "Usage: $0 <version-tag>"
    exit 1
fi

VERSION=$1
VERSION_WITHOUT_V=${VERSION#v}

echo -e "${YELLOW}Updating Homebrew formula for version ${VERSION}${NC}"

# Download the source archive
SOURCE_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/archive/refs/tags/${VERSION}.tar.gz"
echo "Downloading source archive from: ${SOURCE_URL}"

if ! curl -s -L -o "/tmp/${REPO_NAME}-${VERSION}.tar.gz" "${SOURCE_URL}"; then
    echo -e "${RED}Error: Failed to download source archive${NC}"
    exit 1
fi

# Calculate SHA256
SHA256=$(shasum -a 256 "/tmp/${REPO_NAME}-${VERSION}.tar.gz" | cut -d' ' -f1)
echo "Calculated SHA256: ${SHA256}"

# Update the formula file
if [ ! -f "${FORMULA_FILE}" ]; then
    echo -e "${RED}Error: Formula file '${FORMULA_FILE}' not found${NC}"
    exit 1
fi

# Update version and SHA256 in formula
sed -i.bak "s|url \".*\"|url \"${SOURCE_URL}\"|" "${FORMULA_FILE}"
sed -i.bak "s|sha256 \".*\"|sha256 \"${SHA256}\"|" "${FORMULA_FILE}"

# Clean up backup file
rm "${FORMULA_FILE}.bak"

# Clean up downloaded file
rm "/tmp/${REPO_NAME}-${VERSION}.tar.gz"

echo -e "${GREEN}Formula updated successfully!${NC}"
echo "Updated ${FORMULA_FILE} with:"
echo "  - Version: ${VERSION}"
echo "  - SHA256: ${SHA256}"
echo ""
echo "Next steps:"
echo "1. Commit and push the updated formula to your Homebrew tap repository"
echo "2. Test the formula locally: brew install --build-from-source ${FORMULA_FILE}"
echo "3. Create a PR to homebrew-core if you want it in the official repository"
