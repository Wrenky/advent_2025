package main

import (
	"advent/helpers"
	log "log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/samber/lo"
)

type Action struct {
	Direction string
	Value     int
}

// This needs to change to match the input
func parseInput(input string) []Action {
	pattern, _ := regexp.Compile(`(\w)(\d+)`)
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) Action {
		res := pattern.FindAllStringSubmatch(line, -1)
		if len(res[0]) != 3 {
			log.Warn("Got less than three values on line", "line", line)
			return Action{}
		}
		return Action{
			Direction: res[0][1],
			Value:     helpers.Atoi(res[0][2]),
		}
	})
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	parsed := parseInput(data)
	for _, v := range parsed {
		log.Debug("", "line", v)
	}

	// Part 1
	pre1 := time.Now()
	turns, p2 := TurnDial([]int{50}, parsed, 0)
	log.Debug("Turns", "turns", turns)
	p1 := len(lo.Filter(turns, func(t int, _ int) bool {
		if t == 0 {
			return true
		}
		return false
	}))
	post1 := time.Now()
	log.Info("Part1", "answer", p1, "time", post1.Sub(pre1))

	// Part 2
	log.Info("Part2", "answer", p2+p1)
}

func countClicks(value int, act Action) int {
	offset := 1
	clicks := 0
	dial := value
	if act.Direction == "L" {
		offset = -1
	}
	for range act.Value {
		if dial += offset; dial%100 == 0 {
			clicks++
		}
	}
	return clicks
}

func turn(value int, act Action, clicks int) (int, int) {
	log.Debug("turn start", "action", act, "value", value)
	var val int
	if act.Direction == "L" {
		val = (value - act.Value)
	} else {
		val = (value + act.Value)
	}

	val = helpers.Mod(val, 100)
	if val != 0 {
		clicks += countClicks(value, act)
	}

	log.Debug("turn end, click debug", "clicks", clicks, "val", val)
	return val, clicks
}

func TurnDial(acc []int, acts []Action, clicks int) ([]int, int) {
	if len(acts) == 0 {
		return acc, clicks
	}
	act := acts[0]
	val := acc[len(acc)-1]
	next := helpers.RemoveElement(acts, 0)

	res, clicks := turn(val, act, clicks)
	return TurnDial(append(acc, res), next, clicks)
}
