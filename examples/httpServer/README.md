# ```httpServer``` - simple http server to server sensor data

## Usage

Run `main.go` with parameters

```bash
-a string
    http server address (default ":15848")
-b int
    baud rate (default 230400)
-p string
    serial port name
```

## Receive data

Go to `http://localhost:15848` (by default) to get data from sensor.
