package lfu

import (
	"github.com/google/uuid"
	"strconv"
	"testing"
)

type TestData struct {
	Key string
	Val string
}

func TestLFU_Add(t *testing.T) {
	size := 4
	testcase := make([]TestData, 0, size)
	testcase = append(testcase, TestData{
		Key: "key_1",
		Val: uuid.New().String(),
	})
	lfu := NewLFU(size)
	//test add first element.The root and first element are related to each other
	lfu.Add(testcase[0].Key, testcase[0].Val)
	// check map - get *Node
	node := lfu.bucket[testcase[0].Key]
	if node.Prev.typeNode != elementRoot || node.Next.typeNode != elementRoot || node.Value.(Cell).value != testcase[0].Val || *node.Value.(Cell).key != testcase[0].Key {
		t.Logf("Fail.Want: node.Prev.typeNode==0,node.Next.typeNode==0,"+
			"node.value=%s,"+
			"node.key=%s."+
			" got:node.Prev.typeNode==%d,"+
			"node.Next.typeNode==%d,"+
			"node.value=%s,node.key=%s", testcase[0].Val, testcase[0].Key, node.Prev.typeNode, node.Next.typeNode, node.Value.(Cell).value, *node.Value.(Cell).key)
		t.Fail()
	}
	//add other element
	for k := 2; k <= size; k++ {
		key := "key_" + strconv.Itoa(k)
		val := uuid.New().String()
		testcase = append(testcase, TestData{
			Key: key,
			Val: val,
		})
		lfu.Add(key, val)
	}

	// check order
	index := len(testcase) - 1
	for elem := lfu.ring.First(); ; elem = elem.Next {
		if elem.typeNode == elementRoot {
			if elem.Next.Value.(Cell).value != testcase[len(testcase)-1].Val {
				t.Logf("Fail.Want Root element next.val=%s,got %s", testcase[len(testcase)-1].Val, elem.Next.Value.(Cell).value)
				t.Fail()
			}
			if elem.Prev.Value.(Cell).value != testcase[0].Val {
				t.Logf("Fail.Want Root element prev.val=%s,got %s", testcase[0].Val, elem.Next.Value.(Cell).value)
				t.Fail()
			}
			break
		}
		if elem.Value.(Cell).value != testcase[index].Val {
			t.Logf("Fail.Want element val=%s,got %s", testcase[0].Val, elem.Next.Value.(Cell).value)
			t.Fail()
		}
		index--
	}

}

func initTestData(size int) ([]TestData, LFU) {
	testcase := make([]TestData, 0, size)
	lfu := NewLFU(size)
	for k := 1; k <= size; k++ {
		key := "key_" + strconv.Itoa(k)
		val := uuid.New().String()
		testcase = append(testcase, TestData{
			Key: key,
			Val: val,
		})
		lfu.Add(key, val)
	}
	return testcase, lfu
}

// tests overwrite functionality
func TestLFU_AddMoreElement(t *testing.T) {
	size := 4
	testcase, lfu := initTestData(size)
	newKey := "NEW_KEY"
	newVal := "NEW_VAL"
	lfu.Add(newKey, newVal)
	//check new element
	if lfu.ring.First().Value.(Cell).value != newVal || *lfu.ring.First().Value.(Cell).key != newKey {
		t.Logf("Fail.Want element val=%s,key=%s,got val=%s,key=%s", newVal, newKey, lfu.ring.First().Value.(Cell).value, *lfu.ring.First().Value.(Cell).key)
		t.Fail()
	}
	//check last element
	if lfu.ring.First().Prev.Prev.Value.(Cell).value != testcase[1].Val || *lfu.ring.First().Prev.Prev.Value.(Cell).key != testcase[1].Key {
		t.Logf("Fail.Want element val=%s,key=%s,got val=%s,key=%s", testcase[1].Val, testcase[1].Key, lfu.ring.First().Prev.Prev.Value.(Cell).value, *lfu.ring.First().Prev.Prev.Value.(Cell).key)
		t.Fail()
	}
}

