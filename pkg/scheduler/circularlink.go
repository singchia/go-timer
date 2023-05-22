package scheduler

import "sync"

type CircularList struct {
	head   *CircularNode
	tail   *CircularNode
	length int
	cur    *CircularNode
	mutex  sync.RWMutex
}

type CircularNode struct {
	data interface{}
	next *CircularNode
}

type ForEachFunc func(node *CircularNode) error

func NewCircularList() *CircularList {
	return &CircularList{head: nil, tail: nil, length: 0, cur: nil}
}

func (c *CircularList) ForEach(f ForEachFunc) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	itor := c.head
	for i := 0; i < c.length; i++ {
		err := f(itor)
		if err != nil {
			return false
		}
		itor = itor.next
	}
	return true
}

func (c *CircularList) generateNode(data interface{}) *CircularNode {
	return &CircularNode{data: data, next: nil}
}

func (c *CircularList) AppendNode(data interface{}) *CircularNode {
	node := c.generateNode(data)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.length == 0 {
		//first time to append node
		c.head = node
		c.tail = node
		c.tail.next = c.head
		c.cur = node
		c.length++
		return node
	}

	c.tail.next = node
	c.tail = c.tail.next
	c.tail.next = c.head
	c.length++
	return node
}

func (c *CircularList) DeleteNode(node *CircularNode) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var nextItor *CircularNode
	itor := c.tail
	for i := 0; i < c.length; i++ {
		nextItor = itor.next
		if nextItor == node {
			//if the node is current pointed node, then move forward
			//if c.length == 0, don't need to specify corresponding action,
			//it will be repointed at AppendNode
			if c.cur == node {
				c.cur = c.cur.next
			}
			//if the node is tail, tail should move backward
			if c.tail == node {
				c.tail = itor
			}
			itor.next = nextItor.next
			c.length--
			return true
		}
		itor = itor.next
	}
	return false
}

func (c *CircularList) GetCurNode() *CircularNode {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.length == 0 {
		return nil
	}
	node := *c.cur
	return &node
}

func (c *CircularList) GetCurNodeWithNoCopied() *CircularNode {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.length == 0 {
		return nil
	}
	return c.cur
}

func (c *CircularList) RightShiftCurPointer() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.length == 0 {
		return false
	}
	c.cur = c.cur.next
	return true
}

func (c *CircularList) RightShiftCurPointerAndUpdate(data interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.length == 0 {
		return false
	}
	c.cur = c.cur.next
	c.cur.data = data
	return true
}

func (c *CircularList) RightShiftCurPointerToCertainNode(dstNode *CircularNode) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	itor := c.head
	for i := 0; i < c.length; i++ {
		if itor == dstNode {
			c.cur = dstNode
			return true
		}
		itor = itor.next
	}
	return false
}
