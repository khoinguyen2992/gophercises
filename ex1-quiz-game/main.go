package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

var (
	csvFileFlag   = flag.String("csv", "problem.csv", `a csv file in the format of 'question,anwer'`)
	timeLimitFlag = flag.Int("limit", 30, "the time for the quiz in seconds")
	shuffleFlag   = flag.Bool("shuffle", false, "the flag if problems will be shuffled or not")
)

type problem struct {
	question string
	solution string
}

func main() {
	flag.Parse()

	problems, err := getProblems(getCSVFile())
	if err != nil {
		panic(err)
	}

	if getShuffleFlag() {
		problems = shuffleProblems(problems)
	}

	done := make(chan bool, 1)
	results := make(chan bool, len(problems))
	go func() {
		quiz(problems, results, done)
	}()

	select {
	case <-done:
		fmt.Println("Completed")
	case <-time.After(time.Duration(getTimeLimitFlag()) * time.Second):
		fmt.Println()
		fmt.Println("Out of time")
	}

	score := calculateResult(results)
	renderFinalScore(score, len(problems))
}

func quiz(problems []problem, result chan bool, done chan bool) {
	for _, p := range problems {
		renderQuestion(p.question)
		result <- isCorrect(getAnswer(), p.solution)
	}

	done <- true
}

func calculateResult(results chan bool) int {
	close(results) // dont wait anymore
	score := 0
	for result := range results {
		if result {
			score = score + 1
		}
	}

	return score
}

func getAnswer() string {
	var answer string
	fmt.Scanf("%s\n", &answer)
	return answer
}

func renderQuestion(question string) {
	fmt.Printf("%s is: ", question)
}

func renderFinalScore(score, total int) {
	fmt.Printf("You scored %d out of %d\n", score, total)
}

func isCorrect(solution string, answer string) bool {
	return refineString(solution) == refineString(answer)
}

func refineString(in string) string {
	return strings.TrimSpace(strings.ToLower(in))
}

func getProblems(csvFile string) ([]problem, error) {
	var problems []problem
	data, err := ioutil.ReadFile(csvFile)
	if err != nil {
		return problems, err
	}

	r := csv.NewReader(bytes.NewReader(data))
	lines, err := r.ReadAll()
	if err != nil {
		return problems, err
	}

	problems = make([]problem, len(lines))
	for index, line := range lines {
		problems[index] = problem{
			question: line[0],
			solution: line[1],
		}
	}
	return problems, nil
}

func shuffleProblems(problems []problem) []problem {
	shuffledProblems := make([]problem, len(problems))
	perm := rand.Perm(len(problems))
	for i, v := range perm {
		shuffledProblems[v] = problems[i]
	}

	return shuffledProblems
}

func getCSVFile() string {
	return *csvFileFlag
}

func getTimeLimitFlag() int {
	return *timeLimitFlag
}

func getShuffleFlag() bool {
	return *shuffleFlag
}
