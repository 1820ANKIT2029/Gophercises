package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println(os.Args)

	filename := flag.String("f", "problems.csv", "Problem filename flag eg: -f problems.csv")

	flag.Parse()

	fmt.Println(*filename)

	Questions, err := getQuestion(*filename)

	if err != nil {
		fmt.Println("err in getQuestion", err)
		return
	}

	totalMark := len(Questions)
	mark := 0
	var answer int

	for i := 0; i < totalMark; i++ {
		fmt.Printf("Q%d: %s\n", i+1, Questions[i][0])
		fmt.Printf("Your Answer: ")

		fmt.Scan(&answer)

		tmp, err := strconv.Atoi(Questions[i][1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		if answer == tmp {
			mark++
		}
	}

	fmt.Printf("Your score is %d / %d\n", mark, totalMark)

}

func getQuestion(csvfilepath string) ([][]string, error) {
	file, err := os.Open(csvfilepath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)

	rec, _ := r.ReadAll()

	return rec, nil
}