func TestLFU_AddElementUpdatePosition(t *testing.T) {
	size := 4
	_, lfu := initTestData(size)
	//If we add an element with a key present, then its value is updated and it is shifted to the beginning
	key := "key_3"
	newValue := "UPDATE_VAL"
	lfu.Add(key, newValue)
	if lfu.ring.First().Value.(Cell).value != newValue || *lfu.ring.First().Value.(Cell).key != key {
		t.Logf("Fail update element.Want val=%s,key=%s,got val=%s,key=%s", key, newValue, lfu.ring.First().Value.(Cell).value, *lfu.ring.First().Value.(Cell).key)
		t.Fail()
	}
}

func TestLFU_Get(t *testing.T) {
	size := 4
	key := "key_3"
	testcase, lfu := initTestData(size)
	data, err := lfu.Get(key)
	if err != nil {
		t.Logf("Fail:Get() return err %v", err)
		t.Fail()
	}
	var findKey string
	for _, val := range testcase {
		if val.Val == data {
			findKey = val.Key
			break
		}
	}
	if findKey != key {
		t.Logf("Fail Get() element.Want key=%s,got key=%s", key, findKey)
		t.Fail()
	}
}

func getTestData(size int) []TestData {
	testcase := make([]TestData, 0, size)
	for i := 0; i < size; i++ {
		testcase = append(testcase, TestData{
			Key: "key_" + strconv.Itoa(i),
			Val: uuid.New().String(),
		})
	}
	return testcase
}

// test all unique data sequence.Add random unique
func benchmarkSetSequence(testcase []TestData, b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lfu := NewLFU(50)
		for _, v := range testcase {
			b.StartTimer()
			lfu.Add(v.Key, v.Val)
			b.StopTimer()
		}
	}
}

/*
	cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
  	go test -bench=BenchmarkLFU_Set -benchmem -benchtime=10000x

	BenchmarkLFU_Set50-8               10000             60518 ns/op            8994 B/op        206 allocs/op
	BenchmarkLFU_Set100-8              10000            119474 ns/op           14696 B/op        408 allocs/op
	BenchmarkLFU_Set150-8              10000            188976 ns/op           20221 B/op        610 allocs/op
	BenchmarkLFU_Set200-8              10000            260181 ns/op           25787 B/op        811 allocs/op
	BenchmarkLFU_Set250-8              10000            313833 ns/op           31510 B/op       1012 allocs/op
	BenchmarkLFU_Set300-8              10000            379259 ns/op           37321 B/op       1214 allocs/op

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
		lfu.Add(v.Key, v.Val)
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
	cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
	go test -bench=BenchmarkLFU_Get -benchmem -benchtime=10000x

	BenchmarkLFU_Get50-8               10000               						  	1 B/op          0 allocs/op
	BenchmarkLFU_Get100-8              10000               						  	3 B/op          0 allocs/op
	BenchmarkLFU_Get150-8              10000               							5 B/op          0 allocs/op
	BenchmarkLFU_Get200-8              10000              51.48 ns/op            	6 B/op          0 allocs/op
	BenchmarkLFU_Get250-8              10000               							9 B/op          0 allocs/op
	BenchmarkLFU_Get300-8              10000              0.8700 ns/op              10 B/op         0 allocs/op
*/

func BenchmarkLFU_Get50(b *testing.B)  { benchmarkGetSequence(getTestData(50), b) }
func BenchmarkLFU_Get100(b *testing.B) { benchmarkGetSequence(getTestData(100), b) }
func BenchmarkLFU_Get150(b *testing.B) { benchmarkGetSequence(getTestData(150), b) }
func BenchmarkLFU_Get200(b *testing.B) { benchmarkGetSequence(getTestData(200), b) }
func BenchmarkLFU_Get250(b *testing.B) { benchmarkGetSequence(getTestData(250), b) }
func BenchmarkLFU_Get300(b *testing.B) { benchmarkGetSequence(getTestData(300), b) }
