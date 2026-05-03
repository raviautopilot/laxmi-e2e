#!/bin/bash

# setup-go.sh - Force system to use Go from /apps/go/go1.26.1
# Run with: sudo ./setup-go.sh

set -e  # Exit on error

GO_VERSION="1.26.1"
GO_PATH="/apps/go/go${GO_VERSION}"
GO_BIN="${GO_PATH}/bin"

echo "=== Setting up Go ${GO_VERSION} from ${GO_PATH} ==="

# 1. Remove existing Go installations
echo "Removing existing Go installations..."

# Remove system Go packages
if command -v apt &> /dev/null; then
    echo "Removing apt-installed Go..."
    sudo apt-get remove --purge -y golang-go golang-1.2* golang-1.3* golang-1.4* golang-1.5* \
        golang-1.6* golang-1.7* golang-1.8* golang-1.9* golang-1.10* golang-1.11* golang-1.12* \
        golang-1.13* golang-1.14* golang-1.15* golang-1.16* golang-1.17* golang-1.18* \
        golang-1.19* golang-1.20* golang-1.21* golang-1.22* golang-1.23* golang-1.24* \
        golang-1.25* golang-1.26* 2>/dev/null || true
fi

# Remove /usr/local/go if it exists
if [ -d "/usr/local/go" ]; then
    echo "Removing /usr/local/go..."
    sudo rm -rf /usr/local/go
fi

# Remove any Go from /usr/lib
if [ -d "/usr/lib/go" ]; then
    echo "Removing /usr/lib/go..."
    sudo rm -rf /usr/lib/go
fi

# Remove snap versions
if command -v snap &> /dev/null; then
    echo "Removing snap Go versions..."
    sudo snap remove go 2>/dev/null || true
    sudo snap remove goland 2>/dev/null || true
fi

# 2. Create symlinks for system-wide access
echo "Creating symlinks in /usr/local/bin..."

# Remove existing symlinks
sudo rm -f /usr/local/bin/go
sudo rm -f /usr/local/bin/gofmt

# Create new symlinks
sudo ln -sf "${GO_BIN}/go" /usr/local/bin/go
sudo ln -sf "${GO_BIN}/gofmt" /usr/local/bin/gofmt

# 3. Verify the Go installation
echo "Verifying Go installation..."
if [ -x "${GO_BIN}/go" ]; then
    export PATH="/usr/local/bin:${GO_BIN}:$PATH"
    export GOROOT="${GO_PATH}"

    echo "✓ Go binary found at: ${GO_BIN}/go"
    echo "✓ Go version: $(${GO_BIN}/go version)"
    echo "✓ GOROOT: ${GO_PATH}"
else
    echo "✗ Error: Go binary not found at ${GO_BIN}/go"
    exit 1
fi

# 4. Update shell configuration files
echo "Updating shell configuration files..."

# Function to add to shell config
add_to_config() {
    local config_file="$1"
    local backup_file="${config_file}.backup.$(date +%Y%m%d_%H%M%S)"

    if [ -f "$config_file" ]; then
        # Backup original
        cp "$config_file" "$backup_file"

        # Remove any existing Go related lines
        sed -i '/# GO CONFIGURATION - DO NOT EDIT/d' "$config_file"
        sed -i '/export GOROOT=/d' "$config_file"
        sed -i '/export GOPATH=/d' "$config_file"
        sed -i '/export PATH=.*\/go\/bin/d' "$config_file"

        # Add new configuration
        cat >> "$config_file" << 'EOF'

# GO CONFIGURATION - DO NOT EDIT
export GOROOT=/apps/go/go1.26.1
export GOPATH=$HOME/go
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
EOF
        echo "✓ Updated $config_file"
    fi
}

# Update common shell config files
add_to_config "$HOME/.bashrc"
add_to_config "$HOME/.zshrc"
add_to_config "$HOME/.profile"
add_to_config "$HOME/.bash_profile"

# 5. Clean Go caches
echo "Cleaning Go caches..."
rm -rf ~/.cache/go-build
rm -rf ~/go/pkg/mod
go clean -cache -modcache -testcache 2>/dev/null || true

# 6. Create GOPATH directory structure
echo "Creating GOPATH directory structure..."
mkdir -p ~/go/{bin,src,pkg}
chmod 755 ~/go

# 7. Clear GoLand caches if installed
echo "Clearing GoLand caches..."
rm -rf ~/.cache/JetBrains/GoLand*
rm -rf ~/.local/share/JetBrains/GoLand*
rm -rf ~/.config/JetBrains/GoLand*

echo ""
echo "=== Setup Complete! ==="
echo ""
echo "Next steps:"
echo "1. Restart your terminal or run: source ~/.bashrc"
echo "2. Verify installation: go version"
echo "3. For GoLand:"
echo "   - File -> Invalidate Caches and Restart"
echo "   - File -> Settings -> Go -> GOROOT -> Set to: ${GO_PATH}"
echo "   - File -> Settings -> Go -> GOPATH -> Set to: ~/go"
echo ""
echo "Current Go version:"
source ~/.bashrc 2>/dev/null || true
go version