package game

type posWithPriority struct {
	Pos
	priority int
}

type posPriorityQueue []posWithPriority

func (pq posPriorityQueue) push(pos Pos, priority int) posPriorityQueue {
	newNode := posWithPriority{pos, priority}
	pq = append(pq, newNode)

	newNodeIndex := len(pq) - 1
	parentIndex, parent := pq.parent(newNodeIndex)

	for newNode.priority < parent.priority && newNodeIndex > 0 {
		pq.swap(newNodeIndex, parentIndex)
		newNodeIndex = parentIndex
		parentIndex, parent = pq.parent(newNodeIndex)
	}

	return pq
}

func (pq posPriorityQueue) pop() (posPriorityQueue, Pos) {
	result := pq[0].Pos
	pq[0] = pq[len(pq)-1]
	pq = pq[:len(pq)-1]

	if len(pq) == 0 {
		return pq, result
	}

	index := 0
	node := pq[index]

	leftExists, leftIndex, left := pq.left(index)
	rightExists, rightIndex, right := pq.right(index)

	for (leftExists && node.priority > left.priority) || (rightExists && node.priority > right.priority) {
		if !rightExists || left.priority <= right.priority {
			pq.swap(index, leftIndex)
			index = leftIndex
		} else {
			pq.swap(index, rightIndex)
			index = rightIndex
		}

		leftExists, leftIndex, left = pq.left(index)
		rightExists, rightIndex, right = pq.right(index)
	}

	return pq, result
}

func (pq posPriorityQueue) swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq posPriorityQueue) parent(i int) (int, posWithPriority) {
	index := (i - 1) / 2
	return index, pq[index]
}

func (pq posPriorityQueue) left(i int) (bool, int, posWithPriority) {
	index := i*2 + 1

	if index < len(pq) {
		return true, index, pq[index]
	}

	return false, 0, posWithPriority{}
}

func (pq posPriorityQueue) right(i int) (bool, int, posWithPriority) {
	index := i*2 + 2

	if index < len(pq) {
		return true, index, pq[index]
	}

	return false, 0, posWithPriority{}
}
