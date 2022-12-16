package lfu

import (
	"testing"
)

func TestRing_Len(t *testing.T) {
	testcase := []int{1, 2, 3}
	ring := newRing()
	for _, v := range testcase {
		ring.Add(v)
	}
	if ring.Len() != len(testcase) {
		t.Logf("ring.Len():want %d,got %d", len(testcase), ring.Len())
		t.Fail()
	}
}

func TestRing_First(t *testing.T) {
	testcase := []int{1, 2, 3}
	ring := newRing()
	for _, v := range testcase {
		ring.Add(v)
	}
	// first element = last element in test slice
	if ring.First().Value != testcase[len(testcase)-1] {
		t.Logf("ring.First():want %d,got %d", testcase[len(testcase)-1], ring.First().Value)
		t.Fail()
	}
}

func TestRing_NextElement(t *testing.T) {
	testcase := []int{1, 2, 3}
	ring := newRing()
	for _, v := range testcase {
		ring.Add(v)
	}
	cnt := len(testcase) - 1
	for e := ring.First(); e != nil; e = e.NextElement() {
		if e.Value != testcase[cnt] {
			t.Logf("ring.Node.NextElement():want %d,got %d", e.Value, testcase[cnt])
			t.Fail()
		}
		cnt--
	}
}

func TestRing_Add(t *testing.T) {
	ring := newRing()
	want := 1
	ring.Add(want)
	got := ring.First()
	if want != got.Value {
		t.Logf("ring.Add()-first element test:want %d,got %d", want, got.Value)
		t.Fail()
	}

	testcase := []int{1, 2, 3, 4, 5}
	container := make([]*Node, 0, len(testcase))
	ring = newRing()
	for _, v := range testcase {
		container = append(container, ring.Add(v))
	}
	cnt := 0
	for i := 0; i < len(container); i++ {
		if container[i].Value != testcase[cnt] {
			t.Logf("ring.Add():want %d,got %d", container[i].Value, testcase[cnt])
			t.Fail()
		}
		cnt++
	}
}

func TestRing_MoveToFirst(t *testing.T) {
	testcase := []int{1, 2, 3, 4, 5}
	container := make([]*Node, 0, len(testcase))
	ring := newRing()
	for _, v := range testcase {
		container = append(container, ring.Add(v))
	}
	// get last element & move to first position
	data := container[len(container)-1]
	ring.MoveToFirst(data)
	if ring.First().Value != data.Value {
		t.Logf("ring.MoveToFirst():want %d,got %d", data.Value, ring.First().Value)
		t.Fail()
	}
}

func TestRing_Delete(t *testing.T) {
	testcase := []int{1, 2, 3}
	ring := newRing()
	for _, v := range testcase {
		ring.Add(v)
	}
	ring.Delete()
	if ring.First().Prev.Prev.Value != testcase[len(testcase)-2] {
		t.Logf("ring.Delete():want last element: %d,got %d", ring.First().Prev.Prev.Value, testcase[len(testcase)-2])
		t.Fail()
	}
}
