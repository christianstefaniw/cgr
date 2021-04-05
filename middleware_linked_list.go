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
