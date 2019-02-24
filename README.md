# Funny Router [Linux Only]
## Introduction

This project simulate a simple router in Golang. It has no way for finding out destination MAC address so they must be written by hand in routing table.
it uses full matching for its route but it can be improved using packages like [radix](https://github.com/armon/go-radix).

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
