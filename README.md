# `lds` - Go driver for LDS-01 LIDAR

## Installation

```
go get -u github.com/ImpressionableRaccoon/lds
```

## Usage

```go
l, err := lds.New(lds.Config{
    PortName: "/dev/ttyUSB0",
    BaudRate: 230400,
})
if err != nil {
    panic(err)
}
defer l.Close()

go l.Worker(ctx)

for {
    data = l.Get()
    ...
}
```
