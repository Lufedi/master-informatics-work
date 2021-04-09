package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Create a file with 2e14 numbers randombly generated
func generateFile(fileName string, totalNumbers int) {
	file, err := os.Create(fileName)
	defer file.Close()
	check(err)

	max, min := math.MaxInt32, math.MinInt32
	rand.Seed(time.Now().UnixNano())
	//fmt.Println(int(math.Pow(2, 14)))
	for i := 0; i < totalNumbers; i++ {
		n := rand.Intn(max-min) + min
		_, err := file.WriteString(fmt.Sprint(n) + "\n")
		check(err)
	}

}

func loadNumbers(fileName string) []int {
	filename := fileName
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		check(err)
		result = append(result, x)
	}
	return result
}

func aggregateResult(channel *chan int, L int) int {
	var res float64
	for i := 0; i < L; i++ {
		n := <-*channel
		res = math.Max(float64(n), float64(res))
	}
	return int(res)
}

func getMaxValueInRange(channel chan int, low int, high int, numbers []int) {
	//fmt.Println("going from", low, "to", high)
	var result float64
	for i := low; i < high; i++ {
		result = math.Max(float64(result), float64(numbers[i]))
	}
	//fmt.Println("max", result)
	channel <- int(result)
}

func linealMax(numbers []int) int {
	var res float64
	for i := 0; i < len(numbers); i++ {
		res = math.Max(float64(res), float64(numbers[i]))
	}
	return int(res)
}

func runMaxNumberProblem(numberOfCoroutines int) {
	fileName := "labnumbers.txt"
	numbers := loadNumbers(fileName)
	L := len(numbers)

	segmentLength := int(L / numberOfCoroutines)

	channel := make(chan int)
	defer close(channel)

	for i := 0; i < numberOfCoroutines; i++ {
		go getMaxValueInRange(channel, i*segmentLength, segmentLength*(i+1), numbers)
	}

	aggregateResult(&channel, numberOfCoroutines)
	//maxN := aggregateResult(&channel, numberOfCoroutines)
	//fmt.Println("The max number is", maxN)
}

func main() {
	//generateFile("labnumbers20.txt", int(math.Pow(2, 14)))

	const CASES = 19
	var results [CASES]int64

	for i := 1; i <= CASES; i++ {
		var averageTime int64
		for j := 0; j < 10; j++ {
			start := time.Now()
			n_goroutines := int(math.Pow(2, float64(i)))
			runMaxNumberProblem(n_goroutines)
			averageTime += time.Since(start).Nanoseconds()
		}
		averageTime = averageTime / 10
		results[i-1] = averageTime
	}

	for i, v := range results {
		fmt.Printf("N coroutines: %6d   seconds: %10f   nanoseconds %10d \n", int(math.Pow(2, float64(i+1))), float64(v)/1000000000, v)
	}
	fmt.Printf("\n")
}
