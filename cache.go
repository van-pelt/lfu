package lfu

import (
	"errors"
)

type LFU struct {
	maxSize int
	bucket  map[string]*Node
	ring    Ring
}

type Cell struct {
	key   *string
	value interface{}
}

var errNotFound = errors.New("record not found")

func NewLFU(size int) LFU {
	return LFU{
		maxSize: size,
		bucket:  make(map[string]*Node),
		ring:    newRing(),
	}
}

func (l *LFU) Add(key string, value interface{}) {
	cell := Cell{
		key:   &key,
		value: value,
	}
	data, ok := l.bucket[key]
	if ok {
		//update only data & move to first position
		data.Value = cell
		l.ring.MoveToFirst(data)
		return
	}

	if l.ring.Len() == l.maxSize {
		ptr := l.ring.Delete()
		delete(l.bucket, *ptr.Value.(Cell).key)
	}

	l.bucket[key] = l.ring.Add(cell)
	return
}

func (l *LFU) Get(key string) (interface{}, error) {
	data, ok := l.bucket[key]
	if !ok {
		return nil, errNotFound
	}
	return data.Value.(Cell).value, nil
}

/*func (l *LFU) Dump() {
	fmt.Println("=====backet======")
	for k, v := range l.bucket {
		fmt.Println("key=", k, "value=", v)
	}
	fmt.Println("=====RING========")
	for elem := l.ring.First(); ; elem = elem.Next {
		if elem.typeNode == elementRoot {
			fmt.Printf("%v<=[value=ROOT]=>%v\n", elem.Prev.Value.(Cell).value, elem.Next.Value.(Cell).value)
			break
		}
		if elem.Prev.typeNode == elementRoot && elem.typeNode == elementData {
			fmt.Printf("%v<=[value=%v]=>%v		key=%s,value=%v\n", "ROOT", elem.Value.(Cell).value, elem.Next.Value.(Cell).value, *elem.Value.(Cell).key, l.bucket[*elem.Value.(Cell).key].Value.(Cell).value)
			continue
		}
		if elem.Next.typeNode == elementRoot && elem.typeNode == elementData {
			fmt.Printf("%v<=[value=%v]=>%v		key=%s,value=%v\n", elem.Prev.Value.(Cell).value, elem.Value.(Cell).value, "ROOT", *elem.Value.(Cell).key, l.bucket[*elem.Value.(Cell).key].Value.(Cell).value)
			continue
		}
		fmt.Printf("%v<=[value=%v]=>%v		key=%s,value=%v\n", elem.Prev.Value.(Cell).value, elem.Value.(Cell).value, elem.Next.Value.(Cell).value, *elem.Value.(Cell).key, l.bucket[*elem.Value.(Cell).key].Value.(Cell).value)
	}
}*/
