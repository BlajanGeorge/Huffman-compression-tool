package main

import (
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
	if root.left == nil && root.right == nil {
		prefixTable[strconv.Itoa(root.element)] = "0"
	} else {
		traverseTree(&root, prefixTable, "")
	}
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

func writeToFile(fileName, destName string, prefixTable map[string]string) {
	inputFile, err := os.Open(fileName)
	check(err)
	destFile, err := os.Create(destName)
	check(err)
	defer closeF(inputFile)
	defer closeF(destFile)
	buffer := make([]byte, 1)
	compressionByte := make([]byte, 1)
	bitsAvailable := 8

	for {
		number, err := inputFile.Read(buffer)
		checkEOF(err)

		if number == 0 {
			break
		}

		prefixForSymbol := prefixTable[strconv.Itoa(int(buffer[0]))]

		if len(prefixForSymbol) > bitsAvailable {
			_, err := destFile.Write(compressionByte)
			check(err)
			bitsAvailable = 8
			compressionByte[0] = 0
		}

		for _, char := range prefixForSymbol {
			compressionByte[0] <<= 1
			if char == '1' {
				compressionByte[0] += 1
			}
			bitsAvailable--
		}
	}

	if bitsAvailable != 8 {
		_, err := destFile.Write(compressionByte)
		check(err)
		bitsAvailable = 8
		compressionByte[0] = 0
	}
}

func Compress(fileName, destName string) {
	//TODO: write header with prefix table in compressed file
	frequencyTable := computeFrequencyTable(fileName)
	huffmanTree := computeHuffmanTree(frequencyTable)
	prefixTable := computePrefixTable(huffmanTree)
	writeToFile(fileName, destName, prefixTable)
}

func Decompress(fileName, destName string) {
	//TODO: decompress file
}
