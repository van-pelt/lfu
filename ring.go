package lfu

type Node struct {
	Next     *Node
	Prev     *Node
	Value    interface{}
	typeNode state
}

type Ring struct {
	count   int
	nodeseq *Node
}

type state uint8

const (
	elementRoot state = iota
	elementData
)

func newRing() Ring {
	return Ring{
		count: 0,
		nodeseq: &Node{
			Next:     nil,
			Prev:     nil,
			Value:    nil,
			typeNode: elementRoot, //add root-element - border ring
		},
	}
}

func (l *Ring) Add(value interface{}) *Node {

	node := &Node{
		Next:     nil,
		Prev:     nil,
		Value:    value,
		typeNode: elementData,
	}

	if l.count == 0 { // set first element
		node.Prev = l.nodeseq
		node.Next = l.nodeseq
		l.nodeseq.Next = node
		l.nodeseq.Prev = node
	} else { //insert after root-element and before first not-null-element
		node.Next = l.nodeseq.Next
		node.Prev = l.nodeseq
		l.nodeseq.Next.Prev = node
		l.nodeseq.Next = node
	}
	l.count++
	return node
}

func (l *Ring) MoveToFirst(elem *Node) *Node {

	if elem.Prev.typeNode == elementRoot { //check first position
		return elem
	}
	//unlink element
	elem.Next.Prev = elem.Prev
	elem.Prev.Next = elem.Next

	//move to first position
	l.nodeseq.Next.Prev = elem
	ptr := l.nodeseq.Next
	l.nodeseq.Next = elem
	elem.Next = ptr
	elem.Prev = l.nodeseq
	return elem
}

func (l *Ring) Len() int {
	return l.count
}

func (l *Ring) Delete() *Node {

	lastPtr := l.nodeseq.Prev
	l.nodeseq.Prev.Prev.Next = l.nodeseq
	l.nodeseq.Prev = lastPtr.Prev
	l.count--
	return lastPtr

}

func (n *Node) NextElement() *Node {
	if p := n.Next; p != nil && p.typeNode != elementRoot {
		return p
	}
	return nil
}

func (l *Ring) First() *Node {
	if l.count == 0 {
		return nil
	}
	return l.nodeseq.Next
}
