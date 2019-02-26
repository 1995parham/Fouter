# Funny Router [Linux Only]
## Introduction

This project simulate a simple router in Golang.

Correctly implementing this will require the ability to
- read/write to multiple raw sockets
- parse IP headers
- analyze IP addresses and CIDR ranges
- forwarding the packet.

Develop a program that can successfully router packets from any host to any other host in a different subnet. The program must perform the following tasks:
1. Correctly accept and forward all packets based on their destination addresses.
2. Read and decrement the time to live (TTL)
3. Recompute the checksum (after 2)
4. Correctly specific the destination Ethernet address on all forwarded frames.

It has no way for finding out destination MAC address so they must be written by hand in routing table.
It uses full matching for its route but it can be improved using packages like [radix](https://github.com/armon/go-radix).

## Mininet
You can test it on the [Mininet](http://mininet.org/) environment with following instructions:

```sh
go build # build an executable version of funny router
```

```sh
sudo python router.py # run a custom mininet topoloy with a funny router in background
```

```sh
tail -f out.log # follow up funny router logs
```
