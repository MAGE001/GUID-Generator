package random

import (
	"crypto/rand"
	"encoding/binary"
	"log"
)

const (
	readSize = 4096
)

// random impl Generator interface.
type random struct {
	n      int
	buffer chan int64
}

// NewRandomGenerator return a random generator, generate ids by reading /dev/urandom
func NewRandomGenerator() *random {
	r := &random{
		n:      8,
		buffer: make(chan int64, 2^16),
	}
	r.start()
	return r
}

func (r *random) NextIds(n int) (ids []int64) {
	for i := 0; i < n; i++ {
		ids = append(ids, <-r.buffer)
	}
	return
}

// start filling buffer with global unique ids in the background.
func (r *random) start() {
	go func() {
		buffer := make([]byte, readSize)
		for {
			_, err := rand.Read(buffer)
			if err != nil {
				log.Printf("read random failed, err: %s\n", err)
				continue
			}
			for i := 0; i < readSize; i += r.n {
				r.buffer <- int64(binary.LittleEndian.Uint64(buffer[i : i+r.n]))
			}
		}
	}()
}
