package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_random(t *testing.T) {
	generator := NewRandomGenerator()

	ids := generator.NextIds(10)
	assert.Equal(t, len(ids), 10)
	fmt.Println(ids)
}

func Benchmark_10000(b *testing.B) {
	generator := NewRandomGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(10000)
	}
}

func Benchmark_1000(b *testing.B) {
	generator := NewRandomGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(1000)
	}
}

func Benchmark_100(b *testing.B) {
	generator := NewRandomGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(100)
	}
}

func Benchmark_1(b *testing.B) {
	generator := NewRandomGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(1)
	}
}
