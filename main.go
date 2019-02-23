package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/ipv4"

	log "github.com/sirupsen/logrus"
)

func main() {
	// need to take care of machine dependent stuff such as endianness when we use syscall directly
	// https://storage.googleapis.com/go-attachment/7653/3/afpacket.go
	const proto = (syscall.ETH_P_IP<<8)&0xff00 | syscall.ETH_P_IP>>8

	// gets all interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	// creates a raw socket for each interface
	for _, iface := range ifaces {
		if iface.Name != "lo" {
			continue
		}
		log.Infof("Listen on interface %s", iface.Name)
		fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_DGRAM, proto)
		if err != nil {
			log.Fatalf("Socket(): %s", err)
		}

		if err := syscall.Bind(fd, &syscall.SockaddrLinklayer{
			Protocol: proto,
			Ifindex:  iface.Index,
		}); err != nil {
			log.Fatalf("Bind(): %s", err)
		}

		// goroutine for listening
		go func(fd int) {
			log.Infof("Waiting for incoming packet on fd (%d)", fd)
			for {
				b := make([]byte, 65536)
				n, from, err := syscall.Recvfrom(fd, b, 0)
				if err != nil {
					log.Errorf("Recvfrom(): %s", err)
				}
				ll, ok := from.(*syscall.SockaddrLinklayer)
				if !ok {
					log.Error("Invalid LinkLayer Address Structure")
					continue
				}
				log.Infof("Recieved %d bytes from %s", n, net.HardwareAddr(ll.Addr[0:ll.Halen]))
				b = b[:n]
				hdr, err := ipv4.ParseHeader(b[0:20])
				if err != nil {
					log.Errorf("Invalid IPv4 header: %s", err)
				}
				log.Infof("IP packet is comming from %s to %s", hdr.Src.String(), hdr.Dst.String())

			}
		}(fd)
	}

	// listens to close event
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
