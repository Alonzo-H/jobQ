package list

import "sync"

func New() List {
	return List{}
}

type List struct {
	head *Node
	tail *Node
    lock sync.RWMutex
}

func (l *List) Push(i uint64) *Node {
    l.lock.Lock()
    defer l.lock.Unlock()

	n := Node{value: i}

	if l.tail == nil && l.head == nil {
		l.head = &n
		l.tail = &n
		return &n
	}

	n.prev = l.tail
	l.tail.next = &n
	l.tail = &n
	return &n
}

func (l *List) Head() (uint64, error) {
    l.lock.RLock()
    defer l.lock.RUnlock()

	if l.head == nil {
		return 0, ErrEmptyList
	}

	return l.head.value, nil
}

func (l *List) Pop() (uint64, error) {
    l.lock.Lock()
    defer l.lock.Unlock()

	if l.head == nil {
		return 0, ErrEmptyList
	}

	if l.head == l.tail {
		v := l.head.value
		l.head = nil
		l.tail = nil
		return v, nil
	}

	o := l.head
	l.head = o.next
	o.next = nil
	if l.head != nil {
		l.head.prev = nil
	}
	return o.value, nil
}

func (l *List) Remove(n *Node) {
    l.lock.Lock()
    defer l.lock.Unlock()

	if l.head == n {
		l.head = n.next
	}
	if l.tail == n {
		l.tail = n.prev
	}

	if n.next != nil {
		n.next.prev = n.prev
	}

	if n.prev != nil {
		n.prev.next = n.next
	}

	n.prev = nil
	n.next = nil
}

type Node struct {
	value uint64
	next  *Node
	prev  *Node
}
