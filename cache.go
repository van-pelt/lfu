// Package lfu implements a simple LRU caching algorithm implemented on two linked lists in the form of a ring structure
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

// LFU implements the basic cache structure
type LFU struct {
	bucket    map[string]*node
	firstNode *node
	size      int
}

// ErrNotFound an error indicating that the item was not found
var ErrNotFound = errors.New("record not found")

type state uint8

const (
	elementRoot state = iota
	elementData
)

// NewLFU creates a new cache object with the specified size
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

// Len returns the number of items in the cache
func (l *LFU) Len() int {
	return len(l.bucket)
}

// Set Adds the key data to the cache and moves the item to the top of the list.If an element with such a key exists, it will be updated and also moved to the top of the list
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
	l.moveToFirst(elem)
	elem.value = value
}

func (l *LFU) moveToFirst(elem *node) {
	elem.prev.next = elem.next
	elem.next.prev = elem.prev
	l.addElement(elem)
}

// Get Returns data from the cache by key.When requesting data, the item is also moved to the top of the list.If the element was not found, an ErrNotFound error is returned
func (l *LFU) Get(key string) (interface{}, error) {
	data, ok := l.bucket[key]
	if !ok {
		return nil, ErrNotFound
	}
	l.moveToFirst(data)
	return data.value, nil
}

func (l *LFU) first() *node {
	if len(l.bucket) == 0 {
		return nil
	}
	return l.firstNode.next
}

// Clear Clears the cache
func (l *LFU) Clear() {
	l.bucket = make(map[string]*node, l.size)
	l.firstNode = &node{
		next:     nil,
		prev:     nil,
		key:      nil,
		value:    nil,
		flagNode: elementRoot,
	}
}

/*func (l *LFU) Dump() {
	if len(l.bucket) == 0 {
		fmt.Printf("oops... the cache is still empty =)\n")
		return
	}
	fmt.Println("====================backet====================")
	for k, v := range l.bucket {
		fmt.Println("key=", k, "value=", v)
	}
	fmt.Println("=====================RING=====================")

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
