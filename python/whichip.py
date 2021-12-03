from __future__ import print_function

import argparse
import platform
import socket
import subprocess

import time

PORT = 53535
CLIENT_MSG = b"pfg_ip_broadcast_cl"
SERVER_MSG = b"pfg_ip_response_serv"


def print_nothing(*args, **kwargs):
    pass


def get_bind_ip_windows():
    results = subprocess.check_output(["route", "print", "0.0.0.0"])
    return (
        list(filter(lambda x: b"0.0.0.0" in x, results.splitlines()))[0]
            .split()[-2]
            .decode("ascii")
    )


def get_bind_ip_posix():
    results = subprocess.check_output(["ip", "route"]).decode("utf-8", "ignore")
    vals = results.splitlines()[0].split()
    dev = vals[vals.index("dev") + 1]
    results = subprocess.check_output(["ip", "addr", "show", "dev", dev]).decode(
        "utf-8", "ignore"
    )
    return (
        list(filter(lambda x: "inet" in x, results.splitlines()))[0]
            .split()[1]
            .split("/")[0]
    )


if platform.system().lower() == "windows":
    get_bind_ip = get_bind_ip_windows
else:
    get_bind_ip = get_bind_ip_posix


def listen(log=print, timeout=3.0):
    addr = ("", PORT)
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.settimeout(timeout)
    sock.bind(addr)
    while True:
        try:
            data, address = sock.recvfrom(4096)
        except socket.timeout:
            continue
        except KeyboardInterrupt:
            print("Stopped by Ctrl-C")
            break

        if data:
            log("Received " + str(len(data)) + " bytes from " + str(address))
            log("Data:" + str(data))

            if data == CLIENT_MSG:
                log("responding...")
                sent = sock.sendto(SERVER_MSG, address)
                log("Sent confirmation back")


def discover(log=print, timeout=1.0, all=True):
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind((get_bind_ip(), 0))
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    sock.sendto(CLIENT_MSG, ("<broadcast>", PORT))
    log("waiting to receive")
    ips = []
    begin_time = time.time()
    while True:
        try:
            sock.settimeout(max(0, begin_time + timeout - time.time()))
            data, server = sock.recvfrom(4096)
        except socket.timeout:
            break
        if data == SERVER_MSG:
            log("Received response")
            ip = str(server[0])
            print(ip)
            ips.append(ip)
            if not all:
                break
    if not ips:
        print('no listener is found')
    return ips


def main():
    parser = argparse.ArgumentParser(
        description="discovers server IP in local network",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "type",
        choices=["listen", "discover"],
        help=(
            "The type of operation, \n"
            "`listen` will listen to UDP Broadcasts and respond the IP to the sender, \n"
            "`discover` will make broadcast to listeners and print out their IPs"
        ),
    )
    parser.add_argument(
        "--timeout",
        default=1.0,
        type=float,
        help="socket timeout, default to 1 (seconds)",
    )
    parser.add_argument(
        "--debug",
        action="store_true",
        default=False,
        help="whether to print debug messages",
    )
    parser.add_argument(
        "--all",
        action="store_true",
        default=False,
        help=(
            "discover only, whether to find all listeners, \r\n"
            "if false, only first listener's IP will be printed"
        ),
    )

    args = parser.parse_args()
    log = print_nothing if not args.debug else print
    timeout = args.timeout
    all = args.all

    if args.type == "listen":
        listen(log=log, timeout=timeout)
    elif args.type == "discover":
        discover(log=log, timeout=timeout, all=all)


if __name__ == "__main__":
    main()
