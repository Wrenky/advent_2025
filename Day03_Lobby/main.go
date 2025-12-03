package main

import (
	"advent/helpers"
	log "log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

// This needs to change to match the input
func parseInput(input string) [][]int {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []int {
		return lo.Map(strings.Split(line, ""), func(char string, _ int) int {
			return helpers.Atoi(char)
		})
	})
}

func pickTwo(bank []int) int {
	first, second := 0, 0
	firsti := -1
	for i, val := range bank[:len(bank)-1] {
		if val > first {
			first = val
			firsti = i
		}
	}
	for _, valj := range bank[firsti+1:] {
		if valj > second {
			second = valj
		}
	}
	return helpers.Atoi(strconv.Itoa(first) + strconv.Itoa(second))
}

// Convert to and from the stupid format
func number(bank []int) int {
	return helpers.Atoi(lo.Reduce(bank, func(agg string, a int, _ int) string {
		return agg + strconv.Itoa(a)
	}, ""))
}
func reverseNumber(a int) []int {
	str := strconv.Itoa(a)
	var digits []int
	for _, char := range str {
		// ugh runes
		digits = append(digits, int(char-'0'))
	}
	return digits
}

func Best(agg []int, prev []int, rest []int) []int {
	if len(rest) == 0 {
		return agg
	}
	cand := rest[0]
	// Now, take our previous value, and find the best with char from rest.
	candidates := lo.Map(prev, func(elem int, i int) int {
		return number(append(helpers.RemoveElement(prev, i), cand))
	})
	best := lo.Max(candidates)
	return Best(append(agg, best), reverseNumber(best), rest[1:])
}

func findBest(bank []int) int {
	start := bank[:12]
	rest := bank[12:]
	all := Best([]int{number(start)}, start, rest)
	return lo.Max(all)
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	joltages := parseInput(data)
	for _, v := range joltages {
		log.Debug("", "line", v)
	}

	// Part 1
	pre1 := time.Now()
	p1 := lo.Map(joltages, func(bank []int, _ int) int {
		return pickTwo(bank)
	})
	post1 := time.Now()
	log.Info("Part1", "answer", lo.Sum(p1), "time", post1.Sub(pre1))

	// Part 2
	pre2 := time.Now()
	p2 := lo.Map(joltages, func(bank []int, _ int) int {
		return findBest(bank)
	})
	post2 := time.Now()
	log.Info("Part2", "answer", lo.Sum(p2), "time", post2.Sub(pre2))
}
