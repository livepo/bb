#!/usr/bin/env bash

set -euo pipefail

# Simple installer for bb
# Usage:
#   curl -sL https://github.com/livepo/bb/raw/master/install.sh | bash
#   curl -sL https://github.com/livepo/bb/raw/master/install.sh | bash -s -- v0.1.0

REPO_OWNER="livepo"
REPO_NAME="bb"

VERSION="${1:-}"
INSTALL_DIR="/usr/local/bin"

command_exists() { command -v "$1" >/dev/null 2>&1; }

detect_platform() {
	local os
	local arch
	os="$(uname -s | tr '[:upper:]' '[:lower:]')"
	arch="$(uname -m)"
	case "$os" in
		darwin) os="darwin" ;;
		linux) os="linux" ;;
		msys*|mingw*|cygwin*) os="windows" ;;
		*) os="linux" ;;
	esac
	case "$arch" in
		x86_64|amd64) arch="amd64" ;;
		aarch64|arm64) arch="arm64" ;;
		armv7*) arch="armv7" ;;
		*) arch="amd64" ;;
	esac
	echo "$os" "$arch"
}

get_latest_tag() {
	local api
	api="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
	if command_exists curl; then
		curl -sL "$api" | grep -E '"tag_name"' | head -n1 | sed -E 's/.*: "([^"]+)".*/\1/'
	elif command_exists wget; then
		wget -qO- "$api" | grep -E '"tag_name"' | head -n1 | sed -E 's/.*: "([^"]+)".*/\1/'
	else
		echo ""
	fi
}

download() {
	local url dest
	url="$1"; dest="$2"
	if command_exists curl; then
		curl -sSL -o "$dest" "$url"
	else
		wget -qO "$dest" "$url"
	fi
}

if [ -z "$VERSION" ]; then
	echo "未指定版本，正在获取最新 release tag..."
	VERSION="$(get_latest_tag)"
	if [ -z "$VERSION" ]; then
		echo "无法获取最新版本，请手动指定版本号，例如: bash install.sh v0.1.0" >&2
		exit 1
	fi
	echo "使用版本: $VERSION"
fi

read -r OS ARCH <<<"$(detect_platform)"
echo "检测到平台: $OS / $ARCH"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

EXT=""
BIN_NAME="bb"
if [ "$OS" = "windows" ]; then
	EXT=".exe"
fi

BASE_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${VERSION}"

# Try archive produced by goreleaser (tar.gz)
ARCHIVE_NAME="${REPO_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
ARCHIVE_URL="$BASE_URL/${ARCHIVE_NAME}"
ARCHIVE_DEST="$TMPDIR/${ARCHIVE_NAME}"

echo "尝试下载归档: $ARCHIVE_URL"
if download "$ARCHIVE_URL" "$ARCHIVE_DEST"; then
	if [ -s "$ARCHIVE_DEST" ]; then
		echo "归档下载成功，正在解压..."
		tar -xzf "$ARCHIVE_DEST" -C "$TMPDIR"
		# 查找可执行文件
		if [ -f "$TMPDIR/${BIN_NAME}${EXT}" ]; then
			BIN_SRC="$TMPDIR/${BIN_NAME}${EXT}"
		else
			# sometimes goreleaser puts binaries under a folder
			BIN_SRC="$(find "$TMPDIR" -type f -name "${BIN_NAME}${EXT}" | head -n1)"
		fi
		if [ -z "$BIN_SRC" ]; then
			echo "未能在归档中找到可执行文件 ${BIN_NAME}${EXT}" >&2
			exit 1
		fi
		if [ ! -w "$INSTALL_DIR" ]; then
			sudo mv "$BIN_SRC" "$INSTALL_DIR/${BIN_NAME}${EXT}"
		else
			mv "$BIN_SRC" "$INSTALL_DIR/${BIN_NAME}${EXT}"
		fi
		sudo chmod +x "$INSTALL_DIR/${BIN_NAME}${EXT}" || true
		echo "安装完成"
		"$INSTALL_DIR/${BIN_NAME}${EXT}" version || true
		exit 0
	fi
fi

echo "无法找到可下载的归档。请检查版本号或到 https://github.com/${REPO_OWNER}/${REPO_NAME}/releases 手动下载。" >&2
exit 1
