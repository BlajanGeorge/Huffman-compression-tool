package main

import (
	"fmt"
)

type PriorityQueue struct {
	elements []HuffmanNode
}

func (queue *PriorityQueue) heapify(i int) {
	l := 2*i + 1
	r := 2*i + 2
	smallest := i

	if l < len(queue.elements) && queue.elements[l].weight < queue.elements[smallest].weight {
		smallest = l
	}
	if r < len(queue.elements) && queue.elements[r].weight < queue.elements[smallest].weight {
		smallest = r
	}

	if smallest != i {
		tmpNode := queue.elements[i]
		queue.elements[i] = queue.elements[smallest]
		queue.elements[smallest] = tmpNode
		queue.heapify(smallest)
	}
}

func (queue *PriorityQueue) insert(node HuffmanNode) {
	if len(queue.elements) == 0 {
		queue.elements = append(queue.elements, node)
	} else {
		queue.elements = append(queue.elements, node)
		for i := len(queue.elements)/2 - 1; i >= 0; i-- {
			queue.heapify(i)
		}
	}
}

func (queue *PriorityQueue) removeMin() (tmpNode HuffmanNode) {
	if len(queue.elements) == 0 {
		return HuffmanNode{}
	}

	tmpNode = queue.elements[0]
	queue.elements[0] = queue.elements[len(queue.elements)-1]
	queue.elements = append(queue.elements[:len(queue.elements)-1])
	queue.heapify(0)
	return
}

func (queue *PriorityQueue) size() int {
	return len(queue.elements)
}

func (queue *PriorityQueue) print() {
	for _, elem := range queue.elements {
		fmt.Printf(elem.toString())
	}
}

type HuffmanNode struct {
	weight  int
	element int
	left    *HuffmanNode
	right   *HuffmanNode
}

func (node *HuffmanNode) toString() string {
	if node == nil {
		return ""
	}
	var result = fmt.Sprintf("{%d %d}", node.weight, node.element)

	if node.left != nil {
		result += " " + node.left.toString()
	}
	if node.right != nil {
		result += " " + node.right.toString()
	}

	return result
}
