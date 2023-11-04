package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkEOF(err error) bool {
	if err != nil {
		if err == io.EOF {
			return true
		}
		log.Fatal(err)
	}
	return false
}

func closeF(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func computeFrequencyTable(fileName string) map[int]int {
	file, err := os.Open(fileName)
	check(err)
	fileInfo, err := os.Stat(fileName)
	check(err)
	if fileInfo.Size() == 0 {
		log.Fatal("Empty file provided")
	}
	defer closeF(file)

	frequencyTable := make(map[int]int)
	buffer := make([]byte, 1024)

	for {
		number, err := file.Read(buffer)
		checkEOF(err)

		if number == 0 {
			break
		}

		for _, character := range buffer {
			if character != 0 {
				frequencyTable[int(character)]++
			}
		}
	}

	return frequencyTable
}

func computeHuffmanTree(frequencyTable map[int]int) HuffmanNode {
	priorityQueue := PriorityQueue{}
	//priorityQueue.insert(HuffmanNode{weight: 32, element: 'C'})
	//priorityQueue.insert(HuffmanNode{weight: 42, element: 'D'})
	//priorityQueue.insert(HuffmanNode{weight: 120, element: 'E'})
	//priorityQueue.insert(HuffmanNode{weight: 7, element: 'K'})
	//priorityQueue.insert(HuffmanNode{weight: 42, element: 'L'})
	//priorityQueue.insert(HuffmanNode{weight: 24, element: 'M'})
	//priorityQueue.insert(HuffmanNode{weight: 37, element: 'U'})
	//priorityQueue.insert(HuffmanNode{weight: 2, element: 'Z'})

	for element, frequency := range frequencyTable {
		priorityQueue.insert(HuffmanNode{weight: frequency, element: element})
	}

	for priorityQueue.size() > 1 {
		leftNode := priorityQueue.removeMin()
		rightNode := priorityQueue.removeMin()
		parentNode := HuffmanNode{weight: leftNode.weight + rightNode.weight, left: &leftNode, right: &rightNode}
		priorityQueue.insert(parentNode)
	}

	return priorityQueue.removeMin()
}

func computePrefixTable(root HuffmanNode) map[string]string {
	prefixTable := make(map[string]string)
	traverseTree(&root, prefixTable, "")
	return prefixTable
}

func traverseTree(root *HuffmanNode, prefixTable map[string]string, prefix string) {
	if root == nil {
		return
	}
	if root.left == nil && root.right == nil {
		prefixTable[strconv.Itoa(root.element)] = prefix
	} else {
		traverseTree(root.left, prefixTable, prefix+"0")
		traverseTree(root.right, prefixTable, prefix+"1")
	}
}

func compress(fileName string) {
	frequencyTable := computeFrequencyTable(fileName)
	huffmanTree := computeHuffmanTree(frequencyTable)
	prefixTable := computePrefixTable(huffmanTree)
	fmt.Println(prefixTable)
}
