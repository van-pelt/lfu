package lfu

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

func getTestCase(size int) ([]TestData, LFU) {
	testcase := make([]TestData, 0, size)
	for k := 1; k <= size; k++ {
		testcase = append(testcase, TestData{
			Key: "key_" + strconv.Itoa(k),
			Val: "val_" + strconv.Itoa(k) + "_" + strconv.Itoa(time.Now().Nanosecond()),
		})
	}
	lfu := NewLFU(size)
	for _, v := range testcase {
		lfu.Set(v.Key, v.Val)
	}
	return testcase, lfu
}

func TestLFU_Add_FirstElement(t *testing.T) {
	testcase := TestData{
		Key: "key1",
		Val: "val1",
	}
	lfu := NewLFU(3)
	lfu.Set(testcase.Key, testcase.Val)
	t.Logf("	Check first node")
	if lfu.first().value != testcase.Val {
		t.Logf("Fail.Want:%s,got:%v", testcase.Val, lfu.first().value)
		t.Fail()
	}
	t.Logf("	Check first node relation")
	if lfu.first().flagNode != elementData && lfu.first().prev.flagNode != elementRoot && lfu.first().next.flagNode != elementRoot {
		t.Logf("Fail.Relation is broken")
		t.Fail()
	}
}

func TestLFU_Add_MoreElements(t *testing.T) {
	testcase, lfu := getTestCase(4)
	t.Logf("	Check relation")
	testcaseIndex := len(testcase) - 1
	for e := lfu.first(); ; e = e.next {
		if e.flagNode == elementRoot {
			break
		}
		if testcase[testcaseIndex].Val != e.value || testcase[testcaseIndex].Key != *e.key {
			t.Logf("Fail.Want pair [%v,%v],got [%v,%v]", testcase[testcaseIndex].Key, testcase[testcaseIndex].Val, *e.key, e.value)
			t.Fail()
		}
		testcaseIndex--
	}

}

func TestLFU_Add_CheckRemove(t *testing.T) {

	size := 4
	testcase := make([]TestData, 0, size)
	for k := 1; k <= size+1; k++ { //add 5 element - the first comer must be deleted
		testcase = append(testcase, TestData{
			Key: "key_" + strconv.Itoa(k),
			Val: "val_" + strconv.Itoa(k) + "_" + strconv.Itoa(time.Now().Nanosecond()),
		})
	}
	lfu := NewLFU(size)
	for _, v := range testcase {
		lfu.Set(v.Key, v.Val)
	}
	t.Logf("	Check relation & backet")
	if lfu.first().value != testcase[len(testcase)-1].Val || lfu.bucket[testcase[len(testcase)-1].Key].value != testcase[len(testcase)-1].Val {
		t.Logf("Fail.Want pair values [%v,%v],got [%v,%v]", testcase[len(testcase)-1].Val, testcase[len(testcase)-1].Val, lfu.first().value, lfu.bucket[testcase[len(testcase)-1].Key].value)
		t.Fail()
	}
}

func TestLFU_Add_ExistingElement(t *testing.T) {
	testcase, lfu := getTestCase(4)
	t.Logf("	Check existing element & move to first position")
	existingKey := testcase[2].Key
	newValue := "NewVal"
	lfu.Set(existingKey, newValue)
	if *lfu.first().key != existingKey || lfu.first().value != newValue {
		t.Logf("Fail.Want pair key=>val [%v,%v],got [%v,%v]", existingKey, newValue, *lfu.first().key, lfu.first().value)
		t.Fail()
	}
}

func TestLFU_GetExistingElement(t *testing.T) {
	testcase, lfu := getTestCase(4)
	t.Logf("	Get existing element")
	existingElement := testcase[2]
	data, err := lfu.Get(existingElement.Key)
	if err != nil {
		t.Logf("Fail.Got error %v", err)
		t.Fail()
	}
	if data != existingElement.Val {
		t.Logf("Fail.Want %v,got %v", existingElement.Val, data)
		t.Fail()
	}
}

func TestLFU_GetNotFoundError(t *testing.T) {
	_, lfu := getTestCase(4)
	t.Logf("	Get non-existing element")
	nonExistingKey := "bad_key"
	_, err := lfu.Get(nonExistingKey)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Logf("Fail.Got unknown error %v", err)
			t.Fail()
		}
	}
}

func TestLFU_Len(t *testing.T) {
	testcase, lfu := getTestCase(4)
	t.Logf("	check len")
	if len(testcase) != lfu.Len() {
		t.Logf("Fail.Want len %v,got %v", len(testcase), lfu.Len())
		t.Fail()
	}
}
