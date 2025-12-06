package main

import (
	"advent/helpers"
	"advent/helpers/grid"
	log "log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/samber/lo"
)

// This needs to change to match the input
func parseInput(input string) grid.Grid[string] {
	re := regexp.MustCompile(`\s+`)
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []string {
		cleaned := strings.TrimSpace(re.ReplaceAllString(line, " "))
		return lo.Map(strings.Split(cleaned, " "), func(ele string, _ int) string {
			return ele
		})
	})
}

// Easier to just do it as we read it in for part2, but quite messy
func part2(input string) int {

	lines := strings.Split(input, "\n")

	// Extract operations into list
	inst := lines[len(lines)-1]
	re := regexp.MustCompile(`\s+`)
	cleaned := strings.TrimSpace(re.ReplaceAllString(inst, " "))
	operators := lo.Map(strings.Split(cleaned, " "), func(op string, _ int) string {
		return op
	})

	// Spaces matter! just read in the [][]string and trasnpose it- Numbers will be "normalized", one per line.
	g := grid.Transpose(lo.Map(lines[:len(lines)-1], func(line string, i int) []string {
		return lo.Map(strings.Split(line, ""), func(ele string, _ int) string {
			return ele
		})
	}))

	// now make a list of numbers from our data, seperatored by equation by -1
	numbers := lo.Map(g, func(term []string, _ int) int {
		numberString := strings.TrimSpace(lo.Reduce(term, func(agg string, t string, _ int) string {
			return agg + t
		}, ""))
		if numberString == "" {
			// this is our seperators, no negatives in the data.
			return -1
		}
		return helpers.Atoi(numberString)
	})
	// Do the math!
	return lo.Sum(PerformMath(operators, numbers, []int{}))
}

// Pops on operation off the list, finds the terms until -1, perform the operation.
func PerformMath(operators []string, numbers []int, agg []int) []int {
	if len(operators) == 0 {
		return agg
	}
	// "end" is the index of the next -1 (or end), ie our "terms" to operate on
	i := lo.IndexOf(numbers, -1)
	end := i
	if i == -1 {
		end = len(numbers)
	}
	result := PerformOp(operators[0], numbers[:end])

	if len(numbers) == -1 {
		return append(agg, result)
	}
	return PerformMath(operators[1:], numbers[i+1:], append(agg, result))
}

// No dynamic +/*, so do the operation on the list passed in.
func PerformOp(op string, terms []int) int {
	var start int
	var reducer func(agg int, num int, _ int) int

	switch op {
	case "+":
		start = 0
		reducer = func(agg int, num int, _ int) int { return agg + num }
	default:
		start = 1
		reducer = func(agg int, num int, _ int) int { return agg * num }
	}
	return lo.Reduce(terms, reducer, start)
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)

	// Part 1
	g := grid.Rotate90(parseInput(data))
	for _, v := range g {
		log.Debug("", "line", v)
	}
	pre1 := time.Now()
	result := lo.Sum(lo.Map(g, func(equation []string, _ int) int {
		op := equation[0]
		terms := lo.Map(equation[1:], func(t string, _ int) int {
			return helpers.Atoi(t)
		})
		return PerformOp(op, terms)
	}))
	post1 := time.Now()
	log.Info("Part1", "answer", result, "time", post1.Sub(pre1))


	pre2 := time.Now()
	res := part2(data)
	post2 := time.Now()
	log.Info("Part2", "answer", res, "time", post2.Sub(pre2))
}
