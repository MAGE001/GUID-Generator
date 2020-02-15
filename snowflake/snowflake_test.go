package snowflake

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_snowflake(t *testing.T) {
	s := &fakeStorage{}
	generator := NewSnowflakeGenerator(s)

	ids := generator.NextIds(10)
	assert.Equal(t, len(ids), 10)
	fmt.Println(ids)
}

func Benchmark_10000(b *testing.B) {
	s := &fakeStorage{}
	generator := NewSnowflakeGenerator(s)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(10000)
	}
}

func Benchmark_1000(b *testing.B) {
	s := &fakeStorage{}
	generator := NewSnowflakeGenerator(s)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(1000)
	}
}

func Benchmark_100(b *testing.B) {
	s := &fakeStorage{}
	generator := NewSnowflakeGenerator(s)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(100)
	}
}

func Benchmark_1(b *testing.B) {
	s := &fakeStorage{}
	generator := NewSnowflakeGenerator(s)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.NextIds(1)
	}
}

type fakeStorage struct{}

func (f *fakeStorage) NextNodeId() (int64, error) {
	return 10, nil
}
