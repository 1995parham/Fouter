package main

import (
	"encoding/binary"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/ipv4"

	log "github.com/sirupsen/logrus"
)

// map each iterface to its socket
var ifaceSocket map[int]int

// based on the following link
// https://go-review.googlesource.com/c/net/+/112817/2/ipv4/header.go
func headerChecksum(b []byte) uint16 {
	// Algorithm taken from: https://en.wikipedia.org/wiki/IPv4_header_checksum.

	// "First calculate the sum of each 16 bit value within the header,
	// skipping only the checksum field itself."
	var chk uint32
	for i := 0; i < 20; i += 2 {
		// Iterating two bytes at a time; checksum bytes occur at offsets
		// 10 and 11.  Skip them.
		if i == 10 {
			continue
		}

		chk += uint32(binary.BigEndian.Uint16(b[i : i+2]))
	}

	// "The first 4 bits are the carry and will be added to the rest of
	// the value."
	carry := uint16(chk >> 16)
	sum := carry + uint16(chk&0x0ffff)

	// "Next, we flip every bit in that value, to obtain the checksum."
	return uint16(^sum)
}

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

		ifaceSocket[iface.Index] = fd

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
					continue
				}
				log.Infof("IP packet is comming from %s to %s", hdr.Src.String(), hdr.Dst.String())

				// reduces TTL
				hdr.TTL = hdr.TTL - 1
				// removes invalid checksum
				hdr.Checksum = 0

				// brand new packet is here!
				hdrb, err := hdr.Marshal()
				if err != nil {
					log.Errorf("Cannot create IPv4 header: %s", err)
					continue
				}
				binary.BigEndian.PutUint16(hdrb[10:12], headerChecksum(hdrb))
				copy(b[0:20], hdrb)
			}
		}(fd)
	}

	// listens to close event
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
