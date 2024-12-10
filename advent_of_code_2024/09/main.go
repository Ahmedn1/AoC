package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func createDiskLayout(diskMap string) ([]int, map[int][]int, int) {
	// create disk layout
	diskLayout := make([]int, 0)
	filesIdxs := make(map[int][]int, 0)

	fileID := 0
	for i := 0; i < len(diskMap); i++ {
		blocks, _ := strconv.Atoi(string(diskMap[i]))
		if i%2 == 0 {
			filesIdxs[fileID] = make([]int, 0)
			for j := 0; j < blocks; j++ {
				diskLayout = append(diskLayout, fileID)
				filesIdxs[fileID] = append(filesIdxs[fileID], len(diskLayout)-1)
			}
			fileID++
		} else {
			for j := 0; j < blocks; j++ {
				diskLayout = append(diskLayout, -1)
			}
		}
	}

	return diskLayout, filesIdxs, fileID
}

func defragDisk(diskLayout []int) []int {
	diskLayoutCopy := make([]int, len(diskLayout))
	copy(diskLayoutCopy, diskLayout)

	freeSpaceIterator := 0
	fileIdxIterator := len(diskLayoutCopy) - 1

	for {
		for diskLayoutCopy[freeSpaceIterator] != -1 {
			freeSpaceIterator++
		}
		for diskLayoutCopy[fileIdxIterator] == -1 {
			fileIdxIterator--
		}

		if freeSpaceIterator >= fileIdxIterator {
			break
		}

		diskLayoutCopy[freeSpaceIterator] = diskLayoutCopy[fileIdxIterator]
		diskLayoutCopy[fileIdxIterator] = -1
		freeSpaceIterator++
		fileIdxIterator--
	}

	return diskLayoutCopy
}

func findContiguousFreeSpace(diskLayout []int, spacesRequired int, firstFileIdx int) int {
	freeSpaceIdx := -1
	spacesFoundCounter := 0
	for i := 0; i < firstFileIdx && spacesFoundCounter < spacesRequired; i++ {
		if diskLayout[i] == -1 {
			if freeSpaceIdx == -1 {
				freeSpaceIdx = i
			}
			spacesFoundCounter++
		} else {
			freeSpaceIdx = -1
			spacesFoundCounter = 0
		}
	}

	if spacesFoundCounter < spacesRequired {
		freeSpaceIdx = -1
	}
	return freeSpaceIdx
}

func defragDiskByFiles(diskLayout []int, filesIdxs map[int][]int, maxFileID int) []int {
	for fileID := maxFileID; fileID >= 0; fileID-- {
		fileIdxs := filesIdxs[fileID]
		if len(fileIdxs) > 0 {
			// find contiguous free space
			spacesRequired := len(fileIdxs)
			freeSpaceIdx := findContiguousFreeSpace(diskLayout, spacesRequired, fileIdxs[0])
			if freeSpaceIdx != -1 {
				// move file to free space
				for i := 0; i < len(fileIdxs); i++ {
					diskLayout[freeSpaceIdx+i] = diskLayout[fileIdxs[i]]
					diskLayout[fileIdxs[i]] = -1
				}
			}
		}
	}

	return diskLayout
}

func calculateDiskChecksum(diskLayout []int) int {
	checksum := 0
	for i := 0; i < len(diskLayout); i++ {
		if diskLayout[i] != -1 {
			checksum += diskLayout[i] * i
		}
	}

	return checksum
}

func main() {
	// read data from test_input.txt
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var diskMap string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		diskMap = scanner.Text()
		fmt.Println("Disk Map:", diskMap)
	}

	diskLayout, filesIdxs, maxFileID := createDiskLayout(diskMap)
	fmt.Println("Disk Layout:", diskLayout)

	defraggedDisk := defragDisk(diskLayout)
	fmt.Println("Defragged Disk:", defraggedDisk)
	checksum := calculateDiskChecksum(defraggedDisk)
	fmt.Println("Checksum:", checksum)

	defraggedDiskByFile := defragDiskByFiles(diskLayout, filesIdxs, maxFileID)
	fmt.Println("Defragged Disk by Files:", defraggedDiskByFile)
	checksum = calculateDiskChecksum(defraggedDiskByFile)
	fmt.Println("Checksum:", checksum)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
