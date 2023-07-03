# ```httpServer``` - simple http server to server sensor data

## Build

### For Raspberry Pi 4B

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build .
```

## Usage

Run `main.go` with parameters

```bash
-a string
    http server address (default ":15848")
-b int
    baud rate (default 230400)
-p string
    serial port name
-s int
    zero shift (default 0)
```

## Receive data

Go to `http://localhost:15848` (by default) to get data from sensor.
