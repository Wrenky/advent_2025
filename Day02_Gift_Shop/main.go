package main

import (
	"advent/helpers"
	log "log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

type IDRange struct {
	Start int
	End   int
}

// This needs to change to match the input
func parseInput(input string) []IDRange {
	return lo.Map(strings.Split(input, ","), func(line string, _ int) IDRange {
		values := strings.Split(line, "-")
		return IDRange{
			Start: helpers.Atoi(values[0]),
			End:   helpers.Atoi(values[1]),
		}
	})
}

func valid(id int) bool {
	str := strconv.Itoa(id)
	mid := (len(str)) / 2
	return str[mid:] != str[:mid]
}

func validDouble(id int) bool {
	str := strconv.Itoa(id)
	for i := 1; i <= len(str)/2; i++ {
		bad := false
		curr := str[0:i]
		sets := lo.ChunkString(str[i:], len(curr))
		for _, cand := range sets {
			if curr != cand {
				bad = false
				break
			}
			bad = true
		}
		if bad {
			return false
		}
	}
	return true
}

// Samber lo.RangeFrom is insane?
func kRange(a, b int) []int {
	res := []int{}
	for i := a; i <= b; i++ {
		res = append(res, i)
	}
	return res
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	parsed := parseInput(data)
	for _, v := range parsed {
		log.Debug("", "line", v)
	}

	// Function to do the work with  predicate checker
	runner := func(parsed []IDRange, predicate func(int) bool) []int {
		return lo.Flatten(lo.FilterMap(parsed, func(set IDRange, _ int) ([]int, bool) {
			pass := lo.Filter(kRange(set.Start, set.End), func(id int, _ int) bool {
				return !predicate(id)
			})
			return pass, len(pass) != 0
		}))
	}

	// Part 1, timed
	pre1 := time.Now()
	invalidIds := runner(parsed, valid)
	post1 := time.Now()
	log.Info("Part1", "answer", lo.Sum(invalidIds), "time", post1.Sub(pre1))

	// Part 2, timed
	pre2 := time.Now()
	p2 := runner(parsed, validDouble)
	post2 := time.Now()
	log.Info("Part2", "answer", lo.Sum(p2), "time", post2.Sub(pre2))
}
