package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func readLinesFromFile(filename string) <-chan string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	outputChannel := make(chan string)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			outputChannel <- line
		}
		close(outputChannel)

		err := file.Close()
		if err != nil {
			fmt.Printf("Unable to close file: %v\n", err.Error())
			return
		}
	}()

	return outputChannel
}

func mergeChannels(channels ...<-chan string) chan string {
	mergedChannel := make(chan string)
	var waitGroup sync.WaitGroup

	for _, channel := range channels {
		waitGroup.Add(1)
		go func(c <-chan string) {
			for value := range c {
				mergedChannel <- value
			}
			waitGroup.Done()
		}(channel)
	}

	go func() {
		waitGroup.Wait()
		close(mergedChannel)
	}()

	return mergedChannel
}

func main() {
	channel1 := readLinesFromFile("./text1.txt")
	channel2 := readLinesFromFile("./text2.txt")
	mergedChannel := mergeChannels(channel1, channel2)

	for value := range mergedChannel {
		fmt.Println(value)
	}
}
