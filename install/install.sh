#!/usr/bin/env bash

VER="v0.1"

OS=$(uname)
case "$OS" in
Linux)
  OS="linux"
  ;;
Darwin)
  OS="darwin"
  ;;
esac

ARCH=$(arch)
case "$ARCH" in
aarch64)
  ARCH="arm64"
  ;;
x86_64)
  ARCH="amd64"
  ;;
esac

FILENAME=whichip_"${VER}"_"${OS}"_"${ARCH}"

echo "download and install the binary..."
wget -O /usr/local/bin/whichip https://github.com/observerss/whichip/releases/download/${VER}/"${FILENAME}"
#mv /usr/local/bin/"${FILENAME}" /usr/local/bin/whichip
chmod +x /usr/local/bin/whichip

echo "setup systemd"
wget -O /lib/systemd/system/whichip.service https://raw.githubusercontent.com/observerss/whichip/main/install/whichip.service
chmod 644 /lib/systemd/system/whichip.service
systemctl enable whichip.service
systemctl start whichip.service

echo "done"
