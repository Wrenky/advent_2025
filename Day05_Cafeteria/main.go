package main

import (
	"advent/helpers"
	"fmt"
	log "log/slog"
	"strings"
	"time"

	"github.com/samber/lo"
)

type IDRange struct {
	start int
	end   int
}

func (f IDRange) String() string {
	return fmt.Sprintf("(%d-%d)", f.start, f.end)
}

// This needs to change to match the input
func parseInput(input string) ([]IDRange, []int) {
	parts := strings.Split(input, "\n\n")
	fresh := lo.Map(strings.Split(parts[0], "\n"), func(line string, _ int) IDRange {
		r := strings.Split(line, "-")
		return IDRange{
			start: helpers.Atoi(r[0]),
			end:   helpers.Atoi(r[1]),
		}
	})
	ingredients := lo.Map(strings.Split(parts[1], "\n"), func(line string, _ int) int {
		return helpers.Atoi(line)
	})
	return fresh, ingredients
}

// Is this ingredient fresh?
func checkFreshness(fr []IDRange, id int) bool {
	return lo.SomeBy(fr, func(r IDRange) bool {
		return (id >= r.start) && (id <= r.end)
	})
}

// Do these ranges overlap?
func overlap(a, b IDRange) bool {
	return a.start <= b.end && b.start <= a.end
}

// Simple combine of two ranges
func combine(a, b IDRange) IDRange {
	return IDRange{
		start: min(a.start, b.start),
		end:   max(a.end, b.end),
	}
}

// Combine all the ranges! Necessary for part 2.
// First, pop a candidate off the input list
// Compare the candidate to the existing "final" list:
//
//	  If candidate is equal to an existing range, drop it.
//		 If it overlaps with another range:
//			   remove it from final, combine it with the next input and recurse
//		 If there are no overlaps, just add it to final.
//
// If there is no input left, we are done!
func CombineRanges(final, input []IDRange) []IDRange {
	if len(input) == 0 {
		return final
	}
	cand, next := input[0], input[1:]
	for i, existing := range final {
		if cand == existing {
			// Drop it!
			return CombineRanges(final, next)
		}
		if overlap(existing, cand) {
			n := combine(existing, cand)
			nFinal := helpers.RemoveElement(final, i)
			return CombineRanges(nFinal, append([]IDRange{n}, next...))
		}
	}
	return CombineRanges(append(final, cand), next)
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	fresh, ingredients := parseInput(data)
	for _, v := range fresh {
		log.Debug("", "fresh range", v)
	}
	for _, v := range ingredients {
		log.Debug("", "ingredient id", v)
	}

	// Part 1: Just count fresh ingredients
	pre1 := time.Now()
	FreshIDs := lo.Filter(ingredients, func(id int, _ int) bool {
		return checkFreshness(fresh, id)
	})
	post1 := time.Now()
	log.Info("Part1", "answer", len(FreshIDs), "time", post1.Sub(pre1))

	// Part 2. Seed CombineRanges with the start of the freshlist, then run
	pre2 := time.Now()
	combined := lo.Uniq(CombineRanges([]IDRange{fresh[0]}, fresh[1:]))
	ans := lo.Sum(lo.Map(combined, func(f IDRange, _ int) int {
		return f.end - f.start + 1
	}))
	post2 := time.Now()
	log.Info("Part2", "answer", ans, "time", post2.Sub(pre2))
}
