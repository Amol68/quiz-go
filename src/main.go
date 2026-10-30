package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"

	"os"
	"strings"
	"time"
)

var (
	inChan  chan int
	outChan chan int
	correct int
)

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {

		defer f.Close()
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()

	return records, err
}

func getAns(ans string) {

	var inp string

	fmt.Scanln(&inp)

	if sanitize(inp) == sanitize(ans) {
		outChan <- <-inChan + 1
	} else {
		outChan <- <-inChan
	}

}

func updateCorrect() bool {
	select {
	case res := <-outChan:
		correct = res
		return true
	case <-time.After(time.Duration(timeout) * time.Second):
		close(outChan)
		fmt.Println()
		return false
	}
}

func sanitize(s string) string {
	return strings.ToLower(strings.Trim(s, "\n\r\t "))
}

var timeout int

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name to start your quiz:")

	input, err := reader.ReadString('\n')

	fmt.Printf("Input: %s,", input)

	records, err := readCsvFile("./problems.csv")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading your name: %v\n", err)
		os.Exit(1)
	}

	flag.IntVar(&timeout, "timeout", 10, "The amount of time you have for a single question in seconds")
	flag.Parse()

	inChan, outChan = make(chan int, 1), make(chan int, 1)

	for index, curr := range records {
		fmt.Printf("Question %v: %v -> ", index, curr[0])

		go getAns(curr[1])
		inChan <- correct

		if !updateCorrect() {
			break
		}

	}

	fmt.Println("You got", correct, "out of", len(records), "correct.")
}
