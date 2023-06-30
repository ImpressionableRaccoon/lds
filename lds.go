package lds

import (
	"context"
	"sync"
	"time"

	"go.bug.st/serial"
)

type Data struct {
	Config      Config    `json:"config"`
	LastUpdated time.Time `json:"last_updated"`
	Points      []Point   `json:"point"`
}

type Point struct {
	Intensity uint16 `json:"intensity"`
	Range     uint16 `json:"range"`
}

type Config struct {
	PortName string `json:"port_name"`
	BaudRate int    `json:"baud_rate"`
}

type Lidar struct {
	port serial.Port

	points      []Point
	lastUpdated time.Time
	mu          sync.RWMutex

	cfg Config
}

func New(cfg Config) (_ *Lidar, err error) {
	l := Lidar{
		points: make([]Point, 360),

		cfg: cfg,
	}

	mode := &serial.Mode{
		BaudRate: l.cfg.BaudRate,
	}

	l.port, err = serial.Open(l.cfg.PortName, mode)
	if err != nil {
		return nil, err
	}
	err = l.port.SetReadTimeout(serial.NoTimeout)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *Lidar) Get() Data {
	l.mu.RLock()
	defer l.mu.RUnlock()

	points := make([]Point, len(l.points))
	copy(points, l.points)

	return Data{
		Config:      l.cfg,
		LastUpdated: l.lastUpdated,
		Points:      points,
	}
}

func (l *Lidar) Worker(ctx context.Context) error {
	raw := make([]byte, 2520)
	var err error

	startCount := 0
	for {
		if err = ctx.Err(); err != nil {
			return err
		}

		_, err = l.port.Read(raw[startCount : startCount+1])
		if err != nil {
			return err
		}

		if startCount == 0 {
			if raw[startCount] == 0xFA {
				startCount = 1
			}
			continue
		}

		if raw[startCount] != 0xA0 {
			startCount = 0
			continue
		}
		startCount = 0

		currentByteIndex := 2
		var n int
		for currentByteIndex < len(raw) {
			n, err = l.port.Read(raw[currentByteIndex:])
			if err != nil {
				return err
			}
			currentByteIndex += n
		}

		l.updatePoints(raw)
	}
}

func (l *Lidar) Close() error {
	return l.port.Close()
}

func (l *Lidar) updatePoints(raw []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i := 0; i < len(raw); i += 42 {
		if raw[i] != 0xFA || raw[i+1] != byte(0xA0+i/42) {
			continue
		}

		for j := i + 4; j < i+40; j += 6 {
			index := 6*(i/42) + (j-4-i)/6

			byte0 := uint16(raw[j])
			byte1 := uint16(raw[j+1])
			byte2 := uint16(raw[j+2])
			byte3 := uint16(raw[j+3])

			l.points[index].Intensity = (byte1 << 8) + byte0
			l.points[index].Range = (byte3 << 8) + byte2
		}
	}

	l.lastUpdated = time.Now().UTC()
}
