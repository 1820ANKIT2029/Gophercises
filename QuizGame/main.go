package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	// flags parsing
	filename := flag.String("f", "problems.csv", "Problem filename flag eg: -f problems.csv")
	timerValue := flag.Int("t", 30, "Time given to answer a question eg: -t 30")
	flag.Parse()

	// get questions from file
	Questions, err := getQuestion(*filename)
	if err != nil {
		fmt.Println("err in getQuestion", err)
		return
	}

	// Mark, answer temp variable
	totalMark := len(Questions)
	mark := 0
	var answer int

	// channel for timer, input
	outTimerChan := make(chan bool)
	inTimerChan := make(chan bool)
	answerChan := make(chan int)

	go timer(outTimerChan, inTimerChan, *timerValue) // timer goroutine

	for i := 0; i < totalMark; i++ {
		fmt.Printf("Q%d: %s\n", i+1, Questions[i][0]) // Question printing
		fmt.Printf("Your Answer: ")

		go Input(answerChan) // get user answer goroutine
		inTimerChan <- true  // start timer by pass true in inTimeChan channel

		// answer verification and time out checking
		select {
		case <-outTimerChan:
			fmt.Println("\nTimes up!!")
			i = totalMark + 1

		case answer = <-answerChan:
			inTimerChan <- false
			tmp, err := strconv.Atoi(Questions[i][1])
			if err != nil {
				fmt.Println(err)
				continue
			}

			// add 1 mark if answer is correct
			if answer == tmp {
				mark++
			}
		}
	}

	// display final Score
	fmt.Printf("Your score is %d / %d\n", mark, totalMark)
}

// function to get Question List from the Provided filepath
func getQuestion(csvfilepath string) ([][]string, error) {
	file, err := os.Open(csvfilepath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)

	rec, _ := r.ReadAll()

	return rec, nil
}

// timer gorooutine
func timer(outTimerChan, inTimerChan chan bool, t int) {
	ticker := time.NewTicker(time.Duration(t) * time.Second)
	ticker.Stop()

	var b bool

	for {
		select {
		case <-ticker.C:
			outTimerChan <- true
		case b = <-inTimerChan:
			if b {
				ticker.Reset(time.Duration(t) * time.Second)
			} else {
				ticker.Stop()
			}
		}
	}
}

// user input goroutine
func Input(answerChan chan int) {
	var answer int
	fmt.Scan(&answer)

	answerChan <- answer
}
