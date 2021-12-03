package main

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net"
	"runtime"
	"time"
)

func Listen(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	laddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", c.Int("port")))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp4", laddr)
	defer closeUDP(conn)

	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf)
		log.Debug("received ", string(buf[:n]))
		if err != nil {
			log.Warning(err)
			continue
		}
		if bytes.Equal(buf[:n], CLIENT_MSG) {
			_, err = conn.WriteTo(SERVER_MSG, addr)
			if err != nil {
				log.Warning(err)
				continue
			}
		} else {
			log.Warning("received unknown data ", string(buf[:n]))
		}
	}
}

func Discover(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if c.Command.FullName() == "" {
		c.Set("all", "true")
	}

	var (
		err     error = nil
		ch            = make(chan string)
		waitAll       = c.Bool("all")
		timeout       = c.Float64("timeout")
	)

	ipNets := usableIPNets()
	for _, ipNet := range ipNets {
		go func(ipNet *net.IPNet) {
			laddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:0", ipNet.IP.String()))

			baddr, _ := net.ResolveUDPAddr("udp4",
				fmt.Sprintf("255.255.255.255:%s", fmt.Sprintf("%d", c.Int("port"))))

			conn, err := net.ListenUDP("udp4", laddr)
			defer closeUDP(conn)

			if err != nil {
				log.Debug("error when listening on ", laddr, " ", err)
				return
			}

			_, err = conn.WriteTo(CLIENT_MSG, baddr)
			if err != nil {
				log.Debug("error when broadcasting to ", baddr, " ", err)
				return
			}

			buf := make([]byte, 1024)
			conn.SetDeadline(time.Now().Add(time.Duration(timeout * float64(time.Second))))
			for {
				n, addr, err := conn.ReadFromUDP(buf)
				if err != nil {
					if !err.(net.Error).Timeout() {
						log.Debug("unexpected error ", err)
					}
					break
				}
				if bytes.Equal(buf[:n], SERVER_MSG) {
					ch <- addr.IP.String()
					if !waitAll {
						break
					}
				} else {
					log.Debug("received irrelevant data ", buf[:n])
				}
			}
		}(ipNet)
	}

	result := make([]string, 0, len(ipNets))
	timer := time.NewTimer(time.Duration(timeout * float64(time.Second)))
outer:
	for i := 0; i < len(ipNets); i++ {
		select {
		case <-timer.C:
			break outer

		case ip := <-ch:
			if len(ip) > 0 {
				fmt.Printf("%s\n", ip)
				result = append(result, ip)
				if !waitAll {
					break outer
				}
			}
		}

	}

	if len(result) == 0 {
		err = errors.New("not found")
	}

	// we prompt a press any key message
	// only on windows, and being used without discover command
	if runtime.GOOS == "windows" && c.Command.FullName() == "" {
		PressAnyKey()
	}
	return err
}

func closeUDP(conn *net.UDPConn) {
	if conn != nil {
		log.Debug("close udp conn")
		err := (*conn).Close()
		if err != nil {
			log.Debug("close udp error ", err)
		}
	}
}
