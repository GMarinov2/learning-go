package main

import (
	"bufio"
	"os"
	"slices"
	"strconv"
	"strings"
)

func readFileLines() []string {
	file, err := os.OpenFile("./input.dat", os.O_RDONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func getNumArrays(input []string) ([]int, []int) {
	var leftList []int
	var rightList []int

	for _, v := range input {
		nums := strings.Split(v, "   ")

		leftNumber, _ := strconv.Atoi(nums[0])
		rightNumber, _ := strconv.Atoi(nums[1])

		leftList = append(leftList, leftNumber)
		rightList = append(rightList, rightNumber)
	}

	return leftList, rightList
}

func calculateDistance(leftList []int, rightList []int) int {
	result := 0

	for i, _ := range leftList {
		distance := rightList[i] - leftList[i]
		if distance < 0 {
			distance *= -1
		}
		result += distance
	}

	return result
}

func getCounts(arr []int) []int {
	result := make([]int, arr[len(arr)-1]+1)

	for _, v := range arr {
		result[v] = result[v] + 1
	}

	return result
}

func caclulateSimilarity(leftList []int, counts []int) int {
	result := 0

	for _, v := range leftList {
		result += v * counts[v]
	}

	return result
}

func main() {
	lines := readFileLines()
	leftList, rightList := getNumArrays(lines)

	slices.Sort(leftList)
	slices.Sort(rightList)

	distance := calculateDistance(leftList, rightList)
	similarity := caclulateSimilarity(leftList, getCounts(rightList))

	println(distance)
	println(similarity)
}
