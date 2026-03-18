# Portscanner

Claude helped me with this one. The goroutine stuff I'm still rusty with, but the rest is fairly straight forward.

## Build

```shell
go build -o scanner ./cmd/main.go
```

## Usage

```shell
scanner <host> -port <port[,[port]]>
```

### Example

```shell
scanner 192.168.10.1 -port 80,443
```

If you don't specify ports, a common list of ports will be scanned:

```
21, 22, 23, 25, 53, 80, 110,
111, 135, 139,	143, 443, 445,
993, 995, 1723, 3306, 3389, 5900,
8080, 8443, 8888, 27017
```
