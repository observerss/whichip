#!/usr/bin/env bash

VER="0.1"

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

info "download and install the binary..."
cd /usr/local/bin && curl -O https://github.com/observerss/whichip/blob/main/install/"${FILENAME}"
mv /usr/local/bin/"${FILENAME}" /usr/local/bin/whichip
chmod +x /usr/local/bin/whichip

info "setup systemd"
cd /lib/systemd/system && curl -O https://github.com/observerss/whichip/blob/main/install/whichip.service
chmod 644 /lib/systemd/system/whichip.service
systemctl enable whichip.service
systemctl start whichip.service

info "done"
