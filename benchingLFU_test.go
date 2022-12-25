package lfu

import (
	"strconv"
	"testing"
	"time"
)

type TestData struct {
	Key string
	Val string
}

func getTestData(size int) []TestData {
	testcase := make([]TestData, 0, size)
	for i := 0; i < size; i++ {
		testcase = append(testcase, TestData{
			Key: "key_" + strconv.Itoa(i),
			Val: "VAL-" + strconv.Itoa(i) + strconv.Itoa(time.Now().Nanosecond()),
		})
	}
	return testcase
}

func benchmarkSetSequence(testcase []TestData, b *testing.B) {
	b.Run("Benchmark for "+strconv.Itoa(len(testcase))+"elements", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			lfu := NewLFU(50)
			for _, v := range testcase {
				b.StartTimer()
				lfu.Set(v.Key, v.Val)
				b.StopTimer()
			}
		}
	})
}

/*
go test -bench=BenchmarkLFU_Set -benchmem -benchtime=10000x
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
	BenchmarkLFU_Set50/Benchmark_for_50elements-8              10000             52961 ns/op            4278 B/op        151 allocs/op
	BenchmarkLFU_Set100/Benchmark_for_100elements-8            10000            109699 ns/op            8779 B/op        303 allocs/op
	BenchmarkLFU_Set150/Benchmark_for_150elements-8            10000            161257 ns/op           13099 B/op        455 allocs/op
	BenchmarkLFU_Set200/Benchmark_for_200elements-8            10000            215904 ns/op           17469 B/op        606 allocs/op
	BenchmarkLFU_Set250/Benchmark_for_250elements-8            10000            283016 ns/op           21999 B/op        757 allocs/op
	BenchmarkLFU_Set300/Benchmark_for_300elements-8            10000            337132 ns/op           26605 B/op        909 allocs/op
*/

func BenchmarkLFU_Set50(b *testing.B)  { benchmarkSetSequence(getTestData(50), b) }
func BenchmarkLFU_Set100(b *testing.B) { benchmarkSetSequence(getTestData(100), b) }
func BenchmarkLFU_Set150(b *testing.B) { benchmarkSetSequence(getTestData(150), b) }
func BenchmarkLFU_Set200(b *testing.B) { benchmarkSetSequence(getTestData(200), b) }
func BenchmarkLFU_Set250(b *testing.B) { benchmarkSetSequence(getTestData(250), b) }
func BenchmarkLFU_Set300(b *testing.B) { benchmarkSetSequence(getTestData(300), b) }

func benchmarkGetSequence(testcase []TestData, b *testing.B) {

	lfu := NewLFU(len(testcase))
	for _, v := range testcase {
		lfu.Set(v.Key, v.Val)
	}
	for _, v := range testcase {
		b.StartTimer()
		_, err := lfu.Get(v.Key)
		b.StopTimer()
		if err != nil {
			b.Logf("GET ERROR:%v", err)
			b.Failed()
		}
	}
}

/*
go test -bench=BenchmarkLFU_Get -benchmem -benchtime=10000x
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
	BenchmarkLFU_Get50-8               10000               							0 B/op          	0 allocs/op
	BenchmarkLFU_Get100-8              10000               							1 B/op          	0 allocs/op
	BenchmarkLFU_Get150-8              10000               							3 B/op          	0 allocs/op
	BenchmarkLFU_Get200-8              10000               							3 B/op          	0 allocs/op
	BenchmarkLFU_Get250-8              10000               							5 B/op          	0 allocs/op
	BenchmarkLFU_Get300-8              10000               51.10 ns/op            	5 B/op          	0 allocs/op
*/

func BenchmarkLFU_Get50(b *testing.B)  { benchmarkGetSequence(getTestData(50), b) }
func BenchmarkLFU_Get100(b *testing.B) { benchmarkGetSequence(getTestData(100), b) }
func BenchmarkLFU_Get150(b *testing.B) { benchmarkGetSequence(getTestData(150), b) }
func BenchmarkLFU_Get200(b *testing.B) { benchmarkGetSequence(getTestData(200), b) }
func BenchmarkLFU_Get250(b *testing.B) { benchmarkGetSequence(getTestData(250), b) }
func BenchmarkLFU_Get300(b *testing.B) { benchmarkGetSequence(getTestData(300), b) }
