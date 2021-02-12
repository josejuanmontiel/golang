package main

import (
	"strconv"
	"testing"
)

var throuhtput = 1 // 200

func FuncAverage() float64 {
	var v float64
	v = Average([]float64{1, 2})
	return v
}

func FuncViper() int {
	var v int
	v = Viper()
	return v
}

func FuncConfigDb() int {
	var v int
	v = ConfigDb()
	return v
}

func TestAverage(t *testing.T) {
	t.Parallel()
	v := FuncViper()
	// v := FuncConfigDb()
	if v != 2 {
		t.Error("Expected 2, got ", v)
	}
}

func TestTeardownParallel(t *testing.T) {
	// This Run will not return until the parallel tests finish.
	t.Run("group", func(t *testing.T) {
		for i := 0; i < throuhtput; i++ {
			t.Run("Test"+strconv.Itoa(i), TestAverage)
		}
	})
	// <tear-down code>
}

func Benchmark_YourFunc(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			FuncAverage()
		}
	})
}
