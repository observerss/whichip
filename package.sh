#!/usr/bin/env bash
go build -o ./bin/whichip -ldflags="-s -w" ./src
VER=$(./bin/whichip version)
env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./release/whichip_"${VER}"_windows_x86.exe -trimpath -ldflags="-s -w" ./src
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/whichip_"${VER}"_windows_x64.exe -trimpath -ldflags="-s -w" ./src
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/whichip_"${VER}"_linux_arm64 -trimpath -ldflags="-s -w" ./src
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/whichip_"${VER}"_linux_amd64 -trimpath -ldflags="-s -w" ./src
env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./release/whichip_"${VER}"_darwin_arm64 -trimpath -ldflags="-s -w" ./src
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/whichip_"${VER}"_darwin_amd64 -trimpath -ldflags="-s -w" ./src
