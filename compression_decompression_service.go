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
		_, err := file.Read(buffer)
		if checkEOF(err) {
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

func writeHeader(destFile *os.File, prefixTable map[string]string) {
	prefixTableHeader := ""

	for char, prefix := range prefixTable {
		prefixTableHeader += fmt.Sprintf("%s:%s\n", char, prefix)
	}

	_, err := destFile.Write([]byte(fmt.Sprintf("%d\n", rune(len(prefixTableHeader)))))
	check(err)

	_, err = destFile.Write([]byte(prefixTableHeader))
	check(err)
}

func writeToFile(fileName, destName string, prefixTable map[string]string) {
	inputFile, err := os.Open(fileName)
	check(err)
	destFile, err := os.Create(destName)
	check(err)
	defer closeF(inputFile)
	defer closeF(destFile)
	writeHeader(destFile, prefixTable)
	buffer := make([]byte, 512)
	compressionByte := make([]byte, 1)
	bitsAvailable := 5

	for {
		_, err := inputFile.Read(buffer)
		if checkEOF(err) {
			break
		}

		for _, charByte := range buffer {
			prefixForSymbol := prefixTable[strconv.Itoa(int(charByte))]

			if len(prefixForSymbol) > bitsAvailable {
				compressionByte[0] <<= 3
				compressionByte[0] = compressionByte[0] | byte(5-bitsAvailable)
				_, err := destFile.Write(compressionByte)
				check(err)
				bitsAvailable = 5
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
	}

	if bitsAvailable != 5 {
		compressionByte[0] <<= 3
		compressionByte[0] = compressionByte[0] | byte(5-bitsAvailable)
		_, err := destFile.Write(compressionByte)
		check(err)
		bitsAvailable = 5
		compressionByte[0] = 0
	}
}

func extractHeaderSize(filePtr *os.File) int {
	buffer := make([]byte, 1)
	headerSizeString := ""

	for {
		_, err := filePtr.Read(buffer)
		if checkEOF(err) {
			log.Fatalf("Header size could not be extracted.")
		}

		if string(buffer) == "\n" {
			break
		}

		headerSizeString += string(buffer)
	}

	headerSize, _ := strconv.Atoi(headerSizeString)
	return headerSize
}

func extractPrefixTable(filePtr *os.File) map[string]string {
	prefixTable := make(map[string]string)
	buffer := make([]byte, extractHeaderSize(filePtr))

	_, err := filePtr.Read(buffer)
	if checkEOF(err) {
		log.Fatalf("Prefix table could not be extracted for file %s.", filePtr.Name())
	}

	var letter string
	var prefix string
	afterSeparator := false

	for _, readByte := range buffer {
		if string(readByte) == ":" {
			afterSeparator = true
			continue
		}

		if string(readByte) == "\n" {
			afterSeparator = false
			prefixTable[letter] = prefix
			prefix = ""
			letter = ""
			continue
		}

		if afterSeparator {
			prefix += string(readByte)
		} else {
			letter += string(readByte)
		}
	}

	return prefixTable
}

func composeInversePrefixTable(prefixTable map[string]string) map[string]string {
	inversePrefixTable := make(map[string]string)

	for element, prefix := range prefixTable {
		inversePrefixTable[prefix] = element
	}

	return inversePrefixTable
}

func Compress(fileName, destName string) {
	frequencyTable := computeFrequencyTable(fileName)
	huffmanTree := computeHuffmanTree(frequencyTable)
	prefixTable := computePrefixTable(huffmanTree)
	writeToFile(fileName, destName, prefixTable)
}

func Decompress(fileName, destName string) {
	inputFile, err := os.Open(fileName)
	check(err)
	destFile, err := os.Create(destName)
	check(err)
	defer closeF(inputFile)
	defer closeF(destFile)
	prefixTable := composeInversePrefixTable(extractPrefixTable(inputFile))

	for {
		buffer := make([]byte, 1)
		filePrefix := ""
		i := 0

		_, err := inputFile.Read(buffer)
		if checkEOF(err) {
			break
		}

		writtenBits := buffer[0] & 7
		buffer[0] >>= 3

		for i < int(writtenBits) {
			if buffer[0]&1 == 1 {
				filePrefix += "1"
			} else {
				filePrefix += "0"
			}

			letter, ok := prefixTable[filePrefix]
			if ok {
				filePrefix = ""
				letterNumber, _ := strconv.Atoi(letter)
				_, err := destFile.WriteString(string(rune(letterNumber)))
				check(err)
			}

			buffer[0] >>= 1
			i++
		}
	}
}
