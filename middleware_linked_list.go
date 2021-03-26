package cgr

type _node struct {
	next  *_node
	mware *Middleware
}

// singly linked list for storing middlewares
type middlewareLinkedList struct {
	head *_node
}

func (ll *middlewareLinkedList) insert(mware *Middleware) {
	newNode := new(_node)
	newNode.next = nil
	newNode.mware = mware

	if ll.head == nil {
		ll.head = newNode
	} else {
		currNode := ll.head
		for currNode.next != nil {
			currNode = currNode.next
		}
		currNode.next = newNode
	}
}

func (ll *middlewareLinkedList) len() int {
	currNode := ll.head

	if currNode == nil {
		return 0
	}

	var len int

	for currNode != nil {
		len++
		currNode = currNode.next
	}
	return len
}
