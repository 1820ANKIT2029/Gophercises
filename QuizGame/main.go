package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// flags parsing
	filename := flag.String("f", "problems.csv", "Problem filename flag eg: -f problems.csv")
	timerValue := flag.Int("t", 30, "Time given to answer a question eg: -t 30")
	shuffle := flag.Bool("s", false, "question will get shuffle eg: -s")
	flag.Parse()

	// get questions from file
	Questions, err := getQuestion(*filename, *shuffle)
	if err != nil {
		fmt.Println("err in getQuestion", err)
		return
	}

	// fmt.Println(Questions)

	// Mark, answer temp variable
	totalMark := len(Questions)
	mark := 0
	var answer string

	// channel for timer, input
	outTimerChan := make(chan bool)
	inTimerChan := make(chan bool)
	answerChan := make(chan string)

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
			answer = strings.ToLower(strings.TrimSpace(answer))
			// add 1 mark if answer is correct
			if answer == Questions[i][1] {
				mark++
			}
		}
	}

	// display final Score
	fmt.Printf("Your score is %d / %d\n", mark, totalMark)
}

// function to get Question List from the Provided filepath
func getQuestion(csvfilepath string, shuffle bool) ([][]string, error) {
	file, err := os.Open(csvfilepath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)

	rec, _ := r.ReadAll()

	// shuffle question
	if shuffle {
		rand.Shuffle(len(rec), func(i, j int) {
			rec[i], rec[j] = rec[j], rec[i]
		})
	}

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
func Input(answerChan chan string) {
	var answer string
	fmt.Scan(&answer)

	answerChan <- answer
}
