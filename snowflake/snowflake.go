package snowflake

import (
	"fmt"
	"sync"
	"time"

)

/*
 * +------+----------------------+----------------+-----------+
 * | sign |     delta seconds    |    node id     | sequence  |
 * +------+----------------------+----------------+-----------+
 *   1bit          28bits              20bits         15bits
 */

const (
	epoch = 1577836800 // Wednesday, January 1, 2020 12:00:00 AM

	deltaSecBits  = 28 // delta seconds since service start
	deltaSecShift = seqBits + nodeIdBits
	maxDeltaSec   = -1 ^ (-1 << deltaSecBits)

	// support about 1 million node ids, service get a new node id every restart,
	// this could help solve twitter snowflake clock backwards problem with `delta seconds`
	nodeIdBits  = 20
	nodeIdShift = seqBits
	maxNodeId   = -1 ^ (-1 << nodeIdBits)

	seqBits = 15 // sequence within the same second
	maxSeq  = -1 ^ (-1 << seqBits)
)

var (
	instOnce sync.Once
	inst     *snowflake
)

// snowflake impl Generator interface.
type snowflake struct {
	lastSec int64
	seq     int64
	nodeId  int64
	buffer  chan int64
}

// NewSnowflakeGenerator return a singleton snowflake instance.
func NewSnowflakeGenerator(s NodeIdStorager) *snowflake {
	instOnce.Do(func() {
		now := time.Now().Unix()
		if now <= epoch {
			panic("clock backwards")
		}

		nid, err := s.NextNodeId()
		if nid > maxNodeId || err != nil {
			panic(fmt.Sprintf("NextNodeId failed, nid: %d, err: %v", nid, err))
		}

		inst = &snowflake{
			lastSec: now,
			seq:     0,
			nodeId:  nid,
			buffer:  make(chan int64, 2^16),
		}
		inst.start()
	})
	return inst
}

func (s *snowflake) NextIds(n int) (ids []int64) {
	for i := 0; i < n; i++ {
		ids = append(ids, <-s.buffer)
	}
	return
}

// start filling buffer with global unique ids in the background.
func (s *snowflake) start() {
	go func() {
		for {
			if s.seq >= maxSeq {
				if (s.lastSec - epoch) >= maxDeltaSec {
					panic("delta seconds exhausted")
				}
				s.lastSec += 1
				s.seq = 0
			}
			s.seq += 1
			s.buffer <- ((s.lastSec - epoch) << deltaSecShift) | (s.nodeId << nodeIdShift) | s.seq
		}
	}()
}
