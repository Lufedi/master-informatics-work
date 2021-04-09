package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const EPSILON = 0.000001

func gcd(fa int, fb int) int {
    a := fa
    b := fb
    
    for b != 0 {
        t := b
        b = a % b
        a = t
    }
    return a
}

func min(a,b int) int  {
	if a < b {
		return a
	}
	return b
}

func max(a,b int) int  {
	if a > b {
		return a
	}
	return b
}

type Fraction struct {
	a int //num
	b int //den
}

func (f Fraction) checkZero() {
	if f.b == 0 {
		panic("0 in denominator")
	}
}

func NewFraction(a,b int) Fraction{
	if b == 0 {
		panic("cant create fraction with 0 in the denomitaor")
	}
	if a == 0 {
		b = 1
	}

	return Fraction{
		a: a,
		b: b,
	}
}

var MINUS_1 = Fraction{ a: -1, b: 1}

type Matrix [][]Fraction

// Reduce all the fractions in the matrix
func (m Matrix) simplify() {
    for i := 0; i < len(m); i++ {
        for j := 0; j < len(m[i]); j++ {
            m[i][j].simplify()
        }
    }
}

func (f Fraction) g() float64 {
	f.checkZero()
	return float64(f.a) / float64(f.b)
}

func (f Fraction) eq(o Fraction) bool {
	var diff = math.Abs(float64(o.a)/float64(o.b) - float64(f.a)/float64(f.b))
	return diff < EPSILON
}

func (f Fraction) multiply(o Fraction) Fraction {
	return NewFraction(f.a*o.a, f.b*o.b)
}

func (f Fraction) add(o Fraction) Fraction {
	return NewFraction((f.a * o.b) + (f.b * o.a), f.b*o.b)
}

func (f Fraction) minus(o Fraction) Fraction {
	return NewFraction((f.a * o.b) - (f.b * o.a), f.b * o.b)
}

func (f Fraction) divide(o Fraction) Fraction {
	return NewFraction(f.a * o.b, f.b*o.a)
}

func (f Fraction) compare(o Fraction) int {
	if f.eq(o) {
		return 0
	} else if f.g() < o.g() {
		return -1
	} else {
		return 1
	}
}

func (f *Fraction) simplify() { 
    g := gcd(f.a, f.b)
    f.a = f.a / g
    f.b = f.b / g
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(fileName string) (int, Matrix) {
	file, err := os.Open(fileName)
	check(err)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	defer file.Close()

    scanner.Scan()
    n, err := strconv.Atoi(scanner.Text())
    check(err)

	//create initial array
	numbers := make(Matrix, n)
	for i := 0; i < n; i++ {
		numbers[i] = make([]Fraction, n+1)
	}

	var line = 0
	for scanner.Scan() {
		text := scanner.Text()
		data := strings.Fields(text)

		for i, v := range data {
			parsedNumber, err := strconv.Atoi(v)
			check(err)
			numbers[line][i] = NewFraction(parsedNumber,1)
		}
		line += 1
	}
	return n, numbers
}

func printMatrix(matrix Matrix) {

	if len(matrix) == 0 {
		panic("Matrix is empty")
	}
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("(%d, %d)", matrix[i][j].a, matrix[i][j].b)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func swapRows(i int, j int, matrix Matrix) {
	tmp := make([]Fraction, len(matrix[0]))
	copy(tmp, matrix[i])
	copy(matrix[i], matrix[j])
	copy(matrix[j], tmp)
}

func getPivotRow(matrix Matrix, column int) (int, Fraction) {
	var pivotRow int = column
	currentMax := matrix[column][column]
	for i := column; i < len(matrix); i++ {
		if currentMax.compare(matrix[i][column]) <= 0 {
			currentMax = matrix[i][column]
			pivotRow = i
		}
	}
	return pivotRow, currentMax
}

func killColumn(column int, matrix Matrix, numberOfCoroutines int) {
	pivot := matrix[column][column]
	var wg sync.WaitGroup
	L := len(matrix)
	segmentLength := int(L/numberOfCoroutines)
	for j := 0; j < numberOfCoroutines; j++ {
		wg.Add(1)
		go killRow(matrix, pivot,  column, j*segmentLength, segmentLength*(j+1), &wg)
	}
	wg.Wait()
}

func killRow(matrix Matrix, pivot Fraction, column int, low int, high int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := max(column+1, low); i < high; i++{
		current := matrix[i][column]
		x := current.divide(pivot).multiply(MINUS_1)
		changeRows(matrix, i, x, pivot)
	}
}

func changeRows(matrix Matrix, row int, x Fraction, pivot Fraction){
	//r3 <- r3 - x*r2
	for j := 0; j < len(matrix[row]); j++ {
		current := matrix[row][j]
		result := pivot.multiply(x).add(current)
		matrix[row][j] = result
	}
}

func generateRandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min) + min
}

