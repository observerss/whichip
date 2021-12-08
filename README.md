# whichip: discover (IoT) device's IP in local network

![](https://hjc-image-bed.oss-cn-shanghai.aliyuncs.com/img/20211202145436.png)

## Install

### On (IoT) Device

```bash
wget -O install.sh https://raw.githubusercontent.com/observerss/whichip/main/install/install.sh && sudo bash install.sh
```

Use this script, it will

1) download the `whichip` daemon binary file
2) install as a `systemctl` service
3) run service immediately

(tested on ubuntu 18.04 arm64 only)

### On Any Client

or download binary client in the table below

| OS | Arch | Executable File |
| --- | --- | --- |
| Windows |　x86 | [whichip_v0.1.1_windows_x86.exe](https://github.com/observerss/whichip/releases/download/v0.1/whichip_v0.1.1_windows_x86.exe) |
| Windows |　amd64 | [whichip_v0.1.1_windows_x64.exe](https://github.com/observerss/whichip/releases/download/v0.1/whichip_v0.1.1_windows_x64.exe) |
| Linux |　arm64 | [whichip_v0.1.1_linux_arm64](https://github.com/observerss/whichip/releases/download/v0.1/whichip_v0.1.1_linux_arm64) |
| Linux |　amd64 | [whichip_v0.1.1_linux_amd64](https://github.com/observerss/whichip/releases/download/v0.1/whichip_v0.1.1_linux_amd64) |
| Mac |　arm64 | [whichip_v0.1.1_darwin_arm64](https://github.com/observerss/whichip/releases/download/v0.1/whichip_v0.1_darwin_arm64) |
| Mac |　amd64 | [whichip_v0.1.1_darwin_amd64](https://github.com/observerss/whichip/releases/download/v0.1/whichip_v0.1_darwin_amd64) |

## Usage

Discover Device's IP in your local network

```bash
$ whichip
# 10.86.2.99
```

or

double click `whichip.exe`

All Usages

```raw
NAME:
   whichip - discover device's IP(s) in your local network

USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
   version   get version
   listen    listen to udp broadcast and respond accordingly
   discover  make udp broadcast to discover device's IP(s)
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --timeout value  discover timeout (default: 1)
   --all            print all IPs
   --debug          print debug message
   --port value     the port to bind (default: 53535)
   --help, -h       show help
```

## Other Clients

There are [python](./python) and [nodejs](./nodejs) clients with source code, take a look.