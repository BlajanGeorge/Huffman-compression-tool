package main

import (
	"fmt"
)

type Node interface {
	Weight() int
	Element() int
	ToString() string
}

type PriorityQueue[T Node] struct {
	elements []T
}

func (queue *PriorityQueue[Node]) heapify(i int) {
	l := 2*i + 1
	r := 2*i + 2
	smallest := i

	if l < len(queue.elements) && queue.elements[l].Weight() < queue.elements[smallest].Weight() {
		smallest = l
	}
	if r < len(queue.elements) && queue.elements[r].Weight() < queue.elements[smallest].Weight() {
		smallest = r
	}

	if smallest != i {
		tmpNode := queue.elements[i]
		queue.elements[i] = queue.elements[smallest]
		queue.elements[smallest] = tmpNode
		queue.heapify(smallest)
	}
}

func (queue *PriorityQueue[Node]) insert(node Node) {
	if len(queue.elements) == 0 {
		queue.elements = append(queue.elements, node)
	} else {
		queue.elements = append(queue.elements, node)
		for i := len(queue.elements)/2 - 1; i >= 0; i-- {
			queue.heapify(i)
		}
	}
}

func (queue *PriorityQueue[Node]) removeMin() (tmpNode HuffmanNode) {
	if len(queue.elements) == 0 {
		return HuffmanNode{}
	}

	tmpNode = queue.elements[0]
	queue.elements[0] = queue.elements[len(queue.elements)-1]
	queue.elements = append(queue.elements[:len(queue.elements)-1])
	queue.heapify(0)
	return
}

func (queue *PriorityQueue[Node]) size() int {
	return len(queue.elements)
}

func (queue *PriorityQueue[Node]) print() {
	for _, elem := range queue.elements {
		fmt.Printf(elem.ToString())
	}
}

type HuffmanNode struct {
	weight  int
	element int
	left    *HuffmanNode
	right   *HuffmanNode
}

func (node *HuffmanNode) ToString() string {
	if node == nil {
		return ""
	}
	var result = fmt.Sprintf("{%d %d}", node.weight, node.element)

	if node.left != nil {
		result += " " + node.left.ToString()
	}
	if node.right != nil {
		result += " " + node.right.ToString()
	}

	return result
}

func (node *HuffmanNode) Weight() int {
	return node.weight
}

func (node *HuffmanNode) Element() int {
	return node.element
}