func createRandomMatrix(n int) Matrix{
	numbers := make(Matrix, n)
	for i := 0; i < n; i++ {
		numbers[i] = make([]Fraction, n+1)
	}
	const MAX_NUMBER = 4
	for i:=0 ; i < n; i++{
		for j:=0; j < n+1; j++{
			numbers[i][j] = NewFraction(generateRandomInt(0, MAX_NUMBER), 1 )
		}
	}
	return numbers
}

func reduceRow(matrix Matrix, low, high, pivotRow int, wg *sync.WaitGroup) {
    defer wg.Done()

    for i := low; i <= min(len(matrix) - 1, high); i++  {
		r2 := matrix[pivotRow][pivotRow]
		r3 := matrix[i][pivotRow]
		 
		if r2.a == 0 {
			//r3 <- r3*r2
			for j := 0; j < len(matrix[i]); j++ {
				matrix[i][j] = r3.multiply(r2)
			}
		} else {
			//r3 <- r3 - x*r2
			x :=  r3.divide(r2)
			for j := 0; j < len(matrix[i]); j++ {
				r2, r3 := matrix[pivotRow][j], matrix[i][j]
				matrix[i][j] = r3.minus(x.multiply(r2))
			}
		}
		matrix.simplify()
    }
}

func reduceRows(pivot int, matrix Matrix, coroutines int) {
    var wg sync.WaitGroup
	L := len(matrix) - pivot
	segmentLength := int(L/coroutines)
	start := pivot + 1
	
	for start <= len(matrix) {
		wg.Add(1)
		go reduceRow(matrix, start, start +  segmentLength, pivot, &wg)
		start += segmentLength + 1 
	}

	wg.Wait()	
}

func writeFile(matrix Matrix) {
	N := len(matrix)
	f, err := os.Create("matrixrandom.txt")
	check(err)

	defer f.Close()

	_, writeErr := f.WriteString(fmt.Sprintf("%d\n", N))
	check(writeErr)

	for i := 0; i < len(matrix) ; i++ {
		line := ""
		for j := 0; j < len(matrix[i]) ; j++ {
			//line := strings.Trim(strings.Join(strings.Split(fmt.Sprint(matrix[i]), " "), " "), "[]")
			d := fmt.Sprint(matrix[i][j].a)
			if j == 0 {
				line += d
			} else {
				line += " " + d
			}
		}
		_, writeErr := f.WriteString(fmt.Sprintf("%s\n", line))
		check(writeErr)
	}
}

func gaussElminationProblem(numbers Matrix, nCoroutines int) {
	// TODO validate if killColumn and reduceRows can be merge in one single func, look at it later when refactoring
	N := len(numbers)
	for i := 0; i < N-1; i++ {
		if i == 0 {
			pivot, _ := getPivotRow(numbers, i)
			swapRows(i, pivot, numbers)
			killColumn(i, numbers, nCoroutines)
		} else {
			reduceRows(i, numbers, nCoroutines)
		}
	}
}

func main() {
	
	const CASES = 10
	const REPEATS = 1000
	var results [CASES]int64
	//numbers, N := createRandomMatrix(10),  10
	//writeFile(numbers)
    _, numbers := readFile("matrixrandom.txt")
	printMatrix(numbers)
	
	for c := 1; c <= CASES; c++ {

		var averageTime int64	
		start := time.Now()

		for cr := 0; cr < REPEATS; cr++ {
			gaussElminationProblem(numbers, c)
			averageTime += time.Since(start).Nanoseconds()
		}
		averageTime = averageTime / REPEATS
		results[c-1] = averageTime
	}
	
	for i, v := range results {
		fmt.Printf("N coroutines: %6d   seconds: %10f   nanoseconds %10d \n", i+1, float64(v)/1000000000, v)
	}
}
