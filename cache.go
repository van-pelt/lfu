package lfu

import (
	"errors"
)

type node struct {
	next     *node
	prev     *node
	key      *string
	value    interface{}
	flagNode state
}

type LFU struct {
	bucket    map[string]*node
	firstNode *node
	size      int
}

var errNotFound = errors.New("record not found")

type state uint8

const (
	elementRoot state = iota
	elementData
)

func NewLFU(maxSize int) LFU {

	return LFU{
		bucket: make(map[string]*node, maxSize),
		firstNode: &node{
			next:     nil,
			prev:     nil,
			key:      nil,
			value:    nil,
			flagNode: elementRoot,
		},
		size: maxSize,
	}
}

func (l *LFU) Len() int {
	return len(l.bucket)
}

func (l *LFU) Set(key string, value interface{}) {
	node := &node{
		next:     nil,
		prev:     nil,
		key:      &key,
		value:    value,
		flagNode: elementData,
	}

	data, ok := l.bucket[key]
	if ok {
		//move to first position & update value
		l.updateElement(data, value)
		return
	}
	if len(l.bucket) == 0 { // set first element - init relation
		node.next = l.firstNode
		node.prev = l.firstNode
		l.firstNode.next = node
		l.firstNode.prev = node
		l.bucket[key] = node
		return
	}
	if len(l.bucket) == l.size {
		old := l.firstNode.prev
		l.firstNode.prev.prev.next = l.firstNode
		l.firstNode.prev = old.prev
		delete(l.bucket, *old.key)
		old = nil
	}
	l.addElement(node)
	l.bucket[key] = node
	return
}

func (l *LFU) addElement(elem *node) {
	elem.prev = l.firstNode
	elem.next = l.firstNode.next
	l.firstNode.next.prev = elem
	l.firstNode.next = elem
}

func (l *LFU) updateElement(elem *node, value interface{}) {
	elem.prev.next = elem.next
	elem.next.prev = elem.prev
	l.addElement(elem)
	elem.value = value
}

func (l *LFU) Get(key string) (interface{}, error) {
	data, ok := l.bucket[key]
	if !ok {
		return nil, errNotFound
	}
	return data.value, nil
}

func (l *LFU) first() *node {
	if len(l.bucket) == 0 {
		return nil
	}
	return l.firstNode.next
}

/*func (l *LFU) Dump() {
	fmt.Println("=====backet======")
	for k, v := range l.bucket {
		fmt.Println("key=", k, "value=", v)
	}
	fmt.Println("=====RING========")

	for elem := l.first(); ; elem = elem.next {
		switch {
		case elem.prev.flagNode == elementRoot && elem.next.flagNode == elementRoot:
			fmt.Printf("%v<=[value=%v]=>%v		key=%s,value=%v\n", "ROOT", elem.value, "ROOT", *elem.key, l.bucket[*elem.key].value)
		case elem.flagNode == elementRoot:
			fmt.Printf("%v<=[value=ROOT]=>%v\n", elem.prev.value, elem.next.value)
			return
		case elem.prev.flagNode == elementRoot && elem.next.flagNode == elementData:
			fmt.Printf("%v<=[value=%v]=>%v		key=%s,value=%v\n", "ROOT", elem.value, elem.next.value, *elem.key, l.bucket[*elem.key].value)
		case elem.prev.flagNode == elementData && elem.next.flagNode == elementRoot:
			fmt.Printf("%v<=[value=%v]=>%v		key=%s,value=%v\n", elem.prev.value, elem.value, "ROOT", *elem.key, l.bucket[*elem.key].value)
		default:
			fmt.Printf("%v<=[value=%v]=>%v		key=%s,value=%v\n", elem.prev.value, elem.value, elem.next.value, *elem.key, l.bucket[*elem.key].value)
		}
	}
}*/
